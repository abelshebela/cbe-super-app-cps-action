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

func (h *miniAppMerchantAdapter) Create(w http.ResponseWriter, r *http.Request) {
	h.logger.Debugf(">>> Entered Create endpoint")
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
		localization.SendErrorResponse(w, localization.ErrorValidationFailed, nil, nil)
		return
	}

	// Extract User Context
	userContext := local_util.ExtractUserContext(r)
	if local_util.IsIncomplete(userContext) {
		h.logger.Warnf("Incomplete user context: %+v", userContext)
		localization.SendErrorResponse(w, localization.ErrorIncompleteUserInfo, nil, nil)
		return
	}

	// Convert DTO → Domain Model
	merchantDomain := ToMiniAppMerchantDomainFromUpdateDTO(&reqDTO)
	h.logger.Debugf("Converted to domain model: %+v", merchantDomain)
	// Call service to create merchant
	createdMerchant, err := h.miniappMerchantService.Create(r.Context(), merchantDomain)
	if err != nil {
		h.logger.Errorf("Failed to create merchant: %v", err)
		localization.SendErrorResponse(w, localization.ErrorMiniAppMerchantCreateFromActionFailed, nil, nil)
		return
	}
	h.logger.Debugf("Created merchant: %+v", createdMerchant)

	// Convert domain → response DTO
	responseDTO := ToMiniAppMerchantResponseDTO(createdMerchant)
	h.logger.Debugf("Response DTO: %+v", responseDTO)
	// Send success response
	localization.SendSuccessResponse(w, localization.SuccessMiniAppMerchantCreateRequestCreated, responseDTO)
}

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
	if err := reqDTO.Validate(true); err != nil {
		localization.SendErrorResponse(w, localization.ErrorValidationFailed, nil, nil)
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
	updatedMerchant, oldMerchant, err := h.miniappMerchantService.Update(r.Context(), id, merchantReq)
	if err != nil {
		localization.SendErrorByCodeResponse(w, err.Error())
		return
	}

	if oldMerchant == nil {
		localization.SendErrorByCodeResponse(w, localization.ErrorMiniAppMerchantExistsCheckFailed.Code)
		return
	}

	// Map to response DTO
	respDTO := ToMiniAppMerchantResponseDTO(updatedMerchant)

	localization.SendSuccessResponse(w, localization.SuccessMiniAppEnabledStateUpdated, respDTO)
}

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
		localization.SendErrorResponse(w, localization.ErrorMiniAppMerchantDeleteFailed, nil, nil)
		return
	}

	localization.SendSuccessResponse(w, localization.SuccessMiniAppDeletedRequestCreated, nil)
}

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
		localization.SendErrorResponse(w, localization.ErrorMiniAppMerchantEnableFailed, nil, nil)
		return
	}

	localization.SendSuccessResponse(w, localization.SuccessMiniAppEnabledStateUpdated, nil)
}

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
		localization.SendErrorResponse(w, localization.ErrorMiniAppMerchantDisableFailed, nil, nil)
		return
	}

	localization.SendSuccessResponse(w, localization.SuccessMiniAppDesable, nil)
}

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

func (h *miniAppMerchantAdapter) FindAllWithPagination(w http.ResponseWriter, r *http.Request) {
	filterParams := local_util.ExtractFilterParams(r)
	if filterParams.Page < 0 || filterParams.PerPage < 0 {
		localization.SendErrorResponse(w, localization.ErrorInvalidRequest, nil, nil)
		return
	}
	ctx := r.Context()

	miniAppMerchant, err := h.miniappMerchantService.FindAllWithPagination(ctx, filterParams)
	if err != nil {
		localization.SendErrorByCodeResponse(w, err.Error())
		return
	}

	localization.SendSuccessResponse(w, localization.SuccessMiniAppDetailsFetched, miniAppMerchant)
}
