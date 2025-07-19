package services

import (
	"context"

	entities "github.com/CBE-Super-App/cbe-super-app-cps-action/internal/domain/cps_actions/entities"
	repo "github.com/CBE-Super-App/cbe-super-app-cps-action/internal/domain/cps_actions/repository"
	common_util "github.com/CBE-Super-App/cbe-super-app-cps-action/pkgs/utils"
	"gitlab.com/bersufekadgetachew/cbe-super-app-shared/shared/utils"
)

type CPSActionService interface {
	CreateCPSAction(ctx context.Context, action *entities.CPSAction) (*entities.CPSAction, error)
	UpdateCPSAction(ctx context.Context, action *entities.CPSAction) (*entities.CPSAction, error)
	CPSActionExists(ctx context.Context, uniqueID string) (bool, error)
	ApproveCPSAction(ctx context.Context, action *entities.CPSAction) (*entities.CPSAction, error)
	RejectCPSAction(ctx context.Context, action *entities.CPSAction) (*entities.CPSAction, error)
	GetCPSActionsByDepartment(ctx context.Context, department string) (*common_util.PaginatedResponse[[]*entities.CPSAction], error)
	GetCPSActionByID(ctx context.Context, id string) (*entities.CPSAction, error)
	GetCPSActionByUniqueID(ctx context.Context, uniqueID string) (*entities.CPSAction, error)
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
func (s *cpsActionService) CPSActionExists(ctx context.Context, uniqueID string) (bool, error) {
	s.logger.Infof("Checking if CPS Action exists with unique ID: %s", uniqueID)
	return s.repo.CPSActionExists(ctx, uniqueID)
}
func (s *cpsActionService) ApproveCPSAction(ctx context.Context, action *entities.CPSAction) (*entities.CPSAction, error) {
	s.logger.Infof("Approving CPS Action with ID: %s", action.ID)
	return s.repo.ApproveCPSAction(ctx, action)
}
func (s *cpsActionService) RejectCPSAction(ctx context.Context, action *entities.CPSAction) (*entities.CPSAction, error) {
	s.logger.Infof("Rejecting CPS Action with ID: %s", action.ID)
	return s.repo.RejectCPSAction(ctx, action)
}
func (s *cpsActionService) GetCPSActionsByDepartment(ctx context.Context, department string) (*common_util.PaginatedResponse[[]*entities.CPSAction], error) {
	s.logger.Infof("Fetching CPS Actions for department: %s", department)
	return s.repo.GetCPSActionsByDepartment(ctx, department)
}
func (s *cpsActionService) GetCPSActionByID(ctx context.Context, id string) (*entities.CPSAction, error) {
	s.logger.Infof("Fetching CPS Action by ID: %s", id)
	return s.repo.GetCPSActionByID(ctx, id)
}
func (s *cpsActionService) GetCPSActionByUniqueID(ctx context.Context, uniqueID string) (*entities.CPSAction, error) {
	s.logger.Infof("Fetching CPS Action by unique ID: %s", uniqueID)
	return s.repo.GetCPSActionByUniqueID(ctx, uniqueID)
}
