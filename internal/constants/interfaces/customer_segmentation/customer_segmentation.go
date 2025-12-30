package customersegmentation

import "net/http"

type CustomerSegmentation interface {
	CreateCustomerSegmentation(w http.ResponseWriter, r *http.Request)
	UpdateCustomerSegmentation(w http.ResponseWriter, r *http.Request)
	GetAllCustomerSegmentations(w http.ResponseWriter, r *http.Request)
	GetCustomerSegmentation(w http.ResponseWriter, r *http.Request)
	DeleteCustomerSegmentation(w http.ResponseWriter, r *http.Request)
	Enable(w http.ResponseWriter, r *http.Request)
	Disable(w http.ResponseWriter, r *http.Request)
}
