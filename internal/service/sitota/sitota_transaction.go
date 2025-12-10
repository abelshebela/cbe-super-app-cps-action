package sitota

import (
	"cbe-super-app-cps-action/internal/constants/localization"
	"cbe-super-app-cps-action/internal/constants/model"
	"cbe-super-app-cps-action/internal/constants/types"
	"context"
	"errors"

	"cbe-super-app-cps-action/internal/storage"

	"gitlab.com/bersufekadgetachew/cbe-super-app-shared/shared/utils"
)

type SitotaTransactionService struct {
	repo storage.SitotaRepository
	// client transactionpb.TransactionServiceClient
	logger utils.Logger
}

func NewSitotaTransactionService(repo storage.SitotaRepository, logger utils.Logger) *SitotaTransactionService {
	return &SitotaTransactionService{repo: repo, logger: logger}
}

func (s *SitotaTransactionService) GetAllSitotas(ctx context.Context, filterParams *types.Filter) (*types.PaginatedResponse[[]*model.SitotaTransaction], error) {
	if filterParams == nil {
		f := &types.Filter{}
		filterParams = f

	}
	sitotas, err := s.repo.FindAllWithPagination(ctx, *filterParams)
	if err != nil {
		s.logger.Errorf("[GetAllSitotas] failed to fetch sitotas: %v", err)
		return nil, err
	}

	s.logger.Infof("[GetAllSitotas] retrieved %d sitota transactions", len(sitotas.Data))
	return &types.PaginatedResponse[[]*model.SitotaTransaction]{
		Data: sitotas.Data,
		Meta: sitotas.Meta,
	}, nil
}

func (s *SitotaTransactionService) GetSitotaByID(ctx context.Context, id string) (*model.SitotaTransaction, error) {
	if id == "" {
		s.logger.Errorf("[GetSitotaByID] invalid id provided")
		return nil, errors.New(localization.ErrorInvalidID.Code)
	}

	sitota, err := s.repo.Get(ctx, id)
	if err != nil {
		s.logger.Errorf("[GetSitotaByID] failed to get sitota transaction: %v", err)
		return nil, err
	}

	s.logger.Infof("[GetSitotaByID] sitota transaction retrieved successfully for id: %s", id)
	return sitota, nil
}

// func (s *SitotaTransactionService) GetAllSitotas(ctx context.Context) ([]*model.SitotaTransaction, error) {
// 	req := &transactionpb.ListTransactionRequest{
// 		Limit: 100,
// 		Filter: []*transactionpb.Filter{
// 			{
// 				Field:    "transaction_type",
// 				Operator: "=",
// 				Value:    "sitota",
// 			},
// 		},
// 	}
// 	resp, err := s.client.GetTransactionsList(ctx, req)
// 	if err != nil {
// 		s.logger.Errorf("Failed to get transactions list from gRPC service: %v", err)
// 		return nil, errors.New(localization.ErrorUnexpectedError.Code)
// 	}

// 	var sitotas []*model.SitotaTransaction
// 	for _, tx := range resp.GetTransactions() {
// 		sitota := core.MapTransactionToSitota(tx)
// 		sitotas = append(sitotas, sitota)
// 	}

// 	return sitotas, nil
// }

// func (s *SitotaTransactionService) GetSitotaByID(ctx context.Context, id string) (*model.SitotaTransaction, error) {
// 	if id == "" {
// 		return nil, errors.New(localization.ErrorInvalidID.Code)
// 	}

// 	req := &transactionpb.TransactionRequest{
// 		TransactionId: id,
// 	}

// 	resp, err := s.client.GetTransaction(ctx, req)
// 	if err != nil {
// 		s.logger.Errorf("Failed to get transaction from gRPC service: %v", err)
// 		return nil, errors.New(localization.ErrorUnexpectedError.Code)
// 	}

// 	sitota := core.MapTransactionToSitota(resp)

// 	return sitota, nil
// }
