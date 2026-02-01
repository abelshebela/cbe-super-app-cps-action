package bps_action

import "net/http"

type BPSActionAdapter interface {
	ApproveBPSAction(w http.ResponseWriter, r *http.Request)
	RejectBPSAction(w http.ResponseWriter, r *http.Request)
	GetBPSActionsByDepartment(w http.ResponseWriter, r *http.Request)
	GetUserCheckedActions(w http.ResponseWriter, r *http.Request)
	GetBPSActionByID(w http.ResponseWriter, r *http.Request)
	GetBPSActionByActionCode(w http.ResponseWriter, r *http.Request)
	GetActionCounts(w http.ResponseWriter, r *http.Request)
	GetUserAuditorActions(w http.ResponseWriter, r *http.Request)
	GetAuthorizerIndex(w http.ResponseWriter, r *http.Request)
	GetUserApproverActions(w http.ResponseWriter, r *http.Request)
	GetUserApproverApprovedActions(w http.ResponseWriter, r *http.Request)
	AuditorAction(w http.ResponseWriter, r *http.Request)
}
