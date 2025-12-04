package region

import (
	"cbe-super-app-cps-action/internal/constants/model"
	"time"

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
