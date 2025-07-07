package budget_handler

import (
	"net/http"

	route "github.com/CBE-Super-App/cbe-super-app-cps-action/internal/adapter/inbound/http"
	"github.com/CBE-Super-App/cbe-super-app-cps-action/internal/application/middleware"
	inbound "github.com/CBE-Super-App/cbe-super-app-cps-action/internal/port/inbound/budget"
	role "github.com/CBE-Super-App/cbe-super-app-cps-action/utils"

	"github.com/go-chi/chi/v5"
)

func InitBudgetRoutes(router chi.Router, budgetHandler inbound.BudgetPortHandler, authMiddleware middleware.AuthMiddleware) {
	router.Route("/budget", func(r chi.Router) {
		routes := []route.Route{
			{
				Method:  http.MethodPost,
				Path:    "/icon/create",
				Handler: budgetHandler.CreateBudgetIcon,
				Middlewares: []func(next http.Handler) http.Handler{
					authMiddleware.AuthenticateToken,
					authMiddleware.AccessControl([]string{role.Maker}),
				},
			},
			{
				Method:  http.MethodGet,
				Path:    "/icon/fetch/all",
				Handler: budgetHandler.BudgetFetchIcons,
				Middlewares: []func(next http.Handler) http.Handler{
					authMiddleware.AuthenticateToken,
					authMiddleware.AccessControl([]string{role.Maker, role.Checker}),
				},
			},
			{
				Method:  http.MethodPut,
				Path:    "/icon/update/{id}",
				Handler: budgetHandler.BudgetUpdateIcon,
				Middlewares: []func(next http.Handler) http.Handler{
					authMiddleware.AuthenticateToken,
					authMiddleware.AccessControl([]string{role.Maker}),
				},
			},
			{
				Method:  http.MethodPost,
				Path:    "/color/create",
				Handler: budgetHandler.BudgetCreateColor,
				Middlewares: []func(next http.Handler) http.Handler{
					authMiddleware.AuthenticateToken,
					authMiddleware.AccessControl([]string{role.Maker}),
				},
			},
			{
				Method:  http.MethodGet,
				Path:    "/color/fetch/all",
				Handler: budgetHandler.BudgetFetchColors,
				Middlewares: []func(next http.Handler) http.Handler{
					authMiddleware.AuthenticateToken,
					authMiddleware.AccessControl([]string{role.Maker, role.Checker}),
				},
			},

			{
				Method:  http.MethodPut,
				Path:    "/color/update/{id}",
				Handler: budgetHandler.BudgetUpdateColor,
				Middlewares: []func(next http.Handler) http.Handler{
					authMiddleware.AuthenticateToken,
					authMiddleware.AccessControl([]string{role.Maker}),
				},
			},
			{
				Method:  http.MethodPost,
				Path:    "/approve/action/{action_code}",
				Handler: budgetHandler.BudgetCheckerApproval,
				Middlewares: []func(next http.Handler) http.Handler{
					authMiddleware.AuthenticateToken,
					authMiddleware.AccessControl([]string{role.Checker, role.IFBChecker}),
				},
			},
		}

		route.RegisterRoutes(r, routes)
	})
}
