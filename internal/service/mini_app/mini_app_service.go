package miniapp

import (
	"context"
	"encoding/base64"
	"errors"
	"fmt"
	"time"

	"cbe-super-app-cps-action/internal/constants"
	miniappdto "cbe-super-app-cps-action/internal/constants/dto/mini_app"
	"cbe-super-app-cps-action/internal/constants/lib"
	"cbe-super-app-cps-action/internal/constants/localization"
	"cbe-super-app-cps-action/internal/constants/model"
	"cbe-super-app-cps-action/internal/constants/types"
	"cbe-super-app-cps-action/internal/service"
	miniappcore "cbe-super-app-cps-action/internal/service/mini_app/core"
	"cbe-super-app-cps-action/internal/storage"
	"cbe-super-app-cps-action/pkgs/keygen"

	"gitlab.com/bersufekadgetachew/cbe-super-app-shared/shared/config"
	"gitlab.com/bersufekadgetachew/cbe-super-app-shared/shared/utils"
)

type miniAppService struct {
	repo            storage.MiniAppRepository
	cpsService      service.CPSActionService
	merchantService service.MiniAppMerchantService
	userRepo        storage.UserRepository
	keyGenService   keygen.KeyGeneratorService
	bucketName      string
	minio           config.MinioClientInterface
	cfg             *config.VaultConfig
	logger          utils.Logger
}

func NewMiniAppService(
	repo storage.MiniAppRepository,
	cpsService service.CPSActionService,
	merchantService service.MiniAppMerchantService,
	userRepo storage.UserRepository,
	keyGenService keygen.KeyGeneratorService,
	minio config.MinioClientInterface,
	bucketName string,
	cfg *config.VaultConfig,
	logger utils.Logger,
) service.MiniAppService {
	return &miniAppService{
		repo:            repo,
		cpsService:      cpsService,
		merchantService: merchantService,
		userRepo:        userRepo,
		keyGenService:   keyGenService,
		minio:           minio,
		bucketName:      bucketName,
		cfg:             cfg,
		logger:          logger,
	}
}

func (s *miniAppService) CreateMiniApp(ctx context.Context, req *miniappdto.MiniAppCreateRequest) error {
	if err := miniappcore.SetMerchantDetails(ctx, s.merchantService, req); err != nil {
		s.logger.Errorf("SetMerchantDetails failed, error: %v", err)
		return err
	}

	existing, err := s.repo.Find(ctx, req.AppName)
	if err != nil {
		s.logger.Errorf("Find failed: %v", err)
		return errors.New(localization.ErrorUnhandledServer.Code)
	}
	if existing != nil {
		s.logger.Warnf("MiniApp already exists: %s", req.AppName)
		return errors.New(localization.ErrorMiniAppAlreadyExists.Code)
	}

	var appIconURL, bannerImageURL string
	if req.AppIcon != nil {
		appIconURL, err = lib.UploadFileToMinio(ctx, s.minio, s.bucketName, req.AppIcon, "miniapp_app_icon", s.cfg.MinioEndPoint, s.logger)
		if err != nil {
			s.logger.Errorf("UploadFileToMinio failed for app icon, error: %v", err)
			return errors.New(localization.ErrorUnhandledServer.Code)
		}
		s.logger.Infof("Uploaded app icon to URL: %s", appIconURL)
	}

	if req.BannerImage != nil {
		bannerImageURL, err = lib.UploadFileToMinio(ctx, s.minio, s.bucketName, req.BannerImage, "miniapp_banner_image", s.cfg.MinioEndPoint, s.logger)
		if err != nil {
			s.logger.Errorf("UploadFileToMinio failed for banner image, error: %v", err)
			return errors.New(localization.ErrorUnhandledServer.Code)
		}
		s.logger.Infof("Uploaded banner image to URL: %s", bannerImageURL)
	}

	cred, _, err := s.CredentialInformationGenerator(constants.UatEnvironment)
	if err != nil {
		s.logger.Errorf("CredentialInformationGenerator failed for environment %v: %v", constants.UatEnvironment, err)
		return errors.New(localization.ErrorUnhandledServer.Code)
	}

	miniApp := miniappcore.BuildMiniAppFromRequest(*req, true)
	miniApp.AppIcon = appIconURL
	miniApp.BannerImage = bannerImageURL
	miniApp.Credential = *cred
 
	s.logger.Infof("MiniApp created successfully, app_code: %s", req.AppName)
	if err := miniappcore.HandleCPSAction(ctx, s.cpsService, "", constants.RequestCreateMiniApp, miniApp, nil, constants.ActionCreate); err != nil {
		s.logger.Errorf("CPS action failed for MiniApp %s: %v", req.AppName, err)
		return err
	}

	return nil
}

func (s *miniAppService) UpdateMiniApp(ctx context.Context, req *miniappdto.MiniAppCreateRequest) error {
	s.logger.Infof("UpdateMiniApp called, app_id: %s", req.ID)

	if err := miniappcore.SetMerchantDetails(ctx, s.merchantService, req); err != nil {
		s.logger.Errorf("SetMerchantDetails failed, app_id: %s, error: %v", req.ID, err)
		return err
	}

	prevMiniApp, err := s.repo.FindByID(ctx, req.ID)
	if err != nil {
		s.logger.Errorf("FindByID failed, app_id: %s, error: %v", req.ID, err)
		return errors.New(localization.ErrorMiniAppNotFound.Code)
	}
	miniApp := miniappcore.BuildMiniAppFromRequest(*req, false)

	var appIconURL, bannerImageURL string
	if req.AppIcon != nil {
		appIconURL, err = lib.UploadFileToMinio(ctx, s.minio, s.bucketName, req.AppIcon, "miniapp_app_icon", s.cfg.MinioEndPoint, s.logger)
		if err != nil {
			s.logger.Errorf("UploadFileToMinio failed for app icon, app_id: %s, error: %v", req.ID, err)
			return errors.New(localization.ErrorUnhandledServer.Code)
		}
		s.logger.Infof("Updated app icon to URL: %s", appIconURL)
	} else {
		appIconURL = prevMiniApp.AppIcon
	}

	if req.BannerImage != nil {
		bannerImageURL, err = lib.UploadFileToMinio(ctx, s.minio, s.bucketName, req.BannerImage, "miniapp_banner_image", s.cfg.MinioEndPoint, s.logger)
		if err != nil {
			s.logger.Errorf("UploadFileToMinio failed for banner image, app_id: %s, error: %v", req.ID, err)
			return errors.New(localization.ErrorUnhandledServer.Code)
		}
		s.logger.Infof("Updated banner image to URL: %s", bannerImageURL)
	} else {
		bannerImageURL = prevMiniApp.BannerImage
	}

	miniApp.AppIcon = appIconURL
	miniApp.BannerImage = bannerImageURL

	s.logger.Infof("MiniApp updated successfully, app_id: %s", req.ID)
	err = miniappcore.HandleCPSAction(ctx, s.cpsService, req.ID, constants.RequestUpdateMiniApp, miniApp, prevMiniApp, constants.ActionUpdate)
	if err != nil {
		s.logger.Errorf("CPS action failed for MiniApp %s: %v", miniApp.AppName, err)
		return err
	}
	return nil
}

func (s *miniAppService) DeleteMiniApp(ctx context.Context, id string) error {
	s.logger.Infof("DeleteMiniApp called, app_id: %s", id)

	prevMiniApp, err := s.repo.FindByID(ctx, id)
	if err != nil {
		s.logger.Errorf("FindByID failed, app_id: %s, error: %v", id, err)
		return errors.New(localization.ErrorMiniAppNotFound.Code)
	}

	miniApp := *prevMiniApp
	miniApp.IsDeleted = true
	miniApp.DeletedAt = time.Now()

	err = miniappcore.HandleCPSAction(ctx, s.cpsService, id, constants.RequestDeleteMiniApp, miniApp, prevMiniApp, constants.ActionDelete)
	if err != nil {
		s.logger.Errorf("CPS action failed for MiniApp %s: %v", miniApp.AppName, err)
		return err
	}
	return nil
}

func (s *miniAppService) EnableDisableMiniAppByID(ctx context.Context, id string, enable bool) error {
	s.logger.Infof("EnableDisableMiniAppByID called, app_id: %s, enable: %v", id, enable)

	prevMiniApp, err := s.repo.FindByID(ctx, id)
	if err != nil {
		s.logger.Errorf("FindByID failed, app_id: %s, error: %v", id, err)
		return errors.New(localization.ErrorMiniAppNotFound.Code)
	}

	if enable && prevMiniApp.Enabled {
		s.logger.Warnf("MiniApp already enabled, app_id: %s", id)
		return errors.New(localization.ErrorMiniAppAlreadyEnabled.Code)
	}
	if !enable && !prevMiniApp.Enabled {
		s.logger.Warnf("MiniApp already disabled, app_id: %s", id)
		return errors.New(localization.ErrorMiniAppAlreadyDisabled.Code)
	}

	miniApp := *prevMiniApp
	miniApp.Enabled = enable
	miniApp.LastModifiedAt = time.Now()

	var action constants.RequestAction
	if enable {
		action = constants.RequestEnableMiniApp
	} else {
		action = constants.RequestDisableMiniApp
	}

	s.logger.Infof("MiniApp enable/disable action handled, app_id: %s, enable: %v", id, enable)
	err = miniappcore.HandleCPSAction(ctx, s.cpsService, id, action, miniApp, prevMiniApp, constants.ActionUpdate)
	if err != nil {
		s.logger.Errorf("CPS action failed for MiniApp %s: %v", miniApp.AppName, err)
		return err
	}
	return nil
}

func (s *miniAppService) FindByID(ctx context.Context, id string) (*model.MiniApp, error) {
	s.logger.Infof("FindByID called, app_id: %s", id)
	miniApp, err := s.repo.FindByID(ctx, id)
	if err != nil {
		s.logger.Errorf("FindByID failed, app_id: %s, error: %v", id, err)
		return nil, err
	}
	s.logger.Infof("FindByID succeeded, miniApp: %+v", miniApp)
	decryptedSecret, err := s.keyGenService.DecryptAppSecret(miniApp.Credential.AppSecret)
	if err != nil {
		s.logger.Errorf("Failed to decrypt AppSecret for MiniApp %s, Environment %v: %v", id, constants.UatEnvironment, err)
		return nil, err
	}

	miniApp.Credential.AppSecret = decryptedSecret
	return miniApp, nil
}

func (s *miniAppService) ListMiniApp(ctx context.Context, filter types.Filter) (*types.PaginatedResponse[[]*model.MiniApp], error) {
	s.logger.Infof("ListMiniApp called with filter: %+v", filter)
	list, err := s.repo.FindAllWithPagination(ctx, filter)
	if err != nil {
		s.logger.Errorf("FindAllWithPagination failed: %v", err)
		return nil, err
	}
	return list, nil
}

func (s *miniAppService) Authorize(ctx context.Context, cpsAction *model.CPSAction) (*model.CPSAction, error) {
	s.logger.Infof("Authorize called, action: %s", cpsAction.ActionCode)

	miniApp, ok := cpsAction.CurrentAction.(*model.MiniApp)
	if !ok || miniApp == nil {
		s.logger.Errorf("Invalid CPS action data: not a MiniApp")
		return nil, errors.New(localization.ErrorInvalidRequest.Code)
	}
	var err error
	switch cpsAction.RequestAction {
	case string(constants.RequestCreateMiniApp):
		err = s.repo.Create(ctx, miniApp)
	case string(constants.RequestUpdateMiniApp):
		err = s.repo.Update(ctx, miniApp.ID.Hex(), miniApp)
	case string(constants.RequestDeleteMiniApp):
		err = s.repo.Delete(ctx, miniApp.ID.Hex())
	case string(constants.RequestEnableMiniApp):
		err = s.repo.EnableOrDisable(ctx, miniApp.ID.Hex(), true)
	case string(constants.RequestDisableMiniApp):
		err = s.repo.EnableOrDisable(ctx, miniApp.ID.Hex(), false)
	default:
		s.logger.Errorf("Invalid request action: %s", cpsAction.RequestAction)
		return nil, errors.New(localization.ErrorInvalidRequest.Code)
	}

	if err != nil {
		s.logger.Errorf("Repository operation failed for action %s: %v", cpsAction.RequestAction, err)
		return nil, err
	}

	cpsAction.ActionStatus = "APPROVED"
	s.logger.Infof("Action %s approved for MiniApp %s", cpsAction.RequestAction, miniApp.AppName)
	return cpsAction, nil
}

func (s *miniAppService) CredentialInformationGenerator(envType constants.EnvironmentType) (*types.CredentialInformation, *string, error) {
	s.logger.Infof("Generating credential information for environment: %v", envType)

	merchantCode, err := s.keyGenService.GenerateNumericCode(15)
	if err != nil {
		s.logger.Errorf("Failed to generate merchant code: %v", err)
		return nil, nil, err
	}
	s.logger.Debugf("Generated merchant code: %s", merchantCode)

	fabID := s.keyGenService.GenerateFabricID()
	s.logger.Debugf("Generated fabric ID: %s", fabID)

	shortCode, err := s.keyGenService.GenerateNumericCode(6)
	if err != nil {
		s.logger.Errorf("Failed to generate short code: %v", err)
		return nil, nil, err
	}
	s.logger.Debugf("Generated short code: %s", shortCode)

	miniAppCode, err := s.keyGenService.GenerateNumericCode(6)
	if err != nil {
		s.logger.Errorf("Failed to generate mini app code: %v", err)
		return nil, nil, err
	}
	s.logger.Debugf("Generated mini app code: %s", miniAppCode)

	keys, err := s.keyGenService.GenerateKeyPair()
	if err != nil {
		s.logger.Errorf("Failed to generate key pair: %v", err)
		return nil, nil, err
	}
	s.logger.Debugf("Generated key pair for environment: %v", envType)

	rawSecret, err := s.keyGenService.GenerateAppSecret(32)
	if err != nil {
		s.logger.Errorf("Failed to generate app secret: %v", err)
		return nil, nil, err
	}
	s.logger.Debugf("Generated app secret for environment: %v", envType)

	encryptedSecret, err := s.keyGenService.EncryptAppSecret(rawSecret)
	if err != nil {
		s.logger.Errorf("Failed to encrypt app secret: %v", err)
		return nil, nil, err
	}
	s.logger.Debugf("Encrypted app secret for environment: %v", envType)

	timestamp := time.Now().UTC()
	sorted := fmt.Sprintf("%s|%s|%s|%s", merchantCode, fabID, miniAppCode, timestamp.Format(time.RFC3339))
	signature, err := s.keyGenService.Sign([]byte(sorted), keys.PrivateKey)
	if err != nil {
		s.logger.Errorf("Failed to sign credential data: %v", err)
		return nil, nil, err
	}
	s.logger.Debugf("Generated signature for environment: %v", envType)

	res := types.CredentialInformation{
		Environment:   envType,
		MerchantAppID: merchantCode,
		FabricAppID:   fabID,
		ShortCode:     shortCode,
		MiniAppCode:   miniAppCode,
		PrivateKey:    keys.PrivateKey,
		PublicKey:     keys.PublicKey,
		AppSecret:     encryptedSecret,
		Signature:     base64.StdEncoding.EncodeToString(signature),
		Timestamp:     timestamp,
	}

	s.logger.Infof("Successfully generated credential information for environment: %v", envType)
	return &res, &rawSecret, nil
}
