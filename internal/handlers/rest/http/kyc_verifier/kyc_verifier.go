package kyc_verifier

import (
	"cbe-super-app-cps-action/internal/constants/dto/kyc_verifier"
	inbound "cbe-super-app-cps-action/internal/constants/interfaces/kyc_verifier"
	"cbe-super-app-cps-action/internal/constants/localization"
	"cbe-super-app-cps-action/internal/service"
	"encoding/json"
	"net/http"

	local_util "cbe-super-app-cps-action/pkgs/utils"

	"go.opentelemetry.io/otel/attribute"

	"github.com/go-chi/chi/v5"
	"gitlab.com/bersufekadgetachew/cbe-super-app-shared/shared/utils"
)

type kycAdapter struct {
	logger utils.Logger
	app    service.KYCVerifierService
}

func InitKYCAdapter(app service.KYCVerifierService, logger utils.Logger) inbound.KYCVerifierAdapter {
	return &kycAdapter{logger: logger, app: app}
}

func (h *kycAdapter) GetKYCList(w http.ResponseWriter, r *http.Request) {
	ctx, span := local_util.TraceLogger(r.Context(), "handler", "getKycList", "handler", "kyc")
	defer span.End()
	filterParams := local_util.ExtractFilterParams(r)
	if filterParams.Page < 0 || filterParams.PerPage < 0 {
		localization.SendErrorResponse(w, localization.ErrorInvalidRequest, nil, nil)
		return
	}
	res, err := h.app.FetchKYCList(ctx, filterParams)
	if err != nil {
		span.RecordError(err)
		localization.SendErrorByCodeResponse(w, err.Error())
		return
	}
	span.SetAttributes(attribute.Int("kyc.count", len(res.Data)))
	localization.SendSuccessResponse(w, localization.SuccessKYCFetched, res)
}

func (h *kycAdapter) GetKYCByID(w http.ResponseWriter, r *http.Request) {
	ctx, span := local_util.TraceLogger(r.Context(), "handler", "getKycById", "handler", "kyc")
	defer span.End()
	id := chi.URLParam(r, "id")
	if id == "" {
		localization.SendErrorResponse(w, localization.ErrorInvalidRequest, nil, nil)
		return
	}
	span.SetAttributes(attribute.String("kyc.id", id))
	res, err := h.app.FetchKYCByID(ctx, id)
	if err != nil {
		span.RecordError(err)
		localization.SendErrorByCodeResponse(w, err.Error())
		return
	}
	localization.SendSuccessResponse(w, localization.SuccessKYCFetched, res)
}

func (h *kycAdapter) UpdateKYC(w http.ResponseWriter, r *http.Request) {
	ctx, span := local_util.TraceLogger(r.Context(), "handler", "updateKyc", "handler", "kyc")
	defer span.End()
	var req kyc_verifier.UpdateKYCRequest
	id := chi.URLParam(r, "id")

	if id == "" {
		localization.SendErrorResponse(w, localization.ErrorInvalidRequest, nil, nil)
		return
	}

	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		span.RecordError(err)
		localization.SendErrorByCodeResponse(w, err.Error())
		return
	}
	span.SetAttributes(attribute.String("kyc.id", id))
	if err := h.app.UpdateKYC(ctx, id, req); err != nil {
		span.RecordError(err)
		localization.SendErrorByCodeResponse(w, err.Error())
		return
	}
	localization.SendSuccessResponse(w, localization.SuccessKYCUpdatedRequestSent, nil)
}

func (h *kycAdapter) ApproveKYC(w http.ResponseWriter, r *http.Request) {
	ctx, span := local_util.TraceLogger(r.Context(), "handler", "approveKyc", "handler", "kyc")
	defer span.End()
	var req kyc_verifier.ApproveKYCRequest
	id := chi.URLParam(r, "id")
	if id == "" {
		localization.SendErrorResponse(w, localization.ErrorInvalidRequest, nil, nil)
		return
	}

	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		span.RecordError(err)
		localization.SendErrorByCodeResponse(w, err.Error())
		return
	}

	span.SetAttributes(attribute.String("kyc.id", id))
	if err := h.app.ApproveKYC(ctx, id, req); err != nil {
		span.RecordError(err)
		localization.SendErrorByCodeResponse(w, err.Error())
		return
	}
	localization.SendSuccessResponse(w, localization.SuccessKYCApproved, nil)
}
