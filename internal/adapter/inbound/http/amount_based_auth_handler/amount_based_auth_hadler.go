package amount_based_auth_handler

import (
	"encoding/json"
	"github.com/go-chi/chi/v5"
	"gitlab.com/bersufekadgetachew/cbe-super-app-cps-action/internal/application/amount_based_auth_app"
	amount_based_auth_domain "gitlab.com/bersufekadgetachew/cbe-super-app-cps-action/internal/domain/amount_based_auth"
	"gitlab.com/bersufekadgetachew/cbe-super-app-cps-action/internal/port/inbound"

	"gitlab.com/bersufekadgetachew/cbe-super-app-shared/shared/common"
	"gitlab.com/bersufekadgetachew/cbe-super-app-shared/shared/utils"
	"net/http"
)

type AmountBasedAuthHandler struct {
	amountBasedAuthService amount_based_auth_app.ApplicationService
	logger                 utils.Logger
}

type Resp struct {
	message string
}

func (a AmountBasedAuthHandler) ApproveAmountBasedAuth(w http.ResponseWriter, r *http.Request) {
	id := chi.URLParam(r, "id")
	amountBasedAuth, err := a.amountBasedAuthService.ApproveAmountBasedAuth(id)
	if err != nil {
		return
	}

	r2 := Resp{
		message: amountBasedAuth,
	}
	response := common.Response[any]{
		ResponseWriter: w,
		Status:         http.StatusOK,
		Data:           r2.message,
	}

	response.SendJSON()
}

func (a AmountBasedAuthHandler) UpdateAmountBasedAuth(w http.ResponseWriter, r *http.Request) {
	var request amount_based_auth_domain.AmountBasedAuthRequest
	err := json.NewDecoder(r.Body).Decode(&request)
	if err != nil {
		return
	}

	amountBasedAuth, err2 := a.amountBasedAuthService.UpdateAmountBasedAuth(request)
	if err2 != nil {
		return
	}

	response := common.Response[any]{
		ResponseWriter: w,
		Status:         http.StatusOK,
		Data:           amountBasedAuth,
	}

	response.SendJSON()

}

func NewAmountBasedAuthHandler(service amount_based_auth_app.ApplicationService) inbound.AmountBasedAuthHandler {
	return &AmountBasedAuthHandler{
		amountBasedAuthService: service,
	}
}
