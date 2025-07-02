package amount_based_auth_domain_test

import (
	"context"
	"errors"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"go.uber.org/mock/gomock"

	"github.com/CBE-Super-App/cbe-super-app-cps-action/internal/adapter/outbound/model"
	amount_based_auth_domain "github.com/CBE-Super-App/cbe-super-app-cps-action/internal/domain/amount_based_auth"
	"github.com/CBE-Super-App/cbe-super-app-cps-action/internal/domain/amount_based_auth/mocks"
)

// MockLogger implements utils.Logger for testing
// Only methods used in the service are implemented
// (You can expand as needed for more detailed logging tests)
type MockLogger struct{}

func (l *MockLogger) Debug(args ...interface{})                 {}
func (l *MockLogger) Debugf(format string, args ...interface{}) {}
func (l *MockLogger) Info(args ...interface{})                  {}
func (l *MockLogger) Infof(format string, args ...interface{})  {}
func (l *MockLogger) Warn(args ...interface{})                  {}
func (l *MockLogger) Warnf(format string, args ...interface{})  {}
func (l *MockLogger) Error(args ...interface{})                 {}
func (l *MockLogger) Errorf(format string, args ...interface{}) {}
func (l *MockLogger) Fatal(args ...interface{})                 {}
func (l *MockLogger) Fatalf(format string, args ...interface{}) {}
func (l *MockLogger) Sync() error                               { return nil }

func TestService_UpdateAuthTier(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockRepo := mocks.NewMockRepository(ctrl)
	mockLogger := &MockLogger{}
	svc := amount_based_auth_domain.NewAmountBasedAuthService(mockRepo, mockLogger)

	testTime := time.Now()
	testUser := model.User{
		UserCode:    "USER001",
		FullName:    "John Doe",
		PhoneNumber: "+251911234567",
	}
	testRequest := amount_based_auth_domain.UpdateAmountBasedAuth{
		Id:        "AUTH001",
		MinAmount: 1,
		MaxAmount: 1000,
		Method:    amount_based_auth_domain.OPEN,
	}
	testCPSAction := model.CreateCPSAction{
		ActionCode:      "ACT001",
		MakerUser:       testUser,
		Department:      "IT",
		Status:          model.ActionPending,
		RequestAction:   model.RequestAuthTier,
		ActionType:      model.ActionUpdate,
		CurrentData:     testRequest,
		MakerActionTime: testTime,
	}
	testCpsActionRes := &model.CpsAction{
		ID:              "CPS001",
		ActionCode:      "ACT001",
		MakerUser:       testUser,
		Department:      "IT",
		Status:          model.ActionPending,
		RequestAction:   model.RequestAuthTier,
		ActionType:      model.ActionUpdate,
		CurrentData:     testRequest,
		MakerActionTime: testTime,
	}

	tests := []struct {
		name    string
		req     amount_based_auth_domain.UpdateAmountBasedAuth
		cps     model.CreateCPSAction
		mock    func()
		want    *model.CpsAction
		wantErr bool
	}{
		{
			name: "success",
			req:  testRequest,
			cps:  testCPSAction,
			mock: func() {
				mockRepo.EXPECT().
					UpdateAuthTier(gomock.Any(), testRequest, testCPSAction).
					Return(testCpsActionRes, nil)
			},
			want:    testCpsActionRes,
			wantErr: false,
		},
		{
			name: "repo error",
			req:  testRequest,
			cps:  testCPSAction,
			mock: func() {
				mockRepo.EXPECT().
					UpdateAuthTier(gomock.Any(), testRequest, testCPSAction).
					Return(nil, errors.New("db error"))
			},
			want:    nil,
			wantErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tt.mock()
			got, err := svc.UpdateAuthTier(context.Background(), tt.req, tt.cps)
			if tt.wantErr {
				assert.Error(t, err)
				assert.Nil(t, got)
			} else {
				assert.NoError(t, err)
				assert.Equal(t, tt.want, got)
			}
		})
	}
}

func TestService_ApproveAuthTierApprove(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockRepo := mocks.NewMockRepository(ctrl)
	mockLogger := &MockLogger{}
	svc := amount_based_auth_domain.NewAmountBasedAuthService(mockRepo, mockLogger)

	testTime := time.Now()
	testChecker := model.User{
		UserCode:    "CHECKER001",
		FullName:    "Checker User",
		PhoneNumber: "+251911234568",
	}
	testAuthorize := model.AuthorizeCPSAction{
		ActionCode:        "ACT001",
		Department:        "IT",
		CheckerUser:       testChecker,
		CheckerActionTime: testTime,
	}
	testCpsAction := &model.CpsAction{
		ID:                "CPS001",
		ActionCode:        "ACT001",
		CheckerUser:       testChecker,
		Department:        "IT",
		Status:            model.ActionApproved,
		RequestAction:     model.RequestAuthTier,
		ActionType:        model.ActionUpdate,
		CheckerActionTime: testTime,
	}

	tests := []struct {
		name    string
		id      string
		req     model.AuthorizeCPSAction
		mock    func()
		want    *model.CpsAction
		wantErr bool
	}{
		{
			name: "success",
			id:   "AUTH001",
			req:  testAuthorize,
			mock: func() {
				mockRepo.EXPECT().
					ApproveAmountBasedAuth(gomock.Any(), "AUTH001", testAuthorize).
					Return(testCpsAction, nil)
			},
			want:    testCpsAction,
			wantErr: false,
		},
		{
			name: "repo error",
			id:   "AUTH001",
			req:  testAuthorize,
			mock: func() {
				mockRepo.EXPECT().
					ApproveAmountBasedAuth(gomock.Any(), "AUTH001", testAuthorize).
					Return(nil, errors.New("db error"))
			},
			want:    nil,
			wantErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tt.mock()
			got, err := svc.ApproveAuthTierApprove(context.Background(), tt.id, tt.req)
			if tt.wantErr {
				assert.Error(t, err)
				assert.Nil(t, got)
			} else {
				assert.NoError(t, err)
				assert.Equal(t, tt.want, got)
			}
		})
	}
}

func TestService_RejectAuthTier(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockRepo := mocks.NewMockRepository(ctrl)
	mockLogger := &MockLogger{}
	svc := amount_based_auth_domain.NewAmountBasedAuthService(mockRepo, mockLogger)

	testTime := time.Now()
	testMaker := model.User{
		UserCode:    "MAKER001",
		FullName:    "Maker User",
		PhoneNumber: "+251911234567",
	}
	testChecker := model.User{
		UserCode:    "CHECKER001",
		FullName:    "Checker User",
		PhoneNumber: "+251911234568",
	}
	testReject := model.RejectCPSAction{
		CreateCPSAction: model.CreateCPSAction{
			ActionCode:      "ACT001",
			MakerUser:       testMaker,
			Department:      "IT",
			Status:          model.ActionPending,
			RequestAction:   model.RequestAuthTier,
			ActionType:      model.ActionUpdate,
			MakerActionTime: testTime,
		},
		CheckerUser:       testChecker,
		RejectedReason:    "Invalid data provided by the user during the request",
		CheckerActionTime: testTime,
	}
	testCpsAction := &model.CpsAction{
		ID:                "CPS001",
		ActionCode:        "ACT001",
		MakerUser:         testMaker,
		CheckerUser:       testChecker,
		Department:        "IT",
		Status:            model.ActionRejected,
		RequestAction:     model.RequestAuthTier,
		ActionType:        model.ActionUpdate,
		RejectedReason:    "Invalid data",
		MakerActionTime:   testTime,
		CheckerActionTime: testTime,
	}

	tests := []struct {
		name    string
		id      string
		req     model.RejectCPSAction
		mock    func()
		want    *model.CpsAction
		wantErr bool
	}{
		{
			name: "success",
			id:   "AUTH001",
			req:  testReject,
			mock: func() {
				mockRepo.EXPECT().
					RejectAmountBasedAuth(gomock.Any(), "AUTH001", testReject).
					Return(testCpsAction, nil)
			},
			want:    testCpsAction,
			wantErr: false,
		},
		{
			name: "repo error",
			id:   "AUTH001",
			req:  testReject,
			mock: func() {
				mockRepo.EXPECT().
					RejectAmountBasedAuth(gomock.Any(), "AUTH001", testReject).
					Return(nil, errors.New("db error"))
			},
			want:    nil,
			wantErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tt.mock()
			got, err := svc.RejectAuthTier(context.Background(), tt.id, tt.req)
			if tt.wantErr {
				assert.Error(t, err)
				assert.Nil(t, got)
			} else {
				assert.NoError(t, err)
				assert.Equal(t, tt.want, got)
			}
		})
	}
}
