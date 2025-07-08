package account_validation

import (
	"context"
	"encoding/json"
	"fmt"
	"reflect"
	"time"

	// "github.com/CBE-Super-App/cbe-super-app-cps-action/internal/domain/account_validation"
	"github.com/CBE-Super-App/cbe-super-app-cps-action/internal/domain/action"
	utils "github.com/CBE-Super-App/cbe-super-app-cps-action/pkgs/utils"
	sharedutils "gitlab.com/bersufekadgetachew/cbe-super-app-shared/shared/utils"
	"go.mongodb.org/mongo-driver/bson"
)

type Service interface {
	GetAccountValidation(ctx context.Context, id string) (ValidationRule, error)
	UpdateAccountValidationRequest(ctx context.Context, id string, update ValidationRule, makerID string, PhoneNumber string, FullName string, Department string) (string, error)
	UpdateAccountValidation(ctx context.Context, actionID string, decision utils.DecisonEnum, checkerID string, PhoneNumber string, FullName string, rejectedReason string) error
}

type ServiceStore struct {
	repository Repository
	actionRepo action.Repository
	logger     sharedutils.Logger
}

func NewService(repo Repository, actionRepo action.Repository, logger sharedutils.Logger) Service {
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

func (s *ServiceStore) UpdateAccountValidationRequest(
	ctx context.Context,
	id string,
	update ValidationRule,
	makerID, PhoneNumber, FullName, Department string,
) (string, error) {
	if id == "" {
		s.logger.Errorf("ID is empty")
		return "", fmt.Errorf("INVALID_ID")
	}

	originalRule, err := s.repository.GetAccountValidationByID(ctx, id)
	if err != nil {
		s.logger.Errorf("failed to fetch account validation: %v", err)
		return "", fmt.Errorf("NOT_FOUND")
	}

	pendingActions, err := s.repository.FetchPendingActionsByUniqueID(ctx, id)
	if err != nil {
		s.logger.Errorf("failed to fetch pending actions: %v", err)
		return "", fmt.Errorf("FAILED_TO_FETCH_PENDING_ACTIONS")
	}
	if len(pendingActions) > 0 {
		s.logger.Errorf("pending action already exists for validation rule: %s", id)
		return "", fmt.Errorf("PENDING_ACTION_EXISTS")
	}

	// Debug log to check update
	fmt.Printf("update being marshaled: %+v\n", update)

	previousActionJSON, err := json.Marshal(originalRule)
	if err != nil {
		s.logger.Errorf("failed to marshal previous action: %v", err)
		return "", fmt.Errorf("FAILED_TO_MARSHAL_PREVIOUS_ACTION")
	}

	// Ensure the update struct has the ID set
	update.ID = id

	type CurrentAction struct {
		Rule ValidationRule `json:"rule"`
	}
	currentAction := CurrentAction{Rule: update}

	actionID := sharedutils.Random(10, &sharedutils.PreSufix{Prefix: "CPS_"})

	a := action.CPSAction{
		ActionCode:         actionID,
		MakerID:            makerID,
		MakerName:          FullName,
		MakerPhoneNumber:   PhoneNumber,
		CheckerID:          "",
		CheckerName:        "",
		CheckerPhoneNumber: "",
		Department:         Department,
		UniqueId:           id,
		ActionType:         action.ActionUpdate,
		RequestAction:      action.RequestUpdateAccountValidation,
		ActionStatus:       action.ActionPending,
		CurrentAction:      currentAction, // assign struct directly
		PreviosAction:      previousActionJSON,
		CreatedAt:          time.Now(),
		LastModifiedAt:     time.Now(),
		RejectionReason:    nil,
		MakerActionTime:    time.Now(),
		CheckerActionTime:  time.Time{},
	}

	createdAction, err := s.actionRepo.CreateCpsAction(ctx, a)
	if err != nil {
		s.logger.Errorf("failed to create CPS action: %v", err)
		return "", fmt.Errorf("FAILED_TO_CREATE_CPS_ACTION")
	}

	return createdAction.ActionCode, nil
}

func toJSONBytes(val interface{}) ([]byte, error) {
	switch v := val.(type) {
	case nil:
		return nil, fmt.Errorf("value is nil")
	case json.RawMessage:
		return v, nil
	case []byte:
		return v, nil
	case string:
		return []byte(v), nil
	case map[string]interface{}, []interface{}:
		return json.Marshal(v)
	case bson.D:
		// Convert bson.D to bson.M using bson.Marshal + bson.Unmarshal (recommended)
		bsonBytes, err := bson.Marshal(v)
		if err != nil {
			return nil, err
		}
		var m bson.M
		if err := bson.Unmarshal(bsonBytes, &m); err != nil {
			return nil, err
		}
		return json.Marshal(m)
	default:
		rv := reflect.ValueOf(v)
		if rv.Kind() == reflect.Struct {
			return json.Marshal(v)
		}
		return nil, fmt.Errorf("unsupported type: %T", v)
	}
}

func (s *ServiceStore) UpdateAccountValidation(
	ctx context.Context,
	actionID string,
	decision utils.DecisonEnum,
	checkerID, PhoneNumber, FullName, rejectedReason string,
) error {
	if actionID == "" || checkerID == "" {
		s.logger.Errorf("action ID or checker ID is empty")
		return fmt.Errorf("ACTION_ID_OR_CHECKER_ID_EMPTY")
	}

	cpsAction, err := s.actionRepo.FetchCpsActionById(ctx, actionID)
	if err != nil {
		s.logger.Errorf("failed to fetch CPS action: %v", err)
		return fmt.Errorf("ACTION_NOT_FOUND")
	}

	if cpsAction.ActionStatus != action.ActionPending {
		s.logger.Errorf("action is not pending")
		return fmt.Errorf("ACTION_NOT_PENDING")
	}

	cpsAction.CheckerID = checkerID
	cpsAction.CheckerName = FullName
	cpsAction.CheckerPhoneNumber = PhoneNumber
	cpsAction.CheckerActionTime = time.Now()
	cpsAction.LastModifiedAt = time.Now()

	if decision == utils.DecisionApproved {
		var currentAction struct {
			Rule ValidationRule `json:"rule"`
		}

		currentActionBytes, err := toJSONBytes(cpsAction.CurrentAction)
		if err != nil {
			s.logger.Errorf("failed to convert CurrentAction to JSON: %v", err)
			return fmt.Errorf("FAILED_TO_CONVERT_CURRENT_ACTION")
		}
		fmt.Printf("CurrentAction raw JSON: %s\n", string(currentActionBytes))
		if err := json.Unmarshal(currentActionBytes, &currentAction); err != nil {
			s.logger.Errorf("failed to unmarshal current action: %v, bytes: %s", err, string(currentActionBytes))
			return fmt.Errorf("FAILED_TO_UNMARSHAL_CURRENT_ACTION")
		}

		updatedRule := currentAction.Rule
		fmt.Printf("Unmarshalled currentAction: %+v\n", currentAction)
		if updatedRule.ID == "" {
			s.logger.Errorf("invalid validation rule ID: empty")
			return fmt.Errorf("VALIDATION_RULE_ID_EMPTY")
		}

		if err := s.repository.UpdateAccountValidation(ctx, updatedRule.ID, updatedRule); err != nil {
			s.logger.Errorf("failed to update validation rule: %v", err)
			return fmt.Errorf("FAILED_TO_UPDATE_VALIDATION_RULE")
		}

		cpsAction.ActionStatus = action.ActionApproved
		s.logger.Infof("validation rule approved successfully")
	} else if decision == utils.DecisionDenied {
		cpsAction.ActionStatus = action.ActionRejected
		if rejectedReason != "" {
			cpsAction.RejectionReason = &rejectedReason
		} else {
			cpsAction.RejectionReason = stringToPointer("Checker rejected the update")
		}
	} else {
		return fmt.Errorf("INVALID_DECISION")
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
