package hq

import (
	"context"
	"encoding/json"
	"fmt"
	"reflect"
	"time"

	cps_constants "github.com/CBE-Super-App/cbe-super-app-cps-action/internal/domain/cps_actions/constant"
	// cps_entities "github.com/CBE-Super-App/cbe-super-app-cps-action/internal/domain/cps_actions/entities"
	entities "github.com/CBE-Super-App/cbe-super-app-cps-action/internal/domain/cps_actions/entities"
	cps_service "github.com/CBE-Super-App/cbe-super-app-cps-action/internal/domain/cps_actions/services"

	utils "github.com/CBE-Super-App/cbe-super-app-cps-action/pkgs/utils"
	constant "github.com/CBE-Super-App/cbe-super-app-cps-action/utils"

	sharedutils "gitlab.com/bersufekadgetachew/cbe-super-app-shared/shared/utils"
)

type Service interface {
	GetHQ(ctx context.Context, id string) (HQ, error)
	GetAllHQ(ctx context.Context, filerParams *constant.Filter) (*utils.PaginatedResponse[[]*HQ], error)
	GetBlockTime(ctx context.Context) (BlockTimeResponse, error)
	GetArchiveTime(ctx context.Context) (ArchiveTimeResponse, error)
	GetPasswordExpiry(ctx context.Context) (PasswordExpiryResponse, error)
	UpdateBlockTimeRequest(ctx context.Context, request UpdateBlockTimeRequest) (*entities.CPSAction, error)
	UpdateArchiveTimeRequest(ctx context.Context, request UpdateArchiveTimeRequest) (*entities.CPSAction, error)
	UpdatePasswordExpiryRequest(ctx context.Context, request UpdatePasswordExpiryRequest) (*entities.CPSAction, error)
	AuthorizeUpdateBlockTime(ctx context.Context, action *entities.CPSAction) (*entities.CPSAction, error)
	AuthorizeUpdateArchiveTime(ctx context.Context, action *entities.CPSAction) (*entities.CPSAction, error)
	AuthorizeUpdatePasswordExpiry(ctx context.Context, action *entities.CPSAction) (*entities.CPSAction, error)
	Authorize(ctx context.Context, action *entities.CPSAction) (*entities.CPSAction, error)
}

type ServiceStore struct {
	repository Repository
	cpsService cps_service.CPSActionService
	logger     sharedutils.Logger
}

func NewService(repo Repository, cpsService cps_service.CPSActionService, logger sharedutils.Logger) Service {
	return &ServiceStore{
		repository: repo,
		cpsService: cpsService,
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

func (s *ServiceStore) GetBlockTime(ctx context.Context) (BlockTimeResponse, error) {
	fmt.Println("=========Service=========")
	hq, err := s.repository.GetSingleHQ(ctx)
	if err != nil {
		s.logger.Errorf("failed to fetch HQ: %v", err)
		return BlockTimeResponse{}, fmt.Errorf("NOT_FOUND")
	}
	return BlockTimeResponse{
		BlockTime:      hq.BlockTime,
		CreatedAtBlock: hq.CreatedAtBlock,
		UpdatedAtBlock: hq.UpdatedAtBlock,
	}, nil
}

func (s *ServiceStore) GetArchiveTime(ctx context.Context) (ArchiveTimeResponse, error) {
	hq, err := s.repository.GetSingleHQ(ctx)
	if err != nil {
		s.logger.Errorf("failed to fetch HQ: %v", err)
		return ArchiveTimeResponse{}, fmt.Errorf("NOT_FOUND")
	}
	return ArchiveTimeResponse{
		ArchiveTime:      hq.ArchiveTime,
		CreatedAtArchive: hq.CreatedAtArchive,
		UpdatedAtArchive: hq.UpdatedAtArchive,
	}, nil
}

func (s *ServiceStore) GetPasswordExpiry(ctx context.Context) (PasswordExpiryResponse, error) {
	hq, err := s.repository.GetSingleHQ(ctx)
	if err != nil {
		s.logger.Errorf("failed to fetch HQ: %v", err)
		return PasswordExpiryResponse{}, fmt.Errorf("NOT_FOUND")
	}
	return PasswordExpiryResponse{
		PasswordExpiry:          hq.PasswordExpiry,
		CreatedAtPasswordExpiry: hq.CreatedAtPasswordExpiry,
		UpdatedAtPasswordExpiry: hq.UpdatedAtPasswordExpiry,
	}, nil
}

func (s *ServiceStore) UpdateBlockTimeRequest(ctx context.Context, request UpdateBlockTimeRequest) (*entities.CPSAction, error) {
	originalHQ, err := s.repository.GetSingleHQ(ctx)
	if err != nil {
		s.logger.Errorf("failed to fetch HQ: %v", err)
		return nil, fmt.Errorf("NOT_FOUND")
	}

	updatedHQ := originalHQ
	updatedHQ.BlockTime = request.BlockTime

	cpsAction := s.cpsService.BuildCPSAction(ctx, entities.CreateCPSRequest{
		User: entities.User{
			UserCode:    request.MakerID,
			FullName:    request.MakerName,
			PhoneNumber: request.MakerPhone,
			Department:  request.Department,
		},
		CurData:       updatedHQ,
		PrevData:      originalHQ,
		RequestAction: cps_constants.RequestUpdateBlockTime,
		ActionStatus:  cps_constants.ActionPending,
		ActionType:    cps_constants.ActionUpdate,
	})

	// Set the unique ID
	cpsAction.UniqueID = originalHQ.ID.Hex()

	createdAction, err := s.cpsService.CreateCPSAction(ctx, cpsAction)
	if err != nil {
		s.logger.Errorf("failed to create CPS action: %v", err)
		return nil, fmt.Errorf("FAILED_TO_CREATE_CPS_ACTION")
	}

	return createdAction, nil
}

func (s *ServiceStore) UpdateArchiveTimeRequest(ctx context.Context, request UpdateArchiveTimeRequest) (*entities.CPSAction, error) {
	originalHQ, err := s.repository.GetSingleHQ(ctx)
	if err != nil {
		s.logger.Errorf("failed to fetch HQ: %v", err)
		return nil, fmt.Errorf("NOT_FOUND")
	}

	updatedHQ := originalHQ
	updatedHQ.ArchiveTime = request.ArchiveTime

	// Create CPS action using the generic service
	cpsAction := s.cpsService.BuildCPSAction(ctx, entities.CreateCPSRequest{
		User: entities.User{
			UserCode:    request.MakerID,
			FullName:    request.MakerName,
			PhoneNumber: request.MakerPhone,
			Department:  request.Department,
		},
		CurData:       updatedHQ,
		PrevData:      originalHQ,
		RequestAction: cps_constants.RequestUpdateArchiveExpiry,
		ActionStatus:  cps_constants.ActionPending,
		ActionType:    cps_constants.ActionUpdate,
	})

	// Set the unique ID
	cpsAction.UniqueID = originalHQ.ID.Hex()

	createdAction, err := s.cpsService.CreateCPSAction(ctx, cpsAction)
	if err != nil {
		s.logger.Errorf("failed to create CPS action: %v", err)
		return nil, fmt.Errorf("FAILED_TO_CREATE_CPS_ACTION")
	}

	return createdAction, nil
}

func (s *ServiceStore) UpdatePasswordExpiryRequest(ctx context.Context, request UpdatePasswordExpiryRequest) (*entities.CPSAction, error) {
	originalHQ, err := s.repository.GetSingleHQ(ctx)
	if err != nil {
		s.logger.Errorf("failed to fetch HQ: %v", err)
		return nil, fmt.Errorf("NOT_FOUND")
	}

	updatedHQ := originalHQ
	updatedHQ.PasswordExpiry = request.PasswordExpiry

	// Create CPS action using the generic service
	cpsAction := s.cpsService.BuildCPSAction(ctx, entities.CreateCPSRequest{
		User: entities.User{
			UserCode:    request.MakerID,
			FullName:    request.MakerName,
			PhoneNumber: request.MakerPhone,
			Department:  request.Department,
		},
		CurData:       updatedHQ,
		PrevData:      originalHQ,
		RequestAction: cps_constants.RequestUpdatePasswordExpiry,
		ActionStatus:  cps_constants.ActionPending,
		ActionType:    cps_constants.ActionUpdate,
	})

	// Set the unique ID
	cpsAction.UniqueID = originalHQ.ID.Hex()

	createdAction, err := s.cpsService.CreateCPSAction(ctx, cpsAction)
	if err != nil {
		s.logger.Errorf("failed to create CPS action: %v", err)
		return nil, fmt.Errorf("FAILED_TO_CREATE_CPS_ACTION")
	}

	return createdAction, nil
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
	default:
		rv := reflect.ValueOf(v)
		if rv.Kind() == reflect.Struct {
			return json.Marshal(v)
		}
		return nil, fmt.Errorf("unsupported type: %T", v)
	}
}

func (s *ServiceStore) AuthorizeUpdateBlockTime(ctx context.Context, cpsAction *entities.CPSAction) (*entities.CPSAction, error) {
	return s.authorizeUpdateHQ(ctx, cpsAction, "block_time")
}

func (s *ServiceStore) AuthorizeUpdateArchiveTime(ctx context.Context, cpsAction *entities.CPSAction) (*entities.CPSAction, error) {
	return s.authorizeUpdateHQ(ctx, cpsAction, "archive_time")
}

func (s *ServiceStore) AuthorizeUpdatePasswordExpiry(ctx context.Context, cpsAction *entities.CPSAction) (*entities.CPSAction, error) {
	return s.authorizeUpdateHQ(ctx, cpsAction, "password_expiry")
}

func (s *ServiceStore) authorizeUpdateHQ(ctx context.Context, cpsAction *entities.CPSAction, actionType string) (*entities.CPSAction, error) {
	

	currentActionBytes, err := toJSONBytes(cpsAction.CurrentAction)
	if err != nil {
		s.logger.Errorf("failed to convert CurrentAction to JSON for %s: %v", actionType, err)
		return nil, fmt.Errorf("FAILED_TO_CONVERT_CURRENT_ACTION")
	}

	var updatedHQ HQ
	if err := json.Unmarshal(currentActionBytes, &updatedHQ); err != nil {
		s.logger.Errorf("failed to unmarshal CurrentAction for %s: %v, bytes: %s",
			actionType, err, string(currentActionBytes))
		return nil, fmt.Errorf("FAILED_TO_UNMARSHAL_CURRENT_ACTION")
	}

	// fmt.Println("updatedHQqqqqqq", updatedHQ)

	var field string
	var value interface{}
	now := time.Now()

	switch actionType {
	case "block_time":
		field = "block_time"
		value = updatedHQ.BlockTime
	case "archive_time":
		field = "archive_time"
		value = updatedHQ.ArchiveTime
	case "password_expiry":
		field = "password_expiry"
		value = updatedHQ.PasswordExpiry
	default:
		return nil, fmt.Errorf("UNSUPPORTED_FIELD")
	}

	if err := s.repository.UpdateHQField(ctx, field, value, now); err != nil {
		s.logger.Errorf("failed to update HQ for %s: %v", actionType, err)
		return nil, fmt.Errorf("FAILED_TO_UPDATE_HQ")
	}

	cpsAction.ActionStatus = cps_constants.ActionApproved
	s.logger.Infof("HQ %s approved successfully, action_id=%s", actionType, cpsAction.ActionCode)

	return cpsAction, nil
}


func (s *ServiceStore) Authorize(ctx context.Context, action *entities.CPSAction) (*entities.CPSAction, error) {
	switch action.RequestAction {
	case cps_constants.RequestUpdateBlockTime:
		return s.AuthorizeUpdateBlockTime(ctx, action)
	case cps_constants.RequestUpdateArchiveExpiry:
		return s.AuthorizeUpdateArchiveTime(ctx, action)
	case cps_constants.RequestUpdatePasswordExpiry:
		return s.AuthorizeUpdatePasswordExpiry(ctx, action)
	default:
		return nil, fmt.Errorf("UNSUPPORTED_REQUEST_ACTION")
	}
}
