package bank

import "net/http"

type BankHandler interface {
	GetAllBank(w http.ResponseWriter, r *http.Request)
	GetOneBank(w http.ResponseWriter, r *http.Request)
	CreateOneBank(w http.ResponseWriter, r *http.Request)
	UpdateOneBank(w http.ResponseWriter, r *http.Request)
	DeleteOneBank(w http.ResponseWriter, r *http.Request)
	Enable(w http.ResponseWriter, r *http.Request)
	Disable(w http.ResponseWriter, r *http.Request)
	UpdateLogo(w http.ResponseWriter, r *http.Request)
}
