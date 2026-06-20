package deviceversioncontrol

import (
	"time"

	imodel "cbe-super-app-cps-action/internal/constants/model"

	"go.mongodb.org/mongo-driver/v2/bson"
)

func deviceVersionControlMapper(deviceVersionControl imodel.DeviceVersionControl) bson.M {
	return bson.M{
		"platform":            deviceVersionControl.Platform,
		"latest_version":      deviceVersionControl.LatestVersion,
		"enabled":             deviceVersionControl.Enabled,
		"is_maintenance_mode": deviceVersionControl.IsMaintenanceMode,
		"device_state":        deviceVersionControl.DeviceState,
		"updated_at":          time.Now(),
	}
}
