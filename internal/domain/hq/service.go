package hq

import (
	"context"
	"encoding/json"
	"fmt"
	"reflect"
	"time"

	"github.com/CBE-Super-App/cbe-super-app-cps-action/internal/domain/action"
	contexts "github.com/CBE-Super-App/cbe-super-app-cps-action/pkgs/context"
	utils "github.com/CBE-Super-App/cbe-super-app-cps-action/pkgs/utils"
	constant "github.com/CBE-Super-App/cbe-super-app-cps-action/utils"

	sharedutils "gitlab.com/bersufekadgetachew/cbe-super-app-shared/shared/utils"
)

type Service interface {
	GetHQ(ctx context.Context, id string) (HQ, error)
	GetAllHQ(ctx context.Context, filerParams *constant.Filter) (*utils.PaginatedResponse[[]*HQ], error)
	UpdateBlockTimeRequest(ctx context.Context, request UpdateBlockTimeRequest) (string, error)
	UpdateArchiveTimeRequest(ctx context.Context, request UpdateArchiveTimeRequest) (string, error)
	UpdateBlockTime(ctx context.Context, request ApproveRejectRequest) error
	UpdateArchiveTime(ctx context.Context, request ApproveRejectRequest) error
}

type ServiceStore struct {
	repository Repository
	actionRepo action.ActionRepository
	logger     sharedutils.Logger
}

func NewService(repo Repository, actionRepo action.ActionRepository, logger sharedutils.Logger) Service {
	return &ServiceStore{
		repository: repo,
		actionRepo: actionRepo,
		logger:     logger,
	}
}

func (s *ServiceStore) GetAllHQ(ctx context.Context, filerParams *constant.Filter) (*utils.PaginatedResponse[[]*HQ], error) {
	return s.repository.GetAllHQ(ctx, filerParams)
}
func (s *ServiceStore) GetHQ(ctx context.Context, id string) (HQ, error) {
	if id == "" {
		s.logger.Errorf("HQ ID is empty")
		return HQ{}, fmt.Errorf("INVALID_ID")
	}
	s.logger.Infof("fetching HQ", "id", id)

	hq, err := s.repository.GetHQByID(ctx, id)
	if err != nil {
		s.logger.Errorf("failed to fetch HQ: %v", err)
		return HQ{}, fmt.Errorf("NOT_FOUND")
	}
	return hq, nil
}

func (s *ServiceStore) UpdateBlockTimeRequest(ctx context.Context, request UpdateBlockTimeRequest) (string, error) {
	if request.ID == "" {
		return "", fmt.Errorf("INVALID_ID")
	}
	originalHQ, err := s.repository.GetHQByID(ctx, request.ID)
	if err != nil {
		s.logger.Errorf("failed to fetch HQ: %v", err)
		return "", fmt.Errorf("NOT_FOUND")
	}

	// Use actionRepo to check for pending actions by unique ID (outbound/init)
	pendingActions, err := s.actionRepo.(interface {
		FetchPendingActionsByUniqueID(context.Context, string) ([]action.ActionResponse, error)
	}).FetchPendingActionsByUniqueID(ctx, request.ID)
	if err != nil {
		return "", fmt.Errorf("FAILED_TO_FETCH_PENDING_ACTIONS")
	}
	if len(pendingActions) > 0 {
		return "", fmt.Errorf("PENDING_ACTION_EXISTS")
	}

	previousActionJSON, err := json.Marshal(originalHQ)
	if err != nil {
		s.logger.Errorf("failed to marshal previous action: %v", err)
		return "", fmt.Errorf("FAILED_TO_MARSHAL_PREVIOUS_ACTION")
	}

	updatedHQ := originalHQ
	updatedHQ.BlockTime = request.BlockTime
	// Use struct wrapper for current action, matching account validation
	type CurrentAction struct {
		HQ HQ `json:"hq"`
	}
	currentAction := CurrentAction{HQ: updatedHQ}

	actionID := sharedutils.Random(10, &sharedutils.PreSufix{Prefix: "CPS_"})
	a := action.CPSAction{
		ActionCode:         actionID,
		MakerID:            request.MakerID,
		MakerName:          request.MakerName,
		MakerPhoneNumber:   request.MakerPhone,
		CheckerID:          "",
		CheckerName:        "",
		CheckerPhoneNumber: "",
		Department:         originalHQ.Name,
		ActionType:         action.ActionUpdate,
		RequestAction:      action.RequestUpdateHQBlockTime,
		ActionStatus:       action.ActionPending,
		CurrentAction:      currentAction, // pass struct, not marshaled JSON
		PreviosAction:      previousActionJSON,
		CreatedAt:          time.Now(),
		LastModifiedAt:     time.Now(),
		RejectionReason:    nil,
		MakerActionTime:    time.Now(),
		CheckerActionTime:  time.Time{},
		UniqueId:           request.ID,
	}
	createdAction, err := s.actionRepo.CreateCpsAction(ctx, a)
	if err != nil {
		s.logger.Errorf("failed to create CPS action: %v", err)
		return "", fmt.Errorf("FAILED_TO_CREATE_CPS_ACTION")
	}

	return createdAction.ActionCode, nil
}

func (s *ServiceStore) UpdateArchiveTimeRequest(ctx context.Context, request UpdateArchiveTimeRequest) (string, error) {
	if request.ID == "" {
		return "", fmt.Errorf("INVALID_ID")
	}
	originalHQ, err := s.repository.GetHQByID(ctx, request.ID)
	if err != nil {
		s.logger.Errorf("failed to fetch HQ: %v", err)
		return "", fmt.Errorf("NOT_FOUND")
	}

	// Use actionRepo to check for pending actions by unique ID (outbound/init)
	pendingActions, err := s.actionRepo.(interface {
		FetchPendingActionsByUniqueID(context.Context, string) ([]action.ActionResponse, error)
	}).FetchPendingActionsByUniqueID(ctx, request.ID)
	if err != nil {
		return "", fmt.Errorf("FAILED_TO_FETCH_PENDING_ACTIONS")
	}
	if len(pendingActions) > 0 {
		return "", fmt.Errorf("PENDING_ACTION_EXISTS")
	}

	previousActionJSON, err := json.Marshal(originalHQ)
	if err != nil {
		s.logger.Errorf("failed to marshal previous action: %v", err)
		return "", fmt.Errorf("FAILED_TO_MARSHAL_PREVIOUS_ACTION")
	}

	makerUser := contexts.ExtractContext(ctx)
	updatedHQ := originalHQ
	updatedHQ.ArchiveTime = request.ArchiveTime
	// Use struct wrapper for current action, matching account validation
	type CurrentAction struct {
		HQ HQ `json:"hq"`
	}
	currentAction := CurrentAction{HQ: updatedHQ}

	actionID := sharedutils.Random(10, &sharedutils.PreSufix{Prefix: "CPS_"})
	a := action.CPSAction{
		ActionCode:         actionID,
		MakerID:            request.MakerID,
		MakerName:          request.MakerName,
		MakerPhoneNumber:   request.MakerPhone,
		CheckerID:          "",
		CheckerName:        "",
		CheckerPhoneNumber: "",
		Department:         makerUser.Department,
		ActionType:         action.ActionUpdate,
		RequestAction:      action.RequestUpdateHQArchiveTime,
		ActionStatus:       action.ActionPending,
		CurrentAction:      currentAction, // pass struct, not marshaled JSON
		PreviosAction:      previousActionJSON,
		CreatedAt:          time.Now(),
		LastModifiedAt:     time.Now(),
		RejectionReason:    nil,
		MakerActionTime:    time.Now(),
		CheckerActionTime:  time.Time{},
		UniqueId:           request.ID,
	}
	createdAction, err := s.actionRepo.CreateCpsAction(ctx, a)
	if err != nil {
		s.logger.Errorf("failed to create CPS action: %v", err)
		return "", fmt.Errorf("FAILED_TO_CREATE_CPS_ACTION")
	}

	return createdAction.ActionCode, nil
}

// Add helper for robust JSON extraction, matching account validation
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
	default:
		rv := reflect.ValueOf(v)
		if rv.Kind() == reflect.Struct {
			return json.Marshal(v)
		}
		return nil, fmt.Errorf("unsupported type: %T", v)
	}
}

func (s *ServiceStore) UpdateBlockTime(ctx context.Context, request ApproveRejectRequest) error {
	if request.ActionCode == "" {
		s.logger.Errorf("action ID is empty")
		return fmt.Errorf("ACTION_ID_EMPTY")
	}

	if request.CheckerID == "" {
		s.logger.Errorf("checker ID is empty")
		return fmt.Errorf("CHECKER_ID_EMPTY")
	}

	s.logger.Infof("processing HQ block time update", "action_id", request.ActionCode, "decision", request.Decision, "checker_id", request.CheckerID)

	cpsAction, err := s.actionRepo.FetchCpsActionById(ctx, request.ActionCode)
	if err != nil {
		s.logger.Errorf("failed to fetch CPS action: %v", err)
		return fmt.Errorf("ACTION_NOT_FOUND")
	}

	if cpsAction.ActionStatus != action.ActionPending {
		s.logger.Errorf("action is not pending", "action_id", request.ActionCode, "status", cpsAction.ActionStatus)
		return fmt.Errorf("ACTION_NOT_PENDING")
	}

	cpsAction.CheckerID = request.CheckerID
	cpsAction.CheckerName = request.CheckerName
	cpsAction.CheckerPhoneNumber = request.CheckerPhone
	cpsAction.CheckerActionTime = time.Now()
	cpsAction.LastModifiedAt = time.Now()

	switch request.Decision {
	case utils.DecisionApproved:
		// Use struct wrapper for current action, matching account validation
		type CurrentAction struct {
			HQ HQ `json:"hq"`
		}
		var currentAction CurrentAction
		currentActionBytes, err := toJSONBytes(cpsAction.CurrentAction)
		if err != nil {
			s.logger.Errorf("failed to convert CurrentAction to JSON: %v", err)
			return fmt.Errorf("FAILED_TO_CONVERT_CURRENT_ACTION")
		}
		if err := json.Unmarshal(currentActionBytes, &currentAction); err != nil {
			s.logger.Errorf("failed to unmarshal current action: %v, bytes: %s", err, string(currentActionBytes))
			return fmt.Errorf("FAILED_TO_UNMARSHAL_CURRENT_ACTION")
		}
		updatedHQ := currentAction.HQ
		fmt.Println("updated hq  ", updatedHQ.ID)

		if err := s.repository.UpdateHQ(ctx, updatedHQ.ID.Hex(), updatedHQ); err != nil {
			s.logger.Errorf("failed to update HQ: %v", err)
			return fmt.Errorf("FAILED_TO_UPDATE_HQ")
		}
		cpsAction.ActionStatus = action.ActionApproved
		s.logger.Infof("HQ block time approved successfully", "action_id", request.ActionCode)
	case utils.DecisionDenied:
		cpsAction.ActionStatus = action.ActionRejected
		if request.RejectedReason != "" {
			cpsAction.RejectionReason = &request.RejectedReason
		}
		s.logger.Infof("HQ block time update rejected", "action_id", request.ActionCode)
	default:
		return fmt.Errorf("INVALID_DECISION")
	}

	if err := s.actionRepo.UpdateCpsAction(ctx, cpsAction); err != nil {
		s.logger.Errorf("failed to update CPS action: %v", err)
		return fmt.Errorf("FAILED_TO_UPDATE_CPS_ACTION")
	}

	return nil
}

func (s *ServiceStore) UpdateArchiveTime(ctx context.Context, request ApproveRejectRequest) error {
	if request.ActionCode == "" {
		s.logger.Errorf("action ID is empty")
		return fmt.Errorf("ACTION_ID_EMPTY")
	}

	if request.CheckerID == "" {
		s.logger.Errorf("checker ID is empty")
		return fmt.Errorf("CHECKER_ID_EMPTY")
	}

	s.logger.Infof("processing HQ archive time update", "action_id", request.ActionCode, "decision", request.Decision, "checker_id", request.CheckerID)

	cpsAction, err := s.actionRepo.FetchCpsActionById(ctx, request.ActionCode)
	if err != nil {
		s.logger.Errorf("failed to fetch CPS action: %v", err)
		return fmt.Errorf("ACTION_NOT_FOUND")
	}

	if cpsAction.ActionStatus != action.ActionPending {
		s.logger.Errorf("action is not pending", "action_id", request.ActionCode, "status", cpsAction.ActionStatus)
		return fmt.Errorf("ACTION_NOT_PENDING")
	}

	cpsAction.CheckerID = request.CheckerID
	cpsAction.CheckerName = request.CheckerName
	cpsAction.CheckerPhoneNumber = request.CheckerPhone
	cpsAction.CheckerActionTime = time.Now()
	cpsAction.LastModifiedAt = time.Now()

	if request.Decision == utils.DecisionApproved {
		// Use struct wrapper for current action, matching account validation
		type CurrentAction struct {
			HQ HQ `json:"hq"`
		}
		var currentAction CurrentAction
		currentActionBytes, err := toJSONBytes(cpsAction.CurrentAction)
		if err != nil {
			s.logger.Errorf("failed to convert CurrentAction to JSON: %v", err)
			return fmt.Errorf("FAILED_TO_CONVERT_CURRENT_ACTION")
		}
		if err := json.Unmarshal(currentActionBytes, &currentAction); err != nil {
			s.logger.Errorf("failed to unmarshal current action: %v, bytes: %s", err, string(currentActionBytes))
			return fmt.Errorf("FAILED_TO_UNMARSHAL_CURRENT_ACTION")
		}
		updatedHQ := currentAction.HQ
		if err := s.repository.UpdateHQ(ctx, updatedHQ.ID.Hex(), updatedHQ); err != nil {
			s.logger.Errorf("failed to update HQ: %v", err)
			return fmt.Errorf("FAILED_TO_UPDATE_HQ")
		}
		cpsAction.ActionStatus = action.ActionApproved
		s.logger.Infof("HQ archive time approved successfully", "action_id", request.ActionCode)
	} else if request.Decision == utils.DecisionDenied {
		cpsAction.ActionStatus = action.ActionRejected
		if request.RejectedReason != "" {
			cpsAction.RejectionReason = &request.RejectedReason
		}
		s.logger.Infof("HQ archive time update rejected", "action_id", request.ActionCode)
	} else {
		return fmt.Errorf("INVALID_DECISION")
	}

	if err := s.actionRepo.UpdateCpsAction(ctx, cpsAction); err != nil {
		s.logger.Errorf("failed to update CPS action: %v", err)
		return fmt.Errorf("FAILED_TO_UPDATE_CPS_ACTION")
	}

	s.logger.Infof("successfully processed HQ archive time update", "action_id", request.ActionCode)
	return nil
}

func stringToPointer(s string) *string {
	return &s
}
