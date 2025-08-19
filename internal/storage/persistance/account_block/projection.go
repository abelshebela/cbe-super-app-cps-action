package account_block

import (
	"cbe-super-app-cps-action/internal/constants/model"
	"time"

	"go.mongodb.org/mongo-driver/v2/bson"
)

func BranchMapper(branch model.Branch) bson.M {
	return bson.M{
		"branch_code":    branch.BranchCode,
		"branch_name":    branch.BranchName,
		"branch_address": branch.BranchAddress,
		"district_code":  branch.DistrictCode,
		"district_name":  branch.DistrictName,
		"branch_region":  branch.BranchRegion,
		"record_stat":    branch.RecordStat,
		"enabled":        branch.Enabled,
		"updated_at":     time.Now(),
	}
}

func CityMapper(city model.City) bson.M {
	return bson.M{
		"city_code":     city.CityCode,
		"city_name":     city.CityName,
		"city_address":  city.City,
		"district_id":   city.DistrictID,
		"district_name": city.DistrictName,
		"region_id":     city.RegionID,
		"region_name":   city.RegionName,
		"enabled":       city.Enabled,
		"updated_at":    time.Now(),
	}
}
