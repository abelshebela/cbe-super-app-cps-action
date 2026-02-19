package media

import (
	"cbe-super-app-cps-action/internal/constants/localization"
	"cbe-super-app-cps-action/internal/handlers/middleware"
	"cbe-super-app-cps-action/internal/storage"
	"context"
	"errors"
	"time"

	"gitlab.com/bersufekadgetachew/cbe-super-app-shared/shared/entities/model"

	"gitlab.com/bersufekadgetachew/cbe-super-app-shared/shared/config"
	"gitlab.com/bersufekadgetachew/cbe-super-app-shared/shared/dal"
	shared_utils "gitlab.com/bersufekadgetachew/cbe-super-app-shared/shared/utils"
	"go.mongodb.org/mongo-driver/v2/bson"
	"go.mongodb.org/mongo-driver/v2/mongo"
)

var (
	invalidCategoryID = "Invalid article category ID"
)

type articleCategory struct {
	logger             shared_utils.Logger
	articleCategoryDal dal.MongoDal[model.NewsCategoryModel, model.NewsCategoryModel]
	client             *mongo.Client
}

func NewArticleCategoryRepository(logger shared_utils.Logger, client *mongo.Client, cfg *config.VaultConfig, dbName, collectionName string) storage.ArticleCategoryRepository {
	return &articleCategory{
		logger:             logger,
		articleCategoryDal: dal.NewMongoDal[model.NewsCategoryModel, model.NewsCategoryModel](client, cfg, dbName, collectionName),
		client:             client,
	}
}

func (a *articleCategory) CreateArticleCategory(ctx context.Context, category *model.NewsCategoryModel) error {
	_, err := a.articleCategoryDal.InsertOne(ctx, *category)
	if err != nil {
		return errors.New(localization.ErrorUnexpectedError.Code)
	}
	return nil
}

func (a *articleCategory) UpdateArticleCategory(ctx context.Context, category *model.NewsCategoryModel, id string) error {
	objId, err := bson.ObjectIDFromHex(id)

	if err != nil {
		a.logger.Errorf(invalidCategoryID, err)
		return middleware.NewBadRequestError(invalidCategoryID, err)
	}
	filter := bson.M{"_id": objId, "is_deleted": false}
	update := buildCategoryUpdate(*category)
	_, err = a.articleCategoryDal.UpdateOne(ctx, filter, update)
	if err != nil {
		return err
	}
	return nil
}

func (a *articleCategory) DeleteArticleCategory(ctx context.Context, id string) error {
	objId, err := bson.ObjectIDFromHex(id)

	if err != nil {
		a.logger.Errorf("Invalid article category ID:", err)
		return middleware.NewBadRequestError(invalidCategoryID, err)
	}
	err = a.articleCategoryDal.DeleteOne(ctx, bson.M{"_id": objId, "is_deleted": false})
	if err != nil {
		return err
	}
	return nil
}

func (a *articleCategory) EnableOrDisableArticleCategory(ctx context.Context, id string, enable bool) error {
	objId, err := bson.ObjectIDFromHex(id)

	if err != nil {
		a.logger.Errorf("Invalid article category ID:", err)
		return middleware.NewBadRequestError(invalidCategoryID, err)
	}
	filter := bson.M{"_id": objId, "is_deleted": false}
	update := bson.M{"is_active": enable, "updated_at": time.Now()}
	_, err = a.articleCategoryDal.UpdateOne(ctx, filter, update)
	if err != nil {
		return err
	}
	return nil
}

func buildCategoryUpdate(updateFields model.NewsCategoryModel) bson.M {
	update := bson.M{}
	if updateFields.Name != "" {
		update["name"] = updateFields.Name
	}
	if updateFields.IsActive {
		update["is_active"] = updateFields.IsActive
	}

	if updateFields.Color != "" {
		update["color"] = updateFields.Color
	}

	update["updated_at"] = time.Now()

	return update
}
