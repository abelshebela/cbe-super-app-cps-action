package newstag_adaptor

import "net/http"

type NewsTagAdaptor interface {
	FetchNewsTags(w http.ResponseWriter, r *http.Request)
	GetNewsTagByID(w http.ResponseWriter, r *http.Request)
	CreateNewsTags(w http.ResponseWriter, r *http.Request)
	UpdateNewsTag(w http.ResponseWriter, r *http.Request)
	DeleteNewsTag(w http.ResponseWriter, r *http.Request)
}
