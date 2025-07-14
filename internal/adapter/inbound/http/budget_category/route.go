package budget_category

import (
	"net/http"

	route "github.com/CBE-Super-App/cbe-super-app-cps-action/internal/adapter/inbound/http"
	"github.com/CBE-Super-App/cbe-super-app-cps-action/internal/application/middleware"
	inbound "github.com/CBE-Super-App/cbe-super-app-cps-action/internal/port/inbound/budget_category"
	"github.com/go-chi/chi/v5"
)

func InitBudgetCategoryRoute(router chi.Router, budgetCategoryHandler inbound.BudgetCategoryInbound, authMiddleware middleware.AuthMiddleware) {
	router.Route("/budget_category", func(r chi.Router) {
		routes := []route.Route{
			{
				Method:  http.MethodPost,
				Path:    "/create",
				Handler: budgetCategoryHandler.CreateBudgetCategory,
				Middlewares: []func(next http.Handler) http.Handler{
					authMiddleware.AuthenticateToken,
					authMiddleware.AccessControl([]string{"MAKER", "IFB-MAKER"}),
				},
			},
			{
				Method:  http.MethodPut,
				Path:    "/update/{budget_category_id}",
				Handler: budgetCategoryHandler.UpdateBudgetCategory,
				Middlewares: []func(next http.Handler) http.Handler{
					authMiddleware.AuthenticateToken,
					authMiddleware.AccessControl([]string{"MAKER", "IFB-MAKER"}),
				},
			},
			{
				Method:  http.MethodDelete,
				Path:    "/delete/{budget_category_id}",
				Handler: budgetCategoryHandler.DeleteBudgetCategory,
				Middlewares: []func(next http.Handler) http.Handler{
					authMiddleware.AuthenticateToken,
					authMiddleware.AccessControl([]string{"MAKER", "IFB-MAKER"}),
				},
			},
			{
				Method:  http.MethodGet,
				Path:    "/{budget_category_id}",
				Handler: budgetCategoryHandler.GetBudgetCategory,
				Middlewares: []func(next http.Handler) http.Handler{
					authMiddleware.AuthenticateToken,
					// authMiddleware.AccessControl([]string{"MAKER", "IFB-MAKER"}),
				},
			},
			{
				Method:  http.MethodGet,
				Path:    "/budget_categories",
				Handler: budgetCategoryHandler.GetAllBudgetCategory,
				Middlewares: []func(next http.Handler) http.Handler{
					authMiddleware.AuthenticateToken,
					// authMiddleware.AccessControl([]string{"MAKER"}),
				},
			},
			{
				Method:  http.MethodPost,
				Path:    "/approve",
				Handler: budgetCategoryHandler.ApproveBudgetCategoryActionHTTP,
				Middlewares: []func(next http.Handler) http.Handler{
					authMiddleware.AuthenticateToken,
					authMiddleware.AccessControl([]string{"CHECKER", "IFB-CHECKER"}),
				},
			},
		}

		route.RegisterRoutes(r, routes)
	})
}
