package unlink

import (
	"cbe-super-app-cps-action/internal/constants"
	"cbe-super-app-cps-action/internal/constants/lib"
	"cbe-super-app-cps-action/internal/constants/localization"
	"cbe-super-app-cps-action/internal/constants/model"
	"cbe-super-app-cps-action/internal/constants/types"
	"cbe-super-app-cps-action/internal/service"
	"cbe-super-app-cps-action/internal/service/unlink/core"
	"cbe-super-app-cps-action/internal/storage"
	local_util "cbe-super-app-cps-action/pkgs/utils"
	"context"
	"errors"
	"strings"

	"gitlab.com/bersufekadgetachew/cbe-super-app-shared/shared/utils"
	"go.mongodb.org/mongo-driver/v2/bson"
	"go.mongodb.org/mongo-driver/v2/mongo"
)

type unlinkService struct {
	logger utils.Logger
	// repo                      storage.UnlinkAccount
	userRepo                  storage.UserRepository
	archivedUserRepo          storage.ArchivedUserRepository
	linkedAccountRepo         storage.LinkedAccountRepository
	archivedLinkedAccountRepo storage.ArchivedLinkedAccountRepository
	cpsService                service.CPSActionService
}

func NewUnlinkService(client *mongo.Client, userData storage.UserRepository, archivedUserRepo storage.ArchivedUserRepository, linkedAccountRepo storage.LinkedAccountRepository, archivedLinkedAccountRepo storage.ArchivedLinkedAccountRepository, cpsAction service.CPSActionService, logger utils.Logger) service.UnlinkService {
	return &unlinkService{
		userRepo:                  userData,
		archivedUserRepo:          archivedUserRepo,
		linkedAccountRepo:         linkedAccountRepo,
		archivedLinkedAccountRepo: archivedLinkedAccountRepo,
		cpsService:                cpsAction,
		logger:                    logger,
	}
}

func (u *unlinkService) GetUserByAccount(ctx context.Context, accNumber string) (*model.User, error) {
	account, err := u.linkedAccountRepo.FindByAccountNumber(ctx, accNumber)
	if err != nil {
		u.logger.Errorf("[GetUserByAccount] failed to find account: %v", err)
		return nil, err
	}

	user, err := u.userRepo.FindById(ctx, account.UserID.Hex())
	if err != nil {
		u.logger.Errorf("[GetUserByAccount] failed to find user: %v", err)
		return nil, err
	}

	u.logger.Infof("[GetUserByAccount] user retrieved successfully for account number")
	return user, nil
}
func (u *unlinkService) GetAllArchivedUser(ctx context.Context, filterParams *types.Filter) (*types.PaginatedResponse[[]*model.ArchivedUser], error) {
	result, err := u.archivedUserRepo.FindAllWithPagination(ctx, *filterParams)
	if err != nil {
		u.logger.Errorf("[GetAllArchivedUser] failed to fetch archived users: %v", err)
		return nil, err
	}
	u.logger.Infof("[GetAllArchivedUser] retrieved %d archived users", len(result.Data))
	return result, nil
}
func (u *unlinkService) UnlinkUserCif(ctx context.Context, userCode string) error {
	makerData := local_util.ExtractUserFromContext(ctx)
	if incomplet := local_util.IsIncomplete(makerData); incomplet {
		u.logger.Errorf("[UnlinkUserCif] incomplete user data")
		return errors.New(localization.ErrorIncompleteUserInfo.Code)
	}

	user, err := u.userRepo.FindByUserCode(ctx, userCode)
	if err != nil {
		u.logger.Errorf("[UnlinkUserCif] failed to find user: %v", err)
		return err
	}
	cpsAction := lib.CpsModelBuilder(user.UserCode, makerData, user, nil, string(constants.RequestUnlinkUser), constants.DELETE)

	if err := u.cpsService.CreateCPSAction(ctx, &cpsAction); err != nil {
		u.logger.Errorf("[UnlinkUserCif] failed to create CPS action: %v", err)
		return err
	}

	u.logger.Infof("[UnlinkUserCif] unlink user request created successfully for user_code: %s", userCode)
	return nil
}
func (u *unlinkService) Authorize(ctx context.Context, cpsAction *model.CPSAction) (*model.CPSAction, error) {
	u.logger.Infof("[Authorize] authorizing unlink user action for user_code: %s", cpsAction.UniqueId)
	haveAccount := false
	var archivedUserId, archivedLinkedAccountId string
	var linkedAccountOldData *model.LinkedAccount
	if cpsAction.ActionStatus != constants.Approved {
		u.logger.Errorf("[Authorize] CPS action status is not approved: %s", cpsAction.ActionStatus)
		return nil, errors.New(localization.ErrorCPSActionStatusInvalid.Code)
	}

	userOldData, err := u.userRepo.FindByUserCode(ctx, cpsAction.UniqueId)
	if err != nil {
		u.logger.Errorf("[Authorize] failed to find user: %v", err)
		return nil, err
	}

	if strings.EqualFold(userOldData.CustomerNumber, "") {
		u.logger.Infof("[Authorize] user doesn't have any account linked")
	} else {
		haveAccount = true
	}

	if haveAccount {
		linkedAccountOldData, err = u.linkedAccountRepo.FindByCustomerNumber(ctx, userOldData.CustomerNumber)
		if err != nil {
			u.logger.Errorf("[Authorize] failed to find linked account: %v", err)
			return nil, errors.New(localization.ErrorUnlinkFaild.Code)
		}
	}

	archivedUserId = userOldData.ID.Hex()
	archivedLinkedAccountId = linkedAccountOldData.ID.Hex()
	userOldData.ID = bson.NilObjectID
	linkedAccountOldData.ID = bson.NilObjectID

	archUserErr, archLinkedAccErr := core.CreatArchiveUserDataWithLinkedAccount(ctx, u.archivedUserRepo, u.archivedLinkedAccountRepo, userOldData, linkedAccountOldData, haveAccount)

	if archUserErr != nil || archLinkedAccErr != nil {
		u.logger.Errorf("[Authorize] error occurred during archiving user: %v / %v", archUserErr, archLinkedAccErr)
		return nil, archUserErr
	}

	var userErr, linkedAccErr error
	lib.GoRoutinBaker(types.BakerOptions{Sequential: false, UseMutex: false},
		func() {
			u.logger.Infof("[Authorize] deleting user in routine, userid: %s", archivedUserId)
			userErr = u.userRepo.Delete(ctx, archivedUserId)
		},
		func() {
			if haveAccount {
				u.logger.Infof("[Authorize] deleting linked account in routine, account id: %s", archivedLinkedAccountId)
				linkedAccErr = u.linkedAccountRepo.Delete(ctx, archivedLinkedAccountId)
			}
		},
	)
	if userErr != nil || linkedAccErr != nil {
		u.logger.Errorf("[Authorize] error occurred during deleting user: %v / %v", userErr, linkedAccErr)

		return nil, errors.New(localization.ErrorUnlinkFaild.Code)
	}

	u.logger.Infof("[Authorize] user unlink authorized successfully for user_code: %s", cpsAction.UniqueId)
	return nil, nil

}
