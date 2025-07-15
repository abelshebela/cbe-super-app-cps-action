package service

import (
	"encoding/json"
	"errors"
	"net/http"
	"time"

	"github.com/CBE-Super-App/cbe-super-app-cps-action/internal/adapter/inbound/http/service/dto"
	serviceApp "github.com/CBE-Super-App/cbe-super-app-cps-action/internal/application/service"
	"github.com/CBE-Super-App/cbe-super-app-cps-action/internal/domain/action"
	cpsuser "github.com/CBE-Super-App/cbe-super-app-cps-action/internal/domain/cps_user"
	"github.com/CBE-Super-App/cbe-super-app-cps-action/internal/port/inbound"
	"github.com/go-chi/chi/v5"

	contexts "github.com/CBE-Super-App/cbe-super-app-cps-action/pkgs/context"
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

func (s *serviceHandler) GetAllServiceFee(w http.ResponseWriter, r *http.Request) {

	services, err := s.appService.GetAllService(r.Context())
	if err != nil {
		common_util.SendErrorResponse(w, err.Error(), 0, nil)
		return
	}
	common_util.BaseResponseMaker(map[string]interface{}{"services": services}, w, "Successfuly service fetched", 200)
}

func (s *serviceHandler) GetServiceFeeById(w http.ResponseWriter, r *http.Request) {
	id := chi.URLParam(r, "id")
	if id == "" {
		common_util.SendErrorResponse(w, "Service ID is required", 400, nil)
		return
	}
	service, err := s.appService.GetOneService(r.Context(), id)
	if err != nil {
		common_util.SendErrorResponse(w, err.Error(), 404, nil)
		return
	}
	common_util.BaseResponseMaker(service, w, "Service fetched successfully", 200)
}
func (s *serviceHandler) AuthorizeServiceFee(w http.ResponseWriter, r *http.Request) {
	action_code := chi.URLParam(r, "action_code")
	if action_code == "" {
		s.logger.Errorf("action code is required")
		common_util.SendErrorResponse(w, "Required input missing", 403, nil)
		return
	}

	err := s.appService.AuthorizeAction(r.Context(), action_code)
	if err != nil {
		s.logger.Errorf("Failed to Authorize Action Some thing goes wrong err %v", err.Error())
		common_util.SendErrorResponse(w, err.Error(), 500, nil)
		return
	}

	common_util.BaseResponseMaker(nil, w, "Successfuly Action Approved", 200)
}
func (s *serviceHandler) RejectServiceFee(w http.ResponseWriter, r *http.Request) {
	var req dto.RejectActionRequest
	action_code := chi.URLParam(r, "action_code")
	if action_code == "" {
		s.logger.Errorf("Action code is required")
		common_util.SendErrorResponse(w, "action_code", 519, nil)
		return
	}

	if req.Validate() != nil {
		s.logger.Errorf("Input Validatin Failed to For rejection")
		common_util.SendErrorResponse(w, req.Validate().Error(), 500, nil)
		return
	}

	err := s.appService.RejectAction(r.Context(), action_code, req.RejectReason)

	if err != nil {
		common_util.SendErrorResponse(w, err.Error(), 500, nil)
		return
	}

	common_util.BaseResponseMaker(nil, w, "Successfuly Action Rejected", 200)
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

	queryID := chi.URLParam(r, "id")
	if queryID == "" {
		common_util.SendErrorResponse(w, common_util.InvalidID, http.StatusBadRequest, nil)
		return
	}
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
		common_util.SendErrorResponse(w, "service not found", 0, nil)
		return
	}

	// Prepare CPS action data
	actionData := s.buildServiceFeeActionData(req, queryID, service.ServiceName)
	initAction := s.appService.InitCPSAction(user, actionData, "CREATE", "UPDATE SERVICE FEE", service)

	cpsAction, err := s.appService.CreateAction(r.Context(), initAction)
	if err != nil {
		common_util.SendErrorResponse(w, "failde", http.StatusConflict, nil)
		return
	}

	data, err := common_util.StructToMap(cpsAction)
	if err != nil {
		common_util.SendErrorResponse(w, err.Error(), 500, nil)
		return
	}

	common_util.BaseResponseMaker(data, w, "Update Request Sent Successfuly", 200)
}

func (s *serviceHandler) CreateServiceFee(w http.ResponseWriter, r *http.Request) {
	var req dto.CreateServiceRequest

	userContext := contexts.ExtractUserContext(r)

	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		common_util.SendErrorResponse(w, err.Error(), 0, nil)
		return
	}
	if req.Validate() != nil {
		common_util.SendErrorResponse(w, req.Validate(), 404, nil)
		return
	}

	cpsAction := action.CPSAction{
		CurrentAction:    req,
		UniqueId:         "SRV" + common_util.GenerateRandom(5),
		MakerID:          userContext.UserID,
		MakerName:        userContext.FullName,
		MakerPhoneNumber: userContext.PhoneNumber,
		Department:       userContext.Department,
		RequestAction:    action.RequestCreateServiceFee,
	}
	serviceFee, err := s.appService.CreateAction(r.Context(), cpsAction)
	if err != nil {
		common_util.SendErrorResponse(w, err.Error(), 500, nil)
		return
	}

	data, err := common_util.StructToMap(serviceFee)
	if err != nil {
		common_util.SendErrorResponse(w, err.Error(), 500, nil)
		return
	}

	common_util.BaseResponseMaker(data, w, "Successfuly Action Created", 200)
}
