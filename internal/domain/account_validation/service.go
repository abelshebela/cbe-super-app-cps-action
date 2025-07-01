package account_validation

import (
	"context"
	"encoding/json"
	"errors"
	"time"

	// "gitlab.com/bersufekadgetachew/cbe-super-app-cps-action/internal/domain/account_validation"
	"gitlab.com/bersufekadgetachew/cbe-super-app-cps-action/internal/domain/action"
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
		s.logger.Errorf("validation rule ID is empty")
		return ValidationRule{}, errors.New("validation rule ID cannot be empty")
	}
	s.logger.Infof("fetching account validation rule", "id", id)

	rule, err := s.repository.GetAccountValidationByID(ctx, id)
	if err != nil {
		s.logger.Errorf("failed to fetch validation rule: %v", err)
		return ValidationRule{}, err
	}
	return rule, nil
}

func (s *ServiceStore) UpdateAccountValidationRequest(ctx context.Context, id string, update ValidationRule, makerID string, PhoneNumber string, FullName string) (string, error) {
	originalRule, err := s.repository.GetAccountValidationByID(ctx, id)
	if err != nil {
		s.logger.Errorf("failed to fetch account validation: %v", err)
		return "", err
	}

	previousActionJSON, err := json.Marshal(originalRule)
	if err != nil {
		s.logger.Errorf("failed to marshal previous action: %v", err)
		return "", errors.New("failed to marshal previous action")
	}

	type CurrentAction struct {
		Rule ValidationRule `json:"rule"`
	}
	currentAction := CurrentAction{Rule: update}
	currentActionJSON, err := json.Marshal(currentAction)
	if err != nil {
		s.logger.Errorf("failed to marshal current action: %v", err)
		return "", errors.New("failed to marshal current action")
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
		return "", err
	}

	s.logger.Infof("successfully created CPS action", "action_code", createdAction.ActionCode)
	return createdAction.ActionCode, nil
}

func (s *ServiceStore) UpdateAccountValidation(ctx context.Context, actionID string, approve bool, checkerID string, PhoneNumber string, FullName string) error {
	if actionID == "" {
		s.logger.Errorf("action ID is empty")
		return errors.New("action ID cannot be empty")
	}

	if checkerID == "" {
		s.logger.Errorf("checker ID is empty")
		return errors.New("checker ID cannot be empty")
	}

	s.logger.Infof("processing account validation update", "action_id", actionID, "approve", approve, "checker_id", checkerID)

	cpsAction, err := s.actionRepo.FetchCpsActionById(ctx, actionID)
	if err != nil {
		s.logger.Errorf("failed to fetch CPS action: %v", err)
		return errors.New("action not found")
	}

	if cpsAction.ActionStatus != action.ActionPending {
		s.logger.Errorf("action is not pending", "action_id", actionID, "status", cpsAction.ActionStatus)
		return errors.New("action is not pending")
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
			return errors.New("current action is nil")
		default:

			var err error
			currentActionBytes, err = json.Marshal(v)
			if err != nil {
				s.logger.Errorf("failed to marshal current action: %v", err)
				return errors.New("current action is not a valid type")
			}
		}

		if err := json.Unmarshal(currentActionBytes, &currentAction); err != nil {
			s.logger.Errorf("failed to unmarshal current action: %v", err)
			return errors.New("failed to unmarshal current action")
		}

		updatedRule := currentAction.Rule
		if err := validateValidationRule(updatedRule); err != nil {
			s.logger.Errorf("validation rule validation failed: %v", err)
			return errors.New("validation failed: " + err.Error())
		}

		if updatedRule.ID == "" || updatedRule.ID != cpsAction.ID {
			s.logger.Errorf("validation rule ID mismatch", "action_id", actionID, "rule_id", updatedRule.ID, "unique_id", cpsAction.ID)
			return errors.New("validation rule ID mismatch")
		}

		if err := s.repository.UpdateAccountValidation(ctx, updatedRule.ID, updatedRule); err != nil {
			s.logger.Errorf("failed to update validation rule: %v", err)
			return errors.New("failed to update validation rule")
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
		return errors.New("failed to update CPS action")
	}

	return nil
}

func validateValidationRule(rule ValidationRule) error {
	if rule.Identifier == "" {
		return errors.New("identifier cannot be empty")
	}
	if rule.MinLength > rule.MaxLength {
		return errors.New("min length cannot exceed max length")
	}
	if rule.ServiceID == "" {
		return errors.New("service ID cannot be empty")
	}
	return nil
}

func stringToPointer(s string) *string {
	return &s
}
