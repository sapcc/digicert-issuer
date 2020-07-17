package provisioners

import (
	"context"
	"fmt"

	certmanagerv1alpha2 "github.com/jetstack/cert-manager/pkg/apis/certmanager/v1alpha2"
	"github.com/sapcc/digicert-issuer/apis/certmanager/v1beta1"
	certcentral "github.com/sapcc/go-certcentral"
)

const defaultValidityYears = 1

type CertCentral struct {
	name   string
	client *certcentral.Client

	validityYears,
	organizationID int
	caCertID            string
	organizationalUnits []string
	skipApproval,
	disableRenewalNotifications bool
	orderType     certcentral.OrderType
	paymentMethod certcentral.PaymentMethod
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

	validityYears := defaultValidityYears
	if issuer.Spec.Provisioner.ValidityYears != nil {
		validityYears = *issuer.Spec.Provisioner.ValidityYears
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

	orderType := certcentral.OrderTypes.PrivateSSLPlus
	if t, ok := mapToOrderType(issuer.Spec.Provisioner.OrderType); ok {
		orderType = t
	}

	paymentMethod := certcentral.PaymentMethods.Balance
	if m, ok := mapToPaymentMethod(issuer.Spec.Provisioner.PaymentMethod); ok {
		paymentMethod = m
	}

	return &CertCentral{
		name:                        fmt.Sprintf("%s/%s", issuer.GetName(), issuer.GetNamespace()),
		client:                      client,
		validityYears:               validityYears,
		organizationID:              organizationID,
		caCertID:                    issuer.Spec.Provisioner.CACertID,
		organizationalUnits:         orgUnits,
		skipApproval:                skipApproval,
		disableRenewalNotifications: disableRenewalNotifications,
		orderType:                   orderType,
		paymentMethod:               paymentMethod,
	}, nil
}

func (c *CertCentral) Sign(ctx context.Context, cr *certmanagerv1alpha2.CertificateRequest) ([]byte, error) {
	certReq, err := decodeCertificateRequest(cr.Spec.CSRPEM)
	if err != nil {
		return nil, err
	}

	if err := certReq.CheckSignature(); err != nil {
		return nil, err
	}

	sans := certReq.DNSNames
	for _, ipAddr := range certReq.IPAddresses {
		sans = append(sans, ipAddr.String())
	}

	orderResponse, err := c.client.SubmitOrder(certcentral.Order{
		Certificate: certcentral.Certificate{
			CommonName:        getCommonName(certReq),
			DNSNames:          sans,
			CSR:               string(cr.Spec.CSRPEM),
			ServerPlatform:    certcentral.ServerPlatformForType(certcentral.ServerPlatformTypes.Nginx),
			SignatureHash:     certcentral.SignatureHashes.SHA256,
			CaCertID:          c.caCertID,
			OrganizationUnits: c.organizationalUnits,
		},
		ValidityYears:               c.validityYears,
		DisableRenewalNotifications: c.disableRenewalNotifications,
		PaymentMethod:               c.paymentMethod,
		SkipApproval:                c.skipApproval,
		Organization:                &certcentral.Organization{ID: c.organizationID},
	}, c.orderType)
	if err != nil {
		return nil, err
	}

	crtChainPEMs := make([]byte, 0)
	for _, crt := range orderResponse.CertificateChain {
		crtChainPEMs = append(crtChainPEMs, crt.Pem...)
	}

	return crtChainPEMs, nil
}
