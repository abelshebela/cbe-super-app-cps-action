package hq

import (
	"context"

	"gitlab.com/bersufekadgetachew/cbe-super-app-cps-action/internal/domain/hq/models"
	"gitlab.com/bersufekadgetachew/cbe-super-app-shared/shared/utils"
)

type Service interface {
	GetHQ(ctx context.Context, id string) (models.HQ, error)
	UpdateBlockTimeRequest(ctx context.Context, request models.UpdateBlockTimeRequest) (string, error)
	UpdateArchiveTimeRequest(ctx context.Context, request models.UpdateArchiveTimeRequest) (string, error)
	UpdateBlockTime(ctx context.Context, request models.ApproveRejectRequest) error
	UpdateArchiveTime(ctx context.Context, request models.ApproveRejectRequest) error
}

type HQHandler struct {
	service Service
	logger  utils.Logger
}

func NewHQHandler(service Service, logger utils.Logger) *HQHandler {
	return &HQHandler{
		service: service,
		logger:  logger,
	}
}

func (h *HQHandler) GetHQ(ctx context.Context, id string) (models.HQ, error) {
	return h.service.GetHQ(ctx, id)
}

func (h *HQHandler) UpdateBlockTimeRequest(ctx context.Context, request models.UpdateBlockTimeRequest) (string, error) {
	return h.service.UpdateBlockTimeRequest(ctx, request)
}

func (h *HQHandler) UpdateArchiveTimeRequest(ctx context.Context, request models.UpdateArchiveTimeRequest) (string, error) {
	return h.service.UpdateArchiveTimeRequest(ctx, request)
}

func (h *HQHandler) UpdateBlockTime(ctx context.Context, request models.ApproveRejectRequest) error {
	return h.service.UpdateBlockTime(ctx, request)
}

func (h *HQHandler) UpdateArchiveTime(ctx context.Context, request models.ApproveRejectRequest) error {
	return h.service.UpdateArchiveTime(ctx, request)
}
