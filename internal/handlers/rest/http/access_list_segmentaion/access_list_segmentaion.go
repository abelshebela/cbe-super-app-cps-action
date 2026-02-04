package accesslistsegmentaion

import (
	"cbe-super-app-cps-action/internal/constants"
	access_list_segmentation_dto "cbe-super-app-cps-action/internal/constants/dto/access_list_segmentation"
	accesslistsegmentation "cbe-super-app-cps-action/internal/constants/interfaces/access_list_segmentation"
	"cbe-super-app-cps-action/internal/constants/localization"
	"cbe-super-app-cps-action/internal/constants/types"
	"cbe-super-app-cps-action/internal/service"
	local_util "cbe-super-app-cps-action/pkgs/utils"
	"context"
	"encoding/json"
	"net/http"

	"github.com/go-chi/chi/v5"
	"gitlab.com/bersufekadgetachew/cbe-super-app-shared/shared/utils"
)

type accessListSegmentation struct {
	service service.AccessListSegmentationService
	logger  utils.Logger
}

// CreateAccessListSegmentation godoc
//
//	@Summary		Create access list segmentation
//	@Description	Create a new access list segmentation with the provided information
//	@Tags			Access List Segmentation
//	@Accept			json
//	@Produce		json
//	@Param			body	body		access_list_segmentation_dto.CreateAccessListSegmentationRequest	true	"Create access list segmentation request"
//	@Success		200		{object}	localization.StandardResponse{data=nil}	"Access list segmentation created successfully"
//	@Failure		400		{object}	localization.StandardResponse{data=nil}	"Bad request"
//	@Failure		500		{object}	localization.StandardResponse{data=nil}	"Internal server error"
//	@Security		BearerAuth
//	@Router			/access_list_segmentation [post]
func (a *accessListSegmentation) CreateAccessListSegmentation(w http.ResponseWriter, r *http.Request) {
	ctx, span := local_util.TraceLogger(r.Context(), "handler", "CreateAccessListSegmentation", "handler", "accessListSegmentation")
	defer span.End()
	md := &types.ContextMetadata{}
	ctx = context.WithValue(ctx, constants.ContextKeyMetadata, md)

	var req access_list_segmentation_dto.CreateAccessListSegmentationRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		a.logger.Errorf("[CreateAccessListSegmentation] failed to decode request body: %v", err)
		localization.SendErrorByCodeResponse(w, err.Error())
		return
	}
	if err := req.Validate(); err != nil {
		a.logger.Errorf("[CreateAccessListSegmentation] validation error: %v", err)
		localization.SendErrorByCodeResponse(w, err.Error())
		return
	}
	if err := a.service.CreateAccessListSegmentation(ctx, req); err != nil {
		a.logger.Errorf("[CreateAccessListSegmentation] service error: %v", err)
		localization.SendErrorByCodeResponse(w, err.Error())
		return
	}
	if md.IsMakerOnly {
		localization.SendSuccessResponse(w, localization.SuccessAccessListSegmentationCreatedSP, nil)
		return
	}
	localization.SendSuccessResponse(w, localization.SuccessAccessListSegmentationCreated, nil)
}

// DisableAccessListSegmentation godoc
//
//	@Summary		Disable access list segmentation
//	@Description	Disable an access list segmentation by ID with optional access list keys
//	@Tags			Access List Segmentation
//	@Accept			json
//	@Produce		json
//	@Param			id		path		string															true	"Access list segmentation ID"
//	@Param			body	body		access_list_segmentation_dto.EnableDisableAccessListSegmentationRequest	true	"Disable request with access list keys"
//	@Success		200		{object}	localization.StandardResponse{data=nil}	"Access list segmentation disabled successfully"
//	@Failure		400		{object}	localization.StandardResponse{data=nil}	"Bad request"
//	@Failure		404		{object}	localization.StandardResponse{data=nil}	"Access list segmentation not found"
//	@Failure		500		{object}	localization.StandardResponse{data=nil}	"Internal server error"
//	@Security		BearerAuth
//	@Router			/access_list_segmentation/disable/{id} [patch]
func (a *accessListSegmentation) DisableAccessListSegmentation(w http.ResponseWriter, r *http.Request) {
	id := chi.URLParam(r, "id")
	if id == "" {
		a.logger.Errorf("[DisableAccessListSegmentation] missing id parameter")
		localization.SendErrorByCodeResponse(w, localization.ErrorAccessListSegmentationInvalidID.Code)
		return
	}
	var req access_list_segmentation_dto.EnableDisableAccessListSegmentationRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		a.logger.Errorf("[DisableAccessListSegmentation] failed to decode request body: %v", err)
		localization.SendErrorByCodeResponse(w, err.Error())
		return
	}

	if err := req.Validate(); err != nil {
		a.logger.Errorf("[DisableAccessListSegmentation] validation error: %v", err)
		localization.SendErrorByCodeResponse(w, err.Error())
		return
	}

	if err := a.service.EnableDisableAccessListSegmentation(r.Context(), id, false, req.AccessListKeys); err != nil {
		a.logger.Errorf("[DisableAccessListSegmentation] service error: %v", err)
		localization.SendErrorByCodeResponse(w, err.Error())
		return
	}

	localization.SendSuccessResponse(w, localization.SuccessAccessListSegmentationDisabled, nil)
}

// EnableAccessListSegmentation implements accesslistsegmentation.AccessListSegmentationHandler.
func (a *accessListSegmentation) EnableAccessListSegmentation(w http.ResponseWriter, r *http.Request) {
	id := chi.URLParam(r, "id")
	if id == "" {
		a.logger.Errorf("[EnableAccessListSegmentation] missing id parameter")
		localization.SendErrorByCodeResponse(w, localization.ErrorAccessListSegmentationInvalidID.Code)
		return
	}
	if err := a.service.EnableDisableAccessListSegmentation(r.Context(), id, true, nil); err != nil {
		a.logger.Errorf("[EnableAccessListSegmentation] service error: %v", err)
		localization.SendErrorByCodeResponse(w, err.Error())
		return
	}
	localization.SendSuccessResponse(w, localization.SuccessAccessListSegmentationEnabled, nil)
}

// GetAccessListSegmentationByID godoc
//
//	@Summary		Get access list segmentation by ID
//	@Description	Retrieve a single access list segmentation by its identifier
//	@Tags			Access List Segmentation
//	@Accept			json
//	@Produce		json
//	@Param			id	path		string																	true	"Access list segmentation ID"
//	@Success		200	{object}	localization.StandardResponse{data=access_list_segmentation_dto.AccessListSegmentationResponse}	"Access list segmentation retrieved successfully"
//	@Failure		400	{object}	localization.StandardResponse{data=nil}								"Bad request"
//	@Failure		404	{object}	localization.StandardResponse{data=nil}								"Access list segmentation not found"
//	@Failure		500	{object}	localization.StandardResponse{data=nil}								"Internal server error"
//	@Security		BearerAuth
//	@Router			/access_list_segmentation/{id} [get]
func (a *accessListSegmentation) GetAccessListSegmentationByID(w http.ResponseWriter, r *http.Request) {
	id := chi.URLParam(r, "id")
	if id == "" {
		a.logger.Errorf("[GetAccessListSegmentationByID] missing id parameter")
		localization.SendErrorByCodeResponse(w, localization.ErrorAccessListSegmentationInvalidID.Code)
		return
	}
	accessListSegmentation, err := a.service.GetAccessListSegmentationByID(r.Context(), id)
	if err != nil {
		a.logger.Errorf("[GetAccessListSegmentationByID] service error: %v", err)
		localization.SendErrorByCodeResponse(w, err.Error())
		return
	}
	localization.SendSuccessResponse(w, localization.SuccessAccessListSegmentationRetrieved, accessListSegmentation)
}

// GetAllAccessListSegmentation godoc
//
//	@Summary		Get all access list segmentations
//	@Description	Retrieve all access list segmentations with pagination and optional search
//	@Tags			Access List Segmentation
//	@Accept			json
//	@Produce		json
//	@Param			page		query		int									false	"Page number"		default(1)
//	@Param			per_page	query		int									false	"Items per page"	default(10)
//	@Param			search		query		string								false	"Search term"
//	@Success		200			{object}	localization.StandardResponse{data=object}	"Access list segmentations retrieved successfully"
//	@Failure		400			{object}	localization.StandardResponse{data=nil}		"Bad request"
//	@Failure		500			{object}	localization.StandardResponse{data=nil}		"Internal server error"
//	@Security		BearerAuth
//	@Router			/access_list_segmentation [get]
func (a *accessListSegmentation) GetAllAccessListSegmentation(w http.ResponseWriter, r *http.Request) {
	filterParams := local_util.ExtractFilterParams(r)

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

	accessListSegmentations, err := a.service.GetAllAccessListSegmentation(r.Context(), *filterParams)
	if err != nil {
		a.logger.Errorf("[GetAllAccessListSegmentation] service error: %v", err)
		localization.SendErrorByCodeResponse(w, err.Error())
		return
	}
	localization.SendSuccessResponse(w, localization.SuccessAccessListSegmentationRetrieved, accessListSegmentations)
}

// UpdateAccessListSegmentation implements accesslistsegmentation.AccessListSegmentationHandler.
func (a *accessListSegmentation) UpdateAccessListSegmentation(w http.ResponseWriter, r *http.Request) {
	id := chi.URLParam(r, "id")
	if id == "" {
		a.logger.Errorf("[UpdateAccessListSegmentation] missing id parameter")
		localization.SendErrorByCodeResponse(w, localization.ErrorAccessListSegmentationInvalidID.Code)
		return
	}
	var req access_list_segmentation_dto.UpdateAccessListSegmentationRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		a.logger.Errorf("[UpdateAccessListSegmentation] failed to decode request body: %v", err)
		localization.SendErrorByCodeResponse(w, err.Error())
		return
	}
	req.ID = id
	if err := req.Validate(); err != nil {
		a.logger.Errorf("[UpdateAccessListSegmentation] validation error: %v", err)
		localization.SendErrorByCodeResponse(w, err.Error())
		return
	}
	if err := a.service.UpdateAccessListSegmentation(r.Context(), req); err != nil {
		a.logger.Errorf("[UpdateAccessListSegmentation] service error: %v", err)
		localization.SendErrorByCodeResponse(w, err.Error())
		return
	}
	localization.SendSuccessResponse(w, localization.SuccessAccessListSegmentationUpdated, nil)
}

// GetAllAccessListSegmentationBySegmentIDorSegmentCode godoc
//
//	@Summary		Get all access list segmentations by segment ID or code
//	@Description	Retrieve all access list segmentations and access lists filtered by segment ID or segment code
//	@Tags			Access List Segmentation
//	@Accept			json
//	@Produce		json
//	@Param			id	path		string									true	"Segment ID or Segment Code"
//	@Success		200	{object}	localization.StandardResponse{data=object}	"Access list segmentations retrieved successfully"
//	@Failure		400	{object}	localization.StandardResponse{data=nil}		"Bad request"
//	@Failure		404	{object}	localization.StandardResponse{data=nil}		"Not found"
//	@Failure		500	{object}	localization.StandardResponse{data=nil}		"Internal server error"
//	@Security		BearerAuth
//	@Router			/access_list_segmentation/all/{id} [get]
func (a *accessListSegmentation) GetAllAccessListSegmentationBySegmentIDorSegmentCode(w http.ResponseWriter, r *http.Request) {
	id := chi.URLParam(r, "id")
	if id == "" {
		a.logger.Errorf("[GetAllAccessListSegmentationBySegmentIDorSegmentCode] missing id parameter")
		localization.SendErrorByCodeResponse(w, localization.ErrorAccessListSegmentationInvalidID.Code)
		return
	}
	accesssList, accessListSegmentations, err := a.service.GetAllAccessListSegmentationBySegmentIDorSegmentCode(r.Context(), id)
	if err != nil {
		a.logger.Errorf("[GetAllAccessListSegmentationBySegmentIDorSegmentCode] service error: %v", err)
		localization.SendErrorByCodeResponse(w, err.Error())
		return
	}

	localization.SendSuccessResponse(w, localization.SuccessAccessListSegmentationRetrieved, map[string]interface{}{
		"access_lists":              accesssList,
		"access_list_segmentations": accessListSegmentations,
	})
}

func InitAccessListSegmentationAdapter(accessListSegmentationApplication service.AccessListSegmentationService, logger utils.Logger) accesslistsegmentation.AccessListSegmentationHandler {
	return &accessListSegmentation{
		service: accessListSegmentationApplication,
		logger:  logger,
	}
}
