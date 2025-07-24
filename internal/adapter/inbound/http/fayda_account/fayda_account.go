package faydaaccount

import (
	"context"
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

func (f *FaydaAccountAdapter) handleFaydaAccountAction(w http.ResponseWriter, r *http.Request, actionData interface{}, handlerFunc func(context.Context, entities.CPSAction) error,
) {
	userCode, ok := common_util.GetParam(r, "user_code")
	if !ok {
		f.logger.Errorf("missing or invalid parameter 'user_code'")
		common_util.SendErrorResponse(w, common_util.InvalidInputParameters, 0, nil)
		return
	}

	if err := json.NewDecoder(r.Body).Decode(actionData); err != nil {
		f.logger.Errorf("failed to bind action data: %v", err)
		constant_util.SendErrorResponse(w, constant_util.InvalidInput, 0, nil)
		return
	}

	userContext := ctx_util.ExtractUserContext(r)
	if userContext.IsIncomplete() {
		f.logger.Errorf("incomplete user context")
		constant_util.SendErrorResponse(w, constant_util.IncompleteUserInfo, 0, nil)
		return
	}

	switch v := actionData.(type) {
	case *faydaaccount.ActionDisableData:
		v.UseCode = userCode
	case *faydaaccount.ActionEnableData:
		v.UseCode = userCode
	}

	cpsAction := entities.CPSAction{
		MakerID:          userContext.UserID,
		MakerName:        userContext.FullName,
		MakerPhoneNumber: userContext.PhoneNumber,
		Department:       userContext.Department,
		CurrentAction:    actionData,
	}

	err := handlerFunc(r.Context(), cpsAction)
	if err != nil {
		constant_util.SendErrorResponse(w, err.Error(), 500, nil)
		return
	}

	constant_util.WriteSuccessResponse(w, nil, "Successfully Fayda Request initiated")
}

func (f *FaydaAccountAdapter) InitiateDisableFaydaAccount(w http.ResponseWriter, r *http.Request) {
	f.handleFaydaAccountAction(
		w,
		r,
		&faydaaccount.ActionDisableData{},
		f.FaydaAccountHandler.InitiateDisableFaydaAccount,
	)
}

func (f *FaydaAccountAdapter) InitiateEnableFaydaAccount(w http.ResponseWriter, r *http.Request) {
	f.handleFaydaAccountAction(
		w,
		r,
		&faydaaccount.ActionEnableData{},
		f.FaydaAccountHandler.InitiateEnableFaydaAccount,
	)
}
