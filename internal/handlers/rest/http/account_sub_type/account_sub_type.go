package account_sub_type_handler

import (
	"context"
	"encoding/json"
	"net/http"

	"github.com/abelshebela/cbe-super-app-cps-action/internal/constants"
	account_sub_type_dto "github.com/abelshebela/cbe-super-app-cps-action/internal/constants/dto/account_sub_type"
	account_sub_type_interface "github.com/abelshebela/cbe-super-app-cps-action/internal/constants/interfaces/account_sub_type"
	"github.com/abelshebela/cbe-super-app-cps-action/internal/constants/localization"
	"github.com/abelshebela/cbe-super-app-cps-action/internal/constants/types"
	"github.com/abelshebela/cbe-super-app-cps-action/internal/service"
	local_util "github.com/abelshebela/cbe-super-app-cps-action/pkgs/utils"

	"github.com/go-chi/chi/v5"
	"gitlab.com/bersufekadgetachew/cbe-super-app-shared/shared/utils"
	"go.opentelemetry.io/otel/attribute"
)

type accountSubTypeAdapter struct {
	svc    service.AccountSubTypeService
	logger utils.Logger
}

func InitAccountSubTypeAdapter(svc service.AccountSubTypeService, logger utils.Logger) account_sub_type_interface.AccountSubTypeHandler {
	return &accountSubTypeAdapter{
		svc:    svc,
		logger: logger,
	}
}

func (h *accountSubTypeAdapter) CreateOneAccountSubType(w http.ResponseWriter, r *http.Request) {
	ctx, span := local_util.TraceLogger(r.Context(), "handler", "createOneAccountSubType", "handler", "account_sub_type")
	defer span.End()
	log := local_util.LoggerFromCtx(ctx, h.logger)

	md := &types.ContextMetadata{}
	ctx = context.WithValue(ctx, constants.ContextKeyMetadata, md)
	localization.UpdateWriterContext(w, ctx)

	var req account_sub_type_dto.CreateAccountSubTypeRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		log.Errorf("[AccountSubTypeH][Create] decode err: %v", err)
		localization.SendErrorResponse(w, localization.ErrorInvalidRequestBody, nil, nil)
		return
	}

	account_sub_type_dto.SanitizeCreateRequest(&req)

	if rc := account_sub_type_dto.ValidateCreateRequest(&req); rc.Code != "" {
		span.SetAttributes(attribute.String("invalid input", rc.Code))
		log.Errorf("[AccountSubTypeH][Create] validation failed: %s", rc.Code)
		localization.SendErrorResponse(w, rc, nil, nil)
		return
	}

	span.SetAttributes(
		attribute.String("account_sub_type.code", req.AccountSubTypeCode),
		attribute.String("account_sub_type.name", req.AccountSubTypeName),
	)

	if err := h.svc.CreateOneAccountSubType(ctx, req); err != nil {
		span.RecordError(err)
		log.Errorf("[AccountSubTypeH][Create] service err: %v", err)
		localization.SendErrorByCodeResponse(w, err.Error())
		return
	}

	if md.IsMakerOnly {
		log.Infof("[AccountSubTypeH][Create] created successfully code: %s", req.AccountSubTypeCode)
		w = localization.ApplyActionCodeHeaderFromWriter(w, ctx)
		localization.SendSuccessResponse(w, localization.SuccessAccountSubTypeCreated, nil)
		return
	}
	log.Infof("[AccountSubTypeH][Create] request sent code: %s", req.AccountSubTypeCode)
	w = localization.ApplyActionCodeHeaderFromWriter(w, ctx)
	localization.SendSuccessResponse(w, localization.SuccessAccountSubTypeCreateRequestSent, nil)
}

func (h *accountSubTypeAdapter) UpdateOneAccountSubType(w http.ResponseWriter, r *http.Request) {
	ctx, span := local_util.TraceLogger(r.Context(), "handler", "updateOneAccountSubType", "handler", "account_sub_type")
	defer span.End()
	log := local_util.LoggerFromCtx(ctx, h.logger)

	md := &types.ContextMetadata{}
	ctx = context.WithValue(ctx, constants.ContextKeyMetadata, md)
	localization.UpdateWriterContext(w, ctx)

	id := chi.URLParam(r, "id")
	if id == "" {
		log.Errorf("[AccountSubTypeH][Update] missing id")
		localization.SendErrorResponse(w, localization.ErrorRequiredFieldMissing, nil, nil)
		return
	}
	span.SetAttributes(attribute.String("account_sub_type.id", id))

	var req account_sub_type_dto.UpdateAccountSubTypeRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		log.Errorf("[AccountSubTypeH][Update] decode err: %v", err)
		localization.SendErrorResponse(w, localization.ErrorInvalidRequestBody, nil, nil)
		return
	}

	account_sub_type_dto.SanitizeUpdateRequest(&req)

	if rc := account_sub_type_dto.ValidateUpdateRequest(&req); rc.Code != "" {
		span.SetAttributes(attribute.String("invalid input", rc.Code))
		log.Errorf("[AccountSubTypeH][Update] validation failed: %s", rc.Code)
		localization.SendErrorResponse(w, rc, nil, nil)
		return
	}

	if err := h.svc.UpdateOneAccountSubType(ctx, id, req); err != nil {
		span.RecordError(err)
		log.Errorf("[AccountSubTypeH][Update] service err: %v", err)
		localization.SendErrorByCodeResponse(w, err.Error())
		return
	}

	if md.IsMakerOnly {
		log.Infof("[AccountSubTypeH][Update] updated successfully id: %s", id)
		w = localization.ApplyActionCodeHeaderFromWriter(w, ctx)
		localization.SendSuccessResponse(w, localization.SuccessAccountSubTypeUpdated, nil)
		return
	}
	log.Infof("[AccountSubTypeH][Update] request sent id: %s", id)
	w = localization.ApplyActionCodeHeaderFromWriter(w, ctx)
	localization.SendSuccessResponse(w, localization.SuccessAccountSubTypeUpdateRequestSent, nil)
}

func (h *accountSubTypeAdapter) DeleteOneAccountSubType(w http.ResponseWriter, r *http.Request) {
	ctx, span := local_util.TraceLogger(r.Context(), "handler", "deleteOneAccountSubType", "handler", "account_sub_type")
	defer span.End()
	log := local_util.LoggerFromCtx(ctx, h.logger)

	id := chi.URLParam(r, "id")
	if id == "" {
		log.Errorf("[AccountSubTypeH][Delete] missing id")
		localization.SendErrorResponse(w, localization.ErrorRequiredFieldMissing, nil, nil)
		return
	}
	span.SetAttributes(attribute.String("account_sub_type.id", id))

	if err := h.svc.DeleteOneAccountSubType(ctx, id); err != nil {
		span.RecordError(err)
		log.Errorf("[AccountSubTypeH][Delete] service err: %v", err)
		localization.SendErrorByCodeResponse(w, err.Error())
		return
	}

	log.Infof("[AccountSubTypeH][Delete] request sent id: %s", id)
	w = localization.ApplyActionCodeHeaderFromWriter(w, ctx)
	localization.SendSuccessResponse(w, localization.SuccessDeleteRequestCreated, nil)
}

func (h *accountSubTypeAdapter) Enable(w http.ResponseWriter, r *http.Request) {
	ctx, span := local_util.TraceLogger(r.Context(), "handler", "enableAccountSubType", "handler", "account_sub_type")
	defer span.End()
	log := local_util.LoggerFromCtx(ctx, h.logger)

	md := &types.ContextMetadata{}
	ctx = context.WithValue(ctx, constants.ContextKeyMetadata, md)
	localization.UpdateWriterContext(w, ctx)

	id := chi.URLParam(r, "id")
	if id == "" {
		log.Errorf("[AccountSubTypeH][Enable] missing id")
		localization.SendErrorResponse(w, localization.ErrorRequiredFieldMissing, nil, nil)
		return
	}
	span.SetAttributes(attribute.String("account_sub_type.id", id))

	if err := h.svc.EnableOrDisableAccountSubType(ctx, id, true); err != nil {
		span.RecordError(err)
		log.Errorf("[AccountSubTypeH][Enable] service err: %v", err)
		localization.SendErrorByCodeResponse(w, err.Error())
		return
	}

	if md.IsMakerOnly {
		log.Infof("[AccountSubTypeH][Enable] enabled id: %s", id)
		w = localization.ApplyActionCodeHeaderFromWriter(w, ctx)
		localization.SendSuccessResponse(w, localization.SuccessAccountSubTypeEnabled, nil)
		return
	}
	log.Infof("[AccountSubTypeH][Enable] request sent id: %s", id)
	w = localization.ApplyActionCodeHeaderFromWriter(w, ctx)
	localization.SendSuccessResponse(w, localization.SuccessAccountSubTypeEnableRequestSent, nil)
}

func (h *accountSubTypeAdapter) Disable(w http.ResponseWriter, r *http.Request) {
	ctx, span := local_util.TraceLogger(r.Context(), "handler", "disableAccountSubType", "handler", "account_sub_type")
	defer span.End()
	log := local_util.LoggerFromCtx(ctx, h.logger)

	md := &types.ContextMetadata{}
	ctx = context.WithValue(ctx, constants.ContextKeyMetadata, md)
	localization.UpdateWriterContext(w, ctx)

	id := chi.URLParam(r, "id")
	if id == "" {
		log.Errorf("[AccountSubTypeH][Disable] missing id")
		localization.SendErrorResponse(w, localization.ErrorRequiredFieldMissing, nil, nil)
		return
	}
	span.SetAttributes(attribute.String("account_sub_type.id", id))

	if err := h.svc.EnableOrDisableAccountSubType(ctx, id, false); err != nil {
		span.RecordError(err)
		log.Errorf("[AccountSubTypeH][Disable] service err: %v", err)
		localization.SendErrorByCodeResponse(w, err.Error())
		return
	}

	if md.IsMakerOnly {
		log.Infof("[AccountSubTypeH][Disable] disabled id: %s", id)
		w = localization.ApplyActionCodeHeaderFromWriter(w, ctx)
		localization.SendSuccessResponse(w, localization.SuccessAccountSubTypeDisabled, nil)
		return
	}
	log.Infof("[AccountSubTypeH][Disable] request sent id: %s", id)
	w = localization.ApplyActionCodeHeaderFromWriter(w, ctx)
	localization.SendSuccessResponse(w, localization.SuccessAccountSubTypeDisableRequestSent, nil)
}

func (h *accountSubTypeAdapter) GetAllAccountSubTypes(w http.ResponseWriter, r *http.Request) {
	ctx, span := local_util.TraceLogger(r.Context(), "handler", "getAllAccountSubTypes", "handler", "account_sub_type")
	defer span.End()
	log := local_util.LoggerFromCtx(ctx, h.logger)

	filterParams, err := local_util.ExtractFilterParams(r)
	if err != nil {
		localization.SendErrorByCodeResponse(w, err.Error())
		return
	}

	search := r.URL.Query().Get("search")
	if err := local_util.NoSpecialChars(search); err != nil {
		localization.SendErrorByCodeResponse(w, err.Error())
		return
	}

	result, err := h.svc.GetAllAccountSubTypes(ctx, *filterParams)
	if err != nil {
		span.RecordError(err)
		log.Errorf("[AccountSubTypeH][GetAll] service err: %v", err)
		localization.SendErrorByCodeResponse(w, err.Error())
		return
	}

	span.SetAttributes(attribute.Int("account_sub_type.count", len(result.Data)))
	log.Infof("[AccountSubTypeH][GetAll] retrieved %d records", len(result.Data))
	localization.SendSuccessResponse(w, localization.SuccessGetAllAccountSubTypes, result)
}

func (h *accountSubTypeAdapter) GetOneAccountSubType(w http.ResponseWriter, r *http.Request) {
	ctx, span := local_util.TraceLogger(r.Context(), "handler", "getOneAccountSubType", "handler", "account_sub_type")
	defer span.End()
	log := local_util.LoggerFromCtx(ctx, h.logger)

	id := chi.URLParam(r, "id")
	if id == "" {
		log.Errorf("[AccountSubTypeH][GetOne] missing id")
		localization.SendErrorResponse(w, localization.ErrorRequiredFieldMissing, nil, nil)
		return
	}
	span.SetAttributes(attribute.String("account_sub_type.id", id))

	result, err := h.svc.GetOneAccountSubType(ctx, id)
	if err != nil {
		span.RecordError(err)
		log.Errorf("[AccountSubTypeH][GetOne] service err: %v", err)
		localization.SendErrorByCodeResponse(w, err.Error())
		return
	}

	log.Infof("[AccountSubTypeH][GetOne] found id: %s", id)
	localization.SendSuccessResponse(w, localization.SuccessGetOneAccountSubType, result)
}
