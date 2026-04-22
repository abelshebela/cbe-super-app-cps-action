package core

import (
	imodel "cbe-super-app-cps-action/internal/constants/model"

	ab_dto "cbe-super-app-cps-action/internal/constants/dto/account_block"
)

func ToAccountBlockResponse(ab *imodel.AccountBlock) *ab_dto.AccountBlockResponse {
	regionID := ab.RegionID
	districtID := ab.DistrictID

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
		Parent:        parent,
		RegionID:      *regionID,
		DistrictID:    *districtID,
		DisableReason: dr,
		IsEnabled:     ab.IsEnabled,
		CreatedAt:     ab.CreatedAt,
		UpdatedAt:     ab.UpdatedAt,
	}
}
