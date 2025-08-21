package initiator

import (
	"time"

	outboundStore "github.com/CBE-Super-App/cbe-super-app-cps-action/internal/adapter/outbound"
	account_lookup_impl "github.com/CBE-Super-App/cbe-super-app-cps-action/internal/adapter/outbound/bps_calls"
	dal "github.com/CBE-Super-App/cbe-super-app-cps-action/internal/adapter/outbound/infra"
	account_block_repo "github.com/CBE-Super-App/cbe-super-app-cps-action/internal/adapter/outbound/persistence/account_block"
	account_validation "github.com/CBE-Super-App/cbe-super-app-cps-action/internal/adapter/outbound/persistence/account_validation"
	advert "github.com/CBE-Super-App/cbe-super-app-cps-action/internal/adapter/outbound/persistence/ad"
	amount_based_persistence "github.com/CBE-Super-App/cbe-super-app-cps-action/internal/adapter/outbound/persistence/amount_based_auth"
	avatarPersitence "github.com/CBE-Super-App/cbe-super-app-cps-action/internal/adapter/outbound/persistence/avatar"
	bank_repo "github.com/CBE-Super-App/cbe-super-app-cps-action/internal/adapter/outbound/persistence/bank"
	bps_user_persistence "github.com/CBE-Super-App/cbe-super-app-cps-action/internal/adapter/outbound/persistence/bps_user"
	budget_repo "github.com/CBE-Super-App/cbe-super-app-cps-action/internal/adapter/outbound/persistence/budget"
	budget_category_repo "github.com/CBE-Super-App/cbe-super-app-cps-action/internal/adapter/outbound/persistence/budget_category"
	bulk_service_repo "github.com/CBE-Super-App/cbe-super-app-cps-action/internal/adapter/outbound/persistence/bulk_service"
	cps_Actions_repo "github.com/CBE-Super-App/cbe-super-app-cps-action/internal/adapter/outbound/persistence/cps_action"
	customer_repo "github.com/CBE-Super-App/cbe-super-app-cps-action/internal/adapter/outbound/persistence/customer"
	dept_repo "github.com/CBE-Super-App/cbe-super-app-cps-action/internal/adapter/outbound/persistence/department"
	donation_persistence "github.com/CBE-Super-App/cbe-super-app-cps-action/internal/adapter/outbound/persistence/donation"
	event_persistence "github.com/CBE-Super-App/cbe-super-app-cps-action/internal/adapter/outbound/persistence/event"
	faydaaccount "github.com/CBE-Super-App/cbe-super-app-cps-action/internal/adapter/outbound/persistence/fayda_account"
	feedback_repo "github.com/CBE-Super-App/cbe-super-app-cps-action/internal/adapter/outbound/persistence/feedback"
	hq_persistence "github.com/CBE-Super-App/cbe-super-app-cps-action/internal/adapter/outbound/persistence/hq"
	inappnotification_persistence "github.com/CBE-Super-App/cbe-super-app-cps-action/internal/adapter/outbound/persistence/in-app-notification"
	miniApp_persistance "github.com/CBE-Super-App/cbe-super-app-cps-action/internal/adapter/outbound/persistence/mini_app"
	miniApp_merchant_persistance "github.com/CBE-Super-App/cbe-super-app-cps-action/internal/adapter/outbound/persistence/mini_app_merchant"
	notification_persistence "github.com/CBE-Super-App/cbe-super-app-cps-action/internal/adapter/outbound/persistence/notification"
	perm_repo "github.com/CBE-Super-App/cbe-super-app-cps-action/internal/adapter/outbound/persistence/permission"
	portal_card_persistence "github.com/CBE-Super-App/cbe-super-app-cps-action/internal/adapter/outbound/persistence/portal_card"
	productcode_persistence "github.com/CBE-Super-App/cbe-super-app-cps-action/internal/adapter/outbound/persistence/product_code"
	servicePersistence "github.com/CBE-Super-App/cbe-super-app-cps-action/internal/adapter/outbound/persistence/service"
	unlink_repo "github.com/CBE-Super-App/cbe-super-app-cps-action/internal/adapter/outbound/persistence/unlink"
	wallet_repo "github.com/CBE-Super-App/cbe-super-app-cps-action/internal/adapter/outbound/persistence/wallet"
	account_lookup_domain "github.com/CBE-Super-App/cbe-super-app-cps-action/internal/domain/account_lookup"
	"github.com/CBE-Super-App/cbe-super-app-cps-action/internal/domain/avatar"
	cps_actions "github.com/CBE-Super-App/cbe-super-app-cps-action/internal/domain/cps_actions/repository"
	dept_entities "github.com/CBE-Super-App/cbe-super-app-cps-action/internal/domain/department/entities"
	donation "github.com/CBE-Super-App/cbe-super-app-cps-action/internal/domain/donation"
	event_domain "github.com/CBE-Super-App/cbe-super-app-cps-action/internal/domain/event"
	miniApp_domain "github.com/CBE-Super-App/cbe-super-app-cps-action/internal/domain/miniapp"
	miniApp_merchant_domain "github.com/CBE-Super-App/cbe-super-app-cps-action/internal/domain/miniapp_merchant"
	notification_domain "github.com/CBE-Super-App/cbe-super-app-cps-action/internal/domain/notification"
	portalCardRepo "github.com/CBE-Super-App/cbe-super-app-cps-action/internal/domain/portal_card"
	productcode "github.com/CBE-Super-App/cbe-super-app-cps-action/internal/domain/product_code"
	serviceDomain "github.com/CBE-Super-App/cbe-super-app-cps-action/internal/domain/service"
	wallet "github.com/CBE-Super-App/cbe-super-app-cps-action/internal/domain/wallet"
	cpsUserOutbound "github.com/CBE-Super-App/cbe-super-app-cps-action/internal/port/outbound"
	"github.com/CBE-Super-App/cbe-super-app-cps-action/internal/port/outbound/account_block"
	"github.com/CBE-Super-App/cbe-super-app-cps-action/internal/port/outbound/ad"
	amount_based_auth "github.com/CBE-Super-App/cbe-super-app-cps-action/internal/port/outbound/amount_based_auth"
	bank "github.com/CBE-Super-App/cbe-super-app-cps-action/internal/port/outbound/bank"
	bulk_outbound "github.com/CBE-Super-App/cbe-super-app-cps-action/internal/port/outbound/bulk_service"
	outbound "github.com/CBE-Super-App/cbe-super-app-cps-action/internal/port/outbound/customer"
	fayda_account_repo "github.com/CBE-Super-App/cbe-super-app-cps-action/internal/port/outbound/fayda_account"
	inappnotification "github.com/CBE-Super-App/cbe-super-app-cps-action/internal/port/outbound/in-app-notification"
	"gitlab.com/bersufekadgetachew/cbe-super-app-shared/shared/config"
	"gitlab.com/bersufekadgetachew/cbe-super-app-shared/shared/utils"
	"go.mongodb.org/mongo-driver/v2/mongo"
)

type Persitence struct {
	CustomerPersistence          outbound.CustomerRepository
	FeedBackPersistence          *feedback_repo.FeedbackRepo
	UnlinkPersistence            unlink_repo.UnlinkAccount
	BudgetPersistence            *budget_repo.BudgetPersistence
	AccountPersistence           *account_validation.AccountValidationRepo
	CPSUserPersistence           cpsUserOutbound.OutboundInfra
	PasswordRulesPersistence     cpsUserOutbound.OutboundPasswordRuleInfra
	BankPersistance              bank.BankPersistence
	DepartmentPersistence        *dept_repo.DepartmentPersistence
	PermissionPersistence        *perm_repo.PermissionPersistence
	advertPersistence            ad.ADRepository
	avatarPersitence             avatar.AvatarRepository
	AccountBlockPersistance      account_block.AccountBlockOutboundPort
	HQPersistence                *hq_persistence.HQPersistence
	BPSUserPersistence           *bps_user_persistence.BpsPersistence
	AmountBasedAuthPersistence   amount_based_auth.AmountBasedAuthRepo
	PortalCardPersistance        portalCardRepo.PortaCardInterface
	WalletPersistance            wallet.WalletRepository
	FaydaPersistence             fayda_account_repo.FaydaRepository
	miniAppPersistance           miniApp_domain.MiniRepository
	EventPersistence             event_domain.EventRepository
	CPSActionsPersistance        cps_actions.CPSActionRepository
	BudgetCategoryPersistence    budget_category_repo.BudgetCategoryRepoInterface
	MiniAppMerchantPersisitenct  miniApp_merchant_domain.MiniAppMerchantRepository
	AccounLookUp                 account_lookup_domain.UserSearchRepository
	ServicePersistence           serviceDomain.ServiceRepo
	BulkServicePersistence       bulk_outbound.BulkServiceRepository
	NotificationPersisitence     notification_domain.NotificationRepository
	ProductCodePersistenct       productcode.Repository
	DonationPersistence          donation.DonationRepository
	BulkServicesPersistence      bulk_outbound.BulkServiceRepository
	InAppNotificationPersistence inappnotification.InAppNotificationRepository
}

func InitPersistence(client *mongo.Client, databaseName string, logger utils.Logger, cfg *config.VaultConfig) Persitence {
	return Persitence{
		advertPersistence:       advert.InitAD(client, databaseName, "adverts", logger),
		CustomerPersistence:     customer_repo.InitCustomerDetail(client, databaseName, "user", logger),
		FeedBackPersistence:     feedback_repo.InitFeedback(client, databaseName, "feedback", logger),
		UnlinkPersistence:       unlink_repo.NewUnlinkPersistence(client, databaseName, []string{"cps_actions", "user", "archived_users", "archived_linked_account", "linked_account"}, logger),
		BulkServicesPersistence: bulk_service_repo.InitBulkServicePersistence(client, databaseName, []string{"cps_actions", "access_list"}, logger),
		avatarPersitence:        avatarPersitence.InitAvatarPersistence(client, databaseName, "avatars", logger),
		BudgetPersistence:       budget_repo.InitBudget(client, databaseName, []string{"icons", "colors", "cps_actions"}, logger),
		AccountPersistence:      account_validation.InitAccountValidationPersistence(client, databaseName, 5*time.Second, logger),
		CPSUserPersistence: outboundStore.NewCPSUserPersistence(client, databaseName, []string{
			"cps_users",
			"cps_actions",
			"branch_user",
			"CPSServices",
			"portal_cards",
			"validation_rules",
			"department",
		}),
		PasswordRulesPersistence: outboundStore.NewOutboundPasswordRuleInfra(client, databaseName, []string{
			"password_rules",
			"cps_actions",
		}),

		BankPersistance: bank_repo.InitBank(client, databaseName, []string{"banks", "cps_actions"}, logger),
		DepartmentPersistence: dept_repo.NewDepartmentPersistence(
			dal.NewMongoDal[dept_entities.Department, dept_entities.Department](client, databaseName, "department"),
			logger,
		),
		PermissionPersistence: perm_repo.InitPermission(client, databaseName, 5*time.Second, logger),
		AccountBlockPersistance: account_block_repo.NewOutboundAccountBlockStore(
			client,
			databaseName,
			"branches",
			"regions",
			"cps_actions",
			"districts",
			"user",
			"cities",
			logger,
		),
		HQPersistence:              hq_persistence.NewHQPersistence(client, databaseName, 5*time.Second, logger),
		BPSUserPersistence:         bps_user_persistence.NewBpsPersistence(client, databaseName, 5*time.Second, logger),
		AmountBasedAuthPersistence: amount_based_persistence.InitAmountBasedAuth(client, databaseName, []string{"auth_tier", "cps_actions"}, logger),
		PortalCardPersistance:      portal_card_persistence.InitPortalCardPersistence(client, databaseName, "portal_cards", logger),
		WalletPersistance:          wallet_repo.InitWalletPersistence(client, databaseName, "wallets", logger),
		FaydaPersistence:           faydaaccount.InitFaydaAccountPersistence(client, databaseName, []string{"cps_actions", "users"}, logger),
		miniAppPersistance:         miniApp_persistance.InitMiniAppPersistence(client, databaseName, []string{"mini_app", "cps_actions"}, logger),
		EventPersistence:           event_persistence.InitEventPersistence(client, databaseName, "events", logger),
		// BudgetCategoryPersistence:  budget_category_repo.NewBudgetCategoryRepo(client, databaseName, logger),
		CPSActionsPersistance:       cps_Actions_repo.NewOutBoundStore(client, databaseName, "cps_actions", logger),
		BudgetCategoryPersistence:   budget_category_repo.NewBudgetCategoryRepo(client, databaseName, logger),
		MiniAppMerchantPersisitenct: miniApp_merchant_persistance.NewMiniAppMerchantPersistence(client, databaseName, "mini_app_merchant", logger),
		AccounLookUp:                account_lookup_impl.NewCBEUserSearchClient(cfg.CBEBaseURL),
		ServicePersistence:          servicePersistence.NewServicePersistence(client, databaseName, []string{"cps_actions", "services", "hq"}, logger),
		BulkServicePersistence:      bulk_service_repo.InitBulkServicePersistence(client, databaseName, []string{"cps_actions", "access_list"}, logger),
		//ServicePersistence:          service_persist.NewServicePersistence(client, databaseName, []string{"cps_actions", "services", "hq"}, logger),
		NotificationPersisitence:     notification_persistence.InitNotificationPersistence(client, databaseName, "notifications", logger),
		ProductCodePersistenct:       productcode_persistence.InitProductCodePersistence(client, databaseName, "services", logger),
		DonationPersistence:          donation_persistence.InitDonationPersistence(client, databaseName, []string{"donations", "donation_categories", "donation_companies"}, logger),
		InAppNotificationPersistence: inappnotification_persistence.InitInAppNotificationPersistence(client, databaseName, "in_app_notifications", logger),
	}
}
