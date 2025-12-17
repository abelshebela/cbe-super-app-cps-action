package transaction

import "net/http"

type TransactionInterface interface {
	FetchTransactionByID(w http.ResponseWriter, r *http.Request)
	FindTransactionByCifOrAccountNumberOrFT(w http.ResponseWriter, r *http.Request)
	FetchAllTransactions(w http.ResponseWriter, r *http.Request)
}
