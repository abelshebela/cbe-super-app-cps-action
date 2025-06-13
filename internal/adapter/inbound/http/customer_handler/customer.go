package customerhandler

import (
	"net/http"
	"strconv"

	"gitlab.com/bersufekadgetachew/cbe-super-app-cps-action/internal/application/customer"
	"gitlab.com/bersufekadgetachew/cbe-super-app-cps-action/internal/application/middleware"
	"gitlab.com/bersufekadgetachew/cbe-super-app-cps-action/internal/domain/customer/entity"
	"gitlab.com/bersufekadgetachew/cbe-super-app-cps-action/internal/port/inbound"
	constant "gitlab.com/bersufekadgetachew/cbe-super-app-cps-action/utils"

	"github.com/go-chi/chi/v5"
	"gitlab.com/bersufekadgetachew/cbe-super-app-shared/shared/common"
	"gitlab.com/bersufekadgetachew/cbe-super-app-shared/shared/entities/member"
	"gitlab.com/bersufekadgetachew/cbe-super-app-shared/shared/utils"
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
	customers, err := c.applicationService.GetCustomersDeatil(ctx, filterParams)
	if err != nil {
		middleware.ErrorHandler(w, err)
		return
	}

	res := common.Response[*entity.CustomerRespose]{
		ResponseWriter: w,
		Status:         http.StatusOK,
		Data:           customers,
	}
	res.SendJSON()
}

func (c CustomerHTTPHandler) GetCustomerByID(w http.ResponseWriter, r *http.Request) {
	id := chi.URLParam(r, "id")
	ctx := r.Context()
	customer, err := c.applicationService.GetCustomerByID(ctx, id)
	if err != nil {
		middleware.ErrorHandler(w, err)
		return
	}

	res := common.Response[*member.User]{
		ResponseWriter: w,
		Status:         http.StatusOK,
		Data:           customer,
	}
	res.SendJSON()
}
