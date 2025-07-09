package service

import (
	"encoding/json"
	"errors"
	"net/http"
	"time"

	"github.com/CBE-Super-App/cbe-super-app-cps-action/internal/adapter/inbound/http/service/dto"
	serviceApp "github.com/CBE-Super-App/cbe-super-app-cps-action/internal/application/service"
	cpsuser "github.com/CBE-Super-App/cbe-super-app-cps-action/internal/domain/cps_user"
	"github.com/CBE-Super-App/cbe-super-app-cps-action/internal/port/inbound"

	ctx_util "github.com/CBE-Super-App/cbe-super-app-cps-action/pkgs/context"
	common_util "github.com/CBE-Super-App/cbe-super-app-cps-action/pkgs/utils"

	"gitlab.com/bersufekadgetachew/cbe-super-app-shared/shared/utils"
)

type serviceHandler struct {
	appService serviceApp.ServiceApplication
	logger     utils.Logger
}

func NewServiceHandler(domain serviceApp.ServiceApplication, logger utils.Logger) inbound.ServiceBound {
	return &serviceHandler{
		appService: domain,
		logger:     logger,
	}
}

func (s *serviceHandler) AuthorizeServiceFee(w http.ResponseWriter, r *http.Request) {

}
func (s *serviceHandler) RejectServiceFee(w http.ResponseWriter, r *http.Request) {

}

func (s *serviceHandler) buildServiceFeeActionData(req dto.UpdateServiceFeeRequest, serviceID, serviceName string) map[string]any {
	return map[string]any{
		"aboveServiceFee":   req.AboveServiceFee / 100,
		"minAmountVIRTUAL":  req.MinAmountVIRTUAL,
		"dailyCapLevelOne":  req.DailyCapLevelOne,
		"singleCapLevelOne": req.SingleCapLevelOne,
		"service_id":        serviceID,
		"service_name":      serviceName,
		"tiers":             req.Tiers,
		"paymentType":       req.PaymentType,
		"serviceType":       req.ServiceType,
		"lastModified":      time.Now(),
	}
}

func (s *serviceHandler) extractCPSUserFromRequest(r *http.Request) (cpsuser.CPSUser, error) {
	userCtx := ctx_util.ExtractUserContext(r)
	if userCtx.IsIncomplete() {
		return cpsuser.CPSUser{}, errors.New(common_util.Unauthorized)
	}
	return cpsuser.CPSUser{
		ID:          userCtx.UserID,
		FullName:    userCtx.FullName,
		PhoneNumber: userCtx.PhoneNumber,
		Department:  userCtx.Department,
	}, nil
}

func (s *serviceHandler) UpdateServiceFee(w http.ResponseWriter, r *http.Request) {
	var req dto.UpdateServiceFeeRequest

	queryID := r.Header.Get("query_id")
	user, err := s.extractCPSUserFromRequest(r)
	if err != nil {
		common_util.SendErrorResponse(w, err.Error(), http.StatusUnauthorized, nil)
		return
	}

	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		common_util.SendErrorResponse(w, common_util.InvalidJSONPayload, 0, nil)
		return
	}

	if err := s.appService.ValidateTiers(req.Tiers, req.AboveAmount); err != nil {
		common_util.SendErrorResponse(w, err.Error(), 0, nil)
		return
	}

	service, err := s.appService.GetOneService(r.Context(), queryID)
	if err != nil {
		common_util.SendErrorResponse(w, common_util.ServiceNotFound, 0, nil)
		return
	}

	// Prepare CPS action data
	actionData := s.buildServiceFeeActionData(req, queryID, service.ServiceName)
	initAction := s.appService.InitCPSAction(user, actionData, "CREATE", "UPDATE SERVICE FEE", service)

	err = s.appService.CreateAction(r.Context(), initAction)
	if err != nil {
		common_util.SendErrorResponse(w, common_util.FailedToCreateAction, http.StatusConflict, nil)
		return
	}

	common_util.WriteSuccessResponse(w, nil, "Request Sent Successfuly")
}
