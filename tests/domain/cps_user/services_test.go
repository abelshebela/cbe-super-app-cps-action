package cps_user_test

import (
	"context"
	"net/http"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"go.mongodb.org/mongo-driver/v2/bson"

	"github.com/CBE-Super-App/cbe-super-app-cps-action/internal/adapter/outbound/model"
	"github.com/CBE-Super-App/cbe-super-app-cps-action/internal/domain/action"
	dto "github.com/CBE-Super-App/cbe-super-app-cps-action/internal/domain/cps_user/dto"
	cpsuserServices "github.com/CBE-Super-App/cbe-super-app-cps-action/internal/domain/cps_user/services"
	"github.com/CBE-Super-App/cbe-super-app-cps-action/mocks"
	common_util "github.com/CBE-Super-App/cbe-super-app-cps-action/pkgs/utils"
	constant "github.com/CBE-Super-App/cbe-super-app-cps-action/utils"
)

type dummyLogger struct{}

func (d dummyLogger) Errorf(format string, args ...interface{}) {}
func (d dummyLogger) Infof(format string, args ...interface{})  {}
func (d dummyLogger) Warnf(format string, args ...interface{})  {}
func (d dummyLogger) Debugf(format string, args ...interface{}) {}
func (d dummyLogger) Fatalf(format string, args ...interface{}) {}
func (d dummyLogger) Sync() error                               { return nil }

func makeCPSUser() model.CPSUser {
	return model.CPSUser{
		ID:           bson.NewObjectID(),
		UserCode:     "TEST001",
		FullName:     "Test User",
		PhoneNumber:  "1234567890",
		Email:        "test@example.com",
		Role:         "Maker",
		Department:   bson.NewObjectID(),
		Enabled:      true,
		IsDeleted:    false,
		LastModified: &time.Time{},
	}
}

func makeCPSAction() model.CPSAction {
	return model.CPSAction{
		ID:               bson.NewObjectID(),
		ActionCode:       "ACT001",
		UniqueId:         "TEST001",
		MakerID:          "MAKER001",
		MakerName:        "Test Maker",
		MakerPhoneNumber: "1234567890",
		Department:       "IT",
		ActionStatus:     string(model.ActionPending),
		ActionType:       string(model.ActionCreate),
		RequestAction:    string(model.RequestUser),
		CreatedAt:        time.Now(),
		MakerActionTime:  time.Now(),
	}
}

func makeUserContext() context.Context {
	ctx := context.Background()
	ctx = context.WithValue(ctx, constant.ContextKey("user_code"), "TEST001")
	ctx = context.WithValue(ctx, constant.ContextKey("user_id"), "USER001")
	ctx = context.WithValue(ctx, constant.ContextKey("full_name"), "Test User")
	ctx = context.WithValue(ctx, constant.ContextKey("phone_number"), "1234567890")
	ctx = context.WithValue(ctx, constant.ContextKey("department"), "IT")
	ctx = context.WithValue(ctx, constant.ContextKey("user_role"), "Maker")
	return ctx
}

func makeTestRequestWithContext() *http.Request {
	req, _ := http.NewRequest("GET", "/test", nil)
	ctx := makeUserContext()
	return req.WithContext(ctx)
}

func TestCreateUserRequest(t *testing.T) {
	mockRepo := new(mocks.MockCPSActionRepo)
	mockPermissionService := new(mocks.MockPermissionDomainService)
	service := cpsuserServices.NewCPSUserService(mockRepo, mockPermissionService, dummyLogger{})
	ctx := makeUserContext()
	cpsAction := makeCPSAction()

	// Mock permission service calls
	mockPermissionService.On("ValidatePermissionCategories", mock.Anything).Return([]string{}, nil)
	mockPermissionService.On("ValidatePermissionGroups", mock.Anything).Return([]string{}, nil)

	// Mock repository calls
	mockRepo.On("FetchPendingActionsByUniqueID", mock.Anything, mock.Anything).Return([]action.ActionResponse{}, nil)
	mockRepo.On("CreateUserRequest", mock.Anything, mock.Anything).Return(&cpsAction, nil)

	result, err := service.CreateUserRequest(ctx, makeTestRequestWithContext(), dto.CreateUserRequest{})
	assert.NoError(t, err)
	assert.NotNil(t, result)
	mockRepo.AssertExpectations(t)
	mockPermissionService.AssertExpectations(t)
}

func TestUpdateUserRequest(t *testing.T) {
	mockRepo := new(mocks.MockCPSActionRepo)
	mockPermissionService := new(mocks.MockPermissionDomainService)
	service := cpsuserServices.NewCPSUserService(mockRepo, mockPermissionService, dummyLogger{})
	ctx := makeUserContext()
	cpsAction := makeCPSAction()

	// Mock repository calls
	mockRepo.On("FetchPendingActionsByUniqueID", mock.Anything, mock.Anything).Return([]action.ActionResponse{}, nil)
	mockRepo.On("UpdateUserRequest", mock.Anything, mock.Anything, mock.Anything).Return(&cpsAction, nil)

	result, err := service.UpdateUserRequest(ctx, makeTestRequestWithContext(), dto.UpdateUserRequest{}, "TEST001")
	assert.NoError(t, err)
	assert.NotNil(t, result)
	mockRepo.AssertExpectations(t)
}

func TestGetPendingUserActions(t *testing.T) {
	mockRepo := new(mocks.MockCPSActionRepo)
	service := cpsuserServices.NewCPSUserService(mockRepo, nil, dummyLogger{})
	ctx := makeUserContext()
	actions := []model.CPSAction{makeCPSAction()}
	mockRepo.On("GetPendingUserActions", mock.Anything).Return(actions, nil)

	result, err := service.GetPendingUserActions(ctx)
	assert.NoError(t, err)
	assert.NotNil(t, result)
	assert.Len(t, result, 1)
	mockRepo.AssertExpectations(t)
}

func TestApproveUserAction(t *testing.T) {
	mockRepo := new(mocks.MockCPSActionRepo)
	service := cpsuserServices.NewCPSUserService(mockRepo, nil, dummyLogger{})
	ctx := makeUserContext()
	cpsAction := makeCPSAction()
	mockRepo.On("ApproveUserAction", mock.Anything, mock.Anything, mock.Anything).Return(&cpsAction, nil)

	approveData := dto.ApproveCPSAction{Approved: true, Reason: nil}
	result, err := service.ApproveUserAction(ctx, makeTestRequestWithContext(), approveData, "ACT001")
	assert.NoError(t, err)
	assert.NotNil(t, result)
	mockRepo.AssertExpectations(t)
}

func TestFetchUserByUserCode(t *testing.T) {
	mockRepo := new(mocks.MockCPSActionRepo)
	service := cpsuserServices.NewCPSUserService(mockRepo, nil, dummyLogger{})
	ctx := makeUserContext()
	cpsUser := makeCPSUser()
	mockRepo.On("FetchUserByUserCode", mock.Anything, mock.Anything).Return(&cpsUser, nil)

	result, err := service.FetchUserByUserCode(ctx, "TEST001")
	assert.NoError(t, err)
	assert.NotNil(t, result)
	assert.Equal(t, "TEST001", result.UserCode)
	mockRepo.AssertExpectations(t)
}

func TestGetAllCPSUser(t *testing.T) {
	mockRepo := new(mocks.MockCPSActionRepo)
	ctx := makeUserContext()

	cpsUser := makeCPSUser()
	expectedResponse := &common_util.PaginatedResponse[[]*model.CPSUser]{
		Data: []*model.CPSUser{&cpsUser},
		Meta: common_util.PaginationMeta{
			TotalDocs:  1,
			Limit:      10,
			TotalPages: 1,
			Page:       1,
		},
	}

	mockRepo.On("GetAllCPSUsers", mock.Anything, mock.Anything).Return(expectedResponse, nil)

	result, err := mockRepo.GetAllCPSUsers(ctx, &constant.Filter{Page: 1, PerPage: 10})
	assert.NoError(t, err)
	assert.NotNil(t, result)
	assert.Len(t, result.Data, 1)
}
