package service

import (
	"context"
	"fmt"
	"strings"

	"github.com/CBE-Super-App/cbe-super-app-cps-action/pkgs/entities"
	local_utils "github.com/CBE-Super-App/cbe-super-app-cps-action/pkgs/utils"
	validation "github.com/go-ozzo/ozzo-validation/v4"
)

type ApplicationService interface {
	GetAllService(ctx context.Context, filterParams local_utils.Filter) (*local_utils.PaginatedResponse[*any], error)
	GetAllMinimumTransferCap(ctx context.Context, filterParams local_utils.Filter) (*local_utils.PaginatedResponse[*any], error)
	GetAllMaximumTransferCap(ctx context.Context, filterParams local_utils.Filter) (*local_utils.PaginatedResponse[*any], error)
	GetAllServiceFee(ctx context.Context, filterParams local_utils.Filter) (*local_utils.PaginatedResponse[*any], error)
	GetAllTotalTransferCap(ctx context.Context) (*any, error)
	GetServiceFeeDetail(ctx context.Context, id string) (*any, error)
	UpdateServiceFee(ctx context.Context, id string, req any) error
	UpdateSingleMaxTransfer(ctx context.Context, id string, req any) error
	UpdateTotalMaxTransferCap(ctx context.Context, id string, newTotalCap uint64) error
	UpdateMinimumTransferCap(ctx context.Context, id string, req any) error
	DeleteServiceFeeTire(ctx context.Context, id string) error
}

// UpdateServiceFeeRequest represents a request to update service fee tiers
type UpdateServiceFeeRequest struct {
	ID        string `json:"tire_id" bson:"tire_id"`
	MinAmount uint64 `json:"min_amount" bson:"min_amount"`
	MaxAmount uint64 `json:"max_amount" bson:"max_amount"`
}

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

// TierDTO represents a fee tier structure
type TierDTO struct {
	Min       uint64 `json:"min" bson:"min"`
	Max       uint64 `json:"max" bson:"max"`
	FeeAmount uint64 `json:"fee_amount" bson:"fee_amount"`
}

func (t *TierDTO) Validate() error {
	err := validation.ValidateStruct(
		t,
		validation.Field(&t.Min,
			validation.Required.Error("min is required"),
			validation.Min(uint64(1)).Error("min must be greater than 0"),
		),
		validation.Field(&t.Max,
			validation.Required.Error("max is required"),
			validation.Min(uint64(1)).Error("max must be greater than 0"),
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

// CBglEntryDTO represents CBE GL entry details
type CBglEntryDTO struct {
	CBglProductAccount    string `json:"cbgl_product_account" bson:"cbgl_product_account"`
	CBglProductBranchcode string `json:"cbgl_product_branchcode" bson:"cbgl_product_branchcode"`
	CBglServiceAccount    string `json:"cbgl_service_account" bson:"cbgl_service_account"`
	CBglServiceBranchcode string `json:"cbgl_service_branchcode" bson:"cbgl_service_branchcode"`
}

func (c *CBglEntryDTO) Validate() error {
	return validation.ValidateStruct(
		c,
		validation.Field(&c.CBglProductAccount,
			validation.Required.Error("cbgl_product_account is required"),
			validation.By(local_utils.NoSpecialChars),
		),
		validation.Field(&c.CBglProductBranchcode,
			validation.Required.Error("cbgl_product_branchcode is required"),
			validation.By(local_utils.NoSpecialChars),
		),
		validation.Field(&c.CBglServiceAccount,
			validation.Required.Error("cbgl_service_account is required"),
			validation.By(local_utils.NoSpecialChars),
		),
		validation.Field(&c.CBglServiceBranchcode,
			validation.Required.Error("cbgl_service_branchcode is required"),
			validation.By(local_utils.NoSpecialChars),
		),
	)
}

// IFBglEntryDTO represents IFB GL entry details
type IFBglEntryDTO struct {
	IFBglProductAccount    string `json:"ifbgl_product_account" bson:"ifbgl_product_account"`
	IFBglProductBranchcode string `json:"ifbgl_product_branchcode" bson:"ifbgl_product_branchcode"`
	IFBglServiceAccount    string `json:"ifbgl_service_account" bson:"ifbgl_service_account"`
	IFBglServiceBranchcode string `json:"ifbgl_service_branchcode" bson:"ifbgl_service_branchcode"`
}

func (i *IFBglEntryDTO) Validate() error {
	return validation.ValidateStruct(
		i,
		validation.Field(&i.IFBglProductAccount,
			validation.Required.Error("ifbgl_product_account is required"),
			validation.By(local_utils.NoSpecialChars),
		),
		validation.Field(&i.IFBglProductBranchcode,
			validation.Required.Error("ifbgl_product_branchcode is required"),
			validation.By(local_utils.NoSpecialChars),
		),
		validation.Field(&i.IFBglServiceAccount,
			validation.Required.Error("ifbgl_service_account is required"),
			validation.By(local_utils.NoSpecialChars),
		),
		validation.Field(&i.IFBglServiceBranchcode,
			validation.Required.Error("ifbgl_service_branchcode is required"),
			validation.By(local_utils.NoSpecialChars),
		),
	)
}

// ServiceFeeDetailDTO represents the complete service fee structure
type ServiceFeeDetailDTO struct {
	ServiceType       string               `json:"service_type" bson:"service_type"`
	PaymentType       entities.PaymentType `json:"payment_type" bson:"payment_type"`
	SingleCapLevelOne int                  `json:"single_cap_level_one" bson:"single_cap_level_one"`
	DailyCapLevelOne  int                  `json:"daily_cap_level_one" bson:"daily_cap_level_one"`
	MinAmountVIRTUAL  int                  `json:"min_amount_virtual" bson:"min_amount_virtual"`
	AboveAmount       uint64               `json:"above_amount" bson:"above_amount"`
	AboveServiceFee   uint64               `json:"above_service_fee" bson:"above_service_fee"`
	Tiers             []TierDTO            `json:"tiers" bson:"tiers"`
	CBglEntry         CBglEntryDTO         `json:"cbgl_entry" bson:"cbgl_entry"`
	IFBglEntry        IFBglEntryDTO        `json:"ifbgl_entry" bson:"ifbgl_entry"`
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
			validation.By(func(value interface{}) error {
				if pt, ok := value.(entities.PaymentType); ok {
					if !pt.IsValid() {
						return fmt.Errorf("payment_type must be either 'percentage' or 'flat_fee', got: %s", pt)
					}
				} else if str, ok := value.(string); ok {
					pt := entities.PaymentType(str)
					if !pt.IsValid() {
						return fmt.Errorf("payment_type must be either 'percentage' or 'flat_fee', got: %s", str)
					}
				} else {
					return fmt.Errorf("invalid payment_type type")
				}
				return nil
			}),
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

// SingleMaxTransferRequest represents a request to update individual transfer caps
type SingleMaxTransferRequest struct {
	ISingleCap interface{} `json:"individual_single_cap" bson:"individual_single_cap"`
	IDailyCap  interface{} `json:"individual_daily_cap" bson:"individual_daily_cap"`
	CSingleCap interface{} `json:"corporate_single_cap" bson:"corporate_single_cap"`
	CDailyCap  interface{} `json:"corporate_daily_cap" bson:"corporate_daily_cap"`
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

	// Business logic validation: Single caps cannot exceed daily caps
	if iSingleCap > iDailyCap {
		return fmt.Errorf("individual_single_cap (%d) cannot exceed individual_daily_cap (%d)", iSingleCap, iDailyCap)
	}

	if cSingleCap > cDailyCap {
		return fmt.Errorf("corporate_single_cap (%d) cannot exceed corporate_daily_cap (%d)", cSingleCap, cDailyCap)
	}

	return nil
}

// TotalMaxTransferUpdateRequest represents a request to update the total transfer cap
type TotalMaxTransferUpdateRequest struct {
	TotalTransferLimit interface{} `json:"total_cap" bson:"total_cap"`
}

func (r *TotalMaxTransferUpdateRequest) Validate() error {
	converter := local_utils.NewNumericConverter()

	_, err := converter.IsPositiveNumber(r.TotalTransferLimit, "total_cap")
	if err != nil {
		return err
	}

	return nil
}

// MinimumTransferUpdateRequest represents a request to update the minimum transfer amount
type MinimumTransferUpdateRequest struct {
	Minimum interface{} `json:"min_amount" bson:"min_amount"`
}

func (r *MinimumTransferUpdateRequest) Validate() error {
	converter := local_utils.NewNumericConverter()

	_, err := converter.IsPositiveNumber(r.Minimum, "min_amount")
	if err != nil {
		return err
	}

	return nil
}

// DeleteServiceFeeTireRequest represents a request to delete a service fee tier
type DeleteServiceFeeTireRequest struct {
	ID string `json:"tire_id" bson:"tire_id"`
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
