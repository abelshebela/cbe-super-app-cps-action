package donation_company

import (
	"cbe-super-app-cps-action/internal/constants/dto/donation_company"
	"cbe-super-app-cps-action/internal/constants/lib"
	"cbe-super-app-cps-action/internal/constants/localization"
	"cbe-super-app-cps-action/internal/constants/model"
	"cbe-super-app-cps-action/internal/constants/types"
	"cbe-super-app-cps-action/internal/storage"
	local_util "cbe-super-app-cps-action/pkgs/utils"
	"context"
	"errors"
	"fmt"

	"gitlab.com/bersufekadgetachew/cbe-super-app-shared/shared/dal"
	"gitlab.com/bersufekadgetachew/cbe-super-app-shared/shared/utils"
	"go.mongodb.org/mongo-driver/v2/bson"
	"go.mongodb.org/mongo-driver/v2/mongo"
)

type DonationCompanyStorage struct {
	dal    dal.MongoDal[model.DonationCompany, model.DonationCompany]
	client *mongo.Client
	logger utils.Logger
}

func NewDonationCompanyRepository(client *mongo.Client, dbName string, collection string, logger utils.Logger) storage.DonationCompanyRepository {
	return &DonationCompanyStorage{
		dal:    dal.NewMongoDal[model.DonationCompany, model.DonationCompany](client, dbName, collection),
		client: client,
		logger: logger,
	}
}

func (s *DonationCompanyStorage) Create(ctx context.Context, details *model.DonationCompany) error {
	fmt.Println("////////////////CREATE")
	_, err := s.dal.InsertOne(ctx, *details)
	if err != nil {
		return errors.New(localization.ErrorUnexpectedError.Code)
	}
	return nil
}

func (s *DonationCompanyStorage) Update(ctx context.Context, id string, details *model.DonationCompany) error {
	objID, err := bson.ObjectIDFromHex(id)
	if err != nil {
		return errors.New(localization.ErrorUnexpectedError.Code)
	}
	filter := bson.M{"_id": objID, "is_deleted": false}
	updateData := DonationCompanyMapper(*details)
	_, err = s.dal.UpdateOne(ctx, filter, updateData)
	if err != nil {

		if err == mongo.ErrNoDocuments {
			return errors.New(localization.ErrorFileNotFound.Code)
		}
		return errors.New(localization.ErrorUnexpectedError.Code)
	}
	return nil
}
func (s *DonationCompanyStorage) FindByID(ctx context.Context, id string) (*donation_company.DonationCompanyListResponse, error) {
	idObj, ok := local_util.StringToObjectID(id)
	if !ok {
		s.logger.Errorf("Invalid ObjectID for fetch by id: %s", id)
		return nil, errors.New(localization.ErrorInvalidID.Code)
	}
	filter := bson.M{"_id": idObj, "is_deleted": false}
	projection := bson.M{}

	result, err := s.dal.FindOne(ctx, filter, projection)
	if err != nil {
		s.logger.Errorf("Error finding donation company: %v", err)
		code, _ := local_util.HandleMongoError(err)
		return nil, errors.New(code)
	}
	s.logger.Infof("Successfully found donation company: %+v", result)
	return MapToDonationCompanyListResponse(result), nil
}
func (s *DonationCompanyStorage) FindAllWithPagination(ctx context.Context, filterParam types.Filter) (*types.PaginatedResponse[[]donation_company.DonationCompanyListResponse], error) {
	// 1. Base filter (only active records)
	filter := bson.M{"is_deleted": false}
	searchKeys := bson.M{}

	// 2. Allowed filterable/searchable fields
	allowedKeys := []string{"company_name", "account_number"}
	// 3. Add search (if provided)
	if filterParam.Search != "" {
		searchRegex := bson.M{"$regex": filterParam.Search, "$options": "i"}
		searchKeys["$or"] = []bson.M{{"company_name": searchRegex}, {"account_number": searchRegex}}
	}

	projection := bson.M{}
	filter, skip, limit := lib.FilterBuilder(filterParam, searchKeys, allowedKeys)

	// 5. Fetch data
	data, err := s.dal.FindAllWithPagination(ctx, filter, projection, skip, limit)
	if err != nil {
		return nil, errors.New(localization.ErrorUnexpectedError.Message)
	}

	// 6. Count total
	total, err := s.dal.TotalCount(ctx, filter)
	if err != nil {
		return nil, errors.New(localization.ErrorUnexpectedError.Message)
	}

	// 7. Build pagination metadata
	meta := local_util.BuildPaginationMeta(total, filterParam.Page, filterParam.PerPage)

	// 8. Map to DTOs
	dtoData := MapToDonationCompanyListResponses(data)

	// 9. Return standard paginated response
	return &types.PaginatedResponse[[]donation_company.DonationCompanyListResponse]{
		Data: dtoData,
		Meta: meta,
	}, nil
}
func (s *DonationCompanyStorage) Delete(ctx context.Context, id string) error {
	objID, err := bson.ObjectIDFromHex(id)
	if err != nil {
		return errors.New(localization.ErrorUnexpectedError.Code)
	}
	filter := bson.M{"_id": objID, "is_deleted": false}
	return s.dal.DeleteOne(ctx, filter)
}
