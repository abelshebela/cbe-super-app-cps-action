package permission

import (
	"cbe-super-app-cps-action/internal/constants/model"
	"cbe-super-app-cps-action/internal/storage"
	"context"
	"fmt"
	"strings"

	"gitlab.com/bersufekadgetachew/cbe-super-app-shared/shared/dal"
	"gitlab.com/bersufekadgetachew/cbe-super-app-shared/shared/utils"

	// "go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/v2/bson"
	"go.mongodb.org/mongo-driver/v2/mongo"
)

type PermissionGroupStorage struct {
	dal    dal.MongoDal[model.Card, model.Card]
	client *mongo.Client
	logger utils.Logger
}

func NewPermissionGroupRepository(client *mongo.Client, dbName string, collection string, logger utils.Logger) storage.PermissionGroupRepository {
	return &PermissionGroupStorage{
		dal:    dal.NewMongoDal[model.Card, model.Card](client, dbName, collection),
		client: client,
		logger: logger,
	}
}

func (p *PermissionGroupStorage) ValidatePermissionGroupByID(ctx context.Context, ids []string) (bool, error) {
	if len(ids) == 0 {
		return false, fmt.Errorf("PERMISSION_GROUP_ARRAY_EMPTY")
	}

	var cleaned []string
	for _, id := range ids {
		if trimmed := strings.TrimSpace(id); trimmed != "" {
			cleaned = append(cleaned, trimmed)
		}
	}

	if len(cleaned) == 0 {
		return false, fmt.Errorf("NO_VALID_PERMISSION_GROUP_ID")
	}
	// Query DB for all given ids
	filter := bson.M{"_id": bson.M{"$in": cleaned}}
	projection := bson.M{"_id": 1}

	groups, err := p.dal.FindAll(ctx, filter, projection)
	if err != nil {
		return false, fmt.Errorf("DB_ERROR: %w", err)
	}

	if len(cleaned) != len(groups) {
		return false, fmt.Errorf("PERMISSION_GROUP_NOT_FOUND")
	}

	return true, nil
}
