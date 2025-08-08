package service

import (
	"context"
	"fmt"
	"regexp"

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
	UpdateTotalMaxTransferCap(ctx context.Context, id string, req any) error
	UpdateMinimumTransferCap(ctx context.Context, id string, req any) error
	DeleteServiceFeeTire(ctx context.Context, id string) error
}

type UpdateServiceFeeRequest struct {
	ID        string `json:"tire_id"`
	MinAmount uint64 `json:"min_amount"`
	MaxAmount uint64 `json:"max_amount"`
}

// StringToUint64 converts MinAmount and MaxAmount from string to uint64

type TierDTO struct {
	Min       uint64 `json:"min"`
	Max       uint64 `json:"max"`
	FeeAmount uint64 `json:"fee_amount"`
}

type CBglEntryDTO struct {
	CBglProductAccount    string `json:"cbgl_product_account"`
	CBglProductBranchcode string `json:"cbgl_product_branchcode"`
	CBglServiceAccount    string `json:"cbgl_service_account"`
	CBglServiceBranchcode string `json:"cbgl_service_branchcode"`
}

type IFBglEntryDTO struct {
	IFBglProductAccount    string `json:"ifbgl_product_account"`
	IFBglProductBranchcode string `json:"ifbgl_product_branchcode"`
	IFBglServiceAccount    string `json:"ifbgl_service_account"`
	IFBglServiceBranchcode string `json:"ifbgl_service_branchcode"`
}

type ServiceFeeDetailDTO struct {
	ServiceType       string        `json:"service_type"`
	PaymentType       string        `json:"payment_type"`
	SingleCapLevelOne int           `json:"single_cap_level_one"`
	DailyCapLevelOne  int           `json:"daily_cap_level_one"`
	MinAmountVIRTUAL  int           `json:"min_amount_virtual"`
	AboveAmount       uint64        `json:"above_amount"`
	AboveServiceFee   uint64        `json:"above_service_fee"`
	Tiers             []TierDTO     `json:"tiers"`
	CBglEntry         CBglEntryDTO  `json:"cbgl_entry"`
	IFBglEntry        IFBglEntryDTO `json:"ifbgl_entry"`
}

func (dto *ServiceFeeDetailDTO) Validate() error {
	// Basic validation
	err := validation.ValidateStruct(
		dto,
		validation.Field(&dto.ServiceType, validation.Required),
		validation.Field(&dto.PaymentType, validation.Required),
		validation.Field(&dto.SingleCapLevelOne, validation.Min(0)),
		validation.Field(&dto.DailyCapLevelOne, validation.Min(0)),
		validation.Field(&dto.MinAmountVIRTUAL, validation.Min(0)),
		validation.Field(&dto.AboveAmount, validation.Required),
		validation.Field(&dto.AboveServiceFee, validation.Required),
		validation.Field(&dto.Tiers, validation.Required, validation.Length(1, 0)),
	)
	if err != nil {
		return err
	}

	for i := 0; i < len(dto.Tiers); i++ {
		tier := dto.Tiers[i]

		min := tier.Min
		max := tier.Max

		if min >= max {
			return fmt.Errorf("in tier %d, min must be less than max", i+1)
		}
		if i < len(dto.Tiers)-1 {
			nextTier := dto.Tiers[i+1]
			nextMin := nextTier.Min

			if max != nextMin {
				return fmt.Errorf("tier %d max must equal tier %d min", i+1, i+2)
			}
		}
	}
	return nil
}

type SingleMaxTransferRequest struct {
	ISingleCap uint64 `json:"individual_single_cap"`
	IDailyCap  uint64 `json:"individual_daily_cap"`
	CSingleCap uint64 `json:"corporate_single_cap"`
	CDailyCap  uint64 `json:"corporate_daily_cap"`
}

type TotalMaxTransferUpdateRequest struct {
	TotalTransferLimit uint `json:"total_cap"`
}

func (r *TotalMaxTransferUpdateRequest) Validate() error {
	err := validation.ValidateStruct(
		r,
		validation.Field(&r.TotalTransferLimit, validation.Required),
	)
	if err != nil {
		return err
	}
	return nil
}

type MinimumTransferUpdateRequest struct {
	Minimum uint64 `json:"min_amount"`
}

func (r *MinimumTransferUpdateRequest) Validate() error {
	return validation.ValidateStruct(
		r,
		validation.Field(&r.Minimum, validation.Required, validation.Match(regexp.MustCompile(`^\d+(\.\d+)?$`))),
	)
}

type DeleteServiceFeeTireRequest struct {
	ID string `json:"tire_id"`
}
