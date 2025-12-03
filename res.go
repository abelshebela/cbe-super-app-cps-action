package donation

import (
	donation_dto "cbe-super-app-cps-action/internal/constants/dto/donation"
	"cbe-super-app-cps-action/internal/constants/lib"
	"cbe-super-app-cps-action/internal/constants/localization"
	"cbe-super-app-cps-action/internal/constants/model"
	"cbe-super-app-cps-action/internal/constants/types"
	"cbe-super-app-cps-action/internal/storage"
	"context"
	"errors"
	"fmt"
	"strconv"

	local_util "cbe-super-app-cps-action/pkgs/utils"

	"gitlab.com/bersufekadgetachew/cbe-super-app-shared/shared/dal"
	"gitlab.com/bersufekadgetachew/cbe-super-app-shared/shared/utils"
	"go.mongodb.org/mongo-driver/v2/bson"
	"go.mongodb.org/mongo-driver/v2/mongo"
	"go.mongodb.org/mongo-driver/v2/mongo/options"
)

type DonationStorage struct {
	dal                 dal.MongoDal[model.Donation, model.Donation]
	donationCompanyDal  dal.MongoDal[model.DonationCompany, model.DonationCompany]
	donationCategoryDal dal.MongoDal[model.DonationCategory, model.DonationCategory]
	client              *mongo.Client
	logger              utils.Logger
	dbName              string
}

func NewDonationRepository(client *mongo.Client, dbName string, collection string, logger utils.Logger) storage.DonationRepository {
	return &DonationStorage{
		dal:                 dal.NewMongoDal[model.Donation, model.Donation](client, dbName, collection),
		donationCompanyDal:  dal.NewMongoDal[model.DonationCompany, model.DonationCompany](client, dbName, "donation_companies"),
		donationCategoryDal: dal.NewMongoDal[model.DonationCategory, model.DonationCategory](client, dbName, "donation_categories"),
		client:              client,
		logger:              logger,
		dbName:              dbName,
	}
}

func (d *DonationStorage) Create(ctx context.Context, donation *model.Donation) error {
	_, err := d.dal.InsertOne(ctx, *donation)
	if err != nil {
		return errors.New(localization.ErrorUnexpectedError.Code)
	}
	return nil
}

func (d *DonationStorage) Update(ctx context.Context, id string, donation *model.Donation) error {
	objID, err := bson.ObjectIDFromHex(id)
	if err != nil {
		return errors.New(localization.ErrorUnexpectedError.Code)
	}
	filter := bson.M{"_id": objID, "is_deleted": false}
	updatedData := DonationMapper(*donation)

	_, err = d.dal.UpdateOne(ctx, filter, updatedData)
	if err != nil {
		if err == mongo.ErrNoDocuments {
			return errors.New(localization.ErrorFileNotFound.Code)
		}
		return errors.New(localization.ErrorUnexpectedError.Code)
	}
	return nil
}

func (d *DonationStorage) Delete(ctx context.Context, id string) error {
	objID, err := bson.ObjectIDFromHex(id)
	if err != nil {
		return errors.New(localization.ErrorUnexpectedError.Code)
	}
	filter := bson.M{"_id": objID, "is_deleted": false}
	return d.dal.DeleteOne(ctx, filter)
}

func (d *DonationStorage) FindByID(ctx context.Context, id string) (*donation_dto.DonationListResponse, error) {
	objID, err := bson.ObjectIDFromHex(id)
	if err != nil {
		return nil, errors.New(localization.ErrorUnexpectedError.Code)
	}
	filter := bson.M{"_id": objID, "is_deleted": false}

	result, err := d.dal.FindOne(ctx, filter, nil)
	if err != nil {
		if err == mongo.ErrNoDocuments {
			return nil, errors.New(localization.ErrorFileNotFound.Code)
		}
		return nil, err
	}

	companyFilter := bson.M{"_id": result.CompanyID}
	company, err := d.donationCompanyDal.FindOne(ctx, companyFilter, nil)
	if err != nil {
		return nil, fmt.Errorf("failed to fetch company: %v", err)
	}

	categoryFilter := bson.M{"_id": result.CategoryID}
	category, err := d.donationCategoryDal.FindOne(ctx, categoryFilter, nil)
	if err != nil {
		return nil, fmt.Errorf("failed to fetch category: %v", err)
	}

	return MapToDonationListResponse(result, company, category), nil
}

func (d *DonationStorage) FindAllWithPagination(ctx context.Context, filterParam types.Filter) (*types.PaginatedResponse[[]donation_dto.DonationListResponse], error) {

	searchKeys := bson.M{}
	allowedKeys := []string{"title", "is_featured", "enabled", "donation_code", "target", "end_date"}

	if filterParam.Search != "" {
		orFilters, err := d.buildDonationSearchFilters(ctx, filterParam.Search)
		if err != nil {
			return nil, err
		}
		if len(orFilters) > 0 {
			searchKeys["$or"] = orFilters
		}
	}
	filter, skip, limit := lib.FilterBuilder(filterParam, searchKeys, allowedKeys)
	filter["is_deleted"] = false
	data, err := d.dal.FindAllWithPagination(ctx, filter, bson.M{}, skip, limit)
	if err != nil {
		return nil, err
	}

	total, err := d.dal.TotalCount(ctx, filter)
	if err != nil {
		return nil, err
	}

	meta := local_util.BuildPaginationMeta(total, filterParam.Page, filterParam.PerPage)

	var result []donation_dto.DonationListResponse
	for _, donation := range data {

		companyFilter := bson.M{"_id": donation.CompanyID}
		company, err := d.donationCompanyDal.FindOne(ctx, companyFilter, nil)
		if err != nil {
			return nil, fmt.Errorf("failed to fetch company: %v", err)
		}

		categoryFilter := bson.M{"_id": donation.CategoryID}
		category, err := d.donationCategoryDal.FindOne(ctx, categoryFilter, nil)
		if err != nil {
			return nil, fmt.Errorf("failed to fetch category: %v", err)
		}

		result = append(result, *MapToDonationListResponse(donation, company, category))
	}

	return &types.PaginatedResponse[[]donation_dto.DonationListResponse]{
		Data: result,
		Meta: meta,
	}, nil
}

func (d *DonationStorage) buildDonationSearchFilters(ctx context.Context, term string) ([]bson.M, error) {
	regex := bson.M{"$regex": term, "$options": "i"}
	orFilters := []bson.M{
		{"title": regex},
		{"donation_description": regex},
		{"donation_code": regex},
	}

	if amount, err := strconv.ParseInt(term, 10, 64); err == nil {
		orFilters = append(orFilters,
			bson.M{"target": amount},
			bson.M{"current_amount": amount},
		)
	}

	if objID, err := bson.ObjectIDFromHex(term); err == nil {
		orFilters = append(orFilters, bson.M{"_id": objID})
	}

	companyIDs, err := d.lookupCompanyIDsBySearch(ctx, regex, term)
	if err != nil {
		return nil, err
	}
	if len(companyIDs) > 0 {
		orFilters = append(orFilters, bson.M{"company_id": bson.M{"$in": companyIDs}})
	}

	categoryIDs, err := d.lookupCategoryIDsBySearch(ctx, regex, term)
	if err != nil {
		return nil, err
	}
	if len(categoryIDs) > 0 {
		orFilters = append(orFilters, bson.M{"category_id": bson.M{"$in": categoryIDs}})
	}

	return orFilters, nil
}

func (d *DonationStorage) lookupCompanyIDsBySearch(ctx context.Context, regex bson.M, raw string) ([]bson.ObjectID, error) {
	conditions := []bson.M{
		{"company_name": regex},
		{"phone_number": regex},
		{"account_number": regex},
		{"email": regex},
	}

	if objID, err := bson.ObjectIDFromHex(raw); err == nil {
		conditions = append(conditions, bson.M{"_id": objID})
	}

	filter := bson.M{
		"is_deleted": false,
		"$or":        conditions,
	}

	return d.lookupIDs(ctx, "donation_companies", filter)
}

func (d *DonationStorage) lookupCategoryIDsBySearch(ctx context.Context, regex bson.M, raw string) ([]bson.ObjectID, error) {
	conditions := []bson.M{
		{"category_name": regex},
	}

	if objID, err := bson.ObjectIDFromHex(raw); err == nil {
		conditions = append(conditions, bson.M{"_id": objID})
	}

	filter := bson.M{
		"is_deleted": false,
		"$or":        conditions,
	}

	return d.lookupIDs(ctx, "donation_categories", filter)
}

func (d *DonationStorage) lookupIDs(ctx context.Context, collection string, filter bson.M) ([]bson.ObjectID, error) {
	coll := d.client.Database(d.dbName).Collection(collection)
	cursor, err := coll.Find(ctx, filter, options.Find().SetProjection(bson.M{"_id": 1}))
	if err != nil {
		return nil, err
	}
	defer cursor.Close(ctx)

	var ids []bson.ObjectID
	for cursor.Next(ctx) {
		var doc struct {
			ID bson.ObjectID `bson:"_id"`
		}
		if err := cursor.Decode(&doc); err != nil {
			return nil, err
		}
		ids = append(ids, doc.ID)
	}

	if err := cursor.Err(); err != nil {
		return nil, err
	}

	return ids, nil
}
