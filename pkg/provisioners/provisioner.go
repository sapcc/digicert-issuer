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
