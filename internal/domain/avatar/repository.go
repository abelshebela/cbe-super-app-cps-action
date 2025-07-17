package avatar

import (
	"context"

	"github.com/CBE-Super-App/cbe-super-app-cps-action/internal/adapter/outbound/model"
	constant "github.com/CBE-Super-App/cbe-super-app-cps-action/utils"
)

type AvatarRepository interface {
	CreateAvatar(ctx context.Context, req model.CreateCPSAction) (*CPSAction, error)
	CPSActionExists(ctx context.Context, cpsReq model.CreateCPSAction) error
	DeleteAvatar(ctx context.Context, id string, cpsActionReq model.CreateCPSAction) (*CPSAction, error)
	Authorize(ctx context.Context, req model.AuthorizeCPSAction) (*CPSAction, error)
	Reject(ctx context.Context, req model.RejectCPSAction) (*CPSAction, error)
	EnableOrDisableAvatar(ctx context.Context, id string,
		requestAction model.RequestAction, cpsReq model.CreateCPSAction) (*CPSAction, error)
	GetAvatar(ctx context.Context, id string) (*Avatar, error)
	GetAllAvatar(ctx context.Context, filterParams constant.Filter) (*AvatarResponse, error)
	UpdateAvatar(ctx context.Context, id string, req model.CreateCPSAction) (*CPSAction, error)
}
