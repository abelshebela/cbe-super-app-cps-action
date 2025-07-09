package service

import (
	"context"
	"errors"
	"testing"
	"time"

	"github.com/CBE-Super-App/cbe-super-app-cps-action/internal/domain/fayda_account/entity"
	"github.com/CBE-Super-App/cbe-super-app-cps-action/internal/domain/fayda_account/mocks"

	"github.com/stretchr/testify/assert"
	"go.uber.org/mock/gomock"
)

type MockLogger struct{}

func (m *MockLogger) Infof(msg string, args ...interface{})  {}
func (m *MockLogger) Errorf(msg string, args ...interface{}) {}
func (m *MockLogger) Debugf(msg string, args ...interface{}) {}
func (m *MockLogger) Fatalf(msg string, args ...interface{}) {}
func (m *MockLogger) Warnf(msg string, args ...interface{})  {}
func (m *MockLogger) Sync() error                            { return nil }

func sampleInitiateCPSAction() entity.CPSAction {
	return entity.CPSAction{
		ID:         "id1",
		ActionCode: "code1",
		MakerUser: entity.User{
			UserCode:    "maker1",
			FullName:    "Maker User",
			PhoneNumber: "987654321",
		},
		RejectedReason:  "",
		Department:      "IT",
		Status:          entity.ActionPending,
		RequestAction:   entity.RequestDisableFaydaAccount,
		ActionType:      entity.ActionCreate,
		ActionData:      entity.ActionData{UseCode: "maker1", FullName: "Maker User", PhoneNumber: "987654321"},
		MakerActionTime: time.Now(),
	}
}

func sampleAuthorizeCPSAction() entity.CPSAction {
	return entity.CPSAction{
		ID:         "id1",
		ActionCode: "code1",
		CheckerUser: entity.User{
			UserCode:    "checker1",
			FullName:    "Checker User",
			PhoneNumber: "123456789",
		},
		MakerUser: entity.User{
			UserCode:    "maker1",
			FullName:    "Maker User",
			PhoneNumber: "987654321",
		},
		RejectedReason:    "",
		Department:        "IT",
		Status:            entity.ActionApproved,
		RequestAction:     entity.RequestDisableFaydaAccount,
		ActionData:        entity.ActionData{UseCode: "maker1", FullName: "Maker User", PhoneNumber: "987654321"},
		MakerActionTime:   time.Now(),
		CheckerActionTime: time.Now(),
	}
}

func sampleRejectCPSAction() entity.CPSAction {
	return entity.CPSAction{
		ID:         "id1",
		ActionCode: "code1",
		CheckerUser: entity.User{
			UserCode:    "checker1",
			FullName:    "Checker User",
			PhoneNumber: "123456789",
		},
		MakerUser: entity.User{
			UserCode:    "maker1",
			FullName:    "Maker User",
			PhoneNumber: "987654321",
		},
		RejectedReason:    "sucpicious account",
		Department:        "IT",
		Status:            entity.ActionRejected,
		RequestAction:     entity.RequestDisableFaydaAccount,
		ActionData:        entity.ActionData{UseCode: "maker1", FullName: "Maker User", PhoneNumber: "987654321"},
		MakerActionTime:   time.Now(),
		CheckerActionTime: time.Now(),
	}
}

func TestInitiateDisableFaydaAccount(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockRepo := mocks.NewMockRepository(ctrl)
	mockLogger := &MockLogger{}
	service := InitFaydaAccountDomain(mockRepo, mockLogger)

	req := sampleInitiateCPSAction()
	resp := req

	t.Run("success", func(t *testing.T) {
		mockRepo.EXPECT().
			InitiateDisableFaydaAccount(gomock.Any(), req).
			Return(&resp, nil)

		result, err := service.InitiateDisableFaydaAccount(context.Background(), req)
		assert.NoError(t, err)
		assert.Equal(t, &resp, result)
	})

	t.Run("error", func(t *testing.T) {
		mockRepo.EXPECT().
			InitiateDisableFaydaAccount(gomock.Any(), req).
			Return(nil, errors.New("internal server error"))

		result, err := service.InitiateDisableFaydaAccount(context.Background(), req)
		assert.Error(t, err)
		assert.Nil(t, result)
	})
}

func TestAuthorizeFaydaAccountDisable(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockRepo := mocks.NewMockRepository(ctrl)
	mockLogger := &MockLogger{}
	service := InitFaydaAccountDomain(mockRepo, mockLogger)

	req := sampleAuthorizeCPSAction()
	resp := req

	t.Run("success", func(t *testing.T) {
		mockRepo.EXPECT().
			AuthorizeFaydaAccountDisable(gomock.Any(), req).
			Return(&resp, nil)

		result, err := service.AuthorizeFaydaAccountDisable(context.Background(), req)
		assert.NoError(t, err)
		assert.Equal(t, &resp, result)
	})

	t.Run("error", func(t *testing.T) {
		mockRepo.EXPECT().
			AuthorizeFaydaAccountDisable(gomock.Any(), req).
			Return(nil, errors.New("internal server error"))

		result, err := service.AuthorizeFaydaAccountDisable(context.Background(), req)
		assert.Error(t, err)
		assert.Nil(t, result)
	})
}

func TestRejectFaydaAccountDisable(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockRepo := mocks.NewMockRepository(ctrl)
	mockLogger := &MockLogger{}
	service := InitFaydaAccountDomain(mockRepo, mockLogger)

	req := sampleRejectCPSAction()
	resp := req

	t.Run("success", func(t *testing.T) {
		mockRepo.EXPECT().
			RejectFaydaAccountDisable(gomock.Any(), req).
			Return(&resp, nil)

		result, err := service.RejectFaydaAccountDisable(context.Background(), req)
		assert.NoError(t, err)
		assert.Equal(t, &resp, result)
	})

	t.Run("error", func(t *testing.T) {
		mockRepo.EXPECT().
			RejectFaydaAccountDisable(gomock.Any(), req).
			Return(nil, errors.New("internal server error"))

		result, err := service.RejectFaydaAccountDisable(context.Background(), req)
		assert.Error(t, err)
		assert.Nil(t, result)
	})
}
