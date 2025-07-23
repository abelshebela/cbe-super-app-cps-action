package ad

import (

	// "fmt"

	"net/http"
	"time"

	"github.com/CBE-Super-App/cbe-super-app-cps-action/internal/adapter/outbound/model"
	"github.com/CBE-Super-App/cbe-super-app-cps-action/internal/application/ad"
	"github.com/CBE-Super-App/cbe-super-app-cps-action/internal/domain/ad/dto"

	ctx_util "github.com/CBE-Super-App/cbe-super-app-cps-action/pkgs/context"
	common_util "github.com/CBE-Super-App/cbe-super-app-cps-action/pkgs/utils"

	inboundAd "github.com/CBE-Super-App/cbe-super-app-cps-action/internal/port/inbound/ad"
	util "github.com/CBE-Super-App/cbe-super-app-cps-action/pkgs/utils"

	"gitlab.com/bersufekadgetachew/cbe-super-app-shared/shared/utils"
)

type ADAdapter struct {
	adHandler ad.ADHandlers
	logger    utils.Logger
}

func InitADAdapter(adHandler ad.ADHandlers, logger utils.Logger) inboundAd.ADAdapter {
	return ADAdapter{
		adHandler: adHandler,
		logger:    logger,
	}
}

func (a ADAdapter) CreateOneAdvert(w http.ResponseWriter, r *http.Request) {
	var advertReq dto.CreateAdvertRequest

	file, fileHeader, err := util.ParseMultipartFormFile(r, "banner_image", 10<<20)
	if err != nil {
		a.logger.Errorf("error parsing file: %v", err)
		util.SendErrorResponse(w, util.MissingOrInvalidImage, 0, nil)
		return
	}
	defer file.Close()

	startedAt := r.FormValue("started_at")
	startedAtTime, err := time.Parse(time.RFC3339, startedAt)
	if err != nil {
		a.logger.Errorf("failed to parse form data: %v", err)
		util.SendErrorResponse(w, err.Error(), http.StatusBadRequest, nil)
		return
	}

	expiredAt := r.FormValue("expired_at")
	expiredAtTime, err := time.Parse(time.RFC3339, expiredAt)
	if err != nil {
		a.logger.Errorf("failed to parse time", err)
		util.SendErrorResponse(w, err.Error(), http.StatusInternalServerError, nil)
		return
	}
	advertReq.Title = r.FormValue("title")
	advertReq.Description = r.FormValue("description")
	advertReq.AdvertFor = dto.AdvertFor(r.FormValue("advert_for"))
	advertReq.Date.StartedAt = startedAtTime
	advertReq.Date.ExpiredAt = expiredAtTime

	advertReq.BannerImage = fileHeader

	var cpsReq model.CreateCPSAction

	userContext := ctx_util.ExtractUserContext(r)
	if userContext.IsIncomplete() {
		util.SendErrorResponse(w, util.IncompleteUserInfo, 0, nil)
		return
	}

	cpsReq.ActionData = advertReq

	cpsReq.MakerUser = model.User{
		UserCode:    userContext.UserCode,
		FullName:    userContext.FullName,
		PhoneNumber: userContext.PhoneNumber,
	}
	cpsReq.Department = userContext.Department

	ctx := r.Context()
	_, err = a.adHandler.CreateOneAdvert(ctx, cpsReq)
	if err != nil {
		util.SendErrorResponse(w, err.Error(), 0, nil)
		return
	}

	util.WriteSuccessResponse(w, nil, "AD Create Request Created successfully")
}

func (a ADAdapter) DeleteOneAdvert(w http.ResponseWriter, r *http.Request) {
	id, ok := common_util.GetParam(r, "id")
	if !ok {
		a.logger.Errorf("missing or invalid parameter 'id'")
		common_util.SendErrorResponse(w, common_util.InvalidInputParameters, 0, nil)
		return
	}

	var cpsReq model.CreateCPSAction

	userContext := ctx_util.ExtractUserContext(r)
	if userContext.IsIncomplete() {
		util.SendErrorResponse(w, util.IncompleteUserInfo, 0, nil)
		return
	}
	cpsReq.MakerUser = model.User{
		UserCode:    userContext.UserCode,
		FullName:    userContext.FullName,
		PhoneNumber: userContext.PhoneNumber,
	}
	cpsReq.Department = userContext.Department
	ctx := r.Context()

	_, err := a.adHandler.DeleteOneAdvert(ctx, id, cpsReq)
	if err != nil {
		util.SendErrorResponse(w, err.Error(), 0, nil)
		return
	}
	util.WriteSuccessResponse(w, nil, "AD Delete Request Created successfully")

}
func (a ADAdapter) GetAllAdvert(w http.ResponseWriter, r *http.Request) {
	filterParams := common_util.ExtractFilterParams(r)
	adverts, err := a.adHandler.GetAllAdvert(r.Context(), filterParams)
	if err != nil {
		util.SendErrorResponse(w, err.Error(), 0, nil)
		return
	}

	util.WriteSuccessResponse(w, adverts, "AD featch successfully")

}

func (a ADAdapter) GetOneAdvert(w http.ResponseWriter, r *http.Request) {
	id, ok := common_util.GetParam(r, "id")
	if !ok {
		a.logger.Errorf("missing or invalid parameter 'id'")
		common_util.SendErrorResponse(w, common_util.InvalidInputParameters, 0, nil)
		return
	}

	ctx := r.Context()

	advert, err := a.adHandler.GetOneAdvert(ctx, id)
	if err != nil {
		util.SendErrorResponse(w, err.Error(), 0, nil)
		return
	}

	util.WriteSuccessResponse(w, advert, "AD featch successfully")

}

func (a ADAdapter) UpdateOneAdvert(w http.ResponseWriter, r *http.Request) {
	id, ok := common_util.GetParam(r, "id")
	if !ok {
		a.logger.Errorf("missing or invalid parameter 'id'")
		common_util.SendErrorResponse(w, common_util.InvalidInputParameters, 0, nil)
		return
	}

	var advertReq dto.UpdateAdvertRequest

	file, fileHeader, err := util.ParseMultipartFormFile(r, "banner_image", 10<<20)
	if err != nil && err.Error() != common_util.ErrMissingFile {
		a.logger.Errorf("error parsing file: %v", err)
		util.SendErrorResponse(w, util.MissingOrInvalidImage, 0, nil)
		return
	}
	if file != nil {
		defer file.Close()
	}

	advertReq.Title = r.FormValue("title")
	advertReq.Description = r.FormValue("description")
	advertReq.AdvertFor = dto.AdvertFor(r.FormValue("advert_for"))
	startedAt := r.FormValue("started_at")
	expiredAt := r.FormValue("expired_at")

	var startedAtTime, expiredAtTime time.Time
	if startedAt != "" {
		startedAtTime, err = time.Parse(time.RFC3339, startedAt)
		if err != nil {
			a.logger.Errorf("failed to parse started_at: %v", err)
			util.SendErrorResponse(w, "invalid started_at format", http.StatusBadRequest, nil)
			return
		}
	}
	if expiredAt != "" {
		expiredAtTime, err = time.Parse(time.RFC3339, expiredAt)
		if err != nil {
			a.logger.Errorf("failed to parse expired_at: %v", err)
			util.SendErrorResponse(w, "invalid expired_at format", http.StatusBadRequest, nil)
			return
		}
	}

	advertReq.Date.StartedAt = startedAtTime
	advertReq.Date.ExpiredAt = expiredAtTime
	advertReq.BannerImage = fileHeader

	if advertReq.Title == "" && advertReq.Description == "" &&
		advertReq.AdvertFor == "" && fileHeader == nil &&
		startedAt == "" && expiredAt == "" {
		util.SendErrorResponse(w, "NO_DATA_PROVIDED_FOR_UPDATE", http.StatusBadRequest, nil)
		return
	}

	userContext := ctx_util.ExtractUserContext(r)
	if userContext.IsIncomplete() {
		util.SendErrorResponse(w, util.IncompleteUserInfo, 0, nil)
		return
	}

	// Build CPS request
	cpsReq := model.CreateCPSAction{
		MakerUser: model.User{
			UserCode:    userContext.UserCode,
			FullName:    userContext.FullName,
			PhoneNumber: userContext.PhoneNumber,
		},
		Department: userContext.Department,
		ActionData: advertReq,
	}
	advertReq.ID = id

	_, err = a.adHandler.UpdateOneAdvert(r.Context(), id, cpsReq)
	if err != nil {
		util.SendErrorResponse(w, err.Error(), 0, nil)
		return
	}

	util.WriteSuccessResponse(w, nil, "AD Update Request Created successfully")
}
