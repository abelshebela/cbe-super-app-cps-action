package core

import (
	helper "cbe-super-app-cps-action/internal/constants/dto/vaultgroup_category"

	"cbe-super-app-cps-action/internal/constants/model"
	"time"
)

// map dto to model for create

func ToDomainCreateVaultGroupCategoryRequest(req helper.CreateVaultGroupCategoryRequest) *model.VaultGroupCategory {
	return &model.VaultGroupCategory{
		Name:        req.Name,
		Description: req.Description,
		IsActive:    false,
	}
}

// map updaterequest dto to model

func ToDomainUpdateVaultGroupCategoryRequest(req helper.UpdateVaultGroupCategoryRequest) *model.VaultGroupCategory {
	return &model.VaultGroupCategory{
		Name: func() string {
			if req.Name != nil {
				return *req.Name
			}
			return ""
		}(),
		Description: func() string {
			if req.Description != nil {
				return *req.Description
			}
			return ""
		}(),
		UpdatedAt: time.Now().UTC(),
	}
}
