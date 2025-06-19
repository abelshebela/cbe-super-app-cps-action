package bank

import (
	"encoding/json"
	"fmt"
	"net/http"
	"strconv"

	"github.com/go-chi/chi/v5"
	"gitlab.com/bersufekadgetachew/cbe-super-app-cps-action/internal/adapter/outbound/model"
	"gitlab.com/bersufekadgetachew/cbe-super-app-cps-action/internal/application/bank"
	"gitlab.com/bersufekadgetachew/cbe-super-app-cps-action/internal/application/middleware"
	"gitlab.com/bersufekadgetachew/cbe-super-app-cps-action/internal/domain/bank/dto"
	"gitlab.com/bersufekadgetachew/cbe-super-app-cps-action/internal/domain/bank/entity"
	inboundBank "gitlab.com/bersufekadgetachew/cbe-super-app-cps-action/internal/port/inbound/bank"
	constant "gitlab.com/bersufekadgetachew/cbe-super-app-cps-action/utils"
	"gitlab.com/bersufekadgetachew/cbe-super-app-shared/shared/common"
	"gitlab.com/bersufekadgetachew/cbe-super-app-shared/shared/utils"
)

type BankAdapter struct {
	bankHandler bank.BankHandlerService
	logger      utils.Logger
}

func InitBankAdapter(bankHandler bank.BankHandlerService, logger utils.Logger) inboundBank.BankAdapter {
	return &BankAdapter{
		bankHandler: bankHandler,
		logger:      logger,
	}
}

func (b *BankAdapter) CreateOneBank(w http.ResponseWriter, r *http.Request) {
	var bankRequest dto.CreateBankRequest

	if err := r.ParseMultipartForm(10 << 20); err != nil {
		b.logger.Errorf("failed to parse form data: %v", err)
		err = fmt.Errorf("failed to parse multipart form: %w", constant.ErrorDefinition{
			Code:    http.StatusBadRequest,
			Message: "invalid multipart form",
		})
		middleware.ErrorHandler(w, err)
		return
	}

	bankRequest.Name = r.FormValue("name")
	bankRequest.Code = r.FormValue("code")
	bankRequest.BIC = r.FormValue("bic")

	file, fileHeader, err := r.FormFile("logo")
	if err != nil {
		b.logger.Errorf("logo error: %v", err)
		err = fmt.Errorf("failed to read logo: %w", constant.ErrorDefinition{
			Code:    http.StatusBadRequest,
			Message: "missing or invalid logo",
		})
		middleware.ErrorHandler(w, err)
		return
	}
	defer file.Close()

	bankRequest.Logo = fileHeader

	var cpsRequest model.CreateCPSAction
	user_code := r.Context().Value(constant.ContextKey("user_code")).(string)
	full_name := r.Context().Value(constant.ContextKey("full_name")).(string)
	phone_number := r.Context().Value(constant.ContextKey("phone_number")).(string)
	department := r.Context().Value(constant.ContextKey("department")).(string)

	cpsRequest.ActionData = bankRequest
	cpsRequest.MakerUser = model.User{
		UserCode:    user_code,
		FullName:    full_name,
		PhoneNumber: phone_number,
	}
	cpsRequest.Department = department

	ctx := r.Context()
	cpsRes, err := b.bankHandler.CreateOneBank(ctx, cpsRequest)
	if err != nil {
		middleware.ErrorHandler(w, err)
		return
	}

	res := common.Response[*model.CpsAction]{
		ResponseWriter: w,
		Status:         http.StatusOK,
		Data:           cpsRes,
	}

	res.SendJSON()
}

func (b *BankAdapter) UpdateOneBank(w http.ResponseWriter, r *http.Request) {
	id := chi.URLParam(r, "id")

	var updateRequest dto.UpdateBankRequest

	if err := json.NewDecoder(r.Body).Decode(&updateRequest); err != nil {
		b.logger.Errorf("invalid input", err)
		err = fmt.Errorf("failed to decode update bank request error data %w", constant.ErrorDefinition{
			Code:    http.StatusBadRequest,
			Message: "invalid request",
		})
		middleware.ErrorHandler(w, err)
		return
	}

	user_code := r.Context().Value(constant.ContextKey("user_code")).(string)
	full_name := r.Context().Value(constant.ContextKey("full_name")).(string)
	phone_number := r.Context().Value(constant.ContextKey("phone_number")).(string)
	department := r.Context().Value(constant.ContextKey("department")).(string)

	var cpsReq model.CreateCPSAction

	updateRequest.ID = id
	cpsReq.MakerUser = model.User{
		UserCode:    user_code,
		FullName:    full_name,
		PhoneNumber: phone_number,
	}
	cpsReq.Department = department
	cpsReq.ActionData = updateRequest

	ctx := r.Context()
	cpsAction, err := b.bankHandler.UpdateOneBank(ctx, id, cpsReq)
	if err != nil {
		middleware.ErrorHandler(w, err)
		return
	}

	res := common.Response[*model.CpsAction]{
		ResponseWriter: w,
		Status:         http.StatusOK,
		Data:           cpsAction,
	}

	res.SendJSON()
}

func (b *BankAdapter) DeleteOneBank(w http.ResponseWriter, r *http.Request) {
	id := chi.URLParam(r, "id")

	var cpsReq model.CreateCPSAction

	user_code := r.Context().Value(constant.ContextKey("user_code")).(string)
	full_name := r.Context().Value(constant.ContextKey("full_name")).(string)
	phone_number := r.Context().Value(constant.ContextKey("phone_number")).(string)
	department := r.Context().Value(constant.ContextKey("department")).(string)

	cpsReq.ActionData = entity.Bank{
		ID: id,
	}
	cpsReq.MakerUser = model.User{
		UserCode:    user_code,
		FullName:    full_name,
		PhoneNumber: phone_number,
	}
	cpsReq.Department = department

	ctx := r.Context()
	cpsAction, err := b.bankHandler.DeleteOneBank(ctx, id, cpsReq)
	if err != nil {
		middleware.ErrorHandler(w, err)
		return
	}

	res := common.Response[*model.CpsAction]{
		ResponseWriter: w,
		Status:         http.StatusOK,
		Data:           cpsAction,
	}

	res.SendJSON()
}

func (b *BankAdapter) GetAllBank(w http.ResponseWriter, r *http.Request) {
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

	banks, err := b.bankHandler.GetAllBank(ctx, filterParams)
	if err != nil {
		middleware.ErrorHandler(w, err)
		return
	}

	res := common.Response[*entity.BankResponse]{
		ResponseWriter: w,
		Status:         http.StatusOK,
		Data:           banks,
	}
	res.SendJSON()
}

func (b *BankAdapter) GetOneBank(w http.ResponseWriter, r *http.Request) {
	id := chi.URLParam(r, "id")

	ctx := r.Context()

	bank, err := b.bankHandler.GetOneBank(ctx, id)
	if err != nil {
		middleware.ErrorHandler(w, err)
		return
	}

	res := common.Response[*entity.Bank]{
		ResponseWriter: w,
		Status:         http.StatusOK,
		Data:           bank,
	}

	res.SendJSON()
}

func (b *BankAdapter) Authorize(w http.ResponseWriter, r *http.Request) {
	action_code := chi.URLParam(r, "action_code")

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
	cpsReq.ActionCode = action_code

	ctx := r.Context()
	authAction, err := b.bankHandler.Authorize(ctx, cpsReq)
	if err != nil {
		middleware.ErrorHandler(w, err)
		return
	}

	res := common.Response[*model.CpsAction]{
		ResponseWriter: w,
		Status:         http.StatusOK,
		Data:           authAction,
	}
	res.SendJSON()
}

func (b *BankAdapter) Reject(w http.ResponseWriter, r *http.Request) {
	action_code := chi.URLParam(r, "action_code")

	var cpsReq model.RejectCPSAction

	if err := json.NewDecoder(r.Body).Decode(&cpsReq); err != nil {
		b.logger.Errorf("failed to decode bank request", err)
		err = fmt.Errorf("failed to decode bank request error data %w", constant.ErrorDefinition{
			Code:    http.StatusBadRequest,
			Message: "invalid request",
		})
		middleware.ErrorHandler(w, err)
		return
	}

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
	cpsReq.ActionCode = action_code

	ctx := r.Context()
	rejectAction, err := b.bankHandler.Reject(ctx, cpsReq)
	if err != nil {
		middleware.ErrorHandler(w, err)
		return
	}

	res := common.Response[*model.CpsAction]{
		ResponseWriter: w,
		Status:         http.StatusOK,
		Data:           rejectAction,
	}
	res.SendJSON()
}

func (b *BankAdapter) Disable(w http.ResponseWriter, r *http.Request) {
	id := chi.URLParam(r, "id")

	user_code := r.Context().Value(constant.ContextKey("user_code")).(string)
	full_name := r.Context().Value(constant.ContextKey("full_name")).(string)
	phone_number := r.Context().Value(constant.ContextKey("phone_number")).(string)
	department := r.Context().Value(constant.ContextKey("department")).(string)

	var cpsReq model.CreateCPSAction

	cpsReq.MakerUser = model.User{
		UserCode:    user_code,
		FullName:    full_name,
		PhoneNumber: phone_number,
	}
	cpsReq.Department = department

	ctx := r.Context()
	cpsAction, err := b.bankHandler.EnableOrDisableWallet(ctx, id, model.RequestDisableWallet, cpsReq)
	if err != nil {
		middleware.ErrorHandler(w, err)
		return
	}

	res := common.Response[*model.CpsAction]{
		ResponseWriter: w,
		Status:         http.StatusOK,
		Data:           cpsAction,
	}

	res.SendJSON()
}

func (b *BankAdapter) Enable(w http.ResponseWriter, r *http.Request) {
	id := chi.URLParam(r, "id")

	user_code := r.Context().Value(constant.ContextKey("user_code")).(string)
	full_name := r.Context().Value(constant.ContextKey("full_name")).(string)
	phone_number := r.Context().Value(constant.ContextKey("phone_number")).(string)
	department := r.Context().Value(constant.ContextKey("department")).(string)

	var cpsReq model.CreateCPSAction

	cpsReq.MakerUser = model.User{
		UserCode:    user_code,
		FullName:    full_name,
		PhoneNumber: phone_number,
	}
	cpsReq.Department = department

	ctx := r.Context()
	cpsAction, err := b.bankHandler.EnableOrDisableWallet(ctx, id, model.RequestEnableWallet, cpsReq)
	if err != nil {
		middleware.ErrorHandler(w, err)
		return
	}

	res := common.Response[*model.CpsAction]{
		ResponseWriter: w,
		Status:         http.StatusOK,
		Data:           cpsAction,
	}

	res.SendJSON()
}
