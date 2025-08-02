package unlink

import (
	"context"

	local_util "github.com/CBE-Super-App/cbe-super-app-cps-action/pkgs/utils"
)

// type Repository interface {
// 	UnlinkDevice(ctx context.Context, userCode string, cpsAction entities.CPSAction) (string, error)
// 	ApproveOrDecline(ctx context.Context, userCode, decision, reason string, cpsAction entities.CPSAction) error
// }

type UnlinkAccount interface {
	GetUserByAccount(ctx context.Context, accNumber string) (*local_util.PaginatedResponse[*any], error)
	UnlinkUserCif(ctx context.Context, userCode string) error
	Authorize(ctx context.Context, cpsAction any) error
}
