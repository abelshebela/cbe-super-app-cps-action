package customer

import (
	"cbe-super-app-cps-action/internal/constants/interfaces/customer"
	"cbe-super-app-cps-action/internal/constants/localization"

	"cbe-super-app-cps-action/internal/service"
	util "cbe-super-app-cps-action/pkgs/utils"
	"net/http"

	utils "gitlab.com/bersufekadgetachew/cbe-super-app-shared/shared/utils"
	"strconv"
)

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

func (c customerAdapter) GetCustomerDetail(w http.ResponseWriter, r *http.Request) {
	kycLevel := r.URL.Query().Get("kyc_level")
	filterParams := util.ExtractFilterParams(r)

	kycLevelInt, err := strconv.Atoi(kycLevel)
	if err != nil {
		localization.SendBadRequestResponse(w, "Invalid kyc_level parameter")
		return
	}
	customers, err := c.customerService.GetCustomersDetail(r.Context(), kycLevelInt, filterParams)
	if err != nil {
		localization.SendErrorByCodeResponse(w, localization.ErrorFailedToGetCustomerDetail.Code)
		return
	}

	localization.SendSuccessResponse(w, localization.SuccessCustomerDetailSuccessfullyFetched, customers)
}

func (c customerAdapter) GetCustomerByID(w http.ResponseWriter, r *http.Request) {
	id := r.URL.Query().Get("id")
	if id == "" {
		c.logger.Errorf("id not set on param")
		localization.SendBadRequestResponse(w, localization.ErrorIdNotSetOnQueryParam.Code)
		return
	}
	ctx := r.Context()
	userDetail, err := c.customerService.GetCustomerByID(ctx, id)
	if err != nil {
		c.logger.Errorf("error while fetching get customer detail:", err)
		localization.SendErrorByCodeResponse(w, localization.ErrorFailedToGetCustomerDetail.Code)
		return
	}
	localization.SendSuccessResponse(w, localization.CustomerDetailSuccessfullyFetched, userDetail)

}

func (c customerAdapter) GetBlockedCustomer(w http.ResponseWriter, r *http.Request) {
	filterParams := util.ExtractFilterParams(r)

	BlockedCustomer, err := c.customerService.GetBlockedCustomer(r.Context(), filterParams)
	if err != nil {
		localization.SendErrorByCodeResponse(w, localization.ErrorFailedToGetBlockedCustomer.Code)
	}

	localization.SendSuccessResponse(w, localization.SuccessFullyFetchBlockCustomer, BlockedCustomer)
}
