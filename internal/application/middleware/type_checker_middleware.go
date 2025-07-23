package middleware

import (
	common_utils "github.com/CBE-Super-App/cbe-super-app-cps-action/pkgs/utils"
	"net/http"
	"strings"
)

// RequireJSONContentType ensures Content-Type is application/json
func RequireJSONContentType() func(http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			if !strings.HasPrefix(r.Header.Get("Content-Type"), "application/json") {
				common_utils.SendErrorResponse(w, "CONTENT_TYPE_MUST_BE_JSON", 0, nil)
				return
			}
			next.ServeHTTP(w, r)
		})
	}
}

// RequireFormContentType ensures Content-Type is multipart/form-data
func RequireFormContentType() func(http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			if !strings.HasPrefix(r.Header.Get("Content-Type"), "multipart/form-data") {
				common_utils.SendErrorResponse(w, "CONTENT_TYPE_MUST_BE_FORM", 0, nil)
				return
			}
			next.ServeHTTP(w, r)
		})
	}
}
