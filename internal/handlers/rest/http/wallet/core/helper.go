package core

import (
	"cbe-super-app-cps-action/internal/constants"
	walletDto "cbe-super-app-cps-action/internal/constants/dto/wallet"
	"cbe-super-app-cps-action/internal/constants/localization"
	"cbe-super-app-cps-action/pkgs/utils"
	"errors"
	"log"
	"net/http"
	"strconv"
)

func toBoolPtr(s string) (*bool, error) {
	if s == "" {
		return nil, nil // treat empty as null
	}
	b, err := strconv.ParseBool(s)
	if err != nil {
		return nil, err
	}
	return &b, nil
}

func ParseWalletRequestFromMultipartForm(r *http.Request, isCreate bool) (walletDto.WalletRequest, error) {
	var req walletDto.WalletRequest
	req.Name = r.FormValue("name")
	req.UniqueCode = r.FormValue("unique_code")
	req.Self, _ = toBoolPtr(r.FormValue("self"))
	req.Other, _ = toBoolPtr(r.FormValue("other"))
	req.Agent, _ = toBoolPtr(r.FormValue("agent"))
	req.SelfServiceID = r.FormValue("self_service_id")
	req.OtherServiceID = r.FormValue("other_service_id")
	req.AgentServiceID = r.FormValue("agent_service_id")
	// switch req.Type {
	// case constants.Bank, constants.Wallet, constants.MFI:
	// 	// valid type
	// default:
	// 	return req, localization.ErrorInvalidWalletCode
	// }
	_, fileHeader, err := utils.ParseMultipartFormFile(r, "avatar", int64(constants.MaxMemoryForUpload))
	if err != nil {
		if isCreate {
			if err.Error() != localization.ErrorMissingFile.Code {
				log.Println("errror here", err)
				return req, localization.ErrorWalletImageMissingOrInvalid
			}
			return req, errors.New(localization.ErrorInvalidFileUpload.Code)
		}

	}
	req.Avatar = fileHeader

	return req, nil
}
