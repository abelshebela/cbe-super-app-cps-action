package cps_action

import "net/http"

type CPSActionAdapter interface {
	ApproveCPSAction(w http.ResponseWriter, r *http.Request)
	RejectCPSAction(w http.ResponseWriter, r *http.Request)
	CancelCPSAction(w http.ResponseWriter, r *http.Request)
	ReverseCPSAction(w http.ResponseWriter, r *http.Request)
	GetCPSActionsByDepartment(w http.ResponseWriter, r *http.Request)
	GetUserCreatedActions(w http.ResponseWriter, r *http.Request)
	GetUserCheckedActions(w http.ResponseWriter, r *http.Request)
	GetCPSActionByID(w http.ResponseWriter, r *http.Request)
	GetCPSActionByActionCode(w http.ResponseWriter, r *http.Request)
	GetActionCounts(w http.ResponseWriter, r *http.Request)
	ApproverCheckerAllocations(w http.ResponseWriter, r *http.Request)
	ApproverAuditorAllocations(w http.ResponseWriter, r *http.Request)
	GetUserApproverPendingActions(w http.ResponseWriter, r *http.Request)
	GetUserApproverApprovedActions(w http.ResponseWriter, r *http.Request)
}
