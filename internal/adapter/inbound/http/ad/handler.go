package ad

import (
	"encoding/json"
	"fmt"
	"net/http"

	application "github.com/CBE-Super-App/cbe-super-app-cps-action/internal/application/ad"
	"github.com/CBE-Super-App/cbe-super-app-cps-action/internal/port/inbound/ad"
	"github.com/CBE-Super-App/cbe-super-app-cps-action/pkgs/utils"
	shared "gitlab.com/bersufekadgetachew/cbe-super-app-shared/shared/utils"
)

const (
	DefaultPage     = 1
	DefaultPageSize = 20
	MaxPageSize     = 100
)

// AdvertHTTPStore handles HTTP requests for advertisement operations
type AdvertHTTPStore struct {
	Application application.ADHandlers
	logger      shared.Logger
}

// NewAdvertHTTPHandler initializes a new AdvertHTTPStore
func NewAdvertHTTPHandler(app application.ADHandlers, logger shared.Logger) ad.ADAdapter {
	return &AdvertHTTPStore{
		Application: app,
		logger:      logger,
	}
}

// CreateAdvert handles creation of a new advertisement
func (h *AdvertHTTPStore) CreateAdvert(w http.ResponseWriter, r *http.Request) {
	req, ok := h.parseAndValidateAdvertRequest(w, r, false)
	if !ok {
		return
	}

	maker, ok := h.extractUserAndMaker(w, r)
	if !ok {
		return
	}

	domainReq, _ := ToDomainAdvertRequest(req)
	err := h.Application.CreateAdvert(r.Context(), domainReq, maker)
	h.sendResponse(w, err, "Advert creation request submitted successfully", nil)
}

func PrettyPrintJSON(data interface{}) error {
	prettyJSON, err := json.MarshalIndent(data, "", "  ")
	if err != nil {
		return err
	}

	fmt.Println(string(prettyJSON))
	return nil
}

// UpdateAdvert handles updating an advertisement
func (h *AdvertHTTPStore) UpdateAdvert(w http.ResponseWriter, r *http.Request) {
	id, ok := h.extractID(w, r)
	if !ok {
		return
	}

	req, ok := h.parseAndValidateAdvertRequest(w, r, true)
	if !ok {
		return
	}

	maker, ok := h.extractUserAndMaker(w, r)
	if !ok {
		return
	}

	PrettyPrintJSON(req)

	domainReq, _ := ToDomainAdvertRequest(req)

PrettyPrintJSON(domainReq)

	if domainReq.Title == "" && domainReq.Description == "" && domainReq.AdvertFor == "" && domainReq.BannerImage == nil && domainReq.Date.StartedAt.IsZero() && domainReq.Date.ExpiredAt.IsZero() {
		h.logger.Errorf("[event.UpdateAdvert] no data provided for update, id: %s", id)
		utils.SendErrorResponse(w, utils.NoDataProvidedForUpdate, http.StatusBadRequest, nil)
		return
	}

	err := h.Application.UpdateAdvert(r.Context(), id, domainReq, maker)
	h.sendResponse(w, err, "Advert update request submitted successfully", nil)
}

// DeleteAdvert handles deletion of an advertisement
func (h *AdvertHTTPStore) DeleteAdvert(w http.ResponseWriter, r *http.Request) {
	id, ok := h.extractID(w, r)
	if !ok {
		return
	}

	maker, ok := h.extractUserAndMaker(w, r)
	if !ok {
		return
	}

	err := h.Application.DeleteAdvert(r.Context(), id, maker)
	h.sendResponse(w, err, "Advert delete request submitted successfully", nil)
}

// EnableAdvert handles enabling an advertisement
func (h *AdvertHTTPStore) EnableAdvert(w http.ResponseWriter, r *http.Request) {
	id, ok := h.extractID(w, r)
	if !ok {
		return
	}

	maker, ok := h.extractUserAndMaker(w, r)
	if !ok {
		return
	}

	err := h.Application.EnableDisableAdvert(r.Context(), id, maker, true)
	h.sendResponse(w, err, "Advert enable request submitted successfully", nil)
}

// DisableAdvert handles disabling an advertisement
func (h *AdvertHTTPStore) DisableAdvert(w http.ResponseWriter, r *http.Request) {
	id, ok := h.extractID(w, r)
	if !ok {
		return
	}

	maker, ok := h.extractUserAndMaker(w, r)
	if !ok {
		return
	}

	err := h.Application.EnableDisableAdvert(r.Context(), id, maker, false)
	h.sendResponse(w, err, "Advert disable request submitted successfully", nil)
}

// FetchAdvertByID retrieves a single advertisement
func (h *AdvertHTTPStore) FetchAdvertByID(w http.ResponseWriter, r *http.Request) {
	id, ok := h.extractID(w, r)
	if !ok {
		return
	}

	data, err := h.Application.FetchAdvertByID(r.Context(), id)
	if err != nil {
		h.sendResponse(w, err, "", nil)
		return
	}

	res := ToAdvertResponse(*data)
	h.sendResponse(w, nil, "Advert successfully retrieved", res)
}

// FetchAdverts retrieves all advertisements
func (h *AdvertHTTPStore) FetchAdverts(w http.ResponseWriter, r *http.Request) {
	filterParams := utils.ExtractFilterParams(r)
	list, err := h.Application.FetchAdverts(r.Context(), filterParams)
	if err != nil {
		h.sendResponse(w, err, "", nil)
		return
	}

	docs := ToAdvertResponses(list.Data)
	res := utils.PaginatedResponse[[]*AdvertResponse]{
		Data: docs,
		Meta: list.Meta,
	}
	h.sendResponse(w, nil, "Adverts successfully retrieved", res)
}
