package transaction_repo

import (
	transaction_dto "cbe-super-app-cps-action/internal/constants/dto/transaction"
	"cbe-super-app-cps-action/internal/constants/localization"
	"cbe-super-app-cps-action/internal/constants/types"
	"cbe-super-app-cps-action/internal/storage"
	local_util "cbe-super-app-cps-action/pkgs/utils"
	"context"
	"database/sql"
	"errors"
	"fmt"
	"strings"

	"gitlab.com/bersufekadgetachew/cbe-super-app-shared/shared/utils"
)

type TransactionRepository struct {
	db     *sql.DB
	logger utils.Logger
}

const (
	vaultTransactionsTable = "VAULT_TRANSACTIONS"
	maxPaginationDefault   = 50
)

func parseBoolFilter(v interface{}) (bool, bool) {
	switch t := v.(type) {
	case bool:
		return t, true
	case int:
		return t != 0, true
	case int32:
		return t != 0, true
	case int64:
		return t != 0, true
	case float64:
		return t != 0, true
	case string:
		switch strings.ToLower(strings.TrimSpace(t)) {
		case "true", "1", "yes", "y":
			return true, true
		case "false", "0", "no", "n":
			return false, true
		}
	}
	return false, false
}

func (t *TransactionRepository) scanVaultTransaction(rowScanner interface {
	Scan(dest ...interface{}) error
}) (transaction_dto.VaultTransaction, error) {
	var out transaction_dto.VaultTransaction
	var (
		ftNumber, vaultID, vaultTxType sql.NullString
		principalDelta                 sql.NullString
		interestDelta                  sql.NullString
		balanceAfter                   sql.NullString
		referenceType                  sql.NullString
		referenceID                    sql.NullString
		debitUserID                    sql.NullString
		debitAccNo                     sql.NullString
		debitAccHolder                 sql.NullString
		creditAccNo                    sql.NullString
		creditAccHolder                sql.NullString
		serviceFee                     sql.NullString
		paidAmount                     sql.NullString
		vat                            sql.NullString
		amount                         sql.NullString
		totalAmount                    sql.NullString
		externalReference              sql.NullString
		transactionReason              sql.NullString
		receiptLink                    sql.NullString
		status                         sql.NullString
		createdAt                      sql.NullTime
	)

	err := rowScanner.Scan(
		&out.ID,
		&out.TransactionID,
		&ftNumber,
		&vaultID,
		&vaultTxType,
		&principalDelta,
		&interestDelta,
		&balanceAfter,
		&referenceType,
		&referenceID,
		&debitUserID,
		&debitAccNo,
		&debitAccHolder,
		&creditAccNo,
		&creditAccHolder,
		&out.Currency,
		&serviceFee,
		&paidAmount,
		&vat,
		&amount,
		&totalAmount,
		&externalReference,
		&transactionReason,
		&receiptLink,
		&out.TransactionType,
		&out.IsIFB,
		&status,
		&createdAt,
	)
	if err != nil {
		return transaction_dto.VaultTransaction{}, err
	}

	out.FtNumber = ftNumber.String
	out.VaultID = vaultID.String
	out.VaultTxType = vaultTxType.String
	out.PrincipalDelta = principalDelta.String
	out.InterestDelta = interestDelta.String
	out.BalanceAfter = balanceAfter.String
	out.ReferenceType = referenceType.String
	out.ReferenceID = referenceID.String
	out.DebitUserID = debitUserID.String
	out.DebitAccountNumber = debitAccNo.String
	out.DebitAccountHolderName = debitAccHolder.String
	out.CreditAccountNumber = creditAccNo.String
	out.CreditAccountHolderName = creditAccHolder.String
	out.ServiceFee = serviceFee.String
	out.PaidAmount = paidAmount.String
	out.VAT = vat.String
	out.Amount = amount.String
	out.TotalAmount = totalAmount.String
	out.ExternalReference = externalReference.String
	out.TransactionReason = transactionReason.String
	out.ReceiptLink = receiptLink.String
	out.Status = status.String
	if createdAt.Valid {
		out.CreatedAt = createdAt.Time
	}

	return out, nil
}

// FindTransactionByCifOrAccountNumberOrFT implements storage.TransactionRepository.
func (t *TransactionRepository) FindTransactionByCifOrAccountNumberOrFT(ctx context.Context, identifier string) (transaction_dto.VaultTransaction, error) {
	log := local_util.LoggerFromCtx(ctx, t.logger)

	identifier = strings.TrimSpace(identifier)
	if identifier == "" {
		return transaction_dto.VaultTransaction{}, errors.New(localization.ErrorNoDataProvided.Code)
	}

	const q = `
SELECT
  ID,
  TRANSACTION_ID,
  FT_NUMBER,
  VAULT_ID,
  VAULT_TX_TYPE,
  TO_CHAR(PRINCIPAL_DELTA),
  TO_CHAR(INTEREST_DELTA),
  TO_CHAR(BALANCE_AFTER),
  REFERENCE_TYPE,
  REFERENCE_ID,
  DEBIT_USER_ID,
  DEBIT_ACCOUNT_NUMBER,
  DEBIT_ACCOUNT_HOLDER_NAME,
  CREDIT_ACCOUNT_NUMBER,
  CREDIT_ACCOUNT_HOLDER_NAME,
  CURRENCY,
  TO_CHAR(SERVICE_FEE),
  TO_CHAR(PAID_AMOUNT),
  TO_CHAR(VAT),
  TO_CHAR(AMOUNT),
  TO_CHAR(TOTAL_AMOUNT),
  EXTERNAL_REFERENCE,
  TRANSACTION_REASON,
  RECEIPT_LINK,
  TRANSACTION_TYPE,
  IS_IFB,
  STATUS,
  CREATED_AT
FROM VAULT_TRANSACTIONS
WHERE UPPER(TRANSACTION_ID) = UPPER(:1)
   OR UPPER(FT_NUMBER) = UPPER(:1)
   OR UPPER(DEBIT_ACCOUNT_NUMBER) = UPPER(:1)
ORDER BY CREATED_AT DESC
FETCH FIRST 1 ROWS ONLY`

	row := t.db.QueryRowContext(ctx, q, identifier)
	transaction, err := t.scanVaultTransaction(row)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return transaction_dto.VaultTransaction{}, errors.New(localization.ErrorResourceNotFound.Code)
		}
		log.Errorf("failed to find transaction by CIF/AccountNumber/FT: %v", err)
		return transaction_dto.VaultTransaction{}, local_util.HandleDBError(err)
	}
	return transaction, nil
}

// FindAllWithPagination implements storage.TransactionRepository.
func (t *TransactionRepository) FindAllWithPagination(ctx context.Context, filterParam types.Filter) (*types.PaginatedResponse[[]transaction_dto.VaultTransaction], error) {
	log := local_util.LoggerFromCtx(ctx, t.logger)

	limit := int64(maxPaginationDefault)
	page := int64(1)
	if filterParam.PerPage > 0 {
		limit = int64(filterParam.PerPage)
	}
	if filterParam.Page > 0 {
		page = int64(filterParam.Page)
	}
	offset := (page - 1) * limit

	clauses := []string{"1=1"}
	var args []interface{}

	search := strings.TrimSpace(filterParam.Search)
	if search != "" {
		clauses = append(clauses,
			`(
				LOWER(TRANSACTION_ID) LIKE '%' || LOWER(:search) || '%'
				OR LOWER(FT_NUMBER) LIKE '%' || LOWER(:search) || '%'
				OR LOWER(DEBIT_ACCOUNT_NUMBER) LIKE '%' || LOWER(:search) || '%'
				OR LOWER(CREDIT_ACCOUNT_NUMBER) LIKE '%' || LOWER(:search) || '%'
				OR LOWER(STATUS) LIKE '%' || LOWER(:search) || '%'
			)`,
		)
		args = append(args, sql.Named("search", search))
	}

	if filterParam.Filters != nil {
		if v, ok := filterParam.Filters["status"]; ok {
			if s, ok2 := v.(string); ok2 && strings.TrimSpace(s) != "" {
				clauses = append(clauses, "LOWER(STATUS) = LOWER(:status)")
				args = append(args, sql.Named("status", s))
			}
		}
		if v, ok := filterParam.Filters["transaction_type"]; ok {
			if s, ok2 := v.(string); ok2 && strings.TrimSpace(s) != "" {
				clauses = append(clauses, "LOWER(TRANSACTION_TYPE) = LOWER(:transaction_type)")
				args = append(args, sql.Named("transaction_type", s))
			}
		}
		if v, ok := filterParam.Filters["is_ifb"]; ok {
			if b, ok2 := parseBoolFilter(v); ok2 {
				if b {
					clauses = append(clauses, "IS_IFB = 1")
				} else {
					clauses = append(clauses, "IS_IFB = 0")
				}
			}
		}
	}

	where := strings.Join(clauses, " AND ")
	countQ := fmt.Sprintf(`SELECT COUNT(*) FROM %s WHERE %s`, vaultTransactionsTable, where)

	var total int64
	if err := t.db.QueryRowContext(ctx, countQ, args...).Scan(&total); err != nil {
		log.Errorf("failed to count transactions with param: %v", err)
		return nil, local_util.HandleDBError(err)
	}

	listQ := fmt.Sprintf(`
SELECT
  ID,
  TRANSACTION_ID,
  FT_NUMBER,
  VAULT_ID,
  VAULT_TX_TYPE,
  TO_CHAR(PRINCIPAL_DELTA),
  TO_CHAR(INTEREST_DELTA),
  TO_CHAR(BALANCE_AFTER),
  REFERENCE_TYPE,
  REFERENCE_ID,
  DEBIT_USER_ID,
  DEBIT_ACCOUNT_NUMBER,
  DEBIT_ACCOUNT_HOLDER_NAME,
  CREDIT_ACCOUNT_NUMBER,
  CREDIT_ACCOUNT_HOLDER_NAME,
  CURRENCY,
  TO_CHAR(SERVICE_FEE),
  TO_CHAR(PAID_AMOUNT),
  TO_CHAR(VAT),
  TO_CHAR(AMOUNT),
  TO_CHAR(TOTAL_AMOUNT),
  EXTERNAL_REFERENCE,
  TRANSACTION_REASON,
  RECEIPT_LINK,
  TRANSACTION_TYPE,
  IS_IFB,
  STATUS,
  CREATED_AT
FROM %s
WHERE %s
ORDER BY CREATED_AT DESC
OFFSET :offset ROWS FETCH NEXT :limit ROWS ONLY`, vaultTransactionsTable, where)

	listArgs := append(args, sql.Named("offset", offset), sql.Named("limit", limit))
	rows, err := t.db.QueryContext(ctx, listQ, listArgs...)
	if err != nil {
		log.Errorf("failed to list transactions with param: %v", err)
		return nil, local_util.HandleDBError(err)
	}
	defer rows.Close()

	var list []transaction_dto.VaultTransaction
	for rows.Next() {
		item, err := t.scanVaultTransaction(rows)
		if err != nil {
			log.Errorf("failed to scan transaction row: %v", err)
			return nil, local_util.HandleDBError(err)
		}
		list = append(list, item)
	}

	meta := local_util.BuildPaginationMeta(total, int(page), int(limit))
	return &types.PaginatedResponse[[]transaction_dto.VaultTransaction]{Data: list, Meta: meta}, nil
}

// FindTransactionByID implements storage.TransactionRepository.
func (t *TransactionRepository) FindTransactionByID(ctx context.Context, id string) (transaction_dto.VaultTransaction, error) {
	log := local_util.LoggerFromCtx(ctx, t.logger)

	id = strings.TrimSpace(id)
	if id == "" {
		return transaction_dto.VaultTransaction{}, errors.New(localization.ErrorTransactionIDRequired.Code)
	}

	const q = `
SELECT
  ID,
  TRANSACTION_ID,
  FT_NUMBER,
  VAULT_ID,
  VAULT_TX_TYPE,
  TO_CHAR(PRINCIPAL_DELTA),
  TO_CHAR(INTEREST_DELTA),
  TO_CHAR(BALANCE_AFTER),
  REFERENCE_TYPE,
  REFERENCE_ID,
  DEBIT_USER_ID,
  DEBIT_ACCOUNT_NUMBER,
  DEBIT_ACCOUNT_HOLDER_NAME,
  CREDIT_ACCOUNT_NUMBER,
  CREDIT_ACCOUNT_HOLDER_NAME,
  CURRENCY,
  TO_CHAR(SERVICE_FEE),
  TO_CHAR(PAID_AMOUNT),
  TO_CHAR(VAT),
  TO_CHAR(AMOUNT),
  TO_CHAR(TOTAL_AMOUNT),
  EXTERNAL_REFERENCE,
  TRANSACTION_REASON,
  RECEIPT_LINK,
  TRANSACTION_TYPE,
  IS_IFB,
  STATUS,
  CREATED_AT
FROM VAULT_TRANSACTIONS
WHERE ID = :1
FETCH FIRST 1 ROWS ONLY`

	row := t.db.QueryRowContext(ctx, q, id)
	transaction, err := t.scanVaultTransaction(row)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return transaction_dto.VaultTransaction{}, errors.New(localization.ErrorResourceNotFound.Code)
		}
		log.Errorf("failed to find transaction by id: %v", err)
		return transaction_dto.VaultTransaction{}, local_util.HandleDBError(err)
	}
	return transaction, nil
}

func NewTransactionRepository(db *sql.DB, logger utils.Logger) storage.TransactionRepository {
	return &TransactionRepository{
		db:     db,
		logger: logger,
	}
}
