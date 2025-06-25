// SPDX-FileCopyrightText: 2025 SAP SE or an SAP affiliate company
// SPDX-License-Identifier: Apache-2.0

// SPDX-FileCopyrightText: 2022 SAP SE or an SAP affiliate company
// SPDX-License-Identifier: Apache-2.0

package certmanager

import (
	"context"
	"errors"

	"github.com/go-logr/logr"
	"github.com/hashicorp/go-multierror"
	"k8s.io/client-go/tools/record"
	ctrl "sigs.k8s.io/controller-runtime"
	"sigs.k8s.io/controller-runtime/pkg/client"

	certmanagerv1beta1 "github.com/sapcc/digicert-issuer/apis/certmanager/v1beta1"
	"github.com/sapcc/digicert-issuer/pkg/k8sutils"
	"github.com/sapcc/digicert-issuer/pkg/provisioners"
)

// DigicertIssuerReconciler reconciles a DigicertIssuer object
type DigicertIssuerReconciler struct {
	client.Client
	log      logr.Logger
	recorder record.EventRecorder
}

// +kubebuilder:rbac:groups=certmanager.cloud.sap,resources=digicertissuers,verbs=get;list;watch;create;update;patch;delete
// +kubebuilder:rbac:groups=certmanager.cloud.sap,resources=digicertissuers/status,verbs=get;update;patch
// +kubebuilder:rbac:groups=core,resources=secrets,verbs=get;list;watch
// +kubebuilder:rbac:groups="",resources=events,verbs=create;patch

func (r *DigicertIssuerReconciler) Reconcile(ctx context.Context, req ctrl.Request) (ctrl.Result, error) {
	logger := r.log.WithValues("digicertissuer", req.NamespacedName)

	issuer := new(certmanagerv1beta1.DigicertIssuer)
	if err := r.Get(ctx, req.NamespacedName, issuer); err != nil {
		logger.Error(err, "failed to get issuer")
		return ctrl.Result{}, client.IgnoreNotFound(err)
	}

	issuer, err := k8sutils.EnsureDigicertIssuerStatusInitialized(ctx, r.Client, issuer)
	if err != nil {
		logger.Error(err, "failed to initialize issuer status")
	}

	if err := validateDigicertIssuerSpec(issuer.Spec); err != nil {
		if _, err := k8sutils.SetDigicertIssuerStatusConditionType(
			ctx, r.Client, issuer, certmanagerv1beta1.ConditionConfigurationError, certmanagerv1beta1.ConditionTrue,
			certmanagerv1beta1.ConditionReasonInvalidIssuerSpec, err.Error(),
		); err != nil {
			logger.Error(err, "failed to set issuer status condition")
		}
		logger.Error(err, "issuer.spec is invalid")
		return ctrl.Result{}, err
	}
	if _, err := k8sutils.SetDigicertIssuerStatusConditionType(
		ctx, r.Client, issuer, certmanagerv1beta1.ConditionConfigurationError, certmanagerv1beta1.ConditionFalse, "", "",
	); err != nil {
		logger.Error(err, "failed to clear issuer status condition")
	}

	secretRef := issuer.Spec.Provisioner.APITokenReference
	digicertAPIToken, err := k8sutils.GetSecretData(ctx, r.Client, issuer.GetNamespace(), secretRef.Name, secretRef.Key)
	if err != nil {
		logger.Error(err, "failed to get provisioner secret containing the API token")
		if _, err := k8sutils.SetDigicertIssuerStatusConditionType(
			ctx, r.Client, issuer, certmanagerv1beta1.ConditionConfigurationError, certmanagerv1beta1.ConditionTrue,
			certmanagerv1beta1.ConditionReasonSecretNotFoundOrEmpty, err.Error(),
		); err != nil {
			logger.Error(err, "failed to set issuer status condition")
		}
		return ctrl.Result{}, err
	}
	if _, err := k8sutils.SetDigicertIssuerStatusConditionType(
		ctx, r.Client, issuer, certmanagerv1beta1.ConditionConfigurationError, certmanagerv1beta1.ConditionFalse, "", "",
	); err != nil {
		logger.Error(err, "failed to clear issuer status condition")
	}

	prov, err := provisioners.New(issuer, digicertAPIToken)
	if err != nil {
		logger.Error(err, "failed to initialize provisioner")
		return ctrl.Result{}, err
	}

	provisioners.Store(req.NamespacedName, prov)
	logger.Info("provisioner is ready", "name", prov.GetName())

	_, err = k8sutils.SetDigicertIssuerStatusConditionType(
		ctx, r.Client, issuer, certmanagerv1beta1.ConditionReady, certmanagerv1beta1.ConditionTrue, "", "",
	)
	return ctrl.Result{}, err
}

func (r *DigicertIssuerReconciler) SetupWithManager(mgr ctrl.Manager) error {
	r.recorder = mgr.GetEventRecorderFor("digicertIssuer")
	r.log = mgr.GetLogger().WithName("controllers").WithName("DigicertIssuer")
	r.Client = mgr.GetClient()
	return ctrl.NewControllerManagedBy(mgr).
		For(&certmanagerv1beta1.DigicertIssuer{}).
		Complete(r)
}

func validateDigicertIssuerSpec(issuerSpec certmanagerv1beta1.DigicertIssuerSpec) error {
	var errs error

	provisionerSpec := issuerSpec.Provisioner
	if provisionerSpec.APITokenReference.Name == "" {
		errs = multierror.Append(errs, errors.New("spec.provisioner.apiTokenReference.name missing"))
	}
	if provisionerSpec.APITokenReference.Key == "" {
		errs = multierror.Append(errs, errors.New("spec.provisioner.apiTokenReference.key missing"))
	}
	if len(provisionerSpec.OrganizationUnits) == 0 {
		errs = multierror.Append(errs, errors.New("spec.provisioner.organizationalUnits missing"))
	}
	if provisionerSpec.OrganizationID == nil && provisionerSpec.OrganizationName == "" {
		errs = multierror.Append(errs, errors.New("spec.provisioner.organizationID or spec.provisioner.organizationName missing"))
	}

	return errs
}
