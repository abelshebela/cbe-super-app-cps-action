package vault

import (
	"cbe-super-app-cps-action/internal/constants"
	localization "cbe-super-app-cps-action/internal/constants/localization"
	"cbe-super-app-cps-action/internal/constants/types"
	common_utils "cbe-super-app-cps-action/pkgs/utils"
	"context"
	"net/http"

	"go.opentelemetry.io/otel/attribute"
)

// func (h *handler) CreateWithdrawalRequest(w http.ResponseWriter, r *http.Request) {
// 	ctx, span := common_utils.TraceLogger(r.Context(), "handler", "CreateWithdrawalRequest", "handler", "CreateWithdrawalRequest")
// 	defer span.End()
// 	log := common_utils.LoggerFromCtx(ctx, h.logger)

// 	md := &types.ContextMetadata{}
// 	ctx = context.WithValue(ctx, constants.ContextKeyMetadata, md)
// 	localization.UpdateWriterContext(w, ctx)

// 	var req vault_category_dto.CreateWithdrawalRequest
// 	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
// 		span.RecordError(err)
// 		log.Errorf("[CreateWithdrawalRequest] decode: %v", err)
// 		localization.SendBadRequestResponse(w, localization.MsgInvalidInput)
// 		return

// 	}

// 	if err := req.Validate(); err != nil {
// 		span.RecordError(err)
// 		log.Errorf("[CreateWithdrawalRequest] validation: %v", err)
// 		localization.SendBadRequestResponse(w, err.Error())
// 		return
// 	}

// 	err := h.service.CreateWithdrawalRequest(ctx, &req)
// 	if err != nil {
// 		span.RecordError(err)
// 		log.Errorf("[CreateWithdrawalRequest] service: %v", err)
// 		localization.SendErrorByCodeResponse(w, err.Error())
// 		return
// 	}

// 	if md.IsMakerOnly {
// 		userCode, _ := ctx.Value(constants.ContextKey("user_code")).(string)
// 		log.Infof("[CreateWithdrawalRequest] withdrawal request sent successfully for user_code: %s is_maker_only: %v", userCode, md.IsMakerOnly)
// 		localization.SendSuccessResponse(w, localization.SuccessVaultWithdrawalRequest, nil)
// 		return
// 	}

// 	log.Infof("[CreateWithdrawalRequest] withdrawal request submitted successfully")
// 	localization.SendSuccessResponse(w, localization.SuccessVaultWithdrawalRequestSubmitted, nil)
// }

func (h *handler) UlockDeadlockRequest(w http.ResponseWriter, r *http.Request) {
	ctx, span := common_utils.TraceLogger(r.Context(), "handler", "UnlockDeadlockRequest", "handler", "UnlockDeadlockRequest")
	defer span.End()
	log := common_utils.LoggerFromCtx(ctx, h.logger)

	md := &types.ContextMetadata{}
	ctx = context.WithValue(ctx, constants.ContextKeyMetadata, md)
	localization.UpdateWriterContext(w, ctx)

	id, err := common_utils.ExtractID(w, r)
	if err != nil || id == "" {
		span.RecordError(err)
		log.Errorf("[UnlockDeadlockRequest] extract id: %v", err)
		localization.SendBadRequestResponse(w, err.Error())
		return
	}

	err = h.service.UnlockDeadlockRequest(ctx, id)
	if err != nil {
		span.RecordError(err)
		log.Errorf("[UnlockDeadlockRequest] service: %v", err)
		localization.SendErrorByCodeResponse(w, err.Error())
		return
	}

	if md.IsMakerOnly {
		userCode, _ := ctx.Value(constants.ContextKey("user_code")).(string)
		log.Infof("[UnlockDeadlockRequest] withdrawal request sent successfully for user_code: %s is_maker_only: %v", userCode, md.IsMakerOnly)
		w = localization.ApplyActionCodeHeaderFromWriter(w, ctx)
		localization.SendSuccessResponse(w, localization.SuccessVaultWithdrawalUpdateRequest, nil)
		return
	}

	log.Infof("[UnlockDeadlockRequest] withdrawal request submitted successfully")
	w = localization.ApplyActionCodeHeaderFromWriter(w, ctx)
	localization.SendSuccessResponse(w, localization.SuccessVaultDeadlockUnlockUpdateSubmitted, nil)
}

func (h *handler) GetAllDeadlockRequests(w http.ResponseWriter, r *http.Request) {
	ctx, span := common_utils.TraceLogger(r.Context(), "handler", "GetAllDeadlockRequests", "handler", "GetAllDeadlockRequests")
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

	result, err := h.service.GetAllDeadlockRequests(ctx, params)
	if err != nil {
		span.RecordError(err)
		log.Errorf("[GetAllDeadlockRequests] service: %v", err)
		localization.SendErrorByCodeResponse(w, err.Error())
		return
	}

	localization.SendSuccessResponse(w, localization.SuccessVaultDeadlockUnlockRequestsFetchedSuccessfully, result)
}

func (h *handler) GetDeadlockRequestById(w http.ResponseWriter, r *http.Request) {
	ctx, span := common_utils.TraceLogger(r.Context(), "handler", "GetDeadlockRequestById", "handler", "GetDeadlockRequestById")
	defer span.End()
	log := common_utils.LoggerFromCtx(ctx, h.logger)
	id, err := common_utils.ExtractID(w, r)
	if id == "" {
		localization.SendErrorByCodeResponse(w, err.Error())
		return
	}
	if err != nil {
		span.RecordError(err)
		log.Errorf("[GetDeadlockRequestById] extract id: %v", err)
		localization.SendBadRequestResponse(w, err.Error())
		return
	}
	span.SetAttributes(attribute.String("vault_transaction.id", id))
	result, err := h.service.GetDeadlockRequestById(ctx, id)
	if err != nil {
		span.RecordError(err)
		log.Errorf("[GetDeadlockRequestById] service: %v", err)
		localization.SendErrorByCodeResponse(w, err.Error())
		return
	}
	log.Infof("[VaultWithdH][GetByID] ok id: %s", id)
	localization.SendSuccessResponse(w, localization.SuccessVaultDeadlockUnlockRequestFetchedSuccessfully, result)
}
