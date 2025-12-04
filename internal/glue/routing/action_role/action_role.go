package actionrole_routing

import (
	"cbe-super-app-cps-action/internal/constants"
	actionrole_inbound "cbe-super-app-cps-action/internal/constants/interfaces/action_role"
	"cbe-super-app-cps-action/internal/glue"
	"cbe-super-app-cps-action/internal/handlers/middleware"
	"net/http"

	"github.com/go-chi/chi/v5"
)

func Init(router chi.Router, handler actionrole_inbound.ActionRoleHandler, auth middleware.AuthMiddleware) {
	routes := []glue.Route{
		{
			Method:  http.MethodGet,
			Path:    "/action-roles",
			Handler: handler.GetAll,
			Middlewares: []func(next http.Handler) http.Handler{
				auth.AuthenticateToken,
				auth.AccessControl([]string{constants.Maker, constants.IFBMaker, constants.Checker, constants.IFBChecker}),
			},
		},
		{
			Method:  http.MethodGet,
			Path:    "/action-roles/{code}",
			Handler: handler.GetByActionCode,
			Middlewares: []func(next http.Handler) http.Handler{
				auth.AuthenticateToken,
				auth.AccessControl([]string{constants.Maker, constants.IFBMaker, constants.Checker, constants.IFBChecker}),
			},
		},
		{
			Method:  http.MethodPost,
			Path:    "/action-roles",
			Handler: handler.Create,
			Middlewares: []func(next http.Handler) http.Handler{
				auth.AuthenticateToken,
				auth.AccessControl([]string{constants.Maker, constants.IFBMaker}),
			},
		},
		{
			Method:  http.MethodPatch,
			Path:    "/action-roles/{code}",
			Handler: handler.Update,
			Middlewares: []func(next http.Handler) http.Handler{
				auth.AuthenticateToken,
				auth.AccessControl([]string{constants.Maker, constants.IFBMaker}),
			},
		},
		{
			Method:  http.MethodPatch,
			Path:    "/action-roles/{code}/enable",
			Handler: handler.Enable,
			Middlewares: []func(next http.Handler) http.Handler{
				auth.AuthenticateToken,
				auth.AccessControl([]string{constants.Maker, constants.IFBMaker}),
			},
		},
		{
			Method:  http.MethodPatch,
			Path:    "/action-roles/{code}/disable",
			Handler: handler.Disable,
			Middlewares: []func(next http.Handler) http.Handler{
				auth.AuthenticateToken,
				auth.AccessControl([]string{constants.Maker, constants.IFBMaker}),
			},
		},
	}
	glue.RegisterRoutes(router, routes)
}
