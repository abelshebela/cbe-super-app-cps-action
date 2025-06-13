package branch_handler

import (
    "net/http"

    "github.com/go-chi/chi/v5"
    "gitlab.com/bersufekadgetachew/cbe-super-app-cps-ms/internal/application/middleware"
    "gitlab.com/bersufekadgetachew/cbe-super-app-cps-ms/internal/port/inbound"
    sharedhttp "gitlab.com/bersufekadgetachew/cbe-super-app-cps-ms/internal/adapter/inbound/http"
)

func RegisterBranchRoutes(r chi.Router, handler inbound.BranchHandler,authMiddleware middleware.AuthMiddleware) {
    routes := []sharedhttp.Route{
        {
            Method:  http.MethodPost,
            Path:    "/branch/filter-single",
            Handler: handler.FilterSingleBranches,
            Middlewares: []func(next http.Handler) http.Handler{
                authMiddleware.AuthenticateToken,
                authMiddleware.AccessControl([]string{"MAKER", "IFB-MAKER"}),
            },
        },
        {
            Method:  http.MethodPost,
            Path:    "/branch/disable-single",
            Handler: handler.DisableSingleBranch,
            Middlewares: []func(next http.Handler) http.Handler{
                authMiddleware.AuthenticateToken,
                authMiddleware.AccessControl([]string{"MAKER", "IFB-MAKER"}),
            },
        },
        {
            Method:  http.MethodPost,
            Path:    "/branch/approve-disable-single",
            Handler: handler.ApproveSingleBranchDisable,
            Middlewares: []func(next http.Handler) http.Handler{
                authMiddleware.AuthenticateToken,
                authMiddleware.AccessControl([]string{"CHECKER", "IFB-CHECKER"}),
            },
        },
        {
            Method:  http.MethodPost,
            Path:    "/branch/filter-multiple",
            Handler: handler.FilterMultipleBranches,
            Middlewares: []func(next http.Handler) http.Handler{
                authMiddleware.AuthenticateToken,
                authMiddleware.AccessControl([]string{"MAKER", "IFB-MAKER"}),
            },
        },
        {
            Method:  http.MethodPost,
            Path:    "/branch/disable-multiple",
            Handler: handler.DisableMultipleBranches,
            Middlewares: []func(next http.Handler) http.Handler{
                authMiddleware.AuthenticateToken,
                authMiddleware.AccessControl([]string{"MAKER", "IFB-MAKER"}),
            },
        },
        {
            Method:  http.MethodPost,
            Path:    "/branch/approve-disable-multiple",
            Handler: handler.ApproveBulkBranchesDisable,
            Middlewares: []func(next http.Handler) http.Handler{
                authMiddleware.AuthenticateToken,
                authMiddleware.AccessControl([]string{"CHECKER", "IFB-CHECKER"}),
            },
        },
    }

    sharedhttp.RegisterRoutes(r, routes)
}