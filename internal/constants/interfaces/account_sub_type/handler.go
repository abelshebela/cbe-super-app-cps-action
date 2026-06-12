package account_sub_type_interface

import "net/http"

type AccountSubTypeHandler interface {
	GetAllAccountSubTypes(w http.ResponseWriter, r *http.Request)
	GetOneAccountSubType(w http.ResponseWriter, r *http.Request)
	CreateOneAccountSubType(w http.ResponseWriter, r *http.Request)
	UpdateOneAccountSubType(w http.ResponseWriter, r *http.Request)
	DeleteOneAccountSubType(w http.ResponseWriter, r *http.Request)
	Enable(w http.ResponseWriter, r *http.Request)
	Disable(w http.ResponseWriter, r *http.Request)
}
