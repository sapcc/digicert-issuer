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

	new := cur.DeepCopy()
	if new.Status == nil {
		new.Status = &certmanagerv1beta1.DigicertIssuerStatus{
			Conditions: make([]certmanagerv1beta1.DigicertIssuerCondition, 0),
		}
	}

	for idx, curCondition := range new.Status.Conditions {
		if curCondition.Type == newCondition.Type {
			new.Status.Conditions[idx] = newCondition
			continue
		}
		new.Status.Conditions[idx].Status = certmanagerv1beta1.ConditionFalse
	}

	return patchDigicertIssuerStatus(k8sClient, cur, new)
}

func patchDigicertIssuerStatus(k8sClient client.Client, cur, new *certmanagerv1beta1.DigicertIssuer) (*certmanagerv1beta1.DigicertIssuer, error) {
	patch := client.MergeFrom(cur)
	if err := k8sClient.Status().Patch(context.Background(), new, patch); err != nil {
		return cur, err
	}

	return new, nil
}
