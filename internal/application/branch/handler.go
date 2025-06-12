package branch

import (
    "context"
    "gitlab.com/bersufekadgetachew/cbe-super-app-cps-ms/internal/domain/bulkcustomer/services"
    "gitlab.com/bersufekadgetachew/cbe-super-app-cps-ms/internal/domain/bulkcustomer/entities"
)

type ApplicationService interface {
    FilterSingleBranches(ctx context.Context, region, district string) ([]string, error)
    DisableSingleBranch(ctx context.Context, branchCode, cpsData string) (*entities.CPSAction, error)
    ApproveSingleBranchDisable(ctx context.Context, actionID string, approve bool, reason *string) error

    FilterMultipleBranches(ctx context.Context, region, district string) ([]string, error)
    DisableMultipleBranches(ctx context.Context, branchCodes []string, cpsData string) (*entities.CPSAction, error)
    ApproveBulkBranchesDisable(ctx context.Context, actionID string, approve bool, reason *string) error
}

type Handler struct {
    service services.BranchServices
}

func NewApplicationHandler(service services.BranchServices) ApplicationService {
    return &Handler{
        service: service,
    }
}

func (h *Handler) FilterSingleBranches(ctx context.Context, region, district string) ([]string, error) {
    return h.service.FilterSingleBranches(ctx, region, district)
}

func (h *Handler) DisableSingleBranch(ctx context.Context, branchCode, cpsData string) (*entities.CPSAction, error) {
    return h.service.DisableSingleBranch(ctx, branchCode, cpsData)
}

func (h *Handler) ApproveSingleBranchDisable(ctx context.Context, actionID string, approve bool, reason *string) error {
    return h.service.ApproveSingleBranchDisable(ctx, actionID, approve, reason)
}

func (h *Handler) FilterMultipleBranches(ctx context.Context, region, district string) ([]string, error) {
    return h.service.FilterMultipleBranches(ctx, region, district)
}

func (h *Handler) DisableMultipleBranches(ctx context.Context, branchCodes []string, cpsData string) (*entities.CPSAction, error) {
    return h.service.DisableMultipleBranches(ctx, branchCodes, cpsData)
}

func (h *Handler) ApproveBulkBranchesDisable(ctx context.Context, actionID string, approve bool, reason *string) error {
    return h.service.ApproveBulkBranchesDisable(ctx, actionID, approve, reason)
}