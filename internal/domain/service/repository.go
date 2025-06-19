package service

import "context"

type Repository interface {
	GetAllServiceDetails(ctx context.Context) ([]*Service, error)
	GetOneServiceDetail(ctx context.Context, id string) (Service, error)
	//UpdateOneServiceDetail(ctx context.Context, id string, update Service) error
	UpdateOneServiceDetailRequest(ctx context.Context, id string, update Service) error
	//UpdateOneSeviceDeatil(ctx context.Context, pd Service) error

	UpdateCapMinAmount(ctx context.Context, id string, minAmount uint64) error
	ApproveServiceDetails(ctx context.Context, actionID string, approve bool, checkerID string, rejectionReason string) error

	InitiateServiceFeeUpdate(ctx context.Context, cpsAction CPSAction) (UpdateServiceDetailsResponse, error)
	ApproveServiceFeeUpdate(ctx context.Context, cpsAction CPSAction) (UpdateServiceDetailsResponse, error)
	RejectServiceFeeUpdate(ctx context.Context, cpsAction CPSAction) (UpdateServiceDetailsResponse, error)
}
