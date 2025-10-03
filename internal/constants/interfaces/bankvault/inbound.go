package bankvault

import "net/http"

type BankVaultHandler interface {
	FindAllBankVaults(w http.ResponseWriter, r *http.Request)
	GetBankVault(w http.ResponseWriter, r *http.Request)
	CreateBankVault(w http.ResponseWriter, r *http.Request)
	UpdateBankVault(w http.ResponseWriter, r *http.Request)
	DeleteBankVault(w http.ResponseWriter, r *http.Request)
	DisableBankVault(w http.ResponseWriter, r *http.Request)
	EnableBankVault(w http.ResponseWriter, r *http.Request)
}
