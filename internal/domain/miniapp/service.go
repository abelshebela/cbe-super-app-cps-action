package miniapp

import (
	"context"
	"fmt"
	"time"

	"gitlab.com/bersufekadgetachew/cbe-super-app-shared/shared/utils"

	"github.com/CBE-Super-App/cbe-super-app-cps-action/internal/adapter/outbound/model"
	domain "github.com/CBE-Super-App/cbe-super-app-cps-action/internal/domain/action"
)

func (s *MiniAppStore) CreateMiniAppAction(ctx context.Context, miniApp MiniApp, maker model.User) (string, error) {
	actionId := utils.Random(10, &utils.PreSufix{Prefix: "CPS_"})
	previous_action := miniApp
	action := model.CPSAction{
		ActionCode:       actionId,
		MakerID:          maker.UserCode,
		MakerName:        maker.FullName,
		MakerPhoneNumber: maker.PhoneNumber,
		MakerActionTime:  time.Now(),
		ActionType:       string(model.ActionCreate),
		RequestAction:    string(domain.RequestCreateMiniAppMerchant),
		ActionStatus:     string(domain.ActionPending),
		PreviosAction:    previous_action,
	}
	a, err := s.Repository.CreateMiniAppAction(ctx, action)
	if err != nil {
		return "", err
	}
	return a.ActionCode, nil
}

func (s *MiniAppStore) CheckMiniApp(ctx context.Context, actionId string, action bool, checker model.User) error {
	act, err := s.Repository.GetMiniAppActionId(ctx, actionId)
	if err != nil {
		return err
	}
	if action {
		act.ActionType = string(model.ActionApproved)
		act.ActionStatus = string(model.ActionApproved)
	} else {
		act.ActionType = string(domain.ActionRejected)
		act.ActionStatus = string(domain.ActionRejected)
	}
	act.CheckerID = checker.UserCode
	act.CheckerName = checker.FullName
	act.CheckerPhoneNumber = checker.PhoneNumber
	act.CheckerActionTime = time.Now()
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

func (s *MiniAppStore) UpdateMiniAppAction(ctx context.Context, data MiniApp, maker model.User) (string, error) {
	actionId := utils.Random(10, &utils.PreSufix{Prefix: "CPS_UPD_"})
	previous_action := data
	action := model.CPSAction{
		ActionCode:       actionId,
		MakerID:          maker.UserCode,
		MakerName:        maker.FullName,
		MakerPhoneNumber: maker.PhoneNumber,
		MakerActionTime:  time.Now(),
		ActionType:       string(domain.ActionUpdate),
		RequestAction:    string(domain.RequestUpdateMiniAppMerchant),
		ActionStatus:     string(domain.ActionPending),
		PreviosAction:    previous_action,
	}
	a, err := s.Repository.CreateMiniAppAction(ctx, action)
	if err != nil {
		return "", err
	}
	return a.ActionCode, nil
}

func (s *MiniAppStore) DeleteMiniAppAction(ctx context.Context, maker model.User) (string, error) {
	actionId := utils.Random(10, &utils.PreSufix{Prefix: "CPS_DEL_"})
	action := model.CPSAction{
		ActionCode:       actionId,
		MakerID:          maker.UserCode,
		MakerName:        maker.FullName,
		MakerPhoneNumber: maker.PhoneNumber,
		MakerActionTime:  time.Now(),
		ActionType:       string(domain.ActionDelete),
		RequestAction:    string(domain.RequestDeleteMiniAppMerchant),
		ActionStatus:     string(domain.ActionPending),
	}
	a, err := s.Repository.CreateMiniAppAction(ctx, action)
	if err != nil {
		return "", err
	}
	return a.ActionCode, nil
}

func (s *MiniAppStore) ListMiniApp(ctx context.Context) ([]*MiniApp, error) {
	return s.Repository.ListMiniApp(ctx)
}

func (s *MiniAppStore) DetailMiniAppByID(ctx context.Context, id string) (MiniApp, error) {
	return s.Repository.DetailMiniAppByID(ctx, id)
}
