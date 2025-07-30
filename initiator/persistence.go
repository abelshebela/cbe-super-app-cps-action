package initiator

import (
	"time"

	advert "github.com/CBE-Super-App/cbe-super-app-cps-action/internal/adapter/outbound/persistence/ad"
	faydaaccount "github.com/CBE-Super-App/cbe-super-app-cps-action/internal/adapter/outbound/persistence/fayda_account"

	outboundStore "github.com/CBE-Super-App/cbe-super-app-cps-action/internal/adapter/outbound"
	account_validation "github.com/CBE-Super-App/cbe-super-app-cps-action/internal/adapter/outbound/persistence/account_validation"
	bank_repo "github.com/CBE-Super-App/cbe-super-app-cps-action/internal/adapter/outbound/persistence/bank"
	budget_repo "github.com/CBE-Super-App/cbe-super-app-cps-action/internal/adapter/outbound/persistence/budget"
	customer_repo "github.com/CBE-Super-App/cbe-super-app-cps-action/internal/adapter/outbound/persistence/customer"
	outbound "github.com/CBE-Super-App/cbe-super-app-cps-action/internal/port/outbound/customer"

	feedback_repo "github.com/CBE-Super-App/cbe-super-app-cps-action/internal/adapter/outbound/persistence/feedback"
	unlink_repo "github.com/CBE-Super-App/cbe-super-app-cps-action/internal/adapter/outbound/unlink"
	cpsUserOutbound "github.com/CBE-Super-App/cbe-super-app-cps-action/internal/port/outbound"
	"github.com/CBE-Super-App/cbe-super-app-cps-action/internal/port/outbound/ad"

	bank "github.com/CBE-Super-App/cbe-super-app-cps-action/internal/port/outbound/bank"

	dept_repo "github.com/CBE-Super-App/cbe-super-app-cps-action/internal/adapter/outbound/persistence/department"
	perm_repo "github.com/CBE-Super-App/cbe-super-app-cps-action/internal/adapter/outbound/persistence/permission"
	service_repo "github.com/CBE-Super-App/cbe-super-app-cps-action/internal/adapter/outbound/persistence/service"
	wallet_repo "github.com/CBE-Super-App/cbe-super-app-cps-action/internal/adapter/outbound/persistence/wallet"

	portalCardRepo "github.com/CBE-Super-App/cbe-super-app-cps-action/internal/domain/portal_card"

	"github.com/CBE-Super-App/cbe-super-app-cps-action/internal/port/outbound/account_block"
	bulkOutbound "github.com/CBE-Super-App/cbe-super-app-cps-action/internal/port/outbound/bulk_services"

	account_block_repo "github.com/CBE-Super-App/cbe-super-app-cps-action/internal/adapter/outbound/persistence/account_block"
	portal_card_persistence "github.com/CBE-Super-App/cbe-super-app-cps-action/internal/adapter/outbound/persistence/portal_card"
	"gitlab.com/bersufekadgetachew/cbe-super-app-shared/shared/config"
	"gitlab.com/bersufekadgetachew/cbe-super-app-shared/shared/utils"
	"go.mongodb.org/mongo-driver/v2/mongo"

	amount_based_persistence "github.com/CBE-Super-App/cbe-super-app-cps-action/internal/adapter/outbound/persistence/amount_based_auth"
	avatarPersitence "github.com/CBE-Super-App/cbe-super-app-cps-action/internal/adapter/outbound/persistence/avatar"
	hq_persistence "github.com/CBE-Super-App/cbe-super-app-cps-action/internal/adapter/outbound/persistence/hq"
	"github.com/CBE-Super-App/cbe-super-app-cps-action/internal/domain/service"
	amount_based_auth "github.com/CBE-Super-App/cbe-super-app-cps-action/internal/port/outbound/amount_based_auth"
	"github.com/CBE-Super-App/cbe-super-app-cps-action/internal/port/outbound/avatar"
	fayda_account_repo "github.com/CBE-Super-App/cbe-super-app-cps-action/internal/port/outbound/fayda_account"

	event_persistence "github.com/CBE-Super-App/cbe-super-app-cps-action/internal/adapter/outbound/persistence/event"
	miniApp_persistance "github.com/CBE-Super-App/cbe-super-app-cps-action/internal/adapter/outbound/persistence/mini_app"
	miniApp_domain "github.com/CBE-Super-App/cbe-super-app-cps-action/internal/domain/miniapp"

	budget_category_repo "github.com/CBE-Super-App/cbe-super-app-cps-action/internal/adapter/outbound/persistence/budget_category"
	cps_Actions_repo "github.com/CBE-Super-App/cbe-super-app-cps-action/internal/adapter/outbound/persistence/cps_action"
	cps_actions "github.com/CBE-Super-App/cbe-super-app-cps-action/internal/domain/cps_actions/repository"

	miniApp_merchant_persistance "github.com/CBE-Super-App/cbe-super-app-cps-action/internal/adapter/outbound/persistence/mini_app_merchant"
	miniApp_merchant_domain "github.com/CBE-Super-App/cbe-super-app-cps-action/internal/domain/miniapp_merchant"
	wallet "github.com/CBE-Super-App/cbe-super-app-cps-action/internal/domain/wallet/repository"

	account_lookup_impl "github.com/CBE-Super-App/cbe-super-app-cps-action/internal/adapter/outbound/bps_calls"
	bulk_service_repo "github.com/CBE-Super-App/cbe-super-app-cps-action/internal/adapter/outbound/persistence/updated_bulk_service"
	account_lookup_domain "github.com/CBE-Super-App/cbe-super-app-cps-action/internal/domain/account_lookup"
	bulk_outbound "github.com/CBE-Super-App/cbe-super-app-cps-action/internal/port/outbound/updated_bulk_service"
)

type Persitence struct {
	CustomerPersistence         outbound.CustomerRepository
	FeedBackPersistence         *feedback_repo.FeedbackRepo
	UnlinkPersistence           *unlink_repo.UnlinkRepo
	BudgetPersistence           *budget_repo.BudgetPersistence
	AccountPersistence          *account_validation.AccountValidationRepo
	BulkServicesPersistence     bulkOutbound.OutboundInfra
	CPSUserPersistence          cpsUserOutbound.OutboundInfra
	PasswordRulesPersistence    cpsUserOutbound.OutboundPasswordRuleInfra
	BankPersistance             bank.BankPersistence
	DepartmentPersistence       *dept_repo.DepartmentPersistence
	PermissionPersistence       *perm_repo.PermissionPersistence
	advertPersistence           ad.ADRepo
	avatarPersitence            avatar.AvatarOutbound
	AccountBlockPersistance     account_block.AccountBlockOutboundPort
	HQPersistence               *hq_persistence.HQPersistence
	AmountBasedAuthPersistence  amount_based_auth.AmountBasedAuthRepo
	ServiceDetailsStore         service.ServiceRepository
	ServicePersistance          service_repo.ServiceFeePersistence
	PortalCardPersistance       portalCardRepo.PortaCardInterface
	WalletPersistance           wallet.WalletRepository
	FaydaPersistence            fayda_account_repo.FaydaRepository
	miniAppPersistance          miniApp_domain.MiniRepository
	EventPersistence            *event_persistence.EventPersistence
	CPSActionsPersistance       cps_actions.CPSActionRepository
	BudgetCategoryPersistence   budget_category_repo.BudgetCategoryRepoInterface
	MiniAppMerchantPersisitenct miniApp_merchant_domain.MiniAppMerchantRepository
	AccounLookUp                account_lookup_domain.UserSearchRepository
	BulkServicePersistence      bulk_outbound.BulkServiceRepository
}

func InitPersistence(client *mongo.Client, databaseName string, logger utils.Logger, cfg *config.VaultConfig) Persitence {
	collectionNames := []string{
		"bps_user",
		"cps_actions",
		"cps_users",
		"service",
		"member",
		"linked_account",
		"mini_app",
		"portal_card",
	}
	return Persitence{
		advertPersistence: advert.InitAD(client, databaseName, []string{"adverts", "cps_actions"}, logger),
		avatarPersitence:  avatarPersitence.InitAvatarPersistence(client, databaseName, []string{"cps_actions", "avatars"}, logger),

		CustomerPersistence:     customer_repo.InitCustomerDetail(client, databaseName, "user", logger),
		FeedBackPersistence:     feedback_repo.InitFeedback(client, databaseName, "feedbacks", logger),
		UnlinkPersistence:       unlink_repo.NewUnlinkInfrastructure(client, databaseName, []string{"user", "otp", "cps_actions"}, logger),
		BudgetPersistence:       budget_repo.InitBudget(client, databaseName, []string{"icons", "colors", "cps_actions"}, logger),
		AccountPersistence:      account_validation.InitAccountValidationPersistence(client, databaseName, 5*time.Second, logger),
		BulkServicesPersistence: outboundStore.NewOutBoundStore(client, databaseName, collectionNames, logger, cfg),
		CPSUserPersistence: outboundStore.NewCPSUserPersistence(client, databaseName, []string{
			"cps_users",
			"cps_actions",
			"BPSUsers",
			"CPSServices",
			"portal_cards",
			"validation_rules",
			"department",
		}),
		PasswordRulesPersistence: outboundStore.NewOutboundPasswordRuleInfra(client, databaseName, []string{
			"password_rules",
			"cps_actions",
		}),

		BankPersistance:       bank_repo.InitBank(client, databaseName, []string{"banks", "cps_actions"}, logger),
		DepartmentPersistence: dept_repo.InitDepartment(client, databaseName, 5*time.Second, logger),
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
		ServiceDetailsStore:        outboundStore.NewServiceDetailsPersistence(client, databaseName, []string{"cps_actions", "service"}, logger),
		ServicePersistance:         *service_repo.NewServiceFeePersistence(client, databaseName, []string{"cps_actions", "service"}, logger),
		HQPersistence:              hq_persistence.NewHQPersistence(client, databaseName, 5*time.Second, logger),
		AmountBasedAuthPersistence: amount_based_persistence.InitAmountBasedAuth(client, databaseName, []string{"auth_tier", "cps_actions"}, logger),
		PortalCardPersistance:      portal_card_persistence.InitPortalCardPersistence(client, databaseName, "portal_cards", logger),
		WalletPersistance:          wallet_repo.InitWalletPersistence(client, databaseName, []string{"wallets", "cps_actions"}, logger),
		FaydaPersistence:           faydaaccount.InitFaydaAccountPersistence(client, databaseName, []string{"cps_actions", "user"}, logger),
		miniAppPersistance:         miniApp_persistance.InitMiniAppPersistence(client, databaseName, []string{"mini_app", "cps_actions"}, logger),
		EventPersistence:           event_persistence.InitEventPersistence(client, databaseName, []string{"events", "cps_actions"}, logger),
		// BudgetCategoryPersistence:  budget_category_repo.NewBudgetCategoryRepo(client, databaseName, logger),
		CPSActionsPersistance:       cps_Actions_repo.NewOutBoundStore(client, databaseName, "cps_actions", logger),
		BudgetCategoryPersistence:   budget_category_repo.NewBudgetCategoryRepo(client, databaseName, logger),
		MiniAppMerchantPersisitenct: miniApp_merchant_persistance.NewMiniAppMerchantPersistence(client, databaseName, "mini_app_merchant", logger),
		AccounLookUp:                account_lookup_impl.NewCBEUserSearchClient(cfg.CBEBaseURL),
		BulkServicePersistence:      bulk_service_repo.InitBulkServicePersistence(client, databaseName, []string{"cps_actions", "access_list"}, logger),
	}
}
