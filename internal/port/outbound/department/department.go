package department

import "github.com/CBE-Super-App/cbe-super-app-cps-action/internal/domain/department/entities"

type CPSActionRepository interface {
	CheckRequestExists(action entities.CPSAction) (bool, error)
	CreateCPSAction(department string, portalCards []string, action entities.CPSAction) error
	FindByActionCode(code string) (*entities.CPSAction, error)
	UpdateActionStatus(actionCode string, status string) error
}

type DepartmentRepository interface {
	CheckDepartmentExists(department string) (bool, error)
	CreateDepartment(dept entities.Department) error
	UpdateDepartment(code string, department string, portalCards []string) error
}
