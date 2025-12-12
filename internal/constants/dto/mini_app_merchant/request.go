package miniappmerchant

import (
	shared_types "gitlab.com/bersufekadgetachew/cbe-super-app-shared/shared/entities/types"
)

type MiniAppMerchantDTO struct {
	Type                       string                    `json:"type" example:"3-click" enums:"3-click,merchant"`
	MerchantName               string                    `json:"merchant_name"`
	MerchantCode               string                    `json:"merchant_code"`
	MerchantRepresentativeName string                    `json:"merchant_representative_name"`
	PhoneNumber                string                    `json:"phone_number"`
	Email                      string                    `json:"email"`
	AccountNumber              string                    `json:"account_number"`
	Branches                   []shared_types.BranchInformation `json:"branches"`
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
