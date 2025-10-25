package accountblock

import "strings"

func (e *EnableOrDisableBranches) Clean() {
	for i, code := range e.BranchCodes {
		e.BranchCodes[i] = strings.TrimSpace(code)
	}
}

func (e *EnableOrDisableRegions) Clean() {
	for i, code := range e.RegionCodes {
		e.RegionCodes[i] = strings.TrimSpace(code)
	}
}

func (e *EnableOrDisableDistricts) Clean() {
	for i, code := range e.DistrictCodes {
		e.DistrictCodes[i] = strings.TrimSpace(code)
	}
}

func (e *EnableOrDisableCities) Clean() {
	for i, code := range e.CityCodes {
		e.CityCodes[i] = strings.TrimSpace(code)
	}
}
