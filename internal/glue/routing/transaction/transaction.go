package transaction

import (
	"cbe-super-app-cps-action/internal/constants/interfaces/transaction"
	"cbe-super-app-cps-action/internal/glue"
	"cbe-super-app-cps-action/internal/handlers/middleware"
	"net/http"

	"github.com/go-chi/chi/v5"
)

func Init(router chi.Router, transactionHandler transaction.TransactionInterface, authMiddleware middleware.AuthMiddleware) {
	routes := []glue.Route{
		{
			Method:  http.MethodGet,
			Path:    "/transactions/{id}",
			Handler: transactionHandler.FetchTransactionByID,
			Middlewares: []func(http.Handler) http.Handler{
				authMiddleware.AuthenticateToken,
			},
		},
		{
			Method:  http.MethodGet,
			Path:    "/transactions/search/{identifier}",
			Handler: transactionHandler.FindTransactionByCifOrAccountNumberOrFT,
			Middlewares: []func(http.Handler) http.Handler{
				authMiddleware.AuthenticateToken,
			},
		},
		{
			Method:  http.MethodGet,
			Path:    "/transactions",
			Handler: transactionHandler.FetchAllTransactions,
			Middlewares: []func(http.Handler) http.Handler{
				authMiddleware.AuthenticateToken,
			},
		},
	}
	glue.RegisterRoutes(router, routes)
}
