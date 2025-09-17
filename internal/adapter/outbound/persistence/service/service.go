package service

import (
	"context"
	"fmt"
	"reflect"
	"strings"
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
	hqDal      *infra_mongo.MongoDal[model.HQ, model.HQ]
	cpsDal     *infra_mongo.MongoDal[model.CPSAction, model.CPSAction]
}

type ServiceRepository interface {
	GetAllService(ctx context.Context, filterParams local_utils.Filter) (*local_utils.PaginatedResponse[*any], error)
	GetAllMinimumTransferCap(ctx context.Context, filterParams local_utils.Filter) (*local_utils.PaginatedResponse[*any], error)
	GetAllMaximumTransferCap(ctx context.Context, filterParams local_utils.Filter) (*local_utils.PaginatedResponse[*any], error)
	GetAllServiceFee(ctx context.Context, filterParams local_utils.Filter) (*local_utils.PaginatedResponse[*any], error)
	GetAllTotalTransferCap(ctx context.Context) (*any, error)
	GetServiceFeeDetail(ctx context.Context, id string) (*any, error)
	UpdateServiceFee(ctx context.Context, id string, req any) error
	UpdateSingleMaxTransfer(ctx context.Context, id string, req any) error
	UpdateTotalMaxTransferCap(ctx context.Context, id string, newTotalCap uint64) error
	UpdateMinimumTransferCap(ctx context.Context, id string, req any) error
	DeleteServiceFeeTire(ctx context.Context, id string) error
	Authorize(ctx context.Context, cpsAction any) (any, error)
}

func NewServicePersistence(client *mongo.Client, dbName string, collection []string, logger utils.Logger) ServiceRepository {
	return &servicePersistence{
		client:     client,
		logger:     logger,
		serviceDal: infra_mongo.NewMongoDal[model.Service, model.Service](client, dbName, collection[1]),
		hqDal:      infra_mongo.NewMongoDal[model.HQ, model.HQ](client, dbName, collection[2]),
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

	sp.logger.Infof("Fetching all services with filter: %+v, projection: %+v, skip: %d, limit: %d", filter, projection, skip, limit)
	service, err := sp.serviceDal.FindAllWithPagination(ctx, filter, projection, int64(skip), int64(limit))
	if err != nil {
		sp.logger.Errorf("failed to fetch service: %v", err)
		return nil, common.DefineError.General["UNHANDLED_SERVER_ERROR"]
	}

	if len(service) == 0 {
		sp.logger.Warnf("no service found for filter: %+v", filter)
		return nil, common.DefineError.General["SERVICE_NOT_FOUND"]
	}

	total, err := sp.serviceDal.TotalCount(ctx, filter)
	if err != nil {
		sp.logger.Errorf("failed to count total services: %v", err)
		return nil, common.DefineError.General["UNHANDLED_SERVER_ERROR"]
	}

	meta := local_utils.BuildPaginationMeta(total, filterParams.Page, limit)
	var data any = service
	sp.logger.Infof("Successfully fetched %d services", len(service))
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
		"_id":          1,
		"service_name": 1,
		"service_code": 1,
		"service_type": 1,
		"cap":          1,
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

	sp.logger.Infof("Fetching all minimum transfer caps with filter: %+v, projection: %+v, skip: %d, limit: %d", filter, projection, skip, limit)
	service, err := sp.serviceDal.FindAllWithPagination(ctx, filter, projection, int64(skip), int64(limit))
	if err != nil {
		sp.logger.Errorf("failed to fetch service: %v", err)
		return nil, common.DefineError.General["UNHANDLED_SERVER_ERROR"]
	}

	if len(service) == 0 {
		sp.logger.Warnf("no service found for minimum transfer cap with filter: %+v", filter)
		return nil, common.DefineError.General["SERVICE_NOT_FOUND"]
	}

	total, err := sp.serviceDal.TotalCount(ctx, filter)
	if err != nil {
		sp.logger.Errorf("failed to count total services for minimum transfer cap: %v", err)
		return nil, common.DefineError.General["UNHANDLED_SERVER_ERROR"]
	}

	projectedData := ProjectDataArray[model.Service](service, projection)
	meta := local_utils.BuildPaginationMeta(total, filterParams.Page, limit)
	var data any = projectedData
	sp.logger.Infof("Successfully fetched %d minimum transfer cap services", len(service))
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
		"_id":               1,
		"service_name":      1,
		"service_key":       1,
		"service_code":      1,
		"service_type":      1,
		"cap":               1,
		"product_codes":     1,
		"ifb_product_codes": 1,
		"gl_entry":          1,
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

	sp.logger.Infof("Fetching all maximum transfer caps with filter: %+v, projection: %+v, skip: %d, limit: %d", filter, projection, skip, limit)
	service, err := sp.serviceDal.FindAllWithPagination(ctx, filter, projection, int64(skip), int64(limit))
	if err != nil {
		sp.logger.Errorf("failed to fetch service: %v", err)
		return nil, common.DefineError.General["UNHANDLED_SERVER_ERROR"]
	}
	if len(service) == 0 {
		sp.logger.Warnf("no service found for maximum transfer cap with filter: %+v", filter)
		return nil, common.DefineError.General["SERVICE_NOT_FOUND"]
	}

	total, err := sp.serviceDal.TotalCount(ctx, filter)
	if err != nil {
		sp.logger.Errorf("failed to count total services for maximum transfer cap: %v", err)
		return nil, common.DefineError.General["UNHANDLED_SERVER_ERROR"]
	}

	res := ProjectDataArray[model.Service](service, projection)

	meta := local_utils.BuildPaginationMeta(total, filterParams.Page, limit)
	var data any = res
	sp.logger.Infof("Successfully fetched %d maximum transfer cap services", len(service))
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
		"_id":               1,
		"service_name":      1,
		"service_key":       1,
		"service_type":      1,
		"payment_type":      1,
		"tiers":             1,
		"min_amount":        1,
		"created_at":        1,
		"product_codes":     1,
		"ifb_product_codes": 1,
		"gl_entry":          1,
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

	sp.logger.Infof("Fetching all service fees with filter: %+v, projection: %+v, skip: %d, limit: %d", filter, projection, skip, limit)
	service, err := sp.serviceDal.FindAllWithPagination(ctx, filter, projection, int64(skip), int64(limit))
	if err != nil {
		sp.logger.Errorf("failed to fetch service: %v", err)
		return nil, common.DefineError.General["UNHANDLED_SERVER_ERROR"]
	}

	if len(service) == 0 {
		sp.logger.Warnf("no service found for service fee with filter: %+v", filter)
		return nil, common.DefineError.General["SERVICE_NOT_FOUND"]
	}

	total, err := sp.serviceDal.TotalCount(ctx, filter)
	if err != nil {
		sp.logger.Errorf("failed to count total services for service fee: %v", err)
		return nil, common.DefineError.General["UNHANDLED_SERVER_ERROR"]
	}
	projectedData := ProjectDataArray[model.Service](service, projection)
	meta := local_utils.BuildPaginationMeta(total, filterParams.Page, limit)
	var data any = projectedData
	sp.logger.Infof("Successfully fetched %d service fee services", len(service))
	return &local_utils.PaginatedResponse[*any]{
		Data: &data,
		Meta: meta,
	}, nil
}

func (sp *servicePersistence) GetAllTotalTransferCap(ctx context.Context) (*any, error) {
	filter := bson.M{
		"is_deleted": false,
	}
	projection := bson.M{
		"_id":                  1,
		"total_cap":            1,
		"created_at":           1,
		"updated_at_total_cap": 1,
	}

	sp.logger.Infof("Fetching all total transfer caps with filter: %+v, projection: %+v, skip: %d, limit: %d", filter, projection, 0, 0)
	service, err := sp.hqDal.FindAll(ctx, filter, projection)
	if err != nil {
		sp.logger.Errorf("failed to fetch service: %v", err)
		return nil, common.DefineError.General["UNHANDLED_SERVER_ERROR"]
	}

	if len(service) == 0 {
		sp.logger.Warnf("no service found for total transfer cap with filter: %+v", filter)
		return nil, common.DefineError.General["SERVICE_NOT_FOUND"]
	}

	projectedData := ProjectDataArray[model.HQ](service, projection)
	var data any = projectedData[0]
	sp.logger.Infof("Successfully fetched %d total transfer cap services", len(service))
	return &data, nil
}

func (sp *servicePersistence) GetServiceFeeDetail(ctx context.Context, id string) (*any, error) {
	objId, err := bson.ObjectIDFromHex(id)
	if err != nil {
		sp.logger.Errorf("invalid object id for service fee detail: %v", err)
		return nil, common.DefineError.General["UNHANDLED_SERVER_ERROR"]
	}
	filter := bson.M{
		"_id":        objId,
		"is_deleted": false,
	}
	projection := bson.M{}

	sp.logger.Infof("Fetching service fee detail for id: %s", id)
	serviceDetail, err := sp.serviceDal.FindOne(ctx, filter, projection)
	if err != nil {
		sp.logger.Errorf("failed to fetch service fee detail: %v", err)
		return nil, common.DefineError.General["UNHANDLED_SERVER_ERROR"]
	}
	var data any = serviceDetail
	sp.logger.Infof("Successfully fetched service fee detail for id: %s", id)
	return &data, nil
}

func (sp *servicePersistence) UpdateServiceFee(ctx context.Context, id string, req any) error {
	sp.logger.Infof("Updating service fee for id: %s", id)
	projection := bson.M{
		"tire": 1,
	}

	prev, err := sp.prevServiceData(ctx, id, projection)
	if err != nil {
		sp.logger.Errorf("error fetching previous service data for update service fee: %v", err)
		if err.Error() == "mongo: no documents in result" {
			return common.DefineError.General["NO_DOC_FOUND"]
		}
		return err
	}

	err = sp.createCpsAction(ctx, id, prev, req, string(model.RequestUpdateServiceFee))
	if err != nil {
		sp.logger.Errorf("error creating CPS action for update service fee: %v", err)
		return err
	}
	sp.logger.Infof("Successfully created CPS action for update service fee, id: %s", id)
	return nil
}

func (sp *servicePersistence) UpdateSingleMaxTransfer(ctx context.Context, id string, req any) error {
	sp.logger.Infof("Updating single max transfer for id: %s", id)
	projection := bson.M{
		"cap": 1,
	}
	prev, err := sp.prevServiceData(ctx, id, projection)
	if err != nil {
		sp.logger.Errorf("error fetching previous service data for update single max transfer: %v", err)
		return err
	}

	capData, err := local_utils.JsonUnmarshal[model.Cap](req)
	if err != nil {
		return err
	}

	if capData.CorporateDailyCap <= prev.Cap.MinAmount || capData.CorporateSingleCap <= prev.Cap.MinAmount || capData.IDailyCap <= prev.Cap.MinAmount || capData.ISingleCap <= prev.Cap.MinAmount {
		return fmt.Errorf("SINGLE_MAX_TRANSFER_CANNOT_BE_LESS_OR_EQUAL_TO_MIN_AMOUNT")
	}

	capData.MinAmount = prev.Cap.MinAmount
	err = sp.createCpsAction(ctx, id, prev, bson.M{"cap": capData}, string(model.RequestUpdateServiceSingle))
	if err != nil {
		sp.logger.Errorf("error creating CPS action for update single max transfer: %v", err)
		return err
	}

	sp.logger.Infof("Successfully created CPS action for update single max transfer, id: %s", id)
	return nil
}

func (sp *servicePersistence) UpdateTotalMaxTransferCap(ctx context.Context, id string, newTotalCap uint64) error {
	sp.logger.Infof("Updating total max transfer cap for id: %s", id)

	exceedingServices, err := sp.validateTotalCapAgainstServices(ctx, newTotalCap)
	if err != nil {
		sp.logger.Errorf("error validating total cap against services: %v", err)
		return err
	}

	if len(exceedingServices) > 0 {
		serviceNames := make([]string, len(exceedingServices))
		for i, service := range exceedingServices {
			serviceNames[i] = service.ServiceName
		}
		sp.logger.Warnf("New total cap %d is less than caps in services: %v", newTotalCap, serviceNames)
		return fmt.Errorf(" The following services exceed the new total cap: %s", strings.Join(serviceNames, ", "))
	}

	projection := bson.M{
		"total_cap": 1,
		"_id":       1,
	}
	prevHQ, err := sp.hqDal.FindAll(ctx, nil, projection)
	if err != nil {
		sp.logger.Errorf("error fetching HQ data for update total max transfer cap: %v", err)
		return err
	}

	req := map[string]interface{}{"totaltransferlimit": newTotalCap}
	err = sp.createCpsAction(ctx, id, prevHQ[0], req, string(model.RequestUpdateServiceTotal))
	if err != nil {
		sp.logger.Errorf("error creating CPS action for update total max transfer cap: %v", err)
		return err
	}

	sp.logger.Infof("Successfully created CPS action for update total max transfer cap, id: %s", id)
	return nil
}

func (sp *servicePersistence) UpdateMinimumTransferCap(ctx context.Context, id string, req any) error {
	sp.logger.Infof("Updating minimum transfer cap for id: %s", id)
	projection := bson.M{
		"cap": 1,
	}
	prev, err := sp.prevServiceData(ctx, id, projection)
	if err != nil {
		sp.logger.Errorf("error fetching previous service data for update minimum transfer cap: %v", err)
		return err
	}

	// Extract the new minimum amount from the request
	res, err := local_utils.JsonUnmarshal[map[string]interface{}](req)
	if err != nil {
		return err
	}

	// Get the new minimum amount value
	newMinAmount, ok := (*res)["min_amount"]
	if !ok {
		return fmt.Errorf("min_amount field is required")
	}

	// Type assertion for the minimum amount
	minAmount, ok := newMinAmount.(uint64)
	if !ok {
		// Try to convert from float64 if it comes as JSON number
		if floatVal, ok := newMinAmount.(float64); ok {
			minAmount = uint64(floatVal)
		} else {
			return fmt.Errorf("min_amount must be a valid number")
		}
	}

	// Create a complete updated cap structure by copying the existing cap and updating only min_amount
	updatedCap := model.Cap{
		ISingleCap:         prev.Cap.ISingleCap,
		IDailyCap:          prev.Cap.IDailyCap,
		CorporateSingleCap: prev.Cap.CorporateSingleCap,
		CorporateDailyCap:  prev.Cap.CorporateDailyCap,
		MinAmount:          minAmount,
	}

	// Validate that the new minimum amount doesn't exceed any of the existing caps
	if updatedCap.MinAmount >= updatedCap.ISingleCap {
		return fmt.Errorf("min_amount cannot be greater than or equal to individual_single_cap")
	}
	if updatedCap.MinAmount >= updatedCap.IDailyCap {
		return fmt.Errorf("min_amount cannot be greater than or equal to individual_daily_cap")
	}
	if updatedCap.MinAmount >= updatedCap.CorporateSingleCap {
		return fmt.Errorf("min_amount cannot be greater than or equal to corporate_single_cap")
	}
	if updatedCap.MinAmount >= updatedCap.CorporateDailyCap {
		return fmt.Errorf("min_amount cannot be greater than or equal to corporate_daily_cap")
	}

	err = sp.createCpsAction(ctx, id, prev, bson.M{"cap": updatedCap}, string(model.RequestUpdateServiceMinCap))
	if err != nil {
		sp.logger.Errorf("error creating CPS action for update minimum transfer cap: %v", err)
		return err
	}

	sp.logger.Infof("Successfully created CPS action for update minimum transfer cap, id: %s", id)
	return nil
}

func (sp *servicePersistence) DeleteServiceFeeTire(ctx context.Context, id string) error {
	sp.logger.Infof("Deleting service fee tire for id: %s", id)
	projection := bson.M{
		"tier": 1,
	}
	prev, err := sp.prevServiceData(ctx, id, projection)
	if err != nil {
		sp.logger.Errorf("error fetching previous service data for delete service fee tire: %v", err)
		return err
	}
	err = sp.createCpsAction(ctx, id, prev, nil, string(model.RequestDeleteServiceFee))
	if err != nil {
		sp.logger.Errorf("error creating CPS action for delete service fee tire: %v", err)
		return err
	}

	sp.logger.Infof("Successfully created CPS action for delete service fee tire, id: %s", id)
	return nil
}

func (sp *servicePersistence) createCpsAction(ctx context.Context, id string, previousAction any, currentAction any, requestAction string) error {
	userData := contexts.ExtractContext(ctx)

	cpsAction := &model.CPSAction{
		ID:               bson.NewObjectID(),
		ActionCode:       local_utils.GenerateUniqueActionCode(20),
		UniqueId:         id,
		MakerID:          userData.UserID,
		MakerName:        userData.FullName,
		MakerPhoneNumber: userData.PhoneNumber,
		Department:       userData.Department,
		ActionStatus:     "PENDING",
		CurrentAction:    currentAction,
		PreviousAction:   previousAction,
		IsDeleted:        false,
		CreatedAt:        time.Now(),
		LastModifiedAt:   time.Now(),
		RequestAction:    requestAction,
		MakerActionTime:  time.Now(),
	}

	sp.logger.Infof("Creating CPS action: %+v", cpsAction)
	_, err := sp.cpsDal.InsertOne(ctx, *cpsAction)
	if err != nil {
		sp.logger.Errorf("failed to insert CPS action: %v", err)
		return common.DefineError.General["UNHANDLED_SERVER_ERROR"]
	}

	sp.logger.Infof("Successfully created CPS action for id: %s", id)
	return nil
}

func (sp *servicePersistence) prevServiceData(ctx context.Context, id string, projection bson.M) (*model.Service, error) {
	objId, err := bson.ObjectIDFromHex(id)
	if err != nil {
		sp.logger.Errorf("invalid object id for prevServiceData: %v", err)
		return nil, err
	}

	filter := bson.M{
		"_id":        objId,
		"is_deleted": false,
	}
	sp.logger.Infof("Fetching previous service data for id: %s, projection: %+v", id, projection)
	prev, err := sp.serviceDal.FindOne(ctx, filter, projection)
	if err != nil {
		sp.logger.Errorf("error fetching previous service data: %v", err)
		return nil, err
	}

	sp.logger.Infof("Successfully fetched previous service data for id: %s", id)
	return prev, nil
}

func FindTierByID(tiers []model.Tier, id string) (*model.Tier, error) {
	objId, err := bson.ObjectIDFromHex(id)
	if err != nil {
		fmt.Printf("FindTierByID: invalid object id: %v\n", err)
		return nil, err
	}
	for i := range tiers {
		if tiers[i].ID == objId {
			fmt.Printf("FindTierByID: found tier with id: %s\n", id)
			return &tiers[i], nil
		}
	}
	fmt.Printf("FindTierByID: tier not found for id: %s\n", id)
	return nil, common.DefineError.General["INVALID_INPUT_PARAMETERS"]
}

func DeleteTierByID(tiers []model.Tier, tierID string) ([]model.Tier, error) {
	objId, err := bson.ObjectIDFromHex(tierID)
	if err != nil {
		fmt.Printf("DeleteTierByID: invalid object id: %v\n", err)
		return tiers, err
	}
	for i, t := range tiers {
		if t.ID == objId {
			fmt.Printf("DeleteTierByID: deleting tier with id: %s\n", tierID)
			return append(tiers[:i], tiers[i+1:]...), nil
		}
	}
	fmt.Printf("DeleteTierByID: tier not found for id: %s\n", tierID)
	return tiers, common.DefineError.General["INVALID_INPUT_PARAMETERS"]
}

func AddTier(tiers []model.Tier, newTier model.Tier) []model.Tier {
	fmt.Printf("AddTier: adding new tier: %+v\n", newTier)
	return append(tiers, newTier)
}

func ProjectDataArray[T any](data []*T, projection bson.M) []map[string]interface{} {
	// Parse projection into a tree structure to support subfields like "cap.min_amount"
	type projNode struct {
		Include  bool
		Children map[string]*projNode
	}
	root := &projNode{Children: make(map[string]*projNode)}

	for k, v := range projection {
		inc := false
		switch v := v.(type) {
		case int:
			inc = v == 1
		case bool:
			inc = v
		default:
			continue
		}

		parts := strings.Split(k, ".")
		node := root
		for i, part := range parts {
			if node.Children == nil {
				node.Children = make(map[string]*projNode)
			}
			if _, ok := node.Children[part]; !ok {
				node.Children[part] = &projNode{Children: make(map[string]*projNode)}
			}
			node = node.Children[part]
			if i == len(parts)-1 {
				node.Include = inc
			}
		}
	}

	// Helper to project a struct or map according to the projection tree
	var projectValue func(val reflect.Value, node *projNode) interface{}
	projectValue = func(val reflect.Value, node *projNode) interface{} {
		if !val.IsValid() {
			return nil
		}
		if val.Kind() == reflect.Ptr {
			if val.IsNil() {
				return nil
			}
			val = val.Elem()
		}
		switch val.Kind() {
		case reflect.Struct:
			typ := val.Type()
			result := make(map[string]interface{})
			for i := 0; i < typ.NumField(); i++ {
				field := typ.Field(i)
				if field.PkgPath != "" { // unexported
					continue
				}
				jsonTag := field.Tag.Get("json")
				if jsonTag == "" {
					jsonTag = field.Name
				} else if idx := strings.Index(jsonTag, ","); idx != -1 {
					jsonTag = jsonTag[:idx]
				}
				child, hasChild := node.Children[jsonTag]
				if hasChild {
					if len(child.Children) > 0 {
						valField := val.Field(i)
						projected := projectValue(valField, child)
						if projected != nil {
							result[jsonTag] = projected
						}
					} else if child.Include {
						result[jsonTag] = val.Field(i).Interface()
					}
				} else if node.Include && len(node.Children) == 0 {
					// If this node is a leaf include, include all fields
					result[jsonTag] = val.Field(i).Interface()
				}
			}
			if len(result) > 0 {
				return result
			}
			return nil
		case reflect.Map:
			result := make(map[string]interface{})
			iter := val.MapRange()
			for iter.Next() {
				key := iter.Key()
				strKey := ""
				if key.Kind() == reflect.String {
					strKey = key.String()
				} else {
					continue
				}
				child, hasChild := node.Children[strKey]
				if hasChild {
					if len(child.Children) > 0 {
						projected := projectValue(iter.Value(), child)
						if projected != nil {
							result[strKey] = projected
						}
					} else if child.Include {
						result[strKey] = iter.Value().Interface()
					}
				} else if node.Include && len(node.Children) == 0 {
					result[strKey] = iter.Value().Interface()
				}
			}
			if len(result) > 0 {
				return result
			}
			return nil
		default:
			if node.Include {
				return val.Interface()
			}
			return nil
		}
	}

	res := make([]map[string]interface{}, len(data))
	for i, item := range data {
		val := reflect.ValueOf(item)
		projected := projectValue(val, root)
		if m, ok := projected.(map[string]interface{}); ok {
			res[i] = m
		} else {
			res[i] = make(map[string]interface{})
		}
	}
	return res
}

func (sp *servicePersistence) ValidateMaxTotalCap(ctx context.Context, cap model.Cap) (bool, error) {
	sp.logger.Infof("Validating max total cap: %+v", cap)
	data, err := sp.hqDal.FindAll(ctx, nil, nil)
	if err != nil {
		sp.logger.Errorf("error fetching HQ data for ValidateMaxTotalCap: %v", err)
		return false, err
	}

	if len(data) == 0 {
		sp.logger.Errorf("NO_HQ_DATA_FOUND in ValidateMaxTotalCap")
		return false, fmt.Errorf("NO_HQ_DATA_FOUND")
	}

	hqData := data[0]
	totalCap := hqData.TotalCap

	totalCapUint := totalCap

	if cap.ISingleCap > uint64(totalCapUint) {
		sp.logger.Warnf("INDIVIDUAL_SINGLE_CAP_EXCEEDS_TOTAL_CAP: %d > %d", cap.ISingleCap, totalCapUint)
		return false, fmt.Errorf("INDIVIDUAL_SINGLE_CAP_EXCEEDS_TOTAL_CAP")
	}
	if cap.IDailyCap > uint64(totalCapUint) {
		sp.logger.Warnf("INDIVIDUAL_DAILY_CAP_EXCEEDS_TOTAL_CAP: %d > %d", cap.IDailyCap, totalCapUint)
		return false, fmt.Errorf("INDIVIDUAL_DAILY_CAP_EXCEEDS_TOTAL_CAP")
	}
	if cap.CorporateSingleCap > uint64(totalCapUint) {
		sp.logger.Warnf("CORPORATE_SINGLE_CAP_EXCEEDS_TOTAL_CAP: %d > %d", cap.CorporateSingleCap, totalCapUint)
		return false, fmt.Errorf("CORPORATE_SINGLE_CAP_EXCEEDS_TOTAL_CAP")
	}
	if cap.CorporateDailyCap > uint64(totalCapUint) {
		sp.logger.Warnf("CORPORATE_DAILY_CAP_EXCEEDS_TOTAL_CAP: %d > %d", cap.CorporateDailyCap, totalCapUint)
		return false, fmt.Errorf("CORPORATE_DAILY_CAP_EXCEEDS_TOTAL_CAP")
	}

	sp.logger.Infof("Max total cap validated successfully")
	return true, nil
}

func (sp *servicePersistence) Authorize(ctx context.Context, cpsAction any) (any, error) {
	action, err := local_utils.JsonUnmarshal[model.CPSAction](cpsAction)
	if err != nil {
		sp.logger.Errorf("invalid cpsAction type or nil")
		return nil, fmt.Errorf("INVALID_CPS_ACTION")
	}

	castToBsonM := func(input interface{}) (bson.M, error) {
		raw, err := bson.Marshal(input)
		if err != nil {
			return nil, err
		}
		var out bson.M
		err = bson.Unmarshal(raw, &out)
		return out, err
	}

	switch action.RequestAction {
	case string(model.RequestUpdateServiceFee), string(model.RequestUpdateServiceSingle), string(model.RequestUpdateServiceMinCap):
		current, err := castToBsonM(action.CurrentAction)
		if err != nil {
			sp.logger.Errorf("invalid currentAction format for update: %v", err)
			return nil, fmt.Errorf("INVALID_CURRENT_ACTION_FORMAT")
		}
		update := bson.M{}
		for k, v := range current {
			switch val := v.(type) {
			case string:
				if strings.TrimSpace(val) != "" {
					update[k] = v
				}
			case []interface{}:
				// Only include non-empty arrays
				if len(val) > 0 {
					update[k] = v
				}
			case nil:
				// skip
			default:
				update[k] = v
			}
		}

		if len(update) == 0 {
			sp.logger.Errorf("empty update body for service update")
			return nil, fmt.Errorf("EMPTY_UPDATE_BODY")
		}

		update["last_modified_at"] = time.Now()
		prev, err := castToBsonM(action.PreviousAction)

		if err != nil {
			sp.logger.Errorf("invalid previousAction format for update: %v", err)
			return nil, fmt.Errorf("INVALID_PREVIOUS_ACTION_FORMAT")
		}

		id, ok := prev["_id"].(string)
		if !ok {
			sp.logger.Errorf("missing _id in previousAction for update")
			return nil, fmt.Errorf("MISSING_ID")
		}
		objId, err := bson.ObjectIDFromHex(id)
		if err != nil {
			return nil, fmt.Errorf("UNHANDLED_SERVER_ERROR")
		}
		filter := bson.M{"_id": objId, "is_deleted": false}

		_, err = sp.serviceDal.UpdateOne(ctx, filter, update)
		fmt.Println(err)
		if err != nil {
			sp.logger.Errorf("failed to approve service update: %v", err)
			return nil, fmt.Errorf("SERVICE_UPDATE_FAILED")
		}

	case string(model.RequestDeleteServiceFee):
		// Delete logic (soft delete, and also remove value from array if needed)
		current, err := local_utils.JsonUnmarshal[bson.M](action.CurrentAction)
		if err != nil {
			sp.logger.Errorf("invalid currentAction format for delete: %v", err)
			return nil, fmt.Errorf("INVALID_CURRENT_ACTION_FORMAT")
		}
		currentData := *current
		id, ok := currentData["_id"]
		if !ok {
			sp.logger.Errorf("missing _id in currentAction for delete")
			return nil, fmt.Errorf("MISSING_ID")
		}
		filter := bson.M{"_id": id, "is_deleted": false}
		update := bson.M{
			"is_deleted":       true,
			"last_modified_at": time.Now(),
		}

		if remove, ok := currentData["remove_from_array"].(map[string]interface{}); ok {
			field, hasField := remove["field"].(string)
			value, hasValue := remove["value"]
			if hasField && hasValue {
				update["$pull"] = bson.M{field: value}
			}
		}
		_, err = sp.serviceDal.UpdateOne(ctx, filter, update)
		if err != nil {
			sp.logger.Errorf("failed to authorize service fee delete: %v", err)
			return nil, fmt.Errorf("FAILED_TO_AUTHORIZE_DELETE")
		}
	case string(model.RequestUpdateServiceTotal):
		sp.logger.Infof("Authorizing total cap update")
		current, err := castToBsonM(action.CurrentAction)
		if err != nil {
			sp.logger.Errorf("invalid currentAction format for update: %v", err)
			return nil, fmt.Errorf("INVALID_CURRENT_ACTION_FORMAT")
		}

		val, ok := current["totaltransferlimit"].(float64)
		if !ok {
			sp.logger.Errorf("invalid total cap format during authorization")
			return nil, fmt.Errorf("INVALID_TOTAL_CAP_FORMAT")
		}
		newTotalCap := uint64(val)

		exceedingServices, err := sp.validateTotalCapAgainstServices(ctx, newTotalCap)
		if err != nil {
			sp.logger.Errorf("error validating total cap against services during authorization: %v", err)
			return nil, fmt.Errorf("TOTAL_CAP_VALIDATION_FAILED")
		}

		if len(exceedingServices) > 0 {
			serviceNames := make([]string, len(exceedingServices))
			for i, service := range exceedingServices {
				serviceNames[i] = service.ServiceName
			}
			sp.logger.Warnf("Cannot authorize: New total cap %d is less than caps in services: %v", newTotalCap, serviceNames)
			return nil, fmt.Errorf("The following services exceed the new total cap: %s", strings.Join(serviceNames, ","))
		}

		update := bson.M{}
		update["updated_at_total_cap"] = time.Now()
		update["total_cap"] = current["totaltransferlimit"]
		prev, err := castToBsonM(action.PreviousAction)

		if err != nil {
			sp.logger.Errorf("invalid previousAction format for update: %v", err)
			return nil, fmt.Errorf("INVALID_PREVIOUS_ACTION_FORMAT")
		}

		id, ok := prev["_id"].(string)
		if !ok {
			sp.logger.Errorf("missing _id in previousAction for update")
			return nil, fmt.Errorf("MISSING_ID")
		}
		objId, err := bson.ObjectIDFromHex(id)
		if err != nil {
			return nil, fmt.Errorf("UNHANDLED_SERVER_ERROR")
		}
		filter := bson.M{"_id": objId, "is_deleted": false}

		_, err = sp.hqDal.UpdateOne(ctx, filter, update)
		if err != nil {
			sp.logger.Errorf("error approving service update: %v", err)
			return nil, fmt.Errorf("SERVICE_UPDATE_FAILED")
		}

		sp.logger.Infof("Successfully authorized total cap update to %d", newTotalCap)
	default:
		sp.logger.Errorf("failed to authorize action: unknown request action %s", action.RequestAction)
		return nil, fmt.Errorf("FAILED_TO_AUTHORIZE")
	}

	// Update CPSAction status to approved
	filter := bson.M{"_id": action.ID}
	update := bson.M{
		"action_status":        model.ActionApproved,
		"checker_action_time":  action.CheckerActionTime,
		"checker_id":           action.CheckerID,
		"checker_name":         action.CheckerName,
		"checker_phone_number": action.CheckerPhoneNumber,
		"last_modified_at":     time.Now(),
	}
	_, err = sp.cpsDal.UpdateOne(ctx, filter, update)
	if err != nil {
		sp.logger.Errorf("failed to update CPSAction status: %v", err)
		return nil, fmt.Errorf("CPSACTION_STATUS_UPDATE_FAILED")
	}

	return action, nil
}

// validateTotalCapAgainstServices checks if the new total cap is less than any service caps
func (sp *servicePersistence) validateTotalCapAgainstServices(ctx context.Context, newTotalCap uint64) ([]*model.Service, error) {
	sp.logger.Infof("Validating new total cap %d against all service caps", newTotalCap)

	// Get all services with their caps
	filter := bson.M{
		"is_deleted": false,
	}
	projection := bson.M{
		"service_name": 1,
		"service_code": 1,
		"cap":          1,
	}

	services, err := sp.serviceDal.FindAll(ctx, filter, projection)
	if err != nil {
		sp.logger.Errorf("error fetching services for total cap validation: %v", err)
		return nil, err
	}

	var exceedingServices []*model.Service

	for _, service := range services {
		// Check individual caps
		if service.Cap.ISingleCap > newTotalCap {
			sp.logger.Warnf("Service %s (ID: %s) has individual_single_cap %d which exceeds new total cap %d",
				service.ServiceName, service.ID.Hex(), service.Cap.ISingleCap, newTotalCap)
			exceedingServices = append(exceedingServices, service)
			continue
		}

		if service.Cap.IDailyCap > newTotalCap {
			sp.logger.Warnf("Service %s (ID: %s) has individual_daily_cap %d which exceeds new total cap %d",
				service.ServiceName, service.ID.Hex(), service.Cap.IDailyCap, newTotalCap)
			exceedingServices = append(exceedingServices, service)
			continue
		}

		// Check corporate caps
		if service.Cap.CorporateSingleCap > newTotalCap {
			sp.logger.Warnf("Service %s (ID: %s) has corporate_single_cap %d which exceeds new total cap %d",
				service.ServiceName, service.ID.Hex(), service.Cap.CorporateSingleCap, newTotalCap)
			exceedingServices = append(exceedingServices, service)
			continue
		}

		if service.Cap.CorporateDailyCap > newTotalCap {
			sp.logger.Warnf("Service %s (ID: %s) has corporate_daily_cap %d which exceeds new total cap %d",
				service.ServiceName, service.ID.Hex(), service.Cap.CorporateDailyCap, newTotalCap)
			exceedingServices = append(exceedingServices, service)
			continue
		}
	}

	if len(exceedingServices) > 0 {
		sp.logger.Warnf("Found %d services with caps exceeding new total cap %d", len(exceedingServices), newTotalCap)
	} else {
		sp.logger.Infof("All service caps are within the new total cap %d", newTotalCap)
	}

	return exceedingServices, nil
}
