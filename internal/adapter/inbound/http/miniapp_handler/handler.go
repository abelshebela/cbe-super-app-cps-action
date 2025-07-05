package miniapphandler

import (
	"net/http"

	"github.com/go-chi/chi/v5"

	route "github.com/CBE-Super-App/cbe-super-app-cps-action/internal/adapter/inbound/http"
	"github.com/CBE-Super-App/cbe-super-app-cps-action/internal/application/middleware"
	miniapp_application "github.com/CBE-Super-App/cbe-super-app-cps-action/internal/application/mini_app"
	Inbound "github.com/CBE-Super-App/cbe-super-app-cps-action/internal/port/inbound/miniapp"
	role "github.com/CBE-Super-App/cbe-super-app-cps-action/utils"
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
				authMiddleware.AccessControl([]string{role.Maker, role.IFBMaker}),
			},
		},
		{
			Method:  http.MethodPost,
			Path:    "/api/v1/cbesuperapp/cps_config/mini_app_create_checker",
			Handler: handler.CheckerMiniApp,
			Middlewares: []func(next http.Handler) http.Handler{
				authMiddleware.AuthenticateToken,
				authMiddleware.AccessControl([]string{role.Checker, role.IFBChecker}),
			},
		},
	}

	route.RegisterRoutes(router, routes)
}
