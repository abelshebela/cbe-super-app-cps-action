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
	"go.mongodb.org/mongo-driver/v2/bson"
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
	s.logger.Infof("Authorize called. CurrentAction (raw): %+v", cpsAction.CurrentAction)
	s.logger.Infof("RequestAction: %s, UniqueId: %s", cpsAction.RequestAction, cpsAction.UniqueId)

	var miniApp *model.MiniApp

	err := miniappcore.BindAction(cpsAction.CurrentAction, &miniApp)
	if err != nil {
		s.logger.Errorf("Failed to bind CurrentAction to MiniApp: %+v, error: %v", cpsAction.CurrentAction, err)
		return nil, errors.New(localization.ErrorInvalidRequest.Code)
	}

	s.logger.Infof("MiniApp after marshaling: %+v", miniApp)

	id := cpsAction.UniqueId

	// Trace repository operations with more granular logs
	switch cpsAction.RequestAction {
	case string(constants.RequestCreateMiniApp):
		s.logger.Infof("Processing CreateMiniApp for MiniApp: %s", miniApp.AppName)
		miniApp.ID = bson.NewObjectID()
		err = s.repo.RunInTransaction(ctx, func(ctx context.Context) error {
			if repoErr := s.repo.Create(ctx, miniApp); repoErr != nil {
				s.logger.Errorf("Create MiniApp failed: %v", repoErr)
				return repoErr
			}
			if merchantErr := s.merchantService.AddMiniApp(ctx, miniApp.MerchantID, model.MiniApps{
				ID:        miniApp.ID.Hex(),
				Enabled:   miniApp.Enabled,
				IsDeleted: miniApp.IsDeleted,
			}); merchantErr != nil {
				s.logger.Errorf("AddMiniApp to merchant failed: %v", merchantErr)
				return merchantErr
			}
			return nil
		})

	case string(constants.RequestUpdateMiniApp):
		s.logger.Infof("Processing UpdateMiniApp for ID: %s", id)
		err = s.repo.Update(ctx, id, miniApp)

	case string(constants.RequestDeleteMiniApp):
		s.logger.Infof("Processing DeleteMiniApp for ID: %s", id)
		err = s.repo.RunInTransaction(ctx, func(ctx context.Context) error {
			if repoErr := s.repo.Delete(ctx, id); repoErr != nil {
				s.logger.Errorf("Delete MiniApp repo failed: %v", repoErr)
				return repoErr
			}
			if merchantErr := s.merchantService.SoftDeleteMiniApp(ctx, miniApp.MerchantID, id); merchantErr != nil {
				s.logger.Errorf("SoftDeleteMiniApp in merchant service failed: %v", merchantErr)
				return merchantErr
			}
			return nil
		})

	case string(constants.RequestEnableMiniApp):
		s.logger.Infof("Processing EnableMiniApp for ID: %s", miniApp.ID.Hex())
		err = s.repo.RunInTransaction(ctx, func(ctx context.Context) error {
			if repoErr := s.repo.EnableOrDisable(ctx, id, true); repoErr != nil {
				s.logger.Errorf("EnableOrDisable MiniApp repo failed: %v", repoErr)
				return repoErr
			}
			if merchantErr := s.merchantService.UpdateMiniAppEnabledState(ctx, miniApp.MerchantID, id, true); merchantErr != nil {
				s.logger.Errorf("UpdateMiniAppEnabledState failed: %v", merchantErr)
				return merchantErr
			}
			return nil
		})

	case string(constants.RequestDisableMiniApp):
		s.logger.Infof("Processing DisableMiniApp for ID: %s", miniApp.ID.Hex())
		err = s.repo.RunInTransaction(ctx, func(ctx context.Context) error {
			if repoErr := s.repo.EnableOrDisable(ctx, id, false); repoErr != nil {
				s.logger.Errorf("EnableOrDisable MiniApp repo failed: %v", repoErr)
				return repoErr
			}
			if merchantErr := s.merchantService.UpdateMiniAppEnabledState(ctx, miniApp.MerchantID, id, false); merchantErr != nil {
				s.logger.Errorf("UpdateMiniAppEnabledState failed: %v", merchantErr)
				return merchantErr
			}
			return nil
		})

	default:
		s.logger.Errorf("Unsupported request action: %s", cpsAction.RequestAction)
		return nil, errors.New(localization.ErrorInvalidRequest.Code)
	}

	if err != nil {
		s.logger.Errorf("Authorize action failed for RequestAction %s, MiniApp ID %s: %v", cpsAction.RequestAction, miniApp.ID.Hex(), err)
		return nil, err
	}

	cpsAction.ActionStatus = "APPROVED"
	cpsAction.CurrentAction = miniApp
	s.logger.Infof("Action %s approved for MiniApp %s (ID: %s)", cpsAction.RequestAction, miniApp.AppName, miniApp.ID.Hex())

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
		ID:            bson.NewObjectID(),
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
