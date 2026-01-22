package ussd_merchant

import (
	ussd_merchant_dto "cbe-super-app-cps-action/internal/constants/dto/ussd_merchant"
	imodel "cbe-super-app-cps-action/internal/constants/model"
)

func ResponseMapper(data imodel.UssdMerchant) ussd_merchant_dto.UssdMerchantResponse {
	return ussd_merchant_dto.UssdMerchantResponse{
		ID:               data.ID,
		Name:             data.Name,
		PhoneNumber:      data.PhoneNumber,
		Email:            data.Email,
		Service:          data.Service,
		AccountNumber:    data.AccountNumber,
		Logo:             data.Logo,
		SettlementMethod: data.SettlementMethod,
		Enabled:          data.Enabled,
		CreatedAt:        data.CreatedAt,
	}
}
