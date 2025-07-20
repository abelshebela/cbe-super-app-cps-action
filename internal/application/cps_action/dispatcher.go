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

	fmt.Println(action, "action")
	switch {
	// case constants.IsActionInGroup(action, "Block"):
	// 	return d.app.AccountBlockDomain.Authorize(ctx, cpsAction)
	case constants.IsActionInGroup(action, "Bank"):
		return d.app.BankDomain.Authorize(ctx, cpsAction)
	case constants.IsActionInGroup(action, "Avatar"):
		return d.app.AvatarDomian.Authorize(ctx, cpsAction)
	case constants.IsActionInGroup(action, "Advert"):
		return d.app.AdDomain.Authorize(ctx, cpsAction)

	// case constants.IsActionInGroup(action, "BPSUser"):
	// 	return d.app.BPSUserDomain.Authorize(ctx, cpsAction)

	// case constants.IsActionInGroup(action, "Wallet"):
	// 	return d.app.WalletDomain.Authorize(ctx, cpsAction)

	// case constants.IsActionInGroup(action, "Validation"):
	// 	return d.app.ValidationDomain.Authorize(ctx, cpsAction)

	// case constants.IsActionInGroup(action, "ServiceFee"):
	// 	return d.app.ServiceFeeDomain.Authorize(ctx, cpsAction)

	// case constants.IsActionInGroup(action, "DailyLimit"):
	// 	return d.app.DailyLimitDomain.Authorize(ctx, cpsAction)

	// case constants.IsActionInGroup(action, "VAT"):
	// 	return d.app.VATDomain.Authorize(ctx, cpsAction)

	// case constants.IsActionInGroup(action, "AuthTier"):
	// 	return d.app.AmountBasedAuthDomain.Authorize(ctx, cpsAction)

	// case constants.IsActionInGroup(action, "Password"):
	// 	return d.app.PasswordPolicyDomain.Authorize(ctx, cpsAction)

	// case constants.IsActionInGroup(action, "Archive"):
	// 	return d.app.ArchiveDomain.Authorize(ctx, cpsAction)

	// case constants.IsActionInGroup(action, "MinimumService"):
	// 	return d.app.MinimumServiceDomain.Authorize(ctx, cpsAction)

	// case constants.IsActionInGroup(action, "ServiceRule"):
	// 	return d.app.ServiceRuleDomain.Authorize(ctx, cpsAction)

	// case constants.IsActionInGroup(action, "Total"):
	// 	return d.app.TotalDomain.Authorize(ctx, cpsAction)

	// case constants.IsActionInGroup(action, "AccessConfig"):
	// 	return d.app.AccessConfigDomain.Authorize(ctx, cpsAction)

	// case constants.IsActionInGroup(action, "Branch"):
	// 	return d.app.BranchDomain.Authorize(ctx, cpsAction)

	// case constants.IsActionInGroup(action, "Business"):
	// 	return d.app.BusinessDomain.Authorize(ctx, cpsAction)

	// case constants.IsActionInGroup(action, "Event"):
	// 	return d.app.EventDomain.Authorize(ctx, cpsAction)

	// case constants.IsActionInGroup(action, "EventCategory"):
	// 	return d.app.EventCategoryDomain.Authorize(ctx, cpsAction)

	// case constants.IsActionInGroup(action, "MiniAppMerchant"):
	// 	return d.app.MiniAppMerchantDomain.Authorize(ctx, cpsAction)

	// case constants.IsActionInGroup(action, "BlockTime"):
	// 	return d.app.BlockTimeDomain.Authorize(ctx, cpsAction)

	// case constants.IsActionInGroup(action, "Budget"):
	// 	return d.app.BudgetDomain.Authorize(ctx, cpsAction)

	// case constants.IsActionInGroup(action, "Product"):
	// 	return d.app.ProductDomain.Authorize(ctx, cpsAction)

	// case constants.IsActionInGroup(action, "PublicNotification"):
	// 	return d.app.PublicNotificationDomain.Authorize(ctx, cpsAction)

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

	case constants.IsActionInGroup(action, "Wallet"):
		return d.app.WalletDomain.Authorize(ctx, cpsAction)

	default:
		return nil, fmt.Errorf("unsupported request action: %s", action)
	}
}
