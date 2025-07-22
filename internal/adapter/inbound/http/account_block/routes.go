package accountblock_handler

import (
	"net/http"

	sharedhttp "github.com/CBE-Super-App/cbe-super-app-cps-action/internal/adapter/inbound/http"
	"github.com/CBE-Super-App/cbe-super-app-cps-action/internal/application/middleware"
	accountblock "github.com/CBE-Super-App/cbe-super-app-cps-action/internal/port/inbound/account_block"

	role "github.com/CBE-Super-App/cbe-super-app-cps-action/utils"
	"github.com/go-chi/chi/v5"
)

func RegisterAccountBlockRoutes(
	router chi.Router,
	handler accountblock.AccountBlockHandler,
	authMiddleware middleware.AuthMiddleware,
) {
	router.Route("/account_block", func(r chi.Router) {
		routes := []sharedhttp.Route{
			{
				Method:  http.MethodGet,
				Path:    "/branch/filter_single",
				Handler: handler.FilterSingleBranches,
				Middlewares: []func(http.Handler) http.Handler{
					authMiddleware.AuthenticateToken,
					authMiddleware.AccessControl([]string{role.Maker, role.IFBMaker}),
				},
			},
			{
				Method:  http.MethodGet,
				Path:    "/branch/all",
				Handler: handler.GetAllBranches,
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
				Method:  http.MethodGet,
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
				Path:    "/region",
				Handler: handler.GetRegionByCode,
				Middlewares: []func(http.Handler) http.Handler{
					authMiddleware.AuthenticateToken,
					authMiddleware.AccessControl([]string{role.Maker, role.IFBMaker, role.Checker, role.IFBChecker}),
				},
			},
			{
				Method:  http.MethodGet,
				Path:    "/region/all",
				Handler: handler.GetAllRegion,
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
				Path:    "/district",
				Handler: handler.GetDistrictByCode,
				Middlewares: []func(http.Handler) http.Handler{
					authMiddleware.AuthenticateToken,
					authMiddleware.AccessControl([]string{role.Maker, role.IFBMaker, role.Checker, role.IFBChecker}),
				},
			},
			{
				Method:  http.MethodGet,
				Path:    "/district/all",
				Handler: handler.GetAllDistrict,
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
				Path:    "/city",
				Handler: handler.GetCityByCode,
				Middlewares: []func(http.Handler) http.Handler{
					authMiddleware.AuthenticateToken,
					authMiddleware.AccessControl([]string{role.Maker, role.IFBMaker, role.Checker, role.IFBChecker}),
				},
			},
			{
				Method:  http.MethodGet,
				Path:    "/city/all",
				Handler: handler.GetAllCities,
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
				Path:    "/user",
				Handler: handler.GetUserByPhone,
				Middlewares: []func(http.Handler) http.Handler{
					authMiddleware.AuthenticateToken,
					authMiddleware.AccessControl([]string{role.Maker, role.IFBMaker, role.Checker, role.IFBChecker}),
				},
			},
		}

		sharedhttp.RegisterRoutes(r, routes)
	})
}
