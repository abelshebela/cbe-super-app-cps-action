package sqlc

import (
	"time"

	"github.com/jackc/pgx/v5/pgtype"
)

type GroupVaultCategory struct {
	ID          string      `json:"id"`
	Name        string      `json:"name"`
	Description pgtype.Text `json:"description"`
	IsActive    bool        `json:"is_active"`
	CreatedAt   time.Time   `json:"created_at"`
	UpdatedAt   time.Time   `json:"updated_at"`
	DeletedAt   *time.Time  `json:"deleted_at"`
	IsDeleted   bool        `json:"is_deleted"`
}
