package core

import (
	deviceversion "cbe-super-app-cps-action/internal/constants/dto/device_version"
	imodel "cbe-super-app-cps-action/internal/constants/model"
	"context"
	"strings"
	"time"

	"go.mongodb.org/mongo-driver/v2/bson"
)

// DeviceStateToFlags maps the device_state enum to the two boolean flags stored in MongoDB.
func DeviceStateToFlags(deviceState string) (forceUpdate bool, isMaintenanceMode bool) {
	switch deviceState {
	case imodel.DeviceStateForceUpdate:
		return true, false
	case imodel.DeviceStateMaintenance:
		return false, true
	default: // STABLE or empty
		return false, false
	}
}

func IdProvider(ctx context.Context, id string) (bson.ObjectID, error) {
	objID, err := bson.ObjectIDFromHex(id)
	if err != nil {
		return bson.ObjectID{}, err
	}
	return objID, nil
}

func UpdateDeviceVersionBson(req deviceversion.UpdateDeviceVersionRequest, updatedBy string, existing imodel.DeviceVersionControl) (bson.M, error) {
	update := bson.M{}
	if req.LatestVersion != "" {
		update["latest_version"] = req.LatestVersion
	}
	if req.Platform != "" {
		update["platform"] = strings.ToUpper(req.Platform)
	}

	if req.DeviceState != "" {
		forceUpdate, isMaintenance := DeviceStateToFlags(req.DeviceState)
		update["device_state"] = req.DeviceState
		update["force_update"] = forceUpdate
		update["is_maintenance_mode"] = isMaintenance
	} else {
		update["device_state"] = existing.DeviceState
		update["force_update"] = existing.ForceUpdate
		update["is_maintenance_mode"] = existing.IsMaintenanceMode
	}

	if req.ReleaseNotes != "" {
		update["release_notes"] = req.ReleaseNotes
	}
	if req.Enabled != nil {
		update["enabled"] = *req.Enabled
	} else {
		update["enabled"] = existing.Enabled
	}

	update["updated_by"] = updatedBy
	update["updated_at"] = time.Now()
	return update, nil
}

func UpdateDeviceVersionBsonForDb(req imodel.DeviceVersionControl, updatedBy string) (bson.M, error) {
	update := bson.M{}
	if req.LatestVersion != "" {
		update["latest_version"] = req.LatestVersion
	}
	if req.Platform != "" {
		update["platform"] = req.Platform
	}
	if req.ReleaseNotes != "" {
		update["release_notes"] = req.ReleaseNotes
	}
	update["device_state"] = req.DeviceState
	update["force_update"] = req.ForceUpdate
	update["is_maintenance_mode"] = req.IsMaintenanceMode
	update["enabled"] = req.Enabled
	update["updated_by"] = updatedBy

	return update, nil
}
