package wallet

import "net/http"

type WalletAdapter interface {
	GetAllWallet(w http.ResponseWriter, r *http.Request)
	GetWallet(w http.ResponseWriter, r *http.Request)
	CreateWallet(w http.ResponseWriter, r *http.Request)
	UpdateWallet(w http.ResponseWriter, r *http.Request)
	DeleteWallet(w http.ResponseWriter, r *http.Request)
	Authorize(w http.ResponseWriter, r *http.Request)
	Reject(w http.ResponseWriter, r *http.Request)
	Enable(w http.ResponseWriter,r *http.Request)
	Disable(w http.ResponseWriter,r *http.Request)
}
