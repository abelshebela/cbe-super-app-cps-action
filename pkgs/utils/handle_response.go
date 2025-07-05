package utils

import (
	"net/http"

	"gitlab.com/bersufekadgetachew/cbe-super-app-shared/shared/common"
)

// ErrorResponse sends a standardized error response as JSON.
func ErrorResponse(w http.ResponseWriter, status int, code, message string) {
	res := common.Response[map[string]string]{
		ResponseWriter: w,
		Status:         status,
		Data: map[string]string{
			"code":    code,
			"message": message,
		},
	}
	res.SendJSON()
}
