package account_validator_test

import (
	"context"
	"encoding/json"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"go.mongodb.org/mongo-driver/v2/bson"

	"github.com/CBE-Super-App/cbe-super-app-cps-action/internal/domain/account_validation"
	"github.com/CBE-Super-App/cbe-super-app-cps-action/internal/domain/action"
	"github.com/CBE-Super-App/cbe-super-app-cps-action/mocks"
)

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
		ID:            bson.NewObjectID(),
		EntityType:    "account",
		ValidationFor: "phone",
		Identifier:    "phone_number",
		MinLength:     10,
		MaxLength:     15,
		Enabled:       true,
		ServiceID:     "test-service",
	}

	t.Run("successful get", func(t *testing.T) {
		mockRepo.On("GetAccountValidationByID", ctx, expectedRule.ID).Return(expectedRule, nil).Once()

		rule, err := service.GetAccountValidation(ctx, expectedRule.ID.String())
		assert.NoError(t, err)
		assert.Equal(t, expectedRule, rule)
	})

	t.Run("empty id", func(t *testing.T) {
		rule, err := service.GetAccountValidation(ctx, "")
		assert.Error(t, err)
		assert.Equal(t, account_validation.ValidationRule{}, rule)
		assert.Equal(t, "INVALID_ID", err.Error())
	})

	t.Run("not found", func(t *testing.T) {
		mockRepo.On("GetAccountValidationByID", ctx, "not-exist").Return(account_validation.ValidationRule{}, assert.AnError).Once()
		rule, err := service.GetAccountValidation(ctx, "not-exist")
		assert.Error(t, err)
		assert.Equal(t, account_validation.ValidationRule{}, rule)
		assert.Equal(t, "NOT_FOUND", err.Error())
	})
}

func TestUpdateAccountValidationRequest(t *testing.T) {
	mockRepo := new(mocks.AccountValidationRepository)
	mockActionRepo := new(mocks.ActionRepository)
	logger := &NoOpLogger{}
	service := account_validation.NewService(mockRepo, mockActionRepo, logger)

	ctx := context.Background()
	originalRule := account_validation.ValidationRule{
		ID:            bson.NewObjectID(),
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
		mockRepo.On("GetAccountValidationByID", ctx, originalRule.ID).Return(originalRule, nil).Once()

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

		actionID, err := service.UpdateAccountValidationRequest(ctx, originalRule.ID.String(), updatedRule, "maker-123", "", "")
		assert.NoError(t, err)
		assert.NotEmpty(t, actionID)
		assert.Contains(t, actionID, "CPS_")
	})

	t.Run("validation rule not found", func(t *testing.T) {
		mockRepo.On("GetAccountValidationByID", ctx, originalRule.ID).Return(account_validation.ValidationRule{}, assert.AnError).Once()

		actionID, err := service.UpdateAccountValidationRequest(ctx, originalRule.ID.String(), updatedRule, "maker-123", "", "")
		assert.Error(t, err)
		assert.Empty(t, actionID)
		assert.Equal(t, "NOT_FOUND", err.Error())
	})

	t.Run("empty id", func(t *testing.T) {
		actionID, err := service.UpdateAccountValidationRequest(ctx, "", updatedRule, "maker-123", "", "")
		assert.Error(t, err)
		assert.Empty(t, actionID)
		assert.Equal(t, "INVALID_ID", err.Error())
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
		ID:            bson.NewObjectID(),
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
			ID:         validationRule.ID,
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
		mockRepo.On("UpdateAccountValidation", ctx, validationRule.ID.String(), mock.Anything).Return(nil).Once()
		mockActionRepo.On("UpdateCpsAction", ctx, mock.Anything).Return(nil).Once()

		err := service.UpdateAccountValidation(ctx, actionID, true, checkerID, "", "")
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

		err := service.UpdateAccountValidation(ctx, actionID, false, checkerID, "", "")
		assert.NoError(t, err)
	})

	t.Run("empty action ID", func(t *testing.T) {
		err := service.UpdateAccountValidation(ctx, "", true, checkerID, "", "")
		assert.Error(t, err)
		assert.Equal(t, "ACTION_ID_EMPTY", err.Error())
	})

	t.Run("empty checker ID", func(t *testing.T) {
		err := service.UpdateAccountValidation(ctx, actionID, true, "", "", "")
		assert.Error(t, err)
		assert.Equal(t, "CHECKER_ID_EMPTY", err.Error())
	})

	t.Run("action not found", func(t *testing.T) {
		mockActionRepo.On("FetchCpsActionById", ctx, actionID).Return(action.CPSAction{}, assert.AnError).Once()

		err := service.UpdateAccountValidation(ctx, actionID, true, checkerID, "", "")
		assert.Error(t, err)
		assert.Equal(t, "ACTION_NOT_FOUND", err.Error())
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

		err := service.UpdateAccountValidation(ctx, actionID, true, checkerID, "", "")
		assert.Error(t, err)
		assert.Equal(t, "ACTION_NOT_PENDING", err.Error())
	})
}
