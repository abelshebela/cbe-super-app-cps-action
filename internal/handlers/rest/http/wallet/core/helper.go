package core

import (
	"cbe-super-app-cps-action/internal/constants"
	walletDto "cbe-super-app-cps-action/internal/constants/dto/wallet"
	"cbe-super-app-cps-action/internal/constants/localization"
	"cbe-super-app-cps-action/pkgs/utils"
	"errors"
	"log"
	"net/http"
)

func ParseWalletRequestFromMultipartForm(r *http.Request, isCreate bool) (walletDto.WalletRequest, error) {
	var req walletDto.WalletRequest
	req.Name = r.FormValue("name")
	req.Code = r.FormValue("code")
	req.Self = r.FormValue("self") == "true"
	req.Other = r.FormValue("other") == "true"
	req.Agent = r.FormValue("agent") == "true"
	req.Type = constants.FinancialInstitutionType(r.FormValue("type"))
	switch req.Type {
	case constants.Bank, constants.Wallet, constants.MFI:
		// valid type
	default:
		return req, localization.ErrorInvalidWalletCode
	}
	_, fileHeader, err := utils.ParseMultipartFormFile(r, "avatar", int64(constants.MaxMemoryForUpload))
	if err != nil {
		if isCreate {
			if err.Error() != localization.ErrorMissingFile.Code {
				log.Println("errror here", err)
				return req, errors.New(localization.ErrorWalletImageMissingOrInvalid.Code)
			}
			return req, errors.New(localization.ErrorInvalidFileUpload.Code)
		}

	}
	req.Avatar = fileHeader

	return req, nil
}
