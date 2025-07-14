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
	service := account_validation.NewAccountValidationService(mockRepo, mockActionRepo, logger)

	ctx := context.Background()
	expectedRule := account_validation.ValidationRule{
		ID:            bson.NewObjectID().Hex(),
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

		rule, err := service.GetAccountValidation(ctx, expectedRule.ID)
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
	service := account_validation.NewAccountValidationService(mockRepo, mockActionRepo, logger)

	ctx := context.Background()
	originalRule := account_validation.ValidationRule{
		ID:            bson.NewObjectID().Hex(),
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
		mockRepo.On("FetchPendingActionsByUniqueID", ctx, "test-id").Return([]action.ActionResponse{}, nil).Once()


		expectedAction := action.CPSAction{
			ActionCode:         "CPS_123",
			MakerID:            "maker-123",
			MakerName:          "Test Maker",
			MakerPhoneNumber:   "1234567890",
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


		actionID, err := service.UpdateAccountValidationRequest(ctx, "test-id", updatedRule, "maker-123", "1234567890", "Test Maker", "test-department")

		assert.NoError(t, err)
		assert.NotEmpty(t, actionID)
		assert.Contains(t, actionID, "CPS_")
	})

	t.Run("validation rule not found", func(t *testing.T) {
		mockRepo.On("GetAccountValidationByID", ctx, originalRule.ID).Return(account_validation.ValidationRule{}, assert.AnError).Once()


		actionID, err := service.UpdateAccountValidationRequest(ctx, "test-id", updatedRule, "maker-123", "1234567890", "Test Maker", "test-department")

		assert.Error(t, err)
		assert.Empty(t, actionID)
		assert.Equal(t, "NOT_FOUND", err.Error())
	})

	t.Run("empty id", func(t *testing.T) {
		actionID, err := service.UpdateAccountValidationRequest(ctx, "", updatedRule, "maker-123", "1234567890", "Test Maker", "test-department")
		assert.Error(t, err)
		assert.Empty(t, actionID)
		assert.Equal(t, "INVALID_ID", err.Error())
	})

	t.Run("pending action exists", func(t *testing.T) {
		mockRepo.On("GetAccountValidationByID", ctx, "test-id").Return(originalRule, nil).Once()
		mockRepo.On("FetchPendingActionsByUniqueID", ctx, "test-id").Return([]action.ActionResponse{{ID: "1", ActionId: "CPS_456"}}, nil).Once()

		actionID, err := service.UpdateAccountValidationRequest(ctx, "test-id", updatedRule, "maker-123", "1234567890", "Test Maker", "test-department")
		assert.Error(t, err)
		assert.Empty(t, actionID)
		assert.Equal(t, "PENDING_ACTION_EXISTS", err.Error())
	})

	t.Run("failed to fetch pending actions", func(t *testing.T) {
		mockRepo.On("GetAccountValidationByID", ctx, "test-id").Return(originalRule, nil).Once()
		mockRepo.On("FetchPendingActionsByUniqueID", ctx, "test-id").Return([]action.ActionResponse{}, assert.AnError).Once()

		actionID, err := service.UpdateAccountValidationRequest(ctx, "test-id", updatedRule, "maker-123", "1234567890", "Test Maker", "test-department")
		assert.Error(t, err)
		assert.Empty(t, actionID)
		assert.Equal(t, "FAILED_TO_FETCH_PENDING_ACTIONS", err.Error())
	})
}

func TestUpdateAccountValidation(t *testing.T) {
	ctx := context.Background()
	actionID := "CPS_123"
	checkerID := "checker-123"

	validationRule := account_validation.ValidationRule{
		ID:            bson.NewObjectID().Hex(),
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
		service := account_validation.NewAccountValidationService(mockRepo, mockActionRepo, logger)

		expectedAction := action.CPSAction{
			ID:                 validationRule.ID,
			ActionCode:         actionID,
			MakerID:            "maker-123",
			MakerName:          "Test Maker",
			MakerPhoneNumber:   "1234567890",
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
		mockRepo.On("UpdateAccountValidation", ctx, validationRule.ID, mock.Anything).Return(nil).Once()
		mockActionRepo.On("UpdateCpsAction", ctx, mock.Anything).Return(nil).Once()

		err := service.UpdateAccountValidation(ctx, actionID, utils.DecisionApproved, checkerID, "0987654321", "Test Checker", "")
		assert.NoError(t, err)
	})

	t.Run("successful rejection", func(t *testing.T) {
		mockRepo := new(mocks.AccountValidationRepository)
		mockActionRepo := new(mocks.ActionRepository)
		logger := &NoOpLogger{}
		service := account_validation.NewAccountValidationService(mockRepo, mockActionRepo, logger)

		expectedAction := action.CPSAction{
			ActionCode:         actionID,
			MakerID:            "maker-123",
			MakerName:          "Test Maker",
			MakerPhoneNumber:   "1234567890",
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

		err := service.UpdateAccountValidation(ctx, actionID, utils.DecisionDenied, checkerID, "0987654321", "Test Checker", "Rejection reason")
		assert.NoError(t, err)
	})

	t.Run("empty action ID or checker ID", func(t *testing.T) {
		mockRepo := new(mocks.AccountValidationRepository)
		mockActionRepo := new(mocks.ActionRepository)
		logger := &NoOpLogger{}
		service := account_validation.NewAccountValidationService(mockRepo, mockActionRepo, logger)

		err := service.UpdateAccountValidation(ctx, "", utils.DecisionApproved, checkerID, "0987654321", "Test Checker", "")
		assert.Error(t, err)
		assert.Equal(t, "ACTION_ID_OR_CHECKER_ID_EMPTY", err.Error())

		err = service.UpdateAccountValidation(ctx, actionID, utils.DecisionApproved, "", "0987654321", "Test Checker", "")
		assert.Error(t, err)
		assert.Equal(t, "ACTION_ID_OR_CHECKER_ID_EMPTY", err.Error())
	})

	t.Run("action not found", func(t *testing.T) {
		mockRepo := new(mocks.AccountValidationRepository)
		mockActionRepo := new(mocks.ActionRepository)
		logger := &NoOpLogger{}
		service := account_validation.NewAccountValidationService(mockRepo, mockActionRepo, logger)

		mockActionRepo.On("FetchCpsActionById", ctx, actionID).Return(action.CPSAction{}, assert.AnError).Once()

		err := service.UpdateAccountValidation(ctx, actionID, utils.DecisionApproved, checkerID, "0987654321", "Test Checker", "")
		assert.Error(t, err)
		assert.Equal(t, "ACTION_NOT_FOUND", err.Error())
	})

	t.Run("action not pending", func(t *testing.T) {
		mockRepo := new(mocks.AccountValidationRepository)
		mockActionRepo := new(mocks.ActionRepository)
		logger := &NoOpLogger{}
		service := account_validation.NewAccountValidationService(mockRepo, mockActionRepo, logger)

		expectedAction := action.CPSAction{
			ActionCode:         actionID,
			MakerID:            "maker-123",
			MakerName:          "Test Maker",
			MakerPhoneNumber:   "1234567890",
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

		err := service.UpdateAccountValidation(ctx, actionID, utils.DecisionApproved, checkerID, "0987654321", "Test Checker", "")
		assert.Error(t, err)
		assert.Equal(t, "ACTION_NOT_PENDING", err.Error())
	})

	t.Run("CPS action creation error", func(t *testing.T) {
		mockRepo := new(mocks.AccountValidationRepository)
		mockActionRepo := new(mocks.ActionRepository)
		logger := &NoOpLogger{}
		service := account_validation.NewAccountValidationService(mockRepo, mockActionRepo, logger)

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
		mockRepo.On("FetchPendingActionsByUniqueID", ctx, "test-id").Return([]action.ActionResponse{}, nil).Once()
		mockActionRepo.On("CreateCpsAction", ctx, mock.Anything).Return(action.CPSAction{}, assert.AnError).Once()

		actionID, err := service.UpdateAccountValidationRequest(ctx, "test-id", testUpdatedRule, "maker-123", "1234567890", "Test Maker", "test-department")
		assert.Error(t, err)
		assert.Empty(t, actionID)
		assert.Equal(t, "FAILED_TO_CREATE_CPS_ACTION", err.Error())
	})

	t.Run("CPS action update error", func(t *testing.T) {
		mockRepo := new(mocks.AccountValidationRepository)
		mockActionRepo := new(mocks.ActionRepository)
		logger := &NoOpLogger{}
		service := account_validation.NewAccountValidationService(mockRepo, mockActionRepo, logger)

		expectedAction := action.CPSAction{
			ID:                 validationRule.ID, // Set the ID to match the validation rule
			ActionCode:         actionID,
			MakerID:            "maker-123",
			MakerName:          "Test Maker",
			MakerPhoneNumber:   "1234567890",
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

		err := service.UpdateAccountValidation(ctx, actionID, utils.DecisionApproved, checkerID, "0987654321", "Test Checker", "")
		assert.Error(t, err)
		assert.Equal(t, "FAILED_TO_UPDATE_CPS_ACTION", err.Error())
	})

	t.Run("validation rule update error", func(t *testing.T) {
		mockRepo := new(mocks.AccountValidationRepository)
		mockActionRepo := new(mocks.ActionRepository)
		logger := &NoOpLogger{}
		service := account_validation.NewAccountValidationService(mockRepo, mockActionRepo, logger)

		expectedAction := action.CPSAction{
			ID:                 validationRule.ID,
			ActionCode:         actionID,
			MakerID:            "maker-123",
			MakerName:          "Test Maker",
			MakerPhoneNumber:   "1234567890",
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

		err := service.UpdateAccountValidation(ctx, actionID, utils.DecisionApproved, checkerID, "0987654321", "Test Checker", "")
		assert.Error(t, err)
		assert.Equal(t, "FAILED_TO_UPDATE_VALIDATION_RULE", err.Error())
	})

	t.Run("validation rule ID empty", func(t *testing.T) {
		mockRepo := new(mocks.AccountValidationRepository)
		mockActionRepo := new(mocks.ActionRepository)
		logger := &NoOpLogger{}
		service := account_validation.NewAccountValidationService(mockRepo, mockActionRepo, logger)

		// Create a rule with empty ID
		emptyIDRule := validationRule
		emptyIDRule.ID = ""
		emptyIDActionBytes, _ := json.Marshal(CurrentAction{Rule: emptyIDRule})

		expectedAction := action.CPSAction{
			ID:                 "test-id",
			ActionCode:         actionID,
			MakerID:            "maker-123",
			MakerName:          "Test Maker",
			MakerPhoneNumber:   "1234567890",
			CheckerID:          "",
			CheckerName:        "",
			CheckerPhoneNumber: "",
			Department:         "test-service",
			ActionType:         action.ActionUpdate,
			RequestAction:      action.RequestUpdateAccountValidation,
			ActionStatus:       action.ActionPending,
			CurrentAction:      emptyIDActionBytes,
			CreatedAt:          time.Now(),
			LastModifiedAt:     time.Now(),
			RejectionReason:    nil,
		}

		mockActionRepo.On("FetchCpsActionById", ctx, actionID).Return(expectedAction, nil).Once()

		err := service.UpdateAccountValidation(ctx, actionID, utils.DecisionApproved, checkerID, "0987654321", "Test Checker", "")
		assert.Error(t, err)
		assert.Equal(t, "VALIDATION_RULE_ID_EMPTY", err.Error())
	})

	t.Run("invalid decision", func(t *testing.T) {
		mockRepo := new(mocks.AccountValidationRepository)
		mockActionRepo := new(mocks.ActionRepository)
		logger := &NoOpLogger{}
		service := account_validation.NewAccountValidationService(mockRepo, mockActionRepo, logger)

		expectedAction := action.CPSAction{
			ActionCode:         actionID,
			MakerID:            "maker-123",
			MakerName:          "Test Maker",
			MakerPhoneNumber:   "1234567890",
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

		err := service.UpdateAccountValidation(ctx, actionID, "INVALID", checkerID, "0987654321", "Test Checker", "")
		assert.Error(t, err)
		assert.Equal(t, "INVALID_DECISION", err.Error())
	})
}
