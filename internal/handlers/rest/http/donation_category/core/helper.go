package core

import (
	"cbe-super-app-cps-action/internal/constants/localization"
	"cbe-super-app-cps-action/pkgs/utils"
	"errors"
	"net/http"
	dto "cbe-super-app-cps-action/internal/constants/dto/donation_category"

)

func ParseRequestFromMultipartForm(r *http.Request, isCreate bool) (dto.DonationCategoryRequest,error){
var req dto.DonationCategoryRequest
_, fileHeader, err := utils.ParseMultipartFormFile(r, "donation_icon", 2<<20)
if err != nil {
		if err.Error() != localization.ErrorMissingFile.Code || isCreate {
			return req, errors.New(localization.ErrorInvalidFileUpload.Code)
		}
	} else {
		// fileHeader = fileHeader[0]
		req.Icon = fileHeader
	}
	req.CategoryName = r.FormValue("category_name")
	
	

	return req,nil
	}



