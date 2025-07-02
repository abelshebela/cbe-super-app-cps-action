package miniapp

import (
	"context"
	"fmt"

	"gitlab.com/bersufekadgetachew/cbe-super-app-shared/shared/utils"

	domain "github.com/CBE-Super-App/cbe-super-app-cps-action/internal/domain/action"
)

func (s *MiniAppStore) CreateMiniAppAction(ctx context.Context, miniApp MiniApp, makerId string) (string, error) {
	actionId := utils.Random(10, &utils.PreSufix{Prefix: "CPS_"})
	previos_action := miniApp
	action := domain.CPSAction{
		ActionCode: actionId,
		Maker: domain.User{
			UserID: makerId,
		},
		Checker:       domain.User{},
		ActionType:    domain.ActionCreate,
		RequestAction: domain.RequestCreateMiniAppMerchant,
		ActionStatus:  domain.ActionPending,
		PreviosAction: previos_action,
	}
	a, err := s.Repository.CreateMiniAppAction(ctx, action)
	if err != nil {
		return "", err
	}
	return a.ActionCode, nil
}

func (s *MiniAppStore) CheckMiniApp(ctx context.Context, actionId string, action bool, checkerId string) error {
	act, err := s.Repository.GetMiniAppActionId(ctx, actionId)
	if err != nil {
		return err
	}
	if action {
		act.ActionType = domain.ActionType(domain.ActionApproved)
		act.ActionStatus = domain.ActionApproved
	} else {
		act.ActionType = domain.ActionType(domain.ActionRejected)
		act.ActionStatus = domain.ActionRejected
	}
	act.Checker = domain.User{
		UserID: checkerId,
	}
	err = s.Repository.UpdateCpsAction(ctx, act)
	if err != nil {
		return err
	}
	data := act.PreviosAction
	miniApp, ok := data.(MiniApp)
	if !ok {
		return fmt.Errorf("failed to cast previous action data to MiniApp")
	}
	err = s.Repository.CreateMiniApp(ctx, miniApp)
	if err != nil {
		return err
	}
	return nil
}
