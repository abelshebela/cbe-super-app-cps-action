package customerkycself

import (
	customer_kyc "cbe-super-app-cps-action/internal/constants/interfaces/customer_kyc"
	"cbe-super-app-cps-action/internal/glue"
	"cbe-super-app-cps-action/internal/handlers/middleware"
	"net/http"

	"github.com/go-chi/chi/v5"
)

func Init(router chi.Router, handler customer_kyc.SelfActivationKyc, authMiddleware middleware.AuthMiddleware) {
	routes := []glue.Route{
		{
			Method:  http.MethodGet,
			Path:    "/customers/self-activation/kyc",
			Handler: handler.GetAllSelfActivateKYCRequests,
			Middlewares: []func(next http.Handler) http.Handler{
				authMiddleware.AuthenticateToken,
			},
		},
		{
			Method:  http.MethodGet,
			Path:    "/customers/self-activation/kyc/{id}",
			Handler: handler.GetSelfActivateKYCRequest,
			Middlewares: []func(next http.Handler) http.Handler{
				authMiddleware.AuthenticateToken,
			},
		},
		{
			Method:  http.MethodPatch,
			Path:    "/customers/self-activation/kyc/{id}/approve",
			Handler: handler.ApproveSelfActivateKycRequest,
			Middlewares: []func(next http.Handler) http.Handler{
				authMiddleware.AuthenticateToken,
			},
		},
		{
			Method:  http.MethodPatch,
			Path:    "/customers/self-activation/kyc/{id}/reject",
			Handler: handler.RejectSelfActivateKycRequest,
			Middlewares: []func(next http.Handler) http.Handler{
				authMiddleware.AuthenticateToken,
			},
		},
		{
			Method:  http.MethodPost,
			Path:    "/customers/self-activation/kyc/start-review/{kycID}",
			Handler: handler.StartSelfActivateKycReview,
			Middlewares: []func(next http.Handler) http.Handler{
				authMiddleware.AuthenticateToken,
			},
		},
		{
			Method:  http.MethodPatch,
			Path:    "/customers/self-activation/kyc/pick-review/{kycID}",
			Handler: handler.PickSelfActivateKycReview,
			Middlewares: []func(next http.Handler) http.Handler{
				authMiddleware.AuthenticateToken,
			},
		},
	}

	glue.RegisterRoutes(router, routes)
}
