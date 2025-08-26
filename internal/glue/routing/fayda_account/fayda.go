package faydaaccount

import (
	"net/http"

	"cbe-super-app-cps-action/internal/constants"
	fayda_account "cbe-super-app-cps-action/internal/constants/interfaces/fayda"
	"cbe-super-app-cps-action/internal/glue"
	"cbe-super-app-cps-action/internal/handlers/middleware"

	"github.com/go-chi/chi/v5"
)

func Init(router chi.Router, handler fayda_account.FaydaAccount, authMiddleware middleware.AuthMiddleware) {

	routes := []glue.Route{
		{
			Method:  http.MethodPost,
			Path:    "/disable/{user_code}",
			Handler: handler.InitiateDisableFaydaAccount,
			Middlewares: []func(next http.Handler) http.Handler{
				authMiddleware.AuthenticateToken,
				authMiddleware.AccessControl([]string{constants.Maker, constants.IFBMaker}),
			},
		},
		{
			Method:  http.MethodPost,
			Path:    "/enable/{user_code}",
			Handler: handler.InitiateEnableFaydaAccount,
			Middlewares: []func(next http.Handler) http.Handler{
				authMiddleware.AuthenticateToken,
				authMiddleware.AccessControl([]string{constants.Maker, constants.IFBMaker}),
			},
		},
	}

	glue.RegisterRoutes(router, routes)

}
