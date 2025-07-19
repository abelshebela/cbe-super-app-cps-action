package amount_based_auth_handler

import (
	"encoding/json"

	"github.com/CBE-Super-App/cbe-super-app-cps-action/internal/adapter/outbound/model"
	"github.com/CBE-Super-App/cbe-super-app-cps-action/internal/application/amount_based_auth_app"
	amount_based_auth_domain "github.com/CBE-Super-App/cbe-super-app-cps-action/internal/domain/amount_based_auth"
	"github.com/CBE-Super-App/cbe-super-app-cps-action/internal/port/inbound"
	common_util "github.com/CBE-Super-App/cbe-super-app-cps-action/pkgs/utils"

	"net/http"

	ctx_util "github.com/CBE-Super-App/cbe-super-app-cps-action/pkgs/context"

	"gitlab.com/bersufekadgetachew/cbe-super-app-shared/shared/utils"
)

type AmountBasedAuthHandler struct {
	amountBasedAuthService amount_based_auth_app.ApplicationService
	logger                 utils.Logger
}

func NewAmountBasedAuthHandler(service amount_based_auth_app.ApplicationService, logger utils.Logger) inbound.AmountBasedAuthHandler {
	return &AmountBasedAuthHandler{
		amountBasedAuthService: service,
		logger:                 logger,
	}
}

func (a AmountBasedAuthHandler) GetAllAmountBasedAuth(w http.ResponseWriter, r *http.Request) {
	filterParams := common_util.ExtractFilterParams(r)

	customers, err := a.amountBasedAuthService.GetAllAmountBasedDetail(r.Context(), filterParams)
	if err != nil {
		common_util.SendErrorResponse(w, err.Error(), 0, nil)
		return
	}

	data, _ := common_util.StructToMap(customers)
	common_util.BaseResponseMaker(data, w, "Successfuly fetched", http.StatusAccepted)

}
func (a *AmountBasedAuthHandler) UpdateAmountBasedAuth(w http.ResponseWriter, r *http.Request) {
	id, ok := common_util.GetParam(r, "id")
	if !ok {
		a.logger.Errorf("missing or invalid parameter 'id'")
		common_util.SendErrorResponse(w, common_util.InvalidInputParameters, 0, nil)
		return
	}

	var request amount_based_auth_domain.UpdateAmountBasedAuth
	if err := json.NewDecoder(r.Body).Decode(&request); err != nil {
		a.logger.Errorf("failed to decode request body: %v", err)
		common_util.SendErrorResponse(w, common_util.InvalidJSONPayload, 0, nil)
		return
	}

	var cpsActionRequest model.CreateCPSAction

	userContext := ctx_util.ExtractUserContext(r)
	if userContext.IsIncomplete() {
		a.logger.Errorf("incomplete user context: %v", userContext)
		common_util.SendErrorResponse(w, common_util.InvalidToken, 0, nil)
		return
	}

	cpsActionRequest.MakerUser = model.User{
		UserCode:    userContext.UserCode,
		FullName:    userContext.FullName,
		PhoneNumber: userContext.PhoneNumber,
	}
	request.Id = id
	cpsActionRequest.Department = userContext.Department
	cpsActionRequest.CurrentData = request

	amountBasedAuth, err := a.amountBasedAuthService.UpdateAmountBasedAuth(r.Context(), request, cpsActionRequest)
	if err != nil {
		a.logger.Errorf("failed to update amount based auth: %v", err)
		common_util.SendErrorResponse(w, err.Error(), 0, nil)
		return
	}

	data, err := common_util.StructToMap(amountBasedAuth)
	if err != nil {
		a.logger.Errorf("failed to convert amount based auth to map: %v", err)
		common_util.SendErrorResponse(w, common_util.UnhandledServerError, 0, nil)
		return
	}
	common_util.BaseResponseMaker(data, w, "Successfuly updated", 200)

}

// func (a *AmountBasedAuthHandler) ApproveAmountBasedAuth(w http.ResponseWriter, r *http.Request) {
// 	id, ok := common_util.GetParam(r, "id")
// 	if !ok {
// 		a.logger.Errorf("missing or invalid parameter 'id'")
// 		common_util.SendErrorResponse(w, common_util.InvalidInputParameters, 0, nil)
// 		return
// 	}

// 	var cpsReq model.AuthorizeCPSAction

// 	userContext := ctx_util.ExtractUserContext(r)
// 	if userContext.IsIncomplete() {
// 		a.logger.Errorf("incomplete user context: %v", userContext)
// 		common_util.SendErrorResponse(w, common_util.InvalidToken, 0, nil)
// 		return
// 	}

// 	cpsReq.CheckerUser = model.User{
// 		UserCode:    userContext.UserCode,
// 		FullName:    userContext.FullName,
// 		PhoneNumber: userContext.PhoneNumber,
// 	}
// 	cpsReq.Department = userContext.Department

// 	amountBasedAuth, err := a.amountBasedAuthService.ApproveAmountBasedAuth(r.Context(), id, cpsReq)
// 	if err != nil {
// 		common_util.SendErrorResponse(w, err.Error(), 0, nil)
// 		return
// 	}

// 	data, err := common_util.StructToMap(amountBasedAuth)
// 	if err != nil {
// 		common_util.SendErrorResponse(w, err, 500, nil)
// 		return
// 	}
// 	common_util.BaseResponseMaker(data, w, "Successfuly updated", 200)
// }

func (a *AmountBasedAuthHandler) RejectAmountBasedAuth(w http.ResponseWriter, r *http.Request) {
	id, ok := common_util.GetParam(r, "id")
	if !ok {
		a.logger.Errorf("missing or invalid parameter 'id'")
		common_util.SendErrorResponse(w, common_util.InvalidInputParameters, 0, nil)
		return
	}

	var cpsReq model.RejectAuthTierCPSAction

	if err := json.NewDecoder(r.Body).Decode(&cpsReq); err != nil {
		a.logger.Errorf("failed to decode amount based auth request", err)
		common_util.SendErrorResponse(w, common_util.InvalidJSONPayload, 0, nil)
		return
	}

	ctx := r.Context()
	rejectAction, err := a.amountBasedAuthService.RejectAmountBasedAuth(ctx, id, cpsReq)
	if err != nil {
		a.logger.Errorf("failed to reject amount based auth: %v", err)
		common_util.SendErrorResponse(w, err.Error(), 0, nil)
		return
	}

	data, err := common_util.StructToMap(rejectAction)
	if err != nil {
		common_util.SendErrorResponse(w, err, 500, nil)
		return
	}
	common_util.BaseResponseMaker(data, w, "Successfuly updated", 200)
}
