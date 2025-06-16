package action

import (
	"context"
<<<<<<< HEAD
	"fmt"
=======
>>>>>>> b69ee66 (feature added)

	"gitlab.com/bersufekadgetachew/cbe-super-app-shared/shared/utils"
)

func (s *ServiceStore) GetServicePaginated(ctx context.Context, limit, offset int) ([]Service, error) {
	return s.repository.GetAllHqServicesPaginated(ctx, offset, limit)
}
func (s *ServiceStore) UpdateServiceFlagRequest(ctx context.Context, id string, action bool, makerId string) (string, error) {
<<<<<<< HEAD
	last_action, err := s.repository.FetchLastCpsActionByMakerID(ctx, makerId)
	if last_action.ActionStatus == ActionPending {
		return "", fmt.Errorf("You have a pending action, please wait for it to be processed")
	}
=======
>>>>>>> b69ee66 (feature added)
	service, err := s.repository.GetHqServiceById(ctx, id)
	if err != nil {
		return "", err
	}
	service.Flag = action
	actionId := utils.Random(10, &utils.PreSufix{Prefix: "CPS_"})

	a := CPSAction{
		ActionCode: actionId,
		MakerAndChecker: MakerAndChecker{
			Maker: User{
				UserID: makerId,
			},
		},
		ActionType:    ActionCreate,
		RequestAction: RequestUpdateServiceRule,
		ActionStatus:  ActionPending,
		CurrentAction: CurrentAction{
<<<<<<< HEAD
			Id:     []string{id},
=======
			Id:     id,
>>>>>>> b69ee66 (feature added)
			Action: action,
		},
	}
	_, err = s.repository.CreateCpsAction(ctx, a)
	if err != nil {
		return "", err
	}
	return actionId, nil
}

func (s *ServiceStore) UpdateServiceFlag(ctx context.Context, action_id string, action bool, checker_id string) error {
	cps_action, err := s.repository.FetchCpsActionById(ctx, action_id)
	if err != nil {
		return err
	}
	if action {
		cps_action.ActionStatus = ActionApproved
	} else {
		cps_action.ActionStatus = ActionRejected
	}
	cps_action.MakerAndChecker.Checker.UserID = checker_id
	err = s.repository.UpdateCpsAction(ctx, cps_action)
	if err != nil {
		return err
	}
<<<<<<< HEAD
	currentAction, ok := cps_action.CurrentAction.(CurrentAction)
	if !ok {
		return fmt.Errorf("failed to cast CurrentAction to its expected type")
	}
	service, err := s.repository.GetHqServiceById(ctx, currentAction.Id[0])
	if err != nil {
		return err
	}
	service.Flag = currentAction.Action
	return s.repository.UpdateHqService(ctx, service)
}

func (s *ServiceStore) GetAccountByAccount(ctx context.Context, account string) ([]LinkedAccount, error) {
	data, err := s.repository.FetchAccountsByAccountNumber(ctx, account)
	if err != nil {
		return nil, err
	}
	return data, nil
}
func (s *ServiceStore) RemoveCifRequest(ctx context.Context, id []string, action bool, maker_id string) (string, error) {
	//create action
	last_action, err := s.repository.FetchLastCpsActionByMakerID(ctx, maker_id)
	if err != nil {
		return "", err
	}
	if last_action.ActionStatus == ActionPending {
		return "", fmt.Errorf("You have a pending action, please wait for it to be processed")
	}
	actionId := utils.Random(10, &utils.PreSufix{Prefix: "CPS_"})
	a := CPSAction{
		ActionCode: actionId,
		MakerAndChecker: MakerAndChecker{
			Maker: User{
				UserID: maker_id,
			},
		},
		ActionType:    ActionCreate,
		RequestAction: RequestUpdateServiceRule,
		ActionStatus:  ActionPending,
		CurrentAction: CurrentAction{
			Id:     id,
			Action: action,
		},
	}
	result, err := s.repository.CreateCpsAction(ctx, a)
	if err != nil {
		return "", nil
	}
	return result.ID, nil
}
func (s *ServiceStore) RemoveCif(ctx context.Context, action_id string, action bool, checker_id string) error {
	cps_action, err := s.repository.FetchCpsActionById(ctx, action_id)
	if err != nil {
		return err
	}
	if action {
		cps_action.ActionStatus = ActionApproved
	} else {
		cps_action.ActionStatus = ActionRejected
	}
	cps_action.MakerAndChecker.Checker.UserID = checker_id
	err = s.repository.UpdateCpsAction(ctx, cps_action)
	if err != nil {
		return err
	}
	currentAction, ok := cps_action.CurrentAction.(CurrentAction)
	if !ok {
		return fmt.Errorf("failed to cast CurrentAction to its expected type")
	}
	linked_accounts, err := s.repository.FetchLinkedAccountById(ctx, currentAction.Id)
	if err != nil {
		return err
	}
	for _, account := range linked_accounts {
		account.LinkedStatus = false
		account.IsAccountActive = false
		_, err = s.repository.UpdateAccount(ctx, account)
		if err != nil {
			return fmt.Errorf("failed to update account %s: %w", account.AccountNumber, err)
		}
	}
	return nil
}
=======
	service, err := s.repository.GetHqServiceById(ctx, cps_action.CurrentAction.Id)
	if err != nil {
		return err
	}
	service.Flag = cps_action.CurrentAction.Action
	return s.repository.UpdateHqService(ctx, service)
}
>>>>>>> b69ee66 (feature added)
