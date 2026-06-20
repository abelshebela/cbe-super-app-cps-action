package account_product_handler

import (
	"context"
	"net/http"

	"cbe-super-app-cps-action/internal/constants"
	ap_interface "cbe-super-app-cps-action/internal/constants/interfaces/account_product"
	"cbe-super-app-cps-action/internal/constants/localization"
	"cbe-super-app-cps-action/internal/constants/types"
	ap_core "cbe-super-app-cps-action/internal/handlers/rest/http/account_product/core"
	"cbe-super-app-cps-action/internal/service"
	local_util "cbe-super-app-cps-action/pkgs/utils"

	"github.com/go-chi/chi/v5"
	"gitlab.com/bersufekadgetachew/cbe-super-app-shared/shared/utils"
	"go.opentelemetry.io/otel/attribute"
)

type accountProductAdapter struct {
	svc    service.AccountProductService
	logger utils.Logger
}

func InitAccountProductAdapter(svc service.AccountProductService, logger utils.Logger) ap_interface.AccountProductHandler {
	return &accountProductAdapter{svc: svc, logger: logger}
}

func (h *accountProductAdapter) GetAll(w http.ResponseWriter, r *http.Request) {
	ctx, span := local_util.TraceLogger(r.Context(), "handler", "getAll", "handler", "account_product")
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
		log.Errorf("[APHandler][GetAll] service err: %v", err)
		localization.SendErrorByCodeResponse(w, err.Error())
		return
	}

	mapped := make([]interface{}, 0, len(result.Data))
	for i := range result.Data {
		mapped = append(mapped, ap_core.MapToResponse(&result.Data[i]))
	}
	res := types.PaginatedResponse[[]interface{}]{Data: mapped, Meta: result.Meta}

	span.SetAttributes(attribute.Int("ap.count", len(result.Data)))
	log.Infof("[APHandler][GetAll] count=%d", len(result.Data))
	localization.SendSuccessResponse(w, localization.SuccessAPsRetrieved, res)
}

func (h *accountProductAdapter) GetByID(w http.ResponseWriter, r *http.Request) {
	ctx, span := local_util.TraceLogger(r.Context(), "handler", "getByID", "handler", "account_product")
	defer span.End()
	log := local_util.LoggerFromCtx(ctx, h.logger)

	id := chi.URLParam(r, "id")
	if id == "" {
		localization.SendErrorResponse(w, localization.ErrorRequiredFieldMissing, nil, nil)
		return
	}
	span.SetAttributes(attribute.String("ap.id", id))

	result, err := h.svc.GetByID(ctx, id)
	if err != nil {
		span.RecordError(err)
		log.Errorf("[APHandler][GetByID] service err: %v", err)
		localization.SendErrorByCodeResponse(w, err.Error())
		return
	}

	w = localization.ApplyActionCodeHeaderFromWriter(w, ctx)

	log.Infof("[APHandler][GetByID] found id=%s", id)
	localization.SendSuccessResponse(w, localization.SuccessAPRetrieved, ap_core.MapToResponse(result))
}

func (h *accountProductAdapter) Create(w http.ResponseWriter, r *http.Request) {
	ctx, span := local_util.TraceLogger(r.Context(), "handler", "create", "handler", "account_product")
	defer span.End()
	log := local_util.LoggerFromCtx(ctx, h.logger)

	md := &types.ContextMetadata{}
	ctx = context.WithValue(ctx, constants.ContextKeyMetadata, md)
	localization.UpdateWriterContext(w, ctx)

	req, err := ap_core.ParseCreateRequest(r, h.logger)
	if err != nil {
		log.Errorf("[APHandler][Create] parse err: %v", err)
		localization.SendErrorResponse(w, localization.ErrorInvalidRequestBody, nil, nil)
		return
	}

	if rc := ap_core.ValidateCreate(&req); rc.Code != "" {
		span.SetAttributes(attribute.String("invalid input", rc.Code))
		log.Errorf("[APHandler][Create] validation: %s", rc.Code)
		localization.SendErrorResponse(w, localization.ErrorInvalidRequestBody, nil, nil)
		return
	}

	span.SetAttributes(attribute.String("ap.code", req.CBSProductCode))

	result, err := h.svc.Create(ctx, req)
	if err != nil {
		w = local_util.HandlePendingResponseError(ctx, w, err)
		span.RecordError(err)
		log.Errorf("[APHandler][Create] service err: %v", err)
		localization.SendErrorByCodeResponse(w, err.Error())
		return
	}

	if md.IsMakerOnly {
		log.Infof("[APHandler][Create] created code=%s", req.CBSProductCode)
		w = localization.ApplyActionCodeHeaderFromWriter(w, ctx)
		localization.SendSuccessResponse(w, localization.SuccessAPCreated, ap_core.MapToResponse(result))
		return
	}
	log.Infof("[APHandler][Create] request sent code=%s", req.CBSProductCode)
	w = localization.ApplyActionCodeHeaderFromWriter(w, ctx)
	localization.SendSuccessResponse(w, localization.SuccessAPCreateRequestSent, ap_core.MapToResponse(result))
}

func (h *accountProductAdapter) Update(w http.ResponseWriter, r *http.Request) {
	ctx, span := local_util.TraceLogger(r.Context(), "handler", "update", "handler", "account_product")
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
	span.SetAttributes(attribute.String("ap.id", id))

	req, err := ap_core.ParseUpdateRequest(r, h.logger)
	if err != nil {
		log.Errorf("[APHandler][Update] parse err: %v", err)
		localization.SendErrorResponse(w, localization.ErrorInvalidRequestBody, nil, nil)
		return
	}

	if rc := ap_core.ValidateUpdate(&req); rc.Code != "" {
		span.SetAttributes(attribute.String("invalid input", rc.Code))
		log.Errorf("[APHandler][Update] validation: %s", rc.Code)
		localization.SendErrorResponse(w, rc, nil, nil)
		return
	}

	result, err := h.svc.Update(ctx, id, req)
	if err != nil {
		w = local_util.HandlePendingResponseError(ctx, w, err)
		span.RecordError(err)
		log.Errorf("[APHandler][Update] service err: %v", err)
		localization.SendErrorByCodeResponse(w, err.Error())
		return
	}

	if md.IsMakerOnly {
		log.Infof("[APHandler][Update] updated id=%s", id)
		w = localization.ApplyActionCodeHeaderFromWriter(w, ctx)
		localization.SendSuccessResponse(w, localization.SuccessAPUpdated, ap_core.MapToResponse(result))
		return
	}
	log.Infof("[APHandler][Update] request sent id=%s", id)
	w = localization.ApplyActionCodeHeaderFromWriter(w, ctx)
	localization.SendSuccessResponse(w, localization.SuccessAPUpdateRequestSent, ap_core.MapToResponse(result))
}

func (h *accountProductAdapter) Delete(w http.ResponseWriter, r *http.Request) {
	ctx, span := local_util.TraceLogger(r.Context(), "handler", "delete", "handler", "account_product")
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
	span.SetAttributes(attribute.String("ap.id", id))

	if err := h.svc.Delete(ctx, id); err != nil {
		w = local_util.HandlePendingResponseError(ctx, w, err)
		span.RecordError(err)
		log.Errorf("[APHandler][Delete] service err: %v", err)
		localization.SendErrorByCodeResponse(w, err.Error())
		return
	}

	if md.IsMakerOnly {
		w = localization.ApplyActionCodeHeaderFromWriter(w, ctx)
		localization.SendSuccessResponse(w, localization.SuccessBankCreatedSuccessfully, nil)
		return
	}
	log.Infof("[APHandler][Delete] request sent id=%s", id)
	w = localization.ApplyActionCodeHeaderFromWriter(w, ctx)
	localization.SendSuccessResponse(w, localization.SuccessAPDeleteRequestSent, nil)
}

func (h *accountProductAdapter) Enable(w http.ResponseWriter, r *http.Request) {
	ctx, span := local_util.TraceLogger(r.Context(), "handler", "enable", "handler", "account_product")
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
	span.SetAttributes(attribute.String("ap.id", id))

	if err := h.svc.EnableOrDisable(ctx, id, true); err != nil {
		w = local_util.HandlePendingResponseError(ctx, w, err)
		span.RecordError(err)
		log.Errorf("[APHandler][Enable] service err: %v", err)
		localization.SendErrorByCodeResponse(w, err.Error())
		return
	}

	if md.IsMakerOnly {
		log.Infof("[APHandler][Enable] enabled id=%s", id)
		w = localization.ApplyActionCodeHeaderFromWriter(w, ctx)
		localization.SendSuccessResponse(w, localization.SuccessAPEnabled, nil)
		return
	}
	log.Infof("[APHandler][Enable] request sent id=%s", id)
	w = localization.ApplyActionCodeHeaderFromWriter(w, ctx)
	localization.SendSuccessResponse(w, localization.SuccessAPEnableRequestSent, nil)
}

func (h *accountProductAdapter) Disable(w http.ResponseWriter, r *http.Request) {
	ctx, span := local_util.TraceLogger(r.Context(), "handler", "disable", "handler", "account_product")
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
	span.SetAttributes(attribute.String("ap.id", id))

	if err := h.svc.EnableOrDisable(ctx, id, false); err != nil {
		w = local_util.HandlePendingResponseError(ctx, w, err)
		span.RecordError(err)
		log.Errorf("[APHandler][Disable] service err: %v", err)
		localization.SendErrorByCodeResponse(w, err.Error())
		return
	}

	if md.IsMakerOnly {
		log.Infof("[APHandler][Disable] disabled id=%s", id)
		w = localization.ApplyActionCodeHeaderFromWriter(w, ctx)
		localization.SendSuccessResponse(w, localization.SuccessAPDisabled, nil)
		return
	}
	log.Infof("[APHandler][Disable] request sent id=%s", id)
	w = localization.ApplyActionCodeHeaderFromWriter(w, ctx)
	localization.SendSuccessResponse(w, localization.SuccessAPDisableRequestSent, nil)
}
