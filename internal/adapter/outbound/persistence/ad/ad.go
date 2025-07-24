package ad

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"time"

	dal "github.com/CBE-Super-App/cbe-super-app-cps-action/internal/adapter/outbound/infra"
	"github.com/CBE-Super-App/cbe-super-app-cps-action/internal/adapter/outbound/model"
	"github.com/CBE-Super-App/cbe-super-app-cps-action/internal/domain/ad/entity"
	"github.com/CBE-Super-App/cbe-super-app-cps-action/internal/port/outbound/ad"

	"gitlab.com/bersufekadgetachew/cbe-super-app-shared/shared/utils"
	"go.mongodb.org/mongo-driver/v2/bson"
	"go.mongodb.org/mongo-driver/v2/mongo"

	entities "github.com/CBE-Super-App/cbe-super-app-cps-action/internal/domain/cps_actions/entities"
	common_util "github.com/CBE-Super-App/cbe-super-app-cps-action/pkgs/utils"
	constant "github.com/CBE-Super-App/cbe-super-app-cps-action/utils"
)

type ADPersistence struct {
	adDal  dal.MongoDal[model.Advert, model.Advert]
	cpsDal dal.MongoDal[model.CPSAction, model.CPSAction]
	logger utils.Logger
}

var _ ad.ADRepo = (*ADPersistence)(nil)

func InitAD(client *mongo.Client, database string, collections []string, logger utils.Logger) ad.ADRepo {
	adDal := dal.NewMongoDal[model.Advert, model.Advert](client, database, collections[0])
	cpsDal := dal.NewMongoDal[model.CPSAction, model.CPSAction](client, database, collections[1])
	return &ADPersistence{
		adDal:  adDal,
		cpsDal: cpsDal,
		logger: logger,
	}
}

func (a *ADPersistence) GetOneAdvert(ctx context.Context, id string) (*entity.Advert, error) {
	objectID, err := bson.ObjectIDFromHex(id)
	if err != nil {
		a.logger.Errorf("invalid id provided", err)
		return nil, fmt.Errorf("INVALID_ID")
	}

	filter := bson.M{
		"_id":        objectID,
		"is_deleted": false,
	}
	projection := bson.M{}

	advert, err := a.adDal.FindOne(ctx, filter, projection)
	if err != nil {
		if errors.Is(err, mongo.ErrNoDocuments) {
			a.logger.Errorf("ad not found", err)
			return nil, fmt.Errorf("AD_NOT_FOUND")
		}
		a.logger.Errorf("failed to get ad", err)
		return nil, fmt.Errorf("AD_NOT_FOUND")
	}
	return model.ToAdvertDomain(*advert), nil
}

func (a *ADPersistence) GetAllAdvert(ctx context.Context, filterParams *constant.Filter) (*common_util.PaginatedResponse[[]*entity.Advert], error) {
	filter := bson.M{"is_deleted": false}
	projection := bson.M{}

	page := filterParams.Page
	limit := filterParams.PerPage
	skip := (page - 1) * limit

	var ads []*entity.Advert
	ad, err := a.adDal.FindAllWithPagination(ctx, filter, projection, int64(skip), int64(filterParams.PerPage))
	if err != nil {
		a.logger.Errorf("failed to get ad data", err)
		return nil, fmt.Errorf("AD_NOT_FOUND")

	}

	for _, doc := range ad {
		ads = append(ads, model.ToAdvertDomain(*doc))
	}

	total, err := a.adDal.TotalCount(ctx, filter)
	if err != nil {
		a.logger.Errorf("failed to get ad total counts", err)
		return nil, fmt.Errorf("AD_NOT_FOUND")

	}

	meta := common_util.BuildPaginationMeta(total, page, limit)

	return &common_util.PaginatedResponse[[]*entity.Advert]{
		Data: ads,
		Meta: meta,
	}, nil
}

func (a *ADPersistence) HandleAdvertCreate(ctx context.Context, cpsRes *entities.CPSAction) (*entities.CPSAction, error) {
	var actionData entity.Advert

	bytes, err := json.Marshal(cpsRes.CurrentAction)
	if err != nil {
		return nil, fmt.Errorf("failed to marshal CurrentAction: %w", err)
	}

	if err := json.Unmarshal(bytes, &actionData); err != nil {
		return nil, fmt.Errorf("failed to unmarshal CurrentAction into Advert: %w", err)
	}

	req := entity.Advert{
		Title:         actionData.Title,
		Description:   actionData.Description,
		BannerImage:   actionData.BannerImage,
		AdvertFor:     actionData.AdvertFor,
		Date:          actionData.Date,
		CreatedAt:     time.Now(),
		LastUpdatedAt: time.Now(),
	}

	adDoc, err := model.ToAdvertDocument(&req)
	if err != nil {
		return nil, err
	}

	advert, err := a.adDal.InsertOne(ctx, adDoc)
	if err != nil {
		a.logger.Errorf("failed to create advert: %v", err)
		return nil, fmt.Errorf("FAILED_TO_CREATE_ADVERT")
	}

	cpsRes.CurrentAction = advert
	return cpsRes, nil
}

func (a *ADPersistence) HandleAdvertUpdate(ctx context.Context, cpsRes *entities.CPSAction) (*entities.CPSAction, error) {
	var actionData entity.UpdateAdvert

	bytes, err := json.Marshal(cpsRes.CurrentAction)
	if err != nil {
		return nil, fmt.Errorf("failed to marshal CurrentAction: %w", err)
	}

	if err := json.Unmarshal(bytes, &actionData); err != nil {
		return nil, fmt.Errorf("failed to unmarshal CurrentAction into UpdateAdvert: %w", err)
	}
	var prevData entity.Advert

	bytes, err = json.Marshal(cpsRes.PreviousAction)
	if err != nil {
		return nil, fmt.Errorf("failed to marshal PreviousAction: %w", err)
	}

	if err := json.Unmarshal(bytes, &prevData); err != nil {
		return nil, fmt.Errorf("failed to unmarshal PreviousAction into Advert: %w", err)
	}

	objId, err := bson.ObjectIDFromHex(prevData.ID)

	fmt.Println(prevData.ID, "id")

	if err != nil {
		return nil, fmt.Errorf("INVALID_ID")
	}

	filter := bson.M{"_id": objId}
	update := bson.M{}

	if actionData.Title != "" {
		update["title"] = actionData.Title
	}
	if actionData.Description != "" {
		update["description"] = actionData.Description
	}
	if actionData.BannerImage != "" {
		update["banner_image"] = actionData.BannerImage
	}
	if !actionData.Date.StartedAt.IsZero() {
		update["date.started_at"] = actionData.Date.StartedAt
	}
	if !actionData.Date.ExpiredAt.IsZero() {
		update["date.expired_at"] = actionData.Date.ExpiredAt
	}
	update["last_updated_at"] = time.Now()

	advert, err := a.adDal.UpdateOne(ctx, filter, update)
	if err != nil {
		if errors.Is(err, mongo.ErrNoDocuments) {
			return nil, fmt.Errorf("NOT_FOUND")
		}
		a.logger.Errorf("failed to update advert: %v", err)
		return nil, fmt.Errorf("GENERAL_DB_UPDATE_FAILED")
	}

	cpsRes.CurrentAction = advert
	return cpsRes, nil
}

func (a *ADPersistence) HandleAdvertDelete(ctx context.Context, cpsRes *entities.CPSAction) (*entities.CPSAction, error) {
	var prevData entity.DeletedAdvert

	bytes, err := json.Marshal(cpsRes.PreviousAction)
	if err != nil {
		return nil, fmt.Errorf("failed to marshal PreviousAction: %w", err)
	}

	if err := json.Unmarshal(bytes, &prevData); err != nil {
		return nil, fmt.Errorf("failed to unmarshal PreviousAction into DeletedAdvert: %w", err)
	}

	objId, err := bson.ObjectIDFromHex(prevData.ID)

	if err != nil {
		return nil, fmt.Errorf("INVALID_ID")
	}

	filter := bson.M{"_id": objId}
	update := bson.M{
		"is_deleted": true,
		"deleted_at": time.Now(),
	}

	advert, err := a.adDal.UpdateOne(ctx, filter, update)
	if err != nil {
		if errors.Is(err, mongo.ErrNoDocuments) {
			return nil, fmt.Errorf("NOT_FOUND")
		}
		a.logger.Errorf("failed to delete advert: %v", err)
		return nil, fmt.Errorf("GENERAL_DB_UPDATE_FAILED")
	}

	cpsRes.CurrentAction = advert
	return cpsRes, nil
}
