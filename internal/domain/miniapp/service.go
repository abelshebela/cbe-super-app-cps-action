package miniapp

import (
	"context"
	// "encoding/json"
	"fmt"
	"time"

	// bson "go.mongodb.org/mongo-driver/v2/bson"

	"github.com/CBE-Super-App/cbe-super-app-cps-action/internal/adapter/outbound/model"
	domain "github.com/CBE-Super-App/cbe-super-app-cps-action/internal/domain/action"
	entities "github.com/CBE-Super-App/cbe-super-app-cps-action/internal/domain/cps_actions/entities"

	"gitlab.com/bersufekadgetachew/cbe-super-app-shared/shared/utils"
)

type MiniAppStore struct {
	Repository MiniRepository
	logger     utils.Logger
}
type MiniAppService interface {
	CreateMiniAppAction(ctx context.Context, miniApp MiniApp, maker model.User, departmen string) (string, error)
	Authorize(ctx context.Context, action *entities.CPSAction) (*entities.CPSAction, error)
	UpdateMiniAppAction(ctx context.Context, data MiniApp, maker model.User) (string, error)

	DeleteMiniAppAction(ctx context.Context, maker model.User, id string) (string, error)
	ListMiniApp(ctx context.Context) ([]*model.MiniApp, error)
	DetailMiniAppByID(ctx context.Context, id string) (model.MiniApp, error)
}

func NewService(repository MiniRepository, logger utils.Logger) MiniAppService {
	return &MiniAppStore{
		Repository: repository,
		logger:     logger,
	}
}

func (s *MiniAppStore) CreateMiniAppAction(ctx context.Context, miniApp MiniApp, maker model.User, department string) (string, error) {
	actionId := utils.Random(10, &utils.PreSufix{Prefix: "CPS_"})
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
		CurrentAction:    miniApp,
		CreatedAt:        time.Now(),
		LastModifiedAt:   time.Now(),
	}

	a, err := s.Repository.CreateMiniAppAction(ctx, action)
	if err != nil {
		return "", err
	}
	return a.ActionCode, nil
}

func (s *MiniAppStore) Authorize(ctx context.Context, action *entities.CPSAction) (*entities.CPSAction, error) {

	_, err := s.Repository.CreateMiniApp(ctx, action)
	if err != nil {
		fmt.Printf("error form domain chekmiiapp to crate mini app : %v", err)
		return nil, err
	}
	return action, nil
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
		PreviousAction:   previous_action,
	}
	a, err := s.Repository.CreateMiniAppAction(ctx, action)
	if err != nil {
		return "", err
	}
	return a.ActionCode, nil
}

func (s *MiniAppStore) DeleteMiniAppAction(ctx context.Context, maker model.User, id string) (string, error) {
	actionId := utils.Random(10, &utils.PreSufix{Prefix: "CPS_DEL_"})
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
		PreviousAction:   miniApp,
		CurrentAction:    currentAction,
	}
	a, err := s.Repository.CreateMiniAppAction(ctx, action)
	if err != nil {
		return "", err
	}
	return a.ActionCode, nil
}

func (s *MiniAppStore) ListMiniApp(ctx context.Context) ([]*model.MiniApp, error) {
	miniApps, err := s.Repository.ListMiniApp(ctx)
	return miniApps, err
}

func (s *MiniAppStore) DetailMiniAppByID(ctx context.Context, id string) (model.MiniApp, error) {
	return s.Repository.DetailMiniAppByID(ctx, id)
}
