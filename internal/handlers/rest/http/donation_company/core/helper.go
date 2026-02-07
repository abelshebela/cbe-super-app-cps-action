package core

import (
	dto "cbe-super-app-cps-action/internal/constants/dto/donation_company"
	"cbe-super-app-cps-action/internal/constants/localization"
	"cbe-super-app-cps-action/pkgs/utils"
	"errors"
	"net/http"
)

func ParseRequestFromMultipartForm(r *http.Request, isCreate bool) (dto.DonationCompanyRequest, error) {
	var req dto.DonationCompanyRequest

	_, fileHeader, err := utils.ParseMultipartFormFile(r, "company_logo", 15<<20)
	if err != nil {
		if errors.Is(err, http.ErrMissingFile) {
			if isCreate {
				return req, errors.New(localization.ErrorInvalidFileUpload.Code)
			}
		} else {
			return req, nil
		}
	} else {
		req.CompanyLogo = fileHeader
	}

	req.CompanyName = r.FormValue("company_name")
	req.AccountNumber = r.FormValue("account_number")
	req.PhoneNumber = r.FormValue("phone_number")
	req.Email = r.FormValue("email")
	req.Address = r.FormValue("address")
	req.CompanyDescription =r.FormValue("company_description")

	return req, nil
}
