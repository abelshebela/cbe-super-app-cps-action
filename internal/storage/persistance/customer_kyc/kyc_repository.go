package customer

import (
	"cbe-super-app-cps-action/internal/constants"
	"cbe-super-app-cps-action/internal/constants/lib"
	"cbe-super-app-cps-action/internal/constants/localization"
	imodel "cbe-super-app-cps-action/internal/constants/model"
	"cbe-super-app-cps-action/internal/constants/types"
	"cbe-super-app-cps-action/internal/storage"
	local_util "cbe-super-app-cps-action/pkgs/utils"
	"context"
	"errors"
	"time"

	"gitlab.com/bersufekadgetachew/cbe-super-app-shared/shared/config"
	"gitlab.com/bersufekadgetachew/cbe-super-app-shared/shared/dal"
	"gitlab.com/bersufekadgetachew/cbe-super-app-shared/shared/utils"
	"go.mongodb.org/mongo-driver/v2/bson"
	"go.mongodb.org/mongo-driver/v2/mongo"
)

type customerKYCRepository struct {
	dal    dal.MongoDal[imodel.CustomerKYC, imodel.CustomerKYC]
	logger utils.Logger
	coll   *mongo.Collection
}

func NewCustomerKYCRepository(client *mongo.Client, cfg *config.VaultConfig, dbName string, collection string, logger utils.Logger) storage.CustomerKYCRepository {
	return &customerKYCRepository{
		dal:    dal.NewMongoDal[imodel.CustomerKYC, imodel.CustomerKYC](client, cfg, dbName, collection),
		logger: logger,
		coll:   client.Database(dbName).Collection(collection),
	}
}

func (r *customerKYCRepository) Create(ctx context.Context, kyc *imodel.CustomerKYC) error {
	r.logger.Infof("[CustomerKYC][Create] creating kyc for customer: %s", kyc.CustomerCode)
	kyc.CreatedAt = time.Now()
	kyc.UpdatedAt = time.Now()
	kyc.KYCStatus = "PENDING"
	kyc.CustomerStatus = constants.CustomerPending

	if _, err := r.dal.InsertOne(ctx, *kyc); err != nil {
		r.logger.Errorf("[CustomerKYC][Create] failed to create kyc: %v", err)
		return errors.New(localization.ErrorUnexpectedError.Code)
	}

	return nil
}

func (r *customerKYCRepository) FindAllWithPagination(ctx context.Context, filterParam types.Filter) (*types.PaginatedResponse[[]imodel.CustomerKYC], error) {
	r.logger.Infof("[CustomerKYC][FindAllWithPagination] fetching kyc requests")

	allowed := []string{"search"}
	filter, skip, limit := lib.FilterBuilder(filterParam, bson.M{}, allowed)
	if filterParam.Search != "" {
		q := bson.M{"$regex": filterParam.Search, "$options": "i"}
		filter["$or"] = []bson.M{{"service_name": q}, {"service_code": q}, {"service_type": q}}
	}

	results, err := r.dal.FindAllWithPagination(ctx, filter, bson.M{}, skip, limit)
	if err != nil {
		r.logger.Errorf("[CustomerKYC][FindAllWithPagination] failed to fetch data: %v", err)
		return nil, errors.New(localization.ErrorUnexpectedError.Code)
	}

	total, err := r.dal.TotalCount(ctx, filter)
	if err != nil {
		r.logger.Errorf("[CustomerKYC][FindAllWithPagination] failed to get total count: %v", err)
		return nil, errors.New(localization.ErrorUnexpectedError.Code)
	}

	// resultData := make([]*imodel.CustomerKYC, len(results))
	// for i := range data {
	// 	resultData[i] = &data[i]
	// }

	meta := local_util.BuildPaginationMeta(total, filterParam.Page, filterParam.PerPage)

	return &types.PaginatedResponse[[]imodel.CustomerKYC]{
		Data: results,
		Meta: meta,
	}, nil
}

func (r *customerKYCRepository) FindByID(ctx context.Context, id string) (*imodel.CustomerKYC, error) {
	r.logger.Infof("[CustomerKYC][FindByID] fetching kyc by id: %s", id)
	objID, err := bson.ObjectIDFromHex(id)
	if err != nil {
		r.logger.Errorf("[CustomerKYC][FindByID] invalid object id: %v", err)
		return nil, errors.New(localization.ErrorInvalidID.Code)
	}

	filter := bson.M{"_id": objID}
	result, err := r.dal.FindOne(ctx, filter, nil)
	if err != nil {
		r.logger.Errorf("[CustomerKYC][FindByID] failed to find kyc: %v", err)
		return nil, local_util.HandleDBError(err)
	}

	return result, nil
}

func (r *customerKYCRepository) Delete(ctx context.Context, id string) error {
	r.logger.Infof("[CustomerKYC][Delete] hard deleting kyc for id: %s", id)
	objID, err := bson.ObjectIDFromHex(id)
	if err != nil {
		r.logger.Errorf("[CustomerKYC][Delete] invalid object id: %v", err)
		return errors.New(localization.ErrorInvalidID.Code)
	}

	filter := bson.M{"_id": objID}
	if err := r.dal.DeleteOneH(ctx, filter); err != nil {
		r.logger.Errorf("[CustomerKYC][Delete] failed to delete kyc: %v", err)
		return errors.New(localization.ErrorUnexpectedError.Code)
	}

	return nil
}

func (r *customerKYCRepository) UpdateKYCStatus(ctx context.Context, id string, status string) error {
	r.logger.Infof("[CustomerKYC][UpdateKYCStatus] updating kyc status for id: %s to %s", id, status)
	objID, err := bson.ObjectIDFromHex(id)
	if err != nil {
		r.logger.Errorf("[CustomerKYC][UpdateKYCStatus] invalid object id: %v", err)
		return errors.New(localization.ErrorInvalidID.Code)
	}

	filter := bson.M{"_id": objID}
	update := bson.M{"kyc_status": status}

	if _, err := r.dal.UpdateOne(ctx, filter, update); err != nil {
		r.logger.Errorf("[CustomerKYC][UpdateKYCStatus] failed to update kyc status: %v", err)
		return errors.New(localization.ErrorUnexpectedError.Code)
	}

	return nil
}
