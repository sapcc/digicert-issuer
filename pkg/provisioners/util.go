// SPDX-FileCopyrightText: 2022 SAP SE or an SAP affiliate company
// SPDX-License-Identifier: Apache-2.0

package provisioners

import (
	"reflect"
	"strings"

	certcentral "github.com/sapcc/go-certcentral"
)

func mapToOrderType(s string) (certcentral.OrderType, bool) {
	orderTypes := listAvailableOrderTypes()
	for _, t := range orderTypes {
		if strings.ToLower(s) == strings.ToLower(t.String()) {
			return t, true
		}
	}
	return "", false
}

func listAvailableOrderTypes() []certcentral.OrderType {
	var orderTypes []certcentral.OrderType
	v := reflect.ValueOf(certcentral.OrderTypes)
	for i := 0; i < v.NumField(); i++ {
		orderTypes = append(orderTypes, v.Field(i).Interface().(certcentral.OrderType))
	}
	return orderTypes
}

func mapToPaymentMethod(s string) (certcentral.PaymentMethod, bool) {
	paymentMethods := listAvailablePaymentMethods()
	for _, m := range paymentMethods {
		if strings.ToLower(s) == strings.ToLower(m.String()) {
			return m, true
		}
	}
	return "", false
}

func listAvailablePaymentMethods() []certcentral.PaymentMethod {
	var paymentMethods []certcentral.PaymentMethod
	v := reflect.ValueOf(certcentral.PaymentMethods)
	for i := 0; i < v.NumField(); i++ {
		paymentMethods = append(paymentMethods, v.Field(i).Interface().(certcentral.PaymentMethod))
	}
	return paymentMethods
}
