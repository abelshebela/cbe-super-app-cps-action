package deviceversioncontrol

import (
	"time"

	"gitlab.com/bersufekadgetachew/cbe-super-app-shared/shared/entities/model"

	"go.mongodb.org/mongo-driver/v2/bson"
)

func deviceVersionControlMapper(deviceVersionControl model.DeviceVersionControl) bson.M {
	return bson.M{
		"platform":       deviceVersionControl.Platform,
		"latest_version": deviceVersionControl.LatestVersion,
		"enabled":        deviceVersionControl.Enabled,
		"updated_at":     time.Now(),
	}
}
