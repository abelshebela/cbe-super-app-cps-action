package bpsuser

import (
	"cbe-super-app-cps-action/internal/constants"
	"cbe-super-app-cps-action/internal/constants/lib"
	localization "cbe-super-app-cps-action/internal/constants/localization"
	"cbe-super-app-cps-action/internal/constants/types"
	"cbe-super-app-cps-action/internal/service"
	bps_user_core "cbe-super-app-cps-action/internal/service/bps_user/core"
	"cbe-super-app-cps-action/internal/storage"
	local_util "cbe-super-app-cps-action/pkgs/utils"
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"time"

	// local_model "cbe-super-app-cps-action/internal/constants/model"

	bps_model "gitlab.com/bersufekadgetachew/cbe-super-app-shared/shared/entities/bps"
	"gitlab.com/bersufekadgetachew/cbe-super-app-shared/shared/entities/model"

	bpsUserDto "cbe-super-app-cps-action/internal/constants/dto/bps_user"

	"gitlab.com/bersufekadgetachew/cbe-super-app-shared/shared/utils"
	"go.mongodb.org/mongo-driver/v2/bson"
	"go.opentelemetry.io/otel/attribute"
	"go.opentelemetry.io/otel/trace"
)

type bpsUserService struct {
	cpsService service.CPSActionService
	repo       storage.BPSUserRepository
	CPSUserRepo storage.CpsUserRepository
	roles_repo storage.RoleRepository
	logger     utils.Logger
}

func NewBPSUserService(repo storage.BPSUserRepository, rolesRepo storage.RoleRepository, cpsService service.CPSActionService, cpsUserRepo storage.CpsUserRepository,logger utils.Logger) service.BPSUserService {
	return &bpsUserService{
		cpsService: cpsService,
		repo:       repo,
		CPSUserRepo: cpsUserRepo,
		roles_repo: rolesRepo,
		logger:     logger,
	}
}

// Authorize implements service.BPSUserService.
func (b *bpsUserService) Authorize(ctx context.Context, cpsAction *model.CPSAction) (*model.CPSAction, error) {
	ctx, span := local_util.TraceLogger(ctx, "service", "Authorize", "BPS User", "Authorize")
	defer span.End()

	var actionMap interface{}
	marshaled, err := json.Marshal(cpsAction.CurrentAction)
	if err != nil {
		span.AddEvent("failed to marshal CurrentAction", trace.WithAttributes(
			attribute.String("error", err.Error()),
			attribute.String("unique_id", cpsAction.UniqueId),
		))
		fmt.Printf("failed to marshal CurrentAction: %v\n", err)
		return nil, errors.New(localization.ErrorUnexpectedError.Code)
	}
	fmt.Printf("JSON bytes: %s\n", string(marshaled))
	err = json.Unmarshal(marshaled, &actionMap)
	if err != nil {
		span.AddEvent("failed to unmarshal CurrentAction", trace.WithAttributes(
			attribute.String("error", err.Error()),
			attribute.String("unique_id", cpsAction.UniqueId),
		))
		fmt.Printf("failed to unmarshal CurrentAction: %v\n", err)
		return nil, errors.New(localization.ErrorUnexpectedError.Code)
	}

	actionData := bps_user_core.BPSUser_mapper(actionMap.(map[string]interface{}))
	if cpsAction.UniqueId != "" {
		objID, err := bson.ObjectIDFromHex(cpsAction.UniqueId)
		if err != nil {
			span.AddEvent("[Authorize] failed to parse unique id", trace.WithAttributes(
				attribute.String("error", err.Error()),
				attribute.String("unique_id", cpsAction.UniqueId),
			))
			b.logger.Errorf("[BpsUserSvc][Authorize] parse unique id err: %v", err)
			return nil, errors.New(localization.ErrorUnexpectedError.Code)
		}
		actionData.ID = objID
	}

	actionData.LastModifiedAt = time.Now()
	// local_actionData := bps_user_core.MapBPSUserToWithJobTitle(actionData)

	switch cpsAction.RequestAction {
	case string(constants.RequestCreateBPSUser):
		if err := b.repo.Create(ctx, actionData); err != nil {
			span.AddEvent("[Authorize] failed to create BPS user", trace.WithAttributes(
				attribute.String("error", err.Error()),
				attribute.String("unique_id", cpsAction.UniqueId),
			))
			b.logger.Errorf("[BpsUserSvc][Authorize] create err: %v", err)
			return nil, err
		}
		b.logger.Infof("[BpsUserSvc][Authorize] created")
		return nil, nil
	case string(constants.RequestBpsUserUpdate):
		if err := b.repo.Update(ctx, &actionData); err != nil {
			span.AddEvent("[Authorize] failed to update BPS user", trace.WithAttributes(
				attribute.String("error", err.Error()),
				attribute.String("unique_id", cpsAction.UniqueId),
			))
			b.logger.Errorf("[BpsUserSvc][Authorize] update err: %v", err)
			return nil, err
		}
		b.logger.Infof("[BpsUserSvc][Authorize] updated id: %s", cpsAction.UniqueId)
		return nil, nil
	case string(constants.RequestEnableBPSUser):
		actionData.Enabled = true
		b.logger.Infof("[BpsUserSvc][Authorize] enabling id: %s", cpsAction.UniqueId)
	case string(constants.RequestDisableBPSUser):
		actionData.Enabled = false
		b.logger.Infof("[BpsUserSvc][Authorize] disabling id: %s", cpsAction.UniqueId)
	default:
		span.AddEvent("[Authorize] unsupported action", trace.WithAttributes(attribute.String("action", cpsAction.RequestAction)))
		b.logger.Errorf("[BpsUserSvc][Authorize] unsupported action: %s", cpsAction.RequestAction)
		return nil, errors.New(localization.ErrorActionNotFound.Code)
	}
	if err := b.repo.Update(ctx, &actionData); err != nil {
		span.AddEvent("[Authorize] failed to update BPS user", trace.WithAttributes(
			attribute.String("error", err.Error()),
			attribute.String("unique_id", cpsAction.UniqueId),
		))
		b.logger.Errorf("[BpsUserSvc][Authorize] update err: %v", err)
		return nil, err
	}
	b.logger.Infof("[BpsUserSvc][Authorize] authorized id: %s", cpsAction.UniqueId)
	return nil, nil
}

// FetchUserByUserCode implements service.BPSUserService.
func (b *bpsUserService) FetchUserByUserCode(ctx context.Context, userCode string) (*bps_model.BPSUser, error) {
	ctx, span := local_util.TraceLogger(ctx, "service", "FetchUserByUserCode", "BPS User", "FetchUserByUserCode")
	defer span.End()

	user, err := b.repo.GetByUserCode(ctx, userCode)
	if err != nil {
		span.AddEvent("[FetchUserByUserCode] failed to fetch BPS user", trace.WithAttributes(
			attribute.String("error", err.Error()),
			attribute.String("user_code", userCode),
		))
		b.logger.Errorf("[BpsUserSvc][FetchByCode] fetch err: %v", err)
		return nil, err
	}
	b.logger.Infof("[BpsUserSvc][FetchByCode] retrieved code: %s", userCode)
	return user, nil
}

// GetAllBPSUsers implements service.BPSUserService.
func (b *bpsUserService) GetAllBPSUsers(ctx context.Context, filterParams *types.Filter) (*types.PaginatedResponse[[]bpsUserDto.BPSUserResposenDTO], error) {
	ctx, span := local_util.TraceLogger(ctx, "service", "GetAllBPSUsers", "BPS User", "GetAllBPSUsers")
	defer span.End()

	result, err := b.repo.FindAllWithPagination(ctx, *filterParams)
	if err != nil {
		span.AddEvent("[GetAllBPSUsers] failed to fetch BPS users", trace.WithAttributes(
			attribute.String("error", err.Error()),
		))
		b.logger.Errorf("[BpsUserSvc][GetAll] fetch err: %v", err)
		return nil, err
	}
	b.logger.Infof("[BpsUserSvc][GetAll] retrieved %d", len(result.Data))
	return result, nil
}

// UpdateBpsUser implements service.BPSUserService.
func (b *bpsUserService) UpdateStatusBpsUser(ctx context.Context, userCode string, status bool) error {
	ctx, span := local_util.TraceLogger(ctx, "service", "UpdateBpsUser", "BPS User", "UpdateBpsUser")
	defer span.End()

	b.logger.Infof("[BpsUserSvc][UpdateStatus] enabled: %v", status)
	makerData := local_util.ExtractUserFromContext(ctx)
	user, err := b.repo.GetByUserCode(ctx, userCode)
	if err != nil {
		span.AddEvent("[UpdateBpsUser] failed to find BPS user", trace.WithAttributes(
			attribute.String("error", err.Error()),
			attribute.String("user_code", userCode),
		))
		b.logger.Errorf("[BpsUserSvc][UpdateStatus] find err: %v", err)
		return err
	}

	if user == nil {
		span.AddEvent("[UpdateBpsUser] BPS user not found", trace.WithAttributes(attribute.String("user_code", userCode)))
		b.logger.Errorf("[BpsUserSvc][UpdateStatus] not found: %s", userCode)
		return errors.New(localization.ErrorUserNotFound.Code)
	}

	if user.Enabled && status {
		span.AddEvent("[UpdateBpsUser] BPS user already enabled", trace.WithAttributes(attribute.String("user_code", userCode)))
		b.logger.Errorf("[BpsUserSvc][UpdateStatus] already enabled: %s", userCode)
		return errors.New(localization.ErrorUserAlreadyEnabled.Code)
	}

	if !user.Enabled && !status {
		span.AddEvent("[UpdateBpsUser] BPS user already disabled", trace.WithAttributes(attribute.String("user_code", userCode)))
		b.logger.Errorf("[BpsUserSvc][UpdateStatus] already disabled: %s", userCode)
		return errors.New(localization.ErrorUserAlreadyDisabled.Code)
	}

	updatedUser := *user
	updatedUser.Enabled = status

	var requestAction string
	if status {
		requestAction = string(constants.RequestEnableBPSUser)
	} else {
		requestAction = string(constants.RequestDisableBPSUser)

	}

	cpsActionData := lib.CpsModelBuilder(user.ID.Hex(), makerData, user, updatedUser, requestAction, constants.UPDATE)

	if err := b.cpsService.CreateCPSAction(ctx, &cpsActionData); err != nil {
		span.AddEvent("[UpdateBpsUser] failed to create CPS action", trace.WithAttributes(
			attribute.String("error", err.Error()),
			attribute.String("user_code", userCode),
		))
		b.logger.Errorf("[BpsUserSvc][UpdateStatus] cps action err: %v", err)
		return err
	}
	b.logger.Infof("[BpsUserSvc][UpdateStatus] request created code: %s", userCode)
	return nil
}
func (b *bpsUserService) CreateBPSUser(ctx context.Context, req bps_model.BPSUser) error {
	ctx, span := local_util.TraceLogger(ctx, "service", "CreateBPSUser", "BPS User", "CreateBPSUser")
	defer span.End()
	makerData := local_util.ExtractUserFromContext(ctx)

	md := &types.ContextMetadata{}
	ctx = context.WithValue(ctx, constants.ContextKeyMetadata, md)

	existing, err := b.repo.FindByOr(ctx, req.PhoneNumber, req.Email, req.Username)
	if err != nil {
		if err.Error() != localization.ErrorResourceNotFound.Code {
			b.logger.Errorf("[BpsUserSvc][Create] check existing err: %v", err)
			return err
		}
	}

	if err := bps_user_core.ExistingIdentifier(existing, req); err != nil {
		b.logger.Infof("[BpsUserSvc][Create] duplicate data: %v", err)
		return err
	}

	roles, err := b.roles_repo.FindByFilterKey(ctx, "job_title", req.JobTitle)
	if err != nil {
		b.logger.Errorf("[BpsUserSvc][Create] role lookup err for job_title: %s, err: %v", req.JobTitle, err)
		return errors.New(localization.ErrorRoleNotFound.Code)
	}
	if roles == nil || roles.Role == "" {
		b.logger.Errorf("[BpsUserSvc][Create] role not found for job_title: %s", req.JobTitle)
		return errors.New(localization.ErrorRoleNotFound.Code)
	}

	is_exist_on_CPS,err := b.CPSUserRepo.FindByEmailOrPhoneNumberOrUserName(ctx,req.Email,req.PhoneNumber,req.Username)
	if  err != nil {
		b.logger.Errorf("[CreateBPSUser] got error while checking user data exist on cps user ")
		return errors.New(localization.ErrorInternalServerError.Code)
	} 
	if is_exist_on_CPS != nil{
		if  is_exist_on_CPS.UserName != "" && is_exist_on_CPS.UserName == req.Username {
		b.logger.Errorf("[CreateBPSUser] user name already exist")
		return errors.New(localization.ErrorExistUserName.Code)
		}

		if is_exist_on_CPS.Email != "" && is_exist_on_CPS.Email == req.Email {
			b.logger.Errorf("[CreateBPSUser] email already exist")
			return errors.New(localization.ErrorExistEmail.Code)

		}

		if is_exist_on_CPS.PhoneNumber != "" && is_exist_on_CPS.PhoneNumber ==req.PhoneNumber{
			b.logger.Errorf("[CreateBPSUser] phone number already exist")
			return errors.New(localization.ErrorExistPhoneNumber.Code)
		}
	}

	req.UserCode = local_util.UniqueIdGenerator()
	// Build CPS action model for create
	cpsActionModel := lib.CpsModelBuilder(
		"",                                     // unique id
		makerData,                              // maker data
		nil,                                    // old data (nil for create)
		req,                                    // new data
		string(constants.RequestCreateBPSUser), // request action
		constants.CREATE,                       // action type
	)

	// Create CPS action
	if err := b.cpsService.CreateCPSAction(ctx, &cpsActionModel); err != nil {
		span.AddEvent("[CreateBPSUser] failed to create CPS action", trace.WithAttributes(
			attribute.String("error", err.Error()),
			attribute.String("user_code", req.UserCode),
		))
		b.logger.Errorf("[BpsUserSvc][Create] cps action err: %v", err)
		return err
	}
	if md.IsMakerOnly {

	} else {

	}
	b.logger.Infof("[BpsUserSvc][Create] request created code: %s", req.UserCode)
	return nil
}

func (b *bpsUserService) UpdateBPSUser(ctx context.Context, userID string, updatedUser bps_model.BPSUser) error {
	ctx, span := local_util.TraceLogger(ctx, "service", "UpdateBPSUser", "BPS User", "UpdateBPSUser")
	defer span.End()
	makerData := local_util.ExtractUserFromContext(ctx)

	existing, err := b.repo.FindByOr(ctx, updatedUser.PhoneNumber, updatedUser.Email, updatedUser.Username)
	if err != nil {
		if err.Error() != localization.ErrorResourceNotFound.Code {
			b.logger.Errorf("[BpsUserSvc][Update] check existing err: %v", err)
			return err
		}
	}

	if err := bps_user_core.ExistingIdentifierForUpdate(*existing, userID, updatedUser); err != nil {
		b.logger.Infof("[BpsUserSvc][Update] duplicate data: %v", err)
		return err
	}

	is_exist_on_CPS,err := b.CPSUserRepo.FindByEmailOrPhoneNumberOrUserName(ctx,updatedUser.Email,updatedUser.PhoneNumber,updatedUser.Username)
	if  err != nil {
		b.logger.Errorf("[CreateBPSUser] got error while checking user data exist on cps user ")
		return errors.New(localization.ErrorInternalServerError.Code)
	} 
	if is_exist_on_CPS != nil{
		if  is_exist_on_CPS.UserName != "" && is_exist_on_CPS.UserName == updatedUser.Username {
		b.logger.Errorf("[CreateBPSUser] user name already exist")
		return errors.New(localization.ErrorExistUserName.Code)
		}

		if is_exist_on_CPS.Email != "" && is_exist_on_CPS.Email == updatedUser.Email {
			b.logger.Errorf("[CreateBPSUser] email already exist")
			return errors.New(localization.ErrorExistEmail.Code)

		}

		if is_exist_on_CPS.PhoneNumber != "" && is_exist_on_CPS.PhoneNumber == updatedUser.PhoneNumber{
			b.logger.Errorf("[CreateBPSUser] phone number already exist")
			return errors.New(localization.ErrorExistPhoneNumber.Code)
		}
	}

	// updatedUser.Role = roles.Role
	cpsActionModel := lib.CpsModelBuilder(
		existing.ID.Hex(),                      // unique id
		makerData,                              // maker data
		existing,                               // old data
		updatedUser,                            // new data
		string(constants.RequestBpsUserUpdate), // request action
		constants.UPDATE,                       // action type
	)

	// Create CPS action
	if err := b.cpsService.CreateCPSAction(ctx, &cpsActionModel); err != nil {
		span.AddEvent("[UpdateBPSUser] failed to create CPS action", trace.WithAttributes(
			attribute.String("error", err.Error()),
			attribute.String("user_code", userID),
		))
		b.logger.Errorf("[BpsUserSvc][Update] cps action err: %v", err)
		return err
	}

	b.logger.Infof("[BpsUserSvc][Update] request created id: %s", userID)
	return nil
}
