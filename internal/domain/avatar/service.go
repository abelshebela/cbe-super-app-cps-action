package avatar

import (
	"context"
	"fmt"
	"net/http"
	"os"
	"path/filepath"
	"time"

	"gitlab.com/bersufekadgetachew/cbe-super-app-cps-action/internal/adapter/outbound/model"
	constant "gitlab.com/bersufekadgetachew/cbe-super-app-cps-action/utils"
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
	CreateAvatar(ctx context.Context, req model.CreateCPSAction) (model.CpsAction, error)
	DeleteAvatar(ctx context.Context, id string, cpsActionReq model.CreateCPSAction) (model.CpsAction, error)
	Authorize(ctx context.Context, req model.AuthorizeCPSAction) (model.CpsAction, error)
	Reject(ctx context.Context, req model.RejectCPSAction) (model.CpsAction, error)
	EnableOrDisableAvatar(ctx context.Context, id string,
		requestAction model.RequestAction, cpsReq model.CreateCPSAction) (*model.CpsAction, error)
	GetAvatar(ctx context.Context, id string) (*Avatar, error)
	GetAllAvatar(ctx context.Context, filterParams constant.Filter) (AvatarResponse, error)
	UpdateAvatar(ctx context.Context, id string, req model.CreateCPSAction) (model.CpsAction, error)
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

func (a *AvatarDomain) CreateAvatar(ctx context.Context, req model.CreateCPSAction) (model.CpsAction, error) {
	req.RequestAction = model.RequestCreateAvatar
	err := a.avatarRepo.CPSActionExists(ctx, req)
	if err != nil {
		return model.CpsAction{}, err
	}

	actionData, ok := req.ActionData.(CreateAvatar)
	if !ok {
		a.logger.Errorf("failed to cast action data to avatar request")
		return model.CpsAction{}, fmt.Errorf("failed to create avatar bucket: %w", constant.ErrorDefinition{
			Code:    http.StatusBadRequest,
			Message: "invalid action data",
		})
	}

	if err := actionData.Validate(); err != nil {
		a.logger.Errorf("validation error", err)
		return model.CpsAction{}, err
	}

	exist, err := a.minioClient.BucketExist(ctx, a.bucketName)
	if err != nil {
		a.logger.Errorf("failed to check wallet bucket: %v", err)
		return model.CpsAction{}, fmt.Errorf("failed to check wallet bucket: %w", constant.ErrorDefinition{
			Code:    http.StatusInternalServerError,
			Message: "internal server error",
		})
	}

	if !exist {
		created, err := a.minioClient.MakeBucket(ctx, a.bucketName)
		if !created || err != nil {
			a.logger.Errorf("failed to create bank bucket: %v", err)
			return model.CpsAction{}, fmt.Errorf("failed to create bank bucket: %w", constant.ErrorDefinition{
				Code:    http.StatusInternalServerError,
				Message: "internal server error",
			})
		}
	}
	baseName := filepath.Base(actionData.Avatar.Filename)
	fileName := fmt.Sprintf("avatar-%d-%s", time.Now().UnixNano(), baseName)

	tempFile, err := os.CreateTemp("", "avatar-*")
	if err != nil {
		a.logger.Errorf("failed to create temp file: %v", err)
		return model.CpsAction{}, fmt.Errorf("failed to create temp file: %w", constant.ErrorDefinition{
			Code:    http.StatusInternalServerError,
			Message: "internal server error",
		})
	}

	defer func() {
		if err := tempFile.Close(); err != nil {
			a.logger.Errorf("failed to close temp file: %v", err)
		}
		if err := os.Remove(tempFile.Name()); err != nil {
			a.logger.Errorf("failed to remove temp file: %v", err)
		}
	}()

	// Save to MinIO
	saveObj, err := a.minioClient.SaveObject(ctx, config.SaveObjectBody{
		BucketName: a.bucketName,
		ObjectName: fileName,
		File:       tempFile.Name(),
	})
	if err != nil {
		a.logger.Errorf("failed to save object to MinIO: %v", err)
		return model.CpsAction{}, fmt.Errorf("upload to MinIO failed: %w", constant.ErrorDefinition{
			Code:    http.StatusInternalServerError,
			Message: "internal server error",
		})
	}

	cpsRes, err := a.avatarRepo.CreateAvatar(ctx, model.CreateCPSAction{
		MakerUser:  req.MakerUser,
		Department: req.Department,
		ActionData: Avatar{
			Label:  actionData.Label,
			Avatar: fmt.Sprintf("%s/%s", saveObj.Bucket, saveObj.Key),
		},
	})
	if err != nil {
		return model.CpsAction{}, err
	}

	return cpsRes, nil
}

func (a *AvatarDomain) DeleteAvatar(ctx context.Context, id string, req model.CreateCPSAction) (model.CpsAction, error) {
	req.RequestAction = model.RequestDeleteAvatar
	err := a.avatarRepo.CPSActionExists(ctx, req)
	if err != nil {
		return model.CpsAction{}, err
	}
	cpsAction, err := a.avatarRepo.DeleteAvatar(ctx, id, req)
	if err != nil {
		return model.CpsAction{}, err
	}

	return cpsAction, nil
}

func (a *AvatarDomain) Authorize(ctx context.Context, req model.AuthorizeCPSAction) (model.CpsAction, error) {
	cpsAction, err := a.avatarRepo.Authorize(ctx, req)
	if err != nil {
		return model.CpsAction{}, err
	}

	return cpsAction, nil
}

func (a *AvatarDomain) Reject(ctx context.Context, req model.RejectCPSAction) (model.CpsAction, error) {
	if err := req.Validate(); err != nil {
		a.logger.Errorf("validation error", err)
		return model.CpsAction{}, err
	}
	cpsAction, err := a.avatarRepo.Reject(ctx, req)
	if err != nil {
		return model.CpsAction{}, err
	}

	return cpsAction, nil
}

func (a *AvatarDomain) EnableOrDisableAvatar(ctx context.Context, id string,
	requestAction model.RequestAction, cpsReq model.CreateCPSAction) (*model.CpsAction, error) {
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

func (a *AvatarDomain) UpdateAvatar(ctx context.Context, id string, req model.CreateCPSAction) (model.CpsAction, error) {
	req.RequestAction = model.RequestUpdateWallet
	err := a.avatarRepo.CPSActionExists(ctx, req)
	if err != nil {
		return model.CpsAction{}, err
	}

	actionData, ok := req.ActionData.(UpdateAvatar)
	if !ok {
		a.logger.Errorf("failed to cast action data to avatar request")
		return model.CpsAction{}, fmt.Errorf("failed to create avatar bucket: %w", constant.ErrorDefinition{
			Code:    http.StatusBadRequest,
			Message: "invalid action data",
		})
	}

	if err := actionData.Validate(); err != nil {
		a.logger.Errorf("validation error", err)
		return model.CpsAction{}, err
	}

	exist, err := a.minioClient.BucketExist(ctx, a.bucketName)
	if err != nil {
		a.logger.Errorf("failed to check avatar bucket: %v", err)
		return model.CpsAction{}, fmt.Errorf("failed to check avatar bucket: %w", constant.ErrorDefinition{
			Code:    http.StatusInternalServerError,
			Message: "internal server error",
		})
	}

	if !exist {
		a.logger.Errorf("avatar bucket not found: %v", err)
		return model.CpsAction{}, fmt.Errorf("failed to check avatar bucket: %w", constant.ErrorDefinition{
			Code:    http.StatusNotFound,
			Message: "bucket not exist",
		})
	}

	avatar, err := a.avatarRepo.GetAvatar(ctx, id)
	if err != nil {
		return model.CpsAction{}, err
	}

	fileName := fmt.Sprintf("avatar-%d-%s", time.Now().UnixNano(), actionData.Avatar.Filename)
	dir, err := os.Getwd()
	if err != nil {
		a.logger.Errorf("failed to get current working directory", err)
		return model.CpsAction{}, fmt.Errorf("failed to create avatar bucket: %w", constant.ErrorDefinition{
			Code:    http.StatusInternalServerError,
			Message: "internal server error",
		})
	}
	filePath := filepath.Join(dir, fileName)

	tempFile, err := os.Create(filePath)
	if err != nil {
		a.logger.Errorf("failed to create temp file: %v", err)
		return model.CpsAction{}, fmt.Errorf("failed to create temp file: %w", constant.ErrorDefinition{
			Code:    http.StatusInternalServerError,
			Message: "internal server error",
		})
	}
	defer func() {
		tempFile.Close()
		os.Remove(filePath)
	}()

	saveObj, err := a.minioClient.SaveObject(ctx, config.SaveObjectBody{
		BucketName: a.bucketName,
		ObjectName: fileName,
		File:       filePath,
	})
	if err != nil {
		a.logger.Errorf("failed to save object to MinIO: %v", err)
		return model.CpsAction{}, fmt.Errorf("failed to save object to MinIO: %w", constant.ErrorDefinition{
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
		return model.CpsAction{}, err
	}

	return cpsAction, nil
}
