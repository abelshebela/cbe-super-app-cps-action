package factory

import (
	"cbe-super-app-cps-action/internal/service"
	// cpsaction "cbe-super-app-cps-action/internal/service/cpsaction"
	"cbe-super-app-cps-action/internal/storage/persistance"

	"gitlab.com/bersufekadgetachew/cbe-super-app-shared/shared/utils"
	// "cbe-super-app-cps-action/internal/service/account_block"
	// "cbe-super-app-cps-action/internal/service/account_search"
	// "cbe-super-app-cps-action/internal/service/account_validation"
	// "cbe-super-app-cps-action/internal/service/action"
	// "cbe-super-app-cps-action/internal/service/advert"
	// "cbe-super-app-cps-action/internal/service/amount_based_auth"
	// "cbe-super-app-cps-action/internal/service/avatar"
	// "cbe-super-app-cps-action/internal/service/bank"
	// "cbe-super-app-cps-action/internal/service/bps_user"
	// "cbe-super-app-cps-action/internal/service/budget"
	// "cbe-super-app-cps-action/internal/service/budget_category"
	// "cbe-super-app-cps-action/internal/service/bulk"
	// "cbe-super-app-cps-action/internal/service/cps_user"
	// "cbe-super-app-cps-action/internal/service/customer"
	// "cbe-super-app-cps-action/internal/service/department"
	// "cbe-super-app-cps-action/internal/service/donation"
	// "cbe-super-app-cps-action/internal/service/event"
	// "cbe-super-app-cps-action/internal/service/fayda"
	// "cbe-super-app-cps-action/internal/service/feedback"
	// "cbe-super-app-cps-action/internal/service/linked_account"
	// "cbe-super-app-cps-action/internal/service/merchant"
	// "cbe-super-app-cps-action/internal/service/mini_app_merchant"
	// "cbe-super-app-cps-action/internal/service/notification"
	// "cbe-super-app-cps-action/internal/service/password_rule"
	// "cbe-super-app-cps-action/internal/service/permission"
	// "cbe-super-app-cps-action/internal/service/portal_card"
	// "cbe-super-app-cps-action/internal/service/product_code"
	// "cbe-super-app-cps-action/internal/service/service"
	// "cbe-super-app-cps-action/internal/service/unlink"
	// "cbe-super-app-cps-action/internal/service/wallet"
)

type ServiceFactory struct {
	persistence persistance.Persistence
	logger      utils.Logger
}

// ==========================================================

// func (f *ServiceFactory) createAccountBlockService() service.AccountBlockService {
// 	return accountblock.NewAccountBlockService(f.persistence.AccountBlockPersistence, f.logger)
// }

// func (f *ServiceFactory) createAccountValidationService() service.AccountValidationService {
// 	return accountvalidation.NewAccountValidationService(f.persistence.UserPersistence, f.logger)
// }

// func (f *ServiceFactory) createAccountSearchService() service.AccountSearchService {
// 	return accountsearch.NewAccountSearchService(f.persistence.AccountAPIPort, f.logger)
// }

// func (f *ServiceFactory) createActionService() service.ActionService {
// 	return action.NewActionService(f.logger)
// }

// func (f *ServiceFactory) createAdvertService() service.AdvertService {
// 	return advert.NewAdvertService(f.persistence.AdvertPersistence, f.logger)
// }

// func (f *ServiceFactory) createAmountBasedAuthService() service.AmountBasedAuthService {
// 	return amountbasedauth.NewAmountBasedAuthService(f.persistence.AmountBasedAuthPersistence, f.logger)
// }

// func (f *ServiceFactory) createAvatarService() service.AvatarService {
// 	return avatar.NewAvatarService(f.persistence.AvatarPersistence, f.logger)
// }

// func (f *ServiceFactory) createBankService() service.BankService {
// 	return bank.NewBankService(f.persistence.BankPersistence, f.logger)
// }

// func (f *ServiceFactory) createBPSUserService() service.BPSUserService {
// 	return bpsuser.NewBPSUserService(f.persistence.BpsUserPersistence, f.logger)
// }

// func (f *ServiceFactory) createBudgetCategoryService() service.BudgetCategoryService {
// 	return budgetcategory.NewBudgetCategoryService(f.persistence.BudgetCategoryPersistence, f.logger)
// }

// func (f *ServiceFactory) createBudgetService() service.BudgetService {
// 	return budget.NewBudgetService(f.persistence.BudgetPersistence, f.logger)
// }

// func (f *ServiceFactory) createBulkService() service.BulkService {
// 	return bulk.NewBulkService(f.persistence.BulkPersistence, f.logger)
// }

// func (f *ServiceFactory) createCPSActionService() service.CPSActionService {
// 	return cpsaction.NewCPSActionService(f.persistence.CPSAction, f.persistence, f.logger)
// }

// func (f *ServiceFactory) createCPSUserService() service.CPSUserService {
// 	return cpsuser.NewCPSUserService(f.persistence.CpsUserPersistence, f.logger)
// }

// func (f *ServiceFactory) createCustomerService() service.CustomerService {
// 	return customer.NewCustomerService(f.persistence.CustomerPersistence, f.logger)
// }

// func (f *ServiceFactory) createDepartmentService() service.DepartmentService {
// 	return department.NewDepartmentService(f.persistence.DepartmentPersistence, f.logger)
// }

// func (f *ServiceFactory) createDonationService() service.DonationService {
// 	return donation.NewDonationService(f.persistence.DonationPersistence, f.logger)
// }

// func (f *ServiceFactory) createFaydaService() service.FaydaAccountService {
// 	return fayda.NewFaydaService(f.persistence.LinkedAccountPersistence, f.logger)
// }

// func (f *ServiceFactory) createFeedbackService() service.FeedbackService {
// 	return feedback.NewFeedbackService(f.persistence.FeedbackPersistence, f.logger)
// }

// func (f *ServiceFactory) createHQService() service.HQService {
// 	return service.NewHQService(f.persistence.HQPersistence, f.logger)
// }

// func (f *ServiceFactory) createKeyGeneratorService() service.KeyGeneratorService {
// 	return service.NewKeyGeneratorService(f.logger)
// }

// func (f *ServiceFactory) createMiniAppService() service.MiniAppService {
// 	return service.NewMiniAppService(f.persistence.MiniAppPersistence, f.logger)
// }

// func (f *ServiceFactory) createMiniAppMerchantService() service.MiniAppMerchantService {
// 	return miniappmerchant.NewMiniAppMerchantService(f.persistence.MiniAppMerchantPersistence, f.logger)
// }

// func (f *ServiceFactory) createNotificationService() service.NotificationService {
// 	return notification.NewNotificationService(f.persistence.NotificationPersistence, f.logger)
// }

// func (f *ServiceFactory) createPasswordRuleService() service.PasswordRuleService {
// 	return passwordrule.NewPasswordRuleService(f.persistence.PasswordRulePersistence, f.logger)
// }

// func (f *ServiceFactory) createPermissionService() service.PermissionService {
// 	return permission.NewPermissionService(f.logger)
// }

// func (f *ServiceFactory) createPortalCardService() service.PortalCardService {
// 	return portalcard.NewPortalCardService(f.persistence.PortalCardPersistence, f.logger)
// }

// func (f *ServiceFactory) createProductCodeService() service.ProductCodeService {
// 	return productcode.NewProductCodeService(f.logger)
// }

// func (f *ServiceFactory) createServiceService() service.ServiceService {
// 	return service.NewServiceService(f.persistence.ServiceDetailsPersistence, f.logger)
// }

// func (f *ServiceFactory) createUnlinkService() service.UnlinkService {
// 	return unlink.NewUnlinkService(
// 		f.persistence.MongoClient,
// 		f.persistence.UserPersistence,
// 		f.persistence.ArchivedUserPersistence,
// 		f.persistence.LinkedAccountPersistence,
// 		f.persistence.ArchivedLinkedAccountPersistence,
// 		f.persistence.UnlinkAccountPersistence,
// 		f.createCPSActionService(),
// 		f.logger,
// 	)
// }

// func (f *ServiceFactory) createWalletService() service.WalletService {
// 	return wallet.NewWalletService(f.persistence.WalletPersistence, f.logger)
// }

// ==========================================================
// NewServiceFactory creates a new service factory
func NewServiceFactory(persistence persistance.Persistence, logger utils.Logger) *ServiceFactory {
	return &ServiceFactory{
		persistence: persistence,
		logger:      logger,
	}
}

// CreateServiceContainer creates all services and returns a populated ServiceContainer
func (f *ServiceFactory) CreateServiceContainer() service.ServiceContainer {
	return service.ServiceContainer{
		// Account related services
		// AccountBlockContainer: f.createAccountBlockService(),
		// AccountContainer:      f.createAccountValidationService(),
		// AccountLookup:         f.createAccountSearchService(),

		// // Action services
		// ActionContainer: f.createActionService(),

		// // Advertisement services
		// AdContainer: f.createAdvertService(),

		// // Authentication services
		// AmountBasedAuthContainer: f.createAmountBasedAuthService(),
		// AvatarDomian:             f.createAvatarService(),

		// // Banking services
		// BankContainer:    f.createBankService(),
		// BPSUserContainer: f.createBPSUserService(),

		// // Budget services
		// BudgetCategoryContainer: f.createBudgetCategoryService(),
		// BudgetContainer:         f.createBudgetService(),

		// // Bulk services
		// BulkServiceContainer: f.createBulkService(),

		// // CPS services
		// CPSActionContainer: f.createCPSActionService(),
		// CPSUserContainer:   f.createCPSUserService(),

		// // Customer services
		// CustomerContainer: f.createCustomerService(),

		// // Department services
		// DepartmentContainer: f.createDepartmentService(),

		// // Donation services
		// DonationContainer: f.createDonationService(),

		// // Event services
		// EventContainer: f.createEventService(),

		// // Fayda services
		// FaydaContainer: f.createFaydaService(),

		// // Feedback services
		// FeedbackContainer: f.createFeedbackService(),

		// // HQ services
		// HQContainer: f.createHQService(),

		// // Key generation services
		// KeyGenService: f.createKeyGeneratorService(),

		// // Mini app services
		// MiniAppContainer:         f.createMiniAppService(),
		// MiniAppMerchantContainer: f.createMiniAppMerchantService(),

		// // Notification services
		// NotificationService: f.createNotificationService(),

		// // Password services
		// PasswordRuleContainer: f.createPasswordRuleService(),

		// // Permission services
		// PermissionContainer: f.createPermissionService(),

		// // Portal services
		// PortalCardContainer: f.createPortalCardService(),

		// // Product services
		// ProductCodeService: f.createProductCodeService(),

		// // Service services
		// ServiceCheckContainer: f.createServiceService(),

		// // Unlink services
		// UnlinkContainer: f.createUnlinkService(),

		// Wallet services
		// WalletContainer: f.createWalletService(),
	}
}
