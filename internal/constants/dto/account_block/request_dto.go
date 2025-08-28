package accountblock

import "strings"

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
	CitiesCode []string `json:"cities_code"`
}

func (e *EnableOrDisableBranches) Clean() {
	for i, code := range e.BranchCodes {
		e.BranchCodes[i] = strings.TrimSpace(code)
	}
}

func (e *EnableOrDisableRegions) Clean() {
	for i, code := range e.RegionsCodes {
		e.RegionsCodes[i] = strings.TrimSpace(code)
	}
}

func (e *EnableOrDisableDistricts) Clean() {
	for i, code := range e.DistrictCodes {
		e.DistrictCodes[i] = strings.TrimSpace(code)
	}
}

func (e *EnableOrDisableCities) Clean() {
	for i, code := range e.CitiesCode {
		e.CitiesCode[i] = strings.TrimSpace(code)
	}
}
