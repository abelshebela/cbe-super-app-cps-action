package customer

import (
	dto "cbe-super-app-cps-action/internal/constants/dto/customer"
	"cbe-super-app-cps-action/internal/constants/interfaces/customer"
	"cbe-super-app-cps-action/internal/constants/localization"
	"cbe-super-app-cps-action/internal/constants/types"
	"cbe-super-app-cps-action/internal/handlers/rest/http/customer/core"
	"encoding/json"

	"cbe-super-app-cps-action/internal/service"
	util "cbe-super-app-cps-action/pkgs/utils"
	"net/http"

	member "gitlab.com/bersufekadgetachew/cbe-super-app-shared/shared/entities/member"

	"github.com/go-chi/chi/v5"
	utils "gitlab.com/bersufekadgetachew/cbe-super-app-shared/shared/utils"
	"go.opentelemetry.io/otel/attribute"
)

type customer_resp *member.User
type customers_paginated_resp *types.PaginatedResponse[[]*member.User]

type customerAdapter struct {
	customerService service.CustomerService
	logger          utils.Logger
}

func InitCustomerAdapter(customer service.CustomerService, logger utils.Logger) customer.CustomerDetail {
	return &customerAdapter{
		logger:          logger,
		customerService: customer,
	}
}

// SetEnableCustomerSession initiates enabling a customer session by generating an OTP
//
//	@Summary		Initiate enable customer session
//	@Description	Initiates enabling a customer session by generating an OTP for the customer
//	@Tags			Customers
//	@Accept			json
//	@Produce		json
//	@Param			id	path		string																		true	"Customer ID"
//	@Success		200	{object}	localization.StandardResponse{data=customer.CustomerEnableSessionResponse}	"OTP generated successfully"
//	@Failure		400	{object}	localization.StandardResponse{data=nil}										"Bad request - Customer ID required"
//	@Failure		500	{object}	localization.StandardResponse{data=nil}										"Internal server error"
//	@Security		BearerAuth
//	@Router			/customers/{id}/enable-session [post]
func (c *customerAdapter) SetEnableCustomerSession(w http.ResponseWriter, r *http.Request) {
	ctx, span := util.TraceLogger(r.Context(), "handler", "setEnableCustomerSession", "handler", "customer")
	defer span.End()
	id := chi.URLParam(r, "id")
	if id == "" {
		c.logger.Errorf("id not set on param")
		localization.SendBadRequestResponse(w, localization.ErrorIdNotSetOnQueryParam.Code)
		return
	}

	span.SetAttributes(attribute.String("customer.id", id))
	otp, err := c.customerService.CreateEnableCustomerSession(ctx, id)
	if err != nil {
		span.RecordError(err)
		c.logger.Errorf("[SetEnableCustomerSession] service error: %v", err)
		localization.SendErrorByCodeResponse(w, err.Error())
		return
	}
	c.logger.Infof("[SetEnableCustomerSession] OTP generated successfully for customer id: %s", id)
	localization.SendSuccessResponse(w, localization.CustomerEnableRequestSessionCreatedSuccessfully, dto.CustomerEnableSessionResponse{
		Otp: otp,
	})
}

// DisableCustomer disables a customer by ID
//
//	@Summary		Disable customer
//	@Description	Disables a customer by their ID
//	@Tags			Customers
//
// @Accept			json
//
//	@Produce		json
//	@Param			id	path		string									true	"Customer ID"
//	@Param			body	body	dto.CustomerDisableDTO	true	" body (fields: is_temporary, disable_reason)"
//	@Success		200	{object}	localization.StandardResponse{data=nil}	"Customer disabled successfully"
//	@Failure		400	{object}	localization.StandardResponse{data=nil}	"Bad request - Customer ID required"
//	@Failure		500	{object}	localization.StandardResponse{data=nil}	"Internal server error"
//	@Security		BearerAuth
//	@Router			/customers/disable/{id} [patch]
func (c *customerAdapter) DisableCustomer(w http.ResponseWriter, r *http.Request) {
	ctx, span := util.TraceLogger(r.Context(), "handler", "disableCustomer", "handler", "customer")
	defer span.End()
	id := chi.URLParam(r, "id")
	if id == "" {
		c.logger.Errorf("id not set on param")
		localization.SendBadRequestResponse(w, localization.ErrorIdNotSetOnQueryParam.Code)
		return
	}

	var payload dto.CustomerDisableDTO
	if err := json.NewDecoder(r.Body).Decode(&payload); err != nil {
		span.RecordError(err)
		localization.SendErrorByCodeResponse(w, localization.ErrorInvalidRequestBody.Code)
		return
	}

	if payload.IsTemporary == nil {
		localization.SendErrorByCodeResponse(w, localization.ErrorInvalidInputParameter.Code)
		return
	}

	if err := core.ValidateString(payload.DisableReason); err != nil {
		c.logger.Errorf("invalid disable reason: %v", err)
		localization.SendErrorByCodeResponse(w, localization.ErrorInvalidInputParameter.Code)
		return
	}

	span.SetAttributes(attribute.String("customer.id", id))
	err := c.customerService.DisableCustomerByID(ctx, id, payload)
	if err != nil {
		span.RecordError(err)
		c.logger.Errorf("[DisableCustomer] service error: %v", err)
		localization.SendErrorByCodeResponse(w, err.Error())
		return
	}
	c.logger.Infof("[DisableCustomer] request sent successfully for customer id: %s", id)
	localization.SendSuccessResponse(w, localization.CustomerDisableRequestCreatedSuccessfully, nil)
}

// EnableCustomer enables a customer by ID using OTP
//
//	@Summary		Enable customer
//	@Description	Enables a customer by their ID using OTP
//	@Tags			Customers
//	@Accept			json
//	@Produce		json
//	@Param			id		path		string									true	"Customer ID"
//	@Param			body	body	dto.CustomerDisableDTO	true	" body (fields: is_temporary, disable_reason)"
//	@Param			body	body		customer.CustomerEnableDTO				true	"Enable customer payload"
//	@Success		200		{object}	localization.StandardResponse{data=nil}	"Customer enabled successfully"
//	@Failure		400		{object}	localization.StandardResponse{data=nil}	"Bad request - Customer ID required or invalid body"
//	@Failure		500		{object}	localization.StandardResponse{data=nil}	"Internal server error"
//	@Security		BearerAuth
//	@Router			/customers/enable/{id} [patch]
func (c *customerAdapter) EnableCustomer(w http.ResponseWriter, r *http.Request) {
	ctx, span := util.TraceLogger(r.Context(), "handler", "enableCustomer", "handler", "customer")
	defer span.End()
	id := chi.URLParam(r, "id")
	if id == "" {
		c.logger.Errorf("id not set on param")
		localization.SendBadRequestResponse(w, localization.ErrorIdNotSetOnQueryParam.Code)
		return
	}

	var payload dto.CustomerEnableDTO
	if err := json.NewDecoder(r.Body).Decode(&payload); err != nil {
		span.RecordError(err)
		localization.SendErrorByCodeResponse(w, localization.ErrorInvalidRequestBody.Code)
		return
	}

	span.SetAttributes(attribute.String("customer.id", id))
	err := c.customerService.EnableCustomerByID(ctx, id, payload.UserOTP)
	if err != nil {
		span.RecordError(err)
		c.logger.Errorf("[EnableCustomer] service error: %v", err)
		localization.SendErrorByCodeResponse(w, err.Error())
		return
	}
	c.logger.Infof("[EnableCustomer] request sent successfully for customer id: %s", id)
	localization.SendSuccessResponse(w, localization.CustomerEnableRequestCreatedSuccessfully, nil)
}

// GetCustomerDetail
//
//	@Summary		Get Customer Detail
//	@Description	Retrieve customer details with pagination, filtering, and search. Searchable fields: full_name, phone_number, gender, user_name, user_code, is_blocked, kyc_level.
//	@Tags			Customers
//	@Accept			json
//	@Produce		json
//	@Param			page			query	int		false	"Page number"
//	@Param			per_page		query	int		false	"Items per page"
//	@Param			search			query	string	false	"Search term (searches full_name, phone_number, gender, user_name, user_code, is_blocked, kyc_level)"
//	@Param			gender			query	string	false	"Filter by gender"
//	@Param			branch_code		query	string	false	"Filter by branch code"
//	@Param			kyc_level		query	int	false	"Filter by KYC level"
//	@Param			is_blocked		query	bool	false	"Filter by blocked status"
//	@Param			enabled			query	bool	false	"Filter by enabled status"
//	@Param			bps_reject_status	query	string	false	"Filter by BPS reject status"
//	@Success		200	{object}	localization.StandardResponse{data=[]model.User}	"Customer details retrieved successfully"
//	@Failure		400	{object}	localization.StandardResponse{data=nil}	"Bad request"
//	@Failure		500	{object}	localization.StandardResponse{data=nil}	"Internal server error"
//	@Security		BearerAuth
//	@Router			/customers [get]
func (c customerAdapter) GetCustomerDetail(w http.ResponseWriter, r *http.Request) {
	ctx, span := util.TraceLogger(r.Context(), "handler", "getCustomerDetail", "handler", "customer")
	defer span.End()
	filterParams := util.ExtractFilterParams(r)
	customers, err := c.customerService.GetCustomersDetail(ctx, filterParams)
	if err != nil {
		span.RecordError(err)
		c.logger.Errorf("[GetCustomerDetail] service error: %v", err)
		localization.SendErrorByCodeResponse(w, localization.ErrorFailedToGetCustomerDetail.Code)
		return
	}

	span.SetAttributes(attribute.Int("customer.count", len(customers.Data)))
	c.logger.Infof("[GetCustomerDetail] retrieved %d customers", len(customers.Data))
	localization.SendSuccessResponse(w, localization.SuccessCustomerDetailSuccessfullyFetched, customers)
}

// GetCustomerByID retrieves a specific customer by ID
//
//	@Summary		Get customer by ID
//	@Description	Retrieves detailed information for a specific customer by their ID
//	@Tags			Customers
//	@Accept			json
//	@Produce		json
//	@Param			id	path		string												true	"Customer ID"
//	@Success		200	{object}	localization.StandardResponse{data=customer_resp}	"Customer details retrieved successfully"
//	@Failure		400	{object}	localization.StandardResponse{data=nil}				"Bad request - Customer ID required"
//	@Failure		404	{object}	localization.StandardResponse{data=nil}				"Customer not found"
//	@Failure		500	{object}	localization.StandardResponse{data=nil}				"Internal server error"
//	@Security		BearerAuth
//	@Router			/customers/{id} [get]
func (c customerAdapter) GetCustomerByID(w http.ResponseWriter, r *http.Request) {
	ctx, span := util.TraceLogger(r.Context(), "handler", "getCustomerById", "handler", "customer")
	defer span.End()
	id := chi.URLParam(r, "id")
	if id == "" {
		c.logger.Errorf("id not set on param")
		localization.SendBadRequestResponse(w, localization.ErrorIdNotSetOnQueryParam.Code)
		return
	}

	span.SetAttributes(attribute.String("customer.id", id))
	userDetail, err := c.customerService.GetCustomerByID(ctx, id)
	if err != nil {
		span.RecordError(err)
		c.logger.Errorf("[GetCustomerByID] service error: %v", err)
		localization.SendErrorByCodeResponse(w, localization.UserNotFoundWithGivenID.Code)
		return
	}
	c.logger.Infof("[GetCustomerByID] customer retrieved successfully for id: %s", id)
	localization.SendSuccessResponse(w, localization.CustomerDetailSuccessfullyFetched, userDetail)

}

// GetBlockedCustomer retrieves blocked customers
//
//	@Summary		Get blocked customers
//	@Description	Retrieves a paginated list of blocked customers
//	@Tags			Customers
//	@Accept			json
//	@Produce		json
//	@Param			page		query		int																false	"Page number"		default(1)
//	@Param			per_page	query		int																false	"Items per page"	default(10)
//	@Param			search			query	string	false	"Search term (searches full_name, phone_number, gender, user_name, user_code, is_blocked, kyc_level)"
//	@Param			gender			query	string	false	"Filter by gender"
//	@Param			branch_code		query	string	false	"Filter by branch code"
//	@Param			kyc_level		query	int	false	"Filter by KYC level"
//	@Param			enabled			query	bool	false	"Filter by enabled status"
//	@Param			bps_reject_status	query	string	false	"Filter by BPS reject status"
//	@Success		200			{object}	localization.StandardResponse{data=customers_paginated_resp}	"Blocked customers retrieved successfully"
//	@Failure		400			{object}	localization.StandardResponse{data=nil}							"Bad request"
//	@Failure		500			{object}	localization.StandardResponse{data=nil}							"Internal server error"
//	@Security		BearerAuth
//	@Router			/customers/blocked [get]
func (c customerAdapter) GetBlockedCustomer(w http.ResponseWriter, r *http.Request) {
	ctx, span := util.TraceLogger(r.Context(), "handler", "getBlockedCustomer", "handler", "customer")
	defer span.End()
	filterParams := util.ExtractFilterParams(r)

	BlockedCustomer, err := c.customerService.GetBlockedCustomer(ctx, filterParams)
	if err != nil {
		span.RecordError(err)
		c.logger.Errorf("[GetBlockedCustomer] service error: %v", err)
		localization.SendErrorByCodeResponse(w, localization.ErrorFailedToGetBlockedCustomer.Code)
		return
	}

	span.SetAttributes(attribute.Int("customer.blocked.count", len(BlockedCustomer.Data)))
	c.logger.Infof("[GetBlockedCustomer] retrieved %d blocked customers", len(BlockedCustomer.Data))
	localization.SendSuccessResponse(w, localization.SuccessFullyFetchBlockCustomer, BlockedCustomer)
}

// GetLinkedAccount retrieves customer linked account
//
//	@Summary		Get Customer Linked Account
//	@Description	Retrives a list of customer linked account
//	@Tags			Customers
//	@Param			customer_number	path	string	true	"Customer number"
//	@Produce		json
//	@Success		200	{object}	localization.StandardResponse{data=customers_paginated_resp}	"Blocked customers retrieved successfully"
//	@Failure		400	{object}	localization.StandardResponse{data=nil}							"Bad request"
//	@Failure		500	{object}	localization.StandardResponse{data=nil}							"Internal server error"
//	@Security		BearerAuth
//	@Router			/customers/linked_account/{customer_number} [get]
func (c customerAdapter) GetLinkedAccount(w http.ResponseWriter, r *http.Request) {
	ctx, span := util.TraceLogger(r.Context(), "handler", "getLinkedAccount", "handler", "customer")
	defer span.End()
	id := chi.URLParam(r, "customer_number")
	if id == "" {
		c.logger.Errorf("customer_number not set on param")
		localization.SendBadRequestResponse(w, localization.ErrorIdNotSetOnQueryParam.Code)
		return
	}

	span.SetAttributes(attribute.String("customer.number", id))
	userDetail, err := c.customerService.GetLinkedAccount(ctx, id)
	if err != nil {
		span.RecordError(err)
		c.logger.Errorf("[GetLinkedAccount] service error: %v", err)
		localization.SendErrorByCodeResponse(w, localization.UserNotFoundWithGivenID.Code)
		return
	}
	c.logger.Infof("[GetLinkedAccount] linked accounts retrieved successfully for customer_number: %s", id)
	localization.SendSuccessResponse(w, localization.CustomerDetailSuccessfullyFetched, userDetail)
}

// ApproveFaydaCustomer
//
//	@Summary		Approve Fayda Customer
//	@Description	Approves a customer's Fayda application by updating their Fayda risk level
//	@Tags			Customers
//	@Accept			json
//	@Produce		json
//	@Param			id		path		string									true	"Customer ID"
//	@Param			body	body		customer.FaydaApproveRequest			true	"Fayda approval payload"
//	@Success		200		{object}	localization.StandardResponse{data=nil}	"Fayda customer approval request sent successfully"
//	@Failure		400		{object}	localization.StandardResponse{data=nil}	"Bad request - Customer ID required or invalid body"
//	@Failure		500		{object}	localization.StandardResponse{data=nil}	"Internal server error"
//	@Security		BearerAuth
//	@Router			/customers/fayda/enable/{id} [post]
func (c customerAdapter) ApproveFaydaCustomer(w http.ResponseWriter, r *http.Request) {
	ctx, span := util.TraceLogger(r.Context(), "handler", "approveFaydaCustomer", "handler", "customer")
	defer span.End()
	var req dto.FaydaApproveRequest
	id := chi.URLParam(r, "id")
	if id == "" {
		c.logger.Errorf("id not set on param")
		localization.SendBadRequestResponse(w, localization.ErrorInvalidAction.Message)
		return
	}

	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		span.RecordError(err)
		c.logger.Errorf(localization.ErrorInvalidJSONPayload.Message)
		localization.SendBadRequestResponse(w, localization.MsgInvalidJSONPayload)
		return
	}

	if req.Validate() != nil {
		c.logger.Errorf("")
	}
	span.SetAttributes(attribute.String("customer.id", id))
	if err := c.customerService.ApproveFaydaCustomer(ctx, id, req); err != nil {
		span.RecordError(err)
		c.logger.Errorf("[ApproveFaydaCustomer] service error: %v", err)
		localization.SendErrorByCodeResponse(w, err.Error())
		return
	}

	c.logger.Infof("[ApproveFaydaCustomer] request sent successfully for customer id: %s", id)
	localization.SendSuccessResponse(w, localization.FaydaCustomerApprovalRequestSent, nil)
}

func (c customerAdapter) SearchCustomerByCIForAccountNumber(w http.ResponseWriter, r *http.Request) {
	ctx, span := util.TraceLogger(r.Context(), "handler", "searchCustomerByCIForAccountNumber", "handler", "customer")
	defer span.End()

	var req dto.SearchCustomerByCIRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		span.RecordError(err)
		localization.SendErrorByCodeResponse(w, localization.ErrorInvalidRequestBody.Code)
		return
	}

	if err := req.Validate(); err != nil {
		c.logger.Errorf("invalid search customer by ci request")
		localization.SendErrorByCodeResponse(w, err.Error())
		return
	}

	customer, err := c.customerService.SearchCustomerByCIForAccountNumber(ctx, req)
	if err != nil {
		span.RecordError(err)
		c.logger.Errorf("[SearchCustomerByCIForAccountNumber] service error: %v", err)
		localization.SendErrorByCodeResponse(w, err.Error())
		return
	}
	localization.SendSuccessResponse(w, localization.SuccessCustomerDetailSuccessfullyFetched, customer)

}
