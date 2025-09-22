package unlink

import (
	inbound "cbe-super-app-cps-action/internal/constants/interfaces/unlink"
	"cbe-super-app-cps-action/internal/constants/localization"
	"cbe-super-app-cps-action/internal/service"
	"net/http"
	"strconv"

	local_util "cbe-super-app-cps-action/pkgs/utils"

	"github.com/go-chi/chi/v5"
	"gitlab.com/bersufekadgetachew/cbe-super-app-shared/shared/utils"
)

type unlinkAdapter struct {
	logger    utils.Logger
	unlinkApp service.UnlinkService
}

func InitUnlinkAdapter(unlinkApp service.UnlinkService, logger utils.Logger) inbound.UnlinkAdapter {
	return &unlinkAdapter{
		logger:    logger,
		unlinkApp: unlinkApp,
	}
}

func (a *unlinkAdapter) GetArchivedUser(w http.ResponseWriter, r *http.Request) {
	page, err := strconv.Atoi(chi.URLParam(r, "page"))
	per_page, pgerr := strconv.Atoi(chi.URLParam(r, "per_page"))
	if err != nil || pgerr != nil {
		localization.SendInternalServerErrorResponse(w, "Parsing Error Occured")
		return
	}
	if page < 0 || per_page < 0 {
		localization.SendErrorResponse(w, localization.ErrorInvalidPaginationParams, nil, nil)
		return
	}
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
	accNumber := chi.URLParam(r, "account_number")
	if accNumber == "" {
		localization.SendErrorResponse(w, localization.ErrorInvalidRequest, nil, nil)
		return
	}
	user, err := a.unlinkApp.GetUserByAccount(r.Context(), accNumber)
	if err != nil {
		localization.SendErrorByCodeResponse(w, err.Error())
		return
	}

	localization.SendSuccessResponse(w, localization.SuccessUserRetrieved, user)
}

func (a *unlinkAdapter) UnlinkUserCif(w http.ResponseWriter, r *http.Request) {

	userCode := chi.URLParam(r, "user_code")
	if userCode == "" {
		localization.SendErrorResponse(w, localization.ErrorInvalidRequest, nil, nil)
		return
	}
	err := a.unlinkApp.UnlinkUserCif(r.Context(), userCode)
	if err != nil {
		localization.SendErrorByCodeResponse(w, err.Error())
		return
	}
	localization.SendSuccessResponse(w, localization.SuccessUnlinkCifRequestSent, nil)
}
