package accountblock

import (
	"cbe-super-app-cps-action/internal/service"
	"net/http"
	"strings"

	ab_interface "cbe-super-app-cps-action/internal/constants/interfaces/account_block"
	"cbe-super-app-cps-action/internal/constants/localization"

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

	data, err := constant_utils.StructToMap(branch)
	if err != nil {
		constant_utils.SendErrorResponse(w, "UNHANDLED_SERVER_ERROR", 500, nil)
		return
	}
	constant_utils.BaseResponseMaker(data, w, "Branch retrieved successfully", http.StatusOK)
}
func (a *accountBlockAdapter) GetAllBranches(w http.ResponseWriter, r *http.Request)    {}
func (a *accountBlockAdapter) GetRegionByCode(w http.ResponseWriter, r *http.Request)   {}
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
