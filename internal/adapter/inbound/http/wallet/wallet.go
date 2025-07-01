package wallet

import (
	"encoding/json"
	"fmt"
	"net/http"
	"strconv"

	"cbe-super-app-cps-action/internal/adapter/outbound/model"
	"cbe-super-app-cps-action/internal/application/middleware"
	"cbe-super-app-cps-action/internal/application/wallet"
	"cbe-super-app-cps-action/internal/domain/wallet/dto"
	"cbe-super-app-cps-action/internal/domain/wallet/entity"
	inboundWallet "cbe-super-app-cps-action/internal/port/inbound/wallet"
	constant "cbe-super-app-cps-action/utils"

	"github.com/go-chi/chi/v5"
	"gitlab.com/bersufekadgetachew/cbe-super-app-shared/shared/common"
	"gitlab.com/bersufekadgetachew/cbe-super-app-shared/shared/utils"
)

type WalletAdapter struct {
	walletHandler wallet.WalletHandlerService
	logger        utils.Logger
}

func InitWalletAdapter(walletHandler wallet.WalletHandlerService, logger utils.Logger) inboundWallet.WalletAdapter {
	return &WalletAdapter{
		walletHandler: walletHandler,
		logger:        logger,
	}
}

func (wa *WalletAdapter) CreateWallet(w http.ResponseWriter, r *http.Request) {
	var walletRequest dto.CreateWalletRequest

	if err := r.ParseMultipartForm(10 << 20); err != nil {
		wa.logger.Errorf("failed to parse form data: %v", err)
		err = fmt.Errorf("failed to parse multipart form: %w", constant.ErrorDefinition{
			Code:    http.StatusBadRequest,
			Message: "invalid multipart form",
		})
		middleware.ErrorHandler(w, err)
		return
	}

	walletRequest.Name = r.FormValue("name")
	walletRequest.Code = r.FormValue("code")

	file, fileHeader, err := r.FormFile("avatar")
	if err != nil {
		wa.logger.Errorf("logo error: %v", err)
		err = fmt.Errorf("failed to read logo: %w", constant.ErrorDefinition{
			Code:    http.StatusBadRequest,
			Message: "missing or invalid avatar",
		})
		middleware.ErrorHandler(w, err)
		return
	}
	defer file.Close()

	walletRequest.Avatar = fileHeader

	var cpsRequest model.CreateCPSAction
	user_code := r.Context().Value(constant.ContextKey("user_code")).(string)
	full_name := r.Context().Value(constant.ContextKey("full_name")).(string)
	phone_number := r.Context().Value(constant.ContextKey("phone_number")).(string)
	department := r.Context().Value(constant.ContextKey("department")).(string)

	cpsRequest.ActionData = walletRequest
	cpsRequest.MakerUser = model.User{
		UserCode:    user_code,
		FullName:    full_name,
		PhoneNumber: phone_number,
	}
	cpsRequest.Department = department

	ctx := r.Context()
	cpsRes, err := wa.walletHandler.CreateWallet(ctx, cpsRequest)
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

func (wa *WalletAdapter) UpdateWallet(w http.ResponseWriter, r *http.Request) {
	id := chi.URLParam(r, "id")

	var updateRequest dto.UpdateWalletRequest

	if err := json.NewDecoder(r.Body).Decode(&updateRequest); err != nil {
		wa.logger.Errorf("invalid input", err)
		err = fmt.Errorf("failed to decode update wallet request error data %w", constant.ErrorDefinition{
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
	cpsAction, err := wa.walletHandler.UpdateWallet(ctx, id, cpsReq)
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

func (wa *WalletAdapter) DeleteWallet(w http.ResponseWriter, r *http.Request) {
	id := chi.URLParam(r, "id")

	var cpsReq model.CreateCPSAction

	user_code := r.Context().Value(constant.ContextKey("user_code")).(string)
	full_name := r.Context().Value(constant.ContextKey("full_name")).(string)
	phone_number := r.Context().Value(constant.ContextKey("phone_number")).(string)
	department := r.Context().Value(constant.ContextKey("department")).(string)

	cpsReq.ActionData = entity.Wallet{
		ID: id,
	}
	cpsReq.MakerUser = model.User{
		UserCode:    user_code,
		FullName:    full_name,
		PhoneNumber: phone_number,
	}
	cpsReq.Department = department

	ctx := r.Context()
	cpsAction, err := wa.walletHandler.DeleteWallet(ctx, id, cpsReq)
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

func (wa *WalletAdapter) GetAllWallet(w http.ResponseWriter, r *http.Request) {
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

	banks, err := wa.walletHandler.GetAllWallet(ctx, filterParams)
	if err != nil {
		middleware.ErrorHandler(w, err)
		return
	}

	res := common.Response[*entity.WalletResponse]{
		ResponseWriter: w,
		Status:         http.StatusOK,
		Data:           banks,
	}
	res.SendJSON()
}

func (wa *WalletAdapter) GetWallet(w http.ResponseWriter, r *http.Request) {
	id := chi.URLParam(r, "id")

	ctx := r.Context()

	bank, err := wa.walletHandler.GetWallet(ctx, id)
	if err != nil {
		middleware.ErrorHandler(w, err)
		return
	}

	res := common.Response[*entity.Wallet]{
		ResponseWriter: w,
		Status:         http.StatusOK,
		Data:           bank,
	}

	res.SendJSON()
}

func (wa *WalletAdapter) Authorize(w http.ResponseWriter, r *http.Request) {
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
	authAction, err := wa.walletHandler.Authorize(ctx, cpsReq)
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

func (wa *WalletAdapter) Reject(w http.ResponseWriter, r *http.Request) {
	action_code := chi.URLParam(r, "action_code")

	var cpsReq model.RejectCPSAction

	if err := json.NewDecoder(r.Body).Decode(&cpsReq); err != nil {
		wa.logger.Errorf("failed to decode wallet request", err)
		err = fmt.Errorf("failed to decode wallet request error data %w", constant.ErrorDefinition{
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
	rejectAction, err := wa.walletHandler.Reject(ctx, cpsReq)
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

func (wa *WalletAdapter) Disable(w http.ResponseWriter, r *http.Request) {
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
	cpsReq.ActionData = entity.Wallet{
		ID: id,
	}

	ctx := r.Context()
	cpsAction, err := wa.walletHandler.EnableOrDisableWallet(ctx, id, model.RequestDisableWallet, cpsReq)
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

func (wa *WalletAdapter) Enable(w http.ResponseWriter, r *http.Request) {
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
	cpsReq.ActionData = entity.Wallet{
		ID: id,
	}

	ctx := r.Context()
	cpsAction, err := wa.walletHandler.EnableOrDisableWallet(ctx, id, model.RequestEnableWallet, cpsReq)
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
