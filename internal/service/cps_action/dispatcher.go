package cpsaction

import (
	"context"
	"errors"
	"fmt"

	"cbe-super-app-cps-action/internal/constants"
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

	fmt.Printf("Authorizing action: %s\n", cpsAction.RequestAction)

	switch {
	case IsActionInGroup(constants.RequestAction(action), "BANK"):
		return d.app.BankContainer.Authorize(ctx, cpsAction)

	case IsActionInGroup(constants.RequestAction(action), "KYCVERIFIER"):
		return d.app.KYCVerifierContainer.Authorize(ctx, cpsAction)

	// case IsActionInGroup(constants.RequestAction(action), "ACCOUNTBLOCK"):
	// 	return d.app.AccountBlockContainer.Authorize(ctx, cpsAction)

	case IsActionInGroup(constants.RequestAction(action), "SINGLEBRANCHENABLEACCOUNTBLOCK"):
		return d.app.AccountBlockContainer.Authorize(ctx, cpsAction)
	case IsActionInGroup(constants.RequestAction(action), "SINGLEBRANCHDISABLEACCOUNTBLOCK"):
		return d.app.AccountBlockContainer.Authorize(ctx, cpsAction)
	case IsActionInGroup(constants.RequestAction(action), "MULTIBRANCHENABLEACCOUNTBLOCK"):
		return d.app.AccountBlockContainer.Authorize(ctx, cpsAction)
	case IsActionInGroup(constants.RequestAction(action), "MULTIBRANCHDISABLEACCOUNTBLOCK"):
		return d.app.AccountBlockContainer.Authorize(ctx, cpsAction)
	case IsActionInGroup(constants.RequestAction(action), "ACCOUNTVALIDATION"):
		return d.app.AccountContainer.Authorize(ctx, cpsAction)

	case IsActionInGroup(constants.RequestAction(action), "ADVERT"):
		return d.app.AdContainer.Authorize(ctx, cpsAction)

	case IsActionInGroup(constants.RequestAction(action), "SERVICE"):
		return d.app.ServiceContainer.Authorize(ctx, cpsAction)

	case IsActionInGroup(constants.RequestAction(action), "SERVICESCATALOG"):
		return d.app.ServiceContainer.Authorize(ctx, cpsAction)

	case IsActionInGroup(constants.RequestAction(action), "DEVICEVERSION"):
		return d.app.DeviceVersionContainer.Authorize(ctx, cpsAction)

	case IsActionInGroup(constants.RequestAction(action), "FAYDA"):
		return d.app.FaydaContainer.Authorize(ctx, cpsAction)

	case IsActionInGroup(constants.RequestAction(action), "MINIAPPMERCHANT"):
		return d.app.MiniAppMerchantContainer.Authorize(ctx, cpsAction)

	case IsActionInGroup(constants.RequestAction(action), "MINIAPP"):
		return d.app.MiniAppContainer.Authorize(ctx, cpsAction)

	case IsActionInGroup(constants.RequestAction(action), "HQ"):
		return d.app.HQContainer.Authorize(ctx, cpsAction)

	case IsActionInGroup(constants.RequestAction(action), "PASSWORDRULE"):
		return d.app.PasswordRuleContainer.Authorize(ctx, cpsAction)

	case IsActionInGroup(constants.RequestAction(action), "PERMISSIONGROUP"):
		return d.app.PermissionContainer.Authorize(ctx, cpsAction)

	case IsActionInGroup(constants.RequestAction(action), "UNLINKDEVICE"):
		return d.app.UnlinkContainer.Authorize(ctx, cpsAction)

	case IsActionInGroup(constants.RequestAction(action), "WALLET"):
		return d.app.WalletContainer.Authorize(ctx, cpsAction)
	case IsActionInGroup(constants.RequestAction(action), "TOPUP"):
		return d.app.TopupContainer.Authorize(ctx, cpsAction)

	case IsActionInGroup(constants.RequestAction(action), "AMOUNTBASEDAUTH"):
		return d.app.AmountBasedAuthContainer.Authorize(ctx, cpsAction)

	case IsActionInGroup(constants.RequestAction(action), "BULKSERVICEALLUSER"):
		return d.app.BulkServiceContainer.Authorize(ctx, cpsAction)

	case IsActionInGroup(constants.RequestAction(action), "EVENT"):
		return d.app.EventContainer.Authorize(ctx, cpsAction)

	case IsActionInGroup(constants.RequestAction(action), "NOTIFICATION"):
		return d.app.NotificationService.Authorize(ctx, cpsAction)

	case IsActionInGroup(constants.RequestAction(action), "USSDMERCHANT"):
		return d.app.UssdMerchantContainer.Authorize(ctx, cpsAction)

	case IsActionInGroup(constants.RequestAction(action), "BPSUSER"):
		return d.app.BPSUserContainer.Authorize(ctx, cpsAction)

	case IsActionInGroup(constants.RequestAction(action), "AVATAR"):
		return d.app.AvatarDomian.Authorize(ctx, cpsAction)

	case IsActionInGroup(constants.RequestAction(action), "DONATIONCATEGORY"):
		return d.app.DonationCategoryContainer.Authorize(ctx, cpsAction)
	case IsActionInGroup(constants.RequestAction(action), "DONATIONCOMPANY"):
		return d.app.DonationCompanyContainer.Authorize(ctx, cpsAction)
	case IsActionInGroup(constants.RequestAction(action), "DONATION"):
		return d.app.DonationContainer.Authorize(ctx, cpsAction)
	case IsActionInGroup(constants.RequestAction(action), "DEPARTMENT"):
		return d.app.DepartmentContainer.Authorize(ctx, cpsAction)
	case IsActionInGroup(constants.RequestAction(action), "CPSUSER"):
		return d.app.CPSUserContainer.Authorize(ctx, cpsAction)
	case IsActionInGroup(constants.RequestAction(action), "BANKVAULT"):
		return d.app.BankProductContainer.Authorize(ctx, cpsAction)
	case IsActionInGroup(constants.RequestAction(action), "VAULTCATEGORY"):
		return d.app.VaultCategoryContainer.Authorize(ctx, cpsAction)
	case IsActionInGroup(constants.RequestAction(action), "ARTICLE"):
		return d.app.ArticleContainer.Authorize(ctx, cpsAction)
	case IsActionInGroup(constants.RequestAction(action), "ARTICLECATEGORY"):
		return d.app.ArticleCategoryContainer.Authorize(ctx, cpsAction)
	case IsActionInGroup(constants.RequestAction(action), "SHORTVIDEO"):
		return d.app.ShortVideoServiceContainer.Authorize(ctx, cpsAction)
	case IsActionInGroup(constants.RequestAction(action), "CUSTOMER"):
		return d.app.CustomerContainer.Authorize(ctx, cpsAction)
	case IsActionInGroup(constants.RequestAction(action), "NEWSTAG"):
		return d.app.NewsTagsServiceContainer.Authorize(ctx, cpsAction)
	case IsActionInGroup(constants.RequestAction(action), "BPSACTIONROLE"):
		return d.app.BPSActionRoleContainer.Authorize(ctx, cpsAction)
	case IsActionInGroup(constants.RequestAction(action), "NEWSCATEGORY"):
		return d.app.NewsCategoryContainer.Authorize(ctx, cpsAction)
	case IsActionInGroup(constants.RequestAction(action), "BUDGETCATEGORY"):
		return d.app.BudgetCategoryContainer.Authorize(ctx, cpsAction)
	case IsActionInGroup(constants.RequestAction(action), "MINIAPPCATEGORY"):
		return d.app.MiniAppCategoryContainer.Authorize(ctx, cpsAction)

	case IsActionInGroup(constants.RequestAction(action), "CPSACTIONROLE"):
		return d.app.CPSActionRoleContainer.Authorize(ctx, cpsAction)
	case IsActionInGroup(constants.RequestAction(action), "EVENTMERCHANT"):
		return d.app.EventMerchantServiceContainer.Authorize(ctx, cpsAction)
	case IsActionInGroup(constants.RequestAction(action), "LOGISTICSMERCHANT"):
		return d.app.LogisticsMerchantServiceContainer.Authorize(ctx, cpsAction)
	case IsActionInGroup(constants.RequestAction(action), "MINIAPPPRODUCTCODE"):
		return d.app.MiniAppProductCodeContainer.Authorize(ctx, cpsAction)
	case IsActionInGroup(constants.RequestAction(action), "ACCESSLISTSEGMENTATION"):
		return d.app.AccessListSegmentationContainer.Authorize(ctx, cpsAction)
	case IsActionInGroup(constants.RequestAction(action), "JOBROLE"):
		return d.app.JobRoleContainer.Authorize(ctx, cpsAction)
	case IsActionInGroup(constants.RequestAction(action), "ROLE"):
		return d.app.RoleContainer.Authorize(ctx, cpsAction)
	case IsActionInGroup(constants.RequestAction(action), "BULKSERVICECUSTOMERSEGMENT"):
		return d.app.AccessListSegmentationContainer.Authorize(ctx, cpsAction)
	case IsActionInGroup(constants.RequestAction(action), "ECOMMERCEMERCHANT"):
		return d.app.EcommerceMerchantContainer.Authorize(ctx, cpsAction)
	case IsActionInGroup(constants.RequestAction(action), "CPSROLE"):
		return d.app.CPSRolesContainer.Authorize(ctx, cpsAction)
	case IsActionInGroup(constants.RequestAction(action), "CUSTOMERKYC"):
		return d.app.CustomerKYCContainer.Authorize(ctx, cpsAction)
	// case IsActionInGroup(constants.RequestAction(action), "VAULT"):
	// 	return d.app.VaultCategoryContainer.Authorize(ctx, cpsAction)

	case IsActionInGroup(constants.RequestAction(action), "VAULTCATEGORIES"):
		return d.app.VaultCategoryContainer.Authorize(ctx, cpsAction)
	case IsActionInGroup(constants.RequestAction(action), "VAULTEMERGENCYDEADLOCKREQUEST"):
		return d.app.VaultCategoryContainer.Authorize(ctx, cpsAction)
	default:
		span.AddEvent("unsupported action", trace.WithAttributes(attribute.String("action", action)))
		fmt.Printf("Unsupported action---------------------: %s\n", action)
		return nil, errors.New(localization.ErrorUnsupportedAction.Code)
	}
}
