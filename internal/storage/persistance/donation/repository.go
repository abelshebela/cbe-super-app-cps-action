package donation

import (
	"cbe-super-app-cps-action/internal/constants/model"
	"cbe-super-app-cps-action/internal/constants/types"
	"cbe-super-app-cps-action/internal/localization"
	"cbe-super-app-cps-action/internal/storage"
	"context"
	"errors"

	local_util "cbe-super-app-cps-action/pkgs/utils"

	"gitlab.com/bersufekadgetachew/cbe-super-app-shared/shared/dal"
	"gitlab.com/bersufekadgetachew/cbe-super-app-shared/shared/utils"
	"go.mongodb.org/mongo-driver/v2/bson"
	"go.mongodb.org/mongo-driver/v2/mongo"
)

type DonationStorage struct {
	dal    dal.MongoDal[model.Donation, model.Donation]
	client *mongo.Client
	logger utils.Logger
}

func NewDonationRepository(client *mongo.Client, dbName string, collection string, logger utils.Logger) storage.DonationRepository {
	return &DonationStorage{
		dal:    dal.NewMongoDal[model.Donation, model.Donation](client, dbName, collection),
		client: client,
		logger: logger,
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
	updateData := DonationMapper(*donation)

	_, err = d.dal.UpdateOne(ctx, filter, updateData)
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

func (d *DonationStorage) FindByID(ctx context.Context, id string) (*model.Donation, error) {
	objID, err := bson.ObjectIDFromHex(id)
	if err != nil {
		return nil, errors.New(localization.ErrorUnexpectedError.Code)
	}
	filter := bson.M{"_id": objID}

	result, err := d.dal.FindOne(ctx, filter, nil)
	if err != nil {
		return nil, err
	}
	return result, nil
}

func (d *DonationStorage) FindAllWithPagination(ctx context.Context, filterParam types.Filter) (*types.PaginatedResponse[[]*model.Donation], error) {
	filter := bson.M{"is_deleted": false}

	if filterParam.Search != "" {
		searchRegex := bson.M{"$regex": filterParam.Search, "$options": "i"}
		filter["name"] = searchRegex
	}

	skip := int64((filterParam.Page - 1) * filterParam.PerPage)
	limit := int64(filterParam.PerPage)

	data, err := d.dal.FindAllWithPagination(ctx, filter, bson.M{}, skip, limit)
	if err != nil {
		return nil, err
	}

	total, err := d.dal.TotalCount(ctx, filter)
	if err != nil {
		return nil, err
	}

	meta := local_util.BuildPaginationMeta(total, filterParam.Page, filterParam.PerPage)

	return &types.PaginatedResponse[[]*model.Donation]{
		Data: data,
		Meta: meta,
	}, nil
}
