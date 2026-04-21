package accesslistsegmentation

import "net/http"

type AccessListSegmentationHandler interface {
	CreateAccessListSegmentation(w http.ResponseWriter, r *http.Request)

	GetAllAccessListSegmentation(w http.ResponseWriter, r *http.Request)

	GetAccessListSegmentationForAccountByID(w http.ResponseWriter, r *http.Request)
	GetAccessListSegmentationForBlockByID(w http.ResponseWriter, r *http.Request)

	GetAllAccessListSegmentationForAccount(w http.ResponseWriter, r *http.Request)
	GetAllAccessListSegmentationForBlock(w http.ResponseWriter, r *http.Request)

	UpdateAccessListSegmentation(w http.ResponseWriter, r *http.Request)

	EnableAccessListSegmentation(w http.ResponseWriter, r *http.Request)
	DisableAccessListSegmentation(w http.ResponseWriter, r *http.Request)
}
