package action

import (
	"context"
	"fmt"
	"time"

	"gitlab.com/bersufekadgetachew/cbe-super-app-shared/shared/utils"
)

func (s *ServiceStore) GetServicePaginated(ctx context.Context, limit, offset int) ([]ServiceDetails, error) {
	return s.Repository.GetAllHqServicesPaginated(ctx, offset, limit)
}

func (s *ServiceStore) UpdateServiceFlagRequest(ctx context.Context, id string, action bool, makerId string) (string, error) {

	_, err := s.Repository.GetHqServiceById(ctx, id)
	if err != nil {
		return "", err
	}
	// service.Enabled = action
	actionID := utils.Random(10, &utils.PreSufix{Prefix: "CPS_"})
	a := CPSAction{
		ActionCode:         actionID,
		MakerID:            makerId,
		MakerName:          "",
		MakerPhoneNumber:   "",
		CheckerID:          "",
		CheckerName:        "",
		CheckerPhoneNumber: "",
		ActionType:         ActionCreate,
		RequestAction:      RequestUpdateServiceRule,
		ActionStatus:       ActionPending,
		CurrentAction: CurrentAction{
			Id:     []string{id},
			Action: action,
		},
		CreatedAt:         time.Now(),
		LastModifiedAt:    time.Now(),
		MakerActionTime:   time.Now(),
		CheckerActionTime: time.Time{},
	}
	_, err = s.Repository.CreateCpsAction(ctx, a)
	if err != nil {
		return "", err
	}
	return actionID, nil
}

func (s *ServiceStore) UpdateServiceFlag(ctx context.Context, action_id string, action bool, checker_id string) error {
	cps_action, err := s.Repository.FetchCpsActionById(ctx, action_id)
	if err != nil {
		return err
	}
	if action {
		cps_action.ActionStatus = ActionApproved
	} else {
		cps_action.ActionStatus = ActionRejected
	}
	cps_action.CheckerID = checker_id
	cps_action.CheckerName = ""
	cps_action.CheckerPhoneNumber = ""
	cps_action.CheckerActionTime = time.Now()
	cps_action.LastModifiedAt = time.Now()
	err = s.Repository.UpdateCpsAction(ctx, cps_action)
	if err != nil {
		return err
	}
	currentAction, ok := cps_action.CurrentAction.(CurrentAction)
	if !ok {
		return fmt.Errorf("failed to cast CurrentAction to its expected type")
	}
	service, err := s.Repository.GetHqServiceById(ctx, currentAction.Id[0])
	if err != nil {
		return err
	}
	service.Enabled = currentAction.Action
	return s.Repository.UpdateHqService(ctx, service)
}

func (s *ServiceStore) GetAccountByAccount(ctx context.Context, account string) ([]LinkedAccount, error) {
	data, err := s.Repository.FetchAccountsByAccountNumber(ctx, account)
	if err != nil {
		return nil, err
	}
	return data, nil
}

func (s *ServiceStore) RemoveCifRequest(ctx context.Context, id []string, action bool, maker_id string) (string, error) {
	last_action, err := s.Repository.FetchLastCpsActionByMakerID(ctx, maker_id)
	if err != nil {
		return "", err
	}
	if last_action.ActionStatus == ActionPending {
		return "", fmt.Errorf("You have a pending action, please wait for it to be processed")
	}
	actionId := utils.Random(10, &utils.PreSufix{Prefix: "CPS_"})
	a := CPSAction{
		ActionCode:         actionId,
		MakerID:            maker_id,
		MakerName:          "",
		MakerPhoneNumber:   "",
		CheckerID:          "",
		CheckerName:        "",
		CheckerPhoneNumber: "",
		ActionType:         ActionCreate,
		RequestAction:      RequestUpdateServiceRule,
		ActionStatus:       ActionPending,
		CurrentAction: CurrentAction{
			Id:     id,
			Action: action,
		},
		CreatedAt:         time.Now(),
		LastModifiedAt:    time.Now(),
		MakerActionTime:   time.Now(),
		CheckerActionTime: time.Time{},
	}
	result, err := s.Repository.CreateCpsAction(ctx, a)
	if err != nil {
		return "", err
	}
	return result.ID, nil
}

func (s *ServiceStore) RemoveCif(ctx context.Context, action_id string, action bool, checker_id string) error {
	cps_action, err := s.Repository.FetchCpsActionById(ctx, action_id)
	if err != nil {
		return err
	}
	if action {
		cps_action.ActionStatus = ActionApproved
	} else {
		cps_action.ActionStatus = ActionRejected
	}
	cps_action.CheckerID = checker_id
	cps_action.CheckerName = ""
	cps_action.CheckerPhoneNumber = ""
	cps_action.CheckerActionTime = time.Now()
	cps_action.LastModifiedAt = time.Now()
	err = s.Repository.UpdateCpsAction(ctx, cps_action)
	if err != nil {
		return err
	}
	currentAction, ok := cps_action.CurrentAction.(CurrentAction)
	if !ok {
		return fmt.Errorf("failed to cast CurrentAction to its expected type")
	}
	linked_accounts, err := s.Repository.FetchLinkedAccountById(ctx, currentAction.Id)
	if err != nil {
		return err
	}
	for _, account := range linked_accounts {
		account.LinkedStatus = false
		account.IsAccountActive = false
		_, err = s.Repository.UpdateAccount(ctx, account)
		if err != nil {
			return fmt.Errorf("failed to update account %s: %w", account.AccountNumber, err)
		}
	}
	return nil
}

func (s *ServiceStore) CreateCpsAction(ctx context.Context, action CPSAction) (CPSAction, error) {
	createdAction, err := s.Repository.CreateCpsAction(ctx, action)
	if err != nil {
		return CPSAction{}, err
	}
	return createdAction, nil
}
