package unlink

import (
	"context"

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

func (ua *unlinkApplication) GetUserByAccount(ctx context.Context, accNumber string) (*any, error) {
	return ua.unlinkServiceDomain.GetUserByAccount(ctx, accNumber)
}
func (ua *unlinkApplication) UnlinkUserCif(ctx context.Context, userCode string) error {
	return ua.unlinkServiceDomain.UnlinkUserCif(ctx, userCode)
}
func (ua *unlinkApplication) Authorize(ctx context.Context, cpsAction any) (any, error) {
	return ua.unlinkServiceDomain.Authorize(ctx, cpsAction)
}
