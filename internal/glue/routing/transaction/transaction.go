package transaction

import (
	"github.com/abelshebela/cbe-super-app-cps-action/internal/constants/interfaces/transaction"
	"github.com/abelshebela/cbe-super-app-cps-action/internal/glue"
	"github.com/abelshebela/cbe-super-app-cps-action/internal/handlers/middleware"
	"net/http"

	"github.com/go-chi/chi/v5"
)

func Init(router chi.Router, transactionHandler transaction.TransactionInterface, authMiddleware middleware.AuthMiddleware) {
	routes := []glue.Route{
		{
			Method:  http.MethodGet,
			Path:    "/transaction/{customer_no}",
			Handler: transactionHandler.FetchTransactionLimitByUserCode,
			Middlewares: []func(http.Handler) http.Handler{
				authMiddleware.AuthenticateToken,
			},
		},
		{
			Method:  http.MethodGet,
			Path:    "/transaction/{id}",
			Handler: transactionHandler.FetchTransactionByID,
			Middlewares: []func(http.Handler) http.Handler{
				authMiddleware.AuthenticateToken,
			},
		},
		{
			Method:  http.MethodGet,
			Path:    "/transaction/search/{identifier}",
			Handler: transactionHandler.FindTransactionByCifOrAccountNumberOrFT,
			Middlewares: []func(http.Handler) http.Handler{
				authMiddleware.AuthenticateToken,
			},
		},
		{
			Method:  http.MethodGet,
			Path:    "/transaction",
			Handler: transactionHandler.FetchAllTransactionLimits,
			Middlewares: []func(http.Handler) http.Handler{
				authMiddleware.AuthenticateToken,
			},
		},
	}
	glue.RegisterRoutes(router, routes)
}
