package vault

import (
	"context"
	"database/sql"
	"errors"

	"github.com/abelshebela/cbe-super-app-cps-action/internal/constants"
	"github.com/abelshebela/cbe-super-app-cps-action/internal/constants/lib"
	localization "github.com/abelshebela/cbe-super-app-cps-action/internal/constants/localization"
	imodel "github.com/abelshebela/cbe-super-app-cps-action/internal/constants/model"
	"github.com/abelshebela/cbe-super-app-cps-action/internal/constants/types"
	local_util "github.com/abelshebela/cbe-super-app-cps-action/pkgs/utils"

	"go.mongodb.org/mongo-driver/mongo"
	"go.opentelemetry.io/otel/attribute"
	"go.opentelemetry.io/otel/trace"
)

// func (s *vaultCategoryService) CreateWithdrawalRequest(ctx context.Context, req *vault_dto.CreateWithdrawalRequest) error {

// 	ctx, span := local_util.TraceLogger(ctx, "service", "CreateWithdrawalRequest", "vaultCategoryService", "vaultCategoryService")
// 	defer span.End()

// 	withdrawal := &imodel.Withdrawal{
// 		LockedVaultID:         req.LockedVaultID,
// 		Amount:                fmt.Sprintf("%f", req.WithdrawalAmount),
// 		WithdrawerName:        req.WithdrawerName,
// 		WithdrawerPhoneNumber: req.WithdrawerPhoneNumber,
// 		Status:                "PENDING",
// 		IsActive:              true,
// 	}

// 	makerData := local_util.ExtractUserFromContext(ctx)
// 	cpsActionModel := lib.CpsModelBuilder("", makerData, nil, withdrawal, string(constants.RequestCreateWithdrawal), string(constants.CREATE))
// 	if err := s.cpsService.CreateCPSAction(ctx, &cpsActionModel); err != nil {
// 		log.Errorf("[VaultDeadlockSvc][Create] cps action err: %v", err)
// 		return err
// 	}
// 	return nil
// }

func (s *vaultCategoryService) UnlockDeadlockRequest(ctx context.Context, id string) error {
	log := local_util.LoggerFromCtx(ctx, s.logger)

	ctx, span := local_util.TraceLogger(ctx, "service", "UnlockDeadlockRequest", "vaultCategoryService", "vaultCategoryService")
	defer span.End()

	_, err := s.repo.GetDeadlockRequest(ctx, id)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) || errors.Is(err, mongo.ErrNoDocuments) {
			return errors.New(localization.ErrorVaultCategoryNotFound.Code)
		}
		log.Errorf("[VaultDeadlockSvc][Update] fetch err: %v", err)
		return errors.New(localization.ErrorUnexpectedError.Code)
	}

	makerData := local_util.ExtractUserFromContext(ctx)
	payload := map[string]string{"id": id, "is_deadlock": "false"}
	cpsActionModel := lib.CpsModelBuilder(id, makerData, nil, payload, string(constants.RequestUnlockDeadlock), string(constants.UPDATE))
	if err := s.cpsService.CreateCPSAction(ctx, &cpsActionModel); err != nil {
		log.Errorf("[VaultDeadlockSvc][Update] cps action err: %v", err)
		return err
	}
	return nil
}

func (s *vaultCategoryService) GetAllDeadlockRequests(ctx context.Context, filterParams *types.Filter) (*types.PaginatedResponse[[]imodel.DeadlockRequest], error) {
	log := local_util.LoggerFromCtx(ctx, s.logger)

	ctx, span := local_util.TraceLogger(ctx, "service", "GetAllWithdrawalRequests", "vaultCategoryService", "vaultCategoryService")
	defer span.End()

	if filterParams == nil {
		filterParams = &types.Filter{}
	}
	entities, err := s.repo.GetAllDeadlockedRequests(ctx, *filterParams)
	if err != nil {
		span.AddEvent("Service error", trace.WithAttributes(attribute.String("error", err.Error())))
		log.Errorf("[VaultDeadlockSvc][GetAll] fetch err: %v", err)
		return nil, errors.New(localization.ErrorUnexpectedError.Code)
	}
	return entities, nil
}

func (s *vaultCategoryService) GetDeadlockRequestById(ctx context.Context, id string) (*imodel.DeadlockRequest, error) {
	log := local_util.LoggerFromCtx(ctx, s.logger)

	ctx, span := local_util.TraceLogger(ctx, "service", "GetWithdrawalRequest", "vaultCategoryService", "vaultCategoryService")
	defer span.End()

	entity, err := s.repo.GetDeadlockRequest(ctx, id)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) || errors.Is(err, mongo.ErrNoDocuments) || err.Error() == localization.ErrorVaultCategoryNotFound.Code {
			return nil, errors.New(localization.ErrorVaultCategoryNotFound.Code)
		}
		log.Errorf("[VaultDeadlockSvc][GetByID] fetch err: %v", err)
		return nil, errors.New(localization.ErrorUnexpectedError.Code)
	}
	return entity, nil
}

// func (s *vaultCategoryService) AuthorizeWithdrawalCreate(ctx context.Context, withdrawal *imodel.Vault) error {
// 	return s.repo.CreateWithdrawalRequest(ctx, withdrawal)
// }

func (s *vaultCategoryService) AuthorizeDeadlockStatusUpdate(ctx context.Context, id string, status string) error {
	log := local_util.LoggerFromCtx(ctx, s.logger)

	ctx, span := local_util.TraceLogger(ctx, "service", "AuthorizeDeadlockStatusUpdate", "vaultCategoryService", "vaultCategoryService")
	defer span.End()

	if err := s.repo.UpdateVaultDeadlock(ctx, id, status); err != nil {
		log.Errorf("[VaultDeadlockSvc][AuthorizeDeadlockStatusUpdate] update err: %v", err)
		return err
	}
	return nil
}
