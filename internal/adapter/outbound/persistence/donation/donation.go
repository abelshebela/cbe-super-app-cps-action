package donation

import (
	"context"
	"fmt"
	"math/rand"
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

func (d *DonationPersistence) generateDonationCode() string {
	const charset = "ABCDEFGHIJKLMNOPQRSTUVWXYZ0123456789"
	const length = 8

	rand.Seed(time.Now().UnixNano())
	result := make([]byte, length)
	for i := range result {
		result[i] = charset[rand.Intn(len(charset))]
	}

	return "DON" + string(result)
}

func (d *DonationPersistence) DonationNameExists(ctx context.Context, categoryName string) (bool, error) {
	filter := bson.M{"category_name": categoryName, "is_deleted": false}
	count, err := d.donationCategoryDal.CountDocuments(ctx, filter)
	if err != nil {
		return false, fmt.Errorf("CATEGORY_LOOKUP_FAILED")
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
	"$set": bson.M{
		"category_name":    donation.CategoryName,
		"last_modified_at": time.Now(),
	},
}


	_, err = d.donationCategoryDal.UpdateOne(ctx, filter, update)
	if err != nil {
		return nil, fmt.Errorf("failed to update donation category: %v", err)
	}

	return &donation, nil
}

func (d *DonationPersistence) DonationCompanyNameExists(ctx context.Context, companyName string) (bool, error) {
	filter := bson.M{"company_name": companyName, "is_deleted": false}
	count, err := d.donationCompanyDal.CountDocuments(ctx, filter)
	if err != nil {
		return false, fmt.Errorf("COMPANY_LOOKUP_FAILED")
	}
	return count > 0, nil
}

func (d *DonationPersistence) DonationCompanyAccountExists(ctx context.Context, accountNumber string) (bool, error) {
	filter := bson.M{"account_number": accountNumber, "is_deleted": false}
	count, err := d.donationCompanyDal.CountDocuments(ctx, filter)
	if err != nil {
		return false, fmt.Errorf("COMPANY_LOOKUP_FAILED")
	}
	return count > 0, nil
}

func (d *DonationPersistence) CreateDonationCompany(ctx context.Context, company dto.DonationCompanyRequest) error {
	companyModel := mappers.ToDonationCompanyModel(&company)
	companyModel.CreatedAt = time.Now()
	companyModel.LastModifiedAt = time.Now()

	_, err := d.donationCompanyDal.InsertOne(ctx, companyModel)
	if err != nil {
		return fmt.Errorf("failed to create donation company: %v", err)
	}
	return nil
}

func (d *DonationPersistence) CreateDonationCompanyWithURL(ctx context.Context, companyName, logoURL, accountNumber string) error {
	companyModel := &model.DonationCompany{
		CompanyName:    companyName,
		CompanyLogo:    logoURL,
		AccountNumber:  accountNumber,
		IsDeleted:      false,
		CreatedAt:      time.Now(),
		LastModifiedAt: time.Now(),
	}

	_, err := d.donationCompanyDal.InsertOne(ctx, companyModel)
	if err != nil {
		return fmt.Errorf("failed to create donation company: %v", err)
	}
	return nil
}

func (d *DonationPersistence) FetchDonationCompany(ctx context.Context, filterParams *constant.MongoFilter) (*common_util.PaginatedResponse[[]*dto.DonationCompanyListResponse], error) {
	filter := bson.M{"is_deleted": false}

	if filterParams.Search != "" {
		searchRegex := bson.M{"$regex": filterParams.Search, "$options": "i"}
		filter["company_name"] = searchRegex
	}

	if filterParams.Filters != nil {
		for key, value := range filterParams.Filters {
			filter[key] = value
		}
	}

	skip := int64((filterParams.Page - 1) * filterParams.PerPage)
	limit := int64(filterParams.PerPage)

	cursor, err := d.donationCompanyDal.Find(ctx, filter, options.Find().SetSkip(skip).SetLimit(limit))
	if err != nil {
		return nil, fmt.Errorf("failed to fetch donation companies: %v", err)
	}
	defer cursor.Close(ctx)

	var companies []*model.DonationCompany
	if err := cursor.All(ctx, &companies); err != nil {
		return nil, fmt.Errorf("failed to decode donation companies: %v", err)
	}

	total, err := d.donationCompanyDal.CountDocuments(ctx, filter)
	if err != nil {
		return nil, fmt.Errorf("failed to count donation companies: %v", err)
	}

	var result []*dto.DonationCompanyListResponse
	for _, company := range companies {
		result = append(result, mappers.ToDonationCompanyListResponse(company))
	}

	meta := common_util.BuildPaginationMeta(total, filterParams.Page, filterParams.PerPage)
	return &common_util.PaginatedResponse[[]*dto.DonationCompanyListResponse]{
		Data: result,
		Meta: meta,
	}, nil
}

func (d *DonationPersistence) FetchDonationCompanyByID(ctx context.Context, id string) (*dto.DonationCompanyListResponse, error) {
	objID, err := common_util.ParsePrimitiveObjectID(id)
	if err != nil {
		return nil, fmt.Errorf("invalid ID format: %v", err)
	}

	filter := bson.M{"_id": objID, "is_deleted": false}
	var company model.DonationCompany
	err = d.donationCompanyDal.FindOne(ctx, filter).Decode(&company)
	if err != nil {
		if err == mongo.ErrNoDocuments {
			return nil, fmt.Errorf("donation company not found")
		}
		return nil, fmt.Errorf("failed to fetch donation company: %v", err)
	}

	return mappers.ToDonationCompanyListResponse(&company), nil
}

func (d *DonationPersistence) UpdateDonationCompany(ctx context.Context, id string, company dto.DonationCompanyRequest) (*dto.DonationCompanyRequest, error) {
	objID, err := common_util.ParsePrimitiveObjectID(id)
	if err != nil {
		return nil, fmt.Errorf("invalid ID format: %v", err)
	}

	filter := bson.M{"_id": objID, "is_deleted": false}
	update := bson.M{
		"$set": bson.M{
		"company_name":     company.CompanyName,
		"account_number":   company.AccountNumber,
		"last_modified_at": time.Now(),
	}}

	_, err = d.donationCompanyDal.UpdateOne(ctx, filter, update)
	if err != nil {
		return nil, fmt.Errorf("failed to update donation company: %v", err)
	}

	return &company, nil
}

func (d *DonationPersistence) UpdateDonationCompanyWithLogoURL(ctx context.Context, id string, company dto.DonationCompanyRequest, logoURL string) (*dto.DonationCompanyRequest, error) {
	objID, err := common_util.ParsePrimitiveObjectID(id)
	if err != nil {
		return nil, fmt.Errorf("invalid ID format: %v", err)
	}

	filter := bson.M{"_id": objID, "is_deleted": false}
	update := bson.M{
		"last_modified_at": time.Now(),
	}

	// Only update fields that are provided
	if company.CompanyName != "" {
		update["company_name"] = company.CompanyName
	}
	if company.AccountNumber != "" {
		update["account_number"] = company.AccountNumber
	}
	if logoURL != "" {
		update["company_logo"] = logoURL
	}

	_, err = d.donationCompanyDal.UpdateOne(ctx, filter, update)
	if err != nil {
		return nil, fmt.Errorf("failed to update donation company: %v", err)
	}

	return &company, nil
}

func (d *DonationPersistence) DonationTitleExists(ctx context.Context, title string) (bool, error) {
	filter := bson.M{"title": title, "is_deleted": false}
	count, err := d.donationDal.CountDocuments(ctx, filter)
	if err != nil {
		return false, fmt.Errorf("DONATION_LOOKUP_FAILED")
	}
	return count > 0, nil
}

func (d *DonationPersistence) DonationCompanyExists(ctx context.Context, companyID string) (bool, error) {
	objID, err := common_util.ParsePrimitiveObjectID(companyID)
	if err != nil {
		return false, fmt.Errorf("INVALID_ID_FORMAT")
	}

	filter := bson.M{"_id": objID, "is_deleted": false}
	count, err := d.donationCompanyDal.CountDocuments(ctx, filter)
	if err != nil {
		return false, fmt.Errorf("COMPANY_LOOKUP_FAILED")
	}
	return count > 0, nil
}

func (d *DonationPersistence) DonationCategoryExists(ctx context.Context, categoryID string) (bool, error) {
	objID, err := common_util.ParsePrimitiveObjectID(categoryID)
	if err != nil {
		return false, fmt.Errorf("INVALID_ID_FORMAT")
	}

	filter := bson.M{"_id": objID, "is_deleted": false}
	count, err := d.donationCategoryDal.CountDocuments(ctx, filter)
	if err != nil {
		return false, fmt.Errorf("CATEGORY_LOOKUP_FAILED")
	}
	return count > 0, nil
}

func (d *DonationPersistence) CreateDonation(ctx context.Context, donation dto.DonationRequest) error {
	companyObjID, err := d.convertStringToObjectID(donation.CompanyID)
	if err != nil {
		return fmt.Errorf("invalid company ID format: %v", err)
	}

	categoryObjID, err := d.convertStringToObjectID(donation.CategoryID)
	if err != nil {
		return fmt.Errorf("invalid category ID format: %v", err)
	}

	donationCode := d.generateDonationCode()

	donationModel := mappers.ToDonationModel(&donation)
	donationModel.DonationCode = donationCode
	donationModel.CompanyID = companyObjID
	donationModel.CategoryID = categoryObjID
	donationModel.CreatedAt = time.Now()
	donationModel.LastModifiedAt = time.Now()

	_, err = d.donationDal.InsertOne(ctx, donationModel)
	if err != nil {
		return fmt.Errorf("failed to create donation: %v", err)
	}
	return nil
}

func (d *DonationPersistence) CreateDonationWithURLs(ctx context.Context, donation dto.DonationCPSRequest) error {
	startDate, err := time.Parse("2006-01-02T15:04:05Z07:00", donation.StartDate)
	if err != nil {
		return fmt.Errorf("failed to parse start date: %v", err)
	}

	endDate, err := time.Parse("2006-01-02T15:04:05Z07:00", donation.EndDate)
	if err != nil {
		return fmt.Errorf("failed to parse end date: %v", err)
	}

	companyObjID, err := d.convertStringToObjectID(donation.CompanyID)
	if err != nil {
		return fmt.Errorf("invalid company ID format: %v", err)
	}

	categoryObjID, err := d.convertStringToObjectID(donation.CategoryID)
	if err != nil {
		return fmt.Errorf("invalid category ID format: %v", err)
	}

	donationCode := d.generateDonationCode()

	donationModel := &model.Donation{
		DonationCode:        donationCode,
		CompanyID:           companyObjID,
		CategoryID:          categoryObjID,
		Title:               donation.Title,
		IsFeatured:          donation.IsFeatured,
		Target:              donation.Target,
		DonationDescription: donation.DonationDescription,
		DonationImages:      donation.DonationImages,
		EndDate:             endDate,
		StartDate:           startDate,
		IsDeleted:           false,
		CreatedAt:           time.Now(),
		LastModifiedAt:      time.Now(),
	}

	_, err = d.donationDal.InsertOne(ctx, donationModel)
	if err != nil {
		return fmt.Errorf("failed to create donation: %v", err)
	}
	return nil
}

func (d *DonationPersistence) FetchDonation(ctx context.Context, filterParams *constant.MongoFilter) (*common_util.PaginatedResponse[[]*dto.DonationListResponse], error) {
	filter := bson.M{"is_deleted": false}

	if filterParams.Search != "" {
		searchRegex := bson.M{"$regex": filterParams.Search, "$options": "i"}
		filter["title"] = searchRegex
	}

	if filterParams.Filters != nil {
		for key, value := range filterParams.Filters {
			// Handle ObjectID conversion for foreign keys
			if key == "company_id" || key == "category_id" {
				if strValue, ok := value.(string); ok {
					objID, err := common_util.ParsePrimitiveObjectID(strValue)
					if err != nil {
						return nil, fmt.Errorf("invalid %s format: %v", key, err)
					}
					filter[key] = objID
				} else {
					filter[key] = value
				}
			} else {
				filter[key] = value
			}
		}
	}

	skip := int64((filterParams.Page - 1) * filterParams.PerPage)
	limit := int64(filterParams.PerPage)

	cursor, err := d.donationDal.Find(ctx, filter, options.Find().SetSkip(skip).SetLimit(limit))
	if err != nil {
		return nil, fmt.Errorf("failed to fetch donations: %v", err)
	}
	defer cursor.Close(ctx)

	var donations []*model.Donation
	if err := cursor.All(ctx, &donations); err != nil {
		return nil, fmt.Errorf("failed to decode donations: %v", err)
	}

	// Get total count
	total, err := d.donationDal.CountDocuments(ctx, filter)
	if err != nil {
		return nil, fmt.Errorf("failed to count donations: %v", err)
	}

	var result []*dto.DonationListResponse
	for _, donation := range donations {
		// Fetch related Company using ObjectID
		var company model.DonationCompany
		companyFilter := bson.M{"_id": donation.CompanyID}
		err = d.donationCompanyDal.FindOne(ctx, companyFilter).Decode(&company)
		if err != nil {
			return nil, fmt.Errorf("failed to fetch company: %v", err)
		}

		// Fetch related Category using ObjectID
		var category model.DonationCategory
		categoryFilter := bson.M{"_id": donation.CategoryID}
		err = d.donationCategoryDal.FindOne(ctx, categoryFilter).Decode(&category)
		if err != nil {
			return nil, fmt.Errorf("failed to fetch category: %v", err)
		}

		result = append(result, mappers.ToDonationListResponse(donation, &company, &category))
	}

	meta := common_util.BuildPaginationMeta(total, filterParams.Page, filterParams.PerPage)
	return &common_util.PaginatedResponse[[]*dto.DonationListResponse]{
		Data: result,
		Meta: meta,
	}, nil
}

func (d *DonationPersistence) FetchDonationByID(ctx context.Context, id string) (*dto.DonationListResponse, error) {
	objID, err := common_util.ParsePrimitiveObjectID(id)
	if err != nil {
		return nil, fmt.Errorf("invalid ID format: %v", err)
	}

	filter := bson.M{"_id": objID, "is_deleted": false}
	var donation model.Donation
	err = d.donationDal.FindOne(ctx, filter).Decode(&donation)
	if err != nil {
		if err == mongo.ErrNoDocuments {
			return nil, fmt.Errorf("donation not found")
		}
		return nil, fmt.Errorf("failed to fetch donation: %v", err)
	}

	// Fetch related Company using ObjectID
	var company model.DonationCompany
	companyFilter := bson.M{"_id": donation.CompanyID}
	err = d.donationCompanyDal.FindOne(ctx, companyFilter).Decode(&company)
	if err != nil {
		return nil, fmt.Errorf("failed to fetch company: %v", err)
	}

	// Fetch related Category using ObjectID
	var category model.DonationCategory
	categoryFilter := bson.M{"_id": donation.CategoryID}
	err = d.donationCategoryDal.FindOne(ctx, categoryFilter).Decode(&category)
	if err != nil {
		return nil, fmt.Errorf("failed to fetch category: %v", err)
	}

	return mappers.ToDonationListResponse(&donation, &company, &category), nil
}

func (d *DonationPersistence) FetchDonationByCode(ctx context.Context, donationCode string) (*dto.DonationListResponse, error) {
	filter := bson.M{"donation_code": donationCode, "is_deleted": false}
	var donation model.Donation
	err := d.donationDal.FindOne(ctx, filter).Decode(&donation)
	if err != nil {
		if err == mongo.ErrNoDocuments {
			return nil, fmt.Errorf("donation not found")
		}
		return nil, fmt.Errorf("failed to fetch donation: %v", err)
	}

	// Fetch related Company using ObjectID
	var company model.DonationCompany
	companyFilter := bson.M{"_id": donation.CompanyID}
	err = d.donationCompanyDal.FindOne(ctx, companyFilter).Decode(&company)
	if err != nil {
		return nil, fmt.Errorf("failed to fetch company: %v", err)
	}

	// Fetch related Category using ObjectID
	var category model.DonationCategory
	categoryFilter := bson.M{"_id": donation.CategoryID}
	err = d.donationCategoryDal.FindOne(ctx, categoryFilter).Decode(&category)
	if err != nil {
		return nil, fmt.Errorf("failed to fetch category: %v", err)
	}

	return mappers.ToDonationListResponse(&donation, &company, &category), nil
}

func (d *DonationPersistence) UpdateDonation(ctx context.Context, id string, donation dto.DonationRequest) (*dto.DonationRequest, error) {
	objID, err := common_util.ParsePrimitiveObjectID(id)
	if err != nil {
		return nil, fmt.Errorf("INVALID_ID_FORMAT")
	}

	filter := bson.M{"_id": objID, "is_deleted": false}

	// Build update document with $set operator
	setUpdate := bson.M{
		"last_modified_at": time.Now(),
	}

	// Only update fields that are provided and not empty
	if donation.CompanyID != "" {
		// Convert string ID to ObjectID
		companyObjID, err := common_util.ParsePrimitiveObjectID(donation.CompanyID)
		if err != nil {
			return nil, fmt.Errorf("INVALID_COMPANY_ID_FORMAT")
		}
		setUpdate["company_id"] = companyObjID
	}

	if donation.CategoryID != "" {
		// Convert string ID to ObjectID
		categoryObjID, err := common_util.ParsePrimitiveObjectID(donation.CategoryID)
		if err != nil {
			return nil, fmt.Errorf("INVALID_CATEGORY_ID_FORMAT")
		}
		setUpdate["category_id"] = categoryObjID
	}

	if donation.Title != "" {
		setUpdate["title"] = donation.Title
	}

	// Always update is_featured as it's a boolean
	setUpdate["is_featured"] = donation.IsFeatured

	if donation.Target > 0 {
		setUpdate["target"] = donation.Target
	}

	if donation.DonationDescription != "" {
		setUpdate["donation_description"] = donation.DonationDescription
	}

	if !donation.EndDate.IsZero() {
		setUpdate["end_date"] = donation.EndDate
	}

	if !donation.StartDate.IsZero() {
		setUpdate["start_date"] = donation.StartDate
	}

	// Use $set operator for proper MongoDB update syntax
	update := bson.M{"$set": setUpdate}

	_, err = d.donationDal.UpdateOne(ctx, filter, update)
	if err != nil {
		return nil, fmt.Errorf("FAILED_TO_UPDATE_DONATION")
	}

	return &donation, nil
}

func (d *DonationPersistence) UpdateDonationWithImageURLs(ctx context.Context, id string, donation dto.DonationRequest, imageURLs []string) (*dto.DonationRequest, error) {
	objID, err := common_util.ParsePrimitiveObjectID(id)
	if err != nil {
		return nil, fmt.Errorf("INVALID_ID_FORMAT")
	}

	filter := bson.M{"_id": objID, "is_deleted": false}

	// Build update document with $set operator
	setUpdate := bson.M{
		"last_modified_at": time.Now(),
	}

	// Only update fields that are provided and not empty
	if donation.CompanyID != "" {
		// Convert string ID to ObjectID
		companyObjID, err := common_util.ParsePrimitiveObjectID(donation.CompanyID)
		if err != nil {
			return nil, fmt.Errorf("INVALID_COMPANY_ID_FORMAT")
		}
		setUpdate["company_id"] = companyObjID
	}

	if donation.CategoryID != "" {
		// Convert string ID to ObjectID
		categoryObjID, err := common_util.ParsePrimitiveObjectID(donation.CategoryID)
		if err != nil {
			return nil, fmt.Errorf("INVALID_CATEGORY_ID_FORMAT")
		}
		setUpdate["category_id"] = categoryObjID
	}

	if donation.Title != "" {
		setUpdate["title"] = donation.Title
	}

	// Always update is_featured as it's a boolean
	setUpdate["is_featured"] = donation.IsFeatured

	if donation.Target > 0 {
		setUpdate["target"] = donation.Target
	}

	if donation.DonationDescription != "" {
		setUpdate["donation_description"] = donation.DonationDescription
	}

	if !donation.EndDate.IsZero() {
		setUpdate["end_date"] = donation.EndDate
	}

	if !donation.StartDate.IsZero() {
		setUpdate["start_date"] = donation.StartDate
	}

	if len(imageURLs) > 0 {
		setUpdate["donation_images"] = imageURLs
	}

	// Use $set operator for proper MongoDB update syntax
	update := bson.M{"$set": setUpdate}

	_, err = d.donationDal.UpdateOne(ctx, filter, update)
	if err != nil {
		return nil, fmt.Errorf("FAILED_TO_UPDATE_DONATION")
	}

	return &donation, nil
}

// Helper function to convert string IDs to ObjectIDs
func (d *DonationPersistence) convertStringToObjectID(id string) (bson.ObjectID, error) {
	return common_util.ParsePrimitiveObjectID(id)
}

// Helper function to safely convert string to ObjectID with fallback
func (d *DonationPersistence) safeConvertToObjectID(value interface{}) (interface{}, error) {
	if strValue, ok := value.(string); ok {
		objID, err := common_util.ParsePrimitiveObjectID(strValue)
		if err != nil {
			return nil, fmt.Errorf("invalid ObjectID format: %v", err)
		}
		return objID, nil
	}
	return value, nil
}
