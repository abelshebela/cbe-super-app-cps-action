package bank

import (
	"context"
	"errors"
	"fmt"
	"time"

	"github.com/CBE-Super-App/cbe-super-app-cps-action/internal/adapter/outbound/model"
	"github.com/CBE-Super-App/cbe-super-app-cps-action/internal/domain/bank/entity"
	"gitlab.com/bersufekadgetachew/cbe-super-app-shared/shared/utils"
	"go.mongodb.org/mongo-driver/v2/bson"
	"go.mongodb.org/mongo-driver/v2/mongo"

	error_codes "github.com/CBE-Super-App/cbe-super-app-cps-action/pkgs/utils"
)

func (b *Bank) parseObjectID(id string) (bson.ObjectID, error) {
	objectID, err := bson.ObjectIDFromHex(id)
	if err != nil {
		return bson.ObjectID{}, fmt.Errorf(error_codes.InvalidID)
	}
	return objectID, nil
}

func (b *Bank) toDomain(doc *entity.BankDocument) *entity.Bank {
	return &entity.Bank{
		ID:             doc.ID.Hex(),
		Name:           doc.Name,
		Logo:           doc.Logo,
		Code:           doc.Code,
		BIC:            doc.BIC,
		Enabled:        doc.Enabled,
		IsDeleted:      doc.IsDeleted,
		CreatedAt:      doc.CreatedAt,
		LastModifiedAt: doc.LastModifiedAt,
	}
}

func (b *Bank) toDocument(domain *entity.BankDocument) (*entity.BankDocument, error) {

	return &entity.BankDocument{
		ID:             domain.ID,
		Name:           domain.Name,
		Logo:           domain.Logo,
		Code:           domain.Code,
		BIC:            domain.BIC,
		Enabled:        domain.Enabled,
		IsDeleted:      domain.IsDeleted,
		CreatedAt:      domain.CreatedAt,
		LastModifiedAt: domain.LastModifiedAt,
	}, nil
}

func (b *Bank) findBankByID(ctx context.Context, id string, projection bson.M) (*entity.BankDocument, error) {
	objectID, err := b.parseObjectID(id)
	if err != nil {
		b.logger.Errorf("invalid id provided: %v", err)
		return nil, fmt.Errorf(error_codes.InvalidID)
	}

	filter := bson.M{"_id": objectID, "is_deleted": false}

	bank, err := b.bankDal.FindOne(ctx, filter, projection)
	if err != nil {
		if errors.Is(err, mongo.ErrNoDocuments) {
			b.logger.Errorf("bank not found: %v", err)
			return nil, fmt.Errorf(error_codes.BankNotFound)
		}
		b.logger.Errorf("failed to get bank: %v", err)
		return nil, fmt.Errorf(error_codes.UnhandledServerError)
	}

	return bank, nil
}

func (b *Bank) updateBankByID(ctx context.Context, id string, update bson.M) (*entity.BankDocument, error) {
	objectID, err := b.parseObjectID(id)
	if err != nil {
		return nil, fmt.Errorf(error_codes.InvalidID)
	}

	filter := bson.M{"_id": objectID}
	update["last_modified_at"] = time.Now()

	bank, err := b.bankDal.UpdateOne(ctx, filter, update)
	if err != nil {
		b.logger.Errorf("failed to update bank: %v", err)
		return nil, fmt.Errorf(error_codes.UnhandledServerError)
	}

	return &bank, nil
}

func (b *Bank) buildCPSAction(cpsReq model.CreateCPSAction, requestAction model.RequestAction, prevAction any) model.CPSAction {
	return model.CPSAction{
		ID:               bson.NewObjectID(),
		ActionCode:       utils.RandomGenerator(20),
		MakerID:          cpsReq.MakerUser.UserCode,
		MakerName:        cpsReq.MakerUser.FullName,
		MakerPhoneNumber: cpsReq.MakerUser.PhoneNumber,
		Department:       cpsReq.Department,
		ActionStatus:     string(model.ActionPending),
		ActionType:       string(model.ActionUpdate),
		RequestAction:    string(requestAction),
		CurrentAction:    cpsReq.ActionData,
		PreviosAction:    prevAction,
		MakerActionTime:  time.Now(),
		CreatedAt:        time.Now(),
	}
}
