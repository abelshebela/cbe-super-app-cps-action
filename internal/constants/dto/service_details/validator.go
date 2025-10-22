package dto

import (
	"fmt"
	"strings"

	// "cbe-super-app-cps-action/internal/constants"
	local_utils "cbe-super-app-cps-action/pkgs/utils"

	validation "github.com/go-ozzo/ozzo-validation/v4"
	// "github.com/go-ozzo/ozzo-validation/v4/is"
)

func (r *UpdateServiceFeeRequest) Validate() error {
	return validation.ValidateStruct(
		r,
		validation.Field(&r.ID,
			validation.Required.Error("tire_id is required"),
			validation.By(local_utils.NoSpecialChars),
		),
		validation.Field(&r.MinAmount,
			validation.Required.Error("min_amount is required"),
			validation.Min(uint64(1)).Error("min_amount must be greater than 0"),
		),
		validation.Field(&r.MaxAmount,
			validation.Required.Error("max_amount is required"),
			validation.Min(uint64(1)).Error("max_amount must be greater than 0"),
		),
	)
}

func (t *TierDTO) Validate() error {
	err := validation.ValidateStruct(
		t,
		validation.Field(&t.Min,
			validation.Required.Error("min is required"),
			validation.Min(uint64(0)).Error("min must be greater than 0"),
		),
		validation.Field(&t.Max,
			validation.Required.Error("max is required"),
			validation.Min(uint64(0)).Error("max must be greater than 0"),
		),
		validation.Field(&t.FeeAmount,
			validation.Required.Error("fee_amount is required"),
			validation.Min(uint64(0)).Error("fee_amount must be greater than or equal to 0"),
		),
	)
	if err != nil {
		return err
	}

	if t.Min >= t.Max {
		return fmt.Errorf("min amount (%d) must be less than max amount (%d)", t.Min, t.Max)
	}

	return nil
}

func (c *CBglEntryDTO) Validate() error {
	return validation.ValidateStruct(
		c,
		validation.Field(&c.ProductAccount,
			validation.Required.Error("cbgl_product_account is required"),
			validation.By(local_utils.NoSpecialChars),
		),
		validation.Field(&c.ProductBranchCode,
			validation.Required.Error("cbgl_product_branchcode is required"),
			validation.By(local_utils.NoSpecialChars),
		),
		validation.Field(&c.ServiceAccount,
			validation.Required.Error("cbgl_service_account is required"),
			validation.By(local_utils.NoSpecialChars),
		),
		validation.Field(&c.ServiceBranchCode,
			validation.Required.Error("cbgl_service_branchcode is required"),
			validation.By(local_utils.NoSpecialChars),
		),
	)
}

func (i *IFBglEntryDTO) Validate() error {
	return validation.ValidateStruct(
		i,
		validation.Field(&i.ProductAccount,
			validation.Required.Error("ifbgl_product_account is required"),
			validation.By(local_utils.NoSpecialChars),
		),
		validation.Field(&i.ProductBranchCode,
			validation.Required.Error("ifbgl_product_branchcode is required"),
			validation.By(local_utils.NoSpecialChars),
		),
		validation.Field(&i.ServiceAccount,
			validation.Required.Error("ifbgl_service_account is required"),
			validation.By(local_utils.NoSpecialChars),
		),
		validation.Field(&i.ServiceBranchCode,
			validation.Required.Error("ifbgl_service_branchcode is required"),
			validation.By(local_utils.NoSpecialChars),
		),
	)
}

func (dto *ServiceFeeDetailDTO) Validate() error {
	// Basic field validation
	err := validation.ValidateStruct(
		dto,
		validation.Field(&dto.ServiceType,
			validation.Required.Error("service_type is required"),
			validation.By(local_utils.NoSpecialChars),
			validation.By(func(value interface{}) error {
				if str, ok := value.(string); ok {
					if strings.TrimSpace(str) == "" {
						return fmt.Errorf("service_type cannot be empty or whitespace only")
					}
					if strings.ContainsAny(str, "0123456789") {
						return fmt.Errorf("service_type cannot contain numbers")
					}
				}
				return nil
			}),
		),
		validation.Field(&dto.PaymentType,
			validation.Required.Error("payment_type is required"),
			// validation.In("percentage", "flat_fee").Error("payment_type must be either 'percentage' or 'flat_fee'"),
		),
		validation.Field(&dto.SingleCapLevelOne,
			validation.Required.Error("single_cap_level_one is required"),
			validation.Min(0).Error("single_cap_level_one must be greater than or equal to 0"),
		),
		validation.Field(&dto.DailyCapLevelOne,
			validation.Required.Error("daily_cap_level_one is required"),
			validation.Min(0).Error("daily_cap_level_one must be greater than or equal to 0"),
		),
		validation.Field(&dto.MinAmountVIRTUAL,
			validation.Required.Error("min_amount_virtual is required"),
			validation.Min(0).Error("min_amount_virtual must be greater than or equal to 0"),
		),
		validation.Field(&dto.AboveAmount,
			validation.Required.Error("above_amount is required"),
			validation.Min(uint64(1)).Error("above_amount must be greater than 0"),
		),
		validation.Field(&dto.AboveServiceFee,
			validation.Required.Error("above_service_fee is required"),
			validation.Min(uint64(0)).Error("above_service_fee must be greater than or equal to 0"),
		),

		validation.Field(&dto.Tiers,
			validation.Required.Error("tiers are required"),
			validation.Length(1, 0).Error("at least one tier is required"),
		),
		validation.Field(&dto.CBglEntry, validation.Required.Error("cbgl_entry is required")),
		validation.Field(&dto.IFBglEntry, validation.Required.Error("ifbgl_entry is required")),
	)
	if err != nil {
		return err
	}

	// Validate nested DTOs
	if err := dto.CBglEntry.Validate(); err != nil {
		return fmt.Errorf("cbgl_entry validation failed: %w", err)
	}

	if err := dto.IFBglEntry.Validate(); err != nil {
		return fmt.Errorf("ifbgl_entry validation failed: %w", err)
	}

	// Business logic validation: Single cap cannot exceed daily cap
	if dto.SingleCapLevelOne > dto.DailyCapLevelOne {
		return fmt.Errorf("single_cap_level_one (%d) cannot exceed daily_cap_level_one (%d)", dto.SingleCapLevelOne, dto.DailyCapLevelOne)
	}

	// Business logic validation: Min amount should not be greater than max cap amounts
	if dto.MinAmountVIRTUAL > dto.SingleCapLevelOne {
		return fmt.Errorf("min_amount_virtual (%d) cannot exceed single_cap_level_one (%d)", dto.MinAmountVIRTUAL, dto.SingleCapLevelOne)
	}

	if dto.MinAmountVIRTUAL > dto.DailyCapLevelOne {
		return fmt.Errorf("min_amount_virtual (%d) cannot exceed daily_cap_level_one (%d)", dto.MinAmountVIRTUAL, dto.DailyCapLevelOne)
	}

	// Validate each tier individually
	for i, tier := range dto.Tiers {
		if err := tier.Validate(); err != nil {
			return fmt.Errorf("tier %d validation failed: %w", i+1, err)
		}
	}

	// Basic business logic validation for tier relationships
	for i := 0; i < len(dto.Tiers); i++ {
		tier := dto.Tiers[i]

		// Min must be less than max within each tier
		if tier.Min >= tier.Max {
			return fmt.Errorf("in tier %d, min (%d) must be less than max (%d)", i+1, tier.Min, tier.Max)
		}

		// Ensure tier continuity (max of current tier = min of next tier)
		if i < len(dto.Tiers)-1 {
			nextTier := dto.Tiers[i+1]
			if tier.Max != nextTier.Min {
				return fmt.Errorf("tier %d max (%d) must equal tier %d min (%d)", i+1, tier.Max, i+2, nextTier.Min)
			}
		}
	}

	return nil
}

func (r *SingleMaxTransferRequest) Validate() error {
	converter := local_utils.NewNumericConverter()

	// Convert and validate individual caps
	iSingleCap, err := converter.IsPositiveNumber(r.ISingleCap, "individual_single_cap")
	if err != nil {
		return err
	}

	iDailyCap, err := converter.IsPositiveNumber(r.IDailyCap, "individual_daily_cap")
	if err != nil {
		return err
	}

	// Convert and validate corporate caps
	cSingleCap, err := converter.IsPositiveNumber(r.CSingleCap, "corporate_single_cap")
	if err != nil {
		return err
	}

	cDailyCap, err := converter.IsPositiveNumber(r.CDailyCap, "corporate_daily_cap")
	if err != nil {
		return err
	}

	if iSingleCap > iDailyCap {
		return fmt.Errorf("individual_single_cap (%d) cannot exceed individual_daily_cap (%d)", iSingleCap, iDailyCap)
	}

	if cSingleCap > cDailyCap {
		return fmt.Errorf("corporate_single_cap (%d) cannot exceed corporate_daily_cap (%d)", cSingleCap, cDailyCap)
	}

	return nil
}

func (r *TotalMaxTransferUpdateRequest) Validate() error {
	converter := local_utils.NewNumericConverter()

	_, err := converter.IsPositiveNumber(r.TotalTransferLimit, "total_cap")
	if err != nil {
		return err
	}

	return nil
}

func (r *MinimumTransferUpdateRequest) Validate() error {
	converter := local_utils.NewNumericConverter()

	_, err := converter.IsPositiveNumber(r.Minimum, "min_amount")
	if err != nil {
		return err
	}

	return nil
}
func (r *DeleteServiceFeeTireRequest) Validate() error {
	return validation.ValidateStruct(
		r,
		validation.Field(&r.ID,
			validation.Required.Error("tire_id is required"),
			validation.By(local_utils.NoSpecialChars),
		),
	)
}
