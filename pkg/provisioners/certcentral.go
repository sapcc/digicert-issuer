// SPDX-FileCopyrightText: 2022 SAP SE or an SAP affiliate company
// SPDX-License-Identifier: Apache-2.0

package provisioners

import (
	"context"
	"crypto/x509"
	"fmt"

	certmanagerv1 "github.com/cert-manager/cert-manager/pkg/apis/certmanager/v1"
	"github.com/sapcc/digicert-issuer/apis/certmanager/v1beta1"
	certcentral "github.com/sapcc/go-certcentral"
)

const defaultValidityYears = 1

type CertCentral struct {
	name   string
	client *certcentral.Client

	validityDays        *int
	validityYears       *int
	organizationID      int
	caCertID            string
	organizationalUnits []string
	skipApproval,
	disableRenewalNotifications bool
	orderType     certcentral.OrderType
	paymentMethod certcentral.PaymentMethod
	containerID   int
}

func (c CertCentral) GetName() string {
	return c.name
}

func New(issuer *v1beta1.DigicertIssuer, apiToken string) (*CertCentral, error) {
	client, err := certcentral.New(&certcentral.Options{
		Token: apiToken,
	})
	if err != nil {
		return nil, err
	}

	var organizationID int
	if issuer.Spec.Provisioner.OrganizationID != nil {
		organizationID = *issuer.Spec.Provisioner.OrganizationID
	}

	if issuer.Spec.Provisioner.OrganizationName != "" {
		org, err := client.GetOrganizationByName(issuer.Spec.Provisioner.OrganizationName)
		if err != nil {
			return nil, err
		}
		organizationID = org.ID
	}

	validityYears := issuer.Spec.Provisioner.ValidityYears
	validityDays := issuer.Spec.Provisioner.ValidityDays
	if validityYears != nil && validityDays != nil {
		return nil, fmt.Errorf("can not handle both validityYears and validityDays")
	}

	if validityYears == nil && validityDays == nil {
		v := defaultValidityYears
		validityYears = &v
	}

	orgUnits := make([]string, 0)
	if issuer.Spec.Provisioner.OrganizationUnits != nil {
		orgUnits = issuer.Spec.Provisioner.OrganizationUnits
	}

	skipApproval := true
	if issuer.Spec.Provisioner.SkipApproval != nil {
		skipApproval = *issuer.Spec.Provisioner.SkipApproval
	}

	disableRenewalNotifications := true
	if issuer.Spec.Provisioner.DisableRenewalNotifications != nil {
		disableRenewalNotifications = *issuer.Spec.Provisioner.DisableRenewalNotifications
	}

	orderType := certcentral.OrderTypes.SecureSiteOV
	if t, ok := mapToOrderType(issuer.Spec.Provisioner.OrderType); ok {
		orderType = t
	}

	paymentMethod := certcentral.PaymentMethods.Balance
	if m, ok := mapToPaymentMethod(issuer.Spec.Provisioner.PaymentMethod); ok {
		paymentMethod = m
	}

	var containerID int
	if issuer.Spec.Provisioner.ContainerID != nil {
		containerID = *issuer.Spec.Provisioner.ContainerID
	}

	return &CertCentral{
		name:                        fmt.Sprintf("%s/%s", issuer.GetName(), issuer.GetNamespace()),
		client:                      client,
		validityYears:               validityYears,
		validityDays:                validityDays,
		organizationID:              organizationID,
		caCertID:                    issuer.Spec.Provisioner.CACertID,
		organizationalUnits:         orgUnits,
		skipApproval:                skipApproval,
		disableRenewalNotifications: disableRenewalNotifications,
		orderType:                   orderType,
		paymentMethod:               paymentMethod,
		containerID:                 containerID,
	}, nil
}

func (c *CertCentral) Sign(ctx context.Context, cr *certmanagerv1.CertificateRequest) ([]byte, []byte, *certcentral.Order, error) {
	certReq, err := decodeCertificateRequest(cr.Spec.Request)
	if err != nil {
		return nil, nil, nil, err
	}

	if err := certReq.CheckSignature(); err != nil {
		return nil, nil, nil, err
	}

	sans := certReq.DNSNames
	for _, ipAddr := range certReq.IPAddresses {
		sans = append(sans, ipAddr.String())
	}

	orderValidity := certcentral.OrderValidity{}
	if c.validityDays != nil {
		orderValidity.Days = *c.validityDays
	}
	if c.validityYears != nil {
		orderValidity.Years = *c.validityYears
	}

	orderResponse, err := c.client.SubmitOrder(certcentral.Order{
		Certificate: certcentral.Certificate{
			CommonName:        getCommonName(certReq),
			DNSNames:          sans,
			CSR:               string(cr.Spec.Request),
			ServerPlatform:    certcentral.ServerPlatformForType(certcentral.ServerPlatformTypes.Nginx),
			SignatureHash:     certcentral.SignatureHashes.SHA256,
			CaCertID:          c.caCertID,
			OrganizationUnits: c.organizationalUnits,
		},
		OrderValidity:               orderValidity,
		DisableRenewalNotifications: c.disableRenewalNotifications,
		PaymentMethod:               c.paymentMethod,
		SkipApproval:                c.skipApproval,
		Organization: &certcentral.Organization{
			ID: c.organizationID,
		},
		Container: &certcentral.Container{
			ID: c.containerID,
		},
	}, c.orderType)
	if err != nil {
		return nil, nil, nil, err
	}

	// TODO: This currently relies on skipApproval=true, so the certificates are returned immediately.
	// If that's not the case, the certificateIDs needs to be extracted from the response and
	// each element of the certificate chain needs to be downloaded.
	// Bonus points: Cache the CA and intermediate as they won't change and only download the missing certificate.
	crtChain, err := orderResponse.DecodeCertificateChain()
	if err != nil {
		return nil, nil, nil, err
	}

	rootCAPEM, crtChainPEMs, err := encodePem(crtChain)
	if err != nil {
		return nil, nil, nil, err
	}

	return rootCAPEM, crtChainPEMs, orderResponse, nil
}

func (c *CertCentral) Download(ctx context.Context, cr *certmanagerv1.CertificateRequest) ([]byte, []byte, error) {
	certID := cr.GetAnnotations()["certmanager.cloud.sap/digicert-cert-id"]
	if certID == "" {
		// TODO: get cert_id by order_id if missing
		return nil, nil, fmt.Errorf("no cert id given for %s", cr.ObjectMeta.Name)
	}

	chain, err := c.client.GetCertificateChain(certID)
	if err != nil {
		return nil, nil, fmt.Errorf("error receiving certificate chain %s for request %s: %s", certID, cr.ObjectMeta.Name, err)
	}

	crtChain := make([]*x509.Certificate, 0)
	for _, crt := range chain {
		decodedCrt, err := crt.DecodePEM()
		if err != nil {
			return nil, nil, err
		}
		crtChain = append(crtChain, decodedCrt...)
	}

	return encodePem(crtChain)
}

func encodePem(crtChain []*x509.Certificate) ([]byte, []byte, error) {
	rootCAPEM := make([]byte, 0)
	crtChainPEMs := make([]byte, 0)

	for _, crt := range crtChain {
		crtPEM, err := encodeCertificate(crt)
		if err != nil {
			return nil, nil, err
		}

		switch {
		case crt.IsCA && isSelfSigned(crt):
			rootCAPEM = append(rootCAPEM, crtPEM...)
		default:
			crtChainPEMs = append(crtChainPEMs, crtPEM...)
		}
	}

	return rootCAPEM, crtChainPEMs, nil
}
