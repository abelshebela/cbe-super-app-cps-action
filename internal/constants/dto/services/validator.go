package services

import (
	"cbe-super-app-cps-action/pkgs/utils"
	"fmt"

	validation "github.com/go-ozzo/ozzo-validation/v4"
	"github.com/go-ozzo/ozzo-validation/v4/is"
)

func (c CapRequest) Validate() error {
	return validation.ValidateStruct(&c,
		validation.Field(&c.Currency, validation.NotNil, is.CurrencyCode),
		validation.Field(&c.SingleCap, validation.NotNil, validation.Min(0.0)),
		validation.Field(&c.MinimumTransferCap, validation.NotNil, validation.Min(0.0), validation.By(func(value interface{}) error {
			if c.SingleCap == nil || c.MinimumTransferCap == nil {
				return nil
			}
			if *c.MinimumTransferCap > *c.SingleCap {
				return fmt.Errorf("minimum transfer cap must not be greater than the single maximum transfer cap")
			}
			return nil
		})),
	)
}

// func (t AdditionalFees) validate() error {
// 	return validation.ValidateStruct(&t,
// 		validation.Field(&t.FeeName, validation.Required, validation.By(utils.NoSpecialChars)),
// 		validation.Field(&t.FeeType, validation.Required, validation.In("PERCENT", "FLAT").Error("fee type must be either 'PERCENT' or 'FLAT'")),
// 		validation.Field(
// 			&t.FeeAmount,
// 			validation.NotNil,
// 			validation.Min(0.0),
// 			validation.By(func(value interface{}) error {
// 				if *t.FeeType == "PERCENT" && t.FeeAmount != nil && *t.FeeAmount > 100 {
// 					return fmt.Errorf("fee amount cannot be greater than 100 percent when fee type is PERCENT")
// 				}
// 				return nil
// 			}),
// 		),
// 	)
// }

// func (t TierRequest) Validate() error {
// 	return validation.ValidateStruct(&t,
// 		validation.Field(&t.Min, validation.NotNil, validation.Min(0.0)),
// 		validation.Field(&t.Max, validation.NotNil, validation.Min(0.0), validation.By(func(value interface{}) error {
// 			if t.Min == nil || t.Max == nil {
// 				return nil
// 			}
// 			if *t.Max < *t.Min {
// 				return fmt.Errorf("maximum amount cannot be less than minimum amount")
// 			}
// 			return nil
// 		})),
// 		validation.Field(&t.FeeType, validation.Required, validation.In("PERCENT", "FLAT").Error("fee type must be either 'PERCENT' or 'FLAT'")),
// 		validation.Field(
// 			&t.FeeAmount,
// 			validation.NotNil,
// 			validation.Min(0.0),
// 			validation.By(func(value interface{}) error {
// 				if *t.FeeType == "PERCENT" && t.FeeAmount != nil && *t.FeeAmount > 100 {
// 					return fmt.Errorf("fee amount cannot be greater than 100 percent when fee type is PERCENT")
// 				}
// 				return nil
// 			}),
// 		),
// 		// validation.Field(&t.AdditionalFees, validation.When(t.AdditionalFees != nil, validation.By(func(value interface{}) error {
// 		// 	if t.AdditionalFees == nil {
// 		// 		return nil
// 		// 	}
// 		// 	return t.AdditionalFees.validate()
// 		// }))),
// 	)
// }

// func validateCapAndTiers(cap CapRequest, tiers []TierRequest) error {
// 	if err := cap.Validate(); err != nil {
// 		return err
// 	}

// 	for i, tier := range tiers {
// 		if err := tier.Validate(); err != nil {
// 			return fmt.Errorf("tier %d: %w", i, err)
// 		}

// 		// Tier chaining: first tier's max amount is the seconds min amount
// 		if i > 0 {
// 			if *tier.Min != *tiers[i-1].Max {
// 				return fmt.Errorf("tier %d: minimum amount (%f) must be equal to previous tier's maximum amount (%f)", i, *tier.Min, *tiers[i-1].Max)
// 			}
// 		}

// 		// Tiers max amount cannot be greater than the single maximum transfer cap
// 		if cap.SingleCap != nil && *tier.Max > *cap.SingleCap {
// 			return fmt.Errorf("tier %d: max amount (%f) cannot be greater than single maximum transfer cap (%f)", i, *tier.Max, *cap.SingleCap)
// 		}

// 		// Tiers min amount cannot be less than the minimum transfer cap
// 		if i == 0 && cap.MinimumTransferCap != nil && *tier.Min < *cap.MinimumTransferCap {
// 			return fmt.Errorf("first tier: min amount (%f) cannot be less than minimum transfer cap (%f)", *tier.Min, *cap.MinimumTransferCap)
// 		}
// 	}
// 	return nil
// }

func validateCap(cap CapRequest) error {
	if err := cap.Validate(); err != nil {
		return err
	}

	return nil
}

// func (s ServiceList) Validate() error {
// 	err := validation.ValidateStruct(&s,
// 		validation.Field(&s.ServiceName, validation.Required, validation.By(utils.NoSpecialChars)),
// 		validation.Field(&s.ServiceKey, validation.Required, validation.By(utils.NoSpecialChars)),
// 		validation.Field(&s.OverideProductGlAccount),
// 	)
// 	if err != nil {
// 		return err
// 	}

// 	if BoolPointer(s.HaveAnOverideTiers, false) {
// 		err := validation.ValidateStruct(&s,
// 			validation.Field(&s.OverideCap, validation.Required),
// 			// validation.Field(&s.OverideTiers, validation.Required, validation.Length(1, 0)),
// 		)
// 		if err != nil {
// 			return err
// 		}
// 		// return validateCapAndTiers(*s.OverideCap, s.OverideTiers)
// 		return validateCap(*s.OverideCap)
// 	}

// 	return nil
// }

func (r CreateServiceRequest) Validate() error {
	err := validation.ValidateStruct(&r,
		validation.Field(&r.ServiceName, validation.Required, validation.By(utils.NoSpecialChars)),
		validation.Field(&r.ServiceKey, validation.Required, validation.By(utils.NoSpecialChars)),
		validation.Field(&r.ServiceCode, validation.By(utils.NoSpecialChars)),
		validation.Field(&r.ProductGlAccount),
	)
	if err != nil {
		return err
	}

	for _, c := range r.Cap {
		if err := validateCap(c); err != nil {
			return err
		}
	}

	// if r.HaveATier {
	// 	err := validation.ValidateStruct(&r,
	// 		validation.Field(&r.Cap, validation.Required),
	// 		// validation.Field(&r.Tiers, validation.Required, validation.Length(1, 0)),
	// 	)
	// 	if err != nil {
	// 		return err
	// 	}
	// 	// if err := validateCapAndTiers(r.Cap, r.Tiers); err != nil {
	// 	if err := validateCap(r.Cap); err != nil {
	// 		return err
	// 	}
	// }

	// if r.HaveAChild {
	// 	err := validation.ValidateStruct(&r,
	// 		validation.Field(&r.ServiceList, validation.Required, validation.Length(1, 0), validation.Each(validation.Required)),
	// 	)
	// 	if err != nil {
	// 		return err
	// 	}
	// 	for i, sl := range r.ServiceList {
	// 		if err := sl.Validate(); err != nil {
	// 			return fmt.Errorf("service list %d: %w", i, err)
	// 		}
	// 	}
	// }

	return nil
}

func (r UpdateServiceRequest) Validate() error {
	err := validation.ValidateStruct(&r,
		validation.Field(&r.ServiceCode, validation.By(utils.NoSpecialChars)),
	)
	if err != nil {
		return err
	}

	if r.Cap != nil {
		for _, c := range r.Cap {
			if c.SingleCap != nil || c.MinimumTransferCap != nil {
				if err := validateCap(c); err != nil {
					return err
				}
			}
		}
	}

	// if BoolPointer(r.HaveATier, false) {
	// 	// if r.Cap != nil && (r.Cap.SingleCap != nil || r.Cap.MinimumTransferCap != nil) || len(r.Tiers) > 0 {
	// 	if r.Cap != nil && (r.Cap.SingleCap != nil || r.Cap.MinimumTransferCap != nil) {
	// 		// if err := validateCapAndTiers(*r.Cap, r.Tiers); err != nil {
	// 		if err := validateCap(*r.Cap); err != nil {
	// 			return err
	// 		}
	// 	}
	// }

	// if BoolPointer(r.HaveAChild, false) {
	// 	err := validation.ValidateStruct(&r,
	// 		validation.Field(&r.ServiceList, validation.Each(validation.Required)),
	// 	)
	// 	if err != nil {
	// 		return err
	// 	}
	// 	for i, sl := range r.ServiceList {
	// 		if err := sl.Validate(); err != nil {
	// 			return fmt.Errorf("service list %d: %w", i, err)
	// 		}
	// 	}
	// }

	return nil
}

func (r CreateServiceList) Validate() error {
	return validation.ValidateStruct(&r,
		validation.Field(&r.ServiceKey, validation.Required, validation.By(utils.NoSpecialChars)),
		validation.Field(&r.ServiceName, validation.Required, validation.By(utils.NoSpecialChars)),
	)
}

func (r UpdateServiceList) Validate() error {
	return validation.ValidateStruct(&r,
		validation.Field(&r.ServiceKey, validation.By(utils.NoSpecialChars)),
		validation.Field(&r.ServiceName, validation.By(utils.NoSpecialChars)),
	)
}
