package customerhandler

import (
	"gitlab.com/bersufekadgetachew/cbe-super-app-cps-ms/internal/application/customer"
	"gitlab.com/bersufekadgetachew/cbe-super-app-cps-ms/internal/application/middleware"
	"gitlab.com/bersufekadgetachew/cbe-super-app-cps-ms/internal/domain/customer/entity"
	"gitlab.com/bersufekadgetachew/cbe-super-app-cps-ms/internal/port/inbound"
	constant "gitlab.com/bersufekadgetachew/cbe-super-app-cps-ms/utils"
	"fmt"
	"net/http"
	"strconv"

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
	page := query.Get("page")
	pageInt := constant.DefaultPage
	if page != "" {
		var err error
		pageInt, err = strconv.Atoi(page)
		if err != nil {
			c.logger.Errorf("invalid page", err)
			err = fmt.Errorf("invalid page %w", constant.ErrorDefinition{
				Code:    http.StatusBadRequest,
				Message: "invalid page",
			})
			middleware.ErrorHandler(w, err)
			return
		}
		if pageInt <= 0 {
			pageInt = constant.DefaultPage
		}
	}

	perPage := query.Get("per_page")
	perPageInt := constant.DefaultPerPage
	if perPage != "" {
		var err error
		perPageInt, err = strconv.Atoi(perPage)
		if err != nil {
			c.logger.Errorf("invalid per page", err)
			err = fmt.Errorf("invalid per page %w", constant.ErrorDefinition{
				Code:    http.StatusBadRequest,
				Message: "invalid per page",
			})
			middleware.ErrorHandler(w, err)
			return
		}
		if perPageInt > 10 {
			perPageInt = constant.DefaultPerPage
		}
	}

	search := query.Get("search")
	filter := query.Get("filter")

	filterParams := &constant.Filter{
		Page:    pageInt,
		PerPage: perPageInt,
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
