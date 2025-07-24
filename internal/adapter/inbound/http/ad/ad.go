package ad

import (
	"fmt"
	"mime/multipart"
	"net/http"
	"time"

	"github.com/CBE-Super-App/cbe-super-app-cps-action/internal/adapter/outbound/model"
	"github.com/CBE-Super-App/cbe-super-app-cps-action/internal/application/ad"
	"github.com/CBE-Super-App/cbe-super-app-cps-action/internal/domain/ad/dto"
	cps_const "github.com/CBE-Super-App/cbe-super-app-cps-action/internal/domain/cps_actions/constant"
	ctx_util "github.com/CBE-Super-App/cbe-super-app-cps-action/pkgs/context"
	common_util "github.com/CBE-Super-App/cbe-super-app-cps-action/pkgs/utils"

	"gitlab.com/bersufekadgetachew/cbe-super-app-shared/shared/utils"
)

// ADAdapter handles HTTP requests for advertisement operations
type ADAdapter struct {
	adHandler ad.ADHandlers
	logger    utils.Logger
}

// InitADAdapter initializes a new ADAdapter
func InitADAdapter(adHandler ad.ADHandlers, logger utils.Logger) ADAdapter {
	return ADAdapter{
		adHandler: adHandler,
		logger:    logger,
	}
}

// parseTime parses a time string in RFC3339 format
func (a ADAdapter) parseTime(timeStr, fieldName string) (time.Time, error) {
	if timeStr == "" {
		return time.Time{}, fmt.Errorf("%s_REQUIRED", fieldName)
	}
	t, err := time.Parse(time.RFC3339, timeStr)
	if err != nil {
		a.logger.Errorf("invalid %s format: %v", fieldName, err)
		return time.Time{}, fmt.Errorf("invalid %s format", fieldName)
	}
	return t, nil
}

// validateID extracts and validates the ID parameter
func (a ADAdapter) validateID(r *http.Request) (string, bool) {
	id, ok := common_util.GetParam(r, "id")
	if !ok {
		a.logger.Errorf("missing or invalid parameter 'id'")
	}
	return id, ok
}

// getUserContext extracts and validates user context
func (a ADAdapter) getUserContext(r *http.Request) (ctx_util.UserContext, bool) {
	userContext := ctx_util.ExtractUserContext(r)
	if userContext.IsIncomplete() {
		return ctx_util.UserContext{}, false
	}
	return userContext, true
}

// buildCPSRequest constructs a CPS request from user context
func (a ADAdapter) buildCPSRequest(userContext ctx_util.UserContext, actionData interface{}) model.CreateCPSAction {
	return model.CreateCPSAction{
		ActionData: actionData,
		MakerUser: model.User{
			UserCode:    userContext.UserCode,
			FullName:    userContext.FullName,
			PhoneNumber: userContext.PhoneNumber,
		},
		Department: userContext.Department,
	}
}

// parseBannerImage handles banner image parsing with size limit of 10MB
func (a ADAdapter) parseBannerImage(r *http.Request, isUpdate bool) (*multipart.FileHeader, error) {
	file, fileHeader, err := common_util.ParseMultipartFormFile(r, "banner_image", 10<<20)
	if err != nil {
		if isUpdate && err.Error() == common_util.ErrMissingFile {
			return nil, nil
		}
		a.logger.Errorf("error parsing file: %v", err)
		return nil, err
	}
	if file != nil {
		defer file.Close()
	}
	return fileHeader, nil
}

// CreateOneAdvert handles creation of a new advertisement
func (a ADAdapter) CreateOneAdvert(w http.ResponseWriter, r *http.Request) {
	var advertReq dto.CreateAdvertRequest

	// Parse banner image
	bannerImage, err := a.parseBannerImage(r, false)
	if err != nil {
		common_util.SendErrorResponse(w, common_util.MissingOrInvalidImage, 0, nil)
		return
	}
	advertReq.BannerImage = bannerImage

	// Parse form values
	advertReq.Title = r.FormValue("title")
	advertReq.Description = r.FormValue("description")
	advertReq.AdvertFor = dto.AdvertFor(r.FormValue("advert_for"))

	// Parse and validate dates
	startedAt, err := a.parseTime(r.FormValue("started_at"), "START_DATE")
	if err != nil {
		common_util.SendErrorResponse(w, err.Error(), http.StatusBadRequest, nil)
		return
	}
	expiredAt, err := a.parseTime(r.FormValue("expired_at"), "EXPIRE_DATE")
	if err != nil {
		common_util.SendErrorResponse(w, err.Error(), http.StatusBadRequest, nil)
		return
	}
	advertReq.Date.StartedAt = startedAt
	advertReq.Date.ExpiredAt = expiredAt

	// Validate required fields
	if advertReq.Title == "" || advertReq.Description == "" || advertReq.AdvertFor == "" {
		common_util.SendErrorResponse(w, "missing required fields: title, description, advert_for", http.StatusBadRequest, nil)
		return
	}

	// Validate user context
	userContext, ok := a.getUserContext(r)
	if !ok {
		common_util.SendErrorResponse(w, common_util.IncompleteUserInfo, 0, nil)
		return
	}

	// Execute request
	cpsReq := a.buildCPSRequest(userContext, advertReq)
	_, err = a.adHandler.CreateOneAdvert(r.Context(), cpsReq)
	if err != nil {
		common_util.SendErrorResponse(w, err.Error(), 0, nil)
		return
	}

	common_util.WriteSuccessResponse(w, nil, "AD Create Request Created successfully")
}

// DeleteOneAdvert handles deletion of an advertisement
func (a ADAdapter) DeleteOneAdvert(w http.ResponseWriter, r *http.Request) {
	id, ok := a.validateID(r)
	if !ok {
		common_util.SendErrorResponse(w, common_util.InvalidInputParameters, 0, nil)
		return
	}

	userContext, ok := a.getUserContext(r)
	if !ok {
		common_util.SendErrorResponse(w, common_util.IncompleteUserInfo, 0, nil)
		return
	}

	cpsReq := a.buildCPSRequest(userContext, nil)
	_, err := a.adHandler.DeleteOneAdvert(r.Context(), id, cpsReq)
	if err != nil {
		common_util.SendErrorResponse(w, err.Error(), 0, nil)
		return
	}

	common_util.WriteSuccessResponse(w, nil, "AD Delete Request Created successfully")
}

// GetAllAdvert retrieves all advertisements
func (a ADAdapter) GetAllAdvert(w http.ResponseWriter, r *http.Request) {
	filterParams := common_util.ExtractFilterParams(r)
	adverts, err := a.adHandler.GetAllAdvert(r.Context(), filterParams)
	if err != nil {
		common_util.SendErrorResponse(w, err.Error(), 0, nil)
		return
	}

	common_util.WriteSuccessResponse(w, adverts, "AD fetch successfully")
}

// GetOneAdvert retrieves a single advertisement
func (a ADAdapter) GetOneAdvert(w http.ResponseWriter, r *http.Request) {
	id, ok := a.validateID(r)
	if !ok {
		common_util.SendErrorResponse(w, common_util.InvalidInputParameters, 0, nil)
		return
	}

	advert, err := a.adHandler.GetOneAdvert(r.Context(), id)
	if err != nil {
		common_util.SendErrorResponse(w, err.Error(), 0, nil)
		return
	}

	common_util.WriteSuccessResponse(w, advert, "AD fetch successfully")
}

// UpdateOneAdvert handles updating an advertisement
func (a ADAdapter) UpdateOneAdvert(w http.ResponseWriter, r *http.Request) {
	id, ok := a.validateID(r)
	if !ok {
		common_util.SendErrorResponse(w, common_util.InvalidInputParameters, 0, nil)
		return
	}

	var advertReq dto.UpdateAdvertRequest
	advertReq.ID = id

	// Parse banner image
	bannerImage, err := a.parseBannerImage(r, true)
	if err != nil {
		common_util.SendErrorResponse(w, common_util.MissingOrInvalidImage, 0, nil)
		return
	}
	advertReq.BannerImage = bannerImage

	// Parse form values
	advertReq.Title = r.FormValue("title")
	advertReq.Description = r.FormValue("description")
	advertReq.AdvertFor = dto.AdvertFor(r.FormValue("advert_for"))

	// Parse optional dates
	startedAtStr := r.FormValue("started_at")
	expiredAtStr := r.FormValue("expired_at")
	if startedAtStr != "" {
		if startedAt, err := a.parseTime(startedAtStr, "started_at"); err != nil {
			common_util.SendErrorResponse(w, err.Error(), http.StatusBadRequest, nil)
			return
		} else {
			advertReq.Date.StartedAt = startedAt
		}
	}
	if expiredAtStr != "" {
		if expiredAt, err := a.parseTime(expiredAtStr, "expired_at"); err != nil {
			common_util.SendErrorResponse(w, err.Error(), http.StatusBadRequest, nil)
			return
		} else {
			advertReq.Date.ExpiredAt = expiredAt
		}
	}

	// Validate update data
	if advertReq.Title == "" && advertReq.Description == "" &&
		advertReq.AdvertFor == "" && advertReq.BannerImage == nil &&
		startedAtStr == "" && expiredAtStr == "" {
		common_util.SendErrorResponse(w, "NO_DATA_PROVIDED_FOR_UPDATE", http.StatusBadRequest, nil)
		return
	}

	userContext, ok := a.getUserContext(r)
	if !ok {
		common_util.SendErrorResponse(w, common_util.IncompleteUserInfo, 0, nil)
		return
	}

	cpsReq := a.buildCPSRequest(userContext, advertReq)
	_, err = a.adHandler.UpdateOneAdvert(r.Context(), id, cpsReq)
	if err != nil {
		common_util.SendErrorResponse(w, err.Error(), 0, nil)
		return
	}

	common_util.WriteSuccessResponse(w, nil, "AD Update Request Created successfully")
}

// EnableOrDisableAdvert handles enabling or disabling an advertisement
func (a ADAdapter) enableOrDisableAdvert(w http.ResponseWriter, r *http.Request, actionType cps_const.RequestAction, successMessage string) {
	id, ok := a.validateID(r)
	if !ok {
		common_util.SendErrorResponse(w, common_util.InvalidInputParameters, 0, nil)
		return
	}

	userContext, ok := a.getUserContext(r)
	if !ok {
		common_util.SendErrorResponse(w, common_util.IncompleteUserInfo, 0, nil)
		return
	}

	cpsReq := a.buildCPSRequest(userContext, nil)
	_, err := a.adHandler.EnableOrDisableAdvert(r.Context(), id, actionType, cpsReq)
	if err != nil {
		common_util.SendErrorResponse(w, err.Error(), 0, nil)
		return
	}

	common_util.WriteSuccessResponse(w, nil, successMessage)
}

// EnableAdvert handles enabling an advertisement
func (a ADAdapter) EnableAdvert(w http.ResponseWriter, r *http.Request) {
	a.enableOrDisableAdvert(w, r, cps_const.RequestEnableAdvert, "Advert Enable Request Created successfully")
}

// DisableAdvert handles disabling an advertisement
func (a ADAdapter) DisableAdvert(w http.ResponseWriter, r *http.Request) {
	a.enableOrDisableAdvert(w, r, cps_const.RequestDisableAdvert, "Advert Disable Request Created successfully")
}
