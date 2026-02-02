package department

import (
	"cbe-super-app-cps-action/internal/constants/lib"
	"cbe-super-app-cps-action/internal/constants/localization"
	"cbe-super-app-cps-action/internal/constants/types"
	"cbe-super-app-cps-action/internal/storage"
	local_util "cbe-super-app-cps-action/pkgs/utils"
	"context"
	"errors"
	"strings"

	"gitlab.com/bersufekadgetachew/cbe-super-app-shared/shared/entities/model"

	"gitlab.com/bersufekadgetachew/cbe-super-app-shared/shared/config"
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

func NewDepartmentRepository(client *mongo.Client, cfg *config.VaultConfig, dbName string, collection string, logger utils.Logger) storage.DepartmentRepository {
	return &DepartmentStorage{
		dal:        dal.NewMongoDal[model.Department, model.Department](client, cfg, dbName, collection),
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
	b.logger.Infof("[Update] updating department for id: %s", id)
	objID, err := bson.ObjectIDFromHex(id)
	if err != nil {
		b.logger.Errorf("[Update] invalid object id: %v", err)
		return errors.New(localization.ErrorDepartmentInvalidID.Code)
	}
	filter := bson.M{"_id": objID, "is_deleted": false}
	updateData := DepartmentMapper(*Department)

	_, err = b.dal.UpdateOne(ctx, filter, updateData)
	if err != nil {
		if err == mongo.ErrNoDocuments {
			b.logger.Errorf("[Update] department not found")
			return errors.New(localization.ErrorFileNotFound.Code)
		}
		b.logger.Errorf("[Update] failed to update department: %v", err)
		return errors.New(localization.ErrorUnexpectedError.Code)
	}
	b.logger.Infof("[Update] department updated successfully")
	return nil
}

func (b *DepartmentStorage) Delete(ctx context.Context, id string) error {
	b.logger.Infof("[Delete] deleting department for id: %s", id)
	objID, err := bson.ObjectIDFromHex(id)
	if err != nil {
		b.logger.Errorf("[Delete] invalid object id: %v", err)
		return errors.New(localization.ErrorUnexpectedError.Code)
	}
	filter := bson.M{"_id": objID, "is_deleted": false}
	err = b.dal.DeleteOne(ctx, filter)
	if err != nil {
		b.logger.Errorf("[Delete] failed to delete department: %v", err)
		return err
	}
	b.logger.Infof("[Delete] department deleted successfully")
	return nil
}

func (b *DepartmentStorage) EnableOrDisable(ctx context.Context, id string, enable bool) error {
	b.logger.Infof("[EnableOrDisable] processing department enable/disable for id: %s, enabled: %v", id, enable)
	objID, err := bson.ObjectIDFromHex(id)
	if err != nil {
		b.logger.Errorf("[EnableOrDisable] invalid object id: %v", err)
		return errors.New(localization.ErrorUnexpectedError.Code)
	}
	filter := bson.M{"_id": objID}
	update := bson.M{"enabled": enable}
	_, err = b.dal.UpdateOne(ctx, filter, update)
	if err != nil {
		b.logger.Errorf("[EnableOrDisable] failed to enable/disable department: %v", err)
		return err
	}
	b.logger.Infof("[EnableOrDisable] department enable/disable completed successfully")
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
			b.logger.Errorf("[FindByID] department not found")
			return nil, errors.New(localization.ErrorDepartmentNotFound.Code)
		}
		b.logger.Errorf("[FindByID] failed to find department: %v", err)
		return nil, errors.New(localization.ErrorUnexpectedError.Code)
	}
	b.logger.Infof("[FindByID] department retrieved successfully")
	return result, nil
}
func (b *DepartmentStorage) FindByName(ctx context.Context, name string) (*model.Department, error) {
	b.logger.Infof("[FindByName] searching for department by name")
	filter := bson.M{
		"department": bson.M{
			"$regex":   "^" + strings.ToLower(name) + "$",
			"$options": "i",
		},
	}
	result, err := b.dal.FindOne(ctx, filter, nil)
	if err != nil {
		b.logger.Errorf("[FindByName] failed to find department: %v", err)
		return nil, err
	}
	b.logger.Infof("[FindByName] department retrieved successfully")
	return result, nil
}

func (s *DepartmentStorage) FindAllWithPagination(ctx context.Context, filterParam types.Filter) (types.PaginatedResponse[[]model.Department], error) {

	searchKeys := bson.M{}

	// 2. Allowed filterable/searchable fields
	allowedKeys := []string{"search", "department_code", "department", "enabled", "is_deleted", "created_at", "last_modified"}

	// 3. Add search (if provided)
	if filterParam.Search != "" {
		searchRegex := bson.M{"$regex": filterParam.Search, "$options": "i"}
		searchKeys["$or"] = []bson.M{
			{"search": searchRegex},
			{"department_code": searchRegex},
			{"department": searchRegex},
			{"created_at": searchRegex},
			{"last_modified_at": searchRegex},
		}

	}
	filter, skip, limit := lib.FilterBuilder(filterParam, searchKeys, allowedKeys)
	// Filter := dal.FilterOp{
	// 	Filter: filter,
	// 	Limit:  limit,
	// }

	// data, err := s.dal.FindAllWithCursorBasedPagination(ctx, Filter)
	// if err != nil {
	// 	s.logger.Errorf("[FindAllWithPagination] failed to fetch departments: %v", err)
	// 	return types.PaginatedResponse[[]model.Department]{}, errors.New(localization.ErrorUnexpectedError.Message)
	// }
	results, err := s.dal.FindAllWithPagination(ctx, filter, bson.M{}, skip, limit)
	if err != nil {
		return types.PaginatedResponse[[]model.Department]{}, err
	}

	// 6. Count total
	total, err := s.dal.TotalCount(ctx, filter)
	if err != nil {
		s.logger.Errorf("[FindAllWithPagination] failed to count departments: %v", err)
		return types.PaginatedResponse[[]model.Department]{}, errors.New(localization.ErrorUnexpectedError.Message)
	}

	// 7. Build pagination metadata
	meta := local_util.BuildPaginationMeta(total, filterParam.Page, filterParam.PerPage)
	s.logger.Infof("[FindAllWithPagination] retrieved %d departments", len(results))

	// 8. Return standard paginated response
	return types.PaginatedResponse[[]model.Department]{
		Data: results,
		Meta: meta,
	}, nil
}
