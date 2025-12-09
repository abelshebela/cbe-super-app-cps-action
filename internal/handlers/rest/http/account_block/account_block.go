package accountblock

import (
	"cbe-super-app-cps-action/internal/handlers/rest/http/account_block/core"
	"cbe-super-app-cps-action/internal/service"
	"encoding/json"
	"net/http"

	accountblock "cbe-super-app-cps-action/internal/constants/dto/account_block"
	ab_interface "cbe-super-app-cps-action/internal/constants/interfaces/account_block"
	"cbe-super-app-cps-action/internal/constants/localization"
	"cbe-super-app-cps-action/internal/constants/types"

	"gitlab.com/bersufekadgetachew/cbe-super-app-shared/shared/utils"

	local_util "cbe-super-app-cps-action/pkgs/utils"
)

type paginated_account_block_resp types.PaginatedResponse[[]*accountblock.AccountBlockResponse]

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
//	@Description	Retrieve a specific branch by its branch code.
//	@Tags			Account Block - Branches
//	@Accept			json
//	@Produce		json
//	@Param			branch_code	path		string																	true	"Branch Code"	example(BR001)
//	@Success		200			{object}	localization.StandardResponse{data=accountblock.AccountBlockResponse}	"Branch retrieved"
//	@Failure		400			{object}	localization.StandardResponse{data=nil}									"Bad request"
//	@Failure		404			{object}	localization.StandardResponse{data=nil}									"Not found"
//	@Failure		500			{object}	localization.StandardResponse{data=nil}									"Server error"
//	@Security		BearerAuth
//	@Router			/account_block/branches/{branch_code} [get]
func (a *accountBlockAdapter) GetBranchByCode(w http.ResponseWriter, r *http.Request) {
	branchCode, ok := local_util.GetParam(r, "branch_code")
	if !ok {
		a.logger.Errorf("Failed to get branch_code from the param")
		localization.SendBadRequestResponse(w, localization.ErrorBranchCodeRequired.Code)
		return
	}

	branch, err := a.accountBlockApplication.GetBranchByCode(r.Context(), branchCode)
	if err != nil {
		a.logger.Errorf("FetchUserRequest failed: %v", err)
		localization.SendErrorByCodeResponse(w, err.Error())
		return
	}

	data := core.ToAccountBlockResponse(branch)

	localization.SendSuccessResponse(w, localization.SuccessBranchRetrieved, data)
}

// GetAllBranches godoc
//
//	@Summary		Get all branches
//	@Description	Retrieve all branches with pagination and optional filters.
//	@Tags			Account Block - Branches
//	@Accept			json
//	@Produce		json
//	@Param			page		query		int																	false	"Page number"										default(1)	minimum(1)	example(1)
//	@Param			per_page	query		int																	false	"Items per page"									default(10)	minimum(1)	maximum(100)	example(10)
//	@Param			search		query		string																false	"Search term"										example("Addis")
//	@Param			filters		query		string																false	"JSON encoded filters (enabled, region, district)"	example("{\"enabled\":true}")
//	@Success		200			{object}	localization.StandardResponse{data=paginated_account_block_resp}	"Branches retrieved successfully"
//	@Failure		500			{object}	localization.StandardResponse{data=nil}								"internal Server error"
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

	localization.SendSuccessResponse(w, localization.SuccessBranchesRetrieved, branches)
}

// GetRegionByCode godoc
//
//	@Summary		Get region by code
//	@Description	Retrieve a specific region by its region code.
//	@Tags			Account Block - Regions
//	@Accept			json
//	@Produce		json
//	@Param			region_code	path		string																	true	"Region Code"	example(RG001)
//	@Success		200			{object}	localization.StandardResponse{data=accountblock.AccountBlockResponse}	"Region retrieved"
//	@Failure		400			{object}	localization.StandardResponse{data=nil}									"Bad request"
//	@Failure		404			{object}	localization.StandardResponse{data=nil}									"Not found"
//	@Failure		500			{object}	localization.StandardResponse{data=nil}									"Server error"
//	@Security		BearerAuth
//	@Router			/account_block/regions/{region_code} [get]
func (a *accountBlockAdapter) GetRegionByCode(w http.ResponseWriter, r *http.Request) {
	regionCode, ok := local_util.GetParam(r, "region_code")
	if !ok {
		a.logger.Errorf("Failed to get region_code from the param")
		localization.SendBadRequestResponse(w, localization.ErrorRegionCodeRequired.Code)
		return
	}

	region, err := a.accountBlockApplication.GetRegionByCode(r.Context(), regionCode)
	if err != nil {
		a.logger.Errorf("GetRegionByCode failed: %v", err)
		localization.SendErrorByCodeResponse(w, err.Error())
		return
	}

	data := core.ToAccountBlockResponse(region)

	localization.SendSuccessResponse(w, localization.SuccessRegionRetrieved, data)
}

// GetAllRegions godoc
//
//	@Summary		Get all regions
//	@Description	Retrieve all regions with pagination and optional filters.
//	@Tags			Account Block - Regions
//	@Accept			json
//	@Produce		json
//	@Param			page		query		int																	false	"Page number"						default(1)	minimum(1)	example(1)
//	@Param			per_page	query		int																	false	"Items per page"					default(10)	minimum(1)	maximum(100)	example(10)
//	@Param			search		query		string																false	"Search term"						example("Addis")
//	@Param			filters		query		string																false	"JSON encoded filters (enabled)"	example("{\"enabled\":true}")
//	@Success		200			{object}	localization.StandardResponse{data=paginated_account_block_resp}	"Regions retrieved successfully"
//	@Failure		500			{object}	localization.StandardResponse{data=nil}								"internal Server error"
//	@Security		BearerAuth
//	@Router			/account_block/regions [get]
func (a *accountBlockAdapter) GetAllRegions(w http.ResponseWriter, r *http.Request) {
	filterParams := local_util.ExtractFilterParams(r)

	regions, err := a.accountBlockApplication.GetAllRegions(r.Context(), filterParams)
	if err != nil {
		a.logger.Errorf("GetAllRegions failed: %v", err)
		localization.SendErrorByCodeResponse(w, err.Error())
		return
	}

	localization.SendSuccessResponse(w, localization.SuccessRegionsRetrieved, regions)
}

// GetDistrictByCode godoc
//
//	@Summary		Get district by code
//	@Description	Retrieve a specific district by its district code.
//	@Tags			Account Block - Districts
//	@Accept			json
//	@Produce		json
//	@Param			district_code	path		string																	true	"District Code"	example(DS001)
//	@Success		200				{object}	localization.StandardResponse{data=accountblock.AccountBlockResponse}	"District retrieved successfully"
//	@Failure		400				{object}	localization.StandardResponse{data=nil}									"Bad request"
//	@Failure		404				{object}	localization.StandardResponse{data=nil}									"Not found"
//	@Failure		500				{object}	localization.StandardResponse{data=nil}									"internal Server error"
//	@Security		BearerAuth
//	@Router			/account_block/districts/{district_code} [get]
func (a *accountBlockAdapter) GetDistrictByCode(w http.ResponseWriter, r *http.Request) {
	districtCode, ok := local_util.GetParam(r, "district_code")
	if !ok {
		a.logger.Errorf("Failed to get district_code from the param")
		localization.SendBadRequestResponse(w, localization.ErrorDistrictCodeRequired.Code)
		return
	}

	district, err := a.accountBlockApplication.GetDistrictByCode(r.Context(), districtCode)
	if err != nil {
		a.logger.Errorf("GetDistrictByCode failed: %v", err)
		localization.SendErrorByCodeResponse(w, err.Error())
		return
	}

	localization.SendSuccessResponse(w, localization.SuccessDistrictRetrieved, district)
}

// GetAllDistricts godoc
//
//	@Summary		Get all districts
//	@Description	Retrieve all districts with pagination and optional filters.
//	@Tags			Account Block - Districts
//	@Accept			json
//	@Produce		json
//	@Param			page		query		int																	false	"Page number"								default(1)	minimum(1)	example(1)
//	@Param			per_page	query		int																	false	"Items per page"							default(10)	minimum(1)	maximum(100)	example(10)
//	@Param			search		query		string																false	"Search term"								example("Bole")
//	@Param			filters		query		string																false	"JSON encoded filters (enabled, region)"	example("{\"enabled\":true,\"region_id\":\"RG001\"}")
//	@Success		200			{object}	localization.StandardResponse{data=paginated_account_block_resp}	"Districts retrieved successfully"
//	@Failure		500			{object}	localization.StandardResponse{data=nil}								"internal Server error"
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

	localization.SendSuccessResponse(w, localization.SuccessDistrictsRetrieved, districts)
}

// GetCityByCode godoc
//
//	@Summary		Get city by code
//	@Description	Retrieve a specific city by its city code.
//	@Tags			Account Block - Cities
//	@Accept			json
//	@Produce		json
//	@Param			city_code	path		string																	true	"City Code"	example(CT001)
//	@Success		200			{object}	localization.StandardResponse{data=accountblock.AccountBlockResponse}	"City retrieved successfully"
//	@Failure		400			{object}	localization.StandardResponse{data=nil}									"Bad request"
//	@Failure		404			{object}	localization.StandardResponse{data=nil}									"Not found"
//	@Failure		500			{object}	localization.StandardResponse{data=nil}									"internal Server error"
//	@Security		BearerAuth
//	@Router			/account_block/cities/{city_code} [get]
func (a *accountBlockAdapter) GetCityByCode(w http.ResponseWriter, r *http.Request) {
	cityCode, ok := local_util.GetParam(r, "city_code")
	if !ok {
		a.logger.Errorf("Failed to get city_code from the param")
		localization.SendBadRequestResponse(w, localization.ErrorCityCodeRequired.Code)
		return
	}

	city, err := a.accountBlockApplication.GetCityByCode(r.Context(), cityCode)
	if err != nil {
		a.logger.Errorf("GetCityByCode failed: %v", err)
		localization.SendErrorByCodeResponse(w, err.Error())
		return
	}

	localization.SendSuccessResponse(w, localization.SuccessCityRetrieved, city)
}

// GetAllCities godoc
//
//	@Summary		Get all cities
//	@Description	Retrieve all cities with pagination and optional filters.
//	@Tags			Account Block - Cities
//	@Accept			json
//	@Produce		json
//	@Param			page		query		int																	false	"Page number"								default(1)	minimum(1)	example(1)
//	@Param			per_page	query		int																	false	"Items per page"							default(10)	minimum(1)	maximum(100)	example(10)
//	@Param			search		query		string																false	"Search term"								example("Addis")
//	@Param			filters		query		string																false	"JSON encoded filters (enabled, district)"	example("{\"enabled\":true}")
//	@Success		200			{object}	localization.StandardResponse{data=paginated_account_block_resp}	"Cities retrieved"
//	@Failure		500			{object}	localization.StandardResponse{data=nil}								"Server error"
//	@Security		BearerAuth
//	@Router			/account_block/cities [get]
func (a *accountBlockAdapter) GetAllCities(w http.ResponseWriter, r *http.Request) {
	filterParams := local_util.ExtractFilterParams(r)

	cities, err := a.accountBlockApplication.GetAllCities(r.Context(), filterParams)
	if err != nil {
		a.logger.Errorf("GetAllCities failed: %v", err)
		localization.SendErrorByCodeResponse(w, err.Error())
		return
	}

	localization.SendSuccessResponse(w, localization.SuccessCitiesRetrieved, cities)
}

// EnableBranches godoc
//
//	@Summary		Enable multiple branches
//	@Description	Enable multiple branches by codes.
//	@Tags			Account Block - Branches
//	@Accept			json
//	@Produce		json
//	@Param			request	body		accountblock.EnableOrDisableBranches	true	"Branch codes and reason"	example({"branches_code":["BR001","BR002"], "reason": "foo"})
//	@Success		200		{object}	localization.StandardResponse{data=nil}	"Enable request submitted"
//	@Failure		400		{object}	localization.StandardResponse{data=nil}	"Invalid request"
//	@Failure		500		{object}	localization.StandardResponse{data=nil}	"Server error"
//	@Security		BearerAuth
//	@Router			/account_block/branches/enable [post]
func (a *accountBlockAdapter) EnableBranches(w http.ResponseWriter, r *http.Request) {
	var req accountblock.EnableOrDisableBranches
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		a.logger.Errorf("Failed to decode request body: %v", err)
		localization.SendBadRequestResponse(w, localization.ErrorInvalidRequestBody.Code)
		return
	}

	// Trim any whitespace
	req.Clean()

	if err := req.Validate(); err != nil {
		a.logger.Errorf("Failed to validate request: %v", err)
		localization.SendBadRequestResponse(w, err.Error())
		return
	}

	if len(req.BranchIds) == 0 {
		a.logger.Errorf("Branch codes are required")
		localization.SendBadRequestResponse(w, localization.ErrorBranchCodeRequired.Code)
		return
	}

	err := a.accountBlockApplication.EnableOrDisableBranches(r.Context(), req.BranchIds, req.Reason, true)
	if err != nil {
		a.logger.Errorf("EnableOrDisableBranches failed: %v", err)
		localization.SendErrorByCodeResponse(w, err.Error())
		return
	}

	localization.SendSuccessResponse(w, localization.SuccessEnableBranchesRequestSent, nil)
}

// DisableBranches godoc
//
//	@Summary		Disable multiple branches
//	@Description	Disable multiple branches by codes.
//	@Tags			Account Block - Branches
//	@Accept			json
//	@Produce		json
//	@Param			request	body		accountblock.EnableOrDisableBranches	true	"Branch codes and reason"	example({"branches_code":["BR001","BR002"], "reason": "foo"})
//	@Success		200		{object}	localization.StandardResponse{data=nil}	"Disable request submitted"
//	@Failure		400		{object}	localization.StandardResponse{data=nil}	"Invalid request"
//	@Failure		500		{object}	localization.StandardResponse{data=nil}	"Server error"
//	@Security		BearerAuth
//	@Router			/account_block/branches/disable [post]
func (a *accountBlockAdapter) DisableBranches(w http.ResponseWriter, r *http.Request) {
	var req accountblock.EnableOrDisableBranches
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		a.logger.Errorf("Failed to decode request body: %v", err)
		localization.SendBadRequestResponse(w, localization.ErrorInvalidRequestBody.Code)
		return
	}

	// Trim any whitespace
	req.Clean()

	if err := req.Validate(); err != nil {
		a.logger.Errorf("Failed to validate request: %v", err)
		localization.SendBadRequestResponse(w, err.Error())
		return
	}

	if len(req.BranchIds) == 0 {
		a.logger.Errorf("Branch codes are required")
		localization.SendBadRequestResponse(w, localization.ErrorBranchCodeRequired.Code)
		return
	}

	err := a.accountBlockApplication.EnableOrDisableBranches(r.Context(), req.BranchIds, req.Reason, false)
	if err != nil {
		a.logger.Errorf("EnableOrDisableBranches failed: %v", err)
		localization.SendErrorByCodeResponse(w, err.Error())
		return
	}

	localization.SendSuccessResponse(w, localization.SuccessDisableBranchesRequestSent, nil)
}

// EnableRegions godoc
//
//	@Summary		Enable multiple regions
//	@Description	Enable multiple regions by codes.
//	@Tags			Account Block - Regions
//	@Accept			json
//	@Produce		json
//	@Param			request	body		accountblock.EnableOrDisableRegions		true	"Region codes and reason"	example({"regions_code":["RG001","RG002"], "reason": "foo"})
//	@Success		200		{object}	localization.StandardResponse{data=nil}	"Enable request submitted"
//	@Failure		400		{object}	localization.StandardResponse{data=nil}	"Invalid request"
//	@Failure		500		{object}	localization.StandardResponse{data=nil}	"Server error"
//	@Security		BearerAuth
//	@Router			/account_block/regions/enable [post]
func (a *accountBlockAdapter) EnableRegions(w http.ResponseWriter, r *http.Request) {
	var req accountblock.EnableOrDisableRegions
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		a.logger.Errorf("Failed to decode request body: %v", err)
		localization.SendBadRequestResponse(w, localization.ErrorInvalidRequestBody.Code)
		return
	}

	// Trim any whitespace
	req.Clean()

	if err := req.Validate(); err != nil {
		a.logger.Errorf("Failed to validate request: %v", err)
		localization.SendBadRequestResponse(w, err.Error())
		return
	}

	if len(req.RegionIds) == 0 {
		a.logger.Errorf("Region codes are required")
		localization.SendBadRequestResponse(w, localization.ErrorRegionCodeRequired.Code)
		return
	}

	err := a.accountBlockApplication.EnableOrDisableRegions(r.Context(), req.RegionIds, req.Reason, true)
	if err != nil {
		a.logger.Errorf("EnableOrDisableRegions failed: %v", err)
		localization.SendErrorByCodeResponse(w, err.Error())
		return
	}

	localization.SendSuccessResponse(w, localization.SuccessEnableRegionsRequestSent, nil)
}

// DisableRegions godoc
//
//	@Summary		Disable multiple regions
//	@Description	Disable multiple regions by codes.
//	@Tags			Account Block - Regions
//	@Accept			json
//	@Produce		json
//	@Param			request	body		accountblock.EnableOrDisableRegions		true	"Region codes and reason"	example({"regions_code":["RG001","RG002"], "reason": "foo"})
//	@Success		200		{object}	localization.StandardResponse{data=nil}	"Disable request submitted"
//	@Failure		400		{object}	localization.StandardResponse{data=nil}	"Invalid request"
//	@Failure		500		{object}	localization.StandardResponse{data=nil}	"Server error"
//	@Security		BearerAuth
//	@Router			/account_block/regions/disable [post]
func (a *accountBlockAdapter) DisableRegions(w http.ResponseWriter, r *http.Request) {
	var req accountblock.EnableOrDisableRegions
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		a.logger.Errorf("Failed to decode request body: %v", err)
		localization.SendBadRequestResponse(w, localization.ErrorInvalidRequestBody.Code)
		return

	}

	// Trim any whitespace
	req.Clean()

	if err := req.Validate(); err != nil {
		a.logger.Errorf("Failed to validate request: %v", err)
		localization.SendBadRequestResponse(w, err.Error())
		return
	}

	if len(req.RegionIds) == 0 {
		a.logger.Errorf("Region codes are required")
		localization.SendBadRequestResponse(w, localization.ErrorRegionCodeRequired.Code)
		return
	}

	err := a.accountBlockApplication.EnableOrDisableRegions(r.Context(), req.RegionIds, req.Reason, false)
	if err != nil {
		a.logger.Errorf("EnableOrDisableRegions failed: %v", err)
		localization.SendErrorByCodeResponse(w, err.Error())
		return
	}

	localization.SendSuccessResponse(w, localization.SuccessDisableRegionsRequestSent, nil)
}

// EnableDistricts godoc
//
//	@Summary		Enable multiple districts
//	@Description	Enable multiple districts by codes.
//	@Tags			Account Block - Districts
//	@Accept			json
//	@Produce		json
//	@Param			request	body		accountblock.EnableOrDisableDistricts	true	"District codes and reason"	example({"districts_code":["DS001","DS002"], "reason": "foo"})
//	@Success		200		{object}	localization.StandardResponse{data=nil}	"Enable request submitted"
//	@Failure		400		{object}	localization.StandardResponse{data=nil}	"Invalid request"
//	@Failure		500		{object}	localization.StandardResponse{data=nil}	"Server error"
//	@Security		BearerAuth
//	@Router			/account_block/districts/enable [post]
func (a *accountBlockAdapter) EnableDistricts(w http.ResponseWriter, r *http.Request) {
	var req accountblock.EnableOrDisableDistricts
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		a.logger.Errorf("Failed to decode request body: %v", err)
		localization.SendBadRequestResponse(w, localization.ErrorInvalidRequestBody.Code)
		return
	}

	// Trim any whitespace
	req.Clean()

	if err := req.Validate(); err != nil {
		a.logger.Errorf("Failed to validate request: %v", err)
		localization.SendBadRequestResponse(w, err.Error())
		return
	}

	if len(req.DistrictIds) == 0 {
		a.logger.Errorf("District codes are required")
		localization.SendBadRequestResponse(w, localization.ErrorDistrictCodeRequired.Code)
		return
	}

	err := a.accountBlockApplication.EnableOrDisableDistricts(r.Context(), req.DistrictIds, req.Reason, true)
	if err != nil {
		a.logger.Errorf("EnableOrDisableDistricts failed: %v", err)
		localization.SendErrorByCodeResponse(w, err.Error())
		return
	}

	localization.SendSuccessResponse(w, localization.SuccessEnableDistrictsRequestSent, nil)
}

// DisableDistricts godoc
//
//	@Summary		Disable multiple districts
//	@Description	Disable multiple districts by codes.
//	@Tags			Account Block - Districts
//	@Accept			json
//	@Produce		json
//	@Param			request	body		accountblock.EnableOrDisableDistricts	true	"District codes and reason"	example({"districts_code":["DS001","DS002"], "reason": "foo"})
//	@Success		200		{object}	localization.StandardResponse{data=nil}	"Disable request submitted"
//	@Failure		400		{object}	localization.StandardResponse{data=nil}	"Invalid request"
//	@Failure		500		{object}	localization.StandardResponse{data=nil}	"Server error"
//	@Security		BearerAuth
//	@Router			/account_block/districts/disable [post]
func (a *accountBlockAdapter) DisableDistricts(w http.ResponseWriter, r *http.Request) {
	var req accountblock.EnableOrDisableDistricts
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		localization.SendBadRequestResponse(w, localization.ErrorInvalidRequestBody.Code)
		return
	}

	// Trim any whitespace
	req.Clean()

	if err := req.Validate(); err != nil {
		localization.SendBadRequestResponse(w, err.Error())
		return
	}

	if len(req.DistrictIds) == 0 {
		localization.SendBadRequestResponse(w, localization.ErrorDistrictCodeRequired.Code)
		return
	}

	err := a.accountBlockApplication.EnableOrDisableDistricts(r.Context(), req.DistrictIds, req.Reason, false)
	if err != nil {
		localization.SendErrorByCodeResponse(w, err.Error())
		return
	}

	localization.SendSuccessResponse(w, localization.SuccessDisableDistrictsRequestSent, nil)
}

// EnableCities godoc
//
//	@Summary		Enable multiple cities
//	@Description	Enable multiple cities by codes.
//	@Tags			Account Block - Cities
//	@Accept			json
//	@Produce		json
//	@Param			request	body		accountblock.EnableOrDisableCities		true	"City codes and reason"	example({"city_codes":["CT001","CT002"], "reason": "foo"})
//	@Success		200		{object}	localization.StandardResponse{data=nil}	"Enable request submitted"
//	@Failure		400		{object}	localization.StandardResponse{data=nil}	"Invalid request"
//	@Failure		500		{object}	localization.StandardResponse{data=nil}	"Server error"
//	@Security		BearerAuth
//	@Router			/account_block/cities/enable [post]
func (a *accountBlockAdapter) EnableCities(w http.ResponseWriter, r *http.Request) {
	var req accountblock.EnableOrDisableCities
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		localization.SendBadRequestResponse(w, localization.ErrorInvalidRequestBody.Code)
		return
	}

	// Trim any whitespace
	req.Clean()

	if err := req.Validate(); err != nil {
		localization.SendBadRequestResponse(w, err.Error())
		return
	}

	if len(req.CityCodes) == 0 {
		localization.SendBadRequestResponse(w, localization.ErrorCityCodeRequired.Type)
		return
	}

	err := a.accountBlockApplication.EnableOrDisableCities(r.Context(), req.CityCodes, req.Reason, true)
	if err != nil {
		localization.SendErrorByCodeResponse(w, err.Error())
		return
	}

	localization.SendSuccessResponse(w, localization.SuccessEnableCitiesRequestSent, nil)
}

// DisableCities godoc
//
//	@Summary		Disable multiple cities
//	@Description	Disable multiple cities by codes.
//	@Tags			Account Block - Cities
//	@Accept			json
//	@Produce		json
//	@Param			request	body		accountblock.EnableOrDisableCities		true	"City codes and reason"	example({"city_codes":["CT001","CT002"], "reason": "foo"})
//	@Success		200		{object}	localization.StandardResponse{data=nil}	"Cities disabled successfully"
//	@Failure		400		{object}	localization.StandardResponse{data=nil}	"Bad request"
//	@Failure		500		{object}	localization.StandardResponse{data=nil}	"Internal server error"
//	@Security		BearerAuth
//	@Router			/account_block/cities/disable [post]
func (a *accountBlockAdapter) DisableCities(w http.ResponseWriter, r *http.Request) {
	var req accountblock.EnableOrDisableCities
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		localization.SendBadRequestResponse(w, localization.ErrorInvalidRequestBody.Code)
		return
	}

	// Trim any whitespace
	req.Clean()

	if err := req.Validate(); err != nil {
		localization.SendBadRequestResponse(w, err.Error())
		return
	}

	if len(req.CityCodes) == 0 {
		localization.SendBadRequestResponse(w, localization.ErrorCityCodeRequired.Type)
		return
	}

	err := a.accountBlockApplication.EnableOrDisableCities(r.Context(), req.CityCodes, req.Reason, false)
	if err != nil {
		localization.SendErrorByCodeResponse(w, err.Error())
		return
	}

	localization.SendSuccessResponse(w, localization.SuccessDisableCitiesRequestSent, nil)
}
