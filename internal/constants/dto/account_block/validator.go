package accountblock

import (
	"errors"
	"regexp"
	"strings"

	validation "github.com/go-ozzo/ozzo-validation"
	"github.com/go-playground/validator/v10"
)

var mongoIDRegex = regexp.MustCompile(`^[a-fA-F0-9]{24}$`)

func ValidateMongoID(id string) error {
	return validation.Validate(
		id,
		validation.Required,
		validation.Match(mongoIDRegex).Error("invalid MongoDB ObjectID"),
	)
}

var validate = validator.New()

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
	if err := validate.Struct(e); err != nil {
		return err
	}
	for _, id := range e.BranchIds {
		if err := ValidateMongoID(id); err != nil {
			return errors.New("invalid branch ID: " + id)
		}
	}
	return nil
}

func (e *EnableOrDisableRegions) Validate() error {
	if err := validate.Struct(e); err != nil {
		return err
	}
	for _, id := range e.RegionIds {
		if err := ValidateMongoID(id); err != nil {
			return errors.New("invalid region ID: " + id)
		}
	}
	return nil
}

func (e *EnableOrDisableDistricts) Validate() error {
	if err := validate.Struct(e); err != nil {
		return err
	}
	for _, id := range e.DistrictIds {
		if err := ValidateMongoID(id); err != nil {
			return errors.New("invalid district ID: " + id)
		}
	}
	return nil
}

func (e *EnableOrDisableCities) Validate() error {
	if err := validate.Struct(e); err != nil {
		return err
	}
	for _, id := range e.CityIds {
		if err := ValidateMongoID(id); err != nil {
			return errors.New("invalid city ID: " + id)
		}
	}
	return nil
}
