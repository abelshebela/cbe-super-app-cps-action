package account_product_category_handler

import (
	"context"
	"net/http"

	"cbe-super-app-cps-action/internal/constants"
	apc_interface "cbe-super-app-cps-action/internal/constants/interfaces/account_product_category"
	"cbe-super-app-cps-action/internal/constants/localization"
	"cbe-super-app-cps-action/internal/constants/types"
	apc_core "cbe-super-app-cps-action/internal/handlers/rest/http/account_product_category/core"
	"cbe-super-app-cps-action/internal/service"
	local_util "cbe-super-app-cps-action/pkgs/utils"

	"github.com/go-chi/chi/v5"
	"gitlab.com/bersufekadgetachew/cbe-super-app-shared/shared/utils"
	"go.opentelemetry.io/otel/attribute"
)

type accountProductCategoryAdapter struct {
	svc    service.AccountProductCategoryService
	logger utils.Logger
}

func InitAccountProductCategoryAdapter(svc service.AccountProductCategoryService, logger utils.Logger) apc_interface.AccountProductCategoryHandler {
	return &accountProductCategoryAdapter{svc: svc, logger: logger}
}

func (h *accountProductCategoryAdapter) GetAll(w http.ResponseWriter, r *http.Request) {
	ctx, span := local_util.TraceLogger(r.Context(), "handler", "getAll", "handler", "account_product_category")
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
		log.Errorf("[APCHandler][GetAll] service err: %v", err)
		localization.SendErrorByCodeResponse(w, err.Error())
		return
	}

	mapped := make([]interface{}, 0, len(result.Data))
	for i := range result.Data {
		mapped = append(mapped, apc_core.MapToResponse(&result.Data[i]))
	}
	res := types.PaginatedResponse[[]interface{}]{Data: mapped, Meta: result.Meta}

	span.SetAttributes(attribute.Int("apc.count", len(result.Data)))
	log.Infof("[APCHandler][GetAll] count=%d", len(result.Data))
	localization.SendSuccessResponse(w, localization.SuccessAPCsRetrieved, res)
}

func (h *accountProductCategoryAdapter) GetByID(w http.ResponseWriter, r *http.Request) {
	ctx, span := local_util.TraceLogger(r.Context(), "handler", "getByID", "handler", "account_product_category")
	defer span.End()
	log := local_util.LoggerFromCtx(ctx, h.logger)

	id := chi.URLParam(r, "id")
	if id == "" {
		localization.SendErrorResponse(w, localization.ErrorRequiredFieldMissing, nil, nil)
		return
	}
	span.SetAttributes(attribute.String("apc.id", id))

	result, err := h.svc.GetByID(ctx, id)
	if err != nil {
		span.RecordError(err)
		log.Errorf("[APCHandler][GetByID] service err: %v", err)
		localization.SendErrorByCodeResponse(w, err.Error())
		return
	}

	log.Infof("[APCHandler][GetByID] found id=%s", id)
	localization.SendSuccessResponse(w, localization.SuccessAPCRetrieved, apc_core.MapToResponse(result))
}

func (h *accountProductCategoryAdapter) Create(w http.ResponseWriter, r *http.Request) {
	ctx, span := local_util.TraceLogger(r.Context(), "handler", "create", "handler", "account_product_category")
	defer span.End()
	log := local_util.LoggerFromCtx(ctx, h.logger)

	md := &types.ContextMetadata{}
	ctx = context.WithValue(ctx, constants.ContextKeyMetadata, md)
	localization.UpdateWriterContext(w, ctx)

	req, err := apc_core.DecodeCreateRequest(r)
	if err != nil {
		log.Errorf("[APCHandler][Create] decode err: %v", err)
		localization.SendErrorResponse(w, localization.ErrorInvalidRequestBody, nil, nil)
		return
	}

	if rc := apc_core.ValidateCreate(&req); rc.Code != "" {
		span.SetAttributes(attribute.String("invalid input", rc.Code))
		log.Errorf("[APCHandler][Create] validation: %s", rc.Message)
		localization.SendErrorResponse(w, rc, nil, nil)
		return
	}

	span.SetAttributes(attribute.String("apc.code", req.CBSCategoryCode))

	result, err := h.svc.Create(ctx, req)
	if err != nil {
		w = local_util.HandlePendingResponseError(ctx, w, err)
		span.RecordError(err)
		log.Errorf("[APCHandler][Create] service err: %v", err)
		localization.SendErrorByCodeResponse(w, err.Error())
		return
	}

	if md.IsMakerOnly {
		log.Infof("[APCHandler][Create] created code=%s", req.CBSCategoryCode)
		w = localization.ApplyActionCodeHeaderFromWriter(w, ctx)
		localization.SendSuccessResponse(w, localization.SuccessAPCCreated, apc_core.MapToResponse(result))
		return
	}
	log.Infof("[APCHandler][Create] request sent code=%s", req.CBSCategoryCode)
	w = localization.ApplyActionCodeHeaderFromWriter(w, ctx)
	localization.SendSuccessResponse(w, localization.SuccessAPCCreateRequestSent, apc_core.MapToResponse(result))
}

func (h *accountProductCategoryAdapter) Update(w http.ResponseWriter, r *http.Request) {
	ctx, span := local_util.TraceLogger(r.Context(), "handler", "update", "handler", "account_product_category")
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
	span.SetAttributes(attribute.String("apc.id", id))

	req, err := apc_core.DecodeUpdateRequest(r)
	if err != nil {
		log.Errorf("[APCHandler][Update] decode err: %v", err)
		localization.SendErrorResponse(w, localization.ErrorInvalidRequestBody, nil, nil)
		return
	}

	if rc := apc_core.ValidateUpdate(&req); rc.Code != "" {
		span.SetAttributes(attribute.String("invalid input", rc.Code))
		log.Errorf("[APCHandler][Update] validation: %s", rc.Code)
		localization.SendErrorResponse(w, rc, nil, nil)
		return
	}

	result, err := h.svc.Update(ctx, id, req)
	if err != nil {
		w = local_util.HandlePendingResponseError(ctx, w, err)
		span.RecordError(err)
		log.Errorf("[APCHandler][Update] service err: %v", err)
		localization.SendErrorByCodeResponse(w, err.Error())
		return
	}

	if md.IsMakerOnly {
		log.Infof("[APCHandler][Update] updated id=%s", id)
		w = localization.ApplyActionCodeHeaderFromWriter(w, ctx)
		localization.SendSuccessResponse(w, localization.SuccessAPCUpdated, apc_core.MapToResponse(result))
		return
	}
	log.Infof("[APCHandler][Update] request sent id=%s", id)
	w = localization.ApplyActionCodeHeaderFromWriter(w, ctx)
	localization.SendSuccessResponse(w, localization.SuccessAPCUpdateRequestSent, apc_core.MapToResponse(result))
}

func (h *accountProductCategoryAdapter) Delete(w http.ResponseWriter, r *http.Request) {
	ctx, span := local_util.TraceLogger(r.Context(), "handler", "delete", "handler", "account_product_category")
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
	span.SetAttributes(attribute.String("apc.id", id))

	if err := h.svc.Delete(ctx, id); err != nil {
		w = local_util.HandlePendingResponseError(ctx, w, err)
		span.RecordError(err)
		log.Errorf("[APCHandler][Delete] service err: %v", err)
		localization.SendErrorByCodeResponse(w, err.Error())
		return
	}

	if md.IsMakerOnly {
		w = localization.ApplyActionCodeHeaderFromWriter(w, ctx)
		localization.SendSuccessResponse(w, localization.SuccessBankCreatedSuccessfully, nil)
		return
	}
	log.Infof("[APCHandler][Delete] request sent id=%s", id)
	w = localization.ApplyActionCodeHeaderFromWriter(w, ctx)
	localization.SendSuccessResponse(w, localization.SuccessAPCDeleteRequestSent, nil)
}

func (h *accountProductCategoryAdapter) Enable(w http.ResponseWriter, r *http.Request) {
	ctx, span := local_util.TraceLogger(r.Context(), "handler", "enable", "handler", "account_product_category")
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
	span.SetAttributes(attribute.String("apc.id", id))

	if err := h.svc.EnableOrDisable(ctx, id, true); err != nil {
		w = local_util.HandlePendingResponseError(ctx, w, err)
		span.RecordError(err)
		log.Errorf("[APCHandler][Enable] service err: %v", err)
		localization.SendErrorByCodeResponse(w, err.Error())
		return
	}

	if md.IsMakerOnly {
		log.Infof("[APCHandler][Enable] enabled id=%s", id)
		w = localization.ApplyActionCodeHeaderFromWriter(w, ctx)
		localization.SendSuccessResponse(w, localization.SuccessAPCEnabled, nil)
		return
	}
	log.Infof("[APCHandler][Enable] request sent id=%s", id)
	w = localization.ApplyActionCodeHeaderFromWriter(w, ctx)
	localization.SendSuccessResponse(w, localization.SuccessAPCEnableRequestSent, nil)
}

func (h *accountProductCategoryAdapter) Disable(w http.ResponseWriter, r *http.Request) {
	ctx, span := local_util.TraceLogger(r.Context(), "handler", "disable", "handler", "account_product_category")
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
	span.SetAttributes(attribute.String("apc.id", id))

	if err := h.svc.EnableOrDisable(ctx, id, false); err != nil {
		w = local_util.HandlePendingResponseError(ctx, w, err)
		span.RecordError(err)
		log.Errorf("[APCHandler][Disable] service err: %v", err)
		localization.SendErrorByCodeResponse(w, err.Error())
		return
	}

	if md.IsMakerOnly {
		log.Infof("[APCHandler][Disable] disabled id=%s", id)
		w = localization.ApplyActionCodeHeaderFromWriter(w, ctx)
		localization.SendSuccessResponse(w, localization.SuccessAPCDisabled, nil)
		return
	}
	log.Infof("[APCHandler][Disable] request sent id=%s", id)
	w = localization.ApplyActionCodeHeaderFromWriter(w, ctx)
	localization.SendSuccessResponse(w, localization.SuccessAPCDisableRequestSent, nil)
}
