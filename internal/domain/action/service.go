// Package action provides business logic and service layer implementations for CPS actions.
package action

import (
	"context"
	"encoding/json"
	"fmt"
	"time"

	common_util "github.com/CBE-Super-App/cbe-super-app-cps-action/pkgs/utils"

	"gitlab.com/bersufekadgetachew/cbe-super-app-shared/shared/utils"
)

const (
	ActionIDPrefix = "CPS_"
	ActionIDLength = 10
)

type ServiceInterface interface {
	GetServicePaginated(ctx context.Context, limit, offset int) ([]ServiceDetails, error)
	UpdateServiceFlagRequest(ctx context.Context, id string, action bool, maker User) (string, error)
	UpdateServiceFlag(ctx context.Context, actionID string, action bool, checker User) error
	GetAccountByAccount(ctx context.Context, Account string) ([]LinkedAccount, error)
	RemoveCifRequest(ctx context.Context, id []string, action bool, maker User) (string, error)
	RemoveCif(ctx context.Context, actionID string, action bool, checker User) error
	CreateCpsAction(ctx context.Context, Action CPSAction) (CPSAction, error)
}

type ServiceStore struct {
	Repository ActionRepository
	Logger     utils.Logger
}

func NewService(repo ActionRepository, logger utils.Logger) ServiceInterface {
	return &ServiceStore{
		Repository: repo,
		Logger:     logger,
	}
}

func (s *ServiceStore) createCpsAction(maker User, actionID string, currentAction any) CPSAction {
	return CPSAction{
		ActionCode:        actionID,
		MakerID:           maker.UserID,
		MakerName:         maker.FullName,
		MakerPhoneNumber:  maker.PhoneNumber,
		ActionType:        ActionCreate,
		ActionStatus:      ActionPending,
		CurrentAction:     currentAction,
		CreatedAt:         time.Now(),
		LastModifiedAt:    time.Now(),
		MakerActionTime:   time.Now(),
		CheckerActionTime: time.Time{},
	}

}

func (s *ServiceStore) approveOrRejectAction(cpsAction *CPSAction, checker User, approved bool) {
	if approved {
		cpsAction.ActionStatus = ActionApproved
	} else {
		cpsAction.ActionStatus = ActionRejected
	}
	cpsAction.CheckerID = checker.UserID
	cpsAction.CheckerName = checker.FullName
	cpsAction.CheckerPhoneNumber = checker.PhoneNumber
	cpsAction.CheckerActionTime = time.Now()
	cpsAction.LastModifiedAt = time.Now()
}

func extractCurrentAction(input any) (CurrentAction, error) {
	var current CurrentAction

	jsonBytes, err := json.Marshal(input)

	if err != nil {
		return current, fmt.Errorf("failed to marshal current action: %w", err)
	}

	err = json.Unmarshal(jsonBytes, &current)
	if err != nil {
		return current, fmt.Errorf("failed to unmarshal into CurrentAction: %w", err)
	}

	return current, nil
}

func (s *ServiceStore) buildAndSaveCpsAction(ctx context.Context, maker User, current CurrentAction) (string, error) {
	actionID := utils.Random(ActionIDLength, &utils.PreSufix{Prefix: ActionIDPrefix})
	cpsAction := s.createCpsAction(maker, actionID, current)
	data, err := s.Repository.CreateCpsAction(ctx, cpsAction)
	if err != nil {
		return "", err
	}
	return data.ActionCode, nil
}

func (s *ServiceStore) GetServicePaginated(ctx context.Context, limit, offset int) ([]ServiceDetails, error) {
	return s.Repository.GetAllHqServicesPaginated(ctx, offset, limit)
}

func (s *ServiceStore) UpdateServiceFlagRequest(ctx context.Context, id string, action bool, maker User) (string, error) {
	_, err := s.Repository.GetHqServiceById(ctx, id)
	if err != nil {
		s.Logger.Errorf("UpdateServiceFlagRequest: failed to get HQ service by ID", "id", id, "error", err)
		return "", err
	}

	lastAction, err := s.Repository.FetchLastCpsActionByMakerID(ctx, maker.UserID)
	if err == nil && lastAction.ActionStatus == ActionPending {
		s.Logger.Errorf("UpdateServiceFlagRequest: pending CPS action already exists", "makerID", maker.UserID)
		return "", fmt.Errorf(common_util.PendingCPSActionExists)
	}

	return s.buildAndSaveCpsAction(ctx, maker, CurrentAction{
		Id:     []string{id},
		Action: action,
	})
}

func (s *ServiceStore) UpdateServiceFlag(ctx context.Context, actionID string, action bool, checker User) error {
	cpsAction, err := s.Repository.FetchCpsActionById(ctx, actionID)
	if err != nil {
		s.Logger.Errorf("UpdateServiceFlag: failed to fetch CPS action by ID", "actionID", actionID, "error", err)
		return err
	}

	s.approveOrRejectAction(&cpsAction, checker, action)

	err = s.Repository.UpdateCpsAction(ctx, cpsAction)
	if err != nil {
		s.Logger.Errorf("UpdateServiceFlag: failed to update CPS action", "actionID", actionID, "error", err)
		return err
	}

	currentAction, err := extractCurrentAction(cpsAction.CurrentAction)
	if err != nil {
		s.Logger.Errorf("UpdateServiceFlag: invalid type for CurrentAction", "actionID", actionID, "error", err)
		return fmt.Errorf(common_util.UnhandledServerError)
	}
	service, err := s.Repository.GetHqServiceById(ctx, currentAction.Id[0])
	if err != nil {
		s.Logger.Errorf("UpdateServiceFlag: failed to fetch service by ID", "serviceID", currentAction.Id[0], "error", err)
		return err
	}

	service.Enabled = currentAction.Action

	if err := s.Repository.UpdateHqService(ctx, service); err != nil {
		s.Logger.Errorf("UpdateServiceFlag: failed to update service", "serviceID", service.ID, "error", err)
		return err
	}

	return nil
}

func (s *ServiceStore) GetAccountByAccount(ctx context.Context, account string) ([]LinkedAccount, error) {
	data, err := s.Repository.FetchAccountsByAccountNumber(ctx, account)
	if err != nil {
		s.Logger.Errorf("GetAccountByAccount: failed to fetch accounts", "account", account, "error", err)
		return nil, err
	}
	return data, nil
}

func (s *ServiceStore) RemoveCifRequest(ctx context.Context, ids []string, action bool, maker User) (string, error) {
	lastAction, err := s.Repository.FetchLastCpsActionByMakerID(ctx, maker.UserID)
	if err != nil {
		s.Logger.Errorf("RemoveCifRequest: failed to fetch last CPS action", "makerID", maker.UserID, "error", err)
		return "", err
	}
	if lastAction.ActionStatus == ActionPending {
		s.Logger.Errorf("RemoveCifRequest: pending CPS action already exists", "makerID", maker.UserID)
		return "", fmt.Errorf(common_util.PendingCPSActionExists)
	}

	return s.buildAndSaveCpsAction(ctx, maker, CurrentAction{
		Id:     ids,
		Action: action,
	})

}

func (s *ServiceStore) RemoveCif(ctx context.Context, actionID string, action bool, checker User) error {
	cpsAction, err := s.Repository.FetchCpsActionById(ctx, actionID)
	if err != nil {
		s.Logger.Errorf("RemoveCif: failed to fetch CPS action", "actionID", actionID, "error", err)
		return err
	}

	s.approveOrRejectAction(&cpsAction, checker, action)

	if err := s.Repository.UpdateCpsAction(ctx, cpsAction); err != nil {
		s.Logger.Errorf("RemoveCif: failed to update CPS action", "actionID", actionID, "error", err)
		return err
	}

	currentAction, err := extractCurrentAction(cpsAction.CurrentAction)
	if err != nil {
		s.Logger.Errorf("RemoveCif: failed to extract current action", "actionID", actionID, "error", err)
		return err
	}

	fmt.Println(currentAction.Id, "currentAction IDs")

	linkedAccounts, err := s.Repository.FetchLinkedAccountById(ctx, currentAction.Id)
	if err != nil {
		s.Logger.Errorf("RemoveCif: failed to fetch linked accounts", "accountIDs", currentAction.Id, "error", err)
		return err
	}

	for _, account := range linkedAccounts {
		account.LinkedStatus = false
		account.IsAccountActive = false

		if _, err := s.Repository.UpdateAccount(ctx, account); err != nil {
			s.Logger.Errorf("RemoveCif: failed to update account", "accountID", account.ID, "error", err)
			return fmt.Errorf(common_util.GeneralDBUpdateFailed)
		}
	}

	return nil
}

func (s *ServiceStore) CreateCpsAction(ctx context.Context, action CPSAction) (CPSAction, error) {
	lastAction, err := s.Repository.FetchLastCpsActionByMakerID(ctx, action.MakerID)
	if err == nil && lastAction.ActionStatus == ActionPending {
		s.Logger.Errorf("CreateCpsAction: pending CPS action already exists", "makerID", action.MakerID)
		return CPSAction{}, fmt.Errorf(common_util.PendingCPSActionExists)
	}
	actionID := utils.Random(ActionIDLength, &utils.PreSufix{Prefix: ActionIDPrefix})
	cpsAction := s.createCpsAction(
		User{
			UserID:      action.MakerID,
			FullName:    action.MakerName,
			PhoneNumber: action.MakerPhoneNumber,
		},
		actionID,
		action.CurrentAction,
	)
	cpsAction.RequestAction = action.RequestAction
	createdAction, err := s.Repository.CreateCpsAction(ctx, cpsAction)
	if err != nil {
		s.Logger.Errorf("CreateCpsAction: failed to create CPS action", "makerID", action.MakerID, "error", err)
		return CPSAction{}, err
	}
	return createdAction, nil
}
