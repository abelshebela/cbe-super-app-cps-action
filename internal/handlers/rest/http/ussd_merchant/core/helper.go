package core

import (
	"cbe-super-app-cps-action/internal/constants"
	ussd_merchant_dto "cbe-super-app-cps-action/internal/constants/dto/ussd_merchant"
	local_util "cbe-super-app-cps-action/pkgs/utils"
	"context"
	"mime/multipart"
	"net/http"
)

func ParseInputData(ctx context.Context, r *http.Request, req *ussd_merchant_dto.CreateUssdMerchantRequest) error {
	var file multipart.File
	var fileHeader *multipart.FileHeader
	var err error
	var input ussd_merchant_dto.CreateUssdMerchantRequest
	if err := req.Validate(); err != nil {
		return err
	}

	input = *req
	if req.Logo != nil {
		file, fileHeader, err = local_util.ParseMultipartFormFile(r, "logo", int64(constants.MaxMemoryForUpload))
		if err != nil {
			return err
		}
		if file != nil {
			defer file.Close()
			input.Logo = fileHeader
		}
	}

	if input.Validate() != nil {
		return input.Validate()
	}

	return nil
}
