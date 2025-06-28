package avatar

import (
	"context"

	"gitlab.com/bersufekadgetachew/cbe-super-app-cps-action/internal/adapter/outbound/model"
	constant "gitlab.com/bersufekadgetachew/cbe-super-app-cps-action/utils"
)

type AvatarRepository interface {
	CreateAvatar(ctx context.Context, req model.CreateCPSAction) (model.CpsAction, error)
	CPSActionExists(ctx context.Context, cpsReq model.CreateCPSAction) error
	DeleteAvatar(ctx context.Context, id string, cpsActionReq model.CreateCPSAction) (model.CpsAction, error)
	Authorize(ctx context.Context, req model.AuthorizeCPSAction) (model.CpsAction, error)
	Reject(ctx context.Context, req model.RejectCPSAction) (model.CpsAction, error)
	EnableOrDisableAvatar(ctx context.Context, id string,
		requestAction model.RequestAction, cpsReq model.CreateCPSAction) (*model.CpsAction, error)
	GetAvatar(ctx context.Context, id string) (*Avatar, error)
	GetAllAvatar(ctx context.Context, filterParams constant.Filter) (AvatarResponse, error)
	UpdateAvatar(ctx context.Context, id string, req model.CreateCPSAction) (model.CpsAction, error)
}
