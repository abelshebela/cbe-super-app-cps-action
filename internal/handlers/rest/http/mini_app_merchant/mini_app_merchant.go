package miniappmerchant

import (
	"encoding/json"
	"errors"
	"net/http"

	miniappmerchant "cbe-super-app-cps-action/internal/constants/dto/mini_app_merchant"
	miniappmerchat "cbe-super-app-cps-action/internal/constants/interfaces/mini_app_merchant"
	"cbe-super-app-cps-action/internal/constants/localization"
	"cbe-super-app-cps-action/internal/service"
	local_util "cbe-super-app-cps-action/pkgs/utils"

	"github.com/go-chi/chi/v5"
	shared_utils "gitlab.com/bersufekadgetachew/cbe-super-app-shared/shared/utils"
	"go.opentelemetry.io/otel/attribute"
)

type miniAppMerchantAdapter struct {
	miniappMerchantService service.MiniAppMerchantService
	logger                 shared_utils.Logger
}

func NewMiniAppMerchantAdapter(miniappMerchantService service.MiniAppMerchantService, logger shared_utils.Logger) miniappmerchat.MiniAppMerchant {
	return &miniAppMerchantAdapter{miniappMerchantService: miniappMerchantService, logger: logger}
}

// Create Mini App Merchant
//
//	@Summary		Create Mini App Merchant
//	@Description	Creates a new mini app merchant
//	@Tags			MiniAppMerchant
//	@Security		BearerAuth
//	@Accept			json
//	@Produce		json
//	@Param			body			body		miniappmerchant.MiniAppMerchantDTO	true	"Mini App Merchant DTO"
//	@Success		201				{object}	localization.StandardResponse{data=model.MiniAppMerchant}
//	@Failure		400,401,422,500	{object}	localization.StandardResponse{data=nil}
//	@Router			/mini-app-merchants [post]
func (h *miniAppMerchantAdapter) Create(w http.ResponseWriter, r *http.Request) {
	ctx, span := local_util.TraceLogger(r.Context(), "handler", "createMiniAppMerchant", "handler", "miniAppMerchant")
	defer span.End()
	var req miniappmerchant.MiniAppMerchantDTO

	// Decode request body
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		span.RecordError(err)
		h.logger.Errorf("Failed to decode request body: %v", err)
		localization.SendErrorResponse(w, localization.ErrorMiniAppMerchantMarshalFailed, nil, nil)
		return
	}

	// Validate DTO
	if err := req.Validate(true); err != nil {
		span.RecordError(err)
		h.logger.Errorf("Validation failed: %v", err)
		localization.SendBadRequestResponse(w, err.Error())
		return
	}

	formattedPhone := local_util.FormatPhoneNumber(req.PhoneNumber)
	req.PhoneNumber = formattedPhone

	createdMerchant, err := h.miniappMerchantService.Create(ctx, &req)
	if err != nil {
		span.RecordError(err)
		h.logger.Errorf("Failed to create merchant: %v", err)
		localization.SendErrorByCodeResponse(w, err.Error())
		return
	}
	h.logger.Debugf("Created merchant: %+v", createdMerchant)
	localization.SendSuccessResponse(w, localization.SuccessMiniAppMerchantCreateRequestCreated, nil)
}

// Update Mini App Merchant
//
//	@Summary		Update Mini App Merchant
//	@Description	Updates an existing mini app merchant
//	@Tags			MiniAppMerchant
//	@Security		BearerAuth
//	@Accept			json
//	@Produce		json
//	@Param			id					path		string								true	"Merchant ID"
//	@Param			body				body		miniappmerchant.MiniAppMerchantDTO	true	"Mini App Merchant DTO"
//	@Success		200					{object}	localization.StandardResponse
//	@Failure		400,401,404,422,500	{object}	localization.StandardResponse{data=nil}
//	@Router			/mini-app-merchants/{id} [patch]
func (h *miniAppMerchantAdapter) Update(w http.ResponseWriter, r *http.Request) {
	ctx, span := local_util.TraceLogger(r.Context(), "handler", "updateMiniAppMerchant", "handler", "miniAppMerchant")
	defer span.End()
	id := chi.URLParam(r, "id")
	if id == "" {
		localization.SendErrorByCodeResponse(w, localization.ErrorInvalidInputParameters.Code)
		return
	}

	var req miniappmerchant.MiniAppMerchantDTO
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		span.RecordError(err)
		localization.SendErrorResponse(w, localization.ErrorMiniAppMerchantMarshalFailed, nil, nil)
		return
	}

	if req.IsEmpty() {
		span.RecordError(errors.New("no data provided for update"))
		localization.SendErrorResponse(w, localization.ErrorNoDataProvidedForUpdate, nil, nil)
		return
	}

	if err := req.Validate(false); err != nil {
		span.RecordError(err)
		localization.SendBadRequestResponse(w, err.Error())
		return
	}

	span.SetAttributes(attribute.String("mini_app_merchant.id", id))
	_, _, err := h.miniappMerchantService.Update(ctx, id, &req)
	if err != nil {
		span.RecordError(err)
		localization.SendErrorByCodeResponse(w, err.Error())
		return
	}

	localization.SendSuccessResponse(w, localization.SuccessMiniAppMerchantUpdateRequestCreated, nil)
}

// Delete Mini App Merchant
//
//	@Summary		Delete Mini App Merchant
//	@Description	Deletes a mini app merchant by ID
//	@Tags			MiniAppMerchant
//	@Security		BearerAuth
//	@Produce		json
//	@Param			id				path		string	true	"Merchant ID"
//	@Success		200				{object}	localization.StandardResponse
//	@Failure		400,401,404,500	{object}	localization.StandardResponse{data=nil}
//	@Router			/mini-app-merchants/{id} [delete]
func (h *miniAppMerchantAdapter) Delete(w http.ResponseWriter, r *http.Request) {
	ctx, span := local_util.TraceLogger(r.Context(), "handler", "deleteMiniAppMerchant", "handler", "miniAppMerchant")
	defer span.End()
	id := chi.URLParam(r, "id")
	if id == "" {
		span.RecordError(errors.New("missing or invalid parameter 'id'"))
		h.logger.Errorf("missing or invalid parameter 'id'")
		localization.SendErrorResponse(w, localization.ErrorInvalidInputParameters, nil, nil)
		return
	}

	userContext := local_util.ExtractUserContext(r)
	if local_util.IsIncomplete(userContext) {
		span.RecordError(errors.New("incomplete user context"))
		localization.SendErrorResponse(w, localization.ErrorIncompleteUserInfo, nil, nil)
		return
	}

	span.SetAttributes(attribute.String("mini_app_merchant.id", id))
	err := h.miniappMerchantService.Delete(ctx, id)
	if err != nil {
		span.RecordError(err)
		h.logger.Errorf("failed to delete merchant %s: %v", id, err)
		localization.SendErrorByCodeResponse(w, err.Error())
		return
	}

	localization.SendSuccessResponse(w, localization.SuccessMiniAppMerchantDeleteRequestCreated, nil)
}

// Enable Mini App Merchant
//
//	@Summary		Enable Mini App Merchant
//	@Description	Enables a mini app merchant by ID
//	@Tags			MiniAppMerchant
//	@Security		BearerAuth
//	@Produce		json
//	@Param			id				path		string	true	"Merchant ID"
//	@Success		200				{object}	localization.StandardResponse
//	@Failure		400,401,404,500	{object}	localization.StandardResponse{data=nil}
//	@Router			/mini-app-merchants/enable/{id} [patch]
func (h *miniAppMerchantAdapter) Enable(w http.ResponseWriter, r *http.Request) {
	ctx, span := local_util.TraceLogger(r.Context(), "handler", "enableMiniAppMerchant", "handler", "miniAppMerchant")
	defer span.End()
	id := chi.URLParam(r, "id")
	if id == "" {
		span.RecordError(errors.New("missing or invalid parameter 'id'"))
		localization.SendErrorByCodeResponse(w, localization.ErrorInvalidInputParameters.Code)
		return
	}

	userContext := local_util.ExtractUserContext(r)
	if local_util.IsIncomplete(userContext) {
		span.RecordError(errors.New("incomplete user context"))
		localization.SendErrorResponse(w, localization.ErrorIncompleteUserInfo, nil, nil)
		return
	}
	span.SetAttributes(attribute.String("mini_app_merchant.id", id))
	if err := h.miniappMerchantService.EnableOrDisable(ctx, id, true); err != nil {
		span.RecordError(err)
		localization.SendErrorByCodeResponse(w, err.Error())
		return
	}

	localization.SendSuccessResponse(w, localization.SuccessMiniAppMerchantEnableRequestCreated, nil)
}

// Disable Mini App Merchant
//
//	@Summary		Disable Mini App Merchant
//	@Description	Disables a mini app merchant by ID
//	@Tags			MiniAppMerchant
//	@Security		BearerAuth
//	@Produce		json
//	@Param			id				path		string	true	"Merchant ID"
//	@Success		200				{object}	localization.StandardResponse
//	@Failure		400,401,404,500	{object}	localization.StandardResponse{data=nil}
//	@Router			/mini-app-merchants/disable/{id} [patch]
func (h *miniAppMerchantAdapter) Disable(w http.ResponseWriter, r *http.Request) {
	ctx, span := local_util.TraceLogger(r.Context(), "handler", "disableMiniAppMerchant", "handler", "miniAppMerchant")
	defer span.End()
	id := chi.URLParam(r, "id")
	if id == "" {
		span.RecordError(errors.New("missing or invalid parameter 'id'"))
		localization.SendErrorByCodeResponse(w, localization.ErrorInvalidInputParameters.Code)
		return
	}

	userContext := local_util.ExtractUserContext(r)
	if local_util.IsIncomplete(userContext) {
		span.RecordError(errors.New("incomplete user context"))
		localization.SendErrorResponse(w, localization.ErrorIncompleteUserInfo, nil, nil)
		return
	}

	span.SetAttributes(attribute.String("mini_app_merchant.id", id))
	if err := h.miniappMerchantService.EnableOrDisable(ctx, id, false); err != nil {
		span.RecordError(err)
		localization.SendErrorByCodeResponse(w, err.Error())
		return
	}

	localization.SendSuccessResponse(w, localization.SuccessMiniAppMerchantDisableRequestCreated, nil)
}

// {
//   "ok": false,
//   "status": 401,
//   "timestamp": "2025-10-01T16:32:52.1069555+03:00",
//   "message": "User is not authorized"
// }

// Get Mini App Merchant by ID
//
//	@Summary		Get Mini App Merchant by ID
//	@Description	Retrieves a mini app merchant by ID
//	@Tags			MiniAppMerchant
//	@Security		BearerAuth
//	@Produce		json
//	@Param			id				path		string	true	"Merchant ID"
//	@Success		200				{object}	localization.StandardResponse{data=model.MiniAppMerchant}
//	@Failure		400,401,404,500	{object}	localization.StandardResponse{data=nil}
//	@Router			/mini-app-merchants/{id} [get]
func (h *miniAppMerchantAdapter) FindByID(w http.ResponseWriter, r *http.Request) {
	ctx, span := local_util.TraceLogger(r.Context(), "handler", "findMiniAppMerchantById", "handler", "miniAppMerchant")
	defer span.End()
	id := chi.URLParam(r, "id")
	if id == "" {
		localization.SendErrorByCodeResponse(w, localization.ErrorInvalidInputParameters.Code)
		return
	}

	span.SetAttributes(attribute.String("mini_app_merchant.id", id))
	result, err := h.miniappMerchantService.FindByID(ctx, id)
	if err != nil {
		span.RecordError(err)
		localization.SendErrorByCodeResponse(w, err.Error())
		return
	}

	localization.SendSuccessResponse(w, localization.SuccessMiniAppDetailsFetched, result)
}

// List Mini App Merchants with Pagination
//
//	@Summary		List Mini App Merchants
//	@Description	Retrieves a paginated list of mini app merchants. Searchable fields: merchant_id, bank_account_number, merchant_name, merchant_code, phone_number, email.
//	@Tags			MiniAppMerchant
//	@Security		BearerAuth
//	@Produce		json
//	@Param			page				query	int		false	"Page number"
//	@Param			per_page			query	int		false	"Items per page"
//	@Param			merchant_type		query	string	false	"Filter by merchant type"
//	@Param			merchant_code		query	string	false	"Filter by merchant code"
//	@Param			merchant_name		query	string	false	"Filter by merchant name"
//	@Param			email				query	string	false	"Filter by email"
//	@Param			phone_number		query	string	false	"Filter by phone number"
//	@Param			enabled				query	bool	false	"Filter by enabled status"
//	@Param			bank_account_number	query	string	false	"Filter by bank account number"
//	@Param			search				query	string	false	"Search term (searches merchant_id, bank_account_number, merchant_name, merchant_code, phone_number, email)"
//	@Success		200	{object}	localization.StandardResponse{data=[]model.MiniAppMerchant}
//	@Failure		400,401,500	{object}	localization.StandardResponse{data=nil}
//	@Router			/mini-app-merchants [get]
func (h *miniAppMerchantAdapter) FindAllWithPagination(w http.ResponseWriter, r *http.Request) {
	ctx, span := local_util.TraceLogger(r.Context(), "handler", "findAllMiniAppMerchants", "handler", "miniAppMerchant")
	defer span.End()
	filterParams := local_util.ExtractFilterParams(r)
	if filterParams.Page < 0 || filterParams.PerPage < 0 {
		localization.SendErrorResponse(w, localization.ErrorInvalidRequest, nil, nil)
		return
	}

	phoneNumber, ok := filterParams.Filters["phone_number"].(string)
	if !ok {
		h.logger.Errorf("Phone number is not present")
	}
	formattedPhone := local_util.FormatPhoneNumber(phoneNumber)

	filterParams.Filters["phone_number"] = formattedPhone

	miniAppMerchant, err := h.miniappMerchantService.FindAllWithPagination(ctx, filterParams)
	if err != nil {
		span.RecordError(err)
		localization.SendErrorByCodeResponse(w, err.Error())
		return
	}

	localization.SendSuccessResponse(w, localization.SuccessMiniAppDetailsFetched, miniAppMerchant)
}

// // Merchant Lookup
// //@Summary Merchant Lookup
// //@Description Retrieves a merchant by ID
// //@Tags MiniAppMerchant
// //@Security BearerAuth
// // @Produce json
// // @Param merchant_id path string true "Merchant ID"
// // @Success 200 {object} localization.StandardResponse{data=merchant_lookup.MerchantLookUpResponse}
// // @Failure 400,401,404,500 {object} localization.StandardResponse{data=nil}
// // @Router /mini-app-merchants/merchant-lookup/{merchant_id} [get]
func (h *miniAppMerchantAdapter) MerchantLookup(w http.ResponseWriter, r *http.Request) {
	ctx, span := local_util.TraceLogger(r.Context(), "handler", "miniAppMerchantLookup", "handler", "miniAppMerchant")
	defer span.End()
	merchantID := chi.URLParam(r, "merchant_id")
	if merchantID == "" {
		span.RecordError(errors.New("missing or invalid parameter 'merchant_id'"))
		localization.SendErrorByCodeResponse(w, localization.ErrorInvalidInputParameters.Code)
		return
	}

	span.SetAttributes(attribute.String("mini_app_merchant.id", merchantID))
	result, err := h.miniappMerchantService.MerchantLookup(ctx, merchantID)
	if err != nil {
		span.RecordError(err)
		localization.SendErrorByCodeResponse(w, err.Error())
		return
	}

	localization.SendSuccessResponse(w, localization.SuccessMiniAppDetailsFetched, result)
}
