package portalcard

import (
	"context"
)

type PortaCardInterface interface {
	GetAllPortalCard(ctx context.Context) ([]*Card, error)
}

type PortalCardRepository interface {
	GetAllPortalCard(ctx context.Context) ([]*Card, error)
}
