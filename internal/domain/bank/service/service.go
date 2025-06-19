package service

import (
	"context"
	"fmt"
	"net/http"
	"os"
	"path/filepath"
	"time"

	"gitlab.com/bersufekadgetachew/cbe-super-app-cps-action/internal/adapter/outbound/model"
	"gitlab.com/bersufekadgetachew/cbe-super-app-cps-action/internal/domain/bank/dto"
	"gitlab.com/bersufekadgetachew/cbe-super-app-cps-action/internal/domain/bank/entity"
	outbound "gitlab.com/bersufekadgetachew/cbe-super-app-cps-action/internal/port/outbound/bank"
	constant "gitlab.com/bersufekadgetachew/cbe-super-app-cps-action/utils"
	"gitlab.com/bersufekadgetachew/cbe-super-app-shared/shared/config"
	"gitlab.com/bersufekadgetachew/cbe-super-app-shared/shared/utils"
)

type BankDomain struct {
	bankRepo    outbound.BankPersistence
	bucketName  string
	minioClient config.MinioClientInterface
	logger      utils.Logger
}

type BankService interface {
	GetAllBank(ctx context.Context, filterParams *constant.Filter) (*entity.BankResponse, error)
	GetOneBank(ctx context.Context, id string) (*entity.Bank, error)
	CreateOneBank(ctx context.Context, req model.CreateCPSAction) (*model.CpsAction, error)
	UpdateOneBank(ctx context.Context, id string, req model.CreateCPSAction) (*model.CpsAction, error)
	DeleteOneBank(ctx context.Context, id string, req model.CreateCPSAction) (*model.CpsAction, error)
	Authorize(ctx context.Context, req model.AuthorizeCPSAction) (*model.CpsAction, error)
	Reject(ctx context.Context, req model.RejectCPSAction) (*model.CpsAction, error)
	EnableOrDisableWallet(ctx context.Context, id string,
		requestAction model.RequestAction, cpsReq model.CreateCPSAction) (*model.CpsAction, error)
}

func InitBankDomain(bankRepo outbound.BankPersistence, minioClient config.MinioClientInterface,
	bucketName string, logger utils.Logger) BankService {
	return &BankDomain{
		bankRepo:    bankRepo,
		minioClient: minioClient,
		bucketName:  bucketName,
		logger:      logger,
	}
}

func (b *BankDomain) CreateOneBank(ctx context.Context, req model.CreateCPSAction) (*model.CpsAction, error) {
	actionData, ok := req.ActionData.(dto.CreateBankRequest)
	if !ok {
		b.logger.Errorf("failed to cast action data to bank request")
		return nil, fmt.Errorf("failed to create bank bucket: %w", constant.ErrorDefinition{
			Code:    http.StatusBadRequest,
			Message: "invalid action data",
		})
	}

	if err := actionData.Validate(); err != nil {
		b.logger.Errorf("validation error", err)
		return nil, err
	}

	exist, err := b.minioClient.BucketExist(ctx, b.bucketName)
	if err != nil {
		b.logger.Errorf("failed to check bank bucket: %v", err)
		return nil, fmt.Errorf("failed to check bank bucket: %w", constant.ErrorDefinition{
			Code:    http.StatusInternalServerError,
			Message: "internal server error",
		})
	}

	if !exist {
		created, err := b.minioClient.MakeBucket(ctx, b.bucketName)
		if !created || err != nil {
			b.logger.Errorf("failed to create bank bucket: %v", err)
			return nil, fmt.Errorf("failed to create bank bucket: %w", constant.ErrorDefinition{
				Code:    http.StatusInternalServerError,
				Message: "internal server error",
			})
		}
	}

	fileName := fmt.Sprintf("bank-%d-%s", time.Now().UnixNano(), actionData.Logo.Filename)
	dir, err := os.Getwd()
	if err != nil {
		b.logger.Errorf("failed to get current working directory", err)
		return nil, fmt.Errorf("failed to create advert bucket: %w", constant.ErrorDefinition{
			Code:    http.StatusInternalServerError,
			Message: "internal server error",
		})
	}
	filePath := filepath.Join(dir, fileName)

	tempFile, err := os.Create(filePath)
	if err != nil {
		b.logger.Errorf("failed to create temp file: %v", err)
		return nil, fmt.Errorf("failed to create temp file: %w", constant.ErrorDefinition{
			Code:    http.StatusInternalServerError,
			Message: "internal server error",
		})
	}
	defer func() {
		tempFile.Close()
		os.Remove(filePath)
	}()

	// Save to MinIO
	saveObj, err := b.minioClient.SaveObject(ctx, config.SaveObjectBody{
		BucketName: b.bucketName,
		ObjectName: fileName,
		File:       filePath,
	})
	if err != nil {
		b.logger.Errorf("failed to save object to MinIO: %v", err)
		return nil, fmt.Errorf("failed to save object to MinIO: %w", constant.ErrorDefinition{
			Code:    http.StatusInternalServerError,
			Message: "internal server error",
		})
	}

	cpsRes, err := b.bankRepo.CreateBank(ctx, model.CreateCPSAction{
		MakerUser:  req.MakerUser,
		Department: req.Department,
		ActionData: entity.Bank{
			Name: actionData.Name,
			Logo: saveObj.Bucket + "/" + saveObj.Key,
			Code: actionData.Code,
			BIC:  actionData.BIC,
		},
	})

	if err != nil {
		return nil, err
	}

	return cpsRes, nil
}

func (b *BankDomain) DeleteOneBank(ctx context.Context, id string, req model.CreateCPSAction) (*model.CpsAction, error) {
	cpsAction, err := b.bankRepo.DeleteBank(ctx, id, req)
	if err != nil {
		return nil, err
	}

	return cpsAction, nil
}

func (b *BankDomain) GetAllBank(ctx context.Context, filterParams *constant.Filter) (*entity.BankResponse, error) {
	banks, err := b.bankRepo.GetAllBanks(ctx, filterParams)
	if err != nil {
		return nil, err
	}

	return banks, nil
}

func (b *BankDomain) GetOneBank(ctx context.Context, id string) (*entity.Bank, error) {
	bank, err := b.bankRepo.GetBank(ctx, id)
	if err != nil {
		return nil, err
	}

	return bank, nil
}

func (b *BankDomain) UpdateOneBank(ctx context.Context, id string, req model.CreateCPSAction) (*model.CpsAction, error) {
	cpsAction, err := b.bankRepo.UpdateBank(ctx, id, req)
	if err != nil {
		return nil, err
	}

	return cpsAction, nil
}

func (b *BankDomain) Authorize(ctx context.Context, req model.AuthorizeCPSAction) (*model.CpsAction, error) {
	cpsAction, err := b.bankRepo.Authorize(ctx, req)
	if err != nil {
		return nil, err
	}

	return cpsAction, nil
}

func (b *BankDomain) Reject(ctx context.Context, req model.RejectCPSAction) (*model.CpsAction, error) {
	if err := req.Validate(); err != nil {
		b.logger.Errorf("validation error", err)
		return nil, err
	}
	cpsAction, err := b.bankRepo.Reject(ctx, req)
	if err != nil {
		return nil, err
	}

	return cpsAction, nil
}

func (b *BankDomain) EnableOrDisableWallet(ctx context.Context, id string,
	requestAction model.RequestAction, cpsReq model.CreateCPSAction) (*model.CpsAction, error) {

	cpsAction, err := b.bankRepo.EnableOrDisableWallet(ctx, id, requestAction, cpsReq)
	if err != nil {
		return nil, err
	}

	return cpsAction, nil
}
