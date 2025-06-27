package accountblock_handler

import (
	"net/http"

	"github.com/go-chi/chi/v5"

	sharedhttp "gitlab.com/bersufekadgetachew/cbe-super-app-cps-action/internal/adapter/inbound/http"
	"gitlab.com/bersufekadgetachew/cbe-super-app-cps-action/internal/application/middleware"
	accountblock "gitlab.com/bersufekadgetachew/cbe-super-app-cps-action/internal/port/inbound/account_block"
)

func RegisterAccountBlockRoutes(r chi.Router, handler accountblock.AccountBlockHandler, authMiddleware middleware.AuthMiddleware) {
	routes := []sharedhttp.Route{
		{
			Method:  http.MethodPost,
			Path:    "/accountblock/filter_single",
			Handler: handler.FilterSingleBranches,
			Middlewares: []func(next http.Handler) http.Handler{
				authMiddleware.AuthenticateToken,
				authMiddleware.AccessControl([]string{"maker", "ifb-maker"}),
			},
		},
		{
			Method:  http.MethodPost,
			Path:    "/accountblock/disable_single",
			Handler: handler.DisableSingleBranch,
			Middlewares: []func(next http.Handler) http.Handler{
				authMiddleware.AuthenticateToken,
				authMiddleware.AccessControl([]string{"maker", "ifb-maker"}),
			},
		},
		{
			Method:  http.MethodPost,
			Path:    "/accountblock/approve_disable_single",
			Handler: handler.ApproveSingleBranchDisable,
			Middlewares: []func(next http.Handler) http.Handler{
				authMiddleware.AuthenticateToken,
				authMiddleware.AccessControl([]string{"checker", "ifb-checker"}),
			},
		},
		{
			Method:  http.MethodPost,
			Path:    "/accountblock/filter_multiple",
			Handler: handler.FilterMultipleBranches,
			Middlewares: []func(next http.Handler) http.Handler{
				authMiddleware.AuthenticateToken,
				authMiddleware.AccessControl([]string{"maker", "ifb-maker"}),
			},
		},
		{
			Method:  http.MethodPost,
			Path:    "/accountblock/disable_multiple",
			Handler: handler.DisableMultipleBranches,
			Middlewares: []func(next http.Handler) http.Handler{
				authMiddleware.AuthenticateToken,
				authMiddleware.AccessControl([]string{"maker", "ifb-maker"}),
			},
		},
		{
			Method:  http.MethodPost,
			Path:    "/accountblock/approve_disable_multiple",
			Handler: handler.ApproveBulkBranchesDisable,
			Middlewares: []func(next http.Handler) http.Handler{
				authMiddleware.AuthenticateToken,
				authMiddleware.AccessControl([]string{"checker", "ifb-checker"}),
			},
		},
		{
			Method:  http.MethodPost,
			Path:    "/accountblock/block_region",
			Handler: handler.BlockRegion,
			Middlewares: []func(next http.Handler) http.Handler{
				authMiddleware.AuthenticateToken,
				authMiddleware.AccessControl([]string{"maker", "ifb-maker"}),
			},
		},
		{
			Method:  http.MethodPost,
			Path:    "/accountblock/update_region",
			Handler: handler.UpdateRegion,
			Middlewares: []func(next http.Handler) http.Handler{
				authMiddleware.AuthenticateToken,
				authMiddleware.AccessControl([]string{"maker", "ifb-maker"}),
			},
		},
		{
			Method:  http.MethodPost,
			Path:    "/accountblock/approve_block_region",
			Handler: handler.ApproveRegionBlock,
			Middlewares: []func(next http.Handler) http.Handler{
				authMiddleware.AuthenticateToken,
				authMiddleware.AccessControl([]string{"checker", "ifb-checker"}),
			},
		},
		{
            Method:  http.MethodGet,
            Path:    "/accountblock/region/{id}",
            Handler: handler.GetRegionByID,
            Middlewares: []func(next http.Handler) http.Handler{
                authMiddleware.AuthenticateToken,
                authMiddleware.AccessControl([]string{"maker", "ifb-maker", "checker", "ifb-checker"}),
            },
        },

	}

	sharedhttp.RegisterRoutes(r, routes)
}
