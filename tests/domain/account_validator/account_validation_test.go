package account_validator_test

import (
	"context"
	"encoding/json"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"

	"github.com/CBE-Super-App/cbe-super-app-cps-action/internal/domain/account_validation"
	"github.com/CBE-Super-App/cbe-super-app-cps-action/internal/domain/action"
	"github.com/CBE-Super-App/cbe-super-app-cps-action/mocks"
	utils "github.com/CBE-Super-App/cbe-super-app-cps-action/pkgs/utils"
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
			ActionCode:         "CPS_123",
			MakerID:            "maker-123",
			MakerName:          "",
			MakerPhoneNumber:   "",
			CheckerID:          "",
			CheckerName:        "",
			CheckerPhoneNumber: "",
			Department:         "test-department",
			ActionType:         action.ActionUpdate,
			RequestAction:      action.RequestUpdateAccountValidation,
			ActionStatus:       action.ActionPending,
			CreatedAt:          time.Now(),
			LastModifiedAt:     time.Now(),
			RejectionReason:    nil,
		}

		mockActionRepo.On("CreateCpsAction", ctx, mock.Anything).Return(expectedAction, nil).Once()

		actionID, err := service.UpdateAccountValidationRequest(ctx, "test-id", updatedRule, "maker-123", "", "", "test-department")
		assert.NoError(t, err)
		assert.NotEmpty(t, actionID)
		assert.Contains(t, actionID, "CPS_")
	})

	t.Run("validation rule not found", func(t *testing.T) {
		mockRepo.On("GetAccountValidationByID", ctx, "test-id").Return(account_validation.ValidationRule{}, assert.AnError).Once()

		actionID, err := service.UpdateAccountValidationRequest(ctx, "test-id", updatedRule, "maker-123", "", "", "test-department")
		assert.Error(t, err)
		assert.Empty(t, actionID)
		assert.Equal(t, "NOT_FOUND", err.Error())
	})

	t.Run("empty id", func(t *testing.T) {
		actionID, err := service.UpdateAccountValidationRequest(ctx, "", updatedRule, "maker-123", "", "", "test-department")
		assert.Error(t, err)
		assert.Empty(t, actionID)
		assert.Equal(t, "INVALID_ID", err.Error())
	})
}

func TestUpdateAccountValidation(t *testing.T) {
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
		mockRepo := new(mocks.AccountValidationRepository)
		mockActionRepo := new(mocks.ActionRepository)
		logger := &NoOpLogger{}
		service := account_validation.NewService(mockRepo, mockActionRepo, logger)

		expectedAction := action.CPSAction{
			ID:                 validationRule.ID,
			ActionCode:         actionID,
			MakerID:            "maker-123",
			MakerName:          "",
			MakerPhoneNumber:   "",
			CheckerID:          "",
			CheckerName:        "",
			CheckerPhoneNumber: "",
			Department:         "test-service",
			ActionType:         action.ActionUpdate,
			RequestAction:      action.RequestUpdateAccountValidation,
			ActionStatus:       action.ActionPending,
			CurrentAction:      currentActionBytes,
			CreatedAt:          time.Now(),
			LastModifiedAt:     time.Now(),
			RejectionReason:    nil,
		}

		mockActionRepo.On("FetchCpsActionById", ctx, actionID).Return(expectedAction, nil).Once()
		mockRepo.On("UpdateAccountValidation", ctx, "test-id", mock.Anything).Return(nil).Once()
		mockActionRepo.On("UpdateCpsAction", ctx, mock.Anything).Return(nil).Once()

		err := service.UpdateAccountValidation(ctx, actionID, utils.DecisionApproved, checkerID, "", "", "")
		assert.NoError(t, err)
	})

	t.Run("successful rejection", func(t *testing.T) {
		mockRepo := new(mocks.AccountValidationRepository)
		mockActionRepo := new(mocks.ActionRepository)
		logger := &NoOpLogger{}
		service := account_validation.NewService(mockRepo, mockActionRepo, logger)

		expectedAction := action.CPSAction{
			ActionCode:         actionID,
			MakerID:            "maker-123",
			MakerName:          "",
			MakerPhoneNumber:   "",
			CheckerID:          "",
			CheckerName:        "",
			CheckerPhoneNumber: "",
			Department:         "test-service",
			ActionType:         action.ActionUpdate,
			RequestAction:      action.RequestUpdateAccountValidation,
			ActionStatus:       action.ActionPending,
			CurrentAction:      currentActionBytes,
			CreatedAt:          time.Now(),
			LastModifiedAt:     time.Now(),
			RejectionReason:    nil,
		}

		mockActionRepo.On("FetchCpsActionById", ctx, actionID).Return(expectedAction, nil).Once()
		mockActionRepo.On("UpdateCpsAction", ctx, mock.Anything).Return(nil).Once()

		err := service.UpdateAccountValidation(ctx, actionID, utils.DecisionDenied, checkerID, "", "", "Rejection reason")
		assert.NoError(t, err)
	})

	t.Run("empty action ID", func(t *testing.T) {
		mockRepo := new(mocks.AccountValidationRepository)
		mockActionRepo := new(mocks.ActionRepository)
		logger := &NoOpLogger{}
		service := account_validation.NewService(mockRepo, mockActionRepo, logger)

		err := service.UpdateAccountValidation(ctx, "", utils.DecisionApproved, checkerID, "", "", "")
		assert.Error(t, err)
		assert.Equal(t, "ACTION_ID_EMPTY", err.Error())
	})

	t.Run("empty checker ID", func(t *testing.T) {
		mockRepo := new(mocks.AccountValidationRepository)
		mockActionRepo := new(mocks.ActionRepository)
		logger := &NoOpLogger{}
		service := account_validation.NewService(mockRepo, mockActionRepo, logger)

		err := service.UpdateAccountValidation(ctx, actionID, utils.DecisionApproved, "", "", "", "")
		assert.Error(t, err)
		assert.Equal(t, "CHECKER_ID_EMPTY", err.Error())
	})

	t.Run("action not found", func(t *testing.T) {
		mockRepo := new(mocks.AccountValidationRepository)
		mockActionRepo := new(mocks.ActionRepository)
		logger := &NoOpLogger{}
		service := account_validation.NewService(mockRepo, mockActionRepo, logger)

		mockActionRepo.On("FetchCpsActionById", ctx, actionID).Return(action.CPSAction{}, assert.AnError).Once()

		err := service.UpdateAccountValidation(ctx, actionID, utils.DecisionApproved, checkerID, "", "", "")
		assert.Error(t, err)
		assert.Equal(t, "ACTION_NOT_FOUND", err.Error())
	})

	t.Run("action not pending", func(t *testing.T) {
		mockRepo := new(mocks.AccountValidationRepository)
		mockActionRepo := new(mocks.ActionRepository)
		logger := &NoOpLogger{}
		service := account_validation.NewService(mockRepo, mockActionRepo, logger)

		expectedAction := action.CPSAction{
			ActionCode:         actionID,
			MakerID:            "maker-123",
			MakerName:          "",
			MakerPhoneNumber:   "",
			CheckerID:          "",
			CheckerName:        "",
			CheckerPhoneNumber: "",
			Department:         "test-service",
			ActionType:         action.ActionUpdate,
			RequestAction:      action.RequestUpdateAccountValidation,
			ActionStatus:       action.ActionApproved,
			CurrentAction:      currentActionBytes,
			CreatedAt:          time.Now(),
			LastModifiedAt:     time.Now(),
			RejectionReason:    nil,
		}

		mockActionRepo.On("FetchCpsActionById", ctx, actionID).Return(expectedAction, nil).Once()

		err := service.UpdateAccountValidation(ctx, actionID, utils.DecisionApproved, checkerID, "", "", "")
		assert.Error(t, err)
		assert.Equal(t, "ACTION_NOT_PENDING", err.Error())
	})

	t.Run("CPS action creation error", func(t *testing.T) {
		mockRepo := new(mocks.AccountValidationRepository)
		mockActionRepo := new(mocks.ActionRepository)
		logger := &NoOpLogger{}
		service := account_validation.NewService(mockRepo, mockActionRepo, logger)

		// Define test data for this specific test
		testOriginalRule := account_validation.ValidationRule{
			ID:            "test-id",
			EntityType:    "account",
			ValidationFor: "phone",
			Identifier:    "phone_number",
			MinLength:     10,
			MaxLength:     15,
			Enabled:       true,
			ServiceID:     "test-service",
		}
		testUpdatedRule := testOriginalRule
		testUpdatedRule.MinLength = 11

		mockRepo.On("GetAccountValidationByID", ctx, "test-id").Return(testOriginalRule, nil).Once()
		mockActionRepo.On("CreateCpsAction", ctx, mock.Anything).Return(action.CPSAction{}, assert.AnError).Once()

		actionID, err := service.UpdateAccountValidationRequest(ctx, "test-id", testUpdatedRule, "maker-123", "", "", "test-department")
		assert.Error(t, err)
		assert.Empty(t, actionID)
		assert.Equal(t, "FAILED_TO_CREATE_CPS_ACTION", err.Error())
	})

	t.Run("CPS action update error", func(t *testing.T) {
		mockRepo := new(mocks.AccountValidationRepository)
		mockActionRepo := new(mocks.ActionRepository)
		logger := &NoOpLogger{}
		service := account_validation.NewService(mockRepo, mockActionRepo, logger)

		expectedAction := action.CPSAction{
			ID:                 validationRule.ID, // Set the ID to match the validation rule
			ActionCode:         actionID,
			MakerID:            "maker-123",
			MakerName:          "",
			MakerPhoneNumber:   "",
			CheckerID:          "",
			CheckerName:        "",
			CheckerPhoneNumber: "",
			Department:         "test-service",
			ActionType:         action.ActionUpdate,
			RequestAction:      action.RequestUpdateAccountValidation,
			ActionStatus:       action.ActionPending,
			CurrentAction:      currentActionBytes,
			CreatedAt:          time.Now(),
			LastModifiedAt:     time.Now(),
			RejectionReason:    nil,
		}

		mockActionRepo.On("FetchCpsActionById", ctx, actionID).Return(expectedAction, nil).Once()
		mockRepo.On("UpdateAccountValidation", ctx, "test-id", mock.Anything).Return(nil).Once()
		mockActionRepo.On("UpdateCpsAction", ctx, mock.Anything).Return(assert.AnError).Once()

		err := service.UpdateAccountValidation(ctx, actionID, utils.DecisionApproved, checkerID, "", "", "")
		assert.Error(t, err)
		assert.Equal(t, "FAILED_TO_UPDATE_CPS_ACTION", err.Error())
	})

	t.Run("validation rule update error", func(t *testing.T) {
		mockRepo := new(mocks.AccountValidationRepository)
		mockActionRepo := new(mocks.ActionRepository)
		logger := &NoOpLogger{}
		service := account_validation.NewService(mockRepo, mockActionRepo, logger)

		expectedAction := action.CPSAction{
			ID:                 validationRule.ID,
			ActionCode:         actionID,
			MakerID:            "maker-123",
			MakerName:          "",
			MakerPhoneNumber:   "",
			CheckerID:          "",
			CheckerName:        "",
			CheckerPhoneNumber: "",
			Department:         "test-service",
			ActionType:         action.ActionUpdate,
			RequestAction:      action.RequestUpdateAccountValidation,
			ActionStatus:       action.ActionPending,
			CurrentAction:      currentActionBytes,
			CreatedAt:          time.Now(),
			LastModifiedAt:     time.Now(),
			RejectionReason:    nil,
		}

		mockActionRepo.On("FetchCpsActionById", ctx, actionID).Return(expectedAction, nil).Once()
		mockRepo.On("UpdateAccountValidation", ctx, "test-id", mock.Anything).Return(assert.AnError).Once()

		err := service.UpdateAccountValidation(ctx, actionID, utils.DecisionApproved, checkerID, "", "", "")
		assert.Error(t, err)
		assert.Equal(t, "FAILED_TO_UPDATE_VALIDATION_RULE", err.Error())
	})

	t.Run("validation rule ID mismatch", func(t *testing.T) {
		mockRepo := new(mocks.AccountValidationRepository)
		mockActionRepo := new(mocks.ActionRepository)
		logger := &NoOpLogger{}
		service := account_validation.NewService(mockRepo, mockActionRepo, logger)

		// Create a rule with different ID than the action
		mismatchedRule := validationRule
		mismatchedRule.ID = "different-id"
		mismatchedActionBytes, _ := json.Marshal(CurrentAction{Rule: mismatchedRule})

		expectedAction := action.CPSAction{
			ID:                 "test-id", // Different from rule ID
			ActionCode:         actionID,
			MakerID:            "maker-123",
			MakerName:          "",
			MakerPhoneNumber:   "",
			CheckerID:          "",
			CheckerName:        "",
			CheckerPhoneNumber: "",
			Department:         "test-service",
			ActionType:         action.ActionUpdate,
			RequestAction:      action.RequestUpdateAccountValidation,
			ActionStatus:       action.ActionPending,
			CurrentAction:      mismatchedActionBytes,
			CreatedAt:          time.Now(),
			LastModifiedAt:     time.Now(),
			RejectionReason:    nil,
		}

		mockActionRepo.On("FetchCpsActionById", ctx, actionID).Return(expectedAction, nil).Once()

		err := service.UpdateAccountValidation(ctx, actionID, utils.DecisionApproved, checkerID, "", "", "")
		assert.Error(t, err)
		assert.Equal(t, "VALIDATION_RULE_ID_MISMATCH", err.Error())
	})

	t.Run("current action is nil", func(t *testing.T) {
		mockRepo := new(mocks.AccountValidationRepository)
		mockActionRepo := new(mocks.ActionRepository)
		logger := &NoOpLogger{}
		service := account_validation.NewService(mockRepo, mockActionRepo, logger)

		expectedAction := action.CPSAction{
			ActionCode:         actionID,
			MakerID:            "maker-123",
			MakerName:          "",
			MakerPhoneNumber:   "",
			CheckerID:          "",
			CheckerName:        "",
			CheckerPhoneNumber: "",
			Department:         "test-service",
			ActionType:         action.ActionUpdate,
			RequestAction:      action.RequestUpdateAccountValidation,
			ActionStatus:       action.ActionPending,
			CurrentAction:      nil, // Nil current action
			CreatedAt:          time.Now(),
			LastModifiedAt:     time.Now(),
			RejectionReason:    nil,
		}

		mockActionRepo.On("FetchCpsActionById", ctx, actionID).Return(expectedAction, nil).Once()

		err := service.UpdateAccountValidation(ctx, actionID, utils.DecisionApproved, checkerID, "", "", "")
		assert.Error(t, err)
		assert.Equal(t, "CURRENT_ACTION_NIL", err.Error())
	})
}
