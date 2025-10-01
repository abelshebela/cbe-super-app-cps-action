-- name: SaveBankVault :one
INSERT INTO bank_vault_products (
  name,
  description,
  currency,
  rate_bps,
  method,
  frequency,
  lock_period,
  min_amount,
  max_amount,
  early_unlock_fee_bps,
  is_active
) VALUES (
  sqlc.arg(name),
  sqlc.arg(description),
  sqlc.arg(currency),
  sqlc.arg(rate_bps),
  sqlc.arg(method)::accrual_method,
  sqlc.arg(frequency)::accrual_frequency,
  sqlc.arg(lock_period),
  sqlc.arg(min_amount),
  sqlc.arg(max_amount),
  sqlc.arg(early_unlock_fee_bps),
  sqlc.arg(is_active)
)
RETURNING id;

-- name: FindBankVaultById :one
SELECT
  id,
  name,
  description,
  currency,
  rate_bps,
  method,
  frequency,
  lock_period,
  min_amount,
  max_amount,
  early_unlock_fee_bps,
  is_active,
  created_at,
  updated_at,
  deleted_at
FROM bank_vault_products
WHERE id = sqlc.arg(id) AND deleted_at IS NULL;

-- name: FindBankVault :many
SELECT
  id,
  name,
  description,
  currency,
  rate_bps,
  method,
  frequency,
  lock_period,
  min_amount,
  max_amount,
  early_unlock_fee_bps,
  is_active,
  created_at,
  updated_at,
  deleted_at,
  COUNT(*) OVER() AS total_count
FROM bank_vault_products
WHERE deleted_at IS NULL
  AND (sqlc.narg(is_active)::boolean IS NULL OR is_active = sqlc.narg(is_active)::boolean)
  AND (sqlc.narg(name)::text IS NULL OR name ILIKE ('%' || sqlc.narg(name)::text || '%'))
  AND (sqlc.narg(currency)::text   IS NULL OR currency = sqlc.narg(currency)::text)
  AND (sqlc.narg(method)::accrual_method IS NULL OR method = sqlc.narg(method)::accrual_method)
  AND (sqlc.narg(frequency)::accrual_frequency IS NULL OR frequency = sqlc.narg(frequency)::accrual_frequency)
ORDER BY created_at DESC
LIMIT  COALESCE(sqlc.narg(limit_count)::int, 50)
OFFSET COALESCE(sqlc.narg(offset_count)::int, 0);

-- name: UpdateBankVault :one
UPDATE bank_vault_products
SET
  description = COALESCE(sqlc.narg(description), description),
  min_amount  = COALESCE(sqlc.narg(min_amount),  min_amount),
  max_amount  = COALESCE(sqlc.narg(max_amount),  max_amount),
  updated_at  = NOW()
WHERE id = sqlc.arg(id) AND deleted_at IS NULL
RETURNING id;

-- name: DeleteBankVault :one
UPDATE bank_vault_products
SET is_deleted = TRUE, deleted_at = NOW(), updated_at = NOW()
WHERE id = sqlc.arg(id) AND deleted_at IS NULL
RETURNING id;

-- name: ActivateBankVault :one
UPDATE bank_vault_products
SET is_active = TRUE, updated_at = NOW()
WHERE id = sqlc.arg(id) AND deleted_at IS NULL
RETURNING id;

-- name: DeactivateBankVault :one
UPDATE bank_vault_products
SET is_active = FALSE, updated_at = NOW()
WHERE id = sqlc.arg(id) AND deleted_at IS NULL
RETURNING id;

