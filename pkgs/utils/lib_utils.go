package utils

import (
	"fmt"
	"net/http"
	"time"

	constants "github.com/CBE-Super-App/cbe-super-app-member-auth/internal/constants"
	"github.com/CBE-Super-App/cbe-super-app-member-auth/internal/constants/dto"
)

func parseTime(date string) time.Time {
	parsedTime, err := time.Parse(time.RFC3339, date)
	if err != nil {
		fmt.Println("Error parsing date:", err)
		return time.Time{}
	}
	return parsedTime
}

func ExtractHeader(r *http.Request) dto.DeviceLookupRequest {
	installationDate := parseTime(r.Header.Get("application_installation_date"))
	return dto.DeviceLookupRequest{
		Platform:                    constants.Platform(r.Header.Get("platform")),
		AppVersion:                  r.Header.Get("app_version"),
		DeviceUUID:                  r.Header.Get("device_uuid"),
		SourceApp:                   r.Header.Get("source_app"),
		ApplicationInstallationDate: installationDate,
	}
}
