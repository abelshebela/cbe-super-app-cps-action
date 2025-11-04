package initiator

import (
	"cbe-super-app-cps-action/internal/storage/external_call"
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
	"cbe-super-app-cps-action/internal/storage/persistance/bps_user"
	"cbe-super-app-cps-action/internal/storage/persistance/branch"
	"cbe-super-app-cps-action/internal/storage/persistance/city"
	"cbe-super-app-cps-action/internal/storage/persistance/cps_action"
	"cbe-super-app-cps-action/internal/storage/persistance/cps_user"
	"cbe-super-app-cps-action/internal/storage/persistance/department"
	"cbe-super-app-cps-action/internal/storage/persistance/donation_category"
	"cbe-super-app-cps-action/internal/storage/persistance/event"
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
	"cbe-super-app-cps-action/internal/storage/persistance/icon"
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

	"gitlab.com/bersufekadgetachew/cbe-super-app-shared/shared/utils"
	"go.mongodb.org/mongo-driver/v2/mongo"
)

func InitPersistanceLayer(client *mongo.Client, dbName string, logger utils.Logger) persistance.Persistence {
	data := persistance.Persistence{
		UserPersistence:              users.NewUserRepository(client, dbName, "members", logger),
		HQPersistence:                hq.NewHQRepository(client, dbName, "hq", logger),
		OTPPersistence:               otp.NewOtpRepository(client, dbName, "otps", logger),
		DeviceLinkHistoryPersistence: device_history.NewDeviceLinkHistoryRepository(client, dbName, "member_device_histories", logger),
		ResetSessionPersistence:      reset_session.NewResetSessionRepository(client, dbName, "pin_reset_sessions", logger),
		SMSSenderApi:                 *external_call.NewSMSPersistence("https://devcbe.eaglelionsystems.com/api/v1.0/chatbirrapi/ldapnotif/sms/send", logger),
		CPSAction:                    cps_action.NewCPSActionRepository(client, dbName, "cps_actions", logger),
		AmountBasedAuthPersistence:   amount_based_auth.NewAmountBasedAuthRepository(client, dbName, "auth_tier", logger),
		AccountBlockPersistence:      account_block.NewAccountBlockRepository(client, dbName, logger),
		PortalCardPersistence:        portal_card.NewPortalCardRepository(client, dbName, "cards", logger),
		MiniAppPersistence:           mini_app.NewMiniAppRepository(client, dbName, "mini_app", logger),
		CityPersistence:              city.NewCityRepository(client, dbName, "cities", logger),
		RegionPersistence:            region.NewRegionRepository(client, dbName, "regions", logger),
		DistrictPersistence:          district.NewDistrictRepository(client, dbName, "districts", logger),
		BranchPersistence:            branch.NewBranchRepository(client, dbName, "branches", logger),
		// Additional repositories
		AccessListPersistence:       access_list.NewAccessListRepository(client, dbName, "access_lists", logger),
		AvatarPersistence:           avatar.NewAvatarRepository(client, dbName, "avatars", logger),
		BPSUserPersistence:          bps_user.NewBPSUserRepository(client, dbName, "branch_user", logger),
		AdvertRepositoryPersistence: advert.NewAdvertRepository(client, dbName, "adverts", logger),

		ArchivedLinkedAccountPersistence: archived_linked_account.NewArchivedLinkedAccountRepository(client, dbName, "archived_linked_account", logger),
		AuthTierPersistence:              auth_tier.NewAuthTierRepository(client, dbName, "auth_tier", logger),
		BankPersistence:                  bank.NewBankRepository(client, dbName, "banks", logger),
		BudgetCategoryPersistence:        budget_category.NewBudgetCategoryRepository(client, dbName, "budget_category", logger),
		BulkService:                      bulk_service.InitBulkServicePersistence(client, dbName, []string{"cps_actions", "access_list"}, logger),
		CustomerService:                  customer.InitCustomerDetail(client, dbName, []string{"members", "linked_account"}, logger),
		CpsUserPersistence:               cps_user.NewCPSUserRepository(client, dbName, "cps_users", logger),
		DonationPersistence:              donation.NewDonationRepository(client, dbName, "donations", logger),
		DonationCategoryPersistence:      donation_category.NewDonationCategoryRepository(client, dbName, "donation_categories", logger),

		ArchivedUserPersistence: archived_user.NewArchivedUserRepository(client, dbName, "archived_users", logger),

		DonationCompanyPersistence: donation_company.NewDonationCompanyRepository(client, dbName, "donation_companies", logger),
		EventPersistence:           event.NewEventRepository(client, dbName, "events", logger),
		PasswordRulePersistent:     password.NewPasswordRuleRepository(client, dbName, "password_rules", logger),
		FeedbackPersistence:        feedback.NewFeedbackRepository(client, dbName, "feedback", logger),
		IconPersistence:            icon.NewIconRepository(client, dbName, "icons", logger),
		LinkedAccountPersistence:   linked_account.NewLinkedAccountRepository(client, dbName, "linked_account", logger),
		MiniAppMerchantPersistence: mini_app_merchant.NewMiniAppMerchantRepository(client, dbName, "mini_app_merchant", logger),
		NotificationPersistence:    notification.NewNotificationRepository(client, dbName, "notifications", logger),
		PasswordRulePersistence:    password.NewPasswordRuleRepository(client, dbName, "password_rules", logger),
		ServiceDetailsPersistence:  service_details.NewServiceDetailsRepository(client, dbName, "services", logger),
		ValidationRulePersistence:  accountvalidation.NewAccountValidationStore(client, dbName, "validation_rule", logger),
		WalletPersistence:          wallet.NewWalletRepository(client, dbName, "wallets", logger),
		TopupPersistence:           Topup.NewTopupRepository(client, dbName, "topups", logger),

		ProductCodePersistence:     productcode.NewProductCodeRepository(client, dbName, "services", logger),
		DepartmentPersistence:      department.NewDepartmentRepository(client, dbName, "department", logger),
		FaydaPersistence:           fayda.InitFaydaAccountPersistence(client, dbName, "members", logger),
		PermissionPersistence:      permission.InitPermission(client, dbName, 30*time.Second, logger),
		ArticlePersistence:         media.NewsArticleRepository(logger, client, dbName, "news_articles"),
		ArticleCategoryPersistence: media.NewArticleCategoryRepository(logger, client, dbName, "news_categories"),
		ShortVideoPersistence:      media.NewShortVideoRepository(logger, client, dbName, "news_short_videos"),
		NewsTagPersistence:         newstag_repo.NewNewsTagRepository(client, dbName, "news_tags", logger),
		NewsCategoryPersistence:    newscategory_repo.NewNewsCategoryRepository(client, dbName, "news_category", logger),
	}

	return data
}
