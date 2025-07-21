package unlink

import (
	"context"

	"github.com/CBE-Super-App/cbe-super-app-cps-action/internal/domain/unlink/entities"
)

type Repository interface {
	UnlinkDevice(ctx context.Context, userCode string, cpsAction entities.CPSAction) (string, error)
	ApproveOrDecline(ctx context.Context, userCode, decision, reason string, cpsAction entities.CPSAction) error
}
