package initiator

import (
	"cbe-super-app-cps-action/internal/storage"
	"database/sql"

	"cbe-super-app-cps-action/internal/storage/kafka"
	access_list_oracle "cbe-super-app-cps-action/internal/storage/persistance/access_list_oracle"
	access_list_segmentation_oracle "cbe-super-app-cps-action/internal/storage/persistance/access_list_segmentation/oracle"
	account_block_repo "cbe-super-app-cps-action/internal/storage/persistance/account_block"
	amount_based_auth_oracle "cbe-super-app-cps-action/internal/storage/persistance/amount_based_auth_oracle"
	"cbe-super-app-cps-action/internal/storage/persistance/bank/gen/sqlc"
	budget_category_oracle "cbe-super-app-cps-action/internal/storage/persistance/budget_category_oracle"
	cpsroles "cbe-super-app-cps-action/internal/storage/persistance/cps_roles"
	customer_oracle "cbe-super-app-cps-action/internal/storage/persistance/customer/oracle"
	customersegmentaion "cbe-super-app-cps-action/internal/storage/persistance/customer_segmentaion"
	donation_oracle "cbe-super-app-cps-action/internal/storage/persistance/donation/oracle"
	donation_category_oracle "cbe-super-app-cps-action/internal/storage/persistance/donation_category/oracle"
	donation_company_oracle "cbe-super-app-cps-action/internal/storage/persistance/donation_company/oracle"
	ecommerce_merchant "cbe-super-app-cps-action/internal/storage/persistance/ecommerce_merchant"
	services_repo "cbe-super-app-cps-action/internal/storage/persistance/services"
	"cbe-super-app-cps-action/internal/storage/persistance/sitota"
	vaultCategory "cbe-super-app-cps-action/internal/storage/persistance/vault"
	wallet_oracle "cbe-super-app-cps-action/internal/storage/persistance/wallet/oracle"
	event_oracle "cbe-super-app-cps-action/internal/storage/persistance/event_merchant_oracle"


	"gitlab.com/bersufekadgetachew/cbe-super-app-shared/shared/config"
	"gitlab.com/bersufekadgetachew/cbe-super-app-shared/shared/utils"
	"go.mongodb.org/mongo-driver/v2/mongo"
)

type OraclePersistence struct {
	Db         *sql.DB
	BankOracle storage.BankOracleRepository
	// vaultTransaction      storage.TransactionRepository
	Vault                 storage.VaultCategoryRepository
	Sitota                storage.SitotaRepository
	ServicesPersistence   storage.ServicesRepository
	CustomerSegmentation  storage.CustomerSegmentationRepository
	NewCPSRolesStorage    storage.CPSRolesRepository
	WalletOracle          storage.WalletOracleRepository
	BudgetCategoryOracle  storage.BudgetCategoryOracleRepository
	AmountBasedAuthOracle storage.AmountBasedAuthOracleRepository
	AccessListOracle      storage.BulkServiceRepository
	AccessListSegmentaion storage.AccessListSegmentationRepositoryOracle
	AccountBlock          storage.AccountBlockRepository
	Customer              storage.CustomerRepository
	EventMerchant            storage.EventMerchantRepository
	EcommerceMerchant     storage.EcommerceMerchantRepository
	DonationCategory      storage.DonationCategoryRepository
	DonationCompany       storage.DonationCompanyRepository
	Donation              storage.DonationRepository
}

func InitOraclePersistence(client *mongo.Client, db *sql.DB, cfg *config.VaultConfig, clientOrchestrationProducer kafka.ClientOrchestrationProducer, accessListSegmentationProducer *kafka.AccessListSegmentationProducer, redisRepository storage.RedisRepository, log utils.Logger) OraclePersistence {
	return OraclePersistence{
		Db:         db,
		BankOracle: sqlc.NewBankRepository(db, log),
		Vault:      vaultCategory.NewVaultCategoryRepository(db, cfg, clientOrchestrationProducer, log),
		Sitota:     sitota.NewSitotaRepository(db, log),
		// vaultTransaction:      transaction_repo.NewTransactionRepository(db, log),
		ServicesPersistence:   services_repo.NewServicesRepository(db, cfg, clientOrchestrationProducer, redisRepository, log),
		CustomerSegmentation:  customersegmentaion.NewCustomerSegmentationRepository(cfg, db, clientOrchestrationProducer, log),
		NewCPSRolesStorage:    cpsroles.NewCPSRolesStorage(cfg, db, clientOrchestrationProducer, log),
		WalletOracle:          wallet_oracle.NewWalletOracleRepository(db, redisRepository, log),
		BudgetCategoryOracle:  budget_category_oracle.NewBudgetCategoryOracleRepository(db, clientOrchestrationProducer, log),
		AmountBasedAuthOracle: amount_based_auth_oracle.NewAmountBasedAuthOracleRepository(db, log),
		AccessListOracle:      access_list_oracle.NewAccessListOracleRepository(db, redisRepository, log),
		AccessListSegmentaion: access_list_segmentation_oracle.NewAccessListSegmentationOracle(db, *cfg, accessListSegmentationProducer, redisRepository, log),
		AccountBlock:          account_block_repo.NewAccountBlockRepository(client, cfg, CPSActionsCollection, db, redisRepository, log),
		Customer:              customer_oracle.NewCustomerOracleRepository(db, *cfg, log),
		EventMerchant:  event_oracle.NewEventMerchantOracleRepository(db,log),
		EcommerceMerchant:     ecommerce_merchant.NewEcommerceMerchantRepository(db, log),
		DonationCategory: donation_category_oracle.NewPublishingRepository(
			donation_category_oracle.NewRepository(db, log),
			clientOrchestrationProducer,
		),
		DonationCompany: donation_company_oracle.NewPublishingRepository(
			donation_company_oracle.NewRepository(db, log),
			clientOrchestrationProducer,
		),
		Donation: donation_oracle.NewPublishingRepository(
			donation_oracle.NewRepository(db, log),
			clientOrchestrationProducer,
		),
	}
}
