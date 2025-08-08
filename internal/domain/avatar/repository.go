package avatar

import (
	"context"

	common_util "github.com/CBE-Super-App/cbe-super-app-cps-action/pkgs/utils"
	constant "github.com/CBE-Super-App/cbe-super-app-cps-action/utils"
)

type AvatarRepository interface {
	CreateAvatar(ctx context.Context, avatar Avatar) (*Avatar, error)
	UpdateAvatar(ctx context.Context, avatar Avatar) (*Avatar, error)
	DeleteAvatar(ctx context.Context, id string) (*Avatar, error)
	EnableDisableAvatar(ctx context.Context, id string, enable bool) (*Avatar, error)
	GetAvatar(ctx context.Context, id string) (*Avatar, error)
	GetAllAvatar(ctx context.Context, filterParams *constant.Filter) (*common_util.PaginatedResponse[[]*Avatar], error)
}
