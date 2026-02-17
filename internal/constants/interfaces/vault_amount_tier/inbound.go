package vaultamounttier

import "net/http"

type VaultAmountTierHandler interface {
	CreateAmountTier(w http.ResponseWriter, r *http.Request)
	FindAllAmountTiers(w http.ResponseWriter, r *http.Request)
	GetAmountTier(w http.ResponseWriter, r *http.Request)
	UpdateAmountTier(w http.ResponseWriter, r *http.Request)
	DeleteAmountTier(w http.ResponseWriter, r *http.Request)
	DisableAmountTier(w http.ResponseWriter, r *http.Request)
	EnableAmountTier(w http.ResponseWriter, r *http.Request)
}
