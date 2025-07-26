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
	Reject(ctx context.Context, action *entities.CPSAction) (*entities.CPSAction, error)
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
	action := cpsAction.RequestAction

	switch {
	case constants.IsActionInGroup(action, "Bank"):
		return d.app.BankDomain.Authorize(ctx, cpsAction)

	case constants.IsActionInGroup(action, "Block"):
		return d.app.AccountBlockDomain.Authorize(ctx, cpsAction)

	case constants.IsActionInGroup(action, "Department"):
		return d.app.DepartmentDomain.Authorize(ctx, cpsAction)

	case constants.IsActionInGroup(action, "Account"):
		return d.app.AccountDomain.Authorize(ctx, cpsAction)

	case constants.IsActionInGroup(action, "Advert"):
		return d.app.AdDomain.Authorize(ctx, cpsAction)

	case constants.IsActionInGroup(action, "Avatar"):
		return d.app.AvatarDomian.Authorize(ctx, cpsAction)

	case constants.IsActionInGroup(action, "Fayda"):
		return d.app.FaydaDomain.Authorize(ctx, cpsAction)

	case constants.IsActionInGroup(action, "MiniAppMerchant"):
		return d.app.MiniAppMerchantDomain.Authorize(ctx, cpsAction)

	case constants.IsActionInGroup(action, "HQ"):
		return d.app.HQDomain.Authorize(ctx, cpsAction)

	case constants.IsActionInGroup(action, "Password"):
		return d.app.PasswordRuleDomain.Authorize(ctx, cpsAction)

	case constants.IsActionInGroup(action, "Permission"):
		return d.app.PermissionDomain.Authorize(ctx, cpsAction)

	case constants.IsActionInGroup(action, "UnlinkDevice"):
		return d.app.UnlinkDomain.Authorize(ctx, cpsAction)

	case constants.IsActionInGroup(action, "Wallet"):
		return d.app.WalletDomain.Authorize(ctx, cpsAction)

	case constants.IsActionInGroup(action, "AmountBasedAuth"):
		return d.app.AmountBasedAuthDomain.Authorize(ctx, cpsAction)

	case constants.IsActionInGroup(action, "BudgetCategory"):
		return d.app.BudgetCategoryDomain.Authorize(ctx, cpsAction)

	case constants.IsActionInGroup(action, "Budget"):
		return d.app.BudgetDomain.Authorize(ctx, cpsAction)
	case constants.IsActionInGroup(action, "CPSUser"):
		return d.app.CPSUserDomain.Authorize(ctx, cpsAction)
	case constants.IsActionInGroup(action, "BulkService"):
		return d.app.BulkServiceDomain.Authorize(ctx, cpsAction)
	default:
		return nil, fmt.Errorf("UNSUPPORTED_REQUEST_ACTION")
	}
}
