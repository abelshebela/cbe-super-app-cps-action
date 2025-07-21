package inbound

import "net/http"

type FaydaAccount interface {
	InitiateDisableFaydaAccount(w http.ResponseWriter, r *http.Request)
	GetAllFaydaAccounts(w http.ResponseWriter, r *http.Request)
}
