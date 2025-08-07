package donation

import (
	"context"
	"fmt"
	"time"

	"gitlab.com/bersufekadgetachew/cbe-super-app-shared/shared/utils"
	"go.mongodb.org/mongo-driver/v2/bson"
	"go.mongodb.org/mongo-driver/v2/mongo"
	"go.mongodb.org/mongo-driver/v2/mongo/options"

	"github.com/CBE-Super-App/cbe-super-app-cps-action/internal/adapter/outbound/mappers"
	"github.com/CBE-Super-App/cbe-super-app-cps-action/internal/adapter/outbound/model"
	dto "github.com/CBE-Super-App/cbe-super-app-cps-action/internal/application/dto"
	"github.com/CBE-Super-App/cbe-super-app-cps-action/internal/domain/donation"
	common_util "github.com/CBE-Super-App/cbe-super-app-cps-action/pkgs/utils"
	constant "github.com/CBE-Super-App/cbe-super-app-cps-action/utils"
)

type DonationPersistence struct {
	client              *mongo.Client
	donationDal         *mongo.Collection
	donationCategoryDal *mongo.Collection
	donationCompanyDal  *mongo.Collection
	logger              utils.Logger
}

func InitDonationPersistence(client *mongo.Client, dbName string, collections []string, logger utils.Logger) donation.DonationRepository {
	return &DonationPersistence{
		client:              client,
		donationDal:         client.Database(dbName).Collection(collections[0]),
		donationCategoryDal: client.Database(dbName).Collection(collections[1]),
		donationCompanyDal:  client.Database(dbName).Collection(collections[2]),
		logger:              logger,
	}
}

func (d *DonationPersistence) DonationNameExists(ctx context.Context, categoryName string) (bool, error) {
	filter := bson.M{"category_name": categoryName, "is_deleted": false}
	count, err := d.donationCategoryDal.CountDocuments(ctx, filter)
	if err != nil {
		return false, fmt.Errorf("failed to check donation category name: %v", err)
	}
	return count > 0, nil
}

func (d *DonationPersistence) CreateDonationCategory(ctx context.Context, donation dto.DonationCategoryRequest) error {
	donationModel := mappers.ToDonationCategoryModel(&donation)
	donationModel.CreatedAt = time.Now()
	donationModel.LastModifiedAt = time.Now()

	_, err := d.donationCategoryDal.InsertOne(ctx, donationModel)
	if err != nil {
		return fmt.Errorf("failed to create donation category: %v", err)
	}
	return nil
}

func (d *DonationPersistence) CreateDonationCategoryWithURL(ctx context.Context, categoryName, iconURL string) error {
	donationModel := &model.DonationCategory{
		CategoryName:   categoryName,
		Icon:           iconURL,
		IsDeleted:      false,
		CreatedAt:      time.Now(),
		LastModifiedAt: time.Now(),
	}

	_, err := d.donationCategoryDal.InsertOne(ctx, donationModel)
	if err != nil {
		return fmt.Errorf("failed to create donation category: %v", err)
	}
	return nil
}

func (d *DonationPersistence) FetchDonationCategory(ctx context.Context, filterParams *constant.MongoFilter) (*common_util.PaginatedResponse[[]*dto.DonationCategoryListResponse], error) {
	filter := bson.M{"is_deleted": false}

	if filterParams.Search != "" {
		searchRegex := bson.M{"$regex": filterParams.Search, "$options": "i"}
		filter["category_name"] = searchRegex
	}

	if filterParams.Filters != nil {
		for key, value := range filterParams.Filters {
			filter[key] = value
		}
	}

	skip := int64((filterParams.Page - 1) * filterParams.PerPage)
	limit := int64(filterParams.PerPage)

	cursor, err := d.donationCategoryDal.Find(ctx, filter, options.Find().SetSkip(skip).SetLimit(limit))
	if err != nil {
		return nil, fmt.Errorf("failed to fetch donation categories: %v", err)
	}
	defer cursor.Close(ctx)

	var categories []*model.DonationCategory
	if err := cursor.All(ctx, &categories); err != nil {
		return nil, fmt.Errorf("failed to decode donation categories: %v", err)
	}

	total, err := d.donationCategoryDal.CountDocuments(ctx, filter)
	if err != nil {
		return nil, fmt.Errorf("failed to count donation categories: %v", err)
	}

	var result []*dto.DonationCategoryListResponse
	for _, category := range categories {
		result = append(result, mappers.ToDonationCategoryListResponse(category))
	}

	meta := common_util.BuildPaginationMeta(total, filterParams.Page, filterParams.PerPage)
	return &common_util.PaginatedResponse[[]*dto.DonationCategoryListResponse]{
		Data: result,
		Meta: meta,
	}, nil
}

func (d *DonationPersistence) FetchDonationCategoryByID(ctx context.Context, id string) (*dto.DonationCategoryListResponse, error) {
	objID, err := common_util.ParsePrimitiveObjectID(id)
	if err != nil {
		return nil, fmt.Errorf("invalid ID format: %v", err)
	}

	filter := bson.M{"_id": objID, "is_deleted": false}
	var category model.DonationCategory
	err = d.donationCategoryDal.FindOne(ctx, filter).Decode(&category)
	if err != nil {
		if err == mongo.ErrNoDocuments {
			return nil, fmt.Errorf("donation category not found")
		}
		return nil, fmt.Errorf("failed to fetch donation category: %v", err)
	}

	return mappers.ToDonationCategoryListResponse(&category), nil
}

func (d *DonationPersistence) UpdateDonationCategory(ctx context.Context, id string, donation dto.DonationCategoryRequest) (*dto.DonationCategoryRequest, error) {
	objID, err := common_util.ParsePrimitiveObjectID(id)
	if err != nil {
		return nil, fmt.Errorf("invalid ID format: %v", err)
	}

	filter := bson.M{"_id": objID, "is_deleted": false}
	update := bson.M{
		"category_name":    donation.CategoryName,
		"last_modified_at": time.Now(),
	}

	_, err = d.donationCategoryDal.UpdateOne(ctx, filter, update)
	if err != nil {
		return nil, fmt.Errorf("failed to update donation category: %v", err)
	}

	return &donation, nil
}
