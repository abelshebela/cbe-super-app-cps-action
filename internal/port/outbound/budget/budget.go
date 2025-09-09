package ad

import (
	"context"

	"github.com/CBE-Super-App/cbe-super-app-cps-action/internal/domain/budget/entities"
)

type Repo interface {
	CreateIconAction(ctx context.Context, cpsAction entities.CPSAction) (*entities.CPSAction, error)
	FetchIcons(ctx context.Context) ([]*entities.Icon, error)
	UpdateIcon(ctx context.Context, id string, cpsAction entities.CPSAction) (*entities.CPSAction, error)
	CreateColor(ctx context.Context, color string, action entities.CPSAction) (*entities.CPSAction, error)
	CheckColorExist(ctx context.Context, color string) (bool, error)
	ListAllColor(ctx context.Context) ([]*entities.Color, error)
	CreateColorUpdateAction(ctx context.Context, id string, cpsAction entities.CPSAction) (*entities.CPSAction, error)
	ApproveAction(ctx context.Context, action entities.CPSAction) (*entities.CPSAction, error)
	Authorize(ctx context.Context, cpsAction *entities.CPSAction) (*entities.CPSAction, error)

	// CreateAction(ctx context.Context, action entities.CPSAction) (*entities.CPSAction, error)
	// GetByIDColor(ctx context.Context, id string) (*entities.Color, error)
	// UpdateColor(ctx context.Context, color entities.CPSAction) (*entities.CPSAction, error)
}
