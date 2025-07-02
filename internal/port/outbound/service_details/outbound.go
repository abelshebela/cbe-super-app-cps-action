package outbound

import (
	"context"

	"github.com/CBE-Super-App/cbe-super-app-cps-action/internal/domain/service"
	// "github.com/CBE-Super-App/cbe-super-app-cps-action/internal/domain/action"
)

type OutboundServiceDetailInfra interface {
	GetAllServiceDetails(ctx context.Context) ([]service.Service, error)
	GetOneServiceDetail(ctx context.Context, id string) (service.Service, error)
	UpdateOneServiceDetailRequest(ctx context.Context, id string, update any) error
	UpdateOneServiceDetailApprove(ctx context.Context, id string, update any) error

	/* 	InitiateServiceFeeUpdate(ctx context.Context, cpsAction service.CPSAction) (service.UpdateServiceDetailsResponse, error)
	   	ApproveServiceFeeUpdate(ctx context.Context, cpsAction service.CPSAction) (service.UpdateServiceDetailsResponse, error)
	   	RejectServiceFeeUpdate(ctx context.Context, cpsAction service.CPSAction) (service.UpdateServiceDetailsResponse, error) */

	InitiateServiceFeeUpdate(ctx context.Context, cpsAction service.CPSAction) (service.UpdateServiceDetailsResponse, error)
	ApproveServiceFeeUpdate(ctx context.Context, cpsAction service.CPSAction) (service.UpdateServiceDetailsResponse, error)
	RejectServiceFeeUpdate(ctx context.Context, cpsAction service.CPSAction) (service.UpdateServiceDetailsResponse, error)
}

// type ServiceFee interface {
// 	InitiateServiceFeeUpdate(ctx context.Context, cpsAction service.CPSAction) (service.UpdateServiceDetailsResponse, error)
// 	ApproveServiceFeeUpdate(ctx context.Context, cpsAction service.CPSAction) (service.UpdateServiceDetailsResponse, error)
// 	RejectServiceFeeUpdate(ctx context.Context, cpsAction service.CPSAction) (service.UpdateServiceDetailsResponse, error)
// }
