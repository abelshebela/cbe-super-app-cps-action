package term_and_condition_handler

import (
	"context"
	"net/http"

	"cbe-super-app-cps-action/internal/constants"
	tac_interface "cbe-super-app-cps-action/internal/constants/interfaces/term_and_condition"
	"cbe-super-app-cps-action/internal/constants/localization"
	"cbe-super-app-cps-action/internal/constants/types"
	tac_core "cbe-super-app-cps-action/internal/handlers/rest/http/term_and_condition/core"
	"cbe-super-app-cps-action/internal/service"
	local_util "cbe-super-app-cps-action/pkgs/utils"

	"github.com/go-chi/chi/v5"
	"gitlab.com/bersufekadgetachew/cbe-super-app-shared/shared/utils"
	"go.opentelemetry.io/otel/attribute"
)

type termAndConditionAdapter struct {
	svc    service.AccountOpeningTermsService
	logger utils.Logger
}

func InitTermAndConditionAdapter(svc service.AccountOpeningTermsService, logger utils.Logger) tac_interface.TermAndConditionHandler {
	return &termAndConditionAdapter{svc: svc, logger: logger}
}

func (h *termAndConditionAdapter) GetAll(w http.ResponseWriter, r *http.Request) {
	ctx, span := local_util.TraceLogger(r.Context(), "handler", "getAll", "handler", "term_and_condition")
	defer span.End()
	log := local_util.LoggerFromCtx(ctx, h.logger)

	filterParams := local_util.ExtractFilterParams(r)
	search := r.URL.Query().Get("search")
	if err := local_util.NoSpecialChars(search); err != nil {
		localization.SendErrorByCodeResponse(w, err.Error())
		return
	}

	result, err := h.svc.GetAll(ctx, filterParams)
	if err != nil {
		span.RecordError(err)
		log.Errorf("[TACHandler][GetAll] service err: %v", err)
		localization.SendErrorByCodeResponse(w, err.Error())
		return
	}

	mapped := make([]interface{}, 0, len(result.Data))
	for i := range result.Data {
		mapped = append(mapped, tac_core.MapToResponse(&result.Data[i]))
	}
	res := types.PaginatedResponse[[]interface{}]{Data: mapped, Meta: result.Meta}

	span.SetAttributes(attribute.Int("tac.count", len(result.Data)))
	log.Infof("[TACHandler][GetAll] count=%d", len(result.Data))
	localization.SendSuccessResponse(w, localization.SuccessTACsRetrieved, res)
}

func (h *termAndConditionAdapter) GetByID(w http.ResponseWriter, r *http.Request) {
	ctx, span := local_util.TraceLogger(r.Context(), "handler", "getByID", "handler", "term_and_condition")
	defer span.End()
	log := local_util.LoggerFromCtx(ctx, h.logger)

	id := chi.URLParam(r, "id")
	if id == "" {
		localization.SendErrorResponse(w, localization.ErrorRequiredFieldMissing, nil, nil)
		return
	}
	span.SetAttributes(attribute.String("tac.id", id))

	result, err := h.svc.GetByID(ctx, id)
	if err != nil {
		span.RecordError(err)
		log.Errorf("[TACHandler][GetByID] service err: %v", err)
		localization.SendErrorByCodeResponse(w, err.Error())
		return
	}

	log.Infof("[TACHandler][GetByID] found id=%s", id)
	localization.SendSuccessResponse(w, localization.SuccessTACRetrieved, tac_core.MapToResponse(result))
}

func (h *termAndConditionAdapter) Upload(w http.ResponseWriter, r *http.Request) {
	ctx, span := local_util.TraceLogger(r.Context(), "handler", "upload", "handler", "term_and_condition")
	defer span.End()
	log := local_util.LoggerFromCtx(ctx, h.logger)

	md := &types.ContextMetadata{}
	ctx = context.WithValue(ctx, constants.ContextKeyMetadata, md)
	localization.UpdateWriterContext(w, ctx)

	req, err := tac_core.ParseUploadRequest(r, h.logger)
	if err != nil {
		log.Errorf("[TACHandler][Upload] parse err: %v", err)
		localization.SendErrorResponse(w, localization.ErrorInvalidRequestBody, nil, nil)
		return
	}

	if rc := tac_core.ValidateUpload(&req); rc.Code != "" {
		span.SetAttributes(attribute.String("invalid input", rc.Code))
		log.Errorf("[TACHandler][Upload] validation: %s", rc.Code)
		localization.SendErrorResponse(w, rc, nil, nil)
		return
	}

	span.SetAttributes(
		attribute.String("tac.product_id", req.AccountProductID),
		attribute.String("tac.version", req.VersionLabel),
	)

	result, err := h.svc.Upload(ctx, req)
	if err != nil {
		w = local_util.HandlePendingResponseError(ctx, w, err)
		span.RecordError(err)
		log.Errorf("[TACHandler][Upload] service err: %v", err)
		localization.SendErrorByCodeResponse(w, err.Error())
		return
	}

	if md.IsMakerOnly {
		log.Infof("[TACHandler][Upload] uploaded product=%s version=%s", req.AccountProductID, req.VersionLabel)
		w = localization.ApplyActionCodeHeaderFromWriter(w, ctx)
		localization.SendSuccessResponse(w, localization.SuccessTACUploaded, tac_core.MapToResponse(result))
		return
	}
	log.Infof("[TACHandler][Upload] request sent product=%s version=%s", req.AccountProductID, req.VersionLabel)
	w = localization.ApplyActionCodeHeaderFromWriter(w, ctx)
	localization.SendSuccessResponse(w, localization.SuccessTACUploadRequestSent, tac_core.MapToResponse(result))
}

func (h *termAndConditionAdapter) Update(w http.ResponseWriter, r *http.Request) {
	ctx, span := local_util.TraceLogger(r.Context(), "handler", "update", "handler", "term_and_condition")
	defer span.End()
	log := local_util.LoggerFromCtx(ctx, h.logger)

	md := &types.ContextMetadata{}
	ctx = context.WithValue(ctx, constants.ContextKeyMetadata, md)
	localization.UpdateWriterContext(w, ctx)

	id := chi.URLParam(r, "id")
	if id == "" {
		localization.SendErrorResponse(w, localization.ErrorRequiredFieldMissing, nil, nil)
		return
	}
	span.SetAttributes(attribute.String("tac.id", id))

	req, err := tac_core.ParseUpdateRequest(r, h.logger)
	if err != nil {
		log.Errorf("[TACHandler][Update] parse err: %v", err)
		localization.SendErrorResponse(w, localization.ErrorInvalidRequestBody, nil, nil)
		return
	}

	if rc := tac_core.ValidateUpdate(&req); rc.Code != "" {
		span.SetAttributes(attribute.String("invalid input", rc.Code))
		log.Errorf("[TACHandler][Update] validation: %s", rc.Code)
		localization.SendErrorResponse(w, rc, nil, nil)
		return
	}

	result, err := h.svc.Update(ctx, id, req)
	if err != nil {
		w = local_util.HandlePendingResponseError(ctx, w, err)
		span.RecordError(err)
		log.Errorf("[TACHandler][Update] service err: %v", err)
		localization.SendErrorByCodeResponse(w, err.Error())
		return
	}

	if md.IsMakerOnly {
		log.Infof("[TACHandler][Update] request sent id=%s", id)
		w = localization.ApplyActionCodeHeaderFromWriter(w, ctx)
		localization.SendSuccessResponse(w, localization.SuccessTACUpdateRequestSent, tac_core.MapToResponse(result))
		return
	}
	log.Infof("[TACHandler][Update] updated id=%s", id)
	w = localization.ApplyActionCodeHeaderFromWriter(w, ctx)
	localization.SendSuccessResponse(w, localization.SuccessTACUpdated, tac_core.MapToResponse(result))
}

func (h *termAndConditionAdapter) Delete(w http.ResponseWriter, r *http.Request) {
	ctx, span := local_util.TraceLogger(r.Context(), "handler", "delete", "handler", "term_and_condition")
	defer span.End()
	log := local_util.LoggerFromCtx(ctx, h.logger)

	md := &types.ContextMetadata{}
	ctx = context.WithValue(ctx, constants.ContextKeyMetadata, md)
	localization.UpdateWriterContext(w, ctx)

	id := chi.URLParam(r, "id")
	if id == "" {
		localization.SendErrorResponse(w, localization.ErrorRequiredFieldMissing, nil, nil)
		return
	}
	span.SetAttributes(attribute.String("tac.id", id))

	if err := h.svc.Delete(ctx, id); err != nil {

		w = local_util.HandlePendingResponseError(ctx, w, err)
		span.RecordError(err)
		log.Errorf("[TACHandler][Delete] service err: %v", err)
		localization.SendErrorByCodeResponse(w, err.Error())
		return
	}

	log.Infof("[TACHandler][Delete] request sent id=%s", id)
	w = localization.ApplyActionCodeHeaderFromWriter(w, ctx)
	localization.SendSuccessResponse(w, localization.SuccessTACDeleteRequestSent, nil)
}
