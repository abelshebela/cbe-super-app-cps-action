package service_details

import (
	"net/http"
)

type ServiceAdapter interface {
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
