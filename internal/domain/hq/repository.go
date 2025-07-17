package hq

import (
	"context"

	constant "github.com/CBE-Super-App/cbe-super-app-cps-action/utils"
)

type Repository interface {
	GetHQByID(ctx context.Context, id string) (HQ, error)
	GetAllHQ(ctx context.Context, filerParams *constant.Filter) (*HQRespose, error)
	UpdateHQ(ctx context.Context, id string, update HQ) error
}
