package avatar

import (
	"context"

	dto "github.com/CBE-Super-App/cbe-super-app-cps-action/internal/domain/avatar"
	constant "github.com/CBE-Super-App/cbe-super-app-cps-action/utils"

	common_util "github.com/CBE-Super-App/cbe-super-app-cps-action/pkgs/utils"
)

type AvatarRepository interface {
	CreateAvatar(ctx context.Context, avatar dto.Avatar) (*dto.Avatar, error)
	UpdateAvatar(ctx context.Context, avatar dto.Avatar) (*dto.Avatar, error)
	DeleteAvatar(ctx context.Context, id string) (*dto.Avatar, error)
	EnableDisableAvatar(ctx context.Context, id string, enable bool) (*dto.Avatar, error)
	GetAvatar(ctx context.Context, id string) (*dto.Avatar, error)
	GetAllAvatar(ctx context.Context, filterParams *constant.Filter) (*common_util.PaginatedResponse[[]*dto.Avatar], error)
}
