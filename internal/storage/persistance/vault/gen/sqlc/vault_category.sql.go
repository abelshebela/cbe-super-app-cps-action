package sqlc

import (
	"context"
	"database/sql"
	"fmt"
	"strings"
	"time"

	imodel "cbe-super-app-cps-action/internal/constants/model"
	"cbe-super-app-cps-action/pkgs/utils"

	"github.com/google/uuid"
)

func generateUUID() string {
	return uuid.New().String()
}

func boolToInt(b bool) int {
	if b {
		return 1
	}
	return 0
}

func pointerBoolToInt(pb *bool) *int {
	if pb == nil {
		return nil
	}
	v := 0
	if *pb {
		v = 1
	}
	return &v
}

const activateVaultCategory = `-- name: ActivateVaultCategory :one
UPDATE vault_categories
SET is_active = 1, updated_at = SYSTIMESTAMP
WHERE id = :1`

func (q *Queries) ActivateVaultCategory(ctx context.Context, id string) (string, error) {
	var result string
	res, err := q.db.ExecContext(ctx, activateVaultCategory, id)
	if err != nil {
		return "", fmt.Errorf("failed to activate vault  category: %w", err)
	}
	rows, _ := res.RowsAffected()
	if rows == 0 {
		return "", fmt.Errorf("no vault  category found with id %s", id)
	}
	return result, nil
}

const deactivateVaultCategory = `-- name: DeactivateVaultCategory :one
UPDATE vault_categories
SET is_active = 0, updated_at = SYSTIMESTAMP
WHERE id = :1`

func (q *Queries) DeactivateVaultCategory(ctx context.Context, id string) (string, error) {
	var result string
	res, err := q.db.ExecContext(ctx, deactivateVaultCategory, id)
	if err != nil {
		return "", fmt.Errorf("failed to deactivate vault  category: %w", err)
	}
	rows, _ := res.RowsAffected()
	if rows == 0 {
		return "", fmt.Errorf("no vault  category found with id %s", id)
	}
	return result, nil
}

const deleteVaultCategory = `-- name: DeleteVaultCategory :one
UPDATE vault_categories
SET is_deleted = 1, deleted_at = SYSTIMESTAMP, updated_at = SYSTIMESTAMP
WHERE id = :1
RETURNING id INTO :result`

func (q *Queries) DeleteVaultCategory(ctx context.Context, id string) (string, error) {
	var result string
	res, err := q.db.ExecContext(ctx, deleteVaultCategory, id, sql.Out{Dest: &result})
	if err != nil {
		return "", fmt.Errorf("failed to delete vault  category: %w", err)
	}
	rows, _ := res.RowsAffected()
	if rows == 0 {
		return "", fmt.Errorf("no vault  category found with id %s", id)
	}
	return result, nil
}

const findVaultCategories = `-- name: FindVaultCategory :many
SELECT
	id,
	COUNT(*) OVER() AS total_count
FROM vault_categories
WHERE (:is_active IS NULL OR is_active = :is_active)
	AND (:name IS NULL OR UPPER(name) LIKE :name)
ORDER BY created_at DESC
OFFSET NVL(:offset, 0) ROWS
FETCH NEXT NVL(:limit, 50) ROWS ONLY`

type FindVaultCategoryParams struct {
	IsActive  sql.NullBool   `json:"is_active"`
	NameQuery sql.NullString `json:"name"`
	Page      sql.NullInt64  `json:"page"`
	Limit     sql.NullInt64  `json:"limit"`
}

func (q *Queries) FindVaultCategories(ctx context.Context, arg FindVaultCategoryParams) ([]VaultCategory, error) {

	namePtr := utils.NullStringToPtrLike(arg.NameQuery)
	limitPtr := utils.NullInt64ToPtr(arg.Limit)
	var offsetPtr *int64
	if arg.Page.Valid && arg.Limit.Valid {
		off := (arg.Page.Int64 - 1) * arg.Limit.Int64
		offsetPtr = &off
	}

	rows, err := q.db.QueryContext(ctx, findVaultCategories,
		sql.Named("is_active", arg.IsActive),
		sql.Named("name", namePtr),
		sql.Named("offset", offsetPtr),
		sql.Named("limit", limitPtr),
	)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	items := []VaultCategory{}

	for rows.Next() {
		var id string
		var totalCount int64
		if err := rows.Scan(&id, &totalCount); err != nil {
			return nil, err
		}

		category, tiers, err := q.FindVaultCategoryWithTiers(ctx, id)
		if err != nil {
			return nil, err
		}

		vc := VaultCategory{
			ID:               category.ID,
			Name:             category.Name,
			CoverImageURL:    category.CoverImageURL,
			InterestType:     category.InterestType,
			CategoryInterest: category.CategoryInterest,
			Deadlock:         category.Deadlock,
			IsActive:         category.IsActive,
			CreatedAt:        category.CreatedAt,
			UpdatedAt:        category.UpdatedAt,
			TotalCount:       totalCount,
			Tiers:            tiers,
		}

		items = append(items, vc)
	}

	if err := rows.Err(); err != nil {
		return nil, err
	}

	return items, nil
}

const findVaultCategoryWithTiers = `-- name: FindVaultCategoryWithTiers :many
SELECT
  c.id,
  c.name,
  c.cover_image_url,
  c.interest_type,
  c.category_interest,
  c.deadlock,
  c.is_active,
  c.created_at,
  c.updated_at,

  t.id,
  t.name,
  t.tier_interest,
  t.min_amount,
  t.max_amount
FROM vault_categories c
LEFT JOIN vault_tiers t
  ON c.id = t.category_id
WHERE c.id = :1`

func (q *Queries) FindVaultCategoryWithTiers(ctx context.Context, id string) (*imodel.VaultCategory, []imodel.VaultTiers, error) {

	rows, err := q.db.QueryContext(ctx, findVaultCategoryWithTiers, id)
	if err != nil {
		return nil, nil, err
	}
	defer rows.Close()

	var category imodel.VaultCategory
	var tiers []imodel.VaultTiers
	found := false

	for rows.Next() {

		var t imodel.VaultTiers

		err := rows.Scan(
			&category.ID,
			&category.Name,
			&category.CoverImageURL,
			&category.InterestType,
			&category.CategoryInterest,
			&category.Deadlock,
			&category.IsActive,
			&category.CreatedAt,
			&category.UpdatedAt,

			&t.ID,
			&t.Name,
			&t.TierInterest,
			&t.MinAmount,
			&t.MaxAmount,
		)

		if err != nil {
			return nil, nil, err
		}

		// mark that we have at least one row (category exists)
		found = true

		if t.ID != "" {
			tiers = append(tiers, t)
		}
	}

	if err := rows.Err(); err != nil {
		return nil, nil, err
	}

	if !found {
		return nil, nil, sql.ErrNoRows
	}

	return &category, tiers, nil
}

const findVaultCategoryByName = `-- name: FindVaultCategoryByName :one
SELECT
  id,
  name,
  cover_image_url,
  interest_type,
  category_interest,
  deadlock,
  is_active,
  created_at,
  updated_at
FROM vault_categories
WHERE UPPER(name) = UPPER(:name)`

func (q *Queries) FindVaultCategoryByName(ctx context.Context, name string) (*imodel.VaultCategory, error) {
	row := q.db.QueryRowContext(ctx, findVaultCategoryByName, sql.Named("name", name))
	var i imodel.VaultCategory
	err := row.Scan(
		&i.ID,
		&i.Name,
		&i.CoverImageURL,
		&i.InterestType,
		&i.CategoryInterest,
		&i.Deadlock,
		&i.IsActive,
		&i.CreatedAt,
		&i.UpdatedAt,
	)
	return &i, err
}

const saveVaultCategory = `-- name: SaveVaultCategory :one
INSERT INTO vault_categories (
    id,
    name,
    cover_image_url,
    interest_type,
    category_interest,
    deadlock,
    is_active,
    created_at,
    updated_at
) VALUES (
    :1, :2, :3, :4, :5, :6, :7, SYSTIMESTAMP, SYSTIMESTAMP
)`

const saveVaultTier = `-- name: SaveVaultTier :exec
INSERT INTO vault_tiers (
    id,
    category_id,
    name,
    tier_interest,
    min_amount,
    max_amount
) VALUES (
    :1, :2, :3, :4, :5, :6
)`

func (q *Queries) SaveVaultCategory(ctx context.Context, arg *imodel.VaultCategory) (string, error) {
	categoryID := generateUUID()

	_, err := q.db.ExecContext(
		ctx,
		saveVaultCategory,
		categoryID,
		strings.ToUpper(arg.Name),
		arg.CoverImageURL,
		arg.InterestType,
		arg.CategoryInterest,
		boolToInt(arg.Deadlock),
		boolToInt(arg.IsActive),
	)
	if err != nil {
		return "", fmt.Errorf("insert category failed: %w", err)
	}

	return categoryID, nil
}

func (q *Queries) SaveVaultTiers(ctx context.Context, categoryID string, arg *imodel.VaultCategory) error {
	for _, tier := range arg.Tiers {
		tierID := generateUUID()

		_, err := q.db.ExecContext(
			ctx,
			saveVaultTier,
			tierID,
			categoryID,
			tier.Name,
			tier.TierInterest,
			tier.MinAmount,
			tier.MaxAmount,
		)
		if err != nil {
			return fmt.Errorf("insert tier failed: %w", err)
		}
	}
	return nil
}

const updateVaultCategory = `-- name: UpdateVaultCategory :one
UPDATE vault_categories
SET
    name              = COALESCE(UPPER(:1), name),
    cover_image_url   = COALESCE(:2, cover_image_url),
    interest_type     = COALESCE(:3, interest_type),
    category_interest = COALESCE(:4, category_interest),
    deadlock          = COALESCE(:5, deadlock),
    updated_at        = COALESCE(:8, updated_at)
WHERE id = :9`

const updateVaultTier = `-- name: UpdateVaultTier :exec
UPDATE vault_tiers
SET
    name          = COALESCE(:name, name),
    tier_interest = COALESCE(:tier_interest, tier_interest),
    min_amount    = COALESCE(:min_amount, min_amount),
    max_amount    = COALESCE(:max_amount, max_amount)
WHERE id = :id AND category_id = :category_id`

func (q *Queries) UpdateVaultCategory(ctx context.Context, catID string, arg *imodel.VaultCategory) (string, error) {
	var name *string
	if arg.Name != "" {
		v := strings.ToUpper(arg.Name)
		name = &v
	}

	var cover *string
	if arg.CoverImageURL != "" {
		cover = &arg.CoverImageURL
	}

	var interestType *string
	if arg.InterestType != "" {
		interestType = &arg.InterestType
	}

	var categoryInterest *string
	if arg.CategoryInterest != "" {
		categoryInterest = &arg.CategoryInterest
	}

	var deadlock *int
	if arg.Deadlock {
		v := 1
		deadlock = &v
	}

	now := time.Now()

	res, err := q.db.ExecContext(
		ctx,
		updateVaultCategory,

		name,             // :1
		cover,            // :2
		interestType,     // :3
		categoryInterest, // :4
		deadlock,         // :5
		now,              // :8
		catID,            // :9
	)

	if err != nil {
		return "", fmt.Errorf("failed to update vault category: %w", err)
	}

	rows, err := res.RowsAffected()
	if err != nil {
		return "", err
	}

	if rows == 0 {
		return "", fmt.Errorf("vault category not found or inactive")
	}

	return arg.ID, nil
}

func (q *Queries) UpdateVaultTiers(ctx context.Context, categoryID string, arg *imodel.VaultCategory) error {
	for _, tier := range arg.Tiers {
		fmt.Printf("[DEBUG] UpdateVaultTiers: attempting update - category_id=%s tier_id=%s name=%s tier_interest=%s min=%s max=%s\n", categoryID, tier.ID, tier.Name, tier.TierInterest, tier.MinAmount, tier.MaxAmount)
		res, err := q.db.ExecContext(
			ctx,
			updateVaultTier,

			sql.Named("id", tier.ID),
			sql.Named("category_id", categoryID),

			sql.Named("name", tier.Name),
			sql.Named("tier_interest", tier.TierInterest),
			sql.Named("min_amount", tier.MinAmount),
			sql.Named("max_amount", tier.MaxAmount),
		)

		if err != nil {
			return fmt.Errorf("update tier failed: %w", err)
		}

		rows, rerr := res.RowsAffected()
		if rerr != nil {
			return fmt.Errorf("could not determine rows affected when updating tier %s: %w", tier.ID, rerr)
		}

		fmt.Printf("[DEBUG] UpdateVaultTiers: rows affected=%d for tier_id=%s\n", rows, tier.ID)

		if rows == 0 {
			return fmt.Errorf("no tier updated for id %s", tier.ID)
		}
	}

	return nil
}

const findVaultTransactions = `-- name: FindVaultCategory :many
SELECT
	id,
	COUNT(*) OVER() AS total_count
FROM vault_categories
WHERE (:is_active IS NULL OR is_active = :is_active)
	AND (:name IS NULL OR UPPER(name) LIKE :name)
ORDER BY created_at DESC
OFFSET NVL(:offset, 0) ROWS
FETCH NEXT NVL(:limit, 50) ROWS ONLY`

func (q *Queries) FindVaultTransactions(ctx context.Context, arg FindVaultCategoryParams) ([]VaultTransaction, error) {
	limitPtr := utils.NullInt64ToPtr(arg.Limit)
	var offsetPtr *int64
	if arg.Page.Valid && arg.Limit.Valid {
		off := (arg.Page.Int64 - 1) * arg.Limit.Int64
		offsetPtr = &off
	}

	rows, err := q.db.QueryContext(ctx, findVaultTransactions,
		sql.Named("offset", offsetPtr),
		sql.Named("limit", limitPtr),
	)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	items := []VaultTransaction{}

	for rows.Next() {
		var id string
		var totalCount int64
		if err := rows.Scan(&id, &totalCount); err != nil {
			return nil, err
		}

		vc := VaultTransaction{}

		items = append(items, vc)
	}

	if err := rows.Err(); err != nil {
		return nil, err
	}

	return items, nil
}
