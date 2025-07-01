package permission

import (
	"cbe-super-app-cps-action/internal/domain/permission/entities"
	"errors"
	"time"

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
		return errors.New("you have pending request for this action")
	}

	if exists := s.permissionGroupRepo.CheckPermissionGroupExists(groupName); exists {
		return errors.New("permission group already exists")
	}

	validCategories, err := s.permissionCategoryRepo.ValidatePermissionCategories(permissionCategoryLists)
	if err != nil {
		return err
	}

	cpsAction.ActionCode = utils.Random(10, &utils.PreSufix{Prefix: "PER_GROUP_"})
	cpsAction.CurrentAction = map[string]interface{}{
		"group_name":            groupName,
		"permission_categories": validCategories,
		"role":                  role,
		"realm":                 "bank",
	}
	cpsAction.MakerActionTime = time.Now()

	return s.cpsActionRepo.CreatePermissionGroup(cpsAction)
}

func (s *Service) ApprovePermissionGroup(actionCode string, action entities.CPSAction) error {
	action, err := s.cpsActionRepo.ValidateActionRequest(actionCode, action.Department)
	if err != nil {
		return err
	}

	switch action.ActionType {
	case "CREATE":
		return s.permissionGroupRepo.CreatePermissionGroupFromAction(action)
	case "UPDATE":
		return s.permissionGroupRepo.UpdatePermissionGroupFromAction(action)
	default:
		return errors.New("invalid action type")
	}

	if err := s.cpsActionRepo.ApproveActionRequest(actionCode, action); err != nil {
		return err
	}
	return nil
}
