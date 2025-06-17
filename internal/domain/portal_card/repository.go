package portalcard

import (
	"context"
)

type PortaCardInterface interface {
	GetAllPortalCard(ctx context.Context) ([]*Card, error)
}

type Repository interface {
	GetAllPortalCard(ctx context.Context) ([]*Card, error)
}
