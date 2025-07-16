package amount_based_auth_handler

import (
	"encoding/json"
	"fmt"
	"strconv"

	"github.com/CBE-Super-App/cbe-super-app-cps-action/internal/adapter/outbound/model"
	"github.com/CBE-Super-App/cbe-super-app-cps-action/internal/application/amount_based_auth_app"
	"github.com/CBE-Super-App/cbe-super-app-cps-action/internal/application/middleware"
	amount_based_auth_domain "github.com/CBE-Super-App/cbe-super-app-cps-action/internal/domain/amount_based_auth"
	"github.com/CBE-Super-App/cbe-super-app-cps-action/internal/port/inbound"
	common_util "github.com/CBE-Super-App/cbe-super-app-cps-action/pkgs/utils"
	constant "github.com/CBE-Super-App/cbe-super-app-cps-action/utils"

	"github.com/go-chi/chi/v5"

	"net/http"

	"gitlab.com/bersufekadgetachew/cbe-super-app-shared/shared/utils"
)

type AmountBasedAuthHandler struct {
	amountBasedAuthService amount_based_auth_app.ApplicationService
	logger                 utils.Logger
}

type Resp struct {
	message string
}

func NewAmountBasedAuthHandler(service amount_based_auth_app.ApplicationService, logger utils.Logger) inbound.AmountBasedAuthHandler {
	return &AmountBasedAuthHandler{
		amountBasedAuthService: service,
		logger:                 logger,
	}
}

func (a AmountBasedAuthHandler) GetAllAmountBasedAuth(w http.ResponseWriter, r *http.Request) {
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
	customers, err := a.amountBasedAuthService.GetAllAmountBasedDetail(ctx, filterParams)
	if err != nil {

		common_util.SendErrorResponse(w, err.Error(), 0, nil)

		return
	}

	data, _ := common_util.StructToMap(customers)
	common_util.BaseResponseMaker(data, w, "Successfuly fetched", http.StatusAccepted)

}
func (a *AmountBasedAuthHandler) UpdateAmountBasedAuth(w http.ResponseWriter, r *http.Request) {
	id := chi.URLParam(r, "id")

	var request amount_based_auth_domain.UpdateAmountBasedAuth
	if err := json.NewDecoder(r.Body).Decode(&request); err != nil {
		a.logger.Errorf("failed to decode request body: %v", err)
		err := fmt.Errorf("failed to decode request body: %w", constant.ErrorDefinition{
			Code:    http.StatusBadRequest,
			Message: "invalid request body",
		})
		middleware.ErrorHandler(w, err)
		return
	}

	var cpsActionRequest model.CreateCPSAction

	user_code := r.Context().Value(constant.ContextKey("user_code")).(string)
	full_name := r.Context().Value(constant.ContextKey("full_name")).(string)
	phone_number := r.Context().Value(constant.ContextKey("phone_number")).(string)
	department := r.Context().Value(constant.ContextKey("department")).(string)

	cpsActionRequest.MakerUser = model.User{
		UserCode:    user_code,
		FullName:    full_name,
		PhoneNumber: phone_number,
	}
	cpsActionRequest.Department = department
	request.Id = id
	cpsActionRequest.CurrentData = request

	ctx := r.Context()
	amountBasedAuth, err := a.amountBasedAuthService.UpdateAmountBasedAuth(ctx, request, cpsActionRequest)
	if err != nil {
		middleware.ErrorHandler(w, err)
		return
	}

	data, err := common_util.StructToMap(amountBasedAuth)
	if err != nil {
		common_util.SendErrorResponse(w, err, 500, nil)
		return
	}
	common_util.BaseResponseMaker(data, w, "Successfuly updated", 200)
	return

}

func (a *AmountBasedAuthHandler) ApproveAmountBasedAuth(w http.ResponseWriter, r *http.Request) {
	id := chi.URLParam(r, "id")

	ctx := r.Context()

	var cpsReq model.AuthorizeCPSAction

	user_code := r.Context().Value(constant.ContextKey("user_code")).(string)
	full_name := r.Context().Value(constant.ContextKey("full_name")).(string)
	phone_number := r.Context().Value(constant.ContextKey("phone_number")).(string)
	department := r.Context().Value(constant.ContextKey("department")).(string)

	cpsReq.CheckerUser = model.User{
		UserCode:    user_code,
		FullName:    full_name,
		PhoneNumber: phone_number,
	}
	cpsReq.Department = department

	amountBasedAuth, err := a.amountBasedAuthService.ApproveAmountBasedAuth(ctx, id, cpsReq)
	if err != nil {
		common_util.SendErrorResponse(w, err.Error(), 0, nil)
		return
	}

	data, err := common_util.StructToMap(amountBasedAuth)
	if err != nil {
		common_util.SendErrorResponse(w, err, 500, nil)
		return
	}
	common_util.BaseResponseMaker(data, w, "Successfuly updated", 200)
	return
}

func (a *AmountBasedAuthHandler) RejectAmountBasedAuth(w http.ResponseWriter, r *http.Request) {
	id := chi.URLParam(r, "id")

	var cpsReq model.RejectAuthTierCPSAction

	if err := json.NewDecoder(r.Body).Decode(&cpsReq); err != nil {
		a.logger.Errorf("failed to decode amount based auth request", err)
		err = fmt.Errorf("%w", constant.ErrorDefinition{
			Code:    http.StatusBadRequest,
			Message: "invalid request",
		})
		middleware.ErrorHandler(w, err)
		return
	}

	ctx := r.Context()
	rejectAction, err := a.amountBasedAuthService.RejectAmountBasedAuth(ctx, id, cpsReq)
	if err != nil {
		middleware.ErrorHandler(w, err)
		return
	}

	data, err := common_util.StructToMap(rejectAction)
	if err != nil {
		common_util.SendErrorResponse(w, err, 500, nil)
		return
	}
	common_util.BaseResponseMaker(data, w, "Successfuly updated", 200)
	return
}
