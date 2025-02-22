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

package certmanager

import (
	"github.com/prometheus/client_golang/prometheus"
	"sigs.k8s.io/controller-runtime/pkg/metrics"
)

func init() {
	metrics.Registry.MustRegister(
		metricRequestsPending, metricRequestErrors, metricIssuerNotReady,
	)
}

var (
	metricRequestsPending = prometheus.NewCounterVec(
		prometheus.CounterOpts{
			Name: "digicertissuer_request_pending_total",
			Help: "Number of retries of a pending certificate request",
		},
		[]string{
			"certificate_request",
			"certificate",
			"secret",
			"order_id",
		},
	)

	metricRequestErrors = prometheus.NewCounterVec(
		prometheus.CounterOpts{
			Name: "digicertissuer_request_errors_total",
			Help: "Number of errors while issuing a certificate",
		},
		[]string{
			"certificate_request",
			"certificate",
			"secret",
			"reason",
		},
	)

	metricIssuerNotReady = prometheus.NewCounterVec(
		prometheus.CounterOpts{
			Name: "digicertissuer_not_ready_total",
			Help: "Increases when digicert-issuer is not ready",
		},
		[]string{
			"issuer",
			"reason",
		},
	)
)
