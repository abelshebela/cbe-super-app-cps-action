package middleware

import (
	"net/http"
	"runtime/debug"

	common_utils "github.com/CBE-Super-App/cbe-super-app-cps-action/pkgs/utils"

	"gitlab.com/bersufekadgetachew/cbe-super-app-shared/shared/utils"
)

func RecoveryMiddleware(logger utils.Logger) func(http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			defer func() {
				if rec := recover(); rec != nil {
					logger.Errorf("panic recovered: %v\n%s", rec, debug.Stack())
					common_utils.SendErrorResponse(w, "UNHANDLED_SERVER_ERROR", 0, nil)
				}
			}()

			next.ServeHTTP(w, r)
		})
	}
}
