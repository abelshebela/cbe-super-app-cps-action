package customer

// import (
// 	"cbe-super-app-cps-action/internal/storage"
// 	"database/sql"

// 	"gitlab.com/bersufekadgetachew/cbe-super-app-shared/shared/config"
// 	"gitlab.com/bersufekadgetachew/cbe-super-app-shared/shared/utils"
// 	"go.mongodb.org/mongo-driver/v2/mongo"
// )

// type selfActivationRepository struct {
// 	customerKYCRepository
// }

// func NewSelfActivationRepository(client *mongo.Client, oracleDB *sql.DB, cfg *config.VaultConfig, dbName string, custKycCollection string, logger utils.Logger) storage.SelfActivationKYCRepository {
// 	base := NewCustomerKYCRepository(client, oracleDB, cfg, dbName, custKycCollection, logger)
// 	return &selfActivationRepository{
// 		customerKYCRepository: *(base.(*customerKYCRepository)),
// 	}
// }
