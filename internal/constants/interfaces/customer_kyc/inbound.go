package customerkyc

import "net/http"

type CustomerKYC interface {
	CreateCustomerKYC(w http.ResponseWriter, r *http.Request)
	GetAllKYCRequests(w http.ResponseWriter, r *http.Request)
	GetKYCRequest(w http.ResponseWriter, r *http.Request)
	UpdateKYCStatus(w http.ResponseWriter, r *http.Request)
	DeleteKYCRequest(w http.ResponseWriter, r *http.Request)
}
