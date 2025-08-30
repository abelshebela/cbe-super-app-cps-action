package accountblock

import (
	"cbe-super-app-cps-action/internal/handlers/rest/http/account_block/core"
	"cbe-super-app-cps-action/internal/service"
	"encoding/json"
	"fmt"
	"net/http"

	ab_dto "cbe-super-app-cps-action/internal/constants/dto/account_block"
	ab_interface "cbe-super-app-cps-action/internal/constants/interfaces/account_block"
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

func (a *accountBlockAdapter) GetBranchByCode(w http.ResponseWriter, r *http.Request) {
	branchCode, ok := local_util.GetParam(r, "branch_code")
	if !ok {
		localization.SendBadRequestResponse(w, localization.ErrorBranchCodeRequired.Code)
		return
	}

	branch, err := a.accountBlockApplication.GetBranchByCode(r.Context(), branchCode)
	if err != nil {
		a.logger.Errorf("FetchUserRequest failed: %v", err)
		localization.SendErrorByCodeResponse(w, err.Error())
		return
	}

	data := core.ToBranchResponse(branch)

	localization.SendSuccessResponse(w, localization.SuccessBranchRetrieved, data)
}

func (a *accountBlockAdapter) GetAllBranches(w http.ResponseWriter, r *http.Request) {
	filterParams := local_util.ExtractFilterParams(r)

	branches, err := a.accountBlockApplication.GetAllBranches(r.Context(), filterParams)
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

func (a *accountBlockAdapter) GetRegionByCode(w http.ResponseWriter, r *http.Request) {
	regionCode, ok := local_util.GetParam(r, "region_code")
	if !ok {
		localization.SendBadRequestResponse(w, localization.ErrorRegionCodeRequired.Code)
		return
	}

	region, err := a.accountBlockApplication.GetRegionByCode(r.Context(), regionCode)
	if err != nil {
		a.logger.Errorf("GetRegionByCode failed: %v", err)
		localization.SendErrorByCodeResponse(w, err.Error())
		return
	}

	data := core.ToRegionResponse(region)

	localization.SendSuccessResponse(w, localization.SuccessRegionRetrieved, data)
}

func (a *accountBlockAdapter) GetAllRegions(w http.ResponseWriter, r *http.Request) {
	filterParams := local_util.ExtractFilterParams(r)

	regions, err := a.accountBlockApplication.GetAllRegions(r.Context(), filterParams)
	if err != nil {
		a.logger.Errorf("GetAllRegion failed: %v", err)
		localization.SendErrorByCodeResponse(w, err.Error())
		return
	}

	data := core.ToRegionsResponse(regions.Data)
	res := types.PaginatedResponse[[]*ab_dto.RegionResponse]{
		Data: data,
		Meta: regions.Meta,
	}

	localization.SendSuccessResponse(w, localization.SuccessRegionsRetrieved, res)
}

func (a *accountBlockAdapter) GetDistrictByCode(w http.ResponseWriter, r *http.Request) {
	districtCode, ok := local_util.GetParam(r, "district_code")
	if !ok {
		localization.SendBadRequestResponse(w, localization.ErrorDistrictCodeRequired.Type)
		return
	}

	district, err := a.accountBlockApplication.GetDistrictByCode(r.Context(), districtCode)
	if err != nil {
		a.logger.Errorf("GetDistrictByCode failed: %v", err)
		localization.SendErrorByCodeResponse(w, err.Error())
		return
	}

	data := core.ToDistrictResponse(district)

	localization.SendSuccessResponse(w, localization.SuccessDistrictRetrieved, data)
}

func (a *accountBlockAdapter) GetAllDistricts(w http.ResponseWriter, r *http.Request) {
	filterParams := local_util.ExtractFilterParams(r)

	districts, err := a.accountBlockApplication.GetAllDistricts(r.Context(), filterParams)
	if err != nil {
		a.logger.Errorf("GetAllDistrict failed: %v", err)
		localization.SendErrorByCodeResponse(w, err.Error())
		return
	}

	data := core.ToDistrictsResponse(districts.Data)
	res := types.PaginatedResponse[[]*ab_dto.DistrictResponse]{
		Data: data,
		Meta: districts.Meta,
	}

	localization.SendSuccessResponse(w, localization.SuccessDistrictsRetrieved, res)
}

func (a *accountBlockAdapter) GetCityByCode(w http.ResponseWriter, r *http.Request) {
	cityCode, ok := local_util.GetParam(r, "city_code")
	if !ok {
		localization.SendBadRequestResponse(w, localization.ErrorCityCodeRequired.Code)
		return
	}

	city, err := a.accountBlockApplication.GetCityByCode(r.Context(), cityCode)
	if err != nil {
		a.logger.Errorf("GetCityByCode failed: %v", err)
		localization.SendErrorByCodeResponse(w, err.Error())
		return
	}

	data := core.ToCityResponse(city)

	localization.SendSuccessResponse(w, localization.SuccessCityRetrieved, data)
}

func (a *accountBlockAdapter) GetAllCities(w http.ResponseWriter, r *http.Request) {
	filterParams := local_util.ExtractFilterParams(r)
	fmt.Println("****************FindAllCitiesWithPagination****************")

	cities, err := a.accountBlockApplication.GetAllCities(r.Context(), filterParams)
	if err != nil {
		a.logger.Errorf("GetAllCity failed: %v", err)
		localization.SendErrorByCodeResponse(w, err.Error())
		return
	}

	data := core.ToCitiesResponse(cities.Data)
	res := types.PaginatedResponse[[]*ab_dto.CityResponse]{
		Data: data,
		Meta: cities.Meta,
	}

	localization.SendSuccessResponse(w, localization.SuccessCitiesRetrieved, res)
}

func (a *accountBlockAdapter) EnableBranches(w http.ResponseWriter, r *http.Request) {
	var req ab_dto.EnableOrDisableBranches
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		localization.SendBadRequestResponse(w, localization.ErrorInvalidRequestBody.Code)
		return
	}

	// Trim any whitespace
	req.Clean()

	if len(req.BranchCodes) == 0 {
		localization.SendBadRequestResponse(w, localization.ErrorBranchCodeRequired.Code)
		return
	}

	err := a.accountBlockApplication.EnableOrDisableBranches(r.Context(), req.BranchCodes, true)
	if err != nil {
		localization.SendErrorByCodeResponse(w, err.Error())
		return
	}

	localization.SendSuccessResponse(w, localization.SuccessEnableBranchesRequestSent, nil)
}

func (a *accountBlockAdapter) DisableBranches(w http.ResponseWriter, r *http.Request) {
	var req ab_dto.EnableOrDisableBranches
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		localization.SendBadRequestResponse(w, localization.ErrorInvalidRequestBody.Code)
		return
	}

	// Trim any whitespace
	req.Clean()

	if len(req.BranchCodes) == 0 {
		localization.SendBadRequestResponse(w, localization.ErrorBranchCodeRequired.Code)
		return
	}

	err := a.accountBlockApplication.EnableOrDisableBranches(r.Context(), req.BranchCodes, false)
	if err != nil {
		localization.SendErrorByCodeResponse(w, err.Error())
		return
	}

	localization.SendSuccessResponse(w, localization.SuccessDisableBranchesRequestSent, nil)
}

func (a *accountBlockAdapter) EnableRegions(w http.ResponseWriter, r *http.Request) {
	var req ab_dto.EnableOrDisableRegions
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		localization.SendBadRequestResponse(w, localization.ErrorInvalidRequestBody.Code)
		return
	}

	// Trim any whitespace
	req.Clean()

	if len(req.RegionsCodes) == 0 {
		localization.SendBadRequestResponse(w, localization.ErrorRegionCodeRequired.Code)
		return
	}

	err := a.accountBlockApplication.EnableOrDisableRegions(r.Context(), req.RegionsCodes, true)
	if err != nil {
		localization.SendErrorByCodeResponse(w, err.Error())
		return
	}

	localization.SendSuccessResponse(w, localization.SuccessEnableRegionsRequestSent, nil)
}

func (a *accountBlockAdapter) DisableRegions(w http.ResponseWriter, r *http.Request) {
	var req ab_dto.EnableOrDisableRegions
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		localization.SendBadRequestResponse(w, localization.ErrorInvalidRequestBody.Code)
		return

	}

	// Trim any whitespace
	req.Clean()

	if len(req.RegionsCodes) == 0 {
		localization.SendBadRequestResponse(w, localization.ErrorRegionCodeRequired.Code)
		return
	}

	err := a.accountBlockApplication.EnableOrDisableRegions(r.Context(), req.RegionsCodes, false)
	if err != nil {
		localization.SendErrorByCodeResponse(w, err.Error())
		return
	}

	localization.SendSuccessResponse(w, localization.SuccessDisableRegionsRequestSent, nil)
}

func (a *accountBlockAdapter) EnableDistricts(w http.ResponseWriter, r *http.Request) {
	var req ab_dto.EnableOrDisableDistricts
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		localization.SendBadRequestResponse(w, localization.ErrorInvalidRequestBody.Code)
		return
	}

	// Trim any whitespace
	req.Clean()

	if len(req.DistrictCodes) == 0 {
		localization.SendBadRequestResponse(w, localization.ErrorDistrictCodeRequired.Code)
		return
	}

	err := a.accountBlockApplication.EnableOrDisableDistricts(r.Context(), req.DistrictCodes, true)
	if err != nil {
		localization.SendErrorByCodeResponse(w, err.Error())
		return
	}

	localization.SendSuccessResponse(w, localization.SuccessEnableDistrictsRequestSent, nil)
}

func (a *accountBlockAdapter) DisableDistricts(w http.ResponseWriter, r *http.Request) {
	var req ab_dto.EnableOrDisableDistricts
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		localization.SendBadRequestResponse(w, localization.ErrorInvalidRequestBody.Code)
		return
	}

	// Trim any whitespace
	req.Clean()

	if len(req.DistrictCodes) == 0 {
		localization.SendBadRequestResponse(w, localization.ErrorDistrictCodeRequired.Code)
		return
	}

	err := a.accountBlockApplication.EnableOrDisableDistricts(r.Context(), req.DistrictCodes, false)
	if err != nil {
		localization.SendErrorByCodeResponse(w, err.Error())
		return
	}

	localization.SendSuccessResponse(w, localization.SuccessDisableDistrictsRequestSent, nil)
}

func (a *accountBlockAdapter) EnableCities(w http.ResponseWriter, r *http.Request) {
	var req ab_dto.EnableOrDisableCities
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		localization.SendBadRequestResponse(w, localization.ErrorInvalidRequestBody.Code)
		return
	}

	// Trim any whitespace
	req.Clean()

	if len(req.CitiesCode) == 0 {
		localization.SendBadRequestResponse(w, localization.ErrorCityCodeRequired.Type)
		return
	}

	err := a.accountBlockApplication.EnableOrDisableCities(r.Context(), req.CitiesCode, true)
	if err != nil {
		localization.SendErrorByCodeResponse(w, err.Error())
		return
	}

	localization.SendSuccessResponse(w, localization.SuccessEnableCitiesRequestSent, nil)
}

func (a *accountBlockAdapter) DisableCities(w http.ResponseWriter, r *http.Request) {
	var req ab_dto.EnableOrDisableCities
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		localization.SendBadRequestResponse(w, localization.ErrorInvalidRequestBody.Code)
		return
	}

	// Trim any whitespace
	req.Clean()

	if len(req.CitiesCode) == 0 {
		localization.SendBadRequestResponse(w, localization.ErrorCityCodeRequired.Type)
		return
	}
	fmt.Println("*************************************")
	err := a.accountBlockApplication.EnableOrDisableCities(r.Context(), req.CitiesCode, false)
	if err != nil {
		localization.SendErrorByCodeResponse(w, err.Error())
		return
	}

	localization.SendSuccessResponse(w, localization.SuccessDisableCitiesRequestSent, nil)
}
