package avatar

import (
	"context"

	"github.com/CBE-Super-App/cbe-super-app-cps-action/internal/adapter/outbound/model"
	dto "github.com/CBE-Super-App/cbe-super-app-cps-action/internal/domain/avatar"
	constant "github.com/CBE-Super-App/cbe-super-app-cps-action/utils"
	common_util "github.com/CBE-Super-App/cbe-super-app-cps-action/pkgs/utils"
)

type AvatarOutbound interface {
	CreateAvatar(ctx context.Context, cpsActionReq model.CreateCPSAction) (*dto.CPSAction, error)
	CPSActionExists(ctx context.Context, cpsReq model.CreateCPSAction) error
	DeleteAvatar(ctx context.Context, id string, cpsActionReq model.CreateCPSAction) (*dto.CPSAction, error)
	Authorize(ctx context.Context, req model.AuthorizeCPSAction) (*dto.CPSAction, error)
	Reject(ctx context.Context, req model.RejectCPSAction) (*dto.CPSAction, error)
	EnableOrDisableAvatar(ctx context.Context, id string,
		requestAction model.RequestAction, cpsReq model.CreateCPSAction) (*dto.CPSAction, error)
	GetAvatar(ctx context.Context, id string) (*dto.Avatar, error)
	GetAllAvatar(ctx context.Context, filterParams constant.Filter) (*common_util.PaginatedResponse[[]*dto.Avatar], error)
	UpdateAvatar(ctx context.Context, id string, cpsActionReq model.CreateCPSAction) (*dto.CPSAction, error)
}
