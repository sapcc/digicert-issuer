/*
Copyright 2019 The cert-manager authors.
Licensed under the Apache License, Version 2.0 (the "License");
you may not use this file except in compliance with the License.
You may obtain a copy of the License at
    http://www.apache.org/licenses/LICENSE-2.0
Unless required by applicable law or agreed to in writing, software
distributed under the License is distributed on an "AS IS" BASIS,
WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
See the License for the specific language governing permissions and
limitations under the License.
*/

package certmanager

import (
	"context"
	"fmt"
	"time"

	"github.com/go-logr/logr"
	apiutil "github.com/jetstack/cert-manager/pkg/api/util"
	cmapi "github.com/jetstack/cert-manager/pkg/apis/certmanager/v1alpha2"
	cmmeta "github.com/jetstack/cert-manager/pkg/apis/meta/v1"
	"github.com/prometheus/client_golang/prometheus"
	certmanagerv1beta1 "github.com/sapcc/digicert-issuer/apis/certmanager/v1beta1"
	"github.com/sapcc/digicert-issuer/pkg/provisioners"
	core "k8s.io/api/core/v1"
	apierrors "k8s.io/apimachinery/pkg/api/errors"
	"k8s.io/apimachinery/pkg/runtime"
	"k8s.io/apimachinery/pkg/types"
	"k8s.io/client-go/tools/record"
	ctrl "sigs.k8s.io/controller-runtime"
	"sigs.k8s.io/controller-runtime/pkg/client"
	"sigs.k8s.io/controller-runtime/pkg/event"
	"sigs.k8s.io/controller-runtime/pkg/predicate"
)

// CertificateRequestReconciler reconciles a DigicertIssuer object.
type CertificateRequestReconciler struct {
	client.Client
	Log                                logr.Logger
	Scheme                             *runtime.Scheme
	BackoffDurationProvisionerNotReady time.Duration
	BackoffDurationRequestPending      time.Duration
	recorder                           record.EventRecorder
	DefaultProviderNamespace           string
	MetricRequestsPending              *prometheus.CounterVec
	MetricIssuerNotReady               *prometheus.CounterVec
}

// +kubebuilder:rbac:groups=cert-manager.io,resources=certificaterequests,verbs=get;list;watch;update
// +kubebuilder:rbac:groups=cert-manager.io,resources=certificaterequests/status,verbs=get;update;patch

// Reconcile will read and validate a DigicertIssuer resource associated to the
// CertificateRequest resource, and it will sign the CertificateRequest with the
// provisioner in the DigicertIssuer.
func (r *CertificateRequestReconciler) Reconcile(req ctrl.Request) (ctrl.Result, error) {
	ctx := context.Background()
	log := r.Log.WithValues("certificaterequest", req.NamespacedName)

	// Fetch the CertificateRequest resource being reconciled.
	// Just ignore the request if the certificate request has been deleted.
	cr := new(cmapi.CertificateRequest)
	if err := r.Client.Get(ctx, req.NamespacedName, cr); err != nil {
		if apierrors.IsNotFound(err) {
			return ctrl.Result{}, nil
		}

		log.Error(err, "failed to retrieve CertificateRequest resource")
		return ctrl.Result{}, err
	}

	// Check the CertificateRequest's issuerRef and if it does not match the api
	// group name, log a message at a debug level and stop processing.
	if cr.Spec.IssuerRef.Group != "" && cr.Spec.IssuerRef.Group != certmanagerv1beta1.GroupVersion.Group {
		log.V(4).Info("resource does not specify an issuerRef group name that we are responsible for", "group", cr.Spec.IssuerRef.Group)
		return ctrl.Result{}, nil
	}

	// If the certificate data is already set then we skip this request as it
	// has already been completed in the past.
	if len(cr.Status.Certificate) > 0 {
		log.V(4).Info("existing certificate data found in status, skipping already completed CertificateRequest")
		return ctrl.Result{}, nil
	}

	iss := new(certmanagerv1beta1.DigicertIssuer)
	issNamespaceName := types.NamespacedName{
		Namespace: req.Namespace,
		Name:      cr.Spec.IssuerRef.Name,
	}
	if err := r.Client.Get(ctx, issNamespaceName, iss); err != nil {
		log.V(4).Info("Failed to retrieve DigicertIssuer resource, falling back to default namespace", "namespace", req.Namespace, "name", cr.Spec.IssuerRef.Name)
		issNamespaceName.Namespace = r.DefaultProviderNamespace
		err = r.Client.Get(ctx, issNamespaceName, iss)
		if err != nil {
			log.Error(err, "No DigicertIssuer resource found", "namespace", r.DefaultProviderNamespace, "name", cr.Spec.IssuerRef.Name)
			_ = r.setStatus(ctx, cr, cmmeta.ConditionFalse, cmapi.CertificateRequestReasonPending, "Failed to retrieve DigicertIssuer resource %s: %v", issNamespaceName, err)
			return ctrl.Result{}, err
		}
	}

	if !isDigicertIssuerReady(iss) {
		err := fmt.Errorf("resource %s is not ready", issNamespaceName)
		log.Error(err, "issuers is not ready")
		_ = r.setStatus(ctx, cr, cmmeta.ConditionFalse, cmapi.CertificateRequestReasonPending, "DigicertIssuer resource %s is not Ready", issNamespaceName)
		return ctrl.Result{Requeue: true, RequeueAfter: r.BackoffDurationProvisionerNotReady}, err
	}

	// Load the provisioner that will sign the CertificateRequest.
	provisioner, ok := provisioners.Load(issNamespaceName)
	if !ok {
		err := fmt.Errorf("provisioner %s not found", issNamespaceName)
		log.Error(err, "failed to load provisioner for DigicertIssuer resource")
		_ = r.setStatus(ctx, cr, cmmeta.ConditionFalse, cmapi.CertificateRequestReasonPending, "Failed to load provisioner for DigicertIssuer resource %s", issNamespaceName)
		return ctrl.Result{Requeue: true, RequeueAfter: r.BackoffDurationProvisionerNotReady}, err
	}

	// Download pending certificate
	if isCertificateRequestPending(cr) {
		log.V(4).Info("CertificateRequest is in pending state, trying to download certificate.", "name", cr.ObjectMeta.Name)
		caPEM, certPEM, err := provisioner.Download(ctx, cr)

		if err != nil || len(caPEM) < 1 || len(certPEM) < 1 {
			log.V(4).Info("Download of pending certificate failed, reqeueing.", "name", cr.ObjectMeta.Name)
			_ = r.setStatus(ctx, cr, cmmeta.ConditionFalse, cmapi.CertificateRequestReasonPending, "Certificate request pending")

			err2 := increaseCounterMetric(
				r.MetricRequestsPending,
				cr.ObjectMeta.Name,
				cr.ObjectMeta.GetAnnotations()["cert-manager.io/certificate-name"],
				cr.ObjectMeta.GetAnnotations()["cert-manager.io/private-key-secret-name"],
			)
			if err2 != nil {
				log.Error(err2, "Could not increase request pending metric.")
			}

			return ctrl.Result{Requeue: true, RequeueAfter: r.BackoffDurationRequestPending}, err
		}

		cr.Status.CA = caPEM
		cr.Status.Certificate = certPEM
		err = r.setStatus(ctx, cr, cmmeta.ConditionTrue, cmapi.CertificateRequestReasonIssued, "Certificate issued")

		return ctrl.Result{}, err
	}

	// Sign CertificateRequest.
	caPEM, certPEM, order, err := provisioner.Sign(ctx, cr)
	if err != nil {
		log.Error(err, "failed to sign certificate request")
		return ctrl.Result{}, r.setStatus(ctx, cr, cmmeta.ConditionFalse, cmapi.CertificateRequestReasonFailed, "Failed to sign certificate request: %v", err)
	}

	if order.ID > 0 {
		annotations := cr.ObjectMeta.GetAnnotations()
		annotations["cert-manager.io/digicert-order-id"] = fmt.Sprintf("%d", order.ID)
		cr.ObjectMeta.SetAnnotations(annotations)
	}

	if order.CertificateID > 0 {
		annotations := cr.ObjectMeta.GetAnnotations()
		annotations["cert-manager.io/digicert-cert-id"] = fmt.Sprintf("%d", order.CertificateID)
		cr.ObjectMeta.SetAnnotations(annotations)
	}

	if len(caPEM) > 0 && len(certPEM) > 0 {
		cr.Status.CA = caPEM
		cr.Status.Certificate = certPEM
		err = r.setStatus(ctx, cr, cmmeta.ConditionTrue, cmapi.CertificateRequestReasonIssued, "Certificate issued")
	} else if order.CertificateID > 0 {
		err = r.setStatus(ctx, cr, cmmeta.ConditionFalse, cmapi.CertificateRequestReasonPending, "Certificate request pending")
		r.Client.Update(ctx, cr) // Update annontations
		return ctrl.Result{Requeue: true, RequeueAfter: r.BackoffDurationProvisionerNotReady}, err
	} else {
		err = r.setStatus(ctx, cr, cmmeta.ConditionUnknown, cmapi.CertificateRequestReasonFailed, "Certificate request failed")
	}

	// Update annontations
	r.Client.Update(ctx, cr)

	return ctrl.Result{}, err
}

func isDigicertIssuerReady(issuer *certmanagerv1beta1.DigicertIssuer) bool {
	status := issuer.Status
	if status == nil {
		return false
	}

	for _, condition := range status.Conditions {
		if condition.Type == certmanagerv1beta1.ConditionReady && condition.Status == certmanagerv1beta1.ConditionTrue {
			return true
		}
	}

	return false
}

func isCertificateRequestPending(cr *cmapi.CertificateRequest) bool {
	status := cr.Status

	for _, condition := range status.Conditions {
		if condition.Reason == "Pending" {
			return true
		}
	}

	return false
}

func isCertificateRequestIssued(cr *cmapi.CertificateRequest) bool {
	status := cr.Status

	for _, condition := range status.Conditions {
		if condition.Reason == "Issued" {
			return true
		}
	}

	return false
}

func isCertificateRequestStatusTrue(cr *cmapi.CertificateRequest) bool {
	status := cr.Status

	for _, condition := range status.Conditions {
		if condition.Status == "True" {
			return true
		}
	}

	return false
}

func increaseCounterMetric(cv *prometheus.CounterVec, labels ...string) error {
	counter, err := cv.GetMetricWithLabelValues(labels...)

	if err != nil {
		return fmt.Errorf("Unable to increase metric %s: %s", counter, err)
	} else {
		counter.Inc()
	}

	return nil
}

// SetupWithManager initializes the CertificateRequest controller into the
// controller runtime.
func (r *CertificateRequestReconciler) SetupWithManager(mgr ctrl.Manager) error {
	filter := predicate.Funcs{
		UpdateFunc: func(e event.UpdateEvent) bool {
			cr := e.ObjectNew.(*cmapi.CertificateRequest)
			return !isCertificateRequestStatusTrue(cr) ||
				!isCertificateRequestIssued(cr) ||
				len(cr.Status.Certificate) == 0
		},
	}

	r.recorder = mgr.GetEventRecorderFor("certificateRequestController")
	return ctrl.NewControllerManagedBy(mgr).
		For(&cmapi.CertificateRequest{}).
		WithEventFilter(filter).
		Complete(r)
}

func (r *CertificateRequestReconciler) setStatus(ctx context.Context, cr *cmapi.CertificateRequest, status cmmeta.ConditionStatus, reason, message string, args ...interface{}) error {
	completeMessage := fmt.Sprintf(message, args...)
	apiutil.SetCertificateRequestCondition(cr, cmapi.CertificateRequestConditionReady, status, reason, completeMessage)

	// Fire an Event to additionally inform users of the change
	eventType := core.EventTypeNormal
	if status == cmmeta.ConditionFalse {
		eventType = core.EventTypeWarning
	}
	r.recorder.Event(cr, eventType, reason, completeMessage)

	return r.Client.Status().Update(ctx, cr)
}
