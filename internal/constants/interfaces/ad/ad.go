package ad

import "net/http"

type ADAdapter interface {
	CreateAdvert(w http.ResponseWriter, r *http.Request)
	FetchAdverts(w http.ResponseWriter, r *http.Request)
	FetchAdvertByID(w http.ResponseWriter, r *http.Request)
	UpdateAdvert(w http.ResponseWriter, r *http.Request)
	DeleteAdvert(w http.ResponseWriter, r *http.Request)
	EnableAdvert(w http.ResponseWriter, r *http.Request)
	DisableAdvert(w http.ResponseWriter, r *http.Request)
}
