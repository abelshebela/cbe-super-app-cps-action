package portalcard

import (
	"cbe-super-app-cps-action/internal/constants/model"
	"cbe-super-app-cps-action/internal/constants/types"
	"cbe-super-app-cps-action/internal/service"
	"cbe-super-app-cps-action/internal/storage"
	"context"

	"gitlab.com/bersufekadgetachew/cbe-super-app-shared/shared/utils"
)

type portalCardService struct {
	appService storage.PortalCardRepository
	logger     utils.Logger
}

func NewportalCardService(storage storage.PortalCardRepository, logger utils.Logger) service.PortalCardService {
	return &portalCardService{
		appService: storage,
		logger:     logger,
	}
}

func (s *portalCardService) GetAll(ctx context.Context, filterParam *types.Filter) (*types.PaginatedResponse[[]*model.Card], error) {

	cards, err := s.appService.FindAllWithPagination(ctx, *filterParam)
	if err != nil {
		return nil, err
	}
	return cards, nil
}

func (s *portalCardService) ValidatePortalCard(ctx context.Context, names []string) (bool, error) {
	card, err := s.appService.ValidatePortalCard(ctx, names)
	if err != nil {
		return false, err
	}
	return card, nil
}
