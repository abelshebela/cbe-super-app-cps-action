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

type articleCategory struct {
	logger     shared_utils.Logger
	articleDal dal.MongoDal[model.NewsCategoryModel, model.NewsCategoryModel]
	client     *mongo.Client
}

func NewArticleCategoryRepository(logger shared_utils.Logger, client *mongo.Client, dbName, collectionName string) storage.ArticleCategoryRepository {
	return &articleCategory{
		logger:     logger,
		articleDal: dal.NewMongoDal[model.NewsCategoryModel, model.NewsCategoryModel](client, dbName, collectionName),
		client:     client,
	}
}

func (a *articleCategory) CreateArticleCategory(ctx context.Context, category *model.NewsCategoryModel) error {
	_, err := a.articleDal.InsertOne(ctx, *category)
	if err != nil {
		return errors.New(localization.ErrorUnexpectedError.Code)
	}
	return nil
}

func (a *articleCategory) UpdateArticleCategory(ctx context.Context, category *model.NewsCategoryModel, id string) error {
	filter := bson.M{"_id": id, "is_deleted": false}
	update := buildCategoryUpdate(*category)
	_, err := a.articleDal.UpdateOne(ctx, filter, update)
	if err != nil {
		a.logger.Errorf("Error updating article category in database:", err)
		if err == mongo.ErrNoDocuments {
			return errors.New(localization.ErrorFileNotFound.Code)
		}
		return errors.New(localization.ErrorUnexpectedError.Code)
	}
	return nil
}

func (a *articleCategory) DeleteArticleCategory(ctx context.Context, id string) error {
	err := a.articleDal.DeleteOne(ctx, bson.M{"_id": id, "is_deleted": false})
	if err != nil {
		a.logger.Errorf("Error deleting article category in database:", err)
		if err == mongo.ErrNoDocuments {
			return errors.New(localization.ErrorFileNotFound.Code)
		}
		return errors.New(localization.ErrorUnexpectedError.Code)
	}
	return nil
}

func (a *articleCategory) EnableOrDisableArticleCategory(ctx context.Context, id string, enable bool) error {
	filter := bson.M{"_id": id, "is_deleted": false}
	update := bson.M{"is_active": enable, "updated_at": time.Now()}
	_, err := a.articleDal.UpdateOne(ctx, filter, update)
	if err != nil {
		a.logger.Errorf("Error enabling or disabling article category in database:", err)
		if err == mongo.ErrNoDocuments {
			return errors.New(localization.ErrorFileNotFound.Code)
		}
		return errors.New(localization.ErrorUnexpectedError.Code)
	}
	return nil
}

func buildCategoryUpdate(updateFields model.NewsCategoryModel) bson.M {
	update := bson.M{}
	if updateFields.Name != "" {
		update["name"] = updateFields.Name
	}
	if updateFields.Description != "" {
		update["description"] = updateFields.Description
	}
	if updateFields.IsActive {
		update["is_active"] = updateFields.IsActive
	}

	if updateFields.Slug != "" {
		update["slug"] = updateFields.Slug
	}

	update["updated_at"] = time.Now()

	return update
}
