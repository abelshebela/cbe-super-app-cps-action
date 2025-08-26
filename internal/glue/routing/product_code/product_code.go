package productcode

import (
	"cbe-super-app-cps-action/internal/constants"
	"cbe-super-app-cps-action/internal/constants/interfaces/productcode"
	"cbe-super-app-cps-action/internal/glue"

	"cbe-super-app-cps-action/internal/handlers/middleware"
	cpsaction "cbe-super-app-cps-action/internal/service/cps_action"
	"net/http"

	"github.com/go-chi/chi/v5"
)

// InitProductCodeHandlerRoutes sets up HTTP routes for product code-related endpoints
func Init(router chi.Router, handler productcode.ProductCodeAdapter, authMiddleware middleware.AuthMiddleware, cpsGuard *middleware.CPSActionMiddlewareFactory) {
	router.Route("/productcodes", func(r chi.Router) {
		routes := []glue.Route{
			{
				Method:  http.MethodPatch,
				Path:    "/{id}",
				Handler: handler.UpdateProductCode,
				Middlewares: []func(next http.Handler) http.Handler{
					authMiddleware.AuthenticateToken,
					authMiddleware.AccessControl([]string{constants.Maker, constants.IFBMaker}),
					cpsGuard.RequireNoPendingCPSActionGuard(string(cpsaction.RequestUpdateProductCode)),
				},
			},
			{
				Method:  http.MethodGet,
				Path:    "/{id}",
				Handler: handler.FetchProductCodeByID,
				Middlewares: []func(next http.Handler) http.Handler{
					authMiddleware.AuthenticateToken,
					authMiddleware.AccessControl([]string{constants.Checker, constants.IFBChecker, constants.Maker, constants.IFBMaker}),
				},
			},
			{
				Method:  http.MethodGet,
				Path:    "/",
				Handler: handler.FetchProductCodes,
				Middlewares: []func(next http.Handler) http.Handler{
					authMiddleware.AuthenticateToken,
					authMiddleware.AccessControl([]string{constants.Checker, constants.IFBChecker, constants.Maker, constants.IFBMaker}),
				},
			},
		}

		glue.RegisterRoutes(r, routes)
	})
}
