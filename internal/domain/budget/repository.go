package budget

import (
	"context"

	"cbe-super-app-cps-action/internal/domain/budget/entities"
)

// type IconRepository interface {
//     ListAll() ([]entities.Icon, error)
//     GetByID(id string) (*entities.Icon, error)
//     Create(icon *entities.Icon) error
//     Update(icon *entities.Icon) error
// }

// type ColorRepository interface {
//     ListAllColor() ([]entities.Color, error)
//     GetByIDColor(id string) (*entities.Color, error)
//     CreateColor(color *entities.Color) error
//     UpdateColor(color *entities.Color) error
// }

// type CPSRepository interface {
//     CreateAction(action *entities.CPSAction) error
//     GetByCode(code string) (*entities.CPSAction, error)
//     Complete(code string) error
// }

type Repository interface {
	CreateIconAction(ctx context.Context, action entities.CPSAction) (*entities.CPSAction, error)
	FetchIcons(ctx context.Context) ([]*entities.Icon, error)
	UpdateIcon(ctx context.Context, id string, cpsAction entities.CPSAction) (*entities.CPSAction, error)

	CreateAction(ctx context.Context, action entities.CPSAction) (*entities.CPSAction, error)
	CreateColor(ctx context.Context, color string, action entities.CPSAction) (*entities.CPSAction, error)
	ListAllColor(ctx context.Context) ([]*entities.Color, error)
	GetByIDColor(ctx context.Context, id string) (*entities.Color, error)
	UpdateColor(ctx context.Context, color entities.CPSAction) (*entities.CPSAction, error)
	ApproveAction(ctx context.Context, action entities.CPSAction) (*entities.CPSAction, error)
}
