package accountblock

import (
	"cbe-super-app-cps-action/internal/handlers/rest/http/account_block/core"
	"cbe-super-app-cps-action/internal/service"
	"encoding/json"
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

// GetBranchByCode godoc
//
//	@Summary		Get branch by code
//	@Description	Retrieve a specific branch by its branch code
//	@Tags			Account Block - Branches
//	@Accept			json
//	@Produce		json
//	@Param			branch_code	path		string					true	"Branch Code"
//	@Success		200			{object}	map[string]interface{}	"Branch retrieved successfully"
//	@Failure		400			{object}	map[string]interface{}	"Bad request"
//	@Failure		500			{object}	map[string]interface{}	"Internal server error"
//	@Security		BearerAuth
//	@Router			/account_block/branches/{branch_code} [get]
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

// GetAllBranches godoc
//
//	@Summary		Get all branches
//	@Description	Retrieve all branches with pagination
//	@Tags			Account Block - Branches
//	@Accept			json
//	@Produce		json
//	@Param			page		query		int						false	"Page number"		default(1)
//	@Param			per_page	query		int						false	"Items per page"	default(10)
//	@Param			search		query		string					false	"Search term"
//	@Success		200			{object}	map[string]interface{}	"Branches retrieved successfully"
//	@Failure		500			{object}	map[string]interface{}	"Internal server error"
//	@Security		BearerAuth
//	@Router			/account_block/branches [get]
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

// GetRegionByCode godoc
//
//	@Summary		Get region by code
//	@Description	Retrieve a specific region by its region code
//	@Tags			Account Block - Regions
//	@Accept			json
//	@Produce		json
//	@Param			region_code	path		string					true	"Region Code"
//	@Success		200			{object}	map[string]interface{}	"Region retrieved successfully"
//	@Failure		400			{object}	map[string]interface{}	"Bad request"
//	@Failure		500			{object}	map[string]interface{}	"Internal server error"
//	@Security		BearerAuth
//	@Router			/account_block/regions/{region_code} [get]
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

// GetAllRegions godoc
//
//	@Summary		Get all regions
//	@Description	Retrieve all regions with pagination
//	@Tags			Account Block - Regions
//	@Accept			json
//	@Produce		json
//	@Param			page		query		int						false	"Page number"		default(1)
//	@Param			per_page	query		int						false	"Items per page"	default(10)
//	@Param			search		query		string					false	"Search term"
//	@Success		200			{object}	map[string]interface{}	"Regions retrieved successfully"
//	@Failure		500			{object}	map[string]interface{}	"Internal server error"
//	@Security		BearerAuth
//	@Router			/account_block/regions [get]
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

// GetDistrictByCode godoc
//
//	@Summary		Get district by code
//	@Description	Retrieve a specific district by its district code
//	@Tags			Account Block - Districts
//	@Accept			json
//	@Produce		json
//	@Param			district_code	path		string					true	"District Code"
//	@Success		200				{object}	map[string]interface{}	"District retrieved successfully"
//	@Failure		400				{object}	map[string]interface{}	"Bad request"
//	@Failure		500				{object}	map[string]interface{}	"Internal server error"
//	@Security		BearerAuth
//	@Router			/account_block/districts/{district_code} [get]
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

// GetAllDistricts godoc
//
//	@Summary		Get all districts
//	@Description	Retrieve all districts with pagination
//	@Tags			Account Block - Districts
//	@Accept			json
//	@Produce		json
//	@Param			page		query		int						false	"Page number"		default(1)
//	@Param			per_page	query		int						false	"Items per page"	default(10)
//	@Param			search		query		string					false	"Search term"
//	@Success		200			{object}	map[string]interface{}	"Districts retrieved successfully"
//	@Failure		500			{object}	map[string]interface{}	"Internal server error"
//	@Security		BearerAuth
//	@Router			/account_block/districts [get]
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

// GetCityByCode godoc
//
//	@Summary		Get city by code
//	@Description	Retrieve a specific city by its city code
//	@Tags			Account Block - Cities
//	@Accept			json
//	@Produce		json
//	@Param			city_code	path		string					true	"City Code"
//	@Success		200			{object}	map[string]interface{}	"City retrieved successfully"
//	@Failure		400			{object}	map[string]interface{}	"Bad request"
//	@Failure		500			{object}	map[string]interface{}	"Internal server error"
//	@Security		BearerAuth
//	@Router			/account_block/city/{city_code} [get]
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

// GetAllCities godoc
//
//	@Summary		Get all cities
//	@Description	Retrieve all cities with pagination
//	@Tags			Account Block - Cities
//	@Accept			json
//	@Produce		json
//	@Param			page		query		int						false	"Page number"		default(1)
//	@Param			per_page	query		int						false	"Items per page"	default(10)
//	@Param			search		query		string					false	"Search term"
//	@Success		200			{object}	map[string]interface{}	"Cities retrieved successfully"
//	@Failure		500			{object}	map[string]interface{}	"Internal server error"
//	@Security		BearerAuth
//	@Router			/account_block/city [get]
func (a *accountBlockAdapter) GetAllCities(w http.ResponseWriter, r *http.Request) {
	filterParams := local_util.ExtractFilterParams(r)

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

// EnableBranches godoc
//
//	@Summary		Enable branches
//	@Description	Enable one or more branches by their branch codes
//	@Tags			Account Block - Branches
//	@Accept			json
//	@Produce		json
//	@Param			request	body		accountblock.EnableOrDisableBranches	true	"Branch codes to enable"
//	@Success		200		{object}	map[string]interface{}					"Branches enabled successfully"
//	@Failure		400		{object}	map[string]interface{}					"Bad request"
//	@Failure		500		{object}	map[string]interface{}					"Internal server error"
//	@Security		BearerAuth
//	@Router			/account_block/branches/enable [post]
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

// DisableBranches godoc
//
//	@Summary		Disable branches
//	@Description	Disable one or more branches by their branch codes
//	@Tags			Account Block - Branches
//	@Accept			json
//	@Produce		json
//	@Param			request	body		accountblock.EnableOrDisableBranches	true	"Branch codes to disable"
//	@Success		200		{object}	map[string]interface{}					"Branches disabled successfully"
//	@Failure		400		{object}	map[string]interface{}					"Bad request"
//	@Failure		500		{object}	map[string]interface{}					"Internal server error"
//	@Security		BearerAuth
//	@Router			/account_block/branches/disable [post]
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

// EnableRegions godoc
//
//	@Summary		Enable regions
//	@Description	Enable one or more regions by their region codes
//	@Tags			Account Block - Regions
//	@Accept			json
//	@Produce		json
//	@Param			request	body		accountblock.EnableOrDisableRegions	true	"Region codes to enable"
//	@Success		200		{object}	map[string]interface{}				"Regions enabled successfully"
//	@Failure		400		{object}	map[string]interface{}				"Bad request"
//	@Failure		500		{object}	map[string]interface{}				"Internal server error"
//	@Security		BearerAuth
//	@Router			/account_block/regions/enable [post]
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

// DisableRegions godoc
//
//	@Summary		Disable regions
//	@Description	Disable one or more regions by their region codes
//	@Tags			Account Block - Regions
//	@Accept			json
//	@Produce		json
//	@Param			request	body		accountblock.EnableOrDisableRegions	true	"Region codes to disable"
//	@Success		200		{object}	map[string]interface{}				"Regions disabled successfully"
//	@Failure		400		{object}	map[string]interface{}				"Bad request"
//	@Failure		500		{object}	map[string]interface{}				"Internal server error"
//	@Security		BearerAuth
//	@Router			/account_block/regions/disable [post]
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

// EnableDistricts godoc
//
//	@Summary		Enable districts
//	@Description	Enable one or more districts by their district codes
//	@Tags			Account Block - Districts
//	@Accept			json
//	@Produce		json
//	@Param			request	body		accountblock.EnableOrDisableDistricts	true	"District codes to enable"
//	@Success		200		{object}	map[string]interface{}					"Districts enabled successfully"
//	@Failure		400		{object}	map[string]interface{}					"Bad request"
//	@Failure		500		{object}	map[string]interface{}					"Internal server error"
//	@Security		BearerAuth
//	@Router			/account_block/districts/enable [post]
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

// DisableDistricts godoc
//
//	@Summary		Disable districts
//	@Description	Disable one or more districts by their district codes
//	@Tags			Account Block - Districts
//	@Accept			json
//	@Produce		json
//	@Param			request	body		accountblock.EnableOrDisableDistricts	true	"District codes to disable"
//	@Success		200		{object}	map[string]interface{}					"Districts disabled successfully"
//	@Failure		400		{object}	map[string]interface{}					"Bad request"
//	@Failure		500		{object}	map[string]interface{}					"Internal server error"
//	@Security		BearerAuth
//	@Router			/account_block/districts/disable [post]
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

// EnableCities godoc
//
//	@Summary		Enable cities
//	@Description	Enable one or more cities by their city codes
//	@Tags			Account Block - Cities
//	@Accept			json
//	@Produce		json
//	@Param			request	body		accountblock.EnableOrDisableCities	true	"City codes to enable"
//	@Success		200		{object}	map[string]interface{}				"Cities enabled successfully"
//	@Failure		400		{object}	map[string]interface{}				"Bad request"
//	@Failure		500		{object}	map[string]interface{}				"Internal server error"
//	@Security		BearerAuth
//	@Router			/account_block/city/enable [post]
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

// DisableCities godoc
//
//	@Summary		Disable cities
//	@Description	Disable one or more cities by their city codes
//	@Tags			Account Block - Cities
//	@Accept			json
//	@Produce		json
//	@Param			request	body		accountblock.EnableOrDisableCities	true	"City codes to disable"
//	@Success		200		{object}	map[string]interface{}				"Cities disabled successfully"
//	@Failure		400		{object}	map[string]interface{}				"Bad request"
//	@Failure		500		{object}	map[string]interface{}				"Internal server error"
//	@Security		BearerAuth
//	@Router			/account_block/city/disable [post]
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

	err := a.accountBlockApplication.EnableOrDisableCities(r.Context(), req.CitiesCode, false)
	if err != nil {
		localization.SendErrorByCodeResponse(w, err.Error())
		return
	}

	localization.SendSuccessResponse(w, localization.SuccessDisableCitiesRequestSent, nil)
}
