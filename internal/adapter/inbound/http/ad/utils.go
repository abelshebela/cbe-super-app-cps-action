package ad

import (
	"fmt"
	"mime/multipart"
	"net/http"
	"time"

	// "github.com/CBE-Super-App/cbe-super-app-cps-action/internal/domain/ad/dto"
	// "github.com/CBE-Super-App/cbe-super-app-cps-action/internal/domain/ad/mapper"

	cps_entities "github.com/CBE-Super-App/cbe-super-app-cps-action/internal/domain/cps_actions/entities"

	context "github.com/CBE-Super-App/cbe-super-app-cps-action/pkgs/context"
	"github.com/CBE-Super-App/cbe-super-app-cps-action/pkgs/utils"
)

// extractUserAndMaker extracts user context and creates a maker, sending an error response if incomplete
func (h *AdvertHTTPStore) extractUserAndMaker(w http.ResponseWriter, r *http.Request) (cps_entities.User, bool) {
	userContext := context.ExtractUserContext(r)
	if userContext.IsIncomplete() {
		h.logger.Errorf("[event.extractUserAndMaker] incomplete user context")
		utils.SendErrorResponse(w, utils.IncompleteUserInfo, 0, nil)
		return cps_entities.User{}, false
	}
	maker := cps_entities.User{
		UserCode:    userContext.UserID,
		FullName:    userContext.FullName,
		PhoneNumber: userContext.PhoneNumber,
		Department:  userContext.Department,
	}
	return maker, true
}

// extractID extracts and validates the ID parameter, sending an error response if invalid
func (h *AdvertHTTPStore) extractID(w http.ResponseWriter, r *http.Request) (string, bool) {
	id, ok := utils.GetParam(r, "id")
	if !ok {
		h.logger.Errorf("[event.extractID] missing or invalid parameter 'id'")
		utils.SendErrorResponse(w, utils.InvalidInputParameters, 0, nil)
		return "", false
	}
	return id, true
}

// parseTime parses a time string in RFC3339 format
func (h *AdvertHTTPStore) parseTime(timeStr, fieldName string, isOptional bool) (time.Time, error) {
	if timeStr == "" && isOptional {
		return time.Time{}, nil
	}
	if timeStr == "" {
		return time.Time{}, fmt.Errorf("%s_REQUIRED", fieldName)
	}
	t, err := time.Parse(time.RFC3339, timeStr)
	if err != nil {
		h.logger.Errorf("[event.parseTime] invalid %s format: %v", fieldName, err)
		return time.Time{}, fmt.Errorf("INVALID_%s_FORMAT", fieldName)
	}
	return t, nil
}

// parseBannerImage handles banner image parsing with size limit of 2MB
func (h *AdvertHTTPStore) parseBannerImage(r *http.Request, isUpdate bool) (*multipart.FileHeader, error) {
	_, fileHeader, err := utils.ParseMultipartFormFile(r, "banner_image", 2<<20)
	if err != nil && (err.Error() != utils.ErrMissingFile && !isUpdate) {
		h.logger.Errorf("[ad.parseBannerImage] error parsing file: %v", err)
		return nil, err
	}
	return fileHeader, nil
}

// parseAndValidateAdvertRequest parses and validates the advert request from multipart form
func (h *AdvertHTTPStore) parseAndValidateAdvertRequest(w http.ResponseWriter, r *http.Request, isUpdate bool) (AdvertRequest, bool) {
	var req AdvertRequest

	// Parse banner image
	bannerImage, err := h.parseBannerImage(r, isUpdate)
	if err != nil {
		h.logger.Errorf("[event.parseAndValidateAdvertRequest] failed to parse banner image: %v", err)
		utils.SendErrorResponse(w, utils.MissingOrInvalidImage, http.StatusBadRequest, nil)
		return AdvertRequest{}, false
	}
	req.BannerImage = bannerImage

	// Parse form values
	req.ID = r.FormValue("id")
	req.Title = r.FormValue("title")
	req.Description = r.FormValue("description")
	req.AdvertFor = r.FormValue("advert_for")

	// Parse dates
	startedAt, err := h.parseTime(r.FormValue("started_at"), "START_DATE", isUpdate)
	if err != nil {
		h.logger.Errorf("[event.parseAndValidateAdvertRequest] failed to parse started_at: %v", err)
		utils.SendErrorResponse(w, err.Error(), http.StatusBadRequest, nil)
		return AdvertRequest{}, false
	}
	expiredAt, err := h.parseTime(r.FormValue("expired_at"), "EXPIRE_DATE", isUpdate)
	if err != nil {
		h.logger.Errorf("[event.parseAndValidateAdvertRequest] failed to parse expired_at: %v", err)
		utils.SendErrorResponse(w, err.Error(), http.StatusBadRequest, nil)
		return AdvertRequest{}, false
	}
	req.Date = AdvertDate{StartedAt: startedAt, ExpiredAt: expiredAt}

	if err := req.Validate(isUpdate); err != nil {
		h.logger.Errorf("[event.parseAndValidateAdvertRequest] validation failed: %v", err)
		utils.SendErrorResponse(w, err, http.StatusBadRequest, nil)
		return AdvertRequest{}, false
	}
	return req, true
}

// sendResponse sends success or error responses
func (h *AdvertHTTPStore) sendResponse(w http.ResponseWriter, err error, successMessage string, data interface{}) {
	if err != nil {
		h.logger.Errorf("[event.sendResponse] error: %v", err)
		utils.SendErrorResponse(w, err.Error(), 0, nil)
		return
	}
	utils.WriteSuccessResponse(w, data, successMessage)
}
