package permission

import (
	"errors"
	"time"

	"github.com/CBE-Super-App/cbe-super-app-cps-action/internal/domain/permission/entities"

	"gitlab.com/bersufekadgetachew/cbe-super-app-shared/shared/utils"
)

type Service struct {
	cpsActionRepo          CPSActionRepository
	permissionGroupRepo    PermissionGroupRepository
	permissionCategoryRepo PermissionCategoryRepository
	logger                 utils.Logger
}

func InitPermissionDomain(cpsActionRepo CPSActionRepository, permissionGroupRepo PermissionGroupRepository, permissionCategoryRepo PermissionCategoryRepository, logger utils.Logger) *Service {
	return &Service{
		cpsActionRepo:          cpsActionRepo,
		permissionGroupRepo:    permissionGroupRepo,
		permissionCategoryRepo: permissionCategoryRepo,
		logger:                 logger,
	}
}

func (s *Service) CreatePermissionGroup(groupName, role string, permissionCategoryLists []string, cpsAction entities.CPSAction) error {
	if err := s.cpsActionRepo.CheckPendingRequest(cpsAction.MakerID, cpsAction.ActionStatus, cpsAction.RequestAction); err != nil {
		s.logger.Warnf("Pending request check failed for user %s: %v", cpsAction.MakerID, err)
		return errors.New("you have pending request for this action")
	}

	if exists := s.permissionGroupRepo.CheckPermissionGroupExists(groupName); exists {
		s.logger.Warnf("Permission group already exists with name: %s", groupName)
		return errors.New("permission group already exists")
	}

	validCategories, err := s.permissionCategoryRepo.ValidatePermissionCategories(permissionCategoryLists)
	if err != nil {
		s.logger.Errorf("Failed to validate permission categories: %v", err)
		return err
	}

	cpsAction.ActionCode = utils.RandomGenerator(20)
	cpsAction.CurrentAction = map[string]interface{}{
		"group_name":            groupName,
		"permission_categories": validCategories,
		"role":                  role,
		"realm":                 "bank",
	}
	cpsAction.MakerActionTime = time.Now()

	err = s.cpsActionRepo.CreatePermissionGroup(cpsAction)
	if err != nil {
		s.logger.Errorf("Failed to create CPSAction: %v", err)
		return err
	}

	s.logger.Infof("Successfully created permission group request with ActionCode: %s", cpsAction.ActionCode)
	return nil
}

func (s *Service) ApprovePermissionGroup(actionCode string, action entities.CPSAction) error {
	action, err := s.cpsActionRepo.ValidateActionRequest(actionCode, action.Department)
	if err != nil {
		return err
	}

	switch action.ActionType {
	case "CREATE":
		err = s.permissionGroupRepo.CreatePermissionGroupFromAction(action)
	case "UPDATE":
		err = s.permissionGroupRepo.UpdatePermissionGroupFromAction(action)
	default:
		return errors.New("invalid action type")
	}

	if err != nil {
		return err
	}

	if err := s.cpsActionRepo.ApproveActionRequest(actionCode, action); err != nil {
		return err
	}
	return nil
}
