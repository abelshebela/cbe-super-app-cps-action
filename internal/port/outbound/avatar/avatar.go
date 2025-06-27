package avatar

import (
	"context"

	"gitlab.com/bersufekadgetachew/cbe-super-app-cps-action/internal/adapter/outbound/model"
)

type AvatarOutbound interface {
	CreateAvatar(ctx context.Context, cpsActionReq model.CreateCPSAction) (model.CpsAction, error)
	CPSActionExists(ctx context.Context, cpsReq model.CreateCPSAction) error
}
