package department

import (
	"context"

	"github.com/CBE-Super-App/cbe-super-app-cps-action/internal/domain/department/entities"
	common_util "github.com/CBE-Super-App/cbe-super-app-cps-action/pkgs/utils"
	constant "github.com/CBE-Super-App/cbe-super-app-cps-action/utils"
)

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
	GetAllDepartments(ctx context.Context, filterParams *constant.Filter) (*common_util.PaginatedResponse[[]*entities.Department], error)
}
