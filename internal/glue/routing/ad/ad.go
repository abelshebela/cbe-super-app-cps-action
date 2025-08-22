package ad

import (
	"cbe-super-app-cps-action/internal/constants"
	"cbe-super-app-cps-action/internal/constants/interfaces/ad"
	"cbe-super-app-cps-action/internal/glue"
	"cbe-super-app-cps-action/internal/handlers/middleware"
	// local_util "cbe-super-app-cps-action/pkgs/utils"
	"net/http"

	"github.com/go-chi/chi/v5"
)

func Init(router chi.Router, ad ad.ADAdapter, middleware middleware.AuthMiddleware ) {
		routes := []glue.Route{
			{
				Method:  http.MethodPost,
				Path:    "/adverts/",
				Handler: ad.CreateAdvert,
				Middlewares: []func(next http.Handler) http.Handler{
					middleware.AuthenticateToken,
					middleware.AccessControl([]string{constants.Maker, constants.IFBMaker}),
					middleware.RequireFormContentType(),
				},
			},
			{
				Method:  http.MethodPatch,
				Path:    "/adverts/{id}",
				Handler: ad.UpdateAdvert,
				Middlewares: []func(next http.Handler) http.Handler{
					middleware.AuthenticateToken,
					middleware.AccessControl([]string{constants.Maker, constants.IFBMaker}),
				},
			},
			{
				Method:  http.MethodDelete,
				Path:    "/adverts/{id}",
				Handler: ad.DeleteAdvert,
				Middlewares: []func(next http.Handler) http.Handler{
					middleware.AuthenticateToken,
					middleware.AccessControl([]string{constants.Maker, constants.IFBMaker}),
				},
			},
			{
				Method:  http.MethodGet,
				Path:    "/adverts/{id}",
				Handler: ad.FetchAdvertByID,
				Middlewares: []func(next http.Handler) http.Handler{
					middleware.AuthenticateToken,
					middleware.AccessControl([]string{constants.Maker, constants.IFBMaker, constants.Checker, constants.IFBChecker}),
				},
			},
			{
				Method:  http.MethodGet,
				Path:    "/adverts/",
				Handler: ad.FetchAdverts,
				Middlewares: []func(next http.Handler) http.Handler{
					middleware.AuthenticateToken,
					middleware.AccessControl([]string{constants.Maker, constants.IFBMaker, constants.Checker, constants.IFBChecker}),
				},
			},
			{
				Method:  http.MethodPatch,
				Path:    "/adverts/enable/{id}",
				Handler: ad.EnableAdvert,
				Middlewares: []func(next http.Handler) http.Handler{
					middleware.AuthenticateToken,
					middleware.AccessControl([]string{constants.Maker, constants.IFBMaker}),
				},
			},
			{
				Method:  http.MethodPatch,
				Path:    "/adverts/disable/{id}",
				Handler: ad.DisableAdvert,
				Middlewares: []func(next http.Handler) http.Handler{
					middleware.AuthenticateToken,
					middleware.AccessControl([]string{constants.Maker, constants.IFBMaker}),
				},
			},
		}

		glue.RegisterRoutes(router, routes)
	
}
