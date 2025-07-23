package avatar

import (
	"context"

	dto "github.com/CBE-Super-App/cbe-super-app-cps-action/internal/domain/avatar"
	entities "github.com/CBE-Super-App/cbe-super-app-cps-action/internal/domain/cps_actions/entities"
	constant "github.com/CBE-Super-App/cbe-super-app-cps-action/utils"

	common_util "github.com/CBE-Super-App/cbe-super-app-cps-action/pkgs/utils"
)

type AvatarOutbound interface {
	GetAvatar(ctx context.Context, id string) (*dto.Avatar, error)
	GetAllAvatar(ctx context.Context, filterParams constant.Filter) (*common_util.PaginatedResponse[[]*dto.Avatar], error)
	AuthorizeCreateAvatar(ctx context.Context, cpsAction *entities.CPSAction) (*entities.CPSAction, error)
	AuthorizeUpdateAvatar(ctx context.Context, cpsAction *entities.CPSAction) (*entities.CPSAction, error)
	AuthorizeDeleteAvatar(ctx context.Context, cpsAction *entities.CPSAction) (*entities.CPSAction, error)
}
