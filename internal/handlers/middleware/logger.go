package middleware

import (
	"cbe-super-app-budget/internal/constants"
	"cbe-super-app-budget/platform/logger"
	"context"
	"net/http"
	"time"

	"github.com/go-chi/chi/v5/middleware"
	"github.com/google/uuid"
	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
)

func ChiLogger(log logger.Logger) func(next http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			start := time.Now()
			path := r.URL.Path
			query := r.URL.RawQuery

			ctx := context.WithValue(r.Context(), constants.ContextKey("request-start-time"), start)

			requestID := r.Context().Value(middleware.RequestIDKey)
			if requestID == "" {
				ctx = context.WithValue(ctx, middleware.RequestIDKey, uuid.New().String())
			}

			ww := middleware.NewWrapResponseWriter(w, r.ProtoMajor)

			r = r.WithContext(ctx)
			next.ServeHTTP(ww, r)

			latency := time.Since(start)

			fields := []zapcore.Field{
				zap.Int("status", ww.Status()),
				zap.String("method", r.Method),
				zap.String("path", path),
				zap.String("query", query),
				zap.String("id", r.RemoteAddr),
				zap.String("ip", r.RemoteAddr),
				zap.String("user-agent", r.UserAgent()),
				zap.Int64("request-latency", latency.Milliseconds()),
			}

			log.Info(r.Context(), "HTTP", fields...)

		})
	}
}
