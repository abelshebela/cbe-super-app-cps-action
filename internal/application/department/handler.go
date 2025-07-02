package department

import (
	domain "github.com/CBE-Super-App/cbe-super-app-cps-action/internal/domain/department"
	"github.com/CBE-Super-App/cbe-super-app-cps-action/internal/domain/department/entities"

	"gitlab.com/bersufekadgetachew/cbe-super-app-shared/shared/utils"
)

type DepartmentService interface {
	CreateDepartment(input string, portalCards []string, cpsAction entities.CPSAction) error
	UpdateDepartment(code string, req CreateDepartmentRequest) error
	ValidateActionRequest(actionCode string, userDept string) (*entities.CPSAction, error)
	ApproveActionRequest(actionCode string, action entities.CPSAction) error
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
func (h *DepartmentHandler) CreateDepartment(input string, portalCards []string, cpsAction entities.CPSAction) error {
	err := h.service.CreateDepartment(input, portalCards, cpsAction)
	if err != nil {
		return err
	}
	return nil
}
func (h *DepartmentHandler) UpdateDepartment(code string, req CreateDepartmentRequest) error {
	err := h.service.UpdateDepartment(code, req.Department, req.PortalCards)
	if err != nil {
		return err
	}
	return nil
}

func (h *DepartmentHandler) ValidateActionRequest(actionCode string, userDept string) (*entities.CPSAction, error) {
	action, err := h.service.ValidateActionRequest(actionCode, userDept)
	if err != nil {
		h.logger.Errorf("Failed to validate action request: %v", err)
		return nil, err
	}
	return action, nil
}

func (h *DepartmentHandler) ApproveActionRequest(actionCode string, action entities.CPSAction) error {
	err := h.service.ApproveActionRequest(actionCode, action)
	if err != nil {
		h.logger.Errorf("Failed to approve action request: %v", err)
		return err
	}
	return nil
}
