package donation_category

import (
	"cbe-super-app-cps-action/internal/constants"
	"cbe-super-app-cps-action/internal/constants/dto/donation_category"
	"cbe-super-app-cps-action/internal/constants/lib"
	"cbe-super-app-cps-action/internal/constants/localization"
	"cbe-super-app-cps-action/internal/constants/types"
	"cbe-super-app-cps-action/internal/storage"
	"cbe-super-app-cps-action/internal/storage/kafka"
	local_util "cbe-super-app-cps-action/pkgs/utils"
	"context"
	"errors"

	"gitlab.com/bersufekadgetachew/cbe-super-app-shared/shared/entities/model"

	"gitlab.com/bersufekadgetachew/cbe-super-app-shared/shared/dal"
	"gitlab.com/bersufekadgetachew/cbe-super-app-shared/shared/utils"
	"go.mongodb.org/mongo-driver/v2/bson"
	"go.mongodb.org/mongo-driver/v2/mongo"
	"gitlab.com/bersufekadgetachew/cbe-super-app-shared/shared/config"
)

type DonationCategoryStorage struct {
	dal           dal.MongoDal[model.DonationCategory, model.DonationCategory]
	client        *mongo.Client
	kafkaProducer kafka.ClientOrchestrationProducer
	logger        utils.Logger
}

func NewDonationCategoryRepository(client *mongo.Client,cfg *config.VaultConfig, dbName string, collection string, kafkaProducer kafka.ClientOrchestrationProducer, logger utils.Logger) storage.DonationCategoryRepository {
	return &DonationCategoryStorage{
		dal:           dal.NewMongoDal[model.DonationCategory, model.DonationCategory](client, cfg,dbName, collection),
		client:        client,
		kafkaProducer: kafkaProducer,
		logger:        logger,
	}
}

func (s *DonationCategoryStorage) Create(ctx context.Context, details *model.DonationCategory) error {
	s.logger.Infof("[Create] creating donation category")
	newDonationCategory, err := s.dal.InsertOne(ctx, *details)
	if err != nil {
		s.logger.Errorf("[Create] failed to create donation category: %v", err)
		return errors.New(localization.ErrorUnexpectedError.Code)
	}

	s.kafkaProducer.PublishMessage(ctx, newDonationCategory, string(constants.ClientOrchestrationDonationCategoryTopic), string(constants.ClientOrchestrationDonationCategoryTopic), "new donation category created")

	s.logger.Infof("[Create] donation category created successfully")
	return nil
}

func (s *DonationCategoryStorage) Update(ctx context.Context, id string, details *model.DonationCategory) error {
	s.logger.Infof("[Update] updating donation category for id: %s", id)
	objID, err := bson.ObjectIDFromHex(id)
	if err != nil {
		s.logger.Errorf("[Update] invalid object id: %v", err)
		return errors.New(localization.ErrorUnexpectedError.Code)
	}
	filter := bson.M{"_id": objID, "is_deleted": false}
	updateData := DonationCategoryMapper(*details)
	updatedDonationCategory, err := s.dal.UpdateOne(ctx, filter, updateData)
	if err != nil {
		if err == mongo.ErrNoDocuments {
			s.logger.Errorf("[Update] donation category not found")
			return errors.New(localization.ErrorFileNotFound.Code)
		}
		s.logger.Errorf("[Update] failed to update donation category: %v", err)
		return errors.New(localization.ErrorUnexpectedError.Code)
	}

	s.kafkaProducer.PublishMessage(ctx, updatedDonationCategory, string(constants.ClientOrchestrationDonationCategoryTopic), string(constants.ClientOrchestrationDonationCategoryTopic), "new donation category updated")

	s.logger.Infof("[Update] donation category updated successfully")
	return nil
}
func (s *DonationCategoryStorage) FindByID(ctx context.Context, id string) (*donation_category.DonationCategoryListResponse, error) {
	s.logger.Infof("[FindByID] fetching donation category by id: %s", id)
	idObj, ok := local_util.StringToObjectID(id)
	if !ok {
		s.logger.Errorf("[FindByID] invalid object id")
		return nil, errors.New(localization.ErrorInvalidID.Code)
	}
	filter := bson.M{"_id": idObj, "is_deleted": false}
	projection := bson.M{}

	result, err := s.dal.FindOne(ctx, filter, projection)
	if err != nil {
		s.logger.Errorf("[FindByID] failed to find donation category: %v", err)
		code, _ := local_util.HandleMongoError(err)
		return nil, errors.New(code)
	}
	s.logger.Infof("[FindByID] donation category retrieved successfully")
	return MapToDonationCategoryListResponse(result), nil
}

func (s *DonationCategoryStorage) FindByName(ctx context.Context, name string) (*donation_category.DonationCategoryListResponse, error) {
	s.logger.Infof("[FindByName] fetching donation category by name: %s", name)
	filter := bson.M{"category_name": name, "is_deleted": false}
	projection := bson.M{}

	result, err := s.dal.FindOne(ctx, filter, projection)
	if err != nil {
		s.logger.Errorf("[FindByName] failed to find donation category: %v", err)
		code, _ := local_util.HandleMongoError(err)
		return nil, errors.New(code)
	}
	s.logger.Infof("[FindByName] donation category retrieved successfully")
	return MapToDonationCategoryListResponse(result), nil
}

func (s *DonationCategoryStorage) FindAllWithPagination(ctx context.Context, filterParam types.Filter) (*types.PaginatedResponse[[]donation_category.DonationCategoryListResponse], error) {

	filter := bson.M{"is_deleted": false}
	searchKeys := bson.M{}

	// 2. Allowed filterable/searchable fields
	allowedKeys := []string{"category_name", "enabled"}
	// 3. Add search (if provided)
	if filterParam.Search != "" {
		searchRegex := bson.M{"$regex": filterParam.Search, "$options": "i"}
		searchKeys["$or"] = []bson.M{{"category_name": searchRegex, "enabled": searchRegex}}
	}

	projection := bson.M{}
	filter, skip, limit := lib.FilterBuilder(filterParam, searchKeys, allowedKeys)

	// 5. Fetch data
	data, err := s.dal.FindAllWithPagination(ctx, filter, projection, skip, limit)
	if err != nil {
		s.logger.Errorf("[FindAllWithPagination] failed to fetch donation categories: %v", err)
		return nil, errors.New(localization.ErrorUnexpectedError.Message)
	}

	// 6. Count total
	total, err := s.dal.TotalCount(ctx, filter)
	if err != nil {
		s.logger.Errorf("[FindAllWithPagination] failed to count donation categories: %v", err)
		return nil, errors.New(localization.ErrorUnexpectedError.Message)
	}

	// 7. Build pagination metadata
	meta := local_util.BuildPaginationMeta(total, filterParam.Page, filterParam.PerPage)

	// 8. Map to DTOs
	dtoData := MapToDonationCategoryListResponses(data)
	s.logger.Infof("[FindAllWithPagination] retrieved %d donation categories", len(dtoData))

	// 9. Return standard paginated response
	return &types.PaginatedResponse[[]donation_category.DonationCategoryListResponse]{
		Data: dtoData,
		Meta: meta,
	}, nil
}
func (s *DonationCategoryStorage) Delete(ctx context.Context, id string) error {
	s.logger.Infof("[Delete] deleting donation category for id: %s", id)
	objID, err := bson.ObjectIDFromHex(id)
	if err != nil {
		s.logger.Errorf("[Delete] invalid object id: %v", err)
		return errors.New(localization.ErrorUnexpectedError.Code)
	}
	filter := bson.M{"_id": objID, "is_deleted": false}
	err = s.dal.DeleteOne(ctx, filter)
	if err != nil {
		s.logger.Errorf("[Delete] failed to delete donation category: %v", err)
		return err
	}
	s.logger.Infof("[Delete] donation category deleted successfully")
	return nil
}
