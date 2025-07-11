package account_validation

import (
	"context"
	"encoding/json"
	"fmt"
	"time"

	// "github.com/CBE-Super-App/cbe-super-app-cps-action/internal/domain/account_validation"
	"github.com/CBE-Super-App/cbe-super-app-cps-action/internal/domain/action"

	"gitlab.com/bersufekadgetachew/cbe-super-app-shared/shared/utils"
)

type Service interface {
	GetAccountValidation(ctx context.Context, id string) (ValidationRule, error)
	UpdateAccountValidationRequest(ctx context.Context, id string, update ValidationRule, makerID string, PhoneNumber string, FullName string) (string, error)
	UpdateAccountValidation(ctx context.Context, actionID string, approve bool, checkerID string, PhoneNumber string, FullName string) error
}

type ServiceStore struct {
	repository Repository
	actionRepo action.Repository
	logger     utils.Logger
}

func NewService(repo Repository, actionRepo action.Repository, logger utils.Logger) Service {
	return &ServiceStore{
		repository: repo,
		actionRepo: actionRepo,
		logger:     logger,
	}
}

func (s *ServiceStore) GetAccountValidation(ctx context.Context, id string) (ValidationRule, error) {
	if id == "" {
		s.logger.Errorf("ID is empty")
		return ValidationRule{}, fmt.Errorf("INVALID_ID")
	}

	rule, err := s.repository.GetAccountValidationByID(ctx, id)
	if err != nil {
		s.logger.Errorf("Failed to fetch user with ID %s: %s", id, err.Error())
		return ValidationRule{}, fmt.Errorf("NOT_FOUND")
	}
	return rule, nil
}

func (s *ServiceStore) UpdateAccountValidationRequest(ctx context.Context, id string, update ValidationRule, makerID string, PhoneNumber string, FullName string) (string, error) {
	if id == "" {
		s.logger.Errorf("ID is empty")
		return "", fmt.Errorf("INVALID_ID")
	}
	originalRule, err := s.repository.GetAccountValidationByID(ctx, id)
	if err != nil {
		s.logger.Errorf("failed to fetch account validation: %v", err)
		return "", fmt.Errorf("NOT_FOUND")
	}

	previousActionJSON, err := json.Marshal(originalRule)
	if err != nil {
		s.logger.Errorf("failed to marshal previous action: %v", err)
		return "", fmt.Errorf("FAILED_TO_MARSHAL_PREVIOUS_ACTION")
	}

	type CurrentAction struct {
		Rule ValidationRule `json:"rule"`
	}
	currentAction := CurrentAction{Rule: update}
	currentActionJSON, err := json.Marshal(currentAction)
	if err != nil {
		s.logger.Errorf("failed to marshal current action: %v", err)
		return "", fmt.Errorf("FAILED_TO_MARSHAL_CURRENT_ACTION")
	}

	s.logger.Infof("original rule: %+v", originalRule)
	s.logger.Infof("new rule: %+v", update)

	actionID := utils.Random(10, &utils.PreSufix{Prefix: "CPS_"})

	a := action.CPSAction{
		ActionCode: actionID,
		Maker: action.User{
			UserID:      makerID,
			FullName:    FullName,
			PhoneNumber: PhoneNumber,
			Timestamp:   time.Now(),
		},
		Checker:         action.User{},
		Department:      update.ServiceID,
		UniqueId:        id,
		ActionType:      action.ActionUpdate,
		RequestAction:   action.RequestUpdateAccountValidation,
		ActionStatus:    action.ActionPending,
		CurrentAction:   currentActionJSON,
		PreviosAction:   previousActionJSON,
		CreatedAt:       time.Now(),
		LastModifiedAt:  time.Now(),
		RejectionReason: nil,
	}
	createdAction, err := s.actionRepo.CreateCpsAction(ctx, a)
	if err != nil {
		s.logger.Errorf("failed to create CPS action: %v", err)
		return "", fmt.Errorf("FAILED_TO_CREATE_CPS_ACTION")
	}

	s.logger.Infof("successfully created CPS action", "action_code", createdAction.ActionCode)
	return createdAction.ActionCode, nil
}

func (s *ServiceStore) UpdateAccountValidation(ctx context.Context, actionID string, approve bool, checkerID string, PhoneNumber string, FullName string) error {
	if actionID == "" {
		s.logger.Errorf("action ID is empty")
		return fmt.Errorf("ACTION_ID_EMPTY")
	}

	if checkerID == "" {
		s.logger.Errorf("checker ID is empty")
		return fmt.Errorf("CHECKER_ID_EMPTY")
	}

	s.logger.Infof("processing account validation update", "action_id", actionID, "approve", approve, "checker_id", checkerID)

	cpsAction, err := s.actionRepo.FetchCpsActionById(ctx, actionID)
	if err != nil {
		s.logger.Errorf("failed to fetch CPS action: %v", err)
		return fmt.Errorf("ACTION_NOT_FOUND")
	}

	if cpsAction.ActionStatus != action.ActionPending {
		s.logger.Errorf("action is not pending", "action_id", actionID, "status", cpsAction.ActionStatus)
		return fmt.Errorf("ACTION_NOT_PENDING")
	}
	cpsAction.Checker = action.User{
		UserID:      checkerID,
		FullName:    FullName,
		PhoneNumber: PhoneNumber,
		Timestamp:   time.Now(),
	}
	cpsAction.LastModifiedAt = time.Now()

	if approve {
		var currentAction struct {
			Rule ValidationRule `json:"rule"`
		}

		var currentActionBytes []byte
		switch v := cpsAction.CurrentAction.(type) {
		case json.RawMessage:
			currentActionBytes = v
		case []byte:
			currentActionBytes = v
		case string:
			currentActionBytes = []byte(v)
		case nil:
			s.logger.Errorf("current action is nil", "action_id", actionID)
			return fmt.Errorf("CURRENT_ACTION_NIL")
		default:

			var err error
			currentActionBytes, err = json.Marshal(v)
			if err != nil {
				s.logger.Errorf("failed to marshal current action: %v", err)
				return fmt.Errorf("CURRENT_ACTION_INVALID_TYPE")
			}
		}

		if err := json.Unmarshal(currentActionBytes, &currentAction); err != nil {
			s.logger.Errorf("failed to unmarshal current action: %v", err)
			return fmt.Errorf("FAILED_TO_UNMARSHAL_CURRENT_ACTION")
		}

		updatedRule := currentAction.Rule

		if updatedRule.ID.String() == "" || updatedRule.ID.String() != cpsAction.ID.String() {
			s.logger.Errorf("validation rule ID mismatch", "action_id", actionID, "rule_id", updatedRule.ID, "unique_id", cpsAction.ID)
			return fmt.Errorf("VALIDATION_RULE_ID_MISMATCH")
		}

		if err := s.repository.UpdateAccountValidation(ctx, updatedRule.ID.String(), updatedRule); err != nil {
			s.logger.Errorf("failed to update validation rule: %v", err)
			return fmt.Errorf("FAILED_TO_UPDATE_VALIDATION_RULE")
		}

		cpsAction.ActionStatus = action.ActionApproved
		s.logger.Infof("validation rule approved successfully", "action_id", actionID)
	} else {
		cpsAction.ActionStatus = action.ActionRejected
		cpsAction.RejectionReason = stringToPointer("Checker rejected the update")
		s.logger.Infof("validation rule update rejected", "action_id", actionID)
	}

	if err := s.actionRepo.UpdateCpsAction(ctx, cpsAction); err != nil {
		s.logger.Errorf("failed to update CPS action: %v", err)
		return fmt.Errorf("FAILED_TO_UPDATE_CPS_ACTION")
	}

	return nil
}

func stringToPointer(s string) *string {
	return &s
}
