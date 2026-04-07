package telemetry

import (
	"context"
	"net/http"

	"cbe-super-app-cps-action/internal/constants"

	"gitlab.com/bersufekadgetachew/cbe-super-app-shared/shared/utils"
	"go.opentelemetry.io/otel/trace"
)

// SpanPresenceMiddleware logs whether a span is present and its trace id for each request.
// This is a lightweight diagnostic middleware to help locate routes where tracing is missing.
func SpanPresenceMiddleware(log utils.Logger) func(next http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			ctx := r.Context()
			sp := trace.SpanFromContext(ctx)
			if sp == nil || !sp.SpanContext().IsValid() {
				// mark in context for logger if not present
				ctx = context.WithValue(ctx, constants.ContextKey("trace_present"), "false")
				log.Warnf("otel span missing for request")
			} else {
				traceID := sp.SpanContext().TraceID().String()
				ctx = context.WithValue(ctx, constants.ContextKey("trace_present"), "true")
				ctx = context.WithValue(ctx, constants.ContextKey("trace_id"), traceID)
				log.Infof("otel span present for request")
			}
			next.ServeHTTP(w, r.WithContext(ctx))
		})
	}
}
