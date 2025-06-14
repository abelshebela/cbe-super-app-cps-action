package branch_handler

import (
    "net/http"

    "github.com/go-chi/chi/v5"

    sharedhttp "gitlab.com/bersufekadgetachew/cbe-super-app-cps-action/internal/adapter/inbound/http"
    "gitlab.com/bersufekadgetachew/cbe-super-app-cps-action/internal/application/middleware"
    "gitlab.com/bersufekadgetachew/cbe-super-app-cps-action/internal/port/inbound"
)

func RegisterBranchRoutes(r chi.Router, handler inbound.BranchHandler, authMiddleware middleware.AuthMiddleware) {
    routes := []sharedhttp.Route{
        {
            Method:  http.MethodPost,
            Path:    "/branch/filter_single",
            Handler: handler.FilterSingleBranches,
            Middlewares: []func(next http.Handler) http.Handler{
                authMiddleware.AuthenticateToken,
                authMiddleware.AccessControl([]string{"MAKER", "IFB-MAKER"}),
            },
        },
        {
            Method:  http.MethodPost,
            Path:    "/branch/disable_single",
            Handler: handler.DisableSingleBranch,
            Middlewares: []func(next http.Handler) http.Handler{
                authMiddleware.AuthenticateToken,
                authMiddleware.AccessControl([]string{"MAKER", "IFB-MAKER"}),
            },
        },
        {
            Method:  http.MethodPost,
            Path:    "/branch/approve_disable_single",
            Handler: handler.ApproveSingleBranchDisable,
            Middlewares: []func(next http.Handler) http.Handler{
                authMiddleware.AuthenticateToken,
                authMiddleware.AccessControl([]string{"CHECKER", "IFB-CHECKER"}),
            },
        },
        {
            Method:  http.MethodPost,
            Path:    "/branch/filter_multiple",
            Handler: handler.FilterMultipleBranches,
            Middlewares: []func(next http.Handler) http.Handler{
                authMiddleware.AuthenticateToken,
                authMiddleware.AccessControl([]string{"MAKER", "IFB-MAKER"}),
            },
        },
        {
            Method:  http.MethodPost,
            Path:    "/branch/disable_multiple",
            Handler: handler.DisableMultipleBranches,
            Middlewares: []func(next http.Handler) http.Handler{
                authMiddleware.AuthenticateToken,
                authMiddleware.AccessControl([]string{"MAKER", "IFB-MAKER"}),
            },
        },
        {
            Method:  http.MethodPost,
            Path:    "/branch/approve_disable_multiple",
            Handler: handler.ApproveBulkBranchesDisable,
            Middlewares: []func(next http.Handler) http.Handler{
                authMiddleware.AuthenticateToken,
                authMiddleware.AccessControl([]string{"CHECKER", "IFB-CHECKER"}),
            },
        },
    }

    sharedhttp.RegisterRoutes(r, routes)
}