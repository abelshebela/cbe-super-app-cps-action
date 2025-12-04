package newscategory_adaptor

import "net/http"

type NewsCategoryAdaptor interface {
	FetchNewsCategories(w http.ResponseWriter, r *http.Request)
	GetNewsCategoryByID(w http.ResponseWriter, r *http.Request)
	CreateNewsCategory(w http.ResponseWriter, r *http.Request)
	UpdateNewsCategory(w http.ResponseWriter, r *http.Request)
	DeleteNewsCategory(w http.ResponseWriter, r *http.Request)
}
