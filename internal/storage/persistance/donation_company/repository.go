package donation_company

import (
	"cbe-super-app-cps-action/internal/constants"
	"cbe-super-app-cps-action/internal/constants/dto/donation_company"
	"cbe-super-app-cps-action/internal/constants/lib"
	"cbe-super-app-cps-action/internal/constants/localization"
	"cbe-super-app-cps-action/internal/constants/types"
	"cbe-super-app-cps-action/internal/storage"
	"cbe-super-app-cps-action/internal/storage/kafka"
	local_util "cbe-super-app-cps-action/pkgs/utils"
	"context"

	donation_model "gitlab.com/bersufekadgetachew/cbe-super-app-shared/shared/entities/donation"

	"errors"

	"gitlab.com/bersufekadgetachew/cbe-super-app-shared/shared/config"
	"gitlab.com/bersufekadgetachew/cbe-super-app-shared/shared/dal"
	"gitlab.com/bersufekadgetachew/cbe-super-app-shared/shared/utils"
	"go.mongodb.org/mongo-driver/v2/bson"
	"go.mongodb.org/mongo-driver/v2/mongo"
)

type DonationCompanyStorage struct {
	dal           dal.MongoDal[donation_model.DonationCompany, donation_model.DonationCompany]
	client        *mongo.Client
	kafkaProducer kafka.ClientOrchestrationProducer
	logger        utils.Logger
}

func NewDonationCompanyRepository(client *mongo.Client, cfg *config.VaultConfig, dbName string, collection string, kafkaProducer kafka.ClientOrchestrationProducer, logger utils.Logger) storage.DonationCompanyRepository {
	return &DonationCompanyStorage{
		dal:           dal.NewMongoDal[donation_model.DonationCompany, donation_model.DonationCompany](client, cfg, dbName, collection),
		client:        client,
		kafkaProducer: kafkaProducer,
		logger:        logger,
	}
}

func (s *DonationCompanyStorage) Create(ctx context.Context, details *donation_model.DonationCompany) error {
	s.logger.Infof("[DonationCompanyStorage][Create] creating donation company")
	newDonationCompany, err := s.dal.InsertOne(ctx, *details)
	if err != nil {
		s.logger.Errorf("[DonationCompanyStorage][Create] failed to create donation company: %v", err)
		return errors.New(localization.ErrorUnexpectedError.Code)
	}

	s.kafkaProducer.PublishMessage(ctx, newDonationCompany, string(constants.ClientOrchestrationDonationCompanyTopic), string(constants.ClientOrchestrationDonationCompanyTopic), "new donation company created")

	s.logger.Infof("[DonationCompanyStorage][Create] donation company created successfully")
	return nil
}

func (s *DonationCompanyStorage) Update(ctx context.Context, id string, details *donation_model.DonationCompany) error {
	s.logger.Infof("[DonationCompanyStorage][Update] updating donation company for id: %s", id)
	objID, err := bson.ObjectIDFromHex(id)
	if err != nil {
		s.logger.Errorf("[DonationCompanyStorage][Update] invalid object id: %v", err)
		return errors.New(localization.ErrorUnexpectedError.Code)
	}
	filter := bson.M{"_id": objID}

	updateData := DonationCompanyMapper(*details)
	updatedDonationCompany, err := s.dal.UpdateOne(ctx, filter, updateData)
	if err != nil {
		return local_util.HandleDBError(err)
	}

	s.kafkaProducer.PublishMessage(ctx, updatedDonationCompany, string(constants.ClientOrchestrationDonationCompanyTopic), string(constants.ClientOrchestrationDonationCompanyTopic), "donation company updated")

	s.logger.Infof("[DonationCompanyStorage][Update] donation company updated successfully")
	return nil
}

func (s *DonationCompanyStorage) FindByID(ctx context.Context, id string) (*donation_company.DonationCompanyListResponse, error) {
	s.logger.Infof("[DonationCompanyStorage][FindByID] fetching donation company by id: %s", id)
	idObj, ok := local_util.StringToObjectID(id)
	if !ok {
		s.logger.Errorf("[DonationCompanyStorage][FindByID] invalid object id")
		return nil, errors.New(localization.ErrorInvalidID.Code)
	}
	filter := bson.M{"_id": idObj, "is_deleted": false}
	projection := bson.M{}

	result, err := s.dal.FindOne(ctx, filter, projection)
	if err != nil {
		s.logger.Errorf("[DonationCompanyStorage][FindByID] failed to find donation company: %v", err)
		code, _ := local_util.HandleMongoError(err)
		return nil, errors.New(code)
	}
	s.logger.Infof("[DonationCompanyStorage][FindByID] donation company retrieved successfully")
	return MapToDonationCompanyListResponse(result), nil
}

func (s *DonationCompanyStorage) FindByAccountNumber(ctx context.Context, accountNumber string) (*donation_model.DonationCompany, error) {
	s.logger.Infof("[DonationCompanyStorage][FindByAccountNumber] searching for donation company by account number")
	filter := bson.M{
		"account_number": accountNumber,
		"is_deleted":     false,
	}

	result, err := s.dal.FindOne(ctx, filter, bson.M{})
	if err != nil {
		s.logger.Errorf("[DonationCompanyStorage][FindByAccountNumber] failed to find donation company: %v", err)
		return nil, local_util.HandleDBError(err)
	}
	s.logger.Infof("[DonationCompanyStorage][FindByAccountNumber] donation company retrieved successfully")
	return result, nil
}

func (s *DonationCompanyStorage) FindAllWithPagination(ctx context.Context, filterParam types.Filter) (*types.PaginatedResponse[[]donation_company.DonationCompanyListResponse], error) {
	// 1. Base filter (only active records)
	filter := bson.M{"is_deleted": false}
	searchKeys := bson.M{}

	// 2. Allowed filterable/searchable fields
	allowedKeys := []string{"company_name", "account_number", "enabled"}
	// 3. Add search (if provided)
	if filterParam.Search != "" {
		searchRegex := bson.M{"$regex": filterParam.Search, "$options": "i"}
		searchKeys["$or"] = []bson.M{{"company_name": searchRegex}, {"account_number": searchRegex}}
	}

	projection := bson.M{}
	filter, skip, limit := lib.FilterBuilder(filterParam, searchKeys, allowedKeys)

	// 5. Fetch data
	data, err := s.dal.FindAllWithPaginationE(ctx, filter, projection, skip, limit)
	if err != nil {
		s.logger.Errorf("[DonationCompanyStorage][FindAllWithPagination] failed to fetch donation companies: %v", err)
		return nil, errors.New(localization.ErrorUnexpectedError.Code)
	}

	// 6. Count total
	total, err := s.dal.TotalCount(ctx, filter)
	if err != nil {
		s.logger.Errorf("[DonationCompanyStorage][FindAllWithPagination] failed to count donation companies: %v", err)
		return nil, errors.New(localization.ErrorUnexpectedError.Code)
	}

	// 7. Build pagination metadata
	meta := local_util.BuildPaginationMeta(total, filterParam.Page, filterParam.PerPage)

	// 8. Map to DTOs
	dtoData := MapToDonationCompanyListResponses(data)
	s.logger.Infof("[DonationCompanyStorage][FindAllWithPagination] retrieved %d donation companies", len(dtoData))

	// 9. Return standard paginated response
	return &types.PaginatedResponse[[]donation_company.DonationCompanyListResponse]{
		Data: dtoData,
		Meta: meta,
	}, nil
}
func (s *DonationCompanyStorage) Delete(ctx context.Context, id string) error {
	s.logger.Infof("[DonationCompanyStorage][Delete] deleting donation company for id: %s", id)
	objID, err := bson.ObjectIDFromHex(id)
	if err != nil {
		s.logger.Errorf("[DonationCompanyStorage][Delete] invalid object id: %v", err)
		return errors.New(localization.ErrorUnexpectedError.Code)
	}
	filter := bson.M{"_id": objID, "is_deleted": false}
	err = s.dal.DeleteOne(ctx, filter)
	if err != nil {
		s.logger.Errorf("[DonationCompanyStorage][Delete] failed to delete donation company: %v", err)
		return local_util.HandleDBError(err)
	}
	s.logger.Infof("[DonationCompanyStorage][Delete] donation company deleted successfully")
	return nil
}
