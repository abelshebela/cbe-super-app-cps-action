// Package department provides handlers and services for managing departments and their actions.
package department

import (
	"context"
	"encoding/json"
	"fmt"

	"github.com/CBE-Super-App/cbe-super-app-cps-action/internal/domain/department"
	domain "github.com/CBE-Super-App/cbe-super-app-cps-action/internal/domain/department"

	"github.com/CBE-Super-App/cbe-super-app-cps-action/internal/domain/department/entities"

	err_msg "github.com/CBE-Super-App/cbe-super-app-cps-action/pkgs/utils"

	"gitlab.com/bersufekadgetachew/cbe-super-app-shared/shared/utils"
)

type DepartmentData struct {
	Code        string
	Name        string
	PortalCards []string
}

type DepartmentService interface {
	CreateCPSAction(ctx context.Context, department string, portalCards []string, cpsAction entities.CPSAction) (*entities.CPSAction, error)
	CreateDepartment(ctx context.Context, input string, portalCards []string, cpsAction entities.CPSAction) error
	UpdateDepartment(ctx context.Context, code string, req CreateDepartmentRequest) (*entities.Department, error)
	CreateDepartmentUpdateCPSAction(ctx context.Context, req entities.CPSAction) (*entities.CPSAction, error)
	ValidateActionRequest(ctx context.Context, actionCode string, userDept string) (*entities.CPSAction, error)
	ApproveActionRequest(ctx context.Context, actionCode string, action entities.CPSAction) error
	ApproveCreateAction(ctx context.Context, action *entities.CPSAction) error
	ApproveUpdateAction(ctx context.Context, action *entities.CPSAction) error
	ApproveActionByType(ctx context.Context, cpsAction entities.CPSAction) error
	ApproveDepartmentUpdate(ctx context.Context, cpsAction entities.CPSAction) error
	RejectDepartmentUpdate(ctx context.Context, cpsAction entities.CPSAction) error
	RejectActionRequest(ctx context.Context, actionCode string, action entities.CPSAction) error
	GetAllDepartments(ctx context.Context) ([]entities.Department, error)
}

type DepartmentHandler struct {
	service          *domain.Service
	logger           utils.Logger
	departmentDomain department.Service
}

func InitDepartmentHandler(service *domain.Service, logger utils.Logger) DepartmentService {
	return &DepartmentHandler{
		service:          service,
		logger:           logger,
		departmentDomain: *service,
	}
}

func (h *DepartmentHandler) extractDepartmentData(current any) (*DepartmentData, error) {
	var data map[string]any

	switch v := current.(type) {
	case map[string]any:
		data = v
	default:
		bytes, err := json.Marshal(current)
		if err != nil {
			h.logger.Errorf("failed to marshal action data: %v", err)
			return nil, fmt.Errorf(err_msg.InvalidInput)
		}
		if err := json.Unmarshal(bytes, &data); err != nil {
			h.logger.Errorf("failed to unmarshal action data: %v", err)
			return nil, fmt.Errorf(err_msg.InvalidJSONPayload)
		}
	}

	name, ok := data["department"].(string)
	if !ok || name == "" {
		h.logger.Errorf("missing or invalid department name")
		return nil, fmt.Errorf(err_msg.DepartmentNameRequired)
	}

	rawCards, ok := data["portal_cards"].([]any)
	if !ok {
		h.logger.Errorf("missing or invalid portal_cards")
		return nil, fmt.Errorf(err_msg.InvalidInput)
	}

	var cards []string
	for _, card := range rawCards {
		str, ok := card.(string)
		if !ok {
			h.logger.Errorf("portal_cards must contain only strings, got: %v", card)
			return nil, fmt.Errorf(err_msg.PortalCardsInvalid)
		}
		cards = append(cards, str)
	}

	code, _ := data["department_code"].(string)

	return &DepartmentData{
		Code:        code,
		Name:        name,
		PortalCards: cards,
	}, nil
}

func (h *DepartmentHandler) CreateCPSAction(ctx context.Context, department string, portalCards []string, cpsAction entities.CPSAction) (*entities.CPSAction, error) {
	// Check request exists
	if existing, err := h.service.CheckRequestExists(ctx, cpsAction); err != nil {
		return nil, err
	} else if existing != nil {
		return nil, fmt.Errorf(err_msg.PendingRequestExists)
	}

	// check department exist
	if exists, err := h.service.CheckDepartmentExists(ctx, department); err != nil {
		return nil, err
	} else if exists {
		return nil, fmt.Errorf(err_msg.DepartmentAlreadyExists)
	}

	// creating dep
	createdAction, err := h.service.CreateCPSAction(ctx, department, portalCards, cpsAction)
	if err != nil {
		return nil, err
	}
	return createdAction, nil
}
func (h *DepartmentHandler) CreateDepartment(ctx context.Context, department string, portalCards []string, cpsAction entities.CPSAction) error {
	h.logger.Infof("Creating department: %s with portal cards: %v", department, portalCards)
	err := h.service.CreateDepartment(ctx, department, portalCards, cpsAction)
	if err != nil {
		h.logger.Errorf("Failed to create department: %v", err)
		return err
	}
	h.logger.Infof("Department created successfully: %s", department)
	return nil
}

func (h *DepartmentHandler) UpdateDepartment(ctx context.Context, code string, req CreateDepartmentRequest) (*entities.Department, error) {
	if exists, err := h.service.CheckDepartmentExists(ctx, req.Department); err != nil {
		return nil, err
	} else if !exists {
		return nil, fmt.Errorf(err_msg.DepartmentNotFound)
	}

	data, err := h.service.UpdateDepartment(ctx, code, req.Department, req.PortalCards)
	if err != nil {
		return nil, err
	}
	return data, nil
}

func (h *DepartmentHandler) ValidateActionRequest(ctx context.Context, actionCode string, userDept string) (*entities.CPSAction, error) {
	action, err := h.service.ValidateActionRequest(ctx, actionCode, userDept)
	if err != nil {
		h.logger.Errorf("Failed to validate action request: %v", err)
		return nil, err
	}
	return action, nil
}

func (h *DepartmentHandler) ApproveActionRequest(ctx context.Context, actionCode string, action entities.CPSAction) error {
	err := h.service.ApproveActionRequest(ctx, actionCode, action)
	if err != nil {
		h.logger.Errorf("Failed to approve action request: %v", err)
		return err
	}
	return nil
}

func (h *DepartmentHandler) ApproveCreateAction(ctx context.Context, action *entities.CPSAction) error {
	data, err := h.extractDepartmentData(action.CurrentAction)
	if err != nil {
		return err
	}

	return h.CreateDepartment(ctx, data.Name, data.PortalCards, *action)
}

func (h *DepartmentHandler) ApproveUpdateAction(ctx context.Context, action *entities.CPSAction) error {
	data, err := h.extractDepartmentData(action.CurrentAction)
	if err != nil {
		return err
	}
	if data.Code == "" {
		h.logger.Errorf("missing department_code for update")
		return fmt.Errorf(err_msg.DepartmentCodeRequired)
	}
	req := CreateDepartmentRequest{
		Department:  data.Name,
		PortalCards: data.PortalCards,
	}

	res := entities.CPSAction{
		CurrentAction: req,
	}

	err = h.departmentDomain.ApproveActionRequest(ctx, data.Code, res)
	if err != nil {
		return err
	}
	return nil
}

func (h *DepartmentHandler) RejectUpdateAction(ctx context.Context, action *entities.CPSAction) error {
	data, err := h.extractDepartmentData(action.CurrentAction)
	if err != nil {
		return err
	}
	if data.Code == "" {
		h.logger.Errorf("missing department_code for update")
		return fmt.Errorf(err_msg.DepartmentCodeRequired)
	}

	res := entities.CPSAction{
		CheckerID:          action.CheckerID,
		CheckerName:        action.CheckerName,
		CheckerPhoneNumber: action.CheckerPhoneNumber,
	}

	err = h.departmentDomain.RejectActionRequest(ctx, data.Code, res)
	if err != nil {
		return err
	}
	return nil
}

func (h *DepartmentHandler) ApproveActionByType(ctx context.Context, cpsAction entities.CPSAction) error {
	switch cpsAction.ActionType {
	case entities.ActionCreate:
		return h.ApproveCreateAction(ctx, &cpsAction)
	case entities.ActionUpdate:
		return h.ApproveUpdateAction(ctx, &cpsAction)
	default:
		return fmt.Errorf(err_msg.InvalidActionType)
	}
}

// CreateDepartmentUpdateCPSAction creates a CPS action for department updates
func (h *DepartmentHandler) CreateDepartmentUpdateCPSAction(ctx context.Context, req entities.CPSAction) (*entities.CPSAction, error) {
	// Check for pending update action
	if existing, err := h.service.CheckRequestExists(ctx, req); err != nil {
		return nil, err
	} else if existing != nil {
		return nil, fmt.Errorf(err_msg.PendingRequestExists)
	}

	// Check if department exists
	if exists, err := h.service.CheckDepartmentExists(ctx, req.Department); err != nil {
		return nil, err
	} else if !exists {
		return nil, fmt.Errorf(err_msg.DepartmentNotFound)
	}

	cpsAction, err := h.service.CreateDepartmentUpdateCPSAction(ctx, req)
	if err != nil {
		h.logger.Errorf("Failed to create department update CPS action: %v", err)
		return nil, err
	}

	h.logger.Infof("Department update CPS action created successfully with action code: %s", cpsAction.ActionCode)
	return cpsAction, nil
}

// ApproveDepartmentUpdate approves a department update CPS action
func (h *DepartmentHandler) ApproveDepartmentUpdate(ctx context.Context, cpsAction entities.CPSAction) error {
	err := h.service.ApproveDepartmentUpdate(ctx, cpsAction)
	if err != nil {
		h.logger.Errorf("Failed to approve department update: %v", err)
		return err
	}
	h.logger.Infof("Department update approved successfully for action code: %s", cpsAction.ActionCode)
	return nil
}

// RejectDepartmentUpdate rejects a department update CPS action
func (h *DepartmentHandler) RejectDepartmentUpdate(ctx context.Context, cpsAction entities.CPSAction) error {
	err := h.service.RejectDepartmentUpdate(ctx, cpsAction)
	if err != nil {
		h.logger.Errorf("Failed to reject department update: %v", err)
		return err
	}
	h.logger.Infof("Department update rejected successfully for action code: %s", cpsAction.ActionCode)
	return nil
}

func (h *DepartmentHandler) RejectActionRequest(ctx context.Context, actionCode string, action entities.CPSAction) error {
	return h.service.RejectActionRequest(ctx, actionCode, action)
}

func (h *DepartmentHandler) GetAllDepartments(ctx context.Context) ([]entities.Department, error) {
	return h.service.GetAllDepartments(ctx)
}
