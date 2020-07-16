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
