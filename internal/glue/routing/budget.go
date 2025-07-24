package routing

import (
	"cbe-super-app-budget/internal/constants"
	"cbe-super-app-budget/internal/glue"
	"cbe-super-app-budget/internal/handlers/rest"
	"cbe-super-app-budget/platform/logger"
	"net/http"

	"github.com/go-chi/chi/v5"
)

func InitBudget(
	router chi.Router,
	budget rest.BudgetHandler,
	logger logger.Logger,
) {
	r := chi.NewRouter()

	budgetRoute := []glue.Router{
		{
			Method:  http.MethodGet,
			Handler: budget.GetAll,
			Path:    "/",
			Realm: []constants.Realm{
				constants.User,
			},
		},
		{
			Method:  http.MethodGet,
			Handler: budget.GetOne,
			Path:    "/{budget_id}",
			Realm: []constants.Realm{
				constants.User,
			},
		},
		{
			Method:  http.MethodPost,
			Handler: budget.InsertOne,
			Path:    "/add",
			Realm: []constants.Realm{
				constants.User,
			},
		},
		{
			Method:  http.MethodPatch,
			Handler: budget.UpdateOne,
			Path:    "/update",
			Realm: []constants.Realm{
				constants.User,
			},
		},
		{
			Method:  http.MethodDelete,
			Handler: budget.DeleteOne,
			Path:    "/delete/{budget_id}",
			Realm: []constants.Realm{
				constants.User,
			},
		},
	}

	glue.RegisterRoute(r, budgetRoute, logger)
	router.Mount("/budget", r)
}
