package core

import (
	deviceversion "cbe-super-app-cps-action/internal/constants/dto/device_version"
	"cbe-super-app-cps-action/internal/constants/model"
	"context"
	"time"

	"go.mongodb.org/mongo-driver/v2/bson"
)

func IdProvider(ctx context.Context, id string) (bson.ObjectID, error) {
	objID, err := bson.ObjectIDFromHex(id)
	if err != nil {
		return bson.ObjectID{}, err
	}
	return objID, nil
}

func UpdateDeviceVersionBson(req deviceversion.UpdateDeviceVersionRequest, updatedBy string) (bson.M, error) {
	update := bson.M{}
	if req.LatestVersion != "" {
		update["latest_version"] = req.LatestVersion
	}
	if req.Platform != "" {
		update["platform"] = req.Platform
	}
	if &req.ForceUpdate != nil {
		update["force_update"] = req.ForceUpdate
	}
	if req.ReleaseNotes != "" {
		update["release_notes"] = req.ReleaseNotes
	}
	update["updated_by"] = updatedBy
	update["updated_at"] = time.Now()
	return update, nil
}

func UpdateDeviceVersionBsonForDb(req model.DeviceVersionControl, updatedBy string) (bson.M, error) {
	update := bson.M{}
	if req.LatestVersion != "" {
		update["latest_version"] = req.LatestVersion
	}
	if req.Platform != "" {
		update["platform"] = req.Platform
	}
	if req.ForceUpdate {
		update["force_update"] = req.ForceUpdate
	}
	if req.ReleaseNotes != "" {
		update["release_notes"] = req.ReleaseNotes
	}
	update["enabled"] = req.Enabled
	update["updated_by"] = updatedBy

	return update, nil
}
