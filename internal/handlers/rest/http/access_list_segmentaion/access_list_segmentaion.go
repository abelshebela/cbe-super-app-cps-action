package accesslistsegmentaion

import (
	access_list_segmentation_dto "cbe-super-app-cps-action/internal/constants/dto/access_list_segmentation"
	accesslistsegmentation "cbe-super-app-cps-action/internal/constants/interfaces/access_list_segmentation"
	"cbe-super-app-cps-action/internal/constants/localization"
	"cbe-super-app-cps-action/internal/service"
	local_util "cbe-super-app-cps-action/pkgs/utils"
	"encoding/json"
	"net/http"

	"github.com/go-chi/chi/v5"
	"gitlab.com/bersufekadgetachew/cbe-super-app-shared/shared/utils"
)

type accessListSegmentation struct {
	service service.AccessListSegmentationService
	logger  utils.Logger
}

// CreateAccessListSegmentation implements accesslistsegmentation.AccessListSegmentationHandler.
func (a *accessListSegmentation) CreateAccessListSegmentation(w http.ResponseWriter, r *http.Request) {
	var req access_list_segmentation_dto.CreateAccessListSegmentationRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		a.logger.Errorf("[CreateAccessListSegmentation] failed to decode request body: %v", err)
		localization.SendErrorByCodeResponse(w, err.Error())
		return
	}
	if err := req.Validate(); err != nil {
		a.logger.Errorf("[CreateAccessListSegmentation] validation error: %v", err)
		localization.SendErrorByCodeResponse(w, err.Error())
		return
	}
	if err := a.service.CreateAccessListSegmentation(r.Context(), req); err != nil {
		a.logger.Errorf("[CreateAccessListSegmentation] service error: %v", err)
		localization.SendErrorByCodeResponse(w, err.Error())
		return
	}
	localization.SendSuccessResponse(w, localization.SuccessAccessListSegmentationCreated, nil)
}

// DisableAccessListSegmentation implements accesslistsegmentation.AccessListSegmentationHandler.
func (a *accessListSegmentation) DisableAccessListSegmentation(w http.ResponseWriter, r *http.Request) {
	id := chi.URLParam(r, "id")
	if id == "" {
		a.logger.Errorf("[DisableAccessListSegmentation] missing id parameter")
		localization.SendErrorByCodeResponse(w, localization.ErrorAccessListSegmentationInvalidID.Code)
		return
	}
	if err := a.service.EnableDisableAccessListSegmentation(r.Context(), id, false); err != nil {
		a.logger.Errorf("[DisableAccessListSegmentation] service error: %v", err)
		localization.SendErrorByCodeResponse(w, err.Error())
		return
	}
	localization.SendSuccessResponse(w, localization.SuccessAccessListSegmentationDisabled, nil)
}

// EnableAccessListSegmentation implements accesslistsegmentation.AccessListSegmentationHandler.
func (a *accessListSegmentation) EnableAccessListSegmentation(w http.ResponseWriter, r *http.Request) {
	id := chi.URLParam(r, "id")
	if id == "" {
		a.logger.Errorf("[EnableAccessListSegmentation] missing id parameter")
		localization.SendErrorByCodeResponse(w, localization.ErrorAccessListSegmentationInvalidID.Code)
		return
	}
	if err := a.service.EnableDisableAccessListSegmentation(r.Context(), id, true); err != nil {
		a.logger.Errorf("[EnableAccessListSegmentation] service error: %v", err)
		localization.SendErrorByCodeResponse(w, err.Error())
		return
	}
	localization.SendSuccessResponse(w, localization.SuccessAccessListSegmentationEnabled, nil)
}

// GetAccessListSegmentationByID implements accesslistsegmentation.AccessListSegmentationHandler.
func (a *accessListSegmentation) GetAccessListSegmentationByID(w http.ResponseWriter, r *http.Request) {
	id := chi.URLParam(r, "id")
	if id == "" {
		a.logger.Errorf("[GetAccessListSegmentationByID] missing id parameter")
		localization.SendErrorByCodeResponse(w, localization.ErrorAccessListSegmentationInvalidID.Code)
		return
	}
	accessListSegmentation, err := a.service.GetAccessListSegmentationByID(r.Context(), id)
	if err != nil {
		a.logger.Errorf("[GetAccessListSegmentationByID] service error: %v", err)
		localization.SendErrorByCodeResponse(w, err.Error())
		return
	}
	localization.SendSuccessResponse(w, localization.SuccessAccessListSegmentationRetrieved, accessListSegmentation)
}

// GetAllAccessListSegmentation implements accesslistsegmentation.AccessListSegmentationHandler.
func (a *accessListSegmentation) GetAllAccessListSegmentation(w http.ResponseWriter, r *http.Request) {
	filterParams := local_util.ExtractFilterParams(r)

	accessListSegmentations, err := a.service.GetAllAccessListSegmentation(r.Context(), *filterParams)
	if err != nil {
		a.logger.Errorf("[GetAllAccessListSegmentation] service error: %v", err)
		localization.SendErrorByCodeResponse(w, err.Error())
		return
	}
	localization.SendSuccessResponse(w, localization.SuccessAccessListSegmentationRetrieved, accessListSegmentations)
}

// UpdateAccessListSegmentation implements accesslistsegmentation.AccessListSegmentationHandler.
func (a *accessListSegmentation) UpdateAccessListSegmentation(w http.ResponseWriter, r *http.Request) {
	id := chi.URLParam(r, "id")
	if id == "" {
		a.logger.Errorf("[UpdateAccessListSegmentation] missing id parameter")
		localization.SendErrorByCodeResponse(w, localization.ErrorAccessListSegmentationInvalidID.Code)
		return
	}
	var req access_list_segmentation_dto.UpdateAccessListSegmentationRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		a.logger.Errorf("[UpdateAccessListSegmentation] failed to decode request body: %v", err)
		localization.SendErrorByCodeResponse(w, err.Error())
		return
	}
	req.ID = id
	if err := req.Validate(); err != nil {
		a.logger.Errorf("[UpdateAccessListSegmentation] validation error: %v", err)
		localization.SendErrorByCodeResponse(w, err.Error())
		return
	}
	if err := a.service.UpdateAccessListSegmentation(r.Context(), req); err != nil {
		a.logger.Errorf("[UpdateAccessListSegmentation] service error: %v", err)
		localization.SendErrorByCodeResponse(w, err.Error())
		return
	}
	localization.SendSuccessResponse(w, localization.SuccessAccessListSegmentationUpdated, nil)
}

func (a *accessListSegmentation) GetAllAccessListSegmentationBySegmentIDorSegmentCode(w http.ResponseWriter, r *http.Request) {
	id := chi.URLParam(r, "id")
	if id == "" {
		a.logger.Errorf("[GetAllAccessListSegmentationBySegmentIDorSegmentCode] missing id parameter")
		localization.SendErrorByCodeResponse(w, localization.ErrorAccessListSegmentationInvalidID.Code)
		return
	}
	accesssList, accessListSegmentations, err := a.service.GetAllAccessListSegmentationBySegmentIDorSegmentCode(r.Context(), id)
	if err != nil {
		a.logger.Errorf("[GetAllAccessListSegmentationBySegmentIDorSegmentCode] service error: %v", err)
		localization.SendErrorByCodeResponse(w, err.Error())
		return
	}

	localization.SendSuccessResponse(w, localization.SuccessAccessListSegmentationRetrieved, map[string]interface{}{
		"access_lists":              accesssList,
		"access_list_segmentations": accessListSegmentations,
	})
}

func InitAccessListSegmentationAdapter(accessListSegmentationApplication service.AccessListSegmentationService, logger utils.Logger) accesslistsegmentation.AccessListSegmentationHandler {
	return &accessListSegmentation{
		service: accessListSegmentationApplication,
		logger:  logger,
	}
}
