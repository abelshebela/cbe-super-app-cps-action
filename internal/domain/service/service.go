package service

import (
	"context"

	local_utils "github.com/CBE-Super-App/cbe-super-app-cps-action/pkgs/utils"

	"gitlab.com/bersufekadgetachew/cbe-super-app-shared/shared/utils"
)

type serviceDomain struct {
	serviceRepo ServiceRepo
	logger      utils.Logger
}

func NewServiceDomain(repo ServiceRepo, logger utils.Logger) ServiceRepo {
	return &serviceDomain{
		serviceRepo: repo,
		logger:      logger,
	}

}

func (sd *serviceDomain) GetAllService(ctx context.Context, filterParams local_utils.Filter) (*local_utils.PaginatedResponse[*any], error) {
	return sd.serviceRepo.GetAllService(ctx, filterParams)
}
func (sd *serviceDomain) GetAllMinimumTransferCap(ctx context.Context, filterParams local_utils.Filter) (*local_utils.PaginatedResponse[*any], error) {
	return sd.serviceRepo.GetAllMinimumTransferCap(ctx, filterParams)
}
func (sd *serviceDomain) GetAllMaximumTransferCap(ctx context.Context, filterParams local_utils.Filter) (*local_utils.PaginatedResponse[*any], error) {
	return sd.serviceRepo.GetAllMaximumTransferCap(ctx, filterParams)
}
func (sd *serviceDomain) GetAllServiceFee(ctx context.Context, filterParams local_utils.Filter) (*local_utils.PaginatedResponse[*any], error) {
	return sd.serviceRepo.GetAllServiceFee(ctx, filterParams)
}
func (sd *serviceDomain) GetAllTotalTransferCap(ctx context.Context) (*any, error) {
	return sd.serviceRepo.GetAllTotalTransferCap(ctx)
}
func (sd *serviceDomain) GetServiceFeeDetail(ctx context.Context, id string) (*any, error) {
	return sd.serviceRepo.GetServiceFeeDetail(ctx, id)
}
func (sd *serviceDomain) UpdateServiceFee(ctx context.Context, id string, req any) error {
	return sd.serviceRepo.UpdateServiceFee(ctx, id, req)
}
func (sd *serviceDomain) UpdateSingleMaxTransfer(ctx context.Context, id string, req any) error {
	return sd.serviceRepo.UpdateSingleMaxTransfer(ctx, id, req)
}
func (sd *serviceDomain) UpdateTotalMaxTransferCap(ctx context.Context, id string, newTotalCap uint64) error {
	return sd.serviceRepo.UpdateTotalMaxTransferCap(ctx, id, newTotalCap)
}
func (sd *serviceDomain) UpdateMinimumTransferCap(ctx context.Context, id string, req any) error {
	return sd.serviceRepo.UpdateMinimumTransferCap(ctx, id, req)
}
func (sd *serviceDomain) DeleteServiceFeeTire(ctx context.Context, id string) error {
	return sd.serviceRepo.DeleteServiceFeeTire(ctx, id)
}
func (sd *serviceDomain) Authorize(ctx context.Context, cpsAction any) (any, error) {
	return sd.serviceRepo.Authorize(ctx, cpsAction)
}
