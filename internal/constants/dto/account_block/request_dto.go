package accountblock

// EnableOrDisableBranches represents the request to enable or disable branches
type EnableOrDisableBranches struct {
	BranchIds []string `json:"branch_Ids" example:"674003000000000000000001" validate:"required,min=1"`
	Reason    string   `json:"reason" example:"Maintenance" validate:"required"`
}

// EnableOrDisableRegions represents the request to enable or disable regions
type EnableOrDisableRegions struct {
	RegionIds []string `json:"region_ids" example:"674003000000000000000001" validate:"required,min=1"`
	Reason    string   `json:"reason" example:"Policy Update" validate:"required"`
}

// EnableOrDisableDistricts represents the request to enable or disable districts
type EnableOrDisableDistricts struct {
	DistrictIds []string `json:"district_ids" example:"674003000000000000010001" validate:"required,min=1"`
	Reason      string   `json:"reason" example:"Operational Change" validate:"required"`
}

// EnableOrDisableCities represents the request to enable or disable cities
type EnableOrDisableCities struct {
	CityCodes []string `json:"city_codes" example:"CT001,CT002" validate:"required,min=1"`
	Reason    string   `json:"reason" example:"Regulatory Compliance" validate:"required"`
}
