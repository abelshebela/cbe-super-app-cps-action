package inbound

import "net/http"

type AmountBasedAuthHandler interface {
	UpdateAmountBasedAuth(w http.ResponseWriter, r *http.Request)
	ApproveAmountBasedAuth(w http.ResponseWriter, r *http.Request)
}
