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
				Path:    "/branches",
				Handler: handler.GetAllBranches,
				Middlewares: []func(http.Handler) http.Handler{
					authMiddleware.AuthenticateToken,
					authMiddleware.AccessControl([]string{role.Maker, role.IFBMaker}),
				},
			},
			{
				Method:  http.MethodGet,
				Path:    "/branches/{branch_code}",
				Handler: handler.GetBranch,
				Middlewares: []func(http.Handler) http.Handler{
					authMiddleware.AuthenticateToken,
					authMiddleware.AccessControl([]string{role.Maker, role.IFBMaker}),
				},
			},
			{
				Method:  http.MethodGet,
				Path:    "/regions{region_code}",
				Handler: handler.GetRegionByCode,
				Middlewares: []func(http.Handler) http.Handler{
					authMiddleware.AuthenticateToken,
					authMiddleware.AccessControl([]string{role.Maker, role.IFBMaker, role.Checker, role.IFBChecker}),
				},
			},
			{
				Method:  http.MethodGet,
				Path:    "/regions",
				Handler: handler.GetAllRegions,
				Middlewares: []func(http.Handler) http.Handler{
					authMiddleware.AuthenticateToken,
					authMiddleware.AccessControl([]string{role.Maker, role.IFBMaker, role.Checker, role.IFBChecker}),
				},
			},
			{
				Method:  http.MethodGet,
				Path:    "/districts{district_code}",
				Handler: handler.GetDistrictByCode,
				Middlewares: []func(http.Handler) http.Handler{
					authMiddleware.AuthenticateToken,
					authMiddleware.AccessControl([]string{role.Maker, role.IFBMaker, role.Checker, role.IFBChecker}),
				},
			},
			{
				Method:  http.MethodGet,
				Path:    "/districts",
				Handler: handler.GetAllDistricts,
				Middlewares: []func(http.Handler) http.Handler{
					authMiddleware.AuthenticateToken,
					authMiddleware.AccessControl([]string{role.Maker, role.IFBMaker, role.Checker, role.IFBChecker}),
				},
			},
			{
				Method:  http.MethodGet,
				Path:    "/cities{city_code}",
				Handler: handler.GetCityByCode,
				Middlewares: []func(http.Handler) http.Handler{
					authMiddleware.AuthenticateToken,
					authMiddleware.AccessControl([]string{role.Maker, role.IFBMaker, role.Checker, role.IFBChecker}),
				},
			},
			{
				Method:  http.MethodGet,
				Path:    "/cities",
				Handler: handler.GetAllCities,
				Middlewares: []func(http.Handler) http.Handler{
					authMiddleware.AuthenticateToken,
					authMiddleware.AccessControl([]string{role.Maker, role.IFBMaker, role.Checker, role.IFBChecker}),
				},
			},
			{
				Method:  http.MethodPost,
				Path:    "/branches/enable",
				Handler: handler.EnableBranches,
				Middlewares: []func(http.Handler) http.Handler{
					authMiddleware.AuthenticateToken,
					authMiddleware.AccessControl([]string{role.Maker, role.IFBMaker}),
				},
			},
			{
				Method:  http.MethodPost,
				Path:    "/branches/disable",
				Handler: handler.DisableBranches,
				Middlewares: []func(http.Handler) http.Handler{
					authMiddleware.AuthenticateToken,
					authMiddleware.AccessControl([]string{role.Maker, role.IFBMaker}),
				},
			},
			{
				Method:  http.MethodPost,
				Path:    "/regions/enable",
				Handler: handler.EnableRegion,
				Middlewares: []func(http.Handler) http.Handler{
					authMiddleware.AuthenticateToken,
					authMiddleware.AccessControl([]string{role.Maker, role.IFBMaker}),
				},
			},
			{
				Method:  http.MethodPost,
				Path:    "/regions/disable",
				Handler: handler.DisableRegion,
				Middlewares: []func(http.Handler) http.Handler{
					authMiddleware.AuthenticateToken,
					authMiddleware.AccessControl([]string{role.Maker, role.IFBMaker}),
				},
			},
			{
				Method:  http.MethodPost,
				Path:    "/districts/enable",
				Handler: handler.EnableDistrict,
				Middlewares: []func(http.Handler) http.Handler{
					authMiddleware.AuthenticateToken,
					authMiddleware.AccessControl([]string{role.Maker, role.IFBMaker}),
				},
			},
			{
				Method:  http.MethodPost,
				Path:    "/districts/disable",
				Handler: handler.DisableDistrict,
				Middlewares: []func(http.Handler) http.Handler{
					authMiddleware.AuthenticateToken,
					authMiddleware.AccessControl([]string{role.Maker, role.IFBMaker}),
				},
			},
			{
				Method:  http.MethodPost,
				Path:    "/cities/enable",
				Handler: handler.EnableCity,
				Middlewares: []func(http.Handler) http.Handler{
					authMiddleware.AuthenticateToken,
					authMiddleware.AccessControl([]string{role.Maker, role.IFBMaker}),
				},
			},
			{
				Method:  http.MethodPost,
				Path:    "/cities/disable",
				Handler: handler.DisableCity,
				Middlewares: []func(http.Handler) http.Handler{
					authMiddleware.AuthenticateToken,
					authMiddleware.AccessControl([]string{role.Maker, role.IFBMaker}),
				},
			},
		}

		sharedhttp.RegisterRoutes(r, routes)
	})
}
