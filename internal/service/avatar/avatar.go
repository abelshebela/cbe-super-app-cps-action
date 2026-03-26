package avatar

import (
	"cbe-super-app-cps-action/internal/constants"
	"cbe-super-app-cps-action/internal/constants/lib"
	"cbe-super-app-cps-action/internal/constants/localization"
	"cbe-super-app-cps-action/internal/constants/types"
	"cbe-super-app-cps-action/internal/service"
	"cbe-super-app-cps-action/internal/service/avatar/core"
	"cbe-super-app-cps-action/internal/storage"
	local_util "cbe-super-app-cps-action/pkgs/utils"
	"context"
	"errors"
	"mime/multipart"
	"path"
	"time"

	"gitlab.com/bersufekadgetachew/cbe-super-app-shared/shared/entities/model"

	"github.com/aws/aws-sdk-go-v2/service/s3"
	"gitlab.com/bersufekadgetachew/cbe-super-app-shared/shared/config"
	"gitlab.com/bersufekadgetachew/cbe-super-app-shared/shared/utils"
	"go.mongodb.org/mongo-driver/v2/bson"
	"go.opentelemetry.io/otel/attribute"
	"go.opentelemetry.io/otel/trace"
)

type avatarService struct {
	logger        utils.Logger
	avatar        storage.AvatarRepository
	minio         *s3.Client
	cpsService    service.CPSActionService
	bucketName    string
	minioEndPoint string
	cfg           config.VaultConfig
}

func NewAvatarService(avatar storage.AvatarRepository, cpsService service.CPSActionService, logger utils.Logger, minio *s3.Client, buckateName string, minioEndPoint string, cfg config.VaultConfig) service.AvatarService {
	return &avatarService{
		cpsService:    cpsService,
		avatar:        avatar,
		logger:        logger,
		minio:         minio,
		bucketName:    buckateName,
		minioEndPoint: minioEndPoint,
		cfg:           cfg,
	}
}

func (a *avatarService) CreateAvatar(ctx context.Context, avatar *model.Avatar, fileHeader *multipart.FileHeader) error {
	ctx, span := local_util.TraceLogger(ctx, "service", "CreateAvatar", "Avatar", "CreateAvatar")
	defer span.End()

	a.logger.Infof("[AvatarSvc][Create] creating")
	makerData := local_util.ExtractUserFromContext(ctx)
	existing, _ := a.avatar.Find(ctx, bson.M{"label": avatar.Label}, nil)

	if existing != nil {
		span.AddEvent("[CreateAvatar] avatar with label already exists", trace.WithAttributes(attribute.String("label", avatar.Label)))
		a.logger.Errorf("[AvatarSvc][Create] label exists")
		return errors.New(localization.ErrorAvatarAlreadyExist.Code)
	}

	url, err := lib.UploadFileToMinio(ctx, a.minio, a.bucketName, fileHeader, string(constants.AvatarFolderName), a.cfg, "", a.logger)
	if err != nil {
		span.AddEvent("[CreateAvatar] failed to upload file", trace.WithAttributes(
			attribute.String("error", err.Error()),
			attribute.String("label", avatar.Label),
		))
		return err
	}
	newAvatar := model.Avatar{Avatar: url, Label: avatar.Label, CreatedAt: time.Now(), Enable: true}
	cpsModel := lib.CpsModelBuilder("", makerData, nil, newAvatar, string(constants.RequestCreateAvatar), constants.CREATE)

	if err := a.cpsService.CreateCPSAction(ctx, &cpsModel); err != nil {
		span.AddEvent("[CreateAvatar] failed to create CPS action", trace.WithAttributes(
			attribute.String("error", err.Error()),
			attribute.String("label", avatar.Label),
		))
		return err
	}
	return nil
}
func (a *avatarService) UpdateAvatar(ctx context.Context, id string, avatar *model.Avatar, fileHeader *multipart.FileHeader, fromEnabledDisable bool) error {
	ctx, span := local_util.TraceLogger(ctx, "service", "UpdateAvatar", "Avatar", "UpdateAvatar")
	defer span.End()

	a.logger.Infof("[AvatarSvc][Update] id: %s", id)
	makerData := local_util.ExtractUserFromContext(ctx)
	if incomplete := local_util.IsIncomplete(makerData); incomplete {
		span.AddEvent("[UpdateAvatar] incomplete user data", trace.WithAttributes(attribute.String("id", id)))
		a.logger.Errorf("[AvatarSvc][Update] incomplete user")
		return errors.New(localization.ErrorIncompleteUserInfo.Code)
	}

	existing, err := a.avatar.Find(ctx, bson.M{"label": avatar.Label}, nil)
	if err != nil {
		if err.Error() != localization.ErrorFileNotFound.Code {
			span.AddEvent("[UpdateAvatar] failed to find avatar", trace.WithAttributes(
				attribute.String("error", err.Error()),
				attribute.String("id", id),
			))
			a.logger.Errorf("[AvatarSvc][Update] find err: %v", err)
			return errors.New(localization.ErrorAvatarAlreadyExist.Code)
		}
	}

	if existing == nil {
		span.AddEvent("[UpdateAvatar] avatar with label already exists", trace.WithAttributes(attribute.String("label", avatar.Label)))
		a.logger.Errorf("[AvatarSvc][Update] label exists")
		return errors.New(localization.ErrorAvatarAlreadyExist.Code)
	}

	existed, update, err := core.UpdateDataBuilder(ctx, a.avatar, id, avatar.Label, fromEnabledDisable, avatar.Enable)
	if err != nil {
		span.AddEvent("[UpdateAvatar] failed to build update data", trace.WithAttributes(
			attribute.String("error", err.Error()),
			attribute.String("id", id),
		))
		a.logger.Errorf("[AvatarSvc][Update] build data err: %v", err)
		return err
	}

	// update := existed
	if fileHeader != nil {
		var objectkey string
		if existed.Avatar != "" {
			objectkey = path.Base(existed.Avatar)
		}

		url, err := lib.UploadFileToMinio(ctx, a.minio, a.bucketName, fileHeader, string(constants.AvatarFolderName), a.cfg, objectkey, a.logger)
		if err != nil {
			span.AddEvent("[UpdateAvatar] failed to upload avatar image", trace.WithAttributes(
				attribute.String("error", err.Error()),
				attribute.String("id", id),
			))
			a.logger.Errorf("[AvatarSvc][Update] upload err: %v", err)
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
		span.AddEvent("[UpdateAvatar] failed to create CPS action", trace.WithAttributes(
			attribute.String("error", err.Error()),
			attribute.String("id", id),
		))
		a.logger.Errorf("[AvatarSvc][Update] cps action err: %v", err)
		return err
	}
	a.logger.Infof("[AvatarSvc][Update] request created id: %s", id)
	return nil

}

func (a *avatarService) EnableDisable(ctx context.Context, id string, enable bool) error {
	ctx, span := local_util.TraceLogger(ctx, "service", "EnableDisable", "Avatar", "EnableDisable")
	defer span.End()

	a.logger.Infof("[AvatarSvc][EnableDisable] id: %s enabled: %v", id, enable)
	makerData := local_util.ExtractUserFromContext(ctx)
	if incomplete := local_util.IsIncomplete(makerData); incomplete {
		span.AddEvent("[EnableDisable] incomplete user data", trace.WithAttributes(attribute.String("id", id)))
		a.logger.Errorf("[AvatarSvc][EnableDisable] incomplete user")
		return errors.New(localization.ErrorIncompleteUserInfo.Code)
	}
	data, err := a.avatar.FindByID(ctx, id)
	if err != nil {
		span.AddEvent("[EnableDisable] failed to find avatar", trace.WithAttributes(
			attribute.String("error", err.Error()),
			attribute.String("id", id),
		))
		a.logger.Errorf("[AvatarSvc][EnableDisable] find err: %v", err)
		return err
	}

	if data.Enable == enable {
		if enable {
			span.AddEvent("[EnableDisable] avatar already enabled", trace.WithAttributes(attribute.String("id", id)))
			a.logger.Errorf("[AvatarSvc][EnableDisable] already enabled")
			return errors.New(localization.ErrorAvatarAlreadyEnabled.Code)
		}
		span.AddEvent("[EnableDisable] avatar already disabled", trace.WithAttributes(attribute.String("id", id)))
		a.logger.Errorf("[AvatarSvc][EnableDisable] already disabled")
		return errors.New(localization.ErrorAvatarAlreadyDisabled.Code)
	}
	var requestAction string
	if enable {
		requestAction = string(constants.RequestEnableAvatar)
	} else {
		requestAction = string(constants.RequestDisableAvatar)
	}
	update := data

	update.Enable = enable
	cpsModel := lib.CpsModelBuilder(id, makerData, data, update, requestAction, string(constants.UPDATE))
	if err := a.cpsService.CreateCPSAction(ctx, &cpsModel); err != nil {
		span.AddEvent("[EnableDisable] failed to create CPS action", trace.WithAttributes(
			attribute.String("error", err.Error()),
			attribute.String("id", id),
		))
		a.logger.Errorf("[AvatarSvc][EnableDisable] cps action err: %v", err)
		return err
	}
	a.logger.Infof("[AvatarSvc][EnableDisable] request created id: %s", id)
	return nil
}
func (a *avatarService) DeleteAvatar(ctx context.Context, id string) error {
	ctx, span := local_util.TraceLogger(ctx, "service", "DeleteAvatar", "Avatar", "DeleteAvatar")
	defer span.End()

	a.logger.Infof("[AvatarSvc][Delete] id: %s", id)
	makerData := local_util.ExtractUserFromContext(ctx)
	existing, err := a.avatar.FindByID(ctx, id)
	if err != nil {
		span.AddEvent("[DeleteAvatar] failed to find avatar", trace.WithAttributes(
			attribute.String("error", err.Error()),
			attribute.String("id", id),
		))
		a.logger.Errorf("[AvatarSvc][Delete] find err: %v", err)
		return err
	}

	if existing == nil {
		span.AddEvent("[DeleteAvatar] avatar not found", trace.WithAttributes(attribute.String("id", id)))
		a.logger.Errorf("[AvatarSvc][Delete] not found: %s", id)
		return errors.New(localization.ErrorAvatarNotExist.Code)
	}

	cpsModel := lib.CpsModelBuilder(id, makerData, existing, nil, string(constants.RequestDeleteAvatar), constants.DELETE)
	if err := a.cpsService.CreateCPSAction(ctx, &cpsModel); err != nil {
		span.AddEvent("[DeleteAvatar] failed to create CPS action", trace.WithAttributes(
			attribute.String("error", err.Error()),
			attribute.String("id", id),
		))
		a.logger.Errorf("[AvatarSvc][Delete] cps action err: %v", err)
		return err
	}
	a.logger.Infof("[AvatarSvc][Delete] request created id: %s", id)
	return nil
}
func (a *avatarService) FetchAllAvatar(ctx context.Context, filterParams types.Filter) (*types.PaginatedResponse[[]model.Avatar], error) {
	ctx, span := local_util.TraceLogger(ctx, "service", "FetchAllAvatar", "Avatar", "FetchAllAvatar")
	defer span.End()

	data, err := a.avatar.FindAllWithPagination(ctx, filterParams)
	if err != nil {
		span.AddEvent("[FetchAllAvatar] failed to fetch avatars", trace.WithAttributes(
			attribute.String("error", err.Error()),
		))
		a.logger.Errorf("[AvatarSvc][FetchAll] fetch err: %v", err)
		return nil, err
	}
	a.logger.Infof("[AvatarSvc][FetchAll] count: %d", len(data.Data))
	return data, nil
}
func (a *avatarService) FetchAvatarById(ctx context.Context, id string) (*model.Avatar, error) {
	ctx, span := local_util.TraceLogger(ctx, "service", "FetchAvatarById", "Avatar", "FetchAvatarById")
	defer span.End()

	data, err := a.avatar.FindByID(ctx, id)
	if err != nil {
		span.AddEvent("[FetchAvatarById] failed to fetch avatar", trace.WithAttributes(
			attribute.String("error", err.Error()),
			attribute.String("id", id),
		))
		a.logger.Errorf("[AvatarSvc][FetchByID] fetch err: %v", err)
		return nil, err
	}
	a.logger.Infof("[AvatarSvc][FetchByID] found id: %s", id)
	return data, nil
}
func (a *avatarService) Authorize(ctx context.Context, cpsAction *model.CPSAction) (*model.CPSAction, error) {
	ctx, span := local_util.TraceLogger(ctx, "service", "Authorize", "Avatar", "Authorize")
	defer span.End()

	a.logger.Infof("[AvatarSvc][Authorize] action: %s", cpsAction.RequestAction)

	avatar, err := local_util.JsonUnmarshal[model.Avatar](cpsAction.CurrentAction)
	if err != nil {
		span.AddEvent("[Authorize] failed to unmarshal current action", trace.WithAttributes(
			attribute.String("error", err.Error()),
			attribute.String("unique_id", cpsAction.UniqueId),
		))
		a.logger.Errorf("[AvatarSvc][Authorize] unmarshal err: %v", err)
		return nil, errors.New(localization.ErrorUnexpectedError.Code)
	}

	switch cpsAction.RequestAction {
	case string(constants.RequestCreateAvatar):
		err = a.avatar.Create(ctx, &model.Avatar{Avatar: avatar.Avatar, Label: avatar.Label, CreatedAt: avatar.CreatedAt, Enable: true})
		if err != nil {
			span.AddEvent("[Authorize] failed to create avatar", trace.WithAttributes(
				attribute.String("error", err.Error()),
				attribute.String("unique_id", cpsAction.UniqueId),
			))
			a.logger.Errorf("[AvatarSvc][Authorize] create err: %v", err)
			return nil, err
		}
		a.logger.Infof("[AvatarSvc][Authorize] created")
	case string(constants.RequestUpdateAvatar):
		err = a.avatar.Update(ctx, cpsAction.UniqueId, &model.Avatar{Avatar: avatar.Avatar, Label: avatar.Label, Enable: avatar.Enable})
		if err != nil {
			span.AddEvent("[Authorize] failed to update avatar", trace.WithAttributes(
				attribute.String("error", err.Error()),
				attribute.String("unique_id", cpsAction.UniqueId),
			))
			a.logger.Errorf("[AvatarSvc][Authorize] update err: %v", err)
			return nil, err
		}
		a.logger.Infof("[AvatarSvc][Authorize] updated id: %s", cpsAction.UniqueId)
	case string(constants.RequestDeleteAvatar):
		err = a.avatar.Delete(ctx, cpsAction.UniqueId)
		if err != nil {
			span.AddEvent("[Authorize] failed to delete avatar", trace.WithAttributes(
				attribute.String("error", err.Error()),
				attribute.String("unique_id", cpsAction.UniqueId),
			))
			a.logger.Errorf("[AvatarSvc][Authorize] delete err: %v", err)
			return nil, err
		}
		a.logger.Infof("[AvatarSvc][Authorize] deleted id: %s", cpsAction.UniqueId)
	case string(constants.RequestEnableAvatar):
		err = a.avatar.EnableOrDisable(ctx, cpsAction.UniqueId, true)
		if err != nil {
			span.AddEvent("[Authorize] failed to enable avatar", trace.WithAttributes(
				attribute.String("error", err.Error()),
				attribute.String("unique_id", cpsAction.UniqueId),
			))
			a.logger.Errorf("[AvatarSvc][Authorize] enable err: %v", err)
			return nil, err
		}
		a.logger.Infof("[AvatarSvc][Authorize] enabled id: %s", cpsAction.UniqueId)
	case string(constants.RequestDisableAvatar):
		err = a.avatar.EnableOrDisable(ctx, cpsAction.UniqueId, false)
		if err != nil {
			span.AddEvent("[Authorize] failed to disable avatar", trace.WithAttributes(
				attribute.String("error", err.Error()),
				attribute.String("unique_id", cpsAction.UniqueId),
			))
			a.logger.Errorf("[AvatarSvc][Authorize] disable err: %v", err)
			return nil, err
		}
		a.logger.Infof("[AvatarSvc][Authorize] disabled id: %s", cpsAction.UniqueId)
	default:
		span.AddEvent("[Authorize] unsupported action", trace.WithAttributes(attribute.String("action", cpsAction.RequestAction)))
		a.logger.Errorf("[AvatarSvc][Authorize] unsupported: %s", cpsAction.RequestAction)
		return nil, errors.New(localization.ErrorInvalidRequest.Code)
	}
	a.logger.Infof("[AvatarSvc][Authorize] done: %s", cpsAction.RequestAction)
	return cpsAction, nil
}
