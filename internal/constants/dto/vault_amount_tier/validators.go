package vaultamounttier

import (
	"errors"

	validation "github.com/go-ozzo/ozzo-validation"
	"github.com/shopspring/decimal"
)

func (r *VaultAmountTierRequest) Validate(isCreate bool) error {
	return validation.ValidateStruct(r,
		validation.Field(&r.VaultCategoryID,
			validation.Required,
		),
		validation.Field(&r.MinAmount, validation.By(func(value interface{}) error {
			if v, ok := value.(decimal.Decimal); ok {
				if v.Equal(decimal.Zero) {
					return errors.New("min_amount is required")
				}
				if v.LessThanOrEqual(decimal.Zero) {
					return errors.New("min_amount must be > 0")
				}
			}
			return nil
		})),
		validation.Field(&r.MaxAmount, validation.By(func(value interface{}) error {
			if v, ok := value.(decimal.Decimal); ok {
				if v.Equal(decimal.Zero) {
					return errors.New("max_amount is required")
				}
				if v.LessThanOrEqual(decimal.Zero) {
					return errors.New("max_amount must be > 0")
				}
				if !r.MinAmount.IsZero() && r.MinAmount.GreaterThan(v) {
					return errors.New("min_amount cannot be greater than max_amount")
				}
			}
			return nil
		})),
		validation.Field(&r.Interest, validation.By(func(value interface{}) error {
			if v, ok := value.(decimal.Decimal); ok {
				if v.Equal(decimal.Zero) {
					return errors.New("interest cannnot be empty")
				}
				if v.LessThan(decimal.Zero) {
					return errors.New("interest must be >= 0")
				}
			}
			return nil
		})),
	)
}

func (r *UpdateVaultAmountTierRequest) Validate() error {
	return validation.ValidateStruct(r,
		validation.Field(&r.MinAmount, validation.By(func(value interface{}) error {
			if v, ok := value.(decimal.Decimal); ok {
				if v.Equal(decimal.Zero) {
					return errors.New("min_amount is required")
				}
				if v.LessThanOrEqual(decimal.Zero) {
					return errors.New("min_amount must be > 0")
				}
			}
			return nil
		})),
		validation.Field(&r.MaxAmount, validation.By(func(value interface{}) error {
			if v, ok := value.(decimal.Decimal); ok {
				if v.Equal(decimal.Zero) {
					return errors.New("max_amount is required")
				}
				if v.LessThanOrEqual(decimal.Zero) {
					return errors.New("max_amount must be > 0")
				}
				if !r.MinAmount.IsZero() && r.MinAmount.GreaterThan(v) {
					return errors.New("min_amount cannot be greater than max_amount")
				}
			}
			return nil
		})),
		validation.Field(&r.Interest, validation.By(func(value interface{}) error {
			if v, ok := value.(decimal.Decimal); ok {
				if v.Equal(decimal.Zero) {
					return errors.New("interest cannnot be empty")
				}
				if v.LessThan(decimal.Zero) {
					return errors.New("interest must be >= 0")
				}
			}
			return nil
		})),
	)
}
