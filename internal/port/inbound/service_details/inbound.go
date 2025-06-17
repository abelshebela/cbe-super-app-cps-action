package service_details

import (
	"net/http"
)

type ServiceDetailsInbound interface {
	GetAllServiceDetails(w http.ResponseWriter, r *http.Request)
	GetServiceDetailsByID(w http.ResponseWriter, r *http.Request)
	UpdateServiceDetailsMaker(w http.ResponseWriter, r *http.Request)
	UpdateServiceDetailsChecker(w http.ResponseWriter, r *http.Request)

	ServiceDetailsDailyCapMaker(w http.ResponseWriter, r *http.Request)
	ServiceDetailsSingleCapMaker(w http.ResponseWriter, r *http.Request)
	TotalTransferCapMaker(w http.ResponseWriter, r *http.Request)
}
