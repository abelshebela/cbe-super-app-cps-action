package accountblock

import (
	"strings"

	"github.com/go-playground/validator/v10"
)

func (e *EnableOrDisableBranches) Clean() {
	for i, code := range e.BranchIds {
		e.BranchIds[i] = strings.TrimSpace(code)
	}
}

func (e *EnableOrDisableRegions) Clean() {
	for i, code := range e.RegionIds {
		e.RegionIds[i] = strings.TrimSpace(code)
	}
}

func (e *EnableOrDisableDistricts) Clean() {
	for i, code := range e.DistrictIds {
		e.DistrictIds[i] = strings.TrimSpace(code)
	}
}

func (e *EnableOrDisableCities) Clean() {
	for i, code := range e.CityIds {
		e.CityIds[i] = strings.TrimSpace(code)
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
