package miniapp

import (
	"context"
	"fmt"
	"time"

	// bson "go.mongodb.org/mongo-driver/v2/bson"

	"github.com/CBE-Super-App/cbe-super-app-cps-action/internal/adapter/outbound/model"
	"github.com/CBE-Super-App/cbe-super-app-cps-action/internal/domain/cps_actions/constant"
	entities "github.com/CBE-Super-App/cbe-super-app-cps-action/internal/domain/cps_actions/entities"

	outbound "github.com/CBE-Super-App/cbe-super-app-cps-action/internal/port/outbound/mini_app"

	"gitlab.com/bersufekadgetachew/cbe-super-app-shared/shared/utils"
)

type MiniAppStore struct {
	Repository outbound.MiniRepository
	logger     utils.Logger
}
type MiniAppService interface {
	CreateMiniAppAction(ctx context.Context, miniApp MiniApp, maker model.User) (*entities.CPSAction, error)
	Authorize(ctx context.Context, action *entities.CPSAction) (*entities.CPSAction, error)
	UpdateMiniAppAction(ctx context.Context, data MiniApp, maker model.User, id string) (*entities.CPSAction, error)

	DeleteMiniAppAction(ctx context.Context, maker model.User, id string) (*entities.CPSAction, error)
	ListMiniApp(ctx context.Context) ([]*model.MiniApp, error)
	DetailMiniAppByID(ctx context.Context, id string) (model.MiniApp, error)
}

func NewService(repository outbound.MiniRepository, logger utils.Logger) MiniAppService {
	return &MiniAppStore{
		Repository: repository,
		logger:     logger,
	}
}

func (s *MiniAppStore) CreateMiniAppAction(ctx context.Context, miniApp MiniApp, maker model.User) (*entities.CPSAction, error) {
	action := entities.CPSAction{
		ActionCode:       utils.RandomGenerator(20),
		MakerID:          maker.UserCode,
		MakerName:        maker.FullName,
		MakerPhoneNumber: maker.PhoneNumber,
		Department:       maker.Department,
		MakerActionTime:  time.Now(),
		ActionType:       constant.ActionCreate,
		RequestAction:    constant.RequestCreateMiniAppMerchant,
		ActionStatus:     constant.ActionPending,
		CurrentAction:    miniApp,
		CreatedAt:        time.Now(),
		LastModifiedAt:   time.Now(),
	}

	cpsAction, err := s.Repository.CreateMiniAppAction(ctx, &action)
	if err != nil {
		return nil, err
	}
	return cpsAction, nil
}

func (s *MiniAppStore) Authorize(ctx context.Context, action *entities.CPSAction) (*entities.CPSAction, error) {

	prevData, Serr := s.Repository.ExtractActionData(action.PreviousAction)
	if Serr != nil {
		return nil, Serr
	}

	requestedAction := action.RequestAction
	var minApp *model.MiniApp
	var err error
	switch requestedAction {
	case constant.RequestCreateMiniAppMerchant:
		minApp, err = s.Repository.CreateMiniApp(ctx, action)
	case constant.RequestUpdateMiniAppMerchant:
		minApp, err = s.Repository.UpdateMinApp(ctx, action, prevData.ID.Hex())
	case constant.RequestDeleteMiniAppMerchant:
		minApp, err = s.Repository.DeleteMiniAppAction(ctx, prevData.ID.Hex())

	default:
		return nil, fmt.Errorf("UNSUPPORTED_REQUEST_ACTION")
	}
	if err != nil {
		fmt.Printf("error form domain chekmiiapp to crate mini app : %v", err)
		return nil, err
	}

	action.CurrentAction = minApp
	return action, nil
}

func (s *MiniAppStore) UpdateMiniAppAction(ctx context.Context, data MiniApp, maker model.User, id string) (*entities.CPSAction, error) {

	prevData, err := s.Repository.DetailMiniAppByID(ctx, id)
	if err != nil {
		return nil, err
	}
	action := entities.CPSAction{
		ActionCode:       utils.RandomGenerator(20),
		MakerID:          maker.UserCode,
		MakerName:        maker.FullName,
		MakerPhoneNumber: maker.PhoneNumber,
		Department:       maker.Department,
		MakerActionTime:  time.Now(),
		ActionType:       constant.ActionUpdate,
		RequestAction:    constant.RequestUpdateMiniAppMerchant,
		ActionStatus:     constant.ActionPending,
		CurrentAction:    data,
		PreviousAction:   prevData,
	}
	a, err := s.Repository.CreateMiniAppAction(ctx, &action)
	if err != nil {
		return nil, err
	}
	return a, nil
}

func (s *MiniAppStore) DeleteMiniAppAction(ctx context.Context, maker model.User, id string) (*entities.CPSAction, error) {
	actionId := utils.Random(10, &utils.PreSufix{Prefix: "CPS_DEL_"})
	miniApp, err := s.Repository.DetailMiniAppByID(ctx, id)
	currentAction := miniApp
	currentAction.IsDeleted = true
	if err != nil {
		return nil, fmt.Errorf("failed to fetch mini app by id: %w", err)
	}

	action := entities.CPSAction{
		ActionCode:       actionId,
		MakerID:          maker.UserCode,
		MakerName:        maker.FullName,
		MakerPhoneNumber: maker.PhoneNumber,
		MakerActionTime:  time.Now(),
		Department:       maker.Department,
		ActionType:       constant.ActionDelete,
		RequestAction:    constant.RequestDeleteMiniAppMerchant,
		ActionStatus:     constant.ActionPending,
		PreviousAction:   miniApp,
		CurrentAction:    currentAction,
		LastModifiedAt:   time.Now(),
	}
	a, err := s.Repository.CreateMiniAppAction(ctx, &action)
	if err != nil {
		return nil, err
	}
	return a, nil
}

func (s *MiniAppStore) ListMiniApp(ctx context.Context) ([]*model.MiniApp, error) {
	miniApps, err := s.Repository.ListMiniApp(ctx)
	return miniApps, err
}

func (s *MiniAppStore) DetailMiniAppByID(ctx context.Context, id string) (model.MiniApp, error) {
	return s.Repository.DetailMiniAppByID(ctx, id)
}
