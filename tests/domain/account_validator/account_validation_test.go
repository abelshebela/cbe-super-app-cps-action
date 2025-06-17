package account_validator_test

import (
	"context"
	"encoding/json"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"

	"gitlab.com/bersufekadgetachew/cbe-super-app-cps-action/internal/domain/account_validation"
	"gitlab.com/bersufekadgetachew/cbe-super-app-cps-action/internal/domain/action"
	"gitlab.com/bersufekadgetachew/cbe-super-app-cps-action/mocks"
)

// MockRepository is a mock implementation of the account_validation.Repository interface
type MockRepository struct {
	mock.Mock
}

func (m *MockRepository) GetAccountValidationByID(ctx context.Context, id string) (account_validation.ValidationRule, error) {
	args := m.Called(ctx, id)
	return args.Get(0).(account_validation.ValidationRule), args.Error(1)
}

func (m *MockRepository) UpdateAccountValidation(ctx context.Context, id string, update account_validation.ValidationRule) error {
	args := m.Called(ctx, id, update)
	return args.Error(0)
}

func (m *MockRepository) FetchPendingActionsByUniqueID(ctx context.Context, uniqueID string) ([]action.CPSAction, error) {
	args := m.Called(ctx, uniqueID)
	return args.Get(0).([]action.CPSAction), args.Error(1)
}

// MockActionRepository is a mock implementation of the action.Repository interface
type MockActionRepository struct {
	mock.Mock
}

func (m *MockActionRepository) CreateCpsAction(ctx context.Context, cpsAction action.CPSAction) (action.CPSAction, error) {
	args := m.Called(ctx, cpsAction)
	return args.Get(0).(action.CPSAction), args.Error(1)
}

func (m *MockActionRepository) FetchCpsActionById(ctx context.Context, id string) (action.CPSAction, error) {
	args := m.Called(ctx, id)
	return args.Get(0).(action.CPSAction), args.Error(1)
}

func (m *MockActionRepository) UpdateCpsAction(ctx context.Context, cpsAction action.CPSAction) error {
	args := m.Called(ctx, cpsAction)
	return args.Error(0)
}

func (m *MockActionRepository) FetchAccountsByAccountNumber(ctx context.Context, accountNumber string) ([]action.LinkedAccount, error) {
	args := m.Called(ctx, accountNumber)
	return args.Get(0).([]action.LinkedAccount), args.Error(1)
}

func (m *MockActionRepository) FetchLastCpsActionByMakerID(ctx context.Context, makerId string) (action.CPSAction, error) {
	args := m.Called(ctx, makerId)
	return args.Get(0).(action.CPSAction), args.Error(1)
}

func (m *MockActionRepository) GetAllHqServices(ctx context.Context) ([]action.ServiceDetails, error) {
	args := m.Called(ctx)
	return args.Get(0).([]action.ServiceDetails), args.Error(1)
}

func (m *MockActionRepository) GetAllHqServicesPaginated(ctx context.Context, offset, limit int) ([]action.ServiceDetails, error) {
	args := m.Called(ctx, offset, limit)
	return args.Get(0).([]action.ServiceDetails), args.Error(1)
}

func (m *MockActionRepository) GetHqServiceById(ctx context.Context, id string) (action.ServiceDetails, error) {
	args := m.Called(ctx, id)
	return args.Get(0).(action.ServiceDetails), args.Error(1)
}

func (m *MockActionRepository) UpdateHqService(ctx context.Context, service action.ServiceDetails) error {
	args := m.Called(ctx, service)
	return args.Error(0)
}

func (m *MockActionRepository) FetchLinkedAccountById(ctx context.Context, id []string) ([]action.LinkedAccount, error) {
	args := m.Called(ctx, id)
	return args.Get(0).([]action.LinkedAccount), args.Error(1)
}

func (m *MockActionRepository) UpdateAccounts(ctx context.Context, linkedAccounts []action.LinkedAccount) ([]action.LinkedAccount, error) {
	args := m.Called(ctx, linkedAccounts)
	return args.Get(0).([]action.LinkedAccount), args.Error(1)
}

func (m *MockActionRepository) UpdateAccount(ctx context.Context, linkedAccount action.LinkedAccount) (action.LinkedAccount, error) {
	args := m.Called(ctx, linkedAccount)
	return args.Get(0).(action.LinkedAccount), args.Error(1)
}

// NoOpLogger implements the utils.Logger interface but does nothing
// Use this in tests to avoid nil pointer panics

type NoOpLogger struct{}

func (l *NoOpLogger) Infof(format string, args ...interface{})  {}
func (l *NoOpLogger) Errorf(format string, args ...interface{}) {}
func (l *NoOpLogger) Debugf(format string, args ...interface{}) {}
func (l *NoOpLogger) Fatalf(format string, args ...interface{}) {}
func (l *NoOpLogger) Warnf(format string, args ...interface{})  {}
func (l *NoOpLogger) Sync() error                               { return nil }

func TestGetAccountValidation(t *testing.T) {
	mockRepo := new(mocks.AccountValidationRepository)
	mockActionRepo := new(mocks.ActionRepository)
	logger := &NoOpLogger{}
	service := account_validation.NewService(mockRepo, mockActionRepo, logger)

	ctx := context.Background()
	expectedRule := account_validation.ValidationRule{
		ID:            "test-id",
		EntityType:    "account",
		ValidationFor: "phone",
		Identifier:    "phone_number",
		MinLength:     10,
		MaxLength:     15,
		Enabled:       true,
		ServiceID:     "test-service",
	}

	t.Run("successful get", func(t *testing.T) {
		mockRepo.On("GetAccountValidationByID", ctx, "test-id").Return(expectedRule, nil).Once()

		rule, err := service.GetAccountValidation(ctx, "test-id")
		assert.NoError(t, err)
		assert.Equal(t, expectedRule, rule)
	})

	t.Run("empty id", func(t *testing.T) {
		rule, err := service.GetAccountValidation(ctx, "")
		assert.Error(t, err)
		assert.Equal(t, account_validation.ValidationRule{}, rule)
		assert.Equal(t, "validation rule ID cannot be empty", err.Error())
	})
}

func TestUpdateAccountValidationRequest(t *testing.T) {
	mockRepo := new(mocks.AccountValidationRepository)
	mockActionRepo := new(mocks.ActionRepository)
	logger := &NoOpLogger{}
	service := account_validation.NewService(mockRepo, mockActionRepo, logger)

	ctx := context.Background()
	originalRule := account_validation.ValidationRule{
		ID:            "test-id",
		EntityType:    "account",
		ValidationFor: "phone",
		Identifier:    "phone_number",
		MinLength:     10,
		MaxLength:     15,
		Enabled:       true,
		ServiceID:     "test-service",
	}

	updatedRule := originalRule
	updatedRule.MinLength = 11

	t.Run("successful update request", func(t *testing.T) {
		mockRepo.On("GetAccountValidationByID", ctx, "test-id").Return(originalRule, nil).Once()

		expectedAction := action.CPSAction{
			ActionCode: "CPS_123",
			Maker: action.User{
				UserID:      "maker-123",
				FullName:    "",
				PhoneNumber: "",
				Timestamp:   time.Now(),
			},
			Checker:         action.User{},
			Department:      "test-service",
			ActionType:      action.ActionUpdate,
			RequestAction:   action.RequestUpdateAccountValidation,
			ActionStatus:    action.ActionPending,
			CreatedAt:       time.Now(),
			LastModifiedAt:  time.Now(),
			RejectionReason: nil,
		}

		mockActionRepo.On("CreateCpsAction", ctx, mock.Anything).Return(expectedAction, nil).Once()

		actionID, err := service.UpdateAccountValidationRequest(ctx, "test-id", updatedRule, "maker-123")
		assert.NoError(t, err)
		assert.NotEmpty(t, actionID)
		assert.Contains(t, actionID, "CPS_")
	})

	t.Run("validation rule not found", func(t *testing.T) {
		mockRepo.On("GetAccountValidationByID", ctx, "test-id").Return(account_validation.ValidationRule{}, assert.AnError).Once()

		actionID, err := service.UpdateAccountValidationRequest(ctx, "test-id", updatedRule, "maker-123")
		assert.Error(t, err)
		assert.Empty(t, actionID)
	})
}

func TestUpdateAccountValidation(t *testing.T) {
	mockRepo := new(mocks.AccountValidationRepository)
	mockActionRepo := new(mocks.ActionRepository)
	logger := &NoOpLogger{}
	service := account_validation.NewService(mockRepo, mockActionRepo, logger)

	ctx := context.Background()
	actionID := "CPS_123"
	checkerID := "checker-123"

	validationRule := account_validation.ValidationRule{
		ID:            "test-id",
		EntityType:    "account",
		ValidationFor: "phone",
		Identifier:    "phone_number",
		MinLength:     10,
		MaxLength:     15,
		Enabled:       true,
		ServiceID:     "test-service",
	}

	type CurrentAction struct {
		Rule account_validation.ValidationRule `json:"rule"`
	}
	currentAction := CurrentAction{Rule: validationRule}
	currentActionBytes, _ := json.Marshal(currentAction)

	t.Run("successful approval", func(t *testing.T) {
		expectedAction := action.CPSAction{
			ActionCode: actionID,
			Maker: action.User{
				UserID:      "maker-123",
				FullName:    "",
				PhoneNumber: "",
				Timestamp:   time.Now(),
			},
			Checker:         action.User{},
			Department:      "test-service",
			ActionType:      action.ActionUpdate,
			RequestAction:   action.RequestUpdateAccountValidation,
			ActionStatus:    action.ActionPending,
			CurrentAction:   currentActionBytes,
			CreatedAt:       time.Now(),
			LastModifiedAt:  time.Now(),
			RejectionReason: nil,
		}

		mockActionRepo.On("FetchCpsActionById", ctx, actionID).Return(expectedAction, nil).Once()
		mockRepo.On("UpdateAccountValidation", ctx, "test-id", mock.Anything).Return(nil).Once()
		mockActionRepo.On("UpdateCpsAction", ctx, mock.Anything).Return(nil).Once()

		err := service.UpdateAccountValidation(ctx, actionID, true, checkerID)
		assert.NoError(t, err)
	})

	t.Run("successful rejection", func(t *testing.T) {
		expectedAction := action.CPSAction{
			ActionCode: actionID,
			Maker: action.User{
				UserID:      "maker-123",
				FullName:    "",
				PhoneNumber: "",
				Timestamp:   time.Now(),
			},
			Checker:         action.User{},
			Department:      "test-service",
			ActionType:      action.ActionUpdate,
			RequestAction:   action.RequestUpdateAccountValidation,
			ActionStatus:    action.ActionPending,
			CurrentAction:   currentActionBytes,
			CreatedAt:       time.Now(),
			LastModifiedAt:  time.Now(),
			RejectionReason: nil,
		}

		mockActionRepo.On("FetchCpsActionById", ctx, actionID).Return(expectedAction, nil).Once()
		mockActionRepo.On("UpdateCpsAction", ctx, mock.Anything).Return(nil).Once()

		err := service.UpdateAccountValidation(ctx, actionID, false, checkerID)
		assert.NoError(t, err)
	})

	t.Run("empty action ID", func(t *testing.T) {
		err := service.UpdateAccountValidation(ctx, "", true, checkerID)
		assert.Error(t, err)
		assert.Equal(t, "action ID cannot be empty", err.Error())
	})

	t.Run("empty checker ID", func(t *testing.T) {
		err := service.UpdateAccountValidation(ctx, actionID, true, "")
		assert.Error(t, err)
		assert.Equal(t, "checker ID cannot be empty", err.Error())
	})

	t.Run("action not found", func(t *testing.T) {
		mockActionRepo.On("FetchCpsActionById", ctx, actionID).Return(action.CPSAction{}, assert.AnError).Once()

		err := service.UpdateAccountValidation(ctx, actionID, true, checkerID)
		assert.Error(t, err)
		assert.Equal(t, "action not found", err.Error())
	})

	t.Run("action not pending", func(t *testing.T) {
		expectedAction := action.CPSAction{
			ActionCode: actionID,
			Maker: action.User{
				UserID:      "maker-123",
				FullName:    "",
				PhoneNumber: "",
				Timestamp:   time.Now(),
			},
			Checker:         action.User{},
			Department:      "test-service",
			ActionType:      action.ActionUpdate,
			RequestAction:   action.RequestUpdateAccountValidation,
			ActionStatus:    action.ActionApproved,
			CurrentAction:   currentActionBytes,
			CreatedAt:       time.Now(),
			LastModifiedAt:  time.Now(),
			RejectionReason: nil,
		}

		mockActionRepo.On("FetchCpsActionById", ctx, actionID).Return(expectedAction, nil).Once()

		err := service.UpdateAccountValidation(ctx, actionID, true, checkerID)
		assert.Error(t, err)
		assert.Equal(t, "action is not pending", err.Error())
	})
}
