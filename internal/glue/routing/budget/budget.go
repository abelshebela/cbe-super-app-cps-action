package budget

import (
	"net/http"

	"cbe-super-app-cps-action/internal/constants"
	budget "cbe-super-app-cps-action/internal/constants/interfaces/budget"
	"cbe-super-app-cps-action/internal/glue"
	"cbe-super-app-cps-action/internal/handlers/middleware"

	"github.com/go-chi/chi/v5"
)

func Init(router chi.Router, handler budget.BudgetPortHandler, authMiddleware middleware.AuthMiddleware) {

	routes := []glue.Route{
		{
			Method:  http.MethodPost,
			Path:    "/budgets/icons",
			Handler: handler.CreateBudgetIcon,
			Middlewares: []func(next http.Handler) http.Handler{
				authMiddleware.AuthenticateToken,
				authMiddleware.AccessControl([]string{constants.Maker}),
			},
		},
		{
			Method:  http.MethodGet,
			Path:    "/budgets/icons",
			Handler: handler.BudgetFetchIcons,
			Middlewares: []func(next http.Handler) http.Handler{
				authMiddleware.AuthenticateToken,
				authMiddleware.AccessControl([]string{constants.Maker, constants.Checker}),
			},
		},
		{
			Method:  http.MethodPut,
			Path:    "/budgets/icons/{id}",
			Handler: handler.BudgetUpdateIcon,
			Middlewares: []func(next http.Handler) http.Handler{
				authMiddleware.AuthenticateToken,
				authMiddleware.AccessControl([]string{constants.Maker}),
			},
		},
		{
			Method:  http.MethodPost,
			Path:    "/budgets/colors",
			Handler: handler.BudgetCreateColor,
			Middlewares: []func(next http.Handler) http.Handler{
				authMiddleware.AuthenticateToken,
				authMiddleware.AccessControl([]string{constants.Maker}),
			},
		},
		{
			Method:  http.MethodGet,
			Path:    "/budgets/colors",
			Handler: handler.BudgetFetchColors,
			Middlewares: []func(next http.Handler) http.Handler{
				authMiddleware.AuthenticateToken,
				authMiddleware.AccessControl([]string{constants.Maker, constants.Checker}),
			},
		},

		{
			Method:  http.MethodPut,
			Path:    "/budgets/colors/{id}",
			Handler: handler.BudgetUpdateColor,
			Middlewares: []func(next http.Handler) http.Handler{
				authMiddleware.AuthenticateToken,
				authMiddleware.AccessControl([]string{constants.Maker}),
			},
		},
		{
			Method:  http.MethodPost,
			Path:    "/budgets/actions/approve/{action_code}",
			Handler: handler.BudgetCheckerApproval,
			Middlewares: []func(next http.Handler) http.Handler{
				authMiddleware.AuthenticateToken,
				authMiddleware.AccessControl([]string{constants.Checker, constants.IFBChecker}),
			},
		},
	}

	glue.RegisterRoutes(router, routes)

}
