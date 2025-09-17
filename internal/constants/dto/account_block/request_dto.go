package accountblock

type EnableOrDisableBranches struct {
	BranchCodes []string `json:"branches_code"`
}

type EnableOrDisableRegions struct {
	RegionsCodes []string `json:"regions_code"`
}

type EnableOrDisableDistricts struct {
	DistrictCodes []string `json:"districts_code"`
}

type EnableOrDisableCities struct {
	CitiesCode []string `json:"city_code"`
}
