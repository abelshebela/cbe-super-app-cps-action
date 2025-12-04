package cps_action

import "net/http"

type CPSActionAdapter interface {
	ApproveCPSAction(w http.ResponseWriter, r *http.Request)
	RejectCPSAction(w http.ResponseWriter, r *http.Request)
	GetCPSActionsByDepartment(w http.ResponseWriter, r *http.Request)
	GetCPSActionByID(w http.ResponseWriter, r *http.Request)
	GetCPSActionByActionCode(w http.ResponseWriter, r *http.Request)
	GetActionCounts(w http.ResponseWriter, r *http.Request)
}
