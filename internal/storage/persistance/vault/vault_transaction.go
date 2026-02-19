package vault

import (
	imodel "cbe-super-app-cps-action/internal/constants/model"
	"cbe-super-app-cps-action/internal/constants/types"
	local_util "cbe-super-app-cps-action/pkgs/utils"
	"context"
	"database/sql"
)

const (
	countTransactions = `SELECT COUNT(*) FROM vault_transactions`
	listTransactions  = `SELECT
		id,
		amount,
		balance_after,
		created_at,
		credit_account_holder_name,
		credit_account_number,
		currency,
		debit_account_holder_name,
		debit_account_number,
		debit_user_id,
		external_reference,
		ft_number,
		interest_delta,
		is_ifb,
		paid_amount,
		principal_delta,
		receipt_link,
		reference_id,
		reference_type,
		service_fee,
		total_amount,
		transaction_id,
		transaction_reason,
		transaction_type,
		vat,
		vault_id,
		vault_tx_type
	FROM vault_transactions
	ORDER BY created_at DESC
	OFFSET :offset ROWS FETCH NEXT :limit ROWS ONLY`

	selectTransactionByID = `SELECT
		id,
		amount,
		balance_after,
		created_at,
		credit_account_holder_name,
		credit_account_number,
		currency,
		debit_account_holder_name,
		debit_account_number,
		debit_user_id,
		external_reference,
		ft_number,
		interest_delta,
		is_ifb,
		paid_amount,
		principal_delta,
		receipt_link,
		reference_id,
		reference_type,
		service_fee,
		total_amount,
		transaction_id,
		transaction_reason,
		transaction_type,
		vat,
		vault_id,
		vault_tx_type
	FROM vault_transactions
	WHERE id = :1`
)

func (r *VaultCategoryRepository) FindAllTransactionsWithPagination(ctx context.Context, filterParam types.Filter) (*types.PaginatedResponse[[]imodel.VaultTransaction], error) {
	limit := int64(50)
	page := int64(1)
	if filterParam.PerPage > 0 {
		limit = int64(filterParam.PerPage)
	}
	if filterParam.Page > 0 {
		page = int64(filterParam.Page)
	}
	offset := (page - 1) * limit

	var total int64
	if err := r.db.QueryRowContext(ctx, countTransactions).Scan(&total); err != nil {
		r.logger.Errorf("failed to count withdrawal requests: %v", err)
		return nil, local_util.HandleDBError(err)
	}

	rows, err := r.db.QueryContext(ctx, listTransactions, sql.Named("offset", offset), sql.Named("limit", limit))
	if err != nil {
		r.logger.Errorf("failed to list withdrawal requests: %v", err)
		return nil, local_util.HandleDBError(err)
	}
	defer rows.Close()

	var list []imodel.VaultTransaction

	for rows.Next() {
		var w imodel.VaultTransaction

		if err := rows.Scan(
			// 1–5
			&w.ID,
			&w.Amount,
			&w.BalanceAfter,
			&w.CreatedAt,
			&w.CreditAccountHolderName,

			// 6–10
			&w.CreditAccountNumber,
			&w.Currency,
			&w.DebitAccountHolderName,
			&w.DebitAccountNumber,
			&w.DebitUserID,

			// 11–15
			&w.ExternalReference,
			&w.FTNumber,
			&w.InterestDelta,
			&w.IsIFB,
			&w.PaidAmount,

			// 16–20
			&w.PrincipalDelta,
			&w.ReceiptLink,
			&w.ReferenceID,
			&w.ReferenceType,
			&w.ServiceFee,

			// 21–25
			&w.TotalAmount,
			&w.TransactionID,
			&w.TransactionReason,
			&w.TransactionType,
			&w.VAT,

			// 26–27
			&w.VaultID,
			&w.VaultTxType,
		); err != nil {
			r.logger.Errorf("failed to scan transaction row: %v", err)
			return nil, local_util.HandleDBError(err)
		}

		list = append(list, w)
	}

	meta := local_util.BuildPaginationMeta(total, int(page), int(limit))
	return &types.PaginatedResponse[[]imodel.VaultTransaction]{
		Data: list,
		Meta: meta,
	}, nil
}

func (r VaultCategoryRepository) FindVaultTransaction(ctx context.Context, id string) (*imodel.VaultTransaction, error) {
	var w imodel.VaultTransaction
	err := r.db.QueryRowContext(ctx, selectTransactionByID, id).Scan(
		// 1–5
		&w.ID,
		&w.Amount,
		&w.BalanceAfter,
		&w.CreatedAt,
		&w.CreditAccountHolderName,

		// 6–10
		&w.CreditAccountNumber,
		&w.Currency,
		&w.DebitAccountHolderName,
		&w.DebitAccountNumber,
		&w.DebitUserID,

		// 11–15
		&w.ExternalReference,
		&w.FTNumber,
		&w.InterestDelta,
		&w.IsIFB,
		&w.PaidAmount,

		// 16–20
		&w.PrincipalDelta,
		&w.ReceiptLink,
		&w.ReferenceID,
		&w.ReferenceType,
		&w.ServiceFee,

		// 21–25
		&w.TotalAmount,
		&w.TransactionID,
		&w.TransactionReason,
		&w.TransactionType,
		&w.VAT,

		// 26–27
		&w.VaultID,
		&w.VaultTxType,
	)
	if err != nil {
		r.logger.Errorf("failed to get vault transaction for id %s: %v", id, err)
		return nil, local_util.HandleDBError(err)
	}
	return &w, nil
}
