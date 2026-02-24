package customer

import (
	dto "cbe-super-app-cps-action/internal/constants/dto/customer_kyc"
	inbound "cbe-super-app-cps-action/internal/constants/interfaces/customer_kyc"
	"cbe-super-app-cps-action/internal/constants/localization"
	"cbe-super-app-cps-action/internal/handlers/rest/http/customer_kyc/core"
	"cbe-super-app-cps-action/internal/service"
	util "cbe-super-app-cps-action/pkgs/utils"
	"encoding/json"
	"net/http"

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

// CreateCustomerKYC creates a new KYC request for a customer
//
//	@Summary		Create Customer KYC
//	@Description	Create a new KYC request for a customer. This will create a CPS action for approval.
//	@Tags			Customer KYC
//	@Accept			json
//	@Produce		json
//	@Param			request	body		customerkyc.CreateCustomerKYCRequest	true	"Create KYC Request"
//	@Success		201		{object}	localization.StandardResponse{data=nil}	"KYC request created successfully"
//	@Failure		400		{object}	localization.StandardResponse{data=nil}	"Bad request"
//	@Failure		500		{object}	localization.StandardResponse{data=nil}	"Internal server error"
//	@Security		BearerAuth
//	@Router			/customers/kyc [post]
func (c *customerKYCAdapter) CreateCustomerKYC(w http.ResponseWriter, r *http.Request) {
	ctx, span := util.TraceLogger(r.Context(), "handler", "CreateCustomerKYC", "handler", "customer_kyc")
	defer span.End()
	log := util.LoggerFromCtx(ctx, c.logger)

	req, err := core.ParseRequestFromMultipleFormData(r)
	if err != nil {
		span.RecordError(err)
		log.Errorf("[KycH][Create] parse form err: %v", err)
		localization.SendBadRequestResponse(w, err.Error())
		return
	}

	if err := req.Validate(); err != nil {
		log.Errorf("[CreateCustomerKYC Validation) validation failed: %v", err)
		localization.SendBadRequestResponse(w, err.Error())
		return
	}

	if err := c.svc.Create(ctx, req); err != nil {
		log.Errorf("[CreateCustomerKYC] service call failed: %v", err)
		localization.SendErrorByCodeResponse(w, err.Error())
		return
	}

	localization.SendSuccessResponse(w, localization.SuccessOperationCompleted, nil)
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

	filterParam := util.ExtractFilterParams(r)

	search := r.URL.Query().Get("search")
	filter := r.URL.Query().Get("filter")

	if err := util.NoSpecialChars(search); err != nil {
		localization.SendErrorByCodeResponse(w, err.Error())
		return
	}

	if err := util.NoSpecialChars(filter); err != nil {
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

	localization.SendSuccessResponse(w, localization.SuccessDataRetrieved, res)
}

// DeleteKYCRequest deletes a specific KYC request by ID
//
//	@Summary		Delete KYC Request
//	@Description	Deletes a specific KYC request by its ID. This will create a CPS action for approval.
//	@Tags			Customer KYC
//	@Accept			json
//	@Produce		json
//	@Param			id	path		string								true	"KYC Request ID"
//	@Success		200	{object}	localization.StandardResponse{data=nil}	"KYC request deletion initiated"
//	@Failure		400	{object}	localization.StandardResponse{data=nil}	"Bad request - ID required"
//	@Failure		500	{object}	localization.StandardResponse{data=nil}	"Internal server error"
//	@Security		BearerAuth
//	@Router			/customers/kyc/{id} [delete]
func (c *customerKYCAdapter) DeleteKYCRequest(w http.ResponseWriter, r *http.Request) {
	ctx, span := util.TraceLogger(r.Context(), "handler", "DeleteKYCRequest", "handler", "customer_kyc")
	defer span.End()
	log := util.LoggerFromCtx(ctx, c.logger)

	id := chi.URLParam(r, "id")
	if id == "" {
		localization.SendBadRequestResponse(w, localization.ErrorIdNotSetOnQueryParam.Code)
		return
	}

	if err := c.svc.Delete(ctx, id); err != nil {
		log.Errorf("[DeleteKYCRequest] failed to delete KYC request: %v", err)
		localization.SendErrorByCodeResponse(w, err.Error())
		return
	}

	localization.SendSuccessResponse(w, localization.SuccessOperationCompleted, nil)
}

// UpdateKYCStatus updates the KYC status for a specific request ID
//
//	@Summary		Update KYC Status
//	@Description	Updates the KYC status for a specific request ID. This will create a CPS action for approval.
//	@Tags			Customer KYC
//	@Accept			json
//	@Produce		json
//	@Param			id		path		string							true	"KYC Request ID"
//	@Param			request	body		customerkyc.UpdateKYCStatusRequest	true	"Update KYC Status Request"
//	@Success		200		{object}	localization.StandardResponse{data=nil}	"KYC status update initiated"
//	@Failure		400		{object}	localization.StandardResponse{data=nil}	"Bad request"
//	@Failure		500		{object}	localization.StandardResponse{data=nil}	"Internal server error"
//	@Security		BearerAuth
//	@Router			/customers/kyc/{id} [patch]
func (c *customerKYCAdapter) UpdateKYCStatus(w http.ResponseWriter, r *http.Request) {
	ctx, span := util.TraceLogger(r.Context(), "handler", "UpdateKYCStatus", "handler", "customer_kyc")
	defer span.End()
	log := util.LoggerFromCtx(ctx, c.logger)

	id := chi.URLParam(r, "id")
	if id == "" {
		localization.SendBadRequestResponse(w, localization.ErrorIdNotSetOnQueryParam.Code)
		return
	}

	var req dto.UpdateKYCStatusRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		span.RecordError(err)
		log.Errorf("[UpdateKYCStatus] failed to unmarshal request: %v", err)
		localization.SendBadRequestResponse(w, localization.ErrorInvalidJSONPayload.Code)
		return
	}

	if err := req.Validate(); err != nil {
		log.Errorf("[UpdateKYCStatus Validation] validation failed: %v", err)
		localization.SendBadRequestResponse(w, err.Error())
		return
	}

	if err := c.svc.UpdateKYCStatus(ctx, id, req.KYCStatus); err != nil {
		log.Errorf("[UpdateKYCStatus] service call failed: %v", err)
		localization.SendErrorByCodeResponse(w, err.Error())
		return
	}

	localization.SendSuccessResponse(w, localization.SuccessOperationCompleted, nil)
}
