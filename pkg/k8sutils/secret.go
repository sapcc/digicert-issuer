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
	"fmt"
	"strings"

	corev1 "k8s.io/api/core/v1"
	"sigs.k8s.io/controller-runtime/pkg/client"
)

func GetSecretData(k8sClient client.Client, secretNamespace, secretName, secretDataKey string) (string, error) {
	s := new(corev1.Secret)
	if err := k8sClient.Get(context.Background(), client.ObjectKey{
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
