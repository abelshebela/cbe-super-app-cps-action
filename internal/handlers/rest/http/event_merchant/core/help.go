package core

import (
	event_merchant_dto "github.com/abelshebela/cbe-super-app-cps-action/internal/constants/dto/event_merchant"

	"github.com/abelshebela/cbe-super-app-cps-action/internal/constants/model"
)

func CreateEventMerchantRequestToModel(req event_merchant_dto.CreateEventMerchantRequest) model.EventMerchant {
	return model.EventMerchant{
		MerchantID:        req.MerchantID,
		MerchantType:      req.MerchantType,
		SettlementMethod:  req.SettlementMethod,
		MerchantName:      req.MerchantName,
		BankAccountNumber: req.BankAccountNumber,
		Enabled:           true,
		Email: req.Email,
		PhoneNumber: req.PhoneNumber,
	}
}

func UpdateEventMerchantRequestToModel(req event_merchant_dto.UpdateEventMerchantRequest) model.EventMerchant {
	return model.EventMerchant{
		MerchantName:      req.MerchantName,
		BankAccountNumber: req.BankAccountNumber,
	}
}
