package services

import (
	"context"
	"errors"
	"time"

	"cbe-super-app-cps-action/internal/constants"
	"cbe-super-app-cps-action/internal/constants/lib"
	"cbe-super-app-cps-action/internal/constants/localization"
	"cbe-super-app-cps-action/internal/storage/kafka"

	"cbe-super-app-cps-action/internal/constants/types"
	"cbe-super-app-cps-action/internal/storage"
	"cbe-super-app-cps-action/internal/storage/persistance/services/core"
	local_util "cbe-super-app-cps-action/pkgs/utils"

	"gitlab.com/bersufekadgetachew/cbe-super-app-shared/shared/dal"

	"gitlab.com/bersufekadgetachew/cbe-super-app-shared/shared/entities/model"
	"gitlab.com/bersufekadgetachew/cbe-super-app-shared/shared/utils"

	"gitlab.com/bersufekadgetachew/cbe-super-app-shared/shared/config"
	"go.mongodb.org/mongo-driver/v2/bson"
	"go.mongodb.org/mongo-driver/v2/mongo"
)

type ServicesStorage struct {
	cfg           *config.VaultConfig
	dal           dal.MongoDal[model.Service, model.Service]
	serviceDal    dal.MongoDal[model.ServiceList, model.ServiceList]
	kafkaProducer kafka.ClientOrchestrationProducer
	logger        utils.Logger
}

func NewServicesRepository(client *mongo.Client, cfg *config.VaultConfig, dbName, collection string, kafkaProducer kafka.ClientOrchestrationProducer, logger utils.Logger) storage.ServicesRepository {
	return &ServicesStorage{
		cfg:           cfg,
		dal:           dal.NewMongoDal[model.Service, model.Service](client, cfg, dbName, collection),
		serviceDal:    dal.NewMongoDal[model.ServiceList, model.ServiceList](client, cfg, dbName, "service_list"),
		kafkaProducer: kafkaProducer,
		logger:        logger,
	}
}

func (s *ServicesStorage) Create(ctx context.Context, service *model.Service) error {
	if service.ID == bson.NilObjectID {
		service.ID = bson.NewObjectID()
	}
	if service.CreatedAt.IsZero() {
		service.CreatedAt = time.Now()
	}
	service.LastModifiedAt = time.Now()
	service.IsDeleted = false

	createService, err := s.dal.InsertOne(ctx, *service)
	if err != nil {
		s.logger.Errorf("insert service failed: %v", err)
		return errors.New(localization.ErrorUnexpectedError.Code)
	}
	s.kafkaProducer.PublishMessage(
		ctx,
		createService,
		string(constants.ClientOrchestrationServicesTopic),
		s.cfg.CPSServiceUpdate,
		"service authorized and updated",
	)
	return nil
}

func (s *ServicesStorage) Update(ctx context.Context, id string, service *model.Service) error {
	objID, err := bson.ObjectIDFromHex(id)
	if err != nil {
		return errors.New(localization.ErrorInvalidID.Code)
	}

	filter := bson.M{
		"_id":        objID,
		"is_deleted": false,
	}

	prev, err := s.dal.FindOne(ctx, bson.M{"_id": objID}, bson.M{})
	if err != nil {
		return err
	}

	update := core.MapToServiceUpdate(*service, *prev)

	if len(update) == 1 {
		return errors.New(localization.ErrorNoDataProvided.Code)
	}

	updatedService, err := s.dal.UpdateOne(ctx, filter, update)
	if err != nil {
		if errors.Is(err, mongo.ErrNoDocuments) {
			return errors.New(localization.ErrorFileNotFound.Code)
		}
		s.logger.Errorf("update service failed: %v", err)
		return errors.New(localization.ErrorUnexpectedError.Code)
	}

	s.kafkaProducer.PublishMessage(
		ctx,
		updatedService,
		string(constants.ClientOrchestrationServicesTopic),
		s.cfg.CPSServiceUpdate,
		"service authorized and updated",
	)

	return nil
}

func (s *ServicesStorage) Delete(ctx context.Context, id string) error {
	objID, err := bson.ObjectIDFromHex(id)
	if err != nil {
		return errors.New(localization.ErrorInvalidID.Code)
	}
	filter := bson.M{"_id": objID, "is_deleted": false}
	update := bson.M{"is_deleted": true, "deleted_at": time.Now()}
	_, err = s.dal.UpdateOne(ctx, filter, update)
	if err != nil {
		if errors.Is(err, mongo.ErrNoDocuments) {
			return errors.New(localization.ErrorFileNotFound.Code)
		}
		s.logger.Errorf("delete service failed: %v", err)
		return errors.New(localization.ErrorUnexpectedError.Code)
	}
	return nil
}

func (s *ServicesStorage) EnableOrDisable(ctx context.Context, id string, enable bool) error {
	objID, err := bson.ObjectIDFromHex(id)
	if err != nil {
		return errors.New(localization.ErrorInvalidID.Code)
	}
	filter := bson.M{"_id": objID, "is_deleted": false}
	update := bson.M{"enabled": enable, "last_modified_at": time.Now()}
	updatedService, err := s.dal.UpdateOne(ctx, filter, update)
	if err != nil {
		if errors.Is(err, mongo.ErrNoDocuments) {
			return errors.New(localization.ErrorFileNotFound.Code)
		}
		s.logger.Errorf("enable/disable service failed: %v", err)
		return errors.New(localization.ErrorUnexpectedError.Code)
	}

	s.kafkaProducer.PublishMessage(
		ctx,
		updatedService,
		string(constants.ClientOrchestrationServicesTopic),
		s.cfg.CPSServiceUpdate,
		"service authorized and updated",
	)

	return nil
}

func (s *ServicesStorage) FindByID(ctx context.Context, id string) (*model.Service, error) {
	objID, err := bson.ObjectIDFromHex(id)
	if err != nil {
		return nil, errors.New(localization.ErrorInvalidID.Code)
	}
	filter := bson.M{"_id": objID, "is_deleted": false}
	doc, err := s.dal.FindOne(ctx, filter, nil)
	if err != nil {
		if errors.Is(err, mongo.ErrNoDocuments) {
			return nil, errors.New(localization.ErrorFileNotFound.Code)
		}
		s.logger.Errorf("find service by id failed: %v", err)
		return nil, errors.New(localization.ErrorUnexpectedError.Code)
	}
	return doc, nil
}

func (s *ServicesStorage) FindAllWithPagination(ctx context.Context, filterParam types.Filter) (*types.PaginatedResponse[[]model.Service], error) {
	allowed := []string{"service_name", "service_code", "service_key", "enabled"}
	filter, skip, limit := lib.FilterBuilder(filterParam, bson.M{}, allowed)
	if filterParam.Search != "" {
		q := bson.M{"$regex": filterParam.Search, "$options": "i"}
		filter["$or"] = []bson.M{
			{"service_name": q},
			{"service_code": q},
			{"service_key": q},
			{"service_list.service_name": q},
			{"service_list.service_key": q},
		}
	}
	// if filterParam.Search != "" {
	// 	q := bson.M{"$regex": filterParam.Search, "$options": "i"}

	// 	filter["$or"] = []bson.M{
	// 		{"service_name": q},
	// 		{"service_code": q},
	// 		{"service_key": q},
	// 		{
	// 			"service_list": bson.M{
	// 				"$elemMatch": bson.M{
	// 					"$or": []bson.M{
	// 						{"service_name": q},
	// 						{"service_key": q},
	// 					},
	// 				},
	// 			},
	// 		},
	// 	}
	// }

	items, err := s.dal.FindAllWithPaginationE(ctx, filter, bson.M{}, skip, limit)
	if err != nil {
		s.logger.Errorf("failed to get all services: %v", err)
		return nil, errors.New(localization.ErrorUnexpectedError.Code)
	}
	total, err := s.dal.TotalCount(ctx, filter)
	if err != nil {
		return nil, errors.New(localization.ErrorUnexpectedError.Code)
	}
	meta := local_util.BuildPaginationMeta(total, filterParam.Page, filterParam.PerPage)
	return &types.PaginatedResponse[[]model.Service]{
		Data: items,
		Meta: meta,
	}, nil
}

func (s *ServicesStorage) FindAllServiceListWithPagination(ctx context.Context, filterParam types.Filter) (*types.PaginatedResponse[[]model.ServiceList], error) {
	allowed := []string{"service_name", "service_code", "service_type", "enabled"}
	filter, skip, limit := lib.FilterBuilder(filterParam, bson.M{}, allowed)
	if filterParam.Search != "" {
		q := bson.M{"$regex": filterParam.Search, "$options": "i"}
		filter["$or"] = []bson.M{{"service_name": q}, {"service_code": q}, {"service_type": q}}
	}

	// if count < limit {
	// 	limit = count
	// }

	items, err := s.serviceDal.FindAllWithPagination(ctx, filter, bson.M{}, skip, limit)
	if err != nil {
		return nil, errors.New(localization.ErrorUnexpectedError.Code)
	}
	total, err := s.dal.TotalCount(ctx, filter)
	if err != nil {
		return nil, errors.New(localization.ErrorUnexpectedError.Code)
	}

	// count, err := s.dal.TotalCount(ctx, bson.M{})
	// if err != nil {
	// 	return nil, errors.New(localization.ErrorUnexpectedError.Code)
	// }
	meta := local_util.BuildPaginationMeta(total, filterParam.Page, filterParam.PerPage)
	return &types.PaginatedResponse[[]model.ServiceList]{
		Data: items,
		Meta: meta,
	}, nil
}

func (s *ServicesStorage) CheckServiceExistence(ctx context.Context, serviceCode, serviceKey, serviceName string) (bool, error) {
	orConditions := make([]bson.M, 0, 3)
	if serviceKey != "" {
		orConditions = append(orConditions, bson.M{"service_key": serviceKey})
	}
	// if serviceCode != "" {
	// 	orConditions = append(orConditions, bson.M{"service_code": serviceCode})
	// }
	if serviceName != "" {
		orConditions = append(orConditions, bson.M{"service_name": serviceName})
	}

	if len(orConditions) == 0 {
		return false, nil
	}

	filter := bson.M{
		"$or":        orConditions,
		"is_deleted": false,
	}

	_, err := s.dal.FindOne(ctx, filter, nil)
	if err != nil {
		if errors.Is(err, mongo.ErrNoDocuments) {
			return false, nil
		}
		s.logger.Errorf("failed to check service existence in services: %v", err)
		return false, errors.New(localization.ErrorUnexpectedError.Code)
	}

	return true, nil
}
