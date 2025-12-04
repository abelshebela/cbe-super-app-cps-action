package faydaaccount

import "net/http"

type FaydaAccount interface {
	InitiateDisableFaydaAccount(w http.ResponseWriter, r *http.Request)
	InitiateEnableFaydaAccount(w http.ResponseWriter, r *http.Request)
}
