package persistence

import (
    "reflect"
    "context"
    "go.mongodb.org/mongo-driver/v2/bson"
    "go.mongodb.org/mongo-driver/v2/mongo"
    "gitlab.com/bersufekadgetachew/cbe-super-app-shared/shared/dal"
    "cbe-super-app-member-users/internal/port/outbound"
)

type MongoRepository[T any] struct {
    dal dal.MongoDal[T, T]
}

func NewMongoRepository[T any](client *mongo.Client, dbName, collectionName string) outbound.Repository[T] {
    return &MongoRepository[T]{
        dal: dal.NewMongoDal[T, T](client, dbName, collectionName),
    }
}

func (r *MongoRepository[T]) FindByID(ctx context.Context, id string) (*T, error) {
    oid, err := bson.ObjectIDFromHex(id)
    if err != nil {
        return nil, err
    }
    filter := bson.M{"_id": oid}
    return r.dal.FindOne(ctx, filter, nil)
}

func (r *MongoRepository[T]) Save(ctx context.Context, entity T) error {
    _, err := r.dal.InsertOne(ctx, entity)
    return err
}

func (r *MongoRepository[T]) Update(ctx context.Context, entity T) error {
    // Assuming entity has an ID field
    filter := bson.M{"_id": reflect.ValueOf(entity).FieldByName("ID").Interface()}
    update := bson.M{"$set": entity}
    _, err := r.dal.UpdateOne(ctx, filter, update)
    return err
}

func (r *MongoRepository[T]) Delete(ctx context.Context, id string) error {
    oid, err := bson.ObjectIDFromHex(id)
    if err != nil {
        return err
    }
    filter := bson.M{"_id": oid}
    return r.dal.DeleteOne(ctx, filter)
}