package spending

import (
	"cbe-super-app-budget/internal/constants/dto"
	"cbe-super-app-budget/internal/handler/middleware"
	"cbe-super-app-budget/internal/service"
	"cbe-super-app-budget/internal/storage"
	"cbe-super-app-budget/platform/logger"
	"context"
	"time"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.uber.org/zap"
)

type spendingService struct {
	repository storage.SpendingRepository
	logger     logger.Logger
}

func (b spendingService) Add(ctx context.Context, req dto.SpendingRequest) (*dto.Spending, error) {
	spending, err := dto.NewSpending(req)
	if err != nil {
		b.logger.Error(ctx, "failed to create spending entry", zap.Error(err))
		return nil, middleware.NewAppErrorWithCause(middleware.ErrorTypeValidation, "failed to create budger", err).
			WithService("spending").
			WithOperation("create")
	}

	if err := b.repository.Save(ctx, spending); err != nil {
		return nil, err
	}

	b.logger.Info(ctx,
		"spending created succesfully",
		zap.String("spending_id", spending.ID.Hex()),
	)
	return spending, nil
}

func (b spendingService) Get(ctx context.Context) ([]dto.Spending, error) {
	spendings, err := b.repository.Find(ctx)
	if err != nil {
		b.logger.Error(ctx,
			"failed to retrive spending",
			zap.Error(err),
		)

		return nil, err
	}
	return spendings, nil
}

func (b spendingService) GetOne(ctx context.Context, id string) (*dto.Spending, error) {
	oid, err := primitive.ObjectIDFromHex(id)
	if err != nil {
		return nil, middleware.NewInternalError("failed to parese spending id", err).
			WithService("cbe-super-app-spending").
			WithOperation("ObjectIDFromHex").
			WithDetail("spending_id", id)
	}

	filter := bson.M{
		"_id":        oid,
		"is_deleted": false,
	}

	spending, err := b.repository.FindById(ctx, filter)
	if err != nil {
		b.logger.Error(ctx,
			"failed to retrive spending",
			zap.String("spending_id", id),
			zap.Error(err),
		)
		return nil, err
	}

	return spending, nil
}

func (b spendingService) Modify(ctx context.Context, req dto.SpendingRequest) error {
	filter := bson.M{
		"_id":        req.ID,
		"is_deleted": false,
	}

	exists, err := b.repository.FindById(ctx, filter)
	if err != nil {
		b.logger.Error(ctx,
			"spending not found",
			zap.String("spending_id", req.ID.Hex()),
			zap.Error(err),
		)

		return err
	}

	exists.LastModifiedAt = time.Now()
	exists.BudgetID = req.BudgetID
	exists.SpendingAmount = req.SpendingAmount
	exists.TransactionDate = req.TransactionDate
	exists.TransactionReference = req.TransactionReference

	if err := b.repository.Update(ctx, exists); err != nil {
		b.logger.Error(ctx,
			"failed to update spending",
			zap.String("spending_id", req.ID.Hex()),
			zap.Error(err),
		)

		return err
	}

	b.logger.Info(ctx,
		"spending update succesfullly",
		zap.String("spending_id", req.ID.Hex()),
	)

	return nil
}

func (b spendingService) Remove(ctx context.Context, id string) error {
	oid, err := primitive.ObjectIDFromHex(id)
	if err != nil {
		return middleware.NewInternalError("failed to parse spending id", err).
			WithService("cbe-super-app-spending").
			WithOperation("ObjectIDFromHex").
			WithDetail("spending_id", id)
	}
	filter := bson.M{"_id": oid, "is_deleted": false}
	if _, err := b.repository.FindById(ctx, filter); err != nil {
		b.logger.Error(ctx,
			"spending not found",
			zap.String("spending_id", id),
			zap.Error(err),
		)

		return err
	}

	if err := b.repository.Delete(ctx, id); err != nil {
		b.logger.Error(ctx,
			"failed to delete spending",
			zap.String("spending_id", id),
			zap.Error(err),
		)
		return err
	}

	b.logger.Info(ctx,
		"spending delete successfuly",
		zap.String("spending_id", id),
	)

	return nil
}

func NewSpendingService(repository storage.SpendingRepository, logger logger.Logger) service.SpendingService {
	return spendingService{
		repository: repository,
		logger:     logger,
	}
}
