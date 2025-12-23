package servicesdto

import (
	"fmt"
	"strings"

	local_utils "cbe-super-app-cps-action/pkgs/utils"

	validation "github.com/go-ozzo/ozzo-validation/v4"
)

func (c *CapRequest) Validate() error {
	if c == nil {
		return nil
	}
	if err := validation.ValidateStruct(
		c,
		validation.Field(&c.SingleCap,
			validation.Required.Error("cap.single_cap is required"),
			validation.Min(int64(1)).Error("cap.single_cap must be greater than 0"),
		),
		validation.Field(&c.MinimumTransferCap,
			validation.Required.Error("cap.minimum_transfer_cap is required"),
			validation.Min(int64(1)).Error("cap.minimum_transfer_cap must be greater than 0"),
		),
	); err != nil {
		return err
	}
	if c.MinimumTransferCap > c.SingleCap {
		return fmt.Errorf("cap.minimum_transfer_cap (%d) cannot exceed cap.single_cap (%d)", c.MinimumTransferCap, c.SingleCap)
	}
	return nil
}

func (t *TierRequest) Validate() error {
	if t == nil {
		return nil
	}
	if err := validation.ValidateStruct(
		t,
		validation.Field(&t.Min,
			validation.Required.Error("tiers.min is required"),
			validation.Min(int64(0)).Error("tiers.min must be >= 0"),
		),
		validation.Field(&t.Max,
			validation.Required.Error("tiers.max is required"),
			validation.Min(int64(0)).Error("tiers.max must be >= 0"),
		),
		validation.Field(&t.FeeAmount,
			validation.Required.Error("tiers.fee_amount is required"),
			validation.Min(uint8(0)).Error("tiers.fee_amount must be >= 0"),
		),
	); err != nil {
		return err
	}
	if t.Min >= t.Max {
		return fmt.Errorf("tiers.min (%d) must be less than tiers.max (%d)", t.Min, t.Max)
	}
	if strings.EqualFold(t.FeeType, "PERCENT") && t.FeeAmount > 100 {
		return fmt.Errorf("tiers.fee_amount cannot exceed 100 when fee_type is PERCENT")
	}
	return nil
}

func (r *CreateServiceRequest) Validate() error {
	if err := validation.ValidateStruct(
		r,
		validation.Field(&r.ServiceCode,
			validation.Required.Error("service_code is required"),
			validation.By(local_utils.NoSpecialChars),
		),
		validation.Field(&r.ServiceName,
			validation.Required.Error("service_name is required"),
		),
		validation.Field(&r.ChargeCode,
			validation.Required.Error("charge_code is required"),
			validation.By(local_utils.NoSpecialChars),
		),
		// CommissionCode optional
		validation.Field(&r.CbeGLProductAccount,
			validation.Required.Error("cbe_gl_product_account is required"),
			validation.By(local_utils.NoSpecialChars),
		),
		validation.Field(&r.PaymentType,
			validation.Required.Error("payment_type is required"),
		),
		validation.Field(&r.Tiers,
			validation.Required.Error("tiers are required"),
			validation.Length(1, 0).Error("at least one tier is required"),
		),
	); err != nil {
		return err
	}
	if err := r.Cap.Validate(); err != nil {
		return err
	}
	for i := range r.Tiers {
		if err := r.Tiers[i].Validate(); err != nil {
			return fmt.Errorf("tier %d validation failed: %w", i+1, err)
		}
	}
	return nil
}

func (r *UpdateServiceRequest) Validate() error {
	// All fields optional; validate if provided
	// if r.ServiceCode != "" {
	// 	if err := validation.Validate(&r.ServiceCode, validation.By(local_utils.NoSpecialChars)); err != nil {
	// 		return err
	// 	}
	// }
	// if r.ServiceType != "" {
	// 	if err := validation.Validate(&r.ServiceType, validation.By(local_utils.NoSpecialChars)); err != nil {
	// 		return err
	// 	}
	// }
	// if r.ChargeCode != "" {
	// 	if err := validation.Validate(&r.ChargeCode, validation.By(local_utils.NoSpecialChars)); err != nil {
	// 		return err
	// 	}
	// }
	if r.CbeGLProductAccount != "" {
		if err := validation.Validate(&r.CbeGLProductAccount, validation.By(local_utils.NoSpecialChars)); err != nil {
			return err
		}
	}
	// Validate cap if provided (any field set)
	if r.Cap.MinimumTransferCap != 0 || r.Cap.SingleCap != 0 || r.Cap.KYCLevel != "" {
		if err := r.Cap.Validate(); err != nil {
			return err
		}
	}
	// Validate tiers if provided
	if len(r.Tiers) > 0 {
		for i := range r.Tiers {
			if err := r.Tiers[i].Validate(); err != nil {
				return fmt.Errorf("tier %d validation failed: %w", i+1, err)
			}
		}
	}
	return nil
}
