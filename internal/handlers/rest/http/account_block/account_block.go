package accountblock

import (
	"cbe-super-app-cps-action/internal/handlers/rest/http/account_block/core"
	"cbe-super-app-cps-action/internal/service"
	"net/http"

	ab_interface "cbe-super-app-cps-action/internal/constants/interfaces/account_block"
	ab_dto "cbe-super-app-cps-action/internal/constants/dto/account_block"
	"cbe-super-app-cps-action/internal/constants/localization"
	"cbe-super-app-cps-action/internal/constants/types"

	"gitlab.com/bersufekadgetachew/cbe-super-app-shared/shared/utils"

	local_util "cbe-super-app-cps-action/pkgs/utils"
)

type accountBlockAdapter struct {
	accountBlockApplication service.AccountBlockService
	logger                  utils.Logger
}

func InitAccountBlockAdapter(accountBlockApplication service.AccountBlockService, logger utils.Logger) ab_interface.AccountBlockAdapter {
	return &accountBlockAdapter{
		logger:                  logger,
		accountBlockApplication: accountBlockApplication,
	}
}

func (a *accountBlockAdapter) GetBranch(w http.ResponseWriter, r *http.Request)         {
	filterParams := local_util.ExtractFilterParams(r)
	branchCode, ok := local_util.GetParam(r, "branch_code")
	if !ok {
		localization.SendBadRequestResponse(w, "Branch code is required")
	}

	branch, err := a.accountBlockApplication.GetBranch(r.Context(), branchCode, filterParams)
	if err != nil {
		a.logger.Errorf("FetchUserRequest failed: %v", err)
		localization.SendErrorByCodeResponse(w, err.Error())
		return
	}

	data := core.ToBranchResponse(branch)

	localization.SendSuccessResponse(w,localization.SuccessBranchRetrieved, data)
}

func (a *accountBlockAdapter) GetAllBranches(w http.ResponseWriter, r *http.Request)    {
	region, ok := local_util.GetParam(r, "region")
	if !ok {
		localization.SendBadRequestResponse(w, "Invalid region")
		return
	}

	district, ok :=local_util.GetParam(r, "district")
	if !ok {
		localization.SendBadRequestResponse(w, "Invalid district")
		return
	}

	filterParams := local_util.ExtractFilterParams(r)

	if len(region) < 3 || len(district) < 3 {
		localization.SendBadRequestResponse(w, "Region and District codes must be at least 3 characters long")
		return
	}

	branches, err := a.accountBlockApplication.GetAllBranches(r.Context(), region, district, filterParams)
	if err != nil {
		a.logger.Errorf("GetAllBranches failed: %v", err)
		localization.SendErrorByCodeResponse(w, err.Error())
		return
	}	

	data := core.ToBranchesResponse(branches.Data)
	res := types.PaginatedResponse[[]*ab_dto.BranchResponse]{
		Data: data,
		Meta: branches.Meta,
	}
	localization.SendSuccessResponse(w, localization.SuccessBranchesRetrieved, res)
}

func (a *accountBlockAdapter) GetRegionByCode(w http.ResponseWriter, r *http.Request)   {
	regionCode, ok := local_util.GetParam(r, "region_code")
	if !ok {
		localization.SendBadRequestResponse(w, "Region code is required")
		return
	}
	if len(regionCode) < 3 {
		localization.SendBadRequestResponse(w, "Region code must be at least 3 characters long")
		return
	}

	local_util.ExtractUserFromContext()

	userInfo := local_util.ExtractUserFromContext(r.Context())

	a.logger.Infof("GetRegionByCode requested by user: %s (%s)", userInfo.UserID, userInfo.FullName)

	region, err := a.accountBlockApplication.GetRegionByCode(r.Context(), regionCode)
	if err != nil {
		h.logger.Errorf("GetRegionByCode failed: %v", err)
		localization.SendErrorByCodeResponse(w, err.Error())
		return
	}

	data, err := core..StructToMap(region)
	if err != nil {
		constant_utils.SendErrorResponse(w, "UNHANDLED_SERVER_ERROR", http.StatusInternalServerError, nil)
		return
	}
	constant_utils.BaseResponseMaker(data, w, "Region retrieved successfully", http.StatusOK)
}
func (a *accountBlockAdapter) GetAllRegions(w http.ResponseWriter, r *http.Request)     {}
func (a *accountBlockAdapter) GetDistrictByCode(w http.ResponseWriter, r *http.Request) {}
func (a *accountBlockAdapter) GetAllDistricts(w http.ResponseWriter, r *http.Request)   {}
func (a *accountBlockAdapter) GetCityByCode(w http.ResponseWriter, r *http.Request)     {}
func (a *accountBlockAdapter) GetAllCities(w http.ResponseWriter, r *http.Request)      {}
func (a *accountBlockAdapter) EnableBranches(w http.ResponseWriter, r *http.Request)    {}
func (a *accountBlockAdapter) DisableBranches(w http.ResponseWriter, r *http.Request)   {}
func (a *accountBlockAdapter) EnableRegion(w http.ResponseWriter, r *http.Request)      {}
func (a *accountBlockAdapter) DisableRegion(w http.ResponseWriter, r *http.Request)     {}
func (a *accountBlockAdapter) EnableDistrict(w http.ResponseWriter, r *http.Request)    {}
func (a *accountBlockAdapter) DisableDistrict(w http.ResponseWriter, r *http.Request)   {}
func (a *accountBlockAdapter) EnableCity(w http.ResponseWriter, r *http.Request)        {}
func (a *accountBlockAdapter) DisableCity(w http.ResponseWriter, r *http.Request)       {}
