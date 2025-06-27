package accountblock_handler

import (
    "net/http"

    "github.com/go-chi/chi/v5"
    sharedhttp "gitlab.com/bersufekadgetachew/cbe-super-app-cps-action/internal/adapter/inbound/http"
    "gitlab.com/bersufekadgetachew/cbe-super-app-cps-action/internal/application/middleware"
    accountblock "gitlab.com/bersufekadgetachew/cbe-super-app-cps-action/internal/port/inbound/account_block"
)

func RegisterAccountBlockRoutes(
    r chi.Router,
    handler accountblock.AccountBlockHandler,
    authMiddleware middleware.AuthMiddleware,
) {
    routes := []sharedhttp.Route{
        {
            Method:  http.MethodPost,
            Path:    "/accountblock/branch/filter-single",
            Handler: handler.FilterSingleBranches,
            Middlewares: []func(http.Handler) http.Handler{
                authMiddleware.AuthenticateToken,
                authMiddleware.AccessControl([]string{"maker", "ifb-maker"}),
            },
        },
        {
            Method:  http.MethodPost,
            Path:    "/accountblock/branch/disable-single",
            Handler: handler.DisableSingleBranch,
            Middlewares: []func(http.Handler) http.Handler{
                authMiddleware.AuthenticateToken,
                authMiddleware.AccessControl([]string{"maker", "ifb-maker"}),
            },
        },
        {
            Method:  http.MethodPost,
            Path:    "/accountblock/branch/approve-disable-single",
            Handler: handler.ApproveSingleBranchDisable,
            Middlewares: []func(http.Handler) http.Handler{
                authMiddleware.AuthenticateToken,
                authMiddleware.AccessControl([]string{"checker", "ifb-checker"}),
            },
        },

        {
            Method:  http.MethodPost,
            Path:    "/accountblock/branch/filter-multiple",
            Handler: handler.FilterMultipleBranches,
            Middlewares: []func(http.Handler) http.Handler{
                authMiddleware.AuthenticateToken,
                authMiddleware.AccessControl([]string{"maker", "ifb-maker"}),
            },
        },
        {
            Method:  http.MethodPost,
            Path:    "/accountblock/branch/disable-multiple",
            Handler: handler.DisableMultipleBranches,
            Middlewares: []func(http.Handler) http.Handler{
                authMiddleware.AuthenticateToken,
                authMiddleware.AccessControl([]string{"maker", "ifb-maker"}),
            },
        },
        {
            Method:  http.MethodPost,
            Path:    "/accountblock/branch/approve-disable-multiple",
            Handler: handler.ApproveBulkBranchesDisable,
            Middlewares: []func(http.Handler) http.Handler{
                authMiddleware.AuthenticateToken,
                authMiddleware.AccessControl([]string{"checker", "ifb-checker"}),
            },
        },

        {
            Method:  http.MethodPost,
            Path:    "/accountblock/region/block",
            Handler: handler.BlockRegion,
            Middlewares: []func(http.Handler) http.Handler{
                authMiddleware.AuthenticateToken,
                authMiddleware.AccessControl([]string{"maker", "ifb-maker"}),
            },
        },
        {
            Method:  http.MethodPut,
            Path:    "/accountblock/region/update",
            Handler: handler.UpdateRegion,
            Middlewares: []func(http.Handler) http.Handler{
                authMiddleware.AuthenticateToken,
                authMiddleware.AccessControl([]string{"maker", "ifb-maker"}),
            },
        },
        {
            Method:  http.MethodPost,
            Path:    "/accountblock/region/approve-block",
            Handler: handler.ApproveRegionBlock,
            Middlewares: []func(http.Handler) http.Handler{
                authMiddleware.AuthenticateToken,
                authMiddleware.AccessControl([]string{"checker", "ifb-checker"}),
            },
        },
        {
            Method:  http.MethodGet,
            Path:    "/accountblock/region/{id}",
            Handler: handler.GetRegionByID,
            Middlewares: []func(http.Handler) http.Handler{
                authMiddleware.AuthenticateToken,
                authMiddleware.AccessControl([]string{"maker", "ifb-maker", "checker", "ifb-checker"}),
            },
        },

        {
            Method:  http.MethodPost,
            Path:    "/accountblock/district/block",
            Handler: handler.BlockDistrict,
            Middlewares: []func(http.Handler) http.Handler{
                authMiddleware.AuthenticateToken,
                authMiddleware.AccessControl([]string{"maker", "ifb-maker"}),
            },
        },
        {
            Method:  http.MethodPost,
            Path:    "/accountblock/district/approve-block",
            Handler: handler.ApproveBlockDistrict,
            Middlewares: []func(http.Handler) http.Handler{
                authMiddleware.AuthenticateToken,
                authMiddleware.AccessControl([]string{"checker", "ifb-checker"}),
            },
        },
        {
            Method:  http.MethodGet,
            Path:    "/accountblock/district/{id}",
            Handler: handler.GetDistrictByID,
            Middlewares: []func(http.Handler) http.Handler{
                authMiddleware.AuthenticateToken,
                authMiddleware.AccessControl([]string{"maker", "ifb-maker", "checker", "ifb-checker"}),
            },
        },

        {
            Method:  http.MethodPost,
            Path:    "/accountblock/city/block",
            Handler: handler.BlockCity,
            Middlewares: []func(http.Handler) http.Handler{
                authMiddleware.AuthenticateToken,
                authMiddleware.AccessControl([]string{"maker", "ifb-maker"}),
            },
        },
        {
            Method:  http.MethodPost,
            Path:    "/accountblock/city/approve-block",
            Handler: handler.ApproveBlockCity,
            Middlewares: []func(http.Handler) http.Handler{
                authMiddleware.AuthenticateToken,
                authMiddleware.AccessControl([]string{"checker", "ifb-checker"}),
            },
        },
        {
            Method:  http.MethodGet,
            Path:    "/accountblock/city/{id}",
            Handler: handler.GetCityByID,
            Middlewares: []func(http.Handler) http.Handler{
                authMiddleware.AuthenticateToken,
                authMiddleware.AccessControl([]string{"maker", "ifb-maker", "checker", "ifb-checker"}),
            },
        },

        {
            Method:  http.MethodPost,
            Path:    "/accountblock/user/block",
            Handler: handler.BlockUser,
            Middlewares: []func(http.Handler) http.Handler{
                authMiddleware.AuthenticateToken,
                authMiddleware.AccessControl([]string{"maker", "ifb-maker"}),
            },
        },
        {
            Method:  http.MethodPost,
            Path:    "/accountblock/user/approve-block",
            Handler: handler.ApproveBlockUser,
            Middlewares: []func(http.Handler) http.Handler{
                authMiddleware.AuthenticateToken,
                authMiddleware.AccessControl([]string{"checker", "ifb-checker"}),
            },
        },
        {
            Method:  http.MethodGet,
            Path:    "/accountblock/user/{id}",
            Handler: handler.GetUserByID,
            Middlewares: []func(http.Handler) http.Handler{
                authMiddleware.AuthenticateToken,
                authMiddleware.AccessControl([]string{"maker", "ifb-maker", "checker", "ifb-checker"}),
            },
        },
    }

    sharedhttp.RegisterRoutes(r, routes)
}