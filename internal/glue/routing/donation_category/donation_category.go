package donation_category

import (
	"net/http"

	"github.com/abelshebela/cbe-super-app-cps-action/internal/constants/interfaces/donation_category"
	"github.com/abelshebela/cbe-super-app-cps-action/internal/glue"
	"github.com/abelshebela/cbe-super-app-cps-action/internal/handlers/middleware"

	"github.com/go-chi/chi/v5"
)

func Init(router chi.Router, handler donation_category.DonationCategoryAdapter, authMiddleware middleware.AuthMiddleware) {

	routes := []glue.Route{
		{
			Method:  http.MethodPost,
			Path:    "/donation_category",
			Handler: handler.CreateDonationCategory,
			Middlewares: []func(next http.Handler) http.Handler{
				authMiddleware.AuthenticateToken,
			},
		},
		{
			Method:  http.MethodPatch,
			Path:    "/donation_category/{id}",
			Handler: handler.UpdateDonationCategory,
			Middlewares: []func(next http.Handler) http.Handler{
				authMiddleware.AuthenticateToken,
			},
		},
		{
			Method:  http.MethodGet,
			Path:    "/donation_category",
			Handler: handler.FetchDonationCategory,
			Middlewares: []func(next http.Handler) http.Handler{
				authMiddleware.AuthenticateToken,
			},
		},
		{
			Method:  http.MethodGet,
			Path:    "/donation_category/{id}",
			Handler: handler.FetchDonationCategoryByID,
			Middlewares: []func(next http.Handler) http.Handler{
				authMiddleware.AuthenticateToken,
			},
		},
		{
			Method:  http.MethodPatch,
			Path:    "/donation_category/enable/{id}",
			Handler: handler.EnableDonationCategory,
			Middlewares: []func(next http.Handler) http.Handler{
				authMiddleware.AuthenticateToken,
			},
		},
		{
			Method:  http.MethodPatch,
			Path:    "/donation_category/disable/{id}",
			Handler: handler.DisableDonationCategory,
			Middlewares: []func(next http.Handler) http.Handler{
				authMiddleware.AuthenticateToken,
			},
		},//DeleteDonationCategory
		{
			Method:  http.MethodDelete,
			Path:    "/donation_category/{id}",
			Handler: handler.DeleteDonationCategory,
			Middlewares: []func(next http.Handler) http.Handler{
				authMiddleware.AuthenticateToken,
			},
		},

	}

	glue.RegisterRoutes(router, routes)
}
