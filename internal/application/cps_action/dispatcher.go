package cpsaction

import (
	"context"
	"fmt"

	"github.com/CBE-Super-App/cbe-super-app-cps-action/internal/application"
	constants "github.com/CBE-Super-App/cbe-super-app-cps-action/internal/domain/cps_actions/constant"
	entities "github.com/CBE-Super-App/cbe-super-app-cps-action/internal/domain/cps_actions/entities"
	local_util "github.com/CBE-Super-App/cbe-super-app-cps-action/pkgs/utils"
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

	fmt.Println(action, "action")
	switch {
	case constants.IsActionInGroup(action, "Bank"):
		return d.app.BankDomain.Authorize(ctx, cpsAction)

	case constants.IsActionInGroup(action, "Block"):
		return d.app.AccountBlockDomain.Authorize(ctx, cpsAction)

	case constants.IsActionInGroup(action, "Account"):
		return d.app.AccountDomain.Authorize(ctx, cpsAction)

	case constants.IsActionInGroup(action, "Advert"):
		return d.app.AdDomain.Authorize(ctx, cpsAction)

	case constants.IsActionInGroup(action, "Service"):
		// fmt.Println("*********IsActionInGroup*********")
		data, err := d.app.ServiceCheckDomain.Authorize(ctx, cpsAction)
		if err != nil {
			return nil, err
		}
		return marshalBuilder(data)
	case constants.IsActionInGroup(action, "Fayda"):

		return d.app.FaydaDomain.Authorize(ctx, cpsAction)

	case constants.IsActionInGroup(action, "MiniAppMerchant"):
		return d.app.MiniAppMerchantDomain.Authorize(ctx, cpsAction)
	case constants.IsActionInGroup(action, "MiniApp"):
		return d.app.MiniAppDomain.Authorize(ctx, cpsAction)
		// return d.app.MiniAppDomain.Authorize(ctx, cpsAction)

	case constants.IsActionInGroup(action, "HQ"):
		return d.app.HQDomain.Authorize(ctx, cpsAction)

	case constants.IsActionInGroup(action, "Password"):
		return d.app.PasswordRuleDomain.Authorize(ctx, cpsAction)

	case constants.IsActionInGroup(action, "Permission"):
		return d.app.PermissionDomain.Authorize(ctx, cpsAction)

	// case constants.IsActionInGroup(action, "UnlinkDevice"):
	// 	return d.app.UnlinkDomain.Authorize(ctx, cpsAction)

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
	case constants.IsActionInGroup(action, "Event"):
		return d.app.EventDomain.Authorize(ctx, cpsAction)
	case constants.IsActionInGroup(action, "Service"):
		return d.app.BudgetDomain.Authorize(ctx, cpsAction)
	case constants.IsActionInGroup(action, "Notification"):
		return d.app.NotificationService.Authorize(ctx, cpsAction)
	case constants.IsActionInGroup(action, "ProductCode"):
		return d.app.ProductCodeService.Authorize(ctx, cpsAction)
		case constants.IsActionInGroup(action, "BPSUser"):
		return d.app.BPSUserDomain.Authorize(ctx, cpsAction)
	default:
		return nil, fmt.Errorf("UNSUPPORTED_REQUEST_ACTION")
	}
}

func marshalBuilder(data any) (*entities.CPSAction, error) {
	action, err := local_util.JsonUnmarshal[*entities.CPSAction](data)
	if err != nil {
		return nil, err
	}
	return *action, nil
}
