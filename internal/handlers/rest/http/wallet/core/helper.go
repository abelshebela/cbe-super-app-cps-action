package core

import (
	walletDto "cbe-super-app-cps-action/internal/constants/dto/wallet"
	"cbe-super-app-cps-action/internal/constants/localization"
	"cbe-super-app-cps-action/pkgs/utils"
	"errors"
	"log"
	"net/http"
)

func ParseWalletRequestFromMultipartForm(r *http.Request, isCreate bool) (walletDto.WalletRequest, error) {
	var req walletDto.WalletRequest
	_, fileHeader, err := utils.ParseMultipartFormFile(r, "avatar", 5<<20)
	if err != nil {
		if err.Error() != localization.ErrorMissingFile.Code || isCreate {
			log.Println("errror here",err)
			return req, errors.New(localization.ErrorInvalidFileUpload.Code)
		}
	} else {
		req.Avatar = fileHeader
	}

	req.Name = r.FormValue("name")
	req.Code = r.FormValue("code")

	return req, nil
}
