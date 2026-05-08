package core

import (
	logistic_merchant_dto "cbe-super-app-cps-action/internal/constants/dto/logistics_merchant"
	"strings"

	local_model "cbe-super-app-cps-action/internal/constants/model"
)

func CreateLogisticsMerchantRequestToModel(req logistic_merchant_dto.CreateLogisticsMerchantRequest) local_model.LogisticsMerchantOracle {
	return local_model.LogisticsMerchantOracle{
		MerchantCode:          req.MerchantID,
		MerchantType:          "LOGISTICS",
		SettlementMethod:      strings.ToUpper(req.SettlementMethod),
		MerchantName:          req.MerchantName,
		MerchantAccountNumber: req.BankAccountNumber,
	}
}

func UpdateLogisticsMerchantRequestToModel(req logistic_merchant_dto.UpdateLogisticsMerchantRequest) local_model.LogisticsMerchantOracle {
	return local_model.LogisticsMerchantOracle{
		MerchantName:          req.MerchantName,
		MerchantAccountNumber: req.BankAccountNumber,
	}
}
