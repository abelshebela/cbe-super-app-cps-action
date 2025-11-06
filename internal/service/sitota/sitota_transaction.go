package sitota

import (
	"cbe-super-app-cps-action/internal/constants/model"
	"context"

	"gitlab.com/bersufekadgetachew/cbe-super-app-shared/shared/utils"
)

type SitotaTransactionService struct {
	// gRPCClient
	logger utils.Logger
}

func NewSitotaTransactionService(logger utils.Logger) *SitotaTransactionService {
	return &SitotaTransactionService{logger: logger}
}

func (s *SitotaTransactionService) GetAllSitotas(ctx context.Context) ([]*model.SitotaTransaction, error) {
	return nil, nil
}

func (s *SitotaTransactionService) GetSitotaByID(ctx context.Context, id string) (*model.SitotaTransaction, error) {
	return nil, nil
}
