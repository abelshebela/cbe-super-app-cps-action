package survey_sampling_persistence

import (
	"cbe-super-app-cps-action/internal/constants/lib"
	"cbe-super-app-cps-action/internal/constants/localization"
	imodel "cbe-super-app-cps-action/internal/constants/model"
	"cbe-super-app-cps-action/internal/constants/types"
	"cbe-super-app-cps-action/internal/storage"
	local_util "cbe-super-app-cps-action/pkgs/utils"
	"context"
	"errors"
	"time"

	"gitlab.com/bersufekadgetachew/cbe-super-app-shared/shared/config"
	"gitlab.com/bersufekadgetachew/cbe-super-app-shared/shared/dal"
	"gitlab.com/bersufekadgetachew/cbe-super-app-shared/shared/utils"
	"go.mongodb.org/mongo-driver/v2/bson"
	"go.mongodb.org/mongo-driver/v2/mongo"
)

type SurveySamplingStorage struct {
	dal        dal.MongoDal[imodel.SurveySamplingConfig, imodel.SurveySamplingConfig]
	collection *mongo.Collection
	logger     utils.Logger
}

func NewSurveySamplingRepository(client *mongo.Client, cfg *config.VaultConfig, dbName string, collection string, logger utils.Logger) storage.SurveySamplingRepository {
	return &SurveySamplingStorage{
		dal:        dal.NewMongoDal[imodel.SurveySamplingConfig, imodel.SurveySamplingConfig](client, cfg, dbName, collection),
		collection: client.Database(dbName).Collection(collection),
		logger:     logger,
	}
}

func (s *SurveySamplingStorage) Create(ctx context.Context, config imodel.SurveySamplingConfig) error {
	log := local_util.LoggerFromCtx(ctx, s.logger)

	config.ID = bson.NewObjectID()
	config.CreatedAt = time.Now()
	config.UpdatedAt = time.Now()

	_, err := s.dal.InsertOne(ctx, config)
	if err != nil {
		log.Errorf("[SurveySamplingStorage][Create] failed: %v", err)
		return errors.New(localization.ErrorUnexpectedError.Code)
	}
	return nil
}

func (s *SurveySamplingStorage) Update(ctx context.Context, id string, config imodel.SurveySamplingConfig) error {
	log := local_util.LoggerFromCtx(ctx, s.logger)

	objID, err := bson.ObjectIDFromHex(id)
	if err != nil {
		log.Errorf("[SurveySamplingStorage][Update] invalid id: %v", err)
		return errors.New(localization.ErrorUnexpectedError.Code)
	}

	update := bson.M{
		"$set": bson.M{
			"name":       config.Name,
			"config":     config.Config,
			"survey_url": config.SurveyURL,
			"updated_at": time.Now(),
		},
	}

	_, err = s.dal.UpdateOne(ctx, bson.M{"_id": objID}, update)
	if err != nil {
		log.Errorf("[SurveySamplingStorage][Update] failed: %v", err)
		return local_util.HandleDBError(err)
	}
	return nil
}

func (s *SurveySamplingStorage) SetEnabled(ctx context.Context, id string, enabled bool) error {
	log := local_util.LoggerFromCtx(ctx, s.logger)

	objID, err := bson.ObjectIDFromHex(id)
	if err != nil {
		log.Errorf("[SurveySamplingStorage][SetEnabled] invalid id: %v", err)
		return errors.New(localization.ErrorUnexpectedError.Code)
	}

	_, err = s.dal.UpdateOne(ctx, bson.M{"_id": objID}, bson.M{"enabled": enabled, "updated_at": time.Now()})
	if err != nil {
		log.Errorf("[SurveySamplingStorage][SetEnabled] failed: %v", err)
		return local_util.HandleDBError(err)
	}
	return nil
}

func (s *SurveySamplingStorage) Delete(ctx context.Context, id string) error {
	log := local_util.LoggerFromCtx(ctx, s.logger)

	objID, err := bson.ObjectIDFromHex(id)
	if err != nil {
		log.Errorf("[SurveySamplingStorage][Delete] invalid id: %v", err)
		return errors.New(localization.ErrorUnexpectedError.Code)
	}

	err = s.dal.DeleteOne(ctx, bson.M{"_id": objID})
	if err != nil {
		log.Errorf("[SurveySamplingStorage][Delete] failed: %v", err)
		return errors.New(localization.ErrorUnexpectedError.Code)
	}
	return nil
}

func (s *SurveySamplingStorage) FindByID(ctx context.Context, id string) (*imodel.SurveySamplingConfig, error) {
	log := local_util.LoggerFromCtx(ctx, s.logger)

	objID, err := bson.ObjectIDFromHex(id)
	if err != nil {
		log.Errorf("[SurveySamplingStorage][FindByID] invalid id: %v", err)
		return nil, errors.New(localization.ErrorUnexpectedError.Code)
	}

	result, err := s.dal.FindOne(ctx, bson.M{"_id": objID}, bson.M{})
	if err != nil {
		log.Errorf("[SurveySamplingStorage][FindByID] failed: %v", err)
		return nil, local_util.HandleDBError(err)
	}
	return result, nil
}

func (s *SurveySamplingStorage) FindByMethod(ctx context.Context, method string) (*imodel.SurveySamplingConfig, error) {
	log := local_util.LoggerFromCtx(ctx, s.logger)

	result, err := s.dal.FindOne(ctx, bson.M{"method": method}, bson.M{})
	if err != nil {
		if errors.Is(err, mongo.ErrNoDocuments) {
			return nil, nil
		}
		log.Errorf("[SurveySamplingStorage][FindByMethod] failed: %v", err)
		return nil, local_util.HandleDBError(err)
	}
	return result, nil
}

func (s *SurveySamplingStorage) FindAllWithPagination(ctx context.Context, filter types.Filter) (*types.PaginatedResponse[[]imodel.SurveySamplingConfig], error) {
	log := local_util.LoggerFromCtx(ctx, s.logger)

	searchKeys := bson.M{}
	if filter.Search != "" {
		searchRegex := bson.M{"$regex": filter.Search, "$options": "i"}
		searchKeys["$or"] = []bson.M{
			{"name": searchRegex},
			{"method": searchRegex},
		}
	}

	allowedKeys := []string{"search", "name", "method", "enabled"}
	mongoFilter, skip, limit := lib.FilterBuilder(filter, searchKeys, allowedKeys)

	results, err := s.dal.FindAllWithPaginationE(ctx, mongoFilter, bson.M{}, skip, limit)
	if err != nil {
		log.Errorf("[SurveySamplingStorage][FindAllWithPagination] failed: %v", err)
		return nil, errors.New(localization.ErrorUnexpectedError.Code)
	}

	total, err := s.dal.TotalCount(ctx, mongoFilter)
	if err != nil {
		log.Errorf("[SurveySamplingStorage][FindAllWithPagination] count failed: %v", err)
		return nil, errors.New(localization.ErrorUnexpectedError.Code)
	}

	meta := local_util.BuildPaginationMeta(total, filter.Page, filter.PerPage)
	return &types.PaginatedResponse[[]imodel.SurveySamplingConfig]{
		Data: results,
		Meta: meta,
	}, nil
}
