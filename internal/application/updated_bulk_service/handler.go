package updatedbulkservice

import (
	"context"
	"fmt"
	"net/http"
	"strings"

	"github.com/go-chi/chi/v5"

	domain "github.com/CBE-Super-App/cbe-super-app-cps-action/internal/domain/updated_bulk_service"
	common_util "github.com/CBE-Super-App/cbe-super-app-cps-action/pkgs/utils"
	constant "github.com/CBE-Super-App/cbe-super-app-cps-action/utils"
	"gitlab.com/bersufekadgetachew/cbe-super-app-shared/shared/utils"
)

type Application interface {
	GetAllBulkServices(ctx context.Context, filterParams *constant.Filter) (*common_util.PaginatedResponse[[]*domain.ServiceDetails], error)
	EnableBulkService(ctx context.Context, r *http.Request) error
	DisableBulkService(ctx context.Context, r *http.Request) error
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

func (h *BulkServiceApplication) EnableBulkService(ctx context.Context, r *http.Request) error {
	service_code := strings.TrimSpace(chi.URLParam(r, "service_code"))

	if service_code == "" {
		h.logger.Errorf("bulk's service_code is required")
		return fmt.Errorf("BULK_SERVICE_CODE_IS_REQUIRED")
	}

	err := h.service.EnableBulkService(ctx, service_code)
	if err != nil {
		return err
	}

	return nil
}

func (h *BulkServiceApplication) DisableBulkService(ctx context.Context, r *http.Request) error {
	service_id := strings.TrimSpace(chi.URLParam(r, "service_code"))

	if service_id == "" {
		h.logger.Errorf("bulk's service_code is required")
		return fmt.Errorf("BULK_SERVICE_CODE_IS_REQUIRED")
	}

	err := h.service.DisableBulkService(ctx, service_id)
	if err != nil {
		return err
	}

	return nil
}
