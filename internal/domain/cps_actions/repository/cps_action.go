package repository

import (
	"context"

	entities "github.com/CBE-Super-App/cbe-super-app-cps-action/internal/domain/cps_actions/entities"
	common_util "github.com/CBE-Super-App/cbe-super-app-cps-action/pkgs/utils"
)

type CPSActionRepository interface {
	CreateCPSAction(ctx context.Context, action *entities.CPSAction) (*entities.CPSAction, error)
	UpdateCPSAction(ctx context.Context, action *entities.CPSAction) (*entities.CPSAction, error)
	CPSActionExists(ctx context.Context, uniqueID string) (bool, error)
	ApproveCPSAction(ctx context.Context, action *entities.CPSAction) (*entities.CPSAction, error)
	RejectCPSAction(ctx context.Context, action *entities.CPSAction) (*entities.CPSAction, error)
	GetCPSActionsByDepartment(ctx context.Context, department string) (*common_util.PaginatedResponse[[]*entities.CPSAction], error)
	GetCPSActionByID(ctx context.Context, id string) (*entities.CPSAction, error)
	GetCPSActionByUniqueID(ctx context.Context, uniqueID string) (*entities.CPSAction, error)
}
