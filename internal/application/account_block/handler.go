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
)

type ApplicationService interface {
	GetBranch(ctx context.Context, branchCode string, filterParams *constant.Filter) (*model.Branch, error)
	GetAllBranches(ctx context.Context, region, district string, filterParams *constant.Filter) (*constant_utils.PaginatedResponse[[]*model.Branch], error)
	GetRegionByCode(ctx context.Context, regionCode string) (action.Region, error)
	GetAllRegions(ctx context.Context, filter *constant.Filter) (*constant_utils.PaginatedResponse[[]*model.Region], error)
	GetDistrictByCode(ctx context.Context, districtCode string) (action.District, error)
	GetAllDistricts(ctx context.Context, filter *constant.Filter) (*constant_utils.PaginatedResponse[[]*model.District], error)
	GetCityByCode(ctx context.Context, cityCode string) (action.City, error)
	GetAllCities(ctx context.Context, filter *constant.Filter) (*constant_utils.PaginatedResponse[[]*model.City], error)
	EnableBranches(r *http.Request) error
	DisableBranches(r *http.Request) error
	EnableRegion(r *http.Request) error
	DisableRegion(r *http.Request) error
	EnableDistrict(r *http.Request) error
	DisableDistrict(r *http.Request) error
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

func (h *Handler) GetBranch(ctx context.Context, branchCode string, filterParams *constant.Filter) (*model.Branch, error) {
	return h.service.GetBranch(ctx, branchCode, filterParams)
}

func (h *Handler) GetAllBranches(ctx context.Context, region, district string, filterParams *constant.Filter) (*constant_utils.PaginatedResponse[[]*model.Branch], error) {
	return h.service.GetAllBranches(ctx, region, district, filterParams)
}

func (h *Handler) GetRegionByCode(ctx context.Context, regionCode string) (action.Region, error) {
	return h.service.GetRegionByCode(ctx, regionCode)
}

func (h *Handler) GetDistrictByCode(ctx context.Context, districtCode string) (action.District, error) {
	return h.service.GetDistrictByCode(ctx, districtCode)
}

func (h *Handler) GetCityByCode(ctx context.Context, cityCode string) (action.City, error) {
	return h.service.GetCityByCode(ctx, cityCode)
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
	err := h.service.EnableBranches(ctx, req.BranchCodes, true)
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
	err := h.service.DisableBranches(ctx, req.BranchCodes, false)
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
	err := h.service.EnableRegion(ctx, req.RegionsCodes, true)
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
	err := h.service.DisableRegion(ctx, req.RegionsCodes, false)
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
	err := h.service.EnableDistrict(ctx, req.DistrictCodes, true)
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
	err := h.service.DisableDistrict(ctx, req.DistrictCodes, false)
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
	err := h.service.EnableCity(ctx, req.CitiesCode, true)
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

	fmt.Println("Citties code:", req.CitiesCode)
	if len(req.CitiesCode) == 0 {
		return fmt.Errorf("CITY_CODE_IS_REQUIRED")
	}

	ctx := r.Context()
	err := h.service.DisableCity(ctx, req.CitiesCode, false)
	if err != nil {
		return err
	}

	return nil
}
