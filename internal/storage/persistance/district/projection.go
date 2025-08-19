package district

import (
	"cbe-super-app-cps-action/internal/constants/model"
	"time"

	"go.mongodb.org/mongo-driver/v2/bson"
)

// DistrictMapper maps a District model to a bson.M for updates
func DistrictMapper(district model.District) bson.M {
	return bson.M{
		"$set": bson.M{
			"district_code":    district.DistrictCode,
			"district_name":    district.DistrictName,
			"district_address": district.DistrictAddress,
			"region_id":        district.RegionID,
			"region_name":      district.RegionName,
			"enabled":          district.Enabled,
			"updated_at":       time.Now(),
		},
	}
} 