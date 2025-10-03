package vaultgroupcategory

import "net/http"

type VaultGroupCategoryHandler interface {
	CreateVaultGroupCategory(w http.ResponseWriter, r *http.Request)
	FindAllVaultGroupCategories(w http.ResponseWriter, r *http.Request)
	GetVaultGroupCategory(w http.ResponseWriter, r *http.Request)
	UpdateVaultGroupCategory(w http.ResponseWriter, r *http.Request)
	DeleteVaultGroupCategory(w http.ResponseWriter, r *http.Request)
	DisableVaultGroupCategory(w http.ResponseWriter, r *http.Request)
	EnableVaultGroupCategory(w http.ResponseWriter, r *http.Request)
}
