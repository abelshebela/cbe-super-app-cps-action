package faydaaccount

import (
	"encoding/json"
	"fmt"
	"net/http"

	faydaaccount "cbe-super-app-cps-action/internal/application/fayda_account"
	"cbe-super-app-cps-action/internal/application/middleware"
	"cbe-super-app-cps-action/internal/port/inbound"

	constant "cbe-super-app-cps-action/utils"

	"gitlab.com/bersufekadgetachew/cbe-super-app-shared/shared/common"
	"gitlab.com/bersufekadgetachew/cbe-super-app-shared/shared/utils"
)

type FaydaAccountAdapter struct {
	FaydaAccountHandler faydaaccount.ApplicationService
	logger              utils.Logger
}

func InitFaydaAdapter(faydaAccountHandler faydaaccount.ApplicationService, logger utils.Logger) inbound.FaydaAccount {
	return FaydaAccountAdapter{
		FaydaAccountHandler: faydaAccountHandler,
		logger:              logger,
	}
}

func (f FaydaAccountAdapter) InitiateDisableFaydaAccount(w http.ResponseWriter, r *http.Request) {
	var req faydaaccount.ActionData

	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		f.logger.Errorf("failed to bind action data", err)
		err = fmt.Errorf("failed to bind error data %w", constant.ErrorDefinition{
			Code:    http.StatusBadRequest,
			Message: "invalid request",
		})
		middleware.ErrorHandler(w, err)
		return
	}

	var cpsAction faydaaccount.CPSAction

	user_code := r.Context().Value(constant.ContextKey("user_code")).(string)
	full_name := r.Context().Value(constant.ContextKey("full_name")).(string)
	phone_number := r.Context().Value(constant.ContextKey("phone_number")).(string)
	department := r.Context().Value(constant.ContextKey("department")).(string)

	cpsAction.MakerUser = faydaaccount.User{
		UserCode:    user_code,
		FullName:    full_name,
		PhoneNumber: phone_number,
	}

	cpsAction.ActionData = req
	cpsAction.Department = department

	ctx := r.Context()
	disableFaydaRes, err := f.FaydaAccountHandler.InitiateDisableFaydaAccount(ctx, cpsAction)
	if err != nil {
		middleware.ErrorHandler(w, err)
		return
	}

	res := common.Response[*faydaaccount.CPSAction]{
		ResponseWriter: w,
		Status:         http.StatusOK,
		Data:           disableFaydaRes,
	}

	res.SendJSON()
}

func (f FaydaAccountAdapter) AuthorizeFaydaAccountDisable(w http.ResponseWriter, r *http.Request) {
	var req faydaaccount.ActionData

	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		f.logger.Errorf("failed to bind action data", err)
		err = fmt.Errorf("failed to bind error data %w", constant.ErrorDefinition{
			Code:    http.StatusBadRequest,
			Message: "invalid request",
		})
		middleware.ErrorHandler(w, err)
		return
	}

	var cpsAction faydaaccount.CPSAction

	user_code := r.Context().Value(constant.ContextKey("user_code")).(string)
	full_name := r.Context().Value(constant.ContextKey("full_name")).(string)
	phone_number := r.Context().Value(constant.ContextKey("phone_number")).(string)
	department := r.Context().Value(constant.ContextKey("department")).(string)

	cpsAction.CheckerUser = faydaaccount.User{
		UserCode:    user_code,
		FullName:    full_name,
		PhoneNumber: phone_number,
	}

	cpsAction.ActionData = req
	cpsAction.Department = department

	ctx := r.Context()
	authorizeFayda, err := f.FaydaAccountHandler.AuthorizeFaydaAccountDisable(ctx, cpsAction)
	if err != nil {
		middleware.ErrorHandler(w, err)
		return
	}

	res := common.Response[*faydaaccount.CPSAction]{
		ResponseWriter: w,
		Status:         http.StatusOK,
		Data:           authorizeFayda,
	}

	res.SendJSON()
}

func (f FaydaAccountAdapter) RejectFaydaAccountDisable(w http.ResponseWriter, r *http.Request) {
	var req faydaaccount.RejectCPSAction

	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		f.logger.Errorf("failed to bind action data", err)
		err = fmt.Errorf("failed to bind error data %w", constant.ErrorDefinition{
			Code:    http.StatusBadRequest,
			Message: "invalid request",
		})
		middleware.ErrorHandler(w, err)
		return
	}

	var cpsAction faydaaccount.CPSAction

	user_code := r.Context().Value(constant.ContextKey("user_code")).(string)
	full_name := r.Context().Value(constant.ContextKey("full_name")).(string)
	phone_number := r.Context().Value(constant.ContextKey("phone_number")).(string)
	department := r.Context().Value(constant.ContextKey("department")).(string)

	cpsAction.CheckerUser = faydaaccount.User{
		UserCode:    user_code,
		FullName:    full_name,
		PhoneNumber: phone_number,
	}

	cpsAction.ActionData = faydaaccount.ActionData{
		UseCode:     req.UseCode,
		FullName:    req.FullName,
		PhoneNumber: req.PhoneNumber,
	}
	cpsAction.RejectedReason = req.RejectedReason
	cpsAction.Department = department

	ctx := r.Context()
	rejectFayda, err := f.FaydaAccountHandler.RejectFaydaAccountDisable(ctx, cpsAction)
	if err != nil {
		middleware.ErrorHandler(w, err)
		return
	}

	res := common.Response[*faydaaccount.CPSAction]{
		ResponseWriter: w,
		Status:         http.StatusOK,
		Data:           rejectFayda,
	}

	res.SendJSON()
}
