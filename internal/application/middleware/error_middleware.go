package middleware

import (
	"gitlab.com/bersufekadgetachew/cbe-super-app-cps-ms/utils"
	"encoding/json"
	"errors"
	"net/http"

	"gitlab.com/bersufekadgetachew/cbe-super-app-shared/shared/common"
)

func ErrorHandler(w http.ResponseWriter, err error) {
	var errorResponse utils.ErrorDefinition

	err = errors.Unwrap(err)
	
	if err := json.Unmarshal([]byte(err.Error()), &errorResponse); err != nil {
		http.Error(w, "failed to unmarshal error", http.StatusInternalServerError)
		return
	}

	res := common.Response[utils.ErrorDefinition]{
		ResponseWriter: w,
		Status:         errorResponse.Code,
		Data:           errorResponse,
	}

	res.SendJSON()
}
