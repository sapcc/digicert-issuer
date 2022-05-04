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

func SetDigicertIssuerStatusConditionType(k8sClient client.Client, cur *certmanagerv1beta1.DigicertIssuer, statusType certmanagerv1beta1.ConditionType, status certmanagerv1beta1.ConditionStatus, reason certmanagerv1beta1.ConditionReason, message string) (*certmanagerv1beta1.DigicertIssuer, error) {
	ts := metav1.NewTime(time.Now().UTC())
	newCondition := certmanagerv1beta1.DigicertIssuerCondition{
		Type:               statusType,
		Status:             status,
		LastTransitionTime: &ts,
		Reason:             reason,
		Message:            message,
	}

	// Skip the update if the condition is already present and only the timestamp changed.
	if isIssuerHasStatusConditionIgnoreTimestamp(cur.Status, newCondition) {
		return cur, nil
	}

	new := cur.DeepCopy()
	if new.Status == nil {
		new.Status = &certmanagerv1beta1.DigicertIssuerStatus{
			Conditions: make([]certmanagerv1beta1.DigicertIssuerCondition, 0),
		}
	}

	if new.Status.Conditions == nil || len(new.Status.Conditions) == 0 {
		new.Status.Conditions = []certmanagerv1beta1.DigicertIssuerCondition{newCondition}
		return patchDigicertIssuerStatus(k8sClient, cur, new)
	}

	for idx, curCondition := range new.Status.Conditions {
		if curCondition.Type == newCondition.Type {
			new.Status.Conditions[idx] = newCondition
		}
	}

	return patchDigicertIssuerStatus(k8sClient, cur, new)
}

func EnsureDigicertIssuerStatusInitialized(k8sClient client.Client, issuer *certmanagerv1beta1.DigicertIssuer) (*certmanagerv1beta1.DigicertIssuer, error) {
	if isDigicertIssuerReady(issuer) {
		return issuer, nil
	}

	return SetDigicertIssuerStatusConditionType(
		k8sClient, issuer, certmanagerv1beta1.ConditionReady, certmanagerv1beta1.ConditionFalse, "", "",
	)
}

func isDigicertIssuerReady(issuer *certmanagerv1beta1.DigicertIssuer) bool {
	return isIssuerHasStatusConditionIgnoreTimestamp(issuer.Status, certmanagerv1beta1.DigicertIssuerCondition{
		Type:   certmanagerv1beta1.ConditionReady,
		Status: certmanagerv1beta1.ConditionTrue,
	})
}

func patchDigicertIssuerStatus(k8sClient client.Client, cur, new *certmanagerv1beta1.DigicertIssuer) (*certmanagerv1beta1.DigicertIssuer, error) {
	patch := client.MergeFrom(cur)
	if err := k8sClient.Status().Patch(context.Background(), new, patch); err != nil {
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
