package initiator

import (
	// Inbound section
	accountBlockHandlerInterface "cbe-super-app-cps-action/internal/constants/interfaces/account_block"
	accountvalidationInterface "cbe-super-app-cps-action/internal/constants/interfaces/account_validation"
	advertHandlerInterface "cbe-super-app-cps-action/internal/constants/interfaces/ad"
	amountBasedInbound "cbe-super-app-cps-action/internal/constants/interfaces/amount_based_auth"
	avatarHandlerInterface "cbe-super-app-cps-action/internal/constants/interfaces/avatar"
	"cbe-super-app-cps-action/internal/constants/interfaces/bank"
	bankvaultInterface "cbe-super-app-cps-action/internal/constants/interfaces/bankvault"
	bpsInbound "cbe-super-app-cps-action/internal/constants/interfaces/bps_user"
	budgetCategory "cbe-super-app-cps-action/internal/constants/interfaces/budget_category"
	bulk_service_inbound "cbe-super-app-cps-action/internal/constants/interfaces/bulk_service"
	actionInbound "cbe-super-app-cps-action/internal/constants/interfaces/cps_action"
	cpsUserInbound "cbe-super-app-cps-action/internal/constants/interfaces/cps_user"
	customerInbound "cbe-super-app-cps-action/internal/constants/interfaces/customer"
	dviface "cbe-super-app-cps-action/internal/constants/interfaces/device_version"
	ecommerce_merchant "cbe-super-app-cps-action/internal/constants/interfaces/ecommerce_merchant"
	eventInbound "cbe-super-app-cps-action/internal/constants/interfaces/event"
	event_merchant_port "cbe-super-app-cps-action/internal/constants/interfaces/event_merchant"
	FaydaInbound "cbe-super-app-cps-action/internal/constants/interfaces/fayda"
	"cbe-super-app-cps-action/internal/constants/interfaces/transaction"
	vaultAmountTierInbound "cbe-super-app-cps-action/internal/constants/interfaces/vault_amount_tier"

	feedbackinterface "cbe-super-app-cps-action/internal/constants/interfaces/feedback"
	hqInbound "cbe-super-app-cps-action/internal/constants/interfaces/hq"
	kycInbound "cbe-super-app-cps-action/internal/constants/interfaces/kyc_verifier"
	newscategory_adaptor "cbe-super-app-cps-action/internal/constants/interfaces/news_category"
	newstag_adaptor "cbe-super-app-cps-action/internal/constants/interfaces/news_tag"
	notificationInbound "cbe-super-app-cps-action/internal/constants/interfaces/notification"
	passwordInbound "cbe-super-app-cps-action/internal/constants/interfaces/password_rule"
	permissionInbound "cbe-super-app-cps-action/internal/constants/interfaces/permission"
	portalCardInterface "cbe-super-app-cps-action/internal/constants/interfaces/portal_card"
	service_details "cbe-super-app-cps-action/internal/constants/interfaces/service_details"
	servicesInbound "cbe-super-app-cps-action/internal/constants/interfaces/services"
	TopupInbound "cbe-super-app-cps-action/internal/constants/interfaces/topup"
	unlinkInbound "cbe-super-app-cps-action/internal/constants/interfaces/unlink"
	vaultgroupcategory "cbe-super-app-cps-action/internal/constants/interfaces/vaultgroup_category"
	walletInbound "cbe-super-app-cps-action/internal/constants/interfaces/wallet"
	"cbe-super-app-cps-action/internal/service"

	// Handler section
	actionrole_iface "cbe-super-app-cps-action/internal/constants/interfaces/action_role"
	cps_actionrole_iface "cbe-super-app-cps-action/internal/constants/interfaces/cps_action_role"
	cps_actionrole_handler "cbe-super-app-cps-action/internal/handlers/rest/http/cps_action_role"
	event_merchant_handler "cbe-super-app-cps-action/internal/handlers/rest/http/event_merchant"

	donation "cbe-super-app-cps-action/internal/constants/interfaces/donation"
	donation_category "cbe-super-app-cps-action/internal/constants/interfaces/donation_category"
	donation_company "cbe-super-app-cps-action/internal/constants/interfaces/donation_company"
	encryptionInbound "cbe-super-app-cps-action/internal/constants/interfaces/encryption"
	sitotaInbound "cbe-super-app-cps-action/internal/constants/interfaces/sitota"
	accountBlockHandler "cbe-super-app-cps-action/internal/handlers/rest/http/account_block"
	accountValidation "cbe-super-app-cps-action/internal/handlers/rest/http/account_validation"
	advertHandlerImpl "cbe-super-app-cps-action/internal/handlers/rest/http/ad"
	amountBasedAuthHandler "cbe-super-app-cps-action/internal/handlers/rest/http/amount_based_auth"
	avatarHandlerImpl "cbe-super-app-cps-action/internal/handlers/rest/http/avatar"
	bankHandler "cbe-super-app-cps-action/internal/handlers/rest/http/bank"
	bankvaulthandler "cbe-super-app-cps-action/internal/handlers/rest/http/bankvault"
	actionrole_handler "cbe-super-app-cps-action/internal/handlers/rest/http/bps_action_role"
	bpsHandler "cbe-super-app-cps-action/internal/handlers/rest/http/bps_user"
	budgetCategoryHandler "cbe-super-app-cps-action/internal/handlers/rest/http/budget_category"
	bulkServiceHandler "cbe-super-app-cps-action/internal/handlers/rest/http/bulk_service"
	cpsactionhandler "cbe-super-app-cps-action/internal/handlers/rest/http/cps_action_handler"
	cpsUserHandler "cbe-super-app-cps-action/internal/handlers/rest/http/cps_user"
	CustomerHandler "cbe-super-app-cps-action/internal/handlers/rest/http/customer"
	"cbe-super-app-cps-action/internal/handlers/rest/http/department"
	deviceversionhandler "cbe-super-app-cps-action/internal/handlers/rest/http/device_version"
	donationHandler "cbe-super-app-cps-action/internal/handlers/rest/http/donation"
	donationCategoryHandler "cbe-super-app-cps-action/internal/handlers/rest/http/donation_category"
	donationCompanyHandler "cbe-super-app-cps-action/internal/handlers/rest/http/donation_company"
	ecommerce_handler "cbe-super-app-cps-action/internal/handlers/rest/http/ecommerce-merchant"
	encryptionHandler "cbe-super-app-cps-action/internal/handlers/rest/http/encryption"
	eventhandler "cbe-super-app-cps-action/internal/handlers/rest/http/event"
	faydaHandler "cbe-super-app-cps-action/internal/handlers/rest/http/fayda"
	feedbackhandler "cbe-super-app-cps-action/internal/handlers/rest/http/feedback"
	hqHandler "cbe-super-app-cps-action/internal/handlers/rest/http/hq"
	kyc_handler "cbe-super-app-cps-action/internal/handlers/rest/http/kyc_verifier"
	newscategory_handler "cbe-super-app-cps-action/internal/handlers/rest/http/news_category"
	newstag_handler "cbe-super-app-cps-action/internal/handlers/rest/http/news_tag"
	notificationHandler "cbe-super-app-cps-action/internal/handlers/rest/http/notifications"
	passwordHandler "cbe-super-app-cps-action/internal/handlers/rest/http/password_rule"
	permissionHandler "cbe-super-app-cps-action/internal/handlers/rest/http/permission"
	portalcard "cbe-super-app-cps-action/internal/handlers/rest/http/portal_card"
	productCodeHandler "cbe-super-app-cps-action/internal/handlers/rest/http/product_code"
	serviceDetailsHandler "cbe-super-app-cps-action/internal/handlers/rest/http/service_details"
	services_http "cbe-super-app-cps-action/internal/handlers/rest/http/services"
	sitotaHandler "cbe-super-app-cps-action/internal/handlers/rest/http/sitota"
	TopupHandler "cbe-super-app-cps-action/internal/handlers/rest/http/topup"
	transaction_handler "cbe-super-app-cps-action/internal/handlers/rest/http/transaction"
	unlinkHandler "cbe-super-app-cps-action/internal/handlers/rest/http/unlink"
	amount_tier_handler "cbe-super-app-cps-action/internal/handlers/rest/http/vault_amount_tier"
	vaultgroupcategoryhandler "cbe-super-app-cps-action/internal/handlers/rest/http/vaultgroup_category"
	walletHandler "cbe-super-app-cps-action/internal/handlers/rest/http/wallet"

	"gitlab.com/bersufekadgetachew/cbe-super-app-shared/shared/utils"
)

type Handler struct {
	CpsActionHandler          actionInbound.CPSActionAdapter
	UnlinkHandler             unlinkInbound.UnlinkAdapter
	EventHandler              eventInbound.EventAdapter
	WalletHandler             walletInbound.WalletAdapter
	TopupHandler              TopupInbound.TopupAdapter
	PasswordHandler           passwordInbound.PasswordRule
	BpsHandler                bpsInbound.BPSUserHandler
	BankHandler               bank.BankHandler
	FeedbackHandler           feedbackinterface.FeedbackAdapter
	BudgetCategoryHandler     budgetCategory.BudgetCategoryPortHandler
	AdvertHandler             advertHandlerInterface.ADAdapter
	PortalCardHander          portalCardInterface.PortalCardAdapter
	AccountValidation         accountvalidationInterface.AccountValidation
	HqHandler                 hqInbound.HQAdapter
	AccountBlockHandler       accountBlockHandlerInterface.AccountBlockAdapter
	ServiceDetailsHandler     service_details.ServiceAdapter
	ServicesHandler           servicesInbound.ServicesHandler
	AmountBasedAuthHandler    amountBasedInbound.AmountBasedAuthAdapter
	DepartmentHandler         department.DepartmentHandler
	AvatarHandler             avatarHandlerInterface.AvatarInbound
	FaydaHandler              FaydaInbound.FaydaAccount
	bulkServiceHandler        bulk_service_inbound.BulkServiceHandler
	customerHandler           customerInbound.CustomerDetail
	Permission                permissionInbound.PermissionHandler
	CPSUser                   cpsUserInbound.CPSUserHandler
	ProductCodeHandler        productCodeHandler.ProductCodeAdapter
	DonationHandler           donation.DonationHandler
	DonationCategoryHandler   donation_category.DonationCategoryAdapter
	DonationCompanyHandler    donation_company.DonationCompanyAdapter
	NotificationHandler       notificationInbound.NotificationHandler
	BankVaultHandler          bankvaultInterface.BankVaultHandler
	VaultGroupCategoryHandler vaultgroupcategory.VaultGroupCategoryHandler
	NewsCategoryHandler       newscategory_adaptor.NewsCategoryAdaptor
	NewsTagHandler            newstag_adaptor.NewsTagAdaptor
	SitotaHandler             sitotaInbound.SitotaAdapter
	KYCVerifierHandler        kycInbound.KYCVerifierAdapter
	EncryptionHandler         encryptionInbound.EncryptionAdapter
	DeviceVersionHandler      dviface.DeviceVersionHandler
	BPSActionRoleHandler      actionrole_iface.BPSActionRoleHandler
	CPSActionRoleHandler      cps_actionrole_iface.CPSActionRoleHandler
	TransactionHandler        transaction.TransactionInterface
	EventMerchantHandler      event_merchant_port.EventMerchantInboundAdaptor
	AmountTierHandler         vaultAmountTierInbound.VaultAmountTierHandler
	EcommerceMerchantHandler  ecommerce_merchant.EcommerceMerchant
}

func InitHandler(serviceLayer service.ServiceLayer, logger utils.Logger) Handler {
	pcs := serviceLayer.ProductCode
	return Handler{

		BudgetCategoryHandler:     budgetCategoryHandler.InitBudgetCategoryAdapter(serviceLayer.BudgetCategory, logger),
		UnlinkHandler:             unlinkHandler.InitUnlinkAdapter(serviceLayer.Unlink, logger),
		BpsHandler:                bpsHandler.InitBPSUserMakerHandler(serviceLayer.BpsUser, logger),
		BankHandler:               bankHandler.InitBankAdapter(serviceLayer.Bank, logger),
		CpsActionHandler:          cpsactionhandler.InitCPSActionAdapter(serviceLayer.CPSAction, logger),
		EventHandler:              eventhandler.InitEventAdapter(serviceLayer.EventService, logger),
		FeedbackHandler:           feedbackhandler.InitFeedbackAdapter(serviceLayer.Feedback, logger),
		PortalCardHander:          portalcard.InitPortalCardAdapter(serviceLayer.PortalCard, logger),
		AdvertHandler:             advertHandlerImpl.InitAdvertAdapter(serviceLayer.Advert, logger),
		AvatarHandler:             avatarHandlerImpl.InitAvatarAdapter(serviceLayer.Avatar, logger),
		WalletHandler:             walletHandler.InitWalletAdapter(serviceLayer.Wallet, logger),
		TopupHandler:              TopupHandler.InitTopupAdapter(serviceLayer.Topup, logger),
		PasswordHandler:           passwordHandler.InitPasswordRuleHandler(serviceLayer.PasswordRule, logger),
		AccountValidation:         accountValidation.NewHttpAccountValidation(serviceLayer.AccountValidation, logger),
		AccountBlockHandler:       accountBlockHandler.InitAccountBlockAdapter(serviceLayer.AccountBlock, logger),
		DepartmentHandler:         department.NewDepartmentHandler(serviceLayer.Department, logger),
		HqHandler:                 hqHandler.InitHQAdapter(serviceLayer.HQService, logger),
		bulkServiceHandler:        bulkServiceHandler.InitBulkServiceAdapter(serviceLayer.BulkService, logger),
		customerHandler:           CustomerHandler.InitCustomerAdapter(serviceLayer.CustomerService, logger),
		FaydaHandler:              faydaHandler.InitFaydaHandler(serviceLayer.Fayda, logger),
		Permission:                permissionHandler.InitPermissionHandler(serviceLayer.Permission, logger),
		CPSUser:                   cpsUserHandler.InitCPSUserHandler(serviceLayer.CPSUser, logger),
		NotificationHandler:       notificationHandler.InitNotificationHandler(serviceLayer.NotificationService, logger),
		BankVaultHandler:          bankvaulthandler.InitBankVaultHandler(serviceLayer.BankVault, logger),
		VaultGroupCategoryHandler: vaultgroupcategoryhandler.InitVaultGroupCategoryHandler(serviceLayer.VaultGroupCategory, logger),
		AmountBasedAuthHandler:    amountBasedAuthHandler.NewAmountBasedAuthHandler(serviceLayer.AmountBasedAuth, logger),
		ServiceDetailsHandler:     serviceDetailsHandler.InitServiceAdapter(serviceLayer.ServiceDetails, logger),
		ServicesHandler:           services_http.InitServicesAdapter(serviceLayer.Services, logger),
		ProductCodeHandler:        productCodeHandler.InitProductcodeAdapter(pcs, logger),
		DonationHandler:           donationHandler.NewDonationAdapter(serviceLayer.Donation, logger),
		DonationCategoryHandler:   donationCategoryHandler.InitDonationCategoryAdapter(serviceLayer.DonationCategory, logger),
		DonationCompanyHandler:    donationCompanyHandler.InitDonationCompanyAdapter(serviceLayer.DonationCompany, logger),
		NewsCategoryHandler:       newscategory_handler.NewNewsCategoryHandler(serviceLayer.NewsCategoryService, logger),
		NewsTagHandler:            newstag_handler.NewNewsTagHandler(serviceLayer.NewsTagService, logger),
		SitotaHandler:             sitotaHandler.InitSitotaHandler(serviceLayer.Sitota, logger),
		KYCVerifierHandler:        kyc_handler.InitKYCAdapter(serviceLayer.KYCVerifier, logger),
		EncryptionHandler:         encryptionHandler.InitEncryption(serviceLayer.Encryption, logger),
		DeviceVersionHandler:      deviceversionhandler.InitDeviceVersionAdapter(serviceLayer.DeviceVersion, logger),
		BPSActionRoleHandler:      actionrole_handler.NewBPSActionRoleHandler(serviceLayer.BPSActionRole, logger),
		CPSActionRoleHandler:      cps_actionrole_handler.NewCPSActionRoleHandler(serviceLayer.CPSActionRole, logger),
		TransactionHandler:        transaction_handler.NewTransactionHandler(serviceLayer.TransactionService, logger),
		EventMerchantHandler:      event_merchant_handler.NewEventMerchantHandler(serviceLayer.EventMerchantService, logger),
		AmountTierHandler:         amount_tier_handler.NewVaultAmountTierHandler(serviceLayer.VaultAmountTierService, logger),
		EcommerceMerchantHandler:  ecommerce_handler.NewEcommerceMerchantdapter(serviceLayer.EcommerceMerchant, logger),
	}
}
