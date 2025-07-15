package permission

import (
	"errors"
	"fmt"
	"time"

	"github.com/CBE-Super-App/cbe-super-app-cps-action/internal/adapter/outbound/model"
	// "github.com/CBE-Super-App/cbe-super-app-cps-action/internal/domain/permission/entities"

	"gitlab.com/bersufekadgetachew/cbe-super-app-shared/shared/utils"
)

type PermissionDomainService interface {
	CreatePermissionGroup(groupName, role string, permissionCategoryLists []string, cpsAction model.CPSAction) (model.CPSAction, error)
	ApprovePermissionGroup(actionCode string, action model.CPSAction) (model.CPSAction, error)
	RejectPermissionGroup(actionCode string, action model.CPSAction, reason string) (model.CPSAction, error)
}

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

func (s *Service) CreatePermissionGroup(groupName, role string, permissionCategoryLists []string, cpsAction model.CPSAction) (model.CPSAction, error) {
	if err := s.cpsActionRepo.CheckPendingRequest(cpsAction.MakerID, model.ActionStatus(cpsAction.ActionStatus), model.RequestAction(cpsAction.RequestAction)); err != nil {
		s.logger.Warnf("Pending request check failed for user %s: %v", cpsAction.MakerID, err)
		return model.CPSAction{}, fmt.Errorf("PENDING_REQUEST_CHECK_FAILED_FOR_CREATE_PERMISSION")
	}

	if exists := s.permissionGroupRepo.CheckPermissionGroupExists(groupName); exists {
		s.logger.Warnf("Permission group already exists with name: %s", groupName)
		return model.CPSAction{}, fmt.Errorf("PERMISSION_GROUP_ALREADY_EXIXTS")
	}

	validCategories, err := s.permissionCategoryRepo.ValidatePermissionCategories(permissionCategoryLists)
	if err != nil {
		s.logger.Errorf("Failed to validate permission categories: %v", err)
		return model.CPSAction{}, err
	}

	cpsAction.ActionCode = utils.RandomGenerator(20)
	cpsAction.CurrentAction = map[string]interface{}{
		"group_name":            groupName,
		"permission_categories": validCategories,
		"role":                  role,
		"realm":                 "bank",
	}

	cpsAction.MakerActionTime = time.Now()

	cpsAction, err = s.cpsActionRepo.CreatePermissionGroup(cpsAction)
	if err != nil {
		s.logger.Errorf("Failed to create CPSAction: %v", err)
		return model.CPSAction{}, err
	}

	s.logger.Infof("Successfully created permission group request with ActionCode: %s", cpsAction.ActionCode)
	return cpsAction, nil
}

func (s *Service) ApprovePermissionGroup(actionCode string, action model.CPSAction) (model.CPSAction, error) {

	CreatedAction, err := s.cpsActionRepo.ValidateActionRequest(actionCode, action.Department)
	if err != nil {
		return model.CPSAction{}, err
	}

	// Set checker fields from input action to CreatedAction
	CreatedAction.CheckerID = action.CheckerID
	CreatedAction.CheckerName = action.CheckerName
	CreatedAction.CheckerPhoneNumber = action.CheckerPhoneNumber
	// CreatedAction.Department = action.Department

	switch CreatedAction.ActionType {
	case "CREATE":
		err = s.permissionGroupRepo.CreatePermissionGroupFromAction(CreatedAction)
	case "UPDATE":
		err = s.permissionGroupRepo.UpdatePermissionGroupFromAction(CreatedAction)
	default:
		return model.CPSAction{}, errors.New("invalid action type")
	}

	if err != nil {
		return model.CPSAction{}, err
	}

	approvedAction, err := s.cpsActionRepo.ApproveActionRequest(actionCode, CreatedAction)
	if err != nil {
		return model.CPSAction{}, err
	}
	return approvedAction, nil
}

func (s *Service) RejectPermissionGroup(actionCode string, action model.CPSAction, reason string) (model.CPSAction, error) {
	CreatedAction, err := s.cpsActionRepo.ValidateActionRequest(actionCode, action.Department)
	if err != nil {
		return model.CPSAction{}, err
	}
	CreatedAction.CheckerID = action.CheckerID
	CreatedAction.CheckerName = action.CheckerName
	CreatedAction.CheckerPhoneNumber = action.CheckerPhoneNumber

	rejectedAction, err := s.cpsActionRepo.RejectActionRequest(actionCode, CreatedAction, reason)
	if err != nil {
		return model.CPSAction{}, err
	}
	return rejectedAction, nil
}
