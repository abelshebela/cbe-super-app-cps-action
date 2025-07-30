package notification

import (
	"net/http"

	"github.com/go-chi/chi/v5"

	route "github.com/CBE-Super-App/cbe-super-app-cps-action/internal/adapter/inbound/http"
	"github.com/CBE-Super-App/cbe-super-app-cps-action/internal/application/middleware"
	cps_const "github.com/CBE-Super-App/cbe-super-app-cps-action/internal/domain/cps_actions/constant"
	notification_inbound "github.com/CBE-Super-App/cbe-super-app-cps-action/internal/port/inbound/notification"
	role "github.com/CBE-Super-App/cbe-super-app-cps-action/utils"
)

// InitNotificationsHandlerMaker sets up HTTP routes for notification-related endpoints
func InitNotificationsHandlerRoutes(router chi.Router, handler notification_inbound.NotificationHandler, authMiddleware middleware.AuthMiddleware, cpsGuard *middleware.CPSActionMiddlewareFactory) {
	router.Route("/notifications", func(r chi.Router) {
		routes := []route.Route{
			{
				Method:  http.MethodPost,
				Path:    "/",
				Handler: handler.CreateNotification,
				Middlewares: []func(next http.Handler) http.Handler{
					authMiddleware.AuthenticateToken,
					authMiddleware.AccessControl([]string{role.Maker, role.IFBMaker}),
					cpsGuard.RequireNoPendingCPSActionGuard(string(cps_const.RequestCreateNotification)),
				},
			},
			{
				Method:  http.MethodPatch,
				Path:    "/{id}",
				Handler: handler.UpdateNotification,
				Middlewares: []func(next http.Handler) http.Handler{
					authMiddleware.AuthenticateToken,
					authMiddleware.AccessControl([]string{role.Maker, role.IFBMaker}),
					cpsGuard.RequireNoPendingCPSActionGuard(string(cps_const.RequestUpdateNotification)),
				},
			},
			{
				Method:  http.MethodPatch,
				Path:    "/enable/{id}",
				Handler: handler.EnableNotification,
				Middlewares: []func(next http.Handler) http.Handler{
					authMiddleware.AuthenticateToken,
					authMiddleware.AccessControl([]string{role.Maker, role.IFBMaker}),
					cpsGuard.RequireNoPendingCPSActionGuard(string(cps_const.RequestEnableNotification)),
				},
			},
			{
				Method:  http.MethodPatch,
				Path:    "/disable/{id}",
				Handler: handler.DisableNotification,
				Middlewares: []func(next http.Handler) http.Handler{
					authMiddleware.AuthenticateToken,
					authMiddleware.AccessControl([]string{role.Maker, role.IFBMaker}),
					cpsGuard.RequireNoPendingCPSActionGuard(string(cps_const.RequestDisableNotification)),
				},
			},
			{
				Method:  http.MethodDelete,
				Path:    "/{id}",
				Handler: handler.DeleteNotification,
				Middlewares: []func(next http.Handler) http.Handler{
					authMiddleware.AuthenticateToken,
					authMiddleware.AccessControl([]string{role.Maker, role.IFBMaker}),
					cpsGuard.RequireNoPendingCPSActionGuard(string(cps_const.RequestDeleteNotification)),
				},
			},
			{
				Method:  http.MethodGet,
				Path:    "/{id}",
				Handler: handler.FetchNotificationByID,
				Middlewares: []func(next http.Handler) http.Handler{
					authMiddleware.AuthenticateToken,
					authMiddleware.AccessControl([]string{role.Checker, role.IFBChecker, role.Maker, role.IFBMaker}),
				},
			},
			{
				Method:  http.MethodGet,
				Path:    "/",
				Handler: handler.FetchNotifications,
				Middlewares: []func(next http.Handler) http.Handler{
					authMiddleware.AuthenticateToken,
					authMiddleware.AccessControl([]string{role.Checker, role.IFBChecker, role.Maker, role.IFBMaker}),
				},
			},
		}

		route.RegisterRoutes(r, routes)
	})
}