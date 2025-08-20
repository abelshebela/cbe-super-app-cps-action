package unlink

import (
	"cbe-super-app-cps-action/internal/constants/localization"
	"cbe-super-app-cps-action/internal/service"
	"net/http"

	local_util "cbe-super-app-cps-action/pkgs/utils"

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

	archivedUser, err := a.unlinkApp(r.Context(), filterParams)
	if err != nil {
		localization.SendErrorByCodeResponse(w, err.E)
	}
}

func (a *unlinkAdapter) GetUserByAccount(w http.ResponseWriter, r *http.Request) {

}

func (a *unlinkAdapter) UnlinkUserCif(w http.ResponseWriter, r *http.Request) {

}
