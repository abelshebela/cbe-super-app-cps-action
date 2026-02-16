package vault

import (
	"context"
	"database/sql"
	"errors"
	"fmt"

	"cbe-super-app-cps-action/internal/constants"
	vault_dto "cbe-super-app-cps-action/internal/constants/dto/vault"
	"cbe-super-app-cps-action/internal/constants/lib"
	localization "cbe-super-app-cps-action/internal/constants/localization"
	imodel "cbe-super-app-cps-action/internal/constants/model"
	"cbe-super-app-cps-action/internal/constants/types"
	local_util "cbe-super-app-cps-action/pkgs/utils"

	"go.mongodb.org/mongo-driver/mongo"
	"go.opentelemetry.io/otel/attribute"
	"go.opentelemetry.io/otel/trace"
)

func (s *vaultCategoryService) CreateWithdrawalRequest(ctx context.Context, req *vault_dto.CreateWithdrawalRequest) error {
	ctx, span := local_util.TraceLogger(ctx, "service", "CreateWithdrawalRequest", "vaultCategoryService", "vaultCategoryService")
	defer span.End()

	withdrawal := &imodel.Withdrawal{
		LockedVaultID:         req.LockedVaultID,
		Amount:                fmt.Sprintf("%f", req.WithdrawalAmount),
		WithdrawerName:        req.WithdrawerName,
		WithdrawerPhoneNumber: req.WithdrawerPhoneNumber,
		Status:                "PENDING",
		IsActive:              true,
	}

	makerData := local_util.ExtractUserFromContext(ctx)
	cpsActionModel := lib.CpsModelBuilder("", makerData, nil, withdrawal, string(constants.RequestCreateWithdrawal), string(constants.CREATE))
	if err := s.cpsService.CreateCPSAction(ctx, &cpsActionModel); err != nil {
		s.logger.Errorf("failed to create CPS action for vault withdrawal | action=%s | err=%v", constants.RequestCreateWithdrawal, err)
		return err
	}
	return nil
}

func (s *vaultCategoryService) UpdateWithdrawalRequest(ctx context.Context, id string, req *vault_dto.UpdateWithdrawalStatusRequest) error {
	ctx, span := local_util.TraceLogger(ctx, "service", "UpdateWithdrawalRequest", "vaultCategoryService", "vaultCategoryService")
	defer span.End()

	_, err := s.repo.GetWithdrawalRequest(ctx, id)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) || errors.Is(err, mongo.ErrNoDocuments) {
			return errors.New(localization.ErrorVaultCategoryNotFound.Code)
		}
		s.logger.Errorf("failed to fetch withdrawal request by id | err=%v", err)
		return errors.New(localization.ErrorUnexpectedError.Code)
	}

	makerData := local_util.ExtractUserFromContext(ctx)
	payload := map[string]string{"withdrawal_id": id, "withdrawal_status": req.WithdrawalStatus}
	cpsActionModel := lib.CpsModelBuilder(id, makerData, nil, payload, string(constants.RequestUpdateWithdrawal), string(constants.UPDATE))
	if err := s.cpsService.CreateCPSAction(ctx, &cpsActionModel); err != nil {
		s.logger.Errorf("failed to create CPS action for vault withdrawal update | action=%s | err=%v", constants.RequestUpdateWithdrawal, err)
		return err
	}
	return nil
}

func (s *vaultCategoryService) GetAllWithdrawalRequests(ctx context.Context, filterParams *types.Filter) (*types.PaginatedResponse[[]imodel.Withdrawal], error) {
	ctx, span := local_util.TraceLogger(ctx, "service", "GetAllWithdrawalRequests", "vaultCategoryService", "vaultCategoryService")
	defer span.End()

	if filterParams == nil {
		filterParams = &types.Filter{}
	}
	entities, err := s.repo.GetAllWithdrawalRequests(ctx, *filterParams)
	if err != nil {
		span.AddEvent("Service error", trace.WithAttributes(attribute.String("error", err.Error())))
		s.logger.Errorf("failed to fetch withdrawal requests | err=%v", err)
		return nil, errors.New(localization.ErrorUnexpectedError.Code)
	}
	return entities, nil
}

func (s *vaultCategoryService) GetWithdrawalRequest(ctx context.Context, id string) (*imodel.Withdrawal, error) {
	ctx, span := local_util.TraceLogger(ctx, "service", "GetWithdrawalRequest", "vaultCategoryService", "vaultCategoryService")
	defer span.End()

	entity, err := s.repo.GetWithdrawalRequest(ctx, id)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) || errors.Is(err, mongo.ErrNoDocuments) || err.Error() == localization.ErrorVaultCategoryNotFound.Code {
			return nil, errors.New(localization.ErrorVaultCategoryNotFound.Code)
		}
		s.logger.Errorf("failed to fetch withdrawal request by id | err=%v", err)
		return nil, errors.New(localization.ErrorUnexpectedError.Code)
	}
	return entity, nil
}

func (s *vaultCategoryService) AuthorizeWithdrawalCreate(ctx context.Context, withdrawal *imodel.Withdrawal) error {
	return s.repo.CreateWithdrawalRequest(ctx, withdrawal)
}

func (s *vaultCategoryService) AuthorizeWithdrawalUpdate(ctx context.Context, id string, status string) error {
	return s.repo.UpdateWithdrawalRequest(ctx, id, status)
}
