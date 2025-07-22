package customerhandler

import (
	"fmt"
	"net/http"
	"strconv"

	"gitlab.com/bersufekadgetachew/cbe-super-app-shared/shared/utils"

	"github.com/CBE-Super-App/cbe-super-app-cps-action/internal/application/customer"
	"github.com/CBE-Super-App/cbe-super-app-cps-action/internal/port/inbound"
	constant "github.com/CBE-Super-App/cbe-super-app-cps-action/utils"

	common "github.com/CBE-Super-App/cbe-super-app-cps-action/pkgs/common"
	common_util "github.com/CBE-Super-App/cbe-super-app-cps-action/pkgs/utils"
	util "github.com/CBE-Super-App/cbe-super-app-cps-action/pkgs/utils"
)

type CustomerHTTPHandler struct {
	applicationService customer.ApplicationService
	logger             utils.Logger
}

func NewCustomerHTTPHandler(applicationService customer.ApplicationService, logger utils.Logger) inbound.CustomerDetail {
	return CustomerHTTPHandler{
		applicationService: applicationService,
		logger:             logger,
	}
}

func (c CustomerHTTPHandler) GetCustomerDetail(w http.ResponseWriter, r *http.Request) {
	filterParams, kycLevel, err := ExtractKYCFilterParams(r)

	if err != nil {
		util.SendErrorResponse(w, err.Error(), 0, nil)
		return
	}

	customers, err := c.applicationService.GetCustomersDeatil(r.Context(), kycLevel, filterParams)
	if err != nil {
		util.SendErrorResponse(w, err.Error(), 0, nil)

		return
	}

	util.WriteSuccessResponse(w, customers, "customer fetched successfully")
}

func (c CustomerHTTPHandler) GetCustomerByID(w http.ResponseWriter, r *http.Request) {
	id, ok := common_util.GetParam(r, "id")
	if !ok {
		c.logger.Errorf("missing or invalid parameter 'id'")
		common_util.SendErrorResponse(w, common_util.InvalidInputParameters, 0, nil)
		return
	}

	customer, err := c.applicationService.GetCustomerByID(r.Context(), id)
	if err != nil {
		util.SendErrorResponse(w, err.Error(), 0, nil)
		return
	}

	def, _ := common.GetSuccessResponseByKey("SUCCESS")
	data, _ := util.StructToMap(customer)
	util.BaseResponseMaker(data, w, def.Message, http.StatusAccepted)
}

func (c CustomerHTTPHandler) GetBlockedCustomer(w http.ResponseWriter, r *http.Request) {
	filterParams := common_util.ExtractFilterParams(r)

	customers, err := c.applicationService.GetBlockedCustomer(r.Context(), filterParams)
	if err != nil {
		util.SendErrorResponse(w, err.Error(), 0, nil)

		return
	}

	def, _ := common.GetSuccessResponseByKey("SUCCESS")
	data, _ := util.StructToMap(customers)
	util.BaseResponseMaker(data, w, def.Message, http.StatusAccepted)

}

func ExtractKYCFilterParams(r *http.Request) (*constant.Filter, int, error) {
	query := r.URL.Query()

	page := constant.DefaultPage
	if pageInt, err := strconv.Atoi(query.Get("page")); err == nil && pageInt > 0 {
		page = pageInt
	}

	perPage := constant.DefaultPerPage
	if perPageInt, err := strconv.Atoi(query.Get("per_page")); err == nil &&
		perPageInt > 0 {
		perPage = perPageInt
	}

	kycLevelStr := query.Get("kyc_level")
	if kycLevelStr == "" {
		return nil, 0, fmt.Errorf("KYC_LEVEL_REQUIRED")
	}

	kycLevel, err := strconv.Atoi(kycLevelStr)
	if err != nil || kycLevel < 0 || kycLevel > 2 {
		return nil, 0, fmt.Errorf("INVALID_KYC_LEVEL")
	}

	return &constant.Filter{
		Page:    page,
		PerPage: perPage,
		Search:  query.Get("search"),
		Filters: query.Get("filter"),
	}, kycLevel, nil
}
