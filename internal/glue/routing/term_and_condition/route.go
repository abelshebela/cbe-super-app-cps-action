package term_and_condition_routing

import (
	tac_interface "cbe-super-app-cps-action/internal/constants/interfaces/term_and_condition"
	"cbe-super-app-cps-action/internal/glue"
	"cbe-super-app-cps-action/internal/handlers/middleware"
	"net/http"

	"github.com/go-chi/chi/v5"
)

const (
	pathCollection = "/term_and_condition"
	pathResource   = "/term_and_condition/{id}"
)

func Init(router chi.Router, h tac_interface.TermAndConditionHandler, authMiddleware middleware.AuthMiddleware) {
	auth := []func(next http.Handler) http.Handler{authMiddleware.AuthenticateToken}
	routes := []glue.Route{
		{Method: http.MethodGet, Path: pathCollection, Handler: h.GetAll, Middlewares: auth},
		{Method: http.MethodGet, Path: pathResource, Handler: h.GetByID, Middlewares: auth},
		{Method: http.MethodPost, Path: pathCollection, Handler: h.Upload, Middlewares: auth},
		{Method: http.MethodPut, Path: pathResource, Handler: h.Update, Middlewares: auth},
		{Method: http.MethodDelete, Path: pathResource, Handler: h.Delete, Middlewares: auth},
	}

	glue.RegisterRoutes(router, routes)
}
