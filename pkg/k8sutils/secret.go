// SPDX-FileCopyrightText: 2022 SAP SE or an SAP affiliate company
// SPDX-License-Identifier: Apache-2.0

package k8sutils

import (
	"context"
	"fmt"
	"strings"

	corev1 "k8s.io/api/core/v1"
	"sigs.k8s.io/controller-runtime/pkg/client"
)

func GetSecretData(ctx context.Context, k8sClient client.Client, secretNamespace, secretName, secretDataKey string) (string, error) {
	s := new(corev1.Secret)
	if err := k8sClient.Get(ctx, client.ObjectKey{
		Namespace: secretNamespace,
		Name:      secretName,
	}, s); err != nil {
		return "", err
	}

	if s.Data == nil {
		return "", fmt.Errorf("secret %s/%s is empty", secretNamespace, secretName)
	}

	valByte, ok := s.Data[secretDataKey]
	if !ok {
		return "", fmt.Errorf("secret %s/%s does not contain key %s", secretNamespace, secretName, secretDataKey)
	}

	valStr := string(valByte)
	valStr = strings.TrimSpace(valStr)
	valStr = strings.TrimSuffix(valStr, "\n")
	return valStr, nil
}
