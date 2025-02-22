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

package provisioners

import (
	"sync"

	"k8s.io/apimachinery/pkg/types"
)

var collection = new(sync.Map)

func Load(namespacedName types.NamespacedName) (*CertCentral, bool) {
	v, ok := collection.Load(namespacedName)
	if !ok {
		return nil, ok
	}

	p, ok := v.(*CertCentral)
	return p, ok
}

func Store(namespacedName types.NamespacedName, provisioner *CertCentral) {
	collection.Store(namespacedName, provisioner)
}
