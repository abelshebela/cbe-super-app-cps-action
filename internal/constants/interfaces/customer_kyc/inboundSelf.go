package customerkyc

import "net/http"

type SelfActivationKyc interface {
	GetAllSelfActivateKYCRequests(w http.ResponseWriter, r *http.Request)
	GetSelfActivateKYCRequest(w http.ResponseWriter, r *http.Request)
	ApproveSelfActivateKycRequest(w http.ResponseWriter, r *http.Request)
	RejectSelfActivateKycRequest(w http.ResponseWriter, r *http.Request)
	StartSelfActivateKycReview(w http.ResponseWriter, r *http.Request)
	PickSelfActivateKycReview(w http.ResponseWriter, r *http.Request)
	ExportUserSelfActivation(w http.ResponseWriter, r *http.Request)
	GetUsersActionLog(w http.ResponseWriter, r *http.Request)
}
