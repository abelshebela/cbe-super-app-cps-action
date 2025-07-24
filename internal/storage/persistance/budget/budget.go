package budget

import (
	"cbe-super-app-budget/internal/constants/dto"
	"cbe-super-app-budget/internal/handlers/middleware"
	"cbe-super-app-budget/internal/storage"
	"cbe-super-app-budget/platform/logger"
	"context"
	"fmt"
	"time"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"go.uber.org/zap"
)

type budgetRepository struct {
	collection *mongo.Collection
	logger     logger.Logger
}

func (s *budgetRepository) Delete(ctx context.Context, id string) error {
	oid, err := primitive.ObjectIDFromHex(id)
	if err != nil {
		return middleware.NewInternalError("failted to convert budget id", err).
			WithService("Budget").
			WithOperation("ObjectIDFromHex").
			WithDetail("budget_id", id)
	}

	filter := bson.M{"_id": oid}
	update := bson.M{
		"$set": bson.M{
			"is_deleted": true,
			"deleted_at": time.Now(),
		},
	}

	if _, err := s.collection.UpdateOne(ctx, filter, update); err != nil {
		s.logger.Error(ctx,
			"failed to delete budget",
			zap.String("budget_id", id),
			zap.Error(err),
		)

		return middleware.NewNotFoundError("fialed to find budget").
			WithService("budget").
			WithOperation("delete").
			WithDetail("budget_id", id)
	}

	return nil
}

func (s *budgetRepository) Find(ctx context.Context) ([]dto.Budget, error) {
	filter := bson.M{"is_deleted": false}
	cursor, err := s.collection.Find(ctx, filter)
	if err != nil {
		s.logger.Error(ctx,
			"failed to find budget",
			zap.Error(err),
		)
	}

	defer cursor.Close(ctx)

	var budgets []dto.Budget
	for cursor.Next(ctx) {
		var budget dto.Budget
		if err := cursor.Decode(&budget); err != nil {
			s.logger.Error(ctx,
				"failed to decode budget document",
				zap.Error(err),
			)

			return nil, middleware.NewInternalError("failed to decode budget document", err).
				WithService("cbe-super-app-budget").
				WithOperation("find")
		}
		fmt.Println(budget)

		budgets = append(budgets, budget)
	}

	fmt.Println("length: ", len(budgets))
	return budgets, nil
}

func (s *budgetRepository) FindById(ctx context.Context, filter bson.M) (*dto.Budget, error) {

	var budget dto.Budget
	if err := s.collection.FindOne(ctx, filter).Decode(&budget); err != nil {
		s.logger.Error(ctx,
			"failed to find supending by ID",
			zap.Error(err),
		)

		if err == mongo.ErrNoDocuments {
			return nil, middleware.NewNotFoundError("budget").
				WithService("cbe-super-app-budget").
				WithOperation("find_one")
		}

		return nil, middleware.NewDatabaseError("fialed to find budget", err).
			WithService("cbe-super-app-budget").
			WithOperation("find_one")
	}

	return &budget, nil

}

func (s *budgetRepository) Save(ctx context.Context, budget *dto.Budget) error {
	if budget.ID.IsZero() {
		budget.ID = primitive.NewObjectID()
	}

	if _, err := s.collection.InsertOne(ctx, budget); err != nil {
		s.logger.Error(ctx,
			"failed to save budget",
			zap.String("budget_id", budget.ID.Hex()),
			zap.Error(err),
		)

		return middleware.NewDatabaseError("failed to save budget", err).
			WithService("cbe-super-app-budget").
			WithOperation("insert").
			WithDetail("budget_id", budget.ID)
	}

	return nil
}

func (s *budgetRepository) Update(ctx context.Context, budget *dto.Budget) error {
	budget.LastModifiedAt = time.Now()
	filter := bson.M{
		"_id":        budget.ID,
		"is_deleted": false,
	}

	update := bson.M{
		"$set": budget,
	}

	if _, err := s.collection.UpdateOne(ctx, filter, update); err != nil {
		s.logger.Error(ctx,
			"failed to update budget",
			zap.String("budget_id", budget.ID.Hex()),
			zap.Error(err),
		)
		if err == mongo.ErrNoDocuments {
			return middleware.NewNotFoundError("budget").
				WithService("cbe-super-app-budget").
				WithOperation("update_one")
		}

		return middleware.NewDatabaseError("fialed to update budget", err).
			WithService("cbe-super-app-budget").
			WithOperation("update_one")
	}

	return nil
}

func NewBudgetRepository(db *mongo.Database, logger logger.Logger) storage.BudgetRepository {
	return &budgetRepository{
		collection: db.Collection("budget"),
		logger:     logger,
	}
}
