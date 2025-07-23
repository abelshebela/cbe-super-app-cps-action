package avatar

import (
	"context"
	"fmt"
	"strings"
	"time"

	"github.com/CBE-Super-App/cbe-super-app-cps-action/internal/adapter/outbound/model"
	cps_constants "github.com/CBE-Super-App/cbe-super-app-cps-action/internal/domain/cps_actions/constant"
	entities "github.com/CBE-Super-App/cbe-super-app-cps-action/internal/domain/cps_actions/entities"
	common_util "github.com/CBE-Super-App/cbe-super-app-cps-action/pkgs/utils"
	constant "github.com/CBE-Super-App/cbe-super-app-cps-action/utils"
	"gitlab.com/bersufekadgetachew/cbe-super-app-shared/shared/config"
	"gitlab.com/bersufekadgetachew/cbe-super-app-shared/shared/utils"
)

type AvatarDomain struct {
	avatarRepo  AvatarRepository
	bucketName  string
	minioClient config.MinioClientInterface
	logger      utils.Logger
	cfg         *config.VaultConfig
}

type AvatarDomainService interface {
	CreateAvatar(ctx context.Context, req model.CreateCPSAction) (*entities.CPSAction, error)
	UpdateAvatar(ctx context.Context, id string, req model.CreateCPSAction) (*entities.CPSAction, error)
	DeleteAvatar(ctx context.Context, id string, req model.CreateCPSAction) (*entities.CPSAction, error)
	Authorize(ctx context.Context, cpsAction *entities.CPSAction) (*entities.CPSAction, error)
	EnableOrDisableAvatar(ctx context.Context, id string, requestAction model.RequestAction, req model.CreateCPSAction) (*entities.CPSAction, error)
	GetAvatar(ctx context.Context, id string) (*Avatar, error)
	GetAllAvatar(ctx context.Context, filterParams constant.Filter) (*common_util.PaginatedResponse[[]*Avatar], error)
}

func stripFieldPrefix(err error) string {
	if err == nil {
		return ""
	}
	errStr := err.Error()
	if idx := strings.Index(errStr, ":"); idx != -1 {
		errStr = strings.TrimSpace(errStr[idx+1:])
	}
	errStr = strings.TrimSuffix(errStr, ".")
	return errStr
}

func InitAvatarDomain(avatarRepo AvatarRepository, minioClient config.MinioClientInterface, bucketName string, cfg *config.VaultConfig, logger utils.Logger) AvatarDomainService {
	return &AvatarDomain{
		avatarRepo:  avatarRepo,
		bucketName:  bucketName,
		minioClient: minioClient,
		logger:      logger,
		cfg:         cfg,
	}
}

func (a *AvatarDomain) CreateAvatar(ctx context.Context, req model.CreateCPSAction) (*entities.CPSAction, error) {
	req.RequestAction = model.RequestCreateAvatar

	actionData, ok := req.ActionData.(CreateAvatar)
	if !ok {
		a.logger.Errorf("failed to cast action data to avatar request")
		return nil, fmt.Errorf("FAILED_TO_CAST_ACTION_DATA")
	}

	if err := actionData.Validate(); err != nil {
		a.logger.Errorf("validation error: %v", err)
		return nil, fmt.Errorf("%s", stripFieldPrefix(err))
	}

	url, err := common_util.UploadFileToMinio(ctx, a.minioClient, a.bucketName, actionData.Avatar, "avatar", a.cfg.MinioEndPoint, a.logger)
	if err != nil {
		return nil, err
	}

	// Create CPSAction
	cpsActionRes := &entities.CPSAction{
		ActionCode:       utils.RandomGenerator(20),
		MakerID:          req.MakerUser.UserCode,
		MakerName:        req.MakerUser.FullName,
		MakerPhoneNumber: req.MakerUser.PhoneNumber,
		Department:       req.Department,
		CurrentAction: Avatar{
			Label:          actionData.Label,
			Avatar:         url,
			Enable:         true,
			IsDeleted:      false,
			CreatedAt:      time.Now(),
			LastModifiedAt: time.Now(),
		},
		RequestAction:   cps_constants.RequestCreateAvatar,
		ActionStatus:    cps_constants.ActionPending,
		ActionType:      cps_constants.ActionCreate,
		MakerActionTime: time.Now(),
		CreatedAt:       time.Now(),
		LastModifiedAt:  time.Now(),
	}

	return cpsActionRes, nil
}

func (a *AvatarDomain) UpdateAvatar(ctx context.Context, id string, req model.CreateCPSAction) (*entities.CPSAction, error) {

	actionData, ok := req.ActionData.(UpdateAvatar)
	if !ok {
		a.logger.Errorf("failed to cast action data to avatar request")
		return nil, fmt.Errorf("FAILED_TO_CAST_ACTION_DATA")
	}

	if err := actionData.Validate(); err != nil {
		a.logger.Errorf("validation error: %v", err)
		return nil, fmt.Errorf("%s", stripFieldPrefix(err))
	}

	// Check if avatar exists
	existingAvatar, err := a.avatarRepo.GetAvatar(ctx, id)
	if err != nil {
		return nil, err
	}

	var url string
	if actionData.Avatar != nil {
		url, err = common_util.UploadFileToMinio(
			ctx,
			a.minioClient,
			a.bucketName,
			actionData.Avatar,
			"avatar",
			a.cfg.MinioEndPoint,
			a.logger,
		)
		if err != nil {
			return nil, err
		}
	}

	// Construct CPSAction
	cpsActionRes := &entities.CPSAction{
		ActionCode:       utils.RandomGenerator(20),
		MakerID:          req.MakerUser.UserCode,
		MakerName:        req.MakerUser.FullName,
		MakerPhoneNumber: req.MakerUser.PhoneNumber,
		Department:       req.Department,
		ActionStatus:     cps_constants.ActionPending,
		ActionType:       cps_constants.ActionUpdate,
		RequestAction:    cps_constants.RequestUpdateAvatar,
		CurrentAction: Avatar{
			ID:             id,
			Label:          actionData.Label,
			Avatar:         url,
			Enable:         existingAvatar.Enable,
			IsDeleted:      existingAvatar.IsDeleted,
			CreatedAt:      existingAvatar.CreatedAt,
			LastModifiedAt: time.Now(),
			DeletedAt:      existingAvatar.DeletedAt,
		},
		PreviousAction:  *existingAvatar,
		MakerActionTime: time.Now(),
		CreatedAt:       time.Now(),
		LastModifiedAt:  time.Now(),
	}

	return cpsActionRes, nil
}

func (a *AvatarDomain) DeleteAvatar(ctx context.Context, id string, req model.CreateCPSAction) (*entities.CPSAction, error) {

	// Check if avatar exists
	existingAvatar, err := a.avatarRepo.GetAvatar(ctx, id)
	if err != nil {
		return nil, err
	}

	// Construct CPSAction
	cpsActionRes := &entities.CPSAction{
		ActionCode:       utils.RandomGenerator(20),
		MakerID:          req.MakerUser.UserCode,
		MakerName:        req.MakerUser.FullName,
		MakerPhoneNumber: req.MakerUser.PhoneNumber,
		Department:       req.Department,
		ActionStatus:     cps_constants.ActionPending,
		ActionType:       cps_constants.ActionDelete,
		RequestAction:    cps_constants.RequestDeleteAvatar,
		CurrentAction: Avatar{
			ID:             id,
			IsDeleted:      true,
			DeletedAt:      &time.Time{},
			LastModifiedAt: time.Now(),
		},
		PreviousAction:  *existingAvatar,
		MakerActionTime: time.Now(),
		CreatedAt:       time.Now(),
		LastModifiedAt:  time.Now(),
	}

	return cpsActionRes, nil
}

func (a *AvatarDomain) Authorize(ctx context.Context, cpsAction *entities.CPSAction) (*entities.CPSAction, error) {

	switch string(cpsAction.ActionType) {
	case string(model.ActionCreate):
		return a.avatarRepo.AuthorizeCreateAvatar(ctx, cpsAction)

	case string(model.ActionUpdate):
		return a.avatarRepo.AuthorizeUpdateAvatar(ctx, cpsAction)

	case string(model.ActionDelete):
		return a.avatarRepo.AuthorizeDeleteAvatar(ctx, cpsAction)

	default:
		return cpsAction, nil
	}
}

func (a *AvatarDomain) EnableOrDisableAvatar(ctx context.Context, id string, requestAction model.RequestAction, req model.CreateCPSAction) (*entities.CPSAction, error) {

	// Check if avatar exists
	existingAvatar, err := a.avatarRepo.GetAvatar(ctx, id)
	if err != nil {
		return nil, err
	}

	// Construct CPSAction
	enable := requestAction == model.RequestEnableAvatar

	if enable && existingAvatar.Enable {
		return nil, fmt.Errorf("RESOURCE_ALREADY_ENABLED")
	}

	if !enable && !existingAvatar.Enable {
		return nil, fmt.Errorf("RESOURCE_ALREADY_DISABLED")
	}
	cpsActionRes := &entities.CPSAction{
		ActionCode:       utils.RandomGenerator(20),
		MakerID:          req.MakerUser.UserCode,
		MakerName:        req.MakerUser.FullName,
		MakerPhoneNumber: req.MakerUser.PhoneNumber,
		Department:       req.Department,
		ActionStatus:     cps_constants.ActionPending,
		ActionType:       cps_constants.ActionUpdate,
		RequestAction:    cps_constants.RequestAction(requestAction),
		CurrentAction: Avatar{
			ID:             id,
			Avatar:         existingAvatar.Avatar,
			Label:          existingAvatar.Label,
			Enable:         enable,
			IsDeleted:      existingAvatar.IsDeleted,
			CreatedAt:      existingAvatar.CreatedAt,
			LastModifiedAt: time.Now(),
			DeletedAt:      existingAvatar.DeletedAt,
		},
		PreviousAction:  *existingAvatar,
		MakerActionTime: time.Now(),
		CreatedAt:       time.Now(),
		LastModifiedAt:  time.Now(),
	}

	return cpsActionRes, nil
}

func (a *AvatarDomain) GetAllAvatar(ctx context.Context, filterParams constant.Filter) (*common_util.PaginatedResponse[[]*Avatar], error) {
	return a.avatarRepo.GetAllAvatar(ctx, filterParams)
}

func (a *AvatarDomain) GetAvatar(ctx context.Context, id string) (*Avatar, error) {
	return a.avatarRepo.GetAvatar(ctx, id)
}
