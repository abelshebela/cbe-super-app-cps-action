package services

import (
	"context"
	"fmt"

	accountBlockService "github.com/CBE-Super-App/cbe-super-app-cps-action/internal/domain/account_block"
	adService "github.com/CBE-Super-App/cbe-super-app-cps-action/internal/domain/ad"
	amountBasedAuthService "github.com/CBE-Super-App/cbe-super-app-cps-action/internal/domain/amount_based_auth"
	avatarService "github.com/CBE-Super-App/cbe-super-app-cps-action/internal/domain/avatar"
	bankService "github.com/CBE-Super-App/cbe-super-app-cps-action/internal/domain/bank/service"
	bpsUserService "github.com/CBE-Super-App/cbe-super-app-cps-action/internal/domain/bps_user"
	budgetService "github.com/CBE-Super-App/cbe-super-app-cps-action/internal/domain/budget"
	budgetCategoryService "github.com/CBE-Super-App/cbe-super-app-cps-action/internal/domain/budget_category"
	bulkServiceService "github.com/CBE-Super-App/cbe-super-app-cps-action/internal/domain/bulk_service"
	constants "github.com/CBE-Super-App/cbe-super-app-cps-action/internal/domain/cps_actions/constant"
	entities "github.com/CBE-Super-App/cbe-super-app-cps-action/internal/domain/cps_actions/entities"
	cpsActionService "github.com/CBE-Super-App/cbe-super-app-cps-action/internal/domain/cps_actions/services"
	cpsUserService "github.com/CBE-Super-App/cbe-super-app-cps-action/internal/domain/cps_user/services"
	donationService "github.com/CBE-Super-App/cbe-super-app-cps-action/internal/domain/donation"
	eventService "github.com/CBE-Super-App/cbe-super-app-cps-action/internal/domain/event"
	faydaService "github.com/CBE-Super-App/cbe-super-app-cps-action/internal/domain/fayda_account/service"
	hqService "github.com/CBE-Super-App/cbe-super-app-cps-action/internal/domain/hq"
	miniAppService "github.com/CBE-Super-App/cbe-super-app-cps-action/internal/domain/miniapp"
	miniAppMerchantService "github.com/CBE-Super-App/cbe-super-app-cps-action/internal/domain/miniapp_merchant"
	notificationService "github.com/CBE-Super-App/cbe-super-app-cps-action/internal/domain/notification"
	passwordRuleService "github.com/CBE-Super-App/cbe-super-app-cps-action/internal/domain/password_rule/services"
	permissionService "github.com/CBE-Super-App/cbe-super-app-cps-action/internal/domain/permission"
	productCodeService "github.com/CBE-Super-App/cbe-super-app-cps-action/internal/domain/product_code"
	serviceCheckService "github.com/CBE-Super-App/cbe-super-app-cps-action/internal/domain/service"
	unlinkService "github.com/CBE-Super-App/cbe-super-app-cps-action/internal/domain/unlink"
	walletService "github.com/CBE-Super-App/cbe-super-app-cps-action/internal/domain/wallet"
	local_util "github.com/CBE-Super-App/cbe-super-app-cps-action/pkgs/utils"

	accountValidationService "github.com/CBE-Super-App/cbe-super-app-cps-action/internal/domain/account_validation"
)

// DomainInterface defines the interface for domain services that the dispatcher needs
type DomainInterface interface {
	CPSActionDomain() cpsActionService.CPSActionService
	BankDomain() bankService.BankService
	AccountBlockDomain() accountBlockService.AccountService
	AccountDomain() accountValidationService.Service
	AdDomain() adService.AdvertService
	ServiceCheckDomain() serviceCheckService.ServiceRepo
	FaydaDomain() faydaService.FaydaAccount
	MiniAppMerchantDomain() miniAppMerchantService.MiniAppMerchantService
	MiniAppDomain() miniAppService.MiniAppService
	HQDomain() hqService.Service
	PasswordRuleDomain() passwordRuleService.PasswordRuleService
	PermissionDomain() permissionService.Service
	UnlinkDomain() unlinkService.UnlinkAccount
	WalletDomain() walletService.WalletService
	AmountBasedAuthDomain() amountBasedAuthService.AmountBasedAuthDomain
	BudgetCategoryDomain() budgetCategoryService.BudgetCategoryService
	BudgetDomain() *budgetService.BudgetService
	CPSUserDomain() cpsUserService.CPSUserService
	BulkServiceDomain() bulkServiceService.BulkService
	EventDomain() eventService.EventService
	NotificationService() notificationService.NotificationService
	ProductCodeService() productCodeService.Service
	BPSUserDomain() bpsUserService.Service
	AvatarDomain() avatarService.AvatarDomainService
	DonationDomain() donationService.DonationService
}

type Dispatcher struct {
	app DomainInterface
}

func NewDispatcher(app DomainInterface) *Dispatcher {
	return &Dispatcher{
		app: app,
	}
}

func (d *Dispatcher) Authorize(ctx context.Context, cpsAction *entities.CPSAction) (*entities.CPSAction, error) {
	action := cpsAction.RequestAction

	fmt.Println(action, "action")
	switch {
	case constants.IsActionInGroup(action, "Bank"):
		return d.app.BankDomain().Authorize(ctx, cpsAction)

	case constants.IsActionInGroup(action, "Block"):
		return d.app.AccountBlockDomain().Authorize(ctx, cpsAction)

	case constants.IsActionInGroup(action, "Account"):
		return d.app.AccountDomain().Authorize(ctx, cpsAction)

	case constants.IsActionInGroup(action, "Advert"):
		return d.app.AdDomain().Authorize(ctx, cpsAction)

	case constants.IsActionInGroup(action, "Service"):
		// fmt.Println("*********IsActionInGroup*********")
		data, err := d.app.ServiceCheckDomain().Authorize(ctx, cpsAction)
		if err != nil {
			return nil, err
		}
		return marshalBuilder(data)
	case constants.IsActionInGroup(action, "Fayda"):

		return d.app.FaydaDomain().Authorize(ctx, cpsAction)

	case constants.IsActionInGroup(action, "MiniAppMerchant"):
		return d.app.MiniAppMerchantDomain().Authorize(ctx, cpsAction)
	case constants.IsActionInGroup(action, "MiniApp"):
		return d.app.MiniAppDomain().Authorize(ctx, cpsAction)
		// return d.app.MiniAppDomain.Authorize(ctx, cpsAction)

	case constants.IsActionInGroup(action, "HQ"):
		return d.app.HQDomain().Authorize(ctx, cpsAction)

	case constants.IsActionInGroup(action, "Password"):
		return d.app.PasswordRuleDomain().Authorize(ctx, cpsAction)

	case constants.IsActionInGroup(action, "Permission"):
		return d.app.PermissionDomain().Authorize(ctx, cpsAction)

	case constants.IsActionInGroup(action, "UnlinkDevice"):
		data, err := d.app.UnlinkDomain().Authorize(ctx, cpsAction)
		if err != nil {
			return nil, err
		}
		return marshalBuilder(data)

	case constants.IsActionInGroup(action, "Wallet"):
		return d.app.WalletDomain().Authorize(ctx, cpsAction)

	case constants.IsActionInGroup(action, "AmountBasedAuth"):
		return d.app.AmountBasedAuthDomain().Authorize(ctx, cpsAction)

	case constants.IsActionInGroup(action, "BudgetCategory"):
		return d.app.BudgetCategoryDomain().Authorize(ctx, cpsAction)

	case constants.IsActionInGroup(action, "Budget"):
		return d.app.BudgetDomain().Authorize(ctx, cpsAction)
	case constants.IsActionInGroup(action, "CPSUser"):
		return d.app.CPSUserDomain().Authorize(ctx, cpsAction)
	case constants.IsActionInGroup(action, "BulkService"):
		return d.app.BulkServiceDomain().Authorize(ctx, cpsAction)
	case constants.IsActionInGroup(action, "Event"):
		return d.app.EventDomain().Authorize(ctx, cpsAction)
	case constants.IsActionInGroup(action, "Notification"):
		return d.app.NotificationService().Authorize(ctx, cpsAction)
	case constants.IsActionInGroup(action, "ProductCode"):
		return d.app.ProductCodeService().Authorize(ctx, cpsAction)
	case constants.IsActionInGroup(action, "BPSUser"):
		return d.app.BPSUserDomain().Authorize(ctx, cpsAction)
	case constants.IsActionInGroup(action, "Avatar"):
		return d.app.AvatarDomain().Authorize(ctx, cpsAction)
	case constants.IsActionInGroup(action, "Donation"):
		return d.app.DonationDomain().Authorize(ctx, cpsAction)
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
