package budget

import (
	"context"
	"errors"
	"fmt"
	"mime/multipart"
	"net/http"
	"os"
	"path/filepath"
	"time"

	service "github.com/CBE-Super-App/cbe-super-app-cps-action/internal/domain/budget"
	"github.com/CBE-Super-App/cbe-super-app-cps-action/internal/domain/budget/entities"
	constant "github.com/CBE-Super-App/cbe-super-app-cps-action/utils"

	"gitlab.com/bersufekadgetachew/cbe-super-app-shared/shared/config"
	"gitlab.com/bersufekadgetachew/cbe-super-app-shared/shared/utils"
)

type BudgetService interface {
	CreateIcon(ctx context.Context, cpsAction entities.CPSAction) (*entities.CPSAction, error)
	FetchIcons(ctx context.Context) ([]*entities.Icon, error)
	UpdateIcon(ctx context.Context, id string, cpsAction entities.CPSAction) (*entities.CPSAction, error)
	CreateColor(ctx context.Context, hexCode string, cpsAction entities.CPSAction) (*entities.CPSAction, error)
	FetchColors(ctx context.Context) ([]*entities.Color, error)
	UpdateColor(ctx context.Context, id, hexCode string, cpsAction entities.CPSAction) (*entities.CPSAction, error)
	ApproveAction(ctx context.Context, cpsAction entities.CPSAction) (*entities.CPSAction, error)
}

type BudgetHandler struct {
	service    *service.BudgetService
	bucketName string
	minio      config.MinioClientInterface
	logger     utils.Logger
}

func InitBudgetHandler(service *service.BudgetService, minioClinet config.MinioClientInterface, bucketName string, logger utils.Logger) BudgetService {
	return &BudgetHandler{
		service:    service,
		minio:      minioClinet,
		bucketName: bucketName,
		logger:     logger,
	}
}

func (b *BudgetHandler) CreateIcon(ctx context.Context, cpsAction entities.CPSAction) (*entities.CPSAction, error) {
	cpsAction.ActionCode = utils.RandomGenerator(20)
	iconData, ok := cpsAction.CurrentAction.(map[string]interface{})
	if !ok {
		b.logger.Errorf("invalid current action type")
		return nil, fmt.Errorf("invalid currect action")
	}

	imageFileHeader, ok := iconData["icon_url"].(*multipart.FileHeader)
	if !ok {
		b.logger.Errorf("banner image is not present or invalid")
		return nil, fmt.Errorf("banner image missing or invalid")
	}

	file, err := imageFileHeader.Open()
	if err != nil {
		b.logger.Errorf("failed to open uploaded image: %v", err)
		return nil, fmt.Errorf("failed to open uploaded image")
	}
	defer file.Close()

	exist, err := b.minio.BucketExist(ctx, b.bucketName)
	if err != nil {
		b.logger.Errorf("failed to check icon bucket: %v", err)
		return nil, fmt.Errorf("failed to check icon bucket")
	}

	if !exist {
		created, err := b.minio.MakeBucket(ctx, b.bucketName)
		if !created || err != nil {
			b.logger.Errorf("failed to create icon bucket: %v", err)
			return nil, fmt.Errorf("failed to create icon bucker")
		}
	}

	fileName := fmt.Sprintf("budget-icon-%d-%s", time.Now().UnixNano(), imageFileHeader.Filename)
	dir, err := os.Getwd()
	if err != nil {
		b.logger.Errorf("failed to get current working directory", err)
		return nil, fmt.Errorf("failed to create advert bucket")
	}

	filePath := filepath.Join(dir, fileName)

	tempFile, err := os.Create(filePath)
	if err != nil {
		b.logger.Errorf("failed to create temp file: %v", err)
		return nil, fmt.Errorf("failed to create temp file")
	}
	defer func() {
		tempFile.Close()
		os.Remove(filePath)
	}()

	saveObj, err := b.minio.SaveObject(ctx, config.SaveObjectBody{
		BucketName: b.bucketName,
		ObjectName: fileName,
		File:       filePath,
	})
	if err != nil {
		b.logger.Errorf("failed to save object to MinIO: %v", err)
		return nil, fmt.Errorf("failed to save object to MinIo")
	}

	cpsAction.CurrentAction = map[string]interface{}{"icon_url": saveObj.Bucket + "/" + saveObj.Key}

	code, err := b.service.CreateIcon(ctx, cpsAction)
	if err != nil {
		b.logger.Errorf("failed to persist CPS action: %v", err)
		return nil, fmt.Errorf("failed to persist CPS action")
	}

	return code, nil
}

func (b *BudgetHandler) FetchIcons(ctx context.Context) ([]*entities.Icon, error) {
	icons, err := b.service.FetchIcons(ctx)
	if err != nil {
		return nil, err
	}

	return icons, nil
}

func (b *BudgetHandler) UpdateIcon(ctx context.Context, id string, cpsAction entities.CPSAction) (*entities.CPSAction, error) {
	cpsAction.ActionCode = utils.RandomGenerator(20)
	iconData, ok := cpsAction.CurrentAction.(map[string]interface{})
	if !ok {
		b.logger.Errorf("invalid current action type")
		return nil, fmt.Errorf("invalid current action type")
	}

	imageFileHeader, ok := iconData["icon_url"].(*multipart.FileHeader)
	if !ok {
		b.logger.Errorf("banner image is not present or invalid")
		return nil, fmt.Errorf("banner image missing or invalid: %w", constant.ErrorDefinition{
			Code:    http.StatusBadRequest,
			Message: "banner image is required",
		})
	}

	file, err := imageFileHeader.Open()
	if err != nil {
		b.logger.Errorf("failed to open uploaded image: %v", err)
		return nil, fmt.Errorf("failed to open uploaded image: %w", constant.ErrorDefinition{
			Code:    http.StatusBadRequest,
			Message: "invalid image file",
		})
	}
	defer file.Close()

	exist, err := b.minio.BucketExist(ctx, b.bucketName)
	if err != nil {
		b.logger.Errorf("failed to check icon bucket: %v", err)
		return nil, fmt.Errorf("failed to check icon bucket: %w", constant.ErrorDefinition{
			Code:    http.StatusInternalServerError,
			Message: "internal server error",
		})
	}

	if !exist {
		created, err := b.minio.MakeBucket(ctx, b.bucketName)
		if !created || err != nil {
			b.logger.Errorf("failed to create icon bucket: %v", err)
			return nil, fmt.Errorf("failed to create icon bucket: %w", constant.ErrorDefinition{
				Code:    http.StatusInternalServerError,
				Message: "internal server error",
			})
		}
	}

	fileName := fmt.Sprintf("budget-icon-%d-%s", time.Now().UnixNano(), imageFileHeader.Filename)
	dir, err := os.Getwd()
	if err != nil {
		b.logger.Errorf("failed to get current working directory", err)
		return nil, fmt.Errorf("failed to create advert bucket: %w", constant.ErrorDefinition{
			Code:    http.StatusInternalServerError,
			Message: "internal server error",
		})
	}

	filePath := filepath.Join(dir, fileName)

	tempFile, err := os.Create(filePath)
	if err != nil {
		b.logger.Errorf("failed to create temp file: %v", err)
		return nil, fmt.Errorf("failed to create temp file: %w", constant.ErrorDefinition{
			Code:    http.StatusInternalServerError,
			Message: "internal server error",
		})
	}
	defer func() {
		tempFile.Close()
		os.Remove(filePath)
	}()

	saveObj, err := b.minio.SaveObject(ctx, config.SaveObjectBody{
		BucketName: b.bucketName,
		ObjectName: fileName,
		File:       filePath,
	})
	if err != nil {
		b.logger.Errorf("failed to save object to MinIO: %v", err)
		return nil, fmt.Errorf("failed to save object to MinIO: %w", constant.ErrorDefinition{
			Code:    http.StatusInternalServerError,
			Message: "internal server error",
		})
	}

	cpsAction.CurrentAction = map[string]interface{}{"icon_url": saveObj.Bucket + "/" + saveObj.Key}

	action, err := b.service.UpdateIcon(ctx, id, cpsAction)
	if err != nil {
		return nil, err
	}

	return action, nil
}

func (b *BudgetHandler) CreateColor(ctx context.Context, hexCode string, cpsAction entities.CPSAction) (*entities.CPSAction, error) {
	cpsAction.ActionCode = utils.RandomGenerator(20)
	if hexCode == "" {
		return nil, errors.New("hex code cannot be empty")
	}
	action, err := b.service.CreateColor(ctx, hexCode, cpsAction)
	if err != nil {
		return nil, err
	}

	return action, nil
}

func (b *BudgetHandler) FetchColors(ctx context.Context) ([]*entities.Color, error) {
	colors, err := b.service.FetchColors(ctx)
	if err != nil {
		return nil, err
	}

	return colors, nil
}

func (b *BudgetHandler) UpdateColor(ctx context.Context, id, hexCode string, cpsAction entities.CPSAction) (*entities.CPSAction, error) {
	cpsAction.ActionCode = utils.RandomGenerator(20)
	actions, err := b.service.UpdateColor(ctx, id, hexCode, cpsAction)
	if err != nil {
		return nil, err
	}

	return actions, nil
}
func (b *BudgetHandler) ApproveAction(ctx context.Context, cpsAction entities.CPSAction) (*entities.CPSAction, error) {
	cpsActionRes, err := b.service.ApproveAction(ctx, cpsAction)
	if err != nil {
		return nil, err
	}

	return cpsActionRes, nil
}
