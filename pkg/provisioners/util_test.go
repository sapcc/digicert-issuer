// SPDX-FileCopyrightText: 2022 SAP SE or an SAP affiliate company
// SPDX-License-Identifier: Apache-2.0

package provisioners

import (
	"testing"

	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
	certcentral "github.com/sapcc/go-certcentral"
)

func TestProvisioners(t *testing.T) {
	RegisterFailHandler(Fail)
	RunSpecs(t, "Provisioners Suite")
}

var _ = Describe("Map Payment Method", func() {
	It("should map existing payment method string to correct payment method", func() {
		paymentMethod, found := mapToPaymentMethod("wire_transfer")
		Expect(paymentMethod).To(Equal(certcentral.PaymentMethods.WireTransfer))
		Expect(found).To(BeTrue())
	})

	It("should return false for non-existing payment method string", func() {
		paymentMethod, found := mapToPaymentMethod("non_existing_method")
		Expect(paymentMethod).To(BeEmpty())
		Expect(found).To(BeFalse())
	})
})

var _ = Describe("Map Order Type", func() {
	It("should map existing order type string to correct order type", func() {
		orderType, found := mapToOrderType("ssl_ev_securesite_pro")
		Expect(orderType).To(Equal(certcentral.OrderTypes.SecureSiteProEVSSL))
		Expect(found).To(BeTrue())
	})

	It("should return false for non-existing order type string", func() {
		orderType, found := mapToOrderType("non_existing_type")
		Expect(orderType).To(BeEmpty())
		Expect(found).To(BeFalse())
	})
})
