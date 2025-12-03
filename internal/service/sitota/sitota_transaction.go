package sitota

import (
	"cbe-super-app-cps-action/internal/constants/localization"
	"cbe-super-app-cps-action/internal/constants/model"
	"cbe-super-app-cps-action/internal/service/sitota/core"
	"context"
	"errors"

	transactionpb "cbe-super-app-cps-action/grpc/sitota"

	"gitlab.com/bersufekadgetachew/cbe-super-app-shared/shared/utils"
)

type SitotaTransactionService struct {
	client transactionpb.TransactionServiceClient
	logger utils.Logger
}

func NewSitotaTransactionService(client transactionpb.TransactionServiceClient, logger utils.Logger) *SitotaTransactionService {
	return &SitotaTransactionService{client: client, logger: logger}
}

func (s *SitotaTransactionService) GetAllSitotas(ctx context.Context) ([]*model.SitotaTransaction, error) {
	req := &transactionpb.ListTransactionRequest{
		Limit: 100,
		Filter: []*transactionpb.Filter{
			{
				Field:    "transaction_type",
				Operator: "=",
				Value:    "sitota",
			},
		},
	}
	resp, err := s.client.GetTransactionsList(ctx, req)
	if err != nil {
		s.logger.Errorf("Failed to get transactions list from gRPC service: %v", err)
		return nil, errors.New(localization.ErrorUnexpectedError.Code)
	}

	var sitotas []*model.SitotaTransaction
	for _, tx := range resp.GetTransactions() {
		sitota := core.MapTransactionToSitota(tx)
		sitotas = append(sitotas, sitota)
	}

	return sitotas, nil
}

func (s *SitotaTransactionService) GetSitotaByID(ctx context.Context, id string) (*model.SitotaTransaction, error) {
	if id == "" {
		return nil, errors.New(localization.ErrorInvalidID.Code)
	}

	req := &transactionpb.TransactionRequest{
		TransactionId: id,
	}

	resp, err := s.client.GetTransaction(ctx, req)
	if err != nil {
		s.logger.Errorf("Failed to get transaction from gRPC service: %v", err)
		return nil, errors.New(localization.ErrorUnexpectedError.Code)
	}

	sitota := core.MapTransactionToSitota(resp)

	return sitota, nil
}
