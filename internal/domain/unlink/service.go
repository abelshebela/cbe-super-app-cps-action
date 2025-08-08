package unlink

import (
	"context"

	"gitlab.com/bersufekadgetachew/cbe-super-app-shared/shared/utils"
)

type unlinkService struct {
	logger            utils.Logger
	unlinkPersistence UnlinkAccount
}

func NewUnlinkServiceDomain(persistence UnlinkAccount, logger utils.Logger) UnlinkAccount {

	return &unlinkService{
		logger:            logger,
		unlinkPersistence: persistence,
	}
}

func (u *unlinkService) GetUserByAccount(ctx context.Context, accNumber string) (*any, error) {
	return u.unlinkPersistence.GetUserByAccount(ctx, accNumber)
}
func (u *unlinkService) UnlinkUserCif(ctx context.Context, userCode string) error {
	return u.unlinkPersistence.UnlinkUserCif(ctx, userCode)
}
func (u *unlinkService) Authorize(ctx context.Context, cpsAction any) (any, error) {
	return u.unlinkPersistence.Authorize(ctx, cpsAction)
}
