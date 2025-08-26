package faydaaccount

import (
	faydaInbound "cbe-super-app-cps-action/internal/constants/interfaces/fayda"
	"cbe-super-app-cps-action/internal/constants/localization"
	"cbe-super-app-cps-action/internal/handlers/rest/http/ad/core"
	"net/http"

	"cbe-super-app-cps-action/internal/service"

	"gitlab.com/bersufekadgetachew/cbe-super-app-shared/shared/utils"
)

type faydaAccountHandler struct {
	faydaService service.FaydaAccountService
	logger       utils.Logger
}

func InitFaydaHandler(faydaService service.FaydaAccountService, logger utils.Logger) faydaInbound.FaydaAccount {
	return &faydaAccountHandler{
		faydaService: faydaService,
		logger:       logger,
	}
}

func (f *faydaAccountHandler) InitiateEnableFaydaAccount(w http.ResponseWriter, r *http.Request) {
	user_code, ok := core.ExtractID(w, r, f.logger)
	if !ok {
		return
	}

	err := f.faydaService.EnableOrDisableFayda(r.Context(), user_code, true)
	if err != nil {
		localization.SendErrorByCodeResponse(w, err.Error())
		return
	}

	localization.SendSuccessResponse(w, localization.SuccessFaydaEnableActionCreated, nil)
}

func (f *faydaAccountHandler) InitiateDisableFaydaAccount(w http.ResponseWriter, r *http.Request) {
	user_code, ok := core.ExtractID(w, r, f.logger)
	if !ok {
		return
	}

	err := f.faydaService.EnableOrDisableFayda(r.Context(), user_code, false)
	if err != nil {
		localization.SendErrorByCodeResponse(w, err.Error())
		return
	}

	localization.SendSuccessResponse(w, localization.SuccessFaydaDisableActionCreated, nil)
}
