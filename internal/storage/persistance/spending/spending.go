package spending

import (
	"cbe-super-app-budget/internal/constants/dto"
	"cbe-super-app-budget/internal/handler/middleware"
	"cbe-super-app-budget/internal/storage"
	"cbe-super-app-budget/platform/logger"
	"context"
	"time"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"go.uber.org/zap"
)

type spendingRepository struct {
	collection *mongo.Collection
	logger     logger.Logger
}

func (s *spendingRepository) Delete(ctx context.Context, id string) error {
	oid, err := primitive.ObjectIDFromHex(id)
	if err != nil {
		return middleware.NewInternalError("failted to convert spending id", err).
			WithService("Spending").
			WithOperation("ObjectIDFromHex").
			WithDetail("spending_id", id)
	}

	filter := bson.M{"_id": oid}
	update := bson.M{
		"$set": bson.M{
			"is_deleted": true,
		},
	}

	if _, err := s.collection.UpdateOne(ctx, filter, update); err != nil {
		s.logger.Error(ctx,
			"failed to delete spending",
			zap.String("spending_id", id),
			zap.Error(err),
		)

		return middleware.NewNotFoundError("fialed to find spending").
			WithService("spending").
			WithOperation("delete").
			WithDetail("spending_id", id)
	}

	return nil
}

func (s *spendingRepository) Find(ctx context.Context) ([]dto.Spending, error) {
	filter := bson.M{"isDelete": false}
	cursor, err := s.collection.Find(ctx, filter)
	if err != nil {
		s.logger.Error(ctx,
			"failed to find spending",
			zap.Error(err),
		)
	}

	defer cursor.Close(ctx)

	var spendings []dto.Spending
	for cursor.Next(ctx) {
		var spending dto.Spending
		if err := cursor.Decode(&spending); err != nil {
			s.logger.Error(ctx,
				"failed to decode spending document",
				zap.Error(err),
			)

			return nil, middleware.NewInternalError("failed to decode spending document", err).
				WithService("cbe-super-app-budget").
				WithOperation("find")
		}

		spendings = append(spendings, spending)
	}

	return spendings, nil
}

func (s *spendingRepository) FindById(ctx context.Context, filter bson.M) (*dto.Spending, error) {

	var spending dto.Spending
	if err := s.collection.FindOne(ctx, filter).Decode(&spending); err != nil {
		s.logger.Error(ctx,
			"failed to find supending by ID",
			zap.Error(err),
		)

		if err == mongo.ErrNoDocuments {
			return nil, middleware.NewNotFoundError("spending").
				WithService("cbe-super-app-budget").
				WithOperation("find_one")
		}

		return nil, middleware.NewDatabaseError("fialed to find spending", err).
			WithService("cbe-super-app-budget").
			WithOperation("find_one")
	}

	return &spending, nil

}

func (s *spendingRepository) Save(ctx context.Context, spending *dto.Spending) error {
	if spending.ID.IsZero() {
		spending.ID = primitive.NewObjectID()
	}

	if _, err := s.collection.InsertOne(ctx, spending); err != nil {
		s.logger.Error(ctx,
			"failed to save spending",
			zap.String("spending_id", spending.ID.Hex()),
			zap.Error(err),
		)

		return middleware.NewDatabaseError("failed to save spending", err).
			WithService("cbe-super-app-budget").
			WithOperation("insert").
			WithDetail("spending_id", spending.ID)
	}

	return nil
}

func (s *spendingRepository) Update(ctx context.Context, spending *dto.Spending) error {
	spending.LastModifiedAt = time.Now()
	filter := bson.M{
		"_id":        spending.ID,
		"is_deleted": false,
	}

	update := bson.M{
		"$set": spending,
	}
	if _, err := s.collection.UpdateOne(ctx, filter, update); err != nil {
		s.logger.Error(ctx,
			"failed to update spending",
			zap.String("spending_id", spending.ID.Hex()),
			zap.Error(err),
		)
		if err == mongo.ErrNoDocuments {
			return middleware.NewNotFoundError("spending").
				WithService("cbe-super-app-budget").
				WithOperation("update_one")
		}

		return middleware.NewDatabaseError("fialed to update spending", err).
			WithService("cbe-super-app-budget").
			WithOperation("update_one")
	}

	return nil
}

func NewSpendingRepository(db *mongo.Database, logger logger.Logger) storage.SpendingRepository {
	return &spendingRepository{
		collection: db.Collection("spending"),
		logger:     logger,
	}
}
