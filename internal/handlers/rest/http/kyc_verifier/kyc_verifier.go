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

// GetKYCList
//
//	@Summary		Get KYC List
//	@Description	Fetch KYC records with pagination, filtering, and search. Filterable fields: kyc_status, kyc_level, kyc_approved, enabled. Searchable field: kyc_status.
//	@Tags			KYC Verifier
//	@Accept			json
//	@Produce		json
//	@Param			page			query		int																		false	"Page number (default 1)"
//	@Param			per_page		query		int																		false	"Items per page (default 10, max 100)"
//	@Param			kyc_status		query		string																	false	"Filter by KYC status"
//	@Param			kyc_level		query		string																	false	"Filter by KYC level"
//	@Param			kyc_approved	query		bool																	false	"Filter by KYC approved"
//	@Param			enabled			query		bool																	false	"Filter by enabled status"
//	@Param			search			query		string																	false	"Search term (searches kyc_status)"
//	@Success		200				{object}	localization.StandardResponse{data=[]kyc_verifier.KYCVerifierResponse}	"KYC records fetched successfully"
//	@Failure		400				{object}	localization.StandardResponse{data=nil}									"Invalid query params"
//	@Failure		500				{object}	localization.StandardResponse{data=nil}									"Internal server error"
//	@Security		BearerAuth
//	@Router			/kyc_verifier [get]
func (h *kycAdapter) GetKYCList(w http.ResponseWriter, r *http.Request) {
	ctx, span := local_util.TraceLogger(r.Context(), "handler", "getKycList", "handler", "kyc")
	defer span.End()
	filterParams := local_util.ExtractFilterParams(r)

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

// GetKYCByID
//
//	@Summary		Get KYC by ID
//	@Description	Retrieve a KYC record by its ID
//	@Tags			KYC Verifier
//	@Accept			json
//	@Produce		json
//	@Param			id	path		string																	true	"KYC record ID"
//	@Success		200	{object}	localization.StandardResponse{data=kyc_verifier.KYCVerifierResponse}	"KYC record retrieved successfully"
//	@Failure		400	{object}	localization.StandardResponse{data=nil}									"Invalid ID"
//	@Failure		404	{object}	localization.StandardResponse{data=nil}									"KYC record not found"
//	@Failure		500	{object}	localization.StandardResponse{data=nil}									"Internal server error"
//	@Security		BearerAuth
//	@Router			/kyc_verifier/{id} [get]
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

// UpdateKYC
//
//	@Summary		Update KYC
//	@Description	Update a KYC record by its ID
//	@Tags			KYC Verifier
//	@Accept			json
//	@Produce		json
//	@Param			id		path		string									true	"KYC record ID"
//	@Param			body	body		kyc_verifier.UpdateKYCRequest			true	"Update KYC request body"
//	@Success		200		{object}	localization.StandardResponse{data=nil}	"KYC record updated successfully"
//	@Failure		400		{object}	localization.StandardResponse{data=nil}	"Invalid ID or request body"
//	@Failure		500		{object}	localization.StandardResponse{data=nil}	"Internal server error"
//	@Security		BearerAuth
//	@Router			/kyc_verifier/update/{id} [patch]
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

// ApproveKYC
//
//	@Summary		Approve KYC
//	@Description	Approve a KYC record by its ID
//	@Tags			KYC Verifier
//	@Accept			json
//	@Produce		json
//	@Param			id		path		string									true	"KYC record ID"
//	@Param			body	body		kyc_verifier.ApproveKYCRequest			true	"Approve KYC request body"
//	@Success		200		{object}	localization.StandardResponse{data=nil}	"KYC record approved successfully"
//	@Failure		400		{object}	localization.StandardResponse{data=nil}	"Invalid ID or request body"
//	@Failure		500		{object}	localization.StandardResponse{data=nil}	"Internal server error"
//	@Security		BearerAuth
//	@Router			/kyc_verifier/approve/{id} [patch]
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
