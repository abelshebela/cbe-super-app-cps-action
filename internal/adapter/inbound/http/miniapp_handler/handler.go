package miniapphandler

import (
	"net/http"

	"github.com/go-chi/chi/v5"

	route "github.com/CBE-Super-App/cbe-super-app-cps-action/internal/adapter/inbound/http"
	"github.com/CBE-Super-App/cbe-super-app-cps-action/internal/application/middleware"
	miniapp_application "github.com/CBE-Super-App/cbe-super-app-cps-action/internal/application/mini_app"
	Inbound "github.com/CBE-Super-App/cbe-super-app-cps-action/internal/port/inbound/miniapp"
	role "github.com/CBE-Super-App/cbe-super-app-cps-action/utils"
	util "gitlab.com/bersufekadgetachew/cbe-super-app-shared/shared/utils"
)

type HttpStore struct {
	Application miniapp_application.ApplicationAbstracts
	logger      util.Logger
}

func NewMiniAppAdapter(app miniapp_application.ApplicationAbstracts, logger util.Logger) Inbound.MiniAppInbound {
	return &HttpStore{
		Application: app,
		logger:      logger,
	}
}

func InitMiniAppHandlerMaker(router chi.Router, handler Inbound.MiniAppInbound, authMiddleware middleware.AuthMiddleware) {
	router.Route("/mini-apps", func(r chi.Router) {
		routes := []route.Route{
			{
				Method:  http.MethodPost,
				Path:    "/",
				Handler: handler.MakerCreateMiniApp,
				Middlewares: []func(next http.Handler) http.Handler{
					authMiddleware.AuthenticateToken,
					authMiddleware.AccessControl([]string{role.Maker, role.IFBMaker}),
				},
			},
			// {
			// 	Method:  http.MethodPost,
			// 	Path:    "/approve_or_reject",
			// 	Handler: handler.CheckerMiniApp,
			// 	Middlewares: []func(next http.Handler) http.Handler{
			// 		authMiddleware.AuthenticateToken,
			// 		authMiddleware.AccessControl([]string{role.Checker, role.IFBChecker}),
			// 	},
			// },
			{
				Method:  http.MethodPatch,
				Path:    "/update",
				Handler: handler.MakerUpdateMiniApp,
				Middlewares: []func(next http.Handler) http.Handler{
					authMiddleware.AuthenticateToken,
					authMiddleware.AccessControl([]string{role.Maker, role.IFBMaker}),
				},
			},
			{
				Method:  http.MethodPatch,
				Path:    "/delete/{id}",
				Handler: handler.MakerDeleteMiniApp,
				Middlewares: []func(next http.Handler) http.Handler{
					authMiddleware.AuthenticateToken,
					authMiddleware.AccessControl([]string{role.Maker, role.IFBMaker}),
				},
			},
			{
				Method:  http.MethodGet,
				Path:    "/list",
				Handler: handler.ListMiniApp,
				Middlewares: []func(next http.Handler) http.Handler{
					authMiddleware.AuthenticateToken,
					authMiddleware.AccessControl([]string{role.Maker, role.IFBMaker, role.Checker, role.IFBChecker}),
				},
			},
			{
				Method:  http.MethodGet,
				Path:    "/detail/{id}",
				Handler: handler.DetailMiniAppByID,
				Middlewares: []func(next http.Handler) http.Handler{
					authMiddleware.AuthenticateToken,
					authMiddleware.AccessControl([]string{role.Maker, role.IFBMaker, role.Checker, role.IFBChecker}),
				},
			},
		}

		route.RegisterRoutes(r, routes)
	})
}
