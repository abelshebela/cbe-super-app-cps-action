package hq

import (
	"context"
)

type Repository interface {
	GetHQByID(ctx context.Context, id string) (HQ, error)
	UpdateHQ(ctx context.Context, id string, update HQ) error
}
