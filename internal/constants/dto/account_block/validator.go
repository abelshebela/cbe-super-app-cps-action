package accountblock

import (
	"strings"

	"github.com/go-playground/validator/v10"
)

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

var validate = validator.New()

func (e *EnableOrDisableBranches) Validate() error {
	return validate.Struct(e)
}

func (e *EnableOrDisableRegions) Validate() error {
	return validate.Struct(e)
}

func (e *EnableOrDisableDistricts) Validate() error {
	return validate.Struct(e)
}

func (e *EnableOrDisableCities) Validate() error {
	return validate.Struct(e)
}
