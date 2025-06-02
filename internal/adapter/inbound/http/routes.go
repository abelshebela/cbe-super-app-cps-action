package http

import (
    "github.com/go-chi/chi/v5"
)

func RegisterRoutes(r chi.Router, handler *HTTPHandler) {
    r.Get("/{id}/linked-accounts", handler.FetchLinkedAccounts)
}