package unlink

import (
	"net/http"

	"cbe-super-app-cps-action/internal/constants"
	unlink "cbe-super-app-cps-action/internal/constants/interfaces/unlink"
	"cbe-super-app-cps-action/internal/glue"
	"cbe-super-app-cps-action/internal/handlers/middleware"

	"github.com/go-chi/chi/v5"
)

func Init(router chi.Router, handler unlink.UnlinkAdapter, authMiddleware middleware.AuthMiddleware) {

	routes := []glue.Route{
		{
			Method:  http.MethodGet,
			Path:    "/unlink/archived-user",
			Handler: handler.GetArchivedUser,
			Middlewares: []func(next http.Handler) http.Handler{
				authMiddleware.AuthenticateToken,
				authMiddleware.AccessControl([]string{constants.Maker, constants.Maker, constants.Checker, constants.IFBChecker}),
			},
		},
		{
			Method:  http.MethodGet,
			Path:    "/unlink/user-by-account/{account_number}",
			Handler: handler.GetUserByAccount,
			Middlewares: []func(next http.Handler) http.Handler{
				authMiddleware.AuthenticateToken,
				authMiddleware.AccessControl([]string{constants.Maker, constants.Maker, constants.Checker, constants.IFBChecker}),
			},
		},
		{
			Method:  http.MethodPost,
			Path:    "/unlink/user-cif",
			Handler: handler.UnlinkUserCif,
			Middlewares: []func(next http.Handler) http.Handler{
				authMiddleware.AuthenticateToken,
				authMiddleware.AccessControl([]string{constants.Maker, constants.Maker, constants.Checker, constants.IFBChecker}),
			},
		},
	}

	glue.RegisterRoutes(router, routes)

}
