package middleware

import (
	"net/http"
	"strings"
)

func CheckDeviceInfoMiddleware(sourceApp, otpFor string) func(http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			deviceUUID := r.Header.Get("deviceuuid")
			appVersion := r.Header.Get("appversion")
			platform := r.Header.Get("platform")

			otpForH := r.Header.Get("otpfor")
			sourceAppH := r.Header.Get("sourceapp")

			if sourceApp != sourceAppH {
				http.Error(w, `{"message":"invalid souce app"}`, http.StatusBadRequest)
				return
			}

			switch sourceApp {
			case "memberapp", "agentapp":
				if otpForH != "pinset" && otpForH != "pinreset" {
					http.Error(w, `{"message":"invalid otpfor"}`, http.StatusBadRequest)
					return
				}
			case "dashportal", "agentportal":
				if otpForH != "login" && otpForH != "pinreset" {
					http.Error(w, `{"message":"invalid otpfor"}`, http.StatusBadRequest)
					return
				}
			case "ldapportal":
				if otpForH != "accountlink" && otpForH != "accountunlink" {
					http.Error(w, `{"message":"invalid otpfor"}`, http.StatusBadRequest)
					return
				}
			}

			if otpFor != otpForH {
				http.Error(w, `{"message":"invalid otp for"}`, http.StatusBadRequest)
				return
			}

			if strings.TrimSpace(deviceUUID) == "" || strings.TrimSpace(appVersion) == "" || strings.TrimSpace(platform) == "" {
				http.Error(w, `{"message":"Invalid request: Unknown source"}`, http.StatusUnauthorized)
				return
			}

			next.ServeHTTP(w, r)
		})
	}
}
