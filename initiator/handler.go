package initiator

import (
	// Inbound section
	accesslistsegmentation "github.com/abelshebela/cbe-super-app-cps-action/internal/constants/interfaces/access_list_segmentation"
	accountBlockHandlerInterface "github.com/abelshebela/cbe-super-app-cps-action/internal/constants/interfaces/account_block"
	ap_interface "github.com/abelshebela/cbe-super-app-cps-action/internal/constants/interfaces/account_product"
	apc_interface "github.com/abelshebela/cbe-super-app-cps-action/internal/constants/interfaces/account_product_category"
	account_sub_type_interface "github.com/abelshebela/cbe-super-app-cps-action/internal/constants/interfaces/account_sub_type"
	accountvalidationInterface "github.com/abelshebela/cbe-super-app-cps-action/internal/constants/interfaces/account_validation"
	advertHandlerInterface "github.com/abelshebela/cbe-super-app-cps-action/internal/constants/interfaces/ad"
	amountBasedInbound "github.com/abelshebela/cbe-super-app-cps-action/internal/constants/interfaces/amount_based_auth"
	avatarHandlerInterface "github.com/abelshebela/cbe-super-app-cps-action/internal/constants/interfaces/avatar"
	"github.com/abelshebela/cbe-super-app-cps-action/internal/constants/interfaces/bank"
	bankvaultInterface "github.com/abelshebela/cbe-super-app-cps-action/internal/constants/interfaces/bankvault"
	bpsActionInbound "github.com/abelshebela/cbe-super-app-cps-action/internal/constants/interfaces/bps_action"
	bpsInbound "github.com/abelshebela/cbe-super-app-cps-action/internal/constants/interfaces/bps_user"
	budgetCategory "github.com/abelshebela/cbe-super-app-cps-action/internal/constants/interfaces/budget_category"
	bulk_service_inbound "github.com/abelshebela/cbe-super-app-cps-action/internal/constants/interfaces/bulk_service"
	actionInbound "github.com/abelshebela/cbe-super-app-cps-action/internal/constants/interfaces/cps_action"
	cpsUserInbound "github.com/abelshebela/cbe-super-app-cps-action/internal/constants/interfaces/cps_user"
	customerInbound "github.com/abelshebela/cbe-super-app-cps-action/internal/constants/interfaces/customer"
	cg_iface "github.com/abelshebela/cbe-super-app-cps-action/internal/constants/interfaces/customer_group"
	customer_seg "github.com/abelshebela/cbe-super-app-cps-action/internal/constants/interfaces/customer_segmentation"
	dviface "github.com/abelshebela/cbe-super-app-cps-action/internal/constants/interfaces/device_version"
	ecommerce_merchant "github.com/abelshebela/cbe-super-app-cps-action/internal/constants/interfaces/ecommerce_merchant"
	eventInbound "github.com/abelshebela/cbe-super-app-cps-action/internal/constants/interfaces/event"
	event_merchant_port "github.com/abelshebela/cbe-super-app-cps-action/internal/constants/interfaces/event_merchant"
	FaydaInbound "github.com/abelshebela/cbe-super-app-cps-action/internal/constants/interfaces/fayda"
	job_role_interface "github.com/abelshebela/cbe-super-app-cps-action/internal/constants/interfaces/job_role"
	logistics_merchant_adaptor "github.com/abelshebela/cbe-super-app-cps-action/internal/constants/interfaces/logistics_merchant"
	role_delegation_outbound "github.com/abelshebela/cbe-super-app-cps-action/internal/constants/interfaces/role_delegation"
	sar_iface "github.com/abelshebela/cbe-super-app-cps-action/internal/constants/interfaces/superapp_role"
	tac_interface "github.com/abelshebela/cbe-super-app-cps-action/internal/constants/interfaces/term_and_condition"
	"github.com/abelshebela/cbe-super-app-cps-action/internal/constants/interfaces/transaction"
	ussdMerchantInbound "github.com/abelshebela/cbe-super-app-cps-action/internal/constants/interfaces/ussd_merchant"
	role_delegation_handler "github.com/abelshebela/cbe-super-app-cps-action/internal/handlers/rest/http/role_delegation"
	"github.com/abelshebela/cbe-super-app-cps-action/internal/handlers/rest/http/ussd_merchant"

	activeState "github.com/abelshebela/cbe-super-app-cps-action/internal/constants/interfaces/active_state"
	cpsRoleInbound "github.com/abelshebela/cbe-super-app-cps-action/internal/constants/interfaces/cps_roles"
	feedbackinterface "github.com/abelshebela/cbe-super-app-cps-action/internal/constants/interfaces/feedback"
	hqInbound "github.com/abelshebela/cbe-super-app-cps-action/internal/constants/interfaces/hq"
	kycInbound "github.com/abelshebela/cbe-super-app-cps-action/internal/constants/interfaces/kyc_verifier"
	newscategory_adaptor "github.com/abelshebela/cbe-super-app-cps-action/internal/constants/interfaces/news_category"
	newstag_adaptor "github.com/abelshebela/cbe-super-app-cps-action/internal/constants/interfaces/news_tag"
	notificationInbound "github.com/abelshebela/cbe-super-app-cps-action/internal/constants/interfaces/notification"
	passwordInbound "github.com/abelshebela/cbe-super-app-cps-action/internal/constants/interfaces/password_rule"
	permissionInbound "github.com/abelshebela/cbe-super-app-cps-action/internal/constants/interfaces/permission"
	portalCardInterface "github.com/abelshebela/cbe-super-app-cps-action/internal/constants/interfaces/portal_card"
	roleInbound "github.com/abelshebela/cbe-super-app-cps-action/internal/constants/interfaces/roles"
	servicesInbound "github.com/abelshebela/cbe-super-app-cps-action/internal/constants/interfaces/services"
	TopupInbound "github.com/abelshebela/cbe-super-app-cps-action/internal/constants/interfaces/topup"
	unlinkInbound "github.com/abelshebela/cbe-super-app-cps-action/internal/constants/interfaces/unlink"
	utilityInbound "github.com/abelshebela/cbe-super-app-cps-action/internal/constants/interfaces/utility"
	survey_sampling_interface "github.com/abelshebela/cbe-super-app-cps-action/internal/constants/interfaces/survey_sampling"
	vaultCategory "github.com/abelshebela/cbe-super-app-cps-action/internal/constants/interfaces/vault"
	walletInbound "github.com/abelshebela/cbe-super-app-cps-action/internal/constants/interfaces/wallet"
	utility_handler "github.com/abelshebela/cbe-super-app-cps-action/internal/handlers/rest/http/utility"
	survey_sampling_handler "github.com/abelshebela/cbe-super-app-cps-action/internal/handlers/rest/http/survey_sampling"
	"github.com/abelshebela/cbe-super-app-cps-action/internal/service"
	queue "github.com/abelshebela/cbe-super-app-cps-action/internal/storage/queue_system"

	// Handler section
	actionrole_iface "github.com/abelshebela/cbe-super-app-cps-action/internal/constants/interfaces/action_role"
	cps_actionrole_iface "github.com/abelshebela/cbe-super-app-cps-action/internal/constants/interfaces/cps_action_role"
	accesslistsegmentaion "github.com/abelshebela/cbe-super-app-cps-action/internal/handlers/rest/http/access_list_segmentaion"
	cps_actionrole_handler "github.com/abelshebela/cbe-super-app-cps-action/internal/handlers/rest/http/cps_action_role"
	customer_group_handler "github.com/abelshebela/cbe-super-app-cps-action/internal/handlers/rest/http/customer_group"
	customer_hand "github.com/abelshebela/cbe-super-app-cps-action/internal/handlers/rest/http/customer_segmentation"
	event_merchant_handler "github.com/abelshebela/cbe-super-app-cps-action/internal/handlers/rest/http/event_merchant"
	logistics_merchant_handler "github.com/abelshebela/cbe-super-app-cps-action/internal/handlers/rest/http/logistics_merchant"
	"github.com/abelshebela/cbe-super-app-cps-action/internal/handlers/rest/http/roles"
	superapp_role_handler "github.com/abelshebela/cbe-super-app-cps-action/internal/handlers/rest/http/superapp_role"

	customerKycInbound "github.com/abelshebela/cbe-super-app-cps-action/internal/constants/interfaces/customer_kyc"
	selfActivationKYCInbound "github.com/abelshebela/cbe-super-app-cps-action/internal/constants/interfaces/customer_kyc"
	donation "github.com/abelshebela/cbe-super-app-cps-action/internal/constants/interfaces/donation"
	donation_category "github.com/abelshebela/cbe-super-app-cps-action/internal/constants/interfaces/donation_category"
	donation_company "github.com/abelshebela/cbe-super-app-cps-action/internal/constants/interfaces/donation_company"
	encryptionInbound "github.com/abelshebela/cbe-super-app-cps-action/internal/constants/interfaces/encryption"
	sitotaInbound "github.com/abelshebela/cbe-super-app-cps-action/internal/constants/interfaces/sitota"
	accountBlockHandler "github.com/abelshebela/cbe-super-app-cps-action/internal/handlers/rest/http/account_block"
	ap_handler "github.com/abelshebela/cbe-super-app-cps-action/internal/handlers/rest/http/account_product"
	apc_handler "github.com/abelshebela/cbe-super-app-cps-action/internal/handlers/rest/http/account_product_category"
	account_sub_type_handler "github.com/abelshebela/cbe-super-app-cps-action/internal/handlers/rest/http/account_sub_type"
	accountValidation "github.com/abelshebela/cbe-super-app-cps-action/internal/handlers/rest/http/account_validation"
	advertHandlerImpl "github.com/abelshebela/cbe-super-app-cps-action/internal/handlers/rest/http/ad"
	amountBasedAuthHandler "github.com/abelshebela/cbe-super-app-cps-action/internal/handlers/rest/http/amount_based_auth"
	avatarHandlerImpl "github.com/abelshebela/cbe-super-app-cps-action/internal/handlers/rest/http/avatar"
	bankHandler "github.com/abelshebela/cbe-super-app-cps-action/internal/handlers/rest/http/bank"
	bankvaulthandler "github.com/abelshebela/cbe-super-app-cps-action/internal/handlers/rest/http/bankvault"
	bpsActionHandler "github.com/abelshebela/cbe-super-app-cps-action/internal/handlers/rest/http/bps_action_handler"
	actionrole_handler "github.com/abelshebela/cbe-super-app-cps-action/internal/handlers/rest/http/bps_action_role"
	bpsHandler "github.com/abelshebela/cbe-super-app-cps-action/internal/handlers/rest/http/bps_user"
	budgetCategoryHandler "github.com/abelshebela/cbe-super-app-cps-action/internal/handlers/rest/http/budget_category"
	bulkServiceHandler "github.com/abelshebela/cbe-super-app-cps-action/internal/handlers/rest/http/bulk_service"
	tac_handler "github.com/abelshebela/cbe-super-app-cps-action/internal/handlers/rest/http/term_and_condition"

	cpsactionhandler "github.com/abelshebela/cbe-super-app-cps-action/internal/handlers/rest/http/cps_action_handler"
	cpsRoleHandler "github.com/abelshebela/cbe-super-app-cps-action/internal/handlers/rest/http/cps_roles"
	cpsUserHandler "github.com/abelshebela/cbe-super-app-cps-action/internal/handlers/rest/http/cps_user"
	CustomerHandler "github.com/abelshebela/cbe-super-app-cps-action/internal/handlers/rest/http/customer"
	CustomerKYCHandler "github.com/abelshebela/cbe-super-app-cps-action/internal/handlers/rest/http/customer_kyc"
	selfActivationKYCHandler "github.com/abelshebela/cbe-super-app-cps-action/internal/handlers/rest/http/customer_kyc"
	"github.com/abelshebela/cbe-super-app-cps-action/internal/handlers/rest/http/department"
	deviceversionhandler "github.com/abelshebela/cbe-super-app-cps-action/internal/handlers/rest/http/device_version"
	donationHandler "github.com/abelshebela/cbe-super-app-cps-action/internal/handlers/rest/http/donation"
	donationCategoryHandler "github.com/abelshebela/cbe-super-app-cps-action/internal/handlers/rest/http/donation_category"
	donationCompanyHandler "github.com/abelshebela/cbe-super-app-cps-action/internal/handlers/rest/http/donation_company"
	ecommerce_handler "github.com/abelshebela/cbe-super-app-cps-action/internal/handlers/rest/http/ecommerce-merchant"

	// encryptionHandler "github.com/abelshebela/cbe-super-app-cps-action/internal/handlers/rest/http/encryption"
	eventhandler "github.com/abelshebela/cbe-super-app-cps-action/internal/handlers/rest/http/event"
	faydaHandler "github.com/abelshebela/cbe-super-app-cps-action/internal/handlers/rest/http/fayda"
	feedbackhandler "github.com/abelshebela/cbe-super-app-cps-action/internal/handlers/rest/http/feedback"
	hqHandler "github.com/abelshebela/cbe-super-app-cps-action/internal/handlers/rest/http/hq"
	jobRoleHandler "github.com/abelshebela/cbe-super-app-cps-action/internal/handlers/rest/http/job_roles"
	kyc_handler "github.com/abelshebela/cbe-super-app-cps-action/internal/handlers/rest/http/kyc_verifier"
	newscategory_handler "github.com/abelshebela/cbe-super-app-cps-action/internal/handlers/rest/http/news_category"
	newstag_handler "github.com/abelshebela/cbe-super-app-cps-action/internal/handlers/rest/http/news_tag"
	notificationHandler "github.com/abelshebela/cbe-super-app-cps-action/internal/handlers/rest/http/notifications"
	passwordHandler "github.com/abelshebela/cbe-super-app-cps-action/internal/handlers/rest/http/password_rule"
	permissionHandler "github.com/abelshebela/cbe-super-app-cps-action/internal/handlers/rest/http/permission"
	portalcard "github.com/abelshebela/cbe-super-app-cps-action/internal/handlers/rest/http/portal_card"
	services_http "github.com/abelshebela/cbe-super-app-cps-action/internal/handlers/rest/http/services"
	sitotaHandler "github.com/abelshebela/cbe-super-app-cps-action/internal/handlers/rest/http/sitota"
	TopupHandler "github.com/abelshebela/cbe-super-app-cps-action/internal/handlers/rest/http/topup"
	transaction_handler "github.com/abelshebela/cbe-super-app-cps-action/internal/handlers/rest/http/transaction"
	unlinkHandler "github.com/abelshebela/cbe-super-app-cps-action/internal/handlers/rest/http/unlink"
	vaultcategoryhandler "github.com/abelshebela/cbe-super-app-cps-action/internal/handlers/rest/http/vault"
	walletHandler "github.com/abelshebela/cbe-super-app-cps-action/internal/handlers/rest/http/wallet"

	ActiveStateHandler "github.com/abelshebela/cbe-super-app-cps-action/internal/handlers/rest/http/active_state"

	"gitlab.com/bersufekadgetachew/cbe-super-app-shared/shared/utils"
)

type Handler struct {
	ActiveStateHandler            activeState.ActiveState
	RoleHandler                   roleInbound.RolesInbound
	jobRoleHandler                job_role_interface.RolesInbound
	CpsActionHandler              actionInbound.CPSActionAdapter
	BpsActionHandler              bpsActionInbound.BPSActionAdapter
	UnlinkHandler                 unlinkInbound.UnlinkAdapter
	EventHandler                  eventInbound.EventAdapter
	WalletHandler                 walletInbound.WalletAdapter
	TopupHandler                  TopupInbound.TopupAdapter
	PasswordHandler               passwordInbound.PasswordRule
	BpsHandler                    bpsInbound.BPSUserHandler
	AccountSubTypeHandler         account_sub_type_interface.AccountSubTypeHandler
	BankHandler                   bank.BankHandler
	FeedbackHandler               feedbackinterface.FeedbackAdapter
	BudgetCategoryHandler         budgetCategory.BudgetCategoryPortHandler
	AdvertHandler                 advertHandlerInterface.ADAdapter
	PortalCardHander              portalCardInterface.PortalCardAdapter
	AccountValidation             accountvalidationInterface.AccountValidation
	HqHandler                     hqInbound.HQAdapter
	AccountBlockHandler           accountBlockHandlerInterface.AccountBlockAdapter
	ServicesHandler               servicesInbound.ServicesHandler
	AmountBasedAuthHandler        amountBasedInbound.AmountBasedAuthAdapter
	DepartmentHandler             department.DepartmentHandler
	AvatarHandler                 avatarHandlerInterface.AvatarInbound
	FaydaHandler                  FaydaInbound.FaydaAccount
	bulkServiceHandler            bulk_service_inbound.BulkServiceHandler
	customerHandler               customerInbound.CustomerDetail
	Permission                    permissionInbound.PermissionHandler
	CPSUser                       cpsUserInbound.CPSUserHandler
	DonationHandler               donation.DonationHandler
	DonationCategoryHandler       donation_category.DonationCategoryAdapter
	DonationCompanyHandler        donation_company.DonationCompanyAdapter
	NotificationHandler           notificationInbound.NotificationHandler
	BankVaultHandler              bankvaultInterface.BankVaultHandler
	VaultCategoryHandler          vaultCategory.VaultCategoryHandler
	NewsCategoryHandler           newscategory_adaptor.NewsCategoryAdaptor
	NewsTagHandler                newstag_adaptor.NewsTagAdaptor
	SitotaHandler                 sitotaInbound.SitotaAdapter
	KYCVerifierHandler            kycInbound.KYCVerifierAdapter
	EncryptionHandler             encryptionInbound.EncryptionAdapter
	DeviceVersionHandler          dviface.DeviceVersionHandler
	BPSActionRoleHandler          actionrole_iface.BPSActionRoleHandler
	CPSActionRoleHandler          cps_actionrole_iface.CPSActionRoleHandler
	TransactionHandler            transaction.TransactionInterface
	EventMerchantHandler          event_merchant_port.EventMerchantInboundAdaptor
	LogisticsMerchantHandler      logistics_merchant_adaptor.LogisticMerchantInboundAdaptor
	EcommerceMerchantHandler      ecommerce_merchant.EcommerceMerchant
	AccessLostSegmentationHandler accesslistsegmentation.AccessListSegmentationHandler
	CustomerSegmentationHandler   customer_seg.CustomerSegmentation
	CustomerGroupHandler          cg_iface.CustomerGroup
	CPSRolesHandler               cpsRoleInbound.CPSRolesAdapter
	CustomerKYCHandler            customerKycInbound.CustomerKYC
	SelfActivationKYCHandler      selfActivationKYCInbound.SelfActivationKyc
	SuperAppRoleHandler           sar_iface.SuperAppRole
	RoleDelegationHandler         role_delegation_outbound.RoleDelegation
	AccountProductCategoryHandler apc_interface.AccountProductCategoryHandler
	AccountProductHandler         ap_interface.AccountProductHandler
	TermAndConditionHandler       tac_interface.TermAndConditionHandler
	UtilityHandler                utilityInbound.UtilityInbound
	UssdMerchantHandler           ussdMerchantInbound.UssdMerchantInbound
	SurveySamplingHandler         survey_sampling_interface.SurveySamplingInbound
}

func InitHandler(serviceLayer service.ServiceLayer, logger utils.Logger, queueManager *queue.QueueManager) Handler {
	return Handler{
		RoleHandler:             roles.NewRoleHandler(serviceLayer.RoleService, logger),
		jobRoleHandler:          jobRoleHandler.NewJobRoleHandler(serviceLayer.JobRoleService, logger),
		BudgetCategoryHandler:   budgetCategoryHandler.InitBudgetCategoryAdapter(serviceLayer.BudgetCategory, logger),
		UnlinkHandler:           unlinkHandler.InitUnlinkAdapter(serviceLayer.Unlink, logger),
		BpsHandler:              bpsHandler.InitBPSUserMakerHandler(serviceLayer.BpsUser, logger),
		AccountSubTypeHandler:   account_sub_type_handler.InitAccountSubTypeAdapter(serviceLayer.AccountSubType, logger),
		BankHandler:             bankHandler.InitBankAdapter(serviceLayer.Bank, logger),
		CpsActionHandler:        cpsactionhandler.InitCPSActionAdapter(serviceLayer.CPSAction, logger),
		BpsActionHandler:        bpsActionHandler.InitBPSActionAdapter(serviceLayer.BPSActionService, logger),
		EventHandler:            eventhandler.InitEventAdapter(serviceLayer.EventService, logger),
		FeedbackHandler:         feedbackhandler.InitFeedbackAdapter(serviceLayer.Feedback, logger),
		PortalCardHander:        portalcard.InitPortalCardAdapter(serviceLayer.PortalCard, logger),
		AdvertHandler:           advertHandlerImpl.InitAdvertAdapter(serviceLayer.Advert, logger),
		AvatarHandler:           avatarHandlerImpl.InitAvatarAdapter(serviceLayer.Avatar, logger),
		WalletHandler:           walletHandler.InitWalletAdapter(serviceLayer.Wallet, logger),
		TopupHandler:            TopupHandler.InitTopupAdapter(serviceLayer.Topup, logger),
		PasswordHandler:         passwordHandler.InitPasswordRuleHandler(serviceLayer.PasswordRule, logger),
		AccountValidation:       accountValidation.NewHttpAccountValidation(serviceLayer.AccountValidation, logger),
		AccountBlockHandler:     accountBlockHandler.InitAccountBlockAdapter(serviceLayer.AccountBlock, logger),
		DepartmentHandler:       department.NewDepartmentHandler(serviceLayer.Department, logger),
		HqHandler:               hqHandler.InitHQAdapter(serviceLayer.HQService, logger),
		bulkServiceHandler:      bulkServiceHandler.InitBulkServiceAdapter(serviceLayer.BulkService, logger),
		customerHandler:         CustomerHandler.InitCustomerAdapter(serviceLayer.CustomerService, logger),
		FaydaHandler:            faydaHandler.InitFaydaHandler(serviceLayer.Fayda, logger),
		Permission:              permissionHandler.InitPermissionHandler(serviceLayer.Permission, logger),
		CPSUser:                 cpsUserHandler.InitCPSUserHandler(serviceLayer.CPSUser, logger),
		NotificationHandler:     notificationHandler.InitNotificationHandler(serviceLayer.NotificationService, logger),
		BankVaultHandler:        bankvaulthandler.InitBankVaultHandler(serviceLayer.BankVault, logger),
		VaultCategoryHandler:    vaultcategoryhandler.InitVaultCategoryHandler(serviceLayer.VaultCategoryService, logger),
		AmountBasedAuthHandler:  amountBasedAuthHandler.NewAmountBasedAuthHandler(serviceLayer.AmountBasedAuth, logger),
		ServicesHandler:         services_http.InitServicesAdapter(serviceLayer.Services, logger),
		DonationHandler:         donationHandler.NewDonationAdapter(serviceLayer.Donation, logger),
		DonationCategoryHandler: donationCategoryHandler.InitDonationCategoryAdapter(serviceLayer.DonationCategory, logger),
		DonationCompanyHandler:  donationCompanyHandler.InitDonationCompanyAdapter(serviceLayer.DonationCompany, logger),
		NewsCategoryHandler:     newscategory_handler.NewNewsCategoryHandler(serviceLayer.NewsCategoryService, logger),
		NewsTagHandler:          newstag_handler.NewNewsTagHandler(serviceLayer.NewsTagService, logger),
		SitotaHandler:           sitotaHandler.InitSitotaHandler(serviceLayer.Sitota, logger),
		KYCVerifierHandler:      kyc_handler.InitKYCAdapter(serviceLayer.KYCVerifier, logger),
		// EncryptionHandler:        encryptionHandler.InitEncryption(serviceLayer.Encryption, logger),
		DeviceVersionHandler:     deviceversionhandler.InitDeviceVersionAdapter(serviceLayer.DeviceVersion, logger),
		BPSActionRoleHandler:     actionrole_handler.NewBPSActionRoleHandler(serviceLayer.BPSActionRole, logger),
		CPSActionRoleHandler:     cps_actionrole_handler.NewCPSActionRoleHandler(serviceLayer.CPSActionRole, logger),
		TransactionHandler:       transaction_handler.NewTransactionHandler(serviceLayer.TransactionService, logger),
		EventMerchantHandler:     event_merchant_handler.NewEventMerchantHandler(serviceLayer.EventMerchantService, logger),
		LogisticsMerchantHandler: logistics_merchant_handler.NewLogisticsMerchantHandler(serviceLayer.LogisticsMerchantService, logger),
		EcommerceMerchantHandler: ecommerce_handler.NewEcommerceMerchantdapter(serviceLayer.EcommerceMerchant, logger),

		AccessLostSegmentationHandler: accesslistsegmentaion.InitAccessListSegmentationAdapter(serviceLayer.AccessListSegmentationService, logger),
		CustomerSegmentationHandler:   customer_hand.NewCustomerSegmentation(serviceLayer.CustomerSegmentation, logger),
		CustomerGroupHandler:          customer_group_handler.NewCustomerGroupHandler(serviceLayer.CustomerGroup, logger),
		CPSRolesHandler:               cpsRoleHandler.NewCPSRolesHandler(serviceLayer.CPSRoles, logger),
		CustomerKYCHandler:            CustomerKYCHandler.NewCustomerKYCAdapter(serviceLayer.CustomerKYC, logger),
		SelfActivationKYCHandler:      selfActivationKYCHandler.NewSelfActivationAdapter(serviceLayer.SelfActivateKYCService, logger),
		UssdMerchantHandler:           ussd_merchant.NewUssdMerchantHandler(serviceLayer.UssdMerchantService, logger),
		SuperAppRoleHandler:           superapp_role_handler.NewSuperAppRoleHandler(serviceLayer.SuperAppRole, logger),
		RoleDelegationHandler:         role_delegation_handler.NewRoleDelegationHandler(serviceLayer.RoleDelegationService, logger),
		AccountProductCategoryHandler: apc_handler.InitAccountProductCategoryAdapter(serviceLayer.AccountProductCategory, logger),
		AccountProductHandler:         ap_handler.InitAccountProductAdapter(serviceLayer.AccountProduct, logger),
		TermAndConditionHandler:       tac_handler.InitTermAndConditionAdapter(serviceLayer.AccountOpeningTerms, logger),
		UtilityHandler:                utility_handler.NewUtilityAdapter(serviceLayer.Utility, logger),
		SurveySamplingHandler:         survey_sampling_handler.NewSurveySamplingHandler(serviceLayer.SurveySampling, logger),
		ActiveStateHandler:            ActiveStateHandler.NewActiveState(),
	}
}
