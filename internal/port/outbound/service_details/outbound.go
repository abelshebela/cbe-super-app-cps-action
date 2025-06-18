package outbound

import (
	"context"

	"gitlab.com/bersufekadgetachew/cbe-super-app-cps-action/internal/domain/service"
	// "gitlab.com/bersufekadgetachew/cbe-super-app-cps-action/internal/domain/action"
)

type OutboundServiceDetailInfra interface {
	GetAllServiceDetails(ctx context.Context) ([]service.Service, error)
	GetOneServiceDetail(ctx context.Context, id string) (service.Service, error)
	UpdateOneServiceDetailRequest(ctx context.Context, id string, update any) error
	UpdateOneServiceDetailApprove(ctx context.Context, id string, update any) error
}
