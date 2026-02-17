package core

import (
	"cbe-super-app-cps-action/internal/constants"
	ussd_merchant_dto "cbe-super-app-cps-action/internal/constants/dto/ussd_merchant"
	"cbe-super-app-cps-action/internal/constants/localization"
	local_util "cbe-super-app-cps-action/pkgs/utils"
	"context"
	"errors"
	"fmt"
	"gitlab.com/bersufekadgetachew/cbe-super-app-shared/shared/utils"
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

	input.PhoneNumber = local_util.FormatPhoneNumber(input.PhoneNumber)

	if input.Validate() != nil {
		return input.Validate()
	}

	return nil
}

func ParseMultipartFormFile(r *http.Request, key string, maxMemory int64, action string, logger utils.Logger) (multipart.File, *multipart.FileHeader, error) {

	file, fileHeader, err := local_util.ParseMultipartFormFile(r, "logo", int64(constants.MaxMemoryForUpload))
	if err != nil {
		if errors.Is(err, http.ErrMissingFile) {

			if action == constants.CREATE {
				if err == http.ErrMissingFile {
					logger.Errorf("Error logo file is missing error: %v", err)
					return nil, nil, fmt.Errorf("%s", localization.MsgFileNotFound)
				}
			}
			logger.Infof("Logo not provided for update - skipping file update")
			return nil, nil, nil
		}

		return nil, nil, fmt.Errorf("missing or invalid file for key '%s': %w", key, err)
	}

	return file, fileHeader, nil
}
