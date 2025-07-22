// Package department provides handlers and services for managing departments and their actions.
package department

import (
	"context"
	"encoding/json"
	"fmt"

	"github.com/CBE-Super-App/cbe-super-app-cps-action/internal/domain/department"
	domain "github.com/CBE-Super-App/cbe-super-app-cps-action/internal/domain/department"

	cpsactions "github.com/CBE-Super-App/cbe-super-app-cps-action/internal/domain/cps_actions/entities"
	"github.com/CBE-Super-App/cbe-super-app-cps-action/internal/domain/department/entities"

	common_util "github.com/CBE-Super-App/cbe-super-app-cps-action/pkgs/utils"
	err_msg "github.com/CBE-Super-App/cbe-super-app-cps-action/pkgs/utils"
	constant "github.com/CBE-Super-App/cbe-super-app-cps-action/utils"

	"gitlab.com/bersufekadgetachew/cbe-super-app-shared/shared/utils"
)

type DepartmentData struct {
	Code             string
	Name             string
	PortalCards      []string
	PermissionGroups []string
}

type DepartmentService interface {
	CreateCPSAction(ctx context.Context, department string, portalCards []string, permissionGroups []string, cpsAction cpsactions.CPSAction) (string, error)
	CreateDepartment(ctx context.Context, input string, portalCards []string, cpsAction cpsactions.CPSAction) error
	UpdateDepartment(ctx context.Context, code string, req UpdateDepartmentRequest) (*entities.Department, error)
	CreateDepartmentUpdateCPSAction(ctx context.Context, req cpsactions.CPSAction) (*cpsactions.CPSAction, error)
	ValidateActionRequest(ctx context.Context, actionCode string, userDept string) (*cpsactions.CPSAction, error)
	RejectActionRequest(ctx context.Context, actionCode string, action cpsactions.CPSAction) error
	GetAllDepartments(ctx context.Context, filterParams *constant.Filter) (*common_util.PaginatedResponse[[]*entities.Department], error)
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

	// Extract permission_groups
	var permissionGroups []string
	if rawGroups, ok := data["permission_groups"]; ok {
		switch v := rawGroups.(type) {
		case []any:
			for _, g := range v {
				if str, ok := g.(string); ok {
					permissionGroups = append(permissionGroups, str)
				}
			}
		case []string:
			permissionGroups = v
		}
	}

	code, _ := data["department_code"].(string)

	return &DepartmentData{
		Code:             code,
		Name:             name,
		PortalCards:      cards,
		PermissionGroups: permissionGroups,
	}, nil
}

func (h *DepartmentHandler) CreateCPSAction(ctx context.Context, department string, portalCards []string, permissionGroups []string, cpsAction cpsactions.CPSAction) (string, error) {
	// Check request exists
	if existing, err := h.service.CheckRequestExists(ctx, cpsAction); err != nil {
		h.logger.Errorf("failed to check request exists: %v", err)
		return "", err
	} else if existing != nil {
		h.logger.Errorf("pending request exists for action code: %s", existing.ActionCode)
		return "", fmt.Errorf(err_msg.PendingRequestExists)
	}

	// check department exist
	if exists, err := h.service.CheckDepartmentExists(ctx, department); err != nil {
		h.logger.Errorf("failed to check department exists: %v", err)
		return "", err
	} else if exists {
		h.logger.Errorf("department already exists for action code: %s", department)
		return "", fmt.Errorf(err_msg.DepartmentAlreadyExists)
	}

	// creating dep
	createdActionCode, err := h.service.CreateCPSAction(ctx, department, portalCards, permissionGroups, cpsAction)
	if err != nil {
		h.logger.Errorf("failed to create CPS action: %v", err)
		return "", err
	}
	return createdActionCode, nil
}
func (h *DepartmentHandler) CreateDepartment(ctx context.Context, department string, portalCards []string, cpsAction cpsactions.CPSAction) error {
	h.logger.Infof("Creating department: %s with portal cards: %v", department, portalCards)
	err := h.service.CreateDepartment(ctx, department, portalCards, cpsAction)
	if err != nil {
		h.logger.Errorf("Failed to create department: %v", err)
		return err
	}
	h.logger.Infof("Department created successfully: %s", department)
	return nil
}

func (h *DepartmentHandler) UpdateDepartment(ctx context.Context, id string, req UpdateDepartmentRequest) (*entities.Department, error) {
	
	if id != "" {
		if exists, err := h.service.CheckDepartmentExistsByID(ctx, id); err != nil {
			return nil, err
		} else if !exists {
			return nil, fmt.Errorf(err_msg.DepartmentNotFound)
		}
	}
	
	domainReq := domain.UpdateDepartmentRequest{
		Department:       req.Department,
		PortalCards:      req.PortalCards,
		PermissionGroups: req.PermissionGroups,
	}
	data, err := h.service.UpdateDepartment(ctx, id, domainReq)
	if err != nil {
		return nil, err
	}
	return data, nil
}

func (h *DepartmentHandler) ValidateActionRequest(ctx context.Context, actionCode string, userDept string) (*cpsactions.CPSAction, error) {
	action, err := h.service.ValidateActionRequest(ctx, actionCode, userDept)
	if err != nil {
		h.logger.Errorf("Failed to validate action request: %v", err)
		return nil, err
	}
	return action, nil
}

func (h *DepartmentHandler) CreateDepartmentUpdateCPSAction(ctx context.Context, req cpsactions.CPSAction) (*cpsactions.CPSAction, error) {
	// Check for pending update action
	if existing, err := h.service.CheckRequestExists(ctx, req); err != nil {
		h.logger.Errorf("failed to check request exists: %v", err)
		return nil, err
	} else if existing != nil {
		h.logger.Errorf("pending request exists for action code: %s", existing.ActionCode)
		return nil, fmt.Errorf(err_msg.PendingRequestExists)
	}

	_, ok := req.CurrentAction.(CreateDepartmentRequest)
	if !ok {
		// Try to parse from map if type assertion fails
		if m, ok := req.CurrentAction.(map[string]interface{}); ok {
			// Convert map to CreateDepartmentRequest
			jsonData, err := json.Marshal(m)
			if err != nil {
				h.logger.Errorf("failed to marshal department request: %v", err)
				return nil, fmt.Errorf("invalid department request format")
			}

			var createReq CreateDepartmentRequest
			if err := json.Unmarshal(jsonData, &createReq); err != nil {
				h.logger.Errorf("failed to unmarshal department request: %v", err)
				return nil, fmt.Errorf("invalid department request format")
			}

			// Now use createReq
			if exists, err := h.service.CheckDepartmentExists(ctx, createReq.Department); err != nil {
				h.logger.Errorf("failed to check department exists: %v", err)
				return nil, err
			} else if !exists {
				h.logger.Errorf("department not found for action code: %s", createReq.Department)
				return nil, fmt.Errorf(err_msg.DepartmentNotFound)
			}
		} else {
			return nil, fmt.Errorf("invalid request type for current action")
		}
	}

	cpsAction, err := h.service.CreateDepartmentUpdateCPSAction(ctx, req)
	if err != nil {
		h.logger.Errorf("Failed to create department update CPS action: %v", err)
		return nil, err
	}

	h.logger.Infof("Department update CPS action created successfully with action code: %s", cpsAction.ActionCode)
	return cpsAction, nil
}

func (h *DepartmentHandler) RejectActionRequest(ctx context.Context, actionCode string, action cpsactions.CPSAction) error {
	return h.service.RejectActionRequest(ctx, actionCode, action)
}

func (h *DepartmentHandler) GetAllDepartments(ctx context.Context, filterParams *constant.Filter) (*common_util.PaginatedResponse[[]*entities.Department], error) {
	return h.service.GetAllDepartments(ctx, filterParams)
}
