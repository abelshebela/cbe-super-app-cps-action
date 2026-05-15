package wallet

import "net/http"

type WalletAdapter interface {
	GetAllWallet(w http.ResponseWriter, r *http.Request)
	GetWallet(w http.ResponseWriter, r *http.Request)
	CreateWallet(w http.ResponseWriter, r *http.Request)
	UpdateWallet(w http.ResponseWriter, r *http.Request)
	DeleteWallet(w http.ResponseWriter, r *http.Request)
	Enable(w http.ResponseWriter, r *http.Request)
	Disable(w http.ResponseWriter, r *http.Request)
	EnableWalletService(w http.ResponseWriter, r *http.Request)
	DisableWalletService(w http.ResponseWriter, r *http.Request)
}
