package bankvault

import (
	"encoding/json"
	"net/http"

	"cbe-super-app-cps-action/internal/constants/dto/bankvault"
	localization "cbe-super-app-cps-action/internal/constants/localization"
	"cbe-super-app-cps-action/internal/handlers/middleware"
	"cbe-super-app-cps-action/internal/handlers/rest/http/bankvault/core"
	"cbe-super-app-cps-action/internal/service"

	common_utils "cbe-super-app-cps-action/pkgs/utils"

	"gitlab.com/bersufekadgetachew/cbe-super-app-shared/shared/utils"
)

type handler struct {
	service service.BankVaultService
	logger  utils.Logger
}

func InitBankVaultHandler(svc service.BankVaultService, logger utils.Logger) *handler {
	return &handler{service: svc, logger: logger}
}

func (h *handler) CreateBankVault(w http.ResponseWriter, r *http.Request) {
	var req bankvault.CreateBankVaultProductRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		h.logger.Errorf("[CreateBankVault] decode: %v", err)
		localization.SendBadRequestResponse(w, localization.MsgInvalidInput)
		return
	}
	if err := req.Validate(); err != nil {
		h.logger.Errorf("[CreateBankVault] validation: %v", err)
		localization.SendBadRequestResponse(w, err.Error())
		return
	}
	lockPeriod, err := common_utils.ParseLockPeriod(req.LockPeriodDays)
	if err != nil {
		h.logger.Errorf("[CreateBankVault] lock period parse: %v", err)
		localization.SendBadRequestResponse(w, err.Error())
		return
	}
	product := core.ToDomainCreateBankVaultRequest(req, lockPeriod)
	id, err := h.service.CreateBankVault(r.Context(), product)
	if err != nil {
		h.logger.Errorf("[CreateBankVault] service: %v", err)
		localization.SendErrorByCodeResponse(w, err.Error())
		return
	}
	h.logger.Infof("Bank vault created with ID: %s", id)
	localization.SendSuccessResponse(w, localization.SuccessBankVaultCreationRequestSubmitted, nil)

}

func (h *handler) FindAllBankVaults(w http.ResponseWriter, r *http.Request) {
	params := common_utils.ExtractFilterParams(r)
	result, err := h.service.FindAllBankVaults(r.Context(), params)
	if err != nil {
		h.logger.Errorf("[FindAllBankVaults] service: %v", err)
		localization.SendErrorByCodeResponse(w, err.Error())
		return
	}
	localization.SendSuccessResponse(w, localization.SuccessBankVaultsRetrieved, result)
}
func (h *handler) GetBankVault(w http.ResponseWriter, r *http.Request) {
	id, err := common_utils.ExtractID(w, r)
	if id == "" {
		appErr := middleware.NewValidationError("Product ID is required", map[string]interface{}{}).WithService("bankvault_product").
			WithOperation("GetBankVault")
		localization.SendErrorByCodeResponse(w, appErr.Error())
		return
	}

	if err != nil {
		h.logger.Errorf("[GetBankVault] extract ID: %v", err)
		localization.SendErrorByCodeResponse(w, err.Error())
		return
	}

	result, err := h.service.GetBankVault(r.Context(), id)
	if err != nil {
		h.logger.Errorf("[GetBankVault] service: %v", err)
		localization.SendErrorByCodeResponse(w, err.Error())
		return
	}
	h.logger.Infof("Bank vault retrieved with ID: %s", id)
	localization.SendSuccessResponse(w, localization.SuccessBankVaultRetrieved, result)
}

func (h *handler) UpdateBankVault(w http.ResponseWriter, r *http.Request) {
	id, err := common_utils.ExtractID(w, r)
	if id == "" {
		appErr := middleware.NewValidationError("Product ID is required", map[string]interface{}{}).
			WithService("bankvault_product").
			WithOperation("UpdateBankVault")
		localization.SendErrorByCodeResponse(w, appErr.Error())
		return
	}

	if err != nil {
		h.logger.Errorf("[UpdateBankVault] extract ID: %v", err)
		localization.SendErrorByCodeResponse(w, err.Error())
		return
	}

	var req bankvault.UpdateBankVaultProductRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		h.logger.Errorf("[UpdateBankVault] decode: %v", err)
		appErr := middleware.NewValidationError("Invalid JSON input", map[string]interface{}{}).
			WithService("bankvault_product").
			WithOperation("UpdateBankVault")
		localization.SendErrorByCodeResponse(w, appErr.Error())
		return
	}

	if err := req.Validate(); err != nil {
		h.logger.Errorf("[UpdateBankVault] validation: %v", err)
		appErr := middleware.NewValidationError(err.Error(), map[string]interface{}{}).
			WithService("bankvault_product").
			WithOperation("UpdateBankVault")
		localization.SendErrorByCodeResponse(w, appErr.Error())
		return
	}

	product := core.ToDomainUpdateBankVaultRequest(req, nil)

	_, err = h.service.UpdateBankVault(r.Context(), id, product)
	if err != nil {
		h.logger.Errorf("[UpdateBankVault] service: %v", err)
		localization.SendErrorByCodeResponse(w, err.Error())
		return
	}

	h.logger.Infof("Bank vault updated with ID: %s", id)
	localization.SendSuccessResponse(w, localization.SuccessBankVaultUpdateRequestSubmitted, nil)
}

func (h *handler) DeleteBankVault(w http.ResponseWriter, r *http.Request) {
	id, err := common_utils.ExtractID(w, r)
	if id == "" {
		appErr := middleware.NewValidationError("Product ID is required", map[string]interface{}{}).
			WithService("bankvault_product").
			WithOperation("DeleteBankVault")
		localization.SendErrorByCodeResponse(w, appErr.Error())
		return
	}

	if err != nil {
		h.logger.Errorf("[DeleteBankVault] extract ID: %v", err)
		localization.SendErrorByCodeResponse(w, err.Error())
		return
	}

	_, err = h.service.DeleteBankVault(r.Context(), id)
	if err != nil {
		h.logger.Errorf("[DeleteBankVault] service: %v", err)
		localization.SendErrorByCodeResponse(w, err.Error())
		return

	}

	h.logger.Infof("Bank vault deleted with ID: %s", id)
	localization.SendSuccessResponse(w, localization.SuccessBankVaultDeleteRequestSubmitted, nil)
}

func (h *handler) DisableBankVault(w http.ResponseWriter, r *http.Request) {
	id, err := common_utils.ExtractID(w, r)
	if id == "" {
		appErr := middleware.NewValidationError("Product ID is required", map[string]interface{}{}).
			WithService("bankvault_product").
			WithOperation("DisableBankVault")
		localization.SendErrorByCodeResponse(w, appErr.Error())
		return
	}

	if err != nil {
		h.logger.Errorf("[DisableBankVault] extract ID: %v", err)
		localization.SendErrorByCodeResponse(w, err.Error())
		return
	}

	err = h.service.DisableBankVault(r.Context(), id)
	if err != nil {
		h.logger.Errorf("[DisableBankVault] service: %v", err)
		localization.SendErrorByCodeResponse(w, err.Error())
		return
	}

	h.logger.Infof("Bank vault disabled with ID: %s", id)
	localization.SendSuccessResponse(w, localization.SuccessBankVaultDisableRequestSubmitted, nil)
}
func (h *handler) EnableBankVault(w http.ResponseWriter, r *http.Request) {
	id, err := common_utils.ExtractID(w, r)
	if id == "" {
		appErr := middleware.NewValidationError("Product ID is required", map[string]interface{}{}).
			WithService("bankvault_product").
			WithOperation("EnableBankVault")
		localization.SendErrorByCodeResponse(w, appErr.Error())
		return
	}

	if err != nil {
		h.logger.Errorf("[EnableBankVault] extract ID: %v", err)
		localization.SendErrorByCodeResponse(w, err.Error())
		return
	}

	err = h.service.EnableBankVault(r.Context(), id)
	if err != nil {
		h.logger.Errorf("[EnableBankVault] service: %v", err)
		localization.SendErrorByCodeResponse(w, err.Error())
		return
	}

	h.logger.Infof("Bank vault enabled with ID: %s", id)
	localization.SendSuccessResponse(w, localization.SuccessBankVaultEnableRequestSubmitted, nil)
}
