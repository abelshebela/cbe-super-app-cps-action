package productcode

import (
	"github.com/abelshebela/cbe-super-app-cps-action/internal/constants"
	"github.com/abelshebela/cbe-super-app-cps-action/internal/constants/interfaces/productcode"
	"github.com/abelshebela/cbe-super-app-cps-action/internal/glue"

	"github.com/abelshebela/cbe-super-app-cps-action/internal/handlers/middleware"
	"net/http"

	"github.com/go-chi/chi/v5"
)

func Init(router chi.Router, handler productcode.ProductCodeAdapter, authMiddleware middleware.AuthMiddleware) {
	router.Route("/productcodes", func(r chi.Router) {
		routes := []glue.Route{
			{
				Method:  http.MethodPatch,
				Path:    "/{id}",
				Handler: handler.UpdateProductCode,
				Middlewares: []func(next http.Handler) http.Handler{
					authMiddleware.AuthenticateToken,
					authMiddleware.AccessControl([]string{constants.Maker, constants.IFBMaker}),
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
