package mappers

import (
	"github.com/CBE-Super-App/cbe-super-app-cps-action/internal/adapter/outbound/model"
	"github.com/CBE-Super-App/cbe-super-app-cps-action/internal/domain/avatar"

	"go.mongodb.org/mongo-driver/v2/bson"
)

func ToAvatarModel(a *model.AvatarDocument) avatar.Avatar {
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

func ToAvatarDocument(a avatar.Avatar) (*model.AvatarDocument, error) {
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
	return &model.AvatarDocument{
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
