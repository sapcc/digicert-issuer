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
	orderTypes := reflect.Indirect(reflect.ValueOf(certcentral.OrderTypes))
	len := orderTypes.NumField()
	allTypes := make([]certcentral.OrderType, len)
	for i := 0; i < len; i++ {
		allTypes[i] = certcentral.OrderType(orderTypes.Type().Field(i).Name)
	}
	return allTypes
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
	paymentMethods := reflect.Indirect(reflect.ValueOf(certcentral.PaymentMethods))
	len := paymentMethods.NumField()
	allMethods := make([]certcentral.PaymentMethod, len)
	for i := 0; i < len; i++ {
		allMethods[i] = certcentral.PaymentMethod(paymentMethods.Type().Field(i).Name)
	}
	return allMethods
}
