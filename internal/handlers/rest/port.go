package rest

import "net/http"

type SpendingHandler interface {
	GetOne(w http.ResponseWriter, r *http.Request)
	GetAll(w http.ResponseWriter, r *http.Request)
	InsertOne(w http.ResponseWriter, r *http.Request)
	UpdateOne(w http.ResponseWriter, r *http.Request)
	DeleteOne(w http.ResponseWriter, r *http.Request)
}
