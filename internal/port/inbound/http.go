package inbound

import (
    // "chi"
    chi "github.com/go-chi/chi/v5"
    "cbe-super-app-member-users/internal/port/inbound/users"
    "cbe-super-app-member-users/internal/port/inbound/account"
)

type Handler interface {
    RegisterRoutes(r chi.Router) chi.Router
}

type UserPortHandler interface {
    Handler
    users.UserPortHandler
}

type AccountPortHandler interface {
    Handler
    account.AccountPortHandler
}