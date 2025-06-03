package inbound

import (
	"net/http"
	"github.com/go-chi/chi/v5"
)

type Handler interface {
	RegisterRoutes(r chi.Router) chi.Router
}

type UserPortHandler interface {
	Handler
	FetchLinkedAccounts(w http.ResponseWriter, r *http.Request)
}