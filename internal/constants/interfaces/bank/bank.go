package bank

import "net/http"

type BankAdapter interface {
	CreateOneBank(w http.ResponseWriter, r *http.Request)
	UpdateOneBank(w http.ResponseWriter, r *http.Request)
	DeleteOneBank(w http.ResponseWriter, r *http.Request)
	GetOneBank(w http.ResponseWriter, r *http.Request)
	GetAllBank(w http.ResponseWriter, r *http.Request)
	Enable(w http.ResponseWriter, r *http.Request)
	Disable(w http.ResponseWriter, r *http.Request)
	UpdateLogo(w http.ResponseWriter, r *http.Request)
} 