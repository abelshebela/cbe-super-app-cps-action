package wallet

import (
	"cbe-super-app-cps-action/internal/constants"
	walletDto "cbe-super-app-cps-action/internal/constants/dto/wallet"
	"cbe-super-app-cps-action/internal/constants/lib"
	"cbe-super-app-cps-action/internal/constants/localization"
	"cbe-super-app-cps-action/internal/constants/types"
	"cbe-super-app-cps-action/internal/service"
	"cbe-super-app-cps-action/internal/service/wallet/core"
	"cbe-super-app-cps-action/internal/storage"
	local_util "cbe-super-app-cps-action/pkgs/utils"
	"strings"

	"gitlab.com/bersufekadgetachew/cbe-super-app-shared/shared/entities/model"

	"path"

	"context"
	"errors"
	"time"

	"github.com/aws/aws-sdk-go-v2/service/s3"
	"gitlab.com/bersufekadgetachew/cbe-super-app-shared/shared/config"
	"gitlab.com/bersufekadgetachew/cbe-super-app-shared/shared/utils"

	local_model "cbe-super-app-cps-action/internal/constants/model"

	"go.opentelemetry.io/otel/attribute"
	"go.opentelemetry.io/otel/trace"
)

type walletService struct {
	repo        storage.WalletRepository
	cpsService  service.CPSActionService
	serviceRepo storage.ServicesRepository
	logger      utils.Logger
	minio       *s3.Client
	bucketName  string
	minioPubUrl string
	cfg         *config.VaultConfig
}

func NewWalletService(repo storage.WalletRepository, cps service.CPSActionService, serviceRepo storage.ServicesRepository, minio *s3.Client, minioPubUrl string, bucketName string, cfg *config.VaultConfig, logger utils.Logger) service.WalletService {
	return &walletService{
		repo:        repo,
		cpsService:  cps,
		serviceRepo: serviceRepo,
		logger:      logger,
		minio:       minio,
		bucketName:  bucketName,
		cfg:         cfg,
	}
}

func (s *walletService) CreateWallet(ctx context.Context, req walletDto.WalletRequest) error {
	ctx, span := local_util.TraceLogger(ctx, "service", "CreateWallet", "walletService", "walletService")
	defer span.End()
	s.logger.Infof("CreateWallet called", "wallet_name", req.Name)

	exist, err := s.repo.Find(ctx, req.UniqueCode, req.Name)
	if err != nil {
		span.AddEvent("Repo find error", trace.WithAttributes(attribute.String("error", err.Error())))
		return err
	}

	if exist != nil {
		span.AddEvent("Wallet name already exists", trace.WithAttributes(attribute.String("name", req.Name)))
		if strings.TrimSpace(exist.Name) == strings.TrimSpace(req.Name) {
			return errors.New(localization.ErrorWalletNameAlreadyExists.Code)
		}
		if strings.TrimSpace(exist.UniqueCode) == strings.TrimSpace(req.UniqueCode) {
			return errors.New(localization.ErrorWalletCodeAlreadyExists.Code)
		}
		s.logger.Errorf("wallet already exists wallet: %v", exist)
		return errors.New(localization.ErrorWalletServiceIDAlreadyExists.Code)
	}

	code := strings.ToUpper(req.UniqueCode)

	URL, err := lib.UploadFileToMinio(ctx, s.minio, s.bucketName, req.Avatar, string(constants.WalletFolderName), *s.cfg, "", s.logger)
	if err != nil {
		span.AddEvent("UploadFileToMinio error", trace.WithAttributes(attribute.String("error", err.Error())))
		return errors.New(localization.ErrorUnhandledServer.Code)
	}

	wallet := core.ToCreateWalletDoc(req.Name, code, URL, req.ServiceID, req.Self, req.Other, req.Agent)
	wallet.Enabled = false
	//here since the unique id is nil 000.. use other unique id like the code
	// if err := core.HandleCPSAction(ctx, s.cpsService, wallet.ID.Hex(), constants.RequestCreateWallet, wallet, nil, constants.ActionCreate); err != nil {
	if err := core.HandleCPSAction(ctx, s.cpsService, "", constants.RequestCreateWallet, wallet, nil, constants.ActionCreate); err != nil {
		span.AddEvent("CPS action failed", trace.WithAttributes(attribute.String("error", err.Error()), attribute.String("unique_code", wallet.UniqueCode)))
		s.logger.Errorf("CPS action failed for wallet %s: %v", wallet.UniqueCode, err)
		return err
	}

	span.AddEvent("Wallet created", trace.WithAttributes(attribute.String("unique_code", wallet.UniqueCode)))
	return nil
}

func (s *walletService) UpdateWallet(ctx context.Context, id string, req walletDto.WalletRequest) error {
	ctx, span := local_util.TraceLogger(ctx, "service", "UpdateWallet", "walletService", "walletService")
	defer span.End()
	s.logger.Infof("UpdateWallet called", "wallet_id", id)
	prevWallet, err := s.repo.FindByID(ctx, id)
	if err != nil {
		span.AddEvent("FindByID error", trace.WithAttributes(attribute.String("error", err.Error()), attribute.String("id", id)))
		return err
	}

	if req.Name != "" || req.UniqueCode != "" {
		exist, err := s.repo.Find(ctx, req.UniqueCode, req.Name)
		if err != nil {
			span.AddEvent("Repo find error", trace.WithAttributes(attribute.String("error", err.Error())))
			return errors.New(localization.ErrorUnhandledServer.Code)
		}
		if exist != nil && exist.ID.Hex() != id {
			if exist.Name == req.Name {
				span.AddEvent("Wallet name already exists", trace.WithAttributes(attribute.String("name", req.Name)))
				return errors.New(localization.ErrorWalletNameAlreadyExists.Code)
			}
			if exist.UniqueCode == req.UniqueCode {
				span.AddEvent("Wallet code already exists", trace.WithAttributes(attribute.String("unique_code", req.UniqueCode)))
				return errors.New(localization.ErrorWalletCodeAlreadyExists.Code)
			}
		}
	}

	var avatarURL string
	if req.Avatar != nil {
		var objectkey string
		if prevWallet.Avatar != "" {
			objectkey = path.Base(prevWallet.Avatar)
		}

		avatarURL, err = lib.UploadFileToMinio(ctx, s.minio, s.bucketName, req.Avatar, string(constants.WalletFolderName), *s.cfg, objectkey, s.logger)
		if err != nil {
			span.AddEvent("UploadFileToMinio error", trace.WithAttributes(attribute.String("error", err.Error()), attribute.String("id", id)))
			return errors.New(localization.ErrorUnhandledServer.Code)
		}
	} else {
		avatarURL = prevWallet.Avatar
	}
	UpdateWallet, change_count := core.ToUpdateWalletDoc(*prevWallet, req, req.ServiceID)
	if avatarURL != prevWallet.Avatar {
		change_count++
	}
	UpdateWallet.Avatar = avatarURL

	if change_count == 0 {
		span.AddEvent("No changes detected", trace.WithAttributes(attribute.String("id", id)))
		return errors.New(localization.ErrorNoChangesDetected.Code)
	}

	if err := core.HandleCPSAction(ctx, s.cpsService, id, constants.RequestUpdateWallet, UpdateWallet, *prevWallet, constants.ActionUpdate); err != nil {
		span.AddEvent("CPS action failed", trace.WithAttributes(attribute.String("error", err.Error()), attribute.String("unique_code", UpdateWallet.UniqueCode)))
		s.logger.Errorf("CPS action failed for wallet %s: %v", UpdateWallet.UniqueCode, err)
		return err
	}

	span.AddEvent("Wallet updated", trace.WithAttributes(attribute.String("id", id)))
	return nil
}

func (s *walletService) DeleteWallet(ctx context.Context, id string) error {
	ctx, span := local_util.TraceLogger(ctx, "service", "DeleteWallet", "walletService", "walletService")
	defer span.End()
	prevWallet, err := s.repo.FindByID(ctx, id)
	if err != nil {
		span.AddEvent("FindByID error", trace.WithAttributes(attribute.String("error", err.Error()), attribute.String("id", id)))
		return errors.New(localization.ErrorWalletNotFound.Code)
	}

	now := time.Now()
	deletedWallet := *prevWallet
	deletedWallet.IsDeleted = true
	deletedWallet.DeletedAt = now

	if err := core.HandleCPSAction(ctx, s.cpsService, id, constants.RequestDeleteWallet, deletedWallet, *prevWallet, constants.ActionDelete); err != nil {
		span.AddEvent("CPS action failed", trace.WithAttributes(attribute.String("error", err.Error()), attribute.String("unique_code", deletedWallet.UniqueCode)))
		s.logger.Errorf("CPS action failed for wallet %s: %v", deletedWallet.UniqueCode, err)
		return err
	}

	span.AddEvent("Wallet deleted", trace.WithAttributes(attribute.String("id", id)))
	return nil
}

func (s *walletService) EnableOrDisableWallet(ctx context.Context, id string, enable bool) error {
	ctx, span := local_util.TraceLogger(ctx, "service", "EnableOrDisableWallet", "walletService", "walletService")
	defer span.End()
	prevWallet, err := s.repo.FindByID(ctx, id)
	if err != nil {
		span.AddEvent("FindByID error", trace.WithAttributes(attribute.String("error", err.Error()), attribute.String("id", id)))
		return errors.New(localization.ErrorWalletNotFound.Code)
	}

	if enable && prevWallet.Enabled {
		span.AddEvent("Wallet already enabled", trace.WithAttributes(attribute.String("id", id)))
		return errors.New(localization.ErrorWalletAlreadyEnabled.Code)
	}
	if !enable && !prevWallet.Enabled {
		span.AddEvent("Wallet already disabled", trace.WithAttributes(attribute.String("id", id)))
		return errors.New(localization.ErrorWalletAlreadyDisabled.Code)
	}

	updatedWallet := *prevWallet
	updatedWallet.Enabled = enable
	updatedWallet.LastModifiedAt = time.Now()

	var action constants.RequestAction
	if enable {
		action = constants.RequestEnableWallet
	} else {
		action = constants.RequestDisableWallet
	}

	if err := core.HandleCPSAction(ctx, s.cpsService, id, action, updatedWallet, *prevWallet, constants.ActionUpdate); err != nil {
		span.AddEvent("CPS action failed", trace.WithAttributes(attribute.String("error", err.Error()), attribute.String("unique_code", updatedWallet.UniqueCode)))
		s.logger.Errorf("CPS action failed for wallet %s: %v", updatedWallet.UniqueCode, err)
		return err
	}

	span.AddEvent("Wallet enable/disable updated", trace.WithAttributes(attribute.String("id", id), attribute.Bool("enabled", enable)))
	return nil
}

func (s *walletService) GetWallet(ctx context.Context, id string) (*local_model.Wallet, error) {
	return s.repo.FindByID(ctx, id)
}

func (s *walletService) GetAllWallet(ctx context.Context, filterParams types.Filter) (*types.PaginatedResponse[[]local_model.Wallet], error) {
	return s.repo.FindAllWithPaginationForGRPC(ctx, filterParams)
}

// GetAllWalletForGRPC implements service.WalletService.
func (s *walletService) GetAllWalletForGRPC(ctx context.Context, filterParams types.Filter) (*types.PaginatedResponse[[]local_model.Wallet], error) {
	return s.repo.FindAllWithPaginationForGRPC(ctx, filterParams)
}

// GetWalletForGRPC implements service.WalletService.
func (s *walletService) GetWalletForGRPC(ctx context.Context, id string) (*local_model.GRPCWallet, error) {
	return s.repo.FindByIDForGRPC(ctx, id)
}

func (s *walletService) Authorize(ctx context.Context, action *model.CPSAction) (*model.CPSAction, error) {
	ctx, span := local_util.TraceLogger(ctx, "service", "Authorize", "walletService", "walletService")
	defer span.End()

	wallet, err := local_util.JsonUnmarshal[local_model.Wallet](action.CurrentAction)
	if err != nil {
		span.AddEvent("JsonUnmarshal error", trace.WithAttributes(attribute.String("error", err.Error())))
		return nil, errors.New(localization.ErrorInvalidActionData.Code)
	}

	switch action.RequestAction {
	case string(constants.RequestCreateWallet):
		span.AddEvent("Creating wallet", trace.WithAttributes(attribute.String("unique_code", wallet.UniqueCode)))
		err = s.repo.Create(ctx, wallet)
	case string(constants.RequestUpdateWallet):
		span.AddEvent("Updating wallet", trace.WithAttributes(attribute.String("unique_code", wallet.UniqueCode)))
		err = s.repo.Update(ctx, action.UniqueId, wallet)
	case string(constants.RequestDeleteWallet):
		span.AddEvent("Deleting wallet", trace.WithAttributes(attribute.String("unique_code", wallet.UniqueCode)))
		err = s.repo.Delete(ctx, action.UniqueId)
	case string(constants.RequestEnableWallet):
		span.AddEvent("Enabling wallet", trace.WithAttributes(attribute.String("unique_code", wallet.UniqueCode)))
		err = s.repo.EnableOrDisable(ctx, action.UniqueId, true)
	case string(constants.RequestDisableWallet):
		span.AddEvent("Disabling wallet", trace.WithAttributes(attribute.String("unique_code", wallet.UniqueCode)))
		err = s.repo.EnableOrDisable(ctx, action.UniqueId, false)
	default:
		span.AddEvent("Unsupported action", trace.WithAttributes(attribute.String("action", action.RequestAction)))
		return nil, errors.New(localization.ErrorInvalidRequest.Code)
	}

	if err != nil {
		return nil, err
	}

	action.CurrentAction = wallet
	return action, nil
}
