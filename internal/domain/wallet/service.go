// Package service provides business logic for wallet operations in the CBE Super App.
package wallet

import (
	"context"
	"fmt"
	"time"

	entities "github.com/CBE-Super-App/cbe-super-app-cps-action/internal/domain/cps_actions/entities"
	"github.com/CBE-Super-App/cbe-super-app-cps-action/pkgs/utils"
	"gitlab.com/bersufekadgetachew/cbe-super-app-shared/shared/config"

	cps_const "github.com/CBE-Super-App/cbe-super-app-cps-action/internal/domain/cps_actions/constant"
	common_util "github.com/CBE-Super-App/cbe-super-app-cps-action/pkgs/utils"
	constant "github.com/CBE-Super-App/cbe-super-app-cps-action/utils"

	shared_util "gitlab.com/bersufekadgetachew/cbe-super-app-shared/shared/utils"
)

type WalletService interface {
	CreateWallet(ctx context.Context, wallet WalletRequest) (*Wallet, error)
	UpdateWallet(ctx context.Context, id string, wallet WalletRequest) (*Wallet, *Wallet, error)
	DeleteWallet(ctx context.Context, id string) (*Wallet, *Wallet, error)
	EnableDisableWallet(ctx context.Context, id string, enable bool) (*Wallet, *Wallet, error)
	FetchWalletByID(ctx context.Context, id string) (*Wallet, error)
	FetchWallet(ctx context.Context, filterParam *constant.Filter) (*utils.PaginatedResponse[[]*Wallet], error)
	Authorize(ctx context.Context, action *entities.CPSAction) (*entities.CPSAction, error)
}

type Service struct {
	Repository WalletRepository
	logger     shared_util.Logger
	minio      config.MinioClientInterface
	bucketName string
	cfg        *config.VaultConfig
}

func NewWalletService(repository WalletRepository,
	minio config.MinioClientInterface,
	bucketName string,
	cfg *config.VaultConfig, logger shared_util.Logger) WalletService {
	return &Service{
		Repository: repository,
		logger:     logger,
		minio:      minio,
		bucketName: bucketName,
		cfg:        cfg,
	}
}

func (s *Service) CreateWallet(ctx context.Context, req WalletRequest) (*Wallet, error) {
	exist, err := s.Repository.WalletNameExists(ctx, req.Name, req.Code, nil)
	if err != nil {
		return nil, err
	}

	if exist {
		return nil, fmt.Errorf(utils.InformationAlreadyExistst)
	}

	code, err := utils.GeneratePrefixedName("WAL", req.Name, s.logger)
	if err != nil {
		return nil, fmt.Errorf(utils.UnhandledServerError)
	}

	URL, err := common_util.UploadFileToMinio(ctx, s.minio, s.bucketName, req.Avatar, "avatar", s.cfg.MinioEndPoint, s.logger)
	if err != nil {
		return nil, err
	}

	result := Wallet{
		Name:           req.Name,
		Code:           code,
		Avatar:         URL,
		CreatedAt:      time.Now(),
		LastModifiedAt: time.Now(),
	}

	return &result, nil
}

func (s *Service) UpdateWallet(ctx context.Context, id string, req WalletRequest) (*Wallet, *Wallet, error) {
	s.logger.Infof("Updating wallet", "id", id)

	prevWallet, err := s.Repository.FetchWalletByID(ctx, id)
	if err != nil {
		s.logger.Errorf("Failed to fetch wallet", "id", id, "error", err)
		return nil, nil, err
	}

	if req.Name != "" {
		exist, err := s.Repository.WalletNameExists(ctx, req.Name, req.Code, &id)
		if err != nil {
			return nil, nil, err
		}
		if exist {
			return nil, nil, fmt.Errorf(utils.InformationAlreadyExistst)
		}
	}

	var url string
	if req.Avatar != nil {
		URL, err := common_util.UploadFileToMinio(ctx, s.minio, s.bucketName, req.Avatar, "avatar", s.cfg.MinioEndPoint, s.logger)
		if err != nil {
			return nil, nil, err
		}

		url = URL

	}

	curWallet := Wallet{
		ID:             prevWallet.ID,
		Name:           nonEmptyString(req.Name, prevWallet.Name),
		Code:           prevWallet.Code,
		Avatar:         nonEmptyString(url, prevWallet.Avatar),
		CreatedAt:      prevWallet.CreatedAt,
		LastModifiedAt: time.Now(),
	}

	s.logger.Infof("Updated wallet", "id", id)
	return &curWallet, prevWallet, nil
}

func (s *Service) DeleteWallet(ctx context.Context, id string) (*Wallet, *Wallet, error) {
	s.logger.Infof("Deleting wallet", "id", id)

	prevWallet, err := s.Repository.FetchWalletByID(ctx, id)
	if err != nil {
		s.logger.Errorf("Failed to fetch wallet", "id", id, "error", err)
		return nil, nil, err
	}

	now := time.Now()
	deletedWallet := prevWallet
	deletedWallet.IsDeleted = true
	deletedWallet.DeletedAt = &now

	return deletedWallet, prevWallet, nil
}

func (s *Service) EnableDisableWallet(ctx context.Context, id string, enable bool) (*Wallet, *Wallet, error) {
	s.logger.Infof("EnableDisable wallet", "id", id)

	prevWallet, err := s.Repository.FetchWalletByID(ctx, id)
	if err != nil {
		s.logger.Errorf("Failed to fetch wallet", "id", id, "error", err)
		return nil, nil, err
	}

	if enable && prevWallet.Enabled {
		return nil, nil, fmt.Errorf(utils.ErrAlreadyEnabled)
	}

	if !enable && !prevWallet.Enabled {
		return nil, nil, fmt.Errorf(utils.ErrAlreadyDisabled)
	}

	updatedWallet := prevWallet
	updatedWallet.Enabled = enable
	updatedWallet.LastModifiedAt = time.Now()

	return updatedWallet, prevWallet, nil
}

func (s *Service) FetchWalletByID(ctx context.Context, id string) (*Wallet, error) {
	s.logger.Infof("FetchWalletByID", "id", id)
	return s.Repository.FetchWalletByID(ctx, id)
}

func (s *Service) FetchWallet(ctx context.Context, filterParam *constant.Filter) (*utils.PaginatedResponse[[]*Wallet], error) {
	return s.Repository.FetchWallet(ctx, filterParam)
}

func nonEmptyString(new, old string) string {
	if new != "" {
		return new
	}
	return old
}

func (s *Service) Authorize(ctx context.Context, action *entities.CPSAction) (*entities.CPSAction, error) {
	requestedAction := action.RequestAction

	var wallet *Wallet
	var err error

	bindErr := common_util.BindAction(action.CurrentAction, &wallet)
	if bindErr != nil {
		s.logger.Errorf("failed to bind current action to wallet: %v", bindErr)
		return nil, fmt.Errorf(common_util.InvalidActionData)
	}

	switch requestedAction {
	case cps_const.RequestCreateWallet:
		wallet, err = s.Repository.CreateWallet(ctx, *wallet)

		if err != nil {
			return nil, err
		}

	case cps_const.RequestUpdateWallet:
		wallet, err = s.Repository.UpdateWallet(ctx, *wallet)

		if err != nil {
			return nil, err
		}

	case cps_const.RequestDeleteWallet:
		wallet, err = s.Repository.CreateWallet(ctx, *wallet)

		if err != nil {
			return nil, err
		}
	case cps_const.RequestEnableWallet:
		wallet, err = s.Repository.EnableDisableWallet(ctx, wallet.ID, true)

		if err != nil {
			return nil, err
		}
	case cps_const.RequestDisableWallet:
		wallet, err = s.Repository.EnableDisableWallet(ctx, wallet.ID, false)

		if err != nil {
			return nil, err
		}

	default:
		return nil, fmt.Errorf(common_util.ErrUnsupported)
	}

	action.CurrentAction = wallet
	return action, nil
}
