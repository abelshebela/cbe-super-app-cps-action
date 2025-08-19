package department

import (
	"context"

	"github.com/CBE-Super-App/cbe-super-app-cps-action/internal/domain/department/entities"
	common_util "github.com/CBE-Super-App/cbe-super-app-cps-action/pkgs/utils"
	constant "github.com/CBE-Super-App/cbe-super-app-cps-action/utils"
)

type DepartmentRepository interface {
	CheckDepartmentExists(ctx context.Context, department string) (bool, error)
	CheckDepartmentExistsByID(ctx context.Context, id string) (bool, error)
	CreateDepartment(ctx context.Context, dept entities.Department) error
	UpdateDepartment(ctx context.Context, id string, updateData map[string]interface{}) (*entities.Department, error)
	GetAllDepartments(ctx context.Context, filterParams *constant.Filter) (*common_util.PaginatedResponse[[]*entities.Department], error)
	GetDepartmentByID(ctx context.Context, id string) (*entities.Department, error)
}
