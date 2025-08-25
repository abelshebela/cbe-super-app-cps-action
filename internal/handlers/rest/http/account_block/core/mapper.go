package core

import "cbe-super-app-cps-action/internal/constants/model"


func ToBranchResponse(branch *model.Branch) ab_dto.BranchResponse {
	return ab_dto.BranchResponse{
		ID:          branch.ID.Hex(),
		BranchCode:  branch.BranchCode,
		BranchName:  branch.BranchName,
		RegionCode:  branch.RegionCode,
		DistrictCode: branch.DistrictCode,
		CityCode:    branch.CityCode,
		Address:     branch.Address,
		Phone:       branch.Phone,
		Email:       branch.Email,
		Enabled:     branch.Enabled,
		CreatedAt:   branch.CreatedAt,
		LastUpdatedAt: branch.LastUpdatedAt,
	}

}