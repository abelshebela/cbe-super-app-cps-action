package miniappmerchant

import (
	"encoding/json"
	"net/http"

	miniappmerchant "cbe-super-app-cps-action/internal/constants/dto/mini_app_merchant"
	miniappmerchat "cbe-super-app-cps-action/internal/constants/interfaces/mini_app_merchant"
	"cbe-super-app-cps-action/internal/constants/localization"
	"cbe-super-app-cps-action/internal/service"
	local_util "cbe-super-app-cps-action/pkgs/utils"

	"github.com/go-chi/chi/v5"
	shared_utils "gitlab.com/bersufekadgetachew/cbe-super-app-shared/shared/utils"
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
//	@Success		201				{object}	localization.StandardResponse{data=miniappmerchant.MiniAppMerchantDTO}
//	@Failure		400,401,422,500	{object}	localization.StandardResponse{data=nil}
//	@Router			/mini-app-merchants [post]
func (h *miniAppMerchantAdapter) Create(w http.ResponseWriter, r *http.Request) {
	var reqDTO miniappmerchant.MiniAppMerchantDTO

	// Decode request body
	if err := json.NewDecoder(r.Body).Decode(&reqDTO); err != nil {
		h.logger.Errorf("Failed to decode request body: %v", err)
		localization.SendErrorResponse(w, localization.ErrorMiniAppMerchantMarshalFailed, nil, nil)
		return
	}

	// Validate DTO
	if err := reqDTO.Validate(true); err != nil {
		h.logger.Errorf("Validation failed: %v", err)
		localization.SendBadRequestResponse(w, err.Error())
		return
	}

	formattedPhone := local_util.FormatPhoneNumber(reqDTO.PhoneNumber)

	reqDTO.PhoneNumber = formattedPhone
	// Extract User Context
	userContext := local_util.ExtractUserContext(r)
	if local_util.IsIncomplete(userContext) {
		h.logger.Warnf("Incomplete user context: %+v", userContext)
		localization.SendErrorResponse(w, localization.ErrorIncompleteUserInfo, nil, nil)
		return
	}

	merchantDomain := ToMiniAppMerchantDomainFromUpdateDTO(&reqDTO)

	h.logger.Debugf("Converted to domain model: %+v", merchantDomain)
	// Call service to create merchant
	createdMerchant, err := h.miniappMerchantService.Create(r.Context(), merchantDomain)
	if err != nil {
		h.logger.Errorf("Failed to create merchant: %v", err)
		localization.SendErrorByCodeResponse(w, err.Error())
		return
	}
	h.logger.Debugf("Created merchant: %+v", createdMerchant)

	// Send success response
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
//	@Router			/mini-app-merchants/{id} [put]
func (h *miniAppMerchantAdapter) Update(w http.ResponseWriter, r *http.Request) {
	id := chi.URLParam(r, "id")
	if id == "" {
		localization.SendErrorByCodeResponse(w, localization.ErrorInvalidInputParameters.Code)
		return
	}

	var reqDTO miniappmerchant.MiniAppMerchantDTO
	if err := json.NewDecoder(r.Body).Decode(&reqDTO); err != nil {
		localization.SendErrorResponse(w, localization.ErrorMiniAppMerchantMarshalFailed, nil, nil)
		return
	}

	// Check if request body is empty
	if reqDTO.IsEmpty() {
		localization.SendErrorResponse(w, localization.ErrorNoDataProvidedForUpdate, nil, nil)
		return
	}

	// Validate input fields
	if err := reqDTO.Validate(false); err != nil {
		localization.SendBadRequestResponse(w, err.Error())
		return
	}

	// Extract user context
	userContext := local_util.ExtractUserContext(r)
	if local_util.IsIncomplete(userContext) {
		localization.SendErrorResponse(w, localization.ErrorIncompleteUserInfo, nil, nil)
		return
	}

	// Convert DTO → Domain Model
	merchantReq := ToMiniAppMerchantDomainFromUpdateDTO(&reqDTO)

	// Call service update
	_, _, err := h.miniappMerchantService.Update(r.Context(), id, merchantReq)
	if err != nil {
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
	id := chi.URLParam(r, "id")
	if id == "" {
		h.logger.Errorf("missing or invalid parameter 'id'")
		localization.SendErrorResponse(w, localization.ErrorInvalidInputParameters, nil, nil)
		return
	}

	userContext := local_util.ExtractUserContext(r)
	if local_util.IsIncomplete(userContext) {
		localization.SendErrorResponse(w, localization.ErrorIncompleteUserInfo, nil, nil)
		return
	}

	err := h.miniappMerchantService.Delete(r.Context(), id)
	if err != nil {
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
	id := chi.URLParam(r, "id")
	if id == "" {
		localization.SendErrorByCodeResponse(w, localization.ErrorInvalidInputParameters.Code)
		return
	}

	userContext := local_util.ExtractUserContext(r)
	if local_util.IsIncomplete(userContext) {
		localization.SendErrorResponse(w, localization.ErrorIncompleteUserInfo, nil, nil)
		return
	}
	if err := h.miniappMerchantService.EnableOrDisable(r.Context(), id, true); err != nil {
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
	id := chi.URLParam(r, "id")
	if id == "" {
		localization.SendErrorByCodeResponse(w, localization.ErrorInvalidInputParameters.Code)
		return
	}

	userContext := local_util.ExtractUserContext(r)
	if local_util.IsIncomplete(userContext) {
		localization.SendErrorResponse(w, localization.ErrorIncompleteUserInfo, nil, nil)
		return
	}

	if err := h.miniappMerchantService.EnableOrDisable(r.Context(), id, false); err != nil {
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
//	@Success		200				{object}	localization.StandardResponse{data=miniappmerchant.MiniAppMerchantResponseDTO}
//	@Failure		400,401,404,500	{object}	localization.StandardResponse{data=nil}
//	@Router			/mini-app-merchants/{id} [get]
func (h *miniAppMerchantAdapter) FindByID(w http.ResponseWriter, r *http.Request) {
	id := chi.URLParam(r, "id")
	if id == "" {
		localization.SendErrorByCodeResponse(w, localization.ErrorInvalidInputParameters.Code)
		return
	}

	result, err := h.miniappMerchantService.FindByID(r.Context(), id)
	if err != nil {
		localization.SendErrorByCodeResponse(w, err.Error())
		return
	}

	response := ToMiniAppMerchantResponseDTO(result)
	localization.SendSuccessResponse(w, localization.SuccessMiniAppDetailsFetched, response)
}

// List Mini App Merchants with Pagination
//
//	@Summary		List Mini App Merchants
//	@Description	Retrieves a paginated list of mini app merchants
//	@Tags			MiniAppMerchant
//	@Security		BearerAuth
//	@Produce		json
//	@Param			page		query		int	false	"Page number"
//	@Param			per_page	query		int	false	"Items per page"
//	@Success		200			{object}	localization.StandardResponse{data=miniappmerchant.PaginatedMiniAppResponseResponse}
//	@Failure		400,401,500	{object}	localization.StandardResponse{data=nil}
//	@Router			/mini-app-merchants [get]
func (h *miniAppMerchantAdapter) FindAllWithPagination(w http.ResponseWriter, r *http.Request) {
	filterParams := local_util.ExtractFilterParams(r)
	if filterParams.Page < 0 || filterParams.PerPage < 0 {
		localization.SendErrorResponse(w, localization.ErrorInvalidRequest, nil, nil)
		return
	}
	ctx := r.Context()

	phoneNumber, ok := filterParams.Filters["phone_number"].(string)
	if !ok {
		h.logger.Errorf("Phone number is not present")
	}
	formattedPhone := local_util.FormatPhoneNumber(phoneNumber)

	filterParams.Filters["phone_number"] = formattedPhone

	miniAppMerchant, err := h.miniappMerchantService.FindAllWithPagination(ctx, filterParams)
	if err != nil {
		localization.SendErrorByCodeResponse(w, err.Error())
		return
	}

	localization.SendSuccessResponse(w, localization.SuccessMiniAppDetailsFetched, miniAppMerchant)
}

// Merchant Lookup
// @Summary Merchant Lookup
// @Description Retrieves a merchant by ID
// @Tags MiniAppMerchant
// @Security BearerAuth
// @Produce json
// @Param merchant_id path string true "Merchant ID"
// @Success 200 {object} localization.StandardResponse{data=merchantlookup.MerchantLookUpResponse}
// @Failure 400,401,404,500 {object} localization.StandardResponse{data=nil}
// @Router /mini-app-merchants/merchant-lookup/{merchant_id} [get]
func (h *miniAppMerchantAdapter) MerchantLookup(w http.ResponseWriter, r *http.Request) {
	merchantID := chi.URLParam(r, "merchant_id")
	if merchantID == "" {
		localization.SendErrorByCodeResponse(w, localization.ErrorInvalidInputParameters.Code)
		return
	}

	result, err := h.miniappMerchantService.MerchantLookup(r.Context(), merchantID)
	if err != nil {
		localization.SendErrorByCodeResponse(w, err.Error())
		return
	}

	localization.SendSuccessResponse(w, localization.SuccessMiniAppDetailsFetched, result)
}
