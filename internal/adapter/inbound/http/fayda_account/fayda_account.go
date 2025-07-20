package faydaaccount

import (
	"encoding/json"
	"fmt"
	"net/http"

	faydaaccount "github.com/CBE-Super-App/cbe-super-app-cps-action/internal/application/fayda_account"
	"github.com/CBE-Super-App/cbe-super-app-cps-action/internal/application/middleware"
	"github.com/CBE-Super-App/cbe-super-app-cps-action/internal/port/inbound"
	"github.com/go-chi/chi/v5"

	constant_util "github.com/CBE-Super-App/cbe-super-app-cps-action/pkgs/utils"
	constant "github.com/CBE-Super-App/cbe-super-app-cps-action/utils"

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

	fmt.Println("department--------", department)
	cpsAction.MakerID = user_code
	cpsAction.MakerName = full_name
	cpsAction.MakerPhoneNumber = phone_number

	cpsAction.CurrentAction = req
	cpsAction.Department = department

	ctx := r.Context()
	disableFaydaRes, err := f.FaydaAccountHandler.InitiateDisableFaydaAccount(ctx, cpsAction)
	if err != nil {
		constant_util.SendErrorResponse(w, err.Error(), 500, nil)
		return
	}

	data, err := constant_util.StructToMap(disableFaydaRes)
	if err != nil {
		constant_util.SendErrorResponse(w, err.Error(), 500, nil)
		return
	}
	constant_util.BaseResponseMaker(data, w, "Successfuly Fayda Request initiated", 200)
}

func (f FaydaAccountAdapter) AuthorizeFaydaAccountDisable(w http.ResponseWriter, r *http.Request) {
	action_code := chi.URLParam(r, "action_code")

	if action_code == "" {
		constant_util.SendErrorResponse(w, fmt.Errorf("UNHANDLED_SERVER_ERROR"), 500, nil)
		return
	}

	var cpsAction faydaaccount.CPSAction

	user_code := r.Context().Value(constant.ContextKey("user_code")).(string)
	full_name := r.Context().Value(constant.ContextKey("full_name")).(string)
	phone_number := r.Context().Value(constant.ContextKey("phone_number")).(string)
	department := r.Context().Value(constant.ContextKey("department")).(string)

	cpsAction.CheckerID = user_code
	cpsAction.CheckerName = full_name
	cpsAction.CheckerPhoneNumber = phone_number
	cpsAction.ActionCode = action_code
	cpsAction.Department = department

	ctx := r.Context()
	authorizeFayda, err := f.FaydaAccountHandler.AuthorizeFaydaAccountDisable(ctx, cpsAction)
	if err != nil {
		f.logger.Errorf("failed to bind action data", err)
		constant_util.SendErrorResponse(w, err.Error(), 500, nil)
		return
	}

	data, err := constant_util.StructToMap(authorizeFayda)
	if err != nil {
		constant_util.SendErrorResponse(w, err.Error(), 500, nil)
		return
	}

	constant_util.BaseResponseMaker(data, w, "Successfuly Approved", 200)
}

func (f FaydaAccountAdapter) RejectFaydaAccountDisable(w http.ResponseWriter, r *http.Request) {
	var req faydaaccount.RejectCPSAction
	action_code := chi.URLParam(r, "action_code")
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

	cpsAction.CheckerID = user_code
	cpsAction.ActionCode = action_code
	cpsAction.CheckerName = full_name
	cpsAction.CheckerPhoneNumber = phone_number

	cpsAction.RejectionReason = &req.RejectedReason
	cpsAction.Department = department

	_, err := f.FaydaAccountHandler.RejectFaydaAccountDisable(r.Context(), cpsAction)
	if err != nil {
		constant_util.SendErrorResponse(w, err.Error(), 500, nil)
		return
	}

	constant_util.BaseResponseMaker(nil, w, "Successfuly action rejected", 200)

}
