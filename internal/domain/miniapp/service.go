package miniapp

import (
	"context"
	"encoding/json"
	"fmt"
	"time"

	"gitlab.com/bersufekadgetachew/cbe-super-app-shared/shared/utils"
	// "go.mongodb.org/mongo-driver/v2/bson"
	// "go.mongodb.org/mongo-driver/v2/bson"

	"github.com/CBE-Super-App/cbe-super-app-cps-action/internal/adapter/outbound/model"
	domain "github.com/CBE-Super-App/cbe-super-app-cps-action/internal/domain/action"
	// "go.mongodb.org/mongo-driver/bson/primitive"
	// "go.mongodb.org/mongo-driver/v2/bson"
)

func (s *MiniAppStore) CreateMiniAppAction(ctx context.Context, miniApp MiniApp, maker model.User, department string) (string, error) {
	actionId := utils.Random(10, &utils.PreSufix{Prefix: "CPS_"})
	// actionData := miniApp

	fmt.Println("miniApp to maker ++++++++++++++++++++++++++++++++++++++++++++++")
	fmt.Println(miniApp)
	fmt.Println("miniApp to maker ++++++++++++++++++++++++++++++++++++++++++++++")

	// var prevAction json.RawMessage
	// var currAction json.RawMessage

	// prevAction, _ = json.Marshal(miniApp)

	// miniApp.IsDeleted = false
	// currAction, err := json.Marshal(miniApp)
	// if err != nil {
	// 	s.logger.Errorf("failed to marshal miniapp request %v", err)
	// 	return "", err
	// }
	// current := string(currAction)
	currentMiniApp := miniApp
	currentMiniApp.IsDeleted = true
	action := model.CPSAction{
		ActionCode:       actionId,
		MakerID:          maker.UserCode,
		MakerName:        maker.FullName,
		MakerPhoneNumber: maker.PhoneNumber,
		Department:       department,
		MakerActionTime:  time.Now(),
		ActionType:       string(model.ActionCreate),
		RequestAction:    string(domain.RequestCreateMiniAppMerchant),
		ActionStatus:     string(domain.ActionPending),
		CurrentAction:    currentMiniApp,
		PreviosAction:    miniApp,
		CreatedAt:        time.Now(),
		LastModifiedAt:   time.Now(),
	}
	a, err := s.Repository.CreateMiniAppAction(ctx, action)
	if err != nil {
		return "", err
	}
	return a.ActionCode, nil
}

func (s *MiniAppStore) CheckMiniApp(ctx context.Context, actionId string, action bool, checker model.User, department string) error {
	act, err := s.Repository.GetMiniAppActionId(ctx, actionId)
	if err != nil {
		fmt.Println("Error fetching action:", err)
		return err
	}

	if act.Department != department {
		s.logger.Errorf("Department mismatch: action department = %s, checker department = %s\n", act.Department, department)
		return fmt.Errorf("Maker and Chaker department mismatch")
	}

	if action {
		act.ActionType = string(model.ActionUpdate)
		act.ActionStatus = string(model.ActionApproved)
	} else {
		act.ActionType = string(model.ActionUpdate)
		act.ActionStatus = string(model.ActionRejected)
	}
	act.CheckerID = checker.UserCode
	act.CheckerName = checker.FullName
	act.CheckerPhoneNumber = checker.PhoneNumber
	act.CheckerActionTime = time.Now()
	act.LastModifiedAt = time.Now()
	updatedAction, err := s.Repository.UpdateCpsAction(ctx, act)

	// var miniApp MiniApp
	fmt.Printf("Type of updatedAction.CurrentAction: %T\n", updatedAction.CurrentAction)
	var miniApp MiniApp

	currentActionStr, ok := updatedAction.CurrentAction.(string)
	if !ok {
		return fmt.Errorf("invalid type: expected string but got %T", updatedAction.CurrentAction)
	}
	err = json.Unmarshal([]byte(currentActionStr), &miniApp)
	if err != nil {
		return fmt.Errorf("unmarshal error: %v", err)
	}
	fmt.Println("feached miniApp===================")
	fmt.Println(miniApp)
	fmt.Println("feached miniApp===================")

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

func (s *MiniAppStore) DeleteMiniAppAction(ctx context.Context, maker model.User, id string) (string, error) {
	actionId := utils.Random(10, &utils.PreSufix{Prefix: "CPS_DEL_"})
	// objectID, err := bson.ObjectIDFromHex(id)
	// if err != nil {
	// 	return "", fmt.Errorf("invalid id: %w", err)
	// }
	miniApp, err := s.Repository.DetailMiniAppByID(ctx, id)
	currentAction := miniApp
	currentAction.IsDeleted = true
	if err != nil {
		return "", fmt.Errorf("failed to fetch mini app by id: %w", err)
	}

	action := model.CPSAction{
		ActionCode:       actionId,
		MakerID:          maker.UserCode,
		MakerName:        maker.FullName,
		MakerPhoneNumber: maker.PhoneNumber,
		MakerActionTime:  time.Now(),
		ActionType:       string(domain.ActionDelete),
		RequestAction:    string(domain.RequestDeleteMiniAppMerchant),
		ActionStatus:     string(domain.ActionPending),
		PreviosAction:    miniApp,
		CurrentAction:    currentAction,
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
