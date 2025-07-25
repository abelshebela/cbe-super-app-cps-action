package service

import (
	"context"
	"time"

	"github.com/CBE-Super-App/cbe-super-app-cps-action/internal/adapter/outbound/model"
	infra_mongo "github.com/CBE-Super-App/cbe-super-app-cps-action/internal/adapter/outbound/mongo"
	"github.com/CBE-Super-App/cbe-super-app-cps-action/pkgs/common"
	contexts "github.com/CBE-Super-App/cbe-super-app-cps-action/pkgs/context"
	local_utils "github.com/CBE-Super-App/cbe-super-app-cps-action/pkgs/utils"
	"gitlab.com/bersufekadgetachew/cbe-super-app-shared/shared/utils"
	"go.mongodb.org/mongo-driver/v2/bson"
	"go.mongodb.org/mongo-driver/v2/mongo"
)

type servicePersistence struct {
	client     *mongo.Client
	logger     utils.Logger
	serviceDal *infra_mongo.MongoDal[model.Service, model.Service]
	cpsDal     *infra_mongo.MongoDal[model.CPSAction, model.CPSAction]
}

type ServiceRepository interface {
	GetAllService(ctx context.Context, filterParams local_utils.Filter) (*local_utils.PaginatedResponse[*any], error)
	GetAllMinimumTransferCap(ctx context.Context, filterParams local_utils.Filter) (*local_utils.PaginatedResponse[*any], error)
	GetAllMaximumTransferCap(ctx context.Context, filterParams local_utils.Filter) (*local_utils.PaginatedResponse[*any], error)
	GetAllServiceFee(ctx context.Context, filterParams local_utils.Filter) (*local_utils.PaginatedResponse[*any], error)
	GetAllTotalTransferCap(ctx context.Context, filterParams local_utils.Filter) (*local_utils.PaginatedResponse[*any], error)
	GetServiceFeeDetail(ctx context.Context, id string) (*any, error)
	UpdateServiceFee(ctx context.Context, id string, req any) error
	UpdateSingleMaxTransfer(ctx context.Context, id string, req any) error
	UpdateTotalMaxTransferCap(ctx context.Context, id string, req any) error
	UpdateMinimumTransferCap(ctx context.Context, id string, req any) error
	DeleteServiceFeeTire(ctx context.Context, id string) error
}

func NewServicePersistence(client *mongo.Client, dbName string, collection []string, logger utils.Logger) ServiceRepository {
	return &servicePersistence{
		client:     client,
		logger:     logger,
		serviceDal: infra_mongo.NewMongoDal[model.Service, model.Service](client, dbName, collection[1]),
		cpsDal:     infra_mongo.NewMongoDal[model.CPSAction, model.CPSAction](client, dbName, collection[0]),
	}
}

func (sp *servicePersistence) GetAllService(ctx context.Context, filterParams local_utils.Filter) (*local_utils.PaginatedResponse[*any], error) {
	filter := bson.M{
		"is_deleted": false,
	}
	projection := bson.M{}
	if filterParams.Search != "" {
		filter["$or"] = []bson.M{
			{"service_name": bson.M{"$regex": filterParams.Search, "$options": "i"}},
			{"service_code": bson.M{"$regex": filterParams.Search, "$options": "i"}},
			{"service_type": bson.M{"$regex": filterParams.Search, "$options": "i"}},
		}
	}

	if filterParams.Filters != "" {
		filter["$or"] = []bson.M{
			{"service_type": bson.M{"$regex": filterParams.Filters, "$options": "i"}},
			{"payment_type": bson.M{"$regex": filterParams.Filters, "$options": "i"}},
		}
	}

	skip := (filterParams.Page - 1) * filterParams.PerPage
	limit := filterParams.PerPage

	service, err := sp.serviceDal.FindAllWithPagination(ctx, filter, projection, int64(skip), int64(limit))
	if err != nil {
		sp.logger.Errorf("failed to fetch service: %v:", err)
		return nil, common.DefineError.General["UNHANDLED_SERVER_ERROR"]
	}

	if len(service) == 0 {
		sp.logger.Errorf("no service found")
		return nil, common.DefineError.General["SERVICE_NOT_FOUND"]
	}

	total, err := sp.serviceDal.TotalCount(ctx, filter)
	if err != nil {
		return nil, common.DefineError.General["UNHANDLED_SERVER_ERROR"]
	}

	meta := local_utils.BuildPaginationMeta(total, filterParams.Page, limit)
	var data any = service
	return &local_utils.PaginatedResponse[*any]{
		Data: &data,
		Meta: meta,
	}, nil
}
func (sp *servicePersistence) GetAllMinimumTransferCap(ctx context.Context, filterParams local_utils.Filter) (*local_utils.PaginatedResponse[*any], error) {
	filter := bson.M{
		"is_deleted": false,
	}

	projection := bson.M{
		"id":             1,
		"service_name":   1,
		"service_code":   1,
		"service_type":   1,
		"cap.min_amount": 1,
		"created_at":     1,
	}
	if filterParams.Search != "" {
		filter["$or"] = []bson.M{
			{"service_name": bson.M{"$regex": filterParams.Search, "$options": "i"}},
			{"service_code": bson.M{"$regex": filterParams.Search, "$options": "i"}},
			{"service_type": bson.M{"$regex": filterParams.Search, "$options": "i"}},
		}
	}

	if filterParams.Filters != "" {
		filter["$or"] = []bson.M{
			{"service_type": bson.M{"$regex": filterParams.Filters, "$options": "i"}},
			{"payment_type": bson.M{"$regex": filterParams.Filters, "$options": "i"}},
		}
	}

	skip := (filterParams.Page - 1) * filterParams.PerPage
	limit := filterParams.PerPage

	service, err := sp.serviceDal.FindAllWithPagination(ctx, filter, projection, int64(skip), int64(limit))
	if err != nil {
		sp.logger.Errorf("failed to fetch service: %v:", err)
		return nil, common.DefineError.General["UNHANDLED_SERVER_ERROR"]
	}

	if len(service) == 0 {
		sp.logger.Errorf("no service found")
		return nil, common.DefineError.General["SERVICE_NOT_FOUND"]
	}

	total, err := sp.serviceDal.TotalCount(ctx, filter)
	if err != nil {
		return nil, common.DefineError.General["UNHANDLED_SERVER_ERROR"]
	}

	meta := local_utils.BuildPaginationMeta(total, filterParams.Page, limit)
	var data any = service
	return &local_utils.PaginatedResponse[*any]{
		Data: &data,
		Meta: meta,
	}, nil
}
func (sp *servicePersistence) GetAllMaximumTransferCap(ctx context.Context, filterParams local_utils.Filter) (*local_utils.PaginatedResponse[*any], error) {
	filter := bson.M{
		"is_deleted": false,
	}
	projection := bson.M{
		"id":                        1,
		"service_name":              1,
		"service_code":              1,
		"service_type":              1,
		"cap.individual_single_cap": 1,
		"cap.individual_daily_cap":  1,
		"cap.corporate_single_cap":  1,
		"cap.corporate_daily_cap":   1,
		"created_at":                1,
	}
	if filterParams.Search != "" {
		filter["$or"] = []bson.M{
			{"service_name": bson.M{"$regex": filterParams.Search, "$options": "i"}},
			{"service_code": bson.M{"$regex": filterParams.Search, "$options": "i"}},
			{"service_type": bson.M{"$regex": filterParams.Search, "$options": "i"}},
		}
	}

	if filterParams.Filters != "" {
		filter["$or"] = []bson.M{
			{"service_type": bson.M{"$regex": filterParams.Filters, "$options": "i"}},
			{"payment_type": bson.M{"$regex": filterParams.Filters, "$options": "i"}},
		}
	}

	skip := (filterParams.Page - 1) * filterParams.PerPage
	limit := filterParams.PerPage

	service, err := sp.serviceDal.FindAllWithPagination(ctx, filter, projection, int64(skip), int64(limit))
	if err != nil {
		sp.logger.Errorf("failed to fetch service: %v:", err)
		return nil, common.DefineError.General["UNHANDLED_SERVER_ERROR"]
	}

	if len(service) == 0 {
		sp.logger.Errorf("no service found")
		return nil, common.DefineError.General["SERVICE_NOT_FOUND"]
	}

	total, err := sp.serviceDal.TotalCount(ctx, filter)
	if err != nil {
		return nil, common.DefineError.General["UNHANDLED_SERVER_ERROR"]
	}

	meta := local_utils.BuildPaginationMeta(total, filterParams.Page, limit)
	var data any = service
	return &local_utils.PaginatedResponse[*any]{
		Data: &data,
		Meta: meta,
	}, nil
}
func (sp *servicePersistence) GetAllServiceFee(ctx context.Context, filterParams local_utils.Filter) (*local_utils.PaginatedResponse[*any], error) {
	filter := bson.M{
		"is_deleted": false,
	}
	projection := bson.M{
		"id":           1,
		"service_name": 1,
		"service_code": 1,
		"service_type": 1,
		"payment_type": 1,
		"tire":         1,
		"created_at":   1,
	}
	if filterParams.Search != "" {
		filter["$or"] = []bson.M{
			{"service_name": bson.M{"$regex": filterParams.Search, "$options": "i"}},
			{"service_code": bson.M{"$regex": filterParams.Search, "$options": "i"}},
			{"service_type": bson.M{"$regex": filterParams.Search, "$options": "i"}},
		}
	}

	if filterParams.Filters != "" {
		filter["$or"] = []bson.M{
			{"service_type": bson.M{"$regex": filterParams.Filters, "$options": "i"}},
			{"payment_type": bson.M{"$regex": filterParams.Filters, "$options": "i"}},
		}
	}

	skip := (filterParams.Page - 1) * filterParams.PerPage
	limit := filterParams.PerPage

	service, err := sp.serviceDal.FindAllWithPagination(ctx, filter, projection, int64(skip), int64(limit))
	if err != nil {
		sp.logger.Errorf("failed to fetch service: %v:", err)
		return nil, common.DefineError.General["UNHANDLED_SERVER_ERROR"]
	}

	if len(service) == 0 {
		sp.logger.Errorf("no service found")
		return nil, common.DefineError.General["SERVICE_NOT_FOUND"]
	}

	total, err := sp.serviceDal.TotalCount(ctx, filter)
	if err != nil {
		return nil, common.DefineError.General["UNHANDLED_SERVER_ERROR"]
	}

	meta := local_utils.BuildPaginationMeta(total, filterParams.Page, limit)
	var data any = service
	return &local_utils.PaginatedResponse[*any]{
		Data: &data,
		Meta: meta,
	}, nil
}
func (sp *servicePersistence) GetAllTotalTransferCap(ctx context.Context, filterParams local_utils.Filter) (*local_utils.PaginatedResponse[*any], error) {
	filter := bson.M{
		"is_deleted": false,
	}
	projection := bson.M{
		"id":                1,
		"service_name":      1,
		"service_code":      1,
		"service_type":      1,
		"payment_type":      1,
		"above_amount":      1,
		"above_service_fee": 1,
		"cap.min_amount":    1,
		"created_at":        1,
	}
	if filterParams.Search != "" {
		filter["$or"] = []bson.M{
			{"service_name": bson.M{"$regex": filterParams.Search, "$options": "i"}},
			{"service_code": bson.M{"$regex": filterParams.Search, "$options": "i"}},
			{"service_type": bson.M{"$regex": filterParams.Search, "$options": "i"}},
		}
	}

	if filterParams.Filters != "" {
		filter["$or"] = []bson.M{
			{"service_type": bson.M{"$regex": filterParams.Filters, "$options": "i"}},
			{"payment_type": bson.M{"$regex": filterParams.Filters, "$options": "i"}},
		}
	}

	skip := (filterParams.Page - 1) * filterParams.PerPage
	limit := filterParams.PerPage

	service, err := sp.serviceDal.FindAllWithPagination(ctx, filter, projection, int64(skip), int64(limit))
	if err != nil {
		sp.logger.Errorf("failed to fetch service: %v:", err)
		return nil, common.DefineError.General["UNHANDLED_SERVER_ERROR"]
	}

	if len(service) == 0 {
		sp.logger.Errorf("no service found")
		return nil, common.DefineError.General["SERVICE_NOT_FOUND"]
	}

	total, err := sp.serviceDal.TotalCount(ctx, filter)
	if err != nil {
		return nil, common.DefineError.General["UNHANDLED_SERVER_ERROR"]
	}

	meta := local_utils.BuildPaginationMeta(total, filterParams.Page, limit)
	var data any = service
	return &local_utils.PaginatedResponse[*any]{
		Data: &data,
		Meta: meta,
	}, nil
}
func (sp *servicePersistence) GetServiceFeeDetail(ctx context.Context, id string) (*any, error) {
	objId, err := bson.ObjectIDFromHex(id)
	if err != nil {
		return nil, common.DefineError.General["UNHANDLED_SERVER_ERROR"]
	}
	filter := bson.M{
		"_id":        objId,
		"is_deleted": false,
	}
	projection := bson.M{}

	serviceDetail, err := sp.serviceDal.FindOne(ctx, filter, projection)
	if err != nil {
		return nil, common.DefineError.General["UNHANDLED_SERVER_ERROR"]
	}
	var data any = serviceDetail
	return &data, nil
}
func (sp *servicePersistence) UpdateServiceFee(ctx context.Context, id string, req any) error {
	actionExist, err := sp.checkPendingAction(ctx, string(model.RequestUpdateServiceFee))
	if err != nil {
		return err
	}

	if actionExist {
		return common.DefineError.General["PENDING_REQUEST_EXISTS"]
	}
	projection := bson.M{
		"tire": 1,
	}
	prev, err := sp.prevServiceData(ctx, id, projection)
	tiers, err := local_utils.JsonUnmarshalArray[model.Tier](prev.Tiers)
	if err != nil {
		return err
	}

	reqData, err := local_utils.JsonUnmarshal[model.Tier](req)
	if err != nil {
		return err
	}
	tireData, err := FindTierByID(*tiers, reqData.ID.Hex())
	if err != nil {
		return nil
	}
	err = sp.updateCpsAction(ctx, id, tireData, reqData, string(model.RequestUpdateServiceFee))
	if err != nil {
		return err
	}
	return nil
}
func (sp *servicePersistence) UpdateSingleMaxTransfer(ctx context.Context, id string, req any) error {
	actionExist, err := sp.checkPendingAction(ctx, string(model.RequestUpdateServiceFee))
	if err != nil {
		return err
	}

	if actionExist {
		return common.DefineError.General["PENDING_REQUEST_EXISTS"]
	}
	reqData, err := local_utils.JsonUnmarshal[model.Cap](req)
	if err != nil {
		return err
	}

	projection := bson.M{
		"cap.individual_single_cap": 1,
		"cap.individual_daily_cap":  1,
		"cap.corporate_single_cap":  1,
		"cap.corporate_daily_cap":   1,
	}
	prev, err := sp.prevServiceData(ctx, id, projection)
	if err != nil {
		return err
	}
	err := sp.createCpsAction(ctx, reqDat)

	return nil
}
func (sp *servicePersistence) UpdateTotalMaxTransferCap(ctx context.Context, id string, req any) error
func (sp *servicePersistence) UpdateMinimumTransferCap(ctx context.Context, id string, req any) error
func (sp *servicePersistence) DeleteMinimumTransferCap(ctx context.Context, id string) error
func (sp *servicePersistence) DeleteServiceFeeTire(ctx context.Context, id string) error
func (sp *servicePersistence) checkPendingAction(ctx context.Context, requestAction string) (bool, error) {
	userData := contexts.ExtractContext(ctx)
	filter := bson.M{
		"maker_id":       userData.UserID,
		"department":     userData.Department,
		"action_status":  "PENDING",
		"request_action": requestAction,
	}

	dataCpsAction, err := sp.cpsDal.FindOne(ctx, filter, nil)
	if err != nil {
		return false, nil
	}

	if dataCpsAction != nil {
		return false, nil
	}

	return true, nil
}
func (sp *servicePersistence) createCpsAction(ctx context.Context, currentAction any, requestAction string) error {
	userData := contexts.ExtractContext(ctx)

	cpsAction := &model.CPSAction{
		ActionCode:       local_utils.GenerateRandom(20),
		MakerID:          userData.UserID,
		MakerName:        userData.FullName,
		MakerPhoneNumber: userData.PhoneNumber,
		Department:       userData.Department,
		ActionStatus:     "PENDING",
		CurrentAction:    currentAction,
		CreatedAt:        time.Now(),
		LastModifiedAt:   time.Now(),
		RequestAction:    requestAction,
		MakerActionTime:  time.Now(),
	}

	_, err := sp.cpsDal.InsertOne(ctx, *cpsAction)
	if err != nil {
		return common.DefineError.General["UNHANDLED_SERVER_ERROR"]
	}

	return nil
}
func (sp *servicePersistence) updateCpsAction(ctx context.Context, id string, prevAction any, currentAction any, requestAction string) error {
	userData := contexts.ExtractContext(ctx)
	objId, err := bson.ObjectIDFromHex(id)
	if err != nil {
		return err
	}
	filter := bson.M{
		"_id": objId,
	}
	cpsAction := &model.CPSAction{
		ActionCode:       local_utils.GenerateRandom(20),
		MakerID:          userData.UserID,
		MakerName:        userData.FullName,
		MakerPhoneNumber: userData.PhoneNumber,
		Department:       userData.Department,
		ActionStatus:     "PENDING",
		CurrentAction:    currentAction,
		PreviousAction:   prevAction,
		CreatedAt:        time.Now(),
		LastModifiedAt:   time.Now(),
		RequestAction:    requestAction,
		MakerActionTime:  time.Now(),
	}
	data, err := local_utils.JsonUnmarshal[bson.M](cpsAction)
	if err != nil {
		return err
	}
	_, err = sp.cpsDal.UpdateOne(ctx, filter, *data)
	if err != nil {
		return common.DefineError.General["UNHANDLED_SERVER_ERROR"]
	}

	return nil
}
func (sp *servicePersistence) prevServiceData(ctx context.Context, id string, projection bson.M) (*model.Service, error) {
	objId, err := bson.ObjectIDFromHex(id)
	if err != nil {
		return nil, err
	}

	filter := bson.M{
		"_id":        objId,
		"is_deleted": objId,
	}
	prev, err := sp.serviceDal.FindOne(ctx, filter, projection)
	if err != nil {
		return nil, err
	}

	return prev, nil
}
func FindTierByID(tiers []model.Tier, id string) (*model.Tier, error) {
	objId, err := bson.ObjectIDFromHex(id)
	if err != nil {
		return nil, err
	}
	for i := range tiers {
		if tiers[i].ID == objId {
			return &tiers[i], nil
		}
	}
	return nil, common.DefineError.General["INVALID_INPUT_PARAMETERS"]
}
