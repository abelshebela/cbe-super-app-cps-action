package initiator

import (
	"time"

	outboundStore "github.com/CBE-Super-App/cbe-super-app-cps-action/internal/adapter/outbound"
	account_validation "github.com/CBE-Super-App/cbe-super-app-cps-action/internal/adapter/outbound/persistence/account_validation"
	bank_repo "github.com/CBE-Super-App/cbe-super-app-cps-action/internal/adapter/outbound/persistence/bank"
	budget_repo "github.com/CBE-Super-App/cbe-super-app-cps-action/internal/adapter/outbound/persistence/budget"
	customer_repo "github.com/CBE-Super-App/cbe-super-app-cps-action/internal/adapter/outbound/persistence/customer"
	feedback_repo "github.com/CBE-Super-App/cbe-super-app-cps-action/internal/adapter/outbound/persistence/feedback"
	unlink_repo "github.com/CBE-Super-App/cbe-super-app-cps-action/internal/adapter/outbound/unlink"
	cpsUserOutbound "github.com/CBE-Super-App/cbe-super-app-cps-action/internal/port/outbound"
	passwordRuleOutbound "github.com/CBE-Super-App/cbe-super-app-cps-action/internal/port/outbound"
	bank "github.com/CBE-Super-App/cbe-super-app-cps-action/internal/port/outbound/bank"
	bulkOutbound "github.com/CBE-Super-App/cbe-super-app-cps-action/internal/port/outbound/bulk_services"

	dept_repo "github.com/CBE-Super-App/cbe-super-app-cps-action/internal/adapter/outbound/persistence/department"
	"github.com/CBE-Super-App/cbe-super-app-cps-action/internal/port/outbound/department"

	"gitlab.com/bersufekadgetachew/cbe-super-app-shared/shared/utils"
	"go.mongodb.org/mongo-driver/v2/mongo"
)

type Persitence struct {
	CustomerPersistence      *customer_repo.CustomerDetailRepo
	FeedBackPersistence      *feedback_repo.FeedbackRepo
	UnlinkPersistence        *unlink_repo.UnlinkRepo
	BudgetPersistence        *budget_repo.BudgetPersistence
	AccountPersistence       *account_validation.AccountValidationRepo
	BulkServicesPersistence  bulkOutbound.OutboundInfra
	CPSUserPersistence       cpsUserOutbound.OutboundInfra
	PasswordRulesPersistence passwordRuleOutbound.OutboundPasswordRuleInfra
	BankPersistance          bank.BankPersistence
	DepartmentPersistence    *dept_repo.DepartmentPersistence
	CPSActionPersistance     department.CPSActionRepository
}

func InitPersistence(client *mongo.Client, database_name string, logger utils.Logger) Persitence {
	collectionNames := []string{
		"bps_user",
		"cps_action",
		"cps_users",
		"service",
		"member",
		"linked_accounts",
		"mini_app",
		"portal_card",
	}
	return Persitence{
		CustomerPersistence:     customer_repo.InitCustomerDetail(client, database_name, "customers", logger),
		FeedBackPersistence:     feedback_repo.InitFeedback(client, database_name, "feedbacks", logger),
		UnlinkPersistence:       unlink_repo.NewUnlinkInfrastructure(client, database_name, []string{"user", "otp", "cps_action"}, logger),
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

		BankPersistance:       bank_repo.InitBank(client, database_name, []string{"cps_actions", "banks"}, logger),
		DepartmentPersistence: dept_repo.InitDepartment(client, database_name, 5*time.Second, logger),
	}
}
