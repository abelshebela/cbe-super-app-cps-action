package service

import (
	"net/http"

	route "github.com/CBE-Super-App/cbe-super-app-cps-action/internal/adapter/inbound/http"
	"github.com/CBE-Super-App/cbe-super-app-cps-action/internal/application/middleware"
	"github.com/CBE-Super-App/cbe-super-app-cps-action/internal/port/inbound"

	"github.com/go-chi/chi/v5"
)

func InitServiceRoutes(router chi.Router, serviceHandler inbound.ServiceBound, authMiddleware middleware.AuthMiddleware) {
	router.Route("/api/v1/cbesuperapp/cps_action/service", func(r chi.Router) {
		routes := []route.Route{
			{
				Method:  http.MethodPost,
				Path:    "/update",
				Handler: serviceHandler.UpdateServiceFee,
				Middlewares: []func(next http.Handler) http.Handler{
					authMiddleware.AuthenticateToken,
					authMiddleware.AccessControl([]string{"MAKER", "IFB-MAKER"}),
				},
			},
			{
				Method:  http.MethodPost,
				Path:    "/approve",
				Handler: serviceHandler.AuthorizeServiceFee,
				Middlewares: []func(next http.Handler) http.Handler{
					authMiddleware.AuthenticateToken,
					authMiddleware.AccessControl([]string{"CHECKER", "IFB-CHECKER"}),
				},
			},
			{
				Method:  http.MethodPost,
				Path:    "/reject",
				Handler: serviceHandler.RejectServiceFee,
				Middlewares: []func(next http.Handler) http.Handler{
					authMiddleware.AuthenticateToken,
					authMiddleware.AccessControl([]string{"CHECKER", "IFB-CHECKER"}),
				},
			},
		}

		route.RegisterRoutes(r, routes)
	})
}
