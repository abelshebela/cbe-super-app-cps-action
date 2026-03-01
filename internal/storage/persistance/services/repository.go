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
	accessDal     dal.MongoDal[model.APPAccessList, model.APPAccessList]
	kafkaProducer kafka.ClientOrchestrationProducer
	logger        utils.Logger
}

func NewServicesRepository(client *mongo.Client, cfg *config.VaultConfig, dbName, collection string, kafkaProducer kafka.ClientOrchestrationProducer, logger utils.Logger) storage.ServicesRepository {
	return &ServicesStorage{
		cfg:           cfg,
		dal:           dal.NewMongoDal[model.Service, model.Service](client, cfg, dbName, collection),
		serviceDal:    dal.NewMongoDal[model.ServiceList, model.ServiceList](client, cfg, dbName, "service_list"),
		accessDal:     dal.NewMongoDal[model.APPAccessList, model.APPAccessList](client, cfg, dbName, "access_list"),
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
		s.logger.Errorf("[ServicesStorage][Create] failed to insert service: %v", err)
		return errors.New(localization.ErrorUnexpectedError.Code)
	}
	s.kafkaProducer.PublishMessage(ctx, createService, string(constants.ClientOrchestrationServicesTopic), string(constants.ClientOrchestrationServicesTopic), "new service created")
	return nil
}

func (s *ServicesStorage) Update(ctx context.Context, id string, service *model.Service) error {
	objID, err := bson.ObjectIDFromHex(id)
	if err != nil {
		s.logger.Errorf("[ServicesStorage][Update] invalid object id: %v", err)
		return errors.New(localization.ErrorInvalidID.Code)
	}

	filter := bson.M{
		"_id":        objID,
		"is_deleted": false,
	}

	prev, err := s.dal.FindOne(ctx, bson.M{"_id": objID}, bson.M{})
	if err != nil {
		s.logger.Errorf("[ServicesStorage][Update] failed to find previous service: %v", err)
		return local_util.HandleDBError(err)
	}

	update := core.MapToServiceUpdate(*service, *prev)

	if len(update) == 1 {
		return errors.New(localization.ErrorNoDataProvided.Code)
	}

	updatedService, err := s.dal.UpdateOne(ctx, filter, update)
	if err != nil {
		s.logger.Errorf("[ServicesStorage][Update] failed to update service: %v", err)
		return local_util.HandleDBError(err)
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
		s.logger.Errorf("[ServicesStorage][Delete] invalid object id: %v", err)
		return errors.New(localization.ErrorInvalidID.Code)
	}
	filter := bson.M{"_id": objID, "is_deleted": false}
	update := bson.M{"is_deleted": true, "deleted_at": time.Now()}
	_, err = s.dal.UpdateOne(ctx, filter, update)
	if err != nil {
		s.logger.Errorf("[ServicesStorage][Delete] failed to delete service: %v", err)
		return local_util.HandleDBError(err)
	}
	return nil
}

func (s *ServicesStorage) EnableOrDisable(ctx context.Context, id string, enable bool) error {
	objID, err := bson.ObjectIDFromHex(id)
	if err != nil {
		s.logger.Errorf("[ServicesStorage][EnableOrDisable] invalid object id: %v", err)
		return errors.New(localization.ErrorInvalidID.Code)
	}
	filter := bson.M{"_id": objID, "is_deleted": false}
	update := bson.M{"enabled": enable, "last_modified_at": time.Now()}
	updatedService, err := s.dal.UpdateOne(ctx, filter, update)
	if err != nil {
		s.logger.Errorf("[ServicesStorage][EnableOrDisable] failed to enable/disable service: %v", err)
		return local_util.HandleDBError(err)
	}

	s.kafkaProducer.PublishMessage(ctx, updatedService, string(constants.ClientOrchestrationServicesTopic), string(constants.ClientOrchestrationServicesTopic), "service enabled/disabled")

	return nil
}

func (s *ServicesStorage) FindByID(ctx context.Context, id string) (*model.Service, error) {
	objID, err := bson.ObjectIDFromHex(id)
	if err != nil {
		s.logger.Errorf("[ServicesStorage][FindByID] invalid object id: %v", err)
		return nil, errors.New(localization.ErrorInvalidID.Code)
	}
	filter := bson.M{"_id": objID, "is_deleted": false}
	doc, err := s.dal.FindOne(ctx, filter, nil)
	if err != nil {
		s.logger.Errorf("[ServicesStorage][FindByID] failed to find service: %v", err)
		return nil, local_util.HandleDBError(err)
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
		s.logger.Errorf("[ServicesStorage][FindAllWithPagination] failed to fetch services: %v", err)
		return nil, local_util.HandleDBError(err)
	}
	total, err := s.dal.TotalCount(ctx, filter)
	if err != nil {
		s.logger.Errorf("[ServicesStorage][FindAllWithPagination] failed to count services: %v", err)
		return nil, local_util.HandleDBError(err)
	}
	meta := local_util.BuildPaginationMeta(total, filterParam.Page, filterParam.PerPage)
	return &types.PaginatedResponse[[]model.Service]{
		Data: items,
		Meta: meta,
	}, nil
}

func (s *ServicesStorage) FindAllServiceListWithPagination(ctx context.Context, filterParam types.Filter) (*types.PaginatedResponse[[]model.ServiceList], error) {
	allowed := []string{"service_name", "service_key", "is_enabled"}
	filter, skip, limit := lib.FilterBuilder(filterParam, bson.M{}, allowed)
	if filterParam.Search != "" {
		q := bson.M{"$regex": filterParam.Search, "$options": "i"}
		filter["$or"] = []bson.M{{"service_name": q}, {"service_key": q}, {"is_enabled": q}}
	}

	items, err := s.serviceDal.FindAllWithPagination(ctx, filter, bson.M{}, skip, limit)
	if err != nil {
		s.logger.Errorf("[ServicesStorage][FindAllServiceListWithPagination] failed to fetch service list: %v", err)
		return nil, local_util.HandleDBError(err)
	}
	total, err := s.serviceDal.TotalCount(ctx, filter)
	if err != nil {
		s.logger.Errorf("[ServicesStorage][FindAllServiceListWithPagination] failed to count service list: %v", err)
		return nil, local_util.HandleDBError(err)
	}

	meta := local_util.BuildPaginationMeta(total, filterParam.Page, filterParam.PerPage)
	return &types.PaginatedResponse[[]model.ServiceList]{
		Data: items,
		Meta: meta,
	}, nil
}

func (s *ServicesStorage) CheckServiceExistence(ctx context.Context, serviceCode, serviceKey, serviceName string) (bool, error) {
	// Build an OR filter with exact matches only.
	// service_key, service_code and service_name must match exactly (no regex / partial match).
	orConditions := make([]bson.M, 0, 3)
	if serviceKey != "" {
		orConditions = append(orConditions, bson.M{"service_key": serviceKey})
	}
	if serviceCode != "" {
		orConditions = append(orConditions, bson.M{"service_code": serviceCode})
	}
	if serviceName != "" {
		orConditions = append(orConditions, bson.M{"service_name": serviceName})
	}

	// If nothing is provided, there is nothing to check.
	if len(orConditions) == 0 {
		return false, nil
	}

	filter := bson.M{"$or": orConditions}

	// Fetch a single document instead of counting.
	_, err := s.dal.FindOne(ctx, filter, nil)
	if err != nil {
		if errors.Is(err, mongo.ErrNoDocuments) {
			return false, nil
		}
		s.logger.Errorf("[ServicesStorage][CheckServiceExistence] failed to check service existence: %v", err)
		return false, local_util.HandleDBError(err)
	}

	return true, nil
}

func (s *ServicesStorage) FindServiceListByID(ctx context.Context, id string) (*model.ServiceList, error) {
	objID, err := bson.ObjectIDFromHex(id)
	if err != nil {
		s.logger.Errorf("[UpdateServiceList][Update] invalid object id: %v", err)
		return nil, errors.New(localization.ErrorInvalidID.Code)
	}
	filter := bson.M{"_id": objID}
	doc, err := s.serviceDal.FindOne(ctx, filter, bson.M{})
	if err != nil {
		if errors.Is(err, mongo.ErrNoDocuments) {
			return nil, errors.New(localization.ErrorServiceListNotFound.Code)
		}
		s.logger.Errorf("[ServicesStorage][FindServiceListByID] failed to find service list: %v", err)
		return nil, local_util.HandleDBError(err)
	}

	return doc, nil
}

func (s *ServicesStorage) CreateServiceList(ctx context.Context, serviceList *model.ServiceList) error {
	serviceList.ID = bson.NewObjectID()
	serviceList.CreatedAt = time.Now()
	serviceList.LastModifiedAt = time.Now()
	serviceList.IsEnabled = true

	_, err := s.serviceDal.InsertOne(ctx, *serviceList)
	if err != nil {
		s.logger.Errorf("[ServicesStorage][CreateServiceList] failed to insert service list: %v", err)
		return errors.New(localization.ErrorUnexpectedError.Code)
	}
	_, err = s.accessDal.InsertOne(ctx, model.APPAccessList{Key: serviceList.ServiceKey, AccessListName: serviceList.ServiceName})
	if err != nil {
		s.logger.Errorf("[ServicesStorage][CreateServiceList] failed to insert access list: %v", err)
		return errors.New(localization.ErrorUnexpectedError.Code)
	}

	// s.kafkaProducer.PublishMessage(ctx, createServiceList, string(constants.ClientOrchestrationServicesTopic), string(constants.ClientOrchestrationServicesTopic), "new service list created")
	return nil

}

func (s *ServicesStorage) FindServiceListByNameOrKey(ctx context.Context, name, key string) (*model.ServiceList, error) {
	var condition []bson.M

	if name != "" {
		condition = append(condition, bson.M{"service_name": bson.M{"$regex": name, "$options": "i"}})
	}
	if key != "" {
		condition = append(condition, bson.M{"service_key": key})
	}

	if len(condition) == 0 {
		return nil, errors.New(localization.ErrorNoDataProvided.Code)
	}

	filter := bson.M{"$or": condition}

	doc, err := s.serviceDal.FindOne(ctx, filter, bson.M{})
	if err != nil {
		if errors.Is(err, mongo.ErrNoDocuments) {
			return nil, errors.New(localization.ErrorServiceListNotFound.Code)
		}
		s.logger.Errorf("[ServicesStorage][FindServiceListByNameOrKey] failed to find service list: %v", err)
		return nil, local_util.HandleDBError(err)
	}
	return doc, nil
}

func (s *ServicesStorage) UpdateServiceList(ctx context.Context, id, serviceKey string, serviceList *model.ServiceList) error {
	objID, err := bson.ObjectIDFromHex(id)
	if err != nil {
		s.logger.Errorf("[UpdateServiceList][Update] invalid object id: %v", err)
		return errors.New(localization.ErrorInvalidID.Code)
	}
	filter := bson.M{"_id": objID}

	serviceListUpdate := bson.M{}
	if serviceList.ServiceName != "" {
		serviceListUpdate["service_name"] = serviceList.ServiceName
	}
	if serviceList.ServiceKey != "" {
		serviceListUpdate["service_key"] = serviceList.ServiceKey
	}
	if len(serviceListUpdate) == 0 {
		return errors.New(localization.ErrorNoDataProvided.Code)
	}

	// update service_list collection
	_, err = s.serviceDal.UpdateOne(ctx, filter, serviceListUpdate)
	if err != nil {
		s.logger.Errorf("[UpdateServiceList] failed to update service lists: %v", err)
		return errors.New(localization.ErrorUnexpectedError.Code)
	}

	// update access_list
	accessListUpdate := bson.M{}
	if serviceList.ServiceName != "" {
		accessListUpdate["access_list_name"] = serviceList.ServiceName
	}
	if serviceList.ServiceKey != "" {
		accessListUpdate["key"] = serviceList.ServiceKey
	}
	if len(accessListUpdate) > 0 {
		accessListUpdate["last_modified_at"] = time.Now()
		_, err = s.accessDal.UpdateOne(ctx, bson.M{"key": serviceKey}, accessListUpdate)
		if err != nil {
			s.logger.Errorf("[UpdateServiceList] failed to update access lists: %v", err)
			return errors.New(localization.ErrorUnexpectedError.Code)
		}
	}

	return nil
}

func (s *ServicesStorage) EnableOrDisableServiceList(ctx context.Context, id, serviceKey string, enable bool) error {
	objID, err := bson.ObjectIDFromHex(id)
	if err != nil {
		s.logger.Errorf("[UpdateServiceList][Update] invalid object id: %v", err)
		return errors.New(localization.ErrorInvalidID.Code)
	}
	filter := bson.M{"_id": objID}
	update := bson.M{"is_enabled": enable}
	_, err = s.serviceDal.UpdateOne(ctx, filter, update)
	if err != nil {
		s.logger.Errorf("[EnableOrDisableServiceList] failed to update service list: %v", err)
		return errors.New(localization.ErrorUnexpectedError.Code)
	}

	// Update access_list too
	accessListUpdate := bson.M{"enabled": enable}
	_, err = s.accessDal.UpdateOne(ctx, bson.M{"key": serviceKey}, accessListUpdate)
	if err != nil {
		s.logger.Errorf("[EnableOrDisableServiceList] failed to update access list: %v", err)
		return errors.New(localization.ErrorUnexpectedError.Code)
	}

	return nil
}
