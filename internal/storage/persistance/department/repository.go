package department

import (
	"cbe-super-app-cps-action/internal/constants/lib"
	"cbe-super-app-cps-action/internal/constants/localization"
	"cbe-super-app-cps-action/internal/constants/model"
	"cbe-super-app-cps-action/internal/constants/types"
	"cbe-super-app-cps-action/internal/storage"
	local_util "cbe-super-app-cps-action/pkgs/utils"
	"context"
	"errors"
	"strings"

	"gitlab.com/bersufekadgetachew/cbe-super-app-shared/shared/dal"
	"gitlab.com/bersufekadgetachew/cbe-super-app-shared/shared/utils"
	"go.mongodb.org/mongo-driver/v2/bson"
	"go.mongodb.org/mongo-driver/v2/mongo"
)

type DepartmentStorage struct {
	dal        dal.MongoDal[model.Department, model.Department]
	client     *mongo.Client
	logger     utils.Logger
	collection *mongo.Collection
}

func NewDepartmentRepository(client *mongo.Client, dbName string, collection string, logger utils.Logger) storage.DepartmentRepository {
	return &DepartmentStorage{
		dal:        dal.NewMongoDal[model.Department, model.Department](client, dbName, collection),
		client:     client,
		collection: client.Database(dbName).Collection(collection),
		logger:     logger,
	}
}

func (b *DepartmentStorage) Create(ctx context.Context, Department *model.Department) error {
	Department.ID = bson.NewObjectID()
	_, err := b.dal.InsertOne(ctx, *Department)
	if err != nil {
		return errors.New(localization.ErrorUnexpectedError.Code)
	}
	return nil
}

func (b *DepartmentStorage) Update(ctx context.Context, id string, Department *model.Department) error {
	objID, err := bson.ObjectIDFromHex(id)
	if err != nil {
		return errors.New(localization.ErrorDepartmentInvalidID.Code)
	}
	filter := bson.M{"_id": objID, "is_deleted": false}
	updateData := DepartmentMapper(*Department)

	_, err = b.dal.UpdateOne(ctx, filter, updateData)
	if err != nil {
		if err == mongo.ErrNoDocuments {
			return errors.New(localization.ErrorFileNotFound.Code)
		}
		return errors.New(localization.ErrorUnexpectedError.Code)
	}
	return nil
}

func (b *DepartmentStorage) Delete(ctx context.Context, id string) error {
	objID, err := bson.ObjectIDFromHex(id)
	if err != nil {
		return errors.New(localization.ErrorUnexpectedError.Code)
	}
	filter := bson.M{"_id": objID, "is_deleted": false}
	return b.dal.DeleteOne(ctx, filter)
}

func (b *DepartmentStorage) EnableOrDisable(ctx context.Context, id string, enable bool) error {
	objID, err := bson.ObjectIDFromHex(id)
	if err != nil {
		return errors.New(localization.ErrorUnexpectedError.Code)
	}
	filter := bson.M{"_id": objID}
	update := bson.M{"enabled": enable}
	_, err = b.dal.UpdateOne(ctx, filter, update)
	if err != nil {
		return err
	}
	return nil
}

func (b *DepartmentStorage) FindByID(ctx context.Context, id string) (*model.Department, error) {
	objID, err := bson.ObjectIDFromHex(id)
	if err != nil {
		return nil, errors.New(localization.ErrorDepartmentInvalidID.Code)
	}
	filter := bson.M{"_id": objID}

	result, err := b.dal.FindOne(ctx, filter, nil)
	if err != nil {
		if errors.Is(err, mongo.ErrNoDocuments) {
			b.logger.Warnf("department not found for ID: %s,error ", id, err)
			return nil, errors.New(localization.ErrorDepartmentNotFound.Code)
		}
		b.logger.Errorf("FindByID Department failed: %v", err)
		return nil, errors.New(localization.ErrorUnexpectedError.Code)
	}

	return result, nil
}
func (b *DepartmentStorage) FindByName(ctx context.Context, name string) (*model.Department, error) {
	filter := bson.M{
		"department": bson.M{
			"$regex":   "^" + strings.ToLower(name) + "$",
			"$options": "i",
		},
	}
	result, err := b.dal.FindOne(ctx, filter, nil)
	if err != nil {
		return nil, err
	}
	return result, nil
}

func (s *DepartmentStorage) FindAllWithPagination(ctx context.Context, filterParam types.Filter) (types.PaginatedResponse[[]model.Department], error) {

	searchKeys := bson.M{}

	// 2. Allowed filterable/searchable fields
	allowedKeys := []string{"department_code", "department", "enabled", "is_deleted", "created_at", "last_modified"}

	// 3. Add search (if provided)
	if filterParam.Search != "" {
		searchRegex := bson.M{"$regex": filterParam.Search, "$options": "i"}
		searchKeys["$or"] = []bson.M{
			{"department_code": searchRegex},
			{"department": searchRegex},
			{"created_at": searchRegex},
			{"last_modified_at": searchRegex},
		}

	}
	filter, skip, limit := lib.FilterBuilder(filterParam, searchKeys, allowedKeys)
	s.logger.Debugf("Mongo filter: %+v, skip: %d, limit: %d", filter, skip, limit)
	Filter := dal.FilterOp{
		Filter: filter,
		Limit:  limit,
	}

	data, err := s.dal.FindAllWithCursorBasedPagination(ctx, Filter)
	if err != nil {
		s.logger.Errorf("Error fetching Department list")
		return types.PaginatedResponse[[]model.Department]{}, errors.New(localization.ErrorUnexpectedError.Message)
	}

	// 6. Count total
	total, err := s.dal.TotalCount(ctx, filter)
	if err != nil {
		return types.PaginatedResponse[[]model.Department]{}, errors.New(localization.ErrorUnexpectedError.Message)
	}

	// 7. Build pagination metadata
	meta := local_util.BuildPaginationMeta(total, filterParam.Page, filterParam.PerPage)

	// 8. Return standard paginated response
	s.logger.Infof("Successfully fetched paginated department list. Total: %d", total)
	return types.PaginatedResponse[[]model.Department]{
		Data: data,
		Meta: meta,
	}, nil
}
