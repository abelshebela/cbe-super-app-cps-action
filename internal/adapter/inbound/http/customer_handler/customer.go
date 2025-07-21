package customerhandler

import (
	"net/http"
	"strconv"

	"github.com/go-chi/chi/v5"
	// "gitlab.com/bersufekadgetachew/cbe-super-app-shared/shared/common"
	// "gitlab.com/bersufekadgetachew/cbe-super-app-shared/shared/entities/member"
	"gitlab.com/bersufekadgetachew/cbe-super-app-shared/shared/utils"

	"github.com/CBE-Super-App/cbe-super-app-cps-action/internal/application/customer"
	// "github.com/CBE-Super-App/cbe-super-app-cps-action/internal/application/middleware"
	// "github.com/CBE-Super-App/cbe-super-app-cps-action/internal/domain/customer/entity"
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
	filterParams := common_util.ExtractFilterParams(r)

	ctx := r.Context()
	customers, err := c.applicationService.GetCustomersDeatil(ctx, filterParams)
	if err != nil {
		util.SendErrorResponse(w, err.Error(), 0, nil)

		return
	}

	// def, _ := common.GetSuccessResponseByKey("SUCCESS Fetched User")
	// data, _ := util.StructToMap(customers)
	// util.BaseResponseMaker(data, w, def.Message, http.StatusAccepted)

	util.WriteSuccessResponse(w, customers, "customer fetched successfully")

}

func (c CustomerHTTPHandler) GetCustomerByID(w http.ResponseWriter, r *http.Request) {
	id := chi.URLParam(r, "id")
	ctx := r.Context()
	customer, err := c.applicationService.GetCustomerByID(ctx, id)
	if err != nil {
		util.SendErrorResponse(w, err.Error(), 0, nil)
		return
	}

	def, _ := common.GetSuccessResponseByKey("SUCCESS")
	data, _ := util.StructToMap(customer)
	util.BaseResponseMaker(data, w, def.Message, http.StatusAccepted)
}

func (c CustomerHTTPHandler) GetFaydaCustomer(w http.ResponseWriter, r *http.Request) {
	filterParams := common_util.ExtractFilterParams(r)

	ctx := r.Context()
	customers, err := c.applicationService.GetFaydaCustomersDeatil(ctx, filterParams)
	if err != nil {
		util.SendErrorResponse(w, err.Error(), 0, nil)

		return
	}

	// def, _ := common.GetSuccessResponseByKey("SUCCESS")
	// data, _ := util.StructToMap(customers)
	// util.BaseResponseMaker(data, w, def.Message, http.StatusAccepted)

	util.WriteSuccessResponse(w, customers, "level 1 customer fetched successfully")

}

func (c CustomerHTTPHandler) GetFaydaCustomerByID(w http.ResponseWriter, r *http.Request) {
	id := chi.URLParam(r, "id")
	ctx := r.Context()
	customer, err := c.applicationService.GetFaydaCustomerByID(ctx, id)
	if err != nil {
		util.SendErrorResponse(w, err.Error(), 0, nil)
		return
	}

	def, _ := common.GetSuccessResponseByKey("SUCCESS")
	data, _ := util.StructToMap(customer)
	util.BaseResponseMaker(data, w, def.Message, http.StatusAccepted)
}

func (c CustomerHTTPHandler) GetBlockedCustomer(w http.ResponseWriter, r *http.Request) {
	query := r.URL.Query()

	page := constant.DefaultPage
	if pageInt, err := strconv.Atoi(query.Get("page")); err == nil && pageInt > 0 {
		page = pageInt
	}

	per_page := constant.DefaultPerPage
	if perPageInt, err := strconv.Atoi(query.Get("per_page")); err == nil &&
		perPageInt <= 10 && perPageInt > 0 {
		per_page = perPageInt
	}

	search := query.Get("search")
	filter := query.Get("filter")

	filterParams := &constant.Filter{
		Page:    page,
		PerPage: per_page,
		Search:  search,
		Filters: filter,
	}

	ctx := r.Context()
	customers, err := c.applicationService.GetBlockedCustomer(ctx, filterParams)
	if err != nil {
		util.SendErrorResponse(w, err.Error(), 0, nil)

		return
	}

	def, _ := common.GetSuccessResponseByKey("SUCCESS")
	data, _ := util.StructToMap(customers)
	util.BaseResponseMaker(data, w, def.Message, http.StatusAccepted)

}
