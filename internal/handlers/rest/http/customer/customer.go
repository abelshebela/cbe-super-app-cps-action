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

type customers_paginated_resp *types.PaginatedResponse[[]*member.User]

type customerAdapter struct {
	customerService service.CustomerService
	logger          utils.Logger
}

// GetCustomerActionLogByID retrieves action logs for a specific customer
//
//	@Summary		Get customer action log
//	@Description	Retrieves a paginated list of action logs for a specific customer by their ID. Includes pagination support.
//	@Tags			Customers
//	@Accept			json
//	@Produce		json
//	@Param			id		path		string												true	"Customer ID"
//	@Param			page	query		int													false	"Page number"
//	@Param			per_page	query	int												false	"Items per page"
//	@Success		200		{object}	localization.StandardResponse{data=object}		"Customer action logs retrieved successfully"
//	@Failure		400		{object}	localization.StandardResponse{data=nil}				"Bad request - Customer ID required"
//	@Failure		404		{object}	localization.StandardResponse{data=nil}				"Customer not found"
//	@Failure		500		{object}	localization.StandardResponse{data=nil}				"Internal server error"
//	@Security		BearerAuth
//	@Router			/customers/action_log/{id} [get]
func (c *customerAdapter) GetCustomerActionLogByID(w http.ResponseWriter, r *http.Request) {
	ctx, span := util.TraceLogger(r.Context(), "handler", "getCustomerActionLogByID", "handler", "customer")
	defer span.End()
	log := util.LoggerFromCtx(ctx, c.logger)
	id := chi.URLParam(r, "id")
	if id == "" {
		log.Errorf("[CustomerH] id not set")
		localization.SendBadRequestResponse(w, localization.ErrorIdNotSetOnQueryParam.Code)
		return
	}
	filterParams := util.ExtractFilterParams(r)

	span.SetAttributes(attribute.String("customer.id", id))
	actionLogs, err := c.customerService.GetCustomerActionLogByID(ctx, id, *filterParams)
	if err != nil {
		span.RecordError(err)
		log.Errorf("[GetCustomerActionLogByID] service error: %v", err)
		localization.SendErrorByCodeResponse(w, err.Error())
		return
	}
	localization.SendSuccessResponse(w, localization.CustomerActionLogRetrievedSuccessfully, actionLogs)
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
//	@Router			/customers/enable/{id} [patch]
func (c *customerAdapter) SetEnableCustomerSession(w http.ResponseWriter, r *http.Request) {
	ctx, span := util.TraceLogger(r.Context(), "handler", "setEnableCustomerSession", "handler", "customer")
	defer span.End()
	log := util.LoggerFromCtx(ctx, c.logger)
	id := chi.URLParam(r, "id")
	if id == "" {
		log.Errorf("[CustomerH] id not set")
		localization.SendBadRequestResponse(w, localization.ErrorIdNotSetOnQueryParam.Code)
		return
	}

	span.SetAttributes(attribute.String("customer.id", id))
	otp, err := c.customerService.CreateEnableCustomerSession(ctx, id)
	if err != nil {
		span.RecordError(err)
		log.Errorf("[SetEnableCustomerSession] service error: %v", err)
		localization.SendErrorByCodeResponse(w, err.Error())
		return
	}
	log.Infof("[SetEnableCustomerSession] OTP generated successfully for customer id: %s", id)
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
//	@Accept			json
//
//	@Produce		json
//	@Param			id		path		string									true	"Customer ID"
//	@Param			body	body		customer.CustomerDisableDTO	true	"Body (fields: is_temporary, disable_reason)"
//	@Success		200		{object}	localization.StandardResponse{data=nil}	"Customer disabled successfully"
//	@Failure		400		{object}	localization.StandardResponse{data=nil}	"Bad request - Customer ID required"
//	@Failure		500		{object}	localization.StandardResponse{data=nil}	"Internal server error"
//	@Security		BearerAuth
//	@Router			/customers/disable/{id} [patch]
func (c *customerAdapter) DisableCustomer(w http.ResponseWriter, r *http.Request) {
	ctx, span := util.TraceLogger(r.Context(), "handler", "disableCustomer", "handler", "customer")
	defer span.End()
	log := util.LoggerFromCtx(ctx, c.logger)
	id := chi.URLParam(r, "id")
	if id == "" {
		log.Errorf("[CustomerH] id not set")
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
		log.Errorf("[CustomerH][Disable] invalid reason: %v", err)
		localization.SendErrorByCodeResponse(w, localization.ErrorInvalidInputParameter.Code)
		return
	}

	span.SetAttributes(attribute.String("customer.id", id))
	err := c.customerService.DisableCustomerByID(ctx, id, payload)
	if err != nil {
		span.RecordError(err)
		log.Errorf("[DisableCustomer] service error: %v", err)
		localization.SendErrorByCodeResponse(w, err.Error())
		return
	}
	log.Infof("[DisableCustomer] request sent successfully for customer id: %s", id)
	localization.SendSuccessResponse(w, localization.CustomerDisableRequestCreatedSuccessfully, nil)
}

// EnableCustomer enables a customer by ID using OTP verification
//
//	@Summary		Enable customer with OTP verification
//	@Description	Enables a customer by their ID using OTP verification. This endpoint verifies the OTP sent to the customer and enables their account.
//	@Tags			Customers
//	@Accept			json
//	@Produce		json
//	@Param			id		path		string									true	"Customer ID"
//	@Param			body	body		customer.CustomerEnableDTO	true	"Enable customer payload (user_otp)"
//	@Success		200		{object}	localization.StandardResponse{data=nil}	"Customer enabled successfully"
//	@Failure		400		{object}	localization.StandardResponse{data=nil}	"Bad request - Customer ID required or invalid OTP"
//	@Failure		500		{object}	localization.StandardResponse{data=nil}	"Internal server error"
//	@Security		BearerAuth
//	@Router			/customers/enable_otp_verify/{id} [patch]
func (c *customerAdapter) EnableCustomer(w http.ResponseWriter, r *http.Request) {
	ctx, span := util.TraceLogger(r.Context(), "handler", "enableCustomer", "handler", "customer")
	defer span.End()
	log := util.LoggerFromCtx(ctx, c.logger)
	id := chi.URLParam(r, "id")
	if id == "" {
		log.Errorf("[CustomerH] id not set")
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
		log.Errorf("[EnableCustomer] service error: %v", err)
		localization.SendErrorByCodeResponse(w, err.Error())
		return
	}
	log.Infof("[EnableCustomer] request sent successfully for customer id: %s", id)
	localization.SendSuccessResponse(w, localization.CustomerEnableRequestCreatedSuccessfully, nil)
}

// GetCustomerDetail
//
//	@Summary		Get Customer Detail
//	@Description	Retrieve customer details with pagination, filtering, and search. Searchable fields: full_name, phone_number, gender, user_name, user_code, is_blocked, kyc_level.
//	@Tags			Customers
//	@Accept			json
//	@Produce		json
//	@Param			page				query		int													false	"Page number"
//	@Param			per_page			query		int													false	"Items per page"
//	@Param			search				query		string												false	"Search term (searches full_name, phone_number, gender, user_name, user_code, is_blocked, kyc_level)"
//	@Param			gender				query		string												false	"Filter by gender"
//	@Param			branch_code			query		string												false	"Filter by branch code"
//	@Param			kyc_level			query		int													false	"Filter by KYC level"
//	@Param			is_blocked			query		bool												false	"Filter by blocked status"
//	@Param			enabled				query		bool												false	"Filter by enabled status"
//	@Param			bps_reject_status	query		string												false	"Filter by BPS reject status"
//	@Success		200					{object}	localization.StandardResponse{data=[]model.User}	"Customer details retrieved successfully"
//	@Failure		400					{object}	localization.StandardResponse{data=nil}				"Bad request"
//	@Failure		500					{object}	localization.StandardResponse{data=nil}				"Internal server error"
//	@Security		BearerAuth
//	@Router			/customers [get]
func (c customerAdapter) GetCustomerDetail(w http.ResponseWriter, r *http.Request) {
	ctx, span := util.TraceLogger(r.Context(), "handler", "getCustomerDetail", "handler", "customer")
	defer span.End()
	log := util.LoggerFromCtx(ctx, c.logger)
	filterParams := util.ExtractFilterParams(r)

	search := r.URL.Query().Get("search")
	filter := r.URL.Query().Get("filter")

	if err := util.NoSpecialChars(search); err != nil {
		log.Errorf("[GetCustomerDetail] invalid search query: %v", err)
		localization.SendErrorByCodeResponse(w, err.Error())
		return
	}

	if err := util.NoSpecialChars(filter); err != nil {
		log.Errorf("[GetCustomerDetail] invalid filter query: %v", err)
		localization.SendErrorByCodeResponse(w, err.Error())
		return
	}

	customers, err := c.customerService.GetCustomersDetail(ctx, filterParams)
	if err != nil {
		span.RecordError(err)
		log.Errorf("[GetCustomerDetail] service error: %v", err)
		localization.SendErrorByCodeResponse(w, localization.ErrorFailedToGetCustomerDetail.Code)
		return
	}

	span.SetAttributes(attribute.Int("customer.count", len(customers.Data)))
	log.Infof("[GetCustomerDetail] retrieved %d customers", len(customers.Data))
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
//	@Success		200	{object}	localization.StandardResponse{data=model.User}	"Customer details retrieved successfully"
//	@Failure		400	{object}	localization.StandardResponse{data=nil}				"Bad request - Customer ID required"
//	@Failure		404	{object}	localization.StandardResponse{data=nil}				"Customer not found"
//	@Failure		500	{object}	localization.StandardResponse{data=nil}				"Internal server error"
//	@Security		BearerAuth
//	@Router			/customers/{id} [get]
func (c customerAdapter) GetCustomerByID(w http.ResponseWriter, r *http.Request) {
	ctx, span := util.TraceLogger(r.Context(), "handler", "getCustomerById", "handler", "customer")
	defer span.End()
	log := util.LoggerFromCtx(ctx, c.logger)
	id := chi.URLParam(r, "id")
	if id == "" {
		log.Errorf("[CustomerH] id not set")
		localization.SendBadRequestResponse(w, localization.ErrorIdNotSetOnQueryParam.Code)
		return
	}

	span.SetAttributes(attribute.String("customer.id", id))
	userDetail, err := c.customerService.GetCustomerByID(ctx, id)
	if err != nil {
		span.RecordError(err)
		log.Errorf("[GetCustomerByID] service error: %v", err)
		localization.SendErrorByCodeResponse(w, err.Error())
		return
	}
	log.Infof("[GetCustomerByID] customer retrieved successfully for id: %s", id)
	localization.SendSuccessResponse(w, localization.CustomerDetailSuccessfullyFetched, userDetail)

}

// GetBlockedCustomer retrieves blocked customers
//
//	@Summary		Get blocked customers
//	@Description	Retrieves a paginated list of blocked customers
//	@Tags			Customers
//	@Accept			json
//	@Produce		json
//	@Param			page				query		int																false	"Page number"		default(1)
//	@Param			per_page			query		int																false	"Items per page"	default(10)
//	@Param			search				query		string															false	"Search term (searches full_name, phone_number, gender, user_name, user_code, is_blocked, kyc_level)"
//	@Param			gender				query		string															false	"Filter by gender"
//	@Param			branch_code			query		string															false	"Filter by branch code"
//	@Param			kyc_level			query		int																false	"Filter by KYC level"
//	@Param			enabled				query		bool															false	"Filter by enabled status"
//	@Param			bps_reject_status	query		string															false	"Filter by BPS reject status"
//	@Success		200					{object}	localization.StandardResponse{data=customers_paginated_resp}	"Blocked customers retrieved successfully"
//	@Failure		400					{object}	localization.StandardResponse{data=nil}							"Bad request"
//	@Failure		500					{object}	localization.StandardResponse{data=nil}							"Internal server error"
//	@Security		BearerAuth
//	@Router			/customers/blocked [get]
func (c customerAdapter) GetBlockedCustomer(w http.ResponseWriter, r *http.Request) {
	ctx, span := util.TraceLogger(r.Context(), "handler", "getBlockedCustomer", "handler", "customer")
	defer span.End()
	log := util.LoggerFromCtx(ctx, c.logger)
	filterParams := util.ExtractFilterParams(r)

	search := r.URL.Query().Get("search")
	filter := r.URL.Query().Get("filter")

	if err := util.NoSpecialChars(search); err != nil {
		log.Errorf("[GetBlockedCustomer] invalid search query: %v", err)
		localization.SendErrorByCodeResponse(w, err.Error())
		return
	}

	if err := util.NoSpecialChars(filter); err != nil {
		localization.SendErrorByCodeResponse(w, err.Error())
		return
	}

	BlockedCustomer, err := c.customerService.GetBlockedCustomer(ctx, filterParams)
	if err != nil {
		span.RecordError(err)
		log.Errorf("[GetBlockedCustomer] service error: %v", err)
		localization.SendErrorByCodeResponse(w, localization.ErrorFailedToGetBlockedCustomer.Code)
		return
	}

	span.SetAttributes(attribute.Int("customer.blocked.count", len(BlockedCustomer.Data)))
	log.Infof("[GetBlockedCustomer] retrieved %d blocked customers", len(BlockedCustomer.Data))
	localization.SendSuccessResponse(w, localization.SuccessFullyFetchBlockCustomer, BlockedCustomer)
}

// GetLinkedAccount retrieves customer linked account
//
//	@Summary		Get Customer Linked Account
//	@Description	Retrives a list of customer linked account
//	@Tags			Customers
//	@Param			user_id	path	string	true	"Customer number"
//	@Produce		json
//	@Success		200	{object}	localization.StandardResponse{data=customers_paginated_resp}	"Blocked customers retrieved successfully"
//	@Failure		400	{object}	localization.StandardResponse{data=nil}							"Bad request"
//	@Failure		500	{object}	localization.StandardResponse{data=nil}							"Internal server error"
//	@Security		BearerAuth
//	@Router			/customers/linked_account/{user_id} [get]
func (c customerAdapter) GetLinkedAccount(w http.ResponseWriter, r *http.Request) {
	ctx, span := util.TraceLogger(r.Context(), "handler", "getLinkedAccount", "handler", "customer")
	defer span.End()
	log := util.LoggerFromCtx(ctx, c.logger)
	id := chi.URLParam(r, "user_id")
	if id == "" {
		log.Errorf("[CustomerH][GetByUserId] user_id not set")
		localization.SendBadRequestResponse(w, localization.ErrorIdNotSetOnQueryParam.Code)
		return
	}

	span.SetAttributes(attribute.String("customer.number", id))
	userDetail, err := c.customerService.GetLinkedAccount(ctx, id)
	if err != nil {
		span.RecordError(err)
		log.Errorf("[GetLinkedAccount] service error: %v", err)
		localization.SendErrorByCodeResponse(w, err.Error())
		return
	}
	log.Infof("[GetLinkedAccount] linked accounts retrieved successfully for user_id: %s", id)
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
//	@Param			body	body		customer.FaydaApproveRequest			true	"Fayda approval payload (risk_level: LOW, MEDIUM, HIGH)"
//	@Success		200		{object}	localization.StandardResponse{data=nil}	"Fayda customer approval request sent successfully"
//	@Failure		400		{object}	localization.StandardResponse{data=nil}	"Bad request - Customer ID required or invalid body"
//	@Failure		500		{object}	localization.StandardResponse{data=nil}	"Internal server error"
//	@Security		BearerAuth
//	@Router			/customers/fayda/enable/{id} [patch]
func (c customerAdapter) ApproveFaydaCustomer(w http.ResponseWriter, r *http.Request) {
	ctx, span := util.TraceLogger(r.Context(), "handler", "approveFaydaCustomer", "handler", "customer")
	defer span.End()
	log := util.LoggerFromCtx(ctx, c.logger)
	var req dto.FaydaApproveRequest
	id := chi.URLParam(r, "id")
	if id == "" {
		log.Errorf("[CustomerH] id not set")
		localization.SendBadRequestResponse(w, localization.ErrorInvalidAction.Message)
		return
	}

	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		span.RecordError(err)
		log.Errorf(localization.ErrorInvalidJSONPayload.Message)
		localization.SendBadRequestResponse(w, localization.MsgInvalidJSONPayload)
		return
	}

	if req.Validate() != nil {
		log.Errorf("[CustomerH][ApproveFayda] validate err")
	}
	span.SetAttributes(attribute.String("customer.id", id))
	if err := c.customerService.ApproveFaydaCustomer(ctx, id, req); err != nil {
		span.RecordError(err)
		log.Errorf("[ApproveFaydaCustomer] service error: %v", err)
		localization.SendErrorByCodeResponse(w, err.Error())
		return
	}

	log.Infof("[ApproveFaydaCustomer] request sent successfully for customer id: %s", id)
	localization.SendSuccessResponse(w, localization.FaydaCustomerApprovalRequestSent, nil)
}

// SearchCustomerByCIForAccountNumber searches for a customer by CIF or account number
//
//	@Summary		Search customer by CIF or account number
//	@Description	Retrieves customer information by searching with CIF (Customer Information File) or account number
//	@Tags			Customers
//	@Accept			json
//	@Produce		json
//	@Param			number	path		string												true	"CIF or Account Number"
//	@Success		200		{object}	localization.StandardResponse{data=object}		"Customer found successfully"
//	@Failure		400		{object}	localization.StandardResponse{data=nil}				"Bad request - Number parameter required"
//	@Failure		404		{object}	localization.StandardResponse{data=nil}				"Customer not found"
//	@Failure		500		{object}	localization.StandardResponse{data=nil}				"Internal server error"
//	@Security		BearerAuth
//	@Router			/customers/account_lookup/{number} [get]
func (c customerAdapter) SearchCustomerByCIForAccountNumber(w http.ResponseWriter, r *http.Request) {
	ctx, span := util.TraceLogger(r.Context(), "handler", "searchCustomerByCIForAccountNumber", "handler", "customer")
	defer span.End()
	log := util.LoggerFromCtx(ctx, c.logger)

	number := chi.URLParam(r, "number")
	if err := core.ValidateCustomerLookupRequest(number); err != nil {
		log.Errorf("[SearchCustomerByCIForAccountNumber] validation error: %v", err)
		localization.SendBadRequestResponse(w, "lookup value must not be empty and must be either a valid CID or account number or phone number")
		return
	}
	customer, err := c.customerService.SearchCustomerByCIForAccountNumber(ctx, number)
	if err != nil {
		span.RecordError(err)
		log.Errorf("[SearchCustomerByCIForAccountNumber] service error: %v", err)
		localization.SendErrorByCodeResponse(w, err.Error())
		return
	}
	localization.SendSuccessResponse(w, localization.SuccessCustomerDetailSuccessfullyFetched, customer)

}

// GetCustomerDetailByID retrieves detailed customer information by ID
//
//	@Summary		Get customer detail by ID
//	@Description	Retrieves comprehensive customer details including personal information, linked accounts, and other customer data by their ID
//	@Tags			Customers
//	@Accept			json
//	@Produce		json
//	@Param			id	path		string												true	"Customer ID"
//	@Success		200	{object}	localization.StandardResponse{data=customer.CustomerDetailResponse}	"Customer details retrieved successfully"
//	@Failure		400	{object}	localization.StandardResponse{data=nil}				"Bad request - Customer ID required"
//	@Failure		404	{object}	localization.StandardResponse{data=nil}				"Customer not found"
//	@Failure		500	{object}	localization.StandardResponse{data=nil}				"Internal server error"
//	@Security		BearerAuth
//	@Router			/customers/detail/{id} [get]
func (c *customerAdapter) GetCustomerDetailByID(w http.ResponseWriter, r *http.Request) {
	log := util.LoggerFromCtx(r.Context(), c.logger)
	id := chi.URLParam(r, "id")
	if id == "" {
		log.Errorf("[CustomerH] id not set")
		localization.SendBadRequestResponse(w, localization.ErrorIdNotSetOnQueryParam.Code)
		return
	}
	customerDetail, err := c.customerService.GetCustomerDetailByID(r.Context(), id)
	if err != nil {
		localization.SendErrorByCodeResponse(w, err.Error())
		return
	}
	localization.SendSuccessResponse(w, localization.SuccessUserRetrieved, customerDetail)
}
