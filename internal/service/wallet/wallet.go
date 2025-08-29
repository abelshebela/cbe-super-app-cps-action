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
	"context"
	"errors"
	"time"

	"gitlab.com/bersufekadgetachew/cbe-super-app-shared/shared/config"
	"gitlab.com/bersufekadgetachew/cbe-super-app-shared/shared/utils"
)

type walletService struct {
	repo       storage.WalletRepository
	cpsService service.CPSActionService
	logger     utils.Logger
	minio      config.MinioClientInterface
	bucketName string
	cfg        *config.VaultConfig
}

func NewWalletService(repo storage.WalletRepository, cps service.CPSActionService, minio config.MinioClientInterface, bucketName string, cfg *config.VaultConfig, logger utils.Logger) service.WalletService {
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

	exist, err := s.repo.Find(ctx, req.Name)
	if err != nil {
		return errors.New(localization.ErrorUnhandledServer.Code)
	}
	if exist != nil {
		return errors.New(localization.ErrorWalletAlreadyExists.Code)
	}

	code, err := core.GeneratePrefixedName("WAL", req.Name, s.logger)
	if err != nil {
		return errors.New(localization.ErrorUnhandledServer.Code)
	}

	URL, err := lib.UploadFileToMinio(ctx, s.minio, s.bucketName, req.Avatar, "avatar", s.cfg.MinioEndPoint, s.logger)
	if err != nil {
		return errors.New(localization.ErrorUnhandledServer.Code)
	}

	wallet := core.ToWalletDoc(req.Name, code, URL)
	if err := core.HandleCPSAction(ctx, s.cpsService, wallet.ID.Hex(), constants.RequestCreateWallet, wallet, nil, constants.ActionCreate); err != nil {
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
		exist, err := s.repo.Find(ctx, req.Name)
		if err != nil {
			return errors.New(localization.ErrorUnhandledServer.Code)
		}
		if exist != nil {
			return errors.New(localization.ErrorWalletAlreadyExists.Code)
		}
	}

	var avatarURL string
	if req.Avatar != nil {
		avatarURL, err = lib.UploadFileToMinio(ctx, s.minio, s.bucketName, req.Avatar, "avatar", s.cfg.MinioEndPoint, s.logger)
		if err != nil {
			return errors.New(localization.ErrorUnhandledServer.Code)
		}
	} else {
		avatarURL = prevWallet.Avatar
	}

	updatedWallet := model.Wallet{
		ID:        prevWallet.ID,
		Name:      req.Name,
		Code:      req.Code,
		Avatar:    avatarURL,
		Enabled:   prevWallet.Enabled,
		IsDeleted: prevWallet.IsDeleted,
		CreatedAt: prevWallet.CreatedAt,
		DeletedAt: prevWallet.DeletedAt,
	}

	if err := core.HandleCPSAction(ctx, s.cpsService, id, constants.RequestUpdateWallet, updatedWallet, *prevWallet, constants.ActionUpdate); err != nil {
		s.logger.Errorf("CPS action failed for wallet %s: %v", updatedWallet.Code, err)
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
	wallet, ok := action.CurrentAction.(*model.Wallet)
	if !ok || wallet == nil {
		return nil, errors.New(localization.ErrorInvalidRequest.Code)
	}

	var err error
	switch action.RequestAction {
	case string(constants.RequestCreateWallet):
		err = s.repo.Create(ctx, wallet)
	case string(constants.RequestUpdateWallet):
		err = s.repo.Update(ctx, wallet.ID.Hex(), wallet)
	case string(constants.RequestDeleteWallet):
		err = s.repo.Delete(ctx, wallet.ID.Hex())
	case string(constants.RequestEnableWallet):
		err = s.repo.EnableOrDisable(ctx, wallet.ID.Hex(), true)
	case string(constants.RequestDisableWallet):
		err = s.repo.EnableOrDisable(ctx, wallet.ID.Hex(), false)
	default:
		return nil, errors.New(localization.ErrorInvalidRequest.Code)
	}

	if err != nil {
		return nil, err
	}

	action.CurrentAction = wallet
	return action, nil
}
