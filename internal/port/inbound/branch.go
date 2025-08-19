package inbound

import "net/http"

type BranchHandler interface {
	GetBranch(w http.ResponseWriter, r *http.Request)
	DisableSingleBranch(w http.ResponseWriter, r *http.Request)
	ApproveSingleBranchDisable(w http.ResponseWriter, r *http.Request)

	GetAllBranches(w http.ResponseWriter, r *http.Request)
	DisableMultipleBranches(w http.ResponseWriter, r *http.Request)
	ApproveBulkBranchesDisable(w http.ResponseWriter, r *http.Request)
}
