package initiator

import (
	"cbe-super-app-cps-action/internal/storage"
	"database/sql"

	"cbe-super-app-cps-action/internal/storage/kafka"
	access_list_oracle "cbe-super-app-cps-action/internal/storage/persistance/access_list_oracle"
	access_list_segmentation_oracle "cbe-super-app-cps-action/internal/storage/persistance/access_list_segmentation/oracle"
	account_block_repo "cbe-super-app-cps-action/internal/storage/persistance/account_block"
	ap_oracle "cbe-super-app-cps-action/internal/storage/persistance/account_product/oracle"
	apc_oracle "cbe-super-app-cps-action/internal/storage/persistance/account_product_category/oracle"
	account_sub_type_oracle "cbe-super-app-cps-action/internal/storage/persistance/account_sub_type"
	amount_based_auth_oracle "cbe-super-app-cps-action/internal/storage/persistance/amount_based_auth_oracle"
	"cbe-super-app-cps-action/internal/storage/persistance/bank/gen/sqlc"
	budget_category_oracle "cbe-super-app-cps-action/internal/storage/persistance/budget_category_oracle"
	cpsroles "cbe-super-app-cps-action/internal/storage/persistance/cps_roles"
	customer_oracle "cbe-super-app-cps-action/internal/storage/persistance/customer/oracle"
	customergroup "cbe-super-app-cps-action/internal/storage/persistance/customer_group"
	customer_kyc "cbe-super-app-cps-action/internal/storage/persistance/customer_kyc"
	customersegmentaion "cbe-super-app-cps-action/internal/storage/persistance/customer_segmentaion"
	donation_oracle "cbe-super-app-cps-action/internal/storage/persistance/donation/oracle"
	donation_category_oracle "cbe-super-app-cps-action/internal/storage/persistance/donation_category/oracle"
	donation_company_oracle "cbe-super-app-cps-action/internal/storage/persistance/donation_company/oracle"
	ecommerce_merchant "cbe-super-app-cps-action/internal/storage/persistance/ecommerce_merchant"
	event_oracle "cbe-super-app-cps-action/internal/storage/persistance/event_merchant_oracle"
	logistics_merchant_oracle "cbe-super-app-cps-action/internal/storage/persistance/logistics_merchant/oracle"
	mini_app_repo "cbe-super-app-cps-action/internal/storage/persistance/mini_app/oracle"
	services_repo "cbe-super-app-cps-action/internal/storage/persistance/services"
	"cbe-super-app-cps-action/internal/storage/persistance/sitota"
	superapprole "cbe-super-app-cps-action/internal/storage/persistance/superapp_role"
	tac_oracle "cbe-super-app-cps-action/internal/storage/persistance/term_and_condition/oracle"
	"cbe-super-app-cps-action/internal/storage/persistance/ussd_merchant"
	vaultCategory "cbe-super-app-cps-action/internal/storage/persistance/vault"
	wallet_oracle "cbe-super-app-cps-action/internal/storage/persistance/wallet/oracle"

	"gitlab.com/bersufekadgetachew/cbe-super-app-shared/shared/config"
	"gitlab.com/bersufekadgetachew/cbe-super-app-shared/shared/utils"
	"go.mongodb.org/mongo-driver/v2/mongo"
)

type OraclePersistence struct {
	Db         *sql.DB
	BankOracle storage.BankOracleRepository
	// vaultTransaction      storage.TransactionRepository
	Vault                   storage.VaultCategoryRepository
	Sitota                  storage.SitotaRepository
	ServicesPersistence     storage.ServicesRepository
	CustomerSegmentation    storage.CustomerSegmentationRepository
	NewCPSRolesStorage      storage.CPSRolesRepository
	WalletOracle            storage.WalletOracleRepository
	BudgetCategoryOracle    storage.BudgetCategoryOracleRepository
	AmountBasedAuthOracle   storage.AmountBasedAuthOracleRepository
	AccessListOracle        storage.BulkServiceRepository
	LogisticsMerchantOracle storage.LogisticsMerchantOracleRepository
	AccessListSegmentaion   storage.AccessListSegmentationRepositoryOracle
	AccountBlock            storage.AccountBlockRepository
	Customer                storage.CustomerRepository
	EventMerchant           storage.EventMerchantRepository
	EcommerceMerchant       storage.EcommerceMerchantRepository
	UssdMerchant            storage.UssdMerchantRepository
	DonationCategory        storage.DonationCategoryRepository
	DonationCompany         storage.DonationCompanyRepository
	Donation                storage.DonationRepository
	CustomerKYC             storage.CustomerKYCRepository
	SelfActivationKYC       storage.SelfActivationKYCRepository
	CustomerGroup           storage.CustomerGroupRepository
	SuperAppRole            storage.SuperAppRoleRepository
	MiniApp                 storage.MiniAppRepository
	MiniAppMerchant         storage.MiniAppMerchant
	MiniAppCategory         storage.MiniAppCategoryRepository
	AccountSubType          storage.AccountSubTypeOracleRepository
	AccountProductCategory  storage.AccountProductCategoryRepository
	AccountProduct          storage.AccountProductRepository
	AccountOpeningTerms     storage.AccountOpeningTermsRepository
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
		EventMerchant:         event_oracle.NewEventMerchantOracleRepository(db, log),
		EcommerceMerchant:     ecommerce_merchant.NewEcommerceMerchantRepository(db, log),
		UssdMerchant:          ussd_merchant.NewUssdMerchant(db, log),
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
		CustomerKYC:             customer_kyc.NewCustomerKYCRepository(client, db, cfg, cfg.MongoDBDatabase, CustomersKYCCollection, log),
		SelfActivationKYC:       customer_kyc.NewSelfActivationRepository(client, db, cfg, cfg.MongoDBDatabase, SelfActivationKYCCollection, CPSActionsCollection, log),
		CustomerGroup:           customergroup.NewCustomerGroupRepository(cfg, db, log),
		SuperAppRole:            superapprole.NewSuperAppRoleRepository(db, log),
		LogisticsMerchantOracle: logistics_merchant_oracle.NewLogisticsMerchantOracle(db, cfg, log),
		MiniApp:                 mini_app_repo.NewMiniAppOracleRepository(log, db),
		MiniAppMerchant:         mini_app_repo.NewMiniAppMerchantOracleRepository(log, db),
		MiniAppCategory:         mini_app_repo.NewCategoryOracleRepository(log, db),
		AccountSubType:          account_sub_type_oracle.NewAccountSubTypeOracleRepository(db, log),
		AccountProductCategory:  apc_oracle.NewRepository(db, log),
		AccountProduct:          ap_oracle.NewRepository(db, log),
		AccountOpeningTerms:     tac_oracle.NewRepository(db, log),
	}
}
