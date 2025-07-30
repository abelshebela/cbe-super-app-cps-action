package productcode

import (
	"net/http"

	"github.com/go-chi/chi/v5"

	route "github.com/CBE-Super-App/cbe-super-app-cps-action/internal/adapter/inbound/http"
	"github.com/CBE-Super-App/cbe-super-app-cps-action/internal/application/middleware"
	cps_const "github.com/CBE-Super-App/cbe-super-app-cps-action/internal/domain/cps_actions/constant"
	productcode "github.com/CBE-Super-App/cbe-super-app-cps-action/internal/port/inbound/product_code"
	role "github.com/CBE-Super-App/cbe-super-app-cps-action/utils"
)

// InitProductCodeHandlerRoutes sets up HTTP routes for product code-related endpoints
func InitProductCodeHandlerRoutes(router chi.Router, handler productcode.ProductCodeHandler, authMiddleware middleware.AuthMiddleware, cpsGuard *middleware.CPSActionMiddlewareFactory) {
	router.Route("/productcodes", func(r chi.Router) {
		routes := []route.Route{
			{
				Method:  http.MethodPatch,
				Path:    "/{id}",
				Handler: handler.UpdateProductCode,
				Middlewares: []func(next http.Handler) http.Handler{
					authMiddleware.AuthenticateToken,
					authMiddleware.AccessControl([]string{role.Maker, role.IFBMaker}),
					cpsGuard.RequireNoPendingCPSActionGuard(string(cps_const.RequestUpdateProductCode)),
				},
			},
			{
				Method:  http.MethodGet,
				Path:    "/{id}",
				Handler: handler.FetchProductCodeByID,
				Middlewares: []func(next http.Handler) http.Handler{
					authMiddleware.AuthenticateToken,
					authMiddleware.AccessControl([]string{role.Checker, role.IFBChecker, role.Maker, role.IFBMaker}),
				},
			},
			{
				Method:  http.MethodGet,
				Path:    "/",
				Handler: handler.FetchProductCodes,
				Middlewares: []func(next http.Handler) http.Handler{
					authMiddleware.AuthenticateToken,
					authMiddleware.AccessControl([]string{role.Checker, role.IFBChecker, role.Maker, role.IFBMaker}),
				},
			},
		}

		route.RegisterRoutes(r, routes)
	})
}
