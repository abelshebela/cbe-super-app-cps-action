package permission

import (
	// "fmt"
	"context"
	"strings"

	"github.com/CBE-Super-App/cbe-super-app-cps-action/internal/adapter/outbound/model"
	domain "github.com/CBE-Super-App/cbe-super-app-cps-action/internal/domain/permission"

	// "github.com/CBE-Super-App/cbe-super-app-cps-action/pkgs/entities"
	entities "github.com/CBE-Super-App/cbe-super-app-cps-action/internal/domain/permission/entities"

	common_util "github.com/CBE-Super-App/cbe-super-app-cps-action/pkgs/utils"
	constant "github.com/CBE-Super-App/cbe-super-app-cps-action/utils"
	"gitlab.com/bersufekadgetachew/cbe-super-app-shared/shared/utils"
)

type PermissionService interface {
	CreatePermissionGroup(ctx context.Context,oldGroupName, groupName, role string, permissionCategoryIDs []string, cpsAction model.CPSAction) (model.CPSAction, error)
	GetPermissionGroups(ctx context.Context, filterParams *constant.Filter) (*common_util.PaginatedResponse[[]*entities.PermissionGroup], error)
	GetPermissionGroup(groupName string) (entities.PermissionGroup, error)
	UpdatePermissionGroup(groupName string, permissionCategoryIDs []string) (entities.PermissionGroup, error)
	UpdatePermissionGroupRequest(ctx context.Context,oldGroupName, groupName, role string, permissionCategoryIDs []string, cpsAction model.CPSAction) (model.CPSAction, error)
	GetAllPermissionCategoriesWithPermissions(ctx context.Context) ([]*entities.PermissionCategory, error)
}

type PermissionHandler struct {
	service domain.PermissionDomainService
	logger  utils.Logger
}

func InitPermissionHandler(service *domain.Service, logger utils.Logger) PermissionService {
	return &PermissionHandler{
		service: service,
		logger:  logger,
	}
}

func (h *PermissionHandler) CreatePermissionGroup(ctx context.Context,oldGroupName, groupName, role string, permissionCategoryLists []string, cpsAction model.CPSAction) (model.CPSAction, error) {
	oldGroupName = strings.ToUpper(oldGroupName)
	groupName = strings.ToUpper(groupName)
	h.logger.Infof("Handler: Initiating Create cps action for PermissionGroup with groupName: %s, role: %s, makerID: %s", groupName, role, cpsAction.MakerID)

	cpsAction, err := h.service.CreatePermissionGroup(ctx,oldGroupName, groupName, role, permissionCategoryLists, cpsAction)
	if err != nil {
		h.logger.Errorf("Handler: Failed to create cps action for  permission group '%s' for role '%s': %v", groupName, role, err)
		return model.CPSAction{}, err
	}

	h.logger.Infof("Handler: Successfully initiated cps action for permission group creation for '%s'", groupName)
	return cpsAction, nil
}

func (h *PermissionHandler) UpdatePermissionGroupRequest(ctx context.Context,oldGroupName, groupName, role string, permissionCategoryLists []string, cpsAction model.CPSAction) (model.CPSAction, error) {
	oldGroupName = strings.ToUpper(oldGroupName)
	groupName = strings.ToUpper(groupName)
	perv_action, err := h.service.GetPermissionGroup(oldGroupName)
	if err != nil {
		h.logger.Errorf("old group not found", oldGroupName, err)
		return model.CPSAction{}, err
	}
	cpsAction.PreviousAction = perv_action

	
	cpsAction, err = h.service.CreatePermissionGroup(ctx,oldGroupName, groupName, role, permissionCategoryLists, cpsAction)

	if err != nil {
		h.logger.Errorf("Handler: Failed to update permission group '%s' for role '%s': %v", groupName, role, err)
		return model.CPSAction{}, err
	}

	h.logger.Infof("Handler: Successfully initiated cps action for permission group update for '%s'", groupName)
	return cpsAction, nil
}
func (h *PermissionHandler) GetPermissionGroups(ctx context.Context, filterParams *constant.Filter) (*common_util.PaginatedResponse[[]*entities.PermissionGroup], error) {
	return h.service.GetPermissionGroups(ctx, filterParams)
}

func (h *PermissionHandler) GetPermissionGroup(groupName string) (entities.PermissionGroup, error) {
	groupName = strings.ToUpper(groupName)
	return h.service.GetPermissionGroup(groupName)
}

func (h *PermissionHandler) UpdatePermissionGroup(groupName string, permissionCategoryLists []string) (entities.PermissionGroup, error) {
	groupName = strings.ToUpper(groupName)
	return h.service.UpdatePermissionGroup(groupName, permissionCategoryLists)
}

func (h *PermissionHandler) GetAllPermissionCategoriesWithPermissions(ctx context.Context) ([]*entities.PermissionCategory, error) {
	return h.service.GetAllPermissionCategoriesWithPermissions(ctx)
}
