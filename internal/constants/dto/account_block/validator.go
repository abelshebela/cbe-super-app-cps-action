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
	for i, name := range e.FederalRegionNames {
		e.FederalRegionNames[i] = strings.TrimSpace(name)
	}
}

func (e *EnableOrDisableDistricts) Clean() {
	for i, name := range e.DistrictNames {
		e.DistrictNames[i] = strings.TrimSpace(name)
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
	for _, name := range e.FederalRegionNames {
		if strings.TrimSpace(name) == "" {
			return errors.New("invalid federal region name: empty value")
		}
	}
	return nil
}

func (e *EnableOrDisableDistricts) Validate() error {
	if e.Reason == "" {
		return errors.New("reason is required")
	}
	for _, name := range e.DistrictNames {
		if strings.TrimSpace(name) == "" {
			return errors.New("invalid district name: empty value")
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
