package media

import (
	"cbe-super-app-cps-action/internal/constants/model"
	"cbe-super-app-cps-action/internal/handlers/middleware"
	"cbe-super-app-cps-action/internal/storage"
	"context"
	"time"

	"gitlab.com/bersufekadgetachew/cbe-super-app-shared/shared/dal"
	shared_utils "gitlab.com/bersufekadgetachew/cbe-super-app-shared/shared/utils"
	"go.mongodb.org/mongo-driver/v2/bson"
	"go.mongodb.org/mongo-driver/v2/mongo"
)

type article struct {
	logger     shared_utils.Logger
	articleDal dal.MongoDal[model.NewsArticle, model.NewsArticle]
	client     *mongo.Client
}

func NewsArticleRepository(logger shared_utils.Logger, client *mongo.Client, dbName, collectionName string) storage.ArticleRepository {
	return &article{
		logger:     logger,
		articleDal: dal.NewMongoDal[model.NewsArticle, model.NewsArticle](client, dbName, collectionName),
		client:     client,
	}
}

func (a *article) CreateArticle(ctx context.Context, article *model.NewsArticle) error {
	_, err := a.articleDal.InsertOne(ctx, *article)
	if err != nil {
		a.logger.Errorf("Error inserting article into database:", err)
		return middleware.NewDatabaseError("Error inserting article into database", err)
	}

	return nil
}

func (a *article) UpdateArticle(ctx context.Context, article *model.NewsArticle, id string) error {

	objId, err := a.getArticleID(id)
	if err != nil {
		return err
	}
	filter := bson.M{"_id": objId, "is_deleted": false}
	update := buildUpdate(*article)

	_, err = a.articleDal.UpdateOne(ctx, filter, update)
	if err != nil {
		a.logger.Errorf("Error updating article in database:", err)
		if err == mongo.ErrNoDocuments {
			return middleware.NewNotFoundError("article")
		}
		return middleware.NewDatabaseError("Error updating article in database", err)
	}
	return nil
}

func (a *article) DeleteArticle(ctx context.Context, id string) error {
	objId, err := a.getArticleID(id)
	if err != nil {
		return err
	}

	filter := bson.M{"_id": objId, "is_deleted": false}
	if err := a.articleDal.DeleteOne(ctx, filter); err != nil {
		a.logger.Errorf("Error deleting article from database:", err)
		if err == mongo.ErrNoDocuments {
			return middleware.NewNotFoundError("article")
		}
		return middleware.NewDatabaseError("Error deleting article from database", err)
	}
	return nil
}

func (a *article) PublishUnpublishArticle(ctx context.Context, id string, isPublished bool) error {
	objId, err := a.getArticleID(id)
	if err != nil {
		return err
	}

	filter := bson.M{"_id": objId, "is_deleted": false}
	update := bson.M{
		"is_published": isPublished,
		"updated_at":   time.Now(),
	}

	if isPublished {
		update["published_at"] = time.Now()
	}

	_, err = a.articleDal.UpdateOne(ctx, filter, update)
	if err != nil {
		a.logger.Errorf("Error updating article publish status in database:", err)
		if err == mongo.ErrNoDocuments {
			return middleware.NewNotFoundError("article")
		}
		return middleware.NewDatabaseError("Error updating article publish status in database", err)
	}
	return nil
}

func buildUpdate(updateFields model.NewsArticle) bson.M {
	update := bson.M{}

	if updateFields.Title != "" {
		update["title"] = updateFields.Title
	}
	if updateFields.Author != "" {
		update["author"] = updateFields.Author
	}
	if updateFields.CategoryID != bson.NilObjectID {
		update["category_id"] = updateFields.CategoryID
	}
	if updateFields.Thumbnail != "" {
		update["thumbnail"] = updateFields.Thumbnail
	}
	if len(updateFields.Tags) > 0 {
		update["tags"] = updateFields.Tags
	}
	if updateFields.Content != "" {
		update["content"] = updateFields.Content
	}

	if updateFields.Slug != "" {
		update["slug"] = updateFields.Slug
	}

	if updateFields.Language != "" {
		update["language"] = updateFields.Language
	}

	if updateFields.ThumbnailAltText != "" {
		update["thumbnail_alt_text"] = updateFields.ThumbnailAltText
	}

	update["updated_at"] = time.Now()

	return update
}

func (a *article) getArticleID(id string) (bson.ObjectID, error) {
	objId, err := bson.ObjectIDFromHex(id)

	if err != nil {
		a.logger.Errorf("Invalid article ID:", err)
		return bson.NilObjectID, middleware.NewBadRequestError("Invalid article ID", err)
	}
	return objId, nil
}
