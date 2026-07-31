package core

import (
	dto "github.com/abelshebela/cbe-super-app-cps-action/internal/constants/dto/donation_category"
	"github.com/abelshebela/cbe-super-app-cps-action/internal/constants/localization"
	"github.com/abelshebela/cbe-super-app-cps-action/pkgs/utils"
	"errors"
	"net/http"
)

func ParseRequestFromMultipartForm(r *http.Request, isCreate bool) (dto.DonationCategoryRequest, error) {
	var req dto.DonationCategoryRequest

	_, fileHeader, err := utils.ParseMultipartFormFile(r, "donation_icon", 15<<20)
	if err != nil {
		if errors.Is(err, http.ErrMissingFile) {
			if isCreate {
				return req, localization.ErrorInvalidFileUpload
			}

		} else {
			return req, nil
		}
	} else {
		req.Icon = fileHeader
	}

	req.CategoryName = r.FormValue("category_name")

	return req, nil
}
