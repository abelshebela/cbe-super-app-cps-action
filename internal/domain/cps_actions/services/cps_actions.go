package services

import (
	"context"
	"time"

	"github.com/CBE-Super-App/cbe-super-app-cps-action/internal/adapter/outbound/model"
	entities "github.com/CBE-Super-App/cbe-super-app-cps-action/internal/domain/cps_actions/entities"
	repo "github.com/CBE-Super-App/cbe-super-app-cps-action/internal/domain/cps_actions/repository"
	common_util "github.com/CBE-Super-App/cbe-super-app-cps-action/pkgs/utils"
	constant "github.com/CBE-Super-App/cbe-super-app-cps-action/utils"

	"gitlab.com/bersufekadgetachew/cbe-super-app-shared/shared/utils"
)

type CPSActionService interface {
	CreateCPSAction(ctx context.Context, action *entities.CPSAction) (*entities.CPSAction, error)
	UpdateCPSAction(ctx context.Context, action *entities.CPSAction) (*entities.CPSAction, error)
	CPSActionExists(ctx context.Context, user entities.CheckCPSAction) (bool, error)
	ApproveCPSAction(ctx context.Context, action *entities.AuthorizeCPSAction) (*entities.CPSAction, error)
	RejectCPSAction(ctx context.Context, action *entities.AuthorizeCPSAction) (*entities.CPSAction, error)
	GetCPSActionsByDepartment(ctx context.Context, department string, status string, filterParams *constant.Filter) (*common_util.PaginatedResponse[[]*entities.CPSAction], error)
	GetCPSActionByID(ctx context.Context, id string) (*entities.CPSAction, error)
	GetCPSActionByActionCode(ctx context.Context, uniqueID string) (*entities.CPSAction, error)
	BuildCPSAction(ctx context.Context, request entities.CreateCPSRequest) *entities.CPSAction
}

type cpsActionService struct {
	repo   repo.CPSActionRepository
	logger utils.Logger
}

func NewCPSActionService(repo repo.CPSActionRepository, logger utils.Logger) CPSActionService {
	return &cpsActionService{
		repo:   repo,
		logger: logger,
	}
}
func (s *cpsActionService) CreateCPSAction(ctx context.Context, action *entities.CPSAction) (*entities.CPSAction, error) {
	s.logger.Infof("Creating CPS Action with unique ID: %s", action.UniqueID)
	return s.repo.CreateCPSAction(ctx, action)
}
func (s *cpsActionService) UpdateCPSAction(ctx context.Context, action *entities.CPSAction) (*entities.CPSAction, error) {
	s.logger.Infof("Updating CPS Action with ID: %s", action.ID)
	return s.repo.UpdateCPSAction(ctx, action)
}
func (s *cpsActionService) CPSActionExists(ctx context.Context, user entities.CheckCPSAction) (bool, error) {
	s.logger.Infof("Checking if CPS Action exists with unique ID: %s", user.UserCode)
	return s.repo.CPSActionExists(ctx, user)
}
func (s *cpsActionService) ApproveCPSAction(ctx context.Context, action *entities.AuthorizeCPSAction) (*entities.CPSAction, error) {
	return s.repo.ApproveCPSAction(ctx, action)
}
func (s *cpsActionService) RejectCPSAction(ctx context.Context, action *entities.AuthorizeCPSAction) (*entities.CPSAction, error) {
	s.logger.Infof("Rejecting CPS Action with ID: %s", action.ActionCode)
	return s.repo.RejectCPSAction(ctx, action)
}
func (s *cpsActionService) GetCPSActionsByDepartment(ctx context.Context, department string, status string, filterParams *constant.Filter) (*common_util.PaginatedResponse[[]*entities.CPSAction], error) {
	s.logger.Infof("Fetching CPS Actions for department: %s", department)
	if status == "" {
		status = string(model.ActionPending)
	}
	return s.repo.GetCPSActionsByDepartment(ctx, department, status, filterParams)
}
func (s *cpsActionService) GetCPSActionByID(ctx context.Context, id string) (*entities.CPSAction, error) {
	s.logger.Infof("Fetching CPS Action by ID: %s", id)
	return s.repo.GetCPSActionByID(ctx, id)
}
func (s *cpsActionService) GetCPSActionByActionCode(ctx context.Context, uniqueID string) (*entities.CPSAction, error) {
	s.logger.Infof("Fetching CPS Action by unique ID: %s", uniqueID)
	return s.repo.GetCPSActionByActionCode(ctx, uniqueID)
}

func (s *cpsActionService) BuildCPSAction(ctx context.Context, request entities.CreateCPSRequest) *entities.CPSAction {
	return &entities.CPSAction{
		ActionCode:       utils.RandomGenerator(20),
		MakerID:          request.User.UserCode,
		MakerName:        request.User.FullName,
		MakerPhoneNumber: request.User.PhoneNumber,
		Department:       request.User.Department,
		CurrentAction:    request.CurData,
		PreviousAction:   request.PrevData,
		RequestAction:    request.RequestAction,
		ActionStatus:     request.ActionStatus,
		ActionType:       request.ActionType,
		MakerActionTime:  time.Now(),
		CreatedAt:        time.Now(),
		LastModifiedAt:   time.Now(),
	}
}
