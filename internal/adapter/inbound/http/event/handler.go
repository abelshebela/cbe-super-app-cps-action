package eventhandler

import (
	"net/http"

	"github.com/go-chi/chi/v5"

	route "github.com/CBE-Super-App/cbe-super-app-cps-action/internal/adapter/inbound/http"
	event_application "github.com/CBE-Super-App/cbe-super-app-cps-action/internal/application/event"
	"github.com/CBE-Super-App/cbe-super-app-cps-action/internal/application/middleware"
	event_inbound "github.com/CBE-Super-App/cbe-super-app-cps-action/internal/port/inbound/event"
	role "github.com/CBE-Super-App/cbe-super-app-cps-action/utils"
)

type HttpStore struct {
	Application event_application.ApplicationAbstracts
}

func NewHTTPBulkService(app event_application.ApplicationAbstracts) event_inbound.EventHandler {
	return &HttpStore{
		Application: app,
	}
}

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
