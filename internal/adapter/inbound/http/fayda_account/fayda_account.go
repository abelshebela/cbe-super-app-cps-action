package faydaaccount

import (
	"encoding/json"
	"net/http"

	faydaaccount "github.com/CBE-Super-App/cbe-super-app-cps-action/internal/application/fayda_account"
	"github.com/CBE-Super-App/cbe-super-app-cps-action/internal/port/inbound"

	entities "github.com/CBE-Super-App/cbe-super-app-cps-action/internal/domain/cps_actions/entities"
	ctx_util "github.com/CBE-Super-App/cbe-super-app-cps-action/pkgs/context"
	common_util "github.com/CBE-Super-App/cbe-super-app-cps-action/pkgs/utils"
	constant_util "github.com/CBE-Super-App/cbe-super-app-cps-action/pkgs/utils"

	"gitlab.com/bersufekadgetachew/cbe-super-app-shared/shared/utils"
)

type FaydaAccountAdapter struct {
	FaydaAccountHandler faydaaccount.ApplicationService
	logger              utils.Logger
}

func InitFaydaAdapter(faydaAccountHandler faydaaccount.ApplicationService, logger utils.Logger) inbound.FaydaAccount {
	return &FaydaAccountAdapter{
		FaydaAccountHandler: faydaAccountHandler,
		logger:              logger,
	}
}

func (f *FaydaAccountAdapter) InitiateDisableFaydaAccount(w http.ResponseWriter, r *http.Request) {
	var req faydaaccount.ActionDisableData

	userCode, ok := common_util.GetParam(r, "user_code")
	if !ok {
		f.logger.Errorf("missing or invalid parameter 'user_code'")
		common_util.SendErrorResponse(w, common_util.InvalidInputParameters, 0, nil)
		return
	}

	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		f.logger.Errorf("failed to bind action data", err)
		constant_util.SendErrorResponse(w, constant_util.InvalidInput, 0, nil)
		return
	}

	var cpsAction entities.CPSAction
	req.UseCode = userCode

	userContext := ctx_util.ExtractUserContext(r)
	if userContext.IsIncomplete() {
		f.logger.Errorf("incomplete user context", "request_id")
		constant_util.SendErrorResponse(w, constant_util.IncompleteUserInfo, 0, nil)
		return
	}

	cpsAction.MakerID = userContext.UserID
	cpsAction.MakerName = userContext.FullName
	cpsAction.MakerPhoneNumber = userContext.PhoneNumber

	cpsAction.CurrentAction = req
	cpsAction.Department = userContext.Department

	err := f.FaydaAccountHandler.InitiateDisableFaydaAccount(r.Context(), cpsAction)
	if err != nil {
		constant_util.SendErrorResponse(w, err.Error(), 500, nil)
		return
	}

	constant_util.WriteSuccessResponse(w, nil, "Successfuly Fayda Request initiated")
}

func (f *FaydaAccountAdapter) InitiateEnableFaydaAccount(w http.ResponseWriter, r *http.Request) {
	var req faydaaccount.ActionEnableData

	userCode, ok := common_util.GetParam(r, "user_code")
	if !ok {
		f.logger.Errorf("missing or invalid parameter 'user_code'")
		common_util.SendErrorResponse(w, common_util.InvalidInputParameters, 0, nil)
		return
	}

	var cpsAction entities.CPSAction
	req.UseCode = userCode

	userContext := ctx_util.ExtractUserContext(r)
	if userContext.IsIncomplete() {
		f.logger.Errorf("incomplete user context", "request_id")
		constant_util.SendErrorResponse(w, constant_util.IncompleteUserInfo, 0, nil)
		return
	}

	cpsAction.MakerID = userContext.UserID
	cpsAction.MakerName = userContext.FullName
	cpsAction.MakerPhoneNumber = userContext.PhoneNumber
	cpsAction.Department = userContext.Department

	cpsAction.CurrentAction = req
	err := f.FaydaAccountHandler.InitiateEnableFaydaAccount(r.Context(), cpsAction)
	if err != nil {
		constant_util.SendErrorResponse(w, err.Error(), 500, nil)
		return
	}

	constant_util.WriteSuccessResponse(w, nil, "Successfuly Fayda Request initiated")
}
