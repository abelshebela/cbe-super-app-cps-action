package inbound

import "net/http"

type CustomerDetail interface {
	GetCustomerDetail(w http.ResponseWriter, r *http.Request)
	GetCustomerByID(w http.ResponseWriter, r *http.Request)
	GetFaydaCustomer(w http.ResponseWriter, r *http.Request)
	GetFaydaCustomerByID(w http.ResponseWriter, r *http.Request)
	GetBlockedCustomer(w http.ResponseWriter, r *http.Request)
}
