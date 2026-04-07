package portalcard

import (
	"net/http"

	"cbe-super-app-cps-action/internal/constants/interfaces/portal_card"
	"cbe-super-app-cps-action/internal/glue"
	"cbe-super-app-cps-action/internal/handlers/middleware"

	"github.com/go-chi/chi/v5"
)

func Init(router chi.Router, portalHander portal_card.PortalCardAdapter, authMiddleware middleware.AuthMiddleware) {
	routes := []glue.Route{
		{
			Method:  http.MethodGet,
			Path:    "/portal_cards",
			Handler: portalHander.GetAllPortalCard,
			Middlewares: []func(next http.Handler) http.Handler{
				authMiddleware.AuthenticateToken,
			},
		},
	}
	glue.RegisterRoutes(router, routes)
}
