package routing

import (
	"cbe-super-app-budget/internal/constants"
	"cbe-super-app-budget/internal/glue"
	"cbe-super-app-budget/internal/handlers/rest"
	"cbe-super-app-budget/platform/logger"
	"net/http"

	"github.com/go-chi/chi/v5"
)

func InitBudgetCategory(
	router chi.Router,
	budgetCategory rest.BudgetCategoryHandler,
	logger logger.Logger,
) {

	r := chi.NewRouter()

	budgetCategoryRoute := []glue.Router{
		{
			Method:  http.MethodGet,
			Handler: budgetCategory.GetAll,
			Path:    "/",
			Realm: []constants.Realm{
				constants.User,
			},
		},
		{
			Method:  http.MethodGet,
			Handler: budgetCategory.GetOne,
			Path:    "/{category_id}",
			Realm: []constants.Realm{
				constants.User,
			},
		},
		{
			Method:  http.MethodPost,
			Handler: budgetCategory.InsertOne,
			Path:    "/add",
			Realm: []constants.Realm{
				constants.User,
			},
		},
		{
			Method:  http.MethodPatch,
			Handler: budgetCategory.UpdateOne,
			Path:    "/update/{category_id}",
			Realm: []constants.Realm{
				constants.User,
			},
		},
	}

	glue.RegisterRoute(r, budgetCategoryRoute, logger)
	router.Mount("/budget_category", r)
}
