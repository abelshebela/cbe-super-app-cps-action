// Package department provides services and business logic for managing departments and related CPS actions.
package department

import (
	"context"
	"fmt"
	"time"

	"github.com/CBE-Super-App/cbe-super-app-cps-action/internal/domain/department/entities"
	"github.com/CBE-Super-App/cbe-super-app-cps-action/internal/domain/permission"
	portal_card "github.com/CBE-Super-App/cbe-super-app-cps-action/internal/domain/portal_card"
	common_util "github.com/CBE-Super-App/cbe-super-app-cps-action/pkgs/utils"
	constant "github.com/CBE-Super-App/cbe-super-app-cps-action/utils"
	"go.mongodb.org/mongo-driver/v2/bson"

	cps_const "github.com/CBE-Super-App/cbe-super-app-cps-action/internal/domain/cps_actions/constant"
	cpsactions "github.com/CBE-Super-App/cbe-super-app-cps-action/internal/domain/cps_actions/entities"
	"gitlab.com/bersufekadgetachew/cbe-super-app-shared/shared/utils"
)

const (
	CPSPrefix  = "CPS_"
	DeptPrefix = "DEP_"
)

type UpdateDepartmentRequest struct {
	Department       string   `json:"department"`
	PortalCards      []string `json:"portal_cards"`
	PermissionGroups []string `json:"permission_groups"`
}

// Service defines the interface for department operations
type Service interface {
	CheckDepartmentExists(ctx context.Context, department string) (bool, error)
	CheckDepartmentExistsByID(ctx context.Context, id string) (bool, error)
	CreateDepartment(ctx context.Context, department string, portalCards []string, permissionGroups []string, cpsAction interface{}) error
	UpdateDepartment(ctx context.Context, id string, updateData map[string]interface{}) (*entities.Department, error)
	GetAllDepartments(ctx context.Context, filterParams *constant.Filter) (*common_util.PaginatedResponse[[]*entities.Department], error)
	GetDepartmentByID(ctx context.Context, id string) (*entities.Department, error)
	Authorize(ctx context.Context, action *cpsactions.CPSAction) (*cpsactions.CPSAction, error)
}

// ServiceImpl implements the Service interface
type ServiceImpl struct {
	departmentRepo    DepartmentRepository
	permissionService permission.PermissionDomainService
	portalService     portal_card.PortaCardInterface
	logger            utils.Logger
}

func InitDepartmentDomain(
	departmentRepo DepartmentRepository,
	permissionService permission.PermissionDomainService,
	portalService portal_card.PortaCardInterface,
	logger utils.Logger,
) Service {
	return &ServiceImpl{
		departmentRepo:    departmentRepo,
		permissionService: permissionService,
		portalService:     portalService,
		logger:            logger,
	}
}

func (s *ServiceImpl) CheckDepartmentExists(ctx context.Context, department string) (bool, error) {
	return s.departmentRepo.CheckDepartmentExists(ctx, department)
}

func (s *ServiceImpl) CheckDepartmentExistsByID(ctx context.Context, id string) (bool, error) {
	return s.departmentRepo.CheckDepartmentExistsByID(ctx, id)
}

func (s *ServiceImpl) CreateDepartment(ctx context.Context, department string, portalCards []string, permissionGroups []string, cpsAction interface{}) error {

	var permissionGroupIDs []bson.ObjectID
	for _, groupID := range permissionGroups {
		objectID, err := bson.ObjectIDFromHex(groupID)
		if err != nil {
			return fmt.Errorf("DEPARTMENT_INVALID_PERMISSION_GROUP_ID")
		}
		permissionGroupIDs = append(permissionGroupIDs, objectID)
	}

	dept := entities.Department{
		DepartmentCode:   utils.RandomGenerator(20),
		Department:       department,
		PortalCards:      portalCards,
		PermissionGroups: permissionGroupIDs,
		CreatedAt:        time.Now(),
		LastModified:     time.Now(),
	}
	if err := s.departmentRepo.CreateDepartment(ctx, dept); err != nil {
		return fmt.Errorf("DEPARTMENT_CREATION_FAILED")
	}

	return nil
}

// UpdateDepartment updates only the provided fields in the department
func (s *ServiceImpl) UpdateDepartment(ctx context.Context, id string, updateData map[string]interface{}) (*entities.Department, error) {
	// Validate permission groups if provided
	if permissionGroups, exists := updateData["permission_groups"]; exists {
		if groupIDs, ok := permissionGroups.([]string); ok {
			for _, idStr := range groupIDs {
				if _, err := bson.ObjectIDFromHex(idStr); err != nil {
					return nil, fmt.Errorf("DEPARTMENT_INVALID_PERMISSION_GROUP_ID")
				}
			}
			if _, err := s.permissionService.ValidatePermissionGroups(ctx, groupIDs); err != nil {
				return nil, fmt.Errorf("DEPARTMENT_PERMISSION_GROUPS_VALIDATION_FAILED")
			}
		}
	}
	fmt.Println("we are on the authorize update 1st :", updateData)
	// Update only the provided fields
	data, err := s.departmentRepo.UpdateDepartment(ctx, id, updateData)
	if err != nil {
		return nil, fmt.Errorf("DEPARTMENT_UPDATE_FAILED")
	}
	return data, nil
}

func (s *ServiceImpl) GetAllDepartments(ctx context.Context, filterParams *constant.Filter) (*common_util.PaginatedResponse[[]*entities.Department], error) {
	return s.departmentRepo.GetAllDepartments(ctx, filterParams)
}

func (s *ServiceImpl) GetDepartmentByID(ctx context.Context, id string) (*entities.Department, error) {
	return s.departmentRepo.GetDepartmentByID(ctx, id)
}

// Authorize handles CPS action approvals for department operations

func (s *ServiceImpl) Authorize(ctx context.Context, action *cpsactions.CPSAction) (*cpsactions.CPSAction, error) {
	requestedAction := action.RequestAction
	var err error

	switch requestedAction {
	case cps_const.RequestCreateDepartment:
		// Extract department data from CurrentAction
		currentAction, ok := action.CurrentAction.(map[string]interface{})
		if !ok {
			s.logger.Errorf("invalid current action format for create")
			return nil, fmt.Errorf("DEPARTMENT_INVALID_ACTION_DATA_FORMAT")
		}

		s.logger.Infof("Processing create request with data: %+v", currentAction)

		// Extract fields from current action
		department, _ := currentAction["department"].(string)

		permissionGroupsRaw := currentAction["permission_groups"]

		s.logger.Infof("Extracted raw permission groups: %+v (type: %T)", permissionGroupsRaw, permissionGroupsRaw)
		var portalCardsList []string
		for _, item := range currentAction["portal_cards"].([]interface{}) {
			if str, ok := item.(string); ok {
				portalCardsList = append(portalCardsList, str)
			}
		}
		// Convert permission groups to proper []string
		var permissionGroups []string
		if permissionGroupsRaw != nil {
			switch v := permissionGroupsRaw.(type) {
			case []string:
				permissionGroups = v
			case []interface{}:
				// Convert []interface{} to []string
				for _, item := range v {
					if str, ok := item.(string); ok {
						permissionGroups = append(permissionGroups, str)
					} else {
						s.logger.Errorf("invalid permission group item type: %T, value: %+v", item, item)
						return nil, fmt.Errorf("DEPARTMENT_INVALID_PERMISSION_GROUP_FORMAT")
					}
				}
			default:
				s.logger.Errorf("unexpected permission groups type: %T, value: %+v", permissionGroupsRaw, permissionGroupsRaw)
				return nil, fmt.Errorf("DEPARTMENT_INVALID_PERMISSION_GROUPS_FORMAT")
			}
		}

		s.logger.Infof("Converted permission groups: %+v", permissionGroups)

		// Validate that we have permission groups
		if len(permissionGroups) == 0 {
			s.logger.Errorf("permission groups cannot be empty")
			return nil, fmt.Errorf("DEPARTMENT_PERMISSION_GROUPS_REQUIRED")
		}

		// Validate permission groups
		if _, err := s.permissionService.ValidatePermissionGroups(ctx, permissionGroups); err != nil {
			s.logger.Errorf("failed to validate permission groups: %v", err)
			return nil, fmt.Errorf("DEPARTMENT_PERMISSION_GROUPS_VALIDATION_FAILED")
		}

		// Create department
		err = s.CreateDepartment(ctx, department, portalCardsList, permissionGroups, action)
		if err != nil {
			s.logger.Errorf("failed to create department: %v", err)
			return nil, fmt.Errorf("DEPARTMENT_CREATION_FAILED")
		}

	case cps_const.RequestUpdateDepartment:
		// Extract update data from CurrentAction
		updateData, ok := action.CurrentAction.(map[string]interface{})
		if !ok {
			s.logger.Errorf("invalid current action format for update")
			return nil, fmt.Errorf("DEPARTMENT_INVALID_ACTION_DATA_FORMAT")
		}
		fmt.Println("we are on the authorize update 1st :", updateData)
		s.logger.Infof("Processing update request with data: %+v", updateData)

		// Get department ID from the update data
		departmentID, ok := updateData["department_id"].(string)
		if !ok || departmentID == "" {
			s.logger.Errorf("department_id not found in update data")
			return nil, fmt.Errorf("DEPARTMENT_ID_MISSING_IN_UPDATE_DATA")
		}

		s.logger.Infof("Updating department with ID: %s", departmentID)

		// Create a clean update data map with only the fields to update
		cleanUpdateData := make(map[string]interface{})

		// Copy only the fields that should be updated (excluding department_id)
		for key, value := range updateData {
			if key != "department_id" && value != nil {
				// Handle different field types
				switch key {
				case "department":
					if strVal, ok := value.(string); ok && strVal != "" {
						cleanUpdateData[key] = strVal
					}
				case "portal_cards":
					if arrVal, ok := value.([]string); ok && len(arrVal) > 0 {
						cleanUpdateData[key] = arrVal
					}
				case "permission_groups":
					if arrVal, ok := value.([]string); ok && len(arrVal) > 0 {
						cleanUpdateData[key] = arrVal
					}
				case "enable":
					if boolVal, ok := value.(bool); ok {
						fmt.Println("Enabling department:", boolVal)
						cleanUpdateData["enabled"] = boolVal // Map to the correct field name
					}
				case "enabled":
					if boolVal, ok := value.(bool); ok {
						cleanUpdateData["enabled"] = boolVal
					}
				}
			}
		}
		fmt.Println("we are on the authorize update 2nd--cleaned one  :", cleanUpdateData)
		s.logger.Infof("Clean update data: %+v", cleanUpdateData)

		// Update department with the cleaned data
		_, err = s.UpdateDepartment(ctx, departmentID, cleanUpdateData)
		if err != nil {
			s.logger.Errorf("failed to update department: %v", err)
			return nil, fmt.Errorf("DEPARTMENT_UPDATE_FAILED")
		}

		s.logger.Infof("Department updated successfully: %s", departmentID)

	default:
		return nil, fmt.Errorf("DEPARTMENT_UNSUPPORTED_ACTION_TYPE")
	}

	return action, nil
}
