-- name: SaveVaultCategory :one
INSERT INTO vault_categories (
  id,
  name,
  category_type,
  cover_image,
  interest_type,
  interest,
  deadlock,
  is_active
) VALUES (
  $1,
  $2,
  $3,
  $4,
  $5,
  $6,
  $7,
  $8
)
RETURNING id;

-- name: UpdateVaultCategory :one
UPDATE vault_categories
SET
  name          = COALESCE($1, name),
  cover_image   = COALESCE($2, cover_image),
  interest_type = COALESCE($3, interest_type),
  interest      = COALESCE($4, interest),
  deadlock      = COALESCE($5, deadlock),
  updated_at    = NOW()
WHERE id = $6 AND is_deleted = FALSE
RETURNING id;

-- name: ActivateVaultCategory :one
UPDATE vault_categories
SET is_active = TRUE, updated_at = NOW()
WHERE id = $1 AND is_deleted = FALSE
RETURNING id;

-- name: DeactivateVaultCategory :one
UPDATE vault_categories
SET is_active = FALSE, updated_at = NOW()
WHERE id = $1 AND is_deleted = FALSE
RETURNING id;

-- name: DeleteVaultCategory :one
UPDATE vault_categories
SET is_deleted = TRUE, deleted_at = NOW(), updated_at = NOW()
WHERE id = $1 AND is_deleted = FALSE
RETURNING id;

-- name: FindVaultCategory :many
SELECT
  id,
  name,
  category_type,
  cover_image,
  interest_type,
  interest,
  deadlock,
  is_active,
  is_deleted,
  created_at,
  updated_at,
  deleted_at,
  COUNT(*) OVER() AS total_count
FROM vault_categories
WHERE ($1::boolean IS NULL OR is_active = $1::boolean)
  AND ($2::text IS NULL OR name ILIKE ('%' || $2::text || '%'))
ORDER BY created_at DESC
LIMIT  COALESCE($4::int, 50)
OFFSET COALESCE($3::int, 0);

-- name: FindVaultCategoryById :one
SELECT
  id,
  name,
  category_type,
  cover_image,
  interest_type,
  interest,
  deadlock,
  is_active,
  is_deleted,
  created_at,
  updated_at,
  deleted_at
FROM vault_categories
WHERE id = $1;

-- name: FindVaultCategoryByName :one
SELECT id, name, category_type, cover_image, interest_type, interest, deadlock, is_active, is_deleted, created_at, updated_at, deleted_at
FROM vault_categories
WHERE name = $1;



