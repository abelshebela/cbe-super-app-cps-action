package faydaaccount

import (
	"encoding/json"
	"fmt"
	"net/http"

	faydaaccount "github.com/CBE-Super-App/cbe-super-app-cps-action/internal/application/fayda_account"
	"github.com/CBE-Super-App/cbe-super-app-cps-action/internal/application/middleware"
	"github.com/CBE-Super-App/cbe-super-app-cps-action/internal/port/inbound"

	entities "github.com/CBE-Super-App/cbe-super-app-cps-action/internal/domain/cps_actions/entities"
	ctx_util "github.com/CBE-Super-App/cbe-super-app-cps-action/pkgs/context"
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

	var cpsAction entities.CPSAction

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

func (f FaydaAccountAdapter) GetAllFaydaAccounts(w http.ResponseWriter, r *http.Request) {
	filterParams := constant_util.ExtractFilterParams(r)

	accounts, err := f.FaydaAccountHandler.GetAllFaydaAccounts(r.Context(), filterParams)
	if err != nil {
		f.logger.Errorf("GetAllFaydaAccounts request failed: %v", err)
		constant_util.SendErrorResponse(w, err.Error(), 0, nil)
		return
	}

	data, _ := constant_util.StructToMap(accounts)
	constant_util.BaseResponseMaker(data, w, "Fayda accounts fetched successfully", 200)
}
