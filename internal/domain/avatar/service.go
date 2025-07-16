package avatar

import (
	"context"
	"fmt"
	"net/http"
	"time"

	"github.com/CBE-Super-App/cbe-super-app-cps-action/internal/adapter/outbound/model"
	constant "github.com/CBE-Super-App/cbe-super-app-cps-action/utils"
	"go.mongodb.org/mongo-driver/v2/bson"

	"gitlab.com/bersufekadgetachew/cbe-super-app-shared/shared/config"
	"gitlab.com/bersufekadgetachew/cbe-super-app-shared/shared/utils"
)

type AvatarDomain struct {
	avatarRepo  AvatarRepository
	bucketName  string
	minioClient config.MinioClientInterface
	logger      utils.Logger
}

type AvatarDomainService interface {
	CreateAvatar(ctx context.Context, req model.CreateCPSAction) (model.CPSAction, error)
	DeleteAvatar(ctx context.Context, id string, cpsActionReq model.CreateCPSAction) (model.CPSAction, error)
	Authorize(ctx context.Context, req model.AuthorizeCPSAction) (model.CPSAction, error)
	Reject(ctx context.Context, req model.RejectCPSAction) (model.CPSAction, error)
	EnableOrDisableAvatar(ctx context.Context, id string,
		requestAction model.RequestAction, cpsReq model.CreateCPSAction) (*model.CPSAction, error)
	GetAvatar(ctx context.Context, id string) (*Avatar, error)
	GetAllAvatar(ctx context.Context, filterParams constant.Filter) (AvatarResponse, error)
	UpdateAvatar(ctx context.Context, id string, req model.CreateCPSAction) (model.CPSAction, error)
}

func InitAvatarDomain(avatarRepo AvatarRepository,
	minioclient config.MinioClientInterface, bucketName string, logger utils.Logger) AvatarDomainService {
	return &AvatarDomain{
		avatarRepo:  avatarRepo,
		bucketName:  bucketName,
		minioClient: minioclient,
		logger:      logger,
	}
}

func (a *AvatarDomain) CreateAvatar(ctx context.Context, req model.CreateCPSAction) (model.CPSAction, error) {

	req.RequestAction = model.RequestCreateAvatar
	err := a.avatarRepo.CPSActionExists(ctx, req)
	if err != nil {
		return model.CPSAction{}, err
	}

	actionData, ok := req.ActionData.(CreateAvatar)

	if !ok {
		a.logger.Errorf("failed to cast action data to avatar request")
		return model.CPSAction{}, fmt.Errorf("failed to create avatar bucket: %w", constant.ErrorDefinition{
			Code:    http.StatusBadRequest,
			Message: "invalid action data",
		})
	}

	if err := actionData.Validate(); err != nil {
		a.logger.Errorf("validation error", err)
		return model.CPSAction{}, err
	}

	exist, err := a.minioClient.BucketExist(ctx, a.bucketName)
	if err != nil {
		a.logger.Errorf("failed to check avatar bucket: %v", err)
		return model.CPSAction{}, fmt.Errorf("failed to check avatar bucket: %w", constant.ErrorDefinition{
			Code:    http.StatusInternalServerError,
			Message: "internal server error",
		})
	}

	if !exist {
		created, err := a.minioClient.MakeBucket(ctx, a.bucketName)
		if !created || err != nil {
			a.logger.Errorf("failed to create avatar bucket: %v", err)
			return model.CPSAction{}, fmt.Errorf("failed to create avatar bucket: %w", constant.ErrorDefinition{
				Code:    http.StatusInternalServerError,
				Message: "internal server error",
			})
		}
	}

	fileName := fmt.Sprintf("avatar-%d-%s", time.Now().UnixNano(), actionData.Avatar.Filename)
	file, err := actionData.Avatar.Open()
	if err != nil {
		a.logger.Errorf("failed to open uploaded file: %v", err)
		return model.CPSAction{}, fmt.Errorf("failed to open uploaded file: %w", constant.ErrorDefinition{
			Code:    http.StatusInternalServerError,
			Message: "internal server error",
		})
	}
	defer file.Close()

	saveObj, err := a.minioClient.SaveObjectN(ctx, config.SaveObjectBodyN{
		BucketName:  a.bucketName,
		ObjectName:  fileName,
		Reader:      file,
		Size:        actionData.Avatar.Size,
		ContentType: config.ContentType(actionData.Avatar.Header.Get("Content-Type")),
	})

	if err != nil {
		a.logger.Errorf("failed to save object to MinIO: %v", err)
		return model.CPSAction{}, fmt.Errorf("upload to MinIO failed: %w", constant.ErrorDefinition{
			Code:    http.StatusInternalServerError,
			Message: "internal server error",
		})
	}

	cpsRes, err := a.avatarRepo.CreateAvatar(ctx, model.CreateCPSAction{
		MakerUser:  req.MakerUser,
		Department: req.Department,
		ActionData: Avatar{
			ID:             bson.NewObjectID().Hex(),
			Label:          actionData.Label,
			Avatar:         fmt.Sprintf("%s/%s", saveObj.Bucket, saveObj.Key),
			CreatedAt:      time.Now(),
			LastModifiedAt: time.Now(),
		},
	})
	if err != nil {
		return model.CPSAction{}, err
	}

	return cpsRes, nil
}

func (a *AvatarDomain) DeleteAvatar(ctx context.Context, id string, req model.CreateCPSAction) (model.CPSAction, error) {
	req.RequestAction = model.RequestDeleteAvatar
	err := a.avatarRepo.CPSActionExists(ctx, req)
	if err != nil {
		return model.CPSAction{}, err
	}
	cpsAction, err := a.avatarRepo.DeleteAvatar(ctx, id, req)
	if err != nil {
		return model.CPSAction{}, err
	}

	return cpsAction, nil
}

func (a *AvatarDomain) Authorize(ctx context.Context, req model.AuthorizeCPSAction) (model.CPSAction, error) {
	cpsAction, err := a.avatarRepo.Authorize(ctx, req)
	if err != nil {
		return model.CPSAction{}, err
	}

	return cpsAction, nil
}

func (a *AvatarDomain) Reject(ctx context.Context, req model.RejectCPSAction) (model.CPSAction, error) {
	if err := req.Validate(); err != nil {
		a.logger.Errorf("validation error", err)
		return model.CPSAction{}, err
	}
	cpsAction, err := a.avatarRepo.Reject(ctx, req)
	if err != nil {
		return model.CPSAction{}, err
	}

	return cpsAction, nil
}

func (a *AvatarDomain) EnableOrDisableAvatar(ctx context.Context, id string,
	requestAction model.RequestAction, cpsReq model.CreateCPSAction) (*model.CPSAction, error) {
	cpsReq.RequestAction = requestAction
	err := a.avatarRepo.CPSActionExists(ctx, cpsReq)
	if err != nil {
		return nil, err
	}
	cpsAction, err := a.avatarRepo.EnableOrDisableAvatar(ctx, id, requestAction, cpsReq)
	if err != nil {
		return nil, err
	}

	return cpsAction, nil
}

func (a *AvatarDomain) GetAllAvatar(ctx context.Context, filterParams constant.Filter) (AvatarResponse, error) {
	return a.avatarRepo.GetAllAvatar(ctx, filterParams)
}

func (a *AvatarDomain) GetAvatar(ctx context.Context, id string) (*Avatar, error) {
	return a.avatarRepo.GetAvatar(ctx, id)
}

func (a *AvatarDomain) UpdateAvatar(ctx context.Context, id string, req model.CreateCPSAction) (model.CPSAction, error) {
	req.RequestAction = model.RequestUpdateAvatar
	err := a.avatarRepo.CPSActionExists(ctx, req)
	if err != nil {
		return model.CPSAction{}, err
	}

	actionData, ok := req.ActionData.(UpdateAvatar)
	if !ok {
		a.logger.Errorf("failed to cast action data to avatar request")
		return model.CPSAction{}, fmt.Errorf("failed to create avatar bucket: %w", constant.ErrorDefinition{
			Code:    http.StatusBadRequest,
			Message: "invalid action data",
		})
	}

	if err := actionData.Validate(); err != nil {
		a.logger.Errorf("validation error", err)
		return model.CPSAction{}, err
	}

	exist, err := a.minioClient.BucketExist(ctx, a.bucketName)
	if err != nil {
		a.logger.Errorf("failed to check avatar bucket %v", err)
		return model.CPSAction{}, fmt.Errorf("failed to check avatar bucket: %w", constant.ErrorDefinition{
			Code:    http.StatusInternalServerError,
			Message: "internal server error",
		})
	}

	if !exist {
		a.logger.Errorf("avatar bucket not found")
		return model.CPSAction{}, fmt.Errorf("failed to check avatar bucket: %w", constant.ErrorDefinition{
			Code:    http.StatusNotFound,
			Message: "bucket not exist",
		})
	}

	avatar, err := a.avatarRepo.GetAvatar(ctx, id)
	if err != nil {
		return model.CPSAction{}, err
	}

	fileName := fmt.Sprintf("avatar-%d-%s", time.Now().UnixNano(), actionData.Avatar.Filename)
	file, err := actionData.Avatar.Open()
	if err != nil {
		a.logger.Errorf("failed to open uploaded file: %v", err)
		return model.CPSAction{}, fmt.Errorf("failed to open uploaded file: %w", constant.ErrorDefinition{
			Code:    http.StatusInternalServerError,
			Message: "internal server error",
		})
	}
	defer file.Close()

	saveObj, err := a.minioClient.SaveObjectN(ctx, config.SaveObjectBodyN{
		BucketName:  a.bucketName,
		ObjectName:  fileName,
		Reader:      file,
		Size:        actionData.Avatar.Size,
		ContentType: config.ContentType(actionData.Avatar.Header.Get("Content-Type")),
	})

	if err != nil {
		a.logger.Errorf("failed to save object to MinIO: %v", err)
		return model.CPSAction{}, fmt.Errorf("upload to MinIO failed: %w", constant.ErrorDefinition{
			Code:    http.StatusInternalServerError,
			Message: "internal server error",
		})
	}

	req.ActionData = Avatar{
		ID:     avatar.ID,
		Avatar: fmt.Sprintf("%s/%s", saveObj.Bucket, saveObj.Key),
	}
	cpsAction, err := a.avatarRepo.UpdateAvatar(ctx, id, req)
	if err != nil {
		return model.CPSAction{}, err
	}

	return cpsAction, nil
}
