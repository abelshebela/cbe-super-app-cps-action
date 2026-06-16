package term_and_condition_handler_core

import (
	"fmt"
	"net/http"
	"strings"

	"cbe-super-app-cps-action/internal/constants"
	tac_dto "cbe-super-app-cps-action/internal/constants/dto/term_and_condition"
	"cbe-super-app-cps-action/internal/constants/localization"
	imodel "cbe-super-app-cps-action/internal/constants/model"
	"mime/multipart"

	"gitlab.com/bersufekadgetachew/cbe-super-app-shared/shared/utils"
)

func ParsePDFFile(r *http.Request, logger utils.Logger) (*multipart.FileHeader, error) {
	if err := r.ParseMultipartForm(int64(constants.MaxMemoryForUpload)); err != nil {
		logger.Errorf("[TACHandler][ParsePDF] parse form err: %v", err)
		return nil, fmt.Errorf("failed to parse multipart form: %w", err)
	}
	_, fileHeader, err := r.FormFile("term_and_condition")
	if err != nil {
		if err == http.ErrMissingFile {
			return nil, fmt.Errorf("term_and_condition file is required")
		}
		logger.Errorf("[TACHandler][ParsePDF] form file err: %v", err)
		return nil, fmt.Errorf("failed to parse term_and_condition file: %w", err)
	}
	return fileHeader, nil
}

func ParseUploadRequest(r *http.Request, logger utils.Logger) (tac_dto.CreateTACRequest, error) {
	var req tac_dto.CreateTACRequest

	pdf, err := ParsePDFFile(r, logger)
	if err != nil {
		return req, err
	}
	req.TermAndCondition = pdf

	req.AccountProductID = strings.TrimSpace(r.FormValue("product"))
	req.ActivationTime = strings.TrimSpace(r.FormValue("activation_time"))
	req.VersionLabel = strings.TrimSpace(r.FormValue("version_label"))
	logger.Infof("[TACHandler][ParseUpload] parsed request: product=%s, activation_time=%s, version_label=%s", req.AccountProductID, req.ActivationTime, req.VersionLabel)
	return req, nil
}

func ParseUpdateRequest(r *http.Request, logger utils.Logger) (tac_dto.UpdateTACRequest, error) {
	var req tac_dto.UpdateTACRequest

	if err := r.ParseMultipartForm(int64(constants.MaxMemoryForUpload)); err != nil {
		logger.Errorf("[TACHandler][ParseUpdate] parse form err: %v", err)
		return req, fmt.Errorf("failed to parse multipart form: %w", err)
	}

	_, fileHeader, err := r.FormFile("term_and_condition")
	if err != nil && err != http.ErrMissingFile {
		logger.Errorf("[TACHandler][ParseUpdate] form file err: %v", err)
		return req, fmt.Errorf("failed to parse term_and_condition file: %w", err)
	}
	if fileHeader != nil {
		req.TermAndCondition = fileHeader
	}

	req.ActivationTime = strings.TrimSpace(r.FormValue("activation_time"))
	req.VersionLabel = strings.TrimSpace(r.FormValue("version_label"))

	return req, nil
}

func ValidateUpdate(req *tac_dto.UpdateTACRequest) localization.ResponseCode {
	if err := req.Validate(); err != nil {
		return localization.ErrorToResponseCode(err.Error(), 400, err.Error())
	}
	return localization.ResponseCode{}
}

func MapToResponse(m *imodel.AccountOpeningTerms) tac_dto.TACResponse {
	return tac_dto.TACResponse{
		ID: m.ID,
		AccountProduct: tac_dto.TACAccountProduct{
			ID:              m.AccountProductID,
			ProductName:     m.ProductName,
			CBSCategoryCode: m.CBSCategoryCode,
		},
		ActivationTime:         m.ActivationTime,
		VersionLabel:           m.VersionLabel,
		TermsAndConditionsPath: m.TermsAndConditionsPath,
		IsEnabled:              m.IsEnabled,
		IsDeleted:              m.IsDeleted,
		CreatedAt:              m.CreatedAt.String(),
		LastModifiedAt:         m.LastModifiedAt.String(),
	}
}

func ValidateUpload(req *tac_dto.CreateTACRequest) localization.ResponseCode {
	if err := req.Validate(); err != nil {
		return localization.ErrorToResponseCode(err.Error(), 400, err.Error())
	}
	return localization.ResponseCode{}
}
