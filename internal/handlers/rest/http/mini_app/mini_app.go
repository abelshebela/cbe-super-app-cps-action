package miniapphandler

import (
	miniappdto "cbe-super-app-cps-action/internal/constants/dto/mini_app"
	miniAppInbound "cbe-super-app-cps-action/internal/constants/interfaces/mini_app"
	"cbe-super-app-cps-action/internal/constants/localization"
	miniappcore "cbe-super-app-cps-action/internal/handlers/rest/http/mini_app/core"
	"cbe-super-app-cps-action/internal/service"
	local_util "cbe-super-app-cps-action/pkgs/utils"
	"net/http"

	"gitlab.com/bersufekadgetachew/cbe-super-app-shared/shared/utils"

	"github.com/go-chi/chi/v5"
)

type HttpStore struct {
	miniAppSrvc service.MiniAppService
	logger      utils.Logger
}

func InitMiniAppAdapter(app service.MiniAppService, logger utils.Logger) miniAppInbound.MiniAppInbound {
	return &HttpStore{
		miniAppSrvc: app,
		logger:      logger,
	}
}

// CreateMiniApp godoc
//
//	@Summary		Create a new mini app
//	@Description	Creates a new mini app with form-data (supports file upload)
//	@Tags			MiniApps
//	@Accept			multipart/form-data
//	@Produce		json
//	@Param			app_name				formData	string						true	"App name"
//	@Param			app_icon				formData	file						false	"App icon file"
//	@Param			banner_image			formData	file						false	"Banner image file"
//	@Param			commission_gl_account	formData	string						false	"Commission GL account"
//	@Param			merchant_id				formData	string						true	"Merchant ID"
//	@Param			is_event_mini_app		formData	bool						false	"Is event mini app"
//	@Param			is_three_click			formData	bool						false	"Is three click app"
//	@Param			app_view_type			formData	string						false	"App view type"
//	@Param			url						formData	string						false	"Mini app URL"
//	@Param			ifb_product_code		formData	string						false	"IFB product code"
//	@Param			ifb_vat_code			formData	string						false	"IFB VAT code"
//	@Param			ifb_service_fee_code	formData	string						false	"IFB service fee code"
//	@Param			cb_product_code			formData	string						false	"CB product code"
//	@Param			cb_vat_code				formData	string						false	"CB VAT code"
//	@Param			cb_service_fee_code		formData	string						false	"CB service fee code"
//	@Success		201						{object}	localization.ResponseCode	"Mini app created successfully"
//	@Failure		400						{object}	localization.ResponseCode	"Bad request"
//	@Failure		401						{object}	localization.ResponseCode	"Unauthorized"
//	@Failure		500						{object}	localization.ResponseCode	"Internal server error"
//	@Router			/mini-apps [post]
//	@Security		BearerAuth
func (h *HttpStore) CreateMiniApp(w http.ResponseWriter, r *http.Request) {
	req, err := miniappcore.ParseMiniAppRequestFromMultipartForm(r, true)
	if err != nil {
		h.logger.Errorf("failed to parse mini app request from multipart form: %v", err)
		localization.SendErrorByCodeResponse(w, err.Error())
		return
	}

	if err := req.Validate(true); err != nil {
		h.logger.Errorf("mini app request validation failed: %v", err)
		localization.SendBadRequestResponse(w, err.Error())
		return
	}

	dto, err := miniappcore.ToMiniAppCreateRequest(req, true)
	if err != nil {
		h.logger.Errorf("failed to convert request to DTO: %v", err)
		localization.SendBadRequestResponse(w, err.Error())
		return
	}

	if err := h.miniAppSrvc.CreateMiniApp(r.Context(), dto); err != nil {
		h.logger.Errorf("failed to create mini app: %v", err)
		localization.SendErrorByCodeResponse(w, err.Error())
		return
	}

	h.logger.Infof("mini app creation request submitted successfully")
	localization.SendSuccessResponse(w, localization.SuccessMiniAppCreateRequestCreated, nil)
}

// UpdateMiniApp godoc
//
//	@Summary		Update a mini app
//	@Description	Updates an existing mini app with form-data (supports file upload)
//	@Tags			MiniApps
//	@Accept			multipart/form-data
//	@Produce		json
//	@Param			id						path		string						true	"Mini app ID"
//	@Param			app_name				formData	string						false	"App name"
//	@Param			app_icon				formData	file						false	"App icon file"
//	@Param			banner_image			formData	file						false	"Banner image file"
//	@Param			commission_gl_account	formData	string						false	"Commission GL account"
//	@Param			merchant_id				formData	string						false	"Merchant ID"
//	@Param			is_event_mini_app		formData	bool						false	"Is event mini app"
//	@Param			is_three_click			formData	bool						false	"Is three click app"
//	@Param			app_view_type			formData	string						false	"App view type"
//	@Param			url						formData	string						false	"Mini app URL"
//	@Param			ifb_product_code		formData	string						false	"IFB product code"
//	@Param			ifb_vat_code			formData	string						false	"IFB VAT code"
//	@Param			ifb_service_fee_code	formData	string						false	"IFB service fee code"
//	@Param			cb_product_code			formData	string						false	"CB product code"
//	@Param			cb_vat_code				formData	string						false	"CB VAT code"
//	@Param			cb_service_fee_code		formData	string						false	"CB service fee code"
//	@Success		200						{object}	localization.ResponseCode	"Mini app updated successfully"
//	@Failure		400						{object}	localization.ResponseCode	"Bad request"
//	@Failure		401						{object}	localization.ResponseCode	"Unauthorized"
//	@Failure		404						{object}	localization.ResponseCode	"Not found"
//	@Failure		500						{object}	localization.ResponseCode	"Internal server error"
//	@Router			/mini-apps/{id} [patch]
//	@Security		BearerAuth
func (h *HttpStore) UpdateMiniApp(w http.ResponseWriter, r *http.Request) {
	id := chi.URLParam(r, "id")
	if id == "" {
		h.logger.Errorf("mini app ID is required for update")
		localization.SendErrorResponse(w, localization.ErrorMiniAppIDRequired, nil, nil)
		return
	}

	req, err := miniappcore.ParseMiniAppRequestFromMultipartForm(r, false)
	if err != nil {
		h.logger.Errorf("failed to parse mini app update request: %v", err)
		localization.SendBadRequestResponse(w, err.Error())
		return
	}
	if err := req.Validate(false); err != nil {
		h.logger.Errorf("mini app request validation failed: %v", err)
		localization.SendBadRequestResponse(w, err.Error())
		return
	}
	dtoForUpdate, err := miniappcore.ToMiniAppCreateRequest(req, false)
	if err != nil {
		h.logger.Errorf("failed to convert update request to DTO: %v", err)
		localization.SendBadRequestResponse(w, err.Error())
		return
	}
	if miniappdto.IsEmptyUpdate(dtoForUpdate) {
		h.logger.Warnf("no data provided for mini app update, mini app ID: %s", id)
		localization.SendErrorResponse(w, localization.ErrorUpdateMiniAppEmptyPayload, nil, nil)
		return
	}

	dto, err := miniappcore.ToMiniAppCreateRequest(req, false)
	if err != nil {
		h.logger.Errorf("failed to convert update request to DTO: %v", err)
		localization.SendBadRequestResponse(w, err.Error())
		return
	}
	dto.ID = id

	if err := h.miniAppSrvc.UpdateMiniApp(r.Context(), dto); err != nil {
		h.logger.Errorf("failed to update mini app (ID: %s): %v", id, err)
		localization.SendErrorByCodeResponse(w, err.Error())
		return
	}

	h.logger.Infof("mini app update request submitted successfully, mini app ID: %s", id)
	localization.SendSuccessResponse(w, localization.SuccessMiniAppUpdateRequestCreated, nil)
}

// DeleteMiniApp godoc
//
//	@Summary		Delete a mini app
//	@Description	Deletes a mini app by ID
//	@Tags			MiniApps
//	@Accept			json
//	@Produce		json
//	@Param			id	path		string						true	"Mini app ID"
//	@Success		200	{object}	localization.ResponseCode	"Mini app deleted successfully"
//	@Failure		400	{object}	localization.ResponseCode	"Bad request"
//	@Failure		401	{object}	localization.ResponseCode	"Unauthorized"
//	@Failure		404	{object}	localization.ResponseCode	"Not found"
//	@Failure		500	{object}	localization.ResponseCode	"Internal server error"
//	@Router			/mini-apps/{id} [delete]
//	@Security		BearerAuth
func (h *HttpStore) DeleteMiniApp(w http.ResponseWriter, r *http.Request) {
	id := chi.URLParam(r, "id")
	if id == "" {
		h.logger.Errorf("mini app ID is required for delete")
		localization.SendErrorResponse(w, localization.ErrorMiniAppIDRequired, nil, nil)
		return
	}

	if err := h.miniAppSrvc.DeleteMiniApp(r.Context(), id); err != nil {
		h.logger.Errorf("failed to delete mini app (ID: %s): %v", id, err)
		localization.SendErrorByCodeResponse(w, err.Error())
		return
	}

	h.logger.Infof("mini app delete request submitted successfully, mini app ID: %s", id)
	localization.SendSuccessResponse(w, localization.SuccessMiniAppDeletedRequestCreated, nil)
}

// EnableMiniAppByID godoc
//
//	@Summary		Enable a mini app
//	@Description	Enables a mini app by ID
//	@Tags			MiniApps
//	@Accept			json
//	@Produce		json
//	@Param			id	path		string						true	"Mini app ID"
//	@Success		200	{object}	localization.ResponseCode	"Mini app enabled successfully"
//	@Failure		400	{object}	localization.ResponseCode	"Bad request"
//	@Failure		401	{object}	localization.ResponseCode	"Unauthorized"
//	@Failure		404	{object}	localization.ResponseCode	"Not found"
//	@Failure		500	{object}	localization.ResponseCode	"Internal server error"
//	@Router			/mini-apps/enable/{id} [patch]
//	@Security		BearerAuth
func (h *HttpStore) EnableMiniAppByID(w http.ResponseWriter, r *http.Request) {
	id := chi.URLParam(r, "id")
	if id == "" {
		h.logger.Errorf("mini app ID is required for enable")
		localization.SendErrorResponse(w, localization.ErrorMiniAppIDRequired, nil, nil)
		return
	}

	if err := h.miniAppSrvc.EnableDisableMiniAppByID(r.Context(), id, true); err != nil {
		h.logger.Errorf("failed to enable mini app (ID: %s): %v", id, err)
		localization.SendErrorByCodeResponse(w, err.Error())
		return
	}

	h.logger.Infof("mini app enable request submitted successfully, mini app ID: %s", id)
	localization.SendSuccessResponse(w, localization.SuccessMiniAppEnableRequestSubmitted, nil)
}

// DisableMiniAppByID godoc
//
//	@Summary		Disable a mini app
//	@Description	Disables a mini app by ID
//	@Tags			MiniApps
//	@Accept			json
//	@Produce		json
//	@Param			id	path		string						true	"Mini app ID"
//	@Success		200	{object}	localization.ResponseCode	"Mini app disabled successfully"
//	@Failure		400	{object}	localization.ResponseCode	"Bad request"
//	@Failure		401	{object}	localization.ResponseCode	"Unauthorized"
//	@Failure		404	{object}	localization.ResponseCode	"Not found"
//	@Failure		500	{object}	localization.ResponseCode	"Internal server error"
//	@Router			/mini-apps/disable/{id} [patch]
//	@Security		BearerAuth
func (h *HttpStore) DisableMiniAppByID(w http.ResponseWriter, r *http.Request) {
	id := chi.URLParam(r, "id")
	if id == "" {
		h.logger.Errorf("mini app ID is required for disable")
		localization.SendErrorResponse(w, localization.ErrorMiniAppIDRequired, nil, nil)
		return
	}

	if err := h.miniAppSrvc.EnableDisableMiniAppByID(r.Context(), id, false); err != nil {
		h.logger.Errorf("failed to disable mini app (ID: %s): %v", id, err)
		localization.SendErrorByCodeResponse(w, err.Error())
		return
	}

	h.logger.Infof("mini app disable request submitted successfully, mini app ID: %s", id)
	localization.SendSuccessResponse(w, localization.SuccessMiniAppDisableRequestSubmitted, nil)
}

// DetailMiniAppByID godoc
//
//	@Summary		Get mini app by ID
//	@Description	Retrieves the details of a mini app by its ID
//	@Tags			MiniApps
//	@Accept			json
//	@Produce		json
//	@Param			id	path		string						true	"Mini app ID"
//	@Success		200	{object}	miniappdto.MiniAppResponse	"Mini app details"
//	@Failure		400	{object}	localization.ResponseCode	"Bad request"
//	@Failure		401	{object}	localization.ResponseCode	"Unauthorized"
//	@Failure		404	{object}	localization.ResponseCode	"Not found"
//	@Failure		500	{object}	localization.ResponseCode	"Internal server error"
//	@Router			/mini-apps/{id} [get]
//	@Security		BearerAuth
func (h *HttpStore) DetailMiniAppByID(w http.ResponseWriter, r *http.Request) {
	id := chi.URLParam(r, "id")
	if id == "" {
		h.logger.Errorf("mini app ID is required for fetch")
		localization.SendErrorResponse(w, localization.ErrorMiniAppIDRequired, nil, nil)
		return
	}

	miniApp, err := h.miniAppSrvc.FindByID(r.Context(), id)
	if err != nil {
		h.logger.Errorf("failed to fetch mini app (ID: %s): %v", id, err)
		localization.SendErrorByCodeResponse(w, err.Error())
		return
	}

	h.logger.Infof("mini app fetched successfully, mini app ID: %s", id)
	localization.SendSuccessResponse(w, localization.SuccessMiniAppFetchedByID, miniApp)
}

// ListMiniApp godoc
//
//	@Summary		List mini apps
//	@Description	Retrieves a list of mini apps with optional filters
//	@Tags			MiniApps
//	@Accept			json
//	@Produce		json
//	@Param			filter	query		string						false	"Filter params (e.g., status, type)"
//	@Success		200		{array}		miniappdto.MiniAppResponse	"List of mini apps"
//	@Failure		400		{object}	localization.ResponseCode	"Bad request"
//	@Failure		401		{object}	localization.ResponseCode	"Unauthorized"
//	@Failure		500		{object}	localization.ResponseCode	"Internal server error"
//	@Router			/mini-apps [get]
//	@Security		BearerAuth
func (h *HttpStore) ListMiniApp(w http.ResponseWriter, r *http.Request) {
	filter := local_util.ExtractFilterParams(r)
	h.logger.Infof("fetching mini apps with filter: %+v", filter)

	list, err := h.miniAppSrvc.ListMiniApp(r.Context(), *filter)
	if err != nil {
		h.logger.Errorf("failed to fetch mini apps: %v", err)
		localization.SendErrorByCodeResponse(w, err.Error())
		return
	}

	h.logger.Infof("mini apps fetched successfully")
	localization.SendSuccessResponse(w, localization.SuccessMiniAppsRetrieved, list)
}
