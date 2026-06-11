package account_product_category_handler_core

import (
	"encoding/json"
	"net/http"

	apc_dto "cbe-super-app-cps-action/internal/constants/dto/account_product_category"
	"cbe-super-app-cps-action/internal/constants/localization"
	imodel "cbe-super-app-cps-action/internal/constants/model"
)

func DecodeCreateRequest(r *http.Request) (apc_dto.CreateAPCRequest, error) {
	var req apc_dto.CreateAPCRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		return req, err
	}
	return req, nil
}

func DecodeUpdateRequest(r *http.Request) (apc_dto.UpdateAPCRequest, error) {
	var req apc_dto.UpdateAPCRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		return req, err
	}
	return req, nil
}

func MapToResponse(m *imodel.AccountProductCategory) apc_dto.APCResponse {
	return apc_dto.APCResponse{
		ID:              m.ID,
		ProductLine:     m.AccountType,
		CBSCategoryCode: m.CBSCategoryCode,
		CategoryName:    m.CategoryName,
		Description:     m.Description,
		IsEnabled:       m.IsEnabled,
		IsDeleted:       m.IsDeleted,
		CreatedAt:       m.CreatedAt.String(),
		LastModifiedAt:  m.LastModifiedAt.String(),
	}
}

func ValidateCreate(req *apc_dto.CreateAPCRequest) localization.ResponseCode {
	if err := req.Validate(); err != nil {
		return localization.ErrorToResponseCode(err.Error(), 400, err.Error())
	}
	return localization.ResponseCode{}
}

func ValidateUpdate(req *apc_dto.UpdateAPCRequest) localization.ResponseCode {
	if err := req.Validate(); err != nil {
		return localization.ErrorToResponseCode(err.Error(), 400, err.Error())
	}
	return localization.ResponseCode{}
}
