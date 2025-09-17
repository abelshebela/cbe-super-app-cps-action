package inbound

import "net/http"

type Inbound interface {
	FetchAccountValidation(w http.ResponseWriter, r *http.Request)
	FetchAllAccountValidation(w http.ResponseWriter, r *http.Request)
	UpdateAccountValidationMaker(w http.ResponseWriter, r *http.Request)
	UpdateAccountValidationChecker(w http.ResponseWriter, r *http.Request)
}
