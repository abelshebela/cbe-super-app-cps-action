package vaultgroupcategory

import (
	"context"
	"encoding/json"
	"net/http"

	"cbe-super-app-cps-action/internal/constants"
	vault_category_dto "cbe-super-app-cps-action/internal/constants/dto/vault_category"
	localization "cbe-super-app-cps-action/internal/constants/localization"
	"cbe-super-app-cps-action/internal/constants/types"
	"cbe-super-app-cps-action/internal/handlers/middleware"
	"cbe-super-app-cps-action/internal/handlers/rest/http/vault_category/core"
	"cbe-super-app-cps-action/internal/service"

	common_utils "cbe-super-app-cps-action/pkgs/utils"

	"gitlab.com/bersufekadgetachew/cbe-super-app-shared/shared/utils"
	"go.opentelemetry.io/otel/attribute"
)

type handler struct {
	service service.VaultCategoryService
	logger  utils.Logger
}

func InitVaultCategoryHandler(svc service.VaultCategoryService, logger utils.Logger) *handler {
	return &handler{service: svc, logger: logger}
}

// CreateVaultGroupCatagory
//
//	@Summary		Create Vault Group Category
//	@Description	Create a new vault group category with the provided information
//	@Tags			Vault Group Category
//	@Accept			multipart/form-data
//	@Produce		json
//	@Param			name			formData	string									true	"Category name"
//	@Param			category_type	formData	string									true	"Category type (GROUP or PERSONAL)"
//	@Param			cover_image		formData	file									true	"Cover image file"
//	@Success		200		{object}	localization.StandardResponse{data=nil}				"Vault group category creation request submitted successfully"
//	@Failure		400		{object}	localization.StandardResponse{data=nil}				"Bad request"
//	@Failure		500		{object}	localization.StandardResponse{data=nil}				"Internal server error"
//	@Security		BearerAuth
//	@Router			/vaultgroupcategory/create [post]
func (h *handler) CreateVaultCategory(w http.ResponseWriter, r *http.Request) {
	ctx, span := common_utils.TraceLogger(r.Context(), "handler", "createVaultGroupCategory", "handler", "vaultGroupCategory")
	defer span.End()

	md := &types.ContextMetadata{}
	ctx = context.WithValue(ctx, constants.ContextKeyMetadata, md)

	var req vault_category_dto.CreateCategoryRequest

	file, fileHeader, err := core.ParseMultipartFormFile(r, "cover_image", 10<<20, true, h.logger)
	if err != nil {
		span.RecordError(err)
		h.logger.Errorf("[CreateVaultGroupCategory] parse multipart form file: %v", err)
		localization.SendErrorResponse(w, localization.ErrorVaultCoverImageMissedOrInvalid, nil, nil)
		return
	}
	defer file.Close()

	req.Name = r.FormValue("name")
	req.CoverImage = fileHeader
	req.InterestType = r.FormValue("interest_type")
	req.CategoryInterest = r.FormValue("category_interest")

	// Parse deadlock (form value is string, DTO field is *bool)
	deadlockStr := r.FormValue("deadlock")
	if deadlockStr != "" {
		var deadlockBool bool
		if deadlockStr == "true" || deadlockStr == "1" {
			deadlockBool = true
		} else {
			deadlockBool = false
		}
		req.Deadlock = &deadlockBool
	}

	// Parse tiers JSON (from form-data string into []CreateTierDTO)
	tiersStr := r.FormValue("tiers")
	if tiersStr != "" {
		var tiers []vault_category_dto.CreateTierDTO
		if err := json.Unmarshal([]byte(tiersStr), &tiers); err != nil {
			span.RecordError(err)
			h.logger.Errorf("[CreateVaultGroupCategory] parse tiers: %v", err)
			localization.SendBadRequestResponse(w, "Invalid tiers format")
			return
		}
		req.Tiers = tiers
	}

	if err := req.Validate(); err != nil {
		span.RecordError(err)
		h.logger.Errorf("[CreateVaultGroupCategory] validation: %v", err)
		localization.SendBadRequestResponse(w, err.Error())
		return
	}
	// product := core.ToDomainCreateVaultGroupCategoryRequest(req)
	span.SetAttributes(attribute.String("vault_group_category.name", req.Name))
	id, err := h.service.CreateVaultCategory(ctx, &req)
	if err != nil {
		span.RecordError(err)
		h.logger.Errorf("[CreateVaultGroupCategory] service: %v", err)
		localization.SendErrorByCodeResponse(w, err.Error())
		return
	}

	if md.IsMakerOnly {
		userCode, _ := ctx.Value(constants.ContextKey("user_code")).(string)
		h.logger.Infof("[Create] request sent successfully for user_code: %s is_maker_only: %v", userCode, md.IsMakerOnly)
		localization.SendSuccessResponse(w, localization.SuccessVaultCategoryCreatedSuccessfully, nil)
		return
	}

	h.logger.Infof("Vault group category created with ID: %s", id)
	localization.SendSuccessResponse(w, localization.SuccessVaultCategoryCreationRequestSubmitted, nil)
}

// FindAllVaultGroupCategories
//
//	@Summary		Find All Vault Group Categories
//	@Description	Find all vault group categories with the provided filters
//	@Tags			Vault Group Category
//	@Accept			json
//	@Produce		json
//	@Param			filters	query		string									false	"Filters"
//	@Success		200		{object}	localization.StandardResponse{data=object}	"Vault group categories retrieved successfully"
//	@Failure		400		{object}	localization.StandardResponse{data=nil}	"Bad request"
//	@Failure		500		{object}	localization.StandardResponse{data=nil}	"Internal server error"
//	@Security		BearerAuth
//	@Router			/vaultgroupcategory [get]
func (h *handler) FindAllVaultCategories(w http.ResponseWriter, r *http.Request) {
	ctx, span := common_utils.TraceLogger(r.Context(), "handler", "findAllVaultGroupCategories", "handler", "vaultGroupCategory")
	defer span.End()

	params := common_utils.ExtractFilterParams(r)

	search := r.URL.Query().Get("search")
	filter := r.URL.Query().Get("filter")

	if err := common_utils.NoSpecialChars(search); err != nil {
		localization.SendErrorByCodeResponse(w, err.Error())
		return
	}

	if err := common_utils.NoSpecialChars(filter); err != nil {
		localization.SendErrorByCodeResponse(w, err.Error())
		return
	}

	result, err := h.service.FindAllVaultCategories(ctx, params)
	if err != nil {
		span.RecordError(err)
		h.logger.Errorf("[FindAllVaultGroupCategories] service: %v", err)
		localization.SendErrorByCodeResponse(w, err.Error())
		return
	}

	span.SetAttributes(attribute.Int("vault_group_category.count", len(result.Data)))
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
//	@Success		200	{object}	localization.StandardResponse{data=vaultgroupcategory.VaultGroupCategoryResponse}	"Vault group category retrieved successfully"
//	@Failure		400	{object}	localization.StandardResponse{data=nil}	"Bad request"
//	@Failure		500	{object}	localization.StandardResponse{data=nil}	"Internal server error"
//	@Security		BearerAuth
//	@Router			/vaultgroupcategory/{id} [get]
func (h *handler) GetVaultCategory(w http.ResponseWriter, r *http.Request) {
	ctx, span := common_utils.TraceLogger(r.Context(), "handler", "getVaultGroupCategory", "handler", "vaultGroupCategory")
	defer span.End()
	id, err := common_utils.ExtractID(w, r)
	if id == "" {
		appErr := middleware.NewValidationError("vault group category id is required", map[string]interface{}{}).WithService("bankvault_product").
			WithOperation("GetBankVault")
		localization.SendErrorByCodeResponse(w, appErr.Error())
		return
	}
	if err != nil {
		span.RecordError(err)
		h.logger.Errorf("[GetVaultGroupCategory] extract id: %v", err)
		localization.SendBadRequestResponse(w, err.Error())
		return
	}
	span.SetAttributes(attribute.String("vault_group_category.id", id))
	result, err := h.service.GetVaultCategory(ctx, id)
	if err != nil {
		span.RecordError(err)
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
//	@Description	Update a vault group category by the provided ID. Expects multipart/form-data (optional file "cover_image" and form field "name").
//	@Tags			Vault Group Category
//	@Accept			multipart/form-data
//	@Produce		json
//	@Param			id			path		string									true	"Vault group category ID"
//	@Param			cover_image	formData	file									false	"Cover image file"
//	@Param			name		formData	string									false	"Name"
//	@Success		200			{object}	localization.StandardResponse{data=nil}	"Vault group category update request submitted successfully"
//	@Failure		400			{object}	localization.StandardResponse{data=nil}	"Bad request"
//	@Failure		500			{object}	localization.StandardResponse{data=nil}	"Internal server error"
//	@Security		BearerAuth
//	@Router			/vaultgroupcategory/update/{id} [patch]
func (h *handler) UpdateVaultCategory(w http.ResponseWriter, r *http.Request) {
	ctx, span := common_utils.TraceLogger(r.Context(), "handler", "updateVaultGroupCategory", "handler", "vaultGroupCategory")
	defer span.End()

	md := &types.ContextMetadata{}
	ctx = context.WithValue(ctx, constants.ContextKeyMetadata, md)

	id, err := common_utils.ExtractID(w, r)
	if id == "" {
		appErr := middleware.NewValidationError("vault group category id is required", map[string]interface{}{}).WithService("vault_group_category").
			WithOperation("UpdateVaultGroupCategory")
		localization.SendErrorByCodeResponse(w, appErr.Error())
		return
	}
	if err != nil {
		span.RecordError(err)
		h.logger.Errorf("[UpdateVaultGroupCategory] extract id: %v", err)
		localization.SendBadRequestResponse(w, err.Error())
		return
	}
	var req vault_category_dto.UpdateCategoryRequest

	file, fileHeader, err := core.ParseMultipartFormFile(r, "cover_image", 10<<20, false, h.logger)
	if err != nil {
		span.RecordError(err)
		h.logger.Errorf("[CreateVaultGroupCategory] parse multipart form file: %v", err)
		localization.SendErrorResponse(w, localization.ErrorVaultCoverImageMissedOrInvalid, nil, nil)
		return
	}
	if file != nil {
		defer file.Close()
		req.CoverImage = fileHeader
	}
	name := r.FormValue("name")
	intType := r.FormValue("interest_type")
	cateInt := r.FormValue("category_interest")
	deadlockStr := r.FormValue("deadlock")
	tiersStr := r.FormValue("tiers")

	if name != "" {
		req.Name = &name
	}
	if intType != "" {
		req.InterestType = &intType
	}
	if cateInt != "" {
		req.CategoryInterest = &cateInt
	}

	if deadlockStr != "" {
		var deadlockBool bool
		if deadlockStr == "true" || deadlockStr == "1" {
			deadlockBool = true
		} else {
			deadlockBool = false
		}
		req.Deadlock = &deadlockBool
	}

	if tiersStr != "" {
		var tier []vault_category_dto.UpdateTierDTO
		if err := json.Unmarshal([]byte(tiersStr), &tier); err != nil {
			span.RecordError(err)
			h.logger.Errorf("[UpdateVaultCategory] parse tiers: %v", err)
			localization.SendBadRequestResponse(w, "Invalid tiers format")
			return
		}
		req.Tiers = tier
	}

	if err := req.Validate(); err != nil {
		span.RecordError(err)
		h.logger.Errorf("[UpdateVaultGroupCategory] validation: %v", err)
		localization.SendBadRequestResponse(w, err.Error())
		return
	}

	_, err = h.service.UpdateVaultCategory(ctx, id, &req)
	if err != nil {
		span.RecordError(err)
		h.logger.Errorf("[UpdateVaultGroupCategory] service: %v", err)
		localization.SendErrorByCodeResponse(w, err.Error())
		return
	}

	if md.IsMakerOnly {
		userCode, _ := ctx.Value(constants.ContextKey("user_code")).(string)
		h.logger.Infof("[Create] request sent successfully for user_code: %s is_maker_only: %v", userCode, md.IsMakerOnly)
		localization.SendSuccessResponse(w, localization.SuccessVaultCategoryUpdatedSuccessfully, nil)
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
func (h *handler) DeleteVaultCategory(w http.ResponseWriter, r *http.Request) {
	ctx, span := common_utils.TraceLogger(r.Context(), "handler", "deleteVaultGroupCategory", "handler", "vaultGroupCategory")
	defer span.End()
	id, err := common_utils.ExtractID(w, r)
	if id == "" {
		appErr := middleware.NewValidationError("vault group category id is required", map[string]interface{}{}).WithService("vault_group_category").
			WithOperation("DeleteVaultGroupCategory")
		span.RecordError(appErr)
		localization.SendErrorByCodeResponse(w, appErr.Error())
		return
	}
	if err != nil {
		span.RecordError(err)
		h.logger.Errorf("[DeleteVaultGroupCategory] extract id: %v", err)
		localization.SendBadRequestResponse(w, err.Error())
		return
	}
	span.SetAttributes(attribute.String("vault_group_category.id", id))
	_, err = h.service.DeleteVaultCategory(ctx, id)
	if err != nil {
		span.RecordError(err)
		h.logger.Errorf("[DeleteVaultGroupCategory] service: %v", err)
		localization.SendErrorByCodeResponse(w, err.Error())
		return
	}
	h.logger.Infof("Vault group category deleted with ID: %s", id)
	localization.SendSuccessResponse(w, localization.SuccessVaultGroupCategoryDeleteRequestSubmitted, nil)
}

// DisableVaultCategory
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
func (h *handler) DisableVaultCategory(w http.ResponseWriter, r *http.Request) {
	ctx, span := common_utils.TraceLogger(r.Context(), "handler", "disableVaultGroupCategory", "handler", "vaultGroupCategory")
	defer span.End()
	id, err := common_utils.ExtractID(w, r)
	if id == "" {
		appErr := middleware.NewValidationError("vault group category id is required", map[string]interface{}{}).
			WithService("vault_group_category").
			WithOperation("DisableBankVault")
		span.RecordError(appErr)
		localization.SendErrorByCodeResponse(w, appErr.Error())
		return
	}
	if err != nil {
		span.RecordError(err)
		h.logger.Errorf("[DisableVaultGroupCategory] extract ID: %v", err)
		localization.SendErrorByCodeResponse(w, err.Error())
		return
	}
	span.SetAttributes(attribute.String("vault_group_category.id", id))
	err = h.service.DisableVaultCategory(ctx, id)
	if err != nil {
		span.RecordError(err)
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
func (h *handler) EnableVaultCategory(w http.ResponseWriter, r *http.Request) {
	ctx, span := common_utils.TraceLogger(r.Context(), "handler", "enableVaultGroupCategory", "handler", "vaultGroupCategory")
	defer span.End()
	id, err := common_utils.ExtractID(w, r)
	if id == "" {
		appErr := middleware.NewValidationError("vault group category id is required", map[string]interface{}{}).
			WithService("vault_group_category").
			WithOperation("EnableBankVault")
		span.RecordError(appErr)
		localization.SendErrorByCodeResponse(w, appErr.Error())
		return
	}
	if err != nil {
		span.RecordError(err)
		h.logger.Errorf("[EnableVaultGroupCategory] extract ID: %v", err)
		localization.SendErrorByCodeResponse(w, err.Error())
		return
	}
	span.SetAttributes(attribute.String("vault_group_category.id", id))
	err = h.service.EnableVaultCategory(ctx, id)
	if err != nil {
		span.RecordError(err)
		h.logger.Errorf("[EnableVaultGroupCategory] service: %v", err)
		localization.SendErrorByCodeResponse(w, err.Error())
		return
	}
	h.logger.Infof("Vault group category enabled with ID: %s", id)
	localization.SendSuccessResponse(w, localization.SuccessVaultGroupCategoryEnableRequestSubmitted, nil)
}
