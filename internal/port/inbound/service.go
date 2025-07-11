package inbound

import "net/http"

type ServiceBound interface {
	CreateServiceFee(w http.ResponseWriter, r *http.Request)
	UpdateServiceFee(w http.ResponseWriter, r *http.Request)
	AuthorizeServiceFee(w http.ResponseWriter, r *http.Request)
	RejectServiceFee(w http.ResponseWriter, r *http.Request)
}
