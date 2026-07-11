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

	localization.SendSuccessResponse(w, localization.SuccessOperationCompleted, newReview)
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

	localization.SendSuccessResponse(w, localization.SuccessOperationCompleted, nil)
}
