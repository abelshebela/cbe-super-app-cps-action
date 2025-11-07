package kyc_verifier

import (
	"cbe-super-app-cps-action/internal/constants/dto/kyc_verifier"
	inbound "cbe-super-app-cps-action/internal/constants/interfaces/kyc_verifier"
	"cbe-super-app-cps-action/internal/constants/localization"
	"cbe-super-app-cps-action/internal/service"
	"encoding/json"
	"net/http"

	local_util "cbe-super-app-cps-action/pkgs/utils"

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
	filterParams := local_util.ExtractFilterParams(r)
	if filterParams.Page < 0 || filterParams.PerPage < 0 {
		localization.SendErrorResponse(w, localization.ErrorInvalidRequest, nil, nil)
		return
	}
	ctx := r.Context()
	res, err := h.app.FetchKYCList(ctx, filterParams)
	if err != nil {
		localization.SendErrorByCodeResponse(w, err.Error())
		return
	}
	localization.SendSuccessResponse(w, localization.SuccessKYCFetched, res)
}

func (h *kycAdapter) GetKYCByID(w http.ResponseWriter, r *http.Request) {
	id := chi.URLParam(r, "id")
	if id == "" {
		localization.SendErrorResponse(w, localization.ErrorInvalidRequest, nil, nil)
		return
	}
	ctx := r.Context()
	res, err := h.app.FetchKYCByID(ctx, id)
	if err != nil {
		localization.SendErrorByCodeResponse(w, err.Error())
		return
	}
	localization.SendSuccessResponse(w, localization.SuccessKYCFetched, res)
}

func (h *kycAdapter) UpdateKYC(w http.ResponseWriter, r *http.Request) {
	var req kyc_verifier.UpdateKYCRequest
	id := chi.URLParam(r, "id")

	if id == "" {
		localization.SendErrorResponse(w, localization.ErrorInvalidRequest, nil, nil)
		return
	}

	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		localization.SendErrorByCodeResponse(w, err.Error())
		return
	}
	if err := h.app.UpdateKYC(r.Context(), id, req); err != nil {
		localization.SendErrorByCodeResponse(w, err.Error())
		return
	}
	localization.SendSuccessResponse(w, localization.SuccessKYCUpdatedRequestSent, nil)
}

func (h *kycAdapter) ApproveKYC(w http.ResponseWriter, r *http.Request) {
	var req kyc_verifier.ApproveKYCRequest
	id := chi.URLParam(r, "id")
	if id == "" {
		localization.SendErrorResponse(w, localization.ErrorInvalidRequest, nil, nil)
		return
	}

	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		localization.SendErrorByCodeResponse(w, err.Error())
		return
	}

	if err := h.app.ApproveKYC(r.Context(), id, req); err != nil {
		localization.SendErrorByCodeResponse(w, err.Error())
		return
	}
	localization.SendSuccessResponse(w, localization.SuccessKYCApproved, nil)
}
