package cpsactions

import (
	"context"

	entities "github.com/CBE-Super-App/cbe-super-app-cps-action/internal/domain/cps_actions/entities"
	common_util "github.com/CBE-Super-App/cbe-super-app-cps-action/pkgs/utils"
	constant "github.com/CBE-Super-App/cbe-super-app-cps-action/utils"
)

type CPSActionRepository interface {
	CreateCPSAction(ctx context.Context, action *entities.CPSAction) (*entities.CPSAction, error)
	UpdateCPSAction(ctx context.Context, action *entities.CPSAction) (*entities.CPSAction, error)
	CPSActionExists(ctx context.Context, user entities.CheckCPSAction) (bool, error)
	ApproveCPSAction(ctx context.Context, action *entities.AuthorizeCPSAction) (*entities.CPSAction, error)
	RejectCPSAction(ctx context.Context, action *entities.AuthorizeCPSAction) (*entities.CPSAction, error)
	GetCPSActionsByDepartment(ctx context.Context, department string, status string, filterParams *constant.Filter) (*common_util.PaginatedResponse[[]*entities.CPSAction], error)
	GetCPSActionByID(ctx context.Context, id string) (*entities.CPSAction, error)
	GetCPSActionByActionCode(ctx context.Context, uniqueID string) (*entities.CPSAction, error)
}
