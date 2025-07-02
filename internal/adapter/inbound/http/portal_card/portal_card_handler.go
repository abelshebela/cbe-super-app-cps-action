package service

import (
	"encoding/json"
	"net/http"

	portalcardApp "github.com/CBE-Super-App/cbe-super-app-cps-action/internal/application/portal_card"
	"github.com/CBE-Super-App/cbe-super-app-cps-action/internal/port/inbound"

	"gitlab.com/bersufekadgetachew/cbe-super-app-shared/shared/utils"
)

type portalCardHandler struct {
	appService portalcardApp.PortalCardApplication
	logger     utils.Logger
}

func NewportalCardHandler(domain portalcardApp.PortalCardApplication, logger utils.Logger) inbound.PortalCardBound {
	return &portalCardHandler{
		appService: domain,
		logger:     logger,
	}
}

func (s *portalCardHandler) GetAllPortalCard(w http.ResponseWriter, r *http.Request) {

	returndata := make(map[string]interface{})
	userContext := r.Context().Value("user")

	if userContext == nil {
		returndata["message"] = "Unauthorized Access"
		w.WriteHeader(http.StatusUnauthorized)
		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(returndata)
		return
	}

	if r == nil {
		returndata["message"] = "Invalid Message Request"
		w.WriteHeader(http.StatusBadRequest)
		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(returndata)
		return
	}

	data, err := s.appService.GetAll(r.Context())
	if err != nil {
		returndata["message"] = "Unable to fetch the portal cards"
		returndata["status"] = 400
		w.WriteHeader(http.StatusNotFound)
		w.Header().Set("Content-Type", "applicaiton/json")
		json.NewEncoder(w).Encode(returndata)
		return
	}

	returndata["status"] = true
	returndata["message"] = "Successfuly Fetched"
	returndata["data"] = data

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusAccepted)
	json.NewEncoder(w).Encode(returndata)
	return
}
