package bank

import "net/http"

type BankAdapter interface {
	GetAllBank(w http.ResponseWriter, r *http.Request)
	GetOneBank(w http.ResponseWriter, r *http.Request)
	CreateOneBank(w http.ResponseWriter, r *http.Request)
	UpdateOneBank(w http.ResponseWriter, r *http.Request)
	DeleteOneBank(w http.ResponseWriter, r *http.Request)
	Authorize(w http.ResponseWriter, r *http.Request)
	Reject(w http.ResponseWriter, r *http.Request)
	Enable(w http.ResponseWriter, r *http.Request)
	Disable(w http.ResponseWriter, r *http.Request)
}
