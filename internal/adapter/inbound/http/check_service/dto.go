package service

import (
	"context"
	"fmt"
	"regexp"
	"strconv"

	local_utils "github.com/CBE-Super-App/cbe-super-app-cps-action/pkgs/utils"
	validation "github.com/go-ozzo/ozzo-validation/v4"
)

type ApplicationService interface {
	GetAllService(ctx context.Context, filterParams local_utils.Filter) (*local_utils.PaginatedResponse[*any], error)
	GetAllMinimumTransferCap(ctx context.Context, filterParams local_utils.Filter) (*local_utils.PaginatedResponse[*any], error)
	GetAllMaximumTransferCap(ctx context.Context, filterParams local_utils.Filter) (*local_utils.PaginatedResponse[*any], error)
	GetAllServiceFee(ctx context.Context, filterParams local_utils.Filter) (*local_utils.PaginatedResponse[*any], error)
	GetAllTotalTransferCap(ctx context.Context, filterParams local_utils.Filter) (*local_utils.PaginatedResponse[*any], error)
	GetServiceFeeDetail(ctx context.Context, id string) (*any, error)
	UpdateServiceFee(ctx context.Context, id string, req any) error
	UpdateSingleMaxTransfer(ctx context.Context, id string, req any) error
	UpdateTotalMaxTransferCap(ctx context.Context, id string, req any) error
	UpdateMinimumTransferCap(ctx context.Context, id string, req any) error
	DeleteServiceFeeTire(ctx context.Context, id string) error
}

type UpdateServiceFeeRequest struct {
	ID        string `json:"tire_id"`
	MinAmount string `json:"min_amount"`
	MaxAmount string `json:"max_amount"`
}

type TierDTO struct {
	Min       string `json:"min"`
	Max       string `json:"max"`
	FeeAmount string `json:"feeAmount"`
}

type CBglEntryDTO struct {
	CBglProductAccount    string `json:"CBglProductAccount"`
	CBglProductBranchcode string `json:"CBglProductBranchcode"`
	CBglServiceAccount    string `json:"CBglServiceAccount"`
	CBglServiceBranchcode string `json:"CBglServiceBranchcode"`
}

type IFBglEntryDTO struct {
	IFBglProductAccount    string `json:"IFBglProductAccount"`
	IFBglProductBranchcode string `json:"IFBglProductBranchcode"`
	IFBglServiceAccount    string `json:"IFBglServiceAccount"`
	IFBglServiceBranchcode string `json:"IFBglServiceBranchcode"`
}

type ServiceFeeDetailDTO struct {
	ServiceType       string        `json:"serviceType"`
	PaymentType       string        `json:"paymentType"`
	SingleCapLevelOne int           `json:"singleCapLevelOne"`
	DailyCapLevelOne  int           `json:"dailyCapLevelOne"`
	MinAmountVIRTUAL  int           `json:"minAmountVIRTUAL"`
	AboveAmount       string        `json:"aboveAmount"`
	AboveServiceFee   string        `json:"aboveServiceFee"`
	Tiers             []TierDTO     `json:"tiers"`
	CBglEntry         CBglEntryDTO  `json:"CBglEntry"`
	IFBglEntry        IFBglEntryDTO `json:"IFBglEntry"`
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

		min, errMin := strconv.ParseFloat(tier.Min, 64)
		max, errMax := strconv.ParseFloat(tier.Max, 64)
		if errMin != nil || errMax != nil {
			return fmt.Errorf("invalid min or max value in tier %d", i+1)
		}
		if min >= max {
			return fmt.Errorf("in tier %d, min must be less than max", i+1)
		}
		if i < len(dto.Tiers)-1 {
			nextTier := dto.Tiers[i+1]
			nextMin, errNextMin := strconv.ParseFloat(nextTier.Min, 64)
			if errNextMin != nil {
				return fmt.Errorf("invalid min value in tier %d", i+2)
			}
			if max != nextMin {
				return fmt.Errorf("tier %d max must equal tier %d min", i+1, i+2)
			}
		}
	}
	return nil
}

type SingleMaxTransferRequest struct {
	ISingleCap string `json:"individual_single_cap"`
	IDailyCap  string `json:"individual_daily_cap"`
	CSingleCap string `json:"corporate_single_cap"`
	CDailyCap  string `json:"corporate_daily_cap"`
}

type TotalMaxTransferUpdateRequest struct {
	TotalTransferLimit string `json:"total_cap"`
}

func (r *TotalMaxTransferUpdateRequest) Validate() error {
	err := validation.ValidateStruct(
		r,
		validation.Field(&r.TotalTransferLimit, validation.Required, validation.Match(regexp.MustCompile(`^\d+$`))),
	)
	if err != nil {
		return err
	}
	return nil
}

type MinimumTransferUpdateRequest struct {
	Minimum string `json:"min_amount"`
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
