package storage

import (
	"context"
)

type RoleRepository interface {
	Exists(ctx context.Context, id string) (bool, error)
	ExistsMany(ctx context.Context, ids []string) (bool, error)
}
