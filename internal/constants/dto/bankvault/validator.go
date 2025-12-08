package bankvault

import (
	"errors"
	"regexp"
	"strings"

	validation "github.com/go-ozzo/ozzo-validation"
	"github.com/shopspring/decimal"
)

func (r *CreateBankVaultProductRequest) Validate() error {
	if r == nil {
		return errors.New("request is required")
	}

	var err error
	if r.Name, err = sanitizeString(r.Name); err != nil {
		return err
	}
	// if r.Description, err = sanitizeString(r.Description); err != nil {
	// 	return err
	// }
	// if r.Currency, err = sanitizeString(r.Currency); err != nil {
	// 	return err
	// }
	// r.Currency = strings.ToUpper(r.Currency)

	// if r.Method, err = sanitizeString(r.Method); err != nil {
	// 	return err
	// }
	// r.Method = strings.ToUpper(r.Method)

	// if r.Frequency, err = sanitizeString(r.Frequency); err != nil {
	// 	return err
	// }
	// r.Frequency = strings.ToUpper(r.Frequency)

	return validation.ValidateStruct(r,
		validation.Field(&r.Name,
			validation.Required,
			validation.Length(3, 100),
			validation.By(noSpecialChars),
		),
		// validation.Field(&r.Currency,
		// 	validation.Required,
		// 	validation.Match(reISO4217),
		// 	validation.By(noSpecialChars),
		// ),
		validation.Field(&r.LockPeriodDays,
			validation.Required,
			validation.Match(regexp.MustCompile(`^\d+(d|m|y)$`)).Error("must be a number followed by 'd', 'm', or 'y'"),
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
		validation.Field(&r.Frequency, validation.Required, validation.By(func(value interface{}) error {
			if v, ok := value.(int64); ok {
				if v <= 0 {
					return errors.New("frequency must be > 0")
				}
				if v > 365 {
					return errors.New("frequency must be <= 365")
				}
			} else {
				return errors.New("value must be an integer")
			}

			return nil
		})),
		// validation.Field(&r.RateBps, validation.By(func(value interface{}) error {
		// 	if v, ok := value.(decimal.Decimal); ok {
		// 		if v.Equal(decimal.Zero) {
		// 			return errors.New("rate_bps cannnot be empty")
		// 		}
		// 		if v.LessThan(decimal.Zero) {
		// 			return errors.New("rate_bps must be >= 0")
		// 		}
		// 	}
		// 	return nil
		// })),
		// validation.Field(&r.Method, validation.Required, validation.In("SIMPLE", "COMPOUND"), validation.By(noSpecialChars)),
		// validation.Field(&r.Description, validation.Required, validation.By(noSpecialChars)),
	)
}

func (r *UpdateBankVaultProductRequest) Validate() error {
	if r == nil {
		return errors.New("request is required")
	}

	return validation.ValidateStruct(r,
		validation.Field(&r.MinAmount, validation.By(func(value interface{}) error {
			if v, ok := value.(*float64); ok && v != nil {
				if decimal.NewFromFloat(*v).LessThanOrEqual(decimal.Zero) {
					return errors.New("min_amount must be > 0")
				}
			}
			return nil
		})),
		validation.Field(&r.MaxAmount, validation.By(func(value interface{}) error {
			if v, ok := value.(*float64); ok && v != nil {
				if decimal.NewFromFloat(*v).LessThanOrEqual(decimal.Zero) {
					return errors.New("max_amount must be > 0")
				}
			}
			return nil
		})),
		// 👇 Attach cross-field validation to one of the fields (e.g., MaxAmount)
		validation.Field(&r.MaxAmount, validation.By(func(_ interface{}) error {
			if r.MinAmount != nil && r.MaxAmount != nil {
				if decimal.NewFromFloat(*r.MinAmount).GreaterThan(decimal.NewFromFloat(*r.MaxAmount)) {
					return errors.New("min_amount must be <= max_amount")
				}
			}
			return nil
		})),
	)
}

func sanitizeString(s string) (string, error) {
	trimmed := strings.TrimSpace(s)
	return trimmed, nil
}

var reISO4217 = regexp.MustCompile(`^[A-Z]{3}$`)

func noSpecialChars(value interface{}) error {
	var s string

	switch v := value.(type) {
	case string:
		s = v
	case *string:
		if v != nil {
			s = *v
		}
	default:
		return nil
	}

	if s == "" {
		return nil
	}

	// allow letters, numbers, space, dot, underscore, dash
	re := regexp.MustCompile(`^[a-zA-Z0-9 ._-]+$`)
	if !re.MatchString(s) {
		// return validation.NewError("validation_no_special_chars", "contains invalid characters")
		return errors.New("contains invalid characters")
	}
	return nil
}
