// Package event provides HTTP handlers and routing for event-related endpoints.
package event

import (
	"net/http"

	"github.com/go-chi/chi/v5"

	"github.com/CBE-Super-App/cbe-super-app-cps-action/internal/application/middleware"
	route "github.com/CBE-Super-App/cbe-super-app-cps-action/internal/adapter/inbound/http"
	event_inbound "github.com/CBE-Super-App/cbe-super-app-cps-action/internal/port/inbound/event"
	role "github.com/CBE-Super-App/cbe-super-app-cps-action/utils"

)

func InitEventsHandlerMaker(router chi.Router, handler event_inbound.EventHandler, authMiddleware middleware.AuthMiddleware) {
	router.Route("/events", func(r chi.Router) {
		routes := []route.Route{
			{
				Method:  http.MethodPost,
				Path:    "/create",
				Handler: handler.MakerCreateEvent,
				Middlewares: []func(next http.Handler) http.Handler{
					authMiddleware.AuthenticateToken,
					authMiddleware.AccessControl([]string{role.Maker, role.IFBMaker}),
				},
			},
			{
				Method:  http.MethodPost,
				Path:    "/review",
				Handler: handler.CheckerEvent,
				Middlewares: []func(next http.Handler) http.Handler{
					authMiddleware.AuthenticateToken,
					authMiddleware.AccessControl([]string{role.Checker, role.IFBChecker}),
				},
			},
		}

		route.RegisterRoutes(r, routes)
	})
}
