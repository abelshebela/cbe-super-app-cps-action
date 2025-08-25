package core

import (
	"cbe-super-app-cps-action/internal/constants/dto/ad"
	"cbe-super-app-cps-action/internal/constants/localization"
	local_util "cbe-super-app-cps-action/pkgs/utils"
	"fmt"
	"mime/multipart"
	"time"

	local_errors "cbe-super-app-cps-action/internal/constants/errors"
	"net/http"

	cps_entities "cbe-super-app-cps-action/internal/constants/model"

	"gitlab.com/bersufekadgetachew/cbe-super-app-shared/shared/utils"
)

// extractUserAndMaker extracts user context and creates a maker, sending an error response if incomplete
func ExtractUserAndMaker(w http.ResponseWriter, r *http.Request, logger utils.Logger) (cps_entities.CPSUser, bool) {
	userContext := local_util.ExtractUserContext(r)
	if local_util.IsIncomplete(userContext) {
		logger.Errorf("[event.extractUserAndMaker] incomplete user context")
		localization.SendErrorResponse(w, localization.ErrorIncompleteUserInfo, nil, nil)
		return cps_entities.CPSUser{}, false
	}
	maker := cps_entities.CPSUser{
		UserCode:    userContext.UserID,
		FullName:    userContext.FullName,
		PhoneNumber: userContext.PhoneNumber,
		Department:  userContext.Department,
	}
	return maker, true
}

// extractID extracts and validates the ID parameter, sending an error response if invalid
func ExtractID(w http.ResponseWriter, r *http.Request, logger utils.Logger) (string, bool) {
	id, ok := local_util.GetParam(r, "id")
	if !ok {
		logger.Errorf("[event.extractID] missing or invalid parameter 'id'")
		localization.SendBadRequestResponse(w, local_errors.ErrInvalidInputParameters.Error())
		return "", false
	}
	return id, true
}

// parseTime parses a time string in RFC3339 format
func ParseTime(timeStr, fieldName string, isOptional bool, logger utils.Logger) (time.Time, error) {
	if timeStr == "" && isOptional {
		return time.Time{}, nil
	}
	if timeStr == "" {
		return time.Time{}, fmt.Errorf("%s_REQUIRED", fieldName)
	}
	t, err := time.Parse(time.RFC3339, timeStr)
	if err != nil {
		logger.Errorf("[event.parseTime] invalid %s format: %v", fieldName, err)
		return time.Time{}, fmt.Errorf("INVALID_%s_FORMAT", fieldName)
	}
	return t, nil
}

// parseBannerImage handles banner image parsing with size limit of 2MB
func ParseBannerImage(r *http.Request, isUpdate bool, logger utils.Logger) (*multipart.FileHeader, error) {
	_, fileHeader, err := local_util.ParseMultipartFormFile(r, "banner_image", 2<<20)
	if err != nil && (err.Error() != local_errors.ErrMissingFile.Error() && !isUpdate) {
		logger.Errorf("[ad.parseBannerImage] error parsing file: %v", err)
		return nil, err
	}
	return fileHeader, nil
}

// parseAndValidateAdvertRequest parses and validates the advert request from multipart form
func ParseAndValidateAdvertRequest(w http.ResponseWriter, r *http.Request, isUpdate bool, logger utils.Logger) (ad.AdvertRequest, bool) {
	var req ad.AdvertRequest

	// Parse banner image
	bannerImage, err := ParseBannerImage(r, isUpdate, logger)
	if err != nil {
		logger.Errorf("[event.parseAndValidateAdvertRequest] failed to parse banner image: %v", err)
		localization.SendErrorResponse(w, localization.ErrorMissingOrInvalidImage, nil, nil)
		return ad.AdvertRequest{}, false
	}
	req.BannerImage = bannerImage

	// Parse form values
	req.ID = r.FormValue("id")
	req.Title = r.FormValue("title")
	req.Description = r.FormValue("description")
	req.AdvertFor = r.FormValue("advert_for")

	// Parse dates
	startedAt, err := ParseTime(r.FormValue("started_at"), "START_DATE", isUpdate, logger)
	if err != nil {
		logger.Errorf("[event.parseAndValidateAdvertRequest] failed to parse started_at: %v", err)
		localization.SendErrorResponse(w, localization.ErrorInvalidDate, nil, nil)
		return ad.AdvertRequest{}, false
	}
	expiredAt, err := ParseTime(r.FormValue("expired_at"), "EXPIRE_DATE", isUpdate, logger)
	if err != nil {
		logger.Errorf("[event.parseAndValidateAdvertRequest] failed to parse expired_at: %v", err)
		localization.SendErrorResponse(w, localization.ErrorInvalidDate, nil, nil)
		return ad.AdvertRequest{}, false
	}
	req.Date = ad.AdvertDate{StartedAt: startedAt, ExpiredAt: expiredAt}

	if err := req.Validate(isUpdate); err != nil {
		logger.Errorf("[event.parseAndValidateAdvertRequest] validation failed: %v", err)
		localization.SendBadRequestResponse(w, err.Error())
		return ad.AdvertRequest{}, false
	}
	return req, true
}
