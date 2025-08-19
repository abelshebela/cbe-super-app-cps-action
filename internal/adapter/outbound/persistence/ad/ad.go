package ad

import (
	"context"
	"errors"
	"fmt"
	"time"

	"github.com/CBE-Super-App/cbe-super-app-cps-action/internal/adapter/outbound/mappers"
	"github.com/CBE-Super-App/cbe-super-app-cps-action/internal/adapter/outbound/model"

	entity "github.com/CBE-Super-App/cbe-super-app-cps-action/internal/domain/ad"
	"github.com/CBE-Super-App/cbe-super-app-cps-action/internal/port/outbound/ad"
	shared "github.com/CBE-Super-App/cbe-super-app-cps-action/pkgs/utils"
	util_constant "github.com/CBE-Super-App/cbe-super-app-cps-action/utils"
	"gitlab.com/bersufekadgetachew/cbe-super-app-shared/shared/dal"
	"gitlab.com/bersufekadgetachew/cbe-super-app-shared/shared/utils"
	"go.mongodb.org/mongo-driver/v2/bson"
	"go.mongodb.org/mongo-driver/v2/mongo"
)

// ADPersistence implements ADRepository
type ADPersistence struct {
	adDal  dal.MongoDal[model.Advert, model.Advert]
	logger utils.Logger
}

// InitAD initializes the ADPersistence
func InitAD(client *mongo.Client, database string, collection string, logger utils.Logger) ad.ADRepository {
	adDal := dal.NewMongoDal[model.Advert, model.Advert](client, database, collection)
	return &ADPersistence{
		adDal:  adDal,
		logger: logger,
	}
}

// CreateAdvert creates a new advert in the database
func (a *ADPersistence) CreateAdvert(ctx context.Context, advert *entity.Advert) (*entity.Advert, error) {
	a.logger.Infof("Creating advert, title: %s", advert.Title)

	adDoc, err := mappers.ToAdvertDocument(advert)
	if err != nil {
		a.logger.Errorf("Failed to convert advert to document: %v", err)
		return nil, fmt.Errorf(shared.GeneralDBInsertFailed)
	}

	result, err := a.adDal.InsertOne(ctx, *adDoc)
	if err != nil {
		a.logger.Errorf("Failed to insert advert: %v", err)
		return nil, fmt.Errorf(shared.GeneralDBInsertFailed)
	}

	a.logger.Infof("Advert created successfully, id: %s", result.ID)
	return mappers.ToAdvertDomain(result), nil
}

// UpdateAdvert updates an existing advert
func (a *ADPersistence) UpdateAdvert(ctx context.Context, advert *entity.Advert) (*entity.Advert, error) {
	a.logger.Infof("Updating advert, id: %s", advert.ID)

	objectID, err := bson.ObjectIDFromHex(advert.ID)
	if err != nil {
		a.logger.Errorf("Invalid id provided: %v", err)
		return nil, fmt.Errorf(shared.InvalidID)
	}

	update := bson.M{
		"title":           advert.Title,
		"description":     advert.Description,
		"banner_image":    advert.BannerImage,
		"advert_for":      advert.AdvertFor,
		"date.started_at": advert.Date.StartedAt,
		"date.expired_at": advert.Date.ExpiredAt,
		"enabled":         advert.Enabled,
		"last_updated_at": advert.LastUpdatedAt,
	}

	filter := bson.M{"_id": objectID}
	result, err := a.adDal.UpdateOne(ctx, filter, update)
	if err != nil {
		if errors.Is(err, mongo.ErrNoDocuments) {
			a.logger.Errorf("Advert not found, id: %s", advert.ID)
			return nil, fmt.Errorf(shared.NotFound)
		}
		a.logger.Errorf("Failed to update advert: %v", err)
		return nil, fmt.Errorf(shared.GeneralDBUpdateFailed)
	}

	a.logger.Infof("Advert updated successfully, id: %s", advert.ID)
	return mappers.ToAdvertDomain(result), nil
}

// DeleteAdvert soft-deletes an advert
func (a *ADPersistence) DeleteAdvert(ctx context.Context, id string) (*entity.Advert, error) {
	a.logger.Infof("Deleting advert, id: %s", id)

	objectID, err := bson.ObjectIDFromHex(id)
	if err != nil {
		a.logger.Errorf("Invalid id provided: %v", err)
		return nil, fmt.Errorf(shared.InvalidID)
	}

	filter := bson.M{"_id": objectID}
	update := bson.M{
		"is_deleted":      true,
		"deleted_at":      time.Now(),
		"last_updated_at": time.Now(),
	}

	result, err := a.adDal.UpdateOne(ctx, filter, update)
	if err != nil {
		if errors.Is(err, mongo.ErrNoDocuments) {
			a.logger.Errorf("Advert not found, id: %s", id)
			return nil, fmt.Errorf(shared.NotFound)
		}
		a.logger.Errorf("Failed to delete advert: %v", err)
		return nil, fmt.Errorf(shared.GeneralDBUpdateFailed)
	}

	a.logger.Infof("Advert deleted successfully, id: %s", id)
	return mappers.ToAdvertDomain(result), nil
}

// EnableDisableAdvert enables or disables an advert
func (a *ADPersistence) EnableDisableAdvert(ctx context.Context, id string, enable bool) (*entity.Advert, error) {
	a.logger.Infof("EnableDisable advert, id: %s, enable: %v", id, enable)

	objectID, err := bson.ObjectIDFromHex(id)
	if err != nil {
		a.logger.Errorf("Invalid id provided: %v", err)
		return nil, fmt.Errorf(shared.InvalidID)
	}

	filter := bson.M{"_id": objectID}
	update := bson.M{
		"enabled":         enable,
		"last_updated_at": time.Now(),
	}

	result, err := a.adDal.UpdateOne(ctx, filter, update)
	if err != nil {
		if errors.Is(err, mongo.ErrNoDocuments) {
			a.logger.Errorf("Advert not found, id: %s", id)
			return nil, fmt.Errorf(shared.NotFound)
		}
		a.logger.Errorf("Failed to update advert enable status: %v", err)
		return nil, fmt.Errorf(shared.GeneralDBUpdateFailed)
	}

	a.logger.Infof("Advert enable/disable successful, id: %s", id)
	return mappers.ToAdvertDomain(result), nil
}

// FetchAdvertByID retrieves an advert by ID
func (a *ADPersistence) FetchAdvertByID(ctx context.Context, id string) (*entity.Advert, error) {
	a.logger.Infof("Fetching advert by ID, id: %s", id)

	objectID, err := bson.ObjectIDFromHex(id)
	if err != nil {
		a.logger.Errorf("Invalid id provided: %v", err)
		return nil, fmt.Errorf(shared.InvalidID)
	}

	filter := bson.M{
		"_id":        objectID,
		"is_deleted": false,
	}
	projection := bson.M{}

	advert, err := a.adDal.FindOne(ctx, filter, projection)
	if err != nil {
		if errors.Is(err, mongo.ErrNoDocuments) {
			a.logger.Errorf("Advert not found, id: %s", id)
			return nil, fmt.Errorf(shared.NotFound)
		}
		a.logger.Errorf("Failed to fetch advert: %v", err)
		return nil, fmt.Errorf(shared.GeneralDBQueryFailed)
	}

	a.logger.Infof("Advert fetched successfully, id: %s", id)
	return mappers.ToAdvertDomain(*advert), nil
}

// FetchAdverts retrieves adverts with pagination and filtering
func (a *ADPersistence) FetchAdverts(ctx context.Context, filterParams *util_constant.Filter) (*shared.PaginatedResponse[[]*entity.Advert], error) {
	a.logger.Infof("Fetching adverts with filter: %v", filterParams)

	filter := bson.M{"is_deleted": false}
	projection := bson.M{}

	if filterParams.Search != "" {
		searchRegex := bson.M{"$regex": filterParams.Search, "$options": "i"}
		filter["$or"] = []bson.M{
			{"title": searchRegex},
			{"description": searchRegex},
		}
	}

	page := filterParams.Page
	limit := filterParams.PerPage
	skip := (page - 1) * limit

	ads, err := a.adDal.FindAllWithPagination(ctx, filter, projection, int64(skip), int64(limit))
	if err != nil {
		a.logger.Errorf("Failed to fetch adverts: %v", err)
		return nil, fmt.Errorf(shared.GeneralDBQueryFailed)
	}

	var result []*entity.Advert
	for _, doc := range ads {
		result = append(result, mappers.ToAdvertDomain(*doc))
	}

	total, err := a.adDal.TotalCount(ctx, filter)
	if err != nil {
		a.logger.Errorf("Failed to get advert total count: %v", err)
		return nil, fmt.Errorf(shared.GeneralDBQueryFailed)
	}

	meta := shared.BuildPaginationMeta(total, page, limit)
	a.logger.Infof("Adverts fetched successfully, count: %d", len(result))
	return &shared.PaginatedResponse[[]*entity.Advert]{
		Data: result,
		Meta: meta,
	}, nil
}
