package cps_action_role

import "net/http"

type CPSActionRoleHandler interface {
	GetAll(w http.ResponseWriter, r *http.Request)
	GetAllActionList(w http.ResponseWriter, r *http.Request)
	GetByActionCode(w http.ResponseWriter, r *http.Request)
	Create(w http.ResponseWriter, r *http.Request)
	Update(w http.ResponseWriter, r *http.Request)
	Enable(w http.ResponseWriter, r *http.Request)
	Disable(w http.ResponseWriter, r *http.Request)
	GetVersions(w http.ResponseWriter, r *http.Request)
	GetConfiguredRoles(w http.ResponseWriter, r *http.Request)
	UpdateVersionRoleCode(w http.ResponseWriter, r *http.Request)
}
