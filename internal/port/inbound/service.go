package inbound

import "net/http"

type ServiceBound interface {
	GetAllServiceFee(w http.ResponseWriter, r *http.Request)
	GetServiceFeeById(w http.ResponseWriter, r *http.Request)
	CreateServiceFee(w http.ResponseWriter, r *http.Request)
	UpdateServiceFee(w http.ResponseWriter, r *http.Request)
	AuthorizeServiceFee(w http.ResponseWriter, r *http.Request)
	RejectServiceFee(w http.ResponseWriter, r *http.Request)
}
