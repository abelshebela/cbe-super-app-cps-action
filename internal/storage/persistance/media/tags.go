package media

import (
	"cbe-super-app-cps-action/internal/constants/localization"
	"cbe-super-app-cps-action/internal/constants/model"
	"cbe-super-app-cps-action/internal/storage"
	"context"
	"errors"
	"time"

	"gitlab.com/bersufekadgetachew/cbe-super-app-shared/shared/dal"
	shared_utils "gitlab.com/bersufekadgetachew/cbe-super-app-shared/shared/utils"
	"go.mongodb.org/mongo-driver/v2/bson"
	"go.mongodb.org/mongo-driver/v2/mongo"
)

var (
	invalidTagID = "Invalid article tag ID"
)

type newsTags struct {
	logger      shared_utils.Logger
	newsTagsDal dal.MongoDal[model.NewsTags, model.NewsTags]
	client      *mongo.Client
}

func NewNewsTagsRepository(logger shared_utils.Logger, client *mongo.Client, dbName, collectionName string) storage.NewsTagsRepository {
	return &newsTags{
		logger:      logger,
		newsTagsDal: dal.NewMongoDal[model.NewsTags, model.NewsTags](client, dbName, collectionName),
		client:      client,
	}
}

func (n *newsTags) Create(ctx context.Context, newsTag *model.NewsTags) error {
	_, err := n.newsTagsDal.InsertOne(ctx, *newsTag)
	if err != nil {
		return errors.New(localization.ErrorUnexpectedError.Code)
	}
	return nil
}

func (n *newsTags) Update(ctx context.Context, newsTag *model.NewsTags, id string) error {
	objId, err := bson.ObjectIDFromHex(id)

	if err != nil {
		n.logger.Errorf(invalidTagID, err)
		return errors.New(invalidTagID)
	}
	filter := bson.M{"_id": objId, "is_deleted": false}
	update := buildNewsTagUpdate(*newsTag)
	_, err = n.newsTagsDal.UpdateOne(ctx, filter, update)
	if err != nil {
		n.logger.Errorf("Error updating news tag in database:", err)
		if err == mongo.ErrNoDocuments {
			return errors.New(localization.ErrorFileNotFound.Code)
		}
		return errors.New(localization.ErrorUnexpectedError.Code)
	}
	return nil
}

func (n *newsTags) Delete(ctx context.Context, id string) error {
	objId, err := bson.ObjectIDFromHex(id)

	if err != nil {
		n.logger.Errorf(invalidTagID, err)
		return errors.New(localization.ErrorInvalidRequest.Code)
	}
	filter := bson.M{"_id": objId, "is_deleted": false}

	err = n.newsTagsDal.DeleteOne(ctx, filter)
	if err != nil {
		n.logger.Errorf("Error deleting news tag in database:", err)
		if err == mongo.ErrNoDocuments {
			return errors.New(localization.ErrorFileNotFound.Code)
		}
		return errors.New(localization.ErrorUnexpectedError.Code)
	}
	return nil
}

func (n *newsTags) EnableDisable(ctx context.Context, id string, isEnable bool) error {
	objId, err := bson.ObjectIDFromHex(id)

	if err != nil {
		n.logger.Errorf(invalidTagID, err)
		return errors.New(localization.ErrorInvalidRequest.Code)
	}
	filter := bson.M{"_id": objId, "is_deleted": false}
	update := bson.M{
		"is_enabled": isEnable,
		"updated_at": time.Now(),
	}
	_, err = n.newsTagsDal.UpdateOne(ctx, filter, update)
	if err != nil {
		n.logger.Errorf("Error enabling/disabling news tag in database:", err)
		if err == mongo.ErrNoDocuments {
			return errors.New(localization.ErrorFileNotFound.Code)
		}
		return errors.New(localization.ErrorUnexpectedError.Code)
	}
	return nil
}

func buildNewsTagUpdate(newsTag model.NewsTags) bson.M {
	update := bson.M{
		"name":       newsTag.Name,
		"is_enabled": newsTag.IsEnabled,
		"updated_at": newsTag.UpdatedAt,
		"color":      newsTag.Color,
	}
	return update
}
