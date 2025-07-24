package repository

import (
	"context"
	entities "github.com/CBE-Super-App/cbe-super-app-cps-action/internal/domain/cps_actions/entities"

)

type FaydaRepository interface {
	AuthorizeFaydaAccountEnableDisable(ctx context.Context, req *entities.CPSAction) (*entities.CPSAction, error)
}
