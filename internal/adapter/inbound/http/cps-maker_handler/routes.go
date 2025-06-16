package cpsmakerhandler

import (
	"net/http"

	"github.com/go-chi/chi/v5"

	sharedhttp "gitlab.com/bersufekadgetachew/cbe-super-app-cps-action/internal/adapter/inbound/http"
	"gitlab.com/bersufekadgetachew/cbe-super-app-cps-action/internal/application/middleware"
	"gitlab.com/bersufekadgetachew/cbe-super-app-cps-action/internal/port/inbound"
)

func RegisterCPSUserMakerRoutes(r chi.Router, handler inbound.CPSUserMakerHandler, authMiddleware middleware.AuthMiddleware) {
    routes := []sharedhttp.Route{
        {
            Method:  http.MethodPost,
            Path:    "/cps_user_maker/create",
            Handler: handler.CreateUserRequest,
            Middlewares: []func(next http.Handler) http.Handler{
                authMiddleware.AuthenticateToken,
                authMiddleware.AccessControl([]string{"MAKER", "IFB-MAKER"}),
            },
        },
        {
            Method:  http.MethodPut,
            Path:    "/cps_user_maker/update",
            Handler: handler.UpdateUserRequest,
            Middlewares: []func(next http.Handler) http.Handler{
                authMiddleware.AuthenticateToken,
                authMiddleware.AccessControl([]string{"MAKER", "IFB-MAKER"}),
            },
        },
        {
            Method:  http.MethodPost,
            Path:    "/cps_user_maker/approve",
            Handler: handler.ApproveUserAction,
            Middlewares: []func(next http.Handler) http.Handler{
                authMiddleware.AuthenticateToken,
                authMiddleware.AccessControl([]string{"CHECKER", "IFB-CHECKER"}),
            },
        },
        {
            Method:  http.MethodGet,
            Path:    "/cps_user_maker/pending_actions",
            Handler: handler.GetPendingUserActions,
            Middlewares: []func(next http.Handler) http.Handler{
                authMiddleware.AuthenticateToken,
            },
        },
    }

    sharedhttp.RegisterRoutes(r, routes)
}
