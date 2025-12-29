package miniappmerchant

import (
	"gitlab.com/bersufekadgetachew/cbe-super-app-shared/shared/entities/model"
)

type EcommerceMerchant struct {
	MerchantName     string                    `json:"merchant_name"`
	MerchantCode     string                    `json:"merchant_code"`
	AccountNumber    string                    `json:"account_number"`
	SettlementMethod string                    `json:"settlement_method" bson:"settlement_method"`
	Branches         []model.BranchInformation `json:"branches"`
}

type KYCDTO struct {
	Status         string            `json:"status"`
	Representative RepresentativeDTO `json:"representative"`
}

type RepresentativeDTO struct {
	Name  string `json:"name"`
	Phone string `json:"phone"`
	Email string `json:"email"`
}
