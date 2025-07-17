package infra_mongo

import (
	"context"

	"go.mongodb.org/mongo-driver/v2/bson"
	"go.mongodb.org/mongo-driver/v2/mongo"
	"go.mongodb.org/mongo-driver/v2/mongo/options"
)

type MongoDalInterface[T any, K any] interface {
	FindAll(ctx context.Context, filter, projection bson.M) ([]*K, error)
	FindAllWithPagination(ctx context.Context, filter, projection bson.M, skip, limit int64) ([]*K, error)
	TotalCount(ctx context.Context, filter bson.M) (int64, error)
	FindOne(ctx context.Context, filter, projection bson.M) (*K, error)
	InsertOne(ctx context.Context, req T) (T, error)
	UpdateOne(ctx context.Context, filter, update bson.M) (T, error)
	DeleteOne(ctx context.Context, filter bson.M) error
	FindRecentDocument(ctx context.Context, filter, projection bson.M) (*K, error)
}

var _ MongoDalInterface[any, any] = (*MongoDal[any, any])(nil)

type MongoDal[T any, K any] struct {
	collection *mongo.Collection
}

func NewMongoDal[T any, K any](client *mongo.Client, dbName, collectionName string) *MongoDal[T, K] {
	collection := client.Database(dbName).Collection(collectionName)
	return &MongoDal[T, K]{
		collection: collection,
	}
}

func (m *MongoDal[T, K]) FindAll(ctx context.Context, filter, projection bson.M) ([]*K, error) {
	opts := options.Find().SetProjection(projection)
	cursor, err := m.collection.Find(ctx, filter, opts)
	if err != nil {
		return nil, err
	}
	var results []*K
	if err := cursor.All(ctx, &results); err != nil {
		return nil, err
	}
	return results, nil
}

func (m *MongoDal[T, K]) FindAllWithPagination(ctx context.Context, filter, projection bson.M, skip, limit int64) ([]*K, error) {
	opts := options.Find().SetSkip(skip).SetLimit(limit).SetProjection(projection)
	cursor, err := m.collection.Find(ctx, filter, opts)
	if err != nil {
		return nil, err
	}
	var results []*K
	if err := cursor.All(ctx, &results); err != nil {
		return nil, err
	}
	return results, nil
}

func (m *MongoDal[T, K]) TotalCount(ctx context.Context, filter bson.M) (int64, error) {
	return m.collection.CountDocuments(ctx, filter)
}

func (m *MongoDal[T, K]) FindOne(ctx context.Context, filter, projection bson.M) (*K, error) {
	opts := options.FindOne().SetProjection(projection)
	var result K
	if err := m.collection.FindOne(ctx, filter, opts).Decode(&result); err != nil {
		return nil, err
	}
	return &result, nil
}

func (m *MongoDal[T, K]) InsertOne(ctx context.Context, req T) (T, error) {
	if _, err := m.collection.InsertOne(ctx, req); err != nil {
		var zero T
		return zero, err
	}

	return req, nil
}

func (m *MongoDal[T, K]) UpdateOne(ctx context.Context, filter, update bson.M) (T, error) {
	var result T
	err := m.collection.FindOneAndUpdate(
		ctx,
		filter,
		bson.M{"$set": update},
		options.FindOneAndUpdate().SetReturnDocument(options.After),
	).Decode(&result)
	if err != nil {
		var zero T
		return zero, err
	}
	return result, nil
}

func (m *MongoDal[T, K]) DeleteOne(ctx context.Context, filter bson.M) error {
	_, err := m.collection.UpdateOne(
		ctx,
		filter,
		bson.M{
			"$set": bson.M{
				"is_deleted": true,
			},
		},
	)
	return err
}

func (m *MongoDal[T, K]) FindRecentDocument(ctx context.Context, filter, projection bson.M) (*K, error) {
	opts := options.FindOne().SetSort(bson.M{"created_at": -1}).SetProjection(projection)
	var result K
	if err := m.collection.FindOne(ctx, filter, opts).Decode(&result); err != nil {
		return nil, err
	}
	return &result, nil
}
