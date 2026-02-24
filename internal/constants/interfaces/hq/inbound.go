package hq

import "net/http"

type HQAdapter interface {
	GetHQ(w http.ResponseWriter, r *http.Request)
	GetAllHQ(w http.ResponseWriter, r *http.Request)
	GetBlockTime(w http.ResponseWriter, r *http.Request)
	GetArchiveTime(w http.ResponseWriter, r *http.Request)
	GetPasswordExpiry(w http.ResponseWriter, r *http.Request)

	UpdateBlockTimeRequest(w http.ResponseWriter, r *http.Request)
	UpdateArchiveTimeRequest(w http.ResponseWriter, r *http.Request)
	UpdatePasswordExpiryRequest(w http.ResponseWriter, r *http.Request)
}
