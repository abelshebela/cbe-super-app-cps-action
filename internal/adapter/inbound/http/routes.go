package http

import (
    "cbe-super-app-member-users/internal/adapter/inbound/http/account"
    "cbe-super-app-member-users/internal/adapter/inbound/http/users"
    "github.com/go-chi/chi/v5"
)

func RegisterRoutes(r chi.Router, userHandler *users.HTTPHandler, accountHandler *account.HTTPHandler) {
    userRouter := chi.NewRouter()
    userRouter.Get("/{id}/linked-accounts", userHandler.FetchLinkedAccounts)
    userRouter.Post("/{id}/email/generate-otp", userHandler.GenerateEmailOTP)
    userRouter.Post("/{id}/email/verify-otp", userHandler.VerifyEmailOTP)
    userRouter.Post("/{id}/add-account", accountHandler.CreateAccount)
    r.Mount("/user", userRouter)
}