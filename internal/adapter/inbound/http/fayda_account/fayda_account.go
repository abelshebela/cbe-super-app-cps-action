package faydaaccount

import (
	// "context"
	// "encoding/json"
	"net/http"

	// "github.com/CBE-Super-App/cbe-super-app-cps-action/internal/adapter/outbound/model"
	faydaaccount "github.com/CBE-Super-App/cbe-super-app-cps-action/internal/application/fayda_account"
	"github.com/CBE-Super-App/cbe-super-app-cps-action/internal/port/inbound"

	// entities "github.com/CBE-Super-App/cbe-super-app-cps-action/internal/domain/cps_actions/entities"
	ctx_util "github.com/CBE-Super-App/cbe-super-app-cps-action/pkgs/context"
	common_util "github.com/CBE-Super-App/cbe-super-app-cps-action/pkgs/utils"
	constant_util "github.com/CBE-Super-App/cbe-super-app-cps-action/pkgs/utils"

	// ctx_util "github.com/CBE-Super-App/cbe-super-app-cps-action/pkgs/context"
	entities "github.com/CBE-Super-App/cbe-super-app-cps-action/internal/domain/cps_actions/entities"
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

func (f *FaydaAccountAdapter) InitiateEnableFaydaAccount(w http.ResponseWriter, r *http.Request) {

	userID, ok := common_util.GetParam(r, "user_code")
	if !ok {
		f.logger.Errorf("missing or invalid parameter 'user_code'")
		common_util.SendErrorResponse(w, common_util.InvalidInputParameters, 0, nil)
		return
	}
	MakerData := ctx_util.ExtractUserContext(r)

	cpsAction := entities.CPSAction{
		MakerID:          MakerData.UserID,
		MakerName:        MakerData.FullName,
		MakerPhoneNumber: MakerData.PhoneNumber,
		Department:       MakerData.Department,
	}

	actionCode, err := f.FaydaAccountHandler.InitiateEnableFaydaAccount(r.Context(), cpsAction, userID)
	if err != nil {
		constant_util.SendErrorResponse(w, err.Error(), 500, nil)
		return
	}

	constant_util.WriteSuccessResponse(w, actionCode, "Successfully Fayda Enable initiated")
}

func (f *FaydaAccountAdapter) InitiateDisableFaydaAccount(w http.ResponseWriter, r *http.Request) {
	userID, ok := common_util.GetParam(r, "user_code")
	if !ok {
		f.logger.Errorf("missing or invalid parameter 'user_ID'")
		common_util.SendErrorResponse(w, common_util.InvalidInputParameters, 0, nil)
		return
	}
	MakerData := ctx_util.ExtractUserContext(r)

	cpsAction := entities.CPSAction{
		MakerID:          MakerData.UserID,
		MakerName:        MakerData.FullName,
		MakerPhoneNumber: MakerData.PhoneNumber,
		Department:       MakerData.Department,
	}

	actionCode, err := f.FaydaAccountHandler.InitiateDisableFaydaAccount(r.Context(), cpsAction, userID)
	if err != nil {
		constant_util.SendErrorResponse(w, err.Error(), 500, nil)
		return
	}

	constant_util.WriteSuccessResponse(w, actionCode, "Successfully Fayda Enable initiated")
}
