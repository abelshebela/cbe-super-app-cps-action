package miniapp

import (
	"context"
	"encoding/base64"
	"fmt"
	"mime/multipart"
	"time"

	"github.com/CBE-Super-App/cbe-super-app-cps-action/internal/domain/cps_actions/constant"
	"gitlab.com/bersufekadgetachew/cbe-super-app-shared/shared/utils"

	entities "github.com/CBE-Super-App/cbe-super-app-cps-action/internal/domain/cps_actions/entities"
	merchant_service "github.com/CBE-Super-App/cbe-super-app-cps-action/internal/domain/miniapp_merchant"
	keyGen "github.com/CBE-Super-App/cbe-super-app-cps-action/pkgs/keygen"
	common_util "github.com/CBE-Super-App/cbe-super-app-cps-action/pkgs/utils"
	util_constant "github.com/CBE-Super-App/cbe-super-app-cps-action/utils"
	"gitlab.com/bersufekadgetachew/cbe-super-app-shared/shared/config"
)

type MiniAppStore struct {
	repository      MiniRepository
	logger          utils.Logger
	cfg             *config.VaultConfig
	bucketName      string
	minioClient     config.MinioClientInterface
	keyGenService   keyGen.KeyGeneratorService
	merchantService merchant_service.MiniAppMerchantService
}
type MiniAppService interface {
	CreateMiniApp(ctx context.Context, miniApp MiniAppCreateRequest, maker entities.User) (*MiniApp, error)
	UpdateMiniApp(ctx context.Context, req MiniAppCreateRequest, maker entities.User, id string) (*MiniApp, *MiniApp, error)
	DeleteMiniApp(ctx context.Context, maker entities.User, id string) (*MiniApp, *MiniApp, error)
	Authorize(ctx context.Context, action *entities.CPSAction) (*entities.CPSAction, error)
	ListMiniApp(ctx context.Context, filterParam *util_constant.Filter) (*common_util.PaginatedResponse[[]*MiniApp], error)
	DetailMiniAppByID(ctx context.Context, id string) (*MiniApp, error)
	EnableDisableMiniApp(ctx context.Context, maker entities.User, id string, enabled bool) (*MiniApp, *MiniApp, error)
}

func NewService(bucketName string, minioClient config.MinioClientInterface, repository MiniRepository, cfg *config.VaultConfig,
	keyGenService keyGen.KeyGeneratorService,
	merchantService merchant_service.MiniAppMerchantService,
	logger utils.Logger) MiniAppService {
	logger.Infof("Initializing MiniAppStore with bucketName: %s", bucketName)
	return &MiniAppStore{
		repository:      repository,
		logger:          logger,
		cfg:             cfg,
		bucketName:      bucketName,
		minioClient:     minioClient,
		keyGenService:   keyGenService,
		merchantService: merchantService,
	}
}

func (s *MiniAppStore) CreateMiniApp(ctx context.Context, req MiniAppCreateRequest, maker entities.User) (*MiniApp, error) {
	s.logger.Infof("Creating MiniApp for maker: %v", maker)
	miniApp := buildMiniAppFromRequest(req, "", true)

	url, err := s.uploadOrFallback(ctx, req.AppIcon, "miniapp_app_icon", s.cfg.MinioEndPoint, "")
	if err != nil {
		s.logger.Errorf("Failed to upload app icon: %v", err)
		return nil, err
	}
	s.logger.Infof("Uploaded app icon to URL: %s", url)
	miniApp.AppIcon = url

	url, err = s.uploadOrFallback(ctx, req.BannerImage, "miniapp_banner_image", s.cfg.MinioEndPoint, "")
	if err != nil {
		s.logger.Errorf("Failed to upload banner image: %v", err)
		return nil, err
	}
	s.logger.Infof("Uploaded banner image to URL: %s", url)
	miniApp.BannerImage = url

	cred, _, err := s.CredentialInformationGenrator(UatEnvironment)
	if err != nil {
		s.logger.Errorf("Failed to generate credentials for environment %v: %v", UatEnvironment, err)
		return nil, err
	}

	miniApp.Credential = *cred

	return &miniApp, nil
}

func (s *MiniAppStore) UpdateMiniApp(ctx context.Context, req MiniAppCreateRequest, maker entities.User, id string) (*MiniApp, *MiniApp, error) {
	s.logger.Infof("Updating MiniApp with ID: %s for maker: %v", id, maker)
	prevData, err := s.repository.DetailMiniAppByID(ctx, id)
	if err != nil {
		s.logger.Errorf("Failed to fetch MiniApp with ID %s: %v", id, err)
		return nil, nil, err
	}
	s.logger.Debugf("Fetched previous MiniApp data: %+v", prevData)

	miniApp := buildMiniAppFromRequest(req, id, false)

	miniApp.AppIcon, err = s.uploadOrFallback(ctx, req.AppIcon, "app-icon", s.cfg.MinioEndPoint, prevData.AppIcon)
	if err != nil {
		s.logger.Errorf("Failed to upload app icon for MiniApp %s: %v", id, err)
		return nil, nil, err
	}
	s.logger.Infof("Updated app icon to URL: %s", miniApp.AppIcon)

	miniApp.BannerImage, err = s.uploadOrFallback(ctx, req.BannerImage, "banner-image", s.cfg.MinioEndPoint, prevData.BannerImage)
	if err != nil {
		s.logger.Errorf("Failed to upload banner image for MiniApp %s: %v", id, err)
		return nil, nil, err
	}
	s.logger.Infof("Updated banner image to URL: %s", miniApp.BannerImage)

	return &miniApp, prevData, nil
}

func (s *MiniAppStore) DeleteMiniApp(ctx context.Context, maker entities.User, id string) (*MiniApp, *MiniApp, error) {
	s.logger.Infof("Deleting MiniApp with ID: %s for maker: %v", id, maker)
	miniApp, err := s.repository.DetailMiniAppByID(ctx, id)
	if err != nil {
		s.logger.Errorf("Failed to fetch MiniApp with ID %s: %v", id, err)
		return nil, nil, err
	}

	prev := *miniApp
	miniApp.IsDeleted = true
	miniApp.DeletedAt = time.Now()

	return miniApp, &prev, nil
}

func (s *MiniAppStore) Authorize(ctx context.Context, action *entities.CPSAction) (*entities.CPSAction, error) {
	s.logger.Infof("Authorizing action: %s", action.RequestAction)
	var minApp *MiniApp

	bindErr := common_util.BindAction(action.CurrentAction, &minApp)
	if bindErr != nil {
		s.logger.Errorf("Failed to bind current action to MiniApp: %v", bindErr)
		return nil, fmt.Errorf(common_util.InvalidActionData)
	}

	var err error
	switch requestedAction := action.RequestAction; requestedAction {
	case constant.RequestCreateMiniApp:
		err = s.repository.RunInTransaction(ctx, func(ctx context.Context) error {
			minApp, err = s.repository.CreateMiniApp(ctx, minApp)
			if err != nil {
				return err
			}
			err := s.merchantService.AddMiniApp(ctx, minApp.MerchantID, merchant_service.MiniApps{
				ID:        minApp.ID,
				Enabled:   minApp.Enabled,
				IsDeleted: minApp.IsDeleted,
			})

			if err != nil {
				return err
			}

			return nil
		})
	case constant.RequestUpdateMiniApp:
		minApp, err = s.repository.UpdateMinApp(ctx, minApp)
	case constant.RequestDeleteMiniApp:
		err = s.repository.RunInTransaction(ctx, func(ctx context.Context) error {
			minApp, err = s.repository.DeleteMiniAppAction(ctx, minApp)
			if err != nil {
				s.logger.Errorf("Failed to DeleteMiniAppAction repo mini app: %v", err)

				return err
			}

			err := s.merchantService.SoftDeleteMiniApp(ctx, minApp.MerchantID, minApp.ID)
			if err != nil {
				s.logger.Errorf("Failed to delete mini app: %v", err)
				return err
			}
			return nil
		})
	case constant.RequestEnableMiniApp:
		err = s.repository.RunInTransaction(ctx, func(ctx context.Context) error {
			s.logger.Debugf("Processing EnableMiniApp for ID: %s", minApp.ID)
			minApp, err = s.repository.EnableDisableMiniApp(ctx, minApp.ID, true)
			if err != nil {
				s.logger.Errorf("Failed to EnableDisableMiniApp repo mini app: %v", err)
				return err
			}

			err := s.merchantService.UpdateMiniAppEnabledState(ctx, minApp.MerchantID, minApp.ID, true)
			if err != nil {
				s.logger.Errorf("Failed to enable mini app: %v", err)
				return err
			}
			return nil
		})
	case constant.RequestDisableMiniApp:

		err = s.repository.RunInTransaction(ctx, func(ctx context.Context) error {
			s.logger.Debugf("Processing DisableMiniApp for ID: %s", minApp.ID)
			minApp, err = s.repository.EnableDisableMiniApp(ctx, minApp.ID, false)
			if err != nil {
				s.logger.Errorf("Failed to EnableDisableMiniApp repo mini app: %v", err)
				return err
			}

			err := s.merchantService.UpdateMiniAppEnabledState(ctx, minApp.MerchantID, minApp.ID, false)
			if err != nil {
				s.logger.Errorf("Failed to disable mini app: %v", err)
				return err
			}
			return nil
		})
	default:
		s.logger.Errorf("Unsupported action: %s", requestedAction)
		return nil, fmt.Errorf(common_util.ErrUnsupported)
	}

	if err != nil {
		s.logger.Errorf("Failed to process action %s: %v", action.RequestAction, err)
		s.logger.Debugf("error from domain check miniapp to %s: %v", action.RequestAction, err)
		return nil, err
	}

	action.CurrentAction = minApp
	return action, nil
}

func (s *MiniAppStore) ListMiniApp(ctx context.Context, filterParam *util_constant.Filter) (*common_util.PaginatedResponse[[]*MiniApp], error) {
	s.logger.Infof("Listing MiniApps with filter: %+v", filterParam)
	result, err := s.repository.ListMiniApp(ctx, filterParam)
	if err != nil {
		s.logger.Errorf("Failed to list MiniApps: %v", err)
		return nil, err
	}
	return result, nil
}

func (s *MiniAppStore) DetailMiniAppByID(ctx context.Context, id string) (*MiniApp, error) {
	s.logger.Infof("Fetching MiniApp details for ID: %s", id)
	miniApp, err := s.repository.DetailMiniAppByID(ctx, id)
	if err != nil {
		s.logger.Errorf("Failed to fetch MiniApp with ID %s: %v", id, err)
		return nil, err
	}

	decryptedSecret, err := s.keyGenService.DecryptAppSecret(miniApp.Credential.AppSecret)
	if err != nil {
		s.logger.Errorf("Failed to decrypt AppSecret for MiniApp %s, Environment %v: %v", id, UatEnvironment, err)
		return nil, err
	}

	miniApp.Credential.AppSecret = decryptedSecret
	s.logger.Infof("Successfully fetched MiniApp details for ID: %s", id)
	return miniApp, nil
}

func buildMiniAppFromRequest(req MiniAppCreateRequest, id string, withTimestamps bool) MiniApp {
	now := time.Now()
	productCodes := make([]ProductCode, 0, len(req.ProductCode))
	for _, pc := range req.ProductCode {
		productCodes = append(productCodes, ProductCode{
			ID:             utils.RandomGenerator(20),
			BranchType:     BranchType(pc.BranchType),
			ProductCode:    pc.ProductCode,
			VATCode:        pc.VATCode,
			ServiceFeeCode: pc.ServiceFeeCode,
		})
	}

	credentials := CredentialInformation{
		ID:            utils.RandomGenerator(20),
		Environment:   UatEnvironment,
		MerchantAppID: req.Credential.MerchantAppID,
		FabricAppID:   req.Credential.FabricAppID,
		ShortCode:     req.Credential.ShortCode,
		AppSecret:     req.Credential.AppSecret,
		PrivateKey:    req.Credential.PrivateKey,
		PublicKey:     req.Credential.PublicKey,
	}

	miniApp := MiniApp{
		ID:                  id,
		AppName:             req.AppName,
		CommissionGLAccount: req.CommissionGLAccount,
		AppType:             req.AppType,
		MerchantID:          req.MerchantID,
		ProductCode:         productCodes,
		Credential:          credentials,
		IsEventMiniApp:      req.IsEventMiniApp,
		IsThreeClick:        req.IsThreeClick,
		URL:                 req.URL,
		Stage:               req.Stage,
		AppViewType:         req.AppViewType,
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

func (s *MiniAppStore) EnableDisableMiniApp(ctx context.Context, maker entities.User, id string, enabled bool) (*MiniApp, *MiniApp, error) {
	s.logger.Infof("Setting MiniApp ID %s to enabled=%v for maker: %v", id, enabled, maker)
	miniApp, err := s.repository.DetailMiniAppByID(ctx, id)
	if err != nil {
		s.logger.Errorf("Failed to fetch MiniApp with ID %s: %v", id, err)
		return nil, nil, err
	}

	prev := *miniApp
	if miniApp.Enabled && enabled {
		s.logger.Errorf("MiniApp %s is already enabled", id)
		return nil, &prev, fmt.Errorf(common_util.ErrAlreadyEnabled)
	}
	if !miniApp.Enabled && !enabled {
		s.logger.Errorf("MiniApp %s is already disabled", id)
		return nil, &prev, fmt.Errorf(common_util.ErrAlreadyDisabled)
	}
	miniApp.Enabled = enabled
	miniApp.LastModifiedAt = time.Now()
	s.logger.Infof("MiniApp %s enabled status set to %v", id, enabled)

	return miniApp, &prev, nil
}

func (s *MiniAppStore) CredentialInformationGenrator(envType EnvironmentType) (*CredentialInformation, *string, error) {
	s.logger.Infof("Generating credential information for environment: %v", envType)
	merchnatCode, err := s.keyGenService.GenerateNumericCode(15)
	if err != nil {
		s.logger.Errorf("Failed to generate merchant code: %v", err)
		return nil, nil, err
	}
	s.logger.Debugf("Generated merchant code: %s", merchnatCode)

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
	sorted := fmt.Sprintf("%s|%s|%s|%s", merchnatCode, fabID, miniAppCode, timestamp.Format(time.RFC3339))
	signature, err := s.keyGenService.Sign([]byte(sorted), keys.PrivateKey)
	if err != nil {
		s.logger.Errorf("Failed to sign credential data: %v", err)
		return nil, nil, err
	}
	s.logger.Debugf("Generated signature for environment: %v", envType)

	res := CredentialInformation{
		Environment:   envType,
		MerchantAppID: merchnatCode,
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

func (s *MiniAppStore) uploadOrFallback(ctx context.Context, file *multipart.FileHeader, folder string, minioEndpoint string, prevURL string) (string, error) {
	if file == nil {
		s.logger.Debugf("No file provided for %s, using previous URL: %s", folder, prevURL)
		return prevURL, nil
	}

	s.logger.Infof("Uploading file to MinIO for folder: %s", folder)
	url, err := common_util.UploadFileToMinio(ctx, s.minioClient, s.bucketName, file, folder, minioEndpoint, s.logger)
	if err != nil {
		s.logger.Errorf("Failed to upload file to MinIO for folder %s: %v", folder, err)
		return "", err
	}
	s.logger.Infof("Successfully uploaded file to MinIO, URL: %s", url)
	return url, nil
}
