package customer

import (
	"github.com/abelshebela/cbe-super-app-cps-action/internal/constants"
	kyc_dto "github.com/abelshebela/cbe-super-app-cps-action/internal/constants/dto/customer_kyc"
	inbound "github.com/abelshebela/cbe-super-app-cps-action/internal/constants/interfaces/customer_kyc"
	"github.com/abelshebela/cbe-super-app-cps-action/internal/constants/localization"
	"github.com/abelshebela/cbe-super-app-cps-action/internal/constants/types"
	"github.com/abelshebela/cbe-super-app-cps-action/internal/service"
	util "github.com/abelshebela/cbe-super-app-cps-action/pkgs/utils"
	"context"
	"encoding/json"
	"net/http"
	"strings"
	"time"

	"github.com/go-chi/chi/v5"
	"gitlab.com/bersufekadgetachew/cbe-super-app-shared/shared/utils"
)

type selfActivationKYCAdapter struct {
	svc    service.SelfActivateKYCService
	logger utils.Logger
}

func NewSelfActivationAdapter(svc service.SelfActivateKYCService, logger utils.Logger) inbound.SelfActivationKyc {
	return &selfActivationKYCAdapter{
		svc:    svc,
		logger: logger,
	}
}

func (c *selfActivationKYCAdapter) GetAllSelfActivateKYCRequests(w http.ResponseWriter, r *http.Request) {
	ctx, span := util.TraceLogger(r.Context(), "handler", "GetAllSelfActivateKYCRequests", "handler", "self_activate_kyc")
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
		log.Errorf("[GetAllSelfActivateKYCRequests] failed to fetch KYC requests: %v", err)
		localization.SendErrorByCodeResponse(w, err.Error())
		return
	}

	localization.SendSuccessResponse(w, localization.SuccessDataRetrieved, res)
}

func (c *selfActivationKYCAdapter) GetSelfActivateKYCRequest(w http.ResponseWriter, r *http.Request) {
	ctx, span := util.TraceLogger(r.Context(), "handler", "GetSelfActivateKYCRequest", "handler", "self_activate_kyc")
	defer span.End()
	log := util.LoggerFromCtx(ctx, c.logger)

	id := chi.URLParam(r, "id")
	if id == "" {
		localization.SendBadRequestResponse(w, localization.ErrorIdNotSetOnQueryParam.Code)
		return
	}

	res, err := c.svc.FindByID(ctx, id)
	if err != nil {
		log.Errorf("[GetSelfActivateKYCRequest] failed to fetch KYC request: %v", err)
		localization.SendErrorByCodeResponse(w, err.Error())
		return
	}

	res.ServerTime = time.Now().UTC().Format(time.RFC3339)
	localization.SendSuccessResponse(w, localization.SuccessDataRetrieved, res)
}

func (c *selfActivationKYCAdapter) ApproveSelfActivateKycRequest(w http.ResponseWriter, r *http.Request) {
	ctx, span := util.TraceLogger(r.Context(), "handler", "ApproveSelfActivateKycRequest", "handler", "self_activate_kyc")
	defer span.End()
	log := util.LoggerFromCtx(ctx, c.logger)

	md := &types.ContextMetadata{}
	ctx = context.WithValue(ctx, constants.ContextKeyMetadata, md)
	localization.UpdateWriterContext(w, ctx)

	id, err := util.ExtractID(w, r)
	if err != nil {
		log.Errorf("[ApproveSelfActivateKycRequest] extractID: %v", err)
		return
	}

	if err := c.svc.EnableOrDisable(ctx, id, "", true); err != nil {
		w = util.HandlePendingResponseError(ctx, w, err)
		span.RecordError(err)
		log.Errorf("[ApproveSelfActivateKycRequest] service error: %v", err)
		localization.SendErrorByCodeResponse(w, err.Error())
		return
	}

	w = localization.ApplyActionCodeHeaderFromWriter(w, ctx)
	if md.IsMakerOnly {
		userCode, _ := ctx.Value(constants.ContextKey("user_code")).(string)
		log.Infof("[ApproveSelfActivateKycRequest] request submitted for user_code: %s is_maker_only: %v", userCode, md.IsMakerOnly)
		localization.SendSuccessResponse(w, localization.CustomerKycRequestApprovedSuccessfully, nil)
		return
	}

	localization.SendSuccessResponse(w, localization.CustomerKycApprovalRequestSubmittedSuccessfully, nil)
}

func (c *selfActivationKYCAdapter) RejectSelfActivateKycRequest(w http.ResponseWriter, r *http.Request) {
	ctx, span := util.TraceLogger(r.Context(), "handler", "RejectSelfActivateKycRequest", "handler", "self_activate_kyc")
	defer span.End()
	log := util.LoggerFromCtx(ctx, c.logger)

	md := &types.ContextMetadata{}
	ctx = context.WithValue(ctx, constants.ContextKeyMetadata, md)
	localization.UpdateWriterContext(w, ctx)

	var req kyc_dto.ReasonRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		log.Errorf("[RejectSelfActivateKycRequest] failed to decode request body: %v", err)
		localization.SendErrorByCodeResponse(w, "invalid request format")
		return
	}

	if err := req.Validate(); err != nil {
		log.Errorf("[RejectSelfActivateKycRequest] validation failed: %v", err)
		localization.SendBadRequestResponse(w, err.Error())
		return
	}

	id, err := util.ExtractID(w, r)
	if err != nil {
		log.Errorf("[RejectSelfActivateKycRequest] extractID: %v", err)
		return
	}

	if err := c.svc.EnableOrDisable(ctx, id, req.Reason, false); err != nil {
		w = util.HandlePendingResponseError(ctx, w, err)
		span.RecordError(err)
		log.Errorf("[RejectSelfActivateKycRequest] service error: %v", err)
		localization.SendErrorByCodeResponse(w, err.Error())
		return
	}

	w = localization.ApplyActionCodeHeaderFromWriter(w, ctx)
	if md.IsMakerOnly {
		userCode, _ := ctx.Value(constants.ContextKey("user_code")).(string)
		log.Infof("[RejectSelfActivateKycRequest] request submitted for user_code: %s is_maker_only: %v", userCode, md.IsMakerOnly)
		localization.SendSuccessResponse(w, localization.CustomerKycRequestRejectedSuccessfully, nil)
		return
	}

	localization.SendSuccessResponse(w, localization.CustomerKycRejectRequestSubmittedSuccessfully, nil)
}

func (c *selfActivationKYCAdapter) StartSelfActivateKycReview(w http.ResponseWriter, r *http.Request) {
	ctx, span := util.TraceLogger(r.Context(), "handler", "StartSelfActivateKycReview", "handler", "self_activate_kyc")
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
		log.Errorf("[StartSelfActivateKycReview] failed to start KYC review: %v", err)
		localization.SendErrorByCodeResponse(w, err.Error())
		return
	}

	newReview.ServerTime = time.Now().UTC().Format(time.RFC3339)
	localization.SendSuccessResponse(w, localization.SuccessKYCReviewStarted, newReview)
}

func (c *selfActivationKYCAdapter) PickSelfActivateKycReview(w http.ResponseWriter, r *http.Request) {
	ctx, span := util.TraceLogger(r.Context(), "handler", "PickSelfActivateKycReview", "handler", "self_activate_kyc")
	defer span.End()
	log := util.LoggerFromCtx(ctx, c.logger)

	id := chi.URLParam(r, "kycID")
	if id == "" {
		localization.SendBadRequestResponse(w, localization.ErrorIdNotSetOnQueryParam.Code)
		return
	}

	var req kyc_dto.ReasonRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		log.Errorf("[PickSelfActivateKycReview] failed to decode request body: %v", err)
		localization.SendErrorByCodeResponse(w, "invalid request format")
		return
	}

	err := c.svc.PickKycReview(ctx, id, req.Reason)
	if err != nil {
		w = util.HandlePendingResponseError(ctx, w, err)
		log.Errorf("[PickSelfActivateKycReview] failed to pick KYC review: %v", err)
		localization.SendErrorByCodeResponse(w, err.Error())
		return
	}

	localization.SendSuccessResponse(w, localization.SuccessKYCReviewPicked, nil)
}

func (c *selfActivationKYCAdapter) ExportUserSelfActivation(w http.ResponseWriter, r *http.Request) {
	ctx, span := util.TraceLogger(r.Context(), "handler", "PickSelfActivateKycReview", "handler", "self_activate_kyc")
	defer span.End()
	log := util.LoggerFromCtx(ctx, c.logger)

	fileType := strings.ToLower(strings.TrimSpace(r.URL.Query().Get("file_type")))
	startDateRaw := strings.TrimSpace(r.URL.Query().Get("created_at_from"))
	endDateRaw := strings.TrimSpace(r.URL.Query().Get("created_at_to"))
	customerName := strings.TrimSpace(r.URL.Query().Get("name"))
	status := strings.TrimSpace(r.URL.Query().Get("kyc_status"))

	if fileType == "" || startDateRaw == "" || endDateRaw == "" || status == "" {
		log.Warnf("[ExportUserSelfActivation] missing required params: file_type=%q, created_at_from=%q, created_at_to=%q", fileType, startDateRaw, endDateRaw)
		localization.SendBadRequestResponse(w, localization.ErrorRequiredFieldMissing.Message)
		return
	}

	validStatus := map[string]bool{
		"PENDING,IN_REVIEW": true,
		"APPROVED":          true,
		"REJECTED":          true,
	}
	if !validStatus[status] {
		localization.SendBadRequestResponse(w, "status should be either PENDING,IN_REVIEW, APPROVED or REJECTED")
		return
	}

	from, to, err := util.FormatDateRangeToUTCStrings(startDateRaw, endDateRaw)
	if err != nil {
		log.Warnf("[ExportUserSelfActivation] invalid date format: created_at_from=%s, created_at_to=%s, err=%v", startDateRaw, endDateRaw, err)
		localization.SendErrorResponse(w, localization.ErrorInvalidDateFormat, nil, nil)
		return
	}

	if to.Before(from) {
		log.Warnf("[ExportUserSelfActivation] invalid date range: start=%v, end=%v", from, to)
		localization.SendErrorResponse(w, localization.ErrorInvalidDateFormat, nil, nil)
		return
	}

	fileLink, err := c.svc.ExportUserSelfActivation(ctx, from, to, fileType, status, customerName)
	if err != nil {
		log.Errorf("[ExportUserSelfActivation] service error: %v", err)
		localization.SendErrorByCodeResponse(w, err.Error())
		return
	}

	localization.SendSuccessResponse(w, localization.SucccessKYCExportedSuccessfully, fileLink)
}

func (c *selfActivationKYCAdapter) GetUsersActionLog(w http.ResponseWriter, r *http.Request) {
	ctx, span := util.TraceLogger(r.Context(), "handler", "GetUsersActionLog", "handler", "self_activate_kyc")
	defer span.End()
	log := util.LoggerFromCtx(ctx, c.logger)

	customerNumber := r.URL.Query().Get("number")
	if customerNumber == "" {
		localization.SendBadRequestResponse(w, localization.ErrorIdNotSetOnQueryParam.Code)
		return
	}

	res, err := c.svc.GetUsersActionLog(ctx, customerNumber)
	if err != nil {
		log.Errorf("[GetUsersActionLog] failed to fetch KYC request: %v", err)
		localization.SendErrorByCodeResponse(w, err.Error())
		return
	}

	localization.SendSuccessResponse(w, localization.SuccessActionLogRetrieved, res)
}
