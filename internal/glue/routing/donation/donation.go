package donation


import (
	"net/http"

	role "cbe-super-app-cps-action/internal/constants"
	"cbe-super-app-cps-action/internal/constants/interfaces/donation"
	"cbe-super-app-cps-action/internal/glue"
	"cbe-super-app-cps-action/internal/handlers/middleware"

	"github.com/go-chi/chi/v5"
)

func Init(router chi.Router, handler donation.DonationHandler, authMiddleware middleware.AuthMiddleware) {

	routes := []glue.Route{
		
		{
			Method:  http.MethodPost,
			Path:    "/donation",
			Handler: handler.CreateDonation,
			Middlewares: []func(next http.Handler) http.Handler{
				authMiddleware.AuthenticateToken,
				authMiddleware.AccessControl([]string{role.Maker, role.IFBMaker}),
				
			},
		},
		{
			Method:  http.MethodPatch,
			Path:    "/donation/{id}",
			Handler: handler.UpdateDonation,
			Middlewares: []func(next http.Handler) http.Handler{
				authMiddleware.AuthenticateToken,
				authMiddleware.AccessControl([]string{role.Maker, role.IFBMaker}),
				
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
		{
			Method:  http.MethodPatch,
			Path:    "/donation/image/{id}",
			Handler: handler.UpdateDonationImage,
			Middlewares: []func(next http.Handler) http.Handler{
				authMiddleware.AuthenticateToken,
				authMiddleware.AccessControl([]string{role.Maker, role.IFBMaker}),
				
			},
		},
		{
			Method:  http.MethodDelete,
			Path:    "/donation/image/{id}",
			Handler: handler.DeleteDonationImage,
			Middlewares: []func(next http.Handler) http.Handler{
				authMiddleware.AuthenticateToken,
				authMiddleware.AccessControl([]string{role.Maker, role.IFBMaker}),
				
			},
		},
		{
			Method:  http.MethodPost,
			Path:    "/donation/image/{id}",
			Handler: handler.AddDonationImage,
			Middlewares: []func(next http.Handler) http.Handler{
				authMiddleware.AuthenticateToken,
				authMiddleware.AccessControl([]string{role.Maker, role.IFBMaker}),
				
			},
		},
			{
			Method:  http.MethodPatch,
			Path:    "/donation/enable/{id}",
			Handler: handler.EnableDonation,
			Middlewares: []func(next http.Handler) http.Handler{
				authMiddleware.AuthenticateToken,
				authMiddleware.AccessControl([]string{role.Maker, role.IFBMaker}),
			},
		},
		{
			Method:  http.MethodPatch,
			Path:    "/donation/disable/{id}",
			Handler: handler.DisableDonation,
			Middlewares: []func(next http.Handler) http.Handler{
				authMiddleware.AuthenticateToken,
				authMiddleware.AccessControl([]string{role.Maker, role.IFBMaker}),
			},
		},
	}

	glue.RegisterRoutes(router, routes)
}
