package wallet

import (
	"cbe-super-app-cps-action/internal/constants"
	walletDto "cbe-super-app-cps-action/internal/constants/dto/wallet"
	"cbe-super-app-cps-action/internal/constants/lib"
	"cbe-super-app-cps-action/internal/constants/localization"
	"cbe-super-app-cps-action/internal/constants/model"
	"cbe-super-app-cps-action/internal/constants/types"
	"cbe-super-app-cps-action/internal/service"
	"cbe-super-app-cps-action/internal/service/wallet/core"
	"cbe-super-app-cps-action/internal/storage"
	local_util "cbe-super-app-cps-action/pkgs/utils"

	"context"
	"errors"
	"time"

	"gitlab.com/bersufekadgetachew/cbe-super-app-shared/shared/config"
	"gitlab.com/bersufekadgetachew/cbe-super-app-shared/shared/utils"
)

type walletService struct {
	repo        storage.WalletRepository
	cpsService  service.CPSActionService
	logger      utils.Logger
	minio       config.MinioClientInterface
	bucketName  string
	minioPubUrl string
	cfg         *config.VaultConfig
}

func NewWalletService(repo storage.WalletRepository, cps service.CPSActionService, minio config.MinioClientInterface, minioPubUrl string, bucketName string, cfg *config.VaultConfig, logger utils.Logger) service.WalletService {
	return &walletService{
		repo:       repo,
		cpsService: cps,
		logger:     logger,
		minio:      minio,
		bucketName: bucketName,
		cfg:        cfg,
	}
}

func (s *walletService) CreateWallet(ctx context.Context, req walletDto.WalletRequest) error {
	s.logger.Infof("CreateWallet called", "wallet_name", req.Name)

	exist, err := s.repo.Find(ctx, "name", req.Name)
	if err != nil {
		return errors.New(localization.ErrorUnhandledServer.Code)
	}

	if exist != nil {
		return errors.New(localization.ErrorWalletNameAlreadyExists.Code)
	}

	code, err := core.GeneratePrefixedName("WAL", req.Code, s.logger)
	if err != nil {
		return errors.New(localization.ErrorUnhandledServer.Code)
	}

	existCode, err := s.repo.Find(ctx, "code", req.Code)
	if err != nil {
		return errors.New(localization.ErrorUnhandledServer.Code)
	}
	if existCode != nil {
		return errors.New(localization.ErrorWalletCodeAlreadyExists.Code)
	}

	URL, err := lib.UploadFileToMinio(ctx, s.minio, s.bucketName, req.Avatar, "wallet", s.cfg.MinioPublicEndPoint, s.logger)
	if err != nil {
		return errors.New(localization.ErrorUnhandledServer.Code)
	}

	wallet := core.ToCreateWalletDoc(req.Name, code, URL, req.Self, req.Other, req.Agent)
	wallet.Enabled = true
	//here since the unique id is nil 000.. use other unique id like the code
	// if err := core.HandleCPSAction(ctx, s.cpsService, wallet.ID.Hex(), constants.RequestCreateWallet, wallet, nil, constants.ActionCreate); err != nil {
	if err := core.HandleCPSAction(ctx, s.cpsService, "", constants.RequestCreateWallet, wallet, nil, constants.ActionCreate); err != nil {
		s.logger.Errorf("CPS action failed for wallet %s: %v", wallet.Code, err)
		return err
	}

	return nil
}

func (s *walletService) UpdateWallet(ctx context.Context, id string, req walletDto.WalletRequest) error {
	s.logger.Infof("UpdateWallet called", "wallet_id", id)

	prevWallet, err := s.repo.FindByID(ctx, id)
	if err != nil {
		return errors.New(localization.ErrorWalletNotFound.Code)
	}

	if req.Name != "" {
		exist, err := s.repo.Find(ctx, "name", req.Name)
		if err != nil {
			return errors.New(localization.ErrorUnhandledServer.Code)
		}
		if exist != nil && exist.ID.Hex() != id {
			return errors.New(localization.ErrorWalletNameAlreadyExists.Code)
		}
	}

	if req.Code != "" {
		exist, err := s.repo.Find(ctx, "code", req.Code)
		if err != nil {
			return errors.New(localization.ErrorUnhandledServer.Code)
		}
		if exist != nil && exist.ID.Hex() != id {
			return errors.New(localization.ErrorWalletCodeAlreadyExists.Code)
		}
	}
	var avatarURL string
	if req.Avatar != nil {
		avatarURL, err = lib.UploadFileToMinio(ctx, s.minio, s.bucketName, req.Avatar, "avatar", s.cfg.MinioPublicEndPoint, s.logger)
		if err != nil {
			return errors.New(localization.ErrorUnhandledServer.Code)
		}
	} else {
		avatarURL = prevWallet.Avatar
	}
	UpdateWallet, change_count := core.ToUpdateWalletDoc(*prevWallet, req)
	if avatarURL != prevWallet.Avatar {
		change_count++
	}
	UpdateWallet.Avatar = avatarURL

	if change_count == 0 {
		return errors.New(localization.ErrorNoChangesDetected.Code)
	}

	if !(UpdateWallet.Services.Agent || UpdateWallet.Services.Self || UpdateWallet.Services.Other) {
		return errors.New(localization.ErrorWalletRechangeOption.Code)
	}

	if err := core.HandleCPSAction(ctx, s.cpsService, id, constants.RequestUpdateWallet, UpdateWallet, *prevWallet, constants.ActionUpdate); err != nil {
		s.logger.Errorf("CPS action failed for wallet %s: %v", UpdateWallet.Code, err)
		return err
	}

	return nil
}

func (s *walletService) DeleteWallet(ctx context.Context, id string) error {
	prevWallet, err := s.repo.FindByID(ctx, id)
	if err != nil {
		return errors.New(localization.ErrorWalletNotFound.Code)
	}

	now := time.Now()
	deletedWallet := *prevWallet
	deletedWallet.IsDeleted = true
	deletedWallet.DeletedAt = now

	if err := core.HandleCPSAction(ctx, s.cpsService, id, constants.RequestDeleteWallet, deletedWallet, *prevWallet, constants.ActionDelete); err != nil {
		s.logger.Errorf("CPS action failed for wallet %s: %v", deletedWallet.Code, err)
		return err
	}

	return nil
}

func (s *walletService) EnableOrDisableWallet(ctx context.Context, id string, enable bool) error {
	prevWallet, err := s.repo.FindByID(ctx, id)
	if err != nil {
		return errors.New(localization.ErrorWalletNotFound.Code)
	}

	if enable && prevWallet.Enabled {
		return errors.New(localization.ErrorWalletAlreadyEnabled.Code)
	}
	if !enable && !prevWallet.Enabled {
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
		s.logger.Errorf("CPS action failed for wallet %s: %v", updatedWallet.Code, err)
		return err
	}

	return nil
}

func (s *walletService) GetWallet(ctx context.Context, id string) (*model.Wallet, error) {
	return s.repo.FindByID(ctx, id)
}

func (s *walletService) GetAllWallet(ctx context.Context, filterParams types.Filter) (*types.PaginatedResponse[[]*model.Wallet], error) {
	return s.repo.FindAllWithPagination(ctx, filterParams)
}

func (s *walletService) Authorize(ctx context.Context, action *model.CPSAction) (*model.CPSAction, error) {

	wallet, err := local_util.JsonUnmarshal[model.Wallet](action.CurrentAction)
	if err != nil {
		return nil, errors.New(localization.ErrorInvalidRequest.Code)
	}

	switch action.RequestAction {
	case string(constants.RequestCreateWallet):
		err = s.repo.Create(ctx, wallet)
	case string(constants.RequestUpdateWallet):
		err = s.repo.Update(ctx, action.UniqueId, wallet)
	case string(constants.RequestDeleteWallet):
		err = s.repo.Delete(ctx, action.UniqueId)
	case string(constants.RequestEnableWallet):
		err = s.repo.EnableOrDisable(ctx, action.UniqueId, true)
	case string(constants.RequestDisableWallet):
		err = s.repo.EnableOrDisable(ctx, action.UniqueId, false)
	default:

		return nil, errors.New(localization.ErrorInvalidRequest.Code)
	}

	if err != nil {
		return nil, err
	}

	action.CurrentAction = wallet
	return action, nil
}
