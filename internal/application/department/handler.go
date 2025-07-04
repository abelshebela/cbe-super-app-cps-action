// Package department provides handlers and services for managing departments and their actions.
package department

import (
	"context"
	"fmt"

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
	CreateDepartment(ctx context.Context, input string, portalCards []string, cpsAction entities.CPSAction) error
	UpdateDepartment(ctx context.Context, code string, req CreateDepartmentRequest) error
	ValidateActionRequest(ctx context.Context, actionCode string, userDept string) (*entities.CPSAction, error)
	ApproveActionRequest(ctx context.Context, actionCode string, action entities.CPSAction) error
	ApproveCreateAction(ctx context.Context, action *entities.CPSAction) error
	ApproveUpdateAction(ctx context.Context, action *entities.CPSAction) error
	ApproveActionByType(ctx context.Context, cpsAction entities.CPSAction) error
}

type DepartmentHandler struct {
	service *domain.Service
	logger  utils.Logger
}

func InitDepartmentHandler(service *domain.Service, logger utils.Logger) DepartmentService {
	return &DepartmentHandler{
		service: service,
		logger:  logger,
	}
}

func (h *DepartmentHandler) extractDepartmentData(current any) (*DepartmentData, error) {
	var data map[string]any

	switch v := current.(type) {
	case map[string]any:
		data = v
	default:
		h.logger.Errorf("unsupported data type: %T", current)
		return nil, fmt.Errorf(err_msg.InvalidInput)
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

func (h *DepartmentHandler) CreateDepartment(ctx context.Context, department string, portalCards []string, cpsAction entities.CPSAction) error {
	// Check request exists
	if ok, err := h.service.CheckRequestExists(ctx, cpsAction); err != nil {
		return err
	} else if ok {
		return fmt.Errorf(err_msg.PendingRequestExists)
	}

	// check department exist
	if exists, err := h.service.CheckDepartmentExists(ctx, department); err != nil {
		return err
	} else if exists {
		return fmt.Errorf(err_msg.DepartmentAlreadyExists)
	}

	// createing dep
	err := h.service.CreateDepartment(ctx, department, portalCards, cpsAction)
	if err != nil {
		return err
	}
	return nil
}
func (h *DepartmentHandler) UpdateDepartment(ctx context.Context, code string, req CreateDepartmentRequest) error {
	if exists, err := h.service.CheckDepartmentExists(ctx, req.Department); err != nil {
		return err
	} else if !exists {
		return fmt.Errorf(err_msg.DepartmentNotFound)
	}

	err := h.service.UpdateDepartment(ctx, code, req.Department, req.PortalCards)
	if err != nil {
		return err
	}
	return nil
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
	return h.UpdateDepartment(ctx, data.Code, req)
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
