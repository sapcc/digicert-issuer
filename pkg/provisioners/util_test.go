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
