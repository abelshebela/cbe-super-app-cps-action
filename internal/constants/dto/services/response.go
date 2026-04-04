package services

import (
	"time"

	"gitlab.com/bersufekadgetachew/cbe-super-app-shared/shared/entities/model"
)

type ServiceResponse struct {
	ID                       string      `bson:"_id,omitempty" json:"id"`
	ServiceKeyId             string      `json:"service_key_id" bson:"service_key_id"`
	ServiceName              string      `bson:"service_name" json:"service_name"`
	ServiceKey               string      `bson:"service_key" json:"service_key"`
	ServiceCode              string      `bson:"service_code" json:"service_code"`
	Cap                      []model.Cap `bson:"cap" json:"cap"`
	MinimumFraudAmount       string      `bson:"minimum_fraud_amount" json:"minimum_fraud_amount"`
	ProductGlAccount         string      `bson:"product_gl_account" json:"product_gl_account"`
	ProductGlAccountCurrency string      `bson:"product_gl_account_currency" json:"product_gl_account_currency"`
	Enabled                  bool        `bson:"enabled" json:"enabled"`
	IsDeleted                bool        `bson:"is_deleted" json:"is_deleted"`
	CreatedAt                time.Time   `bson:"created_at" json:"created_at"`
	LastModifiedAt           time.Time   `bson:"last_modified_at" json:"last_modified_at"`
	DeletedAt                *time.Time  `bson:"deleted_at,omitempty" json:"deleted_at,omitempty"`
}
