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

	"github.com/CBE-Super-App/cbe-super-app-cps-action/internal/port/outbound/department"

	"github.com/CBE-Super-App/cbe-super-app-cps-action/internal/port/outbound/account_block"
	bulkOutbound "github.com/CBE-Super-App/cbe-super-app-cps-action/internal/port/outbound/bulk_services"

	account_block_repo "github.com/CBE-Super-App/cbe-super-app-cps-action/internal/adapter/outbound/persistence/account_block"
	portal_card_persistence "github.com/CBE-Super-App/cbe-super-app-cps-action/internal/adapter/outbound/persistence/portal_card"
	"gitlab.com/bersufekadgetachew/cbe-super-app-shared/shared/utils"
	"go.mongodb.org/mongo-driver/v2/mongo"

	amount_based_persistence "github.com/CBE-Super-App/cbe-super-app-cps-action/internal/adapter/outbound/persistence/amount_based_auth"
	avatarPersitence "github.com/CBE-Super-App/cbe-super-app-cps-action/internal/adapter/outbound/persistence/avatar"
	hq_persistence "github.com/CBE-Super-App/cbe-super-app-cps-action/internal/adapter/outbound/persistence/hq"
	"github.com/CBE-Super-App/cbe-super-app-cps-action/internal/domain/service"
	amount_based_auth "github.com/CBE-Super-App/cbe-super-app-cps-action/internal/port/outbound/amount_based_auth"
	"github.com/CBE-Super-App/cbe-super-app-cps-action/internal/port/outbound/avatar"
	fayda_account_repo "github.com/CBE-Super-App/cbe-super-app-cps-action/internal/port/outbound/fayda_account"
	"github.com/CBE-Super-App/cbe-super-app-cps-action/internal/port/outbound/wallet"
)

type Persitence struct {
	CustomerPersistence        *customer_repo.CustomerDetailRepo
	FeedBackPersistence        *feedback_repo.FeedbackRepo
	UnlinkPersistence          *unlink_repo.UnlinkRepo
	BudgetPersistence          *budget_repo.BudgetPersistence
	AccountPersistence         *account_validation.AccountValidationRepo
	BulkServicesPersistence    bulkOutbound.OutboundInfra
	CPSUserPersistence         cpsUserOutbound.OutboundInfra
	PasswordRulesPersistence   cpsUserOutbound.OutboundPasswordRuleInfra
	BankPersistance            bank.BankPersistence
	DepartmentPersistence      *dept_repo.DepartmentPersistence
	PermissionPersistence      *perm_repo.PermissionPersistence
	CPSActionPersistance       department.CPSActionRepository
	advertPersistence          ad.ADRepo
	avatarPersitence           avatar.AvatarOutbound
	AccountBlockPersistance    account_block.AccountBlockOutboundPort
	HQPersistence              *hq_persistence.HQPersistence
	AmountBasedAuthPersistence amount_based_auth.AmountBasedAuthRepo
	ServiceDetailsStore        service.ServiceRepository
	ServicePersistance         service_repo.ServiceFeePersistence
	PortalCardPersistance      portalCardRepo.PortaCardInterface
	WalletPersistance          wallet.WalletPersistence
	FaydaPersistence           fayda_account_repo.FaydaRepository
}

func InitPersistence(client *mongo.Client, database_name string, logger utils.Logger) Persitence {
	collectionNames := []string{
		"bps_user",
		"cps_actions",
		"cps_users",
		"service",
		"member",
		"linked_accounts",
		"mini_app",
		"portal_card",
	}
	return Persitence{
		advertPersistence: advert.InitAD(client, database_name, []string{"adverts", "cps_actions"}, logger),
		avatarPersitence:  avatarPersitence.InitAvatarPersistence(client, database_name, []string{"cps_actions", "avatars"}, logger),

		CustomerPersistence:     customer_repo.InitCustomerDetail(client, database_name, "user", logger),
		FeedBackPersistence:     feedback_repo.InitFeedback(client, database_name, "feedbacks", logger),
		UnlinkPersistence:       unlink_repo.NewUnlinkInfrastructure(client, database_name, []string{"user", "otp", "cps_actions"}, logger),
		BudgetPersistence:       budget_repo.InitBudget(client, database_name, []string{"icons", "colors", "cps_actions"}, logger),
		AccountPersistence:      account_validation.InitAccountValidationPersistence(client, database_name, 5*time.Second, logger),
		BulkServicesPersistence: outboundStore.NewOutBoundStore(client, database_name, collectionNames, logger),
		CPSUserPersistence: outboundStore.NewCPSUserPersistence(client, database_name, []string{
			"cps_users",
			"cps_actions",
			"BPSUsers",
			"CPSServices",
			"portal_cards",
			"validation_rules",
		}),
		PasswordRulesPersistence: outboundStore.NewOutboundPasswordRuleInfra(client, database_name, []string{
			"password_rules",
			"cps_actions",
		}),

		BankPersistance:       bank_repo.InitBank(client, database_name, []string{"banks", "cps_actions"}, logger),
		DepartmentPersistence: dept_repo.InitDepartment(client, database_name, 5*time.Second, logger),
		PermissionPersistence: perm_repo.InitPermission(client, database_name, 5*time.Second, logger),
		AccountBlockPersistance: account_block_repo.NewOutboundAccountBlockStore(
			client,
			database_name,
			"branches",
			"regions",
			"cps_actions",
			"districts",
			"user",
			"cities",
			logger,
		),
		ServiceDetailsStore:        outboundStore.NewServiceDetailsPersistence(client, database_name, []string{"service", "cps_actions"}, logger),
		ServicePersistance:         *service_repo.NewServiceFeePersistence(client, database_name, []string{"cps_actions", "service"}, logger),
		HQPersistence:              hq_persistence.NewHQPersistence(client, database_name, 5*time.Second, logger),
		AmountBasedAuthPersistence: amount_based_persistence.InitAmountBasedAuth(client, database_name, []string{"auth_tier", "cps_actions"}, logger),
		PortalCardPersistance:      portal_card_persistence.InitPortalCardPersistence(client, database_name, "portal_cards", logger),
		WalletPersistance:          wallet_repo.InitWalletPersistence(client, database_name, []string{"wallets", "cps_actions"}, logger),
		FaydaPersistence:           faydaaccount.InitFaydaAccountPersistence(client, database_name, []string{"cps_actions", "user"}, logger),
	}
}
