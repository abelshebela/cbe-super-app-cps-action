package media

import (
	"cbe-super-app-cps-action/internal/constants/localization"
	"cbe-super-app-cps-action/internal/constants/model"
	"cbe-super-app-cps-action/internal/handlers/middleware"
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
	invalidID = "Invalid short video ID"
)

type shortVideoRepo struct {
	logger        shared_utils.Logger
	shortVideoDal dal.MongoDal[model.ShortVideo, model.ShortVideo]
	client        *mongo.Client
}

func NewShortVideoRepository(logger shared_utils.Logger, client *mongo.Client, dbName, collectionName string) storage.ShortVideoRepository {
	return &shortVideoRepo{
		logger:        logger,
		shortVideoDal: dal.NewMongoDal[model.ShortVideo, model.ShortVideo](client, dbName, collectionName),
		client:        client,
	}
}

func (s *shortVideoRepo) Create(ctx context.Context, shortVideo *model.ShortVideo) error {
	_, err := s.shortVideoDal.InsertOne(ctx, *shortVideo)
	if err != nil {
		return errors.New(localization.ErrorUnexpectedError.Code)
	}
	return nil
}

func (s *shortVideoRepo) Update(ctx context.Context, shortVideo *model.ShortVideo, id string) error {
	objId, err := bson.ObjectIDFromHex(id)

	if err != nil {
		s.logger.Errorf(invalidID, err)
		return middleware.NewBadRequestError(invalidID, err)
	}
	filter := bson.M{"_id": objId, "is_deleted": false}
	update := buildShortVideoUpdate(*shortVideo)
	_, err = s.shortVideoDal.UpdateOne(ctx, filter, update)
	if err != nil {
		s.logger.Errorf("Error updating short video in database:", err)
		if err == mongo.ErrNoDocuments {
			return errors.New(localization.ErrorFileNotFound.Code)
		}
		return errors.New(localization.ErrorUnexpectedError.Code)
	}
	return nil
}

func (s *shortVideoRepo) Delete(ctx context.Context, id string) error {
	objId, err := bson.ObjectIDFromHex(id)

	if err != nil {
		s.logger.Errorf("Invalid short video ID:", err)
		return middleware.NewBadRequestError(invalidID, err)
	}
	filter := bson.M{"_id": objId, "is_deleted": false}
	update := bson.M{
		"is_deleted": true,
		"updated_at": time.Now(),
		"deleted_at": time.Now(),
	}
	_, err = s.shortVideoDal.UpdateOne(ctx, filter, update)
	if err != nil {
		s.logger.Errorf("Error deleting short video in database:", err)
		if err == mongo.ErrNoDocuments {
			return errors.New(localization.ErrorFileNotFound.Code)
		}
		return errors.New(localization.ErrorUnexpectedError.Code)
	}
	return nil
}

func (s *shortVideoRepo) PublishUnpublish(ctx context.Context, id string, isPublished bool) error {
	objId, err := bson.ObjectIDFromHex(id)

	if err != nil {
		s.logger.Errorf("Invalid short video ID:", err)
		return middleware.NewBadRequestError(invalidID, err)
	}
	filter := bson.M{"_id": objId, "is_deleted": false}
	update := bson.M{
		"is_published": isPublished,
		"updated_at":   time.Now(),
	}

	if isPublished {
		update["published_at"] = time.Now()
	}

	_, err = s.shortVideoDal.UpdateOne(ctx, filter, update)
	if err != nil {
		s.logger.Errorf("Error updating short video publish status in database:", err)
		if err == mongo.ErrNoDocuments {
			return errors.New(localization.ErrorFileNotFound.Code)
		}
		return errors.New(localization.ErrorUnexpectedError.Code)
	}
	return nil
}

func buildShortVideoUpdate(updateFields model.ShortVideo) bson.M {
	update := bson.M{}
	if updateFields.Title != "" {
		update["title"] = updateFields.Title
	}
	if updateFields.Caption != "" {
		update["caption"] = updateFields.Caption
	}
	if len(updateFields.Tags) > 0 {
		update["tags"] = updateFields.Tags
	}
	if updateFields.Thumbnail != "" {
		update["thumbnail"] = updateFields.Thumbnail
	}
	if updateFields.VideoURL != "" {
		update["video_url"] = updateFields.VideoURL
	}
	if updateFields.Duration > 0 {
		update["duration"] = updateFields.Duration
	}
	if updateFields.CategoryID != bson.NilObjectID {
		update["category_id"] = updateFields.CategoryID
	}
	if updateFields.Author != "" {
		update["author"] = updateFields.Author
	}

	if updateFields.Slug != "" {
		update["slug"] = updateFields.Slug
	}

	if updateFields.ThumbnailAltText != "" {
		update["thumbnail_alt_text"] = updateFields.ThumbnailAltText
	}

	update["updated_at"] = time.Now()

	return update
}
