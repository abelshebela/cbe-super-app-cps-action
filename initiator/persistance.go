package initiator

import (
	"cbe-super-app-cps-action/internal/storage/external_call/merchant_lookup"
	"cbe-super-app-cps-action/internal/storage/kafka"

	"cbe-super-app-cps-action/internal/storage/persistance"
	"cbe-super-app-cps-action/internal/storage/persistance/access_list"
	access_list_segmentation_repository "cbe-super-app-cps-action/internal/storage/persistance/access_list_segmentation"
	"cbe-super-app-cps-action/internal/storage/persistance/account_block"
	accountvalidation "cbe-super-app-cps-action/internal/storage/persistance/account_validation"
	"cbe-super-app-cps-action/internal/storage/persistance/advert"
	"cbe-super-app-cps-action/internal/storage/persistance/amount_based_auth"
	"cbe-super-app-cps-action/internal/storage/persistance/auth_tier"
	"cbe-super-app-cps-action/internal/storage/persistance/avatar"
	"cbe-super-app-cps-action/internal/storage/persistance/bank"
	actionrole_repo "cbe-super-app-cps-action/internal/storage/persistance/bps_action_role"
	"cbe-super-app-cps-action/internal/storage/persistance/bps_user"
	"cbe-super-app-cps-action/internal/storage/persistance/cps_action"
	cps_actionrole_repo "cbe-super-app-cps-action/internal/storage/persistance/cps_action_role"
	"cbe-super-app-cps-action/internal/storage/persistance/cps_user"
	customer_segmentation_repo "cbe-super-app-cps-action/internal/storage/persistance/customer_segmentaion"
	"cbe-super-app-cps-action/internal/storage/persistance/department"
	"cbe-super-app-cps-action/internal/storage/persistance/device_history"
	deviceversioncontrol "cbe-super-app-cps-action/internal/storage/persistance/device_version_control"
	"cbe-super-app-cps-action/internal/storage/persistance/donation"
	"cbe-super-app-cps-action/internal/storage/persistance/donation_category"
	"cbe-super-app-cps-action/internal/storage/persistance/event"
	event_merchant_repository "cbe-super-app-cps-action/internal/storage/persistance/event_merchant"
	"cbe-super-app-cps-action/internal/storage/persistance/icon"
	"cbe-super-app-cps-action/internal/storage/persistance/linked_account"
	logistics_merchant_repository "cbe-super-app-cps-action/internal/storage/persistance/logistics_merchant"
	"cbe-super-app-cps-action/internal/storage/persistance/media"
	newscategory_repo "cbe-super-app-cps-action/internal/storage/persistance/news_category"
	newstag_repo "cbe-super-app-cps-action/internal/storage/persistance/news_tag"
	"time"

	"cbe-super-app-cps-action/internal/storage/persistance/archived_linked_account"
	"cbe-super-app-cps-action/internal/storage/persistance/archived_user"
	"cbe-super-app-cps-action/internal/storage/persistance/budget_category"
	"cbe-super-app-cps-action/internal/storage/persistance/donation_company"
	ecommerce_merchant "cbe-super-app-cps-action/internal/storage/persistance/ecommerce_merchant"
	"cbe-super-app-cps-action/internal/storage/persistance/fayda"
	"cbe-super-app-cps-action/internal/storage/persistance/feedback"
	"cbe-super-app-cps-action/internal/storage/persistance/hq"
	kyc_repo "cbe-super-app-cps-action/internal/storage/persistance/kyc_verifier"
	"cbe-super-app-cps-action/internal/storage/persistance/mini_app"
	"cbe-super-app-cps-action/internal/storage/persistance/notification"
	"cbe-super-app-cps-action/internal/storage/persistance/otp"

	job_repo "cbe-super-app-cps-action/internal/storage/persistance/job_role"
	password "cbe-super-app-cps-action/internal/storage/persistance/password_rule"
	permission "cbe-super-app-cps-action/internal/storage/persistance/permission"
	"cbe-super-app-cps-action/internal/storage/persistance/portal_card"
	"cbe-super-app-cps-action/internal/storage/persistance/reset_session"
	role_repo "cbe-super-app-cps-action/internal/storage/persistance/role"
	services_repo "cbe-super-app-cps-action/internal/storage/persistance/services"
	Topup "cbe-super-app-cps-action/internal/storage/persistance/topup"
	"cbe-super-app-cps-action/internal/storage/persistance/users"
	"cbe-super-app-cps-action/internal/storage/persistance/wallet"

	"cbe-super-app-cps-action/internal/storage/persistance/bulk_service"
	cps_roles "cbe-super-app-cps-action/internal/storage/persistance/cps_roles"
	"cbe-super-app-cps-action/internal/storage/persistance/customer"

	"github.com/hugokessem/coreio/core"
	"gitlab.com/bersufekadgetachew/cbe-super-app-shared/shared/config"
	"gitlab.com/bersufekadgetachew/cbe-super-app-shared/shared/utils"

	"go.mongodb.org/mongo-driver/v2/mongo"
)

func InitPersistanceLayer(client *mongo.Client, dbName string, coreConfig core.CBECoreCredential, merchantApi, merchantXAPIKey string, notificationApi string, notificationProducer kafka.NotificationProducer, clientOrchestrationProducer kafka.ClientOrchestrationProducer, cfg *config.VaultConfig, logger utils.Logger) persistance.Persistence {

	data := persistance.Persistence{
		JobRolePersistence:              job_repo.NewJobRoleRepository(client,cfg, dbName, JobRolesCollection, logger),
		DeviceVersionControlPersistence: deviceversioncontrol.NewDeviceVersionControlRepository(client, cfg,dbName, DeviceVersionControllCollection, clientOrchestrationProducer, logger),
		AccountLookup:                   core.NewCBECoreAPI(coreConfig),
		UserPersistence:                 users.NewUserRepository(client,cfg, dbName, MembersCollection, clientOrchestrationProducer, logger),
		HQPersistence:                   hq.NewHQRepository(client,cfg, dbName, HQCollection, logger),
		OTPPersistence:                  otp.NewOtpRepository(client,cfg, dbName, OTPsCollection, logger),
		DeviceLinkHistoryPersistence:    device_history.NewDeviceLinkHistoryRepository(client, cfg,dbName, MemberDevicesHistoryCollection, logger),
		ResetSessionPersistence:         reset_session.NewResetSessionRepository(client,cfg, dbName, PINResetsCollection, logger),
		CPSAction:                       cps_action.NewCPSActionRepository(client,cfg, dbName, CPSActionsCollection, logger),
		AmountBasedAuthPersistence:      amount_based_auth.NewAmountBasedAuthRepository(client,cfg, dbName, AuthTierCollection, logger),
		AccountBlockPersistence:         account_block.NewAccountBlockRepository(client,cfg, dbName, AccountBlockCollection, clientOrchestrationProducer, logger),
		PortalCardPersistence:           portal_card.NewPortalCardRepository(client, cfg,dbName, CardsCollection, logger),
		MiniAppPersistence:              mini_app.NewMiniAppRepository(client,cfg, dbName, MiniAppsCollection, logger),
		MerchantLookup:                  *merchant_lookup.NewMerchantLookupAdapter(merchantApi, *cfg, merchantXAPIKey, logger),
		SMSSenderApi:                    notificationProducer,
		// Additional repositories
		AccessListPersistence:             access_list.NewAccessListRepository(client,cfg, dbName, AccessListCollection, clientOrchestrationProducer, logger),
		AvatarPersistence:                 avatar.NewAvatarRepository(client,cfg, dbName, AvatarsCollection, logger),
		BPSUserPersistence:                bps_user.NewBPSUserRepository(client,cfg, dbName, BranchUserCollection, logger),
		AdvertRepositoryPersistence:       advert.NewAdvertRepository(client,cfg, dbName, AdvertsCollection, logger),
		ArchivedLinkedAccountPersistence:  archived_linked_account.NewArchivedLinkedAccountRepository(client,cfg, dbName, ArchievedLinkedAccountCollection, logger),
		AuthTierPersistence:               auth_tier.NewAuthTierRepository(client,cfg, dbName, AuthTierCollection, logger),
		BankPersistence:                   bank.NewBankRepository(client,cfg, dbName, BanksCollection, logger),
		BudgetCategoryPersistence:         budget_category.NewBudgetCategoryRepository(client,cfg, dbName, BudgetCategoryCollection, clientOrchestrationProducer, logger),
		BulkService:                       bulk_service.InitBulkServicePersistence(client,cfg, dbName, []string{CPSActionsCollection, AccessListCollection}, logger),
		CustomerService:                   customer.InitCustomerDetail(client,cfg, dbName, []string{MembersCollection, LinkedAccountsCollection}, clientOrchestrationProducer, logger),
		CpsUserPersistence:                cps_user.NewCPSUserRepository(client,cfg, dbName, CPSUsersCollection, []string{DepartmentsCollection, PermissionCollection, PermissionCategoryCollection, PermissionGroupsCollection}, logger),
		DonationPersistence:               donation.NewDonationRepository(client,cfg, dbName, DonationsCollection, clientOrchestrationProducer, logger),
		DonationCategoryPersistence:       donation_category.NewDonationCategoryRepository(client,cfg, dbName, DonationCategoriesCollection, clientOrchestrationProducer, logger),
		ArchivedUserPersistence:           archived_user.NewArchivedUserRepository(client,cfg, dbName, ArchievedUsersCollection, logger),
		DonationCompanyPersistence:        donation_company.NewDonationCompanyRepository(client,cfg, dbName, DonationCompaniesCollection, clientOrchestrationProducer, logger),
		EventPersistence:                  event.NewEventRepository(client,cfg, dbName, EventsCollection, logger),
		PasswordRulePersistent:            password.NewPasswordRuleRepository(client,cfg, dbName, PasswordRulesCollection, logger),
		FeedbackPersistence:               feedback.NewFeedbackRepository(client,cfg, dbName, FeedbackCollection, CustomerFeedbackCollection, logger),
		IconPersistence:                   icon.NewIconRepository(client,cfg, dbName, IconsCollection, logger),
		LinkedAccountPersistence:          linked_account.NewLinkedAccountRepository(client, cfg,dbName, LinkedAccountsCollection, clientOrchestrationProducer, logger),
		EcommerceMerchantPersistence:      ecommerce_merchant.NewEcommerceMerchantRepository(client,cfg, dbName, EcommerceMerchantCollection, logger),
		NotificationPersistence:           notification.NewNotificationRepository(client,cfg, dbName, NotificationsCollection, clientOrchestrationProducer, logger),
		PasswordRulePersistence:           password.NewPasswordRuleRepository(client,cfg, dbName, PasswordRulesCollection, logger),
		ServicesPersistence:               services_repo.NewServicesRepository(client,cfg, dbName, ServicesCollection, clientOrchestrationProducer, logger),
		ValidationRulePersistence:         accountvalidation.NewAccountValidationStore(client,cfg, dbName, ValidationRulesCollection, clientOrchestrationProducer, logger),
		WalletPersistence:                 wallet.NewWalletRepository(client,cfg, dbName, WalletsCollection, ServicesCollection, logger),
		TopupPersistence:                  Topup.NewTopupRepository(client, cfg,dbName, TopUpsCollection, logger),
		DepartmentPersistence:             department.NewDepartmentRepository(client,cfg, dbName, DepartmentsCollection, logger),
		FaydaPersistence:                  fayda.InitFaydaAccountPersistence(client,cfg, dbName, MembersCollection, logger),
		PermissionPersistence:             permission.InitPermission(client, cfg,dbName, []string{PermissionGroupsCollection, PermissionCategoryCollection, PermissionCollection, CPSActionsCollection}, 30*time.Second, logger),
		ArticlePersistence:                media.NewsArticleRepository(logger, client,cfg, dbName, NewsArticlesCollection, clientOrchestrationProducer),
		ArticleCategoryPersistence:        media.NewArticleCategoryRepository(logger, client,cfg, dbName, NewsCategoriesCollection),
		ShortVideoPersistence:             media.NewShortVideoRepository(logger, client, cfg,dbName, NewsShortVideosCollection, clientOrchestrationProducer),
		NewsTagPersistence:                newstag_repo.NewNewsTagRepository(client, cfg,dbName, NewsTagsCollection, clientOrchestrationProducer, logger),
		NewsCategoryPersistence:           newscategory_repo.NewNewsCategoryRepository(client,cfg, dbName, NewsCategoryCollection, clientOrchestrationProducer, logger),
		KYCVerifierPersistence:            kyc_repo.NewKYCVerifierRepository(client,cfg, dbName, CustomersKYCCollection, clientOrchestrationProducer, logger),
		NewsTagsServiceContainer:          media.NewNewsTagsRepository(logger, client,cfg, dbName, NewsTagsCollection),
		BPSActionRolePersistence:          actionrole_repo.NewBPSActionRoleRepository(client,cfg, dbName, []string{BPSActionRolesCollection, BPSActionListCollection}, logger),
		BPSActionApproveIndexPersistence:  actionrole_repo.NewBPSActionApproveIndexRepository(client, dbName, BPSActionApproveIndexCollection, logger),
		RolePersistence:                   role_repo.NewRoleRepository(client,cfg, dbName, RolesCollection, logger),
		MiniAppCategoryPersistence:        mini_app.NewMiniAppCategoryRepository(logger, client,cfg, dbName, MiniAppCategoryCollection),
		CPSActionRolePersistence:          cps_actionrole_repo.NewCPSActionRoleRepository(client, cfg,dbName, []string{CPSActionRolesCollection, CPSActionListCollection, CPSActionApproveIndexCollection}, logger),
		CPSActionApproveIndexPersistence:  cps_actionrole_repo.NewCPSActionApproveIndexRepository(client, dbName, CPSActionApproveIndexCollection, logger),
		EventMerchantPersistence:          event_merchant_repository.NewEventMerchantRepository(client,cfg, dbName, EventMerchantsCollection, logger),
		MiniAppProductCodePersistence:     mini_app.NewMiniAppProdutCodeRepository(logger, client, cfg,dbName, MiniAppProductCodes),
		AccessListSegmentationPersistence: access_list_segmentation_repository.NewAccessListSegmentationRepository(client, cfg, dbName, AccessListSegmentationCollection, logger),
		MiniAppMerchant:                   mini_app.NewMiniAppMerchantRepository(client, cfg,dbName, MiniAppMerchantCollection, logger),
		CustomerSegmentation:              customer_segmentation_repo.NewCustomerSegmentationRepository(client,cfg, dbName, CustomerSegmentationCollection, logger),
		CPSRoles:                          cps_roles.NewCPSRolesStorage(client, cfg ,dbName, CPSRolesCollection, logger),
		LogisticsMerchantPersistence:      logistics_merchant_repository.NewLogisticsMerchantRepository(client,cfg, dbName, LogisticsMerchantsCollection, logger),
	}

	return data
}
