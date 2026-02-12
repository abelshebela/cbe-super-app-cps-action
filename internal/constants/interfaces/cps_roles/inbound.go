package cpsroles

import "net/http"

type CPSRolesAdapter interface {
	CreateCPSRole(w http.ResponseWriter, r *http.Request)
	UpdateCPSRole(w http.ResponseWriter, r *http.Request)
	GetAllCPSRoles(w http.ResponseWriter, r *http.Request)
	GetCPSRole(w http.ResponseWriter, r *http.Request)
	EnableCPSRole(w http.ResponseWriter, r *http.Request)
	DisableCPSRole(w http.ResponseWriter, r *http.Request)
}
