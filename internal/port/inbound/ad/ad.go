package ad

import "net/http"

type ADAdapter interface {
	CreateOneAdvert(w http.ResponseWriter, r *http.Request)
	GetAllAdvert(w http.ResponseWriter, r *http.Request)
	GetOneAdvert(w http.ResponseWriter, r *http.Request)
	UpdateOneAdvert(w http.ResponseWriter, r *http.Request)
	DeleteOneAdvert(w http.ResponseWriter, r *http.Request)
	EnableAdvert(w http.ResponseWriter, r *http.Request)
	DisableAdvert(w http.ResponseWriter, r *http.Request)
}
