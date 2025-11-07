package kyc_verifier

import (
	"net/http"

	role "cbe-super-app-cps-action/internal/constants"
	"cbe-super-app-cps-action/internal/constants/interfaces/kyc_verifier"
	"cbe-super-app-cps-action/internal/glue"
	"cbe-super-app-cps-action/internal/handlers/middleware"

	"github.com/go-chi/chi/v5"
)

func Init(router chi.Router, handler kyc_verifier.KYCVerifierAdapter, authMiddleware middleware.AuthMiddleware) {
	routes := []glue.Route{
		{
			Method:  http.MethodGet,
			Path:    "/kyc_verifier",
			Handler: handler.GetKYCList,
			Middlewares: []func(next http.Handler) http.Handler{
				authMiddleware.AuthenticateToken,
				authMiddleware.AccessControl([]string{role.Maker, role.IFBMaker, role.Checker}),
			},
		},
		{
			Method:  http.MethodGet,
			Path:    "/kyc_verifier/{id}",
			Handler: handler.GetKYCByID,
			Middlewares: []func(next http.Handler) http.Handler{
				authMiddleware.AuthenticateToken,
				authMiddleware.AccessControl([]string{role.Maker, role.IFBMaker, role.Checker}),
			},
		},
		{
			Method:  http.MethodPatch,
			Path:    "/kyc_verifier/update/{id}",
			Handler: handler.UpdateKYC,
			Middlewares: []func(next http.Handler) http.Handler{
				authMiddleware.AuthenticateToken,
				authMiddleware.AccessControl([]string{role.Maker, role.IFBMaker}),
			},
		},
		{
			Method:  http.MethodPatch,
			Path:    "/kyc_verifier/approve/{id}",
			Handler: handler.ApproveKYC,
			Middlewares: []func(next http.Handler) http.Handler{
				authMiddleware.AuthenticateToken,
				authMiddleware.AccessControl([]string{role.Maker, role.IFBMaker}),
			},
		},
	}
	glue.RegisterRoutes(router, routes)
}
