// Package initiator provides adapters for initializing inbound HTTP handlers and services for the CBE Super App CPS Action module.
package initiator

import (
	accountblock_handler "github.com/CBE-Super-App/cbe-super-app-cps-action/internal/adapter/inbound/http/account_block"
	accountvalidation_inbound "github.com/CBE-Super-App/cbe-super-app-cps-action/internal/adapter/inbound/http/account_validation"
	"github.com/CBE-Super-App/cbe-super-app-cps-action/internal/adapter/inbound/http/ad"
	//bps_user_handler "github.com/CBE-Super-App/cbe-super-app-cps-action/internal/adapter/inbound/http/bps_user"
	"github.com/CBE-Super-App/cbe-super-app-cps-action/internal/adapter/inbound/http/permission_handler"

	"github.com/CBE-Super-App/cbe-super-app-cps-action/internal/adapter/inbound/http/bank"
	"github.com/CBE-Super-App/cbe-super-app-cps-action/internal/adapter/inbound/http/budget_handler"
	bulkservices_inbound "github.com/CBE-Super-App/cbe-super-app-cps-action/internal/adapter/inbound/http/bulk_service"
	cpsmakerhandler "github.com/CBE-Super-App/cbe-super-app-cps-action/internal/adapter/inbound/http/cps-user"
	customerhandler "github.com/CBE-Super-App/cbe-super-app-cps-action/internal/adapter/inbound/http/customer_handler"

	"github.com/CBE-Super-App/cbe-super-app-cps-action/internal/adapter/inbound/http/department_handler"
	faydaaccount "github.com/CBE-Super-App/cbe-super-app-cps-action/internal/adapter/inbound/http/fayda_account"
	feedbackhandler "github.com/CBE-Super-App/cbe-super-app-cps-action/internal/adapter/inbound/http/feedback_handler"
	passwordrule "github.com/CBE-Super-App/cbe-super-app-cps-action/internal/adapter/inbound/http/password_rule"

	avatar_adapter "github.com/CBE-Super-App/cbe-super-app-cps-action/internal/adapter/inbound/http/avatar"
	portalcard "github.com/CBE-Super-App/cbe-super-app-cps-action/internal/adapter/inbound/http/portal_card"

	// "github.com/CBE-Super-App/cbe-super-app-cps-action/internal/adapter/inbound/http/unlink_device_handler"
	"github.com/CBE-Super-App/cbe-super-app-cps-action/internal/adapter/inbound/http/wallet"
	"github.com/CBE-Super-App/cbe-super-app-cps-action/internal/port/inbound"
	inboudService "github.com/CBE-Super-App/cbe-super-app-cps-action/internal/port/inbound"
	accountblock "github.com/CBE-Super-App/cbe-super-app-cps-action/internal/port/inbound/account_block"
	inboundLookUp "github.com/CBE-Super-App/cbe-super-app-cps-action/internal/port/inbound/account_lookup"
	inboundAccount "github.com/CBE-Super-App/cbe-super-app-cps-action/internal/port/inbound/account_validation"
	inboundAD "github.com/CBE-Super-App/cbe-super-app-cps-action/internal/port/inbound/ad"
	inboundAvatar "github.com/CBE-Super-App/cbe-super-app-cps-action/internal/port/inbound/avatar"
	inboundBank "github.com/CBE-Super-App/cbe-super-app-cps-action/internal/port/inbound/bank"
	inboundBudget "github.com/CBE-Super-App/cbe-super-app-cps-action/internal/port/inbound/budget"
	inboundBulkServices "github.com/CBE-Super-App/cbe-super-app-cps-action/internal/port/inbound/bulk_services"
	inboundDepartment "github.com/CBE-Super-App/cbe-super-app-cps-action/internal/port/inbound/department"
	inboundFeedback "github.com/CBE-Super-App/cbe-super-app-cps-action/internal/port/inbound/feedback"
	inboundMiniAppMerchant "github.com/CBE-Super-App/cbe-super-app-cps-action/internal/port/inbound/mini_app_merchant"

	inboundMiniApp "github.com/CBE-Super-App/cbe-super-app-cps-action/internal/port/inbound/miniapp"

	inboundPermission "github.com/CBE-Super-App/cbe-super-app-cps-action/internal/port/inbound/permission"
	inboundWallet "github.com/CBE-Super-App/cbe-super-app-cps-action/internal/port/inbound/wallet"

	amountBasedAuth "github.com/CBE-Super-App/cbe-super-app-cps-action/internal/adapter/inbound/http/amount_based_auth_handler"
	eventhandler "github.com/CBE-Super-App/cbe-super-app-cps-action/internal/adapter/inbound/http/event"
	hq "github.com/CBE-Super-App/cbe-super-app-cps-action/internal/adapter/inbound/http/hq"
	AmountBasedAuthRepo "github.com/CBE-Super-App/cbe-super-app-cps-action/internal/port/inbound/amount_based_auth"
	bpsRepo "github.com/CBE-Super-App/cbe-super-app-cps-action/internal/port/inbound/bps_user"
	event_inbound "github.com/CBE-Super-App/cbe-super-app-cps-action/internal/port/inbound/event"
	notification_inbound "github.com/CBE-Super-App/cbe-super-app-cps-action/internal/port/inbound/notification"
	productcode_inbound "github.com/CBE-Super-App/cbe-super-app-cps-action/internal/port/inbound/product_code"

	// service_inbound "github.com/CBE-Super-App/cbe-super-app-cps-action/internal/port/inbound/service_details"
	"github.com/CBE-Super-App/cbe-super-app-cps-action/internal/port/inbound/service_details"
	"gitlab.com/bersufekadgetachew/cbe-super-app-shared/shared/config"
	"gitlab.com/bersufekadgetachew/cbe-super-app-shared/shared/utils"

	account_handler "github.com/CBE-Super-App/cbe-super-app-cps-action/internal/adapter/inbound/http/account_lookup"
	budget_category "github.com/CBE-Super-App/cbe-super-app-cps-action/internal/adapter/inbound/http/budget_category"
	cps_actions_handler "github.com/CBE-Super-App/cbe-super-app-cps-action/internal/adapter/inbound/http/cps_action"
	miniapp_handler "github.com/CBE-Super-App/cbe-super-app-cps-action/internal/adapter/inbound/http/mini_app_handler"
	miniapp_merchant_handler "github.com/CBE-Super-App/cbe-super-app-cps-action/internal/adapter/inbound/http/mini_app_merchant"
	notification_handler "github.com/CBE-Super-App/cbe-super-app-cps-action/internal/adapter/inbound/http/notification"
	productcode_handler "github.com/CBE-Super-App/cbe-super-app-cps-action/internal/adapter/inbound/http/product_code"
	service_handler "github.com/CBE-Super-App/cbe-super-app-cps-action/internal/adapter/inbound/http/service"
	unlink_handler "github.com/CBE-Super-App/cbe-super-app-cps-action/internal/adapter/inbound/http/unlink"
	bulk_service_handler "github.com/CBE-Super-App/cbe-super-app-cps-action/internal/adapter/inbound/http/updated_bulk_service"

	service_inbound "github.com/CBE-Super-App/cbe-super-app-cps-action/internal/port/inbound"
	unlink_inbound "github.com/CBE-Super-App/cbe-super-app-cps-action/internal/port/inbound/unlink"

	bps_user_handler "github.com/CBE-Super-App/cbe-super-app-cps-action/internal/adapter/inbound/http/bps_user"

	// budget_category "github.com/CBE-Super-App/cbe-super-app-cps-action/internal/port/inbound/budget_category"
	inboundBudgetCategory "github.com/CBE-Super-App/cbe-super-app-cps-action/internal/port/inbound/budget_category"
	cps_actions "github.com/CBE-Super-App/cbe-super-app-cps-action/internal/port/inbound/cps_actions"
	bulk_service "github.com/CBE-Super-App/cbe-super-app-cps-action/internal/port/inbound/updated_bulk_service"
	//bps_user "github.com/CBE-Super-App/cbe-super-app-cps-action/internal/port/inbound/bps_user"
	//
)

type Adapter struct {
	AvatarAdapter          inboundAvatar.AvatarInbound
	BankAdapter            inboundBank.BankAdapter
	AdAdapter              inboundAD.ADAdapter
	AmountBasedAuthAdapter inbound.AmountBasedAuthHandler
	WalletAdapter          inboundWallet.WalletAdapter
	FaydaAdapter           inbound.FaydaAccount
	BranchAdapter          inbound.BranchHandler
	CustomerAdapter        inbound.CustomerDetail
	FeedbackAdapter        inboundFeedback.Feedback
	DepartmentAdapter      inboundDepartment.DepartmentPortHandler
	PermissionAdapter      inboundPermission.PermissionPortHandler
	UnlinkAdapter          unlink_inbound.UnlinkHandler
	BudgetAdapter          inboundBudget.BudgetPortHandler
	AccountAdapter         inboundAccount.Inbound
	BulkServiceAdapter     inboundBulkServices.Inbound
	CPSUserAdapter         inbound.CPSUserMakerHandler

	BPSUserAdapter bpsRepo.BPSUserHandler

	PasswordRuleAdapter  inbound.PasswordRuleInbound
	PortalCardAdapter    inbound.PortalCardBound
	ServiceDetailAdapter service_details.ServiceDetailsInbound
	AccountBlockAdapter  accountblock.AccountBlockHandler
	ServiceAdapter       inboudService.ServiceBound
	HQAdapter            *hq.HQHTTPHandler
	AmountBasedAuth      AmountBasedAuthRepo.AmountBasedAuthHandler
	MiniAppAdapter       inboundMiniApp.MiniAppInbound
	EventAdapter         event_inbound.EventHandler
	// BudgetCategoryAdapter  budget_category.BudgetCategoryInbound
	CPSActionAdapter          cps_actions.CPSActionAdapter
	BudgetCategoryAdapter     inboundBudgetCategory.BudgetCategoryInbound
	MiniAppMerchantAdapter    inboundMiniAppMerchant.MiniAppMerchantInbound
	AccountLookUp             inboundLookUp.UserSearchAdapter
	UpdatedBulkServiceAdapter bulk_service.BulkServiceHandler

	ServiceCheckAdapter service_inbound.Service
	NotificationAdapter notification_inbound.NotificationHandler
	ProductCodeAdapter  productcode_inbound.ProductCodeHandler
}

func InitAdapter(application Application, minioClient config.MinioClientInterface, logger utils.Logger) Adapter {
	return Adapter{
		AvatarAdapter:       avatar_adapter.InitAvatarHTTPHandler(application.AvatarApplication, logger),
		BankAdapter:         bank.InitBankAdapter(application.BankApplication, logger),
		AdAdapter:           ad.NewAdvertHTTPHandler(application.AdApplication, logger),
		WalletAdapter:       wallet.InitWalletRouter(application.WalletApplication, logger),
		FaydaAdapter:        faydaaccount.InitFaydaAdapter(application.FaydaApplication, logger),
		CustomerAdapter:     customerhandler.NewCustomerHTTPHandler(application.CustomerApplication, logger),
		FeedbackAdapter:     feedbackhandler.NewFeedbackHTTPHandler(application.FeedbackApplication, logger),
		UnlinkAdapter:       unlink_handler.InitAdapterUnlinkService(application.UnlinkApplication, logger),
		BudgetAdapter:       budget_handler.NewBudgetHTTPHandler(application.BudgetApplication, logger),
		AccountAdapter:      accountvalidation_inbound.NewHttpAccountValidation(application.AccountApplication, logger),
		BulkServiceAdapter:  bulkservices_inbound.NewHttpBulkService(application.BulkServicesApplication, logger),
		CPSUserAdapter:      cpsmakerhandler.InitCPSUserMakerHandler(application.CPSUserApplication, logger),
		PasswordRuleAdapter: passwordrule.NewPasswordRuleHTTPHandler(application.PasswordRuleApplication),
		PortalCardAdapter:   portalcard.NewportalCardHandler(application.PortalCardApplication, logger),
		DepartmentAdapter:   department_handler.NewDepartmentHTTPHandler(application.DepartmentApplication, logger),
		PermissionAdapter:   permission_handler.NewPermissionHTTPHandler(application.PermissionApplication, logger),
		AccountBlockAdapter: accountblock_handler.NewAccountBlockHandler(application.AccountBlockApplication, logger),
		HQAdapter:           hq.NewHQHTTPHandler(application.HQApplication),
		AmountBasedAuth:     amountBasedAuth.NewAmountBasedAuthHandler(application.AmountBasedAuthApplication, logger),
		MiniAppAdapter:      miniapp_handler.NewMiniAppAdapter(application.MiniAppApplication, logger),
		EventAdapter:        eventhandler.NewEventHTTPHandler(application.EventApplication, logger),

		BPSUserAdapter: bps_user_handler.InitBPSUserMakerHandler(application.BPSUserApplication, logger),

		// BudgetCategoryAdapter: budget_category_handler.InitBudgetCategoryAdapter(application.BudgetCategoryApplication, logger),
		UpdatedBulkServiceAdapter: bulk_service_handler.InitBulkServiceHandler(application.BulkServiceApplication, logger),
		CPSActionAdapter:          cps_actions_handler.InitCPSActionAdapter(application.CPSActionApplication, logger),
		BudgetCategoryAdapter:     budget_category.InitBudgetCategoryAdapter(application.BudgetCategoryApplication, logger, application.FileService),
		MiniAppMerchantAdapter:    miniapp_merchant_handler.NewMiniAppMerchantAdapter(application.MiniAppMerchantApplication, logger),
		AccountLookUp:             account_handler.NewMiniAppMerchantAdapter(application.AccountLookupApplication, logger),
		ServiceCheckAdapter:       service_handler.NewServiceHandler(application.ServicCheckeApplication, logger),
		NotificationAdapter:       notification_handler.NewNotificationHTTPHandler(application.NotificationApplication, logger),
		ProductCodeAdapter:        productcode_handler.NewProductCodeHTTPHandler(application.ProductCodeApplication, logger),
	}
}
