package factory

import (
	"cbe-super-app-cps-action/internal/service"
	"cbe-super-app-cps-action/internal/storage/persistance"

	"gitlab.com/bersufekadgetachew/cbe-super-app-shared/shared/utils"
)

type ServiceFactory struct {
	persistence persistance.Persistence
	logger      utils.Logger
}

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
		AccountBlockContainer: f.createAccountBlockService(),
		AccountContainer:      f.createAccountValidationService(),
		AccountLookup:         f.createAccountSearchService(),

		// Action services
		ActionContainer: f.createActionService(),

		// Advertisement services
		AdContainer: f.createAdvertService(),

		// Authentication services
		AmountBasedAuthContainer: f.createAmountBasedAuthService(),
		AvatarDomian:             f.createAvatarService(),

		// Banking services
		BankContainer:    f.createBankService(),
		BPSUserContainer: f.createBPSUserService(),

		// Budget services
		BudgetCategoryContainer: f.createBudgetCategoryService(),
		BudgetContainer:         f.createBudgetService(),

		// Bulk services
		BulkServiceContainer: f.createBulkService(),

		// CPS services
		CPSActionContainer: f.createCPSActionService(),
		CPSUserContainer:   f.createCPSUserService(),

		// Customer services
		CustomerContainer: f.createCustomerService(),

		// Department services
		DepartmentContainer: f.createDepartmentService(),

		// Donation services
		DonationContainer: f.createDonationService(),

		// Event services
		EventContainer: f.createEventService(),

		// Fayda services
		FaydaContainer: f.createFaydaService(),

		// Feedback services
		FeedbackContainer: f.createFeedbackService(),

		// HQ services
		HQContainer: f.createHQService(),

		// Key generation services
		KeyGenService: f.createKeyGeneratorService(),

		// Mini app services
		MiniAppContainer:         f.createMiniAppService(),
		MiniAppMerchantContainer: f.createMiniAppMerchantService(),

		// Notification services
		NotificationService: f.createNotificationService(),

		// Password services
		PasswordRuleContainer: f.createPasswordRuleService(),

		// Permission services
		PermissionContainer: f.createPermissionService(),

		// Portal services
		PortalCardContainer: f.createPortalCardService(),

		// Product services
		ProductCodeService: f.createProductCodeService(),

		// Service services
		ServiceCheckContainer: f.createServiceService(),

		// Unlink services
		UnlinkContainer: f.createUnlinkService(),

		// Wallet services
		WalletContainer: f.createWalletService(),
	}
}

// Individual service creation methods
func (f *ServiceFactory) createAccountBlockService() service.AccountBlockService {
	return &accountBlockService{
		repo:   f.persistence.AccountBlockPersistence,
		logger: f.logger,
	}
}

func (f *ServiceFactory) createAccountValidationService() service.AccountValidationService {
	return &accountValidationService{
		repo:   f.persistence.UserPersistence,
		logger: f.logger,
	}
}

func (f *ServiceFactory) createAccountSearchService() service.AccountSearchService {
	return &accountSearchService{
		repo:   f.persistence.AdvertPersistence,
		logger: f.logger,
	}
}

func (f *ServiceFactory) createActionService() service.ActionService {
	return &actionService{
		logger: f.logger,
	}
}

func (f *ServiceFactory) createAdvertService() service.AdvertService {
	return &advertService{
		repo:   f.persistence.AdvertRepositoryPersistence,
		logger: f.logger,
	}
}

func (f *ServiceFactory) createAmountBasedAuthService() service.AmountBasedAuthService {
	return &amountBasedAuthService{
		repo:   f.persistence.AmountBasedAuthPersistence,
		logger: f.logger,
	}
}

func (f *ServiceFactory) createAvatarService() service.AvatarService {
	return &avatarService{
		repo:   f.persistence.AvatarPersistence,
		logger: f.logger,
	}
}

func (f *ServiceFactory) createBankService() service.BankService {
	return &bankService{
		repo:   f.persistence.BankPersistence,
		logger: f.logger,
	}
}

func (f *ServiceFactory) createBPSUserService() service.BPSUserService {
	return &bpsUserService{
		repo:   f.persistence.BPSUserPersistence,
		logger: f.logger,
	}
}

func (f *ServiceFactory) createBudgetCategoryService() service.BudgetCategoryService {
	return &budgetCategoryService{
		repo:   f.persistence.AmountBasedAuthPersistence,
		logger: f.logger,
	}
}

func (f *ServiceFactory) createBudgetService() service.BudgetService {
	return &budgetService{
		repo:   f.persistence.AmountBasedAuthPersistence,
		logger: f.logger,
	}
}

func (f *ServiceFactory) createBulkService() service.BulkService {
	return &bulkService{
		repo:   f.persistence.AccountBlockPersistence,
		logger: f.logger,
	}
}

func (f *ServiceFactory) createCPSActionService() service.CPSActionService {
	// This should be created by the existing CPS action service
	// For now, return a placeholder that will be replaced in the initiator
	return &cpsActionService{
		logger: f.logger,
	}
}

func (f *ServiceFactory) createCPSUserService() service.CPSUserService {
	return &cpsUserService{
		repo:   f.persistence.CpsUserPersistence,
		logger: f.logger,
	}
}

func (f *ServiceFactory) createCustomerService() service.CustomerService {
	return &customerService{
		repo:   f.persistence.UserPersistence,
		logger: f.logger,
	}
}

func (f *ServiceFactory) createDepartmentService() service.DepartmentService {
	return &departmentService{
		repo:   f.persistence.AccessListPersistence,
		logger: f.logger,
	}
}

func (f *ServiceFactory) createDonationService() service.DonationService {
	return &donationService{
		repo:   f.persistence.DonationPersistence,
		logger: f.logger,
	}
}

func (f *ServiceFactory) createEventService() service.EventService {
	return &eventService{
		repo:   f.persistence.EventPersistence,
		logger: f.logger,
	}
}

func (f *ServiceFactory) createFaydaService() service.FaydaAccountService {
	return &faydaService{
		repo:   f.persistence.LinkedAccountPersistence,
		logger: f.logger,
	}
}

func (f *ServiceFactory) createFeedbackService() service.FeedbackService {
	return &feedbackService{
		repo:   f.persistence.FeedbackPersistence,
		logger: f.logger,
	}
}

func (f *ServiceFactory) createHQService() service.HQService {
	return &hqService{
		repo:   f.persistence.HQPersistence,
		logger: f.logger,
	}
}

func (f *ServiceFactory) createKeyGeneratorService() service.KeyGeneratorService {
	return &keyGeneratorService{
		logger: f.logger,
	}
}

func (f *ServiceFactory) createMiniAppService() service.MiniAppService {
	return &miniAppService{
		repo:   f.persistence.MiniAppPersistence,
		logger: f.logger,
	}
}

func (f *ServiceFactory) createMiniAppMerchantService() service.MiniAppMerchantService {
	return &miniAppMerchantService{
		repo:   f.persistence.MiniAppMerchantPersistence,
		logger: f.logger,
	}
}

func (f *ServiceFactory) createNotificationService() service.NotificationService {
	return &notificationService{
		repo:   f.persistence.NotificationPersistence,
		logger: f.logger,
	}
}

func (f *ServiceFactory) createPasswordRuleService() service.PasswordRuleService {
	return &passwordRuleService{
		repo:   f.persistence.PasswordRulePersistence,
		logger: f.logger,
	}
}

func (f *ServiceFactory) createPermissionService() service.PermissionService {
	return &permissionService{
		logger: f.logger,
	}
}

func (f *ServiceFactory) createPortalCardService() service.PortalCardService {
	return &portalCardService{
		repo:   f.persistence.PortalCardPersistence,
		logger: f.logger,
	}
}

func (f *ServiceFactory) createProductCodeService() service.ProductCodeService {
	return &productCodeService{
		logger: f.logger,
	}
}

func (f *ServiceFactory) createServiceService() service.ServiceService {
	return &serviceService{
		repo:   f.persistence.ServiceDetailsPersistence,
		logger: f.logger,
	}
}

func (f *ServiceFactory) createUnlinkService() service.UnlinkService {
	return &unlinkService{
		logger: f.logger,
	}
}

func (f *ServiceFactory) createWalletService() service.WalletService {
	return &walletService{
		repo:   f.persistence.WalletPersistence,
		logger: f.logger,
	}
}
