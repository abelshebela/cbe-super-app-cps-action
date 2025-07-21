package account_block

import (
	"context"

	"github.com/CBE-Super-App/cbe-super-app-cps-action/internal/adapter/outbound/model"
	"github.com/CBE-Super-App/cbe-super-app-cps-action/internal/domain/account_block"
	"github.com/CBE-Super-App/cbe-super-app-cps-action/internal/domain/action"
	constant_utils "github.com/CBE-Super-App/cbe-super-app-cps-action/pkgs/utils"
	constant "github.com/CBE-Super-App/cbe-super-app-cps-action/utils"
	"gitlab.com/bersufekadgetachew/cbe-super-app-shared/shared/entities/member"
)

type ApplicationService interface {
	FilterSingleBranches(ctx context.Context, region, district string, filterParams *constant.Filter) (*constant_utils.PaginatedResponse[[]*model.Branch], error)
	DisableSingleBranch(ctx context.Context, branch action.Branch, maker action.User) (string, error)
	ApproveSingleBranchDisable(ctx context.Context, actionID string, approve bool, reason *string) error

	FilterMultipleBranches(ctx context.Context, region, district string, filterParams *constant.Filter) (*constant_utils.PaginatedResponse[[]*model.Branch], error)
	DisableMultipleBranches(ctx context.Context, branches []action.Branch, maker action.User) (string, error)
	ApproveBulkBranchesDisable(ctx context.Context, actionID string, approve bool, reason *string) error

	GetBranchByCode(ctx context.Context, branchCode string) (action.Branch, error)

	BlockRegion(ctx context.Context, regionCode string, maker action.CPSAction) (string, error)
	UpdateRegion(ctx context.Context, region action.Region) error
	ApproveRegionBlock(ctx context.Context, actionID string, approve bool, reason *string, checker action.User) error
	GetRegionByCode(ctx context.Context, regionCode string) (action.Region, error)

	BlockDistrict(ctx context.Context, districtCode string, maker action.CPSAction) (string, error)
	GetDistrictByCode(ctx context.Context, districtCode string) (action.District, error)
	ApproveBlockDistrict(ctx context.Context, actionID string, approve bool, reason *string, checker action.User) error

	BlockCity(ctx context.Context, cityCode string, maker action.CPSAction) (string, error)
	GetCityByCode(ctx context.Context, cityCode string) (action.City, error)
	ApproveBlockCity(ctx context.Context, actionID string, approve bool, reason *string, checker action.User) error

	BlockUser(ctx context.Context, userID string, maker action.CPSAction) (string, error)
	GetUserByPhone(ctx context.Context, phoneNumber string, maker action.CPSAction) (member.User, error)
	ApproveBlockUser(ctx context.Context, actionID string, approve bool, reason *string, checker action.User) error

	GetAllCities(ctx context.Context, filter *constant.Filter) (*constant_utils.PaginatedResponse[[]*model.City], error)
	GetAllDistricts(ctx context.Context, filter *constant.Filter) (*constant_utils.PaginatedResponse[[]*model.District], error)
	GetAllRegions(ctx context.Context, filter *constant.Filter) (*constant_utils.PaginatedResponse[[]*model.Region], error)
}

type Handler struct {
	service account_block.ApplicationServices
}

func NewApplicationHandler(service account_block.ApplicationServices) ApplicationService {
	return &Handler{
		service: service,
	}
}

func (h *Handler) FilterSingleBranches(ctx context.Context, region, district string, filterParams *constant.Filter) (*constant_utils.PaginatedResponse[[]*model.Branch], error) {
	return h.service.FilterSingleBranches(ctx, region, district, filterParams)
}

func (h *Handler) DisableSingleBranch(ctx context.Context, branch action.Branch, maker action.User) (string, error) {
	return h.service.DisableSingleBranch(ctx, branch, maker)
}

func (h *Handler) ApproveSingleBranchDisable(ctx context.Context, actionID string, approve bool, reason *string) error {
	// return h.service.ApproveSingleBranchDisable(ctx, actionID, approve, reason)
	return nil
}

func (h *Handler) FilterMultipleBranches(ctx context.Context, region, district string, filterParams *constant.Filter) (*constant_utils.PaginatedResponse[[]*model.Branch], error) {
	return h.service.FilterMultipleBranches(ctx, region, district, filterParams)
}

func (h *Handler) DisableMultipleBranches(ctx context.Context, branches []action.Branch, maker action.User) (string, error) {
	return h.service.DisableMultipleBranches(ctx, branches, maker)
}

func (h *Handler) ApproveBulkBranchesDisable(ctx context.Context, actionID string, approve bool, reason *string) error {
	// return h.service.ApproveBulkBranchesDisable(ctx, actionID, approve, reason)
	return nil
}

func (h *Handler) GetBranchByCode(ctx context.Context, branchCode string) (action.Branch, error) {
	return h.service.GetBranchByCode(ctx, branchCode)
}

func (h *Handler) BlockRegion(ctx context.Context, regionCode string, maker action.CPSAction) (string, error) {
	return h.service.BlockRegion(ctx, regionCode, maker)
}
func (h *Handler) UpdateRegion(ctx context.Context, region action.Region) error {
	return h.service.UpdateRegion(ctx, region)
}

func (h *Handler) ApproveRegionBlock(ctx context.Context, actionID string, approve bool, reason *string, checker action.User) error {
	// return h.service.ApproveRegionBlock(ctx, actionID, approve, reason, checker)
	return nil
}

func (h *Handler) GetRegionByCode(ctx context.Context, regionCode string) (action.Region, error) {
	return h.service.GetRegionByCode(ctx, regionCode)
}

func (h *Handler) BlockDistrict(ctx context.Context, districtCode string, maker action.CPSAction) (string, error) {
	return h.service.BlockDistrict(ctx, districtCode, maker)
}
func (h *Handler) GetDistrictByCode(ctx context.Context, districtCode string) (action.District, error) {
	return h.service.GetDistrictByCode(ctx, districtCode)
}

func (h *Handler) ApproveBlockDistrict(ctx context.Context, actionID string, approve bool, reason *string, checker action.User) error {
	// return h.service.ApproveBlockDistrict(ctx, actionID, approve, reason, checker)
	return nil
}
func (h *Handler) BlockCity(ctx context.Context, cityCode string, maker action.CPSAction) (string, error) {
	return h.service.BlockCity(ctx, cityCode, maker)
}
func (h *Handler) GetCityByCode(ctx context.Context, cityCode string) (action.City, error) {
	return h.service.GetCityByCode(ctx, cityCode)
}

func (h *Handler) ApproveBlockCity(ctx context.Context, actionID string, approve bool, reason *string, checker action.User) error {
	// return h.service.ApproveBlockCity(ctx, actionID, approve, reason, checker)
	return nil
}

func (h *Handler) BlockUser(ctx context.Context, userID string, maker action.CPSAction) (string, error) {
	return h.service.BlockUser(ctx, userID, maker)
}

func (h *Handler) GetUserByPhone(ctx context.Context, phoneNumber string, maker action.CPSAction) (member.User, error) {
	return h.service.GetUserByPhone(ctx, phoneNumber, maker)
}
func (h *Handler) ApproveBlockUser(ctx context.Context, actionID string, approve bool, reason *string, checker action.User) error {
	// return h.service.ApproveBlockUser(ctx, actionID, approve, reason, checker)
	return nil
}

func (h *Handler) GetAllCities(ctx context.Context, filter *constant.Filter) (*constant_utils.PaginatedResponse[[]*model.City], error) {
	return h.service.GetAllCities(ctx, filter)
}

func (h *Handler) GetAllDistricts(ctx context.Context, filter *constant.Filter) (*constant_utils.PaginatedResponse[[]*model.District], error) {
	return h.service.GetAllDistricts(ctx, filter)
}

func (h *Handler) GetAllRegions(ctx context.Context, filter *constant.Filter) (*constant_utils.PaginatedResponse[[]*model.Region], error) {
	return h.service.GetAllRegions(ctx, filter)
}
