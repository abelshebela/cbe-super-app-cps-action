package core

import ( "cbe-super-app-cps-action/internal/constants/model"

	ab_dto "cbe-super-app-cps-action/internal/constants/dto/account_block"
)


func ToBranchResponse(branch *model.Branch) *ab_dto.BranchResponse {
	return &ab_dto.BranchResponse{
		ID:          branch.ID.Hex(),
		BranchCode:  branch.BranchCode,
		BranchName:  branch.BranchName,
		BranchAddress: branch.BranchAddress,
		DistrictCode: branch.DistrictCode,
		DistrictName: branch.DistrictName,
		BranchRegion: branch.BranchRegion,
		RecordStat:  branch.RecordStat,
		CreatedAt:   branch.CreatedAt,
		IsDeleted:   branch.IsDeleted,
		UpdatedAt:   branch.UpdatedAt,
		Enabled:     branch.Enabled,		
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
		ID: region.ID.Hex(),
		RegionName: region.RegionName,
		RegionCode: region.RegionCode,
		RegionAddress: region.RegionAddress,
		Enabled: region.Enabled,
		CreatedAt: region.CreatedAt,
		UpdatedAt: region.UpdatedAt,
	}
}

func ToRegionsResponse (regions []*model.)