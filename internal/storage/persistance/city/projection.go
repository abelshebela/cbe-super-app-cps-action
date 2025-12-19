package city

import (
	"time"

	"gitlab.com/bersufekadgetachew/cbe-super-app-shared/shared/entities/model"

	"go.mongodb.org/mongo-driver/v2/bson"
)

// CityMapper maps a City model to a bson.M for updates
func CityMapper(city model.City) bson.M {
	return bson.M{
		"$set": bson.M{
			"city_code":     city.CityCode,
			"city_name":     city.CityName,
			"city_address":  city.CityAddress,
			"district_id":   city.DistrictID,
			"district_name": city.DistrictName,
			"region_id":     city.RegionID,
			"region_name":   city.RegionName,
			"enabled":       city.Enabled,
			"updated_at":    time.Now(),
		},
	}
}
