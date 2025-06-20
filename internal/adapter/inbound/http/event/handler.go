package eventhandler

import (
	"net/http"

	"github.com/go-chi/chi/v5"

	route "gitlab.com/bersufekadgetachew/cbe-super-app-cps-action/internal/adapter/inbound/http"
	event_application "gitlab.com/bersufekadgetachew/cbe-super-app-cps-action/internal/application/event"
	"gitlab.com/bersufekadgetachew/cbe-super-app-cps-action/internal/application/middleware"
	event_inbound "gitlab.com/bersufekadgetachew/cbe-super-app-cps-action/internal/port/inbound/event"
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
				authMiddleware.AccessControl([]string{"maker", "ifb-maker"}),
			},
		},
		{
			Method:  http.MethodPost,
			Path:    "/api/v1/cbesuperapp/cps_config/checker_event",
			Handler: handler.CheckerEvent,
			Middlewares: []func(next http.Handler) http.Handler{
				authMiddleware.AuthenticateToken,
				authMiddleware.AccessControl([]string{"checker", "ifb-checker"}),
			},
		},
	}

	route.RegisterRoutes(router, routes)
}
