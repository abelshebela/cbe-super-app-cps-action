package amountbasedauth

import (
	"net/http"
)

type AmountBasedAuthHandler interface {
	UpdateAmountBasedAuth(w http.ResponseWriter, r *http.Request)
	// ApproveAmountBasedAuth(w http.ResponseWriter, r *http.Request)
	RejectAmountBasedAuth(w http.ResponseWriter, r *http.Request)
	GetAllAmountBasedAuth(w http.ResponseWriter, r *http.Request)
}
