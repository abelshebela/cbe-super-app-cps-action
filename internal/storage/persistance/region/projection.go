package region

import (
	"time"

	"gitlab.com/bersufekadgetachew/cbe-super-app-shared/shared/entities/model"

	"go.mongodb.org/mongo-driver/v2/bson"
)

// RegionMapper maps a Region model to a bson.M for updates
func RegionMapper(region model.Region) bson.M {
	return bson.M{
		"$set": bson.M{
			"region_code":    region.RegionCode,
			"region_name":    region.RegionName,
			"region_address": region.RegionAddress,
			"enabled":        region.Enabled,
			"updated_at":     time.Now(),
		},
	}
}
