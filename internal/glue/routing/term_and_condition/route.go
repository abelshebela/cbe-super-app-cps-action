package term_and_condition_routing

import (
	tac_interface "cbe-super-app-cps-action/internal/constants/interfaces/term_and_condition"
	"cbe-super-app-cps-action/internal/glue"
	"cbe-super-app-cps-action/internal/handlers/middleware"
	"net/http"

	"github.com/go-chi/chi/v5"
)

func Init(router chi.Router, h tac_interface.TermAndConditionHandler, authMiddleware middleware.AuthMiddleware) {
	routes := []glue.Route{
		{
			Method:      http.MethodGet,
			Path:        "/term_and_condition",
			Handler:     h.GetAll,
			Middlewares: []func(next http.Handler) http.Handler{authMiddleware.AuthenticateToken},
		},
		{
			Method:      http.MethodGet,
			Path:        "/term_and_condition/{id}",
			Handler:     h.GetByID,
			Middlewares: []func(next http.Handler) http.Handler{authMiddleware.AuthenticateToken},
		},
		{
			Method:      http.MethodPost,
			Path:        "/term_and_condition",
			Handler:     h.Upload,
			Middlewares: []func(next http.Handler) http.Handler{authMiddleware.AuthenticateToken},
		},
		{
			Method:      http.MethodDelete,
			Path:        "/term_and_condition/{id}",
			Handler:     h.Delete,
			Middlewares: []func(next http.Handler) http.Handler{authMiddleware.AuthenticateToken},
		},
	}

	glue.RegisterRoutes(router, routes)
}
