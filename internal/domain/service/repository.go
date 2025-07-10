package service

import "context"

type ServiceRepository interface {
	GetAllServiceDetails(ctx context.Context) ([]*Service, error)
	GetOneServiceDetail(ctx context.Context, id string) (Service, error)
	UpdateOneServiceDetailRequest(ctx context.Context, id string, update Service) error
	UpdateCapMinAmount(ctx context.Context, id string, minAmount uint64) error
	InitiateServiceFeeUpdate(ctx context.Context, cpsAction CPSAction) (UpdateServiceDetailsResponse, error)
	ApproveServiceFeeUpdate(ctx context.Context, cpsAction CPSAction) (UpdateServiceDetailsResponse, error)
	RejectServiceFeeUpdate(ctx context.Context, cpsAction CPSAction) (UpdateServiceDetailsResponse, error)
}
