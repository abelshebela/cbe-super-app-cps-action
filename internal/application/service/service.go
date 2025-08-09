package service

import (
	"context"

	local_utils "github.com/CBE-Super-App/cbe-super-app-cps-action/pkgs/utils"
	"gitlab.com/bersufekadgetachew/cbe-super-app-shared/shared/utils"
)

type serviceApplication struct {
	serviceDomain ApplicationService
	logger        utils.Logger
}

func NewServiceApplication(domain ApplicationService, logger utils.Logger) ApplicationService {
	return &serviceApplication{
		serviceDomain: domain,
		logger:        logger,
	}
}

func (sa *serviceApplication) GetAllService(ctx context.Context, filterParams local_utils.Filter) (*local_utils.PaginatedResponse[*any], error) {
	data, err := sa.serviceDomain.GetAllService(ctx, filterParams)
	if err != nil {
		return nil, err
	}

	return data, nil
}
func (sa *serviceApplication) GetAllMinimumTransferCap(ctx context.Context, filterParams local_utils.Filter) (*local_utils.PaginatedResponse[*any], error) {
	data, err := sa.serviceDomain.GetAllMinimumTransferCap(ctx, filterParams)
	if err != nil {
		return nil, err
	}

	return data, nil
}
func (sa *serviceApplication) GetAllMaximumTransferCap(ctx context.Context, filterParams local_utils.Filter) (*local_utils.PaginatedResponse[*any], error) {
	data, err := sa.serviceDomain.GetAllMaximumTransferCap(ctx, filterParams)
	if err != nil {
		return nil, err
	}

	return data, nil
}
func (sa *serviceApplication) GetAllServiceFee(ctx context.Context, filterParams local_utils.Filter) (*local_utils.PaginatedResponse[*any], error) {
	data, err := sa.serviceDomain.GetAllServiceFee(ctx, filterParams)
	if err != nil {
		return nil, err
	}

	return data, nil
}
func (sa *serviceApplication) GetAllTotalTransferCap(ctx context.Context) (*any, error) {
	data, err := sa.serviceDomain.GetAllTotalTransferCap(ctx)
	if err != nil {
		return nil, err
	}

	return data, nil
}
func (sa *serviceApplication) GetServiceFeeDetail(ctx context.Context, id string) (*any, error) {
	data, err := sa.serviceDomain.GetServiceFeeDetail(ctx, id)
	if err != nil {
		return nil, err
	}

	return data, nil
}
func (sa *serviceApplication) UpdateServiceFee(ctx context.Context, id string, req any) error {

	err := sa.serviceDomain.UpdateServiceFee(ctx, id, req)
	if err != nil {
		return err
	}
	return nil
}
func (sa *serviceApplication) UpdateSingleMaxTransfer(ctx context.Context, id string, req any) error {
	err := sa.serviceDomain.UpdateSingleMaxTransfer(ctx, id, req)
	if err != nil {
		return err
	}
	return nil
}
func (sa *serviceApplication) UpdateTotalMaxTransferCap(ctx context.Context, id string, req any) error {
	err := sa.serviceDomain.UpdateTotalMaxTransferCap(ctx, id, req)
	if err != nil {
		return err
	}
	return nil
}
func (sa *serviceApplication) UpdateMinimumTransferCap(ctx context.Context, id string, req any) error {
	err := sa.serviceDomain.UpdateMinimumTransferCap(ctx, id, req)
	if err != nil {
		return err
	}
	return nil
}

func (sa *serviceApplication) DeleteServiceFeeTire(ctx context.Context, id string) error {
	err := sa.serviceDomain.DeleteServiceFeeTire(ctx, id)
	if err != nil {
		return err
	}
	return nil
}
