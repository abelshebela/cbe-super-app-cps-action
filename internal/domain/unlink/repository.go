package unlink

import (
	"context"

	local_util "github.com/CBE-Super-App/cbe-super-app-cps-action/pkgs/utils"
)

type UnlinkAccount interface {
	GetUserByAccount(ctx context.Context, accNumber string) (*any, error)
	GetAllArchivedUser(ctx context.Context, accNumber string) (*local_util.PaginatedResponse[*any], error)
	UnlinkUserCif(ctx context.Context, userCode string) error
	Authorize(ctx context.Context, cpsAction any) (any, error)
}
