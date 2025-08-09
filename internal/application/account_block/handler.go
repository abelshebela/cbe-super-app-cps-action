package account_block

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"

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
	EnableSingleBranch(ctx context.Context, branch action.Branch, maker action.User) (string, error)

	ApproveSingleBranchDisable(ctx context.Context, actionID string, approve bool, reason *string) error

	FilterMultipleBranches(ctx context.Context, region, district string, filterParams *constant.Filter) (*constant_utils.PaginatedResponse[[]*model.Branch], error)
	DisableMultipleBranches(ctx context.Context, branches []action.Branch, maker action.User) (string, error)
	EnableMultipleBranches(ctx context.Context, branches []action.Branch, maker action.User) (string, error)

	ApproveBulkBranchesDisable(ctx context.Context, actionID string, approve bool, reason *string) error

	GetBranchByCode(ctx context.Context, branchCode string) (action.Branch, error)

	BlockRegion(ctx context.Context, regionCode string, maker action.CPSAction) (string, error)
	// EnabelRegion(ctx context.Context, regionCode string, maker action.CPSAction) (string, error)
	UpdateRegion(ctx context.Context, region action.Region) error
	ApproveRegionBlock(ctx context.Context, actionID string, approve bool, reason *string, checker action.User) error
	GetRegionByCode(ctx context.Context, regionCode string) (action.Region, error)

	BlockDistrict(ctx context.Context, districtCode string, maker action.CPSAction) (string, error)
	// EnableDistrict(ctx context.Context, districtCode string, maker action.CPSAction) (string, error)
	GetDistrictByCode(ctx context.Context, districtCode string) (action.District, error)
	ApproveBlockDistrict(ctx context.Context, actionID string, approve bool, reason *string, checker action.User) error

	BlockCity(ctx context.Context, cityCode string, maker action.CPSAction) (string, error)
	// EnableCity(ctx context.Context, cityCode string, maker action.CPSAction) (string, error)

	GetCityByCode(ctx context.Context, cityCode string) (action.City, error)
	ApproveBlockCity(ctx context.Context, actionID string, approve bool, reason *string, checker action.User) error

	BlockUser(ctx context.Context, userID string, maker action.CPSAction) (string, error)
	GetUserByPhone(ctx context.Context, phoneNumber string, maker action.CPSAction) (member.User, error)
	ApproveBlockUser(ctx context.Context, actionID string, approve bool, reason *string, checker action.User) error

	GetAllCities(ctx context.Context, filter *constant.Filter) (*constant_utils.PaginatedResponse[[]*model.City], error)
	GetAllDistricts(ctx context.Context, filter *constant.Filter) (*constant_utils.PaginatedResponse[[]*model.District], error)
	GetAllRegions(ctx context.Context, filter *constant.Filter) (*constant_utils.PaginatedResponse[[]*model.Region], error)
	GetAllBranches(ctx context.Context, filter *constant.Filter) (*constant_utils.PaginatedResponse[[]*model.Branch], error)

	// Newly added
	// Branch
	EnableBranches(r *http.Request) error
	DisableBranches(r *http.Request) error

	// Region
	EnableRegion(r *http.Request) error
	DisableRegion(r *http.Request) error

	// District
	EnableDistrict(r *http.Request) error
	DisableDistrict(r *http.Request) error

	// City
	EnableCity(r *http.Request) error
	DisableCity(r *http.Request) error
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

func (h *Handler) EnableSingleBranch(ctx context.Context, branch action.Branch, maker action.User) (string, error) {
	return h.service.EnableSingleBranch(ctx, branch, maker)
}

func (h *Handler) ApproveSingleBranchDisable(ctx context.Context, actionID string, approve bool, reason *string) error {
	// return h.service.ApproveSingleBranchDisable(ctx, actionID, approve, reason)
	return nil
}

func (h *Handler) FilterMultipleBranches(ctx context.Context, region, district string, filterParams *constant.Filter) (*constant_utils.PaginatedResponse[[]*model.Branch], error) {
	return h.service.FilterMultipleBranches(ctx, region, district, filterParams)
}
func (h *Handler) EnableMultipleBranches(ctx context.Context, branches []action.Branch, maker action.User) (string, error) {
	return h.service.EnableMultipleBranches(ctx, branches, maker)
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

//	func (h *Handler) EnabelRegion(ctx context.Context, regionCode string, maker action.CPSAction) (string, error) {
//		return h.service.EnableRegion(ctx, regionCode, maker)
//	}
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

//	func (h *Handler) EnableDistrict(ctx context.Context, districtCode string, maker action.CPSAction) (string, error) {
//		return h.service.EnableDistrict(ctx, districtCode, maker)
//	}
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

// func (h *Handler) EnableCity(ctx context.Context, cityCode string, maker action.CPSAction) (string, error) {
// 	return h.service.EnableCity(ctx, cityCode, maker)
// }

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

func (h *Handler) GetAllBranches(ctx context.Context, filter *constant.Filter) (*constant_utils.PaginatedResponse[[]*model.Branch], error) {
	return h.service.GetAllBranches(ctx, filter)
}

// Newly added
// Branch
func (h *Handler) EnableBranches(r *http.Request) error {
	var req EnableOrDisableBranches
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		return err
	}

	// Trim any whitespace
	req.clean()

	if len(req.BranchCodes) == 0 {
		return fmt.Errorf("BRANCH_CODE_IS_REQUIRED")
	}

	ctx := r.Context()
	err := h.service.EnableBranches(ctx, req.BranchCodes)
	if err != nil {
		return err
	}

	return nil
}
func (h *Handler) DisableBranches(r *http.Request) error {
	var req EnableOrDisableBranches
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		return err
	}

	// Trim any whitespace
	req.clean()

	if len(req.BranchCodes) == 0 {
		return fmt.Errorf("BRANCH_CODE_IS_REQUIRED")
	}

	ctx := r.Context()
	err := h.service.DisableBranches(ctx, req.BranchCodes)
	if err != nil {
		return err
	}

	return nil
}

// Region
func (h *Handler) EnableRegion(r *http.Request) error {
	var req EnableOrDisableRegions
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		return err
	}

	// Trim any whitespace
	req.clean()

	if len(req.RegionsCodes) == 0 {
		return fmt.Errorf("REGION_CODE_IS_REQUIRED")
	}

	ctx := r.Context()
	err := h.service.EnableRegion(ctx, req.RegionsCodes)
	if err != nil {
		return err
	}

	return nil
}

func (h *Handler) DisableRegion(r *http.Request) error {
	var req EnableOrDisableRegions
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		return err
	}

	// Trim any whitespace
	req.clean()

	if len(req.RegionsCodes) == 0 {
		return fmt.Errorf("REGION_CODE_IS_REQUIRED")
	}

	ctx := r.Context()
	err := h.service.DisableRegion(ctx, req.RegionsCodes)
	if err != nil {
		return err
	}

	return nil
}

// District
func (h *Handler) EnableDistrict(r *http.Request) error {
	var req EnableOrDisableDistricts
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		return err
	}

	// Trim any whitespace
	req.clean()

	if len(req.DistrictCodes) == 0 {
		return fmt.Errorf("DISTRICT_CODE_IS_REQUIRED")
	}

	ctx := r.Context()
	err := h.service.EnableDistrict(ctx, req.DistrictCodes)
	if err != nil {
		return err
	}

	return nil
}

func (h *Handler) DisableDistrict(r *http.Request) error {
	var req EnableOrDisableDistricts
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		return err
	}

	// Trim any whitespace
	req.clean()

	if len(req.DistrictCodes) == 0 {
		return fmt.Errorf("DISTRICT_CODE_IS_REQUIRED")
	}

	ctx := r.Context()
	err := h.service.DisableDistrict(ctx, req.DistrictCodes)
	if err != nil {
		return err
	}

	return nil
}

// City
func (h *Handler) EnableCity(r *http.Request) error {
	var req EnableOrDisableCities
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		return err
	}

	// Trim any whitespace
	req.clean()

	if len(req.CitiesCode) == 0 {
		return fmt.Errorf("CITY_CODE_IS_REQUIRED")
	}

	ctx := r.Context()
	err := h.service.EnableCity(ctx, req.CitiesCode)
	if err != nil {
		return err
	}

	return nil
}

func (h *Handler) DisableCity(r *http.Request) error {
	var req EnableOrDisableCities
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		return err
	}

	// Trim any whitespace
	req.clean()

	if len(req.CitiesCode) == 0 {
		return fmt.Errorf("CITY_CODE_IS_REQUIRED")
	}

	ctx := r.Context()
	err := h.service.DisableCity(ctx, req.CitiesCode)
	if err != nil {
		return err
	}

	return nil
}
