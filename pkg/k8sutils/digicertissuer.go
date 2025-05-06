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

package k8sutils

import (
	"context"
	"time"

	certmanagerv1beta1 "github.com/sapcc/digicert-issuer/apis/certmanager/v1beta1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"sigs.k8s.io/controller-runtime/pkg/client"
)

type Issuer interface {
	Get(ctx context.Context, client client.Client, key client.ObjectKey) error
	Kind() string
	Spec() certmanagerv1beta1.DigicertIssuerSpec
	Status() *certmanagerv1beta1.DigicertIssuerStatus
	SetStatus(*certmanagerv1beta1.DigicertIssuerStatus)
	PatchStatus(ctx context.Context, client client.Client, newStatus *certmanagerv1beta1.DigicertIssuerStatus) (Issuer, error)
}

// DigicertIssuer allows an implementation for the Issuer interface without adding unnecessary dependencies to package API
type DigicertIssuer struct {
	certmanagerv1beta1.DigicertIssuer
}

var _ Issuer = &DigicertIssuer{}

func (iss *DigicertIssuer) Get(ctx context.Context, client client.Client, key client.ObjectKey) error {
	return client.Get(ctx, key, &iss.DigicertIssuer)
}

func (iss *DigicertIssuer) Kind() string {
	return "DigicertIssuer"
}

func (iss *DigicertIssuer) Status() *certmanagerv1beta1.DigicertIssuerStatus {
	return iss.DigicertIssuer.Status
}

func (iss *DigicertIssuer) Spec() certmanagerv1beta1.DigicertIssuerSpec {
	return iss.DigicertIssuer.Spec
}

func (iss *DigicertIssuer) SetStatus(status *certmanagerv1beta1.DigicertIssuerStatus) {
	iss.DigicertIssuer.Status = status
}

func (iss *DigicertIssuer) PatchStatus(ctx context.Context, k8sClient client.Client, newStatus *certmanagerv1beta1.DigicertIssuerStatus) (Issuer, error) {
	patch := client.MergeFrom(iss)
	newIssuer := iss.DeepCopy()
	newIssuer.Status = newStatus
	if err := k8sClient.Status().Patch(ctx, newIssuer, patch); err != nil {
		return iss, err
	}

	return &DigicertIssuer{*newIssuer}, nil
}

func NewDigicertIssuer() Issuer {
	return new(DigicertIssuer)
}

// ClusterDigicertIssuer allows an implementation for the Issuer interface without adding unnecessary dependencies to package API
type ClusterDigicertIssuer struct {
	certmanagerv1beta1.ClusterDigicertIssuer
}

var _ Issuer = &ClusterDigicertIssuer{}

func (iss *ClusterDigicertIssuer) Get(ctx context.Context, client client.Client, key client.ObjectKey) error {
	return client.Get(ctx, key, &iss.ClusterDigicertIssuer)
}

func NewClusterDigicertIssuer() Issuer {
	return new(ClusterDigicertIssuer)
}

func (iss *ClusterDigicertIssuer) Kind() string {
	return "ClusterDigicertIssuer"
}

func (iss *ClusterDigicertIssuer) Spec() certmanagerv1beta1.DigicertIssuerSpec {
	return iss.ClusterDigicertIssuer.Spec
}

func (iss *ClusterDigicertIssuer) Status() *certmanagerv1beta1.DigicertIssuerStatus {
	return iss.ClusterDigicertIssuer.Status
}

func (iss *ClusterDigicertIssuer) SetStatus(status *certmanagerv1beta1.DigicertIssuerStatus) {
	iss.ClusterDigicertIssuer.Status = status
}

func (iss *ClusterDigicertIssuer) PatchStatus(ctx context.Context, k8sClient client.Client, newStatus *certmanagerv1beta1.DigicertIssuerStatus) (Issuer, error) {
	patch := client.MergeFrom(iss)
	newIssuer := iss.DeepCopy()
	newIssuer.Status = newStatus
	if err := k8sClient.Status().Patch(ctx, newIssuer, patch); err != nil {
		return iss, err
	}

	return &ClusterDigicertIssuer{*newIssuer}, nil
}

func SetDigicertIssuerStatusConditionType(ctx context.Context, k8sClient client.Client, cur Issuer, statusType certmanagerv1beta1.ConditionType, status certmanagerv1beta1.ConditionStatus, reason certmanagerv1beta1.ConditionReason, message string) (Issuer, error) {
	ts := metav1.NewTime(time.Now().UTC())
	newCondition := certmanagerv1beta1.DigicertIssuerCondition{
		Type:               statusType,
		Status:             status,
		LastTransitionTime: &ts,
		Reason:             reason,
		Message:            message,
	}

	// Skip the update if the condition is already present and only the timestamp changed.
	if isIssuerHasStatusConditionIgnoreTimestamp(cur.Status(), newCondition) {
		return cur, nil
	}

	newStatus := cur.Status().DeepCopy()
	if newStatus == nil {
		newStatus = &certmanagerv1beta1.DigicertIssuerStatus{
			Conditions: make([]certmanagerv1beta1.DigicertIssuerCondition, 0),
		}
	}

	if newStatus.Conditions == nil || len(newStatus.Conditions) == 0 {
		newStatus.Conditions = []certmanagerv1beta1.DigicertIssuerCondition{newCondition}
		return cur.PatchStatus(ctx, k8sClient, newStatus)
	}

	for idx, curCondition := range newStatus.Conditions {
		if curCondition.Type == newCondition.Type {
			newStatus.Conditions[idx] = newCondition
		}
	}

	return cur.PatchStatus(ctx, k8sClient, newStatus)
}

func EnsureDigicertIssuerStatusInitialized(ctx context.Context, k8sClient client.Client, issuer Issuer) (Issuer, error) {
	if isDigicertIssuerReady(issuer) {
		return issuer, nil
	}

	return SetDigicertIssuerStatusConditionType(
		ctx, k8sClient, issuer, certmanagerv1beta1.ConditionReady, certmanagerv1beta1.ConditionFalse, "", "",
	)
}

func isDigicertIssuerReady(issuer Issuer) bool {
	return isIssuerHasStatusConditionIgnoreTimestamp(issuer.Status(), certmanagerv1beta1.DigicertIssuerCondition{
		Type:   certmanagerv1beta1.ConditionReady,
		Status: certmanagerv1beta1.ConditionTrue,
	})
}

func patchDigicertIssuerStatus(ctx context.Context, k8sClient client.Client, cur, new *certmanagerv1beta1.DigicertIssuer) (*certmanagerv1beta1.DigicertIssuer, error) {
	patch := client.MergeFrom(cur)
	if err := k8sClient.Status().Patch(ctx, new, patch); err != nil {
		return cur, err
	}

	return new, nil
}

func isIssuerHasStatusConditionIgnoreTimestamp(issuerStatus *certmanagerv1beta1.DigicertIssuerStatus, condition certmanagerv1beta1.DigicertIssuerCondition) bool {
	if issuerStatus == nil || issuerStatus.Conditions == nil || len(issuerStatus.Conditions) == 0 {
		return false
	}

	for _, cond := range issuerStatus.Conditions {
		if cond.Status == condition.Status && cond.Type == condition.Type && cond.Reason == condition.Reason && cond.Message == condition.Message {
			return true
		}
	}
	return false
}
