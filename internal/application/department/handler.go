// Package department provides handlers and services for managing departments and their actions.
package department

import (
	"context"
	"fmt"

	cps_const "github.com/CBE-Super-App/cbe-super-app-cps-action/internal/domain/cps_actions/constant"
	cpsactions "github.com/CBE-Super-App/cbe-super-app-cps-action/internal/domain/cps_actions/entities"
	cps_service "github.com/CBE-Super-App/cbe-super-app-cps-action/internal/domain/cps_actions/services"
	"github.com/CBE-Super-App/cbe-super-app-cps-action/internal/domain/department"
	"github.com/CBE-Super-App/cbe-super-app-cps-action/internal/domain/department/entities"
	portal_card "github.com/CBE-Super-App/cbe-super-app-cps-action/internal/domain/portal_card"
	common_util "github.com/CBE-Super-App/cbe-super-app-cps-action/pkgs/utils"
	constant "github.com/CBE-Super-App/cbe-super-app-cps-action/utils"
	"gitlab.com/bersufekadgetachew/cbe-super-app-shared/shared/utils"
)

// Maker represents a user who can create CPS actions
type Maker struct {
	UserCode    string `json:"user_code"`
	FullName    string `json:"full_name"`
	PhoneNumber string `json:"phone_number"`
	Department  string `json:"department"`
}

type DepartmentService interface {
	CreateDepartment(ctx context.Context, request CreateDepartmentRequest, maker Maker) error
	UpdateDepartment(ctx context.Context, id string, request DepartmentUpdateCPSActionRequest, maker Maker) error
	EnableDepartment(ctx context.Context, id string, maker Maker) error
	DisableDepartment(ctx context.Context, id string, maker Maker) error
	GetAllDepartments(ctx context.Context, filterParams *constant.Filter) (*common_util.PaginatedResponse[[]*entities.Department], error)
	GetDepartmentByID(ctx context.Context, id string) (*entities.Department, error)
}

type DepartmentHandler struct {
	service           department.Service
	cpsService        cps_service.CPSActionService
	portalCardService portal_card.PortaCardInterface
	logger            utils.Logger
}

func NewDepartmentHandler(
	service department.Service,
	cpsService cps_service.CPSActionService,
	portalCardService portal_card.PortaCardInterface,
	logger utils.Logger,
) *DepartmentHandler {
	return &DepartmentHandler{
		service:           service,
		cpsService:        cpsService,
		portalCardService: portalCardService,
		logger:            logger,
	}
}

// handleCPSAction handles CPS action creation using the generic CPS service
func (h *DepartmentHandler) handleCPSAction(ctx context.Context, maker Maker, requestAction cps_const.RequestAction, curData, prevData interface{}, actionType cps_const.ActionType, uniqueID string) error {
	// Convert Maker to cpsactions.User
	cpsUser := cpsactions.User{
		UserCode:    maker.UserCode,
		FullName:    maker.FullName,
		PhoneNumber: maker.PhoneNumber,
		Department:  maker.Department,
	}

	cpsAction := h.cpsService.BuildCPSAction(ctx, cpsactions.CreateCPSRequest{
		User:          cpsUser,
		CurData:       curData,
		PrevData:      prevData,
		RequestAction: requestAction,
		ActionStatus:  cps_const.ActionPending,
		ActionType:    actionType,
	})

	// Set the unique ID if provided
	if uniqueID != "" {
		cpsAction.UniqueID = uniqueID
	}

	_, err := h.cpsService.CreateCPSAction(ctx, cpsAction)
	if err != nil {
		h.logger.Errorf("failed to create CPS action: %v", err)
		return err
	}
	return nil
}

func (h *DepartmentHandler) CreateDepartment(ctx context.Context, request CreateDepartmentRequest, maker Maker) error {
	h.logger.Infof("Creating department: %s with portal cards: %v", request.Department, request.PortalCards)

	// Validate department name is not empty
	if request.Department == "" {
		h.logger.Errorf("department name cannot be empty")
		return fmt.Errorf("DEPARTMENT_NAME_REQUIRED")
	}

	// Check if department already exists
	exists, err := h.service.CheckDepartmentExists(ctx, request.Department)
	if err != nil {
		h.logger.Errorf("failed to check department existence: %v", err)
		return fmt.Errorf("DEPARTMENT_EXISTENCE_CHECK_FAILED")
	}
	if exists {
		h.logger.Errorf("department already exists: %s", request.Department)
		return fmt.Errorf("DEPARTMENT_ALREADY_EXISTS")
	}

	// Validate portal cards are not empty
	if len(request.PortalCards) == 0 {
		h.logger.Errorf("portal cards cannot be empty")
		return fmt.Errorf("PORTAL_CARDS_REQUIRED")
	}



	if _, err := h.portalCardService.ValidatePortalCard(ctx, request.PortalCards); err != nil {
		h.logger.Errorf("invalid portal cards: %v", err)
		return err
	}
	// Create CPS request data
	cpsRequest := map[string]interface{}{
		"department":   request.Department,
		"portal_cards": request.PortalCards,
	}

	// Handle CPS action creation
	if err := h.handleCPSAction(ctx, maker, cps_const.RequestCreateDepartment, cpsRequest, nil, cps_const.ActionCreate, ""); err != nil {
		h.logger.Errorf("failed to create CPS action: %v", err)
		return fmt.Errorf("DEPARTMENT_CPS_ACTION_CREATION_FAILED")
	}

	h.logger.Infof("Department creation request sent successfully: %s", request.Department)
	return nil
}

func (h *DepartmentHandler) UpdateDepartment(ctx context.Context, id string, request DepartmentUpdateCPSActionRequest, maker Maker) error {
	h.logger.Infof("Updating department with ID: %s", id)

	// Validate ID format
	if id == "" {
		h.logger.Errorf("department ID cannot be empty")
		return fmt.Errorf("DEPARTMENT_ID_REQUIRED")
	}

	// Check if department exists
	exists, err := h.service.CheckDepartmentExistsByID(ctx, id)
	if err != nil {
		h.logger.Errorf("failed to check department existence: %v", err)
		return fmt.Errorf("DEPARTMENT_EXISTENCE_CHECK_FAILED")
	}
	if !exists {
		h.logger.Errorf("department not found: %s", id)
		return fmt.Errorf("DEPARTMENT_NOT_FOUND")
	}

	// Validate that at least one field is being updated
	if request.Department == "" && len(request.PortalCards) == 0 {
		h.logger.Errorf("no fields provided for update")
		return fmt.Errorf("DEPARTMENT_UPDATE_FIELDS_REQUIRED")
	}

	// Validate department name if provided
	if request.Department != "" {
		// Check if the new name conflicts with existing departments (excluding current one)
		existingDept, err := h.service.GetDepartmentByID(ctx, id)
		if err != nil {
			h.logger.Errorf("failed to get existing department: %v", err)
			return fmt.Errorf("DEPARTMENT_FETCH_FAILED")
		}

		if existingDept.Department != request.Department {
			// Check if new name already exists
			nameExists, err := h.service.CheckDepartmentExists(ctx, request.Department)
			if err != nil {
				h.logger.Errorf("failed to check department name conflict: %v", err)
				return fmt.Errorf("DEPARTMENT_NAME_CONFLICT_CHECK_FAILED")
			}
			if nameExists {
				h.logger.Errorf("department name already exists: %s", request.Department)
				return fmt.Errorf("DEPARTMENT_NAME_ALREADY_EXISTS")
			}
		}
	}

	// Validate portal cards if provided
	if len(request.PortalCards) > 0 {
		// Check if portal cards are not empty strings
		for i, card := range request.PortalCards {
			if card == "" {
				h.logger.Errorf("portal card at index %d cannot be empty", i)
				return fmt.Errorf("PORTAL_CARD_EMPTY_VALUE")
			}
		}
		if _, err := h.portalCardService.ValidatePortalCard(ctx, request.PortalCards); err != nil {
			h.logger.Errorf("invalid portal cards: %v", err)
			return err
		}
	}


	// Create CPS request data with only provided fields
	cpsRequest := map[string]interface{}{
		"department_id": id, // Include the department ID for identification
	}
	if request.Department != "" {
		cpsRequest["department"] = request.Department
	}
	if len(request.PortalCards) > 0 {
		cpsRequest["portal_cards"] = request.PortalCards
	}

	// Handle CPS action creation
	if err := h.handleCPSAction(ctx, maker, cps_const.RequestUpdateDepartment, cpsRequest, nil, cps_const.ActionUpdate, id); err != nil {
		h.logger.Errorf("failed to create CPS action: %v", err)
		return fmt.Errorf("DEPARTMENT_UPDATE_CPS_ACTION_CREATION_FAILED")
	}

	h.logger.Infof("Department update request sent successfully: %s", id)
	return nil
}

func (h *DepartmentHandler) EnableDepartment(ctx context.Context, id string, maker Maker) error {
	h.logger.Infof("Enabling department with ID: %s", id)

	// Validate ID format
	if id == "" {
		h.logger.Errorf("department ID cannot be empty")
		return fmt.Errorf("DEPARTMENT_ID_REQUIRED")
	}

	// Check if department exists
	exists, err := h.service.CheckDepartmentExistsByID(ctx, id)
	if err != nil {
		h.logger.Errorf("failed to check department existence: %v", err)
		return fmt.Errorf("DEPARTMENT_EXISTENCE_CHECK_FAILED")
	}
	if !exists {
		h.logger.Errorf("department not found: %s", id)
		return fmt.Errorf("DEPARTMENT_NOT_FOUND")
	}

	// Check if department is already enabled
	existingDept, err := h.service.GetDepartmentByID(ctx, id)
	if err != nil {
		h.logger.Errorf("failed to get existing department: %v", err)
		return fmt.Errorf("DEPARTMENT_FETCH_FAILED")
	}

	if existingDept.Enabled {
		h.logger.Errorf("department is already enabled: %s", id)
		return fmt.Errorf("DEPARTMENT_ALREADY_ENABLED")
	}

	// Create CPS request data for enable
	cpsRequest := map[string]interface{}{
		"department_id": id,   // Include the department ID for identification
		"enabled":       true, // Use "enabled" to match database field
	}

	// Handle CPS action creation
	if err := h.handleCPSAction(ctx, maker, cps_const.RequestUpdateDepartment, cpsRequest, nil, cps_const.ActionUpdate, id); err != nil {
		h.logger.Errorf("failed to create CPS action: %v", err)
		return fmt.Errorf("DEPARTMENT_ENABLE_CPS_ACTION_CREATION_FAILED")
	}

	h.logger.Infof("Department enable request sent successfully: %s", id)
	return nil
}

func (h *DepartmentHandler) DisableDepartment(ctx context.Context, id string, maker Maker) error {
	h.logger.Infof("Disabling department with ID: %s", id)

	// Validate ID format
	if id == "" {
		h.logger.Errorf("department ID cannot be empty")
		return fmt.Errorf("DEPARTMENT_ID_REQUIRED")
	}

	// Check if department exists
	exists, err := h.service.CheckDepartmentExistsByID(ctx, id)
	if err != nil {
		h.logger.Errorf("failed to check department existence: %v", err)
		return fmt.Errorf("DEPARTMENT_EXISTENCE_CHECK_FAILED")
	}
	if !exists {
		h.logger.Errorf("department not found: %s", id)
		return fmt.Errorf("DEPARTMENT_NOT_FOUND")
	}

	// Check if department is already disabled
	existingDept, err := h.service.GetDepartmentByID(ctx, id)
	if err != nil {
		h.logger.Errorf("failed to get existing department: %v", err)
		return fmt.Errorf("DEPARTMENT_FETCH_FAILED")
	}

	if !existingDept.Enabled {
		h.logger.Errorf("department is already disabled: %s", id)
		return fmt.Errorf("DEPARTMENT_ALREADY_DISABLED")
	}

	// Create CPS request data for disable
	cpsRequest := map[string]interface{}{
		"department_id": id,    // Include the department ID for identification
		"enabled":       false, // Use "enabled" to match database field
	}

	// Handle CPS action creation
	if err := h.handleCPSAction(ctx, maker, cps_const.RequestUpdateDepartment, cpsRequest, nil, cps_const.ActionUpdate, id); err != nil {
		h.logger.Errorf("failed to create CPS action: %v", err)
		return fmt.Errorf("DEPARTMENT_DISABLE_CPS_ACTION_CREATION_FAILED")
	}

	h.logger.Infof("Department disable request sent successfully: %s", id)
	return nil
}

func (h *DepartmentHandler) GetAllDepartments(ctx context.Context, filterParams *constant.Filter) (*common_util.PaginatedResponse[[]*entities.Department], error) {
	h.logger.Infof("Fetching all departments with filter params: %+v", filterParams)

	result, err := h.service.GetAllDepartments(ctx, filterParams)
	if err != nil {
		h.logger.Errorf("Failed to fetch departments: %v", err)
		return nil, err
	}

	h.logger.Infof("Successfully fetched %d departments", len(result.Data))
	return result, nil
}

func (h *DepartmentHandler) GetDepartmentByID(ctx context.Context, id string) (*entities.Department, error) {
	h.logger.Infof("Fetching department by ID: %s", id)

	result, err := h.service.GetDepartmentByID(ctx, id)
	if err != nil {
		h.logger.Errorf("Failed to fetch department by ID: %v", err)
		return nil, err
	}

	h.logger.Infof("Successfully fetched department: %s", id)
	return result, nil
}

// InitDepartmentHandler initializes the department application handler
func InitDepartmentHandler(
	departmentDomain department.Service,
	cpsActionDomain cps_service.CPSActionService,
	portalCardDomain portal_card.PortaCardInterface,
	logger utils.Logger,
) DepartmentService {

	return NewDepartmentHandler(departmentDomain, cpsActionDomain, portalCardDomain, logger)

}
