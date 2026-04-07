package telemetry

import (
	"net/http"

	"github.com/go-chi/chi/v5/middleware"
	"go.opentelemetry.io/contrib/instrumentation/net/http/otelhttp"
)

// HTTPMiddleware returns a slice of middlewares to use with Chi router for tracing.
func HTTPMiddleware() []func(next http.Handler) http.Handler {
	return []func(next http.Handler) http.Handler{
		middleware.RequestID,
		middleware.RealIP,
		middleware.Logger,
		middleware.Recoverer,
	}
}

// WrapHandler exposes a convenience wrapper to convert an http.Handler to an otel-instrumented handler
func WrapHandler(handler http.Handler, name string) http.Handler {
	return otelhttp.NewHandler(handler, name)
}
