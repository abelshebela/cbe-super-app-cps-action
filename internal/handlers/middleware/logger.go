package middleware

import (
	"context"
	"net/http"
	"time"

	"github.com/CBE-Super-App/cbe-super-app-member-auth/internal/constants"
	"gitlab.com/bersufekadgetachew/cbe-super-app-shared/shared/utils"

	"github.com/google/uuid"
	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
)

// Wraps http.ResponseWriter to capture status code
type responseWriter struct {
	http.ResponseWriter
	statusCode int
}

func (rw *responseWriter) WriteHeader(code int) {
	rw.statusCode = code
	rw.ResponseWriter.WriteHeader(code)
}

func ChiLogger(log utils.Logger) func(next http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			start := time.Now()
			path := r.URL.Path
			query := r.URL.RawQuery
			if query != "" {
				path = path + "?" + query
			}
			id := uuid.New().String()
			ctx := context.WithValue(r.Context(), constants.ContextKey("x-request-id"), id)
			ctx = context.WithValue(ctx, constants.ContextKey("request-start-time"), start)

			ww := &responseWriter{ResponseWriter: w}

			next.ServeHTTP(ww, r.WithContext(ctx))

			end := time.Now()
			latency := end.Sub(start)

			fields := []zapcore.Field{
				zap.Int("status", ww.statusCode),
				zap.String("method", r.Method),
				zap.String("path", path),
				zap.String("query", query),
				zap.Int64("request-latency(ms)", latency.Milliseconds()),
				zap.String("user-agent", r.UserAgent()),
				zap.String("ip", r.RemoteAddr),
			}
			log.Infof("request completed", fields)
		})
	}
}
