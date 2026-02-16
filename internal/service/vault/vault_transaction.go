package vault

import (
	localization "cbe-super-app-cps-action/internal/constants/localization"
	imodel "cbe-super-app-cps-action/internal/constants/model"
	"cbe-super-app-cps-action/internal/constants/types"
	local_util "cbe-super-app-cps-action/pkgs/utils"
	"context"
	"errors"

	"go.opentelemetry.io/otel/attribute"
	"go.opentelemetry.io/otel/trace"
)

func (s *vaultCategoryService) FindAllVaultTransactions(ctx context.Context, filterParams *types.Filter) (*types.PaginatedResponse[[]*imodel.VaultTransaction], error) {
	ctx, span := local_util.TraceLogger(ctx, "service", "FindAllVaultTransactions", "vaultService", "vaultService")
	defer span.End()

	if filterParams == nil {
		f := &types.Filter{}
		filterParams = f

	}
	entities, err := s.repo.FindAllTransactionsWithPagination(ctx, *filterParams)
	if err != nil {
		span.AddEvent("Service error", trace.WithAttributes(attribute.String("error", err.Error())))
		s.logger.Errorf("failed to fetch vault transactions | err=%v", err)
		return nil, errors.New(localization.ErrorUnexpectedError.Code)
	}
	return &types.PaginatedResponse[[]*imodel.VaultTransaction]{
		Data: entities.Data,
		Meta: entities.Meta,
	}, nil
}
