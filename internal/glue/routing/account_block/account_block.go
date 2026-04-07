package accountblock

import (
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
				},
			},
			{
				Method:  http.MethodGet,
				Path:    "/branches/{branch_id}",
				Handler: handler.GetBranchById,
				Middlewares: []func(http.Handler) http.Handler{
					authMiddleware.AuthenticateToken,
				},
			},
			{
				Method:  http.MethodGet,
				Path:    "/regions",
				Handler: handler.GetAllRegions,
				Middlewares: []func(http.Handler) http.Handler{
					authMiddleware.AuthenticateToken,
				},
			},
			{
				Method:  http.MethodGet,
				Path:    "/regions/{region_id}",
				Handler: handler.GetRegionById,
				Middlewares: []func(http.Handler) http.Handler{
					authMiddleware.AuthenticateToken,
				},
			},
			{
				Method:  http.MethodGet,
				Path:    "/districts",
				Handler: handler.GetAllDistricts,
				Middlewares: []func(http.Handler) http.Handler{
					authMiddleware.AuthenticateToken,
				},
			},
			{
				Method:  http.MethodGet,
				Path:    "/districts/{district_id}",
				Handler: handler.GetDistrictById,
				Middlewares: []func(http.Handler) http.Handler{
					authMiddleware.AuthenticateToken,
				},
			},
			{
				Method:  http.MethodGet,
				Path:    "/cities",
				Handler: handler.GetAllCities,
				Middlewares: []func(http.Handler) http.Handler{
					authMiddleware.AuthenticateToken,
				},
			},
			{
				Method:  http.MethodGet,
				Path:    "/cities/{city_id}",
				Handler: handler.GetCityById,
				Middlewares: []func(http.Handler) http.Handler{
					authMiddleware.AuthenticateToken,
				},
			},
			{
				Method:  http.MethodPost,
				Path:    "/branches/enable",
				Handler: handler.EnableBranches,
				Middlewares: []func(http.Handler) http.Handler{
					authMiddleware.AuthenticateToken,
				},
			},
			{
				Method:  http.MethodPost,
				Path:    "/branches/disable",
				Handler: handler.DisableBranches,
				Middlewares: []func(http.Handler) http.Handler{
					authMiddleware.AuthenticateToken,
				},
			},
			{
				Method:  http.MethodPost,
				Path:    "/regions/enable",
				Handler: handler.EnableRegions,
				Middlewares: []func(http.Handler) http.Handler{
					authMiddleware.AuthenticateToken,
				},
			},
			{
				Method:  http.MethodPost,
				Path:    "/regions/disable",
				Handler: handler.DisableRegions,
				Middlewares: []func(http.Handler) http.Handler{
					authMiddleware.AuthenticateToken,
				},
			},
			{
				Method:  http.MethodPost,
				Path:    "/districts/enable",
				Handler: handler.EnableDistricts,
				Middlewares: []func(http.Handler) http.Handler{
					authMiddleware.AuthenticateToken,
				},
			},
			{
				Method:  http.MethodPost,
				Path:    "/districts/disable",
				Handler: handler.DisableDistricts,
				Middlewares: []func(http.Handler) http.Handler{
					authMiddleware.AuthenticateToken,
				},
			},
			{
				Method:  http.MethodPost,
				Path:    "/cities/enable",
				Handler: handler.EnableCities,
				Middlewares: []func(http.Handler) http.Handler{
					authMiddleware.AuthenticateToken,
				},
			},
			{
				Method:  http.MethodPost,
				Path:    "/cities/disable",
				Handler: handler.DisableCities,
				Middlewares: []func(http.Handler) http.Handler{
					authMiddleware.AuthenticateToken,
				},
			},
			{
				Method:  http.MethodGet,
				Path:    "/details/{id}",
				Handler: handler.GetAccountBlockDetails,
				Middlewares: []func(http.Handler) http.Handler{
					authMiddleware.AuthenticateToken,
				},
			},
		}

		glue.RegisterRoutes(r, routes)
	})
}
