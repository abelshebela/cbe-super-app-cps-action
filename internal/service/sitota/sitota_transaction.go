package sitota

import (
	"github.com/abelshebela/cbe-super-app-cps-action/internal/constants/localization"
	"github.com/abelshebela/cbe-super-app-cps-action/internal/constants/types"
	"context"
	"errors"

	"gitlab.com/bersufekadgetachew/cbe-super-app-shared/shared/entities/model"

	"github.com/abelshebela/cbe-super-app-cps-action/internal/storage"

	local_util "github.com/abelshebela/cbe-super-app-cps-action/pkgs/utils"

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
	log := local_util.LoggerFromCtx(ctx, s.logger)

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
		log.Errorf("[SitotaSvc][GetAll] fetch err: %v", err)
		span.AddEvent("Failed to fetch sitotas", trace.WithAttributes(
			attribute.String("error", err.Error()),
		))
		return nil, err
	}

	log.Infof("[SitotaSvc][GetAll] retrieved %d", len(sitotas.Data))
	return &types.PaginatedResponse[[]*model.SitotaTransaction]{
		Data: sitotas.Data,
		Meta: sitotas.Meta,
	}, nil
}

func (s *SitotaTransactionService) GetSitotaByID(ctx context.Context, id string) (*model.SitotaTransaction, error) {
	log := local_util.LoggerFromCtx(ctx, s.logger)

	ctx, span := local_util.TraceLogger(ctx, "service", "GetSitotaByID", "Sitota", "GetSitotaByID")
	defer span.End()

	if id == "" {
		log.Errorf("[SitotaSvc][GetByID] invalid id")
		span.AddEvent("Invalid id provided", trace.WithAttributes(
			attribute.String("error", localization.ErrorInvalidID.Code),
		))
		return nil, errors.New(localization.ErrorInvalidID.Code)
	}

	sitota, err := s.repo.Get(ctx, id)
	if err != nil {
		log.Errorf("[SitotaSvc][GetByID] fetch err: %v", err)
		span.AddEvent("Failed to get sitota transaction", trace.WithAttributes(
			attribute.String("error", err.Error()),
			attribute.String("id", id),
		))
		return nil, err
	}

	log.Infof("[SitotaSvc][GetByID] retrieved id: %s", id)
	return sitota, nil
}
