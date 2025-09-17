package service_details_test

// import (
// 	"context"
// 	"encoding/json"
// 	"errors"
// 	"fmt"
// 	"testing"
// 	"time"

// 	"github.com/golang/mock/gomock"
// 	"github.com/stretchr/testify/assert"
// 	"github.com/stretchr/testify/mock"
// 	"gitlab.com/bersufekadgetachew/cbe-super-app-shared/shared/utils"

// 	"github.com/CBE-Super-App/cbe-super-app-cps-action/internal/domain/action"
// 	serviceDomain "github.com/CBE-Super-App/cbe-super-app-cps-action/internal/domain/service"
// 	"github.com/CBE-Super-App/cbe-super-app-cps-action/mocks"
// 	mock_domain_action "github.com/CBE-Super-App/cbe-super-app-cps-action/mocks/domain/action"
// 	mock_domain "github.com/CBE-Super-App/cbe-super-app-cps-action/mocks/domain/service"
// )

// // NoOpLogger implements the utils.Logger interface but does nothing
// type NoOpLogger struct{}

// func (l *NoOpLogger) Infof(format string, args ...interface{})  {}
// func (l *NoOpLogger) Errorf(format string, args ...interface{}) {}
// func (l *NoOpLogger) Debugf(format string, args ...interface{}) {}
// func (l *NoOpLogger) Fatalf(format string, args ...interface{}) {}
// func (l *NoOpLogger) Warnf(format string, args ...interface{})  {}
// func (l *NoOpLogger) Sync() error                               { return nil }

// func TestGetAllServiceDetails(t *testing.T) {
// 	ctrl := gomock.NewController(t)
// 	defer ctrl.Finish()
// 	mockRepo := mock_domain.NewMockRepository(ctrl)
// 	mockActionRepo := mock_domain_action.NewMockRepository(ctrl)
// 	logger := utils.NewLogger()
// 	service := serviceDomain.NewServiceStore(mockRepo, mockActionRepo, logger)

// 	ctx := context.Background()
// 	expectedServices := []*serviceDomain.Service{
// 		{
// 			ID:                 "test-id-1",
// 			ServiceCode:        "TEST1",
// 			ServiceName:        "Test Service 1",
// 			ServiceType:        "TYPE1",
// 			Key:                "key1",
// 			Cap:                serviceDomain.Cap{},
// 			CBEProductCodes:    serviceDomain.ProductCodes{},
// 			CBEIFBProductCodes: serviceDomain.ProductCodes{},
// 			AboveAmount:        1000,
// 			AboveServiceFee:    100,
// 			PaymentType:        "CASH",
// 			Tiers:              []serviceDomain.Tier{},
// 			CBEGLEntry:         serviceDomain.GLEntry{},
// 			CBEIFBGLEntry:      serviceDomain.GLEntry{},
// 			Enabled:            true,
// 			IsDeleted:          false,
// 			CreatedAt:          time.Now(),
// 			LastModifiedAt:     time.Now(),
// 		},
// 		{
// 			ID:                 "test-id-2",
// 			ServiceCode:        "TEST2",
// 			ServiceName:        "Test Service 2",
// 			ServiceType:        "TYPE2",
// 			Key:                "key2",
// 			Cap:                serviceDomain.Cap{},
// 			CBEProductCodes:    serviceDomain.ProductCodes{},
// 			CBEIFBProductCodes: serviceDomain.ProductCodes{},
// 			AboveAmount:        2000,
// 			AboveServiceFee:    200,
// 			PaymentType:        "CASH",
// 			Tiers:              []serviceDomain.Tier{},
// 			CBEGLEntry:         serviceDomain.GLEntry{},
// 			CBEIFBGLEntry:      serviceDomain.GLEntry{},
// 			Enabled:            true,
// 			IsDeleted:          false,
// 			CreatedAt:          time.Now(),
// 			LastModifiedAt:     time.Now(),
// 		},
// 	}

// 	t.Run("successful get all", func(t *testing.T) {

// 		mockRepo.EXPECT().GetAllServiceDetails(ctx).Return(expectedServices, nil).Times(1)

// 		services, err := service.GetAllServiceDetails(ctx)
// 		assert.NoError(t, err)
// 		assert.Equal(t, expectedServices, services)
// 	})

// 	t.Run("error getting all", func(t *testing.T) {
// 		mockRepo.EXPECT().GetAllServiceDetails(ctx).Return(nil, errors.New("db error")).Times(1)

// 		services, err := service.GetAllServiceDetails(ctx)
// 		assert.Error(t, err)
// 		assert.Nil(t, services)
// 	})
// }

// func TestGetServiceDetailsByID(t *testing.T) {
// 	ctrl := gomock.NewController(t)
// 	defer ctrl.Finish()
// 	mockRepo := mock_domain.NewMockRepository(ctrl)
// 	mockActionRepo := mock_domain_action.NewMockRepository(ctrl)
// 	logger := utils.NewLogger()
// 	service := serviceDomain.NewServiceStore(mockRepo, mockActionRepo, logger)

// 	ctx := context.Background()
// 	expectedService := serviceDomain.Service{
// 		ID:                 "test-id",
// 		ServiceCode:        "TEST",
// 		ServiceName:        "Test Service",
// 		ServiceType:        "TYPE",
// 		Key:                "key",
// 		Cap:                serviceDomain.Cap{},
// 		CBEProductCodes:    serviceDomain.ProductCodes{},
// 		CBEIFBProductCodes: serviceDomain.ProductCodes{},
// 		AboveAmount:        1000,
// 		AboveServiceFee:    100,
// 		PaymentType:        "CASH",
// 		Tiers:              []serviceDomain.Tier{},
// 		CBEGLEntry:         serviceDomain.GLEntry{},
// 		CBEIFBGLEntry:      serviceDomain.GLEntry{},
// 		Enabled:            true,
// 		IsDeleted:          false,
// 		CreatedAt:          time.Now(),
// 		LastModifiedAt:     time.Now(),
// 	}

// 	t.Run("successful get", func(t *testing.T) {
// 		mockRepo.EXPECT().GetOneServiceDetail(ctx, "test-id").Return(expectedService, nil).Times(1)

// 		svc, err := service.GetServiceDetailsByID(ctx, "test-id")
// 		assert.NoError(t, err)
// 		assert.Equal(t, &expectedService, svc)
// 	})

// 	t.Run("empty id", func(t *testing.T) {
// 		mockRepo.EXPECT().GetOneServiceDetail(ctx, "").Return(serviceDomain.Service{}, fmt.Errorf("service ID cannot be empty")).Times(1)
// 		svc, err := service.GetServiceDetailsByID(ctx, "")
// 		assert.Error(t, err)
// 		assert.Nil(t, svc)
// 		assert.Equal(t, "service ID cannot be empty", err.Error())
// 	})

// 	t.Run("service not found", func(t *testing.T) {
// 		mockRepo.EXPECT().GetOneServiceDetail(ctx, "non-existent").Return(serviceDomain.Service{}, assert.AnError).Times(1)

// 		svc, err := service.GetServiceDetailsByID(ctx, "non-existent")
// 		assert.Error(t, err)
// 		assert.Nil(t, svc)
// 	})
// }

// func TestUpdateServiceDetailsRequest(t *testing.T) {
// 	ctrl := gomock.NewController(t)
// 	defer ctrl.Finish()
// 	mockRepo := mock_domain.NewMockRepository(ctrl)
// 	mockActionRepo := mock_domain_action.NewMockRepository(ctrl)
// 	logger := utils.NewLogger()
// 	service := serviceDomain.NewServiceStore(mockRepo, mockActionRepo, logger)

// 	ctx := context.Background()
// 	originalService := serviceDomain.Service{
// 		ID:                 "test-id",
// 		ServiceCode:        "TEST",
// 		ServiceName:        "Test Service",
// 		ServiceType:        "TYPE",
// 		Key:                "key",
// 		Cap:                serviceDomain.Cap{},
// 		CBEProductCodes:    serviceDomain.ProductCodes{},
// 		CBEIFBProductCodes: serviceDomain.ProductCodes{},
// 		AboveAmount:        1000,
// 		AboveServiceFee:    100,
// 		PaymentType:        "CASH",
// 		Tiers:              []serviceDomain.Tier{},
// 		CBEGLEntry:         serviceDomain.GLEntry{},
// 		CBEIFBGLEntry:      serviceDomain.GLEntry{},
// 		Enabled:            true,
// 		IsDeleted:          false,
// 		CreatedAt:          time.Now(),
// 		LastModifiedAt:     time.Now(),
// 	}

// 	updatedService := originalService
// 	updatedService.ServiceName = "Updated Service"

// 	t.Run("successful update request", func(t *testing.T) {
// 		mockRepo.EXPECT().GetOneServiceDetail(ctx, "test-id").Return(originalService, nil).Times(1)

// 		expectedAction := action.CPSAction{
// 			ActionCode:       "CPS_BNHafpP8De",
// 			MakerID:          "maker-123",
// 			MakerName:        "",
// 			MakerPhoneNumber: "",
// 			MakerActionTime:  time.Now(),
// 			UniqueId:         "test-id",
// 			Department:       "test-service",
// 			ActionType:       action.ActionUpdate,
// 			RequestAction:    action.RequestUpdateServiceDetails,
// 			ActionStatus:     action.ActionPending,
// 			RejectionReason:  nil,
// 		}

// 		//actionID := "CPS_123"
// 		mockActionRepo.EXPECT().CreateCpsAction(ctx, expectedAction).Return(expectedAction, nil).Times(1)

// 		actionID, err := service.UpdateServiceDetailsRequest(ctx, "test-id", &updatedService, "maker-123")
// 		assert.NoError(t, err)
// 		assert.NotEmpty(t, actionID)
// 		assert.Contains(t, actionID, "CPS_")
// 	})

// 	t.Run("service not found", func(t *testing.T) {
// 		mockRepo.EXPECT().GetOneServiceDetail(ctx, "non-existent").Return(serviceDomain.Service{}, assert.AnError).Times(1)

// 		actionID, err := service.UpdateServiceDetailsRequest(ctx, "non-existent", &updatedService, "maker-123")
// 		assert.Error(t, err)
// 		assert.Empty(t, actionID)
// 	})

// 	t.Run("empty maker id", func(t *testing.T) {
// 		actionID, err := service.UpdateServiceDetailsRequest(ctx, "test-id", &updatedService, "")
// 		assert.Error(t, err)
// 		assert.Empty(t, actionID)
// 		assert.Equal(t, "maker ID cannot be empty", err.Error())
// 	})
// }

// func TestUpdateServiceDetails(t *testing.T) {
// 	ctrl := gomock.NewController(t)
// 	defer ctrl.Finish()
// 	mockRepo := new(mock_domain.MockRepository)
// 	mockActionRepo := new(mocks.ActionRepository)
// 	logger := &NoOpLogger{}
// 	service := serviceDomain.NewServiceStore(mockRepo, mockActionRepo, logger)

// 	ctx := context.Background()
// 	actionID := "CPS_123"
// 	checkerID := "checker-123"

// 	serviceDetails := serviceDomain.Service{
// 		ID:                 "test-id",
// 		ServiceCode:        "TEST",
// 		ServiceName:        "Test Service",
// 		ServiceType:        "TYPE",
// 		Key:                "key",
// 		Cap:                serviceDomain.Cap{},
// 		CBEProductCodes:    serviceDomain.ProductCodes{},
// 		CBEIFBProductCodes: serviceDomain.ProductCodes{},
// 		AboveAmount:        1000,
// 		AboveServiceFee:    100,
// 		PaymentType:        "CASH",
// 		Tiers:              []serviceDomain.Tier{},
// 		CBEGLEntry:         serviceDomain.GLEntry{},
// 		CBEIFBGLEntry:      serviceDomain.GLEntry{},
// 		Enabled:            true,
// 		IsDeleted:          false,
// 		CreatedAt:          time.Now(),
// 		LastModifiedAt:     time.Now(),
// 	}

// 	type CurrentAction struct {
// 		Service serviceDomain.Service `json:"service"`
// 	}
// 	currentAction := CurrentAction{Service: serviceDetails}
// 	currentActionBytes, _ := json.Marshal(currentAction)

// 	t.Run("successful approval", func(t *testing.T) {
// 		expectedAction := action.CPSAction{
// 			ActionCode:       actionID,
// 			MakerID:          "maker-123",
// 			MakerName:        "",
// 			MakerPhoneNumber: "",
// 			MakerActionTime:  time.Now(),
// 			Department:       "test-service",
// 			ActionType:       action.ActionUpdate,
// 			RequestAction:    action.RequestUpdateServiceDetails,
// 			ActionStatus:     action.ActionPending,
// 			CurrentAction:    currentActionBytes,
// 			CreatedAt:        time.Now(),
// 			LastModifiedAt:   time.Now(),
// 			RejectionReason:  nil,
// 		}

// 		mockActionRepo.On("FetchCpsActionById", ctx, actionID).Return(expectedAction, nil).Once()
// 		mockRepo.EXPECT().UpdateOneServiceDetailRequest(ctx, "test-id", gomock.Any()).Return(nil).Times(1)
// 		mockActionRepo.On("UpdateCpsAction", ctx, mock.Anything).Return(nil).Once()

// 		err := service.UpdateServiceDetails(ctx, actionID, true, checkerID, "")
// 		assert.NoError(t, err)
// 	})

// 	t.Run("successful rejection with reason", func(t *testing.T) {
// 		expectedAction := action.CPSAction{
// 			ActionCode:      actionID,
// 			MakerID:         "maker-123",
// 			MakerActionTime: time.Now(),
// 			Department:      "test-service",
// 			ActionType:      action.ActionUpdate,
// 			RequestAction:   action.RequestUpdateServiceDetails,
// 			ActionStatus:    action.ActionPending,
// 			CurrentAction:   currentActionBytes,
// 			CreatedAt:       time.Now(),
// 			LastModifiedAt:  time.Now(),
// 			RejectionReason: nil,
// 		}

// 		mockActionRepo.On("FetchCpsActionById", ctx, actionID).Return(expectedAction, nil).Once()
// 		mockActionRepo.On("UpdateCpsAction", ctx, mock.Anything).Return(nil).Once()

// 		err := service.UpdateServiceDetails(ctx, actionID, false, checkerID, "Invalid service configuration")
// 		assert.NoError(t, err)
// 	})

// 	t.Run("successful rejection without reason", func(t *testing.T) {
// 		expectedAction := action.CPSAction{
// 			ActionCode:       actionID,
// 			MakerID:          "maker-123",
// 			MakerName:        "",
// 			MakerPhoneNumber: "",
// 			MakerActionTime:  time.Now(),
// 			Department:       "test-service",
// 			ActionType:       action.ActionUpdate,
// 			RequestAction:    action.RequestUpdateServiceDetails,
// 			ActionStatus:     action.ActionPending,
// 			CurrentAction:    currentActionBytes,
// 			CreatedAt:        time.Now(),
// 			LastModifiedAt:   time.Now(),
// 			RejectionReason:  nil,
// 		}

// 		mockActionRepo.On("FetchCpsActionById", ctx, actionID).Return(expectedAction, nil).Once()
// 		mockActionRepo.On("UpdateCpsAction", ctx, mock.Anything).Return(nil).Once()

// 		err := service.UpdateServiceDetails(ctx, actionID, false, checkerID, "")
// 		assert.NoError(t, err)
// 	})

// 	t.Run("empty checker id", func(t *testing.T) {
// 		err := service.UpdateServiceDetails(ctx, actionID, true, "", "")
// 		assert.Error(t, err)
// 		assert.Equal(t, "checker ID cannot be empty", err.Error())
// 	})

// 	t.Run("empty action id", func(t *testing.T) {
// 		err := service.UpdateServiceDetails(ctx, "", true, checkerID, "")
// 		assert.Error(t, err)
// 		assert.Equal(t, "action ID cannot be empty", err.Error())
// 	})

// 	t.Run("action not found", func(t *testing.T) {
// 		mockActionRepo.On("FetchCpsActionById", ctx, actionID).Return(action.CPSAction{}, assert.AnError).Once()

// 		err := service.UpdateServiceDetails(ctx, actionID, true, checkerID, "")
// 		assert.Error(t, err)
// 		assert.Equal(t, "action not found", err.Error())
// 	})

// 	t.Run("action not pending", func(t *testing.T) {
// 		expectedAction := action.CPSAction{
// 			ActionCode:       actionID,
// 			MakerID:          "maker-123",
// 			MakerName:        "",
// 			MakerPhoneNumber: "",
// 			MakerActionTime:  time.Now(),
// 			Department:       "test-service",
// 			ActionType:       action.ActionUpdate,
// 			RequestAction:    action.RequestUpdateServiceDetails,
// 			ActionStatus:     action.ActionApproved,
// 			CurrentAction:    currentActionBytes,
// 			CreatedAt:        time.Now(),
// 			LastModifiedAt:   time.Now(),
// 			RejectionReason:  nil,
// 		}

// 		mockActionRepo.On("FetchCpsActionById", ctx, actionID).Return(expectedAction, nil).Once()

// 		err := service.UpdateServiceDetails(ctx, actionID, true, checkerID, "")
// 		assert.Error(t, err)
// 		assert.Equal(t, "action is not pending", err.Error())
// 	})
// }
