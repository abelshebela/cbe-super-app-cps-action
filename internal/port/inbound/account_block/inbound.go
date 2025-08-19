package accountblock

import "net/http"

type AccountBlockHandler interface {
	GetBranch(w http.ResponseWriter, r *http.Request)
	GetAllBranches(w http.ResponseWriter, r *http.Request)
	GetRegionByCode(w http.ResponseWriter, r *http.Request)
	GetAllRegions(w http.ResponseWriter, r *http.Request)
	GetDistrictByCode(w http.ResponseWriter, r *http.Request)
	GetAllDistricts(w http.ResponseWriter, r *http.Request)
	GetCityByCode(w http.ResponseWriter, r *http.Request)
	GetAllCities(w http.ResponseWriter, r *http.Request)
	EnableBranches(w http.ResponseWriter, r *http.Request)
	DisableBranches(w http.ResponseWriter, r *http.Request)
	EnableRegion(w http.ResponseWriter, r *http.Request)
	DisableRegion(w http.ResponseWriter, r *http.Request)
	EnableDistrict(w http.ResponseWriter, r *http.Request)
	DisableDistrict(w http.ResponseWriter, r *http.Request)
	EnableCity(w http.ResponseWriter, r *http.Request)
	DisableCity(w http.ResponseWriter, r *http.Request)
}
