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

type paginated_transaction_resp types.PaginatedResponse[[]transaction_dto.VaultTransaction]
type transaction_by_id transaction_dto.VaultTransaction

type TransactionHandler struct {
	service service.TransactionService
	logger  utils.Logger
}

func NewTransactionHandler(service service.TransactionService, logger utils.Logger) transaction.TransactionInterface {
	return &TransactionHandler{
		service: service,
		logger:  logger,
	}
}

// FindTransactionByCifOrAccountNumberOrFT godoc
//
//	@Summary		Search transaction by CIF, account number, or FT number
//	@Description	Search for transactions using CIF number, account number, or FT (Fund Transfer) number
//	@Tags			Transactions
//	@Accept			json
//	@Produce		json
//	@Param			identifier	path		string													true	"CIF number, account number, or FT number"
//	@Success		200			{object}	localization.StandardResponse{data=transaction_dto.FullTransaction}	"Transaction retrieved successfully"
//	@Failure		400			{object}	localization.StandardResponse{data=nil}					"Bad request"
//	@Failure		404			{object}	localization.StandardResponse{data=nil}					"Transaction not found"
//	@Failure		500			{object}	localization.StandardResponse{data=nil}					"Internal server error"
//	@Security		BearerAuth
//	@Router			/transactions/search/{identifier} [get]
func (t *TransactionHandler) FindTransactionByCifOrAccountNumberOrFT(w http.ResponseWriter, r *http.Request) {
	ctx, span := local_utils.TraceLogger(r.Context(), "handler", "transaction", "TransactionHandler", "FindTransactionByCifOrAccountNumberOrFT")
	defer span.End()
	log := local_utils.LoggerFromCtx(ctx, t.logger)

	identifier := chi.URLParam(r, "identifier")
	if identifier == "" {
		span.AddEvent("Missing identifier", trace.WithAttributes(attribute.String("error", "identifier required")))
		localization.SendErrorByCodeResponse(w, localization.ErrorTransactionIdentifierRequired.Code)
		return
	}
	transaction, err := t.service.FindTransactionByCifOrAccountNumberOrFT(ctx, identifier)
	if err != nil {
		span.AddEvent("Service error", trace.WithAttributes(attribute.String("error", err.Error()), attribute.String("identifier", identifier)))
		log.Errorf("[TxnH][FindByCif] svc err: %v", err)
		localization.SendErrorByCodeResponse(w, err.Error())
		return
	}
	localization.SendSuccessResponse(w, localization.SuccessTransactionRetrieved, transaction)
}

// FetchAllTransactions godoc
//
//	@Summary		Get all transactions
//	@Description	Returns a paginated list of transactions. Supports filtering by status and type.
//	@Tags			Transactions
//	@Accept			json
//	@Produce		json
//	@Param			status		query		string								false	"Transaction status (failed, pending, paid)"
//	@Param			type		query		string								false	"Transaction type (cbe, topup, money_request)"
//	@Param			page		query		int									false	"Page number"		default(1)
//	@Param			per_page	query		int									false	"Items per page"	default(10)
//	@Param			search		query		string								false	"Search term"
//	@Success		200			{object}	localization.StandardResponse{data=object}	"Transactions retrieved successfully"
//	@Failure		400			{object}	localization.StandardResponse{data=nil}		"Bad request"
//	@Failure		500			{object}	localization.StandardResponse{data=nil}		"Internal server error"
//	@Security		BearerAuth
//	@Router			/transactions [get]
//
// FetchAllTransactions implements transaction.TransactionInterface.
func (t *TransactionHandler) FetchAllTransactions(w http.ResponseWriter, r *http.Request) {
	ctx, span := local_utils.TraceLogger(r.Context(), "handler", "transaction", "TransactionHandler", "FetchAllTransactions")
	defer span.End()
	log := local_utils.LoggerFromCtx(ctx, t.logger)

	filterParams, err := local_utils.ExtractFilterParams(r)
	if err != nil {
		localization.SendErrorByCodeResponse(w, err.Error())
		return
	}

	search := r.URL.Query().Get("search")
	filter := r.URL.Query().Get("filter")

	if err := local_utils.NoSpecialChars(search); err != nil {
		log.Errorf("[TxnH][GetAll] invalid search query: %v", err)
		localization.SendErrorByCodeResponse(w, err.Error())
		return
	}

	if err := local_utils.NoSpecialChars(filter); err != nil {
		localization.SendErrorByCodeResponse(w, err.Error())
		return
	}

	transactions, err := t.service.FetchAllTransactions(ctx, filterParams)
	if err != nil {
		span.AddEvent("Service error", trace.WithAttributes(attribute.String("error", err.Error())))
		log.Errorf("[TxnH][GetAll] svc err: %v", err)
		localization.SendErrorByCodeResponse(w, err.Error())
		return
	}
	span.AddEvent("Transactions retrieved", trace.WithAttributes(attribute.Int("count", len(transactions.Data))))
	localization.SendSuccessResponse(w, localization.SuccessTransactionRetrieved, transactions)
}

// FetchTransactionByID godoc
//
//	@Summary		Get transaction by ID
//	@Description	Returns a single transaction by its ID.
//	@Tags			Transactions
//	@Accept			json
//	@Produce		json
//	@Param			id		path		string													true	"Transaction ID"
//	@Success		200		{object}	localization.StandardResponse{data=transaction_dto.FullTransaction}	"Transaction retrieved successfully"
//	@Failure		400		{object}	localization.StandardResponse{data=nil}					"Bad request"
//	@Failure		404		{object}	localization.StandardResponse{data=nil}					"Transaction not found"
//	@Failure		500		{object}	localization.StandardResponse{data=nil}					"Internal server error"
//	@Security		BearerAuth
//	@Router			/transactions/{id} [get]
//
// FetchTransactionLimitByUserCode implements transaction.TransactionInterface.
func (t *TransactionHandler) FetchTransactionLimitByUserCode(w http.ResponseWriter, r *http.Request) {
	ctx, span := local_utils.TraceLogger(r.Context(), "handler", "transaction", "TransactionHandler", "FetchTransactionLimitByUserCode")
	defer span.End()
	log := local_utils.LoggerFromCtx(ctx, t.logger)
	filterParams, err := local_utils.ExtractFilterParams(r)
	if err != nil {
		localization.SendErrorByCodeResponse(w, err.Error())
		return
	}

	customerNumber := chi.URLParam(r, "customer_no")
	if customerNumber == "" {
		span.AddEvent("Missing customer_number", trace.WithAttributes(attribute.String("error", "customer_number required")))
		localization.SendErrorByCodeResponse(w, localization.ErrorNoDataProvided.Code)
		return
	}

	limit, err := t.service.FetchTransactionLimitByCustomerNumber(ctx, filterParams, customerNumber)
	if err != nil {
		span.AddEvent("Service error", trace.WithAttributes(attribute.String("error", err.Error()), attribute.String("customer_number", customerNumber)))
		log.Errorf("[TxnH][GetLimit] svc err: %v", err)
		localization.SendErrorByCodeResponse(w, err.Error())
		return
	}
	span.AddEvent("Transaction limit retrieved", trace.WithAttributes(attribute.String("customer_number", customerNumber)))
	localization.SendSuccessResponse(w, localization.SuccessTransactionLimitRetrieved, limit)
}

// FetchAllTransactionLimits godoc
//
//	@Summary		Get transaction limits
//	@Description	When ?search=<customer_number> is provided, fetches real-time limits for that customer from CoreIO.
//	@Description	Without a search value, returns a paginated list from the database.
//	@Tags			Transactions
//	@Produce		json
//	@Param			search		query		string	false	"Customer number — triggers single-customer CoreIO lookup"
//	@Param			page		query		int		false	"Page number (paginated mode)"		default(1)
//	@Param			per_page	query		int		false	"Items per page (paginated mode)"	default(10)
//	@Success		200	{object}	localization.StandardResponse{data=object}	"Transaction limits retrieved successfully"
//	@Failure		400	{object}	localization.StandardResponse{data=nil}		"Bad request"
//	@Failure		500	{object}	localization.StandardResponse{data=nil}		"Internal server error"
//	@Security		BearerAuth
//	@Router			/transaction [get]
func (t *TransactionHandler) FetchAllTransactionLimits(w http.ResponseWriter, r *http.Request) {
	ctx, span := local_utils.TraceLogger(r.Context(), "handler", "transaction", "TransactionHandler", "FetchAllTransactionLimits")
	defer span.End()
	log := local_utils.LoggerFromCtx(ctx, t.logger)

	customerNumber := r.URL.Query().Get("search")
	filterParams, err := local_utils.ExtractFilterParams(r)
	if err != nil {
		localization.SendErrorByCodeResponse(w, err.Error())
		return
	}
	if customerNumber != "" {
		// Single-customer lookup via CoreIO
		limit, err := t.service.FetchTransactionLimitByCustomerNumber(ctx, filterParams, customerNumber)
		if err != nil {
			span.AddEvent("Service error", trace.WithAttributes(attribute.String("error", err.Error()), attribute.String("customer_number", customerNumber)))
			log.Errorf("[TxnH][GetLimit] svc err customer=%s: %v", customerNumber, err)
			localization.SendErrorByCodeResponse(w, err.Error())
			return
		}

		span.AddEvent("Transaction limit retrieved", trace.WithAttributes(attribute.String("customer_number", customerNumber)))
		localization.SendSuccessResponse(w, localization.SuccessTransactionLimitRetrieved, limit)
		return
	}

	limits, err := t.service.FetchAllTransactionLimits(ctx, *filterParams)
	if err != nil {
		span.AddEvent("Service error", trace.WithAttributes(attribute.String("error", err.Error())))
		log.Errorf("[TxnH][GetAllLimits] svc err: %v", err)
		localization.SendErrorByCodeResponse(w, err.Error())
		return
	}

	span.AddEvent("Transaction limits retrieved", trace.WithAttributes(attribute.Int("count", len(limits.Data))))
	localization.SendSuccessResponse(w, localization.SuccessTransactionLimitRetrieved, limits)
}

// FetchTransactionByID implements transaction.TransactionInterface.
func (t *TransactionHandler) FetchTransactionByID(w http.ResponseWriter, r *http.Request) {
	ctx, span := local_utils.TraceLogger(r.Context(), "handler", "transaction", "TransactionHandler", "FetchTransactionByID")
	defer span.End()
	log := local_utils.LoggerFromCtx(ctx, t.logger)
	id := chi.URLParam(r, "id")
	if id == "" {
		span.AddEvent("Missing transaction ID", trace.WithAttributes(attribute.String("error", "transaction ID required")))
		localization.SendErrorByCodeResponse(w, localization.ErrorTransactionIDRequired.Code)
		return
	}
	transaction, err := t.service.FetchTransactionByID(ctx, id)
	if err != nil {
		span.AddEvent("Service error", trace.WithAttributes(attribute.String("error", err.Error()), attribute.String("id", id)))
		log.Errorf("[TxnH][GetByID] svc err: %v", err)
		localization.SendErrorByCodeResponse(w, err.Error())
		return
	}
	span.AddEvent("Transaction retrieved", trace.WithAttributes(attribute.String("id", id)))
	localization.SendSuccessResponse(w, localization.SuccessTransactionRetrieved, transaction)

}
