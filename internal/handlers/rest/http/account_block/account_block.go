package accountblock

import (
	"cbe-super-app-cps-action/internal/constants"
	"cbe-super-app-cps-action/internal/handlers/rest/http/account_block/core"
	"cbe-super-app-cps-action/internal/service"
	"context"
	"encoding/json"
	"net/http"
	"strings"

	accountblock "cbe-super-app-cps-action/internal/constants/dto/account_block"
	ab_interface "cbe-super-app-cps-action/internal/constants/interfaces/account_block"
	"cbe-super-app-cps-action/internal/constants/localization"
	"cbe-super-app-cps-action/internal/constants/types"

	"go.opentelemetry.io/otel/attribute"

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
func (a *accountBlockAdapter) GetBranchById(w http.ResponseWriter, r *http.Request) {
	ctx, span := local_util.TraceLogger(r.Context(), "handler", "getBranchById", "handler", "accountBlock")
	defer span.End()
	log := local_util.LoggerFromCtx(ctx, a.logger)
	branchId, ok := local_util.GetParam(r, "branch_id")
	if !ok {
		span.AddEvent("missing branch_id param")
		log.Errorf("[AccBlockH][GetBranchById] missing param")
		localization.SendErrorByCodeResponse(w, localization.ErrorBranchCodeRequired.Code)
		return
	}
	if branchId == "" {
		log.Errorf("[AccBlockH] invalid id: %s", branchId)
		localization.SendBadRequestResponse(w, "invalid id")
		return
	}

	span.SetAttributes(attribute.String("account_block.branch_id", branchId))

	branch, err := a.accountBlockApplication.GetBranchById(ctx, branchId)
	if err != nil {

		span.RecordError(err)
		log.Errorf("[AccBlockH][GetBranchById] svc err: %v", err)
		localization.SendErrorByCodeResponse(w, err.Error())
		return
	}

	data := core.ToAccountBlockResponse(branch)
	log.Infof("[AccBlockH][GetBranchById] ok: %s", branchId)
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
	ctx, span := local_util.TraceLogger(r.Context(), "handler", "getAllBranches", "handler", "accountBlock")
	defer span.End()
	log := local_util.LoggerFromCtx(ctx, a.logger)

	filterParams, err := local_util.ExtractFilterParams(r)
	if err != nil {
		localization.SendErrorByCodeResponse(w, err.Error())
		return
	}

	search := r.URL.Query().Get("search")
	filter := r.URL.Query().Get("filter")

	if err := local_util.NoSpecialChars(search); err != nil {
		localization.SendErrorByCodeResponse(w, err.Error())
		return
	}

	if err := local_util.NoSpecialChars(filter); err != nil {
		localization.SendErrorByCodeResponse(w, err.Error())
		return
	}

	branches, err := a.accountBlockApplication.GetAllBranches(ctx, filterParams)
	if err != nil {
		span.RecordError(err)
		log.Errorf("[AccBlockH][GetAllBranches] svc err: %v", err)
		localization.SendErrorByCodeResponse(w, err.Error())
		return
	}

	span.SetAttributes(attribute.Int("account_block.branches.count", len(branches.Data)))
	log.Infof("[AccBlockH][GetAllBranches] count: %d", len(branches.Data))
	localization.SendSuccessResponse(w, localization.SuccessBranchesRetrieved, branches)
}

// GetRegionById godoc
//
//	@Summary		Get region by federal region name
//	@Description	Retrieve a specific region by its federal region name.
//	@Tags			Account Block - Regions
//	@Accept			json
//	@Produce		json
//	@Param			region_id	path		string																	true	"Federal Region Name"	example(Oromia)
//	@Success		200			{object}	localization.StandardResponse{data=accountblock.AccountBlockResponse}	"Region retrieved"
//	@Failure		400			{object}	localization.StandardResponse{data=nil}									"Bad request"
//	@Failure		404			{object}	localization.StandardResponse{data=nil}									"Not found"
//	@Failure		500			{object}	localization.StandardResponse{data=nil}									"Server error"
//	@Security		BearerAuth
//	@Router			/account_block/regions/{region_id} [get]
func (a *accountBlockAdapter) GetRegionById(w http.ResponseWriter, r *http.Request) {
	ctx, span := local_util.TraceLogger(r.Context(), "handler", "getRegionById", "handler", "accountBlock")
	defer span.End()
	log := local_util.LoggerFromCtx(ctx, a.logger)

	regionName, ok := local_util.GetParam(r, "region_id")
	if !ok {
		log.Errorf("[AccBlockH][GetRegionById] missing param")
		localization.SendErrorByCodeResponse(w, localization.ErrorRegionCodeRequired.Code)
		return
	}
	if regionName == "" {
		log.Errorf("[AccBlockH] invalid federal region name: %s", regionName)
		localization.SendBadRequestResponse(w, "invalid federal region name")
		return
	}

	span.SetAttributes(attribute.String("account_block.federal_region_name", regionName))

	region, err := a.accountBlockApplication.GetRegionById(ctx, regionName)
	if err != nil {
		span.RecordError(err)
		log.Errorf("[AccBlockH][GetRegionById] svc err: %v", err)
		localization.SendErrorByCodeResponse(w, err.Error())
		return
	}

	data := core.ToAccountBlockResponse(region)
	log.Infof("[AccBlockH][GetRegionById] ok: %s", regionName)
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
	ctx, span := local_util.TraceLogger(r.Context(), "handler", "getAllRegions", "handler", "accountBlock")
	defer span.End()
	log := local_util.LoggerFromCtx(ctx, a.logger)

	filterParams, err := local_util.ExtractFilterParams(r)
	if err != nil {
		localization.SendErrorByCodeResponse(w, err.Error())
		return
	}

	search := r.URL.Query().Get("search")
	filter := r.URL.Query().Get("filter")

	if err := local_util.NoSpecialChars(search); err != nil {
		localization.SendErrorByCodeResponse(w, err.Error())
		return
	}

	if err := local_util.NoSpecialChars(filter); err != nil {
		localization.SendErrorByCodeResponse(w, err.Error())
		return
	}

	regions, err := a.accountBlockApplication.GetAllRegions(ctx, filterParams)
	if err != nil {
		span.RecordError(err)
		log.Errorf("[AccBlockH][GetAllRegions] svc err: %v", err)
		localization.SendErrorByCodeResponse(w, err.Error())
		return
	}

	span.SetAttributes(attribute.Int("account_block.regions.count", len(regions.Data)))
	log.Infof("[AccBlockH][GetAllRegions] count: %d", len(regions.Data))
	localization.SendSuccessResponse(w, localization.SuccessRegionsRetrieved, regions)
}

// GetDistrictById godoc
//
//	@Summary		Get district by name
//	@Description	Retrieve a specific district by its district name.
//	@Tags			Account Block - Districts
//	@Accept			json
//	@Produce		json
//	@Param			district_id	path		string																	true	"District Name"	example(Bole)
//	@Success		200				{object}	localization.StandardResponse{data=accountblock.AccountBlockResponse}	"District retrieved successfully"
//	@Failure		400				{object}	localization.StandardResponse{data=nil}									"Bad request"
//	@Failure		404				{object}	localization.StandardResponse{data=nil}									"Not found"
//	@Failure		500				{object}	localization.StandardResponse{data=nil}									"internal Server error"
//	@Security		BearerAuth
//	@Router			/account_block/districts/{district_id} [get]
func (a *accountBlockAdapter) GetDistrictById(w http.ResponseWriter, r *http.Request) {
	ctx, span := local_util.TraceLogger(r.Context(), "handler", "getDistrictById", "handler", "accountBlock")
	defer span.End()
	log := local_util.LoggerFromCtx(ctx, a.logger)

	districtName, ok := local_util.GetParam(r, "district_id")
	if !ok {
		log.Errorf("[AccBlockH][GetDistrictById] missing param")
		localization.SendErrorByCodeResponse(w, localization.ErrorDistrictCodeRequired.Code)
		return
	}
	if districtName == "" {
		log.Errorf("[AccBlockH] invalid district name: %s", districtName)
		localization.SendBadRequestResponse(w, "invalid district name")
		return
	}

	span.SetAttributes(attribute.String("account_block.district_name", districtName))

	district, err := a.accountBlockApplication.GetDistrictById(ctx, districtName)
	if err != nil {
		span.RecordError(err)
		log.Errorf("[AccBlockH][GetDistrictById] svc err: %v", err)
		localization.SendErrorByCodeResponse(w, err.Error())
		return
	}

	data := core.ToAccountBlockResponse(district)
	log.Infof("[AccBlockH][GetDistrictById] ok: %s", districtName)
	localization.SendSuccessResponse(w, localization.SuccessDistrictRetrieved, data)
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
//	@Param			filters		query		string																false	"JSON encoded filters (enabled, region)"	example("{\"enabled\":true,\"region_id\":\"Oromia\"}")
//	@Success		200			{object}	localization.StandardResponse{data=paginated_account_block_resp}	"Districts retrieved successfully"
//	@Failure		500			{object}	localization.StandardResponse{data=nil}								"internal Server error"
//	@Security		BearerAuth
//	@Router			/account_block/districts [get]
func (a *accountBlockAdapter) GetAllDistricts(w http.ResponseWriter, r *http.Request) {
	ctx, span := local_util.TraceLogger(r.Context(), "handler", "getAllDistricts", "handler", "accountBlock")
	defer span.End()
	log := local_util.LoggerFromCtx(ctx, a.logger)

	filterParams, err := local_util.ExtractFilterParams(r)
	if err != nil {
		localization.SendErrorByCodeResponse(w, err.Error())
		return
	}

	search := r.URL.Query().Get("search")
	filter := r.URL.Query().Get("filter")

	if err := local_util.NoSpecialChars(search); err != nil {
		localization.SendErrorByCodeResponse(w, err.Error())
		return
	}

	if err := local_util.NoSpecialChars(filter); err != nil {
		localization.SendErrorByCodeResponse(w, err.Error())
		return
	}

	districts, err := a.accountBlockApplication.GetAllDistricts(ctx, filterParams)
	if err != nil {
		span.RecordError(err)
		log.Errorf("[AccBlockH][GetAllDistricts] svc err: %v", err)
		localization.SendErrorByCodeResponse(w, err.Error())
		return
	}

	span.SetAttributes(attribute.Int("account_block.districts.count", len(districts.Data)))
	log.Infof("[AccBlockH][GetAllDistricts] count: %d", len(districts.Data))
	localization.SendSuccessResponse(w, localization.SuccessDistrictsRetrieved, districts)
}

// EnableBranches godoc
//
//	@Summary		Enable multiple branches
//	@Description	Enable multiple branches by codes.
//	@Tags			Account Block - Branches
//	@Accept			json
//	@Produce		json
//	@Param			request	body		accountblock.EnableOrDisableBranches	true	"Branch IDs and reason"	example({"branch_ids":["674003000000000000000001","674003000000000000000002"], "reason": "foo"})
//	@Success		200		{object}	localization.StandardResponse{data=nil}	"Enable request submitted"
//	@Failure		400		{object}	localization.StandardResponse{data=nil}	"Invalid request"
//	@Failure		500		{object}	localization.StandardResponse{data=nil}	"Server error"
//	@Security		BearerAuth
//	@Router			/account_block/branches/enable [post]
func (a *accountBlockAdapter) EnableBranches(w http.ResponseWriter, r *http.Request) {
	ctx, span := local_util.TraceLogger(r.Context(), "handler", "enableBranches", "handler", "accountBlock")
	defer span.End()
	log := local_util.LoggerFromCtx(ctx, a.logger)
	md := &types.ContextMetadata{}
	ctx = context.WithValue(ctx, constants.ContextKeyMetadata, md)
	localization.UpdateWriterContext(w, ctx)
	var req accountblock.EnableOrDisableBranches
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		span.RecordError(err)
		log.Errorf("[AccBlockH] decode body err: %v", err)
		localization.SendBadRequestResponse(w, localization.ErrorInvalidRequestBody.Code)
		return
	}

	// Trim any whitespace
	req.Clean()

	if err := req.Validate(); err != nil {
		span.RecordError(err)
		log.Errorf("[AccBlockH] validate err: %v", err)
		localization.SendBadRequestResponse(w, err.Error())
		return
	}

	if len(req.BranchIds) == 0 {
		span.AddEvent("missing branch ids")
		log.Errorf("[AccBlockH] branch ids required")
		localization.SendErrorByCodeResponse(w, localization.ErrorBranchCodeRequired.Code)
		return
	}

	span.SetAttributes(
		attribute.Int("account_block.branches.count", len(req.BranchIds)),
		attribute.String("account_block.reason", req.Reason),
	)

	err := a.accountBlockApplication.EnableOrDisableBranches(ctx, req.BranchIds, req.Reason, true)
	if err != nil {
		w = local_util.HandlePendingResponseError(ctx, w, err)
		span.RecordError(err)
		log.Errorf("[AccBlockH][EnableBranches] svc err: %v", err)
		localization.SendErrorByCodeResponse(w, err.Error())
		return
	}

	if md.IsMakerOnly {
		userCode, _ := ctx.Value(constants.ContextKey("user_code")).(string)
		log.Infof("[AccBlockH][EnableBranches] ok user: %s maker_only: %v", userCode, md.IsMakerOnly)
		w = localization.ApplyActionCodeHeaderFromWriter(w, ctx)
		localization.SendSuccessResponse(w, localization.SuccessBranchEnabled, nil)
		return
	}

	log.Infof("[AccBlockH][EnableBranches] ok count: %d", len(req.BranchIds))
	w = localization.ApplyActionCodeHeaderFromWriter(w, ctx)
	localization.SendSuccessResponse(w, localization.SuccessEnableBranchesRequestSent, nil)
}

// DisableBranches godoc
//
//	@Summary		Disable multiple branches
//	@Description	Disable multiple branches by codes.
//	@Tags			Account Block - Branches
//	@Accept			json
//	@Produce		json
//	@Param			request	body		accountblock.EnableOrDisableBranches	true	"Branch IDs and reason"	example({"branch_ids":["674003000000000000000001","674003000000000000000002"], "reason": "foo"})
//	@Success		200		{object}	localization.StandardResponse{data=nil}	"Disable request submitted"
//	@Failure		400		{object}	localization.StandardResponse{data=nil}	"Invalid request"
//	@Failure		500		{object}	localization.StandardResponse{data=nil}	"Server error"
//	@Security		BearerAuth
//	@Router			/account_block/branches/disable [post]
func (a *accountBlockAdapter) DisableBranches(w http.ResponseWriter, r *http.Request) {
	ctx, span := local_util.TraceLogger(r.Context(), "handler", "disableBranches", "handler", "accountBlock")
	defer span.End()
	log := local_util.LoggerFromCtx(ctx, a.logger)
	md := &types.ContextMetadata{}
	ctx = context.WithValue(ctx, constants.ContextKeyMetadata, md)
	localization.UpdateWriterContext(w, ctx)
	var req accountblock.EnableOrDisableBranches
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		span.RecordError(err)
		log.Errorf("[AccBlockH] decode body err: %v", err)
		localization.SendErrorByCodeResponse(w, localization.ErrorInvalidRequestBody.Code)
		return
	}

	// Trim any whitespace
	req.Clean()

	if err := req.Validate(); err != nil {
		span.RecordError(err)
		log.Errorf("[AccBlockH] validate err: %v", err)
		localization.SendBadRequestResponse(w, err.Error())
		return
	}

	if len(req.BranchIds) == 0 {
		span.AddEvent("missing branch ids")
		log.Errorf("[AccBlockH] branch ids required")
		localization.SendErrorByCodeResponse(w, localization.ErrorBranchCodeRequired.Code)
		return
	}

	span.SetAttributes(
		attribute.Int("account_block.branches.count", len(req.BranchIds)),
		attribute.String("account_block.reason", req.Reason),
	)

	err := a.accountBlockApplication.EnableOrDisableBranches(ctx, req.BranchIds, req.Reason, false)
	if err != nil {
		w = local_util.HandlePendingResponseError(ctx, w, err)
		span.RecordError(err)
		log.Errorf("[AccBlockH][DisableBranches] svc err: %v", err)
		localization.SendErrorByCodeResponse(w, err.Error())
		return
	}

	if md.IsMakerOnly {
		userCode, _ := ctx.Value(constants.ContextKey("user_code")).(string)
		log.Infof("[AccBlockH][DisableBranches] ok user: %s maker_only: %v", userCode, md.IsMakerOnly)
		w = localization.ApplyActionCodeHeaderFromWriter(w, ctx)
		localization.SendSuccessResponse(w, localization.SuccessDisableBranches, nil)
		return
	}

	log.Infof("[AccBlockH][DisableBranches] ok count: %d", len(req.BranchIds))
	w = localization.ApplyActionCodeHeaderFromWriter(w, ctx)
	localization.SendSuccessResponse(w, localization.SuccessDisableBranchesRequestSent, nil)
}

// EnableRegions godoc
//
//	@Summary		Enable multiple regions
//	@Description	Enable multiple regions by federal region name.
//	@Tags			Account Block - Regions
//	@Accept			json
//	@Produce		json
//	@Param			request	body		accountblock.EnableOrDisableRegions		true	"Federal region names and reason"	example({"federal_region_names":["Oromia","Amhara"], "reason": "foo"})
//	@Success		200		{object}	localization.StandardResponse{data=nil}	"Enable request submitted"
//	@Failure		400		{object}	localization.StandardResponse{data=nil}	"Invalid request"
//	@Failure		500		{object}	localization.StandardResponse{data=nil}	"Server error"
//	@Security		BearerAuth
//	@Router			/account_block/regions/enable [post]
func (a *accountBlockAdapter) EnableRegions(w http.ResponseWriter, r *http.Request) {
	ctx, span := local_util.TraceLogger(r.Context(), "handler", "enableRegions", "handler", "accountBlock")
	defer span.End()
	log := local_util.LoggerFromCtx(ctx, a.logger)
	md := &types.ContextMetadata{}
	ctx = context.WithValue(ctx, constants.ContextKeyMetadata, md)
	localization.UpdateWriterContext(w, ctx)
	var req accountblock.EnableOrDisableRegions
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		span.RecordError(err)
		log.Errorf("[AccBlockH] decode body err: %v", err)
		localization.SendErrorByCodeResponse(w, localization.ErrorInvalidRequestBody.Code)
		return
	}

	// Trim any whitespace
	req.Clean()

	if err := req.Validate(); err != nil {
		span.RecordError(err)
		log.Errorf("[AccBlockH] validate err: %v", err)
		localization.SendBadRequestResponse(w, err.Error())
		return
	}

	if len(req.FederalRegionNames) == 0 {
		span.AddEvent("missing federal region names")
		log.Errorf("[AccBlockH] federal region names required")
		localization.SendErrorByCodeResponse(w, localization.ErrorRegionCodeRequired.Code)
		return
	}

	span.SetAttributes(
		attribute.Int("account_block.regions.count", len(req.FederalRegionNames)),
		attribute.String("account_block.reason", req.Reason),
	)

	err := a.accountBlockApplication.EnableOrDisableRegions(ctx, req.FederalRegionNames, req.Reason, true)
	if err != nil {
		w = local_util.HandlePendingResponseError(ctx, w, err)
		span.RecordError(err)
		log.Errorf("[AccBlockH][EnableRegions] svc err: %v", err)
		localization.SendErrorByCodeResponse(w, err.Error())
		return
	}

	if md.IsMakerOnly {
		userCode, _ := ctx.Value(constants.ContextKey("user_code")).(string)
		log.Infof("[AccBlockH][EnableRegions] ok user: %s maker_only: %v", userCode, md.IsMakerOnly)
		w = localization.ApplyActionCodeHeaderFromWriter(w, ctx)
		localization.SendSuccessResponse(w, localization.SuccessEnableRegion, nil)
		return
	}

	log.Infof("[AccBlockH][EnableRegions] ok count: %d", len(req.FederalRegionNames))
	w = localization.ApplyActionCodeHeaderFromWriter(w, ctx)
	localization.SendSuccessResponse(w, localization.SuccessEnableRegionsRequestSent, nil)
}

// DisableRegions godoc
//
//	@Summary		Disable multiple regions
//	@Description	Disable multiple regions by federal region name.
//	@Tags			Account Block - Regions
//	@Accept			json
//	@Produce		json
//	@Param			request	body		accountblock.EnableOrDisableRegions		true	"Federal region names and reason"	example({"federal_region_names":["Oromia","Amhara"], "reason": "foo"})
//	@Success		200		{object}	localization.StandardResponse{data=nil}	"Disable request submitted"
//	@Failure		400		{object}	localization.StandardResponse{data=nil}	"Invalid request"
//	@Failure		500		{object}	localization.StandardResponse{data=nil}	"Server error"
//	@Security		BearerAuth
//	@Router			/account_block/regions/disable [post]
func (a *accountBlockAdapter) DisableRegions(w http.ResponseWriter, r *http.Request) {
	ctx, span := local_util.TraceLogger(r.Context(), "handler", "disableRegions", "handler", "accountBlock")
	defer span.End()
	log := local_util.LoggerFromCtx(ctx, a.logger)
	md := &types.ContextMetadata{}
	ctx = context.WithValue(ctx, constants.ContextKeyMetadata, md)
	localization.UpdateWriterContext(w, ctx)
	var req accountblock.EnableOrDisableRegions
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		span.RecordError(err)
		log.Errorf("[AccBlockH] decode body err: %v", err)
		localization.SendBadRequestResponse(w, localization.ErrorInvalidRequestBody.Code)
		return

	}

	// Trim any whitespace
	req.Clean()

	if err := req.Validate(); err != nil {
		span.RecordError(err)
		log.Errorf("[AccBlockH] validate err: %v", err)
		localization.SendBadRequestResponse(w, err.Error())
		return
	}

	if len(req.FederalRegionNames) == 0 {
		span.AddEvent("missing federal region names")
		log.Errorf("[AccBlockH] federal region names required")
		localization.SendErrorByCodeResponse(w, localization.ErrorRegionCodeRequired.Code)
		return
	}

	span.SetAttributes(
		attribute.Int("account_block.regions.count", len(req.FederalRegionNames)),
		attribute.String("account_block.reason", req.Reason),
	)

	err := a.accountBlockApplication.EnableOrDisableRegions(ctx, req.FederalRegionNames, req.Reason, false)
	if err != nil {
		w = local_util.HandlePendingResponseError(ctx, w, err)
		span.RecordError(err)
		log.Errorf("[AccBlockH][DisableRegions] svc err: %v", err)
		localization.SendErrorByCodeResponse(w, err.Error())
		return
	}

	if md.IsMakerOnly {
		userCode, _ := ctx.Value(constants.ContextKey("user_code")).(string)
		log.Infof("[AccBlockH][DisableRegions] ok user: %s maker_only: %v", userCode, md.IsMakerOnly)
		w = localization.ApplyActionCodeHeaderFromWriter(w, ctx)
		localization.SendSuccessResponse(w, localization.SuccessDisableRegion, nil)
		return
	}

	w = localization.ApplyActionCodeHeaderFromWriter(w, ctx)
	localization.SendSuccessResponse(w, localization.SuccessDisableRegionsRequestSent, nil)
}

// EnableDistricts godoc
//
//	@Summary		Enable multiple districts
//	@Description	Enable multiple districts by district name.
//	@Tags			Account Block - Districts
//	@Accept			json
//	@Produce		json
//	@Param			request	body		accountblock.EnableOrDisableDistricts	true	"District names and reason"	example({"district_names":["Bole","Arada"], "reason": "foo"})
//	@Success		200		{object}	localization.StandardResponse{data=nil}	"Enable request submitted"
//	@Failure		400		{object}	localization.StandardResponse{data=nil}	"Invalid request"
//	@Failure		500		{object}	localization.StandardResponse{data=nil}	"Server error"
//	@Security		BearerAuth
//	@Router			/account_block/districts/enable [post]
func (a *accountBlockAdapter) EnableDistricts(w http.ResponseWriter, r *http.Request) {
	ctx, span := local_util.TraceLogger(r.Context(), "handler", "enableDistricts", "handler", "accountBlock")
	defer span.End()
	log := local_util.LoggerFromCtx(ctx, a.logger)
	md := &types.ContextMetadata{}
	ctx = context.WithValue(ctx, constants.ContextKeyMetadata, md)
	localization.UpdateWriterContext(w, ctx)
	var req accountblock.EnableOrDisableDistricts
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		span.RecordError(err)
		log.Errorf("[AccBlockH] decode body err: %v", err)
		localization.SendBadRequestResponse(w, localization.ErrorInvalidRequestBody.Code)
		return
	}

	// Trim any whitespace
	req.Clean()

	if err := req.Validate(); err != nil {
		span.RecordError(err)
		log.Errorf("[AccBlockH] validate err: %v", err)
		localization.SendBadRequestResponse(w, err.Error())
		return
	}

	if len(req.DistrictNames) == 0 {
		span.AddEvent("missing district names")
		log.Errorf("[AccBlockH] district names required")
		localization.SendErrorByCodeResponse(w, localization.ErrorDistrictCodeRequired.Code)
		return
	}

	span.SetAttributes(
		attribute.Int("account_block.districts.count", len(req.DistrictNames)),
		attribute.String("account_block.reason", req.Reason),
	)

	err := a.accountBlockApplication.EnableOrDisableDistricts(ctx, req.DistrictNames, req.Reason, true)
	if err != nil {
		w = local_util.HandlePendingResponseError(ctx, w, err)
		span.RecordError(err)
		log.Errorf("[AccBlockH][EnableDistricts] svc err: %v", err)
		localization.SendErrorByCodeResponse(w, err.Error())
		return
	}

	if md.IsMakerOnly {
		userCode, _ := ctx.Value(constants.ContextKey("user_code")).(string)
		log.Infof("[AccBlockH][EnableDistricts] ok user: %s maker_only: %v", userCode, md.IsMakerOnly)
		w = localization.ApplyActionCodeHeaderFromWriter(w, ctx)
		localization.SendSuccessResponse(w, localization.SuccessEnableDistricts, nil)
		return
	}

	w = localization.ApplyActionCodeHeaderFromWriter(w, ctx)
	localization.SendSuccessResponse(w, localization.SuccessEnableDistrictsRequestSent, nil)
}

// DisableDistricts godoc
//
//	@Summary		Disable multiple districts
//	@Description	Disable multiple districts by district name.
//	@Tags			Account Block - Districts
//	@Accept			json
//	@Produce		json
//	@Param			request	body		accountblock.EnableOrDisableDistricts	true	"District names and reason"	example({"district_names":["Bole","Arada"], "reason": "foo"})
//	@Success		200		{object}	localization.StandardResponse{data=nil}	"Disable request submitted"
//	@Failure		400		{object}	localization.StandardResponse{data=nil}	"Invalid request"
//	@Failure		500		{object}	localization.StandardResponse{data=nil}	"Server error"
//	@Security		BearerAuth
//	@Router			/account_block/districts/disable [post]
func (a *accountBlockAdapter) DisableDistricts(w http.ResponseWriter, r *http.Request) {
	ctx, span := local_util.TraceLogger(r.Context(), "handler", "disableDistricts", "handler", "accountBlock")
	defer span.End()
	log := local_util.LoggerFromCtx(ctx, a.logger)
	md := &types.ContextMetadata{}
	ctx = context.WithValue(ctx, constants.ContextKeyMetadata, md)
	localization.UpdateWriterContext(w, ctx)
	var req accountblock.EnableOrDisableDistricts
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		span.RecordError(err)
		localization.SendBadRequestResponse(w, localization.ErrorInvalidRequestBody.Code)
		return
	}

	// Trim any whitespace
	req.Clean()

	if err := req.Validate(); err != nil {
		span.RecordError(err)
		localization.SendBadRequestResponse(w, err.Error())
		return
	}

	if len(req.DistrictNames) == 0 {
		span.AddEvent("missing district names")
		localization.SendErrorByCodeResponse(w, localization.ErrorDistrictCodeRequired.Code)
		return
	}

	span.SetAttributes(
		attribute.Int("account_block.districts.count", len(req.DistrictNames)),
		attribute.String("account_block.reason", req.Reason),
	)

	err := a.accountBlockApplication.EnableOrDisableDistricts(ctx, req.DistrictNames, req.Reason, false)
	if err != nil {
		w = local_util.HandlePendingResponseError(ctx, w, err)
		span.RecordError(err)
		localization.SendErrorByCodeResponse(w, err.Error())
		return
	}

	if md.IsMakerOnly {
		userCode, _ := ctx.Value(constants.ContextKey("user_code")).(string)
		log.Infof("[AccBlockH][DisableDistricts] ok user: %s maker_only: %v", userCode, md.IsMakerOnly)
		w = localization.ApplyActionCodeHeaderFromWriter(w, ctx)
		localization.SendSuccessResponse(w, localization.SuccessDisableDistricts, nil)
		return
	}

	w = localization.ApplyActionCodeHeaderFromWriter(w, ctx)
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
	ctx, span := local_util.TraceLogger(r.Context(), "handler", "enableCities", "handler", "accountBlock")
	defer span.End()
	log := local_util.LoggerFromCtx(ctx, a.logger)
	md := &types.ContextMetadata{}
	ctx = context.WithValue(ctx, constants.ContextKeyMetadata, md)
	localization.UpdateWriterContext(w, ctx)
	var req accountblock.EnableOrDisableCities
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		span.RecordError(err)
		localization.SendBadRequestResponse(w, localization.ErrorInvalidRequestBody.Code)
		return
	}

	// Trim any whitespace
	req.Clean()

	if err := req.Validate(); err != nil {
		span.RecordError(err)
		localization.SendBadRequestResponse(w, err.Error())
		return
	}

	if len(req.CityIds) == 0 {
		span.AddEvent("missing city ids")
		localization.SendErrorByCodeResponse(w, localization.ErrorCityCodeRequired.Code)
		return
	}

	span.SetAttributes(
		attribute.Int("account_block.cities.count", len(req.CityIds)),
		attribute.String("account_block.reason", req.Reason),
	)

	// err := a.accountBlockApplication.EnableOrDisableCities(ctx, req.CityIds, req.Reason, true)
	// if err != nil {
	// 	span.RecordError(err)
	// 	localization.SendErrorByCodeResponse(w, err.Error())
	// 	return
	// }

	if md.IsMakerOnly {
		userCode, _ := ctx.Value(constants.ContextKey("user_code")).(string)
		log.Infof("[AccBlockH][EnableCities] ok user: %s maker_only: %v", userCode, md.IsMakerOnly)
		w = localization.ApplyActionCodeHeaderFromWriter(w, ctx)
		localization.SendSuccessResponse(w, localization.SuccessEnableCities, nil)
		return
	}

	w = localization.ApplyActionCodeHeaderFromWriter(w, ctx)
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
	ctx, span := local_util.TraceLogger(r.Context(), "handler", "disableCities", "handler", "accountBlock")
	defer span.End()
	log := local_util.LoggerFromCtx(ctx, a.logger)
	md := &types.ContextMetadata{}
	ctx = context.WithValue(ctx, constants.ContextKeyMetadata, md)
	localization.UpdateWriterContext(w, ctx)
	var req accountblock.EnableOrDisableCities
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		span.RecordError(err)
		localization.SendBadRequestResponse(w, localization.ErrorInvalidRequestBody.Code)
		return
	}

	// Trim any whitespace
	req.Clean()

	if err := req.Validate(); err != nil {
		span.RecordError(err)
		localization.SendBadRequestResponse(w, err.Error())
		return
	}

	if len(req.CityIds) == 0 {
		span.AddEvent("missing city ids")
		localization.SendBadRequestResponse(w, localization.ErrorCityCodeRequired.Code)
		return
	}

	span.SetAttributes(
		attribute.Int("account_block.cities.count", len(req.CityIds)),
		attribute.String("account_block.reason", req.Reason),
	)

	// err := a.accountBlockApplication.EnableOrDisableCities(ctx, req.CityIds, req.Reason, false)
	// if err != nil {
	// 	span.RecordError(err)
	// 	localization.SendErrorByCodeResponse(w, err.Error())
	// 	return
	// }
	if md.IsMakerOnly {
		userCode, _ := ctx.Value(constants.ContextKey("user_code")).(string)
		log.Infof("[AccBlockH][DisableCities] ok user: %s maker_only: %v", userCode, md.IsMakerOnly)
		w = localization.ApplyActionCodeHeaderFromWriter(w, ctx)
		localization.SendSuccessResponse(w, localization.SuccessDisableCities, nil)
		return
	}

	w = localization.ApplyActionCodeHeaderFromWriter(w, ctx)
	localization.SendSuccessResponse(w, localization.SuccessDisableCitiesRequestSent, nil)
}

func (a *accountBlockAdapter) GetAccountBlockDetails(w http.ResponseWriter, r *http.Request) {
	ctx, span := local_util.TraceLogger(r.Context(), "handler", "GetAccountBlockDetails", "handler", "accountBlock")
	defer span.End()
	log := local_util.LoggerFromCtx(ctx, a.logger)

	id, ok := local_util.GetParam(r, "id")
	if !ok {
		log.Errorf("[AccBlockH][GetDetails] missing param")
		localization.SendErrorByCodeResponse(w, localization.ErrorCodeRequired.Code)
		return
	}
	if id == "" {
		log.Errorf("[AccBlockH] invalid id: %s", id)
		localization.SendBadRequestResponse(w, "invalid id")
		return
	}

	filterParams, err := local_util.ExtractFilterParams(r)
	if err != nil {
		localization.SendErrorByCodeResponse(w, err.Error())
		return

	}

	data, err := a.accountBlockApplication.GetAccountBlockDetails(ctx, id, filterParams)
	if err != nil {
		span.RecordError(err)
		log.Errorf("[AccBlockH][GetDetails] svc err: %v", err)
		localization.SendErrorByCodeResponse(w, err.Error())
		return
	}

	log.Infof("[AccBlockH][GetDetails] ok: %s", id)
	localization.SendSuccessResponse(w, localization.DataRetrievedSuccessfully, data)
}

// GetPreviousReasons godoc
//
//	@Summary		Get previous disable reasons for an account block
//	@Description	Returns the disable-reason history for a branch, district, or region. Use type=B with branch oracle ID, type=R with federal region name, or type=D with district name.
//	@Tags			Account Block
//	@Accept			json
//	@Produce		json
//	@Param			type		query		string																	true	"Entity type (B=branch, R=region, D=district)"	example(B)
//	@Param			identifier	query		string																	true	"Entity identifier (branch oracle ID or region/district name)"	example(3F4A2B1C...)
//	@Success		200			{object}	localization.StandardResponse{data=accountblock.PreviousDisableReasonsResponse}	"Disable reasons retrieved"
//	@Failure		400			{object}	localization.StandardResponse{data=nil}									"Bad request"
//	@Failure		404			{object}	localization.StandardResponse{data=nil}									"Not found"
//	@Failure		500			{object}	localization.StandardResponse{data=nil}									"Server error"
//	@Security		BearerAuth
//	@Router			/account_block/previous_reasons [get]
func (a *accountBlockAdapter) GetPreviousReasons(w http.ResponseWriter, r *http.Request) {
	ctx, span := local_util.TraceLogger(r.Context(), "handler", "GetPreviousReasons", "handler", "accountBlock")
	defer span.End()
	log := local_util.LoggerFromCtx(ctx, a.logger)

	entityType := strings.ToUpper(strings.TrimSpace(r.URL.Query().Get("type")))
	identifier := strings.TrimSpace(r.URL.Query().Get("identifier"))

	if entityType == "" {
		log.Errorf("[AccBlockH][GetPreviousReasons] missing type query param")
		localization.SendBadRequestResponse(w, "type is required")
		return
	}
	if identifier == "" {
		log.Errorf("[AccBlockH][GetPreviousReasons] missing identifier query param")
		localization.SendBadRequestResponse(w, "identifier is required")
		return
	}
	switch entityType {
	case "B", "R", "D":
	default:
		log.Errorf("[AccBlockH][GetPreviousReasons] invalid type: %s", entityType)
		localization.SendBadRequestResponse(w, "invalid type: must be B, R, or D")
		return
	}

	span.SetAttributes(
		attribute.String("account_block.entity_type", entityType),
		attribute.String("account_block.identifier", identifier),
	)

	data, err := a.accountBlockApplication.GetPreviousReasons(ctx, entityType, identifier)
	if err != nil {
		span.RecordError(err)
		log.Errorf("[AccBlockH][GetPreviousReasons] svc err: %v", err)
		localization.SendErrorByCodeResponse(w, err.Error())
		return
	}

	log.Infof("[AccBlockH][GetPreviousReasons] ok: type=%s identifier=%s count=%d", entityType, identifier, len(data.DisableReason))
	localization.SendSuccessResponse(w, localization.DataRetrievedSuccessfully, data)
}
