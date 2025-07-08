package account_block

import (
	"context"
	"fmt"
	"time"

	"github.com/CBE-Super-App/cbe-super-app-cps-action/internal/domain/account_block"
	"github.com/CBE-Super-App/cbe-super-app-cps-action/internal/domain/action"

	"gitlab.com/bersufekadgetachew/cbe-super-app-shared/shared/entities/member"
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
	repo account_block.AccountBlockRepo
}

func NewApplicationHandler(repo account_block.AccountBlockRepo) ApplicationService {
	return &Handler{
		repo: repo,
	}
}

func (h *Handler) FilterSingleBranches(ctx context.Context, region, district string) ([]action.Branch, error) {
	return h.repo.FilterSingleBranches(ctx, region, district)
}

func (h *Handler) DisableSingleBranch(ctx context.Context, branch action.Branch, maker action.User) error {
	return h.repo.DisableSingleBranch(ctx, branch, maker)
}

func (h *Handler) ApproveSingleBranchDisable(ctx context.Context, actionID string, approve bool, reason *string) error {
	return h.repo.ApproveSingleBranchDisable(ctx, actionID, approve, reason)
}

func (h *Handler) FilterMultipleBranches(ctx context.Context, region, district string) ([]action.Branch, error) {
	return h.repo.FilterMultipleBranches(ctx, region, district)
}

func (h *Handler) DisableMultipleBranches(ctx context.Context, branches []action.Branch, maker action.User) error {

	return h.repo.DisableMultipleBranches(ctx, branches, maker)
}
func (h *Handler) ApproveBulkBranchesDisable(ctx context.Context, actionID string, approve bool, reason *string) error {
	return h.repo.ApproveBulkBranchesDisable(ctx, actionID, approve, reason)
}

func (h *Handler) GetBranchByCode(ctx context.Context, branchCode string) (action.Branch, error) {
	return h.repo.GetBranchByCode(ctx, branchCode)
}

func (h *Handler) BlockRegion(ctx context.Context, regionCode string, maker action.CPSAction) error {
	region, err := h.repo.GetRegionByCode(ctx, regionCode)
	if err != nil {
		return fmt.Errorf("failed to get region by code: %w", err)
	}

	region.Enabled = false
	region.UpdatedAt = time.Now()
	if err := h.repo.UpdateRegion(ctx, region); err != nil {
		return err
	}
	return h.repo.BlockRegion(ctx, region.RegionCode, maker)
}

func (h *Handler) ApproveRegionBlock(ctx context.Context, actionID string, approve bool, reason *string, checker action.User) error {
	return h.repo.ApproveRegionBlock(ctx, actionID, approve, reason, checker)
}

func (h *Handler) GetRegionByCode(ctx context.Context, regionCode string) (action.Region, error) {
	if regionCode == "" {
		return action.Region{}, fmt.Errorf("regionCode is required")
	}

	region, err := h.repo.GetRegionByCode(ctx, regionCode)
	if err != nil {
		return action.Region{}, fmt.Errorf("failed to get region by code: %w", err)
	}

	return region, nil
}
func (h *Handler) UpdateRegion(ctx context.Context, region action.Region) error {
    if region.RegionCode == "" {
        return fmt.Errorf("regionCode is required")
    }

    region.UpdatedAt = time.Now()
    return h.repo.UpdateRegion(ctx, region)
}
func (h *Handler) BlockDistrict(ctx context.Context, districtID string, maker action.CPSAction) error {
	if districtID == "" {
		return fmt.Errorf("districtID is required")
	}
	return h.repo.BlockDistrict(ctx, districtID, maker)
}

func (h *Handler) GetDistrictByID(ctx context.Context, districtID string) (action.District, error) {
	if districtID == "" {
		return action.District{}, fmt.Errorf("districtID is required")
	}
	return h.repo.GetDistrictByID(ctx, districtID)
}

func (h *Handler) ApproveBlockDistrict(ctx context.Context, districtID string, checker action.CPSAction) error {
	return h.repo.ApproveBlockDistrict(ctx, districtID, checker)
}

func (h *Handler) BlockCity(ctx context.Context, cityID string, maker action.CPSAction) error {
	if cityID == "" {
		return fmt.Errorf("cityID is required")
	}
	return h.repo.BlockCity(ctx, cityID, maker)
}

func (h *Handler) GetCityByID(ctx context.Context, cityID string) (action.City, error) {
	if cityID == "" {
		return action.City{}, fmt.Errorf("cityID is required")
	}
	return h.repo.GetCityByID(ctx, cityID)
}

func (h *Handler) ApproveBlockCity(ctx context.Context, cityID string, checker action.CPSAction) error {
	return h.repo.ApproveBlockCity(ctx, cityID, checker)
}

func (h *Handler) BlockUser(ctx context.Context, userID string, maker action.CPSAction) error {
	if userID == "" {
		return fmt.Errorf("userID is required")
	}
	return h.repo.BlockUser(ctx, userID, maker)
}

func (h *Handler) GetUserByID(ctx context.Context, userID string, maker action.CPSAction) (member.User, error) {
	if userID == "" {
		return member.User{}, fmt.Errorf("userID is required")
	}
	return h.repo.GetUserByID(ctx, userID, maker)
}

func (h *Handler) ApproveBlockUser(ctx context.Context, userID string, checker action.CPSAction) error {
	return h.repo.ApproveBlockUser(ctx, userID, checker)
}
