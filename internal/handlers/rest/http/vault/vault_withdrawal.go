package vault

import (
	"cbe-super-app-cps-action/internal/constants"
	vault_category_dto "cbe-super-app-cps-action/internal/constants/dto/vault"
	localization "cbe-super-app-cps-action/internal/constants/localization"
	"cbe-super-app-cps-action/internal/constants/types"
	common_utils "cbe-super-app-cps-action/pkgs/utils"
	"context"
	"encoding/json"
	"net/http"

	"go.opentelemetry.io/otel/attribute"
)

func (h *handler) CreateWithdrawalRequest(w http.ResponseWriter, r *http.Request) {
	ctx, span := common_utils.TraceLogger(r.Context(), "handler", "CreateWithdrawalRequest", "handler", "CreateWithdrawalRequest")
	defer span.End()

	md := &types.ContextMetadata{}
	ctx = context.WithValue(ctx, constants.ContextKeyMetadata, md)

	var req vault_category_dto.CreateWithdrawalRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		span.RecordError(err)
		h.logger.Errorf("[CreateWithdrawalRequest] decode: %v", err)
		localization.SendBadRequestResponse(w, localization.MsgInvalidInput)
		return

	}

	if err := req.Validate(); err != nil {
		span.RecordError(err)
		h.logger.Errorf("[CreateWithdrawalRequest] validation: %v", err)
		localization.SendBadRequestResponse(w, err.Error())
		return
	}

	err := h.service.CreateWithdrawalRequest(ctx, &req)
	if err != nil {
		span.RecordError(err)
		h.logger.Errorf("[CreateWithdrawalRequest] service: %v", err)
		localization.SendErrorByCodeResponse(w, err.Error())
		return
	}

	if md.IsMakerOnly {
		userCode, _ := ctx.Value(constants.ContextKey("user_code")).(string)
		h.logger.Infof("[CreateWithdrawalRequest] withdrawal request sent successfully for user_code: %s is_maker_only: %v", userCode, md.IsMakerOnly)
		localization.SendSuccessResponse(w, localization.SuccessVaultWithdrawalRequest, nil)
		return
	}

	h.logger.Infof("[CreateWithdrawalRequest] withdrawal request submitted successfully")
	localization.SendSuccessResponse(w, localization.SuccessVaultWithdrawalRequestSubmitted, nil)
}

func (h *handler) UpdateWithDrawalRequest(w http.ResponseWriter, r *http.Request) {
	ctx, span := common_utils.TraceLogger(r.Context(), "handler", "UpdateWithDrawalRequest", "handler", "UpdateWithDrawalRequest")
	defer span.End()

	md := &types.ContextMetadata{}
	ctx = context.WithValue(ctx, constants.ContextKeyMetadata, md)

	var req vault_category_dto.UpdateWithdrawalStatusRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		span.RecordError(err)
		h.logger.Errorf("[UpdateWithDrawalRequest] decode: %v", err)
		localization.SendBadRequestResponse(w, localization.MsgInvalidInput)
		return

	}

	if err := req.Validate(); err != nil {
		span.RecordError(err)
		h.logger.Errorf("[UpdateWithDrawalRequest] validation: %v", err)
		localization.SendBadRequestResponse(w, err.Error())
		return
	}

	id, err := common_utils.ExtractID(w, r)
	if err != nil || id == "" {
		span.RecordError(err)
		h.logger.Errorf("[UpdateWithDrawalRequest] extract id: %v", err)
		localization.SendBadRequestResponse(w, err.Error())
		return
	}

	err = h.service.UpdateWithdrawalRequest(ctx, id, &req)
	if err != nil {
		span.RecordError(err)
		h.logger.Errorf("[UpdateWithDrawalRequest] service: %v", err)
		localization.SendErrorByCodeResponse(w, err.Error())
		return
	}

	if md.IsMakerOnly {
		userCode, _ := ctx.Value(constants.ContextKey("user_code")).(string)
		h.logger.Infof("[UpdateWithDrawalRequest] withdrawal request sent successfully for user_code: %s is_maker_only: %v", userCode, md.IsMakerOnly)
		localization.SendSuccessResponse(w, localization.SuccessVaultWithdrawalUpdateRequest, nil)
		return
	}

	h.logger.Infof("[UpdateWithDrawalRequest] withdrawal request submitted successfully")
	localization.SendSuccessResponse(w, localization.SuccessVaultWithdrawalUpdateSubmitted, nil)
}

func (h *handler) GetAllWithdrawalRequests(w http.ResponseWriter, r *http.Request) {
	ctx, span := common_utils.TraceLogger(r.Context(), "handler", "GetAllWithdrawalRequests", "handler", "GetAllWithdrawalRequests")
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

	result, err := h.service.GetAllWithdrawalRequests(ctx, params)
	if err != nil {
		span.RecordError(err)
		h.logger.Errorf("[GetAllWithdrawalRequests] service: %v", err)
		localization.SendErrorByCodeResponse(w, err.Error())
		return
	}

	localization.SendSuccessResponse(w, localization.SuccessVaultWithdrawalRequestsFetchedSuccessfully, result)
}

func (h *handler) GetWithdrawalRequestById(w http.ResponseWriter, r *http.Request) {
	ctx, span := common_utils.TraceLogger(r.Context(), "handler", "GetWithdrawalRequestById", "handler", "GetWithdrawalRequestById")
	defer span.End()
	id, err := common_utils.ExtractID(w, r)
	if id == "" {
		localization.SendErrorByCodeResponse(w, err.Error())
		return
	}
	if err != nil {
		span.RecordError(err)
		h.logger.Errorf("[GetWithdrawalRequestById] extract id: %v", err)
		localization.SendBadRequestResponse(w, err.Error())
		return
	}
	span.SetAttributes(attribute.String("vault_transaction.id", id))
	result, err := h.service.GetWithdrawalRequest(ctx, id)
	if err != nil {
		span.RecordError(err)
		h.logger.Errorf("[GetWithdrawalRequestById] service: %v", err)
		localization.SendErrorByCodeResponse(w, err.Error())
		return
	}
	h.logger.Infof("Vault withdrawal request retrieved with ID: %s", id)
	localization.SendSuccessResponse(w, localization.SuccessVaultWithdrawalRequestFetchedSuccessfully, result)
}
