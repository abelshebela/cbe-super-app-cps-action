package avatar

import (
	"context"
	"fmt"
	"net/http"
	"time"

	"gitlab.com/bersufekadgetachew/cbe-super-app-cps-action/internal/adapter/outbound/model"
	dto "gitlab.com/bersufekadgetachew/cbe-super-app-cps-action/internal/domain/avatar"
	"gitlab.com/bersufekadgetachew/cbe-super-app-cps-action/internal/port/outbound/avatar"
	constant "gitlab.com/bersufekadgetachew/cbe-super-app-cps-action/utils"
	"gitlab.com/bersufekadgetachew/cbe-super-app-shared/shared/dal"
	"gitlab.com/bersufekadgetachew/cbe-super-app-shared/shared/utils"

	"go.mongodb.org/mongo-driver/v2/bson"
	"go.mongodb.org/mongo-driver/v2/mongo"
)

type AvatarPersistence struct {
	cpsActionDal dal.MongoDal[model.CpsAction, model.CpsAction]
	avatarDal    dal.MongoDal[dto.Avatar, dto.Avatar]
	logger       utils.Logger
}

func InitAvatarPersistence(client *mongo.Client, dbName string, collections []string, logger utils.Logger) avatar.AvatarOutbound {
	cpsActionDal := dal.NewMongoDal[model.CpsAction, model.CpsAction](client, dbName, collections[0])
	avatarDal := dal.NewMongoDal[dto.Avatar, dto.Avatar](client, dbName, collections[1])
	return &AvatarPersistence{
		cpsActionDal: cpsActionDal,
		avatarDal:    avatarDal,
		logger:       logger,
	}
}

func (a *AvatarPersistence) CPSActionExists(ctx context.Context, cpsReq model.CreateCPSAction) error {
	filter := bson.M{
		"maker_user.phone_number": cpsReq.MakerUser.PhoneNumber,
		"status":                  model.ActionPending,
		"department":              cpsReq.Department,
		"request_action":          cpsReq.RequestAction,
	}

	projection := bson.M{
		"action_code": 1,
	}

	exists, err := a.cpsActionDal.FindOne(ctx, filter, projection)
	if err != nil && err != mongo.ErrNoDocuments {
		a.logger.Errorf("failed to get avatar", err)
		err = fmt.Errorf("failed to get bank %w", constant.ErrorDefinition{
			Code:    http.StatusInternalServerError,
			Message: "internal server error",
		})
		return err
	} else if exists != nil {
		err = fmt.Errorf("pending cps already present %w", constant.ErrorDefinition{
			Code:    http.StatusBadRequest,
			Message: "already pending cps action present",
		})
		a.logger.Infof("pending cps action present", cpsReq.MakerUser.FullName,
			cpsReq.MakerUser.UserCode, cpsReq.Department)
		return err
	}

	return nil
}

func (a *AvatarPersistence) CreateAvatar(ctx context.Context, cpsActionReq model.CreateCPSAction) (model.CpsAction, error) {
	cpsAction, err := a.cpsActionDal.InsertOne(ctx, model.CpsAction{
		ID:              bson.NewObjectID().Hex(),
		ActionCode:      utils.RandomGenerator(20),
		MakerUser:       cpsActionReq.MakerUser,
		Department:      cpsActionReq.Department,
		Status:          model.ActionPending,
		RequestAction:   model.RequestCreateAvatar,
		ActionType:      model.ActionCreate,
		ActionData:      cpsActionReq.ActionData,
		MakerActionTime: time.Now(),
	})

	if err != nil {
		a.logger.Errorf("failed to create cps action", err)
		err = fmt.Errorf("failed to create cps action %w", constant.ErrorDefinition{
			Code:    http.StatusInternalServerError,
			Message: "internal server error",
		})
		return model.CpsAction{}, err
	}

	return cpsAction, nil
}
