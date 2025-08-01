package avatar

import (
	"context"
	"fmt"
	"time"

	cps_const "github.com/CBE-Super-App/cbe-super-app-cps-action/internal/domain/cps_actions/constant"
	entities "github.com/CBE-Super-App/cbe-super-app-cps-action/internal/domain/cps_actions/entities"
	common_util "github.com/CBE-Super-App/cbe-super-app-cps-action/pkgs/utils"
	constant "github.com/CBE-Super-App/cbe-super-app-cps-action/utils"
	"gitlab.com/bersufekadgetachew/cbe-super-app-shared/shared/config"
	"gitlab.com/bersufekadgetachew/cbe-super-app-shared/shared/utils"
)

type AvatarDomainService interface {
	CreateAvatar(ctx context.Context, req AvatarRequest) (*Avatar, error)
	UpdateAvatar(ctx context.Context, id string, req AvatarRequest) (*Avatar, *Avatar, error)
	DeleteAvatar(ctx context.Context, id string) (*Avatar, *Avatar, error)
	EnableOrDisableAvatar(ctx context.Context, id string, enable bool) (*Avatar, *Avatar, error)
	Authorize(ctx context.Context, action *entities.CPSAction) (*entities.CPSAction, error)
	GetAvatar(ctx context.Context, id string) (*Avatar, error)
	GetAllAvatar(ctx context.Context, filterParams *constant.Filter) (*common_util.PaginatedResponse[[]*Avatar], error)
}

type AvatarDomain struct {
	avatarRepo  AvatarRepository
	bucketName  string
	minioClient config.MinioClientInterface
	logger      utils.Logger
	cfg         *config.VaultConfig
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

func (a *AvatarDomain) CreateAvatar(ctx context.Context, req AvatarRequest) (*Avatar, error) {

	url, err := common_util.UploadFileToMinio(ctx, a.minioClient, a.bucketName, req.Avatar, "avatar", a.cfg.MinioEndPoint, a.logger)
	if err != nil {
		a.logger.Errorf("Failed to upload avatar to MinIO: %v", err)
		return nil, err
	}

	result := Avatar{
		Label:          req.Label,
		Avatar:         url,
		Enable:         false,
		IsDeleted:      false,
		CreatedAt:      time.Now(),
		LastModifiedAt: time.Now(),
	}

	return &result, nil
}

func (a *AvatarDomain) UpdateAvatar(ctx context.Context, id string, req AvatarRequest) (*Avatar, *Avatar, error) {
	a.logger.Infof("Updating avatar", "id", id)

	existingAvatar, err := a.avatarRepo.GetAvatar(ctx, id)
	if err != nil {
		a.logger.Errorf("Failed to fetch avatar", "id", id, "error", err)
		return nil, nil, err
	}


	var url string
	if req.Avatar != nil {
		url, err = common_util.UploadFileToMinio(ctx, a.minioClient, a.bucketName, req.Avatar, "avatar", a.cfg.MinioEndPoint, a.logger)
		if err != nil {
			a.logger.Errorf("Failed to upload avatar to MinIO: %v", err)
			return nil, nil, err
		}
	} else {
		url = existingAvatar.Avatar
		a.logger.Infof("Using previous avatar URL", "url", url)
	}

	curAction := Avatar{
		ID:             id,
		Label:          common_util.NonEmptyString(req.Label, existingAvatar.Label),
		Avatar:         url,
		Enable:         existingAvatar.Enable,
		IsDeleted:      existingAvatar.IsDeleted,
		CreatedAt:      existingAvatar.CreatedAt,
		LastModifiedAt: time.Now(),
		DeletedAt:      existingAvatar.DeletedAt,
	}

	return &curAction, existingAvatar, nil
}

func (a *AvatarDomain) DeleteAvatar(ctx context.Context, id string) (*Avatar, *Avatar, error) {
	a.logger.Infof("Deleting avatar", "id", id)

	existingAvatar, err := a.avatarRepo.GetAvatar(ctx, id)
	if err != nil {
		a.logger.Errorf("Failed to fetch avatar", "id", id, "error", err)
		return nil, nil, err
	}

	curData := *existingAvatar
	curData.IsDeleted = true
	curData.DeletedAt = &time.Time{}
	curData.LastModifiedAt = time.Now()

	return &curData, existingAvatar, nil
}

func (a *AvatarDomain) EnableOrDisableAvatar(ctx context.Context, id string, enable bool) (*Avatar, *Avatar, error) {
	a.logger.Infof("EnableDisable avatar", "id", id, "enable", enable)

	existingAvatar, err := a.avatarRepo.GetAvatar(ctx, id)
	if err != nil {
		a.logger.Errorf("Failed to fetch avatar", "id", id, "error", err)
		return nil, nil, err
	}

	if enable && existingAvatar.Enable {
		return nil, nil, fmt.Errorf(common_util.ErrAlreadyEnabled)
	}

	if !enable && !existingAvatar.Enable {
		return nil, nil, fmt.Errorf(common_util.ErrAlreadyDisabled)
	}

	curData := *existingAvatar
	curData.Enable = enable
	curData.LastModifiedAt = time.Now()

	return &curData, existingAvatar, nil
}

func (a *AvatarDomain) Authorize(ctx context.Context, action *entities.CPSAction) (*entities.CPSAction, error) {
	var avatar Avatar
	bindErr := common_util.BindAction(action.CurrentAction, &avatar)
	if bindErr != nil {
		a.logger.Errorf("Failed to bind current action to avatar: %v", bindErr)
		return nil, fmt.Errorf(common_util.InvalidActionData)
	}

	var result *Avatar
	var err error

	switch action.RequestAction {
	case cps_const.RequestCreateAvatar:
		result, err = a.avatarRepo.CreateAvatar(ctx, avatar)
	case cps_const.RequestUpdateAvatar:
		result, err = a.avatarRepo.UpdateAvatar(ctx, avatar)
	case cps_const.RequestDeleteAvatar:
		result, err = a.avatarRepo.DeleteAvatar(ctx, avatar.ID)
	case cps_const.RequestEnableAvatar:
		result, err = a.avatarRepo.EnableDisableAvatar(ctx, avatar.ID, true)
	case cps_const.RequestDisableAvatar:
		result, err = a.avatarRepo.EnableDisableAvatar(ctx, avatar.ID, false)
	default:
		a.logger.Warnf("Unsupported request action: %s", action.RequestAction)
		return nil, fmt.Errorf(common_util.ErrUnsupported)
	}

	if err != nil {
		a.logger.Errorf("Failed to authorize action %s: %v", action.RequestAction, err)
		return nil, err
	}

	action.CurrentAction = result
	return action, nil
}

func (a *AvatarDomain) GetAvatar(ctx context.Context, id string) (*Avatar, error) {
	a.logger.Infof("Fetching avatar", "id", id)
	return a.avatarRepo.GetAvatar(ctx, id)
}

func (a *AvatarDomain) GetAllAvatar(ctx context.Context, filterParams *constant.Filter) (*common_util.PaginatedResponse[[]*Avatar], error) {
	return a.avatarRepo.GetAllAvatar(ctx, filterParams)
}