package donation

import (
	"net/http"

	route "github.com/CBE-Super-App/cbe-super-app-cps-action/internal/adapter/inbound/http"
	"github.com/CBE-Super-App/cbe-super-app-cps-action/internal/application/middleware"
	cps_const "github.com/CBE-Super-App/cbe-super-app-cps-action/internal/domain/cps_actions/constant"
	donation_inbound "github.com/CBE-Super-App/cbe-super-app-cps-action/internal/port/inbound/donation"
	role "github.com/CBE-Super-App/cbe-super-app-cps-action/utils"
	"github.com/go-chi/chi/v5"
)

func InitDonationHandlers(router chi.Router, handler donation_inbound.DonationHandler, authMiddleware middleware.AuthMiddleware, cpsGuard *middleware.CPSActionMiddlewareFactory) {
	routes := []route.Route{
		{
			Method:  http.MethodPost,
			Path:    "/donation_category",
			Handler: handler.CreateDonationCategory,
			Middlewares: []func(next http.Handler) http.Handler{
				authMiddleware.AuthenticateToken,
				authMiddleware.AccessControl([]string{role.Maker, role.IFBMaker}),
				cpsGuard.RequireNoPendingCPSActionGuard(string(cps_const.RequestCreateDonationCategory)),
			},
		},
		{
			Method:  http.MethodPatch,
			Path:    "/donation_category/{id}",
			Handler: handler.UpdateDonationCategory,
			Middlewares: []func(next http.Handler) http.Handler{
				authMiddleware.AuthenticateToken,
				authMiddleware.AccessControl([]string{role.Maker, role.IFBMaker}),
				cpsGuard.RequireNoPendingCPSActionGuard(string(cps_const.RequestUpdateDonationCategory)),
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
			Method:  http.MethodPost,
			Path:    "/donation_company",
			Handler: handler.CreateDonationCompany,
			Middlewares: []func(next http.Handler) http.Handler{
				authMiddleware.AuthenticateToken,
				authMiddleware.AccessControl([]string{role.Maker, role.IFBMaker}),
				cpsGuard.RequireNoPendingCPSActionGuard(string(cps_const.RequestCreateDonationCompany)),
			},
		},
		{
			Method:  http.MethodPatch,
			Path:    "/donation_company/{id}",
			Handler: handler.UpdateDonationCompany,
			Middlewares: []func(next http.Handler) http.Handler{
				authMiddleware.AuthenticateToken,
				authMiddleware.AccessControl([]string{role.Maker, role.IFBMaker}),
				cpsGuard.RequireNoPendingCPSActionGuard(string(cps_const.RequestUpdateDonationCompany)),
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
		{
			Method:  http.MethodPost,
			Path:    "/donation",
			Handler: handler.CreateDonation,
			Middlewares: []func(next http.Handler) http.Handler{
				authMiddleware.AuthenticateToken,
				authMiddleware.AccessControl([]string{role.Maker, role.IFBMaker}),
				cpsGuard.RequireNoPendingCPSActionGuard(string(cps_const.RequestCreateDonation)),
			},
		},
		{
			Method:  http.MethodPatch,
			Path:    "/donation/{id}",
			Handler: handler.UpdateDonation,
			Middlewares: []func(next http.Handler) http.Handler{
				authMiddleware.AuthenticateToken,
				authMiddleware.AccessControl([]string{role.Maker, role.IFBMaker}),
				cpsGuard.RequireNoPendingCPSActionGuard(string(cps_const.RequestUpdateDonation)),
			},
		},
		{
			Method:  http.MethodGet,
			Path:    "/donation",
			Handler: handler.FetchDonation,
			Middlewares: []func(next http.Handler) http.Handler{
				authMiddleware.AuthenticateToken,
				authMiddleware.AccessControl([]string{role.Maker, role.IFBMaker, role.Checker}),
			},
		},
		{
			Method:  http.MethodGet,
			Path:    "/donation/{id}",
			Handler: handler.FetchDonationByID,
			Middlewares: []func(next http.Handler) http.Handler{
				authMiddleware.AuthenticateToken,
				authMiddleware.AccessControl([]string{role.Maker, role.IFBMaker, role.Checker}),
			},
		},
	}

	route.RegisterRoutes(router, routes)
}
