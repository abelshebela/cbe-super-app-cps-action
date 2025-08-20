package unlink

import (
	"cbe-super-app-cps-action/internal/constants/dto"
	"cbe-super-app-cps-action/internal/constants/localization"
	"cbe-super-app-cps-action/internal/service"
	"encoding/json"
	"net/http"

	local_util "cbe-super-app-cps-action/pkgs/utils"

	"github.com/go-chi/chi/v5"
	"gitlab.com/bersufekadgetachew/cbe-super-app-shared/shared/utils"
)

type unlinkAdapter struct {
	logger    utils.Logger
	unlinkApp service.UnlinkService
}

func InitUnlinkAdapter(unlinkApp service.UnlinkService, logger utils.Logger) *unlinkAdapter {
	return &unlinkAdapter{
		logger:    logger,
		unlinkApp: unlinkApp,
	}
}

func (a *unlinkAdapter) GetArchivedUser(w http.ResponseWriter, r *http.Request) {
	filterParams := local_util.ExtractFilterParams(r)
	if filterParams.Page < 0 || filterParams.PerPage < 0 {
		localization.SendErrorResponse(w, localization.ErrorInvalidRequest, nil, nil)
		return
	}

	archivedUser, err := a.unlinkApp.GetAllArchivedUser(r.Context(), filterParams)
	if err != nil {
		localization.SendErrorByCodeResponse(w, err.Error())
		return
	}

	localization.SendSuccessResponse(w, localization.SuccessUserRetrieved, archivedUser)
}

func (a *unlinkAdapter) GetUserByAccount(w http.ResponseWriter, r *http.Request) {
	accountNo := chi.URLParam(r, "account_number")
	if accountNo == "" {
		localization.SendErrorResponse(w, localization.ErrorAccountNumberRequired, nil, nil)
		return
	}
	archivedUser, err := a.unlinkApp.GetUserByAccount(r.Context(), accountNo)
	if err != nil {
		localization.SendErrorByCodeResponse(w, err.Error())
		return
	}

	localization.SendSuccessResponse(w, localization.SuccessUserRetrieved, archivedUser)

}

func (a *unlinkAdapter) UnlinkUserCif(w http.ResponseWriter, r *http.Request) {
	var req dto.UnlinkUserRequest

	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		localization.SendBadRequestResponse(w, localization.MsgInvalidInput)
		return
	}

	if req.Validate() != nil {
		localization.SendBadRequestResponse(w, localization.MsgInvalidInput)
		return
	}

	if err := a.unlinkApp.UnlinkUserCif(r.Context(), req.UserCode); err != nil {
		localization.SendErrorByCodeResponse(w, err.Error())
		return
	}

	localization.SendSuccessResponse(w, localization.SuccessUnlinkCifRequestSent, nil)

}
