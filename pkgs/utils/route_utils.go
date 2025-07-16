package utils

import (
	"net/http"

	"github.com/go-chi/chi/v5"
	"gitlab.com/bersufekadgetachew/cbe-super-app-shared/shared/utils"
)

func GetParam(w http.ResponseWriter, r *http.Request, key string, logger utils.Logger) string {
	value := chi.URLParam(r, key)
	if value == "" {
		logger.Errorf("missing or invalid parameter '%s'", key)
		SendErrorResponse(w, InvalidInputParameters, 0, nil)
		return ""
	}
	return value
}
