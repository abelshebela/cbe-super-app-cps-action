package validator

import (
	"encoding/json"
	"fmt"
	"net/http"
	"time"

	"github.com/CBE-Super-App/cbe-super-app-cps-action/pkgs/common"

	"github.com/google/uuid"
)

func ValidateHeaders(next http.Handler) http.Handler {

	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		platform := r.Header.Get("platform")
		appVersion := r.Header.Get("app_version")
		deviceUUID := r.Header.Get("device_uuid")
		sourceApp := r.Header.Get("source_app")
		installationDate := r.Header.Get("installation_date")
		returndata := make(map[string]interface{})

		var missing []string
		if platform == "" {
			missing = append(missing, "platform")
		}
		if appVersion == "" {
			missing = append(missing, "app_version")
		}
		if deviceUUID == "" {
			missing = append(missing, "device_uuid")
		}
		if sourceApp == "" {
			missing = append(missing, "source_app")
		}
		if installationDate == "" {
			missing = append(missing, "installation_date")
		}

		if len(missing) > 0 {
			returndata["code"] = common.DefineError.General["MISSING_HEADER"].Code
			returndata["message"] = common.DefineError.General["MISSING_HEADER"].Message
			returndata["errors"] = fmt.Sprintf("Missing required headers: %v", missing)
			w.Header().Set("Content-Type", "applicaton/json")
			w.WriteHeader(http.StatusBadRequest)
			json.NewEncoder(w).Encode(returndata)
			return
		}

		// // Validate device_uuid format
		if _, err := uuid.Parse(deviceUUID); err != nil {
			returndata["message"] = "Invalid Device UUID"
			w.Header().Set("Content-Type", "applicaton/json")
			w.WriteHeader(http.StatusBadRequest)
			json.NewEncoder(w).Encode(returndata)
			return
		}

		_, err := time.Parse(time.RFC3339, installationDate)
		if err != nil {
			returndata["message"] = "Invalid installation_date format. Expected ISO8601 (e.g. 2025-06-23T00:00:00Z)"
			w.Header().Set("Content-Type", "applicaton/json")
			w.WriteHeader(http.StatusBadRequest)
			json.NewEncoder(w).Encode(returndata)
			return
		}

		next.ServeHTTP(w, r)
	})
}
