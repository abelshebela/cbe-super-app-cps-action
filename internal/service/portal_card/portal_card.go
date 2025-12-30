package portalcard

import (
	"cbe-super-app-cps-action/internal/constants/types"
	"cbe-super-app-cps-action/internal/service"
	"cbe-super-app-cps-action/internal/storage"
	"context"

	local_util "cbe-super-app-cps-action/pkgs/utils"

	"gitlab.com/bersufekadgetachew/cbe-super-app-shared/shared/entities/model"

	"gitlab.com/bersufekadgetachew/cbe-super-app-shared/shared/utils"
	"go.opentelemetry.io/otel/attribute"
	"go.opentelemetry.io/otel/trace"
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

func (s *portalCardService) GetAll(ctx context.Context, filterParam *types.Filter) (*types.PaginatedResponse[[]model.Card], error) {
	ctx, span := local_util.TraceLogger(ctx, "service", "GetAll", "PortalCard", "GetAll")
	defer span.End()

	cards, err := s.appService.FindAllWithPagination(ctx, *filterParam)
	if err != nil {
		s.logger.Errorf("[GetAll] failed to fetch portal cards: %v", err)
		span.AddEvent("Failed to fetch portal cards", trace.WithAttributes(
			attribute.String("error", err.Error()),
		))
		return nil, err
	}
	s.logger.Infof("[GetAll] retrieved %d portal cards", len(cards.Data))
	return cards, nil
}

func (s *portalCardService) ValidatePortalCard(ctx context.Context, names []string) (bool, error) {
	ctx, span := local_util.TraceLogger(ctx, "service", "ValidatePortalCard", "PortalCard", "ValidatePortalCard")
	defer span.End()

	card, err := s.appService.ValidatePortalCard(ctx, names)
	if err != nil {
		s.logger.Errorf("[ValidatePortalCard] failed to validate portal cards: %v", err)
		span.AddEvent("Failed to validate portal cards", trace.WithAttributes(
			attribute.String("error", err.Error()),
		))
		return false, err
	}
	return card, nil
}
