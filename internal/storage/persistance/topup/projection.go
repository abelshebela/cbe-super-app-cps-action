package Topup

import (
	"cbe-super-app-cps-action/internal/constants/model"
	"time"

	"go.mongodb.org/mongo-driver/v2/bson"
)

func ToTopupDocument(Topup model.Topup) (*model.Topup, error) {
	if Topup.ID == bson.NilObjectID {
		Topup.ID = bson.NewObjectID()
	}

	if Topup.CreatedAt.IsZero() {
		Topup.CreatedAt = time.Now()
	}
	if Topup.LastModifiedAt.IsZero() {
		Topup.LastModifiedAt = time.Now()
	}
	Topup.IsDeleted = false
	Topup.Enabled = false
	return &Topup, nil
}

func UpdateMapper(Topup model.Topup) bson.M {

	update := bson.M{"last_modified_at": time.Now()}
	if Topup.Name != "" {
		update["name"] = Topup.Name
	}
	if Topup.Code != "" {
		update["code"] = Topup.Code
	}
	if Topup.Avatar != "" {
		update["avatar"] = Topup.Avatar
	}

	update["services.self"] = Topup.Services.Self
	update["services.other"] = Topup.Services.Other
	update["services.agent"] = Topup.Services.Agent
	return update
}
