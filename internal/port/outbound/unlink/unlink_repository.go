package unlink

import (
	"context"

	cps_entities "github.com/CBE-Super-App/cbe-super-app-cps-action/internal/domain/cps_actions/entities"
)

type UnlinkRepository interface {
	UnlinkDevice(ctx context.Context, userCode string, cpsAction cps_entities.CPSAction) (string, error)
	Authorize(ctx context.Context, cpsAction *cps_entities.CPSAction) (*cps_entities.CPSAction, error)
}
