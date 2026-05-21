package customergroup

import "net/http"

type CustomerGroup interface {
	CreateCustomerGroup(w http.ResponseWriter, r *http.Request)
	UpdateCustomerGroup(w http.ResponseWriter, r *http.Request)
	GetAllCustomerGroups(w http.ResponseWriter, r *http.Request)
	GetCustomerGroup(w http.ResponseWriter, r *http.Request)
	DeleteCustomerGroup(w http.ResponseWriter, r *http.Request)
	Enable(w http.ResponseWriter, r *http.Request)
	Disable(w http.ResponseWriter, r *http.Request)
}
