package customersegmentation

import (
	"cbe-super-app-cps-action/internal/constants"
	segmentation "cbe-super-app-cps-action/internal/constants/interfaces/customer_segmentation"
	"cbe-super-app-cps-action/internal/glue"
	"cbe-super-app-cps-action/internal/handlers/middleware"
	"net/http"

	"github.com/go-chi/chi/v5"
)

func Init(router chi.Router, handler segmentation.CustomerSegmentation, authMiddleware middleware.AuthMiddleware) {
	routes := []glue.Route{
		{
			Method:  http.MethodPost,
			Path:    "/customer-segmentation",
			Handler: handler.CreateCustomerSegmentation,
			Middlewares: []func(next http.Handler) http.Handler{
				authMiddleware.AuthenticateToken,
				authMiddleware.AccessControl([]string{constants.Maker, constants.Checker}),
			},
		},
	}

	glue.RegisterRoutes(router, routes)
}
