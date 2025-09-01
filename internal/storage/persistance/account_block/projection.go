package account_block

import (
	"cbe-super-app-cps-action/internal/constants/model"
	"time"

	"go.mongodb.org/mongo-driver/v2/bson"
)

func BranchMapperForUpdate(branch model.Branch) bson.M {
	update := bson.M{}
	now := time.Now()
	update["updated_at"] = now

	if branch.BranchCode != "" {
		update["branch_code"] = branch.BranchCode
	}
	if branch.BranchName != "" {
		update["branch_name"] = branch.BranchName
	}
	if branch.BranchAddress != "" {
		update["branch_address"] = branch.BranchAddress
	}
	if branch.DistrictCode != "" {
		update["district_code"] = branch.DistrictCode
	}
	if branch.DistrictName != "" {
		update["district_name"] = branch.DistrictName
	}
	if branch.BranchRegion != "" {
		update["branch_region"] = branch.BranchRegion
	}
	if branch.RecordStat != "" {
		update["record_stat"] = branch.RecordStat
	}
	// Booleans are tricky: include them only if they are explicitly meant to be updated
	update["enabled"] = branch.Enabled

	return update
}

func RegionMapperForUpdate(region model.Region) bson.M {
	update := bson.M{}
	now := time.Now()
	update["updated_at"] = now

	if region.RegionCode != "" {
		update["region_code"] = region.RegionCode
	}
	if region.RegionName != "" {
		update["region_name"] = region.RegionName
	}
	if region.RegionAddress != "" {
		update["region_address"] = region.RegionAddress
	}
	update["enabled"] = region.Enabled

	return update
}

func DistrictMapperForUpdate(district model.District) bson.M {
	update := bson.M{}
	now := time.Now()
	update["updated_at"] = now

	if district.DistrictCode != "" {
		update["district_code"] = district.DistrictCode
	}
	if district.DistrictName != "" {
		update["district_name"] = district.DistrictName
	}
	if district.DistrictAddress != "" {
		update["district_address"] = district.DistrictAddress
	}
	if district.RegionID != "" {
		update["region_id"] = district.RegionID
	}
	if district.RegionName != "" {
		update["region_name"] = district.RegionName
	}
	update["enabled"] = district.Enabled

	return update
}

func CityMapperForUpdate(city model.City) bson.M {
	update := bson.M{}
	now := time.Now()
	update["updated_at"] = now

	if city.CityCode != "" {
		update["city_code"] = city.CityCode
	}
	if city.CityName != "" {
		update["city_name"] = city.CityName
	}
	if city.CityAddress != "" {
		update["city_address"] = city.CityAddress
	}
	if city.DistrictID != "" {
		update["district_id"] = city.DistrictID
	}
	if city.DistrictName != "" {
		update["district_name"] = city.DistrictName
	}
	if city.RegionID != "" {
		update["region_id"] = city.RegionID
	}
	if city.RegionName != "" {
		update["region_name"] = city.RegionName
	}
	update["enabled"] = city.Enabled

	return update
}
