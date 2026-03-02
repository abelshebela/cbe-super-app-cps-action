package vaultamounttier

import (
	"cbe-super-app-cps-action/internal/constants"
	amount_tier "cbe-super-app-cps-action/internal/constants/dto/vault_amount_tier"
	"cbe-super-app-cps-action/internal/constants/localization"
	"cbe-super-app-cps-action/internal/service"
	common_utils "cbe-super-app-cps-action/pkgs/utils"
	"context"
	"encoding/json"
	"net/http"

	"cbe-super-app-cps-action/internal/constants/types"

	"gitlab.com/bersufekadgetachew/cbe-super-app-shared/shared/utils"
)

type VaultAmountTierHandler struct {
	service service.VaultAmountBasedTierService
	logger  utils.Logger
}

func NewVaultAmountTierHandler(service service.VaultAmountBasedTierService, logger utils.Logger) *VaultAmountTierHandler {
	return &VaultAmountTierHandler{
		service: service,
		logger:  logger,
	}
}

// CreateAmountTier godoc
//
//	@Summary		Create amount tier for group and personal vaults
//	@Description	Create amount tier for group and personal vaults with the provided information
//	@Tags			Vault Amount Tier
//	@Accept			json
//	@Produce		json
//	@Param			request	body		vaultamounttier.VaultAmountTierRequest	true	"Vault Amount Tier Request"
//	@Success		200		{object}	localization.StandardResponse{data=nil}	"Amount tier update request submitted successfully"
//	@Failure		400		{object}	localization.StandardResponse{data=nil}	"Bad request"
//	@Failure		404		{object}	localization.StandardResponse{data=nil}	"Amount tier not found"
//	@Failure		500		{object}	localization.StandardResponse{data=nil}	"Internal server error"
//	@Security		BearerAuth
//	@Router			/vault-amount-tier/create [post]
func (h *VaultAmountTierHandler) CreateAmountTier(w http.ResponseWriter, r *http.Request) {
	ctx, span := common_utils.TraceLogger(r.Context(), "handler", "CreateAmountTier", "handler", "CreateAmountTier")
	defer span.End()
	log := common_utils.LoggerFromCtx(ctx, h.logger)
	md := &types.ContextMetadata{}
	ctx = context.WithValue(ctx, constants.ContextKeyMetadata, md)

	var req amount_tier.VaultAmountTierRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		span.RecordError(err)
		log.Errorf("[CreateAmountTier] decode: %v", err)
		localization.SendBadRequestResponse(w, localization.MsgInvalidInput)
		return
	}

	if err := req.Validate(true); err != nil {
		span.RecordError(err)
		log.Errorf("[CreateAmountTier] validation: %v", err)
		localization.SendBadRequestResponse(w, err.Error())
		return
	}

	_, err := h.service.CreateAmountTier(ctx, &req)
	if err != nil {
		span.RecordError(err)
		log.Errorf("[CreateAmountTier] service: %v", err)
		localization.SendErrorByCodeResponse(w, err.Error())
		return
	}

	if md.IsMakerOnly {
		userCode, _ := ctx.Value(constants.ContextKey("user_code")).(string)
		log.Infof("[CreateAmountTier] request sent successfully for user_code: %s is_maker_only: %v", userCode, md.IsMakerOnly)
		localization.SendSuccessResponse(w, localization.CustomerSegmentationEnabled, nil)
		return
	}

	log.Infof("[CreateAmountTier] vault amount tier created request submitted successfully")
	localization.SendSuccessResponse(w, localization.SuccessVaultAmountTierCreationRequestSubmitted, nil)
}

// FindAllAmountTiers godoc
//
//	@Summary		List amount tiers
//	@Description	Retrieve amount tiers with pagination and optional search
//	@Tags			Vault Amount Tier
//	@Accept			json
//	@Produce		json
//	@Param			page		query		int		false	"Page number"		default(1)
//	@Param			per_page	query		int		false	"Items per page"	default(10)
//	@Param			search		query		string	false	"Search term"
//	@Success		200			{object}	localization.StandardResponse{data=object}	"Amount tiers retrieved successfully"
//	@Failure		400			{object}	localization.StandardResponse{data=nil}		"Bad request"
//	@Failure		500			{object}	localization.StandardResponse{data=nil}		"Internal server error"
//	@Security		BearerAuth
//	@Router			/vault-amount-tier/find-all [get]
func (h *VaultAmountTierHandler) FindAllAmountTiers(w http.ResponseWriter, r *http.Request) {
	ctx, span := common_utils.TraceLogger(r.Context(), "handler", "FindAllAmountTiers", "handler", "FindAllAmountTiers")
	defer span.End()
	log := common_utils.LoggerFromCtx(ctx, h.logger)

	filterParams := common_utils.ExtractFilterParams(r)

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

	amountTiers, err := h.service.FindAllAmountTiers(ctx, filterParams)
	if err != nil {
		span.RecordError(err)
		log.Errorf("[FindAllAmountTiers] service: %v", err)
		localization.SendErrorByCodeResponse(w, err.Error())
		return
	}

	log.Infof("[FindAllAmountTiers] vault amount tiers fetched successfully")
	localization.SendSuccessResponse(w, localization.SuccessVaultAmountTierFetchedSuccessfully, amountTiers)
}

// GetAmountTier godoc
//
//	@Summary		Get amount tier by ID
//	@Description	Retrieve a amount tier's details by ID
//	@Tags			Vault Amount Tier
//	@Accept			json
//	@Produce		json
//	@Param			id	path		string	true	"Amount Tier ID"
//	@Success		200	{object}	localization.StandardResponse{data=object}
//	@Failure		400	{object}	localization.StandardResponse{data=nil}
//	@Failure		404	{object}	localization.StandardResponse{data=nil}
//	@Failure		500	{object}	localization.StandardResponse{data=nil}
//	@Security		BearerAuth
//	@Router			/vault-amount-tier/{id} [get]
func (h *VaultAmountTierHandler) GetAmountTier(w http.ResponseWriter, r *http.Request) {
	ctx, span := common_utils.TraceLogger(r.Context(), "handler", "GetAmountTier", "handler", "GetAmountTier")
	defer span.End()
	log := common_utils.LoggerFromCtx(ctx, h.logger)

	id, err := common_utils.ExtractID(w, r)
	if err != nil {
		span.RecordError(err)
		log.Errorf("[GetAmountTier] extractID: %v", err)
		return
	}

	amountTier, err := h.service.GetAmountTier(ctx, id)
	if err != nil {
		span.RecordError(err)
		log.Errorf("[GetAmountTier] service: %v", err)
		localization.SendErrorByCodeResponse(w, err.Error())
		return
	}

	log.Infof("[GetAmountTier] vault amount tier fetched successfully")
	localization.SendSuccessResponse(w, localization.SuccessVaultAmountTierFetchedSuccessfully, amountTier)
}

// UpdateAmountTier godoc
//
//	@Summary		Update amount tier
//	@Description	Update amount tier with the provided information
//	@Tags			Vault Amount Tier
//	@Accept			json
//	@Produce		json
//	@Param			id		path		string										true	"Amount Tier ID"
//	@Param			request	body		vaultamounttier.UpdateVaultAmountTierRequest	true	"Update Amount Tier Request"
//	@Success		200		{object}	localization.StandardResponse{data=nil}	"Amount tier update request submitted successfully"
//	@Failure		400		{object}	localization.StandardResponse{data=nil}	"Bad request"
//	@Failure		404		{object}	localization.StandardResponse{data=nil}	"Amount tier not found"
//	@Failure		500		{object}	localization.StandardResponse{data=nil}	"Internal server error"
//	@Security		BearerAuth
//	@Router			/vault-amount-tier/{id}/update [patch]
func (h *VaultAmountTierHandler) UpdateAmountTier(w http.ResponseWriter, r *http.Request) {
	ctx, span := common_utils.TraceLogger(r.Context(), "handler", "UpdateAmountTier", "handler", "UpdateAmountTier")
	defer span.End()
	log := common_utils.LoggerFromCtx(ctx, h.logger)
	md := &types.ContextMetadata{}
	ctx = context.WithValue(ctx, constants.ContextKeyMetadata, md)

	var req amount_tier.UpdateVaultAmountTierRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		span.RecordError(err)
		log.Errorf("[UpdateAmountTier] decode: %v", err)
		localization.SendBadRequestResponse(w, localization.MsgInvalidInput)
		return
	}

	if err := req.Validate(); err != nil {
		span.RecordError(err)
		log.Errorf("[UpdateAmountTier] validation: %v", err)
		localization.SendBadRequestResponse(w, err.Error())
		return
	}

	id, err := common_utils.ExtractID(w, r)
	if err != nil {
		span.RecordError(err)
		log.Errorf("[GetAmountTier] extractID: %v", err)
		return
	}

	_, err = h.service.UpdateAmountTier(ctx, id, &req)
	if err != nil {
		span.RecordError(err)
		log.Errorf("[UpdateAmountTier] service: %v", err)
		localization.SendErrorByCodeResponse(w, err.Error())
		return
	}

	if md.IsMakerOnly {
		userCode, _ := ctx.Value(constants.ContextKey("user_code")).(string)
		log.Infof("[UpdateAmountTier] request sent successfully for user_code: %s is_maker_only: %v", userCode, md.IsMakerOnly)
		localization.SendSuccessResponse(w, localization.CustomerSegmentationEnabled, nil)
		return
	}

	log.Infof("[UpdateAmountTier] vault amount tier update request submitted successfully")
	localization.SendSuccessResponse(w, localization.SuccessVaultAmountTierUpdateRequestSubmitted, nil)
}

// DeleteAmountTier godoc
//
//	@Summary		Delete amount tier
//	@Description	Delete amount tier by ID
//	@Tags			Vault Amount Tier
//	@Accept			json
//	@Produce		json
//	@Param			id	path		string	true	"Amount Tier ID"
//	@Success		200	{object}	localization.StandardResponse{data=nil}
//	@Failure		400	{object}	localization.StandardResponse{data=nil}
//	@Failure		500	{object}	localization.StandardResponse{data=nil}
//	@Security		BearerAuth
//	@Router			/vault-amount-tier/{id}/delete [delete]
func (h *VaultAmountTierHandler) DeleteAmountTier(w http.ResponseWriter, r *http.Request) {
	ctx, span := common_utils.TraceLogger(r.Context(), "handler", "DeleteAmountTier", "handler", "DeleteAmountTier")
	defer span.End()
	log := common_utils.LoggerFromCtx(ctx, h.logger)
	md := &types.ContextMetadata{}
	ctx = context.WithValue(ctx, constants.ContextKeyMetadata, md)

	id, err := common_utils.ExtractID(w, r)
	if err != nil {
		span.RecordError(err)
		log.Errorf("[GetAmountTier] extractID: %v", err)
		return
	}

	_, err = h.service.DeleteAmountTier(ctx, id)
	if err != nil {
		span.RecordError(err)
		log.Errorf("[DeleteAmountTier] service: %v", err)
		localization.SendErrorByCodeResponse(w, err.Error())
		return
	}

	if md.IsMakerOnly {
		userCode, _ := ctx.Value(constants.ContextKey("user_code")).(string)
		log.Infof("[DeleteAmountTier] request sent successfully for user_code: %s is_maker_only: %v", userCode, md.IsMakerOnly)
		localization.SendSuccessResponse(w, localization.CustomerSegmentationEnabled, nil)
		return
	}

	log.Infof("[DeleteAmountTier] vault amount tier delete request submitted successfully")
	localization.SendSuccessResponse(w, localization.SuccessVaultAmountTierDeleteRequestSubmitted, nil)
}

// DisableAmountTier godoc
//
//	@Summary		Disable amount tier
//	@Description	Disable amount tier by ID
//	@Tags			Vault Amount Tier
//	@Accept			json
//	@Produce		json
//	@Param			id	path		string	true	"Amount Tier ID"
//	@Success		200	{object}	localization.StandardResponse{data=nil}
//	@Failure		400	{object}	localization.StandardResponse{data=nil}
//	@Failure		500	{object}	localization.StandardResponse{data=nil}
//	@Security		BearerAuth
//	@Router			/vault-amount-tier/{id}/disable [patch]
func (h *VaultAmountTierHandler) DisableAmountTier(w http.ResponseWriter, r *http.Request) {
	ctx, span := common_utils.TraceLogger(r.Context(), "handler", "DisableAmountTier", "handler", "DisableAmountTier")
	defer span.End()
	log := common_utils.LoggerFromCtx(ctx, h.logger)
	md := &types.ContextMetadata{}
	ctx = context.WithValue(ctx, constants.ContextKeyMetadata, md)

	id, err := common_utils.ExtractID(w, r)
	if err != nil {
		span.RecordError(err)
		log.Errorf("[GetAmountTier] extractID: %v", err)
		return
	}

	_, err = h.service.EnableOrDisableAmountTier(ctx, id, false)
	if err != nil {
		span.RecordError(err)
		log.Errorf("[DisableAmountTier] service: %v", err)
		localization.SendErrorByCodeResponse(w, err.Error())
		return
	}

	if md.IsMakerOnly {
		userCode, _ := ctx.Value(constants.ContextKey("user_code")).(string)
		log.Infof("[DisableAmountTier] request sent successfully for user_code: %s is_maker_only: %v", userCode, md.IsMakerOnly)
		localization.SendSuccessResponse(w, localization.CustomerSegmentationEnabled, nil)
		return
	}

	log.Infof("[DisableAmountTier] vault amount tier disable request submitted successfully")
	localization.SendSuccessResponse(w, localization.SuccessVaultAmountTierDisableRequestSubmitted, nil)
}

// EnableAmountTier godoc
//
//	@Summary		Enable amount tier
//	@Description	Enable amount tier by ID
//	@Tags			Vault Amount Tier
//	@Accept			json
//	@Produce		json
//	@Param			id	path		string	true	"Amount Tier ID"
//	@Success		200	{object}	localization.StandardResponse{data=nil}
//	@Failure		400	{object}	localization.StandardResponse{data=nil}
//	@Failure		500	{object}	localization.StandardResponse{data=nil}
//	@Security		BearerAuth
//	@Router			/vault-amount-tier/{id}/enable [patch]
func (h *VaultAmountTierHandler) EnableAmountTier(w http.ResponseWriter, r *http.Request) {
	ctx, span := common_utils.TraceLogger(r.Context(), "handler", "EnableAmountTier", "handler", "EnableAmountTier")
	defer span.End()
	log := common_utils.LoggerFromCtx(ctx, h.logger)
	md := &types.ContextMetadata{}
	ctx = context.WithValue(ctx, constants.ContextKeyMetadata, md)

	id, err := common_utils.ExtractID(w, r)
	if err != nil {
		span.RecordError(err)
		log.Errorf("[GetAmountTier] extractID: %v", err)
		return
	}

	_, err = h.service.EnableOrDisableAmountTier(ctx, id, true)
	if err != nil {
		span.RecordError(err)
		log.Errorf("[EnableAmountTier] service: %v", err)
		localization.SendErrorByCodeResponse(w, err.Error())
		return
	}

	if md.IsMakerOnly {
		userCode, _ := ctx.Value(constants.ContextKey("user_code")).(string)
		log.Infof("[EnableAmountTier] request sent successfully for user_code: %s is_maker_only: %v", userCode, md.IsMakerOnly)
		localization.SendSuccessResponse(w, localization.CustomerSegmentationEnabled, nil)
		return
	}

	log.Infof("[EnableAmountTier] vault amount tier enable request submitted successfully")
	localization.SendSuccessResponse(w, localization.SuccessVaultAmountTierEnableRequestSubmitted, nil)
}
