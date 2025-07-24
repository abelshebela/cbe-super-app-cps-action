package budget

import (
	"cbe-super-app-budget/internal/constants/dto"
	"cbe-super-app-budget/internal/handlers/middleware"
	"cbe-super-app-budget/internal/storage"
	"cbe-super-app-budget/platform/logger"
	"context"
	"time"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"go.uber.org/zap"
)

type budgetCategoryRepository struct {
	collection *mongo.Collection
	logger     logger.Logger
}

func (s *budgetCategoryRepository) Delete(ctx context.Context, id string) error {
	oid, err := primitive.ObjectIDFromHex(id)
	if err != nil {
		return middleware.NewInternalError("failted to convert budgetCategory id", err).
			WithService("BudgetCategory").
			WithOperation("ObjectIDFromHex").
			WithDetail("budget_category_id", id)
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
			"failed to delete budget_category",
			zap.String("budget_category_id", id),
			zap.Error(err),
		)

		return middleware.NewNotFoundError("fialed to find budgetCategory").
			WithService("budgetCategory").
			WithOperation("delete").
			WithDetail("budget_category_id", id)
	}

	return nil
}

func (s *budgetCategoryRepository) Find(ctx context.Context) ([]dto.BudgetCategory, error) {
	filter := bson.M{"isDelete": false}
	cursor, err := s.collection.Find(ctx, filter)
	if err != nil {
		s.logger.Error(ctx,
			"failed to find budget_category",
			zap.Error(err),
		)
	}

	defer cursor.Close(ctx)

	var budgetCategorys []dto.BudgetCategory
	for cursor.Next(ctx) {
		var budgetCategory dto.BudgetCategory
		if err := cursor.Decode(&budgetCategory); err != nil {
			s.logger.Error(ctx,
				"failed to decode budgetCategory document",
				zap.Error(err),
			)

			return nil, middleware.NewInternalError("failed to decode budgetCategory document", err).
				WithService("cbe-super-app-budget").
				WithOperation("find")
		}

		budgetCategorys = append(budgetCategorys, budgetCategory)
	}

	return budgetCategorys, nil
}

func (s *budgetCategoryRepository) FindById(ctx context.Context, filter bson.M) (*dto.BudgetCategory, error) {

	var budgetCategory dto.BudgetCategory
	if err := s.collection.FindOne(ctx, filter).Decode(&budgetCategory); err != nil {
		s.logger.Error(ctx,
			"failed to find supending by ID",
			zap.Error(err),
		)

		if err == mongo.ErrNoDocuments {
			return nil, middleware.NewNotFoundError("budgetCategory").
				WithService("cbe-super-app-budget").
				WithOperation("find_one")
		}

		return nil, middleware.NewDatabaseError("fialed to find budgetCategory", err).
			WithService("cbe-super-app-budget").
			WithOperation("find_one")
	}

	return &budgetCategory, nil

}

func (s *budgetCategoryRepository) Save(ctx context.Context, budgetCategory *dto.BudgetCategory) error {
	if budgetCategory.ID.IsZero() {
		budgetCategory.ID = primitive.NewObjectID()
	}

	if _, err := s.collection.InsertOne(ctx, budgetCategory); err != nil {
		s.logger.Error(ctx,
			"failed to save budgetCategory",
			zap.String("budget_category_id", budgetCategory.ID.Hex()),
			zap.Error(err),
		)

		return middleware.NewDatabaseError("failed to save budgetCategory", err).
			WithService("cbe-super-app-budget").
			WithOperation("insert").
			WithDetail("budget_category_id", budgetCategory.ID)
	}

	return nil
}

func (s *budgetCategoryRepository) Update(ctx context.Context, budgetCategory *dto.BudgetCategory) error {
	budgetCategory.LastModifiedAt = time.Now()
	filter := bson.M{
		"_id":        budgetCategory.ID,
		"is_deleted": false,
	}

	update := bson.M{
		"$set": budgetCategory,
	}
	if _, err := s.collection.UpdateOne(ctx, filter, update); err != nil {
		s.logger.Error(ctx,
			"failed to update budgetCategory",
			zap.String("budget_category_id", budgetCategory.ID.Hex()),
			zap.Error(err),
		)
		if err == mongo.ErrNoDocuments {
			return middleware.NewNotFoundError("budgetCategory").
				WithService("cbe-super-app-budget").
				WithOperation("update_one")
		}

		return middleware.NewDatabaseError("fialed to update budgetCategory", err).
			WithService("cbe-super-app-budget").
			WithOperation("update_one")
	}

	return nil
}

func NewBudgetCategoryRepository(db *mongo.Database, logger logger.Logger) storage.BudgetCategoryRepository {
	return &budgetCategoryRepository{
		collection: db.Collection("budgetCategory"),
		logger:     logger,
	}
}
