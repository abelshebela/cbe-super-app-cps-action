package inbound

import "net/http"

type FaydaAccount interface {
	InitiateDisableFaydaAccount(w http.ResponseWriter, r *http.Request)
	AuthorizeFaydaAccountDisable(w http.ResponseWriter, r *http.Request)
	RejectFaydaAccountDisable(w http.ResponseWriter, r *http.Request)
	GetAllFaydaAccounts(w http.ResponseWriter, r *http.Request)
}
