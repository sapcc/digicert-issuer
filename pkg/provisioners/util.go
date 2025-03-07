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
