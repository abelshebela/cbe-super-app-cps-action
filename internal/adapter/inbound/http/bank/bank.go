// Package bank provides HTTP handlers and adapters for bank-related operations in the CPS action service.
package bank

import (
	"encoding/json"
	"fmt"
	"net/http"
	"strconv"

	"github.com/CBE-Super-App/cbe-super-app-cps-action/internal/adapter/outbound/model"
	"github.com/CBE-Super-App/cbe-super-app-cps-action/internal/application/bank"
	"github.com/CBE-Super-App/cbe-super-app-cps-action/internal/domain/bank/dto"
	"github.com/CBE-Super-App/cbe-super-app-cps-action/internal/domain/bank/entity"
	inboundBank "github.com/CBE-Super-App/cbe-super-app-cps-action/internal/port/inbound/bank"
	ctx_util "github.com/CBE-Super-App/cbe-super-app-cps-action/pkgs/context"
	common_util "github.com/CBE-Super-App/cbe-super-app-cps-action/pkgs/utils"
	constant "github.com/CBE-Super-App/cbe-super-app-cps-action/utils"

	"github.com/go-chi/chi/v5"
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

func createCPSUserForAuthorize(r *http.Request) (*model.AuthorizeCPSAction, error) {
	userContext := ctx_util.ExtractUserContext(r)
	if userContext.IsIncomplete() {
		return nil, fmt.Errorf(common_util.IncompleteUserInfo)
	}
	return &model.AuthorizeCPSAction{
		CheckerUser: toModelUser(userContext),
		Department:  userContext.Department,
	}, nil
}

func createCPSUserForReject(r *http.Request) (*model.RejectCPSAction, error) {
	userContext := ctx_util.ExtractUserContext(r)

	if userContext.IsIncomplete() {
		return nil, fmt.Errorf(common_util.IncompleteUserInfo)
	}

	v := &model.RejectCPSAction{
		CheckerUser: toModelUser(userContext),
	}

	v.Department = userContext.Department
	return v, nil
}

func createCPSUserForCreate(r *http.Request) (*model.CreateCPSAction, error) {
	userContext := ctx_util.ExtractUserContext(r)

	if userContext.IsIncomplete() {
		return nil, fmt.Errorf(common_util.IncompleteUserInfo)
	}
	return &model.CreateCPSAction{
		MakerUser:  toModelUser(userContext),
		Department: userContext.Department,
	}, nil
}

func toModelUser(userContext ctx_util.UserContext) model.User {
	return model.User{
		UserCode:    userContext.UserCode,
		FullName:    userContext.FullName,
		PhoneNumber: userContext.PhoneNumber,
	}
}

func (b *BankAdapter) CreateOneBank(w http.ResponseWriter, r *http.Request) {
	var bankRequest dto.CreateBankRequest

	if err := r.ParseMultipartForm(10 << 20); err != nil {
		b.logger.Errorf("failed to parse form data: %v", err)
		common_util.SendErrorResponse(w, common_util.InvalidForm, 0, nil)
		return
	}

	bankRequest.Name = r.FormValue("name")
	bankRequest.Code = r.FormValue("code")
	bankRequest.BIC = r.FormValue("bic")

	file, fileHeader, err := r.FormFile("logo")
	if err != nil {
		b.logger.Errorf("logo error: %v", err)
		common_util.SendErrorResponse(w, common_util.MissingOrInvalidLogo, 0, nil)
		return
	}
	defer file.Close()
	bankRequest.Logo = fileHeader

	cpsRequest, err := createCPSUserForCreate(r)
	if err != nil {
		common_util.SendErrorResponse(w, err.Error(), 0, nil)
		return
	}
	cpsRequest.ActionData = bankRequest

	ctx := r.Context()
	cpsRes, err := b.bankHandler.CreateOneBank(ctx, *cpsRequest)
	if err != nil {
		common_util.SendErrorResponse(w, err.Error(), 0, nil)
		return
	}

	common_util.WriteSuccessResponse(w, cpsRes, "")
}

func (b *BankAdapter) UpdateOneBank(w http.ResponseWriter, r *http.Request) {
	id := chi.URLParam(r, "id")

	var updateRequest dto.UpdateBankRequest
	if err := json.NewDecoder(r.Body).Decode(&updateRequest); err != nil {
		b.logger.Errorf("invalid input", err)
		common_util.SendErrorResponse(w, common_util.InvalidReq, 0, nil)
		return
	}

	cpsReq, err := createCPSUserForCreate(r)
	if err != nil {
		common_util.SendErrorResponse(w, err.Error(), 0, nil)
		return
	}
	updateRequest.ID = id
	cpsReq.ActionData = updateRequest

	ctx := r.Context()
	cpsAction, err := b.bankHandler.UpdateOneBank(ctx, id, *cpsReq)
	if err != nil {
		common_util.SendErrorResponse(w, err.Error(), 0, nil)
		return
	}

	common_util.WriteSuccessResponse(w, cpsAction, "")
}

func (b *BankAdapter) DeleteOneBank(w http.ResponseWriter, r *http.Request) {
	id := chi.URLParam(r, "id")

	cpsReq, err := createCPSUserForCreate(r)
	if err != nil {
		common_util.SendErrorResponse(w, err.Error(), 0, nil)
		return
	}
	cpsReq.ActionData = entity.Bank{ID: id}

	ctx := r.Context()
	cpsAction, err := b.bankHandler.DeleteOneBank(ctx, id, *cpsReq)
	if err != nil {
		common_util.SendErrorResponse(w, err.Error(), 0, nil)
		return
	}

	common_util.WriteSuccessResponse(w, cpsAction, "")
}

func (b *BankAdapter) GetAllBank(w http.ResponseWriter, r *http.Request) {
	query := r.URL.Query()

	page := constant.DefaultPage
	if pageInt, err := strconv.Atoi(query.Get("page")); err == nil && pageInt > 0 {
		page = pageInt
	}

	perPage := constant.DefaultPerPage
	if perPageInt, err := strconv.Atoi(query.Get("per_page")); err == nil && perPageInt <= 10 && perPageInt > 0 {
		perPage = perPageInt
	}

	search := query.Get("search")
	filter := query.Get("filter")

	filterParams := &constant.Filter{
		Page:    page,
		PerPage: perPage,
		Search:  search,
		Filters: filter,
	}

	ctx := r.Context()
	banks, err := b.bankHandler.GetAllBank(ctx, filterParams)
	if err != nil {
		common_util.SendErrorResponse(w, err.Error(), 0, nil)
		return
	}

	common_util.WriteSuccessResponse(w, banks, "")
}

func (b *BankAdapter) GetOneBank(w http.ResponseWriter, r *http.Request) {
	id := chi.URLParam(r, "id")

	ctx := r.Context()
	bank, err := b.bankHandler.GetOneBank(ctx, id)
	if err != nil {
		common_util.SendErrorResponse(w, err.Error(), 0, nil)
		return
	}

	common_util.WriteSuccessResponse(w, bank, "")
}

func (b *BankAdapter) Authorize(w http.ResponseWriter, r *http.Request) {
	actionCode := chi.URLParam(r, "action_code")

	cpsReq, err := createCPSUserForAuthorize(r)
	if err != nil {
		common_util.SendErrorResponse(w, err.Error(), 0, nil)
		return
	}
	cpsReq.ActionCode = actionCode

	ctx := r.Context()
	authAction, err := b.bankHandler.Authorize(ctx, *cpsReq)
	if err != nil {
		common_util.SendErrorResponse(w, err.Error(), 0, nil)
		return
	}

	common_util.WriteSuccessResponse(w, authAction, "")
}

func (b *BankAdapter) Reject(w http.ResponseWriter, r *http.Request) {
	actionCode := chi.URLParam(r, "action_code")

	cpsReq, err := createCPSUserForReject(r)
	if err != nil {
		common_util.SendErrorResponse(w, err.Error(), 0, nil)
		return
	}

	var rejectPayload model.RejectCPSAction
	if err := json.NewDecoder(r.Body).Decode(&rejectPayload); err != nil {
		b.logger.Errorf("failed to decode rejection payload: %v", err)
		common_util.SendErrorResponse(w, common_util.InvalidReq, 0, nil)
		return
	}

	cpsReq.ActionCode = actionCode
	cpsReq.RejectedReason = rejectPayload.RejectedReason

	ctx := r.Context()
	rejectAction, err := b.bankHandler.Reject(ctx, *cpsReq)
	if err != nil {
		common_util.SendErrorResponse(w, err.Error(), 0, nil)
		return
	}

	common_util.WriteSuccessResponse(w, rejectAction, "")
}

func (b *BankAdapter) Disable(w http.ResponseWriter, r *http.Request) {
	id := chi.URLParam(r, "id")

	cpsReq, err := createCPSUserForCreate(r)
	if err != nil {
		common_util.SendErrorResponse(w, err.Error(), 0, nil)
		return
	}
	cpsReq.ActionData = dto.UpdateBankRequest{ID: id}

	ctx := r.Context()
	cpsAction, err := b.bankHandler.EnableOrDisableBank(ctx, id, model.RequestDisableBank, *cpsReq)
	if err != nil {
		common_util.SendErrorResponse(w, err.Error(), 0, nil)
		return
	}

	common_util.WriteSuccessResponse(w, cpsAction, "")
}

func (b *BankAdapter) Enable(w http.ResponseWriter, r *http.Request) {
	id := chi.URLParam(r, "id")

	cpsReq, err := createCPSUserForCreate(r)
	if err != nil {
		common_util.SendErrorResponse(w, err.Error(), 0, nil)
		return
	}
	cpsReq.ActionData = dto.UpdateBankRequest{ID: id}

	ctx := r.Context()
	cpsAction, err := b.bankHandler.EnableOrDisableBank(ctx, id, model.RequestEnableBank, *cpsReq)
	if err != nil {
		common_util.SendErrorResponse(w, err.Error(), 0, nil)
		return
	}

	common_util.WriteSuccessResponse(w, cpsAction, "")
}
