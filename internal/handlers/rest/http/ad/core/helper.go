package core

import (
	"cbe-super-app-cps-action/internal/constants/dto/ad"
	"cbe-super-app-cps-action/internal/constants/localization"
	local_util "cbe-super-app-cps-action/pkgs/utils"
	"errors"
	"mime/multipart"
	"time"

	"net/http"

	cps_entities "cbe-super-app-cps-action/internal/constants/model"

	"gitlab.com/bersufekadgetachew/cbe-super-app-shared/shared/utils"
)

// extractUserAndMaker extracts user context and creates a maker, sending an error response if incomplete
func ExtractUserAndMaker(r *http.Request, logger utils.Logger) (*cps_entities.CPSUser, error) {
	userContext := local_util.ExtractUserContext(r)
	if local_util.IsIncomplete(userContext) {
		logger.Errorf("[event.extractUserAndMaker] incomplete user context")
		return nil, errors.New(localization.ErrorIncompleteUserInfo.Code)
	}
	maker := cps_entities.CPSUser{
		UserCode:    userContext.UserID,
		FullName:    userContext.FullName,
		PhoneNumber: userContext.PhoneNumber,
		// Department:  userContext.Department,
	}
	return &maker, nil
}

// extractID extracts and validates the ID parameter, sending an error response if invalid
func ExtractID(r *http.Request, logger utils.Logger) (string, error) {
	id, ok := local_util.GetParam(r, "id")
	if !ok {
		logger.Errorf("[event.extractID] missing or invalid parameter 'id'")
		return "", errors.New(localization.ErrorInvalidInputParameters.Code)
	}
	return id, nil
}

// parseTime parses a time string in RFC3339 format
func ParseTime(timeStr, fieldName string, isOptional bool, logger utils.Logger) (time.Time, error) {
	if timeStr == "" && isOptional {
		return time.Time{}, nil
	}
	if timeStr == "" {
		return time.Time{}, errors.New(localization.ErrorRequiredFieldMissing.Code)
	}
	t, err := time.Parse(time.RFC3339, timeStr)
	if err != nil {
		logger.Errorf("[event.parseTime] invalid %s format: %v", fieldName, err)
		return time.Time{}, errors.New(localization.ErrorInvalidDate.Code)
	}
	return t, nil
}

// parseBannerImage handles banner image parsing with size limit of 2MB
func ParseBannerImage(r *http.Request, isUpdate bool, logger utils.Logger) (*multipart.FileHeader, error) {
	_, fileHeader, err := local_util.ParseMultipartFormFile(r, "banner_image", 2<<20)
	if err != nil && (err.Error() != localization.ErrorMissingFile.Code && !isUpdate) {
		logger.Errorf("[ad.parseBannerImage] error parsing file: %v", err)
		return nil, err
	}
	return fileHeader, nil
}

// parseAndValidateAdvertRequest parses and validates the advert request from multipart form
func ParseAndValidateAdvertRequest(r *http.Request, isUpdate bool, logger utils.Logger) (*ad.AdvertRequest, error) {
	var req ad.AdvertRequest

	// Parse banner image
	bannerImage, err := ParseBannerImage(r, isUpdate, logger)
	if err != nil {
		logger.Errorf("[event.parseAndValidateAdvertRequest] failed to parse banner image: %v", err)
		return nil, localization.ErrorMissingOrInvalidImage
	}
	req.BannerImage = bannerImage

	// Parse form values
	req.ID = r.FormValue("id")
	req.Title = r.FormValue("title")
	req.Description = r.FormValue("description")
	req.AdvertFor = r.FormValue("advert_for")

	if err := req.Validate(isUpdate); err != nil {
		logger.Errorf("[event.parseAndValidateAdvertRequest] validation failed: %v", err)
		return nil, err
	}
	return &req, nil
}
