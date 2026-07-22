package customer

import (
	"cbe-super-app-cps-action/internal/constants"
	kyc_dto "cbe-super-app-cps-action/internal/constants/dto/customer_kyc"
	inbound "cbe-super-app-cps-action/internal/constants/interfaces/customer_kyc"
	"cbe-super-app-cps-action/internal/constants/localization"
	"cbe-super-app-cps-action/internal/constants/types"
	"cbe-super-app-cps-action/internal/service"
	util "cbe-super-app-cps-action/pkgs/utils"
	"context"
	"encoding/json"
	"net/http"
	"strings"
	"time"

	"github.com/go-chi/chi/v5"
	"gitlab.com/bersufekadgetachew/cbe-super-app-shared/shared/utils"
)

type customerKYCAdapter struct {
	svc    service.CustomerKYCService
	logger utils.Logger
}

func NewCustomerKYCAdapter(kycService service.CustomerKYCService, logger utils.Logger) inbound.CustomerKYC {
	return &customerKYCAdapter{
		svc:    kycService,
		logger: logger,
	}
}

// GetAllKYCRequests returns all KYC requests with pagination
//
//	@Summary		Get All KYC Requests
//	@Description	Returns a paginated list of KYC requests.
//	@Tags			Customer KYC
//	@Accept			json
//	@Produce		json
//	@Param			page		query		int		false	"Page number"
//	@Param			per_page	query		int		false	"Items per page"
//	@Param			sort_by		query		string	false	"Field to sort by"
//	@Param			order		query		string	false	"Sort order (asc/desc)"
//	@Success		200			{object}	localization.StandardResponse{data=object}	"KYC requests fetched successfully (data: paginated list with docs and meta)"
//	@Failure		500			{object}	localization.StandardResponse{data=nil}											"Internal server error"
//	@Security		BearerAuth
//	@Router			/customers/kyc [get]
func (c *customerKYCAdapter) GetAllKYCRequests(w http.ResponseWriter, r *http.Request) {
	ctx, span := util.TraceLogger(r.Context(), "handler", "GetAllKYCRequests", "handler", "customer_kyc")
	defer span.End()
	log := util.LoggerFromCtx(ctx, c.logger)

	filterParam, err := util.ExtractFilterParams(r)
	if err != nil {
		localization.SendErrorByCodeResponse(w, err.Error())
		return
	}

	if err := util.NoSpecialChars(filterParam.Search); err != nil {
		localization.SendErrorByCodeResponse(w, err.Error())
		return
	}

	res, err := c.svc.FindAllWithPagination(ctx, filterParam)
	if err != nil {
		log.Errorf("[GetAllKYCRequests] failed to fetch KYC requests: %v", err)
		localization.SendErrorByCodeResponse(w, err.Error())
		return
	}

	localization.SendSuccessResponse(w, localization.SuccessDataRetrieved, res)
}

// GetKYCRequest returns a specific KYC request by ID
//
//	@Summary		Get KYC Request By ID
//	@Description	Returns a single KYC request by its ID.
//	@Tags			Customer KYC
//	@Accept			json
//	@Produce		json
//	@Param			id	path		string																	true	"KYC Request ID"
//	@Success		200	{object}	localization.StandardResponse{data=object}	"KYC request fetched successfully"
//	@Failure		400	{object}	localization.StandardResponse{data=nil}									"Bad request - ID required"
//	@Failure		404	{object}	localization.StandardResponse{data=nil}									"KYC request not found"
//	@Failure		500	{object}	localization.StandardResponse{data=nil}									"Internal server error"
//	@Security		BearerAuth
//	@Router			/customers/kyc/{id} [get]
func (c *customerKYCAdapter) GetKYCRequest(w http.ResponseWriter, r *http.Request) {
	ctx, span := util.TraceLogger(r.Context(), "handler", "GetKYCRequest", "handler", "customer_kyc")
	defer span.End()
	log := util.LoggerFromCtx(ctx, c.logger)

	id := chi.URLParam(r, "id")
	if id == "" {
		localization.SendBadRequestResponse(w, localization.ErrorIdNotSetOnQueryParam.Code)
		return
	}

	res, err := c.svc.FindByID(ctx, id)
	if err != nil {
		log.Errorf("[GetKYCRequest] failed to fetch KYC request: %v", err)
		localization.SendErrorByCodeResponse(w, err.Error())
		return
	}

	res.ServerTime = time.Now().UTC().Format(time.RFC3339)
	localization.SendSuccessResponse(w, localization.SuccessDataRetrieved, res)
}

func (c *customerKYCAdapter) ApproveKycRequest(w http.ResponseWriter, r *http.Request) {
	ctx, span := util.TraceLogger(r.Context(), "handler", "ApproveKycRequest", "handler", "customer_kyc")
	defer span.End()
	log := util.LoggerFromCtx(ctx, c.logger)

	md := &types.ContextMetadata{}
	ctx = context.WithValue(ctx, constants.ContextKeyMetadata, md)
	localization.UpdateWriterContext(w, ctx)

	id, err := util.ExtractID(w, r)
	if err != nil {
		log.Errorf("[ApproveKycRequest] extractID: %v", err)
		return
	}

	if err := c.svc.EnableOrDisable(ctx, id, "", true); err != nil {
		w = util.HandlePendingResponseError(ctx, w, err)
		span.RecordError(err)
		log.Errorf("[ApproveKycRequest] service error: %v", err)
		localization.SendErrorByCodeResponse(w, err.Error())
		return
	}

	w = localization.ApplyActionCodeHeaderFromWriter(w, ctx)
	if md.IsMakerOnly {
		userCode, _ := ctx.Value(constants.ContextKey("user_code")).(string)
		log.Infof("[ApproveKycRequest] request submitted for user_code: %s is_maker_only: %v", userCode, md.IsMakerOnly)
		localization.SendSuccessResponse(w, localization.CustomerKycApprovalRequestSubmittedSuccessfully, nil)
		return
	}

	localization.SendSuccessResponse(w, localization.CustomerKycRequestApprovedSuccessfully, nil)
}

func (c *customerKYCAdapter) RejectKycRequest(w http.ResponseWriter, r *http.Request) {
	ctx, span := util.TraceLogger(r.Context(), "handler", "RejectKycRequest", "handler", "customer_kyc")
	defer span.End()
	log := util.LoggerFromCtx(ctx, c.logger)

	md := &types.ContextMetadata{}
	ctx = context.WithValue(ctx, constants.ContextKeyMetadata, md)
	localization.UpdateWriterContext(w, ctx)

	var req kyc_dto.ReasonRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		log.Errorf("[RejectKycRequest] failed to decode request body: %v", err)
		localization.SendErrorByCodeResponse(w, "invalid request format")
		return
	}

	if err := req.Validate(); err != nil {
		log.Errorf("[RejectKycRequest] validation failed: %v", err)
		localization.SendBadRequestResponse(w, err.Error())
		return
	}

	id, err := util.ExtractID(w, r)
	if err != nil {
		log.Errorf("[RejectKycRequest] extractID: %v", err)
		return
	}

	if err := c.svc.EnableOrDisable(ctx, id, req.Reason, false); err != nil {
		w = util.HandlePendingResponseError(ctx, w, err)
		span.RecordError(err)
		log.Errorf("[RejectKycRequest] service error: %v", err)
		localization.SendErrorByCodeResponse(w, err.Error())
		return
	}

	w = localization.ApplyActionCodeHeaderFromWriter(w, ctx)
	if md.IsMakerOnly {
		userCode, _ := ctx.Value(constants.ContextKey("user_code")).(string)
		log.Infof("[RejectKycRequest] request submitted for user_code: %s is_maker_only: %v", userCode, md.IsMakerOnly)
		localization.SendSuccessResponse(w, localization.CustomerKycRejectRequestSubmittedSuccessfully, nil)
		return
	}

	localization.SendSuccessResponse(w, localization.CustomerKycRequestRejectedSuccessfully, nil)
}

func (c *customerKYCAdapter) StartKycReview(w http.ResponseWriter, r *http.Request) {
	ctx, span := util.TraceLogger(r.Context(), "handler", "StartKycReview", "handler", "StartKycReview")
	defer span.End()
	log := util.LoggerFromCtx(ctx, c.logger)

	id := chi.URLParam(r, "kycID")
	if id == "" {
		localization.SendBadRequestResponse(w, localization.ErrorIdNotSetOnQueryParam.Code)
		return
	}

	newReview, err := c.svc.StartKycReview(ctx, id)
	if err != nil {
		w = util.HandlePendingResponseError(ctx, w, err)
		log.Errorf("[StartKycReview] failed to start KYC review: %v", err)
		localization.SendErrorByCodeResponse(w, err.Error())
		return
	}

	localization.SendSuccessResponse(w, localization.SuccessOperationCompleted, newReview)
}

func (c *customerKYCAdapter) PickKycReview(w http.ResponseWriter, r *http.Request) {
	ctx, span := util.TraceLogger(r.Context(), "handler", "PickKycReview", "handler", "PickKycReview")
	defer span.End()
	log := util.LoggerFromCtx(ctx, c.logger)

	id := chi.URLParam(r, "kycID")
	if id == "" {
		localization.SendBadRequestResponse(w, localization.ErrorIdNotSetOnQueryParam.Code)
		return
	}

	var req kyc_dto.ReasonRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		log.Errorf("[PickKycReview] failed to decode request body: %v", err)
		localization.SendErrorByCodeResponse(w, "invalid request format")
		return
	}

	err := c.svc.PickKycReview(ctx, id, req.Reason)
	if err != nil {
		w = util.HandlePendingResponseError(ctx, w, err)
		log.Errorf("[PickKycReview] failed to pick KYC review: %v", err)
		localization.SendErrorByCodeResponse(w, err.Error())
		return
	}

	localization.SendSuccessResponse(w, localization.SuccessOperationCompleted, nil)
}

func (c *customerKYCAdapter) ExportKYCOnboarding(w http.ResponseWriter, r *http.Request) {
	ctx, span := util.TraceLogger(r.Context(), "handler", "PickSelfActivateKycReview", "handler", "self_activate_kyc")
	defer span.End()
	log := util.LoggerFromCtx(ctx, c.logger)

	fileType := strings.ToLower(strings.TrimSpace(r.URL.Query().Get("file_type")))
	startDateRaw := strings.TrimSpace(r.URL.Query().Get("created_at_from"))
	endDateRaw := strings.TrimSpace(r.URL.Query().Get("created_at_to"))
	customerName := strings.TrimSpace(r.URL.Query().Get("name"))

	if fileType == "" || startDateRaw == "" || endDateRaw == "" {
		log.Warnf("[ExportKYCOnboarding] missing required params: file_type=%q, created_at_from=%q, created_at_to=%q", fileType, startDateRaw, endDateRaw)
		localization.SendBadRequestResponse(w, localization.ErrorRequiredFieldMissing.Message)
		return
	}

	from, to, err := util.FormatDateRangeToUTCStrings(startDateRaw, endDateRaw)
	if err != nil {
		log.Warnf("[ExportKYCOnboarding] invalid date format: created_at_from=%s, created_at_to=%s, err=%v", startDateRaw, endDateRaw, err)
		localization.SendErrorResponse(w, localization.ErrorInvalidDateFormat, nil, nil)
		return
	}

	if to.Before(from) {
		log.Warnf("[ExportKYCOnboarding] invalid date range: start=%v, end=%v", from, to)
		localization.SendErrorResponse(w, localization.ErrorInvalidDateFormat, nil, nil)
		return
	}

	fileLink, err := c.svc.ExportKYCOnboarding(ctx, from, to, fileType, customerName)
	if err != nil {
		log.Errorf("[ExportKYCOnboarding] service error: %v", err)
		localization.SendErrorByCodeResponse(w, err.Error())
		return
	}

	localization.SendSuccessResponse(w, localization.SucccessKYCExportedSuccessfully, fileLink)
}
