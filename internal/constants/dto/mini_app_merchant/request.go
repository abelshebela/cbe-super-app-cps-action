package miniappmerchant

import "cbe-super-app-cps-action/internal/constants/types"

type MiniAppMerchantRequest struct {
	Type                       string
	MerchantName               string
	MerchantRepresentativeName string
	PhoneNumber                string
	Email                      string
	AccountNumber              string
}

type MiniAppMerchantDTO struct {
	Type                       string                    `json:"type"`
	MerchantName               string                    `json:"merchant_name"`
	MerchantRepresentativeName string                    `json:"merchant_representative_name"`
	PhoneNumber                string                    `json:"phone_number"`
	Email                      string                    `json:"email"`
	AccountNumber              string                    `json:"account_number"`
	Branches                   []types.BranchInformation `json:"branches"`
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
