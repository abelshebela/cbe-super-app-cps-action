package core

import (
	event_merchant_dto "cbe-super-app-cps-action/internal/constants/dto/event_merchant"

	"gitlab.com/bersufekadgetachew/cbe-super-app-shared/shared/entities/model"
)

func CreateEventMerchantRequestToModel(req event_merchant_dto.CreateEventMerchantRequest) model.EventMerchant {
	return model.EventMerchant{
		MerchantID:        req.MerchantID,
		MerchantType:      req.MerchantType,
		SettlementMethod:  req.SettlementMethod,
		MerchantName:      req.MerchantName,
		BankAccountNumber: req.BankAccountNumber,
		Email:             req.Email,
		PhoneNumber:       req.PhoneNumber,
	}
}

func UpdateEventMerchantRequestToModel(req event_merchant_dto.UpdateEventMerchantRequest) model.EventMerchant {
	return model.EventMerchant{
		MerchantID:        req.MerchantID,
		MerchantType:      req.MerchantType,
		SettlementMethod:  req.SettlementMethod,
		MerchantName:      req.MerchantName,
		BankAccountNumber: req.BankAccountNumber,
		Email:             req.Email,
		PhoneNumber:       req.PhoneNumber,
	}
}
