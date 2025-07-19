package avatar

import (
	"time"

	"github.com/CBE-Super-App/cbe-super-app-cps-action/internal/adapter/outbound/model"
	"github.com/CBE-Super-App/cbe-super-app-cps-action/internal/domain/avatar"

	"go.mongodb.org/mongo-driver/v2/bson"
)

type AvatarDocument struct {
	ID             bson.ObjectID `json:"id" bson:"_id"`
	Avatar         string        `json:"avatar" bson:"avatar"`
	Label          string        `json:"label" bson:"label"`
	Enable         bool          `json:"enable" bson:"enable"`
	IsDeleted      bool          `json:"is_deleted" bson:"is_deleted"`
	CreatedAt      time.Time     `json:"created_at" bson:"created_at"`
	LastModifiedAt time.Time     `json:"last_modified_at" bson:"last_modified_at"`
	DeletedAt      *time.Time    `json:"deleted_at" bson:"deleted_at"`
}

func (a *AvatarDocument) toModel() avatar.Avatar {
	return avatar.Avatar{
		ID:             a.ID.Hex(),
		Avatar:         a.Avatar,
		Label:          a.Label,
		Enable:         a.Enable,
		IsDeleted:      a.IsDeleted,
		CreatedAt:      a.CreatedAt,
		LastModifiedAt: a.LastModifiedAt,
		DeletedAt:      a.DeletedAt,
	}
}

func ToAvatarDocument(a avatar.Avatar) (*AvatarDocument, error) {
	var objectID bson.ObjectID
	if a.ID != "" {
		id, err := bson.ObjectIDFromHex(a.ID)
		if err != nil {
			return nil, err
		}

		objectID = id

	} else {
		objectID = bson.NewObjectID()
	}
	return &AvatarDocument{
		ID:             objectID,
		Avatar:         a.Avatar,
		Label:          a.Label,
		Enable:         a.Enable,
		IsDeleted:      a.IsDeleted,
		CreatedAt:      a.CreatedAt,
		LastModifiedAt: a.LastModifiedAt,
		DeletedAt:      a.DeletedAt,
	}, nil
}

func ToCPSAction(cpsAction *model.CPSAction) *avatar.CPSAction {
	if cpsAction == nil {
		return nil
	}

	return &avatar.CPSAction{
		ID:                 cpsAction.ID.Hex(),
		ActionCode:         cpsAction.ActionCode,
		UniqueId:           cpsAction.UniqueId,
		MakerID:            cpsAction.MakerID,
		MakerName:          cpsAction.MakerName,
		MakerPhoneNumber:   cpsAction.MakerPhoneNumber,
		CheckerID:          cpsAction.CheckerID,
		CheckerName:        cpsAction.CheckerName,
		CheckerPhoneNumber: cpsAction.CheckerPhoneNumber,
		Department:         cpsAction.Department,
		RejectionReason:    cpsAction.RejectionReason,
		PreviousAction:     cpsAction.PreviousAction,
		CurrentAction:      cpsAction.CurrentAction,
		ActionStatus:       cpsAction.ActionStatus,
		ActionType:         cpsAction.ActionType,
		RequestAction:      cpsAction.RequestAction,
		CreatedAt:          cpsAction.CreatedAt,
		LastModifiedAt:     cpsAction.LastModifiedAt,
		MakerActionTime:    cpsAction.MakerActionTime,
		CheckerActionTime:  cpsAction.CheckerActionTime,
	}
}
