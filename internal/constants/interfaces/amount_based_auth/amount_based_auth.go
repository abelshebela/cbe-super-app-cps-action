package amount_based_auth

import "net/http"


type AmountBasedAuthAdapter interface {
	GetAllAmountBasedAuth(w http.ResponseWriter, r *http.Request)
	UpdateOpenTier(w http.ResponseWriter, r *http.Request)
	UpdatePinTier(w http.ResponseWriter, r *http.Request)
	UpdateOtpPinTier(w http.ResponseWriter, r *http.Request)
	RejectAmountBasedAuth(w http.ResponseWriter, r *http.Request)
}