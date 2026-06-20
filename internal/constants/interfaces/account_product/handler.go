package account_product_interface

import "net/http"

type AccountProductHandler interface {
	GetAll(w http.ResponseWriter, r *http.Request)
	GetByID(w http.ResponseWriter, r *http.Request)
	Create(w http.ResponseWriter, r *http.Request)
	Update(w http.ResponseWriter, r *http.Request)
	Delete(w http.ResponseWriter, r *http.Request)
	Enable(w http.ResponseWriter, r *http.Request)
	Disable(w http.ResponseWriter, r *http.Request)
}
