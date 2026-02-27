package telemetry

import (
	"net/http"

	promhttp "github.com/prometheus/client_golang/prometheus/promhttp"
)

// MetricsHandler returns an http.Handler that serves Prometheus metrics.
func MetricsHandler() http.Handler {
	return promhttp.Handler()
}

// RegisterMetricsEndpoint registers /metrics route on the provided mux.
func RegisterMetricsEndpoint(mux *http.ServeMux) {
	mux.Handle("/metrics", MetricsHandler())
}
