package core

import (
	imodel "github.com/abelshebela/cbe-super-app-cps-action/internal/constants/model"

	ab_dto "github.com/abelshebela/cbe-super-app-cps-action/internal/constants/dto/account_block"
)

func ToAccountBlockResponse(ab *imodel.AccountBlock) *ab_dto.AccountBlockResponse {
	if ab == nil {
		return nil
	}

	regionID := ""
	if ab.RegionID != nil {
		regionID = *ab.RegionID
	} else if ab.FederalRegionName != "" {
		regionID = ab.FederalRegionName
	}

	districtID := ""
	if ab.DistrictID != nil {
		districtID = *ab.DistrictID
	} else if ab.DistrictName != "" {
		districtID = ab.DistrictName
	}

	var parent *imodel.AccountBlock
	if ab.Parent != nil {
		parent = ab.Parent
	}

	// dr := ab.DisableReason
	// if dr == nil {
	// 	dr = []imodel.AccountBlockReason{}
	// }

	return &ab_dto.AccountBlockResponse{
		ID:                ab.ID,
		Name:              ab.Name,
		Code:              ab.Code,
		Type:              string(ab.Type),
		RegionName:        ab.RegionName,
		FederalRegionName: ab.FederalRegionName,
		DistrictName:      ab.DistrictName,
		DaoCode:           ab.DaoCode,
		AccountType:       ab.AccountType,
		Parent:            parent,
		RegionID:          regionID,
		DistrictID:        districtID,
		IsEnabled:         ab.IsEnabled,
		CreatedAt:         ab.CreatedAt,
		UpdatedAt:         ab.UpdatedAt,
	}
}
