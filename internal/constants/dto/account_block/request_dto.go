package accountblock

// EnableOrDisableBranches represents the request to enable or disable branches
type EnableOrDisableBranches struct {
	BranchCodes []string `json:"branches_code" example:"BR001,BR002"`
}

// EnableOrDisableRegions represents the request to enable or disable regions
type EnableOrDisableRegions struct {
	RegionsCodes []string `json:"regions_code" example:"RG001,RG002"`
}

// EnableOrDisableDistricts represents the request to enable or disable districts
type EnableOrDisableDistricts struct {
	DistrictCodes []string `json:"districts_code" example:"DS001,DS002"`
}

// EnableOrDisableCities represents the request to enable or disable cities
type EnableOrDisableCities struct {
	CitiesCode []string `json:"city_code" example:"CT001,CT002"`
}
