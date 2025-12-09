package transaction_handler

import (
	"cbe-super-app-cps-action/internal/constants/interfaces/transaction"
	"cbe-super-app-cps-action/internal/constants/localization"
	"cbe-super-app-cps-action/internal/service"
	local_utils "cbe-super-app-cps-action/pkgs/utils"
	"net/http"

	"github.com/go-chi/chi/v5"
	"gitlab.com/bersufekadgetachew/cbe-super-app-shared/shared/utils"
)

type TransactionHandler struct {
	service service.TransactionService
	logger  utils.Logger
}

// FetchAllTransactions implements transaction.TransactionInterface.
func (t *TransactionHandler) FetchAllTransactions(w http.ResponseWriter, r *http.Request) {
	filterParams := local_utils.ExtractFilterParams(r)

	transactions, err := t.service.FetchAllTransactions(r.Context(), filterParams)
	if err != nil {
		t.logger.Errorf("FetchAllTransactions failed: %v", err)
		localization.SendErrorByCodeResponse(w, err.Error())
		return
	}
	localization.SendSuccessResponse(w, localization.SuccessTransactionRetrieved, transactions)
}

// FetchTransactionByID implements transaction.TransactionInterface.
func (t *TransactionHandler) FetchTransactionByID(w http.ResponseWriter, r *http.Request) {
	id := chi.URLParam(r, "id")
	if id == "" {
		localization.SendErrorByCodeResponse(w, localization.ErrorTransactionIDRequired.Code)
		return
	}
	transaction, err := t.service.FetchTransactionByID(r.Context(), id)
	if err != nil {
		t.logger.Errorf("FetchTransactionByID failed: %v", err)
		localization.SendErrorByCodeResponse(w, err.Error())
		return
	}
	localization.SendSuccessResponse(w, localization.SuccessTransactionRetrieved, transaction)

}

func NewTransactionHandler(service service.TransactionService, logger utils.Logger) transaction.TransactionInterface {
	return &TransactionHandler{
		service: service,
		logger:  logger,
	}
}
