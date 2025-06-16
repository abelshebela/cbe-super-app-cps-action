package inbound

import "net/http"

type CustomerDetail interface {
	GetCustomerDetail(w http.ResponseWriter, r *http.Request)
	GetCustomerByID(w http.ResponseWriter, r *http.Request)
}


