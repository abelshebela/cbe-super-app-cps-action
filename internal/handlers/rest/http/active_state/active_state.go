package active_state

import (
	"cbe-super-app-cps-action/internal/constants/interfaces/active_state"
	"cbe-super-app-cps-action/internal/constants/localization"

	"net/http"
)

type ActiveStateHandler struct{}

func NewActiveState() active_state.ActiveState {
	return &ActiveStateHandler{}
}

func (a *ActiveStateHandler) UpdateStatus(w http.ResponseWriter, r *http.Request) {

	localization.SendSuccessResponse(w, localization.SuccessActiveStatusUpdated, nil)

}
