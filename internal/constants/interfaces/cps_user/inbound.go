package cpsuser

import "net/http"

type CPSUserHandler interface {
	CreateUserRequest(w http.ResponseWriter, r *http.Request)
	UpdateUserRequest(w http.ResponseWriter, r *http.Request)
	FetchUserByUserCode(w http.ResponseWriter, r *http.Request)
	GetAllCPSUsers(w http.ResponseWriter, r *http.Request)
	DeleteUserRequest(w http.ResponseWriter, r *http.Request)
	DisableUser(w http.ResponseWriter, r *http.Request)
	EnableUser(w http.ResponseWriter, r *http.Request)
}
