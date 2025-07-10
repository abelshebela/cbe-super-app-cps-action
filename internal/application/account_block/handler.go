package account_block

import (
    "context"

    "github.com/CBE-Super-App/cbe-super-app-cps-action/internal/domain/action"
    "gitlab.com/bersufekadgetachew/cbe-super-app-shared/shared/entities/member"
	    "github.com/CBE-Super-App/cbe-super-app-cps-action/internal/domain/account_block"

)

type ApplicationService interface {
    FilterSingleBranches(ctx context.Context, region, district string) ([]action.Branch, error)
    DisableSingleBranch(ctx context.Context, branch action.Branch, maker action.User) error
    ApproveSingleBranchDisable(ctx context.Context, actionID string, approve bool, reason *string) error

    FilterMultipleBranches(ctx context.Context, region, district string) ([]action.Branch, error)
    DisableMultipleBranches(ctx context.Context, branches []action.Branch, maker action.User) error
    ApproveBulkBranchesDisable(ctx context.Context, actionID string, approve bool, reason *string) error

    GetBranchByCode(ctx context.Context, branchCode string) (action.Branch, error)

    BlockRegion(ctx context.Context, regionCode string, maker action.CPSAction) error
    UpdateRegion(ctx context.Context, region action.Region) error
    ApproveRegionBlock(ctx context.Context, actionID string, approve bool, reason *string, checker action.User) error
    GetRegionByCode(ctx context.Context, regionCode string) (action.Region, error)

    BlockDistrict(ctx context.Context, districtID string, maker action.CPSAction) error
    GetDistrictByID(ctx context.Context, districtID string) (action.District, error)
    ApproveBlockDistrict(ctx context.Context, districtID string, checker action.CPSAction) error

    BlockCity(ctx context.Context, cityID string, maker action.CPSAction) error
    GetCityByID(ctx context.Context, cityID string) (action.City, error)
    ApproveBlockCity(ctx context.Context, cityID string, checker action.CPSAction) error

    BlockUser(ctx context.Context, userID string, maker action.CPSAction) error
    GetUserByID(ctx context.Context, userID string, maker action.CPSAction) (member.User, error)
    ApproveBlockUser(ctx context.Context, userID string, checker action.CPSAction) error
}

type Handler struct {
    service account_block.ApplicationServices
}

func NewApplicationHandler(service account_block.ApplicationServices) ApplicationService {
    return &Handler{
        service: service,
    }
}


func (h *Handler) FilterSingleBranches(ctx context.Context, region, district string) ([]action.Branch, error) {
    return h.service.FilterSingleBranches(ctx, region, district)
}

func (h *Handler) DisableSingleBranch(ctx context.Context, branch action.Branch, maker action.User) error {
    return h.service.DisableSingleBranch(ctx, branch, maker)
}

func (h *Handler) ApproveSingleBranchDisable(ctx context.Context, actionID string, approve bool, reason *string) error {
    return h.service.ApproveSingleBranchDisable(ctx, actionID, approve, reason)
}

func (h *Handler) FilterMultipleBranches(ctx context.Context, region, district string) ([]action.Branch, error) {
    return h.service.FilterMultipleBranches(ctx, region, district)
}

func (h *Handler) DisableMultipleBranches(ctx context.Context, branches []action.Branch, maker action.User) error {
    return h.service.DisableMultipleBranches(ctx, branches, maker)
}

func (h *Handler) ApproveBulkBranchesDisable(ctx context.Context, actionID string, approve bool, reason *string) error {
    return h.service.ApproveBulkBranchesDisable(ctx, actionID, approve, reason)
}

func (h *Handler) GetBranchByCode(ctx context.Context, branchCode string) (action.Branch, error) {
    return h.service.GetBranchByCode(ctx, branchCode)
}

func (h *Handler) BlockRegion(ctx context.Context, regionCode string, maker action.CPSAction) error {
    return h.service.BlockRegion(ctx, regionCode, maker)
}

func (h *Handler) UpdateRegion(ctx context.Context, region action.Region) error {
    return h.service.UpdateRegion(ctx, region)
}

func (h *Handler) ApproveRegionBlock(ctx context.Context, actionID string, approve bool, reason *string, checker action.User) error {
    return h.service.ApproveRegionBlock(ctx, actionID, approve, reason, checker)
}

func (h *Handler) GetRegionByCode(ctx context.Context, regionCode string) (action.Region, error) {
    return h.service.GetRegionByCode(ctx, regionCode)
}

func (h *Handler) BlockDistrict(ctx context.Context, districtID string, maker action.CPSAction) error {
    return h.service.BlockDistrict(ctx, districtID, maker)
}

func (h *Handler) GetDistrictByID(ctx context.Context, districtID string) (action.District, error) {
    return h.service.GetDistrictByID(ctx, districtID)
}

func (h *Handler) ApproveBlockDistrict(ctx context.Context, districtID string, checker action.CPSAction) error {
    return h.service.ApproveBlockDistrict(ctx, districtID, checker)
}

func (h *Handler) BlockCity(ctx context.Context, cityID string, maker action.CPSAction) error {
    return h.service.BlockCity(ctx, cityID, maker)
}

func (h *Handler) GetCityByID(ctx context.Context, cityID string) (action.City, error) {
    return h.service.GetCityByID(ctx, cityID)
}

func (h *Handler) ApproveBlockCity(ctx context.Context, cityID string, checker action.CPSAction) error {
    return h.service.ApproveBlockCity(ctx, cityID, checker)
}

func (h *Handler) BlockUser(ctx context.Context, userID string, maker action.CPSAction) error {
    return h.service.BlockUser(ctx, userID, maker)
}

func (h *Handler) GetUserByID(ctx context.Context, userID string, maker action.CPSAction) (member.User, error) {
    return h.service.GetUserByID(ctx, userID, maker)
}

func (h *Handler) ApproveBlockUser(ctx context.Context, userID string, checker action.CPSAction) error {
    return h.service.ApproveBlockUser(ctx, userID, checker)
}