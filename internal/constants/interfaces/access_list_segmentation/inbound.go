package accesslistsegmentation

import "net/http"

type AccessListSegmentationHandler interface {
	CreateAccessListSegmentation(w http.ResponseWriter, r *http.Request)
	GetAllAccessListSegmentation(w http.ResponseWriter, r *http.Request)
	GetAccessListSegmentationByID(w http.ResponseWriter, r *http.Request)
	GetAllAccessListSegmentationBySegmentIDorSegmentCode(w http.ResponseWriter, r *http.Request)
	UpdateAccessListSegmentation(w http.ResponseWriter, r *http.Request)
	EnableAccessListSegmentation(w http.ResponseWriter, r *http.Request)
	DisableAccessListSegmentation(w http.ResponseWriter, r *http.Request)
}
