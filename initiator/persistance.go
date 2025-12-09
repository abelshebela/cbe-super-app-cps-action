package initiator

import (
	"cbe-super-app-cps-action/internal/storage/external_call/merchant_lookup"
	"cbe-super-app-cps-action/internal/storage/kafka"

	// "cbe-super-app-cps-action/internal/storage/external_call/account_lookup"
	"cbe-super-app-cps-action/internal/storage/persistance"
	"cbe-super-app-cps-action/internal/storage/persistance/access_list"
	"cbe-super-app-cps-action/internal/storage/persistance/account_block"
	accountvalidation "cbe-super-app-cps-action/internal/storage/persistance/account_validation"
	"cbe-super-app-cps-action/internal/storage/persistance/advert"
	"cbe-super-app-cps-action/internal/storage/persistance/amount_based_auth"
	"cbe-super-app-cps-action/internal/storage/persistance/auth_tier"
	"cbe-super-app-cps-action/internal/storage/persistance/avatar"
	"cbe-super-app-cps-action/internal/storage/persistance/bank"
	actionrole_repo "cbe-super-app-cps-action/internal/storage/persistance/bps_action_role"
	"cbe-super-app-cps-action/internal/storage/persistance/bps_user"
	"cbe-super-app-cps-action/internal/storage/persistance/branch"
	"cbe-super-app-cps-action/internal/storage/persistance/city"
	"cbe-super-app-cps-action/internal/storage/persistance/cps_action"
	"cbe-super-app-cps-action/internal/storage/persistance/cps_user"
	"cbe-super-app-cps-action/internal/storage/persistance/department"
	deviceversioncontrol "cbe-super-app-cps-action/internal/storage/persistance/device_version_control"
	"cbe-super-app-cps-action/internal/storage/persistance/donation_category"
	"cbe-super-app-cps-action/internal/storage/persistance/event"
	"cbe-super-app-cps-action/internal/storage/persistance/icon"
	"cbe-super-app-cps-action/internal/storage/persistance/linked_account"
	"cbe-super-app-cps-action/internal/storage/persistance/media"
	newscategory_repo "cbe-super-app-cps-action/internal/storage/persistance/news_category"
	newstag_repo "cbe-super-app-cps-action/internal/storage/persistance/news_tag"
	"time"

	// "cbe-super-app-cps-action/internal/storage/persistance/cps_user"
	"cbe-super-app-cps-action/internal/storage/persistance/device_history"
	"cbe-super-app-cps-action/internal/storage/persistance/district"
	"cbe-super-app-cps-action/internal/storage/persistance/donation"

	// "cbe-super-app-cps-action/internal/storage/persistance/donation_category"
	"cbe-super-app-cps-action/internal/storage/persistance/donation_company"
	// "cbe-super-app-cps-action/internal/storage/persistance/event"
	"cbe-super-app-cps-action/internal/storage/persistance/archived_linked_account"
	"cbe-super-app-cps-action/internal/storage/persistance/archived_user"
	"cbe-super-app-cps-action/internal/storage/persistance/budget_category"
	"cbe-super-app-cps-action/internal/storage/persistance/fayda"
	"cbe-super-app-cps-action/internal/storage/persistance/feedback"
	"cbe-super-app-cps-action/internal/storage/persistance/hq"
	kyc_repo "cbe-super-app-cps-action/internal/storage/persistance/kyc_verifier"
	"cbe-super-app-cps-action/internal/storage/persistance/mini_app"
	"cbe-super-app-cps-action/internal/storage/persistance/mini_app_merchant"
	"cbe-super-app-cps-action/internal/storage/persistance/notification"
	"cbe-super-app-cps-action/internal/storage/persistance/otp"

	password "cbe-super-app-cps-action/internal/storage/persistance/password_rule"
	permission "cbe-super-app-cps-action/internal/storage/persistance/permission"
	"cbe-super-app-cps-action/internal/storage/persistance/portal_card"
	productcode "cbe-super-app-cps-action/internal/storage/persistance/product_code"
	"cbe-super-app-cps-action/internal/storage/persistance/region"
	"cbe-super-app-cps-action/internal/storage/persistance/reset_session"
	"cbe-super-app-cps-action/internal/storage/persistance/service_details"
	Topup "cbe-super-app-cps-action/internal/storage/persistance/topup"
	"cbe-super-app-cps-action/internal/storage/persistance/users"
	"cbe-super-app-cps-action/internal/storage/persistance/wallet"

	"cbe-super-app-cps-action/internal/storage/persistance/bulk_service"
	"cbe-super-app-cps-action/internal/storage/persistance/customer"

	"github.com/hugokessem/coreio/core"
	"gitlab.com/bersufekadgetachew/cbe-super-app-shared/shared/config"
	"gitlab.com/bersufekadgetachew/cbe-super-app-shared/shared/utils"

	"go.mongodb.org/mongo-driver/v2/mongo"
)

func InitPersistanceLayer(client *mongo.Client, dbName string, coreConfig core.CBECoreCredential, merchantApi, merchantXAPIKey string, notificationApi string, kafkaService kafka.NotificationProducer, cfg *config.VaultConfig, logger utils.Logger) persistance.Persistence {

	data := persistance.Persistence{
		DeviceVersionControlPersistence: deviceversioncontrol.NewDeviceVersionControlRepository(client, dbName, DeviceVersionControllCollection, logger),
		AccountLookup:                   core.NewCBECoreAPI(coreConfig),
		UserPersistence:                 users.NewUserRepository(client, dbName, MembersCollection, logger),
		HQPersistence:                   hq.NewHQRepository(client, dbName, HQCollection, logger),
		OTPPersistence:                  otp.NewOtpRepository(client, dbName, OTPsCollection, logger),
		DeviceLinkHistoryPersistence:    device_history.NewDeviceLinkHistoryRepository(client, dbName, MemberDevicesHistoryCollection, logger),
		ResetSessionPersistence:         reset_session.NewResetSessionRepository(client, dbName, PINResetsCollection, logger),
		CPSAction:                       cps_action.NewCPSActionRepository(client, dbName, CPSActionsCollection, logger),
		AmountBasedAuthPersistence:      amount_based_auth.NewAmountBasedAuthRepository(client, dbName, AuthTierCollection, logger),
		AccountBlockPersistence:         account_block.NewAccountBlockRepository(client, dbName, AccountBlockCollection, logger),
		PortalCardPersistence:           portal_card.NewPortalCardRepository(client, dbName, CardsCollection, logger),
		MiniAppPersistence:              mini_app.NewMiniAppRepository(client, dbName, MiniAppsCollection, logger),
		MerchantLookup:                  *merchant_lookup.NewMerchantLookupAdapter(merchantApi, *cfg, merchantXAPIKey, logger),
		SMSSenderApi:                    kafkaService,
		CityPersistence:                 city.NewCityRepository(client, dbName, "cities", logger),
		RegionPersistence:               region.NewRegionRepository(client, dbName, "regions", logger),
		DistrictPersistence:             district.NewDistrictRepository(client, dbName, "districts", logger),
		BranchPersistence:               branch.NewBranchRepository(client, dbName, "branches", logger),
		// Additional repositories
		AccessListPersistence:            access_list.NewAccessListRepository(client, dbName, AccessListCollection, logger),
		AvatarPersistence:                avatar.NewAvatarRepository(client, dbName, AvatarsCollection, logger),
		BPSUserPersistence:               bps_user.NewBPSUserRepository(client, dbName, BranchUserCollection, logger),
		AdvertRepositoryPersistence:      advert.NewAdvertRepository(client, dbName, AdvertsCollection, logger),
		ArchivedLinkedAccountPersistence: archived_linked_account.NewArchivedLinkedAccountRepository(client, dbName, ArchievedLinkedAccountCollection, logger),
		AuthTierPersistence:              auth_tier.NewAuthTierRepository(client, dbName, AuthTierCollection, logger),
		BankPersistence:                  bank.NewBankRepository(client, dbName, BanksCollection, logger),
		BudgetCategoryPersistence:        budget_category.NewBudgetCategoryRepository(client, dbName, BudgetCategoryCollection, logger),
		BulkService:                      bulk_service.InitBulkServicePersistence(client, dbName, []string{CPSActionsCollection, AccessListCollection}, logger),
		CustomerService:                  customer.InitCustomerDetail(client, dbName, []string{MembersCollection, LinkedAccountsCollection}, logger),
		CpsUserPersistence:               cps_user.NewCPSUserRepository(client, dbName, CPSUsersCollection, []string{DepartmentsCollection, PermissionCollection, PermissionCategoryCollection, PermissionGroupsCollection}, logger),
		DonationPersistence:              donation.NewDonationRepository(client, dbName, DonationsCollection, logger),
		DonationCategoryPersistence:      donation_category.NewDonationCategoryRepository(client, dbName, DonationCategoriesCollection, logger),
		ArchivedUserPersistence:          archived_user.NewArchivedUserRepository(client, dbName, ArchievedUsersCollection, logger),
		DonationCompanyPersistence:       donation_company.NewDonationCompanyRepository(client, dbName, DonationCompaniesCollection, logger),
		EventPersistence:                 event.NewEventRepository(client, dbName, EventsCollection, logger),
		PasswordRulePersistent:           password.NewPasswordRuleRepository(client, dbName, PasswordRulesCollection, logger),
		FeedbackPersistence:              feedback.NewFeedbackRepository(client, dbName, FeedbackCollection, logger),
		IconPersistence:                  icon.NewIconRepository(client, dbName, IconsCollection, logger),
		LinkedAccountPersistence:         linked_account.NewLinkedAccountRepository(client, dbName, LinkedAccountsCollection, logger),
		MiniAppMerchantPersistence:       mini_app_merchant.NewMiniAppMerchantRepository(client, dbName, MiniAppMerchantCollection, logger),
		NotificationPersistence:          notification.NewNotificationRepository(client, dbName, NotificationsCollection, logger),
		PasswordRulePersistence:          password.NewPasswordRuleRepository(client, dbName, PasswordRulesCollection, logger),
		ServiceDetailsPersistence:        service_details.NewServiceDetailsRepository(client, dbName, ServicesCollection, logger),
		ValidationRulePersistence:        accountvalidation.NewAccountValidationStore(client, dbName, ValidationRulesCollection, logger),
		WalletPersistence:                wallet.NewWalletRepository(client, dbName, WalletsCollection, logger),
		TopupPersistence:                 Topup.NewTopupRepository(client, dbName, TopUpsCollection, logger),
		ProductCodePersistence:           productcode.NewProductCodeRepository(client, dbName, ServicesCollection, logger),
		DepartmentPersistence:            department.NewDepartmentRepository(client, dbName, DepartmentsCollection, logger),
		FaydaPersistence:                 fayda.InitFaydaAccountPersistence(client, dbName, MembersCollection, logger),
		PermissionPersistence:            permission.InitPermission(client, dbName, []string{PermissionGroupsCollection, PermissionCategoryCollection, PermissionCollection, CPSActionsCollection}, 30*time.Second, logger),
		ArticlePersistence:               media.NewsArticleRepository(logger, client, dbName, NewsArticlesCollection),
		ArticleCategoryPersistence:       media.NewArticleCategoryRepository(logger, client, dbName, NewsCategoriesCollection),
		ShortVideoPersistence:            media.NewShortVideoRepository(logger, client, dbName, NewsShortVideosCollection),
		NewsTagPersistence:               newstag_repo.NewNewsTagRepository(client, dbName, NewsTagsCollection, logger),
		NewsCategoryPersistence:          newscategory_repo.NewNewsCategoryRepository(client, dbName, NewsCategoryCollection, logger),
		KYCVerifierPersistence:           kyc_repo.NewKYCVerifierRepository(client, dbName, CustomersKYCCollection, logger),
		NewsTagsServiceContainer:         media.NewNewsTagsRepository(logger, client, dbName, NewsTagsCollection),
		BPSActionRolePersistence:         actionrole_repo.NewBPSActionRoleRepository(client, dbName, BPSActionRolesCollection, logger),
	}

	return data
}
