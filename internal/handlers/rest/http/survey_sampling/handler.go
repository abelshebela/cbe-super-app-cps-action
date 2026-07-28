package survey_sampling_handler

import (
	"context"
	"encoding/json"
	"net/http"
	"strings"
	"time"

	survey_sampling_dto "cbe-super-app-cps-action/internal/constants/dto/survey_sampling"
	survey_sampling_interface "cbe-super-app-cps-action/internal/constants/interfaces/survey_sampling"
	"cbe-super-app-cps-action/internal/constants/localization"
	imodel "cbe-super-app-cps-action/internal/constants/model"
	"cbe-super-app-cps-action/internal/constants"
	"cbe-super-app-cps-action/internal/constants/types"
	"cbe-super-app-cps-action/internal/service"
	local_util "cbe-super-app-cps-action/pkgs/utils"

	"github.com/go-chi/chi/v5"
	"gitlab.com/bersufekadgetachew/cbe-super-app-shared/shared/utils"
)

type SurveySamplingHandler struct {
	service service.SurveySamplingService
	logger  utils.Logger
}

func NewSurveySamplingHandler(svc service.SurveySamplingService, logger utils.Logger) survey_sampling_interface.SurveySamplingInbound {
	return &SurveySamplingHandler{service: svc, logger: logger}
}

func (h *SurveySamplingHandler) GetAll(w http.ResponseWriter, r *http.Request) {
	log := local_util.LoggerFromCtx(r.Context(), h.logger)

	filterParams, err := local_util.ExtractFilterParams(r)
	if err != nil {
		localization.SendErrorByCodeResponse(w, err.Error())
		return
	}

	resp, err := h.service.GetAll(r.Context(), *filterParams)
	if err != nil {
		log.Errorf("[SurveySamplingHandler][GetAll] service error: %v", err)
		localization.SendErrorByCodeResponse(w, err.Error())
		return
	}

	localization.SendSuccessResponse(w, localization.SuccessSurveySamplingsFetchedSuccessfully, resp)
}

func (h *SurveySamplingHandler) GetByID(w http.ResponseWriter, r *http.Request) {
	log := local_util.LoggerFromCtx(r.Context(), h.logger)

	id := chi.URLParam(r, "id")
	if strings.TrimSpace(id) == "" {
		localization.SendErrorResponse(w, localization.ErrorRequiredFieldMissing, nil, nil)
		return
	}

	config, err := h.service.GetByID(r.Context(), id)
	if err != nil {
		log.Errorf("[SurveySamplingHandler][GetByID] service error: %v", err)
		localization.SendErrorByCodeResponse(w, err.Error())
		return
	}

	localization.SendSuccessResponse(w, localization.SuccessSurveySamplingFetchedSuccessfully, config)
}

func (h *SurveySamplingHandler) Create(w http.ResponseWriter, r *http.Request) {
	md := &types.ContextMetadata{}
	ctx := context.WithValue(r.Context(), constants.ContextKeyMetadata, md)
	localization.UpdateWriterContext(w, ctx)
	log := local_util.LoggerFromCtx(ctx, h.logger)

	var body survey_sampling_dto.CreateSurveySamplingRequest
	if err := json.NewDecoder(r.Body).Decode(&body); err != nil {
		http.Error(w, "invalid request body", http.StatusBadRequest)
		return
	}

	if strings.TrimSpace(body.Name) == "" || strings.TrimSpace(body.Method) == "" {
		localization.SendBadRequestResponse(w, localization.ErrorRequiredFieldMissing.Code)
		return
	}

	config := imodel.SurveySamplingConfig{
		Name:      strings.TrimSpace(body.Name),
		Method:    imodel.SamplingMethod(strings.ToUpper(strings.TrimSpace(body.Method))),
		Config:    body.Config,
		SurveyURL: strings.TrimSpace(body.SurveyURL),
		CreatedAt: time.Now(),
		UpdatedAt: time.Now(),
	}

	if err := h.service.Create(ctx, config); err != nil {
		w = local_util.HandlePendingResponseError(ctx, w, err)
		localization.SendErrorByCodeResponse(w, err.Error())
		return
	}

	w = localization.ApplyActionCodeHeaderFromWriter(w, ctx)
	if md.IsMakerOnly {
		localization.SendSuccessResponse(w, localization.SuccessSurveySamplingCreatedSP, nil)
	} else {
		localization.SendSuccessResponse(w, localization.SuccessSurveySamplingCreatedRequestSent, nil)
	}

	log.Infof("[SurveySamplingHandler][Create] survey sampling config created for method=%s", config.Method)
}

func (h *SurveySamplingHandler) Update(w http.ResponseWriter, r *http.Request) {
	md := &types.ContextMetadata{}
	ctx := context.WithValue(r.Context(), constants.ContextKeyMetadata, md)
	localization.UpdateWriterContext(w, ctx)
	log := local_util.LoggerFromCtx(ctx, h.logger)

	id := chi.URLParam(r, "id")
	if strings.TrimSpace(id) == "" {
		localization.SendErrorResponse(w, localization.ErrorRequiredFieldMissing, nil, nil)
		return
	}

	var body survey_sampling_dto.UpdateSurveySamplingRequest
	if err := json.NewDecoder(r.Body).Decode(&body); err != nil {
		http.Error(w, "invalid request body", http.StatusBadRequest)
		return
	}

	config := imodel.SurveySamplingConfig{
		Name:      strings.TrimSpace(body.Name),
		Config:    body.Config,
		SurveyURL: strings.TrimSpace(body.SurveyURL),
		UpdatedAt: time.Now(),
	}

	if err := h.service.Update(ctx, id, config); err != nil {
		w = local_util.HandlePendingResponseError(ctx, w, err)
		log.Errorf("[SurveySamplingHandler][Update] service error: %v", err)
		localization.SendErrorByCodeResponse(w, err.Error())
		return
	}

	w = localization.ApplyActionCodeHeaderFromWriter(w, ctx)
	if md.IsMakerOnly {
		localization.SendSuccessResponse(w, localization.SuccessSurveySamplingUpdatedSP, nil)
	} else {
		localization.SendSuccessResponse(w, localization.SuccessSurveySamplingUpdatedRequestSent, nil)
	}
}

func (h *SurveySamplingHandler) Enable(w http.ResponseWriter, r *http.Request) {
	md := &types.ContextMetadata{}
	ctx := context.WithValue(r.Context(), constants.ContextKeyMetadata, md)
	localization.UpdateWriterContext(w, ctx)
	log := local_util.LoggerFromCtx(ctx, h.logger)

	id := chi.URLParam(r, "id")
	if strings.TrimSpace(id) == "" {
		localization.SendErrorResponse(w, localization.ErrorRequiredFieldMissing, nil, nil)
		return
	}

	if err := h.service.Enable(ctx, id); err != nil {
		w = local_util.HandlePendingResponseError(ctx, w, err)
		log.Errorf("[SurveySamplingHandler][Enable] service error: %v", err)
		localization.SendErrorByCodeResponse(w, err.Error())
		return
	}

	w = localization.ApplyActionCodeHeaderFromWriter(w, ctx)
	if md.IsMakerOnly {
		localization.SendSuccessResponse(w, localization.SuccessSurveySamplingEnabledSP, nil)
	} else {
		localization.SendSuccessResponse(w, localization.SuccessSurveySamplingEnabledRequestSent, nil)
	}
}

func (h *SurveySamplingHandler) Disable(w http.ResponseWriter, r *http.Request) {
	md := &types.ContextMetadata{}
	ctx := context.WithValue(r.Context(), constants.ContextKeyMetadata, md)
	localization.UpdateWriterContext(w, ctx)
	log := local_util.LoggerFromCtx(ctx, h.logger)

	id := chi.URLParam(r, "id")
	if strings.TrimSpace(id) == "" {
		localization.SendErrorResponse(w, localization.ErrorRequiredFieldMissing, nil, nil)
		return
	}

	if err := h.service.Disable(ctx, id); err != nil {
		w = local_util.HandlePendingResponseError(ctx, w, err)
		log.Errorf("[SurveySamplingHandler][Disable] service error: %v", err)
		localization.SendErrorByCodeResponse(w, err.Error())
		return
	}

	w = localization.ApplyActionCodeHeaderFromWriter(w, ctx)
	if md.IsMakerOnly {
		localization.SendSuccessResponse(w, localization.SuccessSurveySamplingDisabledSP, nil)
	} else {
		localization.SendSuccessResponse(w, localization.SuccessSurveySamplingDisabledRequestSent, nil)
	}
}

func (h *SurveySamplingHandler) Delete(w http.ResponseWriter, r *http.Request) {
	md := &types.ContextMetadata{}
	ctx := context.WithValue(r.Context(), constants.ContextKeyMetadata, md)
	localization.UpdateWriterContext(w, ctx)
	log := local_util.LoggerFromCtx(ctx, h.logger)

	id := chi.URLParam(r, "id")
	if strings.TrimSpace(id) == "" {
		localization.SendErrorResponse(w, localization.ErrorRequiredFieldMissing, nil, nil)
		return
	}

	if err := h.service.Delete(ctx, id); err != nil {
		w = local_util.HandlePendingResponseError(ctx, w, err)
		log.Errorf("[SurveySamplingHandler][Delete] service error: %v", err)
		localization.SendErrorByCodeResponse(w, err.Error())
		return
	}

	w = localization.ApplyActionCodeHeaderFromWriter(w, ctx)
	if md.IsMakerOnly {
		localization.SendSuccessResponse(w, localization.SuccessSurveySamplingDeletedSP, nil)
	} else {
		localization.SendSuccessResponse(w, localization.SuccessSurveySamplingDeletedRequestSent, nil)
	}
}
