package accountblock

import "strings"

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
