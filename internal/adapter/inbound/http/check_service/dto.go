package service

import (
	"context"

	local_utils "github.com/CBE-Super-App/cbe-super-app-cps-action/pkgs/utils"
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
	DeleteMinimumTransferCap(ctx context.Context, id string) error
	DeleteServiceFeeTire(ctx context.Context, id string) error
}

type UpdateServiceFeeRequest struct {
	ID        string `json:"tire_id"`
	MinAmount string `json:"min_amount"`
	MaxAmount string `json:"max_amount"`
}

type SingleMaxTransferRequest struct {
	ISingleCap string `json:"individual_single_cap"`
	IDailyCap  string `json:"individual_daily_cap"`
	CSingleCap string `json:"corporate_single_cap"`
	CDailyCap  string `json:"corporate_daily_cap"`
}

type TotalMaxTransferUpdateRequest struct {
	TotalTransferLimit string `json:"total_transfer_limit"`
}

type MinimumTransferUpdateRequest struct {
	Minimum string `json:"min_amount"`
}

type DeleteServiceFeeTireRequest struct {
	ID string `json:"tire_id"`
}
