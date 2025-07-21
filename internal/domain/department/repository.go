package department

import (
	"context"

	cpsactions "github.com/CBE-Super-App/cbe-super-app-cps-action/internal/domain/cps_actions/entities"
	"github.com/CBE-Super-App/cbe-super-app-cps-action/internal/domain/department/entities"
	common_util "github.com/CBE-Super-App/cbe-super-app-cps-action/pkgs/utils"
	constant "github.com/CBE-Super-App/cbe-super-app-cps-action/utils"
)

type CPSActionRepository interface {
	CheckRequestExists(ctx context.Context, action cpsactions.CPSAction) (*cpsactions.CPSAction, error)
	CreateCPSAction(ctx context.Context, department string, portalCards []string, permissionGroups []string, action cpsactions.CPSAction) error
	ValidateActionRequest(ctx context.Context, actionCode string, userDept string) (*cpsactions.CPSAction, error)
	FindByActionCode(ctx context.Context, code string) (*cpsactions.CPSAction, error)
	UpdateActionStatus(ctx context.Context, actionCode string, status string) (*cpsactions.CPSAction, error)
	ApproveActionRequest(ctx context.Context, actionCode string, action cpsactions.CPSAction) error
	RejectActionRequest(ctx context.Context, actionCode string, action cpsactions.CPSAction) error
	CreateDepartmentUpdateCPSAction(ctx context.Context, req cpsactions.CPSAction) (*cpsactions.CPSAction, error)
	ApproveDepartmentUpdate(ctx context.Context, cpsAction cpsactions.CPSAction) error
	RejectDepartmentUpdate(ctx context.Context, cpsAction cpsactions.CPSAction) error
}

type DepartmentRepository interface {
	CheckDepartmentExists(ctx context.Context, department string) (bool, error)
	CheckDepartmentExistsByID(ctx context.Context, id string) (bool, error)
	CreateDepartment(ctx context.Context, dept entities.Department) error
	UpdateDepartment(ctx context.Context, id string, department string, portalCards []string, permissionGroups []string) (*entities.Department, error)
	GetAllDepartments(ctx context.Context, filterParams *constant.Filter) (*common_util.PaginatedResponse[[]*entities.Department], error)
}
