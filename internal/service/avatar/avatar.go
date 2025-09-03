package avatar

import (
	"cbe-super-app-cps-action/internal/constants"
	"cbe-super-app-cps-action/internal/constants/lib"
	"cbe-super-app-cps-action/internal/constants/localization"
	"cbe-super-app-cps-action/internal/constants/model"
	"cbe-super-app-cps-action/internal/constants/types"
	"cbe-super-app-cps-action/internal/service"
	"cbe-super-app-cps-action/internal/service/avatar/core"
	"cbe-super-app-cps-action/internal/storage"
	local_util "cbe-super-app-cps-action/pkgs/utils"
	"context"
	"errors"
	"fmt"
	"mime/multipart"
	"time"

	"gitlab.com/bersufekadgetachew/cbe-super-app-shared/shared/config"
	"gitlab.com/bersufekadgetachew/cbe-super-app-shared/shared/utils"
	"go.mongodb.org/mongo-driver/v2/bson"
)

type avatarService struct {
	logger        utils.Logger
	avatar        storage.AvatarRepository
	minio         config.MinioClientInterface
	cpsService    service.CPSActionService
	bucketName    string
	minioEndPoint string
}

func NewAvatarService(avatar storage.AvatarRepository, cpsService service.CPSActionService, logger utils.Logger, minio config.MinioClientInterface, buckateName, minioEndPoint string) service.AvatarService {
	return &avatarService{
		cpsService:    cpsService,
		avatar:        avatar,
		logger:        logger,
		minio:         minio,
		bucketName:    buckateName,
		minioEndPoint: minioEndPoint,
	}
}

func (a *avatarService) CreateAvatar(ctx context.Context, avatar *model.Avatar, fileHeader *multipart.FileHeader) error {
	makerData := local_util.ExtractUserFromContext(ctx)
	existing, err := a.avatar.FindAll(ctx, bson.M{"label": avatar}, nil)
	if len(existing) > 0 {
		return errors.New(localization.ErrorAvatarAlreadyExist.Code)
	}

	url, err := lib.UploadFileToMinio(ctx, a.minio, a.bucketName, fileHeader, string(constants.Avatar), a.minioEndPoint, a.logger)
	if err != nil {
		return err
	}

	cpsModel := lib.CpsModelBuilder("", makerData, nil, model.Avatar{Avatar: url, Label: avatar.Label, CreatedAt: time.Now()}, string(constants.RequestCreateAvatar), constants.CREATE)
	if err := a.cpsService.CreateCPSAction(ctx, &cpsModel); err != nil {
		return err
	}
	return nil
}
func (a *avatarService) UpdateAvatar(ctx context.Context, id string, avatar *model.Avatar, fileHeader *multipart.FileHeader, fromEnabledDisable bool) error {
	makerData := local_util.ExtractUserFromContext(ctx)
	if incomplete := local_util.IsIncomplete(makerData); incomplete {
		return errors.New(localization.ErrorAccountNumberRequired.Code)
	}
	fmt.Println("------------checkpoint 1=======")

	existed, update, err := core.UpdateDataBuilder(ctx, a.avatar, id, avatar.Label, fromEnabledDisable, avatar.Enable)
	if err != nil {
		return err
	}

	// update := existed
	if fileHeader != nil {
		url, err := lib.UploadFileToMinio(ctx, a.minio, a.bucketName, fileHeader, string(constants.Avatar), a.minioEndPoint, a.logger)
		if err != nil {
			return err
		}
		existed.Avatar = url
	}
	if avatar.Label != "" {
		update.Label = avatar.Label
	}

	update.LastModifiedAt = time.Now()

	cpsModel := lib.CpsModelBuilder(id, makerData, existed, update, string(constants.RequestUpdateAvatar), constants.UPDATE)
	if err := a.cpsService.CreateCPSAction(ctx, &cpsModel); err != nil {
		return err
	}
	return nil

}
func (a *avatarService) DeleteAvatar(ctx context.Context, id string) error {
	makerData := local_util.ExtractUserFromContext(ctx)
	existing, err := a.avatar.FindByID(ctx, id)
	if err != nil {
		return err
	}

	if existing == nil {
		return errors.New(localization.ErrorAvatarNotExist.Code)
	}

	cpsModel := lib.CpsModelBuilder(id, makerData, existing, nil, string(constants.RequestDeleteAvatar), constants.DELETE)
	if err := a.cpsService.CreateCPSAction(ctx, &cpsModel); err != nil {
		return err
	}
	return nil
}
func (a *avatarService) FetchAllAvatar(ctx context.Context, filterParams types.Filter) (*types.PaginatedResponse[[]*model.Avatar], error) {
	data, err := a.avatar.FindAllWithPagination(ctx, filterParams)
	if err != nil {
		return nil, err
	}

	return data, nil
}
func (a *avatarService) FetchAvatarById(ctx context.Context, id string) (*model.Avatar, error) {
	data, err := a.avatar.FindByID(ctx, id)
	if err != nil {
		return nil, err
	}
	return data, nil
}
func (a *avatarService) Authorize(ctx context.Context, cpsAction *model.CPSAction) (*model.CPSAction, error) {

	avatar, ok := cpsAction.CurrentAction.(model.Avatar)
	if !ok {
		return nil, errors.New(localization.ErrorUnexpectedError.Code)
	}

	switch cpsAction.RequestAction {
	case string(constants.RequestCreateAvatar):
		a.avatar.Create(ctx, &avatar)
	case string(constants.RequestUpdateAvatar):
		a.avatar.Update(ctx, cpsAction.UniqueId, &avatar)
	case string(constants.RequestDeleteAvatar):
		a.avatar.Delete(ctx, avatar.ID.Hex())
	case string(constants.RequestEnableAdvert):
		a.avatar.EnableOrDisable(ctx, avatar.ID.Hex(), true)
	case string(constants.RequestDisableAvatar):
		a.avatar.EnableOrDisable(ctx, avatar.ID.Hex(), false)
	default:
		return nil, errors.New(localization.ErrorInvalidRequest.Code)

	}
	return cpsAction, nil
}
