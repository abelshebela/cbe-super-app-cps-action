package inbound

import "net/http"

type BPSUserHandler interface {
	GetPendingUserActions(w http.ResponseWriter, r *http.Request)
	FetchUserByUserCode(w http.ResponseWriter, r *http.Request)
	GetAllBPSUsers(w http.ResponseWriter, r *http.Request)
	DisableUser(w http.ResponseWriter, r *http.Request)
	EnableUser(w http.ResponseWriter, r *http.Request)
}
