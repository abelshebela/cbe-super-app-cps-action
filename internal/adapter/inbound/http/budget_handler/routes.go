package budget_handler

import (
	"net/http"

	route "gitlab.com/bersufekadgetachew/cbe-super-app-cps-action/internal/adapter/inbound/http"
	"gitlab.com/bersufekadgetachew/cbe-super-app-cps-action/internal/application/middleware"
	inbound "gitlab.com/bersufekadgetachew/cbe-super-app-cps-action/internal/port/inbound/budget"

	"github.com/go-chi/chi/v5"
)

func InitBudgetRoutes(router chi.Router, budgetHandler inbound.BudgetPortHandler, authMiddleware middleware.AuthMiddleware) {
	router.Route("/api/v1/cbesuperapp/cps_action/budget", func(r chi.Router) {
		routes := []route.Route{
			{
				Method:  http.MethodPost,
				Path:    "/icon/create",
				Handler: budgetHandler.CreateBudgetIcon,
				Middlewares: []func(next http.Handler) http.Handler{
					authMiddleware.AuthenticateToken,
					authMiddleware.AccessControl([]string{"maker"}),
				},
			},
			{
				Method:  http.MethodGet,
				Path:    "/icon/fetch/all",
				Handler: budgetHandler.BudgetFetchIcons,
				Middlewares: []func(next http.Handler) http.Handler{
					authMiddleware.AuthenticateToken,
					authMiddleware.AccessControl([]string{"maker", "checker"}),
				},
			},
			{
				Method:  http.MethodPut,
				Path:    "/icon/update/{id}",
				Handler: budgetHandler.BudgetUpdateIcon,
				Middlewares: []func(next http.Handler) http.Handler{
					authMiddleware.AuthenticateToken,
					authMiddleware.AccessControl([]string{"maker"}),
				},
			},
			{
				Method:  http.MethodPost,
				Path:    "/color/create",
				Handler: budgetHandler.BudgetCreateColor,
				Middlewares: []func(next http.Handler) http.Handler{
					authMiddleware.AuthenticateToken,
					authMiddleware.AccessControl([]string{"maker"}),
				},
			},
			{
				Method:  http.MethodGet,
				Path:    "/color/fetch/all",
				Handler: budgetHandler.BudgetFetchColors,
				Middlewares: []func(next http.Handler) http.Handler{
					authMiddleware.AuthenticateToken,
					authMiddleware.AccessControl([]string{"maker", "checker"}),
				},
			},

			{
				Method:  http.MethodPut,
				Path:    "/color/update/{id}",
				Handler: budgetHandler.BudgetUpdateColor,
				Middlewares: []func(next http.Handler) http.Handler{
					authMiddleware.AuthenticateToken,
					authMiddleware.AccessControl([]string{"maker"}),
				},
			},
			{
				Method:  http.MethodPost,
				Path:    "/approve/action/{action_code}",
				Handler: budgetHandler.BudgetCheckerApproval,
				Middlewares: []func(next http.Handler) http.Handler{
					authMiddleware.AuthenticateToken,
					authMiddleware.AccessControl([]string{"checker", "ifb-checker"}),
				},
			},
		}

		route.RegisterRoutes(r, routes)
	})
}
