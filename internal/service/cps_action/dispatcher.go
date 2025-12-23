package cpsaction

import (
	"context"
	"errors"
	"fmt"

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
	case IsActionInGroup(RequestAction(action), "Bank"):
		return d.app.BankContainer.Authorize(ctx, cpsAction)

	case IsActionInGroup(RequestAction(action), "KYCVerifier"):
		return d.app.KYCVerifierContainer.Authorize(ctx, cpsAction)

	case IsActionInGroup(RequestAction(action), "Block"):
		return d.app.AccountBlockContainer.Authorize(ctx, cpsAction)

	case IsActionInGroup(RequestAction(action), "Account"):
		return d.app.AccountContainer.Authorize(ctx, cpsAction)

	case IsActionInGroup(RequestAction(action), "Advert"):
		return d.app.AdContainer.Authorize(ctx, cpsAction)

	case IsActionInGroup(RequestAction(action), "Service"):
		return d.app.ServiceContainer.Authorize(ctx, cpsAction)

	case IsActionInGroup(RequestAction(action), "ServicesCatalog"):
		return d.app.ServiceContainer.Authorize(ctx, cpsAction)

	case IsActionInGroup(RequestAction(action), "DeviceVersion"):
		return d.app.DeviceVersionContainer.Authorize(ctx, cpsAction)

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
	case IsActionInGroup(RequestAction(action), "Topup"):
		return d.app.TopupContainer.Authorize(ctx, cpsAction)

	case IsActionInGroup(RequestAction(action), "AmountBasedAuth"):
		return d.app.AmountBasedAuthContainer.Authorize(ctx, cpsAction)

	case IsActionInGroup(RequestAction(action), "BulkService"):
		return d.app.BulkServiceContainer.Authorize(ctx, cpsAction)

	case IsActionInGroup(RequestAction(action), "Event"):
		return d.app.EventContainer.Authorize(ctx, cpsAction)

	case IsActionInGroup(RequestAction(action), "Notification"):
		return d.app.NotificationService.Authorize(ctx, cpsAction)

	case IsActionInGroup(RequestAction(action), "ProductCode"):
		return d.app.ProductCodeService.Authorize(ctx, cpsAction)

	case IsActionInGroup(RequestAction(action), "BPSUser"):
		return d.app.BPSUserContainer.Authorize(ctx, cpsAction)

	case IsActionInGroup(RequestAction(action), "Avatar"):
		return d.app.AvatarDomian.Authorize(ctx, cpsAction)

	case IsActionInGroup(RequestAction(action), "donationCategory"):
		return d.app.DonationCategoryContainer.Authorize(ctx, cpsAction)
	case IsActionInGroup(RequestAction(action), "donationCompany"):
		return d.app.DonationCompanyContainer.Authorize(ctx, cpsAction)
	case IsActionInGroup(RequestAction(action), "Donation"):
		return d.app.DonationContainer.Authorize(ctx, cpsAction)
	case IsActionInGroup(RequestAction(action), "Department"):
		return d.app.DepartmentContainer.Authorize(ctx, cpsAction)
	case IsActionInGroup(RequestAction(action), "CPSUser"):
		return d.app.CPSUserContainer.Authorize(ctx, cpsAction)
	case IsActionInGroup(RequestAction(action), "BankVault"):
		return d.app.BankProductContainer.Authorize(ctx, cpsAction)
	case IsActionInGroup(RequestAction(action), "VaultGroupCategory"):
		return d.app.VaultCategoryContainer.Authorize(ctx, cpsAction)
	case IsActionInGroup(RequestAction(action), "article"):
		return d.app.ArticleContainer.Authorize(ctx, cpsAction)
	case IsActionInGroup(RequestAction(action), "articleCategory"):
		return d.app.ArticleCategoryContainer.Authorize(ctx, cpsAction)
	case IsActionInGroup(RequestAction(action), "short_video"):
		return d.app.ShortVideoServiceContainer.Authorize(ctx, cpsAction)
	case IsActionInGroup(RequestAction(action), "customer"):
		return d.app.CustomerContainer.Authorize(ctx, cpsAction)
	case IsActionInGroup(RequestAction(action), "news_tag"):
		return d.app.NewsTagsServiceContainer.Authorize(ctx, cpsAction)
	case IsActionInGroup(RequestAction(action), "ActionRole"):
		return d.app.BPSActionRoleContainer.Authorize(ctx, cpsAction)
	case IsActionInGroup(RequestAction(action), "news_category"):
		return d.app.NewsCategoryContainer.Authorize(ctx, cpsAction)
	case IsActionInGroup(RequestAction(action), "BudgetCategory"):
		return d.app.BudgetCategoryContainer.Authorize(ctx, cpsAction)
	case IsActionInGroup(RequestAction(action), "MiniAppCategory"):
		return d.app.MiniAppCategoryContainer.Authorize(ctx, cpsAction)
	case IsActionInGroup(RequestAction(action), "CpsActionRole"):
		return d.app.CPSActionRoleContainer.Authorize(ctx, cpsAction)
	case IsActionInGroup(RequestAction(action), "eventMerchant"):
		return d.app.EventMerchantServiceContainer.Authorize(ctx, cpsAction)
	case IsActionInGroup(RequestAction(action), "VaultAmountTier"):
		return d.app.VaultAmountTierContainer.Authorize(ctx, cpsAction)
	case IsActionInGroup(RequestAction(action), "MiniAppProductCode"):
		return d.app.MiniAppProductCodeContainer.Authorize(ctx, cpsAction)
	case IsActionInGroup(RequestAction(action), "AccessListSegmentation"):
		return d.app.AccessListSegmentationContainer.Authorize(ctx, cpsAction)
	case IsActionInGroup(RequestAction(action), "JobRole"):
		return d.app.JobRoleContainer.Authorize(ctx, cpsAction)
	case IsActionInGroup(RequestAction(action), "Role"):
		return d.app.RoleContainer.Authorize(ctx, cpsAction)
	case IsActionInGroup(RequestAction(action), "CustomerSegmentations"):
		fmt.Println("=====================Got here")
		return d.app.CustomerSegmentationContainer.Authorize(ctx, cpsAction)

	default:
		span.AddEvent("unsupported action", trace.WithAttributes(attribute.String("action", action)))
		return nil, errors.New(localization.ErrorUnsupportedAction.Code)
	}
}
