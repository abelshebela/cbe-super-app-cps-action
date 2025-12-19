package donation

import (
	"cbe-super-app-cps-action/internal/constants"
	donation_dto "cbe-super-app-cps-action/internal/constants/dto/donation"
	"cbe-super-app-cps-action/internal/constants/localization"
	"cbe-super-app-cps-action/internal/constants/types"
	"cbe-super-app-cps-action/internal/storage"
	"cbe-super-app-cps-action/internal/storage/kafka"
	"context"
	"errors"

	"gitlab.com/bersufekadgetachew/cbe-super-app-shared/shared/entities/model"

	local_util "cbe-super-app-cps-action/pkgs/utils"

	"cbe-super-app-cps-action/internal/constants/lib"

	"gitlab.com/bersufekadgetachew/cbe-super-app-shared/shared/dal"
	"gitlab.com/bersufekadgetachew/cbe-super-app-shared/shared/utils"
	"go.mongodb.org/mongo-driver/v2/bson"
	"go.mongodb.org/mongo-driver/v2/mongo"
)

type DonationStorage struct {
	dal                 dal.MongoDal[model.Donation, model.Donation]
	donationCompanyDal  dal.MongoDal[model.DonationCompany, model.DonationCompany]
	donationCategoryDal dal.MongoDal[model.DonationCategory, model.DonationCategory]
	client              *mongo.Client
	kafkaProducer       kafka.ClientOrchestrationProducer
	logger              utils.Logger
}

func NewDonationRepository(client *mongo.Client, dbName string, collection string, kafkaProducer kafka.ClientOrchestrationProducer, logger utils.Logger) storage.DonationRepository {
	return &DonationStorage{
		dal:                 dal.NewMongoDal[model.Donation, model.Donation](client, dbName, collection),
		donationCompanyDal:  dal.NewMongoDal[model.DonationCompany, model.DonationCompany](client, dbName, "donation_companies"),
		donationCategoryDal: dal.NewMongoDal[model.DonationCategory, model.DonationCategory](client, dbName, "donation_categories"),
		client:              client,
		kafkaProducer:       kafkaProducer,
		logger:              logger,
	}
}

func (d *DonationStorage) Create(ctx context.Context, donation *model.Donation) error {
	d.logger.Infof("[Create] creating donation")
	newDonation, err := d.dal.InsertOne(ctx, *donation)
	if err != nil {
		d.logger.Errorf("[Create] failed to create donation: %v", err)
		return errors.New(localization.ErrorUnexpectedError.Code)
	}
	d.kafkaProducer.PublishMessage(ctx, newDonation, string(constants.ClientOrchestrationDonationTopic), string(constants.ClientOrchestrationDonationTopic), "new donation created")

	d.logger.Infof("[Create] donation created successfully")
	return nil
}

func (d *DonationStorage) Update(ctx context.Context, id string, donation *model.Donation) error {
	d.logger.Infof("[Update] updating donation for id: %s", id)
	objID, err := bson.ObjectIDFromHex(id)
	if err != nil {
		d.logger.Errorf("[Update] invalid object id: %v", err)
		return errors.New(localization.ErrorUnexpectedError.Code)
	}
	filter := bson.M{"_id": objID, "is_deleted": false}
	updatedData := DonationMapper(*donation)

	updatedDonation, err := d.dal.UpdateOne(ctx, filter, updatedData)
	if err != nil {
		if err == mongo.ErrNoDocuments {
			d.logger.Errorf("[Update] donation not found")
			return errors.New(localization.ErrorFileNotFound.Code)
		}
		d.logger.Errorf("[Update] failed to update donation: %v", err)
		return errors.New(localization.ErrorUnexpectedError.Code)
	}

	d.kafkaProducer.PublishMessage(ctx, updatedDonation, string(constants.ClientOrchestrationDonationTopic), string(constants.ClientOrchestrationDonationTopic), "donation updated")

	d.logger.Infof("[Update] donation updated successfully")
	return nil
}

func (d *DonationStorage) Delete(ctx context.Context, id string) error {
	d.logger.Infof("[Delete] deleting donation for id: %s", id)
	objID, err := bson.ObjectIDFromHex(id)
	if err != nil {
		d.logger.Errorf("[Delete] invalid object id: %v", err)
		return errors.New(localization.ErrorUnexpectedError.Code)
	}
	filter := bson.M{"_id": objID, "is_deleted": false}
	err = d.dal.DeleteOne(ctx, filter)
	if err != nil {
		d.logger.Errorf("[Delete] failed to delete donation: %v", err)
		return err
	}
	d.logger.Infof("[Delete] donation deleted successfully")
	return nil
}

func (d *DonationStorage) FindByID(ctx context.Context, id string) (*donation_dto.DonationListResponse, error) {
	d.logger.Infof("[FindByID] fetching donation by id: %s", id)
	objID, err := bson.ObjectIDFromHex(id)
	if err != nil {
		d.logger.Errorf("[FindByID] invalid object id: %v", err)
		return nil, errors.New(localization.ErrorUnexpectedError.Code)
	}
	filter := bson.M{"_id": objID, "is_deleted": false}

	result, err := d.dal.FindOne(ctx, filter, nil)
	if err != nil {
		d.logger.Errorf("[FindByID] failed to find donation: %v", err)
		return nil, err
	}

	camObj, err := bson.ObjectIDFromHex(result.CompanyID.Hex())
	if err != nil {
		d.logger.Errorf("[FindByID] failed to parse company id: %v", err)
		return nil, errors.New(localization.ErrorUnexpectedError.Code)
	}
	companyFilter := bson.M{"_id": camObj}
	company, err := d.donationCompanyDal.FindOne(ctx, companyFilter, nil)
	if err != nil {
		d.logger.Errorf("[FindByID] failed to find donation company: %v", err)
		return nil, err
	}

	catObj, err := bson.ObjectIDFromHex(result.CategoryID.Hex())
	if err != nil {
		d.logger.Errorf("[FindByID] failed to parse category id: %v", err)
		return nil, errors.New(localization.ErrorUnexpectedError.Code)
	}
	categoryFilter := bson.M{"_id": catObj}
	category, err := d.donationCategoryDal.FindOne(ctx, categoryFilter, nil)
	if err != nil {
		d.logger.Errorf("[FindByID] failed to find donation category: %v", err)
		return nil, err
	}
	d.logger.Infof("[FindByID] donation retrieved successfully")
	return MapToDonationListResponse(result, company, category), nil
}

func (d *DonationStorage) FindAllWithPagination(ctx context.Context, filterParam types.Filter) (*types.PaginatedResponse[[]donation_dto.DonationListResponse], error) {

	searchKeys := bson.M{}
	allowedKeys := []string{"search", "title", "is_featured", "enabled", "donation_code", "target", "end_date"}

	if filterParam.Search != "" {
		searchRegex := bson.M{"$regex": filterParam.Search, "$options": "i"}
		searchKeys["$or"] = []bson.M{
			{"title": searchRegex},
			{"is_featured": searchRegex},
			{"donation_description": searchRegex},
			{"donation_code": searchRegex},
			{"target": searchRegex},
			{"end_date": searchRegex},
			{"start_date": searchRegex},
			{"enabled": searchRegex},
		}
	}
	filter, skip, limit := lib.FilterBuilder(filterParam, searchKeys, allowedKeys)
	data, err := d.dal.FindAllWithPagination(ctx, filter, bson.M{}, skip, limit)
	if err != nil {
		d.logger.Errorf("[FindAllWithPagination] failed to fetch donations: %v", err)
		return nil, err
	}

	total, err := d.dal.TotalCount(ctx, filter)
	if err != nil {
		d.logger.Errorf("[FindAllWithPagination] failed to count donations: %v", err)
		return nil, err
	}

	meta := local_util.BuildPaginationMeta(total, filterParam.Page, filterParam.PerPage)

	var result []donation_dto.DonationListResponse
	for _, donation := range data {
		companyFilter := bson.M{"_id": donation.CompanyID}
		company, err := d.donationCompanyDal.FindOne(ctx, companyFilter, nil)
		if err != nil {
			if err == mongo.ErrNoDocuments {
				continue
			}
			d.logger.Errorf("[FindAllWithPagination] failed to find donation company: %v", err)
		}

		categoryFilter := bson.M{"_id": donation.CategoryID}
		category, err := d.donationCategoryDal.FindOne(ctx, categoryFilter, nil)
		if err != nil {
			if err == mongo.ErrNoDocuments {
				continue
			}
			d.logger.Errorf("[FindAllWithPagination] failed to find donation category: %v", err)
		}

		result = append(result, *MapToDonationListResponse(donation, company, category))
	}
	d.logger.Infof("[FindAllWithPagination] retrieved %d donations", len(result))

	return &types.PaginatedResponse[[]donation_dto.DonationListResponse]{
		Data: result,
		Meta: meta,
	}, nil
}
