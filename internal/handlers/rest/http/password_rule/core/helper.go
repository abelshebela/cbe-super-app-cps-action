package core

import (
	"cbe-super-app-cps-action/internal/constants/localization"
	local_util "cbe-super-app-cps-action/pkgs/utils"
	"encoding/json"
	"net/http"

	"gitlab.com/bersufekadgetachew/cbe-super-app-shared/shared/utils"
)

func ExtractID(w http.ResponseWriter, r *http.Request, logger utils.Logger) (string, bool) {
	id, ok := local_util.GetParam(r, "id")
	if !ok {
		logger.Errorf("[passwordRule.extractID] missing or invalid parameter 'id'")
		localization.SendBadRequestResponse(w, localization.ErrorInvalidInputParameters.Code)
		return "", false
	}
	return id, true
}

func StructToMap(data interface{}) (map[string]interface{}, error) {
	bytes, err := json.Marshal(data)
	if err != nil {
		return nil, err
	}

	var result map[string]interface{}
	err = json.Unmarshal(bytes, &result)
	return result, err
}
