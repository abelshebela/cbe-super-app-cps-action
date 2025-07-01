package hq

import (
	"context"
	"encoding/json"
	"errors"
	"time"

	"cbe-super-app-cps-action/internal/domain/action"
	// "cbe-super-app-cps-action/internal/domain/hq/"
	"gitlab.com/bersufekadgetachew/cbe-super-app-shared/shared/utils"
)

type Service interface {
	GetHQ(ctx context.Context, id string) (HQ, error)
	UpdateBlockTimeRequest(ctx context.Context, request UpdateBlockTimeRequest) (string, error)
	UpdateArchiveTimeRequest(ctx context.Context, request UpdateArchiveTimeRequest) (string, error)
	UpdateBlockTime(ctx context.Context, request ApproveRejectRequest) error
	UpdateArchiveTime(ctx context.Context, request ApproveRejectRequest) error
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

func (s *ServiceStore) GetHQ(ctx context.Context, id string) (HQ, error) {
	if id == "" {
		s.logger.Errorf("HQ ID is empty")
		return HQ{}, errors.New("HQ ID cannot be empty")
	}
	s.logger.Infof("fetching HQ", "id", id)

	hq, err := s.repository.GetHQByID(ctx, id)
	if err != nil {
		s.logger.Errorf("failed to fetch HQ: %v", err)
		return HQ{}, err
	}
	return hq, nil
}

func (s *ServiceStore) UpdateBlockTimeRequest(ctx context.Context, request UpdateBlockTimeRequest) (string, error) {
	originalHQ, err := s.repository.GetHQByID(ctx, request.MakerID)
	if err != nil {
		s.logger.Errorf("failed to fetch HQ: %v", err)
		return "", err
	}

	previousActionJSON, err := json.Marshal(originalHQ)
	if err != nil {
		s.logger.Errorf("failed to marshal previous action: %v", err)
		return "", errors.New("failed to marshal previous action")
	}

	updatedHQ := originalHQ
	updatedHQ.BlockTime = request.BlockTime
	currentActionJSON, err := json.Marshal(updatedHQ)
	if err != nil {
		s.logger.Errorf("failed to marshal current action: %v", err)
		return "", errors.New("failed to marshal current action")
	}

	s.logger.Infof("original HQ: %+v", originalHQ)
	s.logger.Infof("updated HQ: %+v", updatedHQ)

	actionID := utils.Random(10, &utils.PreSufix{Prefix: "CPS_"})

	a := action.CPSAction{
		ActionCode: actionID,
		Maker: action.User{
			UserID:      request.MakerID,
			FullName:    request.MakerName,
			PhoneNumber: request.MakerPhone,
			Timestamp:   time.Now(),
		},
		Checker:         action.User{},
		Department:      originalHQ.Name,
		ActionType:      action.ActionUpdate,
		RequestAction:   action.RequestUpdateHQBlockTime,
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

func (s *ServiceStore) UpdateArchiveTimeRequest(ctx context.Context, request UpdateArchiveTimeRequest) (string, error) {
	originalHQ, err := s.repository.GetHQByID(ctx, request.MakerID)
	if err != nil {
		s.logger.Errorf("failed to fetch HQ: %v", err)
		return "", err
	}

	previousActionJSON, err := json.Marshal(originalHQ)
	if err != nil {
		s.logger.Errorf("failed to marshal previous action: %v", err)
		return "", errors.New("failed to marshal previous action")
	}

	updatedHQ := originalHQ
	updatedHQ.ArchiveTime = request.ArchiveTime
	currentActionJSON, err := json.Marshal(updatedHQ)
	if err != nil {
		s.logger.Errorf("failed to marshal current action: %v", err)
		return "", errors.New("failed to marshal current action")
	}

	s.logger.Infof("original HQ: %+v", originalHQ)
	s.logger.Infof("updated HQ: %+v", updatedHQ)

	actionID := utils.Random(10, &utils.PreSufix{Prefix: "CPS_"})

	a := action.CPSAction{
		ActionCode: actionID,
		Maker: action.User{
			UserID:      request.MakerID,
			FullName:    request.MakerName,
			PhoneNumber: request.MakerPhone,
			Timestamp:   time.Now(),
		},
		Checker:         action.User{},
		Department:      originalHQ.ID,
		ActionType:      action.ActionUpdate,
		RequestAction:   action.RequestUpdateHQArchiveTime,
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

func (s *ServiceStore) UpdateBlockTime(ctx context.Context, request ApproveRejectRequest) error {
	if request.ActionCode == "" {
		s.logger.Errorf("action ID is empty")
		return errors.New("action ID cannot be empty")
	}

	if request.CheckerID == "" {
		s.logger.Errorf("checker ID is empty")
		return errors.New("checker ID cannot be empty")
	}

	s.logger.Infof("processing HQ block time update", "action_id", request.ActionCode, "approve", request.Approved, "checker_id", request.CheckerID)

	cpsAction, err := s.actionRepo.FetchCpsActionById(ctx, request.ActionCode)
	if err != nil {
		s.logger.Errorf("failed to fetch CPS action: %v", err)
		return errors.New("action not found")
	}

	if cpsAction.ActionStatus != action.ActionPending {
		s.logger.Errorf("action is not pending", "action_id", request.ActionCode, "status", cpsAction.ActionStatus)
		return errors.New("action is not pending")
	}

	cpsAction.Checker = action.User{
		UserID:      request.CheckerID,
		FullName:    request.CheckerName,
		PhoneNumber: request.CheckerPhone,
		Timestamp:   time.Now(),
	}
	cpsAction.LastModifiedAt = time.Now()

	if request.Approved {
		var updatedHQ HQ
		var currentActionBytes []byte
		switch v := cpsAction.CurrentAction.(type) {
		case json.RawMessage:
			currentActionBytes = v
		case []byte:
			currentActionBytes = v
		case string:
			currentActionBytes = []byte(v)
		case nil:
			s.logger.Errorf("current action is nil", "action_id", request.ActionCode)
			return errors.New("current action is nil")
		default:
			var err error
			currentActionBytes, err = json.Marshal(v)
			if err != nil {
				s.logger.Errorf("failed to marshal current action: %v", err)
				return errors.New("current action is not a valid type")
			}
		}

		if err := json.Unmarshal(currentActionBytes, &updatedHQ); err != nil {
			s.logger.Errorf("failed to unmarshal current action: %v", err)
			return errors.New("failed to unmarshal current action")
		}

		if err := s.repository.UpdateHQ(ctx, updatedHQ.ID, updatedHQ); err != nil {
			s.logger.Errorf("failed to update HQ: %v", err)
			return err
		}

		cpsAction.ActionStatus = action.ActionApproved
	} else {
		cpsAction.ActionStatus = action.ActionRejected
	}

	if err := s.actionRepo.UpdateCpsAction(ctx, cpsAction); err != nil {
		s.logger.Errorf("failed to update CPS action: %v", err)
		return err
	}

	s.logger.Infof("successfully processed HQ block time update", "action_id", request.ActionCode)
	return nil
}

func (s *ServiceStore) UpdateArchiveTime(ctx context.Context, request ApproveRejectRequest) error {
	if request.ActionCode == "" {
		s.logger.Errorf("action ID is empty")
		return errors.New("action ID cannot be empty")
	}

	if request.CheckerID == "" {
		s.logger.Errorf("checker ID is empty")
		return errors.New("checker ID cannot be empty")
	}

	s.logger.Infof("processing HQ archive time update", "action_id", request.ActionCode, "approve", request.Approved, "checker_id", request.CheckerID)

	cpsAction, err := s.actionRepo.FetchCpsActionById(ctx, request.ActionCode)
	if err != nil {
		s.logger.Errorf("failed to fetch CPS action: %v", err)
		return errors.New("action not found")
	}

	if cpsAction.ActionStatus != action.ActionPending {
		s.logger.Errorf("action is not pending", "action_id", request.ActionCode, "status", cpsAction.ActionStatus)
		return errors.New("action is not pending")
	}

	cpsAction.Checker = action.User{
		UserID:      request.CheckerID,
		FullName:    request.CheckerName,
		PhoneNumber: request.CheckerPhone,
		Timestamp:   time.Now(),
	}
	cpsAction.LastModifiedAt = time.Now()

	if request.Approved {
		var updatedHQ HQ
		var currentActionBytes []byte
		switch v := cpsAction.CurrentAction.(type) {
		case json.RawMessage:
			currentActionBytes = v
		case []byte:
			currentActionBytes = v
		case string:
			currentActionBytes = []byte(v)
		case nil:
			s.logger.Errorf("current action is nil", "action_id", request.ActionCode)
			return errors.New("current action is nil")
		default:
			var err error
			currentActionBytes, err = json.Marshal(v)
			if err != nil {
				s.logger.Errorf("failed to marshal current action: %v", err)
				return errors.New("current action is not a valid type")
			}
		}

		if err := json.Unmarshal(currentActionBytes, &updatedHQ); err != nil {
			s.logger.Errorf("failed to unmarshal current action: %v", err)
			return errors.New("failed to unmarshal current action")
		}

		if err := s.repository.UpdateHQ(ctx, updatedHQ.ID, updatedHQ); err != nil {
			s.logger.Errorf("failed to update HQ: %v", err)
			return err
		}

		cpsAction.ActionStatus = action.ActionApproved
	} else {
		cpsAction.ActionStatus = action.ActionRejected
	}

	if err := s.actionRepo.UpdateCpsAction(ctx, cpsAction); err != nil {
		s.logger.Errorf("failed to update CPS action: %v", err)
		return err
	}

	s.logger.Infof("successfully processed HQ archive time update", "action_id", request.ActionCode)
	return nil
}
