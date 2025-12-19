package mini_app

import (
	"cbe-super-app-cps-action/internal/constants/localization"
	// "cbe-super-app-cps-action/internal/constants/model"
	"cbe-super-app-cps-action/internal/handlers/middleware"
	"cbe-super-app-cps-action/internal/storage"
	"context"
	"errors"
	"time"

	"gitlab.com/bersufekadgetachew/cbe-super-app-shared/shared/dal"
	shared_utils "gitlab.com/bersufekadgetachew/cbe-super-app-shared/shared/utils"

	"gitlab.com/bersufekadgetachew/cbe-super-app-shared/shared/entities/model"

	"go.mongodb.org/mongo-driver/v2/bson"
	"go.mongodb.org/mongo-driver/v2/mongo"
)

var (
	invalidCategoryID = "Invalid mini app category ID"
)

type miniAppCategory struct {
	logger             shared_utils.Logger
	miniAppCategoryDal dal.MongoDal[model.MiniAppCategory, model.MiniAppCategory]
	client             *mongo.Client
}

func NewMiniAppCategoryRepository(logger shared_utils.Logger, client *mongo.Client, dbName, collectionName string) storage.MiniAppCategoryRepository {
	return &miniAppCategory{
		logger:             logger,
		miniAppCategoryDal: dal.NewMongoDal[model.MiniAppCategory, model.MiniAppCategory](client, dbName, collectionName),
		client:             client,
	}
}

func (a *miniAppCategory) Create(ctx context.Context, category *model.MiniAppCategory) error {
	_, err := a.miniAppCategoryDal.InsertOne(ctx, *category)
	if err != nil {
		return errors.New(localization.ErrorUnexpectedError.Code)
	}
	return nil
}

func (a *miniAppCategory) Update(ctx context.Context, category *model.MiniAppCategory, id string) error {
	objId, err := bson.ObjectIDFromHex(id)

	if err != nil {
		a.logger.Errorf(invalidCategoryID, err)
		return middleware.NewBadRequestError(invalidCategoryID, err)
	}
	filter := bson.M{"_id": objId, "is_deleted": false}
	update := buildCategoryUpdate(*category)
	_, err = a.miniAppCategoryDal.UpdateOne(ctx, filter, update)
	if err != nil {
		a.logger.Errorf("Error updating mini app category in database:", err)
		if err == mongo.ErrNoDocuments {
			return errors.New(localization.ErrorFileNotFound.Code)
		}
		return errors.New(localization.ErrorUnexpectedError.Code)
	}
	return nil
}

func (a *miniAppCategory) Delete(ctx context.Context, id string) error {
	objId, err := bson.ObjectIDFromHex(id)

	if err != nil {
		a.logger.Errorf("Invalid mini app category ID:", err)
		return middleware.NewBadRequestError(invalidCategoryID, err)
	}
	err = a.miniAppCategoryDal.DeleteOne(ctx, bson.M{"_id": objId, "is_deleted": false})
	if err != nil {
		a.logger.Errorf("Error deleting mini app category in database:", err)
		if err == mongo.ErrNoDocuments {
			return errors.New(localization.ErrorFileNotFound.Code)
		}
		return errors.New(localization.ErrorUnexpectedError.Code)
	}
	return nil
}

func (a *miniAppCategory) EnableOrDisable(ctx context.Context, id string, enable bool) error {
	objId, err := bson.ObjectIDFromHex(id)

	if err != nil {
		a.logger.Errorf("Invalid mini app category ID:", err)
		return middleware.NewBadRequestError(invalidCategoryID, err)
	}
	filter := bson.M{"_id": objId, "is_deleted": false}
	update := bson.M{"is_enabled": enable, "updated_at": time.Now()}
	_, err = a.miniAppCategoryDal.UpdateOne(ctx, filter, update)
	if err != nil {
		a.logger.Errorf("Error enabling or disabling mini app category in database:", err)
		if err == mongo.ErrNoDocuments {
			return errors.New(localization.ErrorFileNotFound.Code)
		}
		return errors.New(localization.ErrorUnexpectedError.Code)
	}
	return nil
}

func buildCategoryUpdate(updateFields model.MiniAppCategory) bson.M {
	update := bson.M{}
	if updateFields.Name != "" {
		update["name"] = updateFields.Name
	}
	if updateFields.IsEnabled {
		update["is_enabled"] = updateFields.IsEnabled
	}

	if updateFields.Icon != "" {
		update["icon"] = updateFields.Icon
	}

	update["updated_at"] = time.Now()

	return update
}
