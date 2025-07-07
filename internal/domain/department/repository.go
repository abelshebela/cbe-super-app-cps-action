package department

import (
	"context"
	"github.com/CBE-Super-App/cbe-super-app-cps-action/internal/domain/department/entities"
)

type CPSActionRepository interface {
	CheckRequestExists(ctx context.Context, action entities.CPSAction) (bool, error)
	CreateCPSAction(ctx context.Context, department string, portalCards []string, action entities.CPSAction) error
	ValidateActionRequest(ctx context.Context, actionCode string, userDept string) (*entities.CPSAction, error)
	FindByActionCode(ctx context.Context, code string) (*entities.CPSAction, error)
	UpdateActionStatus(ctx context.Context, actionCode string, status string) error
	ApproveActionRequest(ctx context.Context, actionCode string, action entities.CPSAction) error
}

type DepartmentRepository interface {
	CheckDepartmentExists(ctx context.Context, department string) (bool, error)
	CreateDepartment(ctx context.Context, dept entities.Department) error
	UpdateDepartment(ctx context.Context, code string, department string, portalCards []string) error
}
