package initiator

import (
	"github.com/abelshebela/cbe-super-app-cps-action/internal/storage"
	"github.com/abelshebela/cbe-super-app-cps-action/internal/storage/external_call/merchant_lookup"
	"github.com/abelshebela/cbe-super-app-cps-action/internal/storage/kafka"

	shared_producer "gitlab.com/bersufekadgetachew/cbe-super-app-shared/shared/notification/producer"

	"github.com/abelshebela/cbe-super-app-cps-action/internal/storage/persistance"
	"github.com/abelshebela/cbe-super-app-cps-action/internal/storage/persistance/access_list"
	access_list_segmentation_repository "github.com/abelshebela/cbe-super-app-cps-action/internal/storage/persistance/access_list_segmentation"
	accountvalidation "github.com/abelshebela/cbe-super-app-cps-action/internal/storage/persistance/account_validation"
	"github.com/abelshebela/cbe-super-app-cps-action/internal/storage/persistance/advert"
	"github.com/abelshebela/cbe-super-app-cps-action/internal/storage/persistance/amount_based_auth"
	"github.com/abelshebela/cbe-super-app-cps-action/internal/storage/persistance/auth_tier"
	"github.com/abelshebela/cbe-super-app-cps-action/internal/storage/persistance/avatar"
	"github.com/abelshebela/cbe-super-app-cps-action/internal/storage/persistance/bank"
	"github.com/abelshebela/cbe-super-app-cps-action/internal/storage/persistance/bps_action"
	actionrole_repo "github.com/abelshebela/cbe-super-app-cps-action/internal/storage/persistance/bps_action_role"
	"github.com/abelshebela/cbe-super-app-cps-action/internal/storage/persistance/bps_user"
	"github.com/abelshebela/cbe-super-app-cps-action/internal/storage/persistance/cps_action"
	cps_actionrole_repo "github.com/abelshebela/cbe-super-app-cps-action/internal/storage/persistance/cps_action_role"
	"github.com/abelshebela/cbe-super-app-cps-action/internal/storage/persistance/cps_user"
	"github.com/abelshebela/cbe-super-app-cps-action/internal/storage/persistance/department"
	"github.com/abelshebela/cbe-super-app-cps-action/internal/storage/persistance/device_history"
	deviceversioncontrol "github.com/abelshebela/cbe-super-app-cps-action/internal/storage/persistance/device_version_control"
	"github.com/abelshebela/cbe-super-app-cps-action/internal/storage/persistance/event"
	role_delegation_repo "github.com/abelshebela/cbe-super-app-cps-action/internal/storage/persistance/role_delegation"

	// event_merchant_repository "github.com/abelshebela/cbe-super-app-cps-action/internal/storage/persistance/event_merchant"
	// event_merchant_repository_oracle "github.com/abelshebela/cbe-super-app-cps-action/internal/storage/persistance/event_merchant_oracle"

	"github.com/abelshebela/cbe-super-app-cps-action/internal/storage/persistance/icon"
	"github.com/abelshebela/cbe-super-app-cps-action/internal/storage/persistance/linked_account"
	logistics_merchant_repository "github.com/abelshebela/cbe-super-app-cps-action/internal/storage/persistance/logistics_merchant"
	"github.com/abelshebela/cbe-super-app-cps-action/internal/storage/persistance/media"
	newscategory_repo "github.com/abelshebela/cbe-super-app-cps-action/internal/storage/persistance/news_category"
	newstag_repo "github.com/abelshebela/cbe-super-app-cps-action/internal/storage/persistance/news_tag"
	"time"

	"github.com/abelshebela/cbe-super-app-cps-action/internal/storage/persistance/archived_linked_account"
	"github.com/abelshebela/cbe-super-app-cps-action/internal/storage/persistance/archived_user"

	// ecommerce_merchant "github.com/abelshebela/cbe-super-app-cps-action/internal/storage/persistance/ecommerce_merchant"

	"github.com/abelshebela/cbe-super-app-cps-action/internal/storage/persistance/fayda"
	"github.com/abelshebela/cbe-super-app-cps-action/internal/storage/persistance/feedback"
	"github.com/abelshebela/cbe-super-app-cps-action/internal/storage/persistance/hq"
	kyc_repo "github.com/abelshebela/cbe-super-app-cps-action/internal/storage/persistance/kyc_verifier"

	// "github.com/abelshebela/cbe-super-app-cps-action/internal/storage/persistance/mini_app"
	"github.com/abelshebela/cbe-super-app-cps-action/internal/storage/persistance/notification"
	"github.com/abelshebela/cbe-super-app-cps-action/internal/storage/persistance/otp"

	job_roles "github.com/abelshebela/cbe-super-app-cps-action/internal/storage/persistance/job_roles"
	password "github.com/abelshebela/cbe-super-app-cps-action/internal/storage/persistance/password_rule"
	permission "github.com/abelshebela/cbe-super-app-cps-action/internal/storage/persistance/permission"
	"github.com/abelshebela/cbe-super-app-cps-action/internal/storage/persistance/portal_card"
	"github.com/abelshebela/cbe-super-app-cps-action/internal/storage/persistance/reset_session"
	roles "github.com/abelshebela/cbe-super-app-cps-action/internal/storage/persistance/roles"
	Topup "github.com/abelshebela/cbe-super-app-cps-action/internal/storage/persistance/topup"
	"github.com/abelshebela/cbe-super-app-cps-action/internal/storage/persistance/users"
	"github.com/abelshebela/cbe-super-app-cps-action/internal/storage/persistance/wallet"

	"github.com/abelshebela/cbe-super-app-cps-action/internal/storage/persistance/bulk_service"
	survey_sampling_persistence "github.com/abelshebela/cbe-super-app-cps-action/internal/storage/persistance/survey_sampling"
	transaction_limit_repo "github.com/abelshebela/cbe-super-app-cps-action/internal/storage/persistance/transaction_limit"
	user_action_log_repo "github.com/abelshebela/cbe-super-app-cps-action/internal/storage/persistance/user_action_log"

	"github.com/hugokessem/coreio/core"
	"gitlab.com/bersufekadgetachew/cbe-super-app-shared/shared/config"
	"gitlab.com/bersufekadgetachew/cbe-super-app-shared/shared/utils"

	"go.mongodb.org/mongo-driver/v2/mongo"
)

func InitPersistanceLayer(client *mongo.Client, dbName string, coreInterface core.CBECoreAPIInterface, notificationApi string, notificationProducer kafka.NotificationProducer, sharedKafkaProducer *shared_producer.NotificationProducer, clientOrchestrationProducer kafka.ClientOrchestrationProducer, accessListSegmentationProducer kafka.AccessListSegmentationProducer, redisRepository storage.RedisRepository, cfg *config.VaultConfig, logger utils.Logger) persistance.Persistence {
	// AccountBlock now uses Oracle — initialized in OraclePersistence, wired in service.go
	data := persistance.Persistence{
		RolePersistence:                 roles.NewRoleRepository(client, cfg, dbName, JobRolesCollection, logger),
		DeviceVersionControlPersistence: deviceversioncontrol.NewDeviceVersionControlRepository(client, cfg, dbName, DeviceVersionControllCollection, clientOrchestrationProducer, logger),
		AccountLookup:                   coreInterface,
		UserPersistence:                 users.NewUserRepository(client, cfg, dbName, MembersCollection, clientOrchestrationProducer, logger),
		HQPersistence:                   hq.NewHQRepository(client, cfg, dbName, HQCollection, logger),
		OTPPersistence:                  otp.NewOtpRepository(client, cfg, dbName, OTPsCollection, logger),
		DeviceLinkHistoryPersistence:    device_history.NewDeviceLinkHistoryRepository(client, cfg, dbName, MemberDevicesHistoryCollection, logger),
		ResetSessionPersistence:         reset_session.NewResetSessionRepository(client, cfg, dbName, PINResetsCollection, logger),
		CPSAction:                       cps_action.NewCPSActionRepository(client, cfg, dbName, CPSActionsCollection, logger),
		AmountBasedAuthPersistence:      amount_based_auth.NewAmountBasedAuthRepository(client, cfg, dbName, AuthTierCollection, logger),
		AccountBlockPersistence:         nil, // Set from OraclePersistence in service.go
		PortalCardPersistence:           portal_card.NewPortalCardRepository(client, cfg, dbName, CardsCollection, logger),
		// MiniAppPersistence:               mini_app.NewMiniAppRepository(client, cfg, dbName, MiniAppsCollection, logger),
		MerchantLookup:                   *merchant_lookup.NewMerchantLookupAdapter(*cfg, logger),
		SMSSenderApi:                     notificationProducer,
		AccessListPersistence:            access_list.NewAccessListRepository(client, cfg, dbName, AccessListCollection, clientOrchestrationProducer, logger),
		AvatarPersistence:                avatar.NewAvatarRepository(client, cfg, dbName, AvatarsCollection, logger),
		BPSUserPersistence:               bps_user.NewBPSUserRepository(client, cfg, dbName, []string{BranchUserCollection, RoleDelegationCollection}, logger),
		AdvertRepositoryPersistence:      advert.NewAdvertRepository(client, cfg, dbName, AdvertsCollection, logger),
		ArchivedLinkedAccountPersistence: archived_linked_account.NewArchivedLinkedAccountRepository(client, cfg, dbName, ArchievedLinkedAccountCollection, logger),
		AuthTierPersistence:              auth_tier.NewAuthTierRepository(client, cfg, dbName, AuthTierCollection, logger),
		BankPersistence:                  bank.NewBankRepository(client, cfg, dbName, BanksCollection, logger),
		// BudgetCategoryPersistence:         budget_category.NewBudgetCategoryRepository(client, cfg, dbName, BudgetCategoryCollection, clientOrchestrationProducer, logger),

		BulkService: bulk_service.InitBulkServicePersistence(client, cfg, dbName, []string{CPSActionsCollection, AccessListCollection}, clientOrchestrationProducer, notificationProducer, logger),
		// CustomerService:          customer.InitCustomerDetail(client, cfg, dbName, []string{MembersCollection, LinkedAccountsCollection}, clientOrchestrationProducer, logger),
		CpsUserPersistence: cps_user.NewCPSUserRepository(client, redisRepository, cfg, dbName, CPSUsersCollection, []string{DepartmentsCollection, PermissionCollection, PermissionCategoryCollection, PermissionGroupsCollection, RolesCollection, JobRolesCollection, RoleDelegationCollection}, logger),

		ArchivedUserPersistence:  archived_user.NewArchivedUserRepository(client, cfg, dbName, ArchievedUsersCollection, logger),
		EventPersistence:         event.NewEventRepository(client, cfg, dbName, EventsCollection, logger),
		PasswordRulePersistent:   password.NewPasswordRuleRepository(client, cfg, dbName, PasswordRulesCollection, logger),
		FeedbackPersistence:      feedback.NewFeedbackRepository(client, cfg, dbName, FeedbackCollection, CustomerFeedbackCollection, SurveyFeedbackCollection, logger),
		IconPersistence:          icon.NewIconRepository(client, cfg, dbName, IconsCollection, logger),
		LinkedAccountPersistence: linked_account.NewLinkedAccountRepository(client, cfg, dbName, LinkedAccountsCollection, clientOrchestrationProducer, logger),
		// EcommerceMerchantPersistence:      ecommerce_merchant.NewEcommerceMerchantRepository(client, cfg, dbName, EcommerceMerchantCollection, logger),
		NotificationPersistence:          notification.NewNotificationRepository(client, cfg, dbName, NotificationsCollection, notificationProducer, logger),
		PasswordRulePersistence:          password.NewPasswordRuleRepository(client, cfg, dbName, PasswordRulesCollection, logger),
		ValidationRulePersistence:        accountvalidation.NewAccountValidationStore(client, cfg, dbName, ValidationRulesCollection, clientOrchestrationProducer, logger),
		WalletPersistence:                wallet.NewWalletRepository(client, cfg, dbName, WalletsCollection, ServicesCollection, logger),
		BpsActionPersistence:             bps_action.NewBPSActionRepository(client, dbName, BPSActionsCollection, logger, cfg),
		TopupPersistence:                 Topup.NewTopupRepository(client, cfg, dbName, TopUpsCollection, logger),
		DepartmentPersistence:            department.NewDepartmentRepository(client, cfg, dbName, DepartmentsCollection, logger),
		FaydaPersistence:                 fayda.InitFaydaAccountPersistence(client, cfg, dbName, MembersCollection, logger),
		PermissionPersistence:            permission.InitPermission(client, cfg, dbName, []string{PermissionGroupsCollection, PermissionCategoryCollection, PermissionCollection, CPSActionsCollection}, 30*time.Second, logger),
		ArticlePersistence:               media.NewsArticleRepository(logger, client, cfg, dbName, NewsArticlesCollection, clientOrchestrationProducer),
		ArticleCategoryPersistence:       media.NewArticleCategoryRepository(logger, client, cfg, dbName, NewsCategoriesCollection),
		ShortVideoPersistence:            media.NewShortVideoRepository(logger, client, cfg, dbName, NewsShortVideosCollection, clientOrchestrationProducer),
		NewsTagPersistence:               newstag_repo.NewNewsTagRepository(client, cfg, dbName, NewsTagsCollection, clientOrchestrationProducer, logger),
		NewsCategoryPersistence:          newscategory_repo.NewNewsCategoryRepository(client, cfg, dbName, NewsCategoryCollection, clientOrchestrationProducer, logger),
		KYCVerifierPersistence:           kyc_repo.NewKYCVerifierRepository(client, cfg, dbName, CustomersKYCCollection, clientOrchestrationProducer, logger),
		NewsTagsServiceContainer:         media.NewNewsTagsRepository(logger, client, cfg, dbName, NewsTagsCollection),
		BPSActionRolePersistence:         actionrole_repo.NewBPSActionRoleRepository(client, cfg, dbName, []string{BPSActionRolesCollection, BPSActionListCollection, BPSActionApproveIndexCollection}, logger),
		BPSActionApproveIndexPersistence: actionrole_repo.NewBPSActionApproveIndexRepository(client, dbName, BPSActionApproveIndexCollection, logger),
		JobRolePersistence:               job_roles.NewJobRoleRepository(client, cfg, dbName, []string{RolesCollection, JobRolesCollection, CPSUsersCollection}, logger),
		// MiniAppCategoryPersistence:        mini_app.NewMiniAppCategoryRepository(logger, client, cfg, dbName, MiniAppCategoryCollection),
		CPSActionRolePersistence:         cps_actionrole_repo.NewCPSActionRoleRepository(client, cfg, dbName, []string{CPSActionRolesCollection, CPSActionListCollection, CPSActionApproveIndexCollection}, logger),
		CPSActionApproveIndexPersistence: cps_actionrole_repo.NewCPSActionApproveIndexRepository(client, dbName, CPSActionApproveIndexCollection, logger),
		// EventMerchantPersistence:          event_merchant_repository.NewEventMerchantRepository(client, cfg, dbName, EventMerchantsCollection, logger),

		// MiniAppProductCodePersistence:     mini_app.NewMiniAppProdutCodeRepository(logger, client, cfg, dbName, MiniAppProductCodes),
		AccessListSegmentationPersistence: access_list_segmentation_repository.NewAccessListSegmentationRepository(client, cfg, dbName, AccessListSegmentationCollection, accessListSegmentationProducer, nil, logger), // AccountBlock injected in service.go
		// MiniAppMerchant:                   mini_app.NewMiniAppMerchantRepository(client, cfg, dbName, MiniAppMerchantCollection, logger),
		// CustomerSegmentation:              customer_segmentation_repo.NewCustomerSegmentationRepository(client, cfg, dbName, CustomerSegmentationCollection, clientOrchestrationProducer, logger),
		// CPSRoles:                     cps_roles.NewCPSRolesStorage(client, cfg, dbName, []string{CPSRolesCollection, AccessListCollection, AccessListSegmentationCollection}, clientOrchestrationProducer, logger),
		LogisticsMerchantPersistence: logistics_merchant_repository.NewLogisticsMerchantRepository(client, cfg, dbName, LogisticsMerchantsCollection, logger),
		// CustomerKYCPersistence:       persistence_kyc.NewCustomerKYCRepository(client, cfg, dbName, CustomersKYCCollection, logger),
		// UssdMerchantPersistence:   ussd_merchant_repo.NewUssdMerchant(client, dbName, UssdMerchantCollection, cfg, logger),
		UserActionLogPersistence:    user_action_log_repo.NewUserActionLogRepository(client, cfg, dbName, UserActionLogsCollection, logger),
		ServicesPersistence:         nil,                                                                                                                                                         // Services repository is commented out - using Oracle instead
		RoleDelegationPersistence:   role_delegation_repo.NewRoleDelegationRepository(client, cfg, dbName, []string{RoleDelegationCollection, CPSUsersCollection, BranchUserCollection}, logger), // Initialized in service.go to avoid circular dependency
		TransactionLimitPersistence:  transaction_limit_repo.NewTransactionLimitRepository(client, cfg, dbName, TransactionLimitsCollection, logger),
		SurveySamplingPersistence:    survey_sampling_persistence.NewSurveySamplingRepository(client, cfg, dbName, SurveySamplingConfigsCollection, logger),
	}

	return data
}
