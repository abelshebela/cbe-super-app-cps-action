package cpsaction

import (
	"context"
	"fmt"

	"github.com/CBE-Super-App/cbe-super-app-cps-action/internal/application"
	constants "github.com/CBE-Super-App/cbe-super-app-cps-action/internal/domain/cps_actions/constant"
	entities "github.com/CBE-Super-App/cbe-super-app-cps-action/internal/domain/cps_actions/entities"
)

type CPSActionModule interface {
	Authorize(ctx context.Context, action *entities.CPSAction) (*entities.CPSAction, error)
}

type Dispatcher struct {
	app application.Domain
}

func NewDispatcher(app application.Domain) *Dispatcher {
	return &Dispatcher{
		app: app,
	}
}

func (d *Dispatcher) Authorize(ctx context.Context, cpsAction *entities.CPSAction) (*entities.CPSAction, error) {
	switch string(cpsAction.RequestAction) {
	case string(constants.RequestCreateBank),
		string(constants.RequestUpdateBank),
		string(constants.RequestEnableBank),
		string(constants.RequestDisableBank):
		return d.app.BankDomain.Authorize(ctx, cpsAction)
	case string(constants.RequestCreateAvatar),
		string(constants.RequestUpdateAvatar),
		string(constants.RequestEnableAvatar),
		string(constants.RequestDisableAvatar):
		return d.app.AvatarDomian.Authorize(ctx, cpsAction)
	case string(constants.RequestAuthTier):
		return d.app.AmountBasedAuthDomain.Authorize(ctx, cpsAction)
	case string(constants.RequestCreateAdvert),
		string(constants.RequestUpdateAdvert),
		string(constants.RequestEnableAdvert),
		string(constants.RequestDisableAdvert),
		string(constants.RequestDeleteAdvert):
		return d.app.AdDomain.Authorize(ctx, cpsAction)

	default:
		return nil, fmt.Errorf("unsupported request action: %s", cpsAction.RequestAction)
	}
}
