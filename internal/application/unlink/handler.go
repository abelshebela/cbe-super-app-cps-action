package unlink

import (
	"context"

	local_util "github.com/CBE-Super-App/cbe-super-app-cps-action/pkgs/utils"
	"gitlab.com/bersufekadgetachew/cbe-super-app-shared/shared/utils"
)

type unlinkApplication struct {
	logger              utils.Logger
	unlinkServiceDomain UnlinkAccount
}

func NewUnlinkApplication(domain UnlinkAccount, logger utils.Logger) UnlinkAccount {

	return &unlinkApplication{
		logger:              logger,
		unlinkServiceDomain: domain,
	}
}

func (ua *unlinkApplication) GetUserByAccount(ctx context.Context, accNumber string) (*local_util.PaginatedResponse[*any], error) {
	return ua.unlinkServiceDomain.GetUserByAccount(ctx, accNumber)
}
func (ua *unlinkApplication) UnlinkUserCif(ctx context.Context, userCode string) error {
	return ua.UnlinkUserCif(ctx, userCode)
}
func (ua *unlinkApplication) Authorize(ctx context.Context, cpsAction any) error {
	return ua.Authorize(ctx, cpsAction)
}
