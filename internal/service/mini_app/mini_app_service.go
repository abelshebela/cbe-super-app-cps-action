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
	local_util "cbe-super-app-cps-action/pkgs/utils"

	"strings"

	"gitlab.com/bersufekadgetachew/cbe-super-app-shared/shared/config"
	"gitlab.com/bersufekadgetachew/cbe-super-app-shared/shared/utils"
	"go.mongodb.org/mongo-driver/v2/bson"
	"go.mongodb.org/mongo-driver/v2/mongo"
)

type miniAppService struct {
	repo            storage.MiniAppRepository
	cpsService      service.CPSActionService
	merchantService service.MiniAppMerchantService
	userRepo        storage.UserRepository
	keyGenService   keygen.KeyGeneratorService
	bucketName      string
	minio           config.MinioClientInterface
	minioPubUrl     string
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
	minioPubUrl string,
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
		minioPubUrl:     minioPubUrl,
		cfg:             cfg,
		logger:          logger,
	}
}

func (s *miniAppService) CreateMiniApp(ctx context.Context, req *miniappdto.MiniAppCreateRequest) error {
	s.logger.Infof("CreateMiniApp called, app_name: %s", req.AppName)
	if err := miniappcore.SetMerchantDetails(ctx, s.merchantService, req); err != nil {
		s.logger.Errorf("SetMerchantDetails failed, error: %v", err)
		return err
	}
	isValidMMiniAppName, err := miniappcore.ValidMiniAppChecker(ctx, s.repo, true, "", req.AppName)
	if err != nil {
		if err != mongo.ErrNoDocuments {
			s.logger.Errorf("IsMiniAppNameUnique check failed, error: %v", err)
			return errors.New(localization.ErrorUnhandledServer.Code)
		}
	}

	if !isValidMMiniAppName {
		s.logger.Warnf("MiniApp name already exists: %s", req.AppName)
		return errors.New(localization.ErrorMiniAppNameAlreadyExists.Code)
	}

	var appIconURL, bannerImageURL string
	if req.AppIcon != nil {
		appIconURL, err = lib.UploadFileToMinio(ctx, s.minio, s.bucketName, req.AppIcon, "miniapp_app_icon", s.minioPubUrl, s.logger)
		if err != nil {
			s.logger.Errorf("UploadFileToMinio failed for app icon, error: %v", err)
			return errors.New(localization.ErrorUnhandledServer.Code)
		}
		s.logger.Infof("Uploaded app icon to URL: %s", appIconURL)
	}

	if req.BannerImage != nil {
		bannerImageURL, err = lib.UploadFileToMinio(ctx, s.minio, s.bucketName, req.BannerImage, "miniapp_banner_image", s.minioPubUrl, s.logger)
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
	miniApp.Enabled= true

	s.logger.Infof("MiniApp created successfully, app_code: %s", req.AppName)
	if err := miniappcore.HandleCPSAction(ctx, s.cpsService, "", constants.RequestCreateMiniApp, miniApp, nil, constants.ActionCreate); err != nil {
		s.logger.Errorf("CPS action failed for MiniApp %s: %v", req.AppName, err)
		return err
	}

	return nil
}

func (s *miniAppService) UpdateMiniApp(ctx context.Context, req *miniappdto.MiniAppCreateRequest) error {
	s.logger.Infof("UpdateMiniApp called, app_id: %s", req.ID)

	prevMiniApp, err := s.repo.FindByID(ctx, req.ID)
	if err != nil {
		s.logger.Errorf("FindByID failed, app_id: %s, error: %v", req.ID, err)
		return errors.New(localization.ErrorMiniAppNotFound.Code)
	}

	if prevMiniApp.IsDeleted {
		s.logger.Errorf("MiniApp is deleted, app_id: %s", req.ID)
		return errors.New(localization.ErrorMiniAppNotFound.Code)
	}

	merchant, mErr := s.merchantService.FindByID(ctx, prevMiniApp.MerchantID)
	if mErr != nil {
		s.logger.Errorf("Parent merchant not found for update, app_id: %s", req.ID)
		return errors.New(localization.ErrorMiniAppMerchantNotFound.Code)
	}
	if err := miniappcore.ValidateParentMerchantForOperationEntity(merchant); err != nil {
		s.logger.Errorf("Parent merchant validation failed, app_id: %s, error: %v", req.ID, err)
		return err
	}
	
	if req.AppName!= "" && req.AppName!=prevMiniApp.AppName {


isValidMMiniAppName, err := miniappcore.ValidMiniAppChecker(ctx, s.repo, false, req.ID, req.AppName)
	if err != nil {
		s.logger.Errorf("IsMiniAppNameUnique check failed, error: %v", err)
		return errors.New(localization.ErrorUnhandledServer.Code)
	}

	if !isValidMMiniAppName {
		s.logger.Warnf("MiniApp name already exists: %s", req.AppName)
		return errors.New(localization.ErrorMiniAppNameAlreadyExists.Code)
	}
	}

	

	if strings.TrimSpace(req.MerchantID) != "" {
		if err := miniappcore.SetMerchantDetails(ctx, s.merchantService, req); err != nil {
			s.logger.Errorf("SetMerchantDetails failed, app_id: %s, error: %v", req.ID, err)
			return err
		}
	}

	miniApp := miniappcore.BuildMiniAppFromRequest(*req, false)

	var appIconURL, bannerImageURL string
	if req.AppIcon != nil {
		appIconURL, err = lib.UploadFileToMinio(ctx, s.minio, s.bucketName, req.AppIcon, "miniapp_app_icon", s.minioPubUrl, s.logger)
		if err != nil {
			s.logger.Errorf("UploadFileToMinio failed for app icon, app_id: %s, error: %v", req.ID, err)
			return errors.New(localization.ErrorUnhandledServer.Code)
		}
		s.logger.Infof("Updated app icon to URL: %s", appIconURL)
	} else {
		appIconURL = prevMiniApp.AppIcon
	}

	if req.BannerImage != nil {
		bannerImageURL, err = lib.UploadFileToMinio(ctx, s.minio, s.bucketName, req.BannerImage, "miniapp_banner_image", s.minioPubUrl, s.logger)
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

	if prevMiniApp.IsDeleted {
		s.logger.Errorf("MiniApp is already deleted, app_id: %s", id)
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

	if prevMiniApp.IsDeleted {
		s.logger.Errorf("MiniApp is deleted, app_id: %s", id)
		return errors.New(localization.ErrorMiniAppNotFound.Code)
	}

	if enable {
		merchant, mErr := s.merchantService.FindByID(ctx, prevMiniApp.MerchantID)
		if mErr != nil {
			return errors.New(localization.ErrorMiniAppMerchantNotFound.Code)
		}
		if err := miniappcore.ValidateParentMerchantForEnableEntity(merchant); err != nil {
			s.logger.Errorf("Parent merchant validation failed for enable, app_id: %s, error: %v", id, err)
			return err
		}
		if prevMiniApp.Enabled {
			s.logger.Warnf("MiniApp already enabled, app_id: %s", id)
			return errors.New(localization.ErrorMiniAppAlreadyEnabled.Code)
		}
	} else {
		merchant, mErr := s.merchantService.FindByID(ctx, prevMiniApp.MerchantID)
		if mErr != nil {
			return errors.New(localization.ErrorMiniAppMerchantNotFound.Code)
		}
		if err := miniappcore.ValidateParentMerchantForOperationEntity(merchant); err != nil {
			s.logger.Errorf("Parent merchant validation failed for disable, app_id: %s, error: %v", id, err)
			return err
		}
		if !prevMiniApp.Enabled {
			s.logger.Warnf("MiniApp already disabled, app_id: %s", id)
			return errors.New(localization.ErrorMiniAppAlreadyDisabled.Code)
		}
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

func (s *miniAppService) FindByID(ctx context.Context, id string) (*miniappdto.MiniAppResponse, error) {
	s.logger.Infof("FindByID called, app_id: %s", id)
	miniApp, err := s.repo.FindByIDWithMerchant(ctx, id)
	if err != nil {
		s.logger.Errorf("FindByID failed, app_id: %s, error: %v", id, err)
		return nil, errors.New(localization.ErrorMiniAppNotFound.Code)
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
	s.logger.Infof("Authorize called, action: %s", cpsAction.RequestAction)

	miniApp, err := local_util.JsonUnmarshal[model.MiniApp](cpsAction.CurrentAction)
	if err != nil {
		return nil, errors.New(localization.ErrorInvalidRequest.Code)
	}
	switch cpsAction.RequestAction {
	case string(constants.RequestCreateMiniApp):
		err = miniappcore.ValidMerchant(miniApp.MerchantID, ctx, s.merchantService)
		if err != nil{
			break
		}
		err = s.repo.RunInTransaction(ctx, func(ctx context.Context) error {
			return s.repo.Create(ctx, miniApp)
		})

	case string(constants.RequestUpdateMiniApp):
		err = miniappcore.ValidMerchant(miniApp.MerchantID, ctx, s.merchantService)
		if err != nil{
			break
		}
		err = s.repo.Update(ctx, cpsAction.UniqueId, miniApp)

	case string(constants.RequestDeleteMiniApp):
		err = s.repo.Delete(ctx, cpsAction.UniqueId)

	case string(constants.RequestEnableMiniApp):
		err = miniappcore.ValidMerchant(miniApp.MerchantID, ctx, s.merchantService)
		if err != nil{
			break
		}
		err = s.repo.EnableOrDisable(ctx, cpsAction.UniqueId, true)

	case string(constants.RequestDisableMiniApp):
		err = s.repo.EnableOrDisable(ctx, cpsAction.UniqueId, false)

	default:
		s.logger.Errorf("Unsupported request action: %s", cpsAction.RequestAction)
		return nil, errors.New(localization.ErrorInvalidRequest.Code)
	}

	if err != nil {
		s.logger.Errorf("Failed to process action %s: %v", cpsAction.RequestAction, err)
		return nil, err
	}

	cpsAction.ActionStatus = "APPROVED"
	cpsAction.CurrentAction = miniApp
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
