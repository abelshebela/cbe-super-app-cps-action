package service

import (
	"context"
	"errors"
	"time"

	"github.com/CBE-Super-App/cbe-super-app-cps-action/internal/adapter/outbound/model"
	"github.com/CBE-Super-App/cbe-super-app-cps-action/internal/domain/service"

	"gitlab.com/bersufekadgetachew/cbe-super-app-shared/shared/dal"
	"gitlab.com/bersufekadgetachew/cbe-super-app-shared/shared/utils"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/v2/bson"
	"go.mongodb.org/mongo-driver/v2/mongo"
)

type ServiceFeePersistence struct {
	client        *mongo.Client
	cpsDal        dal.MongoDal[service.CPSAction, service.CPSAction]
	serviceFeeDal dal.MongoDal[service.Service, service.Service]
	logger        utils.Logger
}

// var _ outbound.OutboundServiceDetailInfra = (*ServiceFeePersistence)(nil)

func NewServiceFeePersistence(client *mongo.Client, dbName string, collections []string, logger utils.Logger) *ServiceFeePersistence {
	return &ServiceFeePersistence{
		client:        client,
		cpsDal:        dal.NewMongoDal[service.CPSAction, service.CPSAction](client, dbName, collections[0]),
		serviceFeeDal: dal.NewMongoDal[service.Service, service.Service](client, dbName, collections[1]),
		logger:        logger,
	}
}

func (s *ServiceFeePersistence) CreateService(ctx context.Context, req service.CPSAction) (*service.CPSAction, error) {
	cpsAction := service.CPSAction{
		ID:               primitive.NewObjectID().Hex(),
		ActionCode:       utils.RandomGenerator(20),
		MakerID:          req.MakerID,
		MakerName:        req.MakerName,
		MakerPhoneNumber: req.MakerPhoneNumber,
		Department:       req.Department,
		ActionStatus:     service.ActionPending,
		RequestAction:    service.RequestServiceFeeCreate,
		ActionType:       service.ActionCreate,
		CurrentAction:    req.CurrentAction,
		MakerActionTime:  time.Now(),
		CreatedAt:        time.Now(),
		LastModifiedAt:   time.Now(),
	}
	inserted, err := s.cpsDal.InsertOne(ctx, cpsAction)
	if err != nil {
		s.logger.Errorf("failed to create cps action: %v", err)
		return nil, err
	}
	return &inserted, nil
}

func (s *ServiceFeePersistence) DeleteService(ctx context.Context, id string, req service.CPSAction) (*service.CPSAction, error) {
	// Find the service to get previous data
	filter := bson.M{"id": id, "is_deleted": false}
	projection := bson.M{"service_code": 1, "service_name": 1, "service_type": 1}
	serviceObj, err := s.serviceFeeDal.FindOne(ctx, filter, projection)
	if err != nil {
		s.logger.Errorf("failed to get service: %v", err)
		return nil, err
	}
	cpsAction := service.CPSAction{
		ID:               "", // set by Mongo if needed
		ActionCode:       utils.RandomGenerator(20),
		MakerID:          req.MakerID,
		MakerName:        req.MakerName,
		MakerPhoneNumber: req.MakerPhoneNumber,
		Department:       req.Department,
		ActionStatus:     service.ActionStatus(service.ActionPending),
		RequestAction:    service.RequestAction(model.RequestDeleteServiceFee),
		ActionType:       service.ActionType(service.ActionDelete),
		CurrentAction:    req.CurrentAction,
		PreviosAction:    serviceObj,
		MakerActionTime:  time.Now(),
		CreatedAt:        time.Now(),
		LastModifiedAt:   time.Now(),
	}
	inserted, err := s.cpsDal.InsertOne(ctx, cpsAction)
	if err != nil {
		s.logger.Errorf("failed to create cps action: %v", err)
		return nil, err
	}
	return &inserted, nil
}

func (s *ServiceFeePersistence) InitiateServiceFeeUpdate(ctx context.Context, req service.CPSAction) (service.UpdateServiceDetailsResponse, error) {
	filter := bson.M{
		"maker_phone_number": req.MakerPhoneNumber,
		"action_status":      service.ActionPending,
		"department":         req.Department,
	}
	projection := bson.M{}
	existing, err := s.cpsDal.FindOne(ctx, filter, projection)
	if err != nil && err != mongo.ErrNoDocuments {
		s.logger.Errorf("failed to get cps action: %v", err)
		return service.UpdateServiceDetailsResponse{}, err
	}
	if existing != nil {
		return service.UpdateServiceDetailsResponse{}, errors.New("pending cps action present")
	}

	req.ActionStatus = service.ActionPending
	req.RequestAction = service.RequestServiceFeeUpdate
	req.ActionType = service.ActionUpdate
	req.MakerActionTime = time.Now()
	req.CreatedAt = time.Now()
	req.LastModifiedAt = time.Now()

	cpsAction, err := s.cpsDal.InsertOne(ctx, req)
	if err != nil {
		s.logger.Errorf("failed to create cps action: %v", err)
		return service.UpdateServiceDetailsResponse{}, err
	}
	return service.UpdateServiceDetailsResponse{ActionID: cpsAction.ActionCode}, nil
}

func (s *ServiceFeePersistence) CreateServiceAction(ctx context.Context, req service.CPSAction) (service.UpdateServiceDetailsResponse, error) {
	// Check for pending create action
	filter := bson.M{
		"maker_phone_number": req.MakerPhoneNumber,
		"action_status":      service.ActionPending,
		"department":         req.Department,
		"request_action":     "CREATE_SERVICE",
	}
	projection := bson.M{}
	existing, err := s.cpsDal.FindOne(ctx, filter, projection)
	if err != nil && err != mongo.ErrNoDocuments {
		s.logger.Errorf("failed to get cps action: %v", err)
		return service.UpdateServiceDetailsResponse{}, err
	}
	if existing != nil {
		return service.UpdateServiceDetailsResponse{}, errors.New("pending create service action present")
	}

	req.ActionStatus = service.ActionPending
	req.RequestAction = service.RequestServiceFeeCreate
	req.ActionType = service.ActionCreate
	req.MakerActionTime = time.Now()
	req.CreatedAt = time.Now()
	req.LastModifiedAt = time.Now()

	cpsAction, err := s.cpsDal.InsertOne(ctx, req)
	if err != nil {
		s.logger.Errorf("failed to create cps action: %v", err)
		return service.UpdateServiceDetailsResponse{}, err
	}
	return service.UpdateServiceDetailsResponse{ActionID: cpsAction.ActionCode}, nil
}

func (s *ServiceFeePersistence) DeleteServiceAction(ctx context.Context, req service.CPSAction) (service.UpdateServiceDetailsResponse, error) {
	filter := bson.M{
		"maker_phone_number": req.MakerPhoneNumber,
		"action_status":      service.ActionPending,
		"department":         req.Department,
		"request_action":     "DELETE_SERVICE",
	}
	projection := bson.M{}
	existing, err := s.cpsDal.FindOne(ctx, filter, projection)
	if err != nil && err != mongo.ErrNoDocuments {
		s.logger.Errorf("failed to get cps action: %v", err)
		return service.UpdateServiceDetailsResponse{}, err
	}
	if existing != nil {
		return service.UpdateServiceDetailsResponse{}, errors.New("pending delete service action present")
	}

	req.ActionStatus = service.ActionPending
	req.RequestAction = service.RequestAction(model.RequestDeleteServiceFee)
	req.ActionType = service.ActionDelete
	req.MakerActionTime = time.Now()
	req.CreatedAt = time.Now()
	req.LastModifiedAt = time.Now()

	cpsAction, err := s.cpsDal.InsertOne(ctx, req)
	if err != nil {
		s.logger.Errorf("failed to create cps action: %v", err)
		return service.UpdateServiceDetailsResponse{}, err
	}
	return service.UpdateServiceDetailsResponse{ActionID: cpsAction.ActionCode}, nil
}

func (s *ServiceFeePersistence) ApproveServiceFeeUpdate(ctx context.Context, cpsAction service.CPSAction) (service.UpdateServiceDetailsResponse, error) {
	filter := bson.M{
		"action_code":   cpsAction.ActionCode,
		"action_status": service.ActionPending,
	}
	update := bson.M{
		"action_status":        service.ActionApproved,
		"checker_id":           cpsAction.CheckerID,
		"checker_name":         cpsAction.CheckerName,
		"checker_phone_number": cpsAction.CheckerPhoneNumber,
		"checker_action_time":  time.Now(),
		"last_modified_at":     time.Now(),
	}
	_, err := s.cpsDal.UpdateOne(ctx, filter, update)
	if err != nil {
		s.logger.Errorf("failed to approve cps action: %v", err)
		return service.UpdateServiceDetailsResponse{}, err
	}

	switch cpsAction.ActionType {
	case service.ActionCreate:
		_, ok := cpsAction.CurrentAction.(map[string]interface{})
		if !ok {
			return service.UpdateServiceDetailsResponse{}, errors.New("invalid service data for create")
		}
		var svc service.Service

		_, err := s.serviceFeeDal.InsertOne(ctx, svc)
		if err != nil {
			s.logger.Errorf("failed to create service: %v", err)
			return service.UpdateServiceDetailsResponse{}, err
		}
	case service.ActionUpdate:
		serviceCode, ok := cpsAction.CurrentAction.(map[string]interface{})["service_code"].(string)
		if !ok {
			return service.UpdateServiceDetailsResponse{}, errors.New("invalid service_code in current action")
		}
		serviceFilter := bson.M{"service_code": serviceCode}
		serviceUpdate := bson.M{"$set": cpsAction.CurrentAction}
		_, err := s.serviceFeeDal.UpdateOne(ctx, serviceFilter, serviceUpdate)
		if err != nil {
			s.logger.Errorf("failed to update service fee: %v", err)
			return service.UpdateServiceDetailsResponse{}, err
		}
	case service.ActionDelete:
		serviceCode, ok := cpsAction.CurrentAction.(map[string]interface{})["service_code"].(string)
		if !ok {
			return service.UpdateServiceDetailsResponse{}, errors.New("invalid service_code in current action for delete")
		}
		serviceFilter := bson.M{"service_code": serviceCode}
		serviceUpdate := bson.M{"is_deleted": true, "deleted_at": time.Now(), "last_modified_at": time.Now()}
		_, err := s.serviceFeeDal.UpdateOne(ctx, serviceFilter, serviceUpdate)
		if err != nil {
			s.logger.Errorf("failed to delete service: %v", err)
			return service.UpdateServiceDetailsResponse{}, err
		}
	}

	return service.UpdateServiceDetailsResponse{ActionID: cpsAction.ActionCode}, nil
}

func (s *ServiceFeePersistence) RejectServiceFeeUpdate(ctx context.Context, cpsAction service.CPSAction) (service.UpdateServiceDetailsResponse, error) {
	filter := bson.M{
		"action_code":   cpsAction.ActionCode,
		"action_status": service.ActionPending,
	}
	update := bson.M{
		"action_status":        service.ActionRejected,
		"rejection_reason":     cpsAction.RejectionReason,
		"checker_id":           cpsAction.CheckerID,
		"checker_name":         cpsAction.CheckerName,
		"checker_phone_number": cpsAction.CheckerPhoneNumber,
		"checker_action_time":  time.Now(),
		"last_modified_at":     time.Now(),
	}
	_, err := s.cpsDal.UpdateOne(ctx, filter, update)
	if err != nil {
		s.logger.Errorf("failed to reject cps action: %v", err)
		return service.UpdateServiceDetailsResponse{}, err
	}
	return service.UpdateServiceDetailsResponse{ActionID: cpsAction.ActionCode}, nil
}