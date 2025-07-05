package amount_based_auth_handler

import (
	"encoding/json"
	"fmt"

	"github.com/CBE-Super-App/cbe-super-app-cps-action/internal/adapter/outbound/model"
	"github.com/CBE-Super-App/cbe-super-app-cps-action/internal/application/amount_based_auth_app"
	"github.com/CBE-Super-App/cbe-super-app-cps-action/internal/application/middleware"
	amount_based_auth_domain "github.com/CBE-Super-App/cbe-super-app-cps-action/internal/domain/amount_based_auth"
	"github.com/CBE-Super-App/cbe-super-app-cps-action/internal/port/inbound"
	ctx_util "github.com/CBE-Super-App/cbe-super-app-cps-action/pkgs/context"
	common_util "github.com/CBE-Super-App/cbe-super-app-cps-action/pkgs/utils"
	constant "github.com/CBE-Super-App/cbe-super-app-cps-action/utils"

	"github.com/go-chi/chi/v5"

	"net/http"

	"gitlab.com/bersufekadgetachew/cbe-super-app-shared/shared/common"
	"gitlab.com/bersufekadgetachew/cbe-super-app-shared/shared/utils"
)

type AmountBasedAuthHandler struct {
	amountBasedAuthService amount_based_auth_app.ApplicationService
	logger                 utils.Logger
}

type Resp struct {
	message string
}

type CurrentUser struct {
	Department  string
	UserCode    string
	FullName    string
	PhoneNumber string
}

func NewAmountBasedAuthHandler(service amount_based_auth_app.ApplicationService, logger utils.Logger) inbound.AmountBasedAuthHandler {
	return &AmountBasedAuthHandler{
		amountBasedAuthService: service,
		logger:                 logger,
	}
}

func (a *AmountBasedAuthHandler) buildUserContext(r *http.Request) (*CurrentUser, error) {
	userContext := ctx_util.ExtractUserContext(r)
	if userContext.IsIncomplete() {
		return nil, fmt.Errorf(common_util.IncompleteUserInfo)
	}

	return &CurrentUser{
		Department:  userContext.Department,
		UserCode:    userContext.UserCode,
		FullName:    userContext.FullName,
		PhoneNumber: userContext.PhoneNumber,
	}, nil
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
	curUser, curErr := a.buildUserContext(r)

	if curErr != nil {
		middleware.ErrorHandler(w, curErr)
		return
	}

	cpsActionRequest.MakerUser = model.User{
		UserCode:    curUser.UserCode,
		FullName:    curUser.FullName,
		PhoneNumber: curUser.PhoneNumber,
	}
	cpsActionRequest.Department = curUser.Department
	request.Id = id
	cpsActionRequest.CurrentData = request

	ctx := r.Context()
	amountBasedAuth, err := a.amountBasedAuthService.UpdateAmountBasedAuth(ctx, request, cpsActionRequest)
	if err != nil {
		middleware.ErrorHandler(w, err)
		return
	}

	response := common.Response[*model.CPSAction]{
		ResponseWriter: w,
		Status:         http.StatusOK,
		Data:           amountBasedAuth,
	}

	response.SendJSON()

}

func (a *AmountBasedAuthHandler) ApproveAmountBasedAuth(w http.ResponseWriter, r *http.Request) {
	id := chi.URLParam(r, "id")

	ctx := r.Context()

	var cpsReq model.AuthorizeCPSAction

	curUser, curErr := a.buildUserContext(r)

	if curErr != nil {
		middleware.ErrorHandler(w, curErr)
		return
	}

	cpsReq.CheckerUser = model.User{
		UserCode:    curUser.UserCode,
		FullName:    curUser.FullName,
		PhoneNumber: curUser.PhoneNumber,
	}
	cpsReq.Department = curUser.Department

	amountBasedAuth, err := a.amountBasedAuthService.ApproveAmountBasedAuth(ctx, id, cpsReq)
	if err != nil {
		middleware.ErrorHandler(w, err)
		return
	}

	response := common.Response[*model.CPSAction]{
		ResponseWriter: w,
		Status:         http.StatusOK,
		Data:           amountBasedAuth,
	}

	response.SendJSON()
}

func (a *AmountBasedAuthHandler) RejectAmountBasedAuth(w http.ResponseWriter, r *http.Request) {
	id := chi.URLParam(r, "id")

	var cpsReq model.RejectCPSAction

	if err := json.NewDecoder(r.Body).Decode(&cpsReq); err != nil {
		a.logger.Errorf("failed to decode amount based auth request", err)
		err = fmt.Errorf("%w", constant.ErrorDefinition{
			Code:    http.StatusBadRequest,
			Message: "invalid request",
		})
		middleware.ErrorHandler(w, err)
		return
	}

	curUser, curErr := a.buildUserContext(r)

	if curErr != nil {
		middleware.ErrorHandler(w, curErr)
		return
	}

	cpsReq.CheckerUser = model.User{
		UserCode:    curUser.UserCode,
		FullName:    curUser.FullName,
		PhoneNumber: curUser.PhoneNumber,
	}
	cpsReq.Department = curUser.Department

	ctx := r.Context()
	rejectAction, err := a.amountBasedAuthService.RejectAmountBasedAuth(ctx, id, cpsReq)
	if err != nil {
		middleware.ErrorHandler(w, err)
		return
	}

	res := common.Response[*model.CPSAction]{
		ResponseWriter: w,
		Status:         http.StatusOK,
		Data:           rejectAction,
	}
	res.SendJSON()
}
