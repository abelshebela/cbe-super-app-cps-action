package ad

import (
	"context"
	"fmt"
	"net/http"
	"os"
	"path/filepath"
	"time"

	"gitlab.com/bersufekadgetachew/cbe-super-app-cps-action/internal/domain/ad/entity"
	"gitlab.com/bersufekadgetachew/cbe-super-app-cps-action/internal/domain/ad/service"
	constant "gitlab.com/bersufekadgetachew/cbe-super-app-cps-action/utils"
	"gitlab.com/bersufekadgetachew/cbe-super-app-shared/shared/config"
	"gitlab.com/bersufekadgetachew/cbe-super-app-shared/shared/utils"
)

type ADHandlers interface {
	CreateOneAdvert(ctx context.Context, adCpsReq CreateCPSAction) (*CPSAction, error)
	GetAllAdvert(ctx context.Context, filterParams *constant.Filter) (*entity.AdvertResponse, error)
	GetOneAdvert(ctx context.Context, id string) (*AdvertResponse, error)
	UpdateOneAdvert(ctx context.Context, cpsAction CreateCPSAction) (*entity.CPSAction, error)
	DeleteOneAdvert(ctx context.Context, adCpsReq CreateCPSAction) error
	Authorize(ctx context.Context, cpsAction entity.CPSAction) (*entity.CPSAction, error)
	Reject(ctx context.Context, cpsAction entity.CPSAction) (*entity.CPSAction, error)
}

type ADHandler struct {
	adDomain   service.AdvertService
	bucketName string
	minio      config.MinioClientInterface
	logger     utils.Logger
}

func InitADHandler(adDomain service.AdvertService, minioClinet config.MinioClientInterface, bucketName string, logger utils.Logger) ADHandlers {
	return ADHandler{
		adDomain:   adDomain,
		minio:      minioClinet,
		logger:     logger,
		bucketName: bucketName,
	}
}

func (a ADHandler) CreateOneAdvert(ctx context.Context, adCpsReq CreateCPSAction) (*CPSAction, error) {
	if err := adCpsReq.ActionData.Validate(); err != nil {
		a.logger.Errorf("invalid data", err)
		return nil, err
	}

	exist, err := a.minio.BucketExist(ctx, a.bucketName)
	if err != nil {
		a.logger.Errorf("failed to check advert bucket: %v", err)
		return nil, fmt.Errorf("failed to check advert bucket: %w", constant.ErrorDefinition{
			Code:    http.StatusInternalServerError,
			Message: "internal server error",
		})
	}

	if !exist {
		created, err := a.minio.MakeBucket(ctx, a.bucketName)
		if !created || err != nil {
			a.logger.Errorf("failed to create advert bucket: %v", err)
			return nil, fmt.Errorf("failed to create advert bucket: %w", constant.ErrorDefinition{
				Code:    http.StatusInternalServerError,
				Message: "internal server error",
			})
		}
	}

	fileName := fmt.Sprintf("upload-%d-%s", time.Now().UnixNano(), adCpsReq.ActionData.BannerImage.Filename)
	dir, err := os.Getwd()
	if err != nil {
		a.logger.Errorf("failed to get current working directory", err)
		return nil, fmt.Errorf("failed to create advert bucket: %w", constant.ErrorDefinition{
			Code:    http.StatusInternalServerError,
			Message: "internal server error",
		})
	}
	filePath := filepath.Join(dir, fileName)

	tempFile, err := os.Create(filePath)
	if err != nil {
		a.logger.Errorf("failed to create temp file: %v", err)
		return nil, fmt.Errorf("failed to create temp file: %w", constant.ErrorDefinition{
			Code:    http.StatusInternalServerError,
			Message: "internal server error",
		})
	}
	defer func() {
		tempFile.Close()
		os.Remove(filePath)
	}()

	// Save to MinIO
	saveObj, err := a.minio.SaveObject(ctx, config.SaveObjectBody{
		BucketName: a.bucketName,
		ObjectName: fileName,
		File:       filePath,
	})
	if err != nil {
		a.logger.Errorf("failed to save object to MinIO: %v", err)
		return nil, fmt.Errorf("failed to save object to MinIO: %w", constant.ErrorDefinition{
			Code:    http.StatusInternalServerError,
			Message: "internal server error",
		})
	}

	cpsAction, err := a.adDomain.CreateOneAdvert(ctx, entity.CPSAction{
		ActionCode: utils.RandomGenerator(20),
		MakerUser:  entity.User(adCpsReq.MakerUser),
		Department: adCpsReq.Department,
		ActionData: entity.Advert{
			Title:       adCpsReq.ActionData.Title,
			Description: adCpsReq.ActionData.Description,
			BannerImage: saveObj.Bucket + "/" + saveObj.Key,
			AdvertFor:   entity.AdvertFor(adCpsReq.ActionData.AdvertFor),
			Date: entity.AdvertDate{
				StartedAt: adCpsReq.ActionData.Date.StartedAt,
				ExpiredAt: adCpsReq.ActionData.Date.ExpiredAt,
			},
		},
	})

	if err != nil {
		return nil, err
	}

	return &CPSAction{
		ID:         cpsAction.ID,
		ActionCode: cpsAction.ActionCode,
		MakerUser: User{
			UserCode:    cpsAction.MakerUser.UserCode,
			FullName:    cpsAction.MakerUser.FullName,
			PhoneNumber: cpsAction.MakerUser.PhoneNumber,
		},
		Department:      cpsAction.Department,
		Status:          ActionStatus(cpsAction.Status),
		RequestAction:   RequestAction(cpsAction.RequestAction),
		ActionType:      ActionType(cpsAction.ActionType),
		MakerActionTime: cpsAction.MakerActionTime,
		ActionData: AdvertResponse{
			ID:          cpsAction.ActionData.ID,
			Title:       cpsAction.ActionData.Title,
			Description: cpsAction.ActionData.Description,
			BannerImage: cpsAction.ActionData.BannerImage,
			AdvertFor:   AdvertFor(cpsAction.ActionData.AdvertFor),
			Date: AdvertDate{
				StartedAt: cpsAction.ActionData.Date.StartedAt,
				ExpiredAt: cpsAction.ActionData.Date.ExpiredAt,
			},
		},
	}, nil
}

func (a ADHandler) DeleteOneAdvert(ctx context.Context, adCpsReq CreateCPSAction) error {
	err := a.adDomain.DeleteOneAdvert(ctx, entity.CPSAction{
		ActionCode: utils.RandomGenerator(20),
		MakerUser:  entity.User(adCpsReq.MakerUser),
		Department: adCpsReq.Department,
		ActionData: entity.Advert{
			ID: adCpsReq.ActionData.ID,
		},
	})
	return err
}

func (a ADHandler) GetAllAdvert(ctx context.Context, filterParams *constant.Filter) (*entity.AdvertResponse, error) {
	adverts, err := a.adDomain.GetAllAdvert(ctx, filterParams)
	if err != nil {
		return nil, err
	}

	return adverts, nil
}

func (a ADHandler) GetOneAdvert(ctx context.Context, id string) (*AdvertResponse, error) {
	advert, err := a.adDomain.GetOneAdvert(ctx, id)
	if err != nil {
		return nil, err
	}

	return &AdvertResponse{
		ID:          advert.ID,
		Title:       advert.Title,
		Description: advert.Description,
		BannerImage: advert.BannerImage,
		AdvertFor:   AdvertFor(advert.AdvertFor),
		Date:        AdvertDate(advert.Date),
	}, nil
}

func (a ADHandler) UpdateOneAdvert(ctx context.Context, cpsAction CreateCPSAction) (*entity.CPSAction, error) {
	advertCpsAction, err := a.adDomain.UpdateOneAdvert(ctx, entity.CPSAction{
		ActionCode: utils.RandomGenerator(20),
		MakerUser:  entity.User(cpsAction.MakerUser),
		Department: cpsAction.Department,
		ActionData: entity.Advert{
			ID:          cpsAction.ActionData.ID,
			Title:       cpsAction.ActionData.Title,
			Description: cpsAction.ActionData.Description,
			AdvertFor:   entity.AdvertFor(cpsAction.ActionData.AdvertFor),
			Date: entity.AdvertDate{
				StartedAt: cpsAction.ActionData.Date.StartedAt,
				ExpiredAt: cpsAction.ActionData.Date.ExpiredAt,
			},
		},
	})
	if err != nil {
		return nil, err
	}

	return advertCpsAction, nil
}

func (a ADHandler) Authorize(ctx context.Context, cpsAction entity.CPSAction) (*entity.CPSAction, error) {
	cpsActionRes, err := a.adDomain.Authorize(ctx, cpsAction)
	if err != nil {
		return nil, err
	}

	return cpsActionRes, nil
}

func (a ADHandler) Reject(ctx context.Context, cpsAction entity.CPSAction) (*entity.CPSAction, error) {
	cpsActionRes, err := a.adDomain.Reject(ctx, cpsAction)
	if err != nil {
		return nil, err
	}

	return cpsActionRes, nil
}
