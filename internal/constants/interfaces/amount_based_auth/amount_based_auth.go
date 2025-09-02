package amount_based_auth

import "net/http"


type AmountBasedAuthAdapter interface {
	GetAllAmountBasedAuth(w http.ResponseWriter, r *http.Request)
	UpdateAmountBasedAuth(w http.ResponseWriter, r *http.Request)
	RejectAmountBasedAuth(w http.ResponseWriter, r *http.Request)
}