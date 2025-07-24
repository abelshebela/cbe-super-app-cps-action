package routing

import (
	"cbe-super-app-budget/internal/constants"
	"cbe-super-app-budget/internal/glue"
	"cbe-super-app-budget/internal/handlers/rest"
	"cbe-super-app-budget/platform/logger"
	"net/http"

	"github.com/go-chi/chi/v5"
)

func InitSpending(
	router chi.Router,
	spending rest.SpendingHandler,
	logger logger.Logger,
) {
	r := chi.NewRouter()

	spendingRoute := []glue.Router{
		{
			Method:  http.MethodGet,
			Path:    "/",
			Handler: spending.GetAll,
			Realm: []constants.Realm{
				constants.User,
			},
		},
		{
			Method:  http.MethodGet,
			Path:    "/{spending_id}",
			Handler: spending.GetOne,
			Realm: []constants.Realm{
				constants.User,
			},
		},
		{
			Method:  http.MethodPost,
			Path:    "/add",
			Handler: spending.InsertOne,
			Realm: []constants.Realm{
				constants.User,
			},
		},
		{
			Method:  http.MethodPatch,
			Path:    "/update/{spending_id}",
			Handler: spending.UpdateOne,
			Realm: []constants.Realm{
				constants.User,
			},
		},
	}

	glue.RegisterRoute(r, spendingRoute, logger)
	router.Mount("/spending", r)
}
