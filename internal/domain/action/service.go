package action

import (
	"context"

	"gitlab.com/bersufekadgetachew/cbe-super-app-shared/shared/utils"
)

func (s *ServiceStore) GetServicePaginated(ctx context.Context, limit, offset int) ([]Service, error) {
	return s.repository.GetAllHqServicesPaginated(ctx, offset, limit)
}
func (s *ServiceStore) UpdateServiceFlagRequest(ctx context.Context, id string, action bool, makerId string) (string, error) {
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
			Id:     id,
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
	service, err := s.repository.GetHqServiceById(ctx, cps_action.CurrentAction.Id)
	if err != nil {
		return err
	}
	service.Flag = cps_action.CurrentAction.Action
	return s.repository.UpdateHqService(ctx, service)
}
