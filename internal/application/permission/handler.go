package permission

import (
	domain "github.com/CBE-Super-App/cbe-super-app-cps-action/internal/domain/permission"
	"github.com/CBE-Super-App/cbe-super-app-cps-action/internal/domain/permission/entities"

	"gitlab.com/bersufekadgetachew/cbe-super-app-shared/shared/utils"
)

type PermissionService interface {
	CreatePermissionGroup(groupName, role string, permissionCategoryIDs []string, cpsAction entities.CPSAction) error
	ApprovePermissionGroup(actionCode string, action entities.CPSAction) error
}

type PermissionHandler struct {
	service *domain.Service
	logger  utils.Logger
}

func InitPermissionHandler(service *domain.Service, logger utils.Logger) PermissionService {
	return &PermissionHandler{
		service: service,
		logger:  logger,
	}
}

func (h *PermissionHandler) CreatePermissionGroup(groupName, role string, permissionCategoryLists []string, cpsAction entities.CPSAction) error {
	err := h.service.CreatePermissionGroup(groupName, role, permissionCategoryLists, cpsAction)
	if err != nil {
		h.logger.Errorf("Failed to create permission group: %v", err)
		return err
	}
	return nil
}

func (h *PermissionHandler) ApprovePermissionGroup(actionCode string, action entities.CPSAction) error {
	err := h.service.ApprovePermissionGroup(actionCode, action)
	if err != nil {
		h.logger.Errorf("Failed to approve permission group: %v", err)
		return err
	}
	return nil
}
