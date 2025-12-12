package services

import (
	"net/http"

	"cbe-super-app-cps-action/internal/constants"
	inbound "cbe-super-app-cps-action/internal/constants/interfaces/services"
	"cbe-super-app-cps-action/internal/glue"
	"cbe-super-app-cps-action/internal/handlers/middleware"

	"github.com/go-chi/chi/v5"
)

func Init(router chi.Router, handler inbound.ServicesHandler, auth middleware.AuthMiddleware) {
	routes := []glue.Route{
		{
			Method:  http.MethodPost,
			Path:    "/services",
			Handler: handler.Create,
			Middlewares: []func(next http.Handler) http.Handler{
				auth.AuthenticateToken,
				auth.AccessControl([]string{constants.Maker, constants.IFBMaker}),
			},
		},
		{
			Method:  http.MethodPatch,
			Path:    "/services/{id}",
			Handler: handler.Update,
			Middlewares: []func(next http.Handler) http.Handler{
				auth.AuthenticateToken,
				auth.AccessControl([]string{constants.Maker, constants.IFBMaker}),
			},
		},
		{
			Method:  http.MethodPatch,
			Path:    "/services/{id}/enable",
			Handler: handler.Enable,
			Middlewares: []func(next http.Handler) http.Handler{
				auth.AuthenticateToken,
				auth.AccessControl([]string{constants.Maker, constants.IFBMaker}),
			},
		},
		{
			Method:  http.MethodPatch,
			Path:    "/services/{id}/disable",
			Handler: handler.Disable,
			Middlewares: []func(next http.Handler) http.Handler{
				auth.AuthenticateToken,
				auth.AccessControl([]string{constants.Maker, constants.IFBMaker}),
			},
		},
		{
			Method:  http.MethodGet,
			Path:    "/services/{id}",
			Handler: handler.GetByID,
			Middlewares: []func(next http.Handler) http.Handler{
				auth.AuthenticateToken,
				auth.AccessControl([]string{constants.Maker, constants.IFBMaker, constants.Checker, constants.IFBChecker}),
			},
		},
		{
			Method:  http.MethodGet,
			Path:    "/services",
			Handler: handler.GetAll,
			Middlewares: []func(next http.Handler) http.Handler{
				auth.AuthenticateToken,
				auth.AccessControl([]string{constants.Maker, constants.IFBMaker, constants.Checker, constants.IFBChecker}),
			},
		},
	}
	glue.RegisterRoutes(router, routes)
}
