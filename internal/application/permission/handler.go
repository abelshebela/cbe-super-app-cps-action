package permission

import (
	"github.com/CBE-Super-App/cbe-super-app-cps-action/internal/adapter/outbound/model"
	domain "github.com/CBE-Super-App/cbe-super-app-cps-action/internal/domain/permission"
	// "github.com/CBE-Super-App/cbe-super-app-cps-action/internal/domain/permission/entities"

	"gitlab.com/bersufekadgetachew/cbe-super-app-shared/shared/utils"
)

type PermissionService interface {
	CreatePermissionGroup(groupName, role string, permissionCategoryIDs []string, cpsAction model.CPSAction) (model.CPSAction, error)
	ApprovePermissionGroup(actionCode string, action model.CPSAction) error
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

func (h *PermissionHandler) CreatePermissionGroup(groupName, role string, permissionCategoryLists []string, cpsAction model.CPSAction) (model.CPSAction, error) {
	h.logger.Infof("Handler: Initiating CreatePermissionGroup with groupName: %s, role: %s, makerID: %s", groupName, role, cpsAction.MakerID)

	cpsAction, err := h.service.CreatePermissionGroup(groupName, role, permissionCategoryLists, cpsAction)
	if err != nil {
		h.logger.Errorf("Handler: Failed to create permission group '%s' for role '%s': %v", groupName, role, err)
		return model.CPSAction{}, err
	}

	h.logger.Infof("Handler: Successfully initiated permission group creation for '%s'", groupName)
	return cpsAction, nil
}

func (h *PermissionHandler) ApprovePermissionGroup(actionCode string, action model.CPSAction) error {
	err := h.service.ApprovePermissionGroup(actionCode, action)
	if err != nil {
		h.logger.Errorf("Failed to approve permission group: %v", err)
		return err
	}
	return nil
}
