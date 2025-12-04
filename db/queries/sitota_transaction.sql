-- name: FindSitotaTransactions :many
SELECT
  t.transaction_id                AS id,
  t.debit_account_holder_name     AS sender_name,
  t.debit_account_number          AS sender_account_number,
  t.debit_user_id                 AS sender_phone_number,
  t.credit_account_holder_name    AS recipient_name,
  t.credit_account_number         AS recipient_account_number,
  t.credit_user_id                AS recipient_phone_number,
  t.amount::float                 AS sitota_amount,
  t.institution_code              AS gl_account_number,
  t.transaction_status            AS status,
  NULL::timestamp                 AS claimed_at,
  t.created_at,
  t.last_modified_at              AS updated_at,
  NULL::timestamp                 AS deleted_at,
  COUNT(*) OVER()                 AS total_count
FROM transactions t
WHERE t.transaction_type = 'sitota'
  -- Optional filters:
  AND ($1::text IS NULL OR t.transaction_status = $1::text)
  AND ($2::text IS NULL OR t.transaction_id ILIKE ('%' || $2 || '%'))
ORDER BY t.created_at DESC
LIMIT  COALESCE($4::int, 50)
OFFSET COALESCE($3::int, 0);

-- name: FindSitotaTransactionById :one
SELECT
  t.transaction_id                AS id,
  t.debit_account_holder_name     AS sender_name,
  t.debit_account_number          AS sender_account_number,
  t.debit_user_id                 AS sender_phone_number,
  t.credit_account_holder_name    AS recipient_name,
  t.credit_account_number         AS recipient_account_number,
  t.credit_user_id                AS recipient_phone_number,
  t.amount::float                 AS sitota_amount,
  t.institution_code              AS gl_account_number,
  t.transaction_status            AS status,
  NULL::timestamp                 AS claimed_at,
  t.created_at,
  t.last_modified_at              AS updated_at,
  NULL::timestamp                 AS deleted_at
FROM transactions t
WHERE t.transaction_id = $1
  AND t.transaction_type = 'sitota';
