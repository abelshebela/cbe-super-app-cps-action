package cpsaction

import (
	"context"
	"errors"
	"fmt"

	"cbe-super-app-cps-action/internal/constants/model"
	"cbe-super-app-cps-action/internal/service"
)

type CPSActionModule interface {
	Authorize(ctx context.Context, action *model.CPSAction) (*model.CPSAction, error)
}

type Dispatcher struct {
	app service.ServiceContainer
}

func NewDispatcher(app service.ServiceContainer) *Dispatcher {
	return &Dispatcher{
		app: app,
	}
}

func (d *Dispatcher) Authorize(ctx context.Context, cpsAction *model.CPSAction) (*model.CPSAction, error) {
	action := cpsAction.RequestAction

	switch {
	case IsActionInGroup(RequestAction(action), "Bank"):
		return d.app.BankContainer.Authorize(ctx, cpsAction)

	case IsActionInGroup(RequestAction(action), "Block"):
		return d.app.AccountBlockContainer.Authorize(ctx, cpsAction)

	case IsActionInGroup(RequestAction(action), "Account"):
		return d.app.AccountContainer.Authorize(ctx, cpsAction)

	case IsActionInGroup(RequestAction(action), "Advert"):
		return d.app.AdContainer.Authorize(ctx, cpsAction)

	case IsActionInGroup(RequestAction(action), "Service"):
		return d.app.ServiceCheckContainer.Authorize(ctx, cpsAction)

	case IsActionInGroup(RequestAction(action), "Fayda"):
		return d.app.FaydaContainer.Authorize(ctx, cpsAction)

	case IsActionInGroup(RequestAction(action), "MiniAppMerchant"):
		return d.app.MiniAppMerchantContainer.Authorize(ctx, cpsAction)

	case IsActionInGroup(RequestAction(action), "MiniApp"):
		return d.app.MiniAppContainer.Authorize(ctx, cpsAction)

	case IsActionInGroup(RequestAction(action), "HQ"):
		return d.app.HQContainer.Authorize(ctx, cpsAction)

	case IsActionInGroup(RequestAction(action), "Password"):
		return d.app.PasswordRuleContainer.Authorize(ctx, cpsAction)

	case IsActionInGroup(RequestAction(action), "Permission"):
		return d.app.PermissionContainer.Authorize(ctx, cpsAction)

	case IsActionInGroup(RequestAction(action), "UnlinkDevice"):
		return d.app.UnlinkContainer.Authorize(ctx, cpsAction)

	case IsActionInGroup(RequestAction(action), "Wallet"):
		return d.app.WalletContainer.Authorize(ctx, cpsAction)

	case IsActionInGroup(RequestAction(action), "AmountBasedAuth"):
		return d.app.AmountBasedAuthContainer.Authorize(ctx, cpsAction)

	case IsActionInGroup(RequestAction(action), "BudgetCategory"):
		return d.app.BudgetCategoryContainer.Authorize(ctx, cpsAction)

	case IsActionInGroup(RequestAction(action), "Budget"):
		return d.app.BudgetContainer.Authorize(ctx, cpsAction)
	case IsActionInGroup(RequestAction(action), "BulkService"):
		return d.app.BulkServiceContainer.Authorize(ctx, cpsAction)

	case IsActionInGroup(RequestAction(action), "Event"):
		return d.app.EventContainer.Authorize(ctx, cpsAction)

	case IsActionInGroup(RequestAction(action), "Notification"):
		return d.app.NotificationService.Authorize(ctx, cpsAction)

	case IsActionInGroup(RequestAction(action), "ProductCode"):
		if d.app.ProductCodeService == nil {
			return nil, errors.New("ProductCode service is not implemented")
		}
		return d.app.ProductCodeService.Authorize(ctx, cpsAction)

	case IsActionInGroup(RequestAction(action), "BPSUser"):
		return d.app.BPSUserContainer.Authorize(ctx, cpsAction)

	case IsActionInGroup(RequestAction(action), "Avatar"):
		return d.app.AvatarDomian.Authorize(ctx, cpsAction)

	case IsActionInGroup(RequestAction(action), "Donation"):
		return d.app.DonationContainer.Authorize(ctx, cpsAction)

	case IsActionInGroup(RequestAction(action), "Department"):
		return d.app.DepartmentContainer.Authorize(ctx, cpsAction)
	default:
		return nil, fmt.Errorf("UNSUPPORTED_REQUEST_ACTION")
	}
}
