package hq

import (
	"context"
	"encoding/json"
	"fmt"
	"reflect"
	"time"

	"github.com/CBE-Super-App/cbe-super-app-cps-action/internal/domain/action"
	cps_constants "github.com/CBE-Super-App/cbe-super-app-cps-action/internal/domain/cps_actions/constant"
	entities "github.com/CBE-Super-App/cbe-super-app-cps-action/internal/domain/cps_actions/entities"

	constant "github.com/CBE-Super-App/cbe-super-app-cps-action/utils"

	sharedutils "gitlab.com/bersufekadgetachew/cbe-super-app-shared/shared/utils"
)

type Service interface {
	GetHQ(ctx context.Context, id string) (HQ, error)
	GetAllHQ(ctx context.Context, filerParams *constant.Filter) (*HQRespose, error)
	UpdateBlockTimeRequest(ctx context.Context, request UpdateBlockTimeRequest) (*action.CPSAction, error)
	UpdateArchiveTimeRequest(ctx context.Context, request UpdateArchiveTimeRequest) (*action.CPSAction, error)
	AuthorizeUpdateBlockTime(ctx context.Context, action *entities.CPSAction) (*entities.CPSAction, error)
	AuthorizeUpdateArchiveTime(ctx context.Context, action *entities.CPSAction) (*entities.CPSAction, error)
	Authorize(ctx context.Context, action *entities.CPSAction) (*entities.CPSAction, error)
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

func (s *ServiceStore) UpdateBlockTimeRequest(ctx context.Context, request UpdateBlockTimeRequest) (*action.CPSAction, error) {
	if request.ID == "" {
		return nil, fmt.Errorf("INVALID_ID")
	}
	originalHQ, err := s.repository.GetHQByID(ctx, request.ID)
	if err != nil {
		s.logger.Errorf("failed to fetch HQ: %v", err)
		return nil, fmt.Errorf("NOT_FOUND")
	}

	// Use actionRepo to check for pending actions by unique ID (outbound/init)
	pendingActions, err := s.actionRepo.(interface {
		FetchPendingActionsByUniqueID(context.Context, string) ([]action.ActionResponse, error)
	}).FetchPendingActionsByUniqueID(ctx, request.ID)
	if err != nil {
		return nil, fmt.Errorf("FAILED_TO_FETCH_PENDING_ACTIONS")
	}
	if len(pendingActions) > 0 {
		return nil, fmt.Errorf("PENDING_ACTION_EXISTS")
	}

	previousActionJSON, err := json.Marshal(originalHQ)
	if err != nil {
		s.logger.Errorf("failed to marshal previous action: %v", err)
		return nil, fmt.Errorf("FAILED_TO_MARSHAL_PREVIOUS_ACTION")
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
		ActionCode:       actionID,
		MakerID:          request.MakerID,
		MakerName:        request.MakerName,
		MakerPhoneNumber: request.MakerPhone,
		Department:       request.Department,
		ActionType:       action.ActionUpdate,
		RequestAction:    action.RequestUpdateHQBlockTime,
		ActionStatus:     action.ActionPending,
		CurrentAction:    currentAction,
		PreviousAction:   previousActionJSON,
		CreatedAt:        time.Now(),
		LastModifiedAt:   time.Now(),
		MakerActionTime:  time.Now(),
		UniqueId:         request.ID,
	}
	createdAction, err := s.actionRepo.CreateCpsAction(ctx, a)
	if err != nil {
		s.logger.Errorf("failed to create CPS action: %v", err)
		return nil, fmt.Errorf("FAILED_TO_CREATE_CPS_ACTION")
	}

	return &createdAction, nil
}

func (s *ServiceStore) UpdateArchiveTimeRequest(ctx context.Context, request UpdateArchiveTimeRequest) (*action.CPSAction, error) {
	if request.ID == "" {
		return nil, fmt.Errorf("INVALID_ID")
	}
	originalHQ, err := s.repository.GetHQByID(ctx, request.ID)
	if err != nil {
		s.logger.Errorf("failed to fetch HQ: %v", err)
		return nil, fmt.Errorf("NOT_FOUND")
	}

	// Use actionRepo to check for pending actions by unique ID (outbound/init)
	pendingActions, err := s.actionRepo.(interface {
		FetchPendingActionsByUniqueID(context.Context, string) ([]action.ActionResponse, error)
	}).FetchPendingActionsByUniqueID(ctx, request.ID)
	if err != nil {
		return nil, fmt.Errorf("FAILED_TO_FETCH_PENDING_ACTIONS")
	}
	if len(pendingActions) > 0 {
		return nil, fmt.Errorf("PENDING_ACTION_EXISTS")
	}

	previousActionJSON, err := json.Marshal(originalHQ)
	if err != nil {
		s.logger.Errorf("failed to marshal previous action: %v", err)
		return nil, fmt.Errorf("FAILED_TO_MARSHAL_PREVIOUS_ACTION")
	}

	updatedHQ := originalHQ
	updatedHQ.ArchiveTime = request.ArchiveTime
	// Use struct wrapper for current action, matching account validation
	type CurrentAction struct {
		HQ HQ `json:"hq"`
	}
	currentAction := CurrentAction{HQ: updatedHQ}

	actionID := sharedutils.Random(10, &sharedutils.PreSufix{Prefix: "CPS_"})
	a := action.CPSAction{
		ActionCode:        actionID,
		MakerID:           request.MakerID,
		MakerName:         request.MakerName,
		MakerPhoneNumber:  request.MakerPhone,
		Department:        request.Department,
		ActionType:        action.ActionUpdate,
		RequestAction:     action.RequestUpdateHQArchiveTime,
		ActionStatus:      action.ActionPending,
		CurrentAction:     currentAction,
		PreviousAction:    previousActionJSON,
		CreatedAt:         time.Now(),
		LastModifiedAt:    time.Now(),
		MakerActionTime:   time.Now(),
		UniqueId:          request.ID,
	}
	createdAction, err := s.actionRepo.CreateCpsAction(ctx, a)
	if err != nil {
		s.logger.Errorf("failed to create CPS action: %v", err)
		return nil, fmt.Errorf("FAILED_TO_CREATE_CPS_ACTION")
	}

	return &createdAction, nil
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

func (s *ServiceStore) AuthorizeUpdateBlockTime(ctx context.Context, cpsAction *entities.CPSAction) (*entities.CPSAction, error) {
	// Define wrapper to decode CurrentAction
	type CurrentAction struct {
		HQ HQ `json:"hq"`
	}

	var currentAction CurrentAction
	currentActionBytes, err := toJSONBytes(cpsAction.CurrentAction)
	if err != nil {
		s.logger.Errorf("failed to convert CurrentAction to JSON: %v", err)
		return nil, fmt.Errorf("FAILED_TO_CONVERT_CURRENT_ACTION")
	}

	if err := json.Unmarshal(currentActionBytes, &currentAction); err != nil {
		s.logger.Errorf("failed to unmarshal current action: %v, bytes: %s", err, string(currentActionBytes))
		return nil, fmt.Errorf("FAILED_TO_UNMARSHAL_CURRENT_ACTION")
	}

	updatedHQ := currentAction.HQ
	fmt.Println("updated hq", updatedHQ.ID)

	if err := s.repository.UpdateHQ(ctx, updatedHQ.ID.Hex(), updatedHQ); err != nil {
		s.logger.Errorf("failed to update HQ: %v", err)
		return nil, fmt.Errorf("FAILED_TO_UPDATE_HQ")
	}

	cpsAction.ActionStatus = cps_constants.ActionApproved
	return cpsAction, nil
}

func (s *ServiceStore) AuthorizeUpdateArchiveTime(ctx context.Context, cpsAction *entities.CPSAction) (*entities.CPSAction, error) {
	type CurrentAction struct {
		HQ HQ `json:"hq"`
	}

	var currentAction CurrentAction
	currentActionBytes, err := toJSONBytes(cpsAction.CurrentAction)
	if err != nil {
		s.logger.Errorf("failed to convert CurrentAction to JSON: %v", err)
		return nil, fmt.Errorf("FAILED_TO_CONVERT_CURRENT_ACTION")
	}

	if err := json.Unmarshal(currentActionBytes, &currentAction); err != nil {
		s.logger.Errorf("failed to unmarshal current action: %v, bytes: %s", err, string(currentActionBytes))
		return nil, fmt.Errorf("FAILED_TO_UNMARSHAL_CURRENT_ACTION")
	}

	updatedHQ := currentAction.HQ
	if err := s.repository.UpdateHQ(ctx, updatedHQ.ID.Hex(), updatedHQ); err != nil {
		s.logger.Errorf("failed to update HQ: %v", err)
		return nil, fmt.Errorf("FAILED_TO_UPDATE_HQ")
	}

	cpsAction.ActionStatus = cps_constants.ActionApproved
	s.logger.Infof("HQ archive time approved successfully", "action_id", cpsAction.ActionCode)
	return cpsAction, nil
}

func (s *ServiceStore) Authorize(ctx context.Context, action *entities.CPSAction) (*entities.CPSAction, error) {
	switch action.RequestAction {
	case cps_constants.RequestUpdateHQBlockTime:
		return s.AuthorizeUpdateBlockTime(ctx, action)
	case cps_constants.RequestUpdateHQArchiveTime:
		return s.AuthorizeUpdateArchiveTime(ctx, action)

	default:
		return nil, fmt.Errorf("unsupported request action: %s", action)
	}
}
