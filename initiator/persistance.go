package initiator

import (
	"github.com/CBE-Super-App/cbe-super-app-member-auth/internal/storage"
	"github.com/CBE-Super-App/cbe-super-app-member-auth/internal/storage/persistance/users"
	"gitlab.com/bersufekadgetachew/cbe-super-app-shared/shared/utils"

	"go.mongodb.org/mongo-driver/v2/mongo"
)

type Persistence struct {
	UserPersistence storage.UserRepository
	HQPersistence   storage.HQRepository
	OTPPersistence  storage.OTPRepository
}

func InitPersistanceLayer(client *mongo.Client, dbName string, logger utils.Logger) Persistence {
	return Persistence{
		UserPersistence: users.NewUserRepository(client, dbName, "users", logger),
	}
}
