package service

import (
	"encoding/json"
	"net/http"
	"time"

	"github.com/CBE-Super-App/cbe-super-app-cps-action/internal/adapter/inbound/http/service/dto"
	serviceApp "github.com/CBE-Super-App/cbe-super-app-cps-action/internal/application/service"
	cpsuser "github.com/CBE-Super-App/cbe-super-app-cps-action/internal/domain/cps_user"
	"github.com/CBE-Super-App/cbe-super-app-cps-action/internal/port/inbound"

	"gitlab.com/bersufekadgetachew/cbe-super-app-shared/shared/utils"
	"go.mongodb.org/mongo-driver/v2/bson"
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

func (s *serviceHandler) UpdateServiceFee(w http.ResponseWriter, r *http.Request) {
	var req dto.UpdateServiceFeeRequest
	returndata := make(map[string]interface{})

	queryId := r.Header.Get("query_id")
	var user cpsuser.CPSUser
	userContext := r.Header.Get("_user")

	err := json.Unmarshal([]byte(userContext), &user)

	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		returndata["success"] = false
		returndata["message"] = "invalid json data"
		returndata["error"] = err.Error()

		w.WriteHeader(http.StatusForbidden)
		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(returndata)
		return
	}

	if err := s.appService.ValidateTiers(req.Tiers, req.AboveAmount); err != nil {
		returndata["success"] = false
		returndata["message"] = err.Error()
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(w).Encode(returndata)
		return
	}

	objectID, err := bson.ObjectIDFromHex(queryId)
	if err != nil {
		returndata["success"] = false
		returndata["message"] = "Invalid message query id"
		w.Header().Set("Content-Type", "applicatoin/json")
		w.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(w).Encode(returndata)
		return
	}

	service, err := s.appService.GetOneService(r.Context(), objectID)
	if err != nil {
		returndata["success"] = false
		returndata["message"] = "Service not found"
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusNotFound)
		return
	}

	serviceName := service.ServiceName

	// Prepare CPS action data
	action := map[string]interface{}{
		"aboveServiceFee":   req.AboveServiceFee / 100,
		"minAmountVIRTUAL":  req.MinAmountVIRTUAL,
		"dailyCapLevelOne":  req.DailyCapLevelOne,
		"singleCapLevelOne": req.SingleCapLevelOne,
		"service_id":        queryId,
		"service_name":      serviceName,
		"tiers":             req.Tiers,
		"paymentType":       req.PaymentType,
		"serviceType":       req.ServiceType,
		"lastModified":      time.Now(),
	}
	initAction := s.appService.InitCPSAction(user, action, "CREATE", "UPDATE SERVICE FEE", service)

	err = s.appService.CreateAction(r.Context(), initAction)
	if err != nil {
		returndata["success"] = false
		returndata["message"] = "Unable to Create  Action"
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusConflict)
		json.NewEncoder(w).Encode(returndata)
		return
	}

	returndata["message"] = "Request Sent Successfuly"
	returndata["status"] = true
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusAccepted)
	json.NewEncoder(w).Encode(returndata)
	return
}
