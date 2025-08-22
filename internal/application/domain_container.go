package application

import (
	"github.com/CBE-Super-App/cbe-super-app-cps-action/internal/domain/account_block"
	account_service "github.com/CBE-Super-App/cbe-super-app-cps-action/internal/domain/account_lookup"
	"github.com/CBE-Super-App/cbe-super-app-cps-action/internal/domain/account_validation"
	"github.com/CBE-Super-App/cbe-super-app-cps-action/internal/domain/action"
	ad_service "github.com/CBE-Super-App/cbe-super-app-cps-action/internal/domain/ad"
	amount_based_auth_domain "github.com/CBE-Super-App/cbe-super-app-cps-action/internal/domain/amount_based_auth"
	avatar_domain "github.com/CBE-Super-App/cbe-super-app-cps-action/internal/domain/avatar"
	bank_service "github.com/CBE-Super-App/cbe-super-app-cps-action/internal/domain/bank/service"
	bps_user_service "github.com/CBE-Super-App/cbe-super-app-cps-action/internal/domain/bps_user"
	budget_service "github.com/CBE-Super-App/cbe-super-app-cps-action/internal/domain/budget"
	budget_category "github.com/CBE-Super-App/cbe-super-app-cps-action/internal/domain/budget_category"
	bulk_service "github.com/CBE-Super-App/cbe-super-app-cps-action/internal/domain/bulk_service"
	cps_action_service "github.com/CBE-Super-App/cbe-super-app-cps-action/internal/domain/cps_actions/services"
	"github.com/CBE-Super-App/cbe-super-app-cps-action/internal/domain/cps_user/services"
	customer_service "github.com/CBE-Super-App/cbe-super-app-cps-action/internal/domain/customer/service"
	department "github.com/CBE-Super-App/cbe-super-app-cps-action/internal/domain/department"
	event_domain "github.com/CBE-Super-App/cbe-super-app-cps-action/internal/domain/event"
	fayda_service "github.com/CBE-Super-App/cbe-super-app-cps-action/internal/domain/fayda_account/service"
	feedback_service "github.com/CBE-Super-App/cbe-super-app-cps-action/internal/domain/feedback"
	hq "github.com/CBE-Super-App/cbe-super-app-cps-action/internal/domain/hq"
	miniApp_domain "github.com/CBE-Super-App/cbe-super-app-cps-action/internal/domain/miniapp"
	mini_app_merchant_service "github.com/CBE-Super-App/cbe-super-app-cps-action/internal/domain/miniapp_merchant"
	notification_domain "github.com/CBE-Super-App/cbe-super-app-cps-action/internal/domain/notification"
	password_service "github.com/CBE-Super-App/cbe-super-app-cps-action/internal/domain/password_rule/services"
	permission "github.com/CBE-Super-App/cbe-super-app-cps-action/internal/domain/permission"
	portalcard "github.com/CBE-Super-App/cbe-super-app-cps-action/internal/domain/portal_card"
	service_domain "github.com/CBE-Super-App/cbe-super-app-cps-action/internal/domain/service"
	unlink_service "github.com/CBE-Super-App/cbe-super-app-cps-action/internal/domain/unlink"
	wallet_service "github.com/CBE-Super-App/cbe-super-app-cps-action/internal/domain/wallet"

	// notification_domain "github.com/CBE-Super-App/cbe-super-app-cps-action/internal/domain/notification"
	donation_domain "github.com/CBE-Super-App/cbe-super-app-cps-action/internal/domain/donation"
	productcode "github.com/CBE-Super-App/cbe-super-app-cps-action/internal/domain/product_code"

	keyGen_service "github.com/CBE-Super-App/cbe-super-app-cps-action/pkgs/keygen"
)

type Domain struct {
	AccountBlockDomain    account_block.ApplicationServices
	AccountDomain         account_validation.Service
	ActionDomain          action.ServiceInterface
	AdDomain              ad_service.AdvertService
	AmountBasedAuthDomain amount_based_auth_domain.AmountBasedAuthDomain
	AvatarDomian          avatar_domain.AvatarDomainService
	BankDomain            bank_service.BankService
	BPSUserDomain         bps_user_service.Service
	BudgetCategoryDomain  budget_category.BudgetCategoryService
	BudgetDomain          *budget_service.BudgetService
	CPSActionDomain       cps_action_service.CPSActionService
	CPSUserDomain         services.CPSUserService
	CustomerDomain        customer_service.CustomerService
	DepartmentDomain      department.Service
	EventDomain           event_domain.EventService // fully not ready
	FaydaDomain           fayda_service.FaydaAccount
	FeedbackDomain        *feedback_service.FeedbackDomain
	HQDomain              hq.Service
	MiniAppDomain         miniApp_domain.MiniAppService // not done
	PasswordRuleDomain    password_service.PasswordRuleService
	PermissionDomain      permission.Service
	PortalCardDomain      portalcard.PortaCardInterface
	UnlinkDomain          unlink_service.UnlinkAccount
	WalletDomain          wallet_service.WalletService
	MiniAppMerchantDomain mini_app_merchant_service.MiniAppMerchantService
	AccountLookup         *account_service.UserSearchService
	BulkServiceDomain     bulk_service.BulkService
	ServiceCheckDomain    service_domain.ServiceRepo
	KeyGenService         keyGen_service.KeyGeneratorService
	NotificationService   notification_domain.NotificationService
	ProductCodeService    productcode.Service
	DonationDomain        donation_domain.DonationService
}
