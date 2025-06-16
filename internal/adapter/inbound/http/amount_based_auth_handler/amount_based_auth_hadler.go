package amount_based_auth_handler

import (
	"encoding/json"
	"gitlab.com/bersufekadgetachew/cbe-super-app-cps-ms/internal/application/amount_based_auth_app"
	amount_based_auth_domain "gitlab.com/bersufekadgetachew/cbe-super-app-cps-ms/internal/domain/amount_based_auth"
	"gitlab.com/bersufekadgetachew/cbe-super-app-cps-ms/internal/port/inbound"
	"gitlab.com/bersufekadgetachew/cbe-super-app-shared/shared/common"
	"gitlab.com/bersufekadgetachew/cbe-super-app-shared/shared/utils"
	"net/http"
)

type AmountBasedAuthHandler struct {
	amountBasedAuthService amount_based_auth_app.ApplicationService
	logger                 utils.Logger
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

	response := common.Response[amount_based_auth_domain.AuthTier]{
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
