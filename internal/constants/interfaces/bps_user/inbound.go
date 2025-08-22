package bps_user

import "net/http"

type BPSUserHandler interface {
	FetchUserByUserCode(w http.ResponseWriter, r *http.Request)
	GetAllBPSUsers(w http.ResponseWriter, r *http.Request)
	DisableUser(w http.ResponseWriter, r *http.Request)
	EnableUser(w http.ResponseWriter, r *http.Request)
}
