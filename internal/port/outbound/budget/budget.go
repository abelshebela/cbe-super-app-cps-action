package ad

import (
	"context"

	"gitlab.com/bersufekadgetachew/cbe-super-app-cps-action/internal/domain/budget/entities"
)

type Repo interface {
	CreateIconAction(ctx context.Context, cpsAction entities.CPSAction) (*entities.CPSAction, error)
	FetchIcons(ctx context.Context) ([]*entities.Icon, error)
	UpdateIcon(ctx context.Context, id string, cpsAction entities.CPSAction) (*entities.CPSAction, error)

	CreateAction(ctx context.Context, action entities.CPSAction) (*entities.CPSAction, error)
	CreateColor(ctx context.Context, color string, action entities.CPSAction) (*entities.CPSAction, error)
	ListAllColor(ctx context.Context) ([]*entities.Color, error)
	GetByIDColor(ctx context.Context, id string) (*entities.Color, error)
	UpdateColor(ctx context.Context, color entities.CPSAction) (*entities.CPSAction, error)
	ApproveAction(ctx context.Context, action entities.CPSAction) (*entities.CPSAction, error)
}
