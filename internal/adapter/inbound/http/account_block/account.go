package accountblock_handler

import (
	"net/http"
	"strings"

	"github.com/go-chi/chi/v5"

	"github.com/CBE-Super-App/cbe-super-app-cps-action/internal/application/account_block"
	"github.com/CBE-Super-App/cbe-super-app-cps-action/pkgs/common"
	constant_utils "github.com/CBE-Super-App/cbe-super-app-cps-action/pkgs/utils"
	constant "github.com/CBE-Super-App/cbe-super-app-cps-action/utils"

	"gitlab.com/bersufekadgetachew/cbe-super-app-shared/shared/utils"
)

type AccountBlockHandler struct {
	service account_block.ApplicationService
	logger  utils.Logger
}

func NewAccountBlockHandler(service account_block.ApplicationService, logger utils.Logger) *AccountBlockHandler {
	return &AccountBlockHandler{service: service, logger: logger}
}

func (h *AccountBlockHandler) GetBranch(w http.ResponseWriter, r *http.Request) {
	filterParams := constant_utils.ExtractFilterParams(r)
	branchCode := strings.TrimSpace(chi.URLParam(r, "branch_code"))
	if branchCode == "" {
		constant_utils.SendErrorResponse(w, "BRANCH_CODE_REQUIRED", 0, nil)
	}

	branch, err := h.service.GetBranch(r.Context(), branchCode, filterParams)
	if err != nil {
		h.logger.Errorf("FetchUserRequest failed: %v", err)
		constant_utils.SendErrorResponse(w, err.Error(), 0, nil)
		return
	}

	data, err := constant_utils.StructToMap(branch)
	if err != nil {
		constant_utils.SendErrorResponse(w, "UNHANDLED_SERVER_ERROR", 500, nil)
		return
	}
	constant_utils.BaseResponseMaker(data, w, "Branch retrieved successfully", http.StatusOK)
}

func (h *AccountBlockHandler) GetAllBranches(w http.ResponseWriter, r *http.Request) {
	region := strings.TrimSpace(r.URL.Query().Get("region"))
	district := strings.TrimSpace(r.URL.Query().Get("district"))

	filterParams := constant_utils.ExtractFilterParams(r)

	// if len(region) < 3 || len(district) < 3 {
	// 	constant_utils.SendErrorResponse(w, "BRANCH_REGION_AND_DISTRICT_MIN_LENGTH", 400, nil)
	// 	return
	// }

	branches, err := h.service.GetAllBranches(r.Context(), region, district, filterParams)
	if err != nil {
		h.logger.Errorf("GetAllBranches failed: %v", err)
		switch err {
		case common.DefineError.Branch["BRANCH_REGION_AND_DISTRICT_REQUIRED"]:
			constant_utils.SendErrorResponse(w, "BRANCH_REGION_AND_DISTRICT_REQUIRED", 400, nil)
		case common.DefineError.Branch["FAILED_TO_FETCH_BRANCHES"]:
			constant_utils.SendErrorResponse(w, "FAILED_TO_FETCH_BRANCHES", 500, nil)
		case common.DefineError.Branch["BRANCH_NOT_FOUND"]:
			constant_utils.SendErrorResponse(w, "BRANCH_NOT_FOUND", 404, nil)
		default:
			constant_utils.SendErrorResponse(w, "UNHANDLED_SERVER_ERROR", 500, nil)
		}
		return
	}

	data, err := constant_utils.StructToMap(branches)
	if err != nil {
		constant_utils.SendErrorResponse(w, "UNHANDLED_SERVER_ERROR", 500, nil)
		return
	}
	constant_utils.BaseResponseMaker(data, w, "Branches retrieved successfully", http.StatusOK)
}

func (h *AccountBlockHandler) GetRegionByCode(w http.ResponseWriter, r *http.Request) {
	regionCode := strings.TrimSpace(chi.URLParam(r, "region_code"))

	userID, _ := r.Context().Value(constant.ContextKey("user_id")).(string)
	fullName, _ := r.Context().Value(constant.ContextKey("full_name")).(string)
	phoneNumber, _ := r.Context().Value(constant.ContextKey("phone_number")).(string)
	department, _ := r.Context().Value(constant.ContextKey("department")).(string)

	if strings.TrimSpace(userID) == "" ||
		strings.TrimSpace(fullName) == "" ||
		strings.TrimSpace(phoneNumber) == "" {
		constant_utils.SendErrorResponse(w, "INCOMPLETE_USER_INFO", http.StatusUnauthorized, nil)
		return
	}
	if strings.TrimSpace(department) == "" {
		constant_utils.SendErrorResponse(w, "DEPARTMENT_REQUIRED_IN_CONTEXT", http.StatusUnauthorized, nil)
		return
	}

	h.logger.Infof("GetRegionByCode requested by user: %s (%s)", userID, fullName)

	region, err := h.service.GetRegionByCode(r.Context(), regionCode)
	if err != nil {
		h.logger.Errorf("GetRegionByCode failed: %v", err)
		switch err {
		case common.DefineError.Branch["REGION_NOT_FOUND"]:
			constant_utils.SendErrorResponse(w, "REGION_NOT_FOUND", http.StatusNotFound, nil)
		default:
			constant_utils.SendErrorResponse(w, "UNHANDLED_SERVER_ERROR", http.StatusInternalServerError, nil)
		}
		return
	}

	data, err := constant_utils.StructToMap(region)
	if err != nil {
		constant_utils.SendErrorResponse(w, "UNHANDLED_SERVER_ERROR", http.StatusInternalServerError, nil)
		return
	}
	constant_utils.BaseResponseMaker(data, w, "Region retrieved successfully", http.StatusOK)
}

func (h *AccountBlockHandler) GetAllRegions(w http.ResponseWriter, r *http.Request) {
	filterParams := constant_utils.ExtractFilterParams(r)
	if filterParams == nil || filterParams.Page < 1 || filterParams.PerPage < 1 {
		constant_utils.SendErrorResponse(w, "INVALID_PAGINATION_PARAMS", http.StatusBadRequest, nil)
		return
	}

	regions, err := h.service.GetAllRegions(r.Context(), filterParams)
	if err != nil {
		h.logger.Errorf("GetAllRegion failed: %v", err)
		constant_utils.SendErrorResponse(w, err.Error(), 0, nil)
		return
	}

	data, err := constant_utils.StructToMap(regions)
	if err != nil {
		constant_utils.SendErrorResponse(w, "UNHANDLED_SERVER_ERROR", http.StatusInternalServerError, nil)
		return
	}
	constant_utils.BaseResponseMaker(data, w, "Regions retrieved successfully", http.StatusOK)
}

func (h *AccountBlockHandler) GetDistrictByCode(w http.ResponseWriter, r *http.Request) {
	districtCode := strings.TrimSpace(chi.URLParam(r, "district_code"))

	userID, _ := r.Context().Value(constant.ContextKey("user_id")).(string)
	fullName, _ := r.Context().Value(constant.ContextKey("full_name")).(string)
	phoneNumber, _ := r.Context().Value(constant.ContextKey("phone_number")).(string)
	department, _ := r.Context().Value(constant.ContextKey("department")).(string)

	if strings.TrimSpace(userID) == "" ||
		strings.TrimSpace(fullName) == "" ||
		strings.TrimSpace(phoneNumber) == "" {

		constant_utils.SendErrorResponse(w, "DEPARTMEN_REQUIRED", 409, nil)
		return
	}
	if strings.TrimSpace(department) == "" {
		constant_utils.SendErrorResponse(w, "DEPARTMEN_REQUIRED", 409, nil)
		return
	}

	h.logger.Infof("GetDistrictByCode requested by user: %s (%s)", userID, fullName)

	district, err := h.service.GetDistrictByCode(r.Context(), districtCode)
	if err != nil {
		h.logger.Errorf("GetDistrictByCode failed: %v", err)
		constant_utils.SendErrorResponse(w, err.Error(), http.StatusNotFound, nil)
		return
	}

	data, err := constant_utils.StructToMap(district)
	if err != nil {
		constant_utils.SendErrorResponse(w, http.StatusInternalServerError, 500, nil)
		return
	}
	constant_utils.BaseResponseMaker(data, w, "District retrieved successfully", http.StatusOK)
}

func (h *AccountBlockHandler) GetAllDistricts(w http.ResponseWriter, r *http.Request) {
	filterParams := constant_utils.ExtractFilterParams(r)
	if filterParams == nil || filterParams.Page < 1 || filterParams.PerPage < 1 {
		constant_utils.SendErrorResponse(w, "INVALID_PAGINATION_PARAMS", http.StatusBadRequest, nil)
		return
	}

	districts, err := h.service.GetAllDistricts(r.Context(), filterParams)
	if err != nil {
		h.logger.Errorf("GetAllDistrict failed: %v", err)
		constant_utils.SendErrorResponse(w, err.Error(), http.StatusInternalServerError, nil)
		return
	}

	data, err := constant_utils.StructToMap(districts)
	if err != nil {
		constant_utils.SendErrorResponse(w, "UNHANDLED_SERVER_ERROR", http.StatusInternalServerError, nil)
		return
	}
	constant_utils.BaseResponseMaker(data, w, "Districts retrieved successfully", http.StatusOK)
}

func (h *AccountBlockHandler) GetCityByCode(w http.ResponseWriter, r *http.Request) {
	cityCode := strings.TrimSpace(chi.URLParam(r, "city_code"))

	userID, _ := r.Context().Value(constant.ContextKey("user_id")).(string)
	fullName, _ := r.Context().Value(constant.ContextKey("full_name")).(string)
	phoneNumber, _ := r.Context().Value(constant.ContextKey("phone_number")).(string)
	department, _ := r.Context().Value(constant.ContextKey("department")).(string)

	if strings.TrimSpace(userID) == "" ||
		strings.TrimSpace(fullName) == "" ||
		strings.TrimSpace(phoneNumber) == "" {

		constant_utils.SendErrorResponse(w, "INCOMPLETE_USER_INFO", 409, nil)
		return
	}
	if strings.TrimSpace(department) == "" {
		constant_utils.SendErrorResponse(w, "DEPARTMEN_REQUIRED", 409, nil)
		return
	}

	h.logger.Infof("GetCityByCode requested by user: %s (%s)", userID, fullName)

	city, err := h.service.GetCityByCode(r.Context(), cityCode)
	if err != nil {
		h.logger.Errorf("GetCityByCode failed: %v", err)
		constant_utils.SendErrorResponse(w, "CITY_NOT_FOUND", 0, nil)
		return
	}

	data, err := constant_utils.StructToMap(city)
	if err != nil {
		constant_utils.SendErrorResponse(w, http.StatusInternalServerError, 500, nil)
		return
	}
	constant_utils.BaseResponseMaker(data, w, "City retrieved successfully", http.StatusOK)
}

func (h *AccountBlockHandler) GetAllCities(w http.ResponseWriter, r *http.Request) {
	filterParams := constant_utils.ExtractFilterParams(r)
	if filterParams == nil || filterParams.Page < 1 || filterParams.PerPage < 1 {
		constant_utils.SendErrorResponse(w, "INVALID_PAGINATION_PARAMS", http.StatusBadRequest, nil)
		return
	}

	cities, err := h.service.GetAllCities(r.Context(), filterParams)
	if err != nil {
		h.logger.Errorf("GetAllCity failed: %v", err)
		constant_utils.SendErrorResponse(w, err.Error(), http.StatusInternalServerError, nil)
		return
	}

	data, err := constant_utils.StructToMap(cities)
	if err != nil {
		constant_utils.SendErrorResponse(w, "UNHANDLED_SERVER_ERROR", http.StatusInternalServerError, nil)
		return
	}
	constant_utils.BaseResponseMaker(data, w, "Cities retrieved successfully", http.StatusOK)
}

func (h *AccountBlockHandler) EnableBranches(w http.ResponseWriter, r *http.Request) {
	err := h.service.EnableBranches(r)
	if err != nil {
		h.logger.Errorf("Enable branch/es request failed: %v\n", err)
		constant_utils.SendErrorResponse(w, err.Error(), 0, nil)
		return
	}

	constant_utils.BaseResponseMaker(nil, w, "Request to enable branch/es submitted successfully", 200)
}

func (h *AccountBlockHandler) DisableBranches(w http.ResponseWriter, r *http.Request) {
	err := h.service.DisableBranches(r)
	if err != nil {
		h.logger.Errorf("Disable branch/es request failed: %v\n", err)
		constant_utils.SendErrorResponse(w, err.Error(), 0, nil)
		return
	}

	constant_utils.BaseResponseMaker(nil, w, "Request to disable branch/es submitted successfully", 200)
}

func (h *AccountBlockHandler) EnableRegion(w http.ResponseWriter, r *http.Request) {
	err := h.service.EnableRegion(r)
	if err != nil {
		h.logger.Errorf("Enable region/es request failed: %v\n", err)
		constant_utils.SendErrorResponse(w, err.Error(), 0, nil)
		return
	}

	constant_utils.BaseResponseMaker(nil, w, "Request to enable region/s submitted successfully", 200)
}

func (h *AccountBlockHandler) DisableRegion(w http.ResponseWriter, r *http.Request) {
	err := h.service.DisableRegion(r)
	if err != nil {
		h.logger.Errorf("Disable region/es request failed: %v\n", err)
		constant_utils.SendErrorResponse(w, err.Error(), 0, nil)
		return
	}

	constant_utils.BaseResponseMaker(nil, w, "Request to disable region/s submitted successfully", 200)
}

func (h *AccountBlockHandler) EnableDistrict(w http.ResponseWriter, r *http.Request) {
	err := h.service.EnableDistrict(r)
	if err != nil {
		h.logger.Errorf("Enable district/es request failed: %v\n", err)
		constant_utils.SendErrorResponse(w, err.Error(), 0, nil)
		return
	}

	constant_utils.BaseResponseMaker(nil, w, "Request to enable district/s submitted successfully", 200)
}

func (h *AccountBlockHandler) DisableDistrict(w http.ResponseWriter, r *http.Request) {
	err := h.service.DisableDistrict(r)
	if err != nil {
		h.logger.Errorf("Disable district/es request failed: %v\n", err)
		constant_utils.SendErrorResponse(w, err.Error(), 0, nil)
		return
	}

	constant_utils.BaseResponseMaker(nil, w, "Request to disable district/s submitted successfully", 200)
}

func (h *AccountBlockHandler) EnableCity(w http.ResponseWriter, r *http.Request) {
	err := h.service.EnableCity(r)
	if err != nil {
		h.logger.Errorf("Enable city/es request failed: %v\n", err)
		constant_utils.SendErrorResponse(w, err.Error(), 0, nil)
		return
	}

	constant_utils.BaseResponseMaker(nil, w, "Request to enable city or cities submitted successfully", 200)
}

func (h *AccountBlockHandler) DisableCity(w http.ResponseWriter, r *http.Request) {
	err := h.service.DisableCity(r)
	if err != nil {
		h.logger.Errorf("Disable city/es request failed: %v\n", err)
		constant_utils.SendErrorResponse(w, err.Error(), 0, nil)
		return
	}

	constant_utils.BaseResponseMaker(nil, w, "Request to disable city or cities submitted successfully", 200)
}
