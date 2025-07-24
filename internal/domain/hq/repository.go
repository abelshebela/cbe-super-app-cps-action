package hq

import (
	"context"
	"time"

	cpsactions "github.com/CBE-Super-App/cbe-super-app-cps-action/internal/domain/cps_actions/entities"
	utils "github.com/CBE-Super-App/cbe-super-app-cps-action/pkgs/utils"
	constant "github.com/CBE-Super-App/cbe-super-app-cps-action/utils"
)

type Repository interface {
	GetHQByID(ctx context.Context, id string) (HQ, error)
	GetAllHQ(ctx context.Context, filerParams *constant.Filter) (*utils.PaginatedResponse[[]*HQ], error)
	UpdateHQField(ctx context.Context, field string, value interface{}, now time.Time) error
	GetSingleHQ(ctx context.Context) (HQ, error)
	FindPendingAction(ctx context.Context, requestAction string, department string) (*cpsactions.CPSAction, error)
}
