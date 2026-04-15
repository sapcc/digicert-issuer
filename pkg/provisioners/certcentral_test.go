// SPDX-FileCopyrightText: 2025 SAP SE or an SAP affiliate company and sapcc contributors
// SPDX-License-Identifier: Apache-2.0

package provisioners

import (
	"context"
	"crypto/rand"
	"crypto/rsa"
	"crypto/x509"
	"crypto/x509/pkix"
	"encoding/pem"
	"math/big"
	"reflect"
	"testing"
	"time"

	certmanagerv1 "github.com/cert-manager/cert-manager/pkg/apis/certmanager/v1"
	certcentral "github.com/sapcc/go-certcentral"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

type mockCertCentralClient struct {
	submitOrder *certcentral.Order
	submitErr   error
	chain       []certcentral.CertificateChain
	chainErr    error
}

func (f *mockCertCentralClient) SubmitOrder(order certcentral.Order, orderType certcentral.OrderType) (*certcentral.Order, error) {
	return f.submitOrder, f.submitErr
}

func (f *mockCertCentralClient) GetCertificateChain(certID string) ([]certcentral.CertificateChain, error) {
	return f.chain, f.chainErr
}

type chainFixture struct {
	requestedCN       string
	preferredRoot     string
	globalRoot        string
	regularBundle     []*x509.Certificate
	crossSignedBundle []*x509.Certificate
}

func TestCertCentralSignPreferredChain(t *testing.T) {
	fixture := buildChainFixture(t)
	csrPEM := createCSR(t, fixture.requestedCN)

	tests := []struct {
		name                    string
		bundle                  []*x509.Certificate
		preferredCN             string
		expectedIntermediateCNs []string
		expectedRootCNs         []string
	}{
		{
			name:                    "regular_select_preferred",
			bundle:                  fixture.regularBundle,
			preferredCN:             fixture.preferredRoot,
			expectedIntermediateCNs: []string{fixture.requestedCN, "Test Intermediate"},
			expectedRootCNs:         []string{fixture.preferredRoot},
		},
		{
			name:                    "regular_preferred_not_found",
			bundle:                  fixture.regularBundle,
			preferredCN:             "Missing Root",
			expectedIntermediateCNs: []string{fixture.requestedCN, "Test Intermediate"},
			expectedRootCNs:         []string{fixture.preferredRoot},
		},
		{
			name:                    "regular_preferred_empty",
			bundle:                  fixture.regularBundle,
			preferredCN:             "",
			expectedIntermediateCNs: []string{fixture.requestedCN, "Test Intermediate"},
			expectedRootCNs:         []string{fixture.preferredRoot},
		},
		{
			name:                    "cross_select_preferred_root_chain",
			bundle:                  fixture.crossSignedBundle,
			preferredCN:             fixture.preferredRoot,
			expectedIntermediateCNs: []string{fixture.requestedCN, "Test Intermediate"},
			expectedRootCNs:         []string{fixture.preferredRoot},
		},
		{
			name:                    "cross_select_global_root_chain",
			bundle:                  fixture.crossSignedBundle,
			preferredCN:             fixture.globalRoot,
			expectedIntermediateCNs: []string{fixture.requestedCN, "Test Intermediate", fixture.preferredRoot},
			expectedRootCNs:         []string{fixture.globalRoot},
		},
		{
			name:                    "cross_preferred_not_found",
			bundle:                  fixture.crossSignedBundle,
			preferredCN:             "Missing Root",
			expectedIntermediateCNs: []string{fixture.requestedCN, "Test Intermediate", fixture.preferredRoot},
			expectedRootCNs:         []string{fixture.preferredRoot, fixture.globalRoot},
		},
		{
			name:                    "cross_preferred_empty",
			bundle:                  fixture.crossSignedBundle,
			preferredCN:             "",
			expectedIntermediateCNs: []string{fixture.requestedCN, "Test Intermediate", fixture.preferredRoot},
			expectedRootCNs:         []string{fixture.preferredRoot, fixture.globalRoot},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			mockCertCentralClient := &mockCertCentralClient{submitOrder: &certcentral.Order{CertificateChain: constructCertificateChain(t, tt.bundle)}}
			provisioner := &CertCentral{client: mockCertCentralClient, preferredChain: tt.preferredCN}

			cr := &certmanagerv1.CertificateRequest{
				ObjectMeta: metav1.ObjectMeta{Name: "sign-test"},
				Spec:       certmanagerv1.CertificateRequestSpec{Request: csrPEM},
			}

			caPEM, tlsPEM, _, err := provisioner.Sign(context.Background(), cr)
			if err != nil {
				t.Fatalf("Sign returned error: %v", err)
			}

			assertChainCNs(t, tlsPEM, caPEM, tt.expectedIntermediateCNs, tt.expectedRootCNs)
		})
	}
}

func TestCertCentralDownloadPreferredChain(t *testing.T) {
	fixture := buildChainFixture(t)
	csrPEM := createCSR(t, fixture.requestedCN)

	tests := []struct {
		name                    string
		bundle                  []*x509.Certificate
		preferredCN             string
		expectedIntermediateCNs []string
		expectedRootCNs         []string
	}{
		{
			name:                    "regular_select_preferred",
			bundle:                  fixture.regularBundle,
			preferredCN:             fixture.preferredRoot,
			expectedIntermediateCNs: []string{fixture.requestedCN, "Test Intermediate"},
			expectedRootCNs:         []string{fixture.preferredRoot},
		},
		{
			name:                    "regular_preferred_not_found",
			bundle:                  fixture.regularBundle,
			preferredCN:             "Missing Root",
			expectedIntermediateCNs: []string{fixture.requestedCN, "Test Intermediate"},
			expectedRootCNs:         []string{fixture.preferredRoot},
		},
		{
			name:                    "regular_preferred_empty",
			bundle:                  fixture.regularBundle,
			preferredCN:             "",
			expectedIntermediateCNs: []string{fixture.requestedCN, "Test Intermediate"},
			expectedRootCNs:         []string{fixture.preferredRoot},
		},
		{
			name:                    "cross_select_preferred_root_chain",
			bundle:                  fixture.crossSignedBundle,
			preferredCN:             fixture.preferredRoot,
			expectedIntermediateCNs: []string{fixture.requestedCN, "Test Intermediate"},
			expectedRootCNs:         []string{fixture.preferredRoot},
		},
		{
			name:                    "cross_select_global_root_chain",
			bundle:                  fixture.crossSignedBundle,
			preferredCN:             fixture.globalRoot,
			expectedIntermediateCNs: []string{fixture.requestedCN, "Test Intermediate", fixture.preferredRoot},
			expectedRootCNs:         []string{fixture.globalRoot},
		},
		{
			name:                    "cross_preferred_not_found",
			bundle:                  fixture.crossSignedBundle,
			preferredCN:             "Missing Root",
			expectedIntermediateCNs: []string{fixture.requestedCN, "Test Intermediate", fixture.preferredRoot},
			expectedRootCNs:         []string{fixture.preferredRoot, fixture.globalRoot},
		},
		{
			name:                    "cross_preferred_empty",
			bundle:                  fixture.crossSignedBundle,
			preferredCN:             "",
			expectedIntermediateCNs: []string{fixture.requestedCN, "Test Intermediate", fixture.preferredRoot},
			expectedRootCNs:         []string{fixture.preferredRoot, fixture.globalRoot},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			mockCertCentralClient := &mockCertCentralClient{chain: constructCertificateChain(t, tt.bundle)}
			provisioner := &CertCentral{client: mockCertCentralClient, preferredChain: tt.preferredCN}

			cr := &certmanagerv1.CertificateRequest{
				ObjectMeta: metav1.ObjectMeta{
					Name:        "download-test",
					Annotations: map[string]string{"certmanager.cloud.sap/digicert-cert-id": "12345"},
				},
				Spec: certmanagerv1.CertificateRequestSpec{Request: csrPEM},
			}

			caPEM, tlsPEM, err := provisioner.Download(context.Background(), cr)
			if err != nil {
				t.Fatalf("Download returned error: %v", err)
			}

			assertChainCNs(t, tlsPEM, caPEM, tt.expectedIntermediateCNs, tt.expectedRootCNs)
		})
	}
}

func buildChainFixture(t *testing.T) chainFixture {
	t.Helper()

	const (
		leafCN        = "leaf.test.local"
		preferredRoot = "Preferred Root"
		globalRootCN  = "Global Root"
	)

	globalRootKey := generateRSAKey(t)
	globalRootCert := generateCert(t, certSpec{cn: globalRootCN, isCA: true}, nil, nil, globalRootKey)

	preferredRootKey := generateRSAKey(t)
	preferredRootSelf := generateCert(t, certSpec{cn: preferredRoot, isCA: true}, nil, nil, preferredRootKey)

	preferredRootCross := generateCert(t, certSpec{cn: preferredRoot, isCA: true}, globalRootCert, globalRootKey, preferredRootKey)

	intermediateKey := generateRSAKey(t)
	intermediate := generateCert(t, certSpec{cn: "Test Intermediate", isCA: true}, preferredRootSelf, preferredRootKey, intermediateKey)

	leafKey := generateRSAKey(t)
	leaf := generateCert(t, certSpec{cn: leafCN, isCA: false}, intermediate, intermediateKey, leafKey)

	return chainFixture{
		requestedCN:       leafCN,
		preferredRoot:     preferredRoot,
		globalRoot:        globalRootCN,
		regularBundle:     []*x509.Certificate{leaf, intermediate, preferredRootSelf},
		crossSignedBundle: []*x509.Certificate{leaf, intermediate, preferredRootCross, preferredRootSelf, globalRootCert},
	}
}

type certSpec struct {
	cn   string
	isCA bool
}

func generateRSAKey(t *testing.T) *rsa.PrivateKey {
	t.Helper()
	key, err := rsa.GenerateKey(rand.Reader, 2048)
	if err != nil {
		t.Fatalf("generate key: %v", err)
	}
	return key
}

func generateCert(t *testing.T, spec certSpec, parent *x509.Certificate, parentKey *rsa.PrivateKey, subjectKey *rsa.PrivateKey) *x509.Certificate {
	t.Helper()

	now := time.Now().UTC()
	tmpl := &x509.Certificate{
		SerialNumber:          big.NewInt(now.UnixNano()),
		Subject:               pkix.Name{CommonName: spec.cn},
		NotBefore:             now.Add(-time.Hour),
		NotAfter:              now.Add(365 * 24 * time.Hour),
		BasicConstraintsValid: true,
		IsCA:                  spec.isCA,
	}

	if spec.isCA {
		tmpl.KeyUsage = x509.KeyUsageCertSign | x509.KeyUsageCRLSign
	} else {
		tmpl.KeyUsage = x509.KeyUsageDigitalSignature | x509.KeyUsageKeyEncipherment
		tmpl.ExtKeyUsage = []x509.ExtKeyUsage{x509.ExtKeyUsageServerAuth}
	}

	if parent == nil {
		parent = tmpl
	}
	if parentKey == nil {
		parentKey = subjectKey
	}

	der, err := x509.CreateCertificate(rand.Reader, tmpl, parent, &subjectKey.PublicKey, parentKey)
	if err != nil {
		t.Fatalf("create certificate %q: %v", spec.cn, err)
	}

	crt, err := x509.ParseCertificate(der)
	if err != nil {
		t.Fatalf("parse certificate %q: %v", spec.cn, err)
	}

	return crt
}

func createCSR(t *testing.T, commonName string) []byte {
	t.Helper()
	key := generateRSAKey(t)
	csrDER, err := x509.CreateCertificateRequest(rand.Reader, &x509.CertificateRequest{
		Subject:  pkix.Name{CommonName: commonName},
		DNSNames: []string{commonName},
	}, key)
	if err != nil {
		t.Fatalf("create csr: %v", err)
	}
	return pem.EncodeToMemory(&pem.Block{Type: blockTypeCertificateRequest, Bytes: csrDER})
}

func constructCertificateChain(t *testing.T, certs []*x509.Certificate) []certcentral.CertificateChain {
	t.Helper()
	out := make([]certcentral.CertificateChain, 0, len(certs))
	for _, cert := range certs {
		pemBytes := pem.EncodeToMemory(&pem.Block{Type: blockTypeCertificate, Bytes: cert.Raw})
		out = append(out, certcentral.CertificateChain{Pem: string(pemBytes), SubjectCommonName: cert.Subject.CommonName})
	}
	return out
}

func parsePEMCerts(t *testing.T, data []byte) []*x509.Certificate {
	t.Helper()
	certs := make([]*x509.Certificate, 0)
	for len(data) > 0 {
		block, rest := pem.Decode(data)
		if block == nil {
			break
		}
		data = rest
		if block.Type != blockTypeCertificate {
			continue
		}
		crt, err := x509.ParseCertificate(block.Bytes)
		if err != nil {
			t.Fatalf("parse PEM cert: %v", err)
		}
		certs = append(certs, crt)
	}
	return certs
}

func certCNs(certs []*x509.Certificate) []string {
	out := make([]string, 0, len(certs))
	for _, cert := range certs {
		out = append(out, cert.Subject.CommonName)
	}
	return out
}

func assertChainCNs(t *testing.T, tlsPEM, caPEM []byte, wantTLSCNs, wantCACNs []string) {
	t.Helper()

	gotTLS := certCNs(parsePEMCerts(t, tlsPEM))
	if !reflect.DeepEqual(gotTLS, wantTLSCNs) {
		t.Fatalf("unexpected intermediate CNs, got=%v expected=%v", gotTLS, wantTLSCNs)
	}

	gotCA := certCNs(parsePEMCerts(t, caPEM))
	if !reflect.DeepEqual(gotCA, wantCACNs) {
		t.Fatalf("unexpected root CNs, got=%v expected=%v", gotCA, wantCACNs)
	}
}
