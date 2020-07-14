package provisioners

import (
	"context"
	"digicert-issuer/apis/certmanager/v1beta1"
	"sync"

	certmanagerv1alpha2 "github.com/jetstack/cert-manager/pkg/apis/certmanager/v1alpha2"
	"k8s.io/apimachinery/pkg/types"
)

var collection = new(sync.Map)

type CertCentral struct {}

func New(iss *v1beta1.DigicertIssuer) (*CertCentral, error) {
	//TODO:
	return nil, nil
}

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

func (c *CertCentral) Sign(ctx context.Context, cr *certmanagerv1alpha2.CertificateRequest) ([]byte, []byte, error) {
	//TODO:
	return nil, nil, nil
}
