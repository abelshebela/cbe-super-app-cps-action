package customersegmentation

import (
	segmentations "cbe-super-app-cps-action/internal/constants/interfaces/customer_segmentation"
	"cbe-super-app-cps-action/internal/glue"
	"cbe-super-app-cps-action/internal/handlers/middleware"
	"net/http"

	"github.com/go-chi/chi/v5"
)

func Init(router chi.Router, handler segmentations.CustomerSegmentation, authMiddleware middleware.AuthMiddleware) {
	routes := []glue.Route{
		{
			Method:  http.MethodPost,
			Path:    "/customer-segmentations",
			Handler: handler.CreateCustomerSegmentation,
			Middlewares: []func(next http.Handler) http.Handler{
				authMiddleware.AuthenticateToken,
				// authMiddleware.AccessControl([]string{constants.Maker, constants.IFBMaker}),
			},
		},
		{
			Method:  http.MethodPatch,
			Path:    "/customer-segmentations/{id}/update",
			Handler: handler.UpdateCustomerSegmentation,
			Middlewares: []func(next http.Handler) http.Handler{
				authMiddleware.AuthenticateToken,
				// authMiddleware.AccessControl([]string{constants.Maker, constants.IFBMaker}),
			},
		},
		{
			Method:  http.MethodGet,
			Path:    "/customer-segmentations",
			Handler: handler.GetAllCustomerSegmentations,
			Middlewares: []func(next http.Handler) http.Handler{
				authMiddleware.AuthenticateToken,
				// authMiddleware.AccessControl([]string{constants.Maker, constants.IFBMaker, constants.Checker, constants.IFBChecker}),
			},
		},
		{
			Method:  http.MethodGet,
			Path:    "/customer-segmentations/{id}",
			Handler: handler.GetCustomerSegmentation,
			Middlewares: []func(next http.Handler) http.Handler{
				authMiddleware.AuthenticateToken,
				// authMiddleware.AccessControl([]string{constants.Maker, constants.IFBMaker, constants.Checker, constants.IFBChecker}),
			},
		},
		{
			Method:  http.MethodDelete,
			Path:    "/customer-segmentations/{id}",
			Handler: handler.DeleteCustomerSegmentation,
			Middlewares: []func(next http.Handler) http.Handler{
				authMiddleware.AuthenticateToken,
				// authMiddleware.AccessControl([]string{constants.Maker, constants.IFBMaker}),
			},
		},
		{
			Method:  http.MethodPatch,
			Path:    "/customer-segmentations/enable/{id}",
			Handler: handler.Enable,
			Middlewares: []func(next http.Handler) http.Handler{
				authMiddleware.AuthenticateToken,
			},
		},
		{
			Method:  http.MethodPatch,
			Path:    "/customer-segmentations/disable/{id}",
			Handler: handler.Disable,
			Middlewares: []func(next http.Handler) http.Handler{
				authMiddleware.AuthenticateToken,
			},
		},
	}

	glue.RegisterRoutes(router, routes)
}
