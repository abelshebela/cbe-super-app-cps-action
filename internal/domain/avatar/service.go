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

	// Save to MinIO
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

	cpsRes, err := a.avatarRepo.CreateAvatar(ctx, model.CreateCPSAction{
		MakerUser:  req.MakerUser,
		Department: req.Department,
		ActionData: Avatar{
			Label:     actionData.Label,
			AvatarURL: fmt.Sprintf("%s/%s", saveObj.Bucket, saveObj.Key),
		},
	})

	if err != nil {
		return model.CpsAction{}, err
	}

	return cpsRes, nil
}
