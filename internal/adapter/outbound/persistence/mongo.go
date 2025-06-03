package persistence
import (
	"context"
	"go.mongodb.org/mongo-driver/v2/bson"
	"go.mongodb.org/mongo-driver/v2/mongo"
	"gitlab.com/bersufekadgetachew/cbe-super-app-shared/shared/dal"
	"gitlab.com/bersufekadgetachew/cbe-super-app-shared/shared/entities"
	"cbe-super-app-member-users/internal/domain/users"
	"cbe-super-app-member-users/internal/port/outbound"
)

type MongoRepository struct {
	userDal          dal.MongoDal[entities.User, entities.User]
	linkedAccountDal dal.MongoDal[entities.LinkedAccount, entities.LinkedAccount]
}

func NewMongoRepository(client *mongo.Client, dbName string) outbound.UserRepository {
	return &MongoRepository{
		userDal:          dal.NewMongoDal[entities.User, entities.User](client, dbName, "users"),
		linkedAccountDal: dal.NewMongoDal[entities.LinkedAccount, entities.LinkedAccount](client, dbName, "linked_accounts"),
	}
}

func (r *MongoRepository) FindByID(ctx context.Context, id string) (*users.User, error) {
	oid, err := bson.ObjectIDFromHex(id)
	if err != nil {
		return nil, err
	}

	userFilter := bson.M{
		"_id":        oid,
		"is_deleted": bson.M{"$ne": true},
	}
	userEntity, err := r.userDal.FindOne(ctx, userFilter, nil)
	if err != nil {
		if err == mongo.ErrNoDocuments {
			return nil, err
		}
		return nil, err
	}
	if userEntity == nil {
		return nil, mongo.ErrNoDocuments
	}

	return &users.User{
		ID: userEntity.ID.Hex(),
		FullName: users.FullName{
			FirstName:  userEntity.FullName.FirstName,
			MiddleName: userEntity.FullName.MiddleName,
			LastName:   userEntity.FullName.LastName,
		},
		IsDeleted: userEntity.IsDeleted,
	}, nil
}

func (r *MongoRepository) FindActiveLinkedAccounts(ctx context.Context, userID string) ([]users.LinkedAccountDetail, error) {
	oid, err := bson.ObjectIDFromHex(userID)
	if err != nil {
		return nil, err
	}

	accountFilter := bson.M{
		"user_id":          oid,
		"is_account_active": true,
		"linked_status":    true,
		"is_deleted":       bson.M{"$ne": true},
	}
	accountEntities, err := r.linkedAccountDal.Find(ctx, accountFilter, nil)
	if err != nil {
		return nil, err
	}

	linkedAccounts := make([]users.LinkedAccountDetail, 0, len(accountEntities))
	for _, acc := range accountEntities {
		linkedAccounts = append(linkedAccounts, users.LinkedAccountDetail{
			AccountNumber:     acc.AccountNumber,
			AccountBranchCode: acc.AccountBranchCode,
			LinkedBranch:      acc.LinkerBranch,
			IsAccountActive:   acc.IsAccountActive,
			LinkedStatus:      acc.LinkedStatus,
			CurrencyCode:      acc.CurrencyCode,
		})
	}

	return linkedAccounts, nil
}