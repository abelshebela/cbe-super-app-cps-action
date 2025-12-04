package vaultgroupcategory

import (
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

// CreateVaultGroupCatagory
//
//	@Summary		Create Vault Group Category
//	@Description	Create a new vault group category with the provided information
//	@Tags			Vault Group Category
//	@Accept			json
//	@Produce		json
//	@Param			request	body		vaultgroupcategory.CreateVaultGroupCategoryRequest	true	"Vault group category request"
//	@Success		200		{object}	localization.StandardResponse{data=nil}				"Vault group category creation request submitted successfully"
//	@Failure		400		{object}	localization.StandardResponse{data=nil}				"Bad request"
//	@Failure		500		{object}	localization.StandardResponse{data=nil}				"Internal server error"
//	@Security		BearerAuth
//	@Router			/vaultgroupcategory/create [post]
func (h *handler) CreateVaultGroupCategory(w http.ResponseWriter, r *http.Request) {
	var req vaultgroup_category.CreateVaultGroupCategoryRequest

	file, fileHeader, err := core.ParseMultipartFormFile(r, "cover_image", 10<<20, true, h.logger)
	if err != nil {
		h.logger.Errorf("[CreateVaultGroupCategory] parse multipart form file: %v", err)
		localization.SendErrorResponse(w, localization.ErrorVaultCoverImageMissedOrInvalid, nil, nil)
		return
	}
	defer file.Close()

	req.Name = r.FormValue("name")
	req.CoverImage = fileHeader

	if err := req.Validate(); err != nil {
		h.logger.Errorf("[CreateVaultGroupCategory] validation: %v", err)
		localization.SendBadRequestResponse(w, err.Error())
		return
	}
	// product := core.ToDomainCreateVaultGroupCategoryRequest(req)
	id, err := h.service.CreateVaultGroupCategory(r.Context(), &req)
	if err != nil {
		h.logger.Errorf("[CreateVaultGroupCategory] service: %v", err)
		localization.SendErrorByCodeResponse(w, err.Error())
		return
	}
	h.logger.Infof("Vault group category created with ID: %s", id)
	localization.SendSuccessResponse(w, localization.SuccessVaultGroupCategoryCreationRequestSubmitted, nil)
}

// FindAllVaultGroupCategories
//
//	@Summary		Find All Vault Group Categories
//	@Description	Find all vault group categories with the provided filters
//	@Tags			Vault Group Category
//	@Accept			json
//	@Produce		json
//	@Param			filters	query		string									false	"Filters"
//	@Success		200		{object}	localization.StandardResponse{data=nil}	"Vault group categories retrieved successfully"
//	@Failure		400		{object}	localization.StandardResponse{data=nil}	"Bad request"
//	@Failure		500		{object}	localization.StandardResponse{data=nil}	"Internal server error"
//	@Security		BearerAuth
//	@Router			/vaultgroupcategory [get]
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

// GetVaultGroupCategory
//
//	@Summary		Get Vault Group Category
//	@Description	Get a vault group category by the provided ID
//	@Tags			Vault Group Category
//	@Accept			json
//	@Produce		json
//	@Param			id	path		string									true	"Vault group category ID"
//	@Success		200	{object}	localization.StandardResponse{data=nil}	"Vault group category retrieved successfully"
//	@Failure		400	{object}	localization.StandardResponse{data=nil}	"Bad request"
//	@Failure		500	{object}	localization.StandardResponse{data=nil}	"Internal server error"
//	@Security		BearerAuth
//	@Router			/vaultgroupcategory/{id} [get]
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

// UpdateVaultGroupCategory
//
//	@Summary		Update Vault Group Category
//	@Description	Update a vault group category by the provided ID
//	@Tags			Vault Group Category
//	@Accept			json
//	@Produce		json
//	@Param			id		path		string												true	"Vault group category ID"
//	@Param			request	body		vaultgroupcategory.UpdateVaultGroupCategoryRequest	true	"Vault group category request"
//	@Success		200		{object}	localization.StandardResponse{data=nil}				"Vault group category update request submitted successfully"
//	@Failure		400		{object}	localization.StandardResponse{data=nil}				"Bad request"
//	@Failure		500		{object}	localization.StandardResponse{data=nil}				"Internal server error"
//	@Security		BearerAuth
//	@Router			/vaultgroupcategory/update/{id} [patch]
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

	file, fileHeader, err := core.ParseMultipartFormFile(r, "cover_image", 10<<20, false, h.logger)
	if err != nil {
		h.logger.Errorf("[CreateVaultGroupCategory] parse multipart form file: %v", err)
		localization.SendErrorResponse(w, localization.ErrorVaultCoverImageMissedOrInvalid, nil, nil)
		return
	}
	if file != nil {
		defer file.Close()
		req.CoverImage = fileHeader
	}

	if name := r.FormValue("name"); name != "" {
		req.Name = &name
	}

	if err := req.Validate(); err != nil {
		h.logger.Errorf("[UpdateVaultGroupCategory] validation: %v", err)
		localization.SendBadRequestResponse(w, err.Error())
		return
	}
	// updateData := core.ToDomainUpdateVaultGroupCategoryRequest(req)
	_, err = h.service.UpdateVaultGroupCategory(r.Context(), id, &req)
	if err != nil {
		h.logger.Errorf("[UpdateVaultGroupCategory] service: %v", err)
		localization.SendErrorByCodeResponse(w, err.Error())
		return
	}
	h.logger.Infof("Vault group category updated with ID: %s", id)
	localization.SendSuccessResponse(w, localization.SuccessVaultGroupCategoryUpdateRequestSubmitted, nil)
}

// DeleteVaultGroupCategory
//
//	@Summary		Delete Vault Group Category
//	@Description	Delete a vault group category by the provided ID
//	@Tags			Vault Group Category
//	@Accept			json
//	@Produce		json
//	@Param			id	path		string									true	"Vault group category ID"
//	@Success		200	{object}	localization.StandardResponse{data=nil}	"Vault group category delete request submitted successfully"
//	@Failure		400	{object}	localization.StandardResponse{data=nil}	"Bad request"
//	@Failure		500	{object}	localization.StandardResponse{data=nil}	"Internal server error"
//	@Security		BearerAuth
//	@Router			/vaultgroupcategory/delete/{id} [delete]
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

// DisableVaultGroupCategory
//
//	@Summary		Disable Vault Group Category
//	@Description	Disable a vault group category by the provided ID
//	@Tags			Vault Group Category
//	@Accept			json
//	@Produce		json
//	@Param			id	path		string									true	"Vault group category ID"
//	@Success		200	{object}	localization.StandardResponse{data=nil}	"Vault group category disable request submitted successfully"
//	@Failure		400	{object}	localization.StandardResponse{data=nil}	"Bad request"
//	@Failure		500	{object}	localization.StandardResponse{data=nil}	"Internal server error"
//	@Security		BearerAuth
//	@Router			/vaultgroupcategory/disable/{id} [patch]
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

// EnableVaultGroupCategory
//
//	@Summary		Enable Vault Group Category
//	@Description	Enable a vault group category by the provided ID
//	@Tags			Vault Group Category
//	@Accept			json
//	@Produce		json
//	@Param			id	path		string									true	"Vault group category ID"
//	@Success		200	{object}	localization.StandardResponse{data=nil}	"Vault group category enable request submitted successfully"
//	@Failure		400	{object}	localization.StandardResponse{data=nil}	"Bad request"
//	@Failure		500	{object}	localization.StandardResponse{data=nil}	"Internal server error"
//	@Security		BearerAuth
//	@Router			/vaultgroupcategory/enable/{id} [patch]
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
