package core

import (
	"net/http"

	"github.com/CBE-Super-App/cbe-super-app-member-auth/internal/constants"
	"github.com/CBE-Super-App/cbe-super-app-member-auth/internal/constants/dto"
)

func ExtractHeader(r *http.Request) dto.DeviceLookupRequest {
	// installationDate := utils.ParseTime(r.Header.Get("application_installation_date"))
	installationDate := r.Header.Get("installation_date")
	return dto.DeviceLookupRequest{
		Platform:                    constants.Platform(r.Header.Get("platform")),
		AppVersion:                  r.Header.Get("app_version"),
		DeviceUUID:                  r.Header.Get("device_uuid"),
		SourceApp:                   r.Header.Get("source_app"),
		ApplicationInstallationDate: installationDate,
	}
}
