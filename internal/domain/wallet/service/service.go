// Package service provides business logic for wallet operations in the CBE Super App.
package service

import (
	"context"
	"fmt"
	"time"

	"github.com/CBE-Super-App/cbe-super-app-cps-action/internal/adapter/outbound/model"
	actions "github.com/CBE-Super-App/cbe-super-app-cps-action/internal/domain/cps_actions/constant"
	entities "github.com/CBE-Super-App/cbe-super-app-cps-action/internal/domain/cps_actions/entities"
	"github.com/CBE-Super-App/cbe-super-app-cps-action/internal/domain/wallet/dto"
	"github.com/CBE-Super-App/cbe-super-app-cps-action/internal/domain/wallet/entity"
	outbound "github.com/CBE-Super-App/cbe-super-app-cps-action/internal/port/outbound/wallet"
	error_codes "github.com/CBE-Super-App/cbe-super-app-cps-action/pkgs/utils"

	// utils "github.com/CBE-Super-App/cbe-super-app-cps-action/pkgs/utils"
	constant "github.com/CBE-Super-App/cbe-super-app-cps-action/utils"
	config "gitlab.com/bersufekadgetachew/cbe-super-app-shared/shared/config"
	"gitlab.com/bersufekadgetachew/cbe-super-app-shared/shared/utils"
)

type WalletDomain struct {
	walletRepo  outbound.WalletPersistence
	bucketName  string
	minioClient config.MinioClientInterface

	logger utils.Logger
	cfg    *config.VaultConfig
}

type WalletService interface {
	GetAllWallet(ctx context.Context, filterParams *constant.Filter) (*error_codes.PaginatedResponse[[]*entity.Wallet], error)

	GetWallet(ctx context.Context, id string) (*entity.Wallet, error)
	CreateWallet(ctx context.Context, req model.CreateCPSAction) (*model.CPSAction, error)
	UpdateWallet(ctx context.Context, id string, req model.CreateCPSAction) (*model.CPSAction, error)
	DeleteWallet(ctx context.Context, id string, req model.CreateCPSAction) (*model.CPSAction, error)
	Authorize(ctx context.Context, action *entities.CPSAction) (*entities.CPSAction, error)
	EnableOrDisableWallet(ctx context.Context, id string, requestAction model.RequestAction, cpsReq model.CreateCPSAction) (*model.CPSAction, error)
}

func InitWalletDomain(walletRepo outbound.WalletPersistence, minioClient config.MinioClientInterface,

	bucketName string, logger utils.Logger, cfg *config.VaultConfig) WalletService {

	return &WalletDomain{
		walletRepo:  walletRepo,
		minioClient: minioClient,
		bucketName:  bucketName,
		logger:      logger,
		cfg:         cfg,
	}
}

func (w *WalletDomain) CreateWallet(ctx context.Context, req model.CreateCPSAction) (*model.CPSAction, error) {
	req.RequestAction = model.RequestCreateWallet
	err := w.walletRepo.CPSActionExists(ctx, req)
	if err != nil {
		return nil, err
	}

	exists, actionData, err := w.walletAlreadyExists(ctx, req.ActionData)
	if err := actionData.Validate(); err != nil {
		w.logger.Errorf("validation error", err)
		return nil, err
	}

	if err != nil {
		return nil, err
	}
	if exists {
		return nil, fmt.Errorf("WALLET_ALREADY_EXISTS")
	}

	exist, err := w.minioClient.BucketExist(ctx, w.bucketName)
	if err != nil {
		w.logger.Errorf("failed to check wallet bucket: %v", err)
		return nil, fmt.Errorf(error_codes.UnhandledServerError)
	}

	if !exist {
		created, err := w.minioClient.MakeBucket(ctx, w.bucketName)
		if !created || err != nil {
			w.logger.Errorf("failed to create bank bucket: %v", err)
			return nil, fmt.Errorf(error_codes.UnhandledServerError)
		}
	}

	fileName := fmt.Sprintf("bank-%d-%s", time.Now().UnixNano(), actionData.Avatar.Filename)
	file, err := actionData.Avatar.Open()
	if err != nil {
		w.logger.Errorf("failed to open uploaded file: %v", err)
		return nil, fmt.Errorf(error_codes.UnhandledServerError)
	}
	defer file.Close()

	saveObj, err := w.minioClient.SaveObjectN(ctx, config.SaveObjectBodyN{
		BucketName:  w.bucketName,
		ObjectName:  fileName,
		Reader:      file,
		Size:        actionData.Avatar.Size,
		ContentType: config.ContentType(actionData.Avatar.Header.Get("Content-Type")),
	})

	if err != nil {
		w.logger.Errorf("failed to save object to MinIO: %v", err)
		return nil, fmt.Errorf(error_codes.UnhandledServerError)
	}

	// _ := fmt.Sprintf("%s/%s/%s", w.cfg.MinioEndPoint, saveObj.Bucket, saveObj.Key)

	cpsRes, err := w.walletRepo.CreateWallet(ctx, model.CreateCPSAction{
		MakerUser:  req.MakerUser,
		Department: req.Department,
		ActionData: entity.Wallet{
			Name:           actionData.Name,
			Avatar:         fmt.Sprintf("%s/%s", saveObj.Bucket, saveObj.Key),
			Code:           actionData.Code,
			CreatedAt:      time.Now(),
			LastModifiedAt: time.Now(),
		},
	})

	if err != nil {
		return nil, err
	}

	return cpsRes, nil
}

func (w *WalletDomain) UpdateWallet(ctx context.Context, id string, req model.CreateCPSAction) (*model.CPSAction, error) {
	req.RequestAction = model.RequestUpdateWallet
	err := w.walletRepo.CPSActionExists(ctx, req)
	if err != nil {
		return nil, err
	}

	_, err = w.walletRepo.GetWallet(ctx, id)
	if err != nil {
		return nil, err
	}

	cpsAction, err := w.walletRepo.UpdateWallet(ctx, id, req)
	if err != nil {
		return nil, err
	}

	return cpsAction, nil
}

func (w *WalletDomain) DeleteWallet(ctx context.Context, id string, req model.CreateCPSAction) (*model.CPSAction, error) {
	req.RequestAction = model.RequestDeleteWallet
	err := w.walletRepo.CPSActionExists(ctx, req)
	if err != nil {
		return nil, err
	}

	_, err = w.walletRepo.GetWallet(ctx, id)
	if err != nil {
		return nil, err
	}

	cpsAction, err := w.walletRepo.DeleteWallet(ctx, id, req)
	if err != nil {
		return nil, err
	}

	return cpsAction, nil
}

func (w *WalletDomain) GetAllWallet(ctx context.Context, filterParams *constant.Filter) (*error_codes.PaginatedResponse[[]*entity.Wallet], error) {
	banks, err := w.walletRepo.GetAllWallet(ctx, filterParams)
	if err != nil {
		return nil, err
	}

	return banks, nil

}

func (w *WalletDomain) GetWallet(ctx context.Context, id string) (*entity.Wallet, error) {
	bank, err := w.walletRepo.GetWallet(ctx, id)
	if err != nil {
		return nil, err
	}

	return bank, nil
}

func (w *WalletDomain) Authorize(ctx context.Context, action *entities.CPSAction) (*entities.CPSAction, error) {
	actionData, prevData, err := w.walletRepo.ExtractActionData(action)
	if err != nil {
		return nil, err
	}

	action.MakerActionTime = time.Now()
	action.LastModifiedAt = action.MakerActionTime

	switch action.ActionType {
	case actions.ActionCreate:
		return w.walletRepo.AuthorizeCreate(ctx, action, actionData)

	case actions.ActionUpdate:
		return w.walletRepo.AuthorizeUpdate(ctx, action, actionData, prevData)

	case actions.ActionDelete:
		return w.walletRepo.AuthorizeDelete(ctx, action, prevData)

	default:
		return nil, fmt.Errorf("UNHANDLED_SERVER_ERROR")
	}
}

func (w *WalletDomain) EnableOrDisableWallet(ctx context.Context, id string, requestAction model.RequestAction, cpsReq model.CreateCPSAction) (*model.CPSAction, error) {
	cpsReq.RequestAction = requestAction
	err := w.walletRepo.CPSActionExists(ctx, cpsReq)
	if err != nil {
		return nil, err
	}

	_, err = w.walletRepo.GetWallet(ctx, id)
	if err != nil {
		return nil, err
	}

	isEnabled, err := w.walletRepo.CheckIsEnabled(ctx, id)

	if err != nil {
		return nil, err
	}

	switch {
	case requestAction == model.RequestEnableWallet && isEnabled:
		return nil, fmt.Errorf(error_codes.WalletAlreadyEnabled)
	case requestAction == model.RequestDisableWallet && !isEnabled:
		return nil, fmt.Errorf(error_codes.WalletAlreadyDisabled)
	case requestAction != model.RequestEnableWallet && requestAction != model.RequestDisableWallet:
		return nil, fmt.Errorf(error_codes.InvalidRequestAction)
	}

	cpsAction, err := w.walletRepo.EnableOrDisableWallet(ctx, id, requestAction, cpsReq)
	if err != nil {
		return nil, err
	}

	return cpsAction, nil
}

func (w *WalletDomain) walletAlreadyExists(ctx context.Context, actionData any) (bool, *dto.CreateWalletRequest, error) {
	walletReq, ok := actionData.(dto.CreateWalletRequest)
	if !ok {
		w.logger.Errorf("failed to cast action data to CreateWalletRequest")
		return false, nil, fmt.Errorf(error_codes.InvalidActionData)
	}

	exists, err := w.walletRepo.CheckWalletExists(ctx, entity.CheckWallet{
		Name: walletReq.Name,
		Code: walletReq.Code,
	})
	if err != nil {
		return false, nil, err
	}

	return exists, &walletReq, nil
}
