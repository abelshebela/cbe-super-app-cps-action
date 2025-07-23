package avatar

import (
	"context"

	entities "github.com/CBE-Super-App/cbe-super-app-cps-action/internal/domain/cps_actions/entities"
	common_util "github.com/CBE-Super-App/cbe-super-app-cps-action/pkgs/utils"
	constant "github.com/CBE-Super-App/cbe-super-app-cps-action/utils"
)

type AvatarRepository interface {
	AuthorizeCreateAvatar(ctx context.Context, cpsAction *entities.CPSAction) (*entities.CPSAction, error)
	AuthorizeUpdateAvatar(ctx context.Context, cpsAction *entities.CPSAction) (*entities.CPSAction, error)
	AuthorizeDeleteAvatar(ctx context.Context, cpsAction *entities.CPSAction) (*entities.CPSAction, error)
	GetAvatar(ctx context.Context, id string) (*Avatar, error)
	GetAllAvatar(ctx context.Context, filterParams constant.Filter) (*common_util.PaginatedResponse[[]*Avatar], error)
}
