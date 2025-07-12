package responseutil

import (
	"net/http"

	common_util "github.com/CBE-Super-App/cbe-super-app-cps-action/pkgs/utils"
)

type ErrorResponse struct {
	Status  int         `json:"status"`
	Message string      `json:"message"`
	Data    interface{} `json:"data"`
}

func NotFoundHandler(w http.ResponseWriter, r *http.Request) {
	common_util.SendErrorResponse(w, common_util.PageNotFound, http.StatusNotFound, nil)
}
