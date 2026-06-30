package core

import (
	imodel "cbe-super-app-cps-action/internal/constants/model"

	ab_dto "cbe-super-app-cps-action/internal/constants/dto/account_block"
)

func ToAccountBlockResponse(ab *imodel.AccountBlock) *ab_dto.AccountBlockResponse {
	if ab == nil {
		return nil
	}

	regionID := ""
	if ab.RegionID != nil {
		regionID = *ab.RegionID
	} else if ab.RegionCode != "" {
		regionID = ab.RegionCode
	}

	districtID := ""
	if ab.DistrictID != nil {
		districtID = *ab.DistrictID
	} else if ab.DistrictCode != "" {
		districtID = ab.DistrictCode
	}

	var parent *imodel.AccountBlock
	if ab.Parent != nil {
		parent = ab.Parent
	}

	dr := ab.DisableReason
	if dr == nil {
		dr = []imodel.AccountBlockReason{}
	}

	return &ab_dto.AccountBlockResponse{
		ID:            ab.ID,
		Name:          ab.Name,
		Code:          ab.Code,
		Address:       ab.Address,
		Slug:          ab.Slug,
		Type:          string(ab.Type),
		BranchType:    ab.BranchType,
		RegionName:    ab.RegionName,
		RegionCode:    ab.RegionCode,
		DistrictName:  ab.DistrictName,
		DistrictCode:  ab.DistrictCode,
		Parent:        parent,
		RegionID:      regionID,
		DistrictID:    districtID,
		DisableReason: dr,
		IsEnabled:     ab.IsEnabled,
		CreatedAt:     ab.CreatedAt,
		UpdatedAt:     ab.UpdatedAt,
	}
}
