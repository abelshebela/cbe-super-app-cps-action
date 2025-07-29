package wallet

import (
	"github.com/CBE-Super-App/cbe-super-app-cps-action/internal/domain/wallet"
)

func ToDomainWalletRequest(dto *WalletRequest) *wallet.WalletRequest {
	return &wallet.WalletRequest{
		Name:   dto.Name,
		Avatar: dto.Avatar,
		Code:   dto.Code,
	}
}
