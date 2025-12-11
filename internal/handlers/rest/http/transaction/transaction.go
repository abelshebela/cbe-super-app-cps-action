package transaction_handler

import (
	transaction_dto "cbe-super-app-cps-action/internal/constants/dto/transaction"
	"cbe-super-app-cps-action/internal/constants/interfaces/transaction"
	"cbe-super-app-cps-action/internal/constants/localization"
	"cbe-super-app-cps-action/internal/constants/types"
	"cbe-super-app-cps-action/internal/service"
	local_utils "cbe-super-app-cps-action/pkgs/utils"
	"net/http"

	"github.com/go-chi/chi/v5"
	"gitlab.com/bersufekadgetachew/cbe-super-app-shared/shared/utils"
	"go.opentelemetry.io/otel/attribute"
	"go.opentelemetry.io/otel/trace"
)

type paginated_transaction_resp types.PaginatedResponse[[]transaction_dto.FullTransaction]
type transaction_by_id transaction_dto.FullTransaction

type TransactionHandler struct {
	service service.TransactionService
	logger  utils.Logger
}

// FetchAllTransactions godoc
// @Summary      Get all transactions
// @Description  Returns a paginated list of transactions. Supports filtering by status and type.
// @Tags         transactions
// @Accept       json
// @Produce      json
// @Param        status     query   string  false  "Transaction status (failed, pending, paid)"
// @Param        type       query   string  false  "Transaction type (cbe, topup, money_request)"
// @Param        page       query   int     false  "Page number"
// @Param        per_page   query   int     false  "Items per page"
// @Success      200        {object}  paginated_transaction_resp
// @Failure      400,404,500  {object}  localization.ErrorResponse
// @Router       /transactions [get]
// FetchAllTransactions implements transaction.TransactionInterface.
func (t *TransactionHandler) FetchAllTransactions(w http.ResponseWriter, r *http.Request) {
	ctx, span := local_utils.TraceLogger(r.Context(), "handler", "transaction", "TransactionHandler", "FetchAllTransactions")
	defer span.End()
	filterParams := local_utils.ExtractFilterParams(r)

	transactions, err := t.service.FetchAllTransactions(ctx, filterParams)
	if err != nil {
		span.AddEvent("Service error", trace.WithAttributes(attribute.String("error", err.Error())))
		t.logger.Errorf("FetchAllTransactions failed: %v", err)
		localization.SendErrorByCodeResponse(w, err.Error())
		return
	}
	span.AddEvent("Transactions retrieved", trace.WithAttributes(attribute.Int("count", len(transactions.Data))))
	localization.SendSuccessResponse(w, localization.SuccessTransactionRetrieved, transactions)
}

// FetchTransactionByID godoc
// @Summary      Get transaction by ID
// @Description  Returns a single transaction by its ID.
// @Tags         transactions
// @Accept       json
// @Produce      json
// @Param        id   path   string  true  "Transaction ID"
// @Success      200  {object}  transaction_by_id
// @Failure      400,404,500  {object}  localization.ErrorResponse
// @Router       /transactions/{id} [get]
// FetchTransactionByID implements transaction.TransactionInterface.
func (t *TransactionHandler) FetchTransactionByID(w http.ResponseWriter, r *http.Request) {
	ctx, span := local_utils.TraceLogger(r.Context(), "handler", "transaction", "TransactionHandler", "FetchTransactionByID")
	defer span.End()
	id := chi.URLParam(r, "id")
	if id == "" {
		span.AddEvent("Missing transaction ID", trace.WithAttributes(attribute.String("error", "transaction ID required")))
		localization.SendErrorByCodeResponse(w, localization.ErrorTransactionIDRequired.Code)
		return
	}
	transaction, err := t.service.FetchTransactionByID(ctx, id)
	if err != nil {
		span.AddEvent("Service error", trace.WithAttributes(attribute.String("error", err.Error()), attribute.String("id", id)))
		t.logger.Errorf("FetchTransactionByID failed: %v", err)
		localization.SendErrorByCodeResponse(w, err.Error())
		return
	}
	span.AddEvent("Transaction retrieved", trace.WithAttributes(attribute.String("id", id)))
	localization.SendSuccessResponse(w, localization.SuccessTransactionRetrieved, transaction)

}

func NewTransactionHandler(service service.TransactionService, logger utils.Logger) transaction.TransactionInterface {
	return &TransactionHandler{
		service: service,
		logger:  logger,
	}
}
