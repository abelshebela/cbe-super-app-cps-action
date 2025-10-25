package accountblock

// EnableOrDisableBranches represents the request to enable or disable branches
type EnableOrDisableBranches struct {
	BranchCodes []string `json:"branch_codes" example:"BR001,BR002"`
	Reason      string   `json:"reason" example:"Maintenance"`
}

// EnableOrDisableRegions represents the request to enable or disable regions
type EnableOrDisableRegions struct {
	RegionCodes []string `json:"region_codes" example:"RG001,RG002"`
	Reason      string   `json:"reason" example:"Policy Update"`
}

// EnableOrDisableDistricts represents the request to enable or disable districts
type EnableOrDisableDistricts struct {
	DistrictCodes []string `json:"district_codes" example:"DS001,DS002"`
	Reason        string   `json:"reason" example:"Operational Change"`
}

// EnableOrDisableCities represents the request to enable or disable cities
type EnableOrDisableCities struct {
	CityCodes []string `json:"city_codes" example:"CT001,CT002"`
	Reason    string   `json:"reason" example:"Regulatory Compliance"`
}
