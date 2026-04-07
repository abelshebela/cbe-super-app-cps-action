package core

import (
	logistic_merchant_dto "cbe-super-app-cps-action/internal/constants/dto/logistics_merchant"

	local_model "cbe-super-app-cps-action/internal/constants/model"
)

func CreateLogisticsMerchantRequestToModel(req logistic_merchant_dto.CreateLogisticsMerchantRequest) local_model.LogisticsMerchant {
	return local_model.LogisticsMerchant{
		MerchantID:        req.MerchantID,
		MerchantType:      req.MerchantType,
		SettlementMethod:  req.SettlementMethod,
		MerchantName:      req.MerchantName,
		BankAccountNumber: req.BankAccountNumber,
	}
}

func UpdateLogisticsMerchantRequestToModel(req logistic_merchant_dto.UpdateLogisticsMerchantRequest) local_model.LogisticsMerchant {
	return local_model.LogisticsMerchant{
		MerchantName:      req.MerchantName,
		BankAccountNumber: req.BankAccountNumber,
	}
}
