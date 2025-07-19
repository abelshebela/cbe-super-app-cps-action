package application

import (
	"github.com/CBE-Super-App/cbe-super-app-cps-action/internal/domain/account_block"
	"github.com/CBE-Super-App/cbe-super-app-cps-action/internal/domain/account_validation"
	"github.com/CBE-Super-App/cbe-super-app-cps-action/internal/domain/action"
	ad_service "github.com/CBE-Super-App/cbe-super-app-cps-action/internal/domain/ad/service"
	amount_based_auth_domain "github.com/CBE-Super-App/cbe-super-app-cps-action/internal/domain/amount_based_auth"
	avatar_domain "github.com/CBE-Super-App/cbe-super-app-cps-action/internal/domain/avatar"
	bank_service "github.com/CBE-Super-App/cbe-super-app-cps-action/internal/domain/bank/service"
	budget_service "github.com/CBE-Super-App/cbe-super-app-cps-action/internal/domain/budget"
	"github.com/CBE-Super-App/cbe-super-app-cps-action/internal/domain/cps_user/services"
	customer_service "github.com/CBE-Super-App/cbe-super-app-cps-action/internal/domain/customer/service"
	department "github.com/CBE-Super-App/cbe-super-app-cps-action/internal/domain/department"
	fayda_service "github.com/CBE-Super-App/cbe-super-app-cps-action/internal/domain/fayda_account/service"
	feedback_service "github.com/CBE-Super-App/cbe-super-app-cps-action/internal/domain/feedback"
	hq "github.com/CBE-Super-App/cbe-super-app-cps-action/internal/domain/hq"
	miniApp_domain "github.com/CBE-Super-App/cbe-super-app-cps-action/internal/domain/miniapp"
	password_service "github.com/CBE-Super-App/cbe-super-app-cps-action/internal/domain/password_rule/services"
	permission "github.com/CBE-Super-App/cbe-super-app-cps-action/internal/domain/permission"
	portalcard "github.com/CBE-Super-App/cbe-super-app-cps-action/internal/domain/portal_card"
	service "github.com/CBE-Super-App/cbe-super-app-cps-action/internal/domain/service"

	budget_category "github.com/CBE-Super-App/cbe-super-app-cps-action/internal/domain/budget_category"
	cps_action_service "github.com/CBE-Super-App/cbe-super-app-cps-action/internal/domain/cps_actions/services"
	event_domain "github.com/CBE-Super-App/cbe-super-app-cps-action/internal/domain/event"
	unlink_service "github.com/CBE-Super-App/cbe-super-app-cps-action/internal/domain/unlink"
	wallet_service "github.com/CBE-Super-App/cbe-super-app-cps-action/internal/domain/wallet/service"
)

type Domain struct {
	AdDomain              ad_service.AdvertService
	AvatarDomian          avatar_domain.AvatarDomainService
	BankDomain            bank_service.BankService
	WalletDomain          wallet_service.WalletService
	FaydaDomain           *fayda_service.FaydaAccountDomain
	CustomerDomain        *customer_service.CustomerDomain
	FeedbackDomain        *feedback_service.FeedbackDomain
	UnlinkDomain          *unlink_service.Service
	BudgetDomain          *budget_service.BudgetService
	AccountDomain         account_validation.Service
	CPSUserDomain         services.CPSUserService
	PasswordRuleDomain    password_service.PasswordRuleService
	PortalCardDomain      portalcard.PortaCardInterface
	DepartmentDomain      department.Service
	HQDomain              hq.Service
	MiniAppDomain         miniApp_domain.MiniAppService
	AccountBlockDomain    account_block.ApplicationServices
	PermissionDomain      permission.Service
	AmountBasedAuthDomain amount_based_auth_domain.AmountBasedAuthRepository
	ServiceDomain         service.ServiceInterface
	ActionDomain          action.ServiceInterface
	EventDomain           event_domain.EventService
	BudgetCategoryDomain  budget_category.BudgetCategoryService
	CPSActionDomain       cps_action_service.CPSActionService
}
