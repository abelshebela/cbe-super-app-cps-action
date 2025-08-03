package routing

import (
	"net/http"

	"github.com/CBE-Super-App/cbe-super-app-member-auth/internal/constants"
	"github.com/CBE-Super-App/cbe-super-app-member-auth/internal/glue"
	"github.com/CBE-Super-App/cbe-super-app-member-auth/internal/handlers/rest"
	"github.com/CBE-Super-App/cbe-super-app-member-auth/platform/logger"

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
