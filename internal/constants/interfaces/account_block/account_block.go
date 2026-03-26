package accountblock

import "net/http"

type AccountBlockAdapter interface {
	GetBranchById(w http.ResponseWriter, r *http.Request)
	GetAllBranches(w http.ResponseWriter, r *http.Request)
	GetRegionById(w http.ResponseWriter, r *http.Request)
	GetAllRegions(w http.ResponseWriter, r *http.Request)
	GetDistrictById(w http.ResponseWriter, r *http.Request)
	GetAllDistricts(w http.ResponseWriter, r *http.Request)
	GetCityById(w http.ResponseWriter, r *http.Request)
	GetAllCities(w http.ResponseWriter, r *http.Request)
	EnableBranches(w http.ResponseWriter, r *http.Request)
	DisableBranches(w http.ResponseWriter, r *http.Request)
	EnableRegions(w http.ResponseWriter, r *http.Request)
	DisableRegions(w http.ResponseWriter, r *http.Request)
	EnableDistricts(w http.ResponseWriter, r *http.Request)
	DisableDistricts(w http.ResponseWriter, r *http.Request)
	EnableCities(w http.ResponseWriter, r *http.Request)
	DisableCities(w http.ResponseWriter, r *http.Request)
	GetAccountBlockDetails(w http.ResponseWriter, r *http.Request)
}
