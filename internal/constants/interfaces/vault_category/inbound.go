package vaultcategory

import "net/http"

type VaultCategoryHandler interface {
	CreateVaultCategory(w http.ResponseWriter, r *http.Request)
	FindAllVaultCategories(w http.ResponseWriter, r *http.Request)
	GetVaultCategory(w http.ResponseWriter, r *http.Request)
	UpdateVaultCategory(w http.ResponseWriter, r *http.Request)
	DeleteVaultCategory(w http.ResponseWriter, r *http.Request)
	DisableVaultCategory(w http.ResponseWriter, r *http.Request)
	EnableVaultCategory(w http.ResponseWriter, r *http.Request)
}
