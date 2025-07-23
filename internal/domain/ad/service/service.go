package service

import (
	"context"
	"encoding/json"
	"fmt"
	"strings"
	"time"

	"github.com/CBE-Super-App/cbe-super-app-cps-action/internal/adapter/outbound/model"
	"github.com/CBE-Super-App/cbe-super-app-cps-action/internal/domain/ad/entity"
	"github.com/CBE-Super-App/cbe-super-app-cps-action/internal/domain/ad/repository"
	constants "github.com/CBE-Super-App/cbe-super-app-cps-action/internal/domain/cps_actions/constant"
	cps_const "github.com/CBE-Super-App/cbe-super-app-cps-action/internal/domain/cps_actions/constant"
	entities "github.com/CBE-Super-App/cbe-super-app-cps-action/internal/domain/cps_actions/entities"
	common_util "github.com/CBE-Super-App/cbe-super-app-cps-action/pkgs/utils"
	constant "github.com/CBE-Super-App/cbe-super-app-cps-action/utils"

	"github.com/CBE-Super-App/cbe-super-app-cps-action/internal/domain/ad/dto"
	"gitlab.com/bersufekadgetachew/cbe-super-app-shared/shared/config"
	"gitlab.com/bersufekadgetachew/cbe-super-app-shared/shared/utils"
)

type ADDomain struct {
	ADRepo      repository.AdRepository
	bucketName  string
	minioClient config.MinioClientInterface
	logger      utils.Logger
	cfg         *config.VaultConfig
}

type AdvertService interface {
	CreateOneAdvert(ctx context.Context, cpsAction model.CreateCPSAction) (*entities.CPSAction, error)
	UpdateOneAdvert(ctx context.Context, id string, cpsAction model.CreateCPSAction) (*entities.CPSAction, error)
	GetAllAdvert(ctx context.Context, filterParams *constant.Filter) (*common_util.PaginatedResponse[[]*entity.AdvertResponse], error)
	GetOneAdvert(ctx context.Context, id string) (*entity.Advert, error)
	DeleteOneAdvert(ctx context.Context, id string, cpsAction model.CreateCPSAction) (*entities.CPSAction, error)
	Authorize(ctx context.Context, cpsAction *entities.CPSAction) (*entities.CPSAction, error)
}

var _ AdvertService = (*ADDomain)(nil)

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

func InitADDomian(bucketName string, minioClient config.MinioClientInterface, adRepo repository.AdRepository, cfg *config.VaultConfig,
	logger utils.Logger) AdvertService {
	return &ADDomain{
		ADRepo:      adRepo,
		bucketName:  bucketName,
		minioClient: minioClient,
		logger:      logger,
		cfg:         cfg,
	}
}

func (a *ADDomain) CreateOneAdvert(ctx context.Context, cpsAction model.CreateCPSAction) (*entities.CPSAction, error) {
	actionData, ok := cpsAction.ActionData.(dto.CreateAdvertRequest)
	if !ok {
		a.logger.Errorf("failed to cast action data to ad request")
		return nil, fmt.Errorf("FAILED_TO_CAST_ACTION_DATA")
	}

	if err := actionData.Validate(); err != nil {
		a.logger.Errorf("validation error: %v", err)
		return nil, err
	}

	// Upload banner image using reusable utility
	url, err := common_util.UploadFileToMinio(
		ctx,
		a.minioClient,
		a.bucketName,
		actionData.BannerImage,
		"advert",
		a.cfg.MinioEndPoint,
		a.logger,
	)
	if err != nil {
		return nil, err
	}

	// Create CPSAction
	cpsActionRes := &entities.CPSAction{
		ActionCode:       utils.RandomGenerator(20),
		MakerID:          cpsAction.MakerUser.UserCode,
		MakerName:        cpsAction.MakerUser.FullName,
		MakerPhoneNumber: cpsAction.MakerUser.PhoneNumber,
		Department:       cpsAction.Department,
		CurrentAction: entity.Advert{
			Title:         actionData.Title,
			Description:   actionData.Description,
			BannerImage:   url,
			AdvertFor:     entity.AdvertFor(actionData.AdvertFor),
			Date:          entity.AdvertDate(actionData.Date),
			CreatedAt:     time.Now(),
			LastUpdatedAt: time.Now(),
		},
		RequestAction:   constants.RequestAction(constants.RequestCreateAdvert),
		MakerActionTime: time.Now(),
		ActionStatus:    constants.ActionPending,
		ActionType:      constants.ActionCreate,
	}

	return cpsActionRes, nil
}

func (a *ADDomain) DeleteOneAdvert(ctx context.Context, id string, cpsAction model.CreateCPSAction) (*entities.CPSAction, error) {

	// This one handles is the advert isn't there
	ad, err := a.ADRepo.GetOneAdvert(ctx, id)
	if err != nil {
		return nil, err
	}

	cpsRes := entities.CPSAction{
		ActionCode:       utils.RandomGenerator(20),
		MakerID:          cpsAction.MakerUser.UserCode,
		MakerName:        cpsAction.MakerUser.FullName,
		MakerPhoneNumber: cpsAction.MakerUser.PhoneNumber,
		Department:       cpsAction.Department,
		ActionStatus:     constants.ActionPending,
		RequestAction:    constants.RequestDeleteAdvert,
		ActionType:       constants.ActionDelete,
		CurrentAction: entity.DeletedAdvert{
			ID:            id,
			IsDeleted:     true,
			DeletedAt:     time.Now(),
			LastUpdatedAt: time.Now(),
		},
		UniqueID: id,
		PreviousAction: entity.Advert{
			ID:            id,
			Title:         ad.Title,
			Description:   ad.Description,
			IsDeleted:     ad.IsDeleted,
			BannerImage:   ad.BannerImage,
			AdvertFor:     entity.AdvertFor(ad.AdvertFor),
			Date:          entity.AdvertDate(ad.Date),
			CreatedAt:     time.Now(),
			LastUpdatedAt: time.Now(),
		},
		MakerActionTime: time.Now(),
	}

	return &cpsRes, nil
}

func (a *ADDomain) GetAllAdvert(ctx context.Context, filterParams *constant.Filter) (*common_util.PaginatedResponse[[]*entity.AdvertResponse], error) {
	advertRes, err := a.ADRepo.GetAllAdvert(ctx, filterParams)
	if err != nil {
		return nil, err
	}

	var ads []*entity.AdvertResponse
	for _, doc := range advertRes.Data {
		converted := entity.ToAdvertResponse(*doc)
		ads = append(ads, &converted)
	}

	return &common_util.PaginatedResponse[[]*entity.AdvertResponse]{
		Data: ads,
		Meta: advertRes.Meta,
	}, nil
}

func (a *ADDomain) GetOneAdvert(ctx context.Context, id string) (*entity.Advert, error) {
	advertRes, err := a.ADRepo.GetOneAdvert(ctx, id)
	if err != nil {
		return nil, err
	}

	return advertRes, nil
}

func PrettyPrintJSON(data interface{}) error {
	prettyJSON, err := json.MarshalIndent(data, "", "  ")
	if err != nil {
		return err
	}

	fmt.Println(string(prettyJSON))
	return nil
}

func (a *ADDomain) UpdateOneAdvert(ctx context.Context, id string, cpsAction model.CreateCPSAction) (*entities.CPSAction, error) {

	actionData, ok := cpsAction.ActionData.(dto.UpdateAdvertRequest)
	if !ok {
		a.logger.Errorf("failed to cast action data to advert request")
		return nil, fmt.Errorf("")
	}

	if err := actionData.Validate(); err != nil {
		a.logger.Errorf("validation error", err)
		return nil, fmt.Errorf("%s", stripFieldPrefix(err))
	}

	// This one handles is the advert isn't there
	ad, err := a.ADRepo.GetOneAdvert(ctx, id)
	if err != nil {
		return nil, err
	}

	PrettyPrintJSON(ad)
	var url string

	if actionData.BannerImage != nil {
		url, err = common_util.UploadFileToMinio(
			ctx,
			a.minioClient,
			a.bucketName,
			actionData.BannerImage,
			"advert",
			a.cfg.MinioEndPoint,
			a.logger,
		)
		if err != nil {
			return nil, err
		}
	}

	date := entity.AdvertDate(actionData.Date)

	advertAction := entities.CPSAction{
		ActionCode:       utils.RandomGenerator(20),
		MakerID:          cpsAction.MakerUser.UserCode,
		MakerName:        cpsAction.MakerUser.FullName,
		MakerPhoneNumber: cpsAction.MakerUser.PhoneNumber,
		Department:       cpsAction.Department,
		ActionStatus:     constants.ActionPending,
		ActionType:       constants.ActionUpdate,
		CurrentAction: entity.UpdateAdvert{
			Title:         actionData.Title,
			Description:   actionData.Description,
			BannerImage:   url,
			AdvertFor:     entity.AdvertFor(actionData.AdvertFor),
			Date:          &date,
			LastUpdatedAt: time.Now(),
		},
		RequestAction: constants.RequestUpdateAdvert,
		PreviousAction: entity.Advert{
			ID:            ad.ID,
			Title:         ad.Title,
			Description:   ad.Description,
			BannerImage:   ad.BannerImage,
			AdvertFor:     ad.AdvertFor,
			Date:          ad.Date,
			LastUpdatedAt: ad.LastUpdatedAt,
			Enabled:       ad.Enabled,
			IsDeleted:     ad.IsDeleted,
			DeletedAt:     ad.DeletedAt,
			CreatedAt:     ad.CreatedAt,
		},
		MakerActionTime: time.Now(),
		CreatedAt:       time.Now(),
		LastModifiedAt:  time.Now(),
	}

	return &advertAction, nil
}

func (a *ADDomain) Authorize(ctx context.Context, cpsAction *entities.CPSAction) (*entities.CPSAction, error) {
	switch cpsAction.ActionType {
	case cps_const.ActionCreate:
		return a.ADRepo.HandleAdvertCreate(ctx, cpsAction)
	case cps_const.ActionUpdate:
		return a.ADRepo.HandleAdvertUpdate(ctx, cpsAction)
	case cps_const.ActionDelete:
		return a.ADRepo.HandleAdvertDelete(ctx, cpsAction)
	default:
		a.logger.Errorf("unsupported action type: %v", cpsAction.ActionType)
		return nil, fmt.Errorf("UNSUPPORTED_ACTION_TYPE")
	}
}
