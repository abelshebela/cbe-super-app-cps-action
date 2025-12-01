package sqlc

import (
	"context"
	"database/sql"
)

const findSitotaTransactions = `-- name: FindSitotaTransactions :many
SELECT
  transaction_id,
  debit_account_holder_name,
  debit_account_number,
  debit_user_id,
  credit_account_holder_name,
  credit_account_number,
  credit_user_id,
  amount,
  institution_code,
  transaction_status,
  paid_at,
  created_at,
  last_modified_at,
  COUNT(*) OVER() AS total_count
FROM transaction
WHERE transaction_type = 'sitota'
ORDER BY created_at DESC
OFFSET NVL(:offset, 0) ROWS
FETCH NEXT NVL(:limit, 50) ROWS ONLY
`

const findSitotaTransactionByID = `-- name: FindSitotaTransactionByID :one
SELECT
  transaction_id,
  debit_account_holder_name,
  debit_account_number,
  debit_user_id,
  credit_account_holder_name,
  credit_account_number,
  credit_user_id,
  amount,
  paid_at,
  institution_code,
  transaction_status,
  created_at,
  last_modified_at
FROM transaction
WHERE transaction_type = 'sitota'
  AND transaction_id = :1`

type FindSitotaTransactionsParams struct {
	// Status sql.NullString `json:"status"`
	// Search sql.NullString `json:"search"`
	Page  sql.NullInt64 `json:"page"`
	Limit sql.NullInt64 `json:"limit"`
}

func (q *Queries) FindSitotaTransactions(ctx context.Context, arg FindSitotaTransactionsParams) ([]FindSitotaTransactionRow, error) {
	offset := (arg.Page.Int64 - 1) * arg.Limit.Int64
	rows, err := q.db.QueryContext(ctx, findSitotaTransactions,
		sql.Named("offset", offset),
		sql.Named("limit", arg.Limit),
	)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var items []FindSitotaTransactionRow
	for rows.Next() {
		var r FindSitotaTransactionRow
		if err := rows.Scan(
			&r.ID,
			&r.SenderName,
			&r.SenderAccountNumber,
			&r.SenderPhoneNumber,
			&r.RecipientName,
			&r.RecipientAccountNumber,
			&r.RecipientPhoneNumber,
			&r.SitotaAmount,
			&r.GLAccountNumber,
			&r.Status,
			&r.PaidAt,
			&r.CreatedAt,
			&r.UpdatedAt,
			&r.TotalCount,
		); err != nil {
			return nil, err
		}
		items = append(items, r)
	}
	return items, nil
}

func (q *Queries) FindSitotaTransactionByID(ctx context.Context, id string) (FindSitotaTransactionRow, error) {
	row := q.db.QueryRowContext(ctx, findSitotaTransactionByID, id)
	var r FindSitotaTransactionRow
	err := row.Scan(
		&r.ID,
		&r.SenderName,
		&r.SenderAccountNumber,
		&r.SenderPhoneNumber,
		&r.RecipientName,
		&r.RecipientAccountNumber,
		&r.RecipientPhoneNumber,
		&r.SitotaAmount,
		&r.GLAccountNumber,
		&r.Status,
		&r.PaidAt,
		&r.CreatedAt,
		&r.UpdatedAt,
	)
	return r, err
}
