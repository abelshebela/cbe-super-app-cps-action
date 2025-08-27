package amount_based_auth

import (
	"encoding/json"
	"net/http"
	"time"

	amount_based "cbe-super-app-cps-action/internal/constants/interfaces/amount_based_auth"
	"cbe-super-app-cps-action/internal/constants/localization"
	"cbe-super-app-cps-action/internal/constants/model"
	"cbe-super-app-cps-action/internal/service"
	common_util "cbe-super-app-cps-action/pkgs/utils"

	"gitlab.com/bersufekadgetachew/cbe-super-app-shared/shared/utils"
	"go.mongodb.org/mongo-driver/v2/bson"
)

type AmountBasedAuthHandler struct {
	Service service.AmountBasedAuthService
	logger  utils.Logger
}

func NewAmountBasedAuthHandler(service service.AmountBasedAuthService, logger utils.Logger) amount_based.AmountBasedAuthAdapter {
	return &AmountBasedAuthHandler{
		Service: service,
		logger:  logger,
	}
}

func (a *AmountBasedAuthHandler) GetAllAmountBasedAuth(w http.ResponseWriter, r *http.Request) {
	filterParams := common_util.ExtractFilterParams(r)

	customers, err := a.Service.FindAllWithPagination(r.Context(), *filterParams)
	if err != nil {
		localization.SendErrorByCodeResponse(w, err.Error())
		return
	}

	localization.SendSuccessResponse(w, localization.SuccessUserRetrieved, customers)
}

func (a *AmountBasedAuthHandler) UpdateAmountBasedAuth(w http.ResponseWriter, r *http.Request) {
	id, ok := common_util.GetParam(r, "id")
	if !ok {
		localization.SendBadRequestResponse(w, localization.ErrorInvalidInputParameters.Code)
		return
	}

	var request model.AuthTier
	if err := json.NewDecoder(r.Body).Decode(&request); err != nil {
		localization.SendBadRequestResponse(w, localization.ErrorUnexpectedError.Code)
		return
	}

	// Set the ID from URL parameter
	objID, err := bson.ObjectIDFromHex(id)
	if err != nil {
		localization.SendBadRequestResponse(w, localization.ErrorInvalidID.Code)
		return
	}
	request.ID = objID
	request.LastModified = time.Now()

	// Update the auth tier through the service
	err = a.Service.Update(r.Context(), id, &request)
	if err != nil {
		localization.SendErrorByCodeResponse(w, err.Error())
		return
	}

	localization.SendSuccessResponse(w, localization.SuccessUserUpdated, nil)
}

func (a *AmountBasedAuthHandler) RejectAmountBasedAuth(w http.ResponseWriter, r *http.Request) {
	_, ok := common_util.GetParam(r, "id")
	if !ok {
		localization.SendBadRequestResponse(w, localization.ErrorInvalidInputParameters.Code)
		return
	}

	var cpsReq model.CPSAction
	if err := json.NewDecoder(r.Body).Decode(&cpsReq); err != nil {
		localization.SendBadRequestResponse(w, localization.ErrorUnexpectedError.Code)
		return
	}

	// For rejection, we just return success since the actual rejection
	// would be handled by the CPS action system
	localization.SendSuccessResponse(w, localization.SuccessUserUpdated, cpsReq)
} 