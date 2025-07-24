package miniapp

import (
	"context"
	"fmt"
	"time"

	"github.com/CBE-Super-App/cbe-super-app-cps-action/internal/application/dto"
	"github.com/CBE-Super-App/cbe-super-app-cps-action/internal/domain/cps_actions/constant"
	"gitlab.com/bersufekadgetachew/cbe-super-app-shared/shared/utils"

	entities "github.com/CBE-Super-App/cbe-super-app-cps-action/internal/domain/cps_actions/entities"
	common_util "github.com/CBE-Super-App/cbe-super-app-cps-action/pkgs/utils"
	util_constant "github.com/CBE-Super-App/cbe-super-app-cps-action/utils"
	"gitlab.com/bersufekadgetachew/cbe-super-app-shared/shared/config"
)

type MiniAppStore struct {
	Repository  MiniRepository
	logger      utils.Logger
	cfg         *config.VaultConfig
	bucketName  string
	minioClient config.MinioClientInterface
}
type MiniAppService interface {
	CreateMiniAppAction(ctx context.Context, miniApp dto.MiniAppCreateRequest, maker entities.User) (*entities.CPSAction, error)
	UpdateMiniAppAction(ctx context.Context, req dto.MiniAppCreateRequest, maker entities.User, id string) (*entities.CPSAction, error)
	DeleteMiniAppAction(ctx context.Context, maker entities.User, id string) (*entities.CPSAction, error)

	Authorize(ctx context.Context, action *entities.CPSAction) (*entities.CPSAction, error)
	ListMiniApp(ctx context.Context, filterParam *util_constant.Filter) (*common_util.PaginatedResponse[[]*MiniApp], error)
	DetailMiniAppByID(ctx context.Context, id string) (*MiniApp, error)
}

func NewService(bucketName string, minioClient config.MinioClientInterface, repository MiniRepository, cfg *config.VaultConfig, logger utils.Logger) MiniAppService {
	return &MiniAppStore{
		Repository:  repository,
		logger:      logger,
		cfg:         cfg,
		bucketName:  bucketName,
		minioClient: minioClient,
	}
}

func (s *MiniAppStore) CreateMiniAppAction(ctx context.Context, req dto.MiniAppCreateRequest, maker entities.User) (*entities.CPSAction, error) {
	data := buildMiniAppFromRequest(req, "", true)
	now := time.Now()

	var URL string
	if req.AppIcon != nil {
		url, err := common_util.UploadFileToMinio(
			ctx,
			s.minioClient,
			s.bucketName,
			req.AppIcon,
			"advert",
			s.cfg.MinioEndPoint,
			s.logger,
		)
		if err != nil {
			return nil, err
		}

		URL = url
	}

	data.AppIcon = URL

	action := entities.CPSAction{
		ActionCode:       utils.RandomGenerator(20),
		MakerID:          maker.UserCode,
		MakerName:        maker.FullName,
		MakerPhoneNumber: maker.PhoneNumber,
		Department:       maker.Department,
		ActionType:       constant.ActionCreate,
		RequestAction:    constant.RequestCreateMiniAppMerchant,
		ActionStatus:     constant.ActionPending,
		CurrentAction:    data,
		MakerActionTime:  now,
		CreatedAt:        now,
		LastModifiedAt:   now,
	}

	return &action, nil
}

func (s *MiniAppStore) Authorize(ctx context.Context, action *entities.CPSAction) (*entities.CPSAction, error) {
	requestedAction := action.RequestAction

	var minApp *MiniApp
	var err error
	switch requestedAction {
	case constant.RequestCreateMiniAppMerchant:
		minApp, err = s.Repository.CreateMiniApp(ctx, action)
	case constant.RequestUpdateMiniAppMerchant:
		minApp, err = s.Repository.UpdateMinApp(ctx, action)
	case constant.RequestDeleteMiniAppMerchant:
		minApp, err = s.Repository.DeleteMiniAppAction(ctx, action)
	default:
		return nil, fmt.Errorf("UNSUPPORTED_REQUEST_ACTION")
	}

	if err != nil {
		fmt.Printf("error form domain chekmiiapp to crate mini app : %v", err)
		return nil, err
	}

	action.CurrentAction = minApp
	return action, nil
}

func (s *MiniAppStore) UpdateMiniAppAction(ctx context.Context, req dto.MiniAppCreateRequest, maker entities.User, id string) (*entities.CPSAction, error) {
	prevData, err := s.Repository.DetailMiniAppByID(ctx, id)
	if err != nil {
		return nil, err
	}

	data := buildMiniAppFromRequest(req, id, false)
	now := time.Now()

	var url string
	if req.AppIcon != nil {
		url, err = common_util.UploadFileToMinio(
			ctx,
			s.minioClient,
			s.bucketName,
			req.AppIcon,
			"advert",
			s.cfg.MinioEndPoint,
			s.logger,
		)
		if err != nil {
			return nil, err
		}
	}

	data.AppIcon = url

	action := entities.CPSAction{
		ActionCode:       utils.RandomGenerator(20),
		MakerID:          maker.UserCode,
		MakerName:        maker.FullName,
		MakerPhoneNumber: maker.PhoneNumber,
		Department:       maker.Department,
		ActionType:       constant.ActionUpdate,
		RequestAction:    constant.RequestUpdateMiniAppMerchant,
		ActionStatus:     constant.ActionPending,
		CurrentAction:    data,
		PreviousAction:   prevData,
		MakerActionTime:  now,
	}

	return &action, nil
}

func (s *MiniAppStore) DeleteMiniAppAction(ctx context.Context, maker entities.User, id string) (*entities.CPSAction, error) {

	miniApp, err := s.Repository.DetailMiniAppByID(ctx, id)
	if err != nil {
		return nil, err
	}

	currentAction := miniApp
	currentAction.IsDeleted = true

	now := time.Now()
	action := entities.CPSAction{
		ActionCode:       utils.RandomGenerator(20),
		MakerID:          maker.UserCode,
		MakerName:        maker.FullName,
		MakerPhoneNumber: maker.PhoneNumber,
		Department:       maker.Department,
		ActionType:       constant.ActionDelete,
		RequestAction:    constant.RequestDeleteMiniAppMerchant,
		ActionStatus:     constant.ActionPending,
		PreviousAction:   miniApp,
		CurrentAction:    currentAction,
		MakerActionTime:  now,
		LastModifiedAt:   now,
	}

	return &action, nil
}

func (s *MiniAppStore) ListMiniApp(ctx context.Context, filterParam *util_constant.Filter) (*common_util.PaginatedResponse[[]*MiniApp], error) {
	return s.Repository.ListMiniApp(ctx, filterParam)
}

func (s *MiniAppStore) DetailMiniAppByID(ctx context.Context, id string) (*MiniApp, error) {
	return s.Repository.DetailMiniAppByID(ctx, id)
}

func buildMiniAppFromRequest(req dto.MiniAppCreateRequest, id string, withTimestamps bool) MiniApp {
	now := time.Now()

	productCodes := make([]ProductCode, 0, len(req.ProductCode))
	for _, pc := range req.ProductCode {
		productCodes = append(productCodes, ProductCode{
			ID:          pc.ID,
			BranchType:  BranchType(pc.BranchType),
			ProductCode: pc.ProductCode,
		})
	}

	credentials := make([]CredentialInformation, 0, len(req.Credential))
	for _, cred := range req.Credential {
		credentials = append(credentials, CredentialInformation{
			Environment:   EnvironmentType(cred.Environment),
			MerchantAppID: cred.MerchantAppID,
			FabricAppID:   cred.FabricAppID,
			ShortCode:     cred.ShortCode,
			AppSecret:     cred.AppSecret,
			PrivateKey:    cred.PrivateKey,
			PublicKey:     cred.PublicKey,
		})
	}

	miniApp := MiniApp{
		ID:                id,
		AppName:           req.AppName,
		CommisonGLAccount: req.CommisonGLAccount,
		AppType: AppType{
			UAT:        req.AppType.UAT,
			Production: req.AppType.Production,
			Test:       req.AppType.Test,
			Dev:        req.AppType.Dev,
		},
		MerchantID:     req.MerchantID,
		ProductCode:    productCodes,
		Credential:     credentials,
		IsEventMiniApp: req.IsEventMiniApp,
		IsThreeClick:   req.IsThreeClick,
		Enabled:        req.Enabled,
	}

	if withTimestamps {
		miniApp.CreatedAt = now
		miniApp.LastModifiedAt = now
		miniApp.DeletedAt = time.Time{}
	} else {
		miniApp.LastModifiedAt = now
	}

	return miniApp
}
