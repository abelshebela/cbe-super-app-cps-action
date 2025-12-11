package cps_action_role

import "net/http"

type CPSActionRoleHandler interface {
	GetAll(w http.ResponseWriter, r *http.Request)
	GetByActionCode(w http.ResponseWriter, r *http.Request)
	Create(w http.ResponseWriter, r *http.Request)
	Update(w http.ResponseWriter, r *http.Request)
	Enable(w http.ResponseWriter, r *http.Request)
	Disable(w http.ResponseWriter, r *http.Request)
}
