package ad

import (
	"net/http"

	"cbe-super-app-cps-action/internal/constants"
	ad "cbe-super-app-cps-action/internal/constants/interfaces/ad"
	"cbe-super-app-cps-action/internal/glue"
	"cbe-super-app-cps-action/internal/handlers/middleware"

	"github.com/go-chi/chi/v5"
)

func Init(router chi.Router, handler ad.ADAdapter, authMiddleware middleware.AuthMiddleware) {

	routes := []glue.Route{
		{
			Method:  http.MethodPost,
			Path:    "/advert",
			Handler: handler.CreateAdvert,
			Middlewares: []func(next http.Handler) http.Handler{
				authMiddleware.AuthenticateToken,
				authMiddleware.AccessControl([]string{constants.Maker, constants.IFBMaker}),
			},
		},
		{
			Method:  http.MethodPatch,
			Path:    "/advert/{id}",
			Handler: handler.UpdateAdvert,
			Middlewares: []func(next http.Handler) http.Handler{
				authMiddleware.AuthenticateToken,
				authMiddleware.AccessControl([]string{constants.Maker, constants.IFBMaker}),
			},
		},
		{
			Method:  http.MethodDelete,
			Path:    "/adverts/{id}",
			Handler: handler.DeleteAdvert,
			Middlewares: []func(next http.Handler) http.Handler{
				authMiddleware.AuthenticateToken,
				authMiddleware.AccessControl([]string{constants.Maker, constants.IFBMaker}),
			},
		},
		{
			Method:  http.MethodGet,
			Path:    "/advert/{id}",
			Handler: handler.FetchAdvertByID,
			Middlewares: []func(next http.Handler) http.Handler{
				authMiddleware.AuthenticateToken,
				authMiddleware.AccessControl([]string{constants.Maker, constants.IFBMaker, constants.Checker, constants.IFBChecker}),
			},
		},
		{
			Method:  http.MethodGet,
			Path:    "/advert",
			Handler: handler.FetchAdverts,
			Middlewares: []func(next http.Handler) http.Handler{
				authMiddleware.AuthenticateToken,
				authMiddleware.AccessControl([]string{constants.Maker, constants.IFBMaker, constants.Checker, constants.IFBChecker}),
			},
		},
		{
			Method:  http.MethodPatch,
			Path:    "/advert/enable/{id}",
			Handler: handler.EnableAdvert,
			Middlewares: []func(next http.Handler) http.Handler{
				authMiddleware.AuthenticateToken,
				authMiddleware.AccessControl([]string{constants.Maker, constants.IFBMaker}),
			},
		},
		{
			Method:  http.MethodPatch,
			Path:    "/advert/disable/{id}",
			Handler: handler.DisableAdvert,
			Middlewares: []func(next http.Handler) http.Handler{
				authMiddleware.AuthenticateToken,
				authMiddleware.AccessControl([]string{constants.Maker, constants.IFBMaker}),
			},
		},
	}

	glue.RegisterRoutes(router, routes)

}
