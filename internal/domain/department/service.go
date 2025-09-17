// Package department provides services and business logic for managing departments and related CPS actions.
package department

import (
	"context"
	"fmt"
	"time"

	"github.com/CBE-Super-App/cbe-super-app-cps-action/internal/domain/department/entities"
	portal_card "github.com/CBE-Super-App/cbe-super-app-cps-action/internal/domain/portal_card"
	common_util "github.com/CBE-Super-App/cbe-super-app-cps-action/pkgs/utils"
	constant "github.com/CBE-Super-App/cbe-super-app-cps-action/utils"

	cps_const "github.com/CBE-Super-App/cbe-super-app-cps-action/internal/domain/cps_actions/constant"
	cpsactions "github.com/CBE-Super-App/cbe-super-app-cps-action/internal/domain/cps_actions/entities"
	"gitlab.com/bersufekadgetachew/cbe-super-app-shared/shared/utils"
)

const (
	CPSPrefix  = "CPS_"
	DeptPrefix = "DEP_"
)

type UpdateDepartmentRequest struct {
	Department  string   `json:"department"`
	PortalCards []string `json:"portal_cards"`
}

type Service interface {
	CheckDepartmentExists(ctx context.Context, department string) (bool, error)
	CheckDepartmentExistsByID(ctx context.Context, id string) (bool, error)
	CreateDepartment(ctx context.Context, department string, portalCards []string, cpsAction interface{}) error
	UpdateDepartment(ctx context.Context, id string, updateData map[string]interface{}) (*entities.Department, error)
	GetAllDepartments(ctx context.Context, filterParams *constant.Filter) (*common_util.PaginatedResponse[[]*entities.Department], error)
	GetDepartmentByID(ctx context.Context, id string) (*entities.Department, error)
	Authorize(ctx context.Context, action *cpsactions.CPSAction) (*cpsactions.CPSAction, error)
}

type ServiceImpl struct {
	departmentRepo DepartmentRepository
	portalService  portal_card.PortaCardInterface
	logger         utils.Logger
}

func InitDepartmentDomain(
	departmentRepo DepartmentRepository,
	portalService portal_card.PortaCardInterface,
	logger utils.Logger,
) Service {
	return &ServiceImpl{
		departmentRepo: departmentRepo,
		portalService:  portalService,
		logger:         logger,
	}
}

func (s *ServiceImpl) CheckDepartmentExists(ctx context.Context, department string) (bool, error) {
	return s.departmentRepo.CheckDepartmentExists(ctx, department)
}

func (s *ServiceImpl) CheckDepartmentExistsByID(ctx context.Context, id string) (bool, error) {
	return s.departmentRepo.CheckDepartmentExistsByID(ctx, id)
}

func (s *ServiceImpl) CreateDepartment(ctx context.Context, department string, portalCards []string, cpsAction interface{}) error {
	dept := entities.Department{
		DepartmentCode: utils.RandomGenerator(20),
		Department:     department,
		PortalCards:    portalCards,
		CreatedAt:      time.Now(),
		LastModified:   time.Now(),
	}
	if err := s.departmentRepo.CreateDepartment(ctx, dept); err != nil {
		return fmt.Errorf("DEPARTMENT_CREATION_FAILED")
	}

	return nil
}

func (s *ServiceImpl) UpdateDepartment(ctx context.Context, id string, updateData map[string]interface{}) (*entities.Department, error) {
	fmt.Println("we are on the authorize update 1st :", updateData)
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

func (s *ServiceImpl) Authorize(ctx context.Context, action *cpsactions.CPSAction) (*cpsactions.CPSAction, error) {
	requestedAction := action.RequestAction
	var err error

	switch requestedAction {
	case cps_const.RequestCreateDepartment:
		currentAction, ok := action.CurrentAction.(map[string]interface{})
		if !ok {
			s.logger.Errorf("invalid current action format for create")
			return nil, fmt.Errorf("DEPARTMENT_INVALID_ACTION_DATA_FORMAT")
		}

		s.logger.Infof("Processing create request with data: %+v", currentAction)

		department, _ := currentAction["department"].(string)

		var portalCardsList []string
		if portalCardsRaw, exists := currentAction["portal_cards"]; exists {
			if portalCardsArr, ok := portalCardsRaw.([]interface{}); ok {
				for _, item := range portalCardsArr {
					if str, ok := item.(string); ok {
						portalCardsList = append(portalCardsList, str)
					}
				}
			} else if portalCardsArr, ok := portalCardsRaw.([]string); ok {
				portalCardsList = portalCardsArr
			}
		}

		if len(portalCardsList) == 0 {
			s.logger.Errorf("portal cards cannot be empty")
			return nil, fmt.Errorf("DEPARTMENT_PORTAL_CARDS_REQUIRED")
		}

		err = s.CreateDepartment(ctx, department, portalCardsList, action)
		if err != nil {
			s.logger.Errorf("failed to create department: %v", err)
			return nil, fmt.Errorf("DEPARTMENT_CREATION_FAILED")
		}

	case cps_const.RequestUpdateDepartment:
		updateData, ok := action.CurrentAction.(map[string]interface{})
		if !ok {
			s.logger.Errorf("invalid current action format for update")
			return nil, fmt.Errorf("DEPARTMENT_INVALID_ACTION_DATA_FORMAT")
		}
		fmt.Println("we are on the authorize update 1st :", updateData)
		s.logger.Infof("Processing update request with data: %+v", updateData)

		departmentID, ok := updateData["department_id"].(string)
		if !ok || departmentID == "" {
			s.logger.Errorf("department_id not found in update data")
			return nil, fmt.Errorf("DEPARTMENT_ID_MISSING_IN_UPDATE_DATA")
		}

		s.logger.Infof("Updating department with ID: %s", departmentID)

		cleanUpdateData := make(map[string]interface{})

		for key, value := range updateData {
			if key != "department_id" && value != nil {
				switch key {
				case "department":
					if strVal, ok := value.(string); ok && strVal != "" {
						cleanUpdateData[key] = strVal
					}
				case "portal_cards":
					s.logger.Infof("Processing portal_cards field with value: %+v (type: %T)", value, value)
					if arrVal, ok := value.([]string); ok && len(arrVal) > 0 {
						s.logger.Infof("Portal cards as []string: %+v", arrVal)
						cleanUpdateData[key] = arrVal
					} else if arrVal, ok := value.([]interface{}); ok && len(arrVal) > 0 {
						s.logger.Infof("Portal cards as []interface{}: %+v", arrVal)
						var portalCards []string
						for _, item := range arrVal {
							if str, ok := item.(string); ok && str != "" {
								portalCards = append(portalCards, str)
							}
						}
						s.logger.Infof("Converted portal cards: %+v", portalCards)
						if len(portalCards) > 0 {
							cleanUpdateData[key] = portalCards
						}
					} else {
						s.logger.Warnf("Portal cards field has unexpected type or is empty: %+v (type: %T)", value, value)
					}
				case "enable":
					if boolVal, ok := value.(bool); ok {
						fmt.Println("Enabling department:", boolVal)
						cleanUpdateData["enabled"] = boolVal
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
