package cpsaction

import (
	"context"
	"errors"

	"cbe-super-app-cps-action/internal/constants/localization"
	"cbe-super-app-cps-action/internal/service"
	local_util "cbe-super-app-cps-action/pkgs/utils"

	"go.opentelemetry.io/otel/attribute"
	"go.opentelemetry.io/otel/trace"

	"gitlab.com/bersufekadgetachew/cbe-super-app-shared/shared/entities/model"
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
	ctx, span := local_util.TraceLogger(ctx, "service", "Authorize", "Dispatcher", "Authorize")
	defer span.End()

	action := cpsAction.RequestAction
	span.SetAttributes(attribute.String("action", action))

	switch {
	case IsActionInGroup(RequestAction(action), "BANK"):
		return d.app.BankContainer.Authorize(ctx, cpsAction)

	case IsActionInGroup(RequestAction(action), "KYCVERIFIER"):
		return d.app.KYCVerifierContainer.Authorize(ctx, cpsAction)

	case IsActionInGroup(RequestAction(action), "ACCOUNTBLOCK"):
		return d.app.AccountBlockContainer.Authorize(ctx, cpsAction)

	case IsActionInGroup(RequestAction(action), "ACCOUNTVALIDATION"):
		return d.app.AccountContainer.Authorize(ctx, cpsAction)

	case IsActionInGroup(RequestAction(action), "ADVERT"):
		return d.app.AdContainer.Authorize(ctx, cpsAction)

	case IsActionInGroup(RequestAction(action), "SERVICE"):
		return d.app.ServiceContainer.Authorize(ctx, cpsAction)

	case IsActionInGroup(RequestAction(action), "ServicesCatalog"):
		return d.app.ServiceContainer.Authorize(ctx, cpsAction)

	case IsActionInGroup(RequestAction(action), "DEVICEVERSION"):
		return d.app.DeviceVersionContainer.Authorize(ctx, cpsAction)

	case IsActionInGroup(RequestAction(action), "FAYDA"):
		return d.app.FaydaContainer.Authorize(ctx, cpsAction)

	case IsActionInGroup(RequestAction(action), "MINIAPPMERCHANT"):
		return d.app.MiniAppMerchantContainer.Authorize(ctx, cpsAction)

	case IsActionInGroup(RequestAction(action), "MINIAPP"):
		return d.app.MiniAppContainer.Authorize(ctx, cpsAction)

	case IsActionInGroup(RequestAction(action), "HQ"):
		return d.app.HQContainer.Authorize(ctx, cpsAction)

	case IsActionInGroup(RequestAction(action), "PASSWORDRULES"):
		return d.app.PasswordRuleContainer.Authorize(ctx, cpsAction)

	case IsActionInGroup(RequestAction(action), "PERMISSION"):
		return d.app.PermissionContainer.Authorize(ctx, cpsAction)

	case IsActionInGroup(RequestAction(action), "UNLINKDEVICE"):
		return d.app.UnlinkContainer.Authorize(ctx, cpsAction)

	case IsActionInGroup(RequestAction(action), "WALLET"):
		return d.app.WalletContainer.Authorize(ctx, cpsAction)
	case IsActionInGroup(RequestAction(action), "TOPUP"):
		return d.app.TopupContainer.Authorize(ctx, cpsAction)

	case IsActionInGroup(RequestAction(action), "AMOUNTBASEDAUTH"):
		return d.app.AmountBasedAuthContainer.Authorize(ctx, cpsAction)

	case IsActionInGroup(RequestAction(action), "BULKSERVICE"):
		return d.app.BulkServiceContainer.Authorize(ctx, cpsAction)

	case IsActionInGroup(RequestAction(action), "EVENT"):
		return d.app.EventContainer.Authorize(ctx, cpsAction)

	case IsActionInGroup(RequestAction(action), "NOTIFICATION"):
		return d.app.NotificationService.Authorize(ctx, cpsAction)

	case IsActionInGroup(RequestAction(action), "PRODUCTCODE"):
		return d.app.ProductCodeService.Authorize(ctx, cpsAction)

	case IsActionInGroup(RequestAction(action), "BPSUSER"):
		return d.app.BPSUserContainer.Authorize(ctx, cpsAction)

	case IsActionInGroup(RequestAction(action), "AVATAR"):
		return d.app.AvatarDomian.Authorize(ctx, cpsAction)

	case IsActionInGroup(RequestAction(action), "DONATIONCATEGORY"):
		return d.app.DonationCategoryContainer.Authorize(ctx, cpsAction)
	case IsActionInGroup(RequestAction(action), "DONATIONCOMPANY"):
		return d.app.DonationCompanyContainer.Authorize(ctx, cpsAction)
	case IsActionInGroup(RequestAction(action), "DONATION"):
		return d.app.DonationContainer.Authorize(ctx, cpsAction)
	case IsActionInGroup(RequestAction(action), "DEPARTMENT"):
		return d.app.DepartmentContainer.Authorize(ctx, cpsAction)
	case IsActionInGroup(RequestAction(action), "CPSUSER"):
		return d.app.CPSUserContainer.Authorize(ctx, cpsAction)
	case IsActionInGroup(RequestAction(action), "BANKVAULT"):
		return d.app.BankProductContainer.Authorize(ctx, cpsAction)
	case IsActionInGroup(RequestAction(action), "VAULTCATEGORY"):
		return d.app.VaultCategoryContainer.Authorize(ctx, cpsAction)
	case IsActionInGroup(RequestAction(action), "ARTICLE"):
		return d.app.ArticleContainer.Authorize(ctx, cpsAction)
	case IsActionInGroup(RequestAction(action), "ARTICLECATEGORY"):
		return d.app.ArticleCategoryContainer.Authorize(ctx, cpsAction)
	case IsActionInGroup(RequestAction(action), "SHORTVIDEO"):
		return d.app.ShortVideoServiceContainer.Authorize(ctx, cpsAction)
	case IsActionInGroup(RequestAction(action), "CUSTOMER"):
		return d.app.CustomerContainer.Authorize(ctx, cpsAction)
	case IsActionInGroup(RequestAction(action), "NEWSTAG"):
		return d.app.NewsTagsServiceContainer.Authorize(ctx, cpsAction)
	case IsActionInGroup(RequestAction(action), "ACTIONROLE"):
		return d.app.BPSActionRoleContainer.Authorize(ctx, cpsAction)
	case IsActionInGroup(RequestAction(action), "NEWSCATEGORY"):
		return d.app.NewsCategoryContainer.Authorize(ctx, cpsAction)
	case IsActionInGroup(RequestAction(action), "BUDGETCATEGORY"):
		return d.app.BudgetCategoryContainer.Authorize(ctx, cpsAction)
	case IsActionInGroup(RequestAction(action), "MINIAPPCATEGORY"):
		return d.app.MiniAppCategoryContainer.Authorize(ctx, cpsAction)
	case IsActionInGroup(RequestAction(action), "CPSACTIONROLE"):
		return d.app.CPSActionRoleContainer.Authorize(ctx, cpsAction)
	case IsActionInGroup(RequestAction(action), "EVENTMERCHANT"):
		return d.app.EventMerchantServiceContainer.Authorize(ctx, cpsAction)
	case IsActionInGroup(RequestAction(action), "VAULTAMOUNTTIER"):
		return d.app.VaultAmountTierContainer.Authorize(ctx, cpsAction)
	case IsActionInGroup(RequestAction(action), "MINIAPPPRODUCTCODE"):
		return d.app.MiniAppProductCodeContainer.Authorize(ctx, cpsAction)
	case IsActionInGroup(RequestAction(action), "ACCESSLISTSEGMENTATION"):
		return d.app.AccessListSegmentationContainer.Authorize(ctx, cpsAction)
	case IsActionInGroup(RequestAction(action), "JOBROLE"):
		return d.app.JobRoleContainer.Authorize(ctx, cpsAction)
	case IsActionInGroup(RequestAction(action), "ROLE"):
		return d.app.RoleContainer.Authorize(ctx, cpsAction)
	case IsActionInGroup(RequestAction(action), "CUSTOMERSEGMENTATIONS"):
		return d.app.CustomerSegmentationContainer.Authorize(ctx, cpsAction)
	case IsActionInGroup(RequestAction(action), "ECOMMERCEMERCHANT"):
		return d.app.EcommerceMerchantContainer.Authorize(ctx, cpsAction)

	default:
		span.AddEvent("unsupported action", trace.WithAttributes(attribute.String("action", action)))
		return nil, errors.New(localization.ErrorUnsupportedAction.Code)
	}
}
