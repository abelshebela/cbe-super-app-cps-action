package unlink

import (
	"net/http"

	route "github.com/CBE-Super-App/cbe-super-app-cps-action/internal/adapter/inbound/http"
	"github.com/CBE-Super-App/cbe-super-app-cps-action/internal/application/middleware"
	Inbound "github.com/CBE-Super-App/cbe-super-app-cps-action/internal/port/inbound/unlink"
	"github.com/go-chi/chi/v5"
)

func InitUnlinkHanldler(router chi.Router, handler Inbound.UnlinkHandler, authMiddleware middleware.AuthMiddleware) {
	router.Route("/unlink/", func(r chi.Router) {
		routes := []route.Route{
			{
				Method:  http.MethodPost,
				Path:    "/fetch/user/byaccount",
				Handler: handler.GetUserByAccount,
				Middlewares: []func(next http.Handler) http.Handler{
					authMiddleware.AuthenticateToken,
					authMiddleware.AccessControl([]string{"maker", "checker"}),
				},
			},
			{
				Method:  http.MethodGet,
				Path:    "/fetch/archived/users",
				Handler: handler.GetArchivedUser,
				Middlewares: []func(next http.Handler) http.Handler{
					authMiddleware.AuthenticateToken,
					authMiddleware.AccessControl([]string{"maker", "checker"}),
				},
			},
			{
				Method:  http.MethodPost,
				Path:    "/request",
				Handler: handler.UnlinkUserCif,
				Middlewares: []func(next http.Handler) http.Handler{
					authMiddleware.AuthenticateToken,
					authMiddleware.AccessControl([]string{"maker"}),
				},
			},
		}

		route.RegisterRoutes(r, routes)
	})
}
