// SPDX-FileCopyrightText: 2025 SAP SE or an SAP affiliate company
// SPDX-License-Identifier: Apache-2.0

// SPDX-FileCopyrightText: 2022 SAP SE or an SAP affiliate company
// SPDX-License-Identifier: Apache-2.0

package k8sutils

import (
	"context"
	"time"

	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"sigs.k8s.io/controller-runtime/pkg/client"

	certmanagerv1beta1 "github.com/sapcc/digicert-issuer/apis/certmanager/v1beta1"
)

func SetDigicertIssuerStatusConditionType(ctx context.Context, k8sClient client.Client, curIssuer *certmanagerv1beta1.DigicertIssuer, statusType certmanagerv1beta1.ConditionType, status certmanagerv1beta1.ConditionStatus, reason certmanagerv1beta1.ConditionReason, message string) (*certmanagerv1beta1.DigicertIssuer, error) {
	ts := metav1.NewTime(time.Now().UTC())
	newCondition := certmanagerv1beta1.DigicertIssuerCondition{
		Type:               statusType,
		Status:             status,
		LastTransitionTime: &ts,
		Reason:             reason,
		Message:            message,
	}

	// Skip the update if the condition is already present and only the timestamp changed.
	if isIssuerHasStatusConditionIgnoreTimestamp(curIssuer.Status, newCondition) {
		return curIssuer, nil
	}

	newIssuer := curIssuer.DeepCopy()
	if newIssuer.Status == nil {
		newIssuer.Status = &certmanagerv1beta1.DigicertIssuerStatus{
			Conditions: make([]certmanagerv1beta1.DigicertIssuerCondition, 0),
		}
	}

	if len(newIssuer.Status.Conditions) == 0 {
		newIssuer.Status.Conditions = []certmanagerv1beta1.DigicertIssuerCondition{newCondition}
		return patchDigicertIssuerStatus(ctx, k8sClient, curIssuer, newIssuer)
	}

	for idx, curCondition := range newIssuer.Status.Conditions {
		if curCondition.Type == newCondition.Type {
			newIssuer.Status.Conditions[idx] = newCondition
		}
	}

	return patchDigicertIssuerStatus(ctx, k8sClient, curIssuer, newIssuer)
}

func EnsureDigicertIssuerStatusInitialized(ctx context.Context, k8sClient client.Client, issuer *certmanagerv1beta1.DigicertIssuer) (*certmanagerv1beta1.DigicertIssuer, error) {
	if isDigicertIssuerReady(issuer) {
		return issuer, nil
	}

	return SetDigicertIssuerStatusConditionType(
		ctx, k8sClient, issuer, certmanagerv1beta1.ConditionReady, certmanagerv1beta1.ConditionFalse, "", "",
	)
}

func isDigicertIssuerReady(issuer *certmanagerv1beta1.DigicertIssuer) bool {
	return isIssuerHasStatusConditionIgnoreTimestamp(issuer.Status, certmanagerv1beta1.DigicertIssuerCondition{
		Type:   certmanagerv1beta1.ConditionReady,
		Status: certmanagerv1beta1.ConditionTrue,
	})
}

func patchDigicertIssuerStatus(ctx context.Context, k8sClient client.Client, curIssuer, newIssuer *certmanagerv1beta1.DigicertIssuer) (*certmanagerv1beta1.DigicertIssuer, error) {
	patch := client.MergeFrom(curIssuer)
	if err := k8sClient.Status().Patch(ctx, newIssuer, patch); err != nil {
		return curIssuer, err
	}

	return newIssuer, nil
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
