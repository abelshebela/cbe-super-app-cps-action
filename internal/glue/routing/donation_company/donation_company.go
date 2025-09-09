package donation_company

import (
	"net/http"

	role "cbe-super-app-cps-action/internal/constants"
	"cbe-super-app-cps-action/internal/constants/interfaces/donation_company"
	"cbe-super-app-cps-action/internal/glue"
	"cbe-super-app-cps-action/internal/handlers/middleware"

	"github.com/go-chi/chi/v5"
)

func Init(router chi.Router, handler donation_company.DonationCompanyAdapter, authMiddleware middleware.AuthMiddleware) {

	routes := []glue.Route{
		{
			Method:  http.MethodPost,
			Path:    "/donation_company",
			Handler: handler.CreateDonationCompany,
			Middlewares: []func(next http.Handler) http.Handler{
				authMiddleware.AuthenticateToken,
				authMiddleware.AccessControl([]string{role.Maker, role.IFBMaker}),
			},
		},
		{
			Method:  http.MethodPatch,
			Path:    "/donation_company/{id}",
			Handler: handler.UpdateDonationCompany,
			Middlewares: []func(next http.Handler) http.Handler{
				authMiddleware.AuthenticateToken,
				authMiddleware.AccessControl([]string{role.Maker, role.IFBMaker}),
			},
		},
		{
			Method:  http.MethodGet,
			Path:    "/donation_company",
			Handler: handler.FetchDonationCompany,
			Middlewares: []func(next http.Handler) http.Handler{
				authMiddleware.AuthenticateToken,
				authMiddleware.AccessControl([]string{role.Maker, role.IFBMaker, role.Checker}),
			},
		},
		{
			Method:  http.MethodGet,
			Path:    "/donation_company/{id}",
			Handler: handler.FetchDonationCompanyByID,
			Middlewares: []func(next http.Handler) http.Handler{
				authMiddleware.AuthenticateToken,
				authMiddleware.AccessControl([]string{role.Maker, role.IFBMaker, role.Checker}),
			},
		},
	}

	glue.RegisterRoutes(router, routes)
}
