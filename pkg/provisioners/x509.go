// SPDX-FileCopyrightText: 2022 SAP SE or an SAP affiliate company
// SPDX-License-Identifier: Apache-2.0

package provisioners

import (
	"bytes"
	"crypto/x509"
	"encoding/pem"
	"errors"
	"fmt"
)

const (
	blockTypeCertificate        = "CERTIFICATE"
	blockTypeCertificateRequest = "CERTIFICATE REQUEST"
)

func decodeCertificateRequest(data []byte) (*x509.CertificateRequest, error) {
	block, _ := pem.Decode(data)
	if block == nil {
		return nil, errors.New("failed to decode CR pem")
	}

	if block.Type != blockTypeCertificateRequest {
		return nil, fmt.Errorf("pem is not of type %s, but: %s", blockTypeCertificateRequest, block.Type)
	}

	cr, err := x509.ParseCertificateRequest(block.Bytes)
	if err != nil {
		return nil, err
	}

	err = cr.CheckSignature()
	return cr, err
}

func encodeCertificate(cert *x509.Certificate) ([]byte, error) {
	block := &pem.Block{
		Type:  blockTypeCertificate,
		Bytes: cert.Raw,
	}

	buf := bytes.NewBuffer([]byte{})
	err := pem.Encode(buf, block)
	return buf.Bytes(), err
}

func getCommonName(cr *x509.CertificateRequest) string {
	if cr.Subject.CommonName != "" {
		return cr.Subject.CommonName
	}

	if len(cr.DNSNames) > 0 {
		return cr.DNSNames[0]
	}

	return "digicert-issuer-certificate"
}

func isSelfSigned(cert *x509.Certificate) bool {
	return cert.CheckSignatureFrom(cert) == nil
}
