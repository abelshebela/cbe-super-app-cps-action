package inbound

import "net/http"

type BranchHandler interface {
	FilterSingleBranches(w http.ResponseWriter, r *http.Request)
	DisableSingleBranch(w http.ResponseWriter, r *http.Request)
	ApproveSingleBranchDisable(w http.ResponseWriter, r *http.Request)

	FilterMultipleBranches(w http.ResponseWriter, r *http.Request)
	DisableMultipleBranches(w http.ResponseWriter, r *http.Request)
	ApproveBulkBranchesDisable(w http.ResponseWriter, r *http.Request)
}
