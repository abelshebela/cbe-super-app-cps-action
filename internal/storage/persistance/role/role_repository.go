package role_repo

import (
	"cbe-super-app-cps-action/internal/constants/model"
	"cbe-super-app-cps-action/internal/storage"
	"context"

	"gitlab.com/bersufekadgetachew/cbe-super-app-shared/shared/dal"
	"gitlab.com/bersufekadgetachew/cbe-super-app-shared/shared/utils"
	"go.mongodb.org/mongo-driver/v2/bson"
	"go.mongodb.org/mongo-driver/v2/mongo"
)

type RoleRepository struct {
	client     *mongo.Client
	mongoDal   dal.MongoDal[model.Role, model.Role]
	logger     utils.Logger
	collection *mongo.Collection
}

func NewRoleRepository(client *mongo.Client, database, collection string, logger utils.Logger) storage.RoleRepository {
	return &RoleRepository{
		client:     client,
		mongoDal:   dal.NewMongoDal[model.Role, model.Role](client, database, collection),
		logger:     logger,
		collection: client.Database(database).Collection(collection),
	}
}

func (r *RoleRepository) Exists(ctx context.Context, id string) (bool, error) {
	oid, err := bson.ObjectIDFromHex(id)
	if err != nil {
		return false, err
	}
	count, err := r.collection.CountDocuments(ctx, bson.M{"_id": oid})
	if err != nil {
		return false, err
	}
	return count > 0, nil
}

func (r *RoleRepository) ExistsMany(ctx context.Context, ids []string) (bool, error) {
	if len(ids) == 0 {
		return true, nil
	}
	oids := make([]bson.ObjectID, 0, len(ids))
	for _, id := range ids {
		oid, err := bson.ObjectIDFromHex(id)
		if err != nil {
			return false, err
		}
		oids = append(oids, oid)
	}

	count, err := r.collection.CountDocuments(ctx, bson.M{"_id": bson.M{"$in": oids}})
	if err != nil {
		return false, err
	}
	return count == int64(len(ids)), nil
}
