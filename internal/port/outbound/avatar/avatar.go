package avatar

import (
	"context"

	"github.com/CBE-Super-App/cbe-super-app-cps-action/internal/adapter/outbound/model"
	dto "github.com/CBE-Super-App/cbe-super-app-cps-action/internal/domain/avatar"
	constant "github.com/CBE-Super-App/cbe-super-app-cps-action/utils"
)

type AvatarOutbound interface {
	CreateAvatar(ctx context.Context, cpsActionReq model.CreateCPSAction) (model.CPSAction, error)
	CPSActionExists(ctx context.Context, cpsReq model.CreateCPSAction) error
	DeleteAvatar(ctx context.Context, id string, cpsActionReq model.CreateCPSAction) (model.CPSAction, error)
	Authorize(ctx context.Context, req model.AuthorizeCPSAction) (model.CPSAction, error)
	Reject(ctx context.Context, req model.RejectCPSAction) (model.CPSAction, error)
	EnableOrDisableAvatar(ctx context.Context, id string,
		requestAction model.RequestAction, cpsReq model.CreateCPSAction) (*model.CPSAction, error)
	GetAvatar(ctx context.Context, id string) (*dto.Avatar, error)
	GetAllAvatar(ctx context.Context, filterParams constant.Filter) (dto.AvatarResponse, error)
	UpdateAvatar(ctx context.Context, id string, cpsActionReq model.CreateCPSAction) (model.CPSAction, error)
}
