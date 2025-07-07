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

func NewHttpBulkService(app event_application.ApplicationAbstracts) event_inbound.Inbound {
	return &HttpStore{
		Application: app,
	}
}

func InitServiceHandlerMaker(router chi.Router, handler event_inbound.Inbound, authMiddleware middleware.AuthMiddleware) {

	routes := []route.Route{
		{
			Method:  http.MethodPost,
			Path:    "/api/v1/cbesuperapp/cps_config/maker_create_event",
			Handler: handler.MakerCreateEvent,
			Middlewares: []func(next http.Handler) http.Handler{
				authMiddleware.AuthenticateToken,
				authMiddleware.AccessControl([]string{role.Maker, role.IFBMaker}),
			},
		},
		{
			Method:  http.MethodPost,
			Path:    "/api/v1/cbesuperapp/cps_config/checker_event",
			Handler: handler.CheckerEvent,
			Middlewares: []func(next http.Handler) http.Handler{
				authMiddleware.AuthenticateToken,
				authMiddleware.AccessControl([]string{role.Checker, role.IFBChecker}),
			},
		},
	}

	route.RegisterRoutes(router, routes)
}
