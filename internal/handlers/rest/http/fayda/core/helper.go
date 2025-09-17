package core

import (
	"cbe-super-app-cps-action/internal/constants/localization"
	local_util "cbe-super-app-cps-action/pkgs/utils"
	"net/http"

	"gitlab.com/bersufekadgetachew/cbe-super-app-shared/shared/utils"
)

func ExtractID(w http.ResponseWriter, r *http.Request, logger utils.Logger) (string, bool) {
	user_code, ok := local_util.GetParam(r, "user_code")
	if !ok {
		logger.Errorf("[fayda.extractID] missing or invalid parameter 'id'")
		localization.SendBadRequestResponse(w, localization.ErrorInvalidInputParameters.Code)
		return "", false
	}
	return user_code, true
}
