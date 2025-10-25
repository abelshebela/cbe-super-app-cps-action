package accountblock

import (
	"cbe-super-app-cps-action/internal/constants"
	account_block "cbe-super-app-cps-action/internal/constants/interfaces/account_block"
	"cbe-super-app-cps-action/internal/glue"
	"cbe-super-app-cps-action/internal/handlers/middleware"
	"net/http"

	"github.com/go-chi/chi/v5"
)

func Init(
	router chi.Router,
	handler account_block.AccountBlockAdapter,
	authMiddleware middleware.AuthMiddleware,
) {
	router.Route("/account_block", func(r chi.Router) {
		routes := []glue.Route{
			{
				Method:  http.MethodGet,
				Path:    "/branches",
				Handler: handler.GetAllBranches,
				Middlewares: []func(http.Handler) http.Handler{
					authMiddleware.AuthenticateToken,
					authMiddleware.AccessControl([]string{constants.Maker, constants.IFBMaker}),
				},
			},
			{
				Method:  http.MethodGet,
				Path:    "/branches/{branch_code}",
				Handler: handler.GetBranchByCode,
				Middlewares: []func(http.Handler) http.Handler{
					authMiddleware.AuthenticateToken,
					authMiddleware.AccessControl([]string{constants.Maker, constants.IFBMaker}),
				},
			},
			{
				Method:  http.MethodGet,
				Path:    "/regions",
				Handler: handler.GetAllRegions,
				Middlewares: []func(http.Handler) http.Handler{
					authMiddleware.AuthenticateToken,
					authMiddleware.AccessControl([]string{constants.Maker, constants.IFBMaker, constants.Checker, constants.IFBChecker}),
				},
			},
			{
				Method:  http.MethodGet,
				Path:    "/regions/{region_code}",
				Handler: handler.GetRegionByCode,
				Middlewares: []func(http.Handler) http.Handler{
					authMiddleware.AuthenticateToken,
					authMiddleware.AccessControl([]string{constants.Maker, constants.IFBMaker, constants.Checker, constants.IFBChecker}),
				},
			},
			{
				Method:  http.MethodGet,
				Path:    "/districts",
				Handler: handler.GetAllDistricts,
				Middlewares: []func(http.Handler) http.Handler{
					authMiddleware.AuthenticateToken,
					authMiddleware.AccessControl([]string{constants.Maker, constants.IFBMaker, constants.Checker, constants.IFBChecker}),
				},
			},
			{
				Method:  http.MethodGet,
				Path:    "/districts/{district_code}",
				Handler: handler.GetDistrictByCode,
				Middlewares: []func(http.Handler) http.Handler{
					authMiddleware.AuthenticateToken,
					authMiddleware.AccessControl([]string{constants.Maker, constants.IFBMaker, constants.Checker, constants.IFBChecker}),
				},
			},
			{
				Method:  http.MethodGet,
				Path:    "/cities",
				Handler: handler.GetAllCities,
				Middlewares: []func(http.Handler) http.Handler{
					authMiddleware.AuthenticateToken,
					authMiddleware.AccessControl([]string{constants.Maker, constants.IFBMaker, constants.Checker, constants.IFBChecker}),
				},
			},
			{
				Method:  http.MethodGet,
				Path:    "/cities/{city_code}",
				Handler: handler.GetCityByCode,
				Middlewares: []func(http.Handler) http.Handler{
					authMiddleware.AuthenticateToken,
					authMiddleware.AccessControl([]string{constants.Maker, constants.IFBMaker, constants.Checker, constants.IFBChecker}),
				},
			},
			{
				Method:  http.MethodPost,
				Path:    "/branches/enable",
				Handler: handler.EnableBranches,
				Middlewares: []func(http.Handler) http.Handler{
					authMiddleware.AuthenticateToken,
					authMiddleware.AccessControl([]string{constants.Maker, constants.IFBMaker}),
				},
			},
			{
				Method:  http.MethodPost,
				Path:    "/branches/disable",
				Handler: handler.DisableBranches,
				Middlewares: []func(http.Handler) http.Handler{
					authMiddleware.AuthenticateToken,
					authMiddleware.AccessControl([]string{constants.Maker, constants.IFBMaker}),
				},
			},
			{
				Method:  http.MethodPost,
				Path:    "/regions/enable",
				Handler: handler.EnableRegions,
				Middlewares: []func(http.Handler) http.Handler{
					authMiddleware.AuthenticateToken,
					authMiddleware.AccessControl([]string{constants.Maker, constants.IFBMaker}),
				},
			},
			{
				Method:  http.MethodPost,
				Path:    "/regions/disable",
				Handler: handler.DisableRegions,
				Middlewares: []func(http.Handler) http.Handler{
					authMiddleware.AuthenticateToken,
					authMiddleware.AccessControl([]string{constants.Maker, constants.IFBMaker}),
				},
			},
			{
				Method:  http.MethodPost,
				Path:    "/districts/enable",
				Handler: handler.EnableDistricts,
				Middlewares: []func(http.Handler) http.Handler{
					authMiddleware.AuthenticateToken,
					authMiddleware.AccessControl([]string{constants.Maker, constants.IFBMaker}),
				},
			},
			{
				Method:  http.MethodPost,
				Path:    "/districts/disable",
				Handler: handler.DisableDistricts,
				Middlewares: []func(http.Handler) http.Handler{
					authMiddleware.AuthenticateToken,
					authMiddleware.AccessControl([]string{constants.Maker, constants.IFBMaker}),
				},
			},
			{
				Method:  http.MethodPost,
				Path:    "/cities/enable",
				Handler: handler.EnableCities,
				Middlewares: []func(http.Handler) http.Handler{
					authMiddleware.AuthenticateToken,
					authMiddleware.AccessControl([]string{constants.Maker, constants.IFBMaker}),
				},
			},
			{
				Method:  http.MethodPost,
				Path:    "/cities/disable",
				Handler: handler.DisableCities,
				Middlewares: []func(http.Handler) http.Handler{
					authMiddleware.AuthenticateToken,
					authMiddleware.AccessControl([]string{constants.Maker, constants.IFBMaker}),
				},
			},
		}

		glue.RegisterRoutes(r, routes)
	})
}
