package sitota

import (
	"cbe-super-app-cps-action/internal/constants/localization"
	"cbe-super-app-cps-action/internal/constants/types"
	"context"
	"errors"

	"gitlab.com/bersufekadgetachew/cbe-super-app-shared/shared/entities/model"

	"cbe-super-app-cps-action/internal/storage"

	local_util "cbe-super-app-cps-action/pkgs/utils"

	"gitlab.com/bersufekadgetachew/cbe-super-app-shared/shared/utils"
	"go.opentelemetry.io/otel/attribute"
	"go.opentelemetry.io/otel/trace"
)

type SitotaTransactionService struct {
	repo   storage.SitotaRepository
	logger utils.Logger
}

func NewSitotaTransactionService(repo storage.SitotaRepository, logger utils.Logger) *SitotaTransactionService {
	return &SitotaTransactionService{repo: repo, logger: logger}
}

func (s *SitotaTransactionService) GetAllSitotas(ctx context.Context, filterParams *types.Filter) (*types.PaginatedResponse[[]*model.SitotaTransaction], error) {
	ctx, span := local_util.TraceLogger(ctx, "service", "GetAllSitotas", "Sitota", "GetAllSitotas")
	defer span.End()

	if filterParams.Search == "" && len(filterParams.Filters) == 0 {
		return &types.PaginatedResponse[[]*model.SitotaTransaction]{
			Data: []*model.SitotaTransaction{},
			Meta: types.PaginationMeta{},
		}, nil
	}

	sitotas, err := s.repo.FindAllWithPagination(ctx, *filterParams)
	if err != nil {
		s.logger.Errorf("[GetAllSitotas] failed to fetch sitotas: %v", err)
		span.AddEvent("Failed to fetch sitotas", trace.WithAttributes(
			attribute.String("error", err.Error()),
		))
		return nil, err
	}

	s.logger.Infof("[GetAllSitotas] retrieved %d sitota transactions", len(sitotas.Data))
	return &types.PaginatedResponse[[]*model.SitotaTransaction]{
		Data: sitotas.Data,
		Meta: sitotas.Meta,
	}, nil
}

func (s *SitotaTransactionService) GetSitotaByID(ctx context.Context, id string) (*model.SitotaTransaction, error) {
	ctx, span := local_util.TraceLogger(ctx, "service", "GetSitotaByID", "Sitota", "GetSitotaByID")
	defer span.End()

	if id == "" {
		s.logger.Errorf("[GetSitotaByID] invalid id provided")
		span.AddEvent("Invalid id provided", trace.WithAttributes(
			attribute.String("error", localization.ErrorInvalidID.Code),
		))
		return nil, errors.New(localization.ErrorInvalidID.Code)
	}

	sitota, err := s.repo.Get(ctx, id)
	if err != nil {
		s.logger.Errorf("[GetSitotaByID] failed to get sitota transaction: %v", err)
		span.AddEvent("Failed to get sitota transaction", trace.WithAttributes(
			attribute.String("error", err.Error()),
			attribute.String("id", id),
		))
		return nil, err
	}

	s.logger.Infof("[GetSitotaByID] sitota transaction retrieved successfully for id: %s", id)
	return sitota, nil
}
