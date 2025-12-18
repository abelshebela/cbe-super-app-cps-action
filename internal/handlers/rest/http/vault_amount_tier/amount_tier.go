package vaultamounttier

import (
	amount_tier "cbe-super-app-cps-action/internal/constants/dto/vault_amount_tier"
	"cbe-super-app-cps-action/internal/constants/localization"
	"cbe-super-app-cps-action/internal/service"
	common_utils "cbe-super-app-cps-action/pkgs/utils"
	"encoding/json"
	"net/http"

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

func (h *VaultAmountTierHandler) CreateAmountTier(w http.ResponseWriter, r *http.Request) {
	ctx, span := common_utils.TraceLogger(r.Context(), "handler", "CreateAmountTier", "handler", "CreateAmountTier")
	defer span.End()

	var req amount_tier.VaultAmountTierRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		span.RecordError(err)
		h.logger.Errorf("[CreateAmountTier] decode: %v", err)
		localization.SendBadRequestResponse(w, localization.MsgInvalidInput)
		return
	}

	if err := req.Validate(true); err != nil {
		span.RecordError(err)
		h.logger.Errorf("[CreateAmountTier] validation: %v", err)
		localization.SendBadRequestResponse(w, err.Error())
		return
	}

	_, err := h.service.CreateAmountTier(ctx, &req)
	if err != nil {
		span.RecordError(err)
		h.logger.Errorf("[CreateAmountTier] service: %v", err)
		localization.SendErrorByCodeResponse(w, err.Error())
		return
	}

	h.logger.Infof("[CreateAmountTier] vault amount tier created request submitted successfully")
	localization.SendSuccessResponse(w, localization.SuccessVaultAmountTierCreationRequestSubmitted, nil)
}

func (h *VaultAmountTierHandler) FindAllAmountTiers(w http.ResponseWriter, r *http.Request) {
	ctx, span := common_utils.TraceLogger(r.Context(), "handler", "FindAllAmountTiers", "handler", "FindAllAmountTiers")
	defer span.End()

	filterParams := common_utils.ExtractFilterParams(r)

	amountTiers, err := h.service.FindAllAmountTiers(ctx, filterParams)
	if err != nil {
		span.RecordError(err)
		h.logger.Errorf("[FindAllAmountTiers] service: %v", err)
		localization.SendErrorByCodeResponse(w, err.Error())
		return
	}

	h.logger.Infof("[FindAllAmountTiers] vault amount tiers fetched successfully")
	localization.SendSuccessResponse(w, localization.SuccessVaultAmountTierFetchedSuccessfully, amountTiers)
}

func (h *VaultAmountTierHandler) GetAmountTier(w http.ResponseWriter, r *http.Request) {
	ctx, span := common_utils.TraceLogger(r.Context(), "handler", "GetAmountTier", "handler", "GetAmountTier")
	defer span.End()

	id, err := common_utils.ExtractID(w, r)
	if err != nil {
		span.RecordError(err)
		h.logger.Errorf("[GetAmountTier] extractID: %v", err)
		return
	}

	amountTier, err := h.service.GetAmountTier(ctx, id)
	if err != nil {
		span.RecordError(err)
		h.logger.Errorf("[GetAmountTier] service: %v", err)
		localization.SendErrorByCodeResponse(w, err.Error())
		return
	}

	h.logger.Infof("[GetAmountTier] vault amount tier fetched successfully")
	localization.SendSuccessResponse(w, localization.SuccessVaultAmountTierFetchedSuccessfully, amountTier)
}

func (h *VaultAmountTierHandler) UpdateAmountTier(w http.ResponseWriter, r *http.Request) {
	ctx, span := common_utils.TraceLogger(r.Context(), "handler", "UpdateAmountTier", "handler", "UpdateAmountTier")
	defer span.End()

	var req amount_tier.UpdateVaultAmountTierRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		span.RecordError(err)
		h.logger.Errorf("[UpdateAmountTier] decode: %v", err)
		localization.SendBadRequestResponse(w, localization.MsgInvalidInput)
		return
	}

	if err := req.Validate(); err != nil {
		span.RecordError(err)
		h.logger.Errorf("[UpdateAmountTier] validation: %v", err)
		localization.SendBadRequestResponse(w, err.Error())
		return
	}

	id, err := common_utils.ExtractID(w, r)
	if err != nil {
		span.RecordError(err)
		h.logger.Errorf("[GetAmountTier] extractID: %v", err)
		return
	}

	_, err = h.service.UpdateAmountTier(ctx, id, &req)
	if err != nil {
		span.RecordError(err)
		h.logger.Errorf("[UpdateAmountTier] service: %v", err)
		localization.SendErrorByCodeResponse(w, err.Error())
		return
	}

	h.logger.Infof("[UpdateAmountTier] vault amount tier update request submitted successfully")
	localization.SendSuccessResponse(w, localization.SuccessVaultAmountTierUpdateRequestSubmitted, nil)
}

func (h *VaultAmountTierHandler) DeleteAmountTier(w http.ResponseWriter, r *http.Request) {
	ctx, span := common_utils.TraceLogger(r.Context(), "handler", "DeleteAmountTier", "handler", "DeleteAmountTier")
	defer span.End()

	id, err := common_utils.ExtractID(w, r)
	if err != nil {
		span.RecordError(err)
		h.logger.Errorf("[GetAmountTier] extractID: %v", err)
		return
	}

	_, err = h.service.DeleteAmountTier(ctx, id)
	if err != nil {
		span.RecordError(err)
		h.logger.Errorf("[DeleteAmountTier] service: %v", err)
		localization.SendErrorByCodeResponse(w, err.Error())
		return
	}

	h.logger.Infof("[DeleteAmountTier] vault amount tier delete request submitted successfully")
	localization.SendSuccessResponse(w, localization.SuccessVaultAmountTierDeleteRequestSubmitted, nil)
}

func (h *VaultAmountTierHandler) DisableAmountTier(w http.ResponseWriter, r *http.Request) {
	ctx, span := common_utils.TraceLogger(r.Context(), "handler", "DisableAmountTier", "handler", "DisableAmountTier")
	defer span.End()

	id, err := common_utils.ExtractID(w, r)
	if err != nil {
		span.RecordError(err)
		h.logger.Errorf("[GetAmountTier] extractID: %v", err)
		return
	}

	_, err = h.service.EnableOrDisableAmountTier(ctx, id, false)
	if err != nil {
		span.RecordError(err)
		h.logger.Errorf("[DisableAmountTier] service: %v", err)
		localization.SendErrorByCodeResponse(w, err.Error())
		return
	}

	h.logger.Infof("[DisableAmountTier] vault amount tier disable request submitted successfully")
	localization.SendSuccessResponse(w, localization.SuccessVaultAmountTierDisableRequestSubmitted, nil)
}

func (h *VaultAmountTierHandler) EnableAmountTier(w http.ResponseWriter, r *http.Request) {
	ctx, span := common_utils.TraceLogger(r.Context(), "handler", "EnableAmountTier", "handler", "EnableAmountTier")
	defer span.End()

	id, err := common_utils.ExtractID(w, r)
	if err != nil {
		span.RecordError(err)
		h.logger.Errorf("[GetAmountTier] extractID: %v", err)
		return
	}

	_, err = h.service.EnableOrDisableAmountTier(ctx, id, true)
	if err != nil {
		span.RecordError(err)
		h.logger.Errorf("[EnableAmountTier] service: %v", err)
		localization.SendErrorByCodeResponse(w, err.Error())
		return
	}

	h.logger.Infof("[EnableAmountTier] vault amount tier enable request submitted successfully")
	localization.SendSuccessResponse(w, localization.SuccessVaultAmountTierEnableRequestSubmitted, nil)
}
