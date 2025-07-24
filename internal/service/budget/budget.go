package budget

import (
	"cbe-super-app-budget/internal/constants/dto"
	"cbe-super-app-budget/internal/handlers/middleware"
	"cbe-super-app-budget/internal/service"
	"cbe-super-app-budget/internal/storage"
	"cbe-super-app-budget/platform/logger"
	"context"
	"time"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.uber.org/zap"
)

type budgetService struct {
	repository storage.BudgetRepository
	logger     logger.Logger
}

func (b budgetService) Add(ctx context.Context, req dto.BudgetRequest) (*dto.Budget, error) {
	budget, err := dto.NewBudget(req)
	if err != nil {
		b.logger.Error(ctx, "failed to create budget entry", zap.Error(err))
		return nil, middleware.NewAppErrorWithCause(middleware.ErrorTypeValidation, "failed to create budger", err).
			WithService("budget").
			WithOperation("create")
	}

	if err := b.repository.Save(ctx, budget); err != nil {
		return nil, err
	}

	b.logger.Info(ctx,
		"budget created succesfully",
		zap.String("budget_id", budget.ID.Hex()),
	)
	return budget, nil
}

func (b budgetService) Get(ctx context.Context) ([]dto.Budget, error) {
	budgets, err := b.repository.Find(ctx)
	if err != nil {
		b.logger.Error(ctx,
			"failed to retrive budget",
			zap.Error(err),
		)

		return nil, err
	}
	return budgets, nil
}

func (b budgetService) GetOne(ctx context.Context, id string) (*dto.Budget, error) {
	oid, err := primitive.ObjectIDFromHex(id)
	if err != nil {
		return nil, middleware.NewInternalError("failed to parese budget id", err).
			WithService("cbe-super-app-budget").
			WithOperation("ObjectIDFromHex").
			WithDetail("budget_id", id)
	}

	filter := bson.M{
		"_id":        oid,
		"is_deleted": false,
	}

	budget, err := b.repository.FindById(ctx, filter)
	if err != nil {
		b.logger.Error(ctx,
			"failed to retrive budget",
			zap.String("budget_id", id),
			zap.Error(err),
		)
		return nil, err
	}

	return budget, nil
}

func (b budgetService) Modify(ctx context.Context, req dto.BudgetRequest) error {
	filter := bson.M{
		"_id":        req.ID,
		"is_deleted": false,
	}

	exists, err := b.repository.FindById(ctx, filter)
	if err != nil {
		b.logger.Error(ctx,
			"budget not found",
			zap.String("budget_id", req.ID.Hex()),
			zap.Error(err),
		)

		return err
	}

	exists.LastModifiedAt = time.Now()
	var hasChange bool
	if exists.Spending != req.Spending && req.Spending > 0.0 {
		hasChange = true
		exists.Spending = req.Spending
	}

	if exists.BudgetCategory != req.BudgetCategory && !req.BudgetCategory.IsZero() {
		hasChange = true
		exists.BudgetCategory = req.BudgetCategory
	}

	if exists.OverspendNotification != req.OverspendNotification {
		hasChange = true
		exists.OverspendNotification = req.OverspendNotification
	}
	if exists.LimitedBudgetExceededNotifiaction != req.LimitedBudgetExceededNotifiaction {
		hasChange = true
		exists.LimitedBudgetExceededNotifiaction = req.LimitedBudgetExceededNotifiaction
	}
	if exists.BudgetType != req.BudgetType && req.BudgetType != "" {
		hasChange = true
		exists.BudgetType = req.BudgetType
	}
	if exists.StartingDate != req.StartingDate && !req.StartingDate.IsZero() {
		hasChange = true
		exists.StartingDate = req.StartingDate
	}
	if exists.EndingDate != req.EndingDate && !req.EndingDate.IsZero() {
		hasChange = true
		exists.EndingDate = req.EndingDate

	}

	if hasChange {
		if err := b.repository.Update(ctx, exists); err != nil {
			b.logger.Error(ctx,
				"failed to update budget",
				zap.String("budget_id", req.ID.Hex()),
				zap.Error(err),
			)

			return err
		}
	}

	b.logger.Info(ctx,
		"budget update succesfullly",
		zap.String("budget_id", req.ID.Hex()),
	)

	return nil
}

func (b budgetService) Remove(ctx context.Context, id string) error {
	oid, err := primitive.ObjectIDFromHex(id)
	if err != nil {
		return middleware.NewInternalError("failed to parse budget id", err).
			WithService("cbe-super-app-budget").
			WithOperation("ObjectIDFromHex").
			WithDetail("budget_id", id)
	}
	filter := bson.M{"_id": oid, "is_deleted": false}
	if _, err := b.repository.FindById(ctx, filter); err != nil {
		b.logger.Error(ctx,
			"budget not found",
			zap.String("budget_id", id),
			zap.Error(err),
		)

		return err
	}

	if err := b.repository.Delete(ctx, id); err != nil {
		b.logger.Error(ctx,
			"failed to delete budget",
			zap.String("budget_id", id),
			zap.Error(err),
		)
		return err
	}

	b.logger.Info(ctx,
		"budget delete successfuly",
		zap.String("budget_id", id),
	)

	return nil
}

func NewBudgetService(repository storage.BudgetRepository, logger logger.Logger) service.BudgetService {
	return budgetService{
		repository: repository,
		logger:     logger,
	}
}
