package vaultgroupcategory

import (
	"encoding/json"
	"net/http"

	vaultgroup_category "cbe-super-app-cps-action/internal/constants/dto/vaultgroup_category"
	localization "cbe-super-app-cps-action/internal/constants/localization"
	"cbe-super-app-cps-action/internal/handlers/middleware"
	"cbe-super-app-cps-action/internal/handlers/rest/http/vaultgroup_category/core"
	"cbe-super-app-cps-action/internal/service"

	common_utils "cbe-super-app-cps-action/pkgs/utils"

	"gitlab.com/bersufekadgetachew/cbe-super-app-shared/shared/utils"
)

type handler struct {
	service service.VaultGroupCategoryService
	logger  utils.Logger
}

func InitVaultGroupCategoryHandler(svc service.VaultGroupCategoryService, logger utils.Logger) *handler {
	return &handler{service: svc, logger: logger}
}

func (h *handler) CreateVaultGroupCategory(w http.ResponseWriter, r *http.Request) {
	var req vaultgroup_category.CreateVaultGroupCategoryRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		h.logger.Errorf("[CreateVaultGroupCategory] decode: %v", err)
		localization.SendBadRequestResponse(w, localization.MsgInvalidInput)
		return
	}
	if err := req.Validate(); err != nil {
		h.logger.Errorf("[CreateVaultGroupCategory] validation: %v", err)
		localization.SendBadRequestResponse(w, err.Error())
		return
	}
	product := core.ToDomainCreateVaultGroupCategoryRequest(req)
	id, err := h.service.CreateVaultGroupCategory(r.Context(), product)
	if err != nil {
		h.logger.Errorf("[CreateVaultGroupCategory] service: %v", err)
		localization.SendErrorByCodeResponse(w, err.Error())
		return
	}
	h.logger.Infof("Vault group category created with ID: %s", id)
	localization.SendSuccessResponse(w, localization.SuccessVaultGroupCategoryCreationRequestSubmitted, nil)
}

func (h *handler) FindAllVaultGroupCategories(w http.ResponseWriter, r *http.Request) {
	params := common_utils.ExtractFilterParams(r)
	result, err := h.service.FindAllVaultGroupCategories(r.Context(), params)
	if err != nil {
		h.logger.Errorf("[FindAllVaultGroupCategories] service: %v", err)
		localization.SendErrorByCodeResponse(w, err.Error())
		return
	}
	localization.SendSuccessResponse(w, localization.SuccessVaultGroupCategoriesRetrieved, result)
}

func (h *handler) GetVaultGroupCategory(w http.ResponseWriter, r *http.Request) {
	id, err := common_utils.ExtractID(w, r)
	if id == "" {
		appErr := middleware.NewValidationError("vault group category id is required", map[string]interface{}{}).WithService("bankvault_product").
			WithOperation("GetBankVault")
		localization.SendErrorByCodeResponse(w, appErr.Error())
		return
	}
	if err != nil {
		h.logger.Errorf("[GetVaultGroupCategory] extract id: %v", err)
		localization.SendBadRequestResponse(w, err.Error())
		return
	}
	result, err := h.service.GetVaultGroupCategory(r.Context(), id)
	if err != nil {
		h.logger.Errorf("[GetVaultGroupCategory] service: %v", err)
		localization.SendErrorByCodeResponse(w, err.Error())
		return
	}
	h.logger.Infof("Vault group category retrieved with ID: %s", id)
	localization.SendSuccessResponse(w, localization.SuccessVaultGroupCategoryRetrieved, result)
}

func (h *handler) UpdateVaultGroupCategory(w http.ResponseWriter, r *http.Request) {
	id, err := common_utils.ExtractID(w, r)
	if id == "" {
		appErr := middleware.NewValidationError("vault group category id is required", map[string]interface{}{}).WithService("vault_group_category").
			WithOperation("UpdateVaultGroupCategory")
		localization.SendErrorByCodeResponse(w, appErr.Error())
		return
	}
	if err != nil {
		h.logger.Errorf("[UpdateVaultGroupCategory] extract id: %v", err)
		localization.SendBadRequestResponse(w, err.Error())
		return
	}
	var req vaultgroup_category.UpdateVaultGroupCategoryRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		h.logger.Errorf("[UpdateVaultGroupCategory] decode: %v", err)
		localization.SendBadRequestResponse(w, localization.MsgInvalidInput)
		return
	}
	if err := req.Validate(); err != nil {
		h.logger.Errorf("[UpdateVaultGroupCategory] validation: %v", err)
		localization.SendBadRequestResponse(w, err.Error())
		return
	}
	updateData := core.ToDomainUpdateVaultGroupCategoryRequest(req)
	_, err = h.service.UpdateVaultGroupCategory(r.Context(), id, updateData)
	if err != nil {
		h.logger.Errorf("[UpdateVaultGroupCategory] service: %v", err)
		localization.SendErrorByCodeResponse(w, err.Error())
		return
	}
	h.logger.Infof("Vault group category updated with ID: %s", id)
	localization.SendSuccessResponse(w, localization.SuccessVaultGroupCategoryUpdateRequestSubmitted, nil)
}

func (h *handler) DeleteVaultGroupCategory(w http.ResponseWriter, r *http.Request) {
	id, err := common_utils.ExtractID(w, r)
	if id == "" {
		appErr := middleware.NewValidationError("vault group category id is required", map[string]interface{}{}).WithService("vault_group_category").
			WithOperation("DeleteVaultGroupCategory")
		localization.SendErrorByCodeResponse(w, appErr.Error())
		return
	}
	if err != nil {
		h.logger.Errorf("[DeleteVaultGroupCategory] extract id: %v", err)
		localization.SendBadRequestResponse(w, err.Error())
		return
	}
	_, err = h.service.DeleteVaultGroupCategory(r.Context(), id)
	if err != nil {
		h.logger.Errorf("[DeleteVaultGroupCategory] service: %v", err)
		localization.SendErrorByCodeResponse(w, err.Error())
		return
	}
	h.logger.Infof("Vault group category deleted with ID: %s", id)
	localization.SendSuccessResponse(w, localization.SuccessVaultGroupCategoryDeleteRequestSubmitted, nil)
}

func (h *handler) DisableVaultGroupCategory(w http.ResponseWriter, r *http.Request) {
	id, err := common_utils.ExtractID(w, r)
	if id == "" {
		appErr := middleware.NewValidationError("vault group category id is required", map[string]interface{}{}).
			WithService("vault_group_category").
			WithOperation("DisableBankVault")
		localization.SendErrorByCodeResponse(w, appErr.Error())
		return
	}
	if err != nil {
		h.logger.Errorf("[DisableVaultGroupCategory] extract ID: %v", err)
		localization.SendErrorByCodeResponse(w, err.Error())
		return
	}
	err = h.service.DisableVaultGroupCategory(r.Context(), id)
	if err != nil {
		h.logger.Errorf("[DisableVaultGroupCategory] service: %v", err)
		localization.SendErrorByCodeResponse(w, err.Error())
		return
	}
	h.logger.Infof("Vault group category disabled with ID: %s", id)
	localization.SendSuccessResponse(w, localization.SuccessVaultGroupCategoryDisableRequestSubmitted, nil)
}

func (h *handler) EnableVaultGroupCategory(w http.ResponseWriter, r *http.Request) {
	id, err := common_utils.ExtractID(w, r)
	if id == "" {
		appErr := middleware.NewValidationError("vault group category id is required", map[string]interface{}{}).
			WithService("vault_group_category").
			WithOperation("EnableBankVault")
		localization.SendErrorByCodeResponse(w, appErr.Error())
		return
	}
	if err != nil {
		h.logger.Errorf("[EnableVaultGroupCategory] extract ID: %v", err)
		localization.SendErrorByCodeResponse(w, err.Error())
		return
	}
	err = h.service.EnableVaultGroupCategory(r.Context(), id)
	if err != nil {
		h.logger.Errorf("[EnableVaultGroupCategory] service: %v", err)
		localization.SendErrorByCodeResponse(w, err.Error())
		return
	}
	h.logger.Infof("Vault group category enabled with ID: %s", id)
	localization.SendSuccessResponse(w, localization.SuccessVaultGroupCategoryEnableRequestSubmitted, nil)
}
