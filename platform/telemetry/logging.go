package telemetry

import (
	"context"
	"net/http"

	"cbe-super-app-cps-action/internal/constants"

	"go.opentelemetry.io/otel/trace"
)

// TraceContextMiddleware extracts trace/span ids from the current span in context and
// injects them into the request context using keys that the app logger extracts.
func TraceContextMiddleware() func(next http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			ctx := r.Context()
			sp := trace.SpanFromContext(ctx)
			if sp != nil {
				sc := sp.SpanContext()
				if sc.IsValid() {
					ctx = context.WithValue(ctx, constants.ContextKey("trace_id"), sc.TraceID().String())
					ctx = context.WithValue(ctx, constants.ContextKey("span_id"), sc.SpanID().String())
				}
			}
			next.ServeHTTP(w, r.WithContext(ctx))
		})
	}
}
