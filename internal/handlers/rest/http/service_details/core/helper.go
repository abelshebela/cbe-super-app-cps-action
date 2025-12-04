package core

import (
	"cbe-super-app-cps-action/internal/constants/localization"
	"encoding/json"
	"net/http"

	"gitlab.com/bersufekadgetachew/cbe-super-app-shared/shared/utils"
)

func DecodeJSONBody(w http.ResponseWriter, r *http.Request, dst interface{}, logger utils.Logger) bool {
	if err := json.NewDecoder(r.Body).Decode(dst); err != nil {
		logger.Errorf("Failed to decode JSON request: %v", err)
		localization.SendErrorResponse(w, localization.ErrorInvalidJSONPayload, nil, nil)
		return false
	}
	return true
}
