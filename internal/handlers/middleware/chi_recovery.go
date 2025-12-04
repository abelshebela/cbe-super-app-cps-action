package middleware

import (
	"net"
	"net/http"
	"net/http/httputil"
	"os"
	"runtime/debug"
	"strings"

	"cbe-super-app-cps-action/internal/constants/localization"
	"cbe-super-app-cps-action/platform/logger"

	"github.com/go-chi/cors"
	"go.uber.org/zap"
)

func ChiCORS() func(http.Handler) http.Handler {
	return cors.Handler(cors.Options{
		AllowedOrigins: []string{"*"},
		AllowedMethods: []string{"GET", "POST", "PUT", "DELETE", "OPTIONS", "PATCH"},
		AllowedHeaders: []string{
			"Accept", "Authorization", "Content-Type", "X-CSRF-Token",
			"access-control-allow-origin", "x-api-applicationid",
			// Custom headers used by the API
		},
		ExposedHeaders:   []string{"Link"},
		AllowCredentials: true,
		MaxAge:           300,
	})
}

func ChiCustomRecovery(logger logger.Logger, stack bool) func(next http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			defer func() {
				if err := recover(); err != nil {
					logger.Error(r.Context(), "panic recovered", zap.Any("error", err))
					var brokenPipe bool
					if ne, ok := err.(*net.OpError); ok {
						if se, ok := ne.Err.(*os.SyscallError); ok {
							if strings.Contains(strings.ToLower(se.Error()), "broken pipe") ||
								strings.Contains(strings.ToLower(se.Error()), "connection reset by peer") {
								brokenPipe = true
							}
						}
					}

					httpRequest, _ := httputil.DumpRequest(r, false)
					if brokenPipe {
						logger.Error(r.Context(), "broken pipe",
							zap.Any("error", err),
							zap.String("request", string(httpRequest)))
						return
					}

					if stack {
						logger.Error(r.Context(), "[Recovery from panic]",
							zap.Any("error", err),
							zap.String("request", string(httpRequest)),
							zap.String("stack", string(debug.Stack())),
						)
					} else {
						logger.Error(r.Context(), "[Recovery from panic]",
							zap.Any("error", err),
							zap.String("request", string(httpRequest)),
						)
					}

					localization.SendInternalServerErrorResponse(w, localization.ErrorUnexpectedError.Message)
					return
				}
			}()
			next.ServeHTTP(w, r)
		})
	}
}
