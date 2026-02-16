package vault

import "net/http"

type VaultCategoryHandler interface {
	CreateVaultCategory(w http.ResponseWriter, r *http.Request)
	FindAllVaultCategories(w http.ResponseWriter, r *http.Request)
	GetVaultCategory(w http.ResponseWriter, r *http.Request)
	UpdateVaultCategory(w http.ResponseWriter, r *http.Request)
	DeleteVaultCategory(w http.ResponseWriter, r *http.Request)
	DisableVaultCategory(w http.ResponseWriter, r *http.Request)
	EnableVaultCategory(w http.ResponseWriter, r *http.Request)

	// Transaction
	GetVaultTransactions(w http.ResponseWriter, r *http.Request)
	GetVaultTransaction(w http.ResponseWriter, r *http.Request)

	// Withdrawal Request
	CreateWithdrawalRequest(w http.ResponseWriter, r *http.Request)
	UpdateWithDrawalRequest(w http.ResponseWriter, r *http.Request)
	GetAllWithdrawalRequests(w http.ResponseWriter, r *http.Request)
	GetWithdrawalRequestById(w http.ResponseWriter, r *http.Request)
}
