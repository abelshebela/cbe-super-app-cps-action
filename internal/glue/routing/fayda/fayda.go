package faydaaccount

import (
	"net/http"

	fayda_account "github.com/abelshebela/cbe-super-app-cps-action/internal/constants/interfaces/fayda"
	"github.com/abelshebela/cbe-super-app-cps-action/internal/glue"
	"github.com/abelshebela/cbe-super-app-cps-action/internal/handlers/middleware"

	"github.com/go-chi/chi/v5"
)

func Init(router chi.Router, handler fayda_account.FaydaAccount, authMiddleware middleware.AuthMiddleware) {

	routes := []glue.Route{
		{
			Method:  http.MethodPost,
			Path:    "/fayda_account/disable/{user_code}",
			Handler: handler.InitiateDisableFaydaAccount,
			Middlewares: []func(next http.Handler) http.Handler{
				authMiddleware.AuthenticateToken,
			},
		},
		{
			Method:  http.MethodPost,
			Path:    "/fayda_account/enable/{user_code}",
			Handler: handler.InitiateEnableFaydaAccount,
			Middlewares: []func(next http.Handler) http.Handler{
				authMiddleware.AuthenticateToken,
			},
		},
	}

	glue.RegisterRoutes(router, routes)

}
