package core

import (
	TopupDto "cbe-super-app-cps-action/internal/constants/dto/topup"
	"cbe-super-app-cps-action/internal/constants/localization"
	"cbe-super-app-cps-action/pkgs/utils"
	"errors"
	"fmt"
	"log"
	"mime/multipart"
	"net/http"
	"strings"
)

func ParseTopupRequestFromMultipartForm(r *http.Request,isCreate bool) (TopupDto.TopupRequest, error) {
	var req TopupDto.TopupRequest
	_, fileHeader, err := utils.ParseMultipartFormFile(r, "avatar", 5<<20)
	if err != nil {
		if isCreate{
			if err.Error() != localization.ErrorMissingFile.Code {
				log.Println("errror here", err)
				return req, errors.New(localization.ErrorTopupImageMissingOrInvalid.Code)
			}
			return req, errors.New(localization.ErrorInvalidFileUpload.Code)
		}
		
	}
	req.Avatar = fileHeader
	req.Name = r.FormValue("name")
	req.Code = r.FormValue("code")
	req.Self = r.FormValue("self") == "true"
	req.Other = r.FormValue("other") == "true"
	req.Agent = r.FormValue("agent") == "true"

	return req, nil
}


func ParseMultipartFormFile(r *http.Request, key string, maxMemory int64) (multipart.File, *multipart.FileHeader, error) {
	if !strings.HasPrefix(r.Header.Get("Content-Type"), "multipart/form-data") {
		return nil, nil, fmt.Errorf("%s", localization.MsgFileNotFound)
	}
	if err := r.ParseMultipartForm(maxMemory); err != nil {
		return nil, nil, fmt.Errorf("failed to parse multipart form: %w", err)
	}

	file, fileHeader, err := r.FormFile(key)
	if err != nil {
		if err == http.ErrMissingFile {
			return nil, nil, fmt.Errorf("%s", localization.MsgFileNotFound)
		}
		return nil, nil, fmt.Errorf("missing or invalid file for key '%s': %w", key, err)
	}

	return file, fileHeader, nil
}
