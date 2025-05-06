// SPDX-FileCopyrightText: 2025 SAP SE or an SAP affiliate company and sapcc contributors
// SPDX-License-Identifier: Apache-2.0

/*
Copyright 2022 SAP SE.

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
	"errors"
	"fmt"

	"github.com/go-logr/logr"
	"github.com/hashicorp/go-multierror"
	certmanagerv1beta1 "github.com/sapcc/digicert-issuer/apis/certmanager/v1beta1"
	"github.com/sapcc/digicert-issuer/pkg/k8sutils"
	"github.com/sapcc/digicert-issuer/pkg/provisioners"
	"k8s.io/client-go/tools/record"
	ctrl "sigs.k8s.io/controller-runtime"
	"sigs.k8s.io/controller-runtime/pkg/client"
)

// DigicertIssuerReconciler reconciles a DigicertIssuer object
type DigicertIssuerReconciler struct {
	client.Client
	log                    logr.Logger
	recorder               record.EventRecorder
	clusterIssuerNamespace string
}

func NewDigicertIssuerReconciler(clusterIssuerNamespace string) *DigicertIssuerReconciler {
	return &DigicertIssuerReconciler{
		clusterIssuerNamespace: clusterIssuerNamespace,
	}
}

// +kubebuilder:rbac:groups=certmanager.cloud.sap,resources=digicertissuers,verbs=get;list;watch;create;update;patch;delete
// +kubebuilder:rbac:groups=certmanager.cloud.sap,resources=digicertissuers/status,verbs=get;update;patch
// +kubebuilder:rbac:groups=certmanager.cloud.sap,resources=clusterdigicertissuers,verbs=get;list;watch;create;update;patch;delete
// +kubebuilder:rbac:groups=certmanager.cloud.sap,resources=clusterdigicertissuers/status,verbs=get;update;patch
// +kubebuilder:rbac:groups=core,resources=secrets,verbs=get;list;watch
// +kubebuilder:rbac:groups="",resources=events,verbs=create;patch

func (r *DigicertIssuerReconciler) Reconcile(ctx context.Context, req ctrl.Request) (ctrl.Result, error) {
	var issuer k8sutils.Issuer
	secretNamespace := req.Namespace
	if req.Namespace == "" {
		issuer = new(k8sutils.ClusterDigicertIssuer)
		secretNamespace = r.clusterIssuerNamespace
	} else {
		issuer = new(k8sutils.DigicertIssuer)
	}
	logger := r.log.WithValues(issuer.Kind(), req.NamespacedName)

	if err := issuer.Get(ctx, r.Client, req.NamespacedName); err != nil {
		logger.Error(err, "failed to get issuer")
		return ctrl.Result{}, client.IgnoreNotFound(err)
	}

	issuer, err := k8sutils.EnsureDigicertIssuerStatusInitialized(ctx, r.Client, issuer)
	if err != nil {
		logger.Error(err, "failed to initialize issuer status")
	}

	if err := validateDigicertIssuerSpec(issuer.Spec()); err != nil {
		k8sutils.SetDigicertIssuerStatusConditionType(
			ctx, r.Client, issuer, certmanagerv1beta1.ConditionConfigurationError, certmanagerv1beta1.ConditionTrue,
			certmanagerv1beta1.ConditionReasonInvalidIssuerSpec, err.Error(),
		)
		logger.Error(err, "issuer.spec is invalid")
		return ctrl.Result{}, err
	}
	k8sutils.SetDigicertIssuerStatusConditionType(
		ctx, r.Client, issuer, certmanagerv1beta1.ConditionConfigurationError, certmanagerv1beta1.ConditionFalse, "", "",
	)

	secretRef := issuer.Spec().Provisioner.APITokenReference
	digicertAPIToken, err := k8sutils.GetSecretData(ctx, r.Client, secretNamespace, secretRef.Name, secretRef.Key)
	if err != nil {
		logger.Error(err, "failed to get provisioner secret containing the API token")
		k8sutils.SetDigicertIssuerStatusConditionType(
			ctx, r.Client, issuer, certmanagerv1beta1.ConditionConfigurationError, certmanagerv1beta1.ConditionTrue,
			certmanagerv1beta1.ConditionReasonSecretNotFoundOrEmpty, err.Error(),
		)
		return ctrl.Result{}, err
	}
	k8sutils.SetDigicertIssuerStatusConditionType(
		ctx, r.Client, issuer, certmanagerv1beta1.ConditionConfigurationError, certmanagerv1beta1.ConditionFalse, "", "",
	)

	prov, err := provisioners.New(fmt.Sprintf("%s/%s", req.Namespace, req.Name), issuer.Spec(), digicertAPIToken)
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

func (r *DigicertIssuerReconciler) SetupWithManagerClusterIssuer(mgr ctrl.Manager) error {
	r.recorder = mgr.GetEventRecorderFor("clusterDigicertIssuer")
	r.log = mgr.GetLogger().WithName("controllers").WithName("ClusterDigicertIssuer")
	r.Client = mgr.GetClient()
	return ctrl.NewControllerManagedBy(mgr).
		For(&certmanagerv1beta1.ClusterDigicertIssuer{}).
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
	if provisionerSpec.OrganizationUnits == nil || len(provisionerSpec.OrganizationUnits) == 0 {
		errs = multierror.Append(errs, errors.New("spec.provisioner.organizationalUnits missing"))
	}
	if provisionerSpec.OrganizationID == nil && provisionerSpec.OrganizationName == "" {
		errs = multierror.Append(errs, errors.New("spec.provisioner.organizationID or spec.provisioner.organizationName missing"))
	}

	return errs
}
