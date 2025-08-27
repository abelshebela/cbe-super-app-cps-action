package amount_based_auth

import (
	"cbe-super-app-cps-action/internal/constants"
	amount_based "cbe-super-app-cps-action/internal/constants/interfaces/amount_based_auth"
	"cbe-super-app-cps-action/internal/glue"
	"cbe-super-app-cps-action/internal/handlers/middleware"
	"net/http"

	"github.com/go-chi/chi/v5"
)

func Init(router chi.Router, handler amount_based.AmountBasedAuthAdapter, authMiddleware middleware.AuthMiddleware) {


	routes := []glue.Route{
			{
				Method:  http.MethodPatch,
				Path:    "/update/{id}",
				Handler: handler.UpdateAmountBasedAuth,
				Middlewares: []func(next http.Handler) http.Handler{
					authMiddleware.AuthenticateToken,
					authMiddleware.AccessControl([]string{constants.Maker}),
				},
			},
			{
				Method:  http.MethodGet,
				Path:    "/",
				Handler: handler.GetAllAmountBasedAuth,
				Middlewares: []func(next http.Handler) http.Handler{
					authMiddleware.AuthenticateToken,
					authMiddleware.AccessControl([]string{constants.Maker}),
				},
			},
			// {
			// 	Method:  http.MethodPatch,
			// 	Path:    "/approve/{id}",
			// 	Handler: handler.ApproveAmountBasedAuth,
			// 	Middlewares: []func(next http.Handler) http.Handler{
			// 		authMiddleware.AuthenticateToken,
			// 		authMiddleware.AccessControl([]string{role.Checker}),
			// 	},
			// },
			// {
			// 	Method:  http.MethodPatch,
			// 	Path:    "/reject/{id}",
			// 	Handler: handler.RejectAmountBasedAuth,
			// 	Middlewares: []func(next http.Handler) http.Handler{
			// 		authMiddleware.AuthenticateToken,
			// 		authMiddleware.AccessControl([]string{role.Checker}),
			// 	},
			// },
		}

	glue.RegisterRoutes(router, routes)
}