package hqRoute

import (
	"net/http"

	"cbe-super-app-cps-action/internal/constants"
	"cbe-super-app-cps-action/internal/constants/interfaces/hq"
	"cbe-super-app-cps-action/internal/glue"
	"cbe-super-app-cps-action/internal/handlers/middleware"

	"github.com/go-chi/chi/v5"
)

func Init(router chi.Router, hqHandler hq.HQAdapter, authMiddleware middleware.AuthMiddleware) {
	routes := []glue.Route{
		{
			Method:  http.MethodGet,
			Path:    "/hq/block_time",
			Handler: hqHandler.GetBlockTime,
			Middlewares: []func(next http.Handler) http.Handler{
				authMiddleware.AuthenticateToken,
				authMiddleware.AccessControl([]string{constants.Maker, constants.IFBMaker, constants.Checker, constants.IFBChecker}),
			},
		},
		{
			Method:  http.MethodGet,
			Path:    "/hq/archive_time",
			Handler: hqHandler.GetArchiveTime,
			Middlewares: []func(next http.Handler) http.Handler{
				authMiddleware.AuthenticateToken,
				// authMiddleware.AccessControl([]string{constants.Maker, constants.IFBMaker, constants.Checker, constants.IFBChecker}),
			},
		},
		{
			Method:  http.MethodGet,
			Path:    "/hq/password_expiry",
			Handler: hqHandler.GetPasswordExpiry,
			Middlewares: []func(next http.Handler) http.Handler{
				authMiddleware.AuthenticateToken,
				// authMiddleware.AccessControl([]string{constants.Maker, constants.IFBMaker, constants.Checker, constants.IFBChecker}),
			},
		},
		{
			Method:  http.MethodPost,
			Path:    "/hq/block_time",
			Handler: hqHandler.UpdateBlockTimeRequest,
			Middlewares: []func(next http.Handler) http.Handler{
				authMiddleware.AuthenticateToken,
				// authMiddleware.AccessControl([]string{constants.Maker, constants.IFBMaker}),
			},
		},
		{
			Method:  http.MethodPost,
			Path:    "/hq/archive_time",
			Handler: hqHandler.UpdateArchiveTimeRequest,
			Middlewares: []func(next http.Handler) http.Handler{
				authMiddleware.AuthenticateToken,
				// authMiddleware.AccessControl([]string{constants.Maker, constants.IFBMaker}),
			},
		},
		{
			Method:  http.MethodPost,
			Path:    "/hq/password_expiry",
			Handler: hqHandler.UpdatePasswordExpiryRequest,
			Middlewares: []func(next http.Handler) http.Handler{
				authMiddleware.AuthenticateToken,
				// authMiddleware.AccessControl([]string{constants.Maker, constants.IFBMaker}),
			},
		},
	}

	glue.RegisterRoutes(router, routes)
}
