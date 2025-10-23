package customer

import (
	dto "cbe-super-app-cps-action/internal/constants/dto/customer"
	"cbe-super-app-cps-action/internal/constants/interfaces/customer"
	"cbe-super-app-cps-action/internal/constants/localization"
	"cbe-super-app-cps-action/internal/constants/model"
	"cbe-super-app-cps-action/internal/constants/types"
	"encoding/json"

	"cbe-super-app-cps-action/internal/service"
	util "cbe-super-app-cps-action/pkgs/utils"
	"net/http"

	"strconv"

	"github.com/go-chi/chi/v5"
	utils "gitlab.com/bersufekadgetachew/cbe-super-app-shared/shared/utils"
)

type customer_resp *model.User
type customers_paginated_resp *types.PaginatedResponse[[]*model.User]

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
// @Summary Initiate enable customer session
// @Description Initiates enabling a customer session by generating an OTP for the customer
// @Tags Customers
// @Accept json
// @Produce json
// @Param id path string true "Customer ID"
// @Success 200 {object} localization.StandardResponse{data=dto.CustomerEnableSessionResponse} "OTP generated successfully"
// @Failure 400 {object} localization.StandardResponse{data=nil} "Bad request - Customer ID required"
// @Failure 500 {object} localization.StandardResponse{data=nil} "Internal server error"
// @Security BearerAuth
// @Router /customers/{id}/enable-session [post]
func (c *customerAdapter) SetEnableCustomerSession(w http.ResponseWriter, r *http.Request) {
	id := chi.URLParam(r, "id")
	if id == "" {
		c.logger.Errorf("id not set on param")
		localization.SendBadRequestResponse(w, localization.ErrorIdNotSetOnQueryParam.Code)
		return
	}
	ctx := r.Context()
	otp, err := c.customerService.CreateEnableCustomerSession(ctx, id)
	if err != nil {
		c.logger.Errorf("error while fetching get customer detail:", err)
		localization.SendErrorByCodeResponse(w, err.Error())
		return
	}
	localization.SendSuccessResponse(w, localization.CustomerEnableRequestSessionCreatedSuccessfully, dto.CustomerEnableSessionResponse{
		Otp: otp,
	})
}

// DisableCustomer disables a customer by ID
// @Summary Disable customer
// @Description Disables a customer by their ID
// @Tags Customers
// @Accept json
// @Produce json
// @Param id path string true "Customer ID"
// @Success 200 {object} localization.StandardResponse{data=nil} "Customer disabled successfully"
// @Failure 400 {object} localization.StandardResponse{data=nil} "Bad request - Customer ID required"
// @Failure 500 {object} localization.StandardResponse{data=nil} "Internal server error"
// @Security BearerAuth
// @Router /customers/{id}/disable [post]
func (c *customerAdapter) DisableCustomer(w http.ResponseWriter, r *http.Request) {
	id := chi.URLParam(r, "id")
	if id == "" {
		c.logger.Errorf("id not set on param")
		localization.SendBadRequestResponse(w, localization.ErrorIdNotSetOnQueryParam.Code)
		return
	}

	var payload dto.CustomerDisableDTO
	if err := json.NewDecoder(r.Body).Decode(&payload); err != nil {
		localization.SendErrorByCodeResponse(w, localization.ErrorInvalidRequestBody.Code)
		return
	}

	ctx := r.Context()
	err := c.customerService.DisableCustomerByID(ctx, id, payload)
	if err != nil {
		c.logger.Errorf("error while fetching get customer detail:", err)
		localization.SendErrorByCodeResponse(w, err.Error())
		return
	}
	localization.SendSuccessResponse(w, localization.CustomerDisableRequestCreatedSuccessfully, nil)
}

// EnableCustomer enables a customer by ID using OTP
// @Summary Enable customer
// @Description Enables a customer by their ID using OTP
// @Tags Customers
// @Accept json
// @Produce json
// @Param id path string true "Customer ID"
// @Param body body dto.CustomerEnableDTO true "Enable customer payload"
// @Success 200 {object} localization.StandardResponse{data=nil} "Customer enabled successfully"
// @Failure 400 {object} localization.StandardResponse{data=nil} "Bad request - Customer ID required or invalid body"
// @Failure 500 {object} localization.StandardResponse{data=nil} "Internal server error"
// @Security BearerAuth
// @Router /customers/{id}/enable [post]
func (c *customerAdapter) EnableCustomer(w http.ResponseWriter, r *http.Request) {
	id := chi.URLParam(r, "id")
	if id == "" {
		c.logger.Errorf("id not set on param")
		localization.SendBadRequestResponse(w, localization.ErrorIdNotSetOnQueryParam.Code)
		return
	}

	var payload dto.CustomerEnableDTO
	if err := json.NewDecoder(r.Body).Decode(&payload); err != nil {
		localization.SendErrorByCodeResponse(w, localization.ErrorInvalidRequestBody.Code)
		return
	}

	ctx := r.Context()
	err := c.customerService.EnableCustomerByID(ctx, id, payload.UserOTP)
	if err != nil {
		c.logger.Errorf("error while fetching get customer detail:", err)
		localization.SendErrorByCodeResponse(w, err.Error())
		return
	}
	localization.SendSuccessResponse(w, localization.CustomerEnableRequestCreatedSuccessfully, nil)
}

// GetCustomerDetail retrieves customer details with optional KYC level filtering
// @Summary Get customer details
// @Description Retrieves a paginated list of customer details with optional KYC level filtering
// @Tags Customers
// @Accept json
// @Produce json
// @Param kyc_level query int false "KYC Level filter" default(0)
// @Param page query int false "Page number" default(1)
// @Param per_page query int false "Items per page" default(10)
// @Param search query string false "Search term"
// @Success 200 {object} localization.StandardResponse{data=customers_paginated_resp} "Customer details retrieved successfully"
// @Failure 400 {object} localization.StandardResponse{data=nil} "Bad request - Invalid KYC level parameter"
// @Failure 500 {object} localization.StandardResponse{data=nil} "Internal server error"
// @Security BearerAuth
// @Router /customers [get]
func (c customerAdapter) GetCustomerDetail(w http.ResponseWriter, r *http.Request) {
	var kycLevelInt int
	var err error
	kycLevel := r.URL.Query().Get("kyc_level")
	filterParams := util.ExtractFilterParams(r)

	if kycLevel != "" {
		kycLevelInt, err = strconv.Atoi(kycLevel)
		if err != nil {
			localization.SendBadRequestResponse(w, "Invalid kyc_level parameter")
			return
		}
	} else {
		kycLevelInt = 0
	}

	customers, err := c.customerService.GetCustomersDetail(r.Context(), kycLevelInt, filterParams)
	if err != nil {
		localization.SendErrorByCodeResponse(w, localization.ErrorFailedToGetCustomerDetail.Code)
		return
	}

	localization.SendSuccessResponse(w, localization.SuccessCustomerDetailSuccessfullyFetched, customers)
}

// GetCustomerByID retrieves a specific customer by ID
// @Summary Get customer by ID
// @Description Retrieves detailed information for a specific customer by their ID
// @Tags Customers
// @Accept json
// @Produce json
// @Param id path string true "Customer ID"
// @Success 200 {object} localization.StandardResponse{data=customer_resp} "Customer details retrieved successfully"
// @Failure 400 {object} localization.StandardResponse{data=nil} "Bad request - Customer ID required"
// @Failure 404 {object} localization.StandardResponse{data=nil} "Customer not found"
// @Failure 500 {object} localization.StandardResponse{data=nil} "Internal server error"
// @Security BearerAuth
// @Router /customers/{id} [get]
func (c customerAdapter) GetCustomerByID(w http.ResponseWriter, r *http.Request) {
	id := chi.URLParam(r, "id")
	if id == "" {
		c.logger.Errorf("id not set on param")
		localization.SendBadRequestResponse(w, localization.ErrorIdNotSetOnQueryParam.Code)
		return
	}
	ctx := r.Context()
	userDetail, err := c.customerService.GetCustomerByID(ctx, id)
	if err != nil {
		c.logger.Errorf("error while fetching get customer detail:", err)
		localization.SendErrorByCodeResponse(w, localization.UserNotFoundWithGivenID.Code)
		return
	}
	localization.SendSuccessResponse(w, localization.CustomerDetailSuccessfullyFetched, userDetail)

}

// GetBlockedCustomer retrieves blocked customers
// @Summary Get blocked customers
// @Description Retrieves a paginated list of blocked customers
// @Tags Customers
// @Accept json
// @Produce json
// @Param page query int false "Page number" default(1)
// @Param per_page query int false "Items per page" default(10)
// @Param search query string false "Search term"
// @Success 200 {object} localization.StandardResponse{data=customers_paginated_resp} "Blocked customers retrieved successfully"
// @Failure 400 {object} localization.StandardResponse{data=nil} "Bad request"
// @Failure 500 {object} localization.StandardResponse{data=nil} "Internal server error"
// @Security BearerAuth
// @Router /customers/blocked [get]
func (c customerAdapter) GetBlockedCustomer(w http.ResponseWriter, r *http.Request) {
	filterParams := util.ExtractFilterParams(r)

	BlockedCustomer, err := c.customerService.GetBlockedCustomer(r.Context(), filterParams)
	if err != nil {
		localization.SendErrorByCodeResponse(w, localization.ErrorFailedToGetBlockedCustomer.Code)
		return
	}

	localization.SendSuccessResponse(w, localization.SuccessFullyFetchBlockCustomer, BlockedCustomer)
}
