package updatedbulkservice

import (
	"context"

	domain "github.com/CBE-Super-App/cbe-super-app-cps-action/internal/domain/updated_bulk_service"
	common_util "github.com/CBE-Super-App/cbe-super-app-cps-action/pkgs/utils"
	constant "github.com/CBE-Super-App/cbe-super-app-cps-action/utils"
	"gitlab.com/bersufekadgetachew/cbe-super-app-shared/shared/utils"
)

type Application interface {
	GetAllBulkServices(ctx context.Context, filterParams *constant.Filter) (*common_util.PaginatedResponse[[]*domain.ServiceDetails], error)
}

type BulkServiceApplication struct {
	service domain.BulkService
	logger  utils.Logger
}

func NewApplicationHandler(service domain.BulkService, logger utils.Logger) *BulkServiceApplication {
	return &BulkServiceApplication{service: service, logger: logger}
}

func (h *BulkServiceApplication) GetAllBulkServices(ctx context.Context, filterParams *constant.Filter) (*common_util.PaginatedResponse[[]*domain.ServiceDetails], error) {
	return h.service.GetAllBulkServices(ctx, filterParams)
}
