package repository

import (
	"context"

	"github.com/CBE-Super-App/cbe-super-app-cps-action/internal/adapter/outbound/model"
	entities "github.com/CBE-Super-App/cbe-super-app-cps-action/internal/domain/cps_actions/entities"
)

type FaydaRepository interface {
	EnableFaydaAccount(ctx context.Context, req *entities.CPSAction, userID string) (string, error)
	DisableFaydaAccount(ctx context.Context, req *entities.CPSAction, userID string) (string, error)
	// AuthorizeFaydaAccountEnableDisable(ctx context.Context, req *entities.CPSAction) (*entities.CPSAction, error)
	AuthorizeBlockFaydaUser(ctx context.Context, req *entities.CPSAction) (*model.CPSAction, error)
	AuthorizeEnableFaydaUser(ctx context.Context, req *entities.CPSAction) (*model.CPSAction, error)
}
