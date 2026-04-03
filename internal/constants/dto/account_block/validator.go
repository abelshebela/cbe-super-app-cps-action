package accountblock

import (
	"errors"
	"strings"
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

func (e *EnableOrDisableBranches) Validate() error {
	if e.Reason == "" {
		return errors.New("reason is required")
	}
	for _, id := range e.BranchIds {
		if strings.TrimSpace(id) == "" {
			return errors.New("invalid branch ID: empty value")
		}
	}
	return nil
}

func (e *EnableOrDisableRegions) Validate() error {
	if e.Reason == "" {
		return errors.New("reason is required")
	}
	for _, id := range e.RegionIds {
		if strings.TrimSpace(id) == "" {
			return errors.New("invalid region ID: empty value")
		}
	}
	return nil
}

func (e *EnableOrDisableDistricts) Validate() error {
	if e.Reason == "" {
		return errors.New("reason is required")
	}
	for _, id := range e.DistrictIds {
		if strings.TrimSpace(id) == "" {
			return errors.New("invalid district ID: empty value")
		}
	}
	return nil
}

func (e *EnableOrDisableCities) Validate() error {
	if e.Reason == "" {
		return errors.New("reason is required")
	}
	for _, id := range e.CityIds {
		if strings.TrimSpace(id) == "" {
			return errors.New("invalid city ID: empty value")
		}
	}
	return nil
}
