package http

import (
	"github.com/go-chi/chi/v5"
)

func RegisterRoutes(r chi.Router, handler *HTTPHandler) {
	userRouter := chi.NewRouter()
	userRouter.Get("/{id}/linked-accounts", handler.FetchLinkedAccounts)
	userRouter.Post("/{id}/email/generate-otp", handler.GenerateEmailOTP)
	userRouter.Post("/{id}/email/verify-otp", handler.VerifyEmailOTP)
	r.Mount("/user", userRouter)
}