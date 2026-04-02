package core

import (
	imodel "cbe-super-app-cps-action/internal/constants/model"

	ab_dto "cbe-super-app-cps-action/internal/constants/dto/account_block"
)

func ToAccountBlockResponse(ab *imodel.AccountBlock) *ab_dto.AccountBlockResponse {
	var parentID string
	if ab.ParentID != nil {
		parentID = *ab.ParentID
	}

	// Convert parent to shared-compatible form if present
	var parent *imodel.AccountBlock
	if ab.Parent != nil {
		parent = ab.Parent
	}

	return &ab_dto.AccountBlockResponse{
		ID:            ab.ID,
		Name:          ab.Name,
		Code:          ab.Code,
		Address:       ab.Address,
		Slug:          ab.Slug,
		Type:          string(ab.Type),
		ParentID:      parentID,
		Parent:        parent,
		DisableReason: ab.DisableReason,
		IsEnabled:     ab.IsEnabled,
		CreatedAt:     ab.CreatedAt,
		UpdatedAt:     ab.UpdatedAt,
	}
}
