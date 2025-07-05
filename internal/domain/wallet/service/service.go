package service

import (
	"context"
	"fmt"
	"net/http"
	"time"

	"github.com/CBE-Super-App/cbe-super-app-cps-action/internal/adapter/outbound/model"
	"github.com/CBE-Super-App/cbe-super-app-cps-action/internal/domain/wallet/dto"
	"github.com/CBE-Super-App/cbe-super-app-cps-action/internal/domain/wallet/entity"
	outbound "github.com/CBE-Super-App/cbe-super-app-cps-action/internal/port/outbound/wallet"
	constant "github.com/CBE-Super-App/cbe-super-app-cps-action/utils"

	"gitlab.com/bersufekadgetachew/cbe-super-app-shared/shared/config"
	"gitlab.com/bersufekadgetachew/cbe-super-app-shared/shared/utils"
)

type WalletDomain struct {
	walletRepo  outbound.WalletPersistence
	bucketName  string
	minioClient config.MinioClientInterface
	logger      utils.Logger
}

type WalletService interface {
	GetAllWallet(ctx context.Context, filterParams *constant.Filter) (*entity.WalletResponse, error)
	GetWallet(ctx context.Context, id string) (*entity.Wallet, error)
	CreateWallet(ctx context.Context, req model.CreateCPSAction) (*model.CPSAction, error)
	UpdateWallet(ctx context.Context, id string, req model.CreateCPSAction) (*model.CPSAction, error)
	DeleteWallet(ctx context.Context, id string, req model.CreateCPSAction) (*model.CPSAction, error)
	Authorize(ctx context.Context, req model.AuthorizeCPSAction) (*model.CPSAction, error)
	Reject(ctx context.Context, req model.RejectCPSAction) (*model.CPSAction, error)
	EnableOrDisableWallet(ctx context.Context, id string, requestAction model.RequestAction, cpsReq model.CreateCPSAction) (*model.CPSAction, error)
}

func InitWalletDomain(walletRepo outbound.WalletPersistence, minioClient config.MinioClientInterface,
	bucketName string, logger utils.Logger) WalletService {
	return &WalletDomain{
		walletRepo:  walletRepo,
		minioClient: minioClient,
		bucketName:  bucketName,
		logger:      logger,
	}
}

func (w *WalletDomain) CreateWallet(ctx context.Context, req model.CreateCPSAction) (*model.CPSAction, error) {
	req.RequestAction = model.RequestCreateWallet
	err := w.walletRepo.CPSActionExists(ctx, req)
	if err != nil {
		return nil, err
	}

	actionData, ok := req.ActionData.(dto.CreateWalletRequest)
	if !ok {
		w.logger.Errorf("failed to cast action data to wallet request")
		return nil, fmt.Errorf("failed to create waalet bucket: %w", constant.ErrorDefinition{
			Code:    http.StatusBadRequest,
			Message: "invalid action data",
		})
	}

	if err := actionData.Validate(); err != nil {
		w.logger.Errorf("validation error", err)
		return nil, err
	}

	exist, err := w.minioClient.BucketExist(ctx, w.bucketName)
	if err != nil {
		w.logger.Errorf("failed to check wallet bucket: %v", err)
		return nil, fmt.Errorf("failed to check wallet bucket: %w", constant.ErrorDefinition{
			Code:    http.StatusInternalServerError,
			Message: "internal server error",
		})
	}

	if !exist {
		created, err := w.minioClient.MakeBucket(ctx, w.bucketName)
		if !created || err != nil {
			w.logger.Errorf("failed to create bank bucket: %v", err)
			return nil, fmt.Errorf("failed to create bank bucket: %w", constant.ErrorDefinition{
				Code:    http.StatusInternalServerError,
				Message: "internal server error",
			})
		}
	}

	fileName := fmt.Sprintf("bank-%d-%s", time.Now().UnixNano(), actionData.Avatar.Filename)
	file, err := actionData.Avatar.Open()
	if err != nil {
		w.logger.Errorf("failed to open uploaded file: %v", err)
		return nil, fmt.Errorf("failed to open uploaded file: %w", constant.ErrorDefinition{
			Code:    http.StatusInternalServerError,
			Message: "internal server error",
		})
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
		return nil, fmt.Errorf("upload to MinIO failed: %w", constant.ErrorDefinition{
			Code:    http.StatusInternalServerError,
			Message: "internal server error",
		})
	}

	cpsRes, err := w.walletRepo.CreateWallet(ctx, model.CreateCPSAction{
		MakerUser:  req.MakerUser,
		Department: req.Department,
		ActionData: entity.Wallet{
			Name:   actionData.Name,
			Avatar: fmt.Sprintf("%s/%s", saveObj.Bucket, saveObj.Key),
			Code:   actionData.Code,
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
	cpsAction, err := w.walletRepo.DeleteWallet(ctx, id, req)
	if err != nil {
		return nil, err
	}

	return cpsAction, nil
}

func (w *WalletDomain) GetAllWallet(ctx context.Context, filterParams *constant.Filter) (*entity.WalletResponse, error) {
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

func (w *WalletDomain) Authorize(ctx context.Context, req model.AuthorizeCPSAction) (*model.CPSAction, error) {
	cpsAction, err := w.walletRepo.Authorize(ctx, req)
	if err != nil {
		return nil, err
	}

	return cpsAction, nil
}

func (w *WalletDomain) Reject(ctx context.Context, req model.RejectCPSAction) (*model.CPSAction, error) {
	if err := req.Validate(); err != nil {
		w.logger.Errorf("validation error", err)
		return nil, err
	}
	cpsAction, err := w.walletRepo.Reject(ctx, req)
	if err != nil {
		return nil, err
	}

	return cpsAction, nil
}

func (w *WalletDomain) EnableOrDisableWallet(ctx context.Context, id string,
	requestAction model.RequestAction, cpsReq model.CreateCPSAction) (*model.CPSAction, error) {
	cpsReq.RequestAction = requestAction
	err := w.walletRepo.CPSActionExists(ctx, cpsReq)
	if err != nil {
		return nil, err
	}
	cpsAction, err := w.walletRepo.EnableOrDisableWallet(ctx, id, requestAction, cpsReq)
	if err != nil {
		return nil, err
	}

	return cpsAction, nil
}
