package vault

import (
	"context"
	"encoding/json"
	"net/http"

	"cbe-super-app-cps-action/internal/constants"
	vault_category_dto "cbe-super-app-cps-action/internal/constants/dto/vault"
	localization "cbe-super-app-cps-action/internal/constants/localization"
	"cbe-super-app-cps-action/internal/constants/types"
	"cbe-super-app-cps-action/internal/handlers/middleware"
	"cbe-super-app-cps-action/internal/handlers/rest/http/vault/core"
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

// CreateVaultCatagory
//
//	@Summary		Create Vault Category
//	@Description	Create a new vault category with the provided information
//	@Tags			Vault Category
//	@Accept			multipart/form-data
//	@Produce		json
//	@Param			name			formData	string									true	"Category name"
//	@Param			cover_image		formData	file									true	"Cover image file"
//	@Success		200		{object}	localization.StandardResponse{data=nil}				"Vault category creation request submitted successfully"
//	@Failure		400		{object}	localization.StandardResponse{data=nil}				"Bad request"
//	@Failure		500		{object}	localization.StandardResponse{data=nil}				"Internal server error"
//	@Security		BearerAuth
//	@Router			/vaultcategory/create [post]
func (h *handler) CreateVaultCategory(w http.ResponseWriter, r *http.Request) {
	ctx, span := common_utils.TraceLogger(r.Context(), "handler", "createVaultCategory", "handler", "vaultCategory")
	defer span.End()
	log := common_utils.LoggerFromCtx(ctx, h.logger)

	md := &types.ContextMetadata{}
	ctx = context.WithValue(ctx, constants.ContextKeyMetadata, md)

	var req vault_category_dto.CreateCategoryRequest

	file, fileHeader, err := core.ParseMultipartFormFile(r, "cover_image", 10<<20, true, h.logger)
	if err != nil {
		span.RecordError(err)
		log.Errorf("[CreateVaultCategory] parse multipart form file: %v", err)
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
			log.Errorf("[CreateVaultCategory] parse tiers: %v", err)
			localization.SendBadRequestResponse(w, "Invalid tiers format")
			return
		}
		req.Tiers = tiers
	}

	if err := req.Validate(); err != nil {
		span.RecordError(err)
		log.Errorf("[CreateVaultCategory] validation: %v", err)
		localization.SendBadRequestResponse(w, err.Error())
		return
	}
	// product := core.ToDomainCreateVaultCategoryRequest(req)
	span.SetAttributes(attribute.String("vault_category.name", req.Name))
	id, err := h.service.CreateVaultCategory(ctx, &req)
	if err != nil {
		span.RecordError(err)
		log.Errorf("[CreateVaultCategory] service: %v", err)
		localization.SendErrorByCodeResponse(w, err.Error())
		return
	}

	if md.IsMakerOnly {
		userCode, _ := ctx.Value(constants.ContextKey("user_code")).(string)
		log.Infof("[Create] request sent successfully for user_code: %s is_maker_only: %v", userCode, md.IsMakerOnly)
		localization.SendSuccessResponse(w, localization.SuccessVaultCategoryCreatedSuccessfully, nil)
		return
	}

	log.Infof("[VaultCatH][Create] ok id: %s", id)
	localization.SendSuccessResponse(w, localization.SuccessVaultCategoryCreationRequestSubmitted, nil)
}

// FindAllVaultCategories
//
//	@Summary		Find All Vault Categories
//	@Description	Find all vault categories with the provided filters
//	@Tags			Vault Category
//	@Accept			json
//	@Produce		json
//	@Param			filters	query		string									false	"Filters"
//	@Success		200		{object}	localization.StandardResponse{data=object}	"Vault categories retrieved successfully"
//	@Failure		400		{object}	localization.StandardResponse{data=nil}	"Bad request"
//	@Failure		500		{object}	localization.StandardResponse{data=nil}	"Internal server error"
//	@Security		BearerAuth
//	@Router			/vaultcategory [get]
func (h *handler) FindAllVaultCategories(w http.ResponseWriter, r *http.Request) {
	ctx, span := common_utils.TraceLogger(r.Context(), "handler", "findAllVaultCategories", "handler", "vaultCategory")
	defer span.End()
	log := common_utils.LoggerFromCtx(ctx, h.logger)

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
		log.Errorf("[FindAllVaultCategories] service: %v", err)
		localization.SendErrorByCodeResponse(w, err.Error())
		return
	}

	span.SetAttributes(attribute.Int("vault_category.count", len(result.Data)))
	localization.SendSuccessResponse(w, localization.SuccessVaultCategoriesRetrieved, result)
}

// GetVaultCategory
//
//	@Summary		Get Vault Category
//	@Description	Get a vault category by the provided ID
//	@Tags			Vault Category
//	@Accept			json
//	@Produce		json
//	@Param			id	path		string									true	"Vault category ID"
//	@Success		200	{object}	localization.StandardResponse{data=vaultcategory.VaultCategoryResponse}	"Vault category retrieved successfully"
//	@Failure		400	{object}	localization.StandardResponse{data=nil}	"Bad request"
//	@Failure		500	{object}	localization.StandardResponse{data=nil}	"Internal server error"
//	@Security		BearerAuth
//	@Router			/vaultcategory/{id} [get]
func (h *handler) GetVaultCategory(w http.ResponseWriter, r *http.Request) {
	ctx, span := common_utils.TraceLogger(r.Context(), "handler", "getVaultCategory", "handler", "vaultCategory")
	defer span.End()
	log := common_utils.LoggerFromCtx(ctx, h.logger)
	id, err := common_utils.ExtractID(w, r)
	if id == "" {
		appErr := middleware.NewValidationError("vault category id is required", map[string]interface{}{}).WithService("bankvault_product").
			WithOperation("GetBankVault")
		localization.SendErrorByCodeResponse(w, appErr.Error())
		return
	}
	if err != nil {
		span.RecordError(err)
		log.Errorf("[GetVaultCategory] extract id: %v", err)
		localization.SendBadRequestResponse(w, err.Error())
		return
	}
	span.SetAttributes(attribute.String("vault_category.id", id))
	result, err := h.service.GetVaultCategory(ctx, id)
	if err != nil {
		span.RecordError(err)
		log.Errorf("[GetVaultCategory] service: %v", err)
		localization.SendErrorByCodeResponse(w, err.Error())
		return
	}
	log.Infof("[VaultCatH][GetByID] ok id: %s", id)
	localization.SendSuccessResponse(w, localization.SuccessVaultCategoryRetrieved, result)

}

// UpdateVaultCategory
//
//	@Summary		Update Vault Category
//	@Description	Update a vault category by the provided ID. Expects multipart/form-data (optional file "cover_image" and form field "name").
//	@Tags			Vault Category
//	@Accept			multipart/form-data
//	@Produce		json
//	@Param			id			path		string									true	"Vault category ID"
//	@Param			cover_image	formData	file									false	"Cover image file"
//	@Param			name		formData	string									false	"Name"
//	@Success		200			{object}	localization.StandardResponse{data=nil}	"Vault category update request submitted successfully"
//	@Failure		400			{object}	localization.StandardResponse{data=nil}	"Bad request"
//	@Failure		500			{object}	localization.StandardResponse{data=nil}	"Internal server error"
//	@Security		BearerAuth
//	@Router			/vaultcategory/update/{id} [patch]
func (h *handler) UpdateVaultCategory(w http.ResponseWriter, r *http.Request) {
	ctx, span := common_utils.TraceLogger(r.Context(), "handler", "updateVaultCategory", "handler", "vaultCategory")
	defer span.End()
	log := common_utils.LoggerFromCtx(ctx, h.logger)

	md := &types.ContextMetadata{}
	ctx = context.WithValue(ctx, constants.ContextKeyMetadata, md)

	id, err := common_utils.ExtractID(w, r)
	if id == "" {
		appErr := middleware.NewValidationError("vault category id is required", map[string]interface{}{}).WithService("vault_category").
			WithOperation("UpdateVaultCategory")
		localization.SendErrorByCodeResponse(w, appErr.Error())
		return
	}
	if err != nil {
		span.RecordError(err)
		log.Errorf("[UpdateVaultCategory] extract id: %v", err)
		localization.SendBadRequestResponse(w, err.Error())
		return
	}
	var req vault_category_dto.UpdateCategoryRequest

	file, fileHeader, err := core.ParseMultipartFormFile(r, "cover_image", 10<<20, false, h.logger)
	if err != nil {
		span.RecordError(err)
		log.Errorf("[CreateVaultCategory] parse multipart form file: %v", err)
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
			log.Errorf("[UpdateVaultCategory] parse tiers: %v", err)
			localization.SendBadRequestResponse(w, "Invalid tiers format")
			return
		}
		req.Tiers = tier
	}

	if err := req.Validate(); err != nil {
		span.RecordError(err)
		log.Errorf("[UpdateVaultCategory] validation: %v", err)
		localization.SendBadRequestResponse(w, err.Error())
		return
	}

	_, err = h.service.UpdateVaultCategory(ctx, id, &req)
	if err != nil {
		span.RecordError(err)
		log.Errorf("[UpdateVaultCategory] service: %v", err)
		localization.SendErrorByCodeResponse(w, err.Error())
		return
	}

	if md.IsMakerOnly {
		userCode, _ := ctx.Value(constants.ContextKey("user_code")).(string)
		log.Infof("[Create] request sent successfully for user_code: %s is_maker_only: %v", userCode, md.IsMakerOnly)
		localization.SendSuccessResponse(w, localization.SuccessVaultCategoryUpdatedSuccessfully, nil)
		return
	}

	log.Infof("[VaultCatH][Update] ok id: %s", id)
	localization.SendSuccessResponse(w, localization.SuccessVaultCategoryUpdateRequestSubmitted, nil)
}

// DeleteVaultCategory
//
//	@Summary		Delete Vault Category
//	@Description	Delete a vault category by the provided ID
//	@Tags			Vault Category
//	@Accept			json
//	@Produce		json
//	@Param			id	path		string									true	"Vault category ID"
//	@Success		200	{object}	localization.StandardResponse{data=nil}	"Vault category delete request submitted successfully"
//	@Failure		400	{object}	localization.StandardResponse{data=nil}	"Bad request"
//	@Failure		500	{object}	localization.StandardResponse{data=nil}	"Internal server error"
//	@Security		BearerAuth
//	@Router			/vaultcategory/delete/{id} [delete]
func (h *handler) DeleteVaultCategory(w http.ResponseWriter, r *http.Request) {
	ctx, span := common_utils.TraceLogger(r.Context(), "handler", "deleteVaultCategory", "handler", "vaultCategory")
	defer span.End()
	log := common_utils.LoggerFromCtx(ctx, h.logger)
	id, err := common_utils.ExtractID(w, r)
	if id == "" {
		appErr := middleware.NewValidationError("vault category id is required", map[string]interface{}{}).WithService("vault_category").
			WithOperation("DeleteVaultCategory")
		span.RecordError(appErr)
		localization.SendErrorByCodeResponse(w, appErr.Error())
		return
	}
	if err != nil {
		span.RecordError(err)
		log.Errorf("[DeleteVaultCategory] extract id: %v", err)
		localization.SendBadRequestResponse(w, err.Error())
		return
	}
	span.SetAttributes(attribute.String("vault_category.id", id))
	_, err = h.service.DeleteVaultCategory(ctx, id)
	if err != nil {
		span.RecordError(err)
		log.Errorf("[DeleteVaultCategory] service: %v", err)
		localization.SendErrorByCodeResponse(w, err.Error())
		return
	}
	log.Infof("[VaultCatH][Delete] ok id: %s", id)
	localization.SendSuccessResponse(w, localization.SuccessVaultCategoryDeleteRequestSubmitted, nil)
}

// DisableVaultCategory
//
//	@Summary		Disable Vault Category
//	@Description	Disable a vault category by the provided ID
//	@Tags			Vault Category
//	@Accept			json
//	@Produce		json
//	@Param			id	path		string									true	"Vault category ID"
//	@Success		200	{object}	localization.StandardResponse{data=nil}	"Vault category disable request submitted successfully"
//	@Failure		400	{object}	localization.StandardResponse{data=nil}	"Bad request"
//	@Failure		500	{object}	localization.StandardResponse{data=nil}	"Internal server error"
//	@Security		BearerAuth
//	@Router			/vaultcategory/disable/{id} [patch]
func (h *handler) DisableVaultCategory(w http.ResponseWriter, r *http.Request) {
	ctx, span := common_utils.TraceLogger(r.Context(), "handler", "disableVaultCategory", "handler", "vaultCategory")
	defer span.End()
	log := common_utils.LoggerFromCtx(ctx, h.logger)

	md := &types.ContextMetadata{}
	ctx = context.WithValue(ctx, constants.ContextKeyMetadata, md)

	id, err := common_utils.ExtractID(w, r)
	if id == "" {
		appErr := middleware.NewValidationError("vault category id is required", map[string]interface{}{}).
			WithService("vault_category").
			WithOperation("DisableBankVault")
		span.RecordError(appErr)
		localization.SendErrorByCodeResponse(w, appErr.Error())
		return
	}
	if err != nil {
		span.RecordError(err)
		log.Errorf("[DisableVaultCategory] extract ID: %v", err)
		localization.SendErrorByCodeResponse(w, err.Error())
		return
	}
	span.SetAttributes(attribute.String("vault_category.id", id))
	err = h.service.DisableVaultCategory(ctx, id)
	if err != nil {
		span.RecordError(err)
		log.Errorf("[DisableVaultCategory] service: %v", err)
		localization.SendErrorByCodeResponse(w, err.Error())
		return
	}

	if md.IsMakerOnly {
		userCode, _ := ctx.Value(constants.ContextKey("user_code")).(string)
		log.Infof("[Create] request sent successfully for user_code: %s is_maker_only: %v", userCode, md.IsMakerOnly)
		localization.SendSuccessResponse(w, localization.SuccessVaultCategoryDisabledSuccessfully, nil)
		return
	}

	log.Infof("[VaultCatH][Disable] ok id: %s", id)
	localization.SendSuccessResponse(w, localization.SuccessVaultCategoryDisableRequestSubmitted, nil)
}

// EnableVaultCategory
//
//	@Summary		Enable Vault Category
//	@Description	Enable a vault category by the provided ID
//	@Tags			Vault Category
//	@Accept			json
//	@Produce		json
//	@Param			id	path		string									true	"Vault category ID"
//	@Success		200	{object}	localization.StandardResponse{data=nil}	"Vault category enable request submitted successfully"
//	@Failure		400	{object}	localization.StandardResponse{data=nil}	"Bad request"
//	@Failure		500	{object}	localization.StandardResponse{data=nil}	"Internal server error"
//	@Security		BearerAuth
//	@Router			/vaultcategory/enable/{id} [patch]
func (h *handler) EnableVaultCategory(w http.ResponseWriter, r *http.Request) {
	ctx, span := common_utils.TraceLogger(r.Context(), "handler", "enableVaultCategory", "handler", "vaultCategory")
	defer span.End()
	log := common_utils.LoggerFromCtx(ctx, h.logger)

	md := &types.ContextMetadata{}
	ctx = context.WithValue(ctx, constants.ContextKeyMetadata, md)

	id, err := common_utils.ExtractID(w, r)
	if id == "" {
		appErr := middleware.NewValidationError("vault category id is required", map[string]interface{}{}).
			WithService("vault_category").
			WithOperation("EnableBankVault")
		span.RecordError(appErr)
		localization.SendErrorByCodeResponse(w, appErr.Error())
		return
	}
	if err != nil {
		span.RecordError(err)
		log.Errorf("[EnableVaultCategory] extract ID: %v", err)
		localization.SendErrorByCodeResponse(w, err.Error())
		return
	}
	span.SetAttributes(attribute.String("vault_category.id", id))
	err = h.service.EnableVaultCategory(ctx, id)
	if err != nil {
		span.RecordError(err)
		log.Errorf("[EnableVaultCategory] service: %v", err)
		localization.SendErrorByCodeResponse(w, err.Error())
		return
	}

	if md.IsMakerOnly {
		userCode, _ := ctx.Value(constants.ContextKey("user_code")).(string)
		log.Infof("[Create] request sent successfully for user_code: %s is_maker_only: %v", userCode, md.IsMakerOnly)
		localization.SendSuccessResponse(w, localization.SuccessVaultCategoryEnabledSuccessfully, nil)
		return
	}

	log.Infof("[VaultCatH][Enable] ok id: %s", id)
	localization.SendSuccessResponse(w, localization.SuccessVaultCategoryEnableRequestSubmitted, nil)
}
