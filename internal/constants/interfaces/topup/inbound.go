package topup

import "net/http"

type TopupAdapter interface {
	GetAllTopup(w http.ResponseWriter, r *http.Request)
	GetTopup(w http.ResponseWriter, r *http.Request)
	CreateTopup(w http.ResponseWriter, r *http.Request)
	UpdateTopup(w http.ResponseWriter, r *http.Request)
	DeleteTopup(w http.ResponseWriter, r *http.Request)
	Enable(w http.ResponseWriter, r *http.Request)
	Disable(w http.ResponseWriter, r *http.Request)
}
