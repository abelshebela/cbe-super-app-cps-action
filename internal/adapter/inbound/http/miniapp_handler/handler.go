package miniapphandler

import (
	"net/http"

	"github.com/go-chi/chi/v5"

	route "gitlab.com/bersufekadgetachew/cbe-super-app-cps-action/internal/adapter/inbound/http"
	"gitlab.com/bersufekadgetachew/cbe-super-app-cps-action/internal/application/middleware"
	miniapp_application "gitlab.com/bersufekadgetachew/cbe-super-app-cps-action/internal/application/mini_app"
	Inbound "gitlab.com/bersufekadgetachew/cbe-super-app-cps-action/internal/port/inbound/miniapp"
)

type HttpStore struct {
	Application miniapp_application.ApplicationAbstracts
}

func NewHttpBulkService(app miniapp_application.ApplicationAbstracts) Inbound.Inbound {
	return &HttpStore{
		Application: app,
	}
}

func InitServiceHandlerMaker(router chi.Router, handler Inbound.Inbound, authMiddleware middleware.AuthMiddleware) {

	routes := []route.Route{
		{
			Method:  http.MethodPost,
			Path:    "/api/v1/cbesuperapp/cps_config/mini_app_create_maker",
			Handler: handler.MakerCreateMiniApp,
			Middlewares: []func(next http.Handler) http.Handler{
				authMiddleware.AuthenticateToken,
				authMiddleware.AccessControl([]string{"maker", "ifb-maker"}),
			},
		},
		{
			Method:  http.MethodPost,
			Path:    "/api/v1/cbesuperapp/cps_config/mini_app_create_checker",
			Handler: handler.CheckerMiniApp,
			Middlewares: []func(next http.Handler) http.Handler{
				authMiddleware.AuthenticateToken,
				authMiddleware.AccessControl([]string{"checker", "ifb-checker"}),
			},
		},
	}

	route.RegisterRoutes(router, routes)
}
