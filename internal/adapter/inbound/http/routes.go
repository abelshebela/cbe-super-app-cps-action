package http

import (
	"github.com/go-chi/chi/v5"
)

func RegisterRoutes(r chi.Router, handler *HTTPHandler) {
	r.Get("/{id}/linked-accounts", handler.FetchLinkedAccounts)
		r.Post("/{id}/email/generate-otp", handler.GenerateEmailOTP)
	r.Post("/{id}/email/verify-otp", handler.VerifyEmailOTP)
}