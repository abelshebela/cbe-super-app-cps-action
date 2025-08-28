package core

import (
	"cbe-super-app-cps-action/internal/constants/model"

	ab_dto "cbe-super-app-cps-action/internal/constants/dto/account_block"
)

func ToBranchResponse(branch *model.Branch) *ab_dto.BranchResponse {
	return &ab_dto.BranchResponse{
		ID:            branch.ID.Hex(),
		BranchCode:    branch.BranchCode,
		BranchName:    branch.BranchName,
		BranchAddress: branch.BranchAddress,
		DistrictCode:  branch.DistrictCode,
		DistrictName:  branch.DistrictName,
		BranchRegion:  branch.BranchRegion,
		RecordStat:    branch.RecordStat,
		CreatedAt:     branch.CreatedAt,
		IsDeleted:     branch.IsDeleted,
		UpdatedAt:     branch.UpdatedAt,
		Enabled:       branch.Enabled,
	}

}

func ToBranchesResponse(branches []*model.Branch) []*ab_dto.BranchResponse {
	var branchesResponse []*ab_dto.BranchResponse
	for _, branch := range branches {
		branchesResponse = append(branchesResponse, ToBranchResponse(branch))
	}
	return branchesResponse
}

func ToRegionResponse(region *model.Region) *ab_dto.RegionResponse {
	return &ab_dto.RegionResponse{
		ID:            region.ID.Hex(),
		RegionName:    region.RegionName,
		RegionCode:    region.RegionCode,
		RegionAddress: region.RegionAddress,
		Enabled:       region.Enabled,
		CreatedAt:     region.CreatedAt,
		UpdatedAt:     region.UpdatedAt,
	}
}

func ToRegionsResponse(regions []*model.Region) []*ab_dto.RegionResponse {
	var regionsResponse []*ab_dto.RegionResponse
	for _, region := range regions {
		regionsResponse = append(regionsResponse, ToRegionResponse(region))
	}

	return regionsResponse
}

func ToDistrictResponse(district *model.District) *ab_dto.DistrictResponse {
	return &ab_dto.DistrictResponse{
		ID:              district.ID.Hex(),
		DistrictCode:    district.DistrictCode,
		DistrictName:    district.DistrictName,
		DistrictAddress: district.DistrictAddress,
		RegionID:        district.RegionID,
		RegionName:      district.RegionName,
		CreatedAt:       district.CreatedAt,
		UpdatedAt:       district.UpdatedAt,
		Enabled:         district.Enabled,
	}
}

func ToDistrictsResponse(districts []*model.District) []*ab_dto.DistrictResponse {
	var districtsResponse []*ab_dto.DistrictResponse
	for _, district := range districts {
		districtsResponse = append(districtsResponse, ToDistrictResponse(district))
	}

	return districtsResponse
}

func ToCityResponse(city *model.City) *ab_dto.CityResponse {
	return &ab_dto.CityResponse{
		ID:           city.ID.Hex(),
		CityCode:     city.CityCode,
		CityName:     city.CityName,
		CityAddress:  city.CityAddress,
		DistrictID:   city.DistrictID,
		DistrictName: city.DistrictName,
		RegionID:     city.RegionID,
		RegionName:   city.RegionName,
		CreatedAt:    city.CreatedAt,
		UpdatedAt:    city.UpdatedAt,
		Enabled:      city.Enabled,
	}
}

func ToCitiesResponse(cities []*model.City) []*ab_dto.CityResponse {
	var citiesResponse []*ab_dto.CityResponse
	for _, city := range cities {
		citiesResponse = append(citiesResponse, ToCityResponse(city))
	}

	return citiesResponse
}
