package accountblock

import "net/http"

type AccountBlockHandler interface {
	FilterSingleBranches(w http.ResponseWriter, r *http.Request)
	GetAllBranches(w http.ResponseWriter, r *http.Request)
	DisableSingleBranch(w http.ResponseWriter, r *http.Request)
	ApproveSingleBranchDisable(w http.ResponseWriter, r *http.Request)

	FilterMultipleBranches(w http.ResponseWriter, r *http.Request)
	DisableMultipleBranches(w http.ResponseWriter, r *http.Request)
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
}
