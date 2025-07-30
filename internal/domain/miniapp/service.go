package miniapp

import (
	"context"
	"encoding/base64"
	"fmt"
	"time"

	"github.com/CBE-Super-App/cbe-super-app-cps-action/internal/domain/cps_actions/constant"
	"gitlab.com/bersufekadgetachew/cbe-super-app-shared/shared/utils"

	entities "github.com/CBE-Super-App/cbe-super-app-cps-action/internal/domain/cps_actions/entities"
	keyGen "github.com/CBE-Super-App/cbe-super-app-cps-action/pkgs/keygen"
	common_util "github.com/CBE-Super-App/cbe-super-app-cps-action/pkgs/utils"

	util_constant "github.com/CBE-Super-App/cbe-super-app-cps-action/utils"
	"gitlab.com/bersufekadgetachew/cbe-super-app-shared/shared/config"
)

type MiniAppStore struct {
	repository    MiniRepository
	logger        utils.Logger
	cfg           *config.VaultConfig
	bucketName    string
	minioClient   config.MinioClientInterface
	keyGenService keyGen.KeyGeneratorService
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
	keyGenService keyGen.KeyGeneratorService, logger utils.Logger) MiniAppService {
	return &MiniAppStore{
		repository:    repository,
		logger:        logger,
		cfg:           cfg,
		bucketName:    bucketName,
		minioClient:   minioClient,
		keyGenService: keyGenService,
	}
}

func (s *MiniAppStore) CreateMiniApp(ctx context.Context, req MiniAppCreateRequest, maker entities.User) (*MiniApp, error) {
	miniApp := buildMiniAppFromRequest(req, "", true)
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
		miniApp.AppIcon = url
	}
	// Generate credentials for all environments
	envs := []EnvironmentType{DevEnvironment, TestEnvironment, UatEnvironment, ProductionEnvironment}
	miniApp.Credential = make([]CredentialInformation, 0, len(envs))
	for _, env := range envs {
		cred, _, err := s.CredentialInformationGenrator(env)
		if err != nil {
			return nil, err
		}
		miniApp.Credential = append(miniApp.Credential, *cred)
	}

	return &miniApp, nil
}

func (s *MiniAppStore) UpdateMiniApp(ctx context.Context, req MiniAppCreateRequest, maker entities.User, id string) (*MiniApp, *MiniApp, error) {
	prevData, err := s.repository.DetailMiniAppByID(ctx, id)
	if err != nil {
		return nil, nil, err
	}
	miniApp := buildMiniAppFromRequest(req, id, false)
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
			return nil, nil, err
		}
		miniApp.AppIcon = url
	} else {
		miniApp.AppIcon = prevData.AppIcon
	}

	return &miniApp, prevData, nil
}

func (s *MiniAppStore) DeleteMiniApp(ctx context.Context, maker entities.User, id string) (*MiniApp, *MiniApp, error) {
	miniApp, err := s.repository.DetailMiniAppByID(ctx, id)
	if err != nil {
		return nil, nil, err
	}
	prev := *miniApp
	miniApp.IsDeleted = true
	miniApp.DeletedAt = time.Now()

	return miniApp, &prev, nil
}
func (s *MiniAppStore) Authorize(ctx context.Context, action *entities.CPSAction) (*entities.CPSAction, error) {
	requestedAction := action.RequestAction
	var minApp *MiniApp

	bindErr := bindAction(action.CurrentAction, &minApp)
	if bindErr != nil {
		s.logger.Errorf("failed to bind current action to MiniApp: %v", bindErr)
		return nil, fmt.Errorf(common_util.InvalidActionData)
	}

	var err error
	switch requestedAction {
	case constant.RequestCreateMiniApp:
		minApp, err = s.repository.CreateMiniApp(ctx, minApp)
	case constant.RequestUpdateMiniApp:
		minApp, err = s.repository.UpdateMinApp(ctx, minApp)
	case constant.RequestDeleteMiniApp:
		minApp, err = s.repository.DeleteMiniAppAction(ctx, minApp)
	case constant.RequestEnableMiniApp:
		minApp, err = s.repository.EnableDisableMiniApp(ctx, minApp.ID, true)
	case constant.RequestDisableMiniApp:
		minApp, err = s.repository.EnableDisableMiniApp(ctx, minApp.ID, false)
	default:
		return nil, fmt.Errorf(common_util.ErrUnsupported)
	}

	if err != nil {
		fmt.Printf("error form domain chekmiiapp to crate mini app : %v", err)
		return nil, err
	}

	action.CurrentAction = minApp
	return action, nil
}
func (s *MiniAppStore) ListMiniApp(ctx context.Context, filterParam *util_constant.Filter) (*common_util.PaginatedResponse[[]*MiniApp], error) {
	return s.repository.ListMiniApp(ctx, filterParam)
}

func (s *MiniAppStore) DetailMiniAppByID(ctx context.Context, id string) (*MiniApp, error) {
	miniApp, err := s.repository.DetailMiniAppByID(ctx, id)
	if err != nil {
		return nil, err
	}

	for _, cred := range miniApp.Credential {
		decryptedSecret, err := s.keyGenService.DecryptAppSecret(cred.AppSecret)
		if err != nil {
			return nil, err
		}
		cred.AppSecret = decryptedSecret
	}

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

	credentials := make([]CredentialInformation, 0, len(req.Credential))
	for _, cred := range req.Credential {
		credentials = append(credentials, CredentialInformation{
			ID:            utils.RandomGenerator(20),
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
		MPAASID:             req.MPAASID,
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
	miniApp, err := s.repository.DetailMiniAppByID(ctx, id)
	if err != nil {
		return nil, nil, err
	}
	prev := *miniApp
	if miniApp.Enabled && enabled {
		return nil, &prev, fmt.Errorf(common_util.ErrAlreadyEnabled)
	}
	if !miniApp.Enabled && !enabled {
		return nil, &prev, fmt.Errorf(common_util.ErrAlreadyDisabled)
	}
	miniApp.Enabled = enabled
	miniApp.LastModifiedAt = time.Now()

	return miniApp, &prev, nil
}

func (s *MiniAppStore) CredentialInformationGenrator(envType EnvironmentType) (*CredentialInformation, *string, error) {
	merchnatCode, err := s.keyGenService.GenerateNumericCode(15)
	if err != nil {
		return nil, nil, err
	}
	fabID := s.keyGenService.GenerateFabricID()
	shortCode, err := s.keyGenService.GenerateNumericCode(6)
	if err != nil {
		return nil, nil, err
	}

	miniAppCode, err := s.keyGenService.GenerateNumericCode(6)
	if err != nil {
		return nil, nil, err
	}
	keys, err := s.keyGenService.GenerateKeyPair()
	if err != nil {
		return nil, nil, err
	}

	rawSecret, err := s.keyGenService.GenerateAppSecret(32)
	if err != nil {
		return nil, nil, err
	}

	encryptedSecret, err := s.keyGenService.EncryptAppSecret(rawSecret)
	if err != nil {
		return nil, nil, err
	}

	timestamp := time.Now().UTC()
	sorted := fmt.Sprintf("%s|%s|%s|%s", merchnatCode, fabID, miniAppCode, timestamp.Format(time.RFC3339))
	signature, err := s.keyGenService.Sign([]byte(sorted), keys.PrivateKey)

	if err != nil {
		return nil, nil, err
	}

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

	return &res, &rawSecret, nil
}
