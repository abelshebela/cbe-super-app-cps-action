package sqlc

import (
	transaction_dto "cbe-super-app-cps-action/internal/constants/dto/transaction"
	"cbe-super-app-cps-action/internal/constants/model"
	"cbe-super-app-cps-action/internal/constants/types"
	transaction_core "cbe-super-app-cps-action/internal/storage/persistance/transaction/core"
	"context"
	"fmt"
)

const findTransactionByID = `SELECT 
	id,
	transaction_id,
	ft_number,
	debit_branch_code,
	debit_district_code,
	debit_user_id,
	debit_account_number,
	debit_account_holder_name,
	credit_user_id,
	credit_account_number,
	credit_account_holder_name,
	institution_code,
	institution_name,
	currency,
	service_fee,
	tip_amount,
	paid_amount,
	vat,
	amount,
	total_amount,
	external_reference,
	transaction_reason,
	transaction_type,
	transaction_status,
	is_ifb,
	paid_at,
	reversed_at,
	metadata,
	created_at,
	last_modified_at
	FROM transaction
	`

var findTransactionWithParam = `SELECT 
	id,
	transaction_id,
	ft_number,
	debit_branch_code,
	debit_district_code,
	debit_user_id,
	debit_account_number,
	debit_account_holder_name,
	credit_user_id,
	credit_account_number,
	credit_account_holder_name,
	institution_code,
	institution_name,
	currency,
	service_fee,
	tip_amount,
	paid_amount,
	vat,
	amount,
	total_amount,
	external_reference,
	transaction_reason,
	transaction_type,
	transaction_status,
	is_ifb,
	paid_at,
	reversed_at,
	metadata,
	created_at,
	last_modified_at,

	FROM transaction
	WHERE 1=1
	`

func (q *Queries) FindTransactionByID(ctx context.Context, id string) (model.TransactionModel, error) {
	query := findTransactionByID + " WHERE transaction_id = :1"
	row := q.db.QueryRowContext(ctx, query, id)
	var transaction model.TransactionModel
	err := row.Scan(
		&transaction.ID,
		&transaction.TransactionID,
		&transaction.FTNumber,
		&transaction.DebitBranchCode,
		&transaction.DebitDistrictCode,
		&transaction.DebitUserID,
		&transaction.DebitAccountNumber,
		&transaction.DebitAccountHolderName,
		&transaction.CreditUserID,
		&transaction.CreditAccountNumber,
		&transaction.CreditAccountHolderName,
		&transaction.InstitutionCode,
		&transaction.InstitutionName,
		&transaction.Currency,
		&transaction.ServiceFee,
		&transaction.TipAmount,
		&transaction.PaidAmount,
		&transaction.VAT,
		&transaction.Amount,
		&transaction.TotalAmount,
		&transaction.ExternalReference,
		&transaction.TransactionReason,
		&transaction.TransactionType,
		&transaction.TransactionStatus,
		&transaction.IsIFB,
		&transaction.PaidAt,
		&transaction.ReversedAt,
		&transaction.Metadata,
		&transaction.CreatedAt,
		&transaction.LastModifiedAt,
	)
	if err != nil {
		return model.TransactionModel{}, err
	}
	return transaction, nil
}

func (q *Queries) FindTransactionWithParam(ctx context.Context, filterParams *types.Filter) (*types.PaginatedResponse[[]transaction_dto.FullTransaction], error) {
	// Build base query and args
	baseQuery := " FROM transaction WHERE 1=1"
	filters := ""
	args := []interface{}{}
	argIdx := 1

	if vFrom, okFrom := filterParams.Filters["date_from"]; okFrom && vFrom != "" {
		if vTo, okTo := filterParams.Filters["date_to"]; okTo && vTo != "" {
			filters += fmt.Sprintf(" AND created_at BETWEEN :%d AND :%d", argIdx, argIdx+1)
			args = append(args, vFrom, vTo)
			argIdx += 2
		}
	}
	if vType, okType := filterParams.Filters["type"]; okType && vType != "" {
		filters += fmt.Sprintf(" AND transaction_type = :%d", argIdx)
		args = append(args, vType)
		argIdx++
	}
	if vStatus, okStatus := filterParams.Filters["status"]; okStatus && vStatus != "" {
		filters += fmt.Sprintf(" AND transaction_status = :%d", argIdx)
		args = append(args, vStatus)
		argIdx++
	}
	if filterParams.Search != "" {
		filters += fmt.Sprintf(" AND (debit_account_holder_name LIKE :%d OR debit_account_number LIKE :%d OR transaction_id LIKE :%d)", argIdx, argIdx+1, argIdx+2)
		like := "%" + filterParams.Search + "%"
		args = append(args, like, like, like)
		argIdx += 3
	}

	// Count query for total docs
	countQuery := "SELECT COUNT(*)" + baseQuery + filters
	var totalDocs int64
	err := q.db.QueryRowContext(ctx, countQuery, args...).Scan(&totalDocs)
	if err != nil {
		return nil, err
	}

	// Pagination
	page := filterParams.Page
	perPage := filterParams.PerPage
	if page <= 0 {
		page = 1
	}
	if perPage <= 0 {
		perPage = 10
	}
	offset := (page - 1) * perPage

	// Data query
	dataQuery := `SELECT 
			id, transaction_id, ft_number, debit_branch_code, debit_district_code, debit_user_id, debit_account_number, debit_account_holder_name, credit_user_id, credit_account_number, credit_account_holder_name, institution_code, institution_name, currency, service_fee, tip_amount, paid_amount, vat, amount, total_amount, external_reference, transaction_reason, transaction_type, transaction_status, is_ifb, paid_at, reversed_at, metadata, created_at, last_modified_at
			` + baseQuery + filters +
		fmt.Sprintf(" ORDER BY created_at DESC OFFSET %d ROWS FETCH NEXT %d ROWS ONLY", offset, perPage)

	rows, err := q.db.QueryContext(ctx, dataQuery, args...)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var transactions []transaction_dto.FullTransaction
	for rows.Next() {
		var transaction model.TransactionModel
		err := rows.Scan(
			&transaction.ID,
			&transaction.TransactionID,
			&transaction.FTNumber,
			&transaction.DebitBranchCode,
			&transaction.DebitDistrictCode,
			&transaction.DebitUserID,
			&transaction.DebitAccountNumber,
			&transaction.DebitAccountHolderName,
			&transaction.CreditUserID,
			&transaction.CreditAccountNumber,
			&transaction.CreditAccountHolderName,
			&transaction.InstitutionCode,
			&transaction.InstitutionName,
			&transaction.Currency,
			&transaction.ServiceFee,
			&transaction.TipAmount,
			&transaction.PaidAmount,
			&transaction.VAT,
			&transaction.Amount,
			&transaction.TotalAmount,
			&transaction.ExternalReference,
			&transaction.TransactionReason,
			&transaction.TransactionType,
			&transaction.TransactionStatus,
			&transaction.IsIFB,
			// &transaction.IsReversed,
			&transaction.PaidAt,
			&transaction.ReversedAt,
			&transaction.Metadata,
			&transaction.CreatedAt,
			&transaction.LastModifiedAt,
			// &transaction.TotalCount,
		)
		if err != nil {
			return nil, err
		}
		transactions = append(transactions, transaction_core.MapFullTransactionToNative(transaction))
	}

	if err = rows.Err(); err != nil {
		return nil, err
	}

	// Build PaginationMeta
	totalPages := 0
	if totalDocs > 0 {
		totalPages = int((totalDocs + int64(perPage) - 1) / int64(perPage))
	}
	hasPrev := page > 1
	hasNext := page < totalPages
	var prevPage *int
	var nextPage *int
	if hasPrev {
		p := page - 1
		prevPage = &p
	}
	if hasNext {
		n := page + 1
		nextPage = &n
	}
	pagingCounter := 0
	if totalDocs > 0 {
		pagingCounter = offset + 1
	}

	meta := types.PaginationMeta{
		TotalDocs:     totalDocs,
		Limit:         perPage,
		TotalPages:    totalPages,
		Page:          page,
		PagingCounter: pagingCounter,
		HasPrevPage:   hasPrev,
		HasNextPage:   hasNext,
		PrevPage:      prevPage,
		NextPage:      nextPage,
	}

	return &types.PaginatedResponse[[]transaction_dto.FullTransaction]{
		Data: transactions,
		Meta: meta,
	}, nil
}
