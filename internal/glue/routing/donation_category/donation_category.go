package donation_category

import (
	"net/http"

	role "cbe-super-app-cps-action/internal/constants"
	"cbe-super-app-cps-action/internal/constants/interfaces/donation_category"
	"cbe-super-app-cps-action/internal/glue"
	"cbe-super-app-cps-action/internal/handlers/middleware"

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
				authMiddleware.AccessControl([]string{role.Maker, role.IFBMaker}),
			},
		},
		{
			Method:  http.MethodPatch,
			Path:    "/donation_category/{id}",
			Handler: handler.UpdateDonationCategory,
			Middlewares: []func(next http.Handler) http.Handler{
				authMiddleware.AuthenticateToken,
				authMiddleware.AccessControl([]string{role.Maker, role.IFBMaker}),
			},
		},
		{
			Method:  http.MethodGet,
			Path:    "/donation_category",
			Handler: handler.FetchDonationCategory,
			Middlewares: []func(next http.Handler) http.Handler{
				authMiddleware.AuthenticateToken,
				authMiddleware.AccessControl([]string{role.Maker, role.IFBMaker, role.Checker}),
			},
		},
		{
			Method:  http.MethodGet,
			Path:    "/donation_category/{id}",
			Handler: handler.FetchDonationCategoryByID,
			Middlewares: []func(next http.Handler) http.Handler{
				authMiddleware.AuthenticateToken,
				authMiddleware.AccessControl([]string{role.Maker, role.IFBMaker, role.Checker}),
			},
		},
		{
			Method:  http.MethodPatch,
			Path:    "/donation_category/enable/{id}",
			Handler: handler.EnableDonationCategory,
			Middlewares: []func(next http.Handler) http.Handler{
				authMiddleware.AuthenticateToken,
				authMiddleware.AccessControl([]string{role.Maker, role.IFBMaker}),
			},
		},
		{
			Method:  http.MethodPatch,
			Path:    "/donation_category/disable/{id}",
			Handler: handler.DisableDonationCategory,
			Middlewares: []func(next http.Handler) http.Handler{
				authMiddleware.AuthenticateToken,
				authMiddleware.AccessControl([]string{role.Maker, role.IFBMaker}),
			},
		},
	}

	glue.RegisterRoutes(router, routes)
}
