package service

// import (
// 	"context"
// 	"fmt"
// 	"net/http"
// 	"time"

// 	"gitlab.com/bersufekadgetachew/cbe-super-app-cps-action/internal/domain/service"
// 	// "gitlab.com/bersufekadgetachew/cbe-super-app-cps-action/internal/domain/service/entity"
// 	outbound "gitlab.com/bersufekadgetachew/cbe-super-app-cps-action/internal/port/outbound/service_details"

// 	constant "gitlab.com/bersufekadgetachew/cbe-super-app-cps-action/utils"

// 	"gitlab.com/bersufekadgetachew/cbe-super-app-shared/shared/dal"
// 	// "gitlab.com/bersufekadgetachew/cbe-super-app-shared/shared/entities/member"
// 	"gitlab.com/bersufekadgetachew/cbe-super-app-shared/shared/utils"
// 	"go.mongodb.org/mongo-driver/v2/bson"
// 	"go.mongodb.org/mongo-driver/v2/mongo"
// )

// type serviceFeeRepo struct {
// 	client        *mongo.Client
// 	logger        utils.Logger
// 	cpsDal        dal.MongoDal[service.CPSAction, service.CPSAction]
// 	serviceFeeDal dal.MongoDal[service.Service, service.Service]
// }

// var _ outbound.ServiceFee = (*serviceFeeRepo)(nil)

// func InitFaydaAccountPersistence(client *mongo.Client, database string,
// 	cpsCollection []string, logger utils.Logger) *serviceFeeRepo {
// 	cpsDal := dal.NewMongoDal[service.CPSAction, service.CPSAction](client, database, cpsCollection[0])
// 	serviceFeeDal := dal.NewMongoDal[service.Service, service.Service](client, database, cpsCollection[1])
// 	return &serviceFeeRepo{
// 		client:        client,
// 		logger:        logger,
// 		cpsDal:        cpsDal,
// 		serviceFeeDal: serviceFeeDal,
// 	}
// }

// func (s *serviceFeeRepo) InitiateServiceFeeUpdate(ctx context.Context, req service.CPSAction) (service.UpdateServiceDetailsResponse, error) {
// 	filter := bson.M{
// 		"maker_user.phone_number": req.MakerUser.PhoneNumber,
// 		"status":                  "PENDING",
// 		"department":              req.Department,
// 	}

// 	projection := bson.M{}

// 	FindAction, err := s.cpsDal.FindOne(ctx, filter, projection)
// 	if err != nil && err != mongo.ErrNoDocuments {
// 		s.logger.Errorf("failed to get cps Action", err)
// 		err = fmt.Errorf("failed to get get cps action %w", constant.ErrorDefinition{
// 			Code:    http.StatusInternalServerError,
// 			Message: "internal server error",
// 		})
// 		return service.UpdateServiceDetailsResponse{}, err
// 	}

// 	if FindAction != nil {
// 		s.logger.Infof("pending cps action present", req.MakerUser.FullName, req.MakerUser.UserCode, req.Department)
// 		err = fmt.Errorf("failed to get cps action  %w", constant.ErrorDefinition{
// 			Code:    http.StatusBadRequest,
// 			Message: "pending cps action present",
// 		})
// 		return service.UpdateServiceDetailsResponse{}, err
// 	}

// 	serviceFilter := bson.M{
// 		"action_code": req.ActionCode,
// 	}

// 	serviceProjection := bson.M{}

// 	var serviceFee *service.Service
// 	serviceFee, err = s.serviceFeeDal.FindOne(ctx, serviceFilter, serviceProjection)
// 	if err != nil {
// 		s.logger.Errorf("failed to get service", err)
// 		err = fmt.Errorf("failed to get service %w", constant.ErrorDefinition{
// 			Code:    http.StatusInternalServerError,
// 			Message: "internal server error",
// 		})
// 		return service.UpdateServiceDetailsResponse{}, err
// 	}
// 	if serviceFee == nil {
// 		s.logger.Errorf("service not found")
// 		err = fmt.Errorf("service not found: %w", constant.ErrorDefinition{
// 			Code:    http.StatusNotFound,
// 			Message: "service not found",
// 		})
// 		return service.UpdateServiceDetailsResponse{}, err
// 	}

// 	req.Status = "PENDING"
// 	req.RequestAction = "UPDATE"
// 	req.ActionType = "UPDATE_SERVICE_FEE"
// 	req.PreviousData = map[string]any{

// 		"tiers": req.PreviousData,
// 	}
// 	req.CurrentData = map[string]any{
// 		"is_account_blocked": true,
// 	}

// 	req.MakerActionTime = time.Now()
// 	cpsAction, err := s.cpsDal.InsertOne(ctx, req)
// 	if err != nil {
// 		s.logger.Errorf("failed to create cps action", err)
// 		err = fmt.Errorf("failed to create cps action %w", constant.ErrorDefinition{
// 			Code:    http.StatusInternalServerError,
// 			Message: "internal server error",
// 		})
// 		return service.UpdateServiceDetailsResponse{}, err
// 	}

// 	return service.UpdateServiceDetailsResponse{
// 		ActionID: cpsAction.ActionCode,
// 	}, nil

// }

// func (s *serviceFeeRepo) ApproveServiceFeeUpdate(ctx context.Context, cpsAction service.CPSAction) (service.UpdateServiceDetailsResponse, error) {
// 	filter := bson.M{
// 		"action_code": cpsAction.ActionCode,
// 		"status":      "PENDING",
// 	}

// 	projection := bson.M{}
// 	cpsActionPtr, err := s.cpsDal.FindOne(ctx, filter, projection)
// 	if err != nil {
// 		s.logger.Errorf("failed to find cpsActon", err)
// 		err = fmt.Errorf("failed to findCpsAction %w", constant.ErrorDefinition{
// 			Code:    http.StatusInternalServerError,
// 			Message: "internal server error",
// 		})
// 		return service.UpdateServiceDetailsResponse{}, err
// 	}
// 	// Apply the update to the actual service
// 	serviceFilter := bson.M{
// 		"action_code": cpsActionPtr.ActionCode,
// 	}
// 	serviceUpdate := bson.M{
// 		"set": cpsActionPtr.CurrentData,
// 	}
// 	_, err = s.serviceFeeDal.UpdateOne(ctx, serviceFilter, serviceUpdate)
// 	if err != nil {
// 		s.logger.Errorf("failed to update service fee", err)
// 		err = fmt.Errorf("failed to update service fee %w", constant.ErrorDefinition{
// 			Code:    http.StatusInternalServerError,
// 			Message: "internal server error",
// 		})
// 		return service.UpdateServiceDetailsResponse{}, err
// 	}
// 	return service.UpdateServiceDetailsResponse{
// 		ActionID: cpsActionPtr.ActionCode,
// 	}, nil
// }

// func (s *serviceFeeRepo) RejectServiceFeeUpdate(ctx context.Context, cpsAction service.CPSAction) (service.UpdateServiceDetailsResponse, error) {
// 	filter := bson.M{
// 		"action_code": cpsAction.ActionCode,
// 		"status":      "PENDING",
// 	}

// 	update := bson.M{
// 		"checker_user": bson.M{
// 			"full_name":    cpsAction.CheckerUser.FullName,
// 			"phone_number": cpsAction.CheckerUser.PhoneNumber,
// 			"user_code":    cpsAction.CheckerUser.UserCode,
// 		},
// 		"status":                 "REJECTED",
// 		"rejected_action_reason": cpsAction.RejectedReason,
// 		"checker_action_time":    time.Now(),
// 	}

// 	_, err := s.cpsDal.UpdateOne(ctx, filter, update)
// 	if err != nil {
// 		s.logger.Errorf("failed to reject service fee update", err)
// 		err = fmt.Errorf("failed to reject service fee update %w", constant.ErrorDefinition{
// 			Code:    http.StatusInternalServerError,
// 			Message: "internal server error",
// 		})
// 		return service.UpdateServiceDetailsResponse{}, err
// 	}
// 	return service.UpdateServiceDetailsResponse{
// 		ActionID: cpsAction.ActionCode,
// 	}, nil
// }
