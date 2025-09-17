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

type Service interface {
	GetAllService(w http.ResponseWriter, r *http.Request)
	GetAllMinimumTransferCap(w http.ResponseWriter, r *http.Request)
	GetAllMaximumTransferCap(w http.ResponseWriter, r *http.Request)
	GetAllServiceFee(w http.ResponseWriter, r *http.Request)
	GetAllTotalTransferCap(w http.ResponseWriter, r *http.Request)
	GetServiceFeeDetail(w http.ResponseWriter, r *http.Request)
	UpdateServiceFee(w http.ResponseWriter, r *http.Request)
	UpdateSingleMaxTransfer(w http.ResponseWriter, r *http.Request)
	UpdateTotalMaxTransferCap(w http.ResponseWriter, r *http.Request)
	UpdateMinimumTransferCap(w http.ResponseWriter, r *http.Request)
	DeleteServiceFeeTire(w http.ResponseWriter, r *http.Request)
}
