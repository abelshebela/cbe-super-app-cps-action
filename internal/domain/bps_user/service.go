package bps_user

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"time"

	// "github.com/CBE-Super-App/cbe-super-app-cps-action/internal/domain/action"
	cps_constants "github.com/CBE-Super-App/cbe-super-app-cps-action/internal/domain/cps_actions/constant"
	cps_entities "github.com/CBE-Super-App/cbe-super-app-cps-action/internal/domain/cps_actions/entities"
	entities "github.com/CBE-Super-App/cbe-super-app-cps-action/internal/domain/cps_actions/entities"
	cps_repo "github.com/CBE-Super-App/cbe-super-app-cps-action/internal/domain/cps_actions/repository"
	"go.mongodb.org/mongo-driver/v2/mongo"

	utils "github.com/CBE-Super-App/cbe-super-app-cps-action/pkgs/utils"
	constant "github.com/CBE-Super-App/cbe-super-app-cps-action/utils"
	sharedutils "gitlab.com/bersufekadgetachew/cbe-super-app-shared/shared/utils"
)

type Service interface {
	// Basic CRUD operations
	GetBPSUserByUserCode(ctx context.Context, userCode string) (*BPSUser, error)
	GetAllBPSUsers(ctx context.Context, filerParams *constant.MongoFilter) (*utils.PaginatedResponse[[]*BPSUser], error)

	// Enable/Disable operations with maker-checker flow
	EnableBPSUserRequest(ctx context.Context, request EnableBPSUserRequest) (*cps_entities.CPSAction, error)
	DisableBPSUserRequest(ctx context.Context, request DisableBPSUserRequest) (*cps_entities.CPSAction, error)

	// Authorization methods
	AuthorizeEnableBPSUser(ctx context.Context, action *entities.CPSAction) (*cps_entities.CPSAction, error)
	AuthorizeDisableBPSUser(ctx context.Context, action *entities.CPSAction) (*cps_entities.CPSAction, error)
	Authorize(ctx context.Context, action *entities.CPSAction) (*entities.CPSAction, error)
}

type EnableBPSUserRequest struct {
	UserCode   string `json:"user_code" validate:"required"`
	MakerID    string `json:"maker_id"`
	MakerName  string `json:"maker_name"`
	MakerPhone string `json:"maker_phone"`
	Department string `json:"department"`
}

type DisableBPSUserRequest struct {
	UserCode   string `json:"user_code" validate:"required"`
	MakerID    string `json:"maker_id"`
	MakerName  string `json:"maker_name"`
	MakerPhone string `json:"maker_phone"`
	Department string `json:"department"`
}

type ServiceStore struct {
	repository Repository
	actionRepo cps_repo.CPSActionRepository
	logger     sharedutils.Logger
}

func NewBPSUserService(repo Repository, actionRepo cps_repo.CPSActionRepository, logger sharedutils.Logger) Service {
	return &ServiceStore{
		repository: repo,
		actionRepo: actionRepo,
		logger:     logger,
	}
}

func (s *ServiceStore) GetBPSUserByUserCode(ctx context.Context, userCode string) (*BPSUser, error) {
	if userCode == "" {
		return nil, fmt.Errorf("INVALID_USER_CODE")
	}

	user, err := s.repository.GetBPSUserByUserCode(ctx, userCode)
	if err != nil {
		s.logger.Errorf("failed to fetch BPS user: %v", err)
		return nil, fmt.Errorf("NOT_FOUND")
	}

	return user, nil
}

func (s *ServiceStore) GetAllBPSUsers(ctx context.Context, filterParams *constant.MongoFilter) (*utils.PaginatedResponse[[]*BPSUser], error) {
	users, err := s.repository.GetAllBPSUsers(ctx, filterParams)
	if err != nil {
		s.logger.Errorf("failed to fetch BPS users: %v", err)
		return nil, err
	}

	return users, nil
}

func (s *ServiceStore) hasPendingAction(ctx context.Context, requestAction string, department string) (bool, error) {
	// Check for pending actions using the repository
	pending, err := s.repository.CheckPendingAction(ctx, requestAction, department)
	if err != nil && !errors.Is(err, mongo.ErrNoDocuments) {
		return false, err
	}
	return pending, nil
}

func (s *ServiceStore) EnableBPSUserRequest(ctx context.Context, request EnableBPSUserRequest) (*cps_entities.CPSAction, error) {
	// Check if user exists
	originalUser, err := s.repository.GetBPSUserByUserCode(ctx, request.UserCode)
	if err != nil {
		s.logger.Errorf("failed to fetch BPS user: %v", err)
		return nil, fmt.Errorf("NOT_FOUND")
	}
	if originalUser.Enabled == true {
		return nil, fmt.Errorf("USER_ALREADY_ENABLED")
	}

	pending, err := s.hasPendingAction(ctx, string(cps_constants.RequestEnableBPSUser), request.Department)
	if err != nil && !errors.Is(err, mongo.ErrNoDocuments) {
		return nil, err
	}
	if pending {
		s.logger.Infof("pending action found for enable BPS user")
		return nil, fmt.Errorf("PENDING_ACTION_EXISTS")
	}

	// Create updated user data
	updatedUser := *originalUser
	updatedUser.Enabled = true

	type CurrentAction struct {
		BPSUser BPSUser `json:"bps_user"`
	}
	previousAction := CurrentAction{BPSUser: *originalUser}
	currentAction := CurrentAction{BPSUser: updatedUser}

	actionID := sharedutils.Random(10, &sharedutils.PreSufix{Prefix: "CPS_"})
	a := &cps_entities.CPSAction{
		ActionCode:       actionID,
		MakerID:          request.MakerID,
		MakerName:        request.MakerName,
		MakerPhoneNumber: request.MakerPhone,
		Department:       request.Department,
		ActionType:       cps_constants.ActionUpdate,
		RequestAction:    cps_constants.RequestAction(cps_constants.RequestEnableBPSUser),
		ActionStatus:     cps_constants.ActionPending,
		CurrentAction:    currentAction,
		PreviousAction:   previousAction,
		CreatedAt:        time.Now(),
		LastModifiedAt:   time.Now(),
		MakerActionTime:  time.Now(),
		UniqueID:         originalUser.ID.Hex(),
	}

	createdAction, err := s.actionRepo.CreateCPSAction(ctx, a)
	if err != nil {
		s.logger.Errorf("failed to create CPS action: %v", err)
		return nil, fmt.Errorf("FAILED_TO_CREATE_CPS_ACTION")
	}

	return createdAction, nil
}

func (s *ServiceStore) DisableBPSUserRequest(ctx context.Context, request DisableBPSUserRequest) (*cps_entities.CPSAction, error) {
	// Check if user exists
	originalUser, err := s.repository.GetBPSUserByUserCode(ctx, request.UserCode)
	if err != nil {

		s.logger.Errorf("failed to fetch BPS user: %v", err)
		return nil, fmt.Errorf("NOT_FOUND")
	}
	if originalUser.Enabled == false {
		return nil, fmt.Errorf("USER_ALREADY_DISABLED")
	}

	// Check for pending actions
	pending, err := s.hasPendingAction(ctx, string(cps_constants.RequestDisableBPSUser), request.Department)
	if err != nil && !errors.Is(err, mongo.ErrNoDocuments) {
		return nil, err
	}
	if pending {
		s.logger.Infof("pending action found for disable BPS user")
		return nil, fmt.Errorf("PENDING_ACTION_EXISTS")
	}

	// Create updated user data
	updatedUser := *originalUser
	updatedUser.Enabled = false

	type CurrentAction struct {
		BPSUser BPSUser `json:"bps_user"`
	}
	previousAction := CurrentAction{BPSUser: *originalUser}
	currentAction := CurrentAction{BPSUser: updatedUser}

	actionID := sharedutils.Random(10, &sharedutils.PreSufix{Prefix: "CPS_"})
	a := &cps_entities.CPSAction{
		ActionCode:       actionID,
		MakerID:          request.MakerID,
		MakerName:        request.MakerName,
		MakerPhoneNumber: request.MakerPhone,
		Department:       request.Department,
		ActionType:       cps_constants.ActionUpdate,
		RequestAction:    cps_constants.RequestAction(cps_constants.RequestDisableBPSUser),
		ActionStatus:     cps_constants.ActionPending,
		CurrentAction:    currentAction,
		PreviousAction:   previousAction,
		CreatedAt:        time.Now(),
		LastModifiedAt:   time.Now(),
		MakerActionTime:  time.Now(),
		UniqueID:         originalUser.ID.Hex(),
	}

	createdAction, err := s.actionRepo.CreateCPSAction(ctx, a)
	if err != nil {
		s.logger.Errorf("failed to create CPS action: %v", err)
		return nil, fmt.Errorf("FAILED_TO_CREATE_CPS_ACTION")
	}

	return createdAction, nil
}

func (s *ServiceStore) authorizeUpdateBPSUser(ctx context.Context, cpsAction *entities.CPSAction, actionType string) (*entities.CPSAction, error) {
	type CurrentAction struct {
		BPSUser BPSUser `json:"bps_user"`
	}

	var currentAction CurrentAction
	currentActionBytes, err := toJSONBytes(cpsAction.CurrentAction)
	if err != nil {
		s.logger.Errorf("failed to convert CurrentAction to JSON for %s: %v", actionType, err)
		return nil, fmt.Errorf("FAILED_TO_CONVERT_CURRENT_ACTION")
	}

	if err := json.Unmarshal(currentActionBytes, &currentAction); err != nil {
		s.logger.Errorf("failed to unmarshal CurrentAction for %s: %v, bytes: %s", actionType, err, string(currentActionBytes))
		return nil, fmt.Errorf("FAILED_TO_UNMARSHAL_CURRENT_ACTION")
	}

	updatedUser := currentAction.BPSUser
	now := time.Now()

	// Update the user in database
	if err := s.repository.UpdateBPSUserStatus(ctx, updatedUser.UserCode, updatedUser.Enabled, now); err != nil {
		s.logger.Errorf("failed to update BPS user for %s: %v", actionType, err)
		return nil, fmt.Errorf("FAILED_TO_UPDATE_USER")
	}

	cpsAction.ActionStatus = cps_constants.ActionApproved
	s.logger.Infof("BPS user %s approved successfully", actionType)
	return cpsAction, nil
}

func (s *ServiceStore) AuthorizeEnableBPSUser(ctx context.Context, action *entities.CPSAction) (*entities.CPSAction, error) {
	return s.authorizeUpdateBPSUser(ctx, action, "enable")
}

func (s *ServiceStore) AuthorizeDisableBPSUser(ctx context.Context, action *entities.CPSAction) (*entities.CPSAction, error) {
	return s.authorizeUpdateBPSUser(ctx, action, "disable")
}

func (s *ServiceStore) Authorize(ctx context.Context, action *entities.CPSAction) (*entities.CPSAction, error) {
	switch action.RequestAction {
	case cps_constants.RequestEnableBPSUser:
		return s.AuthorizeEnableBPSUser(ctx, action)
	case cps_constants.RequestDisableBPSUser:
		return s.AuthorizeDisableBPSUser(ctx, action)
	default:
		return nil, fmt.Errorf("UNSUPPORTED_REQUEST_ACTION")
	}
}

func toJSONBytes(data interface{}) ([]byte, error) {
	return json.Marshal(data)
}
