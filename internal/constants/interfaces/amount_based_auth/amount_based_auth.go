<<<<<<< HEAD
package amount_based_auth

import "net/http"

type AmountBasedAuthAdapter interface {
	GetAllAmountBasedAuth(w http.ResponseWriter, r *http.Request)
	UpdateAmountBasedAuth(w http.ResponseWriter, r *http.Request)
	RejectAmountBasedAuth(w http.ResponseWriter, r *http.Request)
	AddCurrency(w http.ResponseWriter, r *http.Request)
	ResetConfig(w http.ResponseWriter, r *http.Request)
}
=======
package amount_based_auth

import "net/http"

type AmountBasedAuthAdapter interface {
	GetAllAmountBasedAuth(w http.ResponseWriter, r *http.Request)
	UpdateAmountBasedAuth(w http.ResponseWriter, r *http.Request)
	RejectAmountBasedAuth(w http.ResponseWriter, r *http.Request)
}
>>>>>>> 3640c5b5dde77249222ec7dd9f90a5770ab0dbc4
