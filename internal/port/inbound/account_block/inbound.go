package accountblock

import "net/http"

type AccountBlockHandler interface {
	FilterSingleBranches(w http.ResponseWriter, r *http.Request)
	GetAllBranches(w http.ResponseWriter, r *http.Request)
	DisableSingleBranch(w http.ResponseWriter, r *http.Request)
	EnableSingleBranch(w http.ResponseWriter, r *http.Request)
	ApproveSingleBranchDisable(w http.ResponseWriter, r *http.Request)

	FilterMultipleBranches(w http.ResponseWriter, r *http.Request)
	DisableMultipleBranches(w http.ResponseWriter, r *http.Request)
	EnableMultipleBranches(w http.ResponseWriter, r *http.Request)
	ApproveBulkBranchesDisable(w http.ResponseWriter, r *http.Request)

	BlockRegion(w http.ResponseWriter, r *http.Request)
	UpdateRegion(w http.ResponseWriter, r *http.Request)
	ApproveRegionBlock(w http.ResponseWriter, r *http.Request)
	GetRegionByCode(w http.ResponseWriter, r *http.Request)
	GetAllRegion(w http.ResponseWriter, r *http.Request)

	BlockDistrict(w http.ResponseWriter, r *http.Request)
	GetDistrictByCode(w http.ResponseWriter, r *http.Request)
	GetAllDistrict(w http.ResponseWriter, r *http.Request)
	ApproveBlockDistrict(w http.ResponseWriter, r *http.Request)

	BlockCity(w http.ResponseWriter, r *http.Request)
	GetCityByCode(w http.ResponseWriter, r *http.Request)
	GetAllCities(w http.ResponseWriter, r *http.Request)
	ApproveBlockCity(w http.ResponseWriter, r *http.Request)

	BlockUser(w http.ResponseWriter, r *http.Request)
	GetUserByPhone(w http.ResponseWriter, r *http.Request)
	ApproveBlockUser(w http.ResponseWriter, r *http.Request)

	// Newly added
	// Branch
	EnableBranches(w http.ResponseWriter, r *http.Request)
	DisableBranches(w http.ResponseWriter, r *http.Request)

	// Region
	EnableRegion(w http.ResponseWriter, r *http.Request)
	DisableRegion(w http.ResponseWriter, r *http.Request)

	// District
	EnableDistrict(w http.ResponseWriter, r *http.Request)
	DisableDistrict(w http.ResponseWriter, r *http.Request)

	// City
	EnableCity(w http.ResponseWriter, r *http.Request)
	DisableCity(w http.ResponseWriter, r *http.Request)
}
