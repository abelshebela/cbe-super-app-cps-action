package unlink

import (
	"cbe-super-app-cps-action/internal/constants"
	"cbe-super-app-cps-action/internal/constants/dto/customer"
	"cbe-super-app-cps-action/internal/constants/lib"
	"cbe-super-app-cps-action/internal/constants/localization"
	"cbe-super-app-cps-action/internal/constants/types"
	"cbe-super-app-cps-action/internal/service"
	"cbe-super-app-cps-action/internal/service/unlink/core"
	"cbe-super-app-cps-action/internal/storage"
	local_util "cbe-super-app-cps-action/pkgs/utils"
	"context"
	"errors"
	"strings"

	"gitlab.com/bersufekadgetachew/cbe-super-app-shared/shared/entities/model"
	"gitlab.com/bersufekadgetachew/cbe-super-app-shared/shared/utils"
	"go.mongodb.org/mongo-driver/v2/bson"
	"go.mongodb.org/mongo-driver/v2/mongo"
	"go.opentelemetry.io/otel/attribute"
	"go.opentelemetry.io/otel/trace"
)

type unlinkService struct {
	logger utils.Logger
	// repo                      storage.UnlinkAccount
	userRepo                  storage.UserRepository
	archivedUserRepo          storage.ArchivedUserRepository
	linkedAccountRepo         storage.LinkedAccountRepository
	archivedLinkedAccountRepo storage.ArchivedLinkedAccountRepository
	accountBlockRepo          storage.AccountBlockRepository
	cpsService                service.CPSActionService
}

func NewUnlinkService(client *mongo.Client,
	userData storage.UserRepository,
	archivedUserRepo storage.ArchivedUserRepository,
	linkedAccountRepo storage.LinkedAccountRepository,
	archivedLinkedAccountRepo storage.ArchivedLinkedAccountRepository,
	accountBlockRepo storage.AccountBlockRepository,
	cpsAction service.CPSActionService,
	logger utils.Logger,
) service.UnlinkService {
	return &unlinkService{
		userRepo:                  userData,
		archivedUserRepo:          archivedUserRepo,
		linkedAccountRepo:         linkedAccountRepo,
		archivedLinkedAccountRepo: archivedLinkedAccountRepo,
		accountBlockRepo:          accountBlockRepo,
		cpsService:                cpsAction,
		logger:                    logger,
	}
}

func (u *unlinkService) GetUserByAccount(ctx context.Context, accNumber string) (customer.FindCustomerByIDResponse, error) {
	ctx, span := local_util.TraceLogger(ctx, "service", "GetUserByAccount", "unlinkService", "unlinkService")
	defer span.End()

	account, err := u.linkedAccountRepo.FindByAccountNumber(ctx, accNumber)
	if err != nil {
		span.AddEvent("FindByAccountNumber error", trace.WithAttributes(attribute.String("error", err.Error()), attribute.String("accountNumber", accNumber)))
		u.logger.Errorf("[UnlinkSvc][GetUserByAccount] find account err: %v", err)
		return customer.FindCustomerByIDResponse{}, err
	}

	user, err := u.userRepo.FindById(ctx, account.UserID.Hex())
	if err != nil {
		span.AddEvent("FindById error", trace.WithAttributes(attribute.String("error", err.Error()), attribute.String("userID", account.UserID.Hex())))
		u.logger.Errorf("[UnlinkSvc][GetUserByAccount] find user err: %v", err)
		return customer.FindCustomerByIDResponse{}, err
	}
	res := core.MapToDto(user)
	u.logger.Infof("[UnlinkSvc][GetUserByAccount] retrieved acc: %s", accNumber)
	return res, nil
}

func (u *unlinkService) GetAllArchivedUser(ctx context.Context, filterParams *types.Filter) (*types.PaginatedResponse[[]model.ArchivedUser], error) {
	ctx, span := local_util.TraceLogger(ctx, "service", "GetAllArchivedUser", "unlinkService", "unlinkService")
	defer span.End()
	result, err := u.archivedUserRepo.FindAllWithPagination(ctx, *filterParams)
	if err != nil {
		span.AddEvent("FindAllWithPagination error", trace.WithAttributes(attribute.String("error", err.Error())))
		u.logger.Errorf("[UnlinkSvc][GetAllArchived] fetch err: %v", err)
		return nil, err
	}
	u.logger.Infof("[UnlinkSvc][GetAllArchived] retrieved %d", len(result.Data))
	return result, nil
}

func (u *unlinkService) UnlinkUserCif(ctx context.Context, userCode string) error {
	ctx, span := local_util.TraceLogger(ctx, "service", "UnlinkUserCif", "unlinkService", "unlinkService")
	defer span.End()
	makerData := local_util.ExtractUserFromContext(ctx)
	if incomplet := local_util.IsIncomplete(makerData); incomplet {
		span.AddEvent("Incomplete user data", trace.WithAttributes(attribute.String("userCode", userCode)))
		u.logger.Errorf("[UnlinkSvc][UnlinkCif] incomplete user")
		return errors.New(localization.ErrorIncompleteUserInfo.Code)
	}

	user, err := u.userRepo.FindByUserCode(ctx, userCode)
	if err != nil {
		span.AddEvent("FindByUserCode error", trace.WithAttributes(attribute.String("error", err.Error()), attribute.String("userCode", userCode)))
		u.logger.Errorf("[UnlinkSvc][UnlinkCif] find user err: %v", err)
		return err
	}
	cpsAction := lib.CpsModelBuilder(user.UserCode, makerData, user, nil, string(constants.RequestUnlinkUser), constants.DELETE)

	if err := u.cpsService.CreateCPSAction(ctx, &cpsAction); err != nil {
		span.AddEvent("CPS action failed", trace.WithAttributes(attribute.String("error", err.Error()), attribute.String("userCode", userCode)))
		u.logger.Errorf("[UnlinkSvc][UnlinkCif] cps action err: %v", err)
		return err
	}

	span.AddEvent("Unlink user request created", trace.WithAttributes(attribute.String("userCode", userCode)))
	u.logger.Infof("[UnlinkSvc][UnlinkCif] request created code: %s", userCode)
	return nil
}

func (u *unlinkService) Authorize(ctx context.Context, cpsAction *model.CPSAction) (*model.CPSAction, error) {
	ctx, span := local_util.TraceLogger(ctx, "service", "AuthorizeUnlink", "unlinkService", "unlinkService")
	defer span.End()
	u.logger.Infof("[UnlinkSvc][Authorize] user_code: %s", cpsAction.UniqueId)
	haveAccount := false
	var archivedUserId, archivedLinkedAccountId string
	var linkedAccountOldData *model.LinkedAccount
	if cpsAction.ActionStatus != constants.Approved {
		span.AddEvent("CPS action status not approved", trace.WithAttributes(attribute.String("status", string(cpsAction.ActionStatus)), attribute.String("userCode", cpsAction.UniqueId)))
		u.logger.Errorf("[UnlinkSvc][Authorize] invalid status: %s", cpsAction.ActionStatus)
		return nil, errors.New(localization.ErrorCPSActionStatusInvalid.Code)
	}

	userOldData, err := u.userRepo.FindByUserCode(ctx, cpsAction.UniqueId)
	if err != nil {
		span.AddEvent("FindByUserCode error", trace.WithAttributes(attribute.String("error", err.Error()), attribute.String("userCode", cpsAction.UniqueId)))
		u.logger.Errorf("[UnlinkSvc][Authorize] find user err: %v", err)
		return nil, err
	}

	if strings.EqualFold(userOldData.CustomerNumber, "") {
		span.AddEvent("No account linked", trace.WithAttributes(attribute.String("userCode", cpsAction.UniqueId)))
		u.logger.Infof("[UnlinkSvc][Authorize] no account linked")
	} else {
		haveAccount = true
	}

	if haveAccount {
		linkedAccountOldData, err = u.linkedAccountRepo.FindByCustomerNumber(ctx, userOldData.CustomerNumber)
		if err != nil {
			span.AddEvent("FindByCustomerNumber error", trace.WithAttributes(attribute.String("error", err.Error()), attribute.String("customerNumber", userOldData.CustomerNumber)))
			u.logger.Errorf("[UnlinkSvc][Authorize] find linked acc err: %v", err)
			return nil, errors.New(localization.ErrorUnlinkFaild.Code)
		}
	}

	archivedUserId = userOldData.ID.Hex()
	archivedLinkedAccountId = linkedAccountOldData.ID.Hex()
	userOldData.ID = bson.NilObjectID
	linkedAccountOldData.ID = bson.NilObjectID

	archUserErr, archLinkedAccErr := core.CreatArchiveUserDataWithLinkedAccount(ctx, u.archivedUserRepo, u.archivedLinkedAccountRepo, userOldData, linkedAccountOldData, haveAccount)

	if archUserErr != nil || archLinkedAccErr != nil {
		span.AddEvent("Archiving error", trace.WithAttributes(attribute.String("archUserErr", errorString(archUserErr)), attribute.String("archLinkedAccErr", errorString(archLinkedAccErr))))
		u.logger.Errorf("[UnlinkSvc][Authorize] archive err: %v / %v", archUserErr, archLinkedAccErr)
		return nil, archUserErr
	}

	var userErr, linkedAccErr error
	lib.GoRoutinBaker(types.BakerOptions{Sequential: false, UseMutex: false},
		func() {
			u.logger.Infof("[UnlinkSvc][Authorize] deleting user: %s", archivedUserId)
			userErr = u.userRepo.Delete(ctx, archivedUserId)
		},
		func() {
			if haveAccount {
				u.logger.Infof("[UnlinkSvc][Authorize] deleting linked acc: %s", archivedLinkedAccountId)
				linkedAccErr = u.linkedAccountRepo.Delete(ctx, archivedLinkedAccountId)
			}
		},
	)
	if userErr != nil || linkedAccErr != nil {
		span.AddEvent("Delete error", trace.WithAttributes(attribute.String("userErr", errorString(userErr)), attribute.String("linkedAccErr", errorString(linkedAccErr))))
		u.logger.Errorf("[UnlinkSvc][Authorize] delete err: %v / %v", userErr, linkedAccErr)
		return nil, errors.New(localization.ErrorUnlinkFaild.Code)
	}

	span.AddEvent("User unlink authorized", trace.WithAttributes(attribute.String("userCode", cpsAction.UniqueId)))
	u.logger.Infof("[UnlinkSvc][Authorize] authorized code: %s", cpsAction.UniqueId)
	return nil, nil
}

// errorString safely returns the error string or empty if nil
func errorString(err error) string {
	if err != nil {
		return err.Error()
	}
	return ""
}
