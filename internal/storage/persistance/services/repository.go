package services

import (
	"context"
	"errors"
	"time"

	"cbe-super-app-cps-action/internal/constants/lib"
	"cbe-super-app-cps-action/internal/constants/localization"
	// "cbe-super-app-cps-action/internal/constants/model"
	"cbe-super-app-cps-action/internal/constants/types"
	"cbe-super-app-cps-action/internal/storage"
	local_util "cbe-super-app-cps-action/pkgs/utils"

	"gitlab.com/bersufekadgetachew/cbe-super-app-shared/shared/dal"
	"gitlab.com/bersufekadgetachew/cbe-super-app-shared/shared/entities/model"
	"gitlab.com/bersufekadgetachew/cbe-super-app-shared/shared/utils"

	"go.mongodb.org/mongo-driver/v2/bson"
	"go.mongodb.org/mongo-driver/v2/mongo"
)

type ServicesStorage struct {
	dal    dal.MongoDal[model.Services, model.Services]
	logger utils.Logger
}

func NewServicesRepository(client *mongo.Client, dbName, collection string, logger utils.Logger) storage.ServicesRepository {
	return &ServicesStorage{
		dal:    dal.NewMongoDal[model.Services, model.Services](client, dbName, collection),
		logger: logger,
	}
}

func (s *ServicesStorage) Create(ctx context.Context, service *model.Services) error {
	// set defaults
	if service.ID == bson.NilObjectID {
		service.ID = bson.NewObjectID()
	}
	if service.CreatedAt.IsZero() {
		service.CreatedAt = time.Now()
	}
	service.LastModifiedAt = time.Now()
	service.IsDeleted = false

	_, err := s.dal.InsertOne(ctx, *service)
	if err != nil {
		s.logger.Errorf("insert service failed: %v", err)
		return errors.New(localization.ErrorUnexpectedError.Code)
	}
	return nil
}

func (s *ServicesStorage) Update(ctx context.Context, id string, service *model.Services) error {
	objID, err := bson.ObjectIDFromHex(id)
	if err != nil {
		return errors.New(localization.ErrorInvalidID.Code)
	}
	filter := bson.M{"_id": objID, "is_deleted": false}
	update := bson.M{
		"last_modified_at": time.Now(),
	}
	if service.ServiceCode != "" {
		update["service_code"] = service.ServiceCode
	}
	if service.ServiceName != "" {
		update["service_name"] = service.ServiceName
	}
	if service.ServiceType != "" {
		update["service_type"] = service.ServiceType
	}
	if service.Key != "" {
		update["key"] = service.Key
	}
	if service.PaymentType != "" {
		update["payment_type"] = service.PaymentType
	}
	if service.AboveAmount != 0 {
		update["above_amount"] = service.AboveAmount
	}
	if service.AboveServiceFee != 0 {
		update["above_service_fee"] = service.AboveServiceFee
	}
	if service.CbeIFBProductAccount != "" {
		update["cbe_ifb_product_account"] = service.CbeIFBProductAccount
	}
	if service.CbeGLProductAccount != "" {
		update["cbe_gl_product_account"] = service.CbeGLProductAccount
	}
	if len(service.Tiers) > 0 {
		update["tiers"] = service.Tiers
	}
	if (service.Cap != model.Cap{}) {
		update["cap"] = service.Cap
	}
	if (service.CbeProductCodes != model.ProductCodes{}) {
		update["cbe_product_codes"] = service.CbeProductCodes
	}
	if (service.CbeIfbProductCodes != model.ProductCodes{}) {
		update["cbe_ifb_product_codes"] = service.CbeIfbProductCodes
	}

	if len(update) == 1 { // only last_modified_at
		return errors.New(localization.ErrorNoDataProvided.Code)
	}
	_, err = s.dal.UpdateOne(ctx, filter, update)
	if err != nil {
		if errors.Is(err, mongo.ErrNoDocuments) {
			return errors.New(localization.ErrorFileNotFound.Code)
		}
		s.logger.Errorf("update service failed: %v", err)
		return errors.New(localization.ErrorUnexpectedError.Code)
	}
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
	_, err = s.dal.UpdateOne(ctx, filter, update)
	if err != nil {
		if errors.Is(err, mongo.ErrNoDocuments) {
			return errors.New(localization.ErrorFileNotFound.Code)
		}
		s.logger.Errorf("enable/disable service failed: %v", err)
		return errors.New(localization.ErrorUnexpectedError.Code)
	}
	return nil
}

func (s *ServicesStorage) FindByID(ctx context.Context, id string) (*model.Services, error) {
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

func (s *ServicesStorage) FindAllWithPagination(ctx context.Context, filterParam types.Filter) (*types.PaginatedResponse[[]*model.Services], error) {
	allowed := []string{"service_name", "service_code", "service_type", "enabled"}
	filter, skip, limit := lib.FilterBuilder(filterParam, bson.M{}, allowed)
	if filterParam.Search != "" {
		q := bson.M{"$regex": filterParam.Search, "$options": "i"}
		filter["$or"] = []bson.M{{"service_name": q}, {"service_code": q}, {"service_type": q}}
	}
	items, err := s.dal.FindAllWithPagination(ctx, filter, bson.M{}, skip, limit)
	if err != nil {
		return nil, errors.New(localization.ErrorUnexpectedError.Code)
	}
	total, err := s.dal.TotalCount(ctx, filter)
	if err != nil {
		return nil, errors.New(localization.ErrorUnexpectedError.Code)
	}
	meta := local_util.BuildPaginationMeta(total, filterParam.Page, filterParam.PerPage)
	return &types.PaginatedResponse[[]*model.Services]{
		Data: items,
		Meta: meta,
	}, nil
}
