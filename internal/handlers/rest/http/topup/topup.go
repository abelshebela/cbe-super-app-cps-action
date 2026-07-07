package topup

import (
	topupDto "cbe-super-app-cps-action/internal/constants/dto/topup"
	topupInbound "cbe-super-app-cps-action/internal/constants/interfaces/topup"
	"cbe-super-app-cps-action/internal/constants/types"
	"context"
	"net/http"

	"gitlab.com/bersufekadgetachew/cbe-super-app-shared/shared/entities/model"

	"cbe-super-app-cps-action/internal/constants/localization"
	topupcore "cbe-super-app-cps-action/internal/handlers/rest/http/topup/core"
	"cbe-super-app-cps-action/internal/service"
	local_util "cbe-super-app-cps-action/pkgs/utils"

	"cbe-super-app-cps-action/internal/constants"

	"github.com/go-chi/chi/v5"
	"gitlab.com/bersufekadgetachew/cbe-super-app-shared/shared/utils"
	"go.opentelemetry.io/otel/attribute"
	"go.opentelemetry.io/otel/trace"
)

type PaginatedTopupResponse types.PaginatedResponse[[]*model.Topup]

// TopupRequestForm mirrors topupDto.TopupRequest for Swagger formData only (topupDto.TopupRequest has *multipart.FileHeader which swag cannot parse).
type TopupRequestForm struct {
	Name  string `form:"name" json:"name"`
	Code  string `form:"code" json:"code"`
	Self  bool   `form:"self" json:"self"`
	Other bool   `form:"other" json:"other"`
	Agent bool   `form:"agent" json:"agent"`
}

type topupAdapter struct {
	topupApp service.TopupService
	logger   utils.Logger
}

func InitTopupAdapter(topupApp service.TopupService, logger utils.Logger) topupInbound.TopupAdapter {
	return &topupAdapter{
		topupApp: topupApp,
		logger:   logger,
	}
}

func (a *topupAdapter) CreateTopup(w http.ResponseWriter, r *http.Request) {
	ctx, span := local_util.TraceLogger(r.Context(), "handler", "topup", "topupAdapter", "CreateTopup")
	defer span.End()
	log := local_util.LoggerFromCtx(ctx, a.logger)
	md := &types.ContextMetadata{}
	ctx = context.WithValue(ctx, constants.ContextKeyMetadata, md)
	localization.UpdateWriterContext(w, ctx)

	log.Infof("[TopupH][Create] called")

	var req topupDto.TopupRequest
	req, err := topupcore.ParseTopupRequestFromMultipartForm(r, true)
	if err != nil {
		span.AddEvent("Failed to parse multipart form", trace.WithAttributes(attribute.String("error", err.Error())))
		log.Errorf("[TopupH][Create] parse form err")
	}
	log.Infof("[TopupH] req: %+v", req)

	if err := req.AggregatedValidate(true); err != nil {
		span.AddEvent("Validation error", trace.WithAttributes(attribute.String("error", err.Error())))
		log.Errorf("[TopupH] validate err: %v", err)
		localization.SendBadRequestResponse(w, err.Error())
		return
	}

	if err := a.topupApp.CreateTopup(ctx, req); err != nil {
		w = local_util.HandlePendingResponseError(ctx, w, err)
		span.AddEvent("Service error", trace.WithAttributes(attribute.String("error", err.Error())))
		log.Errorf("[TopupH][Create] svc err: %v", err)
		localization.SendErrorByCodeResponse(w, err.Error())
		return
	}

	if md.IsMakerOnly {
		w = localization.ApplyActionCodeHeaderFromWriter(w, ctx)
		localization.SendSuccessResponse(w, localization.SuccessAccessListSegmentationCreatedSP, nil)
	} else {

		span.AddEvent("Topup creation request submitted")
		log.Infof("[TopupH][Create] request submitted")
		localization.SendSuccessResponse(w, localization.SuccessTopupCreationRequestSent, nil)
	}
}

func (a *topupAdapter) UpdateTopup(w http.ResponseWriter, r *http.Request) {
	ctx, span := local_util.TraceLogger(r.Context(), "handler", "topup", "topupAdapter", "UpdateTopup")
	defer span.End()
	log := local_util.LoggerFromCtx(ctx, a.logger)
	md := &types.ContextMetadata{}
	ctx = context.WithValue(ctx, constants.ContextKeyMetadata, md)
	localization.UpdateWriterContext(w, ctx)

	id := chi.URLParam(r, "id")
	log.Infof("[TopupH][Update] id: %s", id)

	if id == "" {
		span.AddEvent("Missing topup ID", trace.WithAttributes(attribute.String("error", "topup ID required")))
		log.Errorf("[TopupH][Update] id required")
		localization.SendErrorResponse(w, localization.ErrorTopupIDRequired, nil, nil)
		return
	}

	req, err := topupcore.ParseTopupRequestFromMultipartForm(r, false)
	if err != nil {
		span.AddEvent("Failed to parse multipart form", trace.WithAttributes(attribute.String("error", err.Error()), attribute.String("id", id)))
		log.Errorf("[TopupH][Update] parse form err: %v", err)
		localization.SendBadRequestResponse(w, err.Error())
		return
	}
	log.Infof("[TopupH] req: %+v", req)

	if err := req.AggregatedValidate(false); err != nil {
		span.AddEvent("Validation error", trace.WithAttributes(attribute.String("error", err.Error()), attribute.String("id", id)))
		log.Errorf("[TopupH] validate err: %v", err)
		localization.SendBadRequestResponse(w, err.Error())
		return
	}

	if req.IsEmpty() {
		span.AddEvent("Empty payload", trace.WithAttributes(attribute.String("id", id)))
		log.Warnf("[TopupH][Update] empty payload id: %s", id)
		localization.SendErrorResponse(w, localization.ErrorTopupUpdateEmptyPayload, nil, nil)
		return
	}

	if err := a.topupApp.UpdateTopup(ctx, id, req); err != nil {
		w = local_util.HandlePendingResponseError(ctx, w, err)
		span.AddEvent("Service error", trace.WithAttributes(attribute.String("error", err.Error()), attribute.String("id", id)))
		log.Errorf("[TopupH][Update] svc err id: %s: %v", id, err)
		localization.SendErrorByCodeResponse(w, err.Error())
		return
	}
	if md.IsMakerOnly {
		w = localization.ApplyActionCodeHeaderFromWriter(w, ctx)
		localization.SendSuccessResponse(w, localization.SuccessTopupUpdatedSP, nil)
	} else {
		span.AddEvent("Topup update request submitted", trace.WithAttributes(attribute.String("id", id)))
		log.Infof("[TopupH][Update] request submitted id: %s", id)
		w = localization.ApplyActionCodeHeaderFromWriter(w, ctx)
		localization.SendSuccessResponse(w, localization.SuccessTopupUpdateRequestSent, nil)

	}
}

func (a *topupAdapter) DeleteTopup(w http.ResponseWriter, r *http.Request) {
	ctx, span := local_util.TraceLogger(r.Context(), "handler", "topup", "topupAdapter", "DeleteTopup")
	defer span.End()
	log := local_util.LoggerFromCtx(ctx, a.logger)
	md := &types.ContextMetadata{}
	ctx = context.WithValue(ctx, constants.ContextKeyMetadata, md)
	localization.UpdateWriterContext(w, ctx)

	id := chi.URLParam(r, "id")
	log.Infof("[TopupH][Delete] id: %s", id)

	if id == "" {
		span.AddEvent("Missing topup ID", trace.WithAttributes(attribute.String("error", "topup ID required")))
		localization.SendErrorResponse(w, localization.ErrorTopupIDRequired, nil, nil)
		return
	}

	if err := a.topupApp.DeleteTopup(ctx, id); err != nil {
		w = local_util.HandlePendingResponseError(ctx, w, err)
		span.AddEvent("Service error", trace.WithAttributes(attribute.String("error", err.Error()), attribute.String("id", id)))
		localization.SendErrorByCodeResponse(w, err.Error())
		return
	}
	if md.IsMakerOnly {
		w = localization.ApplyActionCodeHeaderFromWriter(w, ctx)
		localization.SendSuccessResponse(w, localization.SuccessTopupDeletedSP, nil)
	} else {
		w = localization.ApplyActionCodeHeaderFromWriter(w, ctx)
		localization.SendSuccessResponse(w, localization.SuccessTopupDeletedRequestSent, nil)
		span.AddEvent("Topup deleted", trace.WithAttributes(attribute.String("id", id)))
		w = localization.ApplyActionCodeHeaderFromWriter(w, ctx)
		localization.SendSuccessResponse(w, localization.SuccessTopupDeletedRequestSent, nil)
	}
}

func (a *topupAdapter) Enable(w http.ResponseWriter, r *http.Request) {
	ctx, span := local_util.TraceLogger(r.Context(), "handler", "topup", "topupAdapter", "Enable")
	defer span.End()
	log := local_util.LoggerFromCtx(ctx, a.logger)
	md := &types.ContextMetadata{}
	ctx = context.WithValue(ctx, constants.ContextKeyMetadata, md)
	localization.UpdateWriterContext(w, ctx)

	id := chi.URLParam(r, "id")
	log.Infof("[TopupH][Enable] id: %s", id)

	if id == "" {
		span.AddEvent("Missing topup ID", trace.WithAttributes(attribute.String("error", "topup ID required")))
		localization.SendErrorResponse(w, localization.ErrorTopupIDRequired, nil, nil)
		return
	}

	if err := a.topupApp.EnableOrDisableTopup(ctx, id, true); err != nil {
		w = local_util.HandlePendingResponseError(ctx, w, err)
		span.AddEvent("Service error", trace.WithAttributes(attribute.String("error", err.Error()), attribute.String("id", id)))
		localization.SendErrorByCodeResponse(w, err.Error())
		return
	}
	if md.IsMakerOnly {
		w = localization.ApplyActionCodeHeaderFromWriter(w, ctx)
		localization.SendSuccessResponse(w, localization.SuccessTopupEnabledSP, nil)
	} else {
		span.AddEvent("Topup enable request submitted", trace.WithAttributes(attribute.String("id", id)))
		w = localization.ApplyActionCodeHeaderFromWriter(w, ctx)
		localization.SendSuccessResponse(w, localization.SuccessTopupEnableRequestSubmitted, nil)
	}
}

func (a *topupAdapter) Disable(w http.ResponseWriter, r *http.Request) {
	ctx, span := local_util.TraceLogger(r.Context(), "handler", "topup", "topupAdapter", "Disable")
	defer span.End()
	log := local_util.LoggerFromCtx(ctx, a.logger)
	md := &types.ContextMetadata{}
	ctx = context.WithValue(ctx, constants.ContextKeyMetadata, md)
	localization.UpdateWriterContext(w, ctx)

	id := chi.URLParam(r, "id")
	log.Infof("[TopupH][Disable] id: %s", id)

	if id == "" {
		span.AddEvent("Missing topup ID", trace.WithAttributes(attribute.String("error", "topup ID required")))
		localization.SendErrorResponse(w, localization.ErrorTopupIDRequired, nil, nil)
		return
	}

	if err := a.topupApp.EnableOrDisableTopup(ctx, id, false); err != nil {
		w = local_util.HandlePendingResponseError(ctx, w, err)
		span.AddEvent("Service error", trace.WithAttributes(attribute.String("error", err.Error()), attribute.String("id", id)))
		localization.SendErrorByCodeResponse(w, err.Error())
		return
	}
	if md.IsMakerOnly {
		w = localization.ApplyActionCodeHeaderFromWriter(w, ctx)
		localization.SendSuccessResponse(w, localization.SuccessTopupDisabledSP, nil)
	} else {
		span.AddEvent("Topup disable request submitted", trace.WithAttributes(attribute.String("id", id)))
		w = localization.ApplyActionCodeHeaderFromWriter(w, ctx)
		localization.SendSuccessResponse(w, localization.SuccessTopupDisableRequestSubmitted, nil)
	}
}

func (a *topupAdapter) GetTopup(w http.ResponseWriter, r *http.Request) {
	ctx, span := local_util.TraceLogger(r.Context(), "handler", "topup", "topupAdapter", "GetTopup")
	defer span.End()
	log := local_util.LoggerFromCtx(ctx, a.logger)

	id := chi.URLParam(r, "id")
	log.Infof("[TopupH][GetByID] id: %s", id)

	if id == "" {
		span.AddEvent("Missing topup ID", trace.WithAttributes(attribute.String("error", "topup ID required")))
		localization.SendErrorResponse(w, localization.ErrorTopupIDRequired, nil, nil)
		return
	}

	topup, err := a.topupApp.GetTopup(ctx, id)
	if err != nil {
		span.AddEvent("Service error", trace.WithAttributes(attribute.String("error", err.Error()), attribute.String("id", id)))
		localization.SendErrorByCodeResponse(w, err.Error())
		return
	}

	span.AddEvent("Topup retrieved", trace.WithAttributes(attribute.String("id", id)))
	localization.SendSuccessResponse(w, localization.SuccessTopupRetrieved, topup)
}

func (a *topupAdapter) GetAllTopup(w http.ResponseWriter, r *http.Request) {
	ctx, span := local_util.TraceLogger(r.Context(), "handler", "topup", "topupAdapter", "GetAllTopup")
	defer span.End()
	log := local_util.LoggerFromCtx(ctx, a.logger)

	filterParams, err := local_util.ExtractFilterParams(r)
	if err != nil {
		localization.SendErrorByCodeResponse(w, err.Error())
		return
	}

	search := r.URL.Query().Get("search")
	filter := r.URL.Query().Get("filter")

	if err := local_util.NoSpecialChars(search); err != nil {
		localization.SendErrorByCodeResponse(w, err.Error())
		return
	}

	if err := local_util.NoSpecialChars(filter); err != nil {
		localization.SendErrorByCodeResponse(w, err.Error())
		return
	}

	log.Infof("[TopupH][GetAll] filter: %+v", filter)

	list, err := a.topupApp.GetAllTopup(ctx, *filterParams)
	if err != nil {
		span.AddEvent("Service error", trace.WithAttributes(attribute.String("error", err.Error())))
		log.Errorf("[TopupH][GetAll] svc err: %v", err)
		localization.SendErrorByCodeResponse(w, err.Error())
		return
	}

	span.AddEvent("Topups retrieved", trace.WithAttributes(attribute.Int("count", len(list.Data))))
	log.Infof("[TopupH][GetAll] ok")
	localization.SendSuccessResponse(w, localization.SuccessTopupsRetrieved, list)
}
