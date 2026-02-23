package ecommercemerchant

import (
	"context"
	"encoding/json"
	"errors"
	"net/http"

	"cbe-super-app-cps-action/internal/constants"
	miniappmerchant "cbe-super-app-cps-action/internal/constants/dto/ecommerce-merchant"
	ecommerce_merchant "cbe-super-app-cps-action/internal/constants/interfaces/ecommerce_merchant"
	"cbe-super-app-cps-action/internal/constants/localization"
	"cbe-super-app-cps-action/internal/constants/types"
	"cbe-super-app-cps-action/internal/service"
	local_util "cbe-super-app-cps-action/pkgs/utils"

	"github.com/go-chi/chi/v5"
	shared_utils "gitlab.com/bersufekadgetachew/cbe-super-app-shared/shared/utils"
	"go.opentelemetry.io/otel/attribute"
)

// EcommerceMerchantRequest mirrors miniappmerchant.EcommerceMerchant for Swagger @Param only
type EcommerceMerchantRequest struct {
	MerchantName     string `json:"merchant_name"`
	MerchantCode     string `json:"merchant_code"`
	AccountNumber    string `json:"account_number"`
	SettlementMethod string `json:"settlement_method"`
}

type ecommerceMerchantAdapter struct {
	srv    service.EcommerceMerchantService
	logger shared_utils.Logger
}

func NewEcommerceMerchantdapter(srv service.EcommerceMerchantService, logger shared_utils.Logger) ecommerce_merchant.EcommerceMerchant {
	return &ecommerceMerchantAdapter{srv: srv, logger: logger}
}

// Create Ecommerce Merchant
//
//	@Summary		Create  Ecommerce Merchant
//	@Description	Creates a new  Ecommerce Merchant
//	@Tags			ecommerce-merchant
//	@Security		BearerAuth
//	@Accept			json
//	@Produce		json
//	@Param			body			body		ecommercemerchant.EcommerceMerchantRequest	true	"Mini App Merchant (see miniappmerchant.EcommerceMerchant)"
//	@Success		201				{object}	localization.StandardResponse{data=nil}
//	@Failure		400,401,422,500	{object}	localization.StandardResponse{data=nil}
//	@Router			/ecommerce-merchant [post]
func (h *ecommerceMerchantAdapter) Create(w http.ResponseWriter, r *http.Request) {
	ctx, span := local_util.TraceLogger(r.Context(), "", "createMiniAppMerchant", "handler", "miniAppMerchant")
	defer span.End()
	log := local_util.LoggerFromCtx(ctx, h.logger)

	md := &types.ContextMetadata{}
	ctx = context.WithValue(ctx, constants.ContextKeyMetadata, md)

	var reqDTO miniappmerchant.EcommerceMerchant

	if err := json.NewDecoder(r.Body).Decode(&reqDTO); err != nil {
		span.RecordError(err)
		log.Errorf("[EcomMerchH][Create] decode body err: %v", err)
		localization.SendErrorResponse(w, localization.ErrorMiniAppMerchantMarshalFailed, nil, nil)
		return
	}

	if err := reqDTO.Validate(true); err != nil {
		span.RecordError(err)
		log.Errorf("[EcomMerchH][Create] validate err: %v", err)
		localization.SendBadRequestResponse(w, err.Error())
		return
	}

	// formattedPhone := local_util.FormatPhoneNumber(reqDTO.PhoneNumber)

	// reqDTO.PhoneNumber = formattedPhone
	userContext := local_util.ExtractUserContext(r)
	if local_util.IsIncomplete(userContext) {
		span.RecordError(errors.New("incomplete user context"))
		log.Warnf("[EcomMerchH][Create] incomplete user info: %+v", userContext)
		localization.SendErrorResponse(w, localization.ErrorIncompleteUserInfo, nil, nil)
		return
	}

	_, err := h.srv.Create(ctx, &reqDTO)
	if err != nil {
		span.RecordError(err)
		log.Errorf("[EcomMerchH][Create] svc err: %v", err)
		localization.SendErrorByCodeResponse(w, err.Error())
		return
	}

	if md.IsMakerOnly {
		userCode, _ := ctx.Value(constants.ContextKey("user_code")).(string)
		log.Infof("[Create] request sent successfully for user_code: %s is_maker_only: %v", userCode, md.IsMakerOnly)
		localization.SendSuccessResponse(w, localization.SuccessEcommerceMerchantCreated, nil)
		return
	}

	localization.SendSuccessResponse(w, localization.SuccessEcommerceMerchantCreatedSuccessfully, nil)
}

// Update  Ecommerce Merchant
//
//	@Summary		Update  Ecommerce Merchant
//	@Description	Updates an existing  Ecommerce Merchant
//	@Tags			ecommerce-merchant
//	@Security		BearerAuth
//	@Accept			json
//	@Produce		json
//	@Param			id					path		string								true	"Merchant ID"
//	@Param			body				body		ecommercemerchant.EcommerceMerchantRequest	true	"Mini App Merchant (see miniappmerchant.EcommerceMerchant)"
//	@Success		200					{object}	localization.StandardResponse
//	@Failure		400,401,404,422,500	{object}	localization.StandardResponse{data=nil}
//	@Router			/ecommerce-merchant/{id} [patch]
func (h *ecommerceMerchantAdapter) Update(w http.ResponseWriter, r *http.Request) {
	ctx, span := local_util.TraceLogger(r.Context(), "", "updateMiniAppMerchant", "handler", "miniAppMerchant")
	defer span.End()
	log := local_util.LoggerFromCtx(ctx, h.logger)

	md := &types.ContextMetadata{}
	ctx = context.WithValue(ctx, constants.ContextKeyMetadata, md)

	id := chi.URLParam(r, "id")
	if id == "" {
		localization.SendErrorByCodeResponse(w, localization.ErrorInvalidInputParameters.Code)
		return
	}

	var reqDTO miniappmerchant.EcommerceMerchant
	if err := json.NewDecoder(r.Body).Decode(&reqDTO); err != nil {
		span.RecordError(err)
		localization.SendErrorResponse(w, localization.ErrorMiniAppMerchantMarshalFailed, nil, nil)
		return
	}

	if reqDTO.IsEmpty() {
		span.RecordError(errors.New("no data provided for update"))
		localization.SendErrorResponse(w, localization.ErrorNoDataProvidedForUpdate, nil, nil)
		return
	}

	if err := reqDTO.Validate(false); err != nil {
		span.RecordError(err)
		localization.SendBadRequestResponse(w, err.Error())
		return
	}

	userContext := local_util.ExtractUserContext(r)
	if local_util.IsIncomplete(userContext) {
		span.RecordError(errors.New("incomplete user context"))
		localization.SendErrorResponse(w, localization.ErrorIncompleteUserInfo, nil, nil)
		return
	}

	span.SetAttributes(attribute.String("ecommerce_merchant_dto.id", id))
	_, _, err := h.srv.Update(ctx, id, &reqDTO)
	if err != nil {
		span.RecordError(err)
		localization.SendErrorByCodeResponse(w, err.Error())
		return
	}

	if md.IsMakerOnly {
		userCode, _ := ctx.Value(constants.ContextKey("user_code")).(string)
		log.Infof("[Update] request sent successfully for user_code: %s is_maker_only: %v", userCode, md.IsMakerOnly)
		localization.SendSuccessResponse(w, localization.SuccessEcommerceMerchantUpdated, nil)
		return
	}

	localization.SendSuccessResponse(w, localization.SuccessEcommerceMerchantUpdatedSuccessfully, nil)
}

// Delete  Ecommerce Merchant
//
//	@Summary		Delete  Ecommerce Merchant
//	@Description	Deletes a mini app merchant by ID
//	@Tags			ecommerce-merchant
//	@Security		BearerAuth
//	@Produce		json
//	@Param			id				path		string	true	"Merchant ID"
//	@Success		200				{object}	localization.StandardResponse
//	@Failure		400,401,404,500	{object}	localization.StandardResponse{data=nil}
//	@Router			/ecommerce-merchant/{id} [delete]
func (h *ecommerceMerchantAdapter) Delete(w http.ResponseWriter, r *http.Request) {
	ctx, span := local_util.TraceLogger(r.Context(), "", "deleteMiniAppMerchant", "handler", "miniAppMerchant")
	defer span.End()
	log := local_util.LoggerFromCtx(ctx, h.logger)

	md := &types.ContextMetadata{}
	ctx = context.WithValue(ctx, constants.ContextKeyMetadata, md)

	id := chi.URLParam(r, "id")
	if id == "" {
		span.RecordError(errors.New("missing or invalid parameter 'id'"))
		log.Errorf("[EcomMerchH][Delete] missing id param")
		localization.SendErrorResponse(w, localization.ErrorInvalidInputParameters, nil, nil)
		return
	}

	userContext := local_util.ExtractUserContext(r)
	if local_util.IsIncomplete(userContext) {
		span.RecordError(errors.New("incomplete user context"))
		localization.SendErrorResponse(w, localization.ErrorIncompleteUserInfo, nil, nil)
		return
	}

	span.SetAttributes(attribute.String("ecommerce_merchant_dto.id", id))
	err := h.srv.Delete(ctx, id)
	if err != nil {
		span.RecordError(err)
		log.Errorf("[EcomMerchH][Delete] svc err id: %s: %v", id, err)
		localization.SendErrorByCodeResponse(w, err.Error())
		return
	}
	if md.IsMakerOnly {
		userCode, _ := ctx.Value(constants.ContextKey("user_code")).(string)
		log.Infof("[Delete] request sent successfully for user_code: %s is_maker_only: %v", userCode, md.IsMakerOnly)
		localization.SendSuccessResponse(w, localization.SuccessEcommerceMerchantDeleted, nil)
		return
	}

	localization.SendSuccessResponse(w, localization.SuccessEcommerceMerchantDeletedSuccessfully, nil)
}

// Enable  Ecommerce Merchant
//
//	@Summary		Enable  Ecommerce Merchant
//	@Description	Enables a mini app merchant by ID
//	@Tags			ecommerce-merchant
//	@Security		BearerAuth
//	@Produce		json
//	@Param			id				path		string	true	"Merchant ID"
//	@Success		200				{object}	localization.StandardResponse
//	@Failure		400,401,404,500	{object}	localization.StandardResponse{data=nil}
//	@Router			/ecommerce-merchant/enable/{id} [patch]
func (h *ecommerceMerchantAdapter) Enable(w http.ResponseWriter, r *http.Request) {
	ctx, span := local_util.TraceLogger(r.Context(), "", "enableMiniAppMerchant", "handler", "miniAppMerchant")
	defer span.End()
	log := local_util.LoggerFromCtx(ctx, h.logger)

	md := &types.ContextMetadata{}
	ctx = context.WithValue(ctx, constants.ContextKeyMetadata, md)

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
	span.SetAttributes(attribute.String("ecommerce_merchant_dto.id", id))
	if err := h.srv.EnableOrDisable(ctx, id, true); err != nil {
		span.RecordError(err)
		localization.SendErrorByCodeResponse(w, err.Error())
		return
	}
	if md.IsMakerOnly {
		userCode, _ := ctx.Value(constants.ContextKey("user_code")).(string)
		log.Infof("[Enable] request sent successfully for user_code: %s is_maker_only: %v", userCode, md.IsMakerOnly)
		localization.SendSuccessResponse(w, localization.SuccessEcommerceMerchantEnable, nil)
		return
	}

	localization.SendSuccessResponse(w, localization.SuccessEcommerceMerchantEnableSuccessfully, nil)
}

// Disable  Ecommerce Merchant
//
//	@Summary		Disable  Ecommerce Merchant
//	@Description	Disables a mini app merchant by ID
//	@Tags			ecommerce-merchant
//	@Security		BearerAuth
//	@Produce		json
//	@Param			id				path		string	true	"Merchant ID"
//	@Success		200				{object}	localization.StandardResponse
//	@Failure		400,401,404,500	{object}	localization.StandardResponse{data=nil}
//	@Router			/ecommerce-merchant/disable/{id} [patch]
func (h *ecommerceMerchantAdapter) Disable(w http.ResponseWriter, r *http.Request) {
	ctx, span := local_util.TraceLogger(r.Context(), "", "disableMiniAppMerchant", "handler", "miniAppMerchant")
	defer span.End()
	log := local_util.LoggerFromCtx(ctx, h.logger)

	md := &types.ContextMetadata{}
	ctx = context.WithValue(ctx, constants.ContextKeyMetadata, md)

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

	span.SetAttributes(attribute.String("ecommerce_merchant_dto.id", id))
	if err := h.srv.EnableOrDisable(ctx, id, false); err != nil {
		span.RecordError(err)
		localization.SendErrorByCodeResponse(w, err.Error())
		return
	}
	if md.IsMakerOnly {
		userCode, _ := ctx.Value(constants.ContextKey("user_code")).(string)
		log.Infof("[Disable] request sent successfully for user_code: %s is_maker_only: %v", userCode, md.IsMakerOnly)
		localization.SendSuccessResponse(w, localization.SuccessEcommerceMerchantDisable, nil)
		return
	}

	localization.SendSuccessResponse(w, localization.SuccessEcommerceMerchantDisableSuccessfully, nil)
}

// {
//   "ok": false,
//   "status": 401,
//   "timestamp": "2025-10-01T16:32:52.1069555+03:00",
//   "message": "User is not authorized"
// }

// Get  Ecommerce Merchant by ID
//
//	@Summary		Get  Ecommerce Merchant by ID
//	@Description	Retrieves a mini app merchant by ID
//	@Tags			ecommerce-merchant
//	@Security		BearerAuth
//	@Produce		json
//	@Param			id				path		string	true	"Merchant ID"
//	@Success		200				{object}	localization.StandardResponse{data=object}
//	@Failure		400,401,404,500	{object}	localization.StandardResponse{data=nil}
//	@Router			/ecommerce-merchant/{id} [get]
func (h *ecommerceMerchantAdapter) FindByID(w http.ResponseWriter, r *http.Request) {
	ctx, span := local_util.TraceLogger(r.Context(), "", "findMiniAppMerchantById", "handler", "miniAppMerchant")
	defer span.End()
	log := local_util.LoggerFromCtx(ctx, h.logger)
	_ = log
	id := chi.URLParam(r, "id")
	if id == "" {
		localization.SendErrorByCodeResponse(w, localization.ErrorInvalidInputParameters.Code)
		return
	}

	span.SetAttributes(attribute.String("ecommerce_merchant_dto.id", id))
	result, err := h.srv.FindByID(ctx, id)
	if err != nil {
		span.RecordError(err)
		localization.SendErrorByCodeResponse(w, err.Error())
		return
	}

	localization.SendSuccessResponse(w, localization.SuccessEcommerceMerchantFetchedSuccessfully, result)
}

// List  Ecommerce Merchants with Pagination
//
//	@Summary		List  Ecommerce Merchants
//	@Description	Retrieves a paginated list of  Ecommerce Merchants
//	@Tags			ecommerce-merchant
//	@Security		BearerAuth
//	@Produce		json
//	@Param			page		query		int	false	"Page number"
//	@Param			per_page	query		int	false	"Items per page"
//	@Success		200			{object}	localization.StandardResponse{data=object}
//	@Failure		400,401,500	{object}	localization.StandardResponse{data=nil}
//	@Router			/ecommerce-merchant [get]
func (h *ecommerceMerchantAdapter) FindAllWithPagination(w http.ResponseWriter, r *http.Request) {
	ctx, span := local_util.TraceLogger(r.Context(), "", "findAllMiniAppMerchants", "handler", "miniAppMerchant")
	defer span.End()
	log := local_util.LoggerFromCtx(ctx, h.logger)
	_ = log
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

	if filterParams.Page < 0 || filterParams.PerPage < 0 {
		localization.SendErrorResponse(w, localization.ErrorInvalidRequest, nil, nil)
		return
	}

	phoneNumber, ok := filterParams.Filters["phone_number"].(string)
	if ok {
		formattedPhone := local_util.FormatPhoneNumber(phoneNumber)
		filterParams.Filters["phone_number"] = formattedPhone
	}

	miniAppMerchant, err := h.srv.FindAllWithPagination(ctx, filterParams)
	if err != nil {
		span.RecordError(err)
		localization.SendErrorByCodeResponse(w, err.Error())
		return
	}

	localization.SendSuccessResponse(w, localization.SuccessEcommerceMerchantsFetchedSuccessfully, miniAppMerchant)
}

// Merchant Lookup
// @Summary Merchant Lookup
// @Description Retrieves a merchant by ID
// @Tags ecommerce-merchant
// @Security BearerAuth
// @Produce json
// @Param merchant_id path string true "Merchant ID"
// @Success 200 {object} localization.StandardResponse{data=object}
// @Failure 400,401,404,500 {object} localization.StandardResponse{data=nil}
// @Router /ecommerce-merchant/merchant-lookup/{merchant_id} [get]
func (h *ecommerceMerchantAdapter) MerchantLookup(w http.ResponseWriter, r *http.Request) {
	ctx, span := local_util.TraceLogger(r.Context(), "", "miniAppMerchantLookup", "handler", "miniAppMerchant")
	defer span.End()
	log := local_util.LoggerFromCtx(ctx, h.logger)
	_ = log
	merchantID := chi.URLParam(r, "merchant_id")
	if merchantID == "" {
		span.RecordError(errors.New("missing or invalid parameter 'merchant_id'"))
		localization.SendErrorByCodeResponse(w, localization.ErrorInvalidInputParameters.Code)
		return
	}
	span.SetAttributes(attribute.String("ecommerce_merchant_dto.id", merchantID))
	result, err := h.srv.MerchantLookup(ctx, merchantID)
	if err != nil {
		span.RecordError(err)
		localization.SendErrorByCodeResponse(w, err.Error())
		return
	}

	localization.SendSuccessResponse(w, localization.SuccessMiniAppDetailsFetched, result)
}
