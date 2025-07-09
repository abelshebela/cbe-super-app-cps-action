package accountblock_handler

import (
	"net/http"

	sharedhttp "github.com/CBE-Super-App/cbe-super-app-cps-action/internal/adapter/inbound/http"
	"github.com/CBE-Super-App/cbe-super-app-cps-action/internal/application/middleware"
	accountblock "github.com/CBE-Super-App/cbe-super-app-cps-action/internal/port/inbound/account_block"

	"github.com/go-chi/chi/v5"

	role "github.com/CBE-Super-App/cbe-super-app-cps-action/utils"
)

func RegisterAccountBlockRoutes(
	router chi.Router,
	handler accountblock.AccountBlockHandler,
	authMiddleware middleware.AuthMiddleware,
) {
	router.Route("/api/v1/cbesuperapp/cps_action/account_block", func(r chi.Router) {
		routes := []sharedhttp.Route{
			{
				Method:  http.MethodPost,
				Path:    "/branch/filter_single",
				Handler: handler.FilterSingleBranches,
				Middlewares: []func(http.Handler) http.Handler{
					authMiddleware.AuthenticateToken,
					authMiddleware.AccessControl([]string{role.Maker, role.IFBMaker}),
				},
			},
			{
				Method:  http.MethodPost,
				Path:    "/branch/disable_single",
				Handler: handler.DisableSingleBranch,
				Middlewares: []func(http.Handler) http.Handler{
					authMiddleware.AuthenticateToken,
					authMiddleware.AccessControl([]string{role.Maker, role.IFBMaker}),
				},
			},
			{
				Method:  http.MethodPost,
				Path:    "/branch/approve_disable_single",
				Handler: handler.ApproveSingleBranchDisable,
				Middlewares: []func(http.Handler) http.Handler{
					authMiddleware.AuthenticateToken,
					authMiddleware.AccessControl([]string{role.Checker, role.IFBChecker}),
				},
			},
			{
				Method:  http.MethodPost,
				Path:    "/branch/filter_multiple",
				Handler: handler.FilterMultipleBranches,
				Middlewares: []func(http.Handler) http.Handler{
					authMiddleware.AuthenticateToken,
					authMiddleware.AccessControl([]string{role.Maker, role.IFBMaker}),
				},
			},
			{
				Method:  http.MethodPost,
				Path:    "/branch/disable_multiple",
				Handler: handler.DisableMultipleBranches,
				Middlewares: []func(http.Handler) http.Handler{
					authMiddleware.AuthenticateToken,
					authMiddleware.AccessControl([]string{role.Maker, role.IFBMaker}),
				},
			},
			{
				Method:  http.MethodPost,
				Path:    "/branch/approve_disable_multiple",
				Handler: handler.ApproveBulkBranchesDisable,
				Middlewares: []func(http.Handler) http.Handler{
					authMiddleware.AuthenticateToken,
					authMiddleware.AccessControl([]string{role.Checker, role.IFBChecker}),
				},
			},
			{
				Method:  http.MethodPost,
				Path:    "/region/block",
				Handler: handler.BlockRegion,
				Middlewares: []func(http.Handler) http.Handler{
					authMiddleware.AuthenticateToken,
					authMiddleware.AccessControl([]string{role.Maker, role.IFBMaker}),
				},
			},
			{
				Method:  http.MethodPut,
				Path:    "/region/update",
				Handler: handler.UpdateRegion,
				Middlewares: []func(http.Handler) http.Handler{
					authMiddleware.AuthenticateToken,
					authMiddleware.AccessControl([]string{role.Maker, role.IFBMaker}),
				},
			},
			{
				Method:  http.MethodPost,
				Path:    "/region/approve_block",
				Handler: handler.ApproveRegionBlock,
				Middlewares: []func(http.Handler) http.Handler{
					authMiddleware.AuthenticateToken,
					authMiddleware.AccessControl([]string{role.Checker, role.IFBChecker}),
				},
			},
			{
				Method:  http.MethodGet,
				Path:    "/region/{id}",
				Handler: handler.GetRegionByID,
				Middlewares: []func(http.Handler) http.Handler{
					authMiddleware.AuthenticateToken,
					authMiddleware.AccessControl([]string{role.Maker, role.IFBMaker, role.Checker, role.IFBChecker}),
				},
			},
			{
				Method:  http.MethodPost,
				Path:    "/district/block",
				Handler: handler.BlockDistrict,
				Middlewares: []func(http.Handler) http.Handler{
					authMiddleware.AuthenticateToken,
					authMiddleware.AccessControl([]string{role.Maker, role.IFBMaker}),
				},
			},
			{
				Method:  http.MethodPost,
				Path:    "/district/approve_block",
				Handler: handler.ApproveBlockDistrict,
				Middlewares: []func(http.Handler) http.Handler{
					authMiddleware.AuthenticateToken,
					authMiddleware.AccessControl([]string{role.Checker, role.IFBChecker}),
				},
			},
			{
				Method:  http.MethodGet,
				Path:    "/district/{id}",
				Handler: handler.GetDistrictByID,
				Middlewares: []func(http.Handler) http.Handler{
					authMiddleware.AuthenticateToken,
					authMiddleware.AccessControl([]string{role.Maker, role.IFBMaker, role.Checker, role.IFBChecker}),
				},
			},
			{
				Method:  http.MethodPost,
				Path:    "/city/block",
				Handler: handler.BlockCity,
				Middlewares: []func(http.Handler) http.Handler{
					authMiddleware.AuthenticateToken,
					authMiddleware.AccessControl([]string{role.Maker, role.IFBMaker}),
				},
			},
			{
				Method:  http.MethodPost,
				Path:    "/city/approve_block",
				Handler: handler.ApproveBlockCity,
				Middlewares: []func(http.Handler) http.Handler{
					authMiddleware.AuthenticateToken,
					authMiddleware.AccessControl([]string{role.Checker, role.IFBChecker}),
				},
			},
			{
				Method:  http.MethodGet,
				Path:    "/city/{id}",
				Handler: handler.GetCityByID,
				Middlewares: []func(http.Handler) http.Handler{
					authMiddleware.AuthenticateToken,
					authMiddleware.AccessControl([]string{role.Maker, role.IFBMaker, role.Checker, role.IFBChecker}),
				},
			},
			{
				Method:  http.MethodPost,
				Path:    "/user/block",
				Handler: handler.BlockUser,
				Middlewares: []func(http.Handler) http.Handler{
					authMiddleware.AuthenticateToken,
					authMiddleware.AccessControl([]string{role.Maker, role.IFBMaker}),
				},
			},
			{
				Method:  http.MethodPost,
				Path:    "/user/approve_block",
				Handler: handler.ApproveBlockUser,
				Middlewares: []func(http.Handler) http.Handler{
					authMiddleware.AuthenticateToken,
					authMiddleware.AccessControl([]string{role.Checker, role.IFBChecker}),
				},
			},
			{
				Method:  http.MethodGet,
				Path:    "/user/{id}",
				Handler: handler.GetUserByID,
				Middlewares: []func(http.Handler) http.Handler{
					authMiddleware.AuthenticateToken,
					authMiddleware.AccessControl([]string{role.Maker, role.IFBMaker, role.Checker, role.IFBChecker}),
				},
			},
		}

		sharedhttp.RegisterRoutes(r, routes)
	})
}
