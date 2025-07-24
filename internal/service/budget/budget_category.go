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

type budgetCategoryService struct {
	repository storage.BudgetCategoryRepository
	logger     logger.Logger
}

func (b budgetCategoryService) Add(ctx context.Context, req dto.BudgetCategoryRequest) (*dto.BudgetCategory, error) {
	budget, err := dto.NewBudgetCategory(req)
	if err != nil {
		b.logger.Error(ctx, "failed to create budget_category entry", zap.Error(err))
		return nil, middleware.NewAppErrorWithCause(middleware.ErrorTypeValidation, "failed to create budger", err).
			WithService("budget_category").
			WithOperation("create")
	}

	if err := b.repository.Save(ctx, budget); err != nil {
		return nil, err
	}

	b.logger.Info(ctx,
		"budget_category created succesfully",
		zap.String("budget_id", budget.ID.Hex()),
	)
	return budget, nil
}

func (b budgetCategoryService) Get(ctx context.Context) ([]dto.BudgetCategory, error) {
	budgets, err := b.repository.Find(ctx)
	if err != nil {
		b.logger.Error(ctx,
			"failed to retrive budget_category",
			zap.Error(err),
		)

		return nil, err
	}
	return budgets, nil
}

func (b budgetCategoryService) GetOne(ctx context.Context, id string) (*dto.BudgetCategory, error) {
	oid, err := primitive.ObjectIDFromHex(id)
	if err != nil {
		return nil, middleware.NewInternalError("failed to parese budget_category id", err).
			WithService("cbe-super-app-budget_category").
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
			"failed to retrive budget_category",
			zap.String("budget_id", id),
			zap.Error(err),
		)
		return nil, err
	}

	return budget, nil
}

func (b budgetCategoryService) Modify(ctx context.Context, req dto.BudgetCategoryRequest) error {
	filter := bson.M{
		"_id":        req.ID,
		"is_deleted": false,
	}

	exists, err := b.repository.FindById(ctx, filter)
	if err != nil {
		b.logger.Error(ctx,
			"budget_category not found",
			zap.String("budget_id", req.ID.Hex()),
			zap.Error(err),
		)

		return err
	}

	var hasChange bool
	exists.LastModifiedAt = time.Now()
	if exists.Logo != req.Logo && len(req.Logo) > 0 {
		hasChange = true
		exists.Logo = req.Logo
	}
	if exists.Name != req.Name && len(req.Name) > 0 {
		hasChange = true
		exists.Name = req.Name
	}
	if exists.Color != req.Color && len(req.Color) > 0 {
		hasChange = true
		exists.Color = req.Color
	}

	if hasChange {
		if err := b.repository.Update(ctx, exists); err != nil {
			b.logger.Error(ctx,
				"failed to update budget_category",
				zap.String("budget_id", req.ID.Hex()),
				zap.Error(err),
			)

			return err
		}
	}

	b.logger.Info(ctx,
		"budget_category update succesfullly",
		zap.String("budget_id", req.ID.Hex()),
	)

	return nil
}

func (b budgetCategoryService) Remove(ctx context.Context, id string) error {
	oid, err := primitive.ObjectIDFromHex(id)
	if err != nil {
		return middleware.NewInternalError("failed to parse budget_category id", err).
			WithService("cbe-super-app-budget_category").
			WithOperation("ObjectIDFromHex").
			WithDetail("budget_id", id)
	}
	filter := bson.M{"_id": oid, "is_deleted": false}
	if _, err := b.repository.FindById(ctx, filter); err != nil {
		b.logger.Error(ctx,
			"budget_category not found",
			zap.String("budget_id", id),
			zap.Error(err),
		)

		return err
	}

	if err := b.repository.Delete(ctx, id); err != nil {
		b.logger.Error(ctx,
			"failed to delete budget_category",
			zap.String("budget_id", id),
			zap.Error(err),
		)
		return err
	}

	b.logger.Info(ctx,
		"budget_category delete successfuly",
		zap.String("budget_id", id),
	)

	return nil
}

func NewBudgetCategoryService(repository storage.BudgetCategoryRepository, logger logger.Logger) service.BudgetCategoryService {
	return budgetCategoryService{
		repository: repository,
		logger:     logger,
	}
}
