package vault

import (
	localization "cbe-super-app-cps-action/internal/constants/localization"
	imodel "cbe-super-app-cps-action/internal/constants/model"
	"cbe-super-app-cps-action/internal/constants/types"
	local_util "cbe-super-app-cps-action/pkgs/utils"
	"context"
	"database/sql"
	"errors"

	"go.opentelemetry.io/otel/attribute"
	"go.opentelemetry.io/otel/trace"
)

func (s *vaultCategoryService) FindAllVaultTransactions(ctx context.Context, filterParams *types.Filter) (*types.PaginatedResponse[[]imodel.VaultTransaction], error) {
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
	return &types.PaginatedResponse[[]imodel.VaultTransaction]{
		Data: entities.Data,
		Meta: entities.Meta,
	}, nil
}

func (s *vaultCategoryService) FindVaultTransaction(ctx context.Context, id string) (*imodel.VaultTransaction, error) {
	ctx, span := local_util.TraceLogger(ctx, "service", "FindVaultTransaction", "vaultService", "FindVaultTransaction")
	defer span.End()

	entity, err := s.repo.FindVaultTransaction(ctx, id)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) || err.Error() == localization.ErrorResourceNotFound.Code {
			return nil, errors.New(localization.ErrorVaultCategoryNotFound.Code)
		}
		s.logger.Errorf("failed to fetch vault transaction by id | err=%v", err)
		return nil, errors.New(localization.ErrorUnexpectedError.Code)
	}
	return entity, nil
}
