package service

import (
	"context"
	"fmt"
	"strings"
	"time"

	"github.com/CBE-Super-App/cbe-super-app-cps-action/internal/adapter/outbound/model"
	"github.com/CBE-Super-App/cbe-super-app-cps-action/internal/domain/ad/entity"
	"github.com/CBE-Super-App/cbe-super-app-cps-action/internal/domain/ad/repository"
	constant "github.com/CBE-Super-App/cbe-super-app-cps-action/utils"

	"github.com/CBE-Super-App/cbe-super-app-cps-action/internal/domain/ad/dto"
	// "github.com/mitchellh/mapstructure"
	"gitlab.com/bersufekadgetachew/cbe-super-app-shared/shared/config"
	"gitlab.com/bersufekadgetachew/cbe-super-app-shared/shared/utils"
)

type ADDomain struct {
	ADRepo      repository.Repository
	bucketName  string
	minioClient config.MinioClientInterface
	logger      utils.Logger
}

type AdvertService interface {
	CreateOneAdvert(ctx context.Context, cpsAction model.CreateCPSAction) (*model.CPSAction, error)
	UpdateOneAdvert(ctx context.Context, id string, cpsAction model.CreateCPSAction) (*model.CPSAction, error)
	GetAllAdvert(ctx context.Context, filterParams *constant.Filter) (*entity.AdvertResponse, error)
	GetOneAdvert(ctx context.Context, id string) (*entity.Advert, error)
	DeleteOneAdvert(ctx context.Context, id string, cpsAction model.CreateCPSAction) (*model.CPSAction, error)
	Authorize(ctx context.Context, cpsAction model.AuthorizeCPSAction) (*model.CPSAction, error)
	Reject(ctx context.Context, cpsAction model.RejectCPSAction) (*model.CPSAction, error)
}

var _ AdvertService = (*ADDomain)(nil)

func stripFieldPrefix(err error) string {
	// Removes "field: message." format from ozzo-validation errors
	// Example: "logo: INVALID_FILE_TYPE." → "INVALID_FILE_TYPE"
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

func InitADDomian(bucketName string, minioClient config.MinioClientInterface, adRepo repository.Repository, logger utils.Logger) AdvertService {
	return &ADDomain{
		ADRepo:      adRepo,
		bucketName:  bucketName,
		minioClient: minioClient,
		logger:      logger,
	}
}

func (a *ADDomain) CreateOneAdvert(ctx context.Context, cpsAction model.CreateCPSAction) (*model.CPSAction, error) {

	actionData, ok := cpsAction.ActionData.(dto.CreateAdvertRequest)

	if !ok {
		a.logger.Errorf("failed to cast action data to ad request")
		return nil, fmt.Errorf("FAILED_TO_CAST_ACTION_DATA")
	}

	if err := actionData.Validate(); err != nil {
		a.logger.Errorf("validation error", err)
		return nil, err
	}

	exist, err := a.minioClient.BucketExist(ctx, a.bucketName)
	if err != nil {
		a.logger.Errorf("failed to check ad bucket: %v", err)
		return nil, fmt.Errorf("FAILED_TO_CHECK_AD_BUCKET")
	}

	if !exist {

		created, err := a.minioClient.MakeBucket(ctx, a.bucketName)
		if !created || err != nil {
			a.logger.Errorf("failed to create ad bucket: %v", err)
			return nil, fmt.Errorf("UNHANDLED_SERVER_ERROR")
		}
	}

	fileName := fmt.Sprintf("upload-%d-%s", time.Now().UnixNano(), actionData.BannerImage.Filename)

	file, err := actionData.BannerImage.Open()
	if err != nil {
		a.logger.Errorf("failed to open file", err)
		return nil, fmt.Errorf("FAILED_TO_OPEN_FILE")
	}
	defer file.Close()
	// Save to MinIO
	saveObj, err := a.minioClient.SaveObjectN(ctx, config.SaveObjectBodyN{
		BucketName:  a.bucketName,
		ObjectName:  fileName,
		Reader:      file,
		Size:        actionData.BannerImage.Size,
		ContentType: config.ContentType(actionData.BannerImage.Header.Get("Content-Type")),
	})
	if err != nil {
		a.logger.Errorf("failed to save object to MinIO: %v", err)
		return nil, fmt.Errorf("UNHANDLED_SERVER_ERROR")
	}

	cpsActionRes, err := a.ADRepo.CreateOneAdvert(ctx, model.CreateCPSAction{
		MakerUser:  cpsAction.MakerUser,
		Department: cpsAction.Department,
		ActionData: entity.Advert{
			Title:       actionData.Title,
			Description: actionData.Description,
			BannerImage: fmt.Sprintf("%s/%s", saveObj.Bucket, saveObj.Key),
			AdvertFor:   entity.AdvertFor(actionData.AdvertFor),
			Date:        entity.AdvertDate(actionData.Date),
		},
		// Map other fields from cpsAction as needed
	})
	if err != nil {
		return nil, err
	}

	return cpsActionRes, nil
}

func (a *ADDomain) DeleteOneAdvert(ctx context.Context, id string, cpsAction model.CreateCPSAction) (*model.CPSAction, error) {
	cpsRes, err := a.ADRepo.DeleteOneAdvert(ctx, id, cpsAction)
	if err != nil {
		return nil, err
	}
	return cpsRes, nil
}

func (a *ADDomain) GetAllAdvert(ctx context.Context, filterParams *constant.Filter) (*entity.AdvertResponse, error) {
	advertRes, err := a.ADRepo.GetAllAdvert(ctx, filterParams)
	if err != nil {
		return nil, err
	}

	return advertRes, nil
}

func (a *ADDomain) GetOneAdvert(ctx context.Context, id string) (*entity.Advert, error) {
	advertRes, err := a.ADRepo.GetOneAdvert(ctx, id)
	if err != nil {
		return nil, err
	}

	return advertRes, nil
}

func (a *ADDomain) UpdateOneAdvert(ctx context.Context, id string, cpsAction model.CreateCPSAction) (*model.CPSAction, error) {

	actionData, ok := cpsAction.ActionData.(dto.UpdateAdvertRequest)
	if !ok {
		a.logger.Errorf("failed to cast action data to advert request")
		return nil, fmt.Errorf("")
	}

	if err := actionData.Validate(); err != nil {
		a.logger.Errorf("validation error", err)
		return nil, fmt.Errorf("%s", stripFieldPrefix(err))
	}
	advertAction, err := a.ADRepo.UpdateOneAdvert(ctx, id, cpsAction)
	if err != nil {
		return nil, err
	}

	return advertAction, nil
}

func (a *ADDomain) Authorize(ctx context.Context, cpsAction model.AuthorizeCPSAction) (*model.CPSAction, error) {
	cpsActionRes, err := a.ADRepo.Authorize(ctx, cpsAction)
	if err != nil {
		return nil, err
	}

	return cpsActionRes, nil
}

func (a *ADDomain) Reject(ctx context.Context, cpsAction model.RejectCPSAction) (*model.CPSAction, error) {

	if err := cpsAction.Validate(); err != nil {
		a.logger.Errorf("validation error", err)
		return nil, fmt.Errorf("%s", stripFieldPrefix(err))
	}

	cpsActionRes, err := a.ADRepo.Reject(ctx, cpsAction)
	if err != nil {
		return nil, err
	}

	return cpsActionRes, nil
}
