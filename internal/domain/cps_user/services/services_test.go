package services

import (
	"context"
	"errors"
	"net/http"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"

	"github.com/CBE-Super-App/cbe-super-app-cps-action/internal/adapter/outbound/model"
	dto "github.com/CBE-Super-App/cbe-super-app-cps-action/internal/domain/cps_user/dto"
	repomock "github.com/CBE-Super-App/cbe-super-app-cps-action/internal/domain/cps_user/repository/mocks"
	"go.mongodb.org/mongo-driver/v2/bson"
)

type dummyLogger struct{}

func (d dummyLogger) Errorf(format string, args ...interface{}) {}
func (d dummyLogger) Infof(format string, args ...interface{})  {}
func (d dummyLogger) Warnf(format string, args ...interface{})  {}
func (d dummyLogger) Debugf(format string, args ...interface{}) {}

func makeUserContext() context.Context {
	ctx := context.Background()
	ctx = context.WithValue(ctx, "user_code", "U123")
	ctx = context.WithValue(ctx, "user_id", "ID123")
	ctx = context.WithValue(ctx, "full_name", "Test User")
	ctx = context.WithValue(ctx, "phone_number", "1234567890")
	ctx = context.WithValue(ctx, "department", "Dept1")
	return ctx
}

func makeCPSAction() model.CPSAction {
	return model.CPSAction{
		ActionCode:    "A1",
		UniqueId:      "U123",
		MakerID:       "ID123",
		MakerName:     "Test User",
		Department:    "Dept1",
		ActionStatus:  string(model.ActionPending),
		ActionType:    string(model.ActionCreate),
		RequestAction: string(model.RequestUser),
		CreatedAt:     time.Now(),
	}
}

func makeCPSUser() model.CPSUser {
	return model.CPSUser{
		UserCode:    "U123",
		FullName:    "Test User",
		Role:        "MAKER",
		Department:  bson.ObjectID{},
		PhoneNumber: "1234567890",
	}
}

func makeTestRequestWithContext() *http.Request {
	req, _ := http.NewRequest("GET", "/", nil)
	ctx := req.Context()
	ctx = context.WithValue(ctx, "user_code", "U123")
	ctx = context.WithValue(ctx, "user_id", "ID123")
	ctx = context.WithValue(ctx, "full_name", "Test User")
	ctx = context.WithValue(ctx, "phone_number", "1234567890")
	ctx = context.WithValue(ctx, "department", "Dept1")
	req = req.WithContext(ctx)
	return req
}

func TestCreateUserRequest(t *testing.T) {
	mockRepo := new(repomock.MockCPSActionRepo)
	service := NewCPSUserService(mockRepo)
	ctx := makeUserContext()
	cpsAction := makeCPSAction()
	mockRepo.On("CreateUserRequest", mock.Anything, mock.Anything).Return(&cpsAction, nil)

	result, err := service.CreateUserRequest(ctx, makeTestRequestWithContext(), dto.CreateUserRequest{})
	assert.NoError(t, err)
	assert.NotNil(t, result)
	mockRepo.AssertExpectations(t)

	// Clear previous expectations before setting new ones
	mockRepo.ExpectedCalls = nil
	mockRepo.On("CreateUserRequest", mock.Anything, mock.Anything).Return(nil, errors.New("fail"))
	result, err = service.CreateUserRequest(ctx, makeTestRequestWithContext(), dto.CreateUserRequest{})
	assert.Error(t, err)
	assert.Nil(t, result)
}

func TestUpdateUserRequest(t *testing.T) {
	mockRepo := new(repomock.MockCPSActionRepo)
	service := NewCPSUserService(mockRepo)
	ctx := makeUserContext()
	cpsAction := makeCPSAction()
	mockRepo.On("UpdateUserRequest", mock.Anything, mock.Anything, mock.Anything).Return(&cpsAction, nil)

	result, err := service.UpdateUserRequest(ctx, makeTestRequestWithContext(), dto.UpdateUserRequest{}, "U123")
	assert.NoError(t, err)
	assert.NotNil(t, result)
	mockRepo.AssertExpectations(t)

	// Clear previous expectations before setting new ones
	mockRepo.ExpectedCalls = nil
	mockRepo.On("UpdateUserRequest", mock.Anything, mock.Anything, mock.Anything).Return(nil, errors.New("fail"))
	result, err = service.UpdateUserRequest(ctx, makeTestRequestWithContext(), dto.UpdateUserRequest{}, "U123")
	assert.Error(t, err)
	assert.Nil(t, result)
}

func TestApproveUserAction(t *testing.T) {
	mockRepo := new(repomock.MockCPSActionRepo)
	service := NewCPSUserService(mockRepo)
	ctx := makeUserContext()
	cpsAction := makeCPSAction()
	mockRepo.On("ApproveUserAction", mock.Anything, mock.Anything, mock.Anything).Return(&cpsAction, nil)

	result, err := service.ApproveUserAction(ctx, makeTestRequestWithContext(), dto.ApproveCPSAction{Approved: true, Reason: nil}, "A1")
	assert.NoError(t, err)
	assert.NotNil(t, result)
	mockRepo.AssertExpectations(t)

	// For rejection, Reason must not be nil
	reason := "rejected for test"
	mockRepo.ExpectedCalls = nil
	mockRepo.On("ApproveUserAction", mock.Anything, mock.Anything, mock.Anything).Return(nil, errors.New("fail"))
	result, err = service.ApproveUserAction(ctx, makeTestRequestWithContext(), dto.ApproveCPSAction{Approved: false, Reason: &reason}, "A1")
	assert.Error(t, err)
	assert.Nil(t, result)
}

func TestGetPendingUserActions(t *testing.T) {
	mockRepo := new(repomock.MockCPSActionRepo)
	service := NewCPSUserService(mockRepo)
	ctx := makeUserContext()
	cpsAction := makeCPSAction()
	mockRepo.On("GetPendingUserActions", mock.Anything).Return([]model.CPSAction{cpsAction}, nil)

	result, err := service.GetPendingUserActions(ctx)
	assert.NoError(t, err)
	assert.NotNil(t, result)
	assert.Len(t, result, 1)
	mockRepo.AssertExpectations(t)

	// Error case
	mockRepo.ExpectedCalls = nil
	mockRepo.On("GetPendingUserActions", mock.Anything).Return(nil, errors.New("fail"))
	result, err = service.GetPendingUserActions(ctx)
	assert.Error(t, err)
	assert.Nil(t, result)

	// Empty slice case (should return error)
	mockRepo.ExpectedCalls = nil
	mockRepo.On("GetPendingUserActions", mock.Anything).Return([]model.CPSAction{}, nil)
	result, err = service.GetPendingUserActions(ctx)
	assert.Error(t, err)
	assert.Nil(t, result)
}

func TestFetchUserByUserCode(t *testing.T) {
	mockRepo := new(repomock.MockCPSActionRepo)
	service := NewCPSUserService(mockRepo)
	ctx := makeUserContext()
	cpsUser := makeCPSUser()
	mockRepo.On("FetchUserByUserCode", mock.Anything, "U123").Return(&cpsUser, nil)

	result, err := service.FetchUserByUserCode(ctx, "U123")
	assert.NoError(t, err)
	assert.NotNil(t, result)
	mockRepo.AssertExpectations(t)

	// Error case
	mockRepo.ExpectedCalls = nil
	mockRepo.On("FetchUserByUserCode", mock.Anything, "U123").Return(nil, errors.New("fail"))
	result, err = service.FetchUserByUserCode(ctx, "U123")
	assert.Error(t, err)
	assert.Nil(t, result)

	// Not found case (should return error)
	mockRepo.ExpectedCalls = nil
	mockRepo.On("FetchUserByUserCode", mock.Anything, "U123").Return(nil, nil)
	result, err = service.FetchUserByUserCode(ctx, "U123")
	assert.Error(t, err)
	assert.Nil(t, result)
}
