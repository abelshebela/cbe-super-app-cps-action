package sqlc

import (
	"cbe-super-app-cps-action/pkgs/utils"
	"context"
	"database/sql"
	"fmt"
	"time"

	"github.com/godror/godror"
	"github.com/google/uuid"
	"github.com/shopspring/decimal"
)

func generateUUID() string {
	return uuid.New().String()
}

const createVaultAmountTier = `--name: CreateAmountTier :one
INSERT INTO amount_tier (
	id,
	vault_category_id,
	min_amount,
	max_amount,
	interest,
	is_active,
	created_at,
	updated_at
) VALUES (:1, :2, :3, :4, :5, :6, :7, :8)
 RETURNING id INTO :result
`

func (r *Queries) SaveAmountTier(ctx context.Context, arg VaultAmountTier) (string, error) {
	var id string
	_, err := r.db.ExecContext(ctx, createVaultAmountTier,
		generateUUID(),
		arg.VaultCategoryID,
		arg.MinAmount,
		arg.MaxAmount,
		arg.Interest,
		utils.NullBoolToInt(arg.IsActive),
		time.Now(),
		time.Now(),
		sql.Out{Dest: &id})
	if err != nil {
		return "", err
	}
	return id, nil
}

const findAmountTier = `--name: FindAmountTier :many
SELECT
	id,
	vault_category_id,
	min_amount,
	max_amount,
	interest,
	is_active,
	created_at,
	updated_at,
	COUNT(*) OVER() AS total_count
FROM amount_tier
WHERE (:is_active IS NULL OR is_active = :is_active AND deleted_at IS NULL)
	AND (:vault_category_id IS NULL OR vault_category_id = :vault_category_id)
	AND (:min_amount IS NULL OR min_amount >= :min_amount)
	AND (:max_amount IS NULL OR max_amount <= :max_amount)
	AND (:is_active IS NULL OR is_active = :is_active)
ORDER BY created_at DESC
OFFSET NVL(:offset, 0) ROWS
FETCH NEXT NVL(:limit, 50) ROWS ONLY
`

func (r *Queries) FindAmountTier(ctx context.Context, arg VaultAmountTierParams) ([]VaultAmountTier, error) {
	limitPtr := utils.NullInt64ToPtr(arg.Limit)
	offsetPtr := (arg.Page.Int64 - 1) * arg.Limit.Int64
	minAmountPtr := utils.NullStringToPtr(arg.MinAmount)
	maxAmountPtr := utils.NullStringToPtr(arg.MaxAmount)

	rows, err := r.db.QueryContext(ctx, findAmountTier,
		sql.Named("vault_category_id", arg.VaultCategoryID),
		sql.Named("min_amount", minAmountPtr),
		sql.Named("max_amount", maxAmountPtr),
		sql.Named("is_active", arg.IsActive),
		sql.Named("limit", limitPtr),
		sql.Named("offset", offsetPtr),
	)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var (
		items         []VaultAmountTier
		minAmount_num godror.Number
		maxAmount_num godror.Number
		interest_num  godror.Number
	)

	for rows.Next() {
		var i VaultAmountTier
		if err := rows.Scan(
			&i.ID,
			&i.VaultCategoryID,
			&minAmount_num,
			&maxAmount_num,
			&interest_num,
			&i.IsActive,
			&i.CreatedAt,
			&i.UpdatedAt,
			&i.TotalCount,
		); err != nil {
			return nil, err
		}

		i.MinAmount, _ = decimal.NewFromString(minAmount_num.String())
		i.MaxAmount, _ = decimal.NewFromString(maxAmount_num.String())
		i.Interest, _ = decimal.NewFromString(interest_num.String())

		items = append(items, i)
	}

	if err := rows.Close(); err != nil {
		return nil, err
	}
	if err := rows.Err(); err != nil {
		return nil, err
	}
	return items, nil
}

const findAmountTierById = `--name: FindAmountTierById :one
SELECT
	id,
	vault_category_id,
	min_amount,
	max_amount,
	interest,
	is_active,
	created_at,
	updated_at
FROM amount_tier
WHERE id = :1 AND deleted_at IS NULL
`

func (r *Queries) FindAmountTierById(ctx context.Context, id string) (VaultAmountTier, error) {
	row := r.db.QueryRowContext(ctx, findAmountTierById, id)
	var i VaultAmountTier

	var (
		minAmount_num godror.Number
		maxAmount_num godror.Number
		interest_num  godror.Number
	)

	err := row.Scan(
		&i.ID,
		&i.VaultCategoryID,
		&minAmount_num,
		&maxAmount_num,
		&interest_num,
		&i.IsActive,
		&i.CreatedAt,
		&i.UpdatedAt,
	)
	if err != nil {
		return VaultAmountTier{}, err
	}

	i.MinAmount, _ = decimal.NewFromString(minAmount_num.String())
	i.MaxAmount, _ = decimal.NewFromString(maxAmount_num.String())
	i.Interest, _ = decimal.NewFromString(interest_num.String())

	return i, nil
}

const enableOrDisableAmountTier = `--name: EnableOrDisableAmountTier :one
UPDATE amount_tier
SET
	is_active = :is_active,
	updated_at = SYSTIMESTAMP
WHERE id = :id
`

func (r *Queries) EnableOrDisableAmountTier(ctx context.Context, id string, isEnabled bool) (string, error) {
	var active int
	if isEnabled {
		active = 1
	} else {
		active = 0
	}

	result, err := r.db.ExecContext(ctx, enableOrDisableAmountTier,
		sql.Named("id", id),
		sql.Named("is_active", active),
	)
	if err != nil {
		return "", err
	}
	rowsAffected, err := result.RowsAffected()
	if err != nil {
		return "", err
	}
	if rowsAffected == 0 {
		return "", fmt.Errorf("no rows updated")
	}
	return id, nil
}

const deleteAmountTier = `--name: DeleteAmountTier :one
UPDATE amount_tier
SET
	is_deleted = 1,
	deleted_at = SYSTIMESTAMP,
	updated_at = SYSTIMESTAMP
WHERE id = :id
`

func (r *Queries) DeleteAmountTier(ctx context.Context, id string) (string, error) {
	result, err := r.db.ExecContext(ctx, deleteAmountTier,
		sql.Named("id", id),
	)
	if err != nil {
		fmt.Println("============E", err)
		return "", err
	}
	rowsAffected, err := result.RowsAffected()
	if err != nil {
		return "", err
	}
	if rowsAffected == 0 {
		return "", fmt.Errorf("no rows updated")
	}
	return id, nil
}

const updateAmountTier = `--name: UpdateAmountTier :one
UPDATE amount_tier
SET
	min_amount = :min_amount,
	max_amount = :max_amount,
	interest = :interest,
	updated_at = SYSTIMESTAMP
WHERE id = :id
`

func (r *Queries) UpdateAmountTier(ctx context.Context, arg VaultAmountTier) (string, error) {
	result, err := r.db.ExecContext(ctx, updateAmountTier,
		sql.Named("id", arg.ID),
		sql.Named("min_amount", arg.MinAmount),
		sql.Named("max_amount", arg.MaxAmount),
		sql.Named("interest", arg.Interest),
	)
	if err != nil {
		return "", err
	}
	rowsAffected, err := result.RowsAffected()
	if err != nil {
		return "", err
	}
	if rowsAffected == 0 {
		return "", fmt.Errorf("no rows updated")
	}
	return arg.ID, nil
}
