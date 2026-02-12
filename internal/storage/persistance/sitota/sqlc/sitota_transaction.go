package sqlc

import (
	"context"
	"database/sql"
	"fmt"
)

// const findSitotaTransactions = `-- name: FindSitotaTransactions :many
// SELECT
//   transaction_id,
//   debit_account_holder_name,
//   debit_account_number,
//   debit_user_id,
//   credit_account_holder_name,
//   credit_account_number,
//   credit_user_id,
//   amount,
//   institution_code,
//   transaction_status,
//   paid_at,
//   created_at,
//   last_modified_at,
//   COUNT(*) OVER() AS total_count
// FROM transaction
// WHERE transaction_type = 'sitota'
// ORDER BY created_at DESC
// OFFSET NVL(:offset, 0) ROWS
// FETCH NEXT NVL(:limit, 50) ROWS ONLY
// `

// const findSitotaTransactionByID = `-- name: FindSitotaTransactionByID :one
// SELECT
//   transaction_id,
//   debit_account_holder_name,
//   debit_account_number,
//   debit_user_id,
//   credit_account_holder_name,
//   credit_account_number,
//   credit_user_id,
//   amount,
//   institution_code,
//   transaction_status,
//   paid_at,
//   created_at,
//   last_modified_at
// FROM transaction
// WHERE transaction_type = 'sitota'
//   AND transaction_id = :1`

// type FindSitotaTransactionsParams struct {
// 	Page  sql.NullInt64 `json:"page"`
// 	Limit sql.NullInt64 `json:"limit"`
// }

// func (q *Queries) FindSitotaTransactions(ctx context.Context, arg FindSitotaTransactionsParams) ([]FindSitotaTransactionRow, error) {
// 	offset := (arg.Page.Int64 - 1) * arg.Limit.Int64
// 	rows, err := q.db.QueryContext(ctx, findSitotaTransactions,
// 		sql.Named("offset", offset),
// 		sql.Named("limit", arg.Limit),
// 	)
// 	if err != nil {
// 		return nil, err
// 	}
// 	defer rows.Close()

// 	var items []FindSitotaTransactionRow
// 	for rows.Next() {
// 		var r FindSitotaTransactionRow
// 		if err := rows.Scan(
// 			&r.ID,
// 			&r.SenderName,
// 			&r.SenderAccountNumber,
// 			&r.SenderPhoneNumber,
// 			&r.RecipientName,
// 			&r.RecipientAccountNumber,
// 			&r.RecipientPhoneNumber,
// 			&r.SitotaAmount,
// 			&r.GLAccountNumber,
// 			&r.Status,
// 			&r.PaidAt,
// 			&r.CreatedAt,
// 			&r.UpdatedAt,
// 			&r.TotalCount,
// 		); err != nil {
// 			return nil, err
// 		}
// 		items = append(items, r)
// 	}
// 	return items, nil
// }

// func (q *Queries) FindSitotaTransactionByID(ctx context.Context, id string) (FindSitotaTransactionRow, error) {
// 	row := q.db.QueryRowContext(ctx, findSitotaTransactionByID, id)
// 	var r FindSitotaTransactionRow
// 	err := row.Scan(
// 		&r.ID,
// 		&r.SenderName,
// 		&r.SenderAccountNumber,
// 		&r.SenderPhoneNumber,
// 		&r.RecipientName,
// 		&r.RecipientAccountNumber,
// 		&r.RecipientPhoneNumber,
// 		&r.SitotaAmount,
// 		&r.GLAccountNumber,
// 		&r.Status,
// 		&r.PaidAt,
// 		&r.CreatedAt,
// 		&r.UpdatedAt,
// 	)
// 	return r, err
// }

const findSitotaTransactions = `-- name: FindSitotaTransactions :many
SELECT
  id,                       -- r.ID
  debit_account_number,     -- r.SenderAccountNumber
  sender_full_name,         -- r.SenderName
  reciever_full_name,       -- r.RecipientName
  credit_account_number,    -- r.RecipientAccountNumber
  locked_amount,            -- r.SitotaAmount
  locked_id,                -- r.GLAccountNumber
  status,                   -- r.Status
  created_at,               -- r.CreatedAt
  updated_at,               -- r.UpdatedAt
  COUNT(*) OVER()           -- r.TotalCount
FROM sitota_sessions
WHERE (:search IS NULL
    OR debit_account_number LIKE '%' || :search || '%'
    OR credit_account_number LIKE '%' || :search || '%'
    OR locked_id LIKE '%' || :search || '%')
ORDER BY created_at DESC
OFFSET NVL(:offset, 0) ROWS
FETCH NEXT NVL(:limit, 50) ROWS ONLY
`

const findSitotaTransactionByID = `-- name: FindSitotaTransactionByID :one
SELECT
  id,                       -- r.ID
  debit_account_number,     -- r.SenderAccountNumber
  sender_full_name,         -- r.SenderName
  reciever_full_name,       -- r.RecipientName
  credit_account_number,    -- r.RecipientAccountNumber
  locked_amount,            -- r.SitotaAmount
  locked_id,                -- r.GLAccountNumber
  status,                   -- r.Status
  created_at,               -- r.CreatedAt
  updated_at               -- r.UpdatedAt
FROM sitota_sessions
WHERE id = :1
`

type FindSitotaTransactionsParams struct {
	Search string        `json:"search"`
	Page   sql.NullInt64 `json:"page"`
	Limit  sql.NullInt64 `json:"limit"`
}

func (q *Queries) FindSitotaTransactions(ctx context.Context, arg FindSitotaTransactionsParams) ([]FindSitotaTransactionRow, error) {
	offset := (arg.Page.Int64 - 1) * arg.Limit.Int64

	rows, err := q.db.QueryContext(ctx, findSitotaTransactions,
		sql.Named("search", arg.Search),
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
			&r.SenderAccountNumber,
			&r.SenderName,
			&r.RecipientName,
			&r.RecipientAccountNumber,
			&r.SitotaAmount,
			&r.GLAccountNumber,
			&r.Status,
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
		&r.SenderAccountNumber,
		&r.SenderName,
		&r.RecipientName,
		&r.RecipientAccountNumber,
		&r.SitotaAmount,
		&r.GLAccountNumber,
		&r.Status,
		&r.CreatedAt,
		&r.UpdatedAt,
	)
	fmt.Println("SQLC ERR:", err)
	return r, err
}
