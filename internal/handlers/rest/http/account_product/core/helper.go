package account_product_handler_core

import (
	"encoding/json"
	"fmt"
	"net/http"
	"strconv"
	"strings"

	"cbe-super-app-cps-action/internal/constants"
	ap_dto "cbe-super-app-cps-action/internal/constants/dto/account_product"
	"cbe-super-app-cps-action/internal/constants/localization"
	imodel "cbe-super-app-cps-action/internal/constants/model"
	local_util "cbe-super-app-cps-action/pkgs/utils"
	"mime/multipart"

	"gitlab.com/bersufekadgetachew/cbe-super-app-shared/shared/utils"
)

func ParseCoverImageFile(r *http.Request, isRequired bool, logger utils.Logger) (*multipart.FileHeader, error) {
	_, fileHeader, err := local_util.ParseMultipartFormFile(r, "cover_image", int64(constants.MaxMemoryForUpload))
	if err != nil {
		if err == http.ErrMissingFile {
			if isRequired {
				return nil, fmt.Errorf("cover image is required")
			}
			return nil, nil
		}
		logger.Errorf("[APHandler][ParseCoverImage] parse err: %v", err)
		return nil, fmt.Errorf("failed to parse cover image file: %w", err)
	}
	return fileHeader, nil
}
func ParseIconFile(r *http.Request, isRequired bool, logger utils.Logger) (*multipart.FileHeader, error) {
	_, fileHeader, err := local_util.ParseMultipartFormFile(r, "icon", int64(constants.MaxMemoryForUpload))
	if err != nil {
		if err == http.ErrMissingFile {
			if isRequired {
				return nil, fmt.Errorf("icon image is required")
			}
			return nil, nil
		}
		logger.Errorf("[APHandler][ParseIcon] parse err: %v", err)
		return nil, fmt.Errorf("failed to parse icon file: %w", err)
	}
	return fileHeader, nil
}

func ParseCreateRequest(r *http.Request, logger utils.Logger) (ap_dto.CreateAPRequest, error) {
	var req ap_dto.CreateAPRequest

	// icon, err := ParseIconFile(r, true, logger)
	// if err != nil {
	// 	return req, err
	// }
	// req.Icon = icon

	coverImage, err := ParseCoverImageFile(r, true, logger)
	if err != nil {
		return req, err
	}
	req.CoverImage = coverImage

	req.CBSProductCode = strings.TrimSpace(r.FormValue("cps_product_code"))
	req.ProductName = strings.TrimSpace(r.FormValue("product_name"))
	req.ProductTagLine = strings.TrimSpace(r.FormValue("product_tagline"))
	req.AccountCategoryID = strings.TrimSpace(r.FormValue("account_category"))
	req.AccountCurrency = strings.ToUpper(strings.TrimSpace(r.FormValue("account_currency")))
	req.FaqURL = strings.TrimSpace(r.FormValue("faq_url"))
	var features []string

	err = json.Unmarshal(
		[]byte(r.FormValue("product_features")),
		&features,
	)

	req.ProductFeatures = features
	// req.Eligibility = strings.TrimSpace(r.FormValue("eligibility"))
	req.ProductDescription = strings.TrimSpace(r.FormValue("product_description"))
	

	if v := r.FormValue("interest_rate"); v != "" {
		f, err := strconv.ParseFloat(v, 64)
		if err != nil {
			return req, fmt.Errorf("interest_rate must be a number")
		}
		req.InterestRate = f
	}

	if v := r.FormValue("is_available_for_onbording"); v != "" {
		b, err := strconv.ParseBool(v)
		if err != nil {
			return req, fmt.Errorf("is_avalilabe_for_onbording must be true or false")
		}
		req.IsAvailableForOnbording = b
	}

	if v := r.FormValue("minimum_opening_balance"); v != "" {
		f, err := strconv.ParseFloat(v, 64)
		if err != nil {
			return req, fmt.Errorf("minimum_opening_balance must be a number")
		}
		req.MinimumOpeningBalance = f
	}
	if maintenanceFeeStr := r.FormValue("minimum_maintenance_fee"); maintenanceFeeStr != "" {
		f, err := strconv.ParseFloat(maintenanceFeeStr, 64)
		if err != nil {
			return req, fmt.Errorf("minimum_maintenance_fee must be a number")
		}
		req.MinimumMaintenanceFee = f
	} else if v := r.FormValue("minimum_balance_to_maintain_account"); v != "" {
		f, err := strconv.ParseFloat(v, 64)
		if err != nil {
			return req, fmt.Errorf("minimum_balance_to_maintain_account must be a number")
		}
		req.MinimumMaintenanceFee = f
	}
	if v := r.FormValue("interest_fee"); v != "" {
		f, err := strconv.ParseFloat(v, 64)
		if err != nil {
			return req, fmt.Errorf("interest_fee must be a number")
		}
		req.InterestFee = f
	}


	if v := r.FormValue("has_atm_and_debit_card"); v != "" {
		b, err := strconv.ParseBool(v)
		if err != nil {
			return req, fmt.Errorf("has_atm_and_debit_card must be true or false")
		}
		req.HasPhysicalCard = b
	}
	if v := r.FormValue("has_virtual_debit_card"); v != "" {
		b, err := strconv.ParseBool(v)
		if err != nil {
			return req, fmt.Errorf("has_virtual_debit_card must be true or false")
		}
		req.HasVirtualCard = b
	}

	return req, nil
}

func ParseUpdateRequest(r *http.Request, logger utils.Logger) (ap_dto.UpdateAPRequest, error) {
	var req ap_dto.UpdateAPRequest

	// icon, err := ParseIconFile(r, false, logger)
	// if err != nil {
	// 	return req, err
	// }
	// req.Icon = icon

	coverImage, err := ParseCoverImageFile(r, false, logger)
	if err != nil {
		return req, err
	}
	req.CoverImage = coverImage

	req.CBSProductCode = strings.TrimSpace(r.FormValue("cps_product_code"))
	req.ProductName = strings.TrimSpace(r.FormValue("product_name"))
	req.ProductTagLine = strings.TrimSpace(r.FormValue("product_tagline"))
	req.AccountCategoryID = strings.TrimSpace(r.FormValue("account_category"))
	req.AccountCurrency = strings.ToUpper(strings.TrimSpace(r.FormValue("account_currency")))
	req.FaqURL = strings.TrimSpace(r.FormValue("faq_url"))
	var features []string

	err = json.Unmarshal(
		[]byte(r.FormValue("product_features")),
		&features,
	)

	req.ProductFeatures = features
	// req.Eligibility = strings.TrimSpace(r.FormValue("eligibility"))
	req.ProductDescription = strings.TrimSpace(r.FormValue("product_description"))
	
	if v := r.FormValue("interest_rate"); v != "" {
		f, err := strconv.ParseFloat(v, 64)
		if err != nil {
			return req, fmt.Errorf("interest_rate must be a number")
		}
		req.InterestRate = f
	}

	if v := r.FormValue("is_available_for_onbording"); v != "" {
		b, err := strconv.ParseBool(v)
		if err != nil {
			return req, fmt.Errorf("is_avalilabe_for_onbording must be true or false")
		}
		req.IsAvailableForOnbording = &b
	}

	if v := r.FormValue("minimum_opening_balance"); v != "" {
		f, err := strconv.ParseFloat(v, 64)
		if err != nil {
			return req, fmt.Errorf("minimum_opening_balance must be a number")
		}
		req.MinimumOpeningBalance = f
	}
	if maintenanceFeeStr := r.FormValue("minimum_maintenance_fee"); maintenanceFeeStr != "" {
		f, err := strconv.ParseFloat(maintenanceFeeStr, 64)
		if err != nil {
			return req, fmt.Errorf("minimum_maintenance_fee must be a number")
		}
		req.MinimumMaintenanceFee = f
	} else if v := r.FormValue("minimum_balance_to_maintain_account"); v != "" {
		f, err := strconv.ParseFloat(v, 64)
		if err != nil {
			return req, fmt.Errorf("minimum_balance_to_maintain_account must be a number")
		}
		req.MinimumMaintenanceFee = f
	}
	if v := r.FormValue("interest_fee"); v != "" {
		f, err := strconv.ParseFloat(v, 64)
		if err != nil {
			return req, fmt.Errorf("interest_fee must be a number")
		}
		req.InterestFee = f
	}

	if v := r.FormValue("has_atm_and_debit_card"); v != "" {
		b, err := strconv.ParseBool(v)
		if err != nil {
			return req, fmt.Errorf("has_atm_and_debit_card must be true or false")
		}
		req.HasPhysicalCard = &b
	}
	if v := r.FormValue("has_virtual_debit_card"); v != "" {
		b, err := strconv.ParseBool(v)
		if err != nil {
			return req, fmt.Errorf("has_virtual_debit_card must be true or false")
		}
		req.HasVirtualCard = &b
	}

	return req, nil
}

func MapToResponse(m *imodel.AccountProduct) ap_dto.APResponse {
	return ap_dto.APResponse{
		ID:             m.ID,
		CBSProductCode: m.CBSProductCode,
		ProductName:    m.ProductName,
		ProductTagLine: m.ProductTagLine,
		AccountCategory: ap_dto.APAccountCategory{
			ID:              m.AccountCategoryID,
			CategoryName:    m.CategoryName,
			CBSCategoryCode: m.CBSCategoryCode,
			AccountType:     m.AccountType,
		},
		AccountCurrency:       m.AccountCurrency,
		MinimumOpeningBalance: m.MinimumOpeningBalance,
		MinimumMaintenanceFee: m.MinimumMaintenanceFee,
		InterestFee:           m.InterestFee,
		FaqURL:                m.FaqURL,
		ProductFeatures:       m.ProductFeatures,
		HasPhysicalCard:       m.HasPhysicalCard,
		HasVirtualCard:        m.HasVirtualCard,
		// ProductIcon:           m.ProductIcon,
		IsAvailableForOnbording: m.IsAvailableForOnbording,
		InterestRate: m.InterestRate,
		ProductCoverImage: m.ProductCoverImage,
		IsEnabled:         m.IsEnabled,
		IsDeleted:         m.IsDeleted,
		CreatedAt:         m.CreatedAt.String(),
		LastModifiedAt:    m.LastModifiedAt.String(),
	}
}

func ValidateCreate(req *ap_dto.CreateAPRequest) localization.ResponseCode {
	if err := req.Validate(); err != nil {
		return localization.ErrorToResponseCode(err.Error(), 400, err.Error())
	}
	return localization.ResponseCode{}
}

func ValidateUpdate(req *ap_dto.UpdateAPRequest) localization.ResponseCode {
	if err := req.Validate(); err != nil {
		return localization.ErrorToResponseCode(err.Error(), 400, err.Error())
	}
	return localization.ResponseCode{}
}
