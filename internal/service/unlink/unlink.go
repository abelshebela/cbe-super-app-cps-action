package unlink

import (
	"cbe-super-app-cps-action/internal/constants"
	"cbe-super-app-cps-action/internal/constants/lib"
	"cbe-super-app-cps-action/internal/constants/localization"
	"cbe-super-app-cps-action/internal/constants/model"
	"cbe-super-app-cps-action/internal/constants/types"
	"cbe-super-app-cps-action/internal/service"
	"cbe-super-app-cps-action/internal/storage"
	local_util "cbe-super-app-cps-action/pkgs/utils"
	"context"
	"errors"

	"gitlab.com/bersufekadgetachew/cbe-super-app-shared/shared/utils"
	"go.mongodb.org/mongo-driver/v2/mongo"
)

type unlinkService struct {
	logger                    utils.Logger
	repo                      storage.UnlinkAccount
	userRepo                  storage.UserRepository
	archivedUserRepo          storage.ArchivedUserRepository
	linkedAccountRepo         storage.LinkedAccountRepository
	archivedLinkedAccountRepo storage.ArchivedLinkedAccountRepository
	cpsService                service.CPSActionService
}

func NewUnlinkService(client *mongo.Client, userData storage.UserRepository, archivedUserRepo storage.ArchivedUserRepository, linkedAccountRepo storage.LinkedAccountRepository, archivedLinkedAccountRepo storage.ArchivedLinkedAccountRepository, repo storage.UnlinkAccount, cpsAction service.CPSActionService, logger utils.Logger) service.UnlinkService {
	return &unlinkService{
		repo:                      repo,
		userRepo:                  userData,
		archivedUserRepo:          archivedUserRepo,
		linkedAccountRepo:         linkedAccountRepo,
		archivedLinkedAccountRepo: archivedLinkedAccountRepo,
		cpsService:                cpsAction,
		logger:                    logger,
	}
}

func (u *unlinkService) GetUserByAccount(ctx context.Context, accNumber string) (*model.ArchivedUser, error) {

	return u.repo.GetUserByAccount(ctx, accNumber)
}
func (u *unlinkService) GetAllArchivedUser(ctx context.Context, filterParams *types.Filter) (*types.PaginatedResponse[*model.ArchivedUser], error) {
	return u.repo.GetAllArchivedUser(ctx, filterParams)
}
func (u *unlinkService) UnlinkUserCif(ctx context.Context, userCode string) error {
	userData := local_util.ExtractUserFromContext(ctx)
	if incomplet := local_util.IsIncomplete(userData); incomplet {
		return errors.New(localization.ErrorAccountNumberRequired.Code)
	}

	user, err := u.userRepo.FindByUserCode(ctx, userCode)
	if err != nil {
		return err
	}

	cpsAction := lib.CpsModelBuilder(userCode, userData, user, nil, string(constants.RequestUnlinkUser), constants.Delete)

	if u.cpsService.CreateCPSAction(ctx, &cpsAction); err != nil {
		return err
	}

	return nil
}
func (u *unlinkService) Authorize(ctx context.Context, cpsAction *model.CPSAction) error {

	if cpsAction.ActionStatus != constants.Approved {
		u.logger.Errorf("Try to authorize the collection without cps action approval")
		return errors.New(localization.ErrorCPSActionStatusInvalid.Code)
	}

	lib.GoRoutinBaker(types.BakerOptions{Sequential: false, UseMutex: false})
	userOldData, err := u.userRepo.FindByUserCode(ctx, cpsAction.UniqueId)
	if err != nil {
		return err
	}

	linkedAccountOldData, err := u.linkedAccountRepo.FindByCustomerNumber(ctx, userOldData.CustomerNumber)
	if err != nil {
		return err
	}

	if err := u.archivedUserRepo.Create(ctx, userOldData); err != nil {
		return err
	}

	if err := u.archivedLinkedAccountRepo.Create(ctx, linkedAccountOldData); err != nil {
		return err
	}

	if err := u.userRepo.Delete(ctx, userOldData.ID.Hex()); err != nil {
		return err
	}

	if err := u.linkedAccountRepo.Delete(ctx, linkedAccountOldData.ID.Hex()); err != nil {
		return err
	}

	return nil

}
