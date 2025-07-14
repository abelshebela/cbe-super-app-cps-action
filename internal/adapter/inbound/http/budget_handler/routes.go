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
	router.Route("/budgets", func(r chi.Router) {
		routes := []route.Route{
			{
				Method:  http.MethodPost,
				Path:    "/icons",
				Handler: budgetHandler.CreateBudgetIcon,
				Middlewares: []func(next http.Handler) http.Handler{
					authMiddleware.AuthenticateToken,
					authMiddleware.AccessControl([]string{role.Maker}),
				},
			},
			{
				Method:  http.MethodGet,
				Path:    "/icons",
				Handler: budgetHandler.BudgetFetchIcons,
				Middlewares: []func(next http.Handler) http.Handler{
					authMiddleware.AuthenticateToken,
					authMiddleware.AccessControl([]string{role.Maker, role.Checker}),
				},
			},
			{
				Method:  http.MethodPut,
				Path:    "/icons/{id}",
				Handler: budgetHandler.BudgetUpdateIcon,
				Middlewares: []func(next http.Handler) http.Handler{
					authMiddleware.AuthenticateToken,
					authMiddleware.AccessControl([]string{role.Maker}),
				},
			},
			{
				Method:  http.MethodPost,
				Path:    "/colors",
				Handler: budgetHandler.BudgetCreateColor,
				Middlewares: []func(next http.Handler) http.Handler{
					authMiddleware.AuthenticateToken,
					authMiddleware.AccessControl([]string{role.Maker}),
				},
			},
			{
				Method:  http.MethodGet,
				Path:    "/colors",
				Handler: budgetHandler.BudgetFetchColors,
				Middlewares: []func(next http.Handler) http.Handler{
					authMiddleware.AuthenticateToken,
					authMiddleware.AccessControl([]string{role.Maker, role.Checker}),
				},
			},
			
			{
				Method:  http.MethodPut,
				Path:    "/colors/{id}",
				Handler: budgetHandler.BudgetUpdateColor,
				Middlewares: []func(next http.Handler) http.Handler{
					authMiddleware.AuthenticateToken,
					authMiddleware.AccessControl([]string{role.Maker}),
				},
			},
			{
				Method:  http.MethodPost,
				Path:    "/actions/approve/{action_code}",
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
