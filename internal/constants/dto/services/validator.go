package services

import (
	"cbe-super-app-cps-action/pkgs/utils"
	"fmt"

	validation "github.com/go-ozzo/ozzo-validation/v4"
)

func (c CapRequest) Validate() error {
	return validation.ValidateStruct(&c,
		validation.Field(&c.SingleCap, validation.Required, validation.Min(0.0)),
		validation.Field(&c.MinimumTransferCap, validation.Required, validation.Min(0.0), validation.Max(c.SingleCap).Error("minimum transfer cap must not be greater than the single maximum transfer cap")),
	)
}

func (t TierRequest) Validate() error {
	return validation.ValidateStruct(&t,
		validation.Field(&t.Min, validation.Required, validation.Min(0.0)),
		validation.Field(&t.Max, validation.Required, validation.Min(0.0), validation.Min(t.Min).Error("maximum amount cannot be less than minimum amount")),
		validation.Field(&t.FeeType, validation.Required, validation.In("PERCENT", "FLAT").Error("fee type must be either 'PERCENT' or 'FLAT'")),
		validation.Field(&t.FeeAmount, validation.Required, validation.Min(0.0)),
	)
}

func validateCapAndTiers(cap CapRequest, tiers []TierRequest) error {
	if err := cap.Validate(); err != nil {
		return err
	}

	for i, tier := range tiers {
		if err := tier.Validate(); err != nil {
			return fmt.Errorf("tier %d: %w", i, err)
		}

		// Tier chaining: first tier's max amount is the seconds min amount
		if i > 0 {
			if tier.Min != tiers[i-1].Max {
				return fmt.Errorf("tier %d: minimum amount (%f) must be equal to previous tier's maximum amount (%f)", i, tier.Min, tiers[i-1].Max)
			}
		}

		// Tiers max amount cannot be greater than the single maximum transfer cap
		if tier.Max > cap.SingleCap {
			return fmt.Errorf("tier %d: max amount (%f) cannot be greater than single maximum transfer cap (%f)", i, tier.Max, cap.SingleCap)
		}

		// Tiers min amount cannot be less than the minimum transfer cap
		if i == 0 && tier.Min < cap.MinimumTransferCap {
			return fmt.Errorf("first tier: min amount (%f) cannot be less than minimum transfer cap (%f)", tier.Min, cap.MinimumTransferCap)
		}
	}
	return nil
}

func (s ServiceList) Validate() error {
	err := validation.ValidateStruct(&s,
		validation.Field(&s.ServiceName, validation.Required, validation.By(utils.NoSpecialChars)),
		validation.Field(&s.ServiceKey, validation.Required, validation.By(utils.NoSpecialChars)),
		validation.Field(&s.OverideProductGlAccount, validation.Required),
	)
	if err != nil {
		return err
	}

	return validateCapAndTiers(s.OverideCap, s.OverideTiers)
}

func (r CreateServiceRequest) Validate() error {
	err := validation.ValidateStruct(&r,
		validation.Field(&r.ServiceName, validation.Required, validation.By(utils.NoSpecialChars)),
		validation.Field(&r.ServiceKey, validation.Required, validation.By(utils.NoSpecialChars)),
		validation.Field(&r.ServiceCode, validation.Required, validation.By(utils.NoSpecialChars)),
		validation.Field(&r.ProductGlAccount, validation.Required),
	)
	if err != nil {
		return err
	}

	if r.HaveATier {
		err := validation.ValidateStruct(&r,
			validation.Field(&r.Cap, validation.Required),
			validation.Field(&r.Tiers, validation.Required, validation.Length(1, 0)),
		)
		if err != nil {
			return err
		}
		if err := validateCapAndTiers(r.Cap, r.Tiers); err != nil {
			return err
		}
	}

	if r.HaveAChild {
		err := validation.ValidateStruct(&r,
			validation.Field(&r.ServiceList, validation.Required, validation.Length(1, 0), validation.Each(validation.Required)),
		)
		if err != nil {
			return err
		}
		for i, sl := range r.ServiceList {
			if err := sl.Validate(); err != nil {
				return fmt.Errorf("service list %d: %w", i, err)
			}
		}
	}

	return nil
}

func (r UpdateServiceRequest) Validate() error {
	err := validation.ValidateStruct(&r,
		validation.Field(&r.ServiceCode, validation.By(utils.NoSpecialChars)),
	)
	if err != nil {
		return err
	}

	if r.HaveATier {
		if r.Cap.SingleCap != 0 || r.Cap.MinimumTransferCap != 0 || len(r.Tiers) > 0 {
			if err := validateCapAndTiers(r.Cap, r.Tiers); err != nil {
				return err
			}
		}
	}

	if r.HaveAChild {
		err := validation.ValidateStruct(&r,
			validation.Field(&r.ServiceList, validation.Each(validation.Required)),
		)
		if err != nil {
			return err
		}
		for i, sl := range r.ServiceList {
			if err := sl.Validate(); err != nil {
				return fmt.Errorf("service list %d: %w", i, err)
			}
		}
	}

	return nil
}
