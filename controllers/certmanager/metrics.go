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
